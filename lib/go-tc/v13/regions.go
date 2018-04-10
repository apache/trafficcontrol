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

// A List of Regions Response
// swagger:response RegionsResponse
// in: body
type RegionsResponse struct {
	// in: body
	Response []Region `json:"response"`
}

// A Single Region Response for Update and Create to depict what changed
// swagger:response RegionResponse
// in: body
type RegionResponse struct {
	// in: body
	Response Region `json:"response"`
}

type Region struct {

	// The Region to retrieve

	// DivisionName - Name of the Division associated to this Region
	//
	// required: true
	DivisionName string `json:"divisionName"`

	// DivisionName of the Division
	//
	// required: true
	Division int `json:"division" db:"division"`

	// Region ID
	//
	ID int `json:"id" db:"id"`

	// LastUpdated
	//
	LastUpdated tc.TimeNoMod `json:"lastUpdated" db:"last_updated"`

	// Region Name
	//
	// required: true
	Name string `json:"name" db:"name"`
}
