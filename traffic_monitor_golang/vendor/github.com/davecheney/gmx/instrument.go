package gmx

import (
	"sync/atomic"
)

type Counter struct {
	value uint64
}

func (c *Counter) Inc() {
	atomic.AddUint64(&c.value, 1)
}

func NewCounter(name string) *Counter {
	c := new(Counter)
	Publish(name, func() interface{} {
		return c.value
	})
	return c
}

type Gauge struct {
	value int64
}

func (g *Gauge) Inc() {
	atomic.AddInt64(&g.value, 1)
}

func (g *Gauge) Dec() {
	atomic.AddInt64(&g.value, -1)
}

func NewGauge(name string) *Gauge {
	g := new(Gauge)
	Publish(name, func() interface{} {
		return g.value
	})
	return g
}
