package crstates

import (
	"sync"
)

// Ths provides threadsafe access to a ThsT pointer. Note the object itself is not safe for multiple access, and must not be mutated, either by the original owner after calling Set, or by future users who call Get. If you need to mutate, perform a deep copy.
type Ths struct {
	v *ThsT
	m *sync.RWMutex
}

func NewThs() Ths {
	v := ThsT(nil)
	return Ths{m: &sync.RWMutex{}, v: &v}
}

func (t Ths) Set(v ThsT) {
	t.m.Lock()
	defer t.m.Unlock()
	*t.v = v
}

func (t Ths) Get() ThsT {
	t.m.RLock()
	defer t.m.RUnlock()
	return *t.v
}
