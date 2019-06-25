package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dianelooney/rave/pitches"

	. "github.com/dianelooney/rave/out"
)

func main() {
	var err error

	if err != nil {
		log.Fatalf("Unable to initialize oto:\n%v", err)
	}

	go Note{"sine", pitches.A4, 2 * time.Second, 0.1}.Play()
	go Note{"sine", pitches.C4, 2 * time.Second, 0.1}.Play()
	go Note{"sine", pitches.F4, 2 * time.Second, 0.1}.Play()
	time.Sleep(2 * time.Second)

	go Note{"sine", pitches.A4, 2 * time.Second, 0.1}.Play()
	go Note{"sine", pitches.D4, 2 * time.Second, 0.1}.Play()
	go Note{"sine", pitches.Fs4, 2 * time.Second, 0.1}.Play()
	time.Sleep(2 * time.Second)
}

type Note struct {
	Pattern   string
	Pitch     float64
	Duration  time.Duration
	Intensity float64
}

func (n Note) Play() {
	p := Ctx.NewPlayer()

	var s Sound
	switch n.Pattern {
	case "sine":
		s = NewSineWave(n.Pitch, n.Duration).Generate()
	default:
		fmt.Fprintf(os.Stderr, "Unsupported wave pattern '%v'\n", n.Pattern)
		return
	}

	s.ScaleAmplitude(n.Intensity)

	if _, err := p.Write(s.ToByteStream()); err != nil {
		log.Fatalf("Unable to write to buffer:\n%v", err)
	}
	if err := p.Close(); err != nil {
		log.Fatalf("Unable to close player:\n%v", err)
	}
}
