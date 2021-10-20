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
	"fmt"
	"net/http"
	"net/url"

	"github.com/apache/trafficcontrol/v6/lib/go-log"
	"github.com/apache/trafficcontrol/v6/lib/go-tc"
	"github.com/apache/trafficcontrol/v6/traffic_monitor/peer"
)

func srvTRState(params url.Values, localStates peer.CRStatesThreadsafe, combinedStates peer.CRStatesThreadsafe, peerStates peer.CRStatesPeersThreadsafe) ([]byte, int, error) {
	// local state requested (peer polling case)
	if _, raw := params["raw"]; raw {
		data, err := srvTRStateSelf(localStates)
		return data, http.StatusOK, err
	}

	// This covers the case where we have lost connectivity to all peers, but multiple peers exist. In this case, it is
	// more likely that the local machine has lost all connectivity than both peers losing connectivity or crashing. If
	// the peers really did crash, the health protocol is essentially broken, and serving a 503 will cause Traffic Router
	// to use the last good state fetched from a Traffic Monitor within the CDN. If the peers are simply unreachable from
	// this Traffic Monitor, serving 503s until connectivity is restored will cause Traffic Router to ignore this instance
	// until the health protocol can be relied upon once again.
	if peerStates.OptimisticQuorumEnabled() {
		optimisticQuorum, peersAvailable, peerCount, minimum := peerStates.HasOptimisticQuorum()
		log.Debugf("optimisticQuorum=%v, peerCount=%v, peersAvailable=%v, minimum=%v", optimisticQuorum, peerCount, peersAvailable, minimum)

		if !optimisticQuorum {
			return nil, http.StatusServiceUnavailable, fmt.Errorf("number of peers available (%d/%d) is less than the minimum number of %d required for optimistic peer quorum", peersAvailable, peerCount, minimum)
		}
	}

	data, err := srvTRStateDerived(combinedStates, peerStates)

	return data, http.StatusOK, err
}

func srvTRStateDerived(combinedStates peer.CRStatesThreadsafe, peerStates peer.CRStatesPeersThreadsafe) ([]byte, error) {
	return tc.CRStatesMarshall(combinedStates.Get())
}

func srvTRStateSelf(localStates peer.CRStatesThreadsafe) ([]byte, error) {
	return tc.CRStatesMarshall(localStates.Get())
}
