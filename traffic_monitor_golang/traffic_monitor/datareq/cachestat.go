package datareq

import (
	"net/http"
	"net/url"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/cache"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/peer"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/threadsafe"
	todata "github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/trafficopsdata"
)

func srvCacheStats(params url.Values, errorCount threadsafe.Uint, path string, toData todata.TODataThreadsafe, statResultHistory threadsafe.ResultStatHistory, statInfoHistory threadsafe.ResultInfoHistory, monitorConfig threadsafe.TrafficMonitorConfigMap, combinedStates peer.CRStatesThreadsafe, statMaxKbpses threadsafe.CacheKbpses) ([]byte, int) {
	filter, err := NewCacheStatFilter(path, params, toData.Get().ServerTypes)
	if err != nil {
		HandleErr(errorCount, path, err)
		return []byte(err.Error()), http.StatusBadRequest
	}
	bytes, err := cache.StatsMarshall(statResultHistory.Get(), statInfoHistory.Get(), combinedStates.Get(), monitorConfig.Get(), statMaxKbpses.Get(), filter, params)
	return WrapErrCode(errorCount, path, bytes, err)
}
