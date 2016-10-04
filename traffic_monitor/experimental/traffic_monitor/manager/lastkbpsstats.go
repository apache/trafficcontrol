package manager

import (
	ds "github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/deliveryservice"
	"sync"
)

type StatsLastKbpsThreadsafe struct {
	stats *ds.StatsLastKbps
	m     *sync.RWMutex
}

func NewStatsLastKbpsThreadsafe() StatsLastKbpsThreadsafe {
	s := ds.NewStatsLastKbps()
	return StatsLastKbpsThreadsafe{m: &sync.RWMutex{}, stats: &s}
}

// Get returns the last KBPS stats object. Callers MUST NOT modify the object. It is not threadsafe for writing. If the object must be modified, callers must call StatsLastKbps.Copy() and modify the copy.
func (o *StatsLastKbpsThreadsafe) Get() ds.StatsLastKbps {
	o.m.RLock()
	defer o.m.RUnlock()
	return *o.stats
}

func (o *StatsLastKbpsThreadsafe) Set(s ds.StatsLastKbps) {
	o.m.Lock()
	*o.stats = s
	o.m.Unlock()
}
