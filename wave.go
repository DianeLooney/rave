package rave

import (
	"fmt"

	"github.com/dianelooney/rave/music"

	"github.com/dianelooney/rave/out"
	"github.com/dianelooney/rave/waves"
)

type Wave struct {
	Name      string
	Sync      string
	Volume    float64
	BaseFreq  float64
	Harmonics []float64
	Pattern   string
	FadeIn    float64
	FadeOut   float64
	chord     music.Chord
	loops     []*WaveLoop
	done      chan bool
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
	l := &WaveLoop{Times: -1}
	w.loops = append(w.loops, l)
	return l
}

func (w *Wave) Harmonic(f float64) {
	w.Harmonics = append(w.Harmonics, f)
}

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
		sync = &ctx.globalBeat
	}
	beat := ctx.beats[w.ID()]

	for _, loop := range w.loops {
		for {
			if loop.playCount >= loop.Times && loop.Times > 0 {
				break
			}
			for i, m := range loop.Measures {
				if i == 0 {
					sync.Wait()
					beat.Done()
				} else {
					ctx.globalBeat.Wait()
				}
				for i, pulse := range m.Pulses {
					go func(m *WaveMeasure, i int, pulse float64) {
						waveform, ok := waves.Patterns[w.Pattern]
						if !ok {
							waveform = waves.Sin
							fmt.Printf("Unrecognized pattern '%s'\n", w.Pattern)
						}
						newWaveform := waveform
						sum := 1.0
						for _, h := range w.Harmonics {
							newWaveform = newWaveform.Add(waveform.Shrink(h).Amplitude(1 / h))
							sum += 1 / h
						}
						waveform = newWaveform.Amplitude(1 / sum)

						waveform = waveform.Amplitude(w.Volume)
						if i < len(m.Weights) {
							waveform = waveform.Amplitude(m.Weights[i])
						}

						var freq = w.BaseFreq
						if i < len(m.Notes) {
							m := w.chord.Multiplier(m.Notes[i])
							freq *= m
						}
						length := 1.0
						if i < len(m.Lengths) {
							length = m.Lengths[i]
						}
						snd := out.Generate(waveform, freq, beatLength(ctx.doc.Tempo, length))
						if w.FadeIn != 0 {
							snd.FadeIn(w.FadeIn)
						}
						if w.FadeOut != 0 {
							snd.FadeOut(w.FadeOut)
						}
						waitForBeat(ctx.doc.Tempo, pulse)
						snd.Play()
					}(m, i, pulse)
				}
			}
			loop.playCount++

			if loop.Times < 0 {
				break
			}
		}
	}
}

type WaveLoop struct {
	playCount int
	Times     int
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
