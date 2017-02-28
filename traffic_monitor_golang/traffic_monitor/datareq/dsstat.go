package datareq

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/threadsafe"
	todata "github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/trafficopsdata"
)

func srvDSStats(params url.Values, errorCount threadsafe.Uint, path string, toData todata.TODataThreadsafe, dsStats threadsafe.DSStatsReader) ([]byte, int) {
	filter, err := NewDSStatFilter(path, params, toData.Get().DeliveryServiceTypes)
	if err != nil {
		HandleErr(errorCount, path, err)
		return []byte(err.Error()), http.StatusBadRequest
	}
	bytes, err := json.Marshal(dsStats.Get().JSON(filter, params))
	return WrapErrCode(errorCount, path, bytes, err)
}
