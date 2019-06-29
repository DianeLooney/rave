package waves

import (
	"math"
)

type WaveFunc func(offset float64) float64

func (w WaveFunc) Shift(x float64) WaveFunc {
	return func(offset float64) float64 {
		return w(offset + x)
	}
}

func (w WaveFunc) Expand(r float64) WaveFunc {
	return func(offset float64) float64 {
		return w(offset / r)
	}
}

var Triangle WaveFunc = func(offset float64) float64 {
	div := int(offset) % 2
	rem := offset - float64(int(offset))
	if div == 0 {
		return -1 + 2*rem
	}

	return 1 - 2*rem
}

var Sin WaveFunc = func(offset float64) float64 {
	return math.Sin(2 * math.Pi * offset)
}

var Square WaveFunc = func(offset float64) float64 {
	if int(offset)%2 == 0 {
		return 1
	}

	return -1
}

//var Experiment = Sin.Expand(1)

var Experiment = Mult(Square, Sin.Expand(2))
