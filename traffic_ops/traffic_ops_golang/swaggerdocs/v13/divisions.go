package v13

import "github.com/apache/trafficcontrol/v8/lib/go-tc"

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

// Divisions -  DivisionsResponse to get the "response" top level key
// swagger:response Divisions
// in: body
type Divisions struct {
	// Division Response Body
	// in: body
	DivisionsResponse tc.DivisionsResponse `json:"response"`
}

// Division -  DivisionResponse to get the "response" top level key
// swagger:response Division
// in: body
type Division struct {
	// Division Response Body
	// in: body
	DivisionResponse tc.DivisionResponse
}

// DivisionQueryParams
//
// swagger:parameters GetDivisions
type DivisionQueryParams struct {

	// DivisionsQueryParams

	// Name for this Division
	//
	Name string `json:"name"`

	// Unique identifier for the Division
	//
	ID string `json:"id"`

	// The field in the response to sort the response by
	//
	Orderby string `json:"orderby"`
}

// swagger:parameters PostDivision
type DivisionPostParam struct {
	// Division Request Body
	//
	// in: body
	// required: true
	Division tc.Division
}

// swagger:parameters GetDivisionById DeleteDivision
type DivisionPathParams struct {

	// Id associated to the Division
	// in: path
	ID int `json:"id"`
}

// PostDivision swagger:route POST /divisions Division PostDivision
//
// # Create a Division
//
// # A Division is a group of regions
//
// Responses:
//
//	200: Alerts
func PostDivision(entity DivisionPostParam) (Division, Alerts) {
	return Division{}, Alerts{}
}

// GetDivisions swagger:route GET /divisions Division GetDivisions
//
// # Retrieve a list of Divisions
//
// # List of Divisions
//
// Responses:
//
//	200: Divisions
//	400: Alerts
func GetDivisions() (Divisions, Alerts) {
	return Divisions{}, Alerts{}
}

// swagger:parameters PutDivision
type DivisionPutParam struct {

	// ID
	// in: path
	ID int `json:"id"`

	// Division Request Body
	//
	// in: body
	// required: true
	Division tc.Division
}

// PutDivision swagger:route PUT /divisions/{id} Division PutDivision
//
// # Update a Division by Id
//
// # Update a single Division
//
// Responses:
//
//	200: Division
func PutDivision(entity DivisionPutParam) (Division, Alerts) {
	return Division{}, Alerts{}
}

// GetDivisionById swagger:route GET /divisions/{id} Division GetDivisionById
//
// # Retrieve a specific Division by Id
//
// # Retrieve a single division
//
// Responses:
//
//	200: Divisions
//	400: Alerts
func GetDivisionById() (Divisions, Alerts) {
	return Divisions{}, Alerts{}
}

// DeleteDivision swagger:route DELETE /divisions/{id} Division DeleteDivision
//
// # Delete a Division by Id
//
// # Delete a single Division
//
// Responses:
//
//	200: Alerts
func DeleteDivision(entityId int) Alerts {
	return Alerts{}
}
