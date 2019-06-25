package main

import (
	"fmt"
	"github.com/dianelooney/rave/pitches"
	"io"
	"log"
	"os"
	"time"

	. "github.com/dianelooney/rave/out"
)

func main() {
	var err error

	if err != nil {
		log.Fatalf("Unable to initialize oto:\n%v", err)
	}

	go func() {
		Note{"sine", pitches.A4, 6 * time.Second, 0.0001}.Play()
	}()

	go Note{"sine", pitches.A4, 2 * time.Second, 0.3}.Play()
	go Note{"sine", pitches.C4, 2 * time.Second, 0.3}.Play()
	go Note{"sine", pitches.E4, 2 * time.Second, 0.3}.Play()
	time.Sleep(2 * time.Second)

	go Note{"sine", pitches.A4, 2 * time.Second, 0.3}.Play()
	go Note{"sine", pitches.C4, 2 * time.Second, 0.3}.Play()
	go Note{"sine", pitches.F4, 2 * time.Second, 0.3}.Play()
	time.Sleep(2 * time.Second)

	go Note{"sine", pitches.A4, 2 * time.Second, 0.3}.Play()
	go Note{"sine", pitches.D4, 2 * time.Second, 0.3}.Play()
	go Note{"sine", pitches.Fs4, 2 * time.Second, 0.3}.Play()
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

	var s io.Reader
	switch n.Pattern {
	case "sine":
		s = NewSineWave(n.Pitch, n.Duration, n.Intensity)
	default:
		fmt.Fprintf(os.Stderr, "Unsupported wave pattern '%v'\n", n.Pattern)
		return
	}

	if _, err := io.Copy(p, s); err != nil {
		log.Fatalf("Unable to copy to buffer:\n%v", err)
	}
	if err := p.Close(); err != nil {
		log.Fatalf("Unable to close player:\n%v", err)
	}
}
