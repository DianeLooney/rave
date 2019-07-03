package waves

import (
	"math"
)

type WaveFunc func(offset float64) float64

func (w WaveFunc) ShiftX(x float64) WaveFunc {
	return func(offset float64) float64 {
		return w(offset + x)
	}
}

func (w WaveFunc) ShiftY(y float64) WaveFunc {
	return func(offset float64) float64 {
		return w(offset) + y
	}
}

func (w WaveFunc) Compose(w2 WaveFunc) WaveFunc {
	return func(offset float64) float64 {
		return w(w2(offset))
	}
}

func (w WaveFunc) Expand(r float64) WaveFunc {
	return func(offset float64) float64 {
		return w(offset / r)
	}
}

func (w WaveFunc) Shrink(r float64) WaveFunc {
	return func(offset float64) float64 {
		return w(offset * r)
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
	if int(2*offset)%2 == 0 {
		return 1
	}

	return -1
}

var Nil WaveFunc = func(offset float64) float64 { return 0 }

//var Experiment = Sin.Expand(1)
//var Experiment = Sin.Mult(Sin.Amplitude(0.1).ShiftY(0.6).Shrink(9))

var Experiment WaveFunc

func init() {
	p1 := Triangle.Mult(Square).Expand(0.5).Amplitude(0.01)
	p2 := Sin.Shrink(2).Amplitude(0.1)
	p3 := Square.Shrink(12).Amplitude(0.07)
	Experiment = p1.Add(p2).Add(p3)
	/*
		p1 := Sin.Shrink(3)
		p2 := Sin.Shrink(7)
		p3 := p1.Add(p2).Amplitude(0.3).ShiftY(0.6)
		Experiment = Sin.Mult(p3).Amplitude(0.6)
	*/
}
