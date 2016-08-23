package manager

import (
	ds "github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/deliveryservice"
	"sync"
)

type DSStatsThreadsafe struct {
	dsStats *ds.Stats
	m       *sync.Mutex
}

func NewDSStatsThreadsafe() DSStatsThreadsafe {
	s := ds.NewStats()
	return DSStatsThreadsafe{m: &sync.Mutex{}, dsStats: &s}
}

func (o *DSStatsThreadsafe) Get() ds.Stats {
	o.m.Lock()
	defer func() {
		o.m.Unlock()
	}()
	return o.dsStats.Copy()
}

func (o *DSStatsThreadsafe) Set(newDsStats ds.Stats) {
	o.m.Lock()
	*o.dsStats = newDsStats
	o.m.Unlock()
}
