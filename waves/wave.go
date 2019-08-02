package waves

import (
	"math"

	"github.com/dianelooney/rave/common"
)

type Amplitude float64
type Pulse float64
type Time float64

type PreFilter interface {
	Init(d Descriptor) PreFilter
	Apply(c1 Pulse, t1 Time) (c2 Pulse, t2 Time)
}
type Generator func(c Pulse) (a Amplitude)
type PostFilter interface {
	Init(d Descriptor) PostFilter
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

func (p Pipeline) Chordify(f float64) Pipeline {
	return Pipeline{
		PreFilters: append(
			p.PreFilters,
		),
		Generator: p.Generator,
		PostFilters: append(
			p.PostFilters,
			ScaleAmplitude{R: 1 / f},
		),
	}
}

func pair(i int, freq float64) [2]float64 {
	x := float64(i) / float64(*common.SampleRate)
	return [2]float64{x * freq, x}
}

// Generate executes the pipeline to create a sound
func (d Descriptor) Generate() (snd common.Sound) {
	pres := make([]PreFilter, len(d.Pipeline.PreFilters))
	for i, f := range d.Pipeline.PreFilters {
		pres[i] = f.Init(d)
	}
	posts := make([]PostFilter, len(d.Pipeline.PostFilters))
	for i, f := range d.Pipeline.PostFilters {
		posts[i] = f.Init(d)
	}

	sampleCount := int(float64(*common.SampleRate) * float64(d.Duration))
	store := make([][2]float64, sampleCount)

	for i := 0; i < sampleCount; i++ {
		store[i] = pair(i, d.Frequency)
	}
	for _, f := range pres {
		for i, p := range store {
			x, y := f.Apply(Pulse(p[0]), Time(p[1]))
			store[i] = [2]float64{float64(x), float64(y)}
		}
	}
	for i, p := range store {
		x := d.Pipeline.Generator(Pulse(p[0]))
		store[i][0] = float64(x)
	}
	for _, f := range posts {
		for i, p := range store {
			x, y := f.Apply(Amplitude(p[0]), Time(p[1]))
			store[i] = [2]float64{float64(x), float64(y)}
		}
	}

	snd.Waveform = make([]float64, sampleCount)
	for i, f := range store {
		snd.Waveform[i] = f[0] / 10
	}

	return
}

var PreFilters = map[string]PreFilter{}
var PostFilters = map[string]PostFilter{}

func Sin(p Pulse) (a Amplitude) {
	return Amplitude(math.Sin(2 * math.Pi * float64(p)))
}
func Square(p Pulse) (a Amplitude) {
	if int(2*p)%2 == 0 {
		a = 1
	} else {
		a = -1
	}
	return a
}
func Triangle(p Pulse) (a Amplitude) {
	iPart := math.Floor(float64(p))
	remainder := float64(p) - iPart
	if remainder < 0.5 {
		return Amplitude(1 - 2*remainder)
	}
	return Amplitude(-2 + 2*remainder)
}
func Saw(p Pulse) (a Amplitude) {
	iPart := math.Floor(float64(p))
	remainder := float64(p) - iPart
	return Amplitude(-1 + 2*remainder)
}
func Experiment1(p Pulse) (a Amplitude) {
	p1 := Saw(p*2) * 0.2
	p2 := Sin(p) * 0.7
	p3 := Square(p*3) * 0.1

	return p1 + p2 + p3
}
func Organ(p Pulse) (a Amplitude) {
	return Sin(p) * ((Sin(3*p)+Sin(5*p)+Sin(7*p))*0.3 - 0.6) * 0.6
}
func Experiment2(p Pulse) (a Amplitude) {
	p1 := Triangle(p) * Square(p) * 0.3
	p2 := Sin(p/4) * 0.2
	p3 := Square(p) * 0.5
	return p1 + p2 + p3
}

// Generators includes the built-in generator functions
var Generators = map[string]Generator{
	"organ":    Organ,
	"sin":      Sin,
	"square":   Square,
	"triangle": Triangle,
	"saw":      Saw,
	"exp1":     Experiment1,
	"exp2":     Experiment2,
}

type PitchUp struct {
	R float64
}

func (f PitchUp) Init(d Descriptor) PreFilter {
	return f
}

func (f PitchUp) Apply(c1 Pulse, t1 Time) (c2 Pulse, t2 Time) {
	return (c1 * Pulse(f.R)), t1
}

type Vibrato struct {
	R  float64
	Hz float64

	c1 *Pulse
	c2 *Pulse
}

func (f Vibrato) Init(d Descriptor) PreFilter {
	var c1, c2 Pulse
	f.c1 = &c1
	f.c2 = &c2
	return f
}

func (f Vibrato) Apply(c1 Pulse, t1 Time) (c2 Pulse, t2 Time) {
	r := 1 + f.R*math.Sin(2*math.Pi*float64(t1)*f.Hz)
	diffIn := c1 - *f.c1
	diffOut := diffIn * Pulse(r)
	*f.c1 = c1
	*f.c2 += diffOut
	return *f.c2, t1
}

// ScaleAmplitude does what it says
type ScaleAmplitude struct {
	R float64
}

func (f ScaleAmplitude) Init(d Descriptor) PostFilter {
	return f
}

// Apply applies the filter
func (f ScaleAmplitude) Apply(a1 Amplitude, t1 Time) (a2 Amplitude, t2 Time) {
	return Amplitude(float64(a1) * f.R), t1
}

// FadeIn is a PostFilter that fades the clip in over time
type FadeIn struct {
	Over float64
}

// Init initializes the filter
func (f FadeIn) Init(d Descriptor) PostFilter {
	return f
}

// Apply applies the filter
func (f FadeIn) Apply(a1 Amplitude, t1 Time) (a2 Amplitude, t2 Time) {
	if t1 < Time(f.Over/100) {
		return Amplitude(float64(a1) * float64(t1) / (f.Over / 100)), t1
	}
	return a1, t1
}

// FadeOut is a PostFilter that fades the clip in over time
type FadeOut struct {
	startAt Time
	Over    float64
}

// Init initializes the filter
func (f FadeOut) Init(d Descriptor) PostFilter {
	f.startAt = d.Duration - Time(f.Over/100)
	return f
}

// Apply applies the filter
func (f FadeOut) Apply(a1 Amplitude, t1 Time) (a2 Amplitude, t2 Time) {
	if t1 < f.startAt {
		return a1, t2
	}
	fac := (f.startAt + Time(f.Over/100) - t1) / Time(f.Over/100)
	return Amplitude(float64(fac) * float64(a1)), t1
}

type Bend struct {
	R     float64
	Start float64
	End   float64

	p1 *Pulse
	p2 *Pulse
}

func (f Bend) Init(d Descriptor) PreFilter {
	var p1, p2 Pulse
	f.p1, f.p2 = &p1, &p2
	return f
}
func (f Bend) Apply(c1 Pulse, t1 Time) (c2 Pulse, t2 Time) {
	if t1 < Time(f.Start) {
		*f.p1 = c1
		*f.p2 = c1
		return c1, t1
	}

	x := 1.0
	if t1 < Time(f.End) {
		x = (float64(t1) - f.Start) / (f.End - f.Start)
	}
	r := f.R*x + (1.0)*(1-x)
	inDiff := c1 - *f.p1
	outDiff := inDiff * Pulse(r)
	*f.p1 = c1
	*f.p2 += outDiff
	return *f.p2, t1
}

type Trill struct {
	Hz float64
	R  float64

	p1 *Pulse
	p2 *Pulse
}

func (f Trill) Init(d Descriptor) PreFilter {
	var p1, p2 Pulse
	f.p1, f.p2 = &p1, &p2
	return f
}

func (f Trill) Apply(c1 Pulse, t1 Time) (c2 Pulse, t2 Time) {
	n := int(float64(t1) * f.Hz)
	var r float64
	if n%2 == 0 {
		r = 1
	} else {
		r = f.R
	}
	inDiff := c1 - *f.p1
	outDiff := inDiff * Pulse(r)
	*f.p1 = c1
	*f.p2 += outDiff
	return *f.p2, t1
}
