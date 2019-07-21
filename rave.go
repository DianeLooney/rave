package rave

import (
	"sync"
	"time"

	"github.com/dianelooney/ferry"
)

type Context struct {
	mtx            sync.Mutex
	globalBeat     *ferry.Value
	globalBeatOnce sync.Once
	doc            *Doc
	insts          map[string]Inst
	beats          map[string]*ferry.Value
}

func (p *Context) Init() {
	p.insts = make(map[string]Inst)
	p.beats = make(map[string]*ferry.Value)
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
			p.beats[k.ID()] = ferry.NewValue()
		}
	}
	for _, i := range p.doc.Insts {
		go p.spawnInst(i)
	}
	go p.globalBeatOnce.Do(p.spawnGlobalBeat)
}

func (p *Context) spawnGlobalBeat() {
	p.globalBeat = ferry.NewValue()

	time.Sleep(250 * time.Millisecond)
	next := time.Now()
	for {
		p.globalBeat.Done(next)

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
func (p *Context) spawnInst(i Inst) {
	go func() {
		p.mtx.Lock()
		old := p.insts[i.ID()]
		p.mtx.Unlock()

		if old != nil {
			<-old.Done()
		}

		p.mtx.Lock()
		p.insts[i.ID()] = i
		p.mtx.Unlock()

		for {
			select {
			case i.Done() <- true:
				return
			default:
			}
			i.PlayLoop(p)
		}
	}()
}

func beatLength(bpm float64, count float64) time.Duration {
	seconds := (60 / float64(bpm)) * count
	return time.Duration(float64(time.Second) * seconds)
}

func waitForBeat(t time.Time, bpm float64, count float64) {
	elapsed := beatLength(bpm, count)
	until := time.Until(t.Add(elapsed))
	<-time.NewTimer(until).C
}
