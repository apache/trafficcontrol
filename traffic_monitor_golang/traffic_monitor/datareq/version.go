package datareq

import (
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/config"
)

func srvAPIVersion(staticAppData config.StaticAppData) []byte {
	s := "traffic_monitor-" + staticAppData.Version + "."
	if len(staticAppData.GitRevision) > 6 {
		s += staticAppData.GitRevision[:6]
	} else {
		s += staticAppData.GitRevision
	}
	return []byte(s)
}
