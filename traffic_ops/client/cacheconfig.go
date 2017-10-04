package client

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
	"encoding/json"
)

type CacheConfigResponse struct {
	Response CacheConfig `json:"response"`
}

type CacheConfig struct {
	Domain           string                       `json:"domain"`
	CDN              string                       `json:"cdn"`
	DeliveryServices []CacheConfigDeliveryService `json:"delivery_services"`
	Parents          []CacheConfigParent          `json:"parents"`
	AllowIP          []string                     `json:"allow_ip"`
}

type CacheConfigDeliveryService struct {
	Protocol          int      `json:"protocol"`
	QueryStringIgnore int      `json:"query_string_ignore"`
	XMLID             string   `json:"xml_id"`
	Type              string   `json:"type"`
	OriginFQDN        string   `json:"origin_fdqn"`
	Regexes           []string `json:"regexes"`
	DSCP              int      `json:"dscp"`
}

type CacheConfigParent struct {
	Host   string `json:"host"`
	Domain string `json:"domain"`
	Port   uint   `json:"port"`
}

// CacheConfig gets the configuration data for a cache
func (to *Session) CacheConfig(cache string) (*CacheConfig, error) {
	url := "/api/1.3/configs/cache/" + cache
	resp, err := to.request("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data CacheConfigResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data.Response, nil
}
