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

package client

import (
	"fmt"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

func (to *Session) TopologiesQueueUpdate(topologyName tc.TopologyName, req tc.TopologiesQueueUpdateRequest) (tc.TopologiesQueueUpdateResponse, toclientlib.ReqInf, error) {
	path := fmt.Sprintf("%s/%s/queue_update", APITopologies, topologyName)
	var resp tc.TopologiesQueueUpdateResponse
	reqInf, err := to.post(path, req, nil, &resp)
	return resp, reqInf, err
}
