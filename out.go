package main

import (
	"time"

	"github.com/dianelooney/rave/music"
	"github.com/dianelooney/rave/notes"
)

func main() {
	set := NewContext(Background)
	set.Set(Tempo, 120)

	music.Measure{
		TimeSignature: TimeSignature{4, 4},
		Tempo:         120,
		Beats:         []Beat{},
	}

	note := notes.Note{
		Pattern:   "experiment",
		Pitch:     440,
		Duration:  time.Second * 2,
		Intensity: 0.1,
	}
	note.Play()
}
