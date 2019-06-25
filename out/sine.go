package out

import (
	"math"
	"time"
)

type SineWave struct {
	Frequency float64
	Length    time.Duration
}

func (s SineWave) Generate() (snd Sound) {
	sampleCount := int(float64(*SampleRate) * float64(s.Length) / float64(time.Second))
	snd.Waveform = make([]float64, sampleCount)

	for i := 0; i < sampleCount; i++ {
		t := 2 * math.Pi * s.Frequency * float64(i) / float64(*SampleRate)
		sin := math.Sin(t)
		snd.Waveform[i] = sin
	}

	return
}
