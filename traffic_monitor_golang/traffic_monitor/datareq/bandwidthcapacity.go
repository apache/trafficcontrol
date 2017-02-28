package datareq

import (
	"fmt"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/threadsafe"
)

func srvAPIBandwidthCapacityKbps(statMaxKbpses threadsafe.CacheKbpses) []byte {
	maxKbpses := statMaxKbpses.Get()
	cap := int64(0)
	for _, kbps := range maxKbpses {
		cap += kbps
	}
	return []byte(fmt.Sprintf("%d", cap))
}
