package datareq

import (
	"strconv"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/enum"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/peer"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/threadsafe"
	to "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
)

func srvAPICacheDownCount(localStates peer.CRStatesThreadsafe, monitorConfig threadsafe.TrafficMonitorConfigMap) []byte {
	return []byte(strconv.Itoa(cacheDownCount(localStates.Get().Caches, monitorConfig.Get().TrafficServer)))
}

// cacheOfflineCount returns the total reported caches marked down, excluding status offline and admin_down.
func cacheDownCount(caches map[enum.CacheName]peer.IsAvailable, toServers map[string]to.TrafficServer) int {
	count := 0
	for cache, available := range caches {
		if !available.IsAvailable && enum.CacheStatusFromString(toServers[string(cache)].Status) == enum.CacheStatusReported {
			count++
		}
	}
	return count
}
