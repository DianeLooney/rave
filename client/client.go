package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"github.com/dianelooney/directive"
	"github.com/dianelooney/rave/kits"
)

const SockAddr = "/tmp/rave.sock"

func main() {
	bytes, err := ioutil.ReadFile("client/src.rave")
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
	p := Player{}
	p.Init()
	p.Load(&d)
	<-make(chan bool)
}

type Doc struct {
	TimeTop float64
	TimeBot float64
	Tempo   float64
	Kits    []*Kit
}

func (d *Doc) Kit() *Kit {
	fmt.Println("*Doc.Kit()")
	i := &Kit{}
	d.Kits = append(d.Kits, i)
	return i
}

type Kit struct {
	Name   string
	Sample string
	loop   *Loop
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
	mtx  sync.Mutex
	doc  *Doc
	kits map[string]chan bool
	kit  kits.Kit
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
	p.kits = make(map[string]chan bool)
	p.kit = kits.LoadManifest("kits/manifest.yml").Load()
}

func (p *Player) Load(doc *Doc) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if p.doc != nil {
		for _, k := range p.doc.Kits {
			if p.doc.hasKit(k.Name) {
				continue
			}
			ch := p.kits[k.Name]
			go func() { <-ch }()
		}
	}

	p.doc = doc
	fmt.Printf("Spawning %v kits\n", len(p.doc.Kits))
	for _, k := range p.doc.Kits {
		p.spawnKit(k)
	}
	fmt.Printf("Kits are all spawned")
}

func (p *Player) spawnKit(k *Kit) {
	fmt.Printf("Spawning %s\n", k.Name)
	ch, ok := p.kits[k.Name]
	if !ok {
		ch = make(chan bool)
		p.kits[k.Name] = ch
		go func() { ch <- true }()
	}

	go func() {
		<-ch
		loop := k.loop
		tempo := p.doc.Tempo
		timeTop := p.doc.TimeTop
		sample := p.kit.Sample(k.Sample)
		for {
			select {
			case ch <- true:
				return
			default:
			}

			for _, m := range loop.Measures {
				for _, p := range m.Pulses {
					go func(p float64) {
						waitForBeat(tempo, p)
						sample.Play()
					}(p)
				}
				waitForBeat(tempo, timeTop)
			}
		}
	}()
}

func waitForBeat(bpm float64, count float64) {
	seconds := (60 / float64(bpm)) * count
	elapsed := time.Duration(float64(time.Second) * seconds)
	<-time.NewTimer(elapsed).C
}
