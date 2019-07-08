package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"github.com/dianelooney/directive"
	"github.com/dianelooney/ferry"
	"github.com/dianelooney/rave/kits"
	"github.com/fsnotify/fsnotify"
)

const src = "client/src.rave"

func main() {
	p := Player{}
	p.Init()

	loadFile := func() {
		bytes, err := ioutil.ReadFile(src)
		if err != nil {
			log.Fatalf("Unable to read file: %v\n", err)
		}
		e, err := directive.Prepare(bytes)
		if err != nil {
			log.Fatalf("Unable to prepare directives: %v\n", err)
		}
		d := Doc{}
		err = e.Execute(&d)
		if err != nil {
			log.Fatalf("Unable to execute directive: %v\n", err)
		}
		p.Load(&d)
	}

	loadFile()

	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			fmt.Println("ERROR", err)
		}
		defer watcher.Close()

		watcher.Add(src)

		for {
			select {
			case e := <-watcher.Events:
				if e.Op != fsnotify.Write {
					continue
				}
				loadFile()
			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}

	}()
	<-make(chan bool)
}

type Doc struct {
	TimeTop float64
	TimeBot float64
	Tempo   float64
	Kits    []*Kit
}

func (d *Doc) Kit() *Kit {
	i := &Kit{
		Name:   "",
		Volume: 0.1,
		Sync:   "global",
		done:   make(chan bool),
	}
	d.Kits = append(d.Kits, i)
	return i
}

type Kit struct {
	Name   string
	Sample string
	Sync   string
	Volume float64
	loop   *Loop
	done   chan bool
}

func (k *Kit) Loop() *Loop {
	if k.loop == nil {
		k.loop = &Loop{}
	}

	return k.loop
}

type Loop struct {
	Measures []*Measure
}

func (l *Loop) Measure() *Measure {
	m := &Measure{}
	l.Measures = append(l.Measures, m)
	return m
}

type Measure struct {
	Pulses []float64
}

func (m *Measure) Pulse(t float64) {
	m.Pulses = append(m.Pulses, t)
}

type Player struct {
	mtx            sync.Mutex
	globalBeat     ferry.Ferry
	globalBeatOnce sync.Once
	doc            *Doc
	kits           map[string]*Kit
	beats          map[string]*ferry.Ferry
	kit            kits.Kit
}

func (d *Doc) hasKit(name string) bool {
	for _, k := range d.Kits {
		if k.Name == name {
			return true
		}
	}
	return false
}

func (p *Player) Init() {
	p.kits = make(map[string]*Kit)
	p.kit = kits.LoadManifest("kits/manifest.yml").Load()
	p.beats = make(map[string]*ferry.Ferry)
}

func (p *Player) Load(doc *Doc) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if p.doc != nil {
		for _, k := range p.doc.Kits {
			if doc.hasKit(k.Name) {
				continue
			}
			kit := p.kits[k.Name]
			go p.despawnKit(kit)
		}
	}

	p.doc = doc
	for _, k := range p.doc.Kits {
		_, ok := p.beats[k.Name]
		if !ok {
			p.beats[k.Name] = &ferry.Ferry{}
		}
	}
	for _, k := range p.doc.Kits {
		go p.spawnKit(k)
	}
	go p.globalBeatOnce.Do(p.spawnGlobalBeat)
}
func (p *Player) spawnGlobalBeat() {
	next := time.Now()
	for {
		p.globalBeat.Done()

		seconds := (60 / float64(p.doc.Tempo)) * p.doc.TimeTop
		elapsed := time.Duration(float64(time.Second) * seconds)
		next = next.Add(elapsed)
		<-time.After(time.Until(next))
	}
}
func (p *Player) despawnKit(k *Kit) {
	<-k.done
	p.mtx.Lock()
	defer p.mtx.Unlock()
	delete(p.kits, k.Name)
}
func (p *Player) spawnKit(k *Kit) {
	loop := k.loop
	sample := p.kit.Sample(k.Sample)
	sample = sample.ScaleAmplitude(k.Volume)
	sync, ok := p.beats[k.Sync]
	if !ok {
		sync = &p.globalBeat
	}
	beat := p.beats[k.Name]
	go func() {
		if old := p.kits[k.Name]; old != nil {
			<-old.done
		}
		p.kits[k.Name] = k
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
				for _, pulse := range m.Pulses {
					go func(pulse float64) {
						waitForBeat(p.doc.Tempo, pulse)
						sample.Play()
					}(pulse)
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
