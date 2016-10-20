package manager

import (
	ds "github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/deliveryservice"
	dsdata "github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/deliveryservicedata"
	"sync"
)

type DSStatsThreadsafe struct {
	dsStats *ds.Stats
	m       *sync.RWMutex
}

type DSStatsReader interface {
	Get() dsdata.StatsReadonly
}

func NewDSStatsThreadsafe() DSStatsThreadsafe {
	s := ds.NewStats()
	return DSStatsThreadsafe{m: &sync.RWMutex{}, dsStats: &s}
}

func (o *DSStatsThreadsafe) Get() dsdata.StatsReadonly {
	o.m.RLock()
	defer o.m.RUnlock()
	return *o.dsStats
}

func (o *DSStatsThreadsafe) Set(newDsStats ds.Stats) {
	o.m.Lock()
	*o.dsStats = newDsStats
	o.m.Unlock()
}
