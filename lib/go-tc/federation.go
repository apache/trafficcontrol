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
)

type CDNFederationResponse struct {
	Response []CDNFederation `json:"response"`
}

type CreateCDNFederationResponse struct {
	Response CDNFederation `json:"response"`
	Alerts
}

type UpdateCDNFederationResponse struct {
	Response CDNFederation `json:"response"`
	Alerts
}

type DeleteCDNFederationResponse struct {
	Alerts
}

type CDNFederation struct {
	ID          *int       `json:"id" db:"id"`
	CName       *string    `json:"cname" db:"cname"`
	TTL         *int       `json:"ttl" db:"ttl"`
	Description *string    `json:"description" db:"description"`
	LastUpdated *TimeNoMod `json:"lastUpdated" db:"last_updated"`

	// omitempty only works with primitive types and pointers
	*DeliveryServiceIDs `json:"deliveryService,omitempty"`
}

type DeliveryServiceIDs struct {
	DsId  *int    `json:"id,omitempty" db:"ds_id"`
	XmlId *string `json:"xmlId,omitempty" db:"xml_id"`
}

type FederationNullable struct {
	Mappings        []FederationMapping `json:"mappings"`
	DeliveryService *string             `json:"deliveryService"`
}

type FederationMapping struct {
	CName string `json:"cname"`
	TTL   int    `json:"ttl"`
}

// AllFederation is the JSON object returned by /api/1.x/federations?all
type AllFederation struct {
	Mappings        []AllFederationMapping `json:"mappings"`
	DeliveryService DeliveryServiceName    `json:"deliveryService"`
}

func (a AllFederation) IsAllFederations() bool { return true }

// AllFederation is the JSON object returned by /api/1.x/federations?all&cdnName=my-cdn-name
type AllFederationCDN struct {
	CDNName *CDNName `json:"cdnName"`
}

func (a AllFederationCDN) IsAllFederations() bool { return true }

type AllFederationMapping struct {
	TTL      *int     `json:"ttl"`
	CName    *string  `json:"cname"`
	Resolve4 []string `json:"resolve4,omitempty"`
	Resolve6 []string `json:"resolve6,omitempty"`
}

// IAllFederation is an interface for the disparate objects returned by /api/1.x/federations?all.
// Adds additional safety, allowing functions to only return one of the valid object types for the endpoint.
type IAllFederation interface {
	IsAllFederations() bool
}

type FederationDSPost struct {
	DSIDs   []int `json:"dsIds"`
	Replace *bool `json:"replace"`
}

func (f *FederationDSPost) Validate(tx *sql.Tx) error {
	return nil
}
