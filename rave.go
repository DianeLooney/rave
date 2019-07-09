package rave

import (
	"log"
	"sync"
	"time"

	"github.com/dianelooney/ferry"
	"github.com/dianelooney/rave/common"
)

type Context struct {
	mtx            sync.Mutex
	globalBeat     ferry.Ferry
	globalBeatOnce sync.Once
	doc            *Doc
	insts          map[string]Inst
	beats          map[string]*ferry.Ferry
}

func (p *Context) Init() {
	p.insts = make(map[string]Inst)
	p.beats = make(map[string]*ferry.Ferry)
}

func (p *Context) Load(doc *Doc) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if p.doc != nil {
		for _, k := range p.doc.Insts {
			if doc.hasInst(k.ID()) {
				continue
			}
			inst := p.insts[k.ID()]
			go p.despawnInst(inst)
		}
	}

	p.doc = doc
	for _, k := range p.doc.Insts {
		_, ok := p.beats[k.ID()]
		if !ok {
			v := ferry.New()
			p.beats[k.ID()] = &v
		}
	}
	for _, i := range p.doc.Insts {
		switch v := i.(type) {
		case *Kit:
			go p.spawnKit(v)
		default:
			log.Fatalf("Unsupported type in Context.Load: '%T'\n", i)
		}
	}
	go p.globalBeatOnce.Do(p.spawnGlobalBeat)
}

func (p *Context) spawnGlobalBeat() {
	p.globalBeat = ferry.New()

	next := time.Now()
	for {
		p.globalBeat.Done()

		seconds := (60 / float64(p.doc.Tempo)) * p.doc.TimeTop
		elapsed := time.Duration(float64(time.Second) * seconds)
		next = next.Add(elapsed)
		<-time.After(time.Until(next))
	}
}

func (p *Context) despawnInst(i Inst) {
	<-i.Done()
	p.mtx.Lock()
	defer p.mtx.Unlock()
	delete(p.insts, i.ID())
}

func (p *Context) spawnKit(k *Kit) {
	loop := k.loop
	samples := make([]common.Sound, len(k.Samples))
	for i, s := range k.Samples {
		x := kit.Sample(s)
		samples[i] = x.ScaleAmplitude(k.Volume)
	}
	sync, ok := p.beats[k.Sync]
	if !ok {
		sync = &p.globalBeat
	}
	beat := p.beats[k.Name]
	go func() {
		if old := p.insts[k.ID()]; old != nil {
			<-old.Done()
		}
		p.mtx.Lock()
		p.insts[k.ID()] = k
		p.mtx.Unlock()
		for {
			select {
			case k.done <- true:
				return
			default:
			}

			for i, m := range loop.Measures {
				if i == 0 {
					sync.Wait()
					beat.Done()
				} else {
					p.globalBeat.Wait()
				}
				for i, pulse := range m.Pulses {
					go func(m *Measure, i int, pulse float64) {
						sampleI := 0
						if i < len(m.Samples) {
							sampleI = m.Samples[i]
						}
						x := samples[0]
						if sampleI < len(samples) {
							x = samples[sampleI]
						}

						if i < len(m.Weights) {
							x = x.ScaleAmplitude(m.Weights[i])
						}
						waitForBeat(p.doc.Tempo, pulse)
						x.Play()
					}(m, i, pulse)
				}
			}
		}
	}()
}

func waitForBeat(bpm float64, count float64) {
	seconds := (60 / float64(bpm)) * count
	elapsed := time.Duration(float64(time.Second) * seconds)
	<-time.NewTimer(elapsed).C
}
