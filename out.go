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

	snd := kit.Samples["snares_1"]
	mes := music.Measure{
		TimeSignature: music.TimeSignature{4, 4},
		Tempo:         80,
		Beats: []music.Beat{
			{
				Offset: 0,
				Player: &snd,
			},
			{
				Offset: 1,
				Player: &snd,
			},
			{
				Offset: 2,
				Player: &snd,
			},
			{
				Offset: 3,
				Player: &snd,
			},
		},
	}
	mes.Play()
}
