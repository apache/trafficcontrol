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

// ServiceCategoriesResponse is a list of Service Categories as a response.
// swagger:response ServiceCategoriesResponse
type ServiceCategoriesResponse struct {
	// in: body
	Response []ServiceCategory `json:"response"`
}

// ServiceCategoryResponse is a single Service Category response for Update and Create to
// depict what changed.
// swagger:response ServiceCategoryResponse
// in: body
type ServiceCategoryResponse struct {
	// in: body
	Response ServiceCategory `json:"response"`
}

// ServiceCategory ...
type ServiceCategory struct {

	// ServiceCategory ID
	//
	ID int `json:"id" db:"id"`

	// LastUpdated
	//
	LastUpdated TimeNoMod `json:"lastUpdated" db:"last_updated"`

	// Division Name
	//
	// required: true
	Name string `json:"name" db:"name"`

	// ServiceCategory Tenant Id
	//
	TenantID int `json:"tenantId" db:"tenant_id"`

	// ServiceCategory Tenant Name
	//
	TenantName string `json:"tenant"`
}

// ServiceCategoryNullable is a nullable struct that holds info about a service category
type ServiceCategoryNullable struct {
	ID          *int       `json:"id" db:"id"`
	LastUpdated *TimeNoMod `json:"lastUpdated" db:"last_updated"`
	Name        *string    `json:"name" db:"name"`
	TenantID    *int       `json:"tenantId" db:"tenant_id"`
	TenantName  *string    `json:"tenant"`
}
