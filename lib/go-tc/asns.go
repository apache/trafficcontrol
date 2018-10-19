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

// Autonomous System Numbers
//
// A List of ASNs Response
// swagger:response ASNsResponse
// in: body
type ASNsResponse struct {
	// in: body
	Response []ASN `json:"response"`
}

// A Single ASN Response for Update and Create to depict what changed
// swagger:response ASNResponse
// in: body
type ASNResponse struct {
	// in: body
	Response ASN `json:"response"`
}

type ASN struct {
	// The ASN to retrieve
	//
	// Autonomous System Numbers per APNIC for identifying a service provider
	// required: true
	ASN int `json:"asn" db:"asn"`

	// Related cachegroup name
	//
	Cachegroup string `json:"cachegroup" db:"cachegroup"`

	// Related cachegroup id
	//
	CachegroupID int `json:"cachegroupId" db:"cachegroup_id"`

	// ID of the ASN
	//
	// required: true
	ID int `json:"id" db:"id"`

	// LastUpdated
	//
	LastUpdated string `json:"lastUpdated" db:"last_updated"`
}

type ASNNullable struct {
	// The ASN to retrieve
	//
	// Autonomous System Numbers per APNIC for identifying a service provider
	// required: true
	ASN *int `json:"asn" db:"asn"`

	// Related cachegroup name
	//
	Cachegroup *string `json:"cachegroup" db:"cachegroup"`

	// Related cachegroup id
	//
	CachegroupID *int `json:"cachegroupId" db:"cachegroup_id"`

	// ID of the ASN
	//
	// required: true
	ID *int `json:"id" db:"id"`

	// LastUpdated
	//
	LastUpdated *TimeNoMod `json:"lastUpdated" db:"last_updated"`
}

type ASNsV11 struct {
	ASNs []interface{} `json:"asns"`
}
