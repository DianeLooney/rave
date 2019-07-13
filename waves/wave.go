package waves

import (
	"math"

	"github.com/dianelooney/rave/common"
)

type Amplitude float64
type Pulse float64
type Time float64

type PreFilter interface {
	Init(d Descriptor)
	Apply(c1 Pulse, t1 Time) (c2 Pulse, t2 Time)
}
type Generator func(c Pulse) (a Amplitude)
type PostFilter interface {
	Init(d Descriptor)
	Apply(a1 Amplitude, t1 Time) (a2 Amplitude, t2 Time)
}

// A Descriptor gives all the necessary information to create a sound
type Descriptor struct {
	Pipeline  Pipeline
	Frequency float64
	Duration  Time
}

// Pipeline represents a series of filters and generators used to create a sound
type Pipeline struct {
	PreFilters  []PreFilter
	Generator   Generator
	PostFilters []PostFilter
}

func pair(i int, freq float64) [2]float64 {
	x := float64(i) / float64(*common.SampleRate)
	return [2]float64{x * freq, x}
}

// Generate executes the pipeline to create a sound
func (d Descriptor) Generate() (snd common.Sound) {
	for _, f := range d.Pipeline.PreFilters {
		f.Init(d)
	}
	for _, f := range d.Pipeline.PostFilters {
		f.Init(d)
	}

	sampleCount := int(float64(*common.SampleRate) * float64(d.Duration))
	store := make([][2]float64, sampleCount)

	for i := 0; i < sampleCount; i++ {
		store[i] = pair(i, d.Frequency)
	}
	for _, f := range d.Pipeline.PreFilters {
		for i, p := range store {
			x, y := f.Apply(Pulse(p[0]), Time(p[1]))
			store[i] = [2]float64{float64(x), float64(y)}
		}
	}
	for i, p := range store {
		x := d.Pipeline.Generator(Pulse(p[0]))
		store[i][0] = float64(x)
	}
	for _, f := range d.Pipeline.PostFilters {
		for i, p := range store {
			x, y := f.Apply(Amplitude(p[0]), Time(p[1]))
			store[i] = [2]float64{float64(x), float64(y)}
		}
	}

	snd.Waveform = make([]float64, sampleCount)
	for i, f := range store {
		snd.Waveform[i] = f[0]
	}

	return
}

var PreFilters = map[string]PreFilter{}
var PostFilters = map[string]PostFilter{}

// Generators includes the built-in generator functions
var Generators = map[string]Generator{
	"sin": func(p Pulse) (a Amplitude) {
		return Amplitude(math.Sin(2 * math.Pi * float64(p)))
	},
	"square": func(p Pulse) (a Amplitude) {
		if int(2*p)%2 == 0 {
			a = 1
		} else {
			a = -1
		}
		return a
	},
	"triangle": func(p Pulse) (a Amplitude) {
		iPart := math.Floor(float64(p))
		remainder := float64(p) - iPart
		if remainder < 0.5 {
			return Amplitude(1 - 2*remainder)
		}
		return Amplitude(-2 + 2*remainder)
	},
	"saw": func(p Pulse) (a Amplitude) {
		iPart := math.Floor(float64(p))
		remainder := float64(p) - iPart
		return Amplitude(-1 + 2*remainder)
	},
}

// ScaleAmplitude does what it says
type ScaleAmplitude struct {
	R float64
}

func (f *ScaleAmplitude) Init(d Descriptor) {}

// Apply applies the filter
func (f *ScaleAmplitude) Apply(a1 Amplitude, t1 Time) (a2 Amplitude, t2 Time) {
	return Amplitude(float64(a1) * f.R), t1
}

// FadeIn is a PostFilter that fades the clip in over time
type FadeIn struct {
	Over float64
}

// Init initializes the filter
func (f *FadeIn) Init(d Descriptor) {}

// Apply applies the filter
func (f FadeIn) Apply(a1 Amplitude, t1 Time) (a2 Amplitude, t2 Time) {
	if t1 < Time(f.Over) {
		return Amplitude(float64(a1) * float64(t1) / f.Over), t1
	}
	return a1, t2
}

// FadeOut is a PostFilter that fades the clip in over time
type FadeOut struct {
	startAt Time
	Over    float64
}

// Init initializes the filter
func (f *FadeOut) Init(d Descriptor) {
	f.startAt = d.Duration - Time(f.Over)
}

// Apply applies the filter
func (f *FadeOut) Apply(a1 Amplitude, t1 Time) (a2 Amplitude, t2 Time) {
	return Amplitude(float64(a1) * float64(t1) / f.Over), t1
}

type Bend struct {
	R     float64
	Start float64
	End   float64

	p1 Pulse
	p2 Pulse
}

func (f *Bend) Init(d Descriptor) {
	f.p1 = 0
	f.p2 = 0
}
func (f *Bend) Apply(c1 Pulse, t1 Time) (c2 Pulse, t2 Time) {
	if t1 < Time(f.Start) {
		return c1, t1
	}
	r := f.R
	if t2 < Time(f.End) {
		r += (1 - r) * (f.End - float64(t1)) / (f.End - f.Start)
	}
	inDiff := c1 - f.p1
	outDiff := inDiff * Pulse(r)
	f.p1 = c1
	f.p2 += outDiff
	return f.p2, t1
}

type Trill struct {
	Hz float64
	R  float64

	p1 Pulse
	p2 Pulse
}

func (f *Trill) Init(d Descriptor) {
	f.p1 = 0
	f.p2 = 0
}

func (f *Trill) Apply(c1 Pulse, t1 Time) (c2 Pulse, t2 Time) {
	n := int(float64(t1) * f.Hz)
	var r float64
	if n%2 == 0 {
		r = 1
	} else {
		r = f.R
	}
	inDiff := c1 - f.p1
	outDiff := inDiff * Pulse(r)
	f.p1 = c1
	f.p2 += outDiff
	return f.p2, t1
}
