package notes

import (
	"fmt"
	"log"
	"os"
	"time"

	. "github.com/dianelooney/rave/out"
)

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
		s = SineWave{
			Frequency: n.Pitch,
			Length:    n.Duration,
		}.Generate()
	case "triangle":
		s = TriangleWave{
			Frequency: n.Pitch,
			Length:    n.Duration,
		}.Generate()
	case "square":
		s = SquareWave{
			Frequency: n.Pitch,
			Length:    n.Duration,
		}.Generate()
	default:
		fmt.Fprintf(os.Stderr, "Unsupported wave pattern '%v'\n", n.Pattern)
		return
	}

	s.FadeIn(0.05)
	s.ScaleAmplitude(n.Intensity)
	s.TaperOff(0.05)

	if _, err := p.Write(s.ToByteStream()); err != nil {
		log.Fatalf("Unable to write to buffer:\n%v", err)
	}
	if err := p.Close(); err != nil {
		log.Fatalf("Unable to close player:\n%v", err)
	}
}
