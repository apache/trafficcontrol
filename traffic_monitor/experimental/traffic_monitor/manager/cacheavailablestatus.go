package manager

import (
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/enum"
	"sync"
)

// CacheAvailableStatusReported is the status string returned by caches set to "reported" in Traffic Ops.
// TODO put somewhere more generic
const CacheAvailableStatusReported = "REPORTED"

// CacheAvailableStatus is the available status of the given cache. It includes a boolean available/unavailable flag, and a descriptive string.
type CacheAvailableStatus struct {
	Available bool
	Status    string
}

// CacheAvailableStatuses is the available status of each cache.
type CacheAvailableStatuses map[enum.CacheName]CacheAvailableStatus

// CacheAvailableStatusThreadsafe wraps a map of cache available statuses to be safe for multiple reader goroutines and one writer.
type CacheAvailableStatusThreadsafe struct {
	caches *CacheAvailableStatuses
	m      *sync.RWMutex
}

// Copy copies this CacheAvailableStatuses. It does not modify, and thus is safe for multiple reader goroutines.
func (a CacheAvailableStatuses) Copy() CacheAvailableStatuses {
	b := CacheAvailableStatuses(map[enum.CacheName]CacheAvailableStatus{})
	for k, v := range a {
		b[k] = v
	}
	return b
}

// NewCacheAvailableStatusThreadsafe creates and returns a new CacheAvailableStatusThreadsafe, initializing internal pointer values.
func NewCacheAvailableStatusThreadsafe() CacheAvailableStatusThreadsafe {
	c := CacheAvailableStatuses(map[enum.CacheName]CacheAvailableStatus{})
	return CacheAvailableStatusThreadsafe{m: &sync.RWMutex{}, caches: &c}
}

// Get returns the internal map of cache statuses. The returned map MUST NOT be modified. If modification is necessary, copy.
func (o *CacheAvailableStatusThreadsafe) Get() CacheAvailableStatuses {
	o.m.RLock()
	defer o.m.RUnlock()
	return *o.caches
}

// Set sets the internal map of cache availability. This MUST NOT be called by multiple goroutines.
func (o *CacheAvailableStatusThreadsafe) Set(v CacheAvailableStatuses) {
	o.m.Lock()
	*o.caches = v
	o.m.Unlock()
}
