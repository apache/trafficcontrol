package thread

import (
	"sync"
)

type Throttler interface {
	Throttle(f func())
}

type throttler struct {
	c chan struct{}
}

func NewThrottler(max uint64) Throttler {
	if max < 1 {
		return NewNoThrottler()
	}
	return &throttler{c: make(chan struct{}, max)}
}

func (l *throttler) Throttle(f func()) {
	l.c <- struct{}{}
	f()
	<-l.c
}

func NewThrottlers(max uint64) Throttlers {
	return &throttlers{max: max, throttlers: map[string]Throttler{}}
}

type Throttlers interface {
	Throttle(k string, f func())
}

// Throttlers provides a threadsafe way to map throttlers to string keys, and delete the throttler when it is no longer being used (i.e. there's nothing requesting that key).
type throttlers struct {
	max           uint64
	throttlers    map[string]Throttler
	mutex         sync.Mutex
	checkoutCount uint64
}

func (t *throttlers) checkoutThrottler(k string) Throttler {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if throttler, ok := t.throttlers[k]; !ok {
		throttler = NewThrottler(t.max)
		t.throttlers[k] = throttler
		return throttler
	} else {
		return throttler
	}
}

func (t *throttlers) checkinThrottler(k string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.checkoutCount -= 1
	if t.checkoutCount == 0 {
		delete(t.throttlers, k)
	}
}

func (t *throttlers) Throttle(k string, f func()) {
	throttler := t.checkoutThrottler(k)
	throttler.Throttle(f)
	t.checkinThrottler(k)
}

type nothrottler struct{}

// NewNoThrottler creates and returns a Throttler which doesn't actually throttle, but Throttle(f) immediately calls f.
func NewNoThrottler() Throttler {
	return &nothrottler{}
}

func (l *nothrottler) Throttle(f func()) {
	f()
}
