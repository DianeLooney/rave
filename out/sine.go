package out

import (
	"math"
	"time"
)

type SineWave struct {
	freq   float64
	length time.Duration
}

func NewSineWave(freq float64, duration time.Duration) *SineWave {
	return &SineWave{
		freq:   freq,
		length: duration,
	}
}

func (s *SineWave) Generate() (snd Sound) {
	sampleCount := int(float64(*SampleRate) * float64(s.length) / float64(time.Second))
	snd.Waveform = make([]float64, sampleCount)

	for i := 0; i < sampleCount; i++ {
		t := 2 * math.Pi * s.freq * float64(i) / float64(*SampleRate)
		sin := math.Sin(t)
		snd.Waveform[i] = sin
	}

	return
}
