package datareq

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/enum"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/peer"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/srvhttp"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/threadsafe"
	todata "github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/trafficopsdata"
)

// APIPeerStates contains the data to be returned for an API call to get the peer states of a Traffic Monitor. This contains common API data returned by most endpoints, and a map of peers, to caches' states.
type APIPeerStates struct {
	srvhttp.CommonAPIData
	Peers map[enum.TrafficMonitorName]map[enum.CacheName][]CacheState `json:"peers"`
}

// CacheState represents the available state of a cache.
type CacheState struct {
	Value bool `json:"value"`
}

func srvPeerStates(params url.Values, errorCount threadsafe.Uint, path string, toData todata.TODataThreadsafe, peerStates peer.CRStatesPeersThreadsafe) ([]byte, int) {
	filter, err := NewPeerStateFilter(path, params, toData.Get().ServerTypes)
	if err != nil {
		HandleErr(errorCount, path, err)
		return []byte(err.Error()), http.StatusBadRequest
	}
	bytes, err := json.Marshal(createAPIPeerStates(peerStates.GetCrstates(), filter, params))
	return WrapErrCode(errorCount, path, bytes, err)
}

func createAPIPeerStates(peerStates map[enum.TrafficMonitorName]peer.Crstates, filter *PeerStateFilter, params url.Values) APIPeerStates {
	apiPeerStates := APIPeerStates{
		CommonAPIData: srvhttp.GetCommonAPIData(params, time.Now()),
		Peers:         map[enum.TrafficMonitorName]map[enum.CacheName][]CacheState{},
	}

	for peer, state := range peerStates {
		if !filter.UsePeer(peer) {
			continue
		}
		if _, ok := apiPeerStates.Peers[peer]; !ok {
			apiPeerStates.Peers[peer] = map[enum.CacheName][]CacheState{}
		}
		peerState := apiPeerStates.Peers[peer]
		for cache, available := range state.Caches {
			if !filter.UseCache(cache) {
				continue
			}
			peerState[cache] = []CacheState{CacheState{Value: available.IsAvailable}}
		}
		apiPeerStates.Peers[peer] = peerState
	}
	return apiPeerStates
}
