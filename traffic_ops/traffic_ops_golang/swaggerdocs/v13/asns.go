package v13

import "github.com/apache/trafficcontrol/v8/lib/go-tc"

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

// ASNs -  ASNsResponse to get the "response" top level key
// swagger:response ASNs
// in: body
type ASNs struct {
	// ASN Response Body
	// in: body
	ASNsResponse tc.ASNsResponse `json:"response"`
}

// ASN -  ASNResponse to get the "response" top level key
// swagger:response ASN
// in: body
type ASN struct {
	// ASN Response Body
	// in: body
	ASNResponse tc.ASNResponse
}

// ASNQueryParams
//
// swagger:parameters GetASNs
type ASNQueryParams struct {

	// ASNsQueryParams

	// Autonomous System Numbers per APNIC for identifying a service provider
	//
	Asn string `json:"asn"`

	// Related cachegroup name
	//
	Cachegroup string `json:"cachegroup"`

	// Related cachegroup id
	//
	CachegroupID string `json:"cachegroupId"`

	// Unique identifier for the CDN
	//
	ID string `json:"id"`

	//
	//
	Orderby string `json:"orderby"`
}

// swagger:parameters PostASN
type ASNPostParam struct {
	// ASN Request Body
	//
	// in: body
	// required: true
	ASN tc.ASN
}

// swagger:parameters GetASNById DeleteASN
type ASNPathParams struct {

	// Id associated to the ASN
	// in: path
	ID int `json:"id"`
}

// PostASN swagger:route POST /asns ASN PostASN
//
// # Create a ASN
//
// # An Autonomous System Number
//
// Responses:
//
//	200: Alerts
func PostASN(entity ASNPostParam) (ASN, Alerts) {
	return ASN{}, Alerts{}
}

// GetASNs swagger:route GET /asns ASN GetASNs
//
// # Retrieve a list of ASNs
//
// # A list of ASNs
//
// Responses:
//
//	200: ASNs
//	400: Alerts
func GetASNs() (ASNs, Alerts) {
	return ASNs{}, Alerts{}
}

// swagger:parameters PutASN
type ASNPutParam struct {

	// ID
	// in: path
	ID int `json:"id"`

	// ASN Request Body
	//
	// in: body
	// required: true
	ASN tc.ASN
}

// PutASN swagger:route PUT /asns/{id} ASN PutASN
//
// # Update an ASN by Id
//
// # Update an ASN
//
// Responses:
//
//	200: ASN
func PutASN(entity ASNPutParam) (ASN, Alerts) {
	return ASN{}, Alerts{}
}

// GetASNById swagger:route GET /asns/{id} ASN GetASNById
//
// # Retrieve a specific ASN by Id
//
// # Retrieve an ASN
//
// Responses:
//
//	200: ASNs
//	400: Alerts
func GetASNById() (ASNs, Alerts) {
	return ASNs{}, Alerts{}
}

// DeleteASN swagger:route DELETE /asns/{id} ASN DeleteASN
//
// # Delete an ASN by Id
//
// # Delete an ASN
//
// Responses:
//
//	200: Alerts
func DeleteASN(entityId int) Alerts {
	return Alerts{}
}
