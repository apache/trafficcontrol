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

// RolesResponse is a list of Roles as a response.
// swagger:response RolesResponse
// in: body
type RolesResponse struct {
	// in: body
	Response []Role `json:"response"`
	Alerts
}

// RoleResponse is a single Role response for Update and Create to depict what
// changed.
// swagger:response RoleResponse
// in: body
type RoleResponse struct {
	// in: body
	Response Role `json:"response"`
	Alerts
}

// A Role is a definition of the permissions afforded to a user with that Role.
type Role struct {
	RoleV11

	// Capabilities associated with the Role
	//
	// required: true
	Capabilities *[]string `json:"capabilities" db:"-"`
}

// RoleV11 is a representation of a Role as it appeared in version 1.1 of the
// Traffic Ops API.
//
// Deprecated: Traffic Ops API version 1.1 no longer exists - the ONLY reason
// this structure still exists is because it is nested in newer structures - DO
// NOT USE THIS!
type RoleV11 struct {
	// ID of the Role
	//
	// required: true
	ID *int `json:"id" db:"id"`

	// Name of the Role
	//
	// required: true
	Name *string `json:"name" db:"name"`

	// Description of the Role
	//
	// required: true
	Description *string `json:"description" db:"description"`

	// Priv Level of the Role
	//
	// required: true
	PrivLevel *int `json:"privLevel" db:"priv_level"`
}
