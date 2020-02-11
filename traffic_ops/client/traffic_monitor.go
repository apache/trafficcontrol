package client

import (
	"encoding/json"
	"fmt"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

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

func (to *Session) GetTrafficMonitorConfigMap(cdn string) (*tc.TrafficMonitorConfigMap, ReqInf, error) {
	tmConfig, reqInf, err := to.GetTrafficMonitorConfig(cdn)
	if err != nil {
		return nil, reqInf, err
	}
	tmConfigMap, err := tc.TrafficMonitorTransformToMap(tmConfig)
	if err != nil {
		return nil, reqInf, err
	}
	return tmConfigMap, reqInf, nil
}

func (to *Session) GetTrafficMonitorConfig(cdn string) (*tc.TrafficMonitorConfig, ReqInf, error) {
	url := fmt.Sprintf("/api/1.3/cdns/%s/configs/monitoring.json", cdn)
	resp, remoteAddr, err := to.request("GET", url, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.TMConfigResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return &data.Response, reqInf, nil
}
