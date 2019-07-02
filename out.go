package main

import (
	"io/ioutil"
	"log"

	"github.com/dianelooney/rave/music"

	"github.com/dianelooney/rave/kits"
	"gopkg.in/yaml.v2"
)

// . "github.com/dianelooney/rave/notes"
// . "github.com/dianelooney/rave/pitches"

func main() {

	data, err := ioutil.ReadFile("kits/manifest.yml")
	if err != nil {
		log.Fatalf("Failed to load manifest: %v\n", err)
	}

	m := kits.Manifest{}
	err = yaml.Unmarshal(data, &m)
	if err != nil {
		log.Fatalf("Failed to parse manifest: %v\n", err)
	}
	kit := m.Load()
	bass := kit.Sample("kicks_1")
	bass.ScaleAmplitude(0.1)
	snare := kit.Sample("snares_1")
	snare.ScaleAmplitude(0.1)
	mes := music.Measure{
		TimeSignature: music.TimeSignature{9, 8},
		Tempo:         240,
		Beats: []music.Beat{
			{
				Offset: 0,
				Player: &bass,
			},
			{
				Offset: 1,
				Player: &snare,
			},
			{
				Offset: 2,
				Player: &bass,
			},
			{
				Offset: 3,
				Player: &snare,
			},
			{
				Offset: 4,
				Player: &bass,
			},
			{
				Offset: 5,
				Player: &snare,
			},
			{
				Offset: 6,
				Player: &bass,
			},
			{
				Offset: 7,
				Player: &snare,
			},
			{
				Offset: 8,
				Player: &snare,
			},
		},
	}
	for {
		mes.Play()
	}
}
