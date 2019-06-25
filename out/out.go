package out

import (
	"flag"
	"log"

	"github.com/hajimehoshi/oto"
)

func init() {
	flag.Parse()
	var err error
	Ctx, err = oto.NewContext(*SampleRate, *ChannelNum, *BitDepthInBytes, 4096)
	if err != nil {
		log.Fatalf("Unable to oto.NewContext: %v\n", err)
	}
}

var (
	SampleRate      = flag.Int("samplerate", 48000, "sample rate")
	ChannelNum      = flag.Int("channelnum", 2, "number of channel")
	BitDepthInBytes = flag.Int("bitdepthinbytes", 2, "bit depth in bytes")

	Ctx *oto.Context
)

type Sound struct {
	Waveform []float64
}

func (s *Sound) Clone() (out Sound) {
	out.Waveform = make([]float64, len(s.Waveform))
	copy(out.Waveform, s.Waveform)
	return
}

func (s *Sound) ScaleAmplitude(r float64) *Sound {
	for i, v := range s.Waveform {
		s.Waveform[i] = v * r
	}
	return s
}

func (s *Sound) ToByteStream() (out []byte) {
	stride := (*BitDepthInBytes) * (*ChannelNum)
	sampleCount := len(s.Waveform)
	out = make([]byte, stride*sampleCount)

	switch *BitDepthInBytes {
	case 1:
		for i, sin := range s.Waveform {
			for ch := 0; ch < *ChannelNum; ch++ {
				out[i*stride+ch] = byte(int(sin*127) + 128)
			}
		}
	case 2:
		for i, sin := range s.Waveform {
			for ch := 0; ch < *ChannelNum; ch++ {
				v := int16(sin * 32767)
				out[i*stride+2*ch] = byte(v)
				out[i*stride+2*ch+1] = byte(v >> 8)
			}
		}
	}
	return
}
