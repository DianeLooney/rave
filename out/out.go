package out

import (
	"flag"
	"github.com/hajimehoshi/oto"
	"log"
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
	SampleRate      = flag.Int("samplerate", 44100, "sample rate")
	ChannelNum      = flag.Int("channelnum", 2, "number of channel")
	BitDepthInBytes = flag.Int("bitdepthinbytes", 2, "bit depth in bytes")

	Ctx *oto.Context
)
