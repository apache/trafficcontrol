package manager

import (
	ds "github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/deliveryservice"
	"sync"
)

type StatsLastKbpsThreadsafe struct {
	stats *ds.StatsLastKbps
	m     *sync.Mutex
}

func NewStatsLastKbpsThreadsafe() StatsLastKbpsThreadsafe {
	s := ds.NewStatsLastKbps()
	return StatsLastKbpsThreadsafe{m: &sync.Mutex{}, stats: &s}
}

func (o *StatsLastKbpsThreadsafe) Get() ds.StatsLastKbps {
	o.m.Lock()

	defer func() {
		o.m.Unlock()
	}()
	return o.stats.Copy()
}

func (o *StatsLastKbpsThreadsafe) Set(s ds.StatsLastKbps) {
	o.m.Lock()
	*o.stats = s // TODO copy?
	o.m.Unlock()
}
