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

// DivisionsResponse is a list of Divisions as a response.
// swagger:response DivisionsResponse
type DivisionsResponse struct {
	// in: body
	Response []Division `json:"response"`
	Alerts
}

// DivisionResponse is a single Division response for Update and Create to
// depict what changed.
// swagger:response DivisionResponse
// in: body
type DivisionResponse struct {
	// in: body
	Response Division `json:"response"`
}

// A Division is a named collection of Regions.
type Division struct {

	// Division ID
	//
	ID int `json:"id" db:"id"`

	// LastUpdated
	//
	LastUpdated TimeNoMod `json:"lastUpdated" db:"last_updated"`

	// Division Name
	//
	// required: true
	Name string `json:"name" db:"name"`
}

// DivisionNullable is a nullable struct that holds info about a division, which
// is a group of regions.
type DivisionNullable struct {
	ID          *int       `json:"id" db:"id"`
	LastUpdated *TimeNoMod `json:"lastUpdated" db:"last_updated"`
	Name        *string    `json:"name" db:"name"`
}
