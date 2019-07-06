package common

func NewContext(parent *Context) Context {
	return Context{
		parent: parent,
		data:   make(map[FieldID]interface{}),
	}
}

type Context struct {
	parent *Context
	data   map[FieldID]interface{}
}

type FieldID uint8

const (
	Tempo FieldID = iota
)

func (c *Context) Get(f FieldID) (value interface{}, ok bool) {
	if c == nil {
		return nil, false
	}

	v, ok := c.data[f]
	if ok {
		return v, true
	}

	return c.parent.Get(f)
}

func (c *Context) Set(f FieldID, v interface{}) {
	if c == nil {
		return
	}

	c.data[f] = v
}

var Background = NewContext(nil)
