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

import "time"

// CoordinatesResponse is a list of Coordinates as a response.
// Deprecated: In newer API versions, coordinates are represented by the
// CoordinatesResponseV5 structures.
// swagger:response CoordinatesResponse
// in: body
type CoordinatesResponse struct {
	// in: body
	Response []Coordinate `json:"response"`
	Alerts
}

// CoordinateResponse is a single Coordinate response for Update and Create to
// depict what changed.
// Deprecated: In newer API versions, coordinates are represented by the
// CoordinateResponseV5 structures.
// swagger:response CoordinateResponse
// in: body
type CoordinateResponse struct {
	// in: body
	Response Coordinate `json:"response"`
	Alerts
}

// Coordinate is a representation of a Coordinate as it relates to the Traffic
// Ops data model.
// Deprecated: In newer API versions, coordinates are represented by the
// CoordinateV5 structures.
type Coordinate struct {

	// The Coordinate to retrieve
	//
	// ID of the Coordinate
	//
	// required: true
	ID int `json:"id" db:"id"`

	// Name of the Coordinate
	//
	// required: true
	Name string `json:"name" db:"name"`

	// the latitude of the Coordinate
	//
	// required: true
	Latitude float64 `json:"latitude" db:"latitude"`

	// the latitude of the Coordinate
	//
	// required: true
	Longitude float64 `json:"longitude" db:"longitude"`

	// LastUpdated
	//
	LastUpdated TimeNoMod `json:"lastUpdated" db:"last_updated"`
}

// CoordinateNullable is identical to Coordinate except that its fields are
// reference values, which allows them to be nil.
// Deprecated: In newer API versions, coordinates are represented by the
// CoordinateV5 structures.
type CoordinateNullable struct {

	// The Coordinate to retrieve
	//
	// ID of the Coordinate
	//
	// required: true
	ID *int `json:"id" db:"id"`

	// Name of the Coordinate
	//
	// required: true
	Name *string `json:"name" db:"name"`

	// the latitude of the Coordinate
	//
	// required: true
	Latitude *float64 `json:"latitude" db:"latitude"`

	// the latitude of the Coordinate
	//
	// required: true
	Longitude *float64 `json:"longitude" db:"longitude"`

	// LastUpdated
	//
	LastUpdated *TimeNoMod `json:"lastUpdated" db:"last_updated"`
}

// CoordinateV5 is the representation of a Coordinate used in the latest minor
// version of APIv5.
type CoordinateV5 = CoordinateV50

// CoordinateV50 is the representation of a Coordinate used in API v5.0.
type CoordinateV50 struct {
	// The integral, unique identifier of a Coordinate.
	ID *int `json:"id" db:"id"`
	// The Coordinate's name.
	Name string `json:"name" db:"name"`
	// The latitude of the Coordinate.
	Latitude float64 `json:"latitude" db:"latitude"`
	// The longitude of the Coordinate.
	Longitude float64 `json:"longitude" db:"longitude"`
	// The time and date at which the Coordinate was last modified.
	LastUpdated time.Time `json:"lastUpdated" db:"last_updated"`
}

// CoordinateResponseV5 is the type of a response from the /coordinates endpoint
// in the latest minor version of APIv5.
type CoordinateResponseV5 struct {
	Alerts
	Response CoordinateV5
}

// CoordinatesResponseV5 is the type of a response from the /coordinates
// endpoint in the latest minor version of APIv5.
type CoordinatesResponseV5 struct {
	Alerts
	Response []CoordinateV5
}
