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
// A Single Servercheck Response for Update and Create to depict what changed
// swagger:response ServercheckResponse
// in: body
type ServercheckResponse struct {
	// in: body
	Response Servercheck `json:"response"`
}

// A list of Servercheck Responses
// swagger:response ServerchecksResponse
// in: body
type ServerchecksResponse struct {
	// in: body
	Response []Servercheck `json:"response"`
}

// A Single Servercheck struct for GET response
// swagger:model Servercheck
type Servercheck struct {

	// The Servercheck response data
	//
	// Admin state of the checked server
	AdminState string `json:"adminState"`

	// Cache group the checked server belongs to
	CacheGroup string `json:"cacheGroup"`

	// ID number of the checked server
	ID int `json:"id"`

	// Hostname of the checked server
	HostName string `json:"hostName"`

	// Reval pending flag for checked server
	RevalPending bool `json:"revalPending"`

	// Profile name of checked server
	Profile string `json:"profile"`

	// Traffic Control type of the checked server
	Type string `json:"type"`

	// Update pending flag for the checked server
	UpdPending bool `json:"updPending"`

	// Various check types
	Checks struct {

		// IPv4 production interface (legacy name)
		Iface10G int `json:"10G"`

		// IPv6 production interface (legacy name)
		Iface10G6 int `json:"10G6"`

		// Cache Disk Usage
		CDU int `json:"CDU"`

		// Cache Hit Ratio
		CHR int `json:"CHR"`

		// DSCP check
		DSCP int `json:"DSCP"`

		// DNS check
		FQDN int `json:"FQDN"`

		// Out-of-band (BMC) interface check
		ILO int `json:"ILO"`

		// IPv4 production interface (new name)
		IPv4 int `json:"IPv4"`

		// IPv6 production interface (new name)
		IPv6 int `json:"IPv6"`

		// MTU check
		MTU int `json:"MTU"`

		// ORT check
		ORT int `json:"ORT"`

		// Traffic Router status for checked server
		RTR int `json:"RTR"`
	} `json:"checks"`
}

// A Single Servercheck struct for Update and Create to depict what changed
// swagger:model ServercheckPost
type ServercheckPost struct {

	// The Servercheck data to submit
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

type ServercheckPostNullable struct {
	Name  string `json:"servercheck_short_name" db:"servercheck_short_name"`
	ID    int    `json:"id" db:"id"`
	Value int    `json:"value" db:"value"`
}

type ServercheckPostResponse struct {
	Alerts   []Alert                 `json:"alerts"`
	Response DeliveryServiceUserPost `json:"response"`
}
