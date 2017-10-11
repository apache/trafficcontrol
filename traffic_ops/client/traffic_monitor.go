package client

import (
	"encoding/json"
	"fmt"

	tc "github.com/apache/incubator-trafficcontrol/lib/go-tc"
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

// TrafficMonitorConfigMap ...
func (to *Session) TrafficMonitorConfigMap(cdn string) (*tc.TrafficMonitorConfigMap, error) {
	tmConfig, err := to.TrafficMonitorConfig(cdn)
	if err != nil {
		return nil, err
	}
	tmConfigMap, err := tc.TrafficMonitorTransformToMap(tmConfig)
	if err != nil {
		return nil, err
	}
	return tmConfigMap, nil
}

// TrafficMonitorConfig ...
func (to *Session) TrafficMonitorConfig(cdn string) (*tc.TrafficMonitorConfig, error) {
	url := fmt.Sprintf("/api/1.2/cdns/%s/configs/monitoring.json", cdn)
	resp, err := to.request("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data tc.TMConfigResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data.Response, nil
}
