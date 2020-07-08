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

// PhysLocationsResponse is a list of PhysLocations as a response.
type PhysLocationsResponse struct {
	Response []PhysLocation `json:"response"`
	Alerts
}

// PhysLocationResponse is a single PhysLocationNullable as a response.
type PhysLocationResponse struct {
	Response PhysLocationNullable `json:"response"`
	Alerts
}

// PhysLocation contains the physical location of a cache group.
type PhysLocation struct {

	//
	// The Street Address of the physical location
	//
	// required: true
	Address string `json:"address" db:"address"`

	//
	// The Address of the physical location
	//
	// required: true
	City string `json:"city" db:"city"`

	//
	// comments are additional details about the physical location
	//
	Comments string `json:"comments" db:"comments"`

	//
	// The email address for the Point of Contact at the physical location
	//
	Email string `json:"email" db:"email"`

	//
	// The name of the physical location
	//
	// required: true
	ID int `json:"id" db:"id"`

	// Timestamp of the last time this row was updated
	//
	LastUpdated TimeNoMod `json:"lastUpdated" db:"last_updated"`

	//
	// The name of the physical location
	//
	// required: true
	Name string `json:"name" db:"name"`

	//
	// The phone number of the physical location
	//
	// required: true
	Phone string `json:"phone" db:"phone"`

	//
	// The Point Of Contact at the physical location
	//
	// required: true
	POC string `json:"poc" db:"poc"`

	//
	// The RegionID associated to this physical location
	//
	// required: true
	RegionID int `json:"regionId" db:"region"`

	//
	// The Region Name for the region associated to this physical location
	//
	RegionName string `json:"region" db:"region_name"`

	//
	// The shortName for the physical location (like an alias)
	//
	// required: true
	ShortName string `json:"shortName" db:"short_name"`

	//
	// The State for the physical location
	//
	// required: true
	State string `json:"state" db:"state"`

	//
	// The Zipcode for the physical location
	//
	// required: true
	Zip string `json:"zip" db:"zip"`
}

// PhysLocationNullable contains the physical location of a cache group. It
// allows for all fields to be null.
type PhysLocationNullable struct {
	//
	// The Street Address of the physical location
	//
	// required: true
	Address *string `json:"address" db:"address"`

	//
	// The Address of the physical location
	//
	// required: true
	City *string `json:"city" db:"city"`

	//
	// comments are additional details about the physical location
	//
	Comments *string `json:"comments" db:"comments"`

	//
	// The email address for the Point of Contact at the physical location
	//
	Email *string `json:"email" db:"email"`

	//
	// The name of the physical location
	//
	// required: true
	ID *int `json:"id" db:"id"`

	// Timestamp of the last time this row was updated
	//
	LastUpdated *TimeNoMod `json:"lastUpdated" db:"last_updated"`

	//
	// The name of the physical location
	//
	// required: true
	Name *string `json:"name" db:"name"`

	//
	// The phone number of the physical location
	//
	// required: true
	Phone *string `json:"phone" db:"phone"`

	//
	// The Point Of Contact at the physical location
	//
	// required: true
	POC *string `json:"poc" db:"poc"`

	//
	// The RegionID associated to this physical location
	//
	// required: true
	RegionID *int `json:"regionId" db:"region"`

	//
	// The Region Name for the region associated to this physical location
	//
	RegionName *string `json:"region" db:"region_name"`

	//
	// The shortName for the physical location (like an alias)
	//
	// required: true
	ShortName *string `json:"shortName" db:"short_name"`

	//
	// The State for the physical location
	//
	// required: true
	State *string `json:"state" db:"state"`

	//
	// The Zipcode for the physical location
	//
	// required: true
	Zip *string `json:"zip" db:"zip"`
}

// PhysLocationTrimmed contains only the name of a physical location.
type PhysLocationTrimmed struct {
	Name string `json:"name"`
}
