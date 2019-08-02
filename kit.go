package rave

import (
	"time"

	"github.com/dianelooney/rave/common"
)

type Kit struct {
	Name    string
	Samples []string
	Sync    string
	Volume  float64
	Poly    int

	loop *KitLoop
	done chan bool
}

func (k *Kit) ID() string {
	return k.Name
}

func (k *Kit) Done() chan bool {
	return k.done
}

func (k *Kit) SyncID() string {
	return k.Sync
}

func (k *Kit) Sample(s string) {
	k.Samples = append(k.Samples, s)
}

func (k *Kit) Loop() *KitLoop {
	if k.loop == nil {
		k.loop = &KitLoop{}
	}

	return k.loop
}
func (k *Kit) playPolyLoop(ctx *Context) {
	samples := make([]common.Sound, len(k.Samples))
	for i, s := range k.Samples {
		x := kit.Sample(s)
		samples[i] = x.ScaleAmplitude(k.Volume)
	}
	m := k.loop.Measures[0]
	mStart := ctx.globalTick.Wait().(time.Time)
	for i, pulse := range m.Pulses {
		go func(m *KitMeasure, i int, pulse float64) {
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
			waitForBeat(mStart, ctx.doc.Tempo, pulse)
			x.Play()
		}(m, i, pulse)
	}

	for i := 1; i < int(m.Size); i++ {
		ctx.globalTick.Wait()
	}
}
func (k *Kit) PlayLoop(ctx *Context) {
	if k.Poly != 0 {
		k.playPolyLoop(ctx)
		return
	}

	samples := make([]common.Sound, len(k.Samples))
	for i, s := range k.Samples {
		x := kit.Sample(s)
		samples[i] = x.ScaleAmplitude(k.Volume)
	}
	sync, ok := ctx.beats[k.SyncID()]
	if !ok {
		sync = ctx.globalBeat
	}
	beat := ctx.beats[k.ID()]
	for i, m := range k.loop.Measures {
		var mStart time.Time
		if i == 0 {
			v := sync.Wait()
			mStart = v.(time.Time)
			beat.Done(mStart)
		} else {
			mStart = ctx.globalBeat.Wait().(time.Time)
		}

		for i, pulse := range m.Pulses {
			go func(m *KitMeasure, i int, pulse float64) {
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
				waitForBeat(mStart, ctx.doc.Tempo, pulse)
				x.Play()
			}(m, i, pulse)
		}

		for i := 1; i < int(m.Size); i++ {
			ctx.globalBeat.Wait()
		}
	}
}

type KitLoop struct {
	Measures []*KitMeasure
}

func (l *KitLoop) Measure() *KitMeasure {
	m := &KitMeasure{Size: 1}
	l.Measures = append(l.Measures, m)
	return m
}

type KitMeasure struct {
	Size    float64
	Weights []float64
	Samples []int
	Pulses  []float64
}

func (m *KitMeasure) Pulse(t float64) {
	m.Pulses = append(m.Pulses, t)
}

func (m *KitMeasure) Sample(t int) {
	m.Samples = append(m.Samples, t)
}

func (m *KitMeasure) Weight(f float64) {
	m.Weights = append(m.Weights, f)
}
