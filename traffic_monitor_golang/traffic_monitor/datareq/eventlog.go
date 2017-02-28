package datareq

import (
	"encoding/json"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/health"
)

// JSONEvents represents the structure we wish to serialize to JSON, for Events.
type JSONEvents struct {
	Events []health.Event `json:"events"`
}

func srvEventLog(events health.ThreadsafeEvents) ([]byte, error) {
	return json.Marshal(JSONEvents{Events: events.Get()})
}
