package datareq

import (
	"encoding/json"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/threadsafe"
)

func srvMonitorConfig(mcThs threadsafe.TrafficMonitorConfigMap) ([]byte, error) {
	return json.Marshal(mcThs.Get())
}
