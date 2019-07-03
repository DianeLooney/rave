package main

import (
	"time"

	"github.com/dianelooney/rave/notes"
)

// . "github.com/dianelooney/rave/notes"
// . "github.com/dianelooney/rave/pitches"

func main() {

	/*
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

		var set interpreter.Set
		d, err := ioutil.ReadFile("interpreter/sample.yml")
		if err != nil {
			log.Fatalf("Unable to read file: %v\n", err)
		}
		err = yaml.Unmarshal(d, &set)
		if err != nil {
			log.Fatalf("Unable to unmarshal yaml: %v\n", err)
		}

		fmt.Printf("%+v\n", set)

		for _, inst := range set.Kits {
			go func(inst interpreter.Inst) {
				x := kit.Sample(inst.N)
				measures := make([]music.Measure, len(inst.M))
				for i, m := range inst.M {
					beats := make([]music.Beat, len(m))
					for j, y := range m {
						beats[j] = music.Beat{
							Offset: y - 1,
							Player: &x,
						}
					}
					measures[i] = music.Measure{
						TimeSignature: music.TimeSignature{set.Tstop, set.Tsbot},
						Tempo:         set.BPM(),
						Beats:         beats,
					}
				}
					for {
						for _, m := range measures {
							m.Play()
						}
					}
			}(inst)
		}
	*/

	note := notes.Note{
		Pattern:   "experiment",
		Pitch:     440,
		Duration:  2 * time.Second,
		Intensity: 0.1,
	}
	note.Play()
	//<-make(chan bool)
}
