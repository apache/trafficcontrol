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

// RegionsResponse is the type of responses from Traffic Ops to GET requests
// made to its /regions API endpoint.

import "time"

type RegionsResponse struct {
	Response []Region `json:"response"`
	Alerts
}

// A Region is a named collection of Physical Locations within a Division.
type Region struct {
	DivisionName string    `json:"divisionName"`
	Division     int       `json:"division" db:"division"`
	ID           int       `json:"id" db:"id"`
	LastUpdated  TimeNoMod `json:"lastUpdated" db:"last_updated"`
	Name         string    `json:"name" db:"name"`
}

// RegionsResponseV5 is an alias for the latest minor version of the major version 5.
type RegionsResponseV5 = RegionsResponseV50

// RegionsResponseV50 is response made from /regions API endpoint - in the latest minor version APIv50.
type RegionsResponseV50 struct {
	Response []RegionV5 `json:"response"`
	Alerts
}

// RegionV5 is an alias for the latest minor version of the major version 5.
type RegionV5 = RegionV50

// RegionV50 is a named collection of Physical Locations within a Division in the latest minor version APIv50.
type RegionV50 struct {
	DivisionName string    `json:"divisionName"`
	Division     int       `json:"division" db:"division"`
	ID           int       `json:"id" db:"id"`
	LastUpdated  time.Time `json:"lastUpdated" db:"last_updated"`
	Name         string    `json:"name" db:"name"`
}

// RegionName is a response to a request to get a region by its name. It
// includes the division that the region is in.
type RegionName struct {
	ID       int                `json:"id"`
	Name     string             `json:"name"`
	Division RegionNameDivision `json:"division"`
}

// RegionNameDivision is the division that contains the region that a request
// is trying to query by name.
type RegionNameDivision struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// RegionNameResponse models the structure of a response to a request to get a
// region by its name.
type RegionNameResponse struct {
	Response []RegionName `json:"response"`
}
