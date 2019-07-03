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

func (w WaveFunc) Amplitude(r float64) WaveFunc {
	return func(offset float64) float64 {
		return w(offset) * r
	}
}

func (w WaveFunc) Add(w2 WaveFunc) WaveFunc {
	return func(offset float64) float64 {
		return w(offset) + w2(offset)
	}
}

func (w WaveFunc) Mult(w2 WaveFunc) WaveFunc {
	return func(offset float64) float64 {
		return (w(offset) * w2(offset))
	}
}

var Triangle = Saw(0.25)

func Saw(peak float64) WaveFunc {
	return func(offset float64) float64 {
		iPart := math.Floor(offset)
		remainder := offset - iPart
		if remainder < peak {
			return (1.0 / peak) * remainder
		}
		if remainder > 1-peak {
			return -1 + (remainder-(1-peak))*(1/peak)
		}
		return 1 + (-2/(1-2*peak))*(remainder-peak)
	}
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

var Experiment = Triangle
