package manager

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

import (
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/health"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/peer"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/threadsafe"
)

// StartDistributedPeerManager listens for distributed peer results and updates the localStates
// directly (because combining peerStates is unnecessary since these results are already combined
// among the distributed TM group they came from).
func StartDistributedPeerManager(
	distributedPeerChan <-chan peer.Result,
	localStates peer.CRStatesThreadsafe,
	distributedPeerStates peer.CRStatesPeersThreadsafe,
	events health.ThreadsafeEvents,
	unpolledCaches threadsafe.UnpolledCaches,
) {
	go func() {
		for distributedPeerResult := range distributedPeerChan {
			compareDistributedPeerState(events, distributedPeerResult, distributedPeerStates)
			distributedPeerStates.Set(distributedPeerResult)
			for name, availability := range distributedPeerResult.PeerStates.Caches {
				localStates.SetCache(name, availability)
			}
			if len(distributedPeerResult.Errors) == 0 {
				unpolledCaches.SetRemotePolled(distributedPeerResult.PeerStates.Caches)
			}
			distributedPeerResult.PollFinished <- distributedPeerResult.PollID
		}
	}()
}

func compareDistributedPeerState(events health.ThreadsafeEvents, result peer.Result, distributedPeerStates peer.CRStatesPeersThreadsafe) {
	if result.Available != distributedPeerStates.GetPeerAvailability(result.ID) {
		description := util.JoinErrsStr(result.Errors)

		if description == "" && result.Available {
			description = "Distributed peer group is reachable"
		} else if description == "" && !result.Available {
			description = "Distributed peer group is unreachable"
		}

		events.Add(
			health.Event{
				Time:        health.Time(result.Time),
				Description: description,
				Name:        result.ID.String(),
				Hostname:    result.ID.String(),
				Type:        "DISTRIBUTED_PEER",
				Available:   result.Available})
	}
}
