package rave

import (
	"fmt"
	"time"

	"github.com/dianelooney/rave/music"
	"github.com/dianelooney/rave/waves"
)

type Wave struct {
	Name     string
	Sync     string
	BaseFreq float64
	Pattern  string
	chord    music.Chord
	loops    []*WaveLoop
	done     chan bool
	pipe     *Pipe
}

func (w *Wave) Pipe() *Pipe {
	if w.pipe == nil {
		w.pipe = &Pipe{}
	}
	return w.pipe
}

type Pipe struct {
	Pipeline waves.Pipeline
}

func (p *Pipe) Pre(s string) {
	f, ok := waves.PreFilters[s]
	if !ok {
		fmt.Printf("Unable to find pre-filter '%s'\n", s)
		return
	}
	p.Pipeline.PreFilters = append(p.Pipeline.PreFilters, f)
}
func (p *Pipe) Wave(s string) {
	g, ok := waves.Generators[s]
	if !ok {
		fmt.Printf("Unable to find generator '%s'\n", s)
		return
	}
	p.Pipeline.Generator = g
}

func (p *Pipe) Vibrato() *waves.Vibrato {
	b := &waves.Vibrato{}
	p.Pipeline.PreFilters = append(p.Pipeline.PreFilters, b)
	return b
}

func (p *Pipe) Trill() *waves.Trill {
	b := &waves.Trill{}
	p.Pipeline.PreFilters = append(p.Pipeline.PreFilters, b)
	return b
}
func (p *Pipe) Bend() *waves.Bend {
	b := &waves.Bend{}
	p.Pipeline.PreFilters = append(p.Pipeline.PreFilters, b)
	return b
}

func (p *Pipe) FadeIn() *waves.FadeIn {
	f := &waves.FadeIn{}
	p.Pipeline.PostFilters = append(p.Pipeline.PostFilters, f)
	return f
}

func (p *Pipe) FadeOut() *waves.FadeOut {
	f := &waves.FadeOut{}
	p.Pipeline.PostFilters = append(p.Pipeline.PostFilters, f)
	return f
}

func (p *Pipe) ScaleAmplitude() *waves.ScaleAmplitude {
	f := &waves.ScaleAmplitude{}
	p.Pipeline.PostFilters = append(p.Pipeline.PostFilters, f)
	return f
}

func (w *Wave) ID() string {
	return w.Name
}

func (w *Wave) Done() chan bool {
	return w.done
}

func (w *Wave) SyncID() string {
	return w.Sync
}

func (w *Wave) Loop() *WaveLoop {
	l := &WaveLoop{EndTimes: -1}
	w.loops = append(w.loops, l)
	return l
}

/*
func (w *Wave) Harmonic(f float64) {
	w.Harmonics = append(w.Harmonics, f)
}
*/

var predefinedChords = map[string]music.Chord{
	"major": {0, 2, 4, 5, 7, 9, 11},
	"minor": {0, 2, 3, 5, 7, 8, 10},
}

func (w *Wave) Chord(s string) {
	c, ok := predefinedChords[s]
	if !ok {
		return
	}
	w.chord = c
}

func (w *Wave) PlayLoop(ctx *Context) {
	sync, ok := ctx.beats[w.SyncID()]
	if !ok {
		sync = ctx.globalBeat
	}
	beat := ctx.beats[w.ID()]

	for _, loop := range w.loops {
		for {
			if loop.playCount >= loop.EndTimes && loop.EndTimes > 0 {
				break
			}
			for i, m := range loop.Measures {
				var mStart time.Time
				if i == 0 {
					mStart = sync.Wait().(time.Time)
					beat.Done(mStart)
				} else {
					mStart = ctx.globalBeat.Wait().(time.Time)
				}

				for i, pulse := range m.Pulses {
					go func(m *WaveMeasure, i int, pulse float64) {
						length := 1.0
						if i < len(m.Lengths) {
							length = m.Lengths[i]
						}
						var freq = w.BaseFreq
						if i < len(m.Notes) {
							m := w.chord.Multiplier(m.Notes[i])
							freq *= m
						}
						t := waves.Time(length * (60 / ctx.doc.Tempo))
						desc := waves.Descriptor{
							Pipeline:  w.pipe.Pipeline,
							Frequency: freq,
							Duration:  t,
						}
						snd := desc.Generate()

						waitForBeat(mStart, ctx.doc.Tempo, pulse)
						snd.Play()
					}(m, i, pulse)
				}
			}
			loop.playCount++

			if loop.EndTimes < 0 {
				break
			}
		}
	}
}

type WaveLoop struct {
	playCount int
	EndTimes  int
	Measures  []*WaveMeasure
}

func (l *WaveLoop) Measure() *WaveMeasure {
	m := &WaveMeasure{}
	l.Measures = append(l.Measures, m)
	return m
}

type WaveMeasure struct {
	Weights []float64
	Notes   []string
	Lengths []float64
	Pulses  []float64
}

func (m *WaveMeasure) Len(f float64) {
	m.Lengths = append(m.Lengths, f)
}
func (m *WaveMeasure) Pulse(t float64) {
	m.Pulses = append(m.Pulses, t)
}

func (m *WaveMeasure) Note(t string) {
	m.Notes = append(m.Notes, t)
}

func (m *WaveMeasure) Weight(f float64) {
	m.Weights = append(m.Weights, f)
}
