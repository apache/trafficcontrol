package datareq

import (
	"fmt"
	"time"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/threadsafe"
	towrap "github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/trafficopswrapper"
)

func srvTRConfig(opsConfig threadsafe.OpsConfig, toSession towrap.ITrafficOpsSession) ([]byte, time.Time, error) {
	cdnName := opsConfig.Get().CdnName
	if toSession == nil {
		return nil, time.Time{}, fmt.Errorf("Unable to connect to Traffic Ops")
	}
	if cdnName == "" {
		return nil, time.Time{}, fmt.Errorf("No CDN Configured")
	}
	return toSession.LastCRConfig(cdnName)
}
