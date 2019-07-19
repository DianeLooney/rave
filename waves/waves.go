package waves

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

//var Experiment = Sin.Expand(1)
//var Experiment = Sin.Mult(Sin.Amplitude(0.1).ShiftY(0.6).Shrink(9))
// var Experiment = Sin.Compose(Bend(0, 2200, 1/(music.HalfStep*music.HalfStep)))

// var Experiment WaveFunc

func init() {
	/*
		p1 := Triangle.Mult(Square).Expand(0.5).Amplitude(0.01)
		p2 := Sin.Shrink(2).Amplitude(0.1)
		p3 := Square.Shrink(12).Amplitude(0.07)
		Experiment = p1.Add(p2).Add(p3)
	*/
	/*
		p1 := Sin.Shrink(3)
		p2 := Sin.Shrink(5)
		p3 := Sin.Shrink(7)
		p4 := p1.Add(p2).Add(p3).Amplitude(0.3).ShiftY(0.6)
		Experiment = Sin.Mult(p4).Amplitude(0.6)
	*/
	/*
		p1 := Triangle.Mult(Square).Expand(0.5).Amplitude(0.01)
		p2 := Sin.Shrink(2).Amplitude(0.1)
		p3 := Square.Shrink(12).Amplitude(0.07)
		Experiment = p1.Add(p2).Add(p3)
	*/
}
