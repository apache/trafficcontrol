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

// ProfileParameters -  ProfileParametersResponse to get the "response" top level key
// swagger:response ProfileParameters
// in: body
type ProfileParameters struct {
	// ProfileParameter Response Body
	// in: body
	ProfileParametersResponse tc.ProfileParametersResponse `json:"response"`
}

// ProfileParameter -  ProfileParameterResponse to get the "response" top level key
// swagger:response ProfileParameter
// in: body
type ProfileParameter struct {
	// ProfileParameter Response Body
	// in: body
	ProfileParameterResponse tc.ProfileParameterResponse
}

// ProfileParameterQueryParams
//
// swagger:parameters GetProfileParameters
type ProfileParameterQueryParams struct {

	// ProfileParametersQueryParams

	// Unique identifier for the ProfileParameter
	//
	ProfileID string `json:"profileId"`

	// Unique identifier for the ProfileParameter
	//
	ParameterID string `json:"parameterId"`

	// The field in the response to sort the response by
	//
	Orderby string `json:"orderby"`
}

// swagger:parameters PostProfileParameter
type ProfileParameterPostParam struct {
	// ProfileParameter Request Body
	//
	// in: body
	// required: true
	ProfileParameter tc.ProfileParameter
}

// PostProfileParameter swagger:route POST /profileparameters ProfileParameter PostProfileParameter
//
// # Create a ProfileParameter
//
// # A ProfileParameter is a join of the Profile and Parameters
//
// Responses:
//
//	200: Alerts
func PostProfileParameter(entity ProfileParameterPostParam) (ProfileParameter, Alerts) {
	return ProfileParameter{}, Alerts{}
}

// GetProfileParameters swagger:route GET /profileparameters ProfileParameter GetProfileParameters
//
// # Retrieve a list of ProfileParameters by narrowing down with query parameters
//
// # List of ProfileParameters
//
// Responses:
//
//	200: ProfileParameters
//	400: Alerts
func GetProfileParameters() (ProfileParameters, Alerts) {
	return ProfileParameters{}, Alerts{}
}

// GetProfileParameterById swagger:route GET /profileparameters?id={id} ProfileParameter GetProfileParameterById
//
// # Retrieve a specific ProfileParameter by Id
//
// # Retrieve a single division
//
// Responses:
//
//	200: ProfileParameters
//	400: Alerts
func GetProfileParameterById() (ProfileParameters, Alerts) {
	return ProfileParameters{}, Alerts{}
}

// DeleteProfileParameter swagger:route DELETE /profileparameters/{id} ProfileParameter DeleteProfileParameter
//
// # Delete a ProfileParameter by Id
//
// # Delete a single ProfileParameter
//
// Responses:
//
//	200: Alerts
func DeleteProfileParameter(entityId int) Alerts {
	return Alerts{}
}
