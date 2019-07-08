package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"github.com/dianelooney/directive"
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
	globalBeat     chan bool
	globalBeatOnce sync.Once
	doc            *Doc
	kits           map[string]chan bool
	beats          map[string]chan bool
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

func awakeAll(ch chan bool) {
	for {
		select {
		case ch <- true:
			continue
		default:
			return
		}
	}
}

func (p *Player) Init() {
	p.globalBeat = make(chan bool)
	p.kits = make(map[string]chan bool)
	p.kit = kits.LoadManifest("kits/manifest.yml").Load()
	p.beats = make(map[string]chan bool)
}

func (p *Player) Load(doc *Doc) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if p.doc != nil {
		for _, k := range p.doc.Kits {
			if doc.hasKit(k.Name) {
				continue
			}
			ch := p.kits[k.Name]
			go func() {
				<-ch
				p.mtx.Lock()
				defer p.mtx.Unlock()
				delete(p.kits, k.Name)
			}()
		}
	}

	p.doc = doc
	for _, k := range p.doc.Kits {
		p.spawnKit(k)
	}
	go p.globalBeatOnce.Do(func() {
		for {
			awakeAll(p.globalBeat)
			waitForBeat(p.doc.Tempo, p.doc.TimeBot)
		}
	})
}

func (p *Player) spawnKit(k *Kit) {
	ch, ok := p.kits[k.Name]
	if !ok {
		ch = make(chan bool)
		p.kits[k.Name] = ch
		go func() {
			b, ok := p.beats[k.Sync]
			if ok {
				<-b
			} else {
				<-p.globalBeat
			}
			ch <- true
		}()
	}
	beat, ok := p.beats[k.Name]
	if !ok {
		beat = make(chan bool)
		p.beats[k.Name] = beat
	}

	go func() {
		<-ch
		loop := k.loop
		tempo := p.doc.Tempo
		timeTop := p.doc.TimeTop
		sample := p.kit.Sample(k.Sample)
		sample = sample.ScaleAmplitude(k.Volume)
		for {
			select {
			case ch <- true:
				return
			default:
				awakeAll(beat)
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
