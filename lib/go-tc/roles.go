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

import (
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/v8/lib/go-util"

	validation "github.com/go-ozzo/ozzo-validation"
)

// AdminRoleName is the Name of the special "admin" Role.
//
// This Role must always exist; it cannot be modified or deleted, and is
// guaranteed to exist in all valid ATC environments.
const AdminRoleName = "admin"

// RoleV4 is an alias for the latest minor version for the major version 4.
type RoleV4 RoleV40

// RolesResponseV4 is a list of RoleV4 as a response.
type RolesResponseV4 struct {
	Response []RoleV4 `json:"response"`
	Alerts
}

// RoleResponseV4 is a RoleV4 as a response.
type RoleResponseV4 struct {
	Response RoleV4 `json:"response"`
	Alerts
}

// RoleV40 is the structure used to depict roles in API v4.0.
type RoleV40 struct {
	Name        string     `json:"name" db:"name"`
	Permissions []string   `json:"permissions" db:"permissions"`
	Description string     `json:"description" db:"description"`
	LastUpdated *time.Time `json:"lastUpdated,omitempty" db:"last_updated"`
}

// Validate will validate and make sure all that the fields in the supplied RoleV4 struct are semantically correct.
func (role RoleV4) Validate() error {
	errs := validation.Errors{
		"name":        validation.Validate(role.Name, validation.Required),
		"description": validation.Validate(role.Description, validation.Required),
	}
	return util.JoinErrs(tovalidate.ToErrors(errs))
}

// Upgrade will convert the passed in instance of Role struct into an instance of RoleV4 struct.
func (role Role) Upgrade() RoleV4 {
	var roleV4 RoleV4
	if role.Name != nil {
		roleV4.Name = *role.Name
	}
	if role.Description != nil {
		roleV4.Description = *role.Description
	}
	if role.Capabilities == nil {
		roleV4.Permissions = nil
	} else {
		roleV4.Permissions = make([]string, len(*role.Capabilities))
		copy(roleV4.Permissions, *role.Capabilities)
	}
	return roleV4
}

// Downgrade will convert the passed in instance of RoleV4 struct into an instance of Role struct.
func (role RoleV4) Downgrade() Role {
	var downgraded Role
	downgraded.Name = &role.Name
	downgraded.Description = &role.Description
	if len(role.Permissions) == 0 {
		downgraded.Capabilities = nil
	} else {
		caps := make([]string, len(role.Permissions))
		copy(caps, role.Permissions)
		downgraded.Capabilities = &caps
	}
	return downgraded
}

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
