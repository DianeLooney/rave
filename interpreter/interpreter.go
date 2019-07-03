package interpreter

import "github.com/dianelooney/rave/music"

type Set struct {
	Tstop      float64
	Tsbot      float64
	Tempo      float64
	BeatedNote float64
	Kits       []Inst
}

func (s Set) BPM() music.BPM {
	if s.BeatedNote == 0 {
		return music.BPM(s.Tempo)
	}
	return music.BPM(s.Tempo * s.Tsbot / s.BeatedNote)
}

type Inst struct {
	N string
	M [][]float64
}
