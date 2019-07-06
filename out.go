package main

import (
	"time"

	"github.com/dianelooney/rave/notes"
)

// . "github.com/dianelooney/rave/notes"
// . "github.com/dianelooney/rave/pitches"

func main() {
	note := notes.Note{
		Pattern:   "experiment",
		Pitch:     440,
		Duration:  time.Second * 2,
		Intensity: 0.1,
	}
	note.Play()
}
