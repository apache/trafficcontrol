package datareq

import (
	"encoding/json"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/threadsafe"
)

func srvConfigDoc(opsConfig threadsafe.OpsConfig) ([]byte, error) {
	opsConfigCopy := opsConfig.Get()
	// if the password is blank, leave it blank, so callers can see it's missing.
	if opsConfigCopy.Password != "" {
		opsConfigCopy.Password = "*****"
	}
	return json.Marshal(opsConfigCopy)
}
