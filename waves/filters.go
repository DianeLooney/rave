package waves

import (
	"math/rand"
)

func Add(w1, w2 WaveFunc) WaveFunc {
	return func(offset float64) float64 {
		return (w1(offset) + w2(offset)) / 2
	}
}

func Mult(w1, w2 WaveFunc) WaveFunc {
	return func(offset float64) float64 {
		return (w1(offset) * w2(offset))
	}
}

func WhiteNoise(w1, w2 WaveFunc) WaveFunc {
	return func(offset float64) float64 {
		r := rand.Float64() / 2
		return (0.25+r)*w1(offset) + (0.75-r)*w2(offset)
	}
}
