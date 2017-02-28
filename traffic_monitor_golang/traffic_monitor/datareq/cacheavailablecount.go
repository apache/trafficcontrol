package datareq

import (
	"strconv"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/enum"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/peer"
)

func srvAPICacheAvailableCount(localStates peer.CRStatesThreadsafe) []byte {
	return []byte(strconv.Itoa(cacheAvailableCount(localStates.Get().Caches)))
}

// cacheOfflineCount returns the total caches not available, including marked unavailable, status offline, and status admin_down
func cacheOfflineCount(caches map[enum.CacheName]peer.IsAvailable) int {
	count := 0
	for _, available := range caches {
		if !available.IsAvailable {
			count++
		}
	}
	return count
}

// cacheAvailableCount returns the total caches available, including marked available and status online
func cacheAvailableCount(caches map[enum.CacheName]peer.IsAvailable) int {
	return len(caches) - cacheOfflineCount(caches)
}
