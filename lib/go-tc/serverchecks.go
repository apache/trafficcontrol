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

// Health monitor checks for servers
// A Single Status Response for Update and Create to depict what changed
// swagger:response StatusResponse
// in: body
type ServercheckResponse struct {
	// in: body
	Response Status `json:"response"`
}

// A Single Statuses Response for Update and Create to depict what changed
// swagger:model Statuses
type Servercheck struct {

	// The Statuses to retrieve
	//
	// Name of the server check type
	//
	// required: true
	Name string `json:"servercheck_short_name" db:"servercheck_short_name"`

	// ID of the server
	//
	// required: true
	ID int `json:"id" db:"id"`

	// Name of the server
	HostName string `json:"name" db:"name"`

	// Value of the check result
    //
	// required: true
	Value int `json:"value" db:"value"`
}

type ServercheckNullable struct {
	Name string `json:"servercheck_short_name" db:"servercheck_short_name"`
	ID int `json:"id" db:"id"`
	Value int `json:"value" db:"value"`
}

type ServercheckPostResponse struct {
        Alerts   []Alert                 `json:"alerts"`
        Response DeliveryServiceUserPost `json:"response"`
}

