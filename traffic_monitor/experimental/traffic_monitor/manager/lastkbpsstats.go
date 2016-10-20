package manager

import (
	ds "github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/deliveryservice"
	"sync"
)

type LastStatsThreadsafe struct {
	stats *ds.LastStats
	m     *sync.RWMutex
}

func NewLastStatsThreadsafe() LastStatsThreadsafe {
	s := ds.NewLastStats()
	return LastStatsThreadsafe{m: &sync.RWMutex{}, stats: &s}
}

// Get returns the last KBPS stats object. Callers MUST NOT modify the object. It is not threadsafe for writing. If the object must be modified, callers must call LastStats.Copy() and modify the copy.
func (o *LastStatsThreadsafe) Get() ds.LastStats {
	o.m.RLock()
	defer o.m.RUnlock()
	return *o.stats
}

func (o *LastStatsThreadsafe) Set(s ds.LastStats) {
	o.m.Lock()
	*o.stats = s
	o.m.Unlock()
}
