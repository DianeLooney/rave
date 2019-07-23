package music

import (
	"time"

	. "github.com/dianelooney/rave/common"
)

type BPM float64

func (b BPM) SecondsPerBeat() float64 {
	return 60 / float64(b)
}

type Measure struct {
	Size float64
	TimeSignature
	Tempo BPM
	Beats []Beat
}

type TimeSignature [2]float64

func (t TimeSignature) Top() float64 {
	return t[0]
}

func (t TimeSignature) Bottom() float64 {
	return t[1]
}

type Beat struct {
	Offset float64
	Player
}

func (m Measure) Play() {
	m.PlayAt(time.Now())
}

func (m Measure) PlayAt(t time.Time) {
	time.Sleep(time.Until(t))

	for _, b := range m.Beats {
		go func(b Beat) {
			waitForBeat(t, m.Tempo, b.Offset)
			b.Play()
		}(b)
	}
	waitForBeat(t, m.Tempo, m.TimeSignature.Top())
}

func waitForBeat(t time.Time, b BPM, count float64) {
	seconds := b.SecondsPerBeat() * count
	elapsed := time.Duration(float64(time.Second) * seconds)
	<-time.NewTimer(elapsed).C
}
