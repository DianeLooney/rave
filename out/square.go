package out

import (
	"time"
)

type SquareWave struct {
	Frequency float64
	Length    time.Duration
}

func (s SquareWave) Generate() (snd Sound) {
	sampleCount := int(float64(*SampleRate) * float64(s.Length) / float64(time.Second))
	snd.Waveform = make([]float64, sampleCount)

	for i := 0; i < sampleCount; i++ {
		t := s.Frequency * float64(i) / float64(*SampleRate)
		if int(t)%2 == 0 {
			snd.Waveform[i] = 1
		} else {
			snd.Waveform[i] = -1
		}
	}

	return
}
