package thread

/*
   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

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
	throttler, ok := t.throttlers[k]
	if !ok {
		throttler = NewThrottler(t.max)
		t.throttlers[k] = throttler
	}
	return throttler
}

func (t *throttlers) checkinThrottler(k string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.checkoutCount--
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
