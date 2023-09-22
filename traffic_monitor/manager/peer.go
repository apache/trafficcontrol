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
)

// StartPeerManager listens for peer results, and when it gets one, it adds it to the peerStates list, and optimistically combines the good results into combinedStates
func StartPeerManager(
	peerChan <-chan peer.Result,
	peerStates peer.CRStatesPeersThreadsafe,
	events health.ThreadsafeEvents,
	combineState func(),
) {
	go func() {
		for peerResult := range peerChan {
			comparePeerState(events, peerResult, peerStates)
			peerStates.Set(peerResult)
			combineState()
			peerResult.PollFinished <- peerResult.PollID
		}
	}()
}

func comparePeerState(events health.ThreadsafeEvents, result peer.Result, peerStates peer.CRStatesPeersThreadsafe) {
	if result.Available != peerStates.GetPeerAvailability(result.ID) {
		description := util.JoinErrsStr(result.Errors)

		if description == "" && result.Available {
			description = "Peer is reachable"
		} else if description == "" && !result.Available {
			description = "Peer is unreachable"
		}

		events.Add(
			health.Event{
				Time:        health.Time(result.Time),
				Description: description,
				Name:        result.ID.String(),
				Hostname:    result.ID.String(),
				Type:        "PEER",
				Available:   result.Available})
	}
}
