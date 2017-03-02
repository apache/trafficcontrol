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
	"net/url"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/peer"
)

func srvTRState(params url.Values, localStates peer.CRStatesThreadsafe, combinedStates peer.CRStatesThreadsafe) ([]byte, error) {
	if _, raw := params["raw"]; raw {
		return srvTRStateSelf(localStates)
	}
	return srvTRStateDerived(combinedStates)
}

func srvTRStateDerived(combinedStates peer.CRStatesThreadsafe) ([]byte, error) {
	return peer.CrstatesMarshall(combinedStates.Get())
}

func srvTRStateSelf(localStates peer.CRStatesThreadsafe) ([]byte, error) {
	return peer.CrstatesMarshall(localStates.Get())
}
