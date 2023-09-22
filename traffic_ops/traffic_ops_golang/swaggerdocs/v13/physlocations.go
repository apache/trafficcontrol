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

import "github.com/apache/trafficcontrol/v8/lib/go-tc"

// PhysLocations -  PhysLocationsResponse to get the "response" top level key
// swagger:response PhysLocations
// in: body
type PhysLocations struct {
	// PhysLocation Response Body
	// in: body
	PhysLocationsResponse tc.PhysLocationsResponse `json:"response"`
}

// PhysLocation -  PhysLocationResponse to get the "response" top level key
// swagger:response PhysLocation
// in: body
type PhysLocation struct {
	// PhysLocation Response Body
	// in: body
	PhysLocationResponse tc.PhysLocationResponse
}

// PhysLocationQueryParams
//
// swagger:parameters GetPhysLocations
type PhysLocationQueryParams struct {

	// PhysLocationsQueryParams

	// The ID of the region associated with this Physical Location
	//
	RegionID int `json:"regionId"`

	//
	//
	Orderby string `json:"orderby"`
}

// swagger:parameters PostPhysLocation
type PhysLocationPostParam struct {
	// PhysLocation Request Body
	//
	// in: body
	// required: true
	PhysLocation tc.PhysLocationNullable
}

// swagger:parameters GetPhysLocationById DeletePhysLocation
type PhysLocationPathParams struct {

	// Id associated to the PhysLocation
	// in: path
	ID int `json:"id"`
}

// PostPhysLocation swagger:route POST /phys_locations PhysLocation PostPhysLocation
//
// # Create a PhysLocation
//
// # A PhysLocation is a collection of Delivery Services
//
// Responses:
//
//	200: Alerts
func PostPhysLocation(entity PhysLocationPostParam) (PhysLocation, Alerts) {
	return PhysLocation{}, Alerts{}
}

// GetPhysLocations swagger:route GET /phys_locations PhysLocation GetPhysLocations
//
// # Retrieve a list of PhysLocations
//
// # List of PhysLocations
//
// Responses:
//
//	200: PhysLocations
//	400: Alerts
func GetPhysLocations() (PhysLocations, Alerts) {
	return PhysLocations{}, Alerts{}
}

// swagger:parameters PutPhysLocation
type PhysLocationPutParam struct {

	// ID
	// in: path
	ID int `json:"id"`

	// PhysLocation Request Body
	//
	// in: body
	// required: true
	PhysLocation tc.PhysLocationNullable
}

// PutPhysLocation swagger:route PUT /phys_locations/{id} PhysLocation PutPhysLocation
//
// # Update a PhysLocation by Id
//
// # Update a PhysLocation
//
// Responses:
//
//	200: PhysLocation
func PutPhysLocation(entity PhysLocationPutParam) (PhysLocation, Alerts) {
	return PhysLocation{}, Alerts{}
}

// GetPhysLocationById swagger:route GET /phys_locations/{id} PhysLocation GetPhysLocationById
//
// # Retrieve a specific PhysLocation by Id
//
// # Retrieve a specific PhysLocation
//
// Responses:
//
//	200: PhysLocations
//	400: Alerts
func GetPhysLocationById() (PhysLocations, Alerts) {
	return PhysLocations{}, Alerts{}
}

// DeletePhysLocation swagger:route DELETE /phys_locations/{id} PhysLocation DeletePhysLocation
//
// # Delete a PhysLocation by Id
//
// # Delete a PhysLocation
//
// Responses:
//
//	200: Alerts
func DeletePhysLocation(entityId int) Alerts {
	return Alerts{}
}
