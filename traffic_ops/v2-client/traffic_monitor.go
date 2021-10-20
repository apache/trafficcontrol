package client

import (
	"encoding/json"
	"fmt"

	"github.com/apache/trafficcontrol/v6/lib/go-tc"
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

// API_CDN_MONITORING_CONFIG is the API path on which Traffic Ops serves the CDN monitoring
// configuration information. It is meant to be used with fmt.Sprintf to insert the necessary
// path parameters (namely the Name of the CDN of interest).
// See Also: https://traffic-control-cdn.readthedocs.io/en/latest/api/v2/cdns_name_configs_monitoring.html
const API_CDN_MONITORING_CONFIG = apiBase + "/cdns/%s/configs/monitoring"

// GetTrafficMonitorConfigMap is functionally identical to GetTrafficMonitorConfig, except that it
// coerces the value returned by the API to the tc.LegacyTrafficMonitorConfigMap structure.
func (to *Session) GetTrafficMonitorConfigMap(cdn string) (*tc.LegacyTrafficMonitorConfigMap, ReqInf, error) {
	tmConfig, reqInf, err := to.GetTrafficMonitorConfig(cdn)
	if err != nil {
		return nil, reqInf, err
	}
	tmConfigMap, err := tc.LegacyTrafficMonitorTransformToMap(tmConfig)
	if err != nil {
		return tmConfigMap, reqInf, err
	}
	return tmConfigMap, reqInf, nil
}

// GetTrafficMonitorConfig returns the monitoring configuration for the CDN named by 'cdn'.
func (to *Session) GetTrafficMonitorConfig(cdn string) (*tc.LegacyTrafficMonitorConfig, ReqInf, error) {
	url := fmt.Sprintf(API_CDN_MONITORING_CONFIG, cdn)
	resp, remoteAddr, err := to.request("GET", url, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.LegacyTMConfigResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return &data.Response, reqInf, nil
}
