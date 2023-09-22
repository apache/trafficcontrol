/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package datareq

import (
	"net/http"
	"net/url"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/peer"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/srvhttp"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/threadsafe"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/todata"
	jsoniter "github.com/json-iterator/go"
)

// APIPeerStates contains the data to be returned for an API call to get the peer states of a Traffic Monitor. This contains common API data returned by most endpoints, and a map of peers, to caches' states.
type APIPeerStates struct {
	tc.CommonAPIData
	Peers map[tc.TrafficMonitorName]map[tc.CacheName][]CacheState `json:"peers"`
}

// CacheState represents the available state of a cache.
type CacheState struct {
	Value         bool `json:"value"`
	Ipv4Available bool `json:"ipv4Available"`
	Ipv6Available bool `json:"ipv6Available"`
}

func srvPeerStates(params url.Values, errorCount threadsafe.Uint, path string, toData todata.TODataThreadsafe, peerStates peer.CRStatesPeersThreadsafe) ([]byte, int) {
	filter, err := NewPeerStateFilter(path, params, toData.Get().ServerTypes)
	if err != nil {
		HandleErr(errorCount, path, err)
		return []byte(err.Error()), http.StatusBadRequest
	}
	json := jsoniter.ConfigFastest
	bytes, err := json.Marshal(createAPIPeerStates(peerStates.GetCrstates(), peerStates.GetPeersOnline(), filter, params))
	return WrapErrCode(errorCount, path, bytes, err)
}

func createAPIPeerStates(peerStates map[tc.TrafficMonitorName]tc.CRStates, peersOnline map[tc.TrafficMonitorName]bool, filter *PeerStateFilter, params url.Values) APIPeerStates {
	apiPeerStates := APIPeerStates{
		CommonAPIData: srvhttp.GetCommonAPIData(params, time.Now()),
		Peers:         map[tc.TrafficMonitorName]map[tc.CacheName][]CacheState{},
	}

	for peer, state := range peerStates {
		if !peersOnline[peer] {
			continue
		}
		if !filter.UsePeer(peer) {
			continue
		}
		if _, ok := apiPeerStates.Peers[peer]; !ok {
			apiPeerStates.Peers[peer] = map[tc.CacheName][]CacheState{}
		}
		peerState := apiPeerStates.Peers[peer]
		for cache, available := range state.Caches {
			if !filter.UseCache(cache) {
				continue
			}
			peerState[cache] = []CacheState{CacheState{Value: available.IsAvailable, Ipv4Available: available.Ipv4Available, Ipv6Available: available.Ipv6Available}}
		}
		apiPeerStates.Peers[peer] = peerState
	}
	return apiPeerStates
}
