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

// GetTenantsResponse is the response for a request for a group of tenants.
type GetTenantsResponse struct {
	Response []Tenant `json:"response"`
	Alerts
}

// TenantResponse is the type of a response from Traffic Ops to a PUT, POST,
// or DELETE request made to its /tenants.
type TenantResponse struct {
	Response Tenant `json:"response"`
	Alerts
}

// A Tenant is a scope that can be applied to groups of users to limit their
// Delivery Service-related actions to specific sets of similarly "Tenanted"
// Delivery Services.
type Tenant struct {
	Active      bool      `json:"active"`
	ID          int       `json:"id"`
	LastUpdated TimeNoMod `json:"lastUpdated" db:"last_updated"`
	Name        string    `json:"name"`
	ParentID    int       `json:"parentId"`
	ParentName  string    `json:"parentName,omitempty" db:"parent_name"`
}

// TenantNullable is identical to Tenant, but its fields are reference values,
// which allows them to be nil.
type TenantNullable struct {
	ID          *int       `json:"id" db:"id"`
	Name        *string    `json:"name" db:"name"`
	Active      *bool      `json:"active" db:"active"`
	LastUpdated *TimeNoMod `json:"lastUpdated" db:"last_updated"`
	ParentID    *int       `json:"parentId" db:"parent_id"`
	ParentName  *string    `json:"parentName,omitempty" db:"parent_name"`
}

// DeleteTenantResponse is a legacy structure used to represent responses to
// DELETE requests made to Traffic Ops's /tenants API endpoint.
//
// Deprecated: This uses a deprecated type for its Alerts property and drops
// information returned by the TO API - new code should use TenantResponse
// instead.
type DeleteTenantResponse struct {
	Alerts []TenantAlert `json:"alerts"`
}

// TenantAlert is an unnecessary and less type-safe duplicate of Alert.
//
// Deprecated: Use Alert instead.
type TenantAlert struct {
	Level string `json:"level"`
	Text  string `json:"text"`
}

// GetTenantsResponseV50 is the response for a request for a group of tenants.
type GetTenantsResponseV50 struct {
	Response []TenantV5 `json:"response"`
	Alerts
}

// TenantResponseV50 is the type of response from Traffic Ops to a PUT, POST,
// or DELETE request made to its /tenants.
type TenantResponseV50 struct {
	Response TenantV5 `json:"response"`
	Alerts
}

// A TenantV50 is a scope that can be applied to groups of users to limit their
// Delivery Service-related actions to specific sets of similarly "Tenanted"
// Delivery Services.
type TenantV50 struct {
	ID          *int       `json:"id" db:"id"`
	Name        *string    `json:"name" db:"name"`
	Active      *bool      `json:"active" db:"active"`
	LastUpdated *time.Time `json:"lastUpdated" db:"last_updated"`
	ParentID    *int       `json:"parentId" db:"parent_id"`
	ParentName  *string    `json:"parentName,omitempty" db:"parent_name"`
}

// GetTenantsResponseV5 is the type of response from the tenants
// Traffic Ops endpoint.
// It always points to the type for the latest minor version of GetTenantsResponseV5x APIv5.
type GetTenantsResponseV5 = GetTenantsResponseV50

// TenantResponseV5 is the type of response from Traffic Ops to a PUT, POST,
// // or DELETE request made to its /tenants.
// It always points to the type for the latest minor version of TenantResponseV5x APIv5.
type TenantResponseV5 = TenantResponseV50

// TenantV5 always points to the type for the latest minor version of TenantV5x APIv5.
type TenantV5 = TenantV50
