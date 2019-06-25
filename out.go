package main

import (
	"log"
	"time"

	. "github.com/dianelooney/rave/out"
	. "github.com/dianelooney/rave/pitches"
)

func main() {
	var err error

	if err != nil {
		log.Fatalf("Unable to initialize oto:\n%v", err)
	}
	sound := TriangleWave{
		Frequency: C4,
		Length:    time.Second * 2,
	}.Generate()
	sound.ScaleAmplitude(0.1)
	sound.Play()
	sound = SquareWave{
		Frequency: C4,
		Length:    time.Second * 2,
	}.Generate()
	sound.ScaleAmplitude(0.1)
	sound.Play()
	time.Sleep(time.Second)
}
