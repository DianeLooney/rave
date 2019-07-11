package music

import (
	"math"
	"strconv"
	"strings"
)

type Chord []float64

var HalfStep = math.Pow(2, 1.0/12)

func (c Chord) Multiplier(s string) float64 {
	var sharpness float64
	for {
		if strings.HasSuffix(s, "#") {
			s = strings.TrimSuffix(s, "#")
			sharpness++
		} else if strings.HasSuffix(s, "b") {
			s = strings.TrimSuffix(s, "b")
			sharpness--
		} else {
			break
		}
	}

	i, _ := strconv.Atoi(s)
	var x float64
	if len(c) != 0 {
		for i < 0 {
			i += len(c)
			x -= 12
		}
		for i >= len(c) {
			i -= len(c)
			x += 12
		}
		x += c[i]
	}
	return math.Pow(HalfStep, x+sharpness)
}
