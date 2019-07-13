package rave

type Inst interface {
	ID() string
	Done() chan bool
	SyncID() string
	PlayLoop(ctx *Context)
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

func (d *Doc) Wave() *Wave {
	w := &Wave{
		Name: "",
		Sync: "global",
		done: make(chan bool),
	}
	d.Insts = append(d.Insts, w)
	return w
}

func (d *Doc) hasInst(name string) bool {
	for _, k := range d.Insts {
		if k.ID() == name {
			return true
		}
	}
	return false
}
