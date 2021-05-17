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

// Tenant ...
type Tenant struct {
	Active      bool      `json:"active"`
	ID          int       `json:"id"`
	LastUpdated TimeNoMod `json:"lastUpdated" db:"last_updated"`
	Name        string    `json:"name"`
	ParentID    int       `json:"parentId"`
	ParentName  string    `json:"parentName,omitempty" db:"parent_name"`
}

// TenantNullable ...
type TenantNullable struct {
	ID          *int       `json:"id" db:"id"`
	Name        *string    `json:"name" db:"name"`
	Active      *bool      `json:"active" db:"active"`
	LastUpdated *TimeNoMod `json:"lastUpdated" db:"last_updated"`
	ParentID    *int       `json:"parentId" db:"parent_id"`
	ParentName  *string    `json:"parentName,omitempty" db:"parent_name"`
}

// DeleteTenantResponse ...
type DeleteTenantResponse struct {
	Alerts []TenantAlert `json:"alerts"`
}

// TenantAlert ...
type TenantAlert struct {
	Level string `json:"level"`
	Text  string `json:"text"`
}
