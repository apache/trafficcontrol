package v13

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

import "github.com/apache/trafficcontrol/v8/lib/go-tc"

// CDNs -  CDNsResponse to get the "response" top level key
// swagger:response CDNs
// in: body
type CDNs struct {
	// CDN Response Body
	// in: body
	CDNsResponse tc.CDNsResponse `json:"response"`
}

// CDN -  CDNResponse to get the "response" top level key
// swagger:response CDN
// in: body
type CDN struct {
	// CDN Response Body
	// in: body
	CDNResponse tc.CDNResponse
}

// CDNQueryParams
//
// swagger:parameters GetCDNs
type CDNQueryParams struct {

	// CDNsQueryParams

	// Enables Domain Name System Security Extensions (DNSSEC) for the CDN
	//
	DNSSecEnabled string `json:"dnssecEnabled"`

	// The domain name for the CDN
	//
	DomainName string `json:"domainName"`

	// Unique identifier for the CDN
	//
	ID string `json:"id"`

	// The CDN name
	//
	Name string `json:"name"`

	//
	//
	Orderby string `json:"orderby"`
}

// swagger:parameters PostCDN
type CDNPostParam struct {
	// CDN Request Body
	//
	// in: body
	// required: true
	CDN tc.CDN
}

// swagger:parameters GetCDNById DeleteCDN
type CDNPathParams struct {

	// Id associated to the CDN
	// in: path
	ID int `json:"id"`
}

// PostCDN swagger:route POST /cdns CDN PostCDN
//
// # Create a CDN
//
// # A CDN is a collection of Delivery Services
//
// Responses:
//
//	200: Alerts
func PostCDN(entity CDNPostParam) (CDN, Alerts) {
	return CDN{}, Alerts{}
}

// GetCDNs swagger:route GET /cdns CDN GetCDNs
//
// # Retrieve a list of CDNs
//
// # List of CDNs
//
// Responses:
//
//	200: CDNs
//	400: Alerts
func GetCDNs() (CDNs, Alerts) {
	return CDNs{}, Alerts{}
}

// swagger:parameters PutCDN
type CDNPutParam struct {

	// ID
	// in: path
	ID int `json:"id"`

	// CDN Request Body
	//
	// in: body
	// required: true
	CDN tc.CDN
}

// PutCDN swagger:route PUT /cdns/{id} CDN PutCDN
//
// # Update a CDN by Id
//
// # Update a CDN
//
// Responses:
//
//	200: CDN
func PutCDN(entity CDNPutParam) (CDN, Alerts) {
	return CDN{}, Alerts{}
}

// GetCDNById swagger:route GET /cdns/{id} CDN GetCDNById
//
// # Retrieve a specific CDN by Id
//
// # Retrieve a specific CDN
//
// Responses:
//
//	200: CDNs
//	400: Alerts
func GetCDNById() (CDNs, Alerts) {
	return CDNs{}, Alerts{}
}

// DeleteCDN swagger:route DELETE /cdns/{id} CDN DeleteCDN
//
// # Delete a CDN by Id
//
// # Delete a CDN
//
// Responses:
//
//	200: Alerts
func DeleteCDN(entityId int) Alerts {
	return Alerts{}
}
