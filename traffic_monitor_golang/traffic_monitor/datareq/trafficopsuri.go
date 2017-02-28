package datareq

import (
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/threadsafe"
)

func srvAPITrafficOpsURI(opsConfig threadsafe.OpsConfig) []byte {
	return []byte(opsConfig.Get().Url)
}
