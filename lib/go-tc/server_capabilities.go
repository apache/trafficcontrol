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

// ServerCapabilitiesResponse contains the result data from a GET /server_capabilities request.
type ServerCapabilitiesResponse struct {
	Response []ServerCapability `json:"response"`
	Alerts
}

// ServerCapability contains information about a given ServerCapability in Traffic Ops.
type ServerCapability struct {
	Name        string     `json:"name" db:"name"`
	Description string     `json:"description" db:"description"`
	LastUpdated *TimeNoMod `json:"lastUpdated" db:"last_updated"`
}

// ServerCapabilityDetailResponse contains the result data from a POST /server_capabilities request.
type ServerCapabilityDetailResponse struct {
	Response ServerCapability `json:"response"`
	Alerts
}
