package datareq

import (
	"fmt"

	ds "github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/deliveryservice"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/threadsafe"
	todata "github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/trafficopsdata"
)

func srvAPIBandwidthKbps(toData todata.TODataThreadsafe, lastStats threadsafe.LastStats) []byte {
	kbpsStats := lastStats.Get()
	sum := float64(0.0)
	for _, data := range kbpsStats.Caches {
		sum += data.Bytes.PerSec / ds.BytesPerKilobit
	}
	return []byte(fmt.Sprintf("%f", sum))
}
