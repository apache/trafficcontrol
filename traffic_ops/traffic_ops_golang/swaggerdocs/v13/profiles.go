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

// Profiles -  ProfilesResponse to get the "response" top level key
// swagger:response Profiles
// in: body
type Profiles struct {
	// Profile Response Body
	// in: body
	ProfilesResponse tc.ProfilesResponse `json:"response"`
}

// Profile -  ProfileResponse to get the "response" top level key
// swagger:response Profile
// in: body
type Profile struct {
	// Profile Response Body
	// in: body
	ProfileResponse tc.ProfileResponse
}

// ProfileQueryParams
//
// swagger:parameters GetProfiles
type ProfileQueryParams struct {

	// ProfilesQueryParams

	// Enables Domain Name System Security Extensions (DNSSEC) for the Profile
	//
	DNSSecEnabled string `json:"dnssecEnabled"`

	// The domain name for the Profile
	//
	DomainName string `json:"domainName"`

	// Unique identifier for the Profile
	//
	ID string `json:"id"`

	// The Profile name for the Profile
	//
	Name string `json:"name"`

	//
	//
	Orderby string `json:"orderby"`
}

// swagger:parameters PostProfile
type ProfilePostParam struct {
	// Profile Request Body
	//
	// in: body
	// required: true
	Profile tc.Profile
}

// swagger:parameters GetProfileById DeleteProfile
type ProfilePathParams struct {

	// Id associated to the Profile
	// in: path
	ID int `json:"id"`
}

// PostProfile swagger:route POST /cdns Profile PostProfile
//
// # Create a Profile
//
// # A Profile is a collection of Delivery Services
//
// Responses:
//
//	200: Alerts
func PostProfile(entity ProfilePostParam) (Profile, Alerts) {
	return Profile{}, Alerts{}
}

// GetProfiles swagger:route GET /cdns Profile GetProfiles
//
// # Retrieve a list of Profiles
//
// # List of Profiles
//
// Responses:
//
//	200: Profiles
//	400: Alerts
func GetProfiles() (Profiles, Alerts) {
	return Profiles{}, Alerts{}
}

// swagger:parameters PutProfile
type ProfilePutParam struct {

	// ID
	// in: path
	ID int `json:"id"`

	// Profile Request Body
	//
	// in: body
	// required: true
	Profile tc.Profile
}

// PutProfile swagger:route PUT /cdns/{id} Profile PutProfile
//
// # Update a Profile by Id
//
// # Update a Profile
//
// Responses:
//
//	200: Profile
func PutProfile(entity ProfilePutParam) (Profile, Alerts) {
	return Profile{}, Alerts{}
}

// GetProfileById swagger:route GET /cdns/{id} Profile GetProfileById
//
// # Retrieve a specific Profile by Id
//
// # Retrieve a specific Profile
//
// Responses:
//
//	200: Profiles
//	400: Alerts
func GetProfileById() (Profiles, Alerts) {
	return Profiles{}, Alerts{}
}

// DeleteProfile swagger:route DELETE /cdns/{id} Profile DeleteProfile
//
// # Delete a Profile by Id
//
// # Delete a Profile
//
// Responses:
//
//	200: Alerts
func DeleteProfile(entityId int) Alerts {
	return Alerts{}
}
