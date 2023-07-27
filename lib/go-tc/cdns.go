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

import (
	"database/sql"
	"time"
)

// CDNsV5Response is a list of CDNs as a response.
// swagger:response CDNsResponse
// in: body
type CDNsV5Response struct {
	// in: body
	Response []CDNV5 `json:"response"`
	Alerts
}

// CDNsResponse is a list of CDNs as a response.
// swagger:response CDNsResponse
// in: body
type CDNsResponse struct {
	// in: body
	Response []CDN `json:"response"`
	Alerts
}

// CDNResponse is a single CDN response for Update and Create to depict what
// changed.
// swagger:response CDNResponse
// in: body
type CDNResponse struct {
	// in: body
	Response CDN `json:"response"`
	Alerts
}

// A CDNV5 represents a set of configuration and hardware that can be used to
// serve content within a specific top-level domain.
type CDNV5 struct {

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
	LastUpdated time.Time `json:"lastUpdated" db:"last_updated"`

	// Name of the CDN
	//
	// required: true
	Name string `json:"name" db:"name"`

	// TTLOverride
	//
	TTLOverride *int `json:"ttlOverride,omitempty" db:"ttl_override"`
}

// A CDN represents a set of configuration and hardware that can be used to
// serve content within a specific top-level domain.
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

	// TTLOverride
	//
	TTLOverride int `json:"ttlOverride,omitempty" db:"ttl_override"`
}

// CDNNullable is identical to CDN except that its fields are reference values,
// which allows them to be nil.
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

	// TTLOverride
	//
	TTLOverride *int `json:"ttlOverride,omitempty" db:"ttl_override"`
}

// CDNSSLKeysResponse is the structure of the Traffic Ops API's response to
// requests made to its /cdns/name/{{name}}/sslkeys endpoint.
type CDNSSLKeysResponse struct {
	Response []CDNSSLKeys `json:"response"`
	Alerts
}

// CDNSSLKeys is an SSL key/certificate pair for a certain Delivery Service.
type CDNSSLKeys struct {
	DeliveryService string                `json:"deliveryservice"`
	Certificate     CDNSSLKeysCertificate `json:"certificate"`
	Hostname        string                `json:"hostname"`
}

// CDNSSLKeysCertificate is an SSL key/certificate pair.
type CDNSSLKeysCertificate struct {
	Crt string `json:"crt"`
	Key string `json:"key"`
}

// CDNConfig includes the name and ID of a single CDN configuration.
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

// CDNQueueUpdateRequest encodes the request data for the POST
// cdns/{{ID}}/queue_update endpoint.
type CDNQueueUpdateRequest ServerQueueUpdateRequest

// CDNQueueUpdateResponse encodes the response data for the POST
// cdns/{{ID}}/queue_update endpoint.
type CDNQueueUpdateResponse struct {
	Action string `json:"action"`
	CDNID  int64  `json:"cdnId"`
}
