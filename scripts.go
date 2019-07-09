package rave

type Inst interface {
	ID() string
	Done() <-chan bool
}

type Doc struct {
	TimeTop float64
	TimeBot float64
	Tempo   float64
	Insts   []Inst
}

func (d *Doc) Kit() *Kit {
	i := &Kit{
		Name:   "",
		Volume: 0.1,
		Sync:   "global",
		done:   make(chan bool),
	}
	d.Insts = append(d.Insts, i)
	return i
}

func (d *Doc) hasInst(name string) bool {
	for _, k := range d.Insts {
		if k.ID() == name {
			return true
		}
	}
	return false
}

type Kit struct {
	Name    string
	Samples []string
	Sync    string
	Volume  float64
	loop    *Loop
	done    chan bool
}

func (k *Kit) ID() string {
	return k.Name
}

func (k *Kit) Done() <-chan bool {
	return k.done
}

func (k *Kit) Sample(s string) {
	k.Samples = append(k.Samples, s)
}

func (k *Kit) Loop() *Loop {
	if k.loop == nil {
		k.loop = &Loop{}
	}

	return k.loop
}

type Loop struct {
	Measures []*Measure
}

func (l *Loop) Measure() *Measure {
	m := &Measure{}
	l.Measures = append(l.Measures, m)
	return m
}

type Measure struct {
	Weights []float64
	Samples []int
	Pulses  []float64
}

func (m *Measure) Pulse(t float64) {
	m.Pulses = append(m.Pulses, t)
}

func (m *Measure) Sample(t int) {
	m.Samples = append(m.Samples, t)
}

func (m *Measure) Weight(f float64) {
	m.Weights = append(m.Weights, f)
}
