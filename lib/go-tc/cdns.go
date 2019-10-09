package tc

import (
	"database/sql"
)

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

// A List of CDNs Response
// swagger:response CDNsResponse
// in: body
type CDNsResponse struct {
	// in: body
	Response []CDN `json:"response"`
}

// A Single CDN Response for Update and Create to depict what changed
// swagger:response CDNResponse
// in: body
type CDNResponse struct {
	// in: body
	Response CDN `json:"response"`
}

// CDN ...
type CDN struct {

	// The CDN to retrieve
	//
	// enables Domain Name Security Extensions on the specified CDN
	//
	// required: true
	DNSSECEnabled bool `json:"dnssecEnabled" db:"dnssec_enabled"`

	// DomainName of the CDN
	//
	// required: true
	DomainName string `json:"domainName" db:"domain_name"`

	// ID of the CDN
	//
	// required: true
	ID int `json:"id" db:"id"`

	// LastUpdated
	//
	LastUpdated TimeNoMod `json:"lastUpdated" db:"last_updated"`

	// Name of the CDN
	//
	// required: true
	Name string `json:"name" db:"name"`
}

// CDNNullable ...
type CDNNullable struct {

	// The CDN to retrieve
	//
	// enables Domain Name Security Extensions on the specified CDN
	//
	// required: true
	DNSSECEnabled *bool `json:"dnssecEnabled" db:"dnssec_enabled"`

	// DomainName of the CDN
	//
	// required: true
	DomainName *string `json:"domainName" db:"domain_name"`

	// ID of the CDN
	//
	// required: true
	ID *int `json:"id" db:"id"`

	// LastUpdated
	//
	LastUpdated *TimeNoMod `json:"lastUpdated" db:"last_updated"`

	// Name of the CDN
	//
	// required: true
	Name *string `json:"name" db:"name"`
}

// CDNSSLKeysResponse ...
type CDNSSLKeysResponse struct {
	Response []CDNSSLKeys `json:"response"`
}

// CDNSSLKeys ...
type CDNSSLKeys struct {
	DeliveryService string                `json:"deliveryservice"`
	Certificate     CDNSSLKeysCertificate `json:"certificate"`
	Hostname        string                `json:"hostname"`
}

// CDNSSLKeysCertificate ...
type CDNSSLKeysCertificate struct {
	Crt string `json:"crt"`
	Key string `json:"key"`
}

type CDNConfig struct {
	Name *string `json:"name"`
	ID   *int    `json:"id"`
}

// CDNExistsByName returns whether a cdn with the given name exists, and any error.
// TODO move to helper package.
func CDNExistsByName(name string, tx *sql.Tx) (bool, error) {
	exists := false
	err := tx.QueryRow(`SELECT EXISTS(SELECT * FROM cdn WHERE name = $1)`, name).Scan(&exists)
	return exists, err
}
