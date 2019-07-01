package out

import (
	"time"

	. "github.com/dianelooney/rave/common"
	"github.com/dianelooney/rave/waves"
)

func Generate(f waves.WaveFunc, freq float64, l time.Duration) (snd Sound) {
	sampleCount := int(float64(*SampleRate) * float64(l) / float64(time.Second))
	snd.Waveform = make([]float64, sampleCount)

	for i := 0; i < sampleCount; i++ {
		snd.Waveform[i] = f(freq * float64(i) / float64(*SampleRate))
	}

	return
}
