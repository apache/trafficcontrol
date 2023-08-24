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

// ProfileParametersResponse is the type of the response from Traffic Ops to
// GET requests made to its /profileparameters API endpoint.
type ProfileParametersResponse struct {
	Response []ProfileParameter `json:"response"`
	Alerts
}

// ProfileParameterResponse is a single ProfileParameter response for Create to
// depict what changed.
// swagger:response ProfileParameterResponse
// in: body
type ProfileParameterResponse struct {
	// in: body
	Response ProfileParameter `json:"response"`
	Alerts
}

// ProfileParameter is a representation of a relationship between a Parameter
// and a Profile to which it is assigned.
//
// Note that not all unique identifiers for each represented object in this
// relationship structure are guaranteed to be populated by the Traffic Ops
// API.
type ProfileParameter struct {
	LastUpdated TimeNoMod `json:"lastUpdated"`
	Profile     string    `json:"profile"`
	ProfileID   int       `json:"profileId"`
	Parameter   string    `json:"parameter"`
	ParameterID int       `json:"parameterId"`
}

// ProfileParameterNullable is identical to ProfileParameter, except that its
// fields are reference values, which allows them to be nil.
type ProfileParameterNullable struct {
	LastUpdated *TimeNoMod `json:"lastUpdated" db:"last_updated"`
	Profile     *string    `json:"profile" db:"profile"`
	ProfileID   *int       `json:"profileId" db:"profile_id"`
	Parameter   *string    `json:"parameter" db:"parameter"`
	ParameterID *int       `json:"parameterId" db:"parameter_id"`
}

// ProfileParametersResponseV5 is the type of the response from Traffic Ops to
// GET requests made to its /profileparameters API endpoint.
type ProfileParametersResponseV5 struct {
	Response []ProfileParameterV5 `json:"response"`
	Alerts
}

// ProfileParameterResponseV5 is a single ProfileParameter response for Create to
// depict what changed.
// swagger:response ProfileParameterResponse
// in: body
type ProfileParameterResponseV5 struct {
	// in: body
	Response ProfileParameterV5 `json:"response"`
	Alerts
}

// ProfileParameterV5 is the latest minor version of the major version 5
type ProfileParameterV5 ProfileParameterV50

// ProfileParameterV50 is a representation of a relationship between a Parameter
// and a Profile to which it is assigned.
//
// Note that not all unique identifiers for each represented object in this
// relationship structure are guaranteed to be populated by the Traffic Ops
// API.
type ProfileParameterV50 struct {
	LastUpdated time.Time `json:"lastUpdated"`
	Profile     string    `json:"profile"`
	ProfileID   int       `json:"profileId"`
	Parameter   string    `json:"parameter"`
	ParameterID int       `json:"parameterId"`
}
