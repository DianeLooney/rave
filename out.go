package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/dianelooney/rave/kits"
	"gopkg.in/yaml.v2"
)

// . "github.com/dianelooney/rave/notes"
// . "github.com/dianelooney/rave/pitches"

func main() {
	/*
		var err error

		if err != nil {
			log.Fatalf("Unable to initialize oto:\n%v", err)
		}
		Note{
			Pitch:     C4,
			Duration:  time.Second * 1,
			Pattern:   "experiment",
			Intensity: 0.2,
		}.Play()
		Note{
			Pitch:     E4,
			Duration:  time.Second * 1,
			Pattern:   "experiment",
			Intensity: 0.2,
		}.Play()
		Note{
			Pitch:     G4,
			Duration:  time.Second * 1,
			Pattern:   "experiment",
			Intensity: 0.2,
		}.Play()
		go Note{
			Pitch:     C4,
			Duration:  time.Second * 2,
			Pattern:   "experiment",
			Intensity: 0.2,
		}.Play()
		go Note{
			Pitch:     E4,
			Duration:  time.Second * 2,
			Pattern:   "experiment",
			Intensity: 0.2,
		}.Play()
		go Note{
			Pitch:     G4,
			Duration:  time.Second * 2,
			Pattern:   "experiment",
			Intensity: 0.2,
		}.Play()
		time.Sleep(2 * time.Second)
	*/

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
	fmt.Println(kit)
	snd := kit.Samples["snares_21"]
	snd.Play()
	snd.Play()
	snd.Play()
	snd.Play()
	snd.Play()
}
