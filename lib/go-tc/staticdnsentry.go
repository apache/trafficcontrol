package tc

import "time"

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

// StaticDNSEntriesResponse is a list of StaticDNSEntry as a response.
type StaticDNSEntriesResponse struct {
	Response []StaticDNSEntry `json:"response"`
	Alerts
}

// StaticDNSEntriesResponseV5 is a list of StaticDNSEntry as a response, for api version 5.0., for the latest
// minor version of 5.x.
type StaticDNSEntriesResponseV5 = StaticDNSEntriesResponseV50

// StaticDNSEntriesResponseV50 is a list of StaticDNSEntry as a response, for api version 5.0.
type StaticDNSEntriesResponseV50 struct {
	Response []StaticDNSEntryV5 `json:"response"`
	Alerts
}

// StaticDNSEntryV5 holds information about a static DNS entry, for the latest minor version of 5.x.
type StaticDNSEntryV5 = StaticDNSEntryV50

// StaticDNSEntryV50 holds information about a static DNS entry, for api version 5.0.
type StaticDNSEntryV50 struct {

	// The static IP Address or fqdn of the static dns entry
	//
	// required: true
	Address *string `json:"address" db:"address"`

	// The Cachegroup Name associated
	//
	CacheGroupName *string `json:"cachegroup" db:"cachegroup"`

	// The Cachegroup ID associated
	//
	CacheGroupID *int `json:"cachegroupId" db:"cachegroup_id"`

	// The DeliveryService associated
	//
	DeliveryService *string `json:"deliveryservice" db:"dsname"`

	// The DeliveryService associated
	//
	// required: true
	DeliveryServiceID *int `json:"deliveryserviceId" db:"deliveryservice_id"`

	// The host of the static dns entry
	//
	// required: true
	Host *string `json:"host" db:"host"`

	// ID of the StaticDNSEntry
	//
	// required: true
	ID *int `json:"id" db:"id"`

	// LastUpdated
	//
	LastUpdated *time.Time `json:"lastUpdated" db:"last_updated"`

	// The Time To Live for the static dns entry
	//
	// required: true
	TTL *int64 `json:"ttl" db:"ttl"`

	// The type of the static DNS entry
	//
	// enum: ["A_RECORD", "AAAA_RECORD", "CNAME_RECORD"]
	Type *string `json:"type"`

	// The type id of the static DNS entry
	//
	// required: true
	TypeID *int `json:"typeId" db:"type_id"`
}

// StaticDNSEntry holds information about a static DNS entry.
type StaticDNSEntry struct {

	// The static IP Address or fqdn of the static dns entry
	//
	// required: true
	Address string `json:"address" db:"address"`

	// The Cachegroup Name associated
	//
	CacheGroupName string `json:"cachegroup" db:"cachegroup"`

	// The Cachegroup ID associated
	//
	CacheGroupID int `json:"cachegroupId" db:"cachegroup_id"`

	// The DeliveryService associated
	//
	DeliveryService string `json:"deliveryservice" db:"dsname"`

	// The DeliveryService associated
	//
	// required: true
	DeliveryServiceID int `json:"deliveryserviceId" db:"deliveryservice_id"`

	// The host of the static dns entry
	//
	// required: true
	Host string `json:"host" db:"host"`

	// ID of the StaticDNSEntry
	//
	// required: true
	ID int `json:"id" db:"id"`

	// LastUpdated
	//
	LastUpdated TimeNoMod `json:"lastUpdated" db:"last_updated"`

	// The Time To Live for the static dns entry
	//
	// required: true
	TTL int64 `json:"ttl" db:"ttl"`

	// The type of the static DNS entry
	//
	// enum: ["A_RECORD", "AAAA_RECORD", "CNAME_RECORD"]
	Type string `json:"type"`

	// The type id of the static DNS entry
	//
	// required: true
	TypeID int `json:"typeId" db:"type_id"`
}

// StaticDNSEntryNullable holds information about a static DNS entry. Its fields
// are nullable.
type StaticDNSEntryNullable struct {

	// The static IP Address or fqdn of the static dns entry
	//
	// required: true
	Address *string `json:"address" db:"address"`

	// The Cachegroup Name associated
	//
	CacheGroupName *string `json:"cachegroup" db:"cachegroup"`

	// The Cachegroup ID associated
	//
	CacheGroupID *int `json:"cachegroupId" db:"cachegroup_id"`

	// The DeliveryService Name associated
	//
	DeliveryService *string `json:"deliveryservice" db:"dsname"`

	// DeliveryService ID of the StaticDNSEntry
	//
	// required: true
	DeliveryServiceID *int `json:"deliveryserviceId" db:"deliveryservice_id"`

	// The host of the static dns entry
	//
	// required: true
	Host *string `json:"host" db:"host"`

	// ID of the StaticDNSEntry
	//
	// required: true
	ID *int `json:"id" db:"id"`

	// LastUpdated
	//
	LastUpdated *TimeNoMod `json:"lastUpdated" db:"last_updated"`

	// The Time To Live for the static dns entry
	//
	// required: true
	TTL *int64 `json:"ttl" db:"ttl"`

	// The type of the static DNS entry
	//
	// enum: ["A_RECORD", "AAAA_RECORD", "CNAME_RECORD"]
	Type *string `json:"type"`

	// The type id of the static DNS entry
	//
	// required: true
	TypeID int `json:"typeId" db:"type_id"`
}
