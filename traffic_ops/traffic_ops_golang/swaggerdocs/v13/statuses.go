package v13

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
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

// Statuses -  StatusesResponse to get the "response" top level key
// swagger:response Statuses
// in: body
type Statuses struct {
	// Status Response Body
	// in: body
	StatusesResponse tc.StatusesResponse `json:"response"`
}

// Status -  StatusResponse to get the "response" top level key
// swagger:response Status
// in: body
type Status struct {
	// Status Response Body
	// in: body
	StatusResponse tc.StatusResponse
}

// StatusQueryParams
//
// swagger:parameters GetStatuses
type StatusQueryParams struct {

	// StatusesQueryParams

	// The name that refers to this Status
	//
	Name string `json:"name"`

	// A short description of the status
	//
	Description string `json:"description"`

	// Unique identifier for the Status
	//
	ID string `json:"id"`

	//
	//
	Orderby string `json:"orderby"`
}

// swagger:parameters PostStatus
type StatusPostParam struct {
	// Status Request Body
	//
	// in: body
	// required: true
	Status tc.Status
}

// swagger:parameters GetStatusById DeleteStatus
type StatusPathParams struct {

	// Id associated to the Status
	// in: path
	ID int `json:"id"`
}

// PostStatus swagger:route POST /statuses Status PostStatus
//
// # Create a Status
//
// Responses:
//
//	200: Alerts
func PostStatus(entity StatusPostParam) (Status, Alerts) {
	return Status{}, Alerts{}
}

// GetStatuses swagger:route GET /statuses Status GetStatuses
//
// # Retrieve a list of Statuses
//
// Responses:
//
//	200: Statuses
//	400: Alerts
func GetStatuses() (Statuses, Alerts) {
	return Statuses{}, Alerts{}
}

// swagger:parameters PutStatus
type StatusPutParam struct {

	// ID
	// in: path
	ID int `json:"id"`

	// Status Request Body
	//
	// in: body
	// required: true
	Status tc.Status
}

// PutStatus swagger:route PUT /statuses/{id} Status PutStatus
//
// # Update a Status
//
// Responses:
//
//	200: Status
func PutStatus(entity StatusPutParam) (Status, Alerts) {
	return Status{}, Alerts{}
}

// GetStatusById swagger:route GET /statuses/{id} Status GetStatusById
//
// # Retrieve a specific Status
//
// Responses:
//
//	200: Statuses
//	400: Alerts
func GetStatusById() (Statuses, Alerts) {
	return Statuses{}, Alerts{}
}

// DeleteStatus swagger:route DELETE /statuses/{id} Status DeleteStatus
//
// # Delete a Status
//
// Responses:
//
//	200: Alerts
func DeleteStatus(entityId int) Alerts {
	return Alerts{}
}
