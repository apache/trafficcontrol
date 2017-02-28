package datareq

import (
	"strconv"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/peer"
)

// TODO determine if this should use peerStates
func srvAPICacheCount(localStates peer.CRStatesThreadsafe) []byte {
	return []byte(strconv.Itoa(len(localStates.Get().Caches)))
}
