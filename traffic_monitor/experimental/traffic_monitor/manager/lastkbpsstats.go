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

func (o *StatsLastKbpsThreadsafe) Get() ds.StatsLastKbps {
	o.m.RLock()
	defer o.m.RUnlock()
	return o.stats.Copy()
}

func (o *StatsLastKbpsThreadsafe) Set(s ds.StatsLastKbps) {
	o.m.Lock()
	*o.stats = s // TODO copy?
	o.m.Unlock()
}
