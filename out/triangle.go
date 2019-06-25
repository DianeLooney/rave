package out

import (
	"time"
)

type TriangleWave struct {
	Frequency float64
	Length    time.Duration
}

func (s TriangleWave) Generate() (snd Sound) {
	sampleCount := int(float64(*SampleRate) * float64(s.Length) / float64(time.Second))
	snd.Waveform = make([]float64, sampleCount)

	for i := 0; i < sampleCount; i++ {
		t := 2 * s.Frequency * float64(i) / float64(*SampleRate)
		div := int(t) % 2
		rem := t - float64(int(t))
		if div == 0 {
			snd.Waveform[i] = -1 + 2*rem
		} else {
			snd.Waveform[i] = 1 - 2*rem
		}
	}

	return
}
