package threadsafe

import (
	"sync"
	"time"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/enum"
)

// DurationMap wraps a map[enum.CacheName]time.Duration in an object safe for a single writer and multiple readers
type DurationMap struct {
	durationMap *map[enum.CacheName]time.Duration
	m           *sync.RWMutex
}

// Copy copies this duration map.
func CopyDurationMap(a map[enum.CacheName]time.Duration) map[enum.CacheName]time.Duration {
	b := map[enum.CacheName]time.Duration{}
	for k, v := range a {
		b[k] = v
	}
	return b
}

// NewDurationMap returns a new DurationMap safe for multiple readers and a single writer goroutine.
func NewDurationMap() DurationMap {
	m := map[enum.CacheName]time.Duration{}
	return DurationMap{m: &sync.RWMutex{}, durationMap: &m}
}

// Get returns the duration map. Callers MUST NOT mutate. If mutation is necessary, call DurationMap.Copy().
func (o *DurationMap) Get() map[enum.CacheName]time.Duration {
	o.m.RLock()
	defer o.m.RUnlock()
	return *o.durationMap
}

// Set sets the internal duration map. This MUST NOT be called by multiple goroutines.
func (o *DurationMap) Set(d map[enum.CacheName]time.Duration) {
	o.m.Lock()
	*o.durationMap = d
	o.m.Unlock()
}
