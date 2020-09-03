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
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

func (to *Session) TopologiesQueueUpdate(topologyName tc.TopologyName, req tc.TopologiesQueueUpdateRequest) (tc.TopologiesQueueUpdateResponse, ReqInf, error) {
	path := fmt.Sprintf(ApiTopologies + "/%s/queue_update", topologyName)
	var reqInf ReqInf
	var resp tc.TopologiesQueueUpdateResponse

	reqBody, err := json.Marshal(req)
	if err != nil {
		return resp, reqInf, err
	}

	httpResp, remoteAddr, err := to.request(http.MethodPost, path, reqBody, nil)
	if httpResp != nil {
		reqInf.StatusCode = httpResp.StatusCode
		reqInf.RemoteAddr = remoteAddr
	}
	if err != nil {
		return resp, reqInf, err
	}
	defer log.Close(httpResp.Body, "unable to close response")

	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	return resp, reqInf, err
}
