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

// A List of Divisions Response
// swagger:response DivisionsResponse
type DivisionsResponse struct {
	// in: body
	Response []Division `json:"response"`
}

// A Single Division Response for Update and Create to depict what changed
// swagger:response DivisionResponse
// in: body
type DivisionResponse struct {
	// in: body
	Response Division `json:"response"`
}

// Division ...
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

type DivisionNullable struct {
	ID          *int       `json:"id" db:"id"`
	LastUpdated *TimeNoMod `json:"lastUpdated" db:"last_updated"`
	Name        *string    `json:"name" db:"name"`
}
