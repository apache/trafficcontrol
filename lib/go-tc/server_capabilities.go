package tc

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

import "time"

// ServerCapabilitiesResponse contains the result data from a GET /server_capabilities request.
type ServerCapabilitiesResponse struct {
	Response []ServerCapability `json:"response"`
	Alerts
}

// ServerCapabilitiesResponseV41 contains the result data from a GET(v4.1 and above) /server_capabilities request.
type ServerCapabilitiesResponseV41 struct {
	Response []ServerCapabilityV41 `json:"response"`
	Alerts
}

// ServerCapability contains information about a given ServerCapability in Traffic Ops.
type ServerCapability struct {
	Name        string     `json:"name" db:"name"`
	LastUpdated *TimeNoMod `json:"lastUpdated" db:"last_updated"`
}

// ServerCapabilityV4 is an alias for the latest minor version for the major version 4.
type ServerCapabilityV4 ServerCapabilityV41

// ServerCapabilityV41 contains information (in-addition to description) about a given ServerCapability  in Traffic Ops.
type ServerCapabilityV41 struct {
	ServerCapability
	Description string `json:"description" db:"description"`
}

// ServerCapabilityDetailResponse contains the result data from a POST /server_capabilities request.
type ServerCapabilityDetailResponse struct {
	Response ServerCapability `json:"response"`
	Alerts
}

// ServerCapabilityDetailResponseV41 contains the result data from a POST(v4.1 and above) /server_capabilities request.
type ServerCapabilityDetailResponseV41 struct {
	Response ServerCapabilityV41 `json:"response"`
	Alerts
}

// ServerCapabilityV5 is an alias for the latest minor version for the major version 5.
type ServerCapabilityV5 ServerCapabilityV50

// ServerCapabilityV50 contains information about a given serverCapability in Traffic Ops V5.
type ServerCapabilityV50 struct {
	Name        string    `json:"name" db:"name"`
	LastUpdated time.Time `json:"lastUpdated" db:"last_updated"`
	Description string    `json:"description" db:"description"`
}

// ServerCapabilitiesResponseV5 is an alias for the latest minor version for the major version 5.
type ServerCapabilitiesResponseV5 ServerCapabilitiesResponseV50

// ServerCapabilitiesResponseV50 contains the result data from a GET(v5.1 and above) /server_capabilities request.
type ServerCapabilitiesResponseV50 struct {
	Response []ServerCapabilityV5 `json:"response"`
	Alerts
}

// ServerCapabilityDetailResponseV5 is an alias for the latest minor version for the major version 5.
type ServerCapabilityDetailResponseV5 ServerCapabilityDetailResponseV50

// ServerCapabilityDetailResponseV50 contains the result data from a POST(v5.1 and above) /server_capabilities request.
type ServerCapabilityDetailResponseV50 struct {
	Response ServerCapabilityV5 `json:"response"`
	Alerts
}
