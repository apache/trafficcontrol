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
	m      *sync.Mutex
}

func copyCacheAvailableStatus(a map[enum.CacheName]CacheAvailableStatus) map[enum.CacheName]CacheAvailableStatus {
	b := map[enum.CacheName]CacheAvailableStatus{}
	for k, v := range a {
		b[k] = v
	}
	return b
}

func NewCacheAvailableStatusThreadsafe() CacheAvailableStatusThreadsafe {
	return CacheAvailableStatusThreadsafe{m: &sync.Mutex{}, caches: map[enum.CacheName]CacheAvailableStatus{}}
}

func (o *CacheAvailableStatusThreadsafe) Get() map[enum.CacheName]CacheAvailableStatus {
	o.m.Lock()
	defer func() {
		o.m.Unlock()
	}()
	return copyCacheAvailableStatus(o.caches)
}

func (o *CacheAvailableStatusThreadsafe) Set(cache enum.CacheName, status CacheAvailableStatus) {
	o.m.Lock()
	o.caches[cache] = status
	o.m.Unlock()
}
