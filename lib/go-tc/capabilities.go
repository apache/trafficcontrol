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

// Capability reflects the ability of a user in ATC to perform some operation.
//
// In practice, they are assigned to relevant Traffic Ops API endpoints - to describe the
// capabilities of said endpoint - and to user permission Roles - to describe the capabilities
// afforded by said Role. Note that enforcement of Capability-based permissions is not currently
// implemented.
type Capability struct {
	Description string    `json:"description" db:"description"`
	LastUpdated TimeNoMod `json:"lastUpdated" db:"last_updated"`
	Name        string    `json:"name" db:"name"`
}

// CapabilitiesResponse models the structure of a minimal response from the Capabilities API in
// Traffic Ops.
type CapabilitiesResponse struct {
	Response []Capability `json:"response"`
	Alerts
}
