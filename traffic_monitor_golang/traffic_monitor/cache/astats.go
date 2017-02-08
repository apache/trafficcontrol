package cache

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

// Astats contains ATS data returned from the Astats ATS plugin. This includes generic stats, as well as fixed system stats.
type Astats struct {
	Ats    map[string]interface{} `json:"ats"`
	System AstatsSystem           `json:"system"`
}

// AstatsSystem represents fixed system stats returne from ATS by the Astats plugin.
type AstatsSystem struct {
	InfName           string `json:"inf.name"`
	InfSpeed          int    `json:"inf.speed"`
	ProcNetDev        string `json:"proc.net.dev"`
	ProcLoadavg       string `json:"proc.loadavg"`
	ConfigLoadRequest int    `json:"configReloadRequests,omitempty"`
	LastReloadRequest int    `json:"lastReloadRequest,omitempty"`
	ConfigReloads     int    `json:"configReloads,omitempty"`
	LastReload        int    `json:"lastReload,omitempty"`
	AstatsLoad        int    `json:"astatsLoad,omitempty"`
	NotAvailable      bool   `json:"notAvailable,omitempty"`
}

// Unmarshal unmarshalls the given bytes, which must be JSON Astats data, into an Astats object.
func Unmarshal(body []byte) (Astats, error) {
	var aStats Astats
	err := json.Unmarshal(body, &aStats)
	return aStats, err
}
