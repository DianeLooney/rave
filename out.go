package main

import (
	"log"
	"time"

	. "github.com/dianelooney/rave/notes"
	. "github.com/dianelooney/rave/pitches"
)

func main() {
	var err error

	if err != nil {
		log.Fatalf("Unable to initialize oto:\n%v", err)
	}
	go Note{
		Pitch:     C4,
		Duration:  time.Second * 2,
		Pattern:   "sine",
		Intensity: 0.2,
	}.Play()
	go Note{
		Pitch:     E4,
		Duration:  time.Second * 2,
		Pattern:   "square",
		Intensity: 0.2,
	}.Play()
	go Note{
		Pitch:     G4,
		Duration:  time.Second * 2,
		Pattern:   "triangle",
		Intensity: 0.2,
	}.Play()
	time.Sleep(2 * time.Second)
}
