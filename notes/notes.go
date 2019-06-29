package notes

import (
	"fmt"
	"log"
	"os"
	"time"

	. "github.com/dianelooney/rave/out"
	"github.com/dianelooney/rave/waves"
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
		s = Generate(
			waves.Sin,
			n.Pitch,
			n.Duration,
		)
	case "triangle":
		s = Generate(
			waves.Triangle,
			n.Pitch,
			n.Duration,
		)
	case "square":
		s = Generate(
			waves.Square,
			n.Pitch,
			n.Duration,
		)
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
