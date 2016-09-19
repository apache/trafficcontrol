package manager

import (
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/enum"
	"sync"
)

// TODO put somewhere more generic
const CacheAvailableStatusReported = "REPORTED"

type CacheAvailableStatus struct {
	Available bool
	Status    string
}

type CacheAvailableStatusThreadsafe struct {
	caches map[enum.CacheName]CacheAvailableStatus // TODO change string -> CacheName
	m      *sync.RWMutex
}

func copyCacheAvailableStatus(a map[enum.CacheName]CacheAvailableStatus) map[enum.CacheName]CacheAvailableStatus {
	b := map[enum.CacheName]CacheAvailableStatus{}
	for k, v := range a {
		b[k] = v
	}
	return b
}

func NewCacheAvailableStatusThreadsafe() CacheAvailableStatusThreadsafe {
	return CacheAvailableStatusThreadsafe{m: &sync.RWMutex{}, caches: map[enum.CacheName]CacheAvailableStatus{}}
}

func (o *CacheAvailableStatusThreadsafe) Get() map[enum.CacheName]CacheAvailableStatus {
	o.m.RLock()
	defer o.m.RUnlock()
	return copyCacheAvailableStatus(o.caches)
}

func (o *CacheAvailableStatusThreadsafe) Set(cache enum.CacheName, status CacheAvailableStatus) {
	o.m.Lock()
	o.caches[cache] = status
	o.m.Unlock()
}
