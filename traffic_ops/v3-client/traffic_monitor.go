package client

import (
	"fmt"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
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
// See Also: https://traffic-control-cdn.readthedocs.io/en/latest/api/v3/cdns_name_configs_monitoring.html
//
// DEPRECATED: will be removed in the next major version. Be aware this may not be the URI being requested, for clients created with Login and ClientOps.ForceLatestAPI false.
const API_CDN_MONITORING_CONFIG = apiBase + "/cdns/%s/configs/monitoring"

const APICDNMonitoringConfig = "/cdns/%s/configs/monitoring"

// GetTrafficMonitorConfigMap is functionally identical to GetTrafficMonitorConfig, except that it
// coerces the value returned by the API to the tc.TrafficMonitorConfigMap structure.
func (to *Session) GetTrafficMonitorConfigMap(cdn string) (*tc.TrafficMonitorConfigMap, toclientlib.ReqInf, error) {
	tmConfig, reqInf, err := to.GetTrafficMonitorConfig(cdn)
	if err != nil {
		return nil, reqInf, err
	}
	tmConfigMap, err := tc.TrafficMonitorTransformToMap(tmConfig)
	if err != nil {
		return tmConfigMap, reqInf, err
	}
	return tmConfigMap, reqInf, nil
}

// GetTrafficMonitorConfig returns the monitoring configuration for the CDN named by 'cdn'.
func (to *Session) GetTrafficMonitorConfig(cdn string) (*tc.TrafficMonitorConfig, toclientlib.ReqInf, error) {
	route := fmt.Sprintf(APICDNMonitoringConfig, cdn)
	var data tc.TMConfigResponse
	reqInf, err := to.get(route, nil, &data)
	return &data.Response, reqInf, err
}
