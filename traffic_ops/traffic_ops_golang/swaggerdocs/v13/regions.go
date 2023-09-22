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

// Regions -  RegionsResponse to get the "response" top level key
// swagger:response Regions
// in: body
type Regions struct {
	// Region Response Body
	// in: body
	RegionsResponse tc.RegionsResponse `json:"response"`
}

// Region -  RegionResponse to get the "response" top level key
// swagger:response Region
// in: body
type Region struct {
	// Region Response Body
	// in: body
	RegionResponse tc.RegionsResponse
}

// RegionQueryParams
//
// swagger:parameters GetRegions
type RegionQueryParams struct {

	// RegionsQueryParams

	// Division ID that refers to this Region
	//
	Division string `json:"division"`

	// Division Name that refers to this Region
	//
	DivisionName string `json:"divisionName"`

	// Unique identifier for the Region
	//
	ID string `json:"id"`

	//
	//
	Orderby string `json:"orderby"`
}

// swagger:parameters PostRegion
type RegionPostParam struct {
	// Region Request Body
	//
	// in: body
	// required: true
	Region tc.Region
}

// swagger:parameters GetRegionById DeleteRegion
type RegionPathParams struct {

	// Id associated to the Region
	// in: path
	ID int `json:"id"`
}

// PostRegion swagger:route POST /regions Region PostRegion
//
// # Create a Region
//
// Responses:
//
//	200: Alerts
func PostRegion(entity RegionPostParam) (Region, Alerts) {
	return Region{}, Alerts{}
}

// GetRegions swagger:route GET /regions Region GetRegions
//
// # Retrieve a list of Regions
//
// Responses:
//
//	200: Regions
//	400: Alerts
func GetRegions() (Regions, Alerts) {
	return Regions{}, Alerts{}
}

// swagger:parameters PutRegion
type RegionPutParam struct {

	// ID
	// in: path
	ID int `json:"id"`

	// Region Request Body
	//
	// in: body
	// required: true
	Region tc.Region
}

// PutRegion swagger:route PUT /regions/{id} Region PutRegion
//
// # Update a Region
//
// Responses:
//
//	200: Region
func PutRegion(entity RegionPutParam) (Region, Alerts) {
	return Region{}, Alerts{}
}

// GetRegionById swagger:route GET /regions/{id} Region GetRegionById
//
// # Retrieve a specific Region
//
// Responses:
//
//	200: Regions
//	400: Alerts
func GetRegionById() (Regions, Alerts) {
	return Regions{}, Alerts{}
}

// DeleteRegion swagger:route DELETE /regions/{id} Region DeleteRegion
//
// # Delete a Region
//
// Responses:
//
//	200: Alerts
func DeleteRegion(entityId int) Alerts {
	return Alerts{}
}
