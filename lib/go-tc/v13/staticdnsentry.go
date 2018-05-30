package v13

import tc "github.com/apache/trafficcontrol/lib/go-tc"

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

// StatisDNSEntry ...
type StaticDNSEntry struct {

	// The static IP Address or fqdn of the static dns entry
	//
	// required: true
	Address string `json:"address" db:"address"`

	// The Cachegroup associated
	//
	CacheGroup string `json:"cachegroup" db:"cachegroup"`

	// The DeliveryService associated
	//
	// required: true
	DeliveryService string `json:"deliveryservice" db:"dsname"`

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
	LastUpdated tc.TimeNoMod `json:"lastUpdated" db:"last_updated"`

	// The Time To Live for the static dns entry
	//
	// required: true
	TTL int64 `json:"ttl" db:"ttl"`

	// The type of the static DNS entry
	//
	// required: true
	// enum: ["A_RECORD", "AAAA_RECORD", "CNAME_RECORD"]
	Type string `json:"type" db:"type"`
}

// StatisDNSEntryNullable ...
type StaticDNSEntryNullable struct {

	// The static IP Address or fqdn of the static dns entry
	//
	// required: true
	Address *string `json:"address" db:"address"`

	// The Cachegroup associated
	//
	CacheGroup *string `json:"cachegroup" db:"cachegroup"`

	// The DeliveryService associated
	//
	// required: true
	DeliveryService *string `json:"deliveryservice" db:"dsname"`

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
	LastUpdated *tc.TimeNoMod `json:"lastUpdated" db:"last_updated"`

	// The Time To Live for the static dns entry
	//
	// required: true
	TTL *int64 `json:"ttl" db:"ttl"`

	// The type of the static DNS entry
	//
	// required: true
	// enum: ["A_RECORD", "AAAA_RECORD", "CNAME_RECORD"]
	Type *string `json:"type" db:"type"`
}
