package kits

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"

	. "github.com/dianelooney/rave/common"
	"github.com/youpy/go-wav"
)

type Kit struct {
	Samples map[string]Sound
}

func (k Kit) Sample(s string) Sound {
	snd, ok := k.Samples[s]
	if !ok {
		fmt.Printf("Sound '%s' missing from kit\n", s)
	}
	return snd
}

type Manifest map[string][]string

func (m Manifest) Load() (k Kit) {
	k.Samples = make(map[string]Sound)
	wg := sync.WaitGroup{}

	mtx := sync.Mutex{}
	for key, grp := range m {
		wg.Add(len(grp))
		for i, path := range grp {
			go func(i int, key, path string) {
				defer wg.Done()

				file, err := os.Open(path)
				if err != nil {
					mtx.Lock()
					fmt.Printf("Unable to open sound file '%v': %v\n", path, err)
					mtx.Unlock()
				}
				reader := wav.NewReader(file)

				defer file.Close()

				snd := Sound{
					Waveform: make([]float64, 0),
				}

				for {
					samples, err := reader.ReadSamples()
					if err == io.EOF {
						break
					}

					for _, sample := range samples {
						f1 := reader.FloatValue(sample, 0)
						snd.Waveform = append(snd.Waveform, f1)
						// TODO: Support stereo at some point
					}
				}

				mtx.Lock()
				defer mtx.Unlock()

				k.Samples[fmt.Sprintf("%v_%v", key, i)] = snd
			}(i, key, path)
		}
	}
	wg.Wait()

	return
}

func LoadManifest(path string) (m Manifest) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Unable to read manifest '%v': %v", path, err)
	}
	err = json.Unmarshal(data, &m)
	if err != nil {
		log.Fatalf("Unable to parse manifest: %v", err)
	}

	return
}
