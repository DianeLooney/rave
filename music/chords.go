package music

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
)

type Scale []float64

var HalfStep = math.Pow(2, 1.0/12)

func (c Scale) Multiplier(s string) float64 {
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
	if s == "?" {
		i = rand.Intn(len(c))
		modifier := rand.Intn(8)
		if modifier == 0 {
			i--
		} else if modifier == 1 {
			i++
		}
	}
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
