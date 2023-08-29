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

// ServerServerCapabilityV5 is a ServerServerCapability as it appears in version 5 of the
// Traffic Ops API - it always points to the highest minor version in APIv5.
type ServerServerCapabilityV5 = ServerServerCapabilityV50

// ServerServerCapabilityV50 represents an association between a server capability and a server.
type ServerServerCapabilityV50 struct {
	LastUpdated      *time.Time `json:"lastUpdated" db:"last_updated"`
	Server           *string    `json:"serverHostName,omitempty" db:"host_name"`
	ServerID         *int       `json:"serverId" db:"server"`
	ServerCapability *string    `json:"serverCapability" db:"server_capability"`
}

// ServerServerCapability represents an association between a server capability and a server.
type ServerServerCapability struct {
	LastUpdated      *TimeNoMod `json:"lastUpdated" db:"last_updated"`
	Server           *string    `json:"serverHostName,omitempty" db:"host_name"`
	ServerID         *int       `json:"serverId" db:"server"`
	ServerCapability *string    `json:"serverCapability" db:"server_capability"`
}

// MultipleServersCapabilities represents an association between a server capability and list of servers
// and an association between a server and list of server capabilities.
type MultipleServersCapabilities struct {
	ServerCapabilities []string `json:"serverCapabilities" db:"server_capability"`
	ServerIDs          []int64  `json:"serverIds" db:"server"`
	PageType           string   `json:"pageType"`
}

// ServerServerCapabilitiesResponseV5 is the type of a response from the
// /api/5.x/server_server_capabilities Traffic Ops endpoint.
// It always points to the type for the latest minor version of APIv5.
type ServerServerCapabilitiesResponseV5 = ServerServerCapabilitiesResponseV50

// ServerServerCapabilitiesResponseV50 is the type of a response from Traffic
// Ops to a request made to its /api/5.0/server_server_capabilities.
type ServerServerCapabilitiesResponseV50 struct {
	Response []ServerServerCapabilityV5 `json:"response"`
	Alerts
}

// ServerServerCapabilitiesResponse is the type of a response from Traffic
// Ops to a request made to its /server_server_capabilities.
type ServerServerCapabilitiesResponse struct {
	Response []ServerServerCapability `json:"response"`
	Alerts
}
