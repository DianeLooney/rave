package chords

import (
	"math"
	"sync"

	"github.com/dianelooney/rave/pitches"

	. "github.com/dianelooney/rave/notes"
)

type Descriptor struct {
	Name       string
	ShortName  string
	Components []Component
}
type Component struct {
	Offset int
}

var MajorTriad = Descriptor{
	Name:       "Major Triad",
	Components: []Component{{0}, {4}, {7}},
}
var MajorSixth = Descriptor{
	Name:       "Major Sixth",
	Components: []Component{{0}, {4}, {7}, {9}},
}
var DominantSeventh = Descriptor{
	Name:       "Dominant Seventh",
	Components: []Component{{0}, {4}, {7}, {10}},
}
var MajorSeventh = Descriptor{
	Name:       "Major Seventh",
	Components: []Component{{0}, {4}, {7}, {11}},
}
var AugmentedTriad = Descriptor{
	Components: []Component{{0}, {4}, {8}},
}
var AugmentedSeventh = Descriptor{
	Components: []Component{{0}, {4}, {8}, {10}},
}

type Chord struct {
	Base Note
	Descriptor
}

func (c Chord) Play() {
	desc := c.Descriptor
	count := len(desc.Components)
	notes := make([]Note, count)

	for i, d := range desc.Components {
		notes[i] = Note{
			c.Base.Pattern,
			c.Base.Pitch * math.Pow(pitches.HalfStep, float64(d.Offset)),
			c.Base.Duration,
			c.Base.Intensity / float64(count),
		}
	}

	wg := sync.WaitGroup{}
	wg.Add(len(notes))
	for _, n := range notes {
		go func(n Note) {
			defer wg.Done()
			n.Play()
		}(n)
	}
	wg.Wait()
}
