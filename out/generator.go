package out

import "time"

type WaveFunc func(offset float64) float64

func Generate(f WaveFunc, freq float64, l time.Duration) (snd Sound) {
	sampleCount := int(float64(*SampleRate) * float64(l) / float64(time.Second))
	snd.Waveform = make([]float64, sampleCount)

	for i := 0; i < sampleCount; i++ {
		snd.Waveform[i] = f(freq * float64(i) / float64(*SampleRate))
	}

	return
}
