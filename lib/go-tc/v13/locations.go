package v13

import tc "github.com/apache/incubator-trafficcontrol/lib/go-tc"

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

// A List of Locations Response
// swagger:response LocationsResponse
// in: body
type LocationsResponse struct {
	// in: body
	Response []Location `json:"response"`
}

// A Single Location Response for Update and Create to depict what changed
// swagger:response LocationResponse
// in: body
type LocationResponse struct {
	// in: body
	Response Location `json:"response"`
}

// Location ...
type Location struct {

	// The Location to retrieve
	//
	// ID of the Location
	//
	// required: true
	ID int `json:"id" db:"id"`

	// Name of the Location
	//
	// required: true
	Name string `json:"name" db:"name"`

	// the latitude of the Location
	//
	// required: true
	Latitude float64 `json:"latitude" db:"latitude"`

	// the latitude of the Location
	//
	// required: true
	Longitude float64 `json:"longitude" db:"longitude"`

	// LastUpdated
	//
	LastUpdated tc.TimeNoMod `json:"lastUpdated" db:"last_updated"`
}

// LocationNullable ...
type LocationNullable struct {

	// The Location to retrieve
	//
	// ID of the Location
	//
	// required: true
	ID *int `json:"id" db:"id"`

	// Name of the Location
	//
	// required: true
	Name *string `json:"name" db:"name"`

	// the latitude of the Location
	//
	// required: true
	Latitude *float64 `json:"latitude" db:"latitude"`

	// the latitude of the Location
	//
	// required: true
	Longitude *float64 `json:"longitude" db:"longitude"`

	// LastUpdated
	//
	LastUpdated *tc.TimeNoMod `json:"lastUpdated" db:"last_updated"`
}
