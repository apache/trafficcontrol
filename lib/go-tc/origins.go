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

import "time"

// OriginsResponse is a list of Origins as a response.
type OriginsResponse struct {
	Response []Origin `json:"response"`
	Alerts
}

// OriginDetailResponse is the JSON object returned for a single origin.
type OriginDetailResponse struct {
	Response Origin `json:"response"`
	Alerts
}

// Origin contains information relating to an Origin, which is NOT, in general,
// the same as an origin *server*.
type Origin struct {
	Cachegroup        *string    `json:"cachegroup" db:"cachegroup"`
	CachegroupID      *int       `json:"cachegroupId" db:"cachegroup_id"`
	Coordinate        *string    `json:"coordinate" db:"coordinate"`
	CoordinateID      *int       `json:"coordinateId" db:"coordinate_id"`
	DeliveryService   *string    `json:"deliveryService" db:"deliveryservice"`
	DeliveryServiceID *int       `json:"deliveryServiceId" db:"deliveryservice_id"`
	FQDN              *string    `json:"fqdn" db:"fqdn"`
	ID                *int       `json:"id" db:"id"`
	IP6Address        *string    `json:"ip6Address" db:"ip6_address"`
	IPAddress         *string    `json:"ipAddress" db:"ip_address"`
	IsPrimary         *bool      `json:"isPrimary" db:"is_primary"`
	LastUpdated       *TimeNoMod `json:"lastUpdated" db:"last_updated"`
	Name              *string    `json:"name" db:"name"`
	Port              *int       `json:"port" db:"port"`
	Profile           *string    `json:"profile" db:"profile"`
	ProfileID         *int       `json:"profileId" db:"profile_id"`
	Protocol          *string    `json:"protocol" db:"protocol"`
	Tenant            *string    `json:"tenant" db:"tenant"`
	TenantID          *int       `json:"tenantId" db:"tenant_id"`
}

// OriginsResponseV5 is an alias for the latest minor version of the major version 5.
type OriginsResponseV5 = OriginsResponseV50

// OriginsResponseV50 is a list of Origins as a response for APIv5.
type OriginsResponseV50 struct {
	Response []OriginV5 `json:"response"`
	Alerts
}

// OriginDetailResponseV5 is an alias for the latest minor version of the major version 5.
type OriginDetailResponseV5 = OriginDetailResponseV50

// OriginDetailResponseV50 is the JSON object returned for a single origin in APIv5.
type OriginDetailResponseV50 struct {
	Response OriginV5 `json:"response"`
	Alerts
}

// OriginV5 is an alias for the latest minor version of the major version 5.
type OriginV5 = OriginV50

// OriginV50 contains information relating to an Origin, in the latest minor version APIv50.
type OriginV50 struct {
	Cachegroup        *string   `json:"cachegroup" db:"cachegroup"`
	CachegroupID      *int      `json:"cachegroupId" db:"cachegroup_id"`
	Coordinate        *string   `json:"coordinate" db:"coordinate"`
	CoordinateID      *int      `json:"coordinateId" db:"coordinate_id"`
	DeliveryService   string    `json:"deliveryService" db:"deliveryservice"`
	DeliveryServiceID int       `json:"deliveryServiceId" db:"deliveryservice_id"`
	FQDN              string    `json:"fqdn" db:"fqdn"`
	ID                int       `json:"id" db:"id"`
	IP6Address        *string   `json:"ip6Address" db:"ip6_address"`
	IPAddress         *string   `json:"ipAddress" db:"ip_address"`
	IsPrimary         bool      `json:"isPrimary" db:"is_primary"`
	LastUpdated       time.Time `json:"lastUpdated" db:"last_updated"`
	Name              string    `json:"name" db:"name"`
	Port              *int      `json:"port" db:"port"`
	Profile           *string   `json:"profile" db:"profile"`
	ProfileID         *int      `json:"profileId" db:"profile_id"`
	Protocol          string    `json:"protocol" db:"protocol"`
	Tenant            string    `json:"tenant" db:"tenant"`
	TenantID          int       `json:"tenantId" db:"tenant_id"`
}

// ToOriginV5 upgrades from Origin to APIv5.
func (old Origin) ToOriginV5() OriginV5 {
	r := time.Unix(old.LastUpdated.Unix(), 0)

	var originV5 OriginV5
	originV5.Cachegroup = old.Cachegroup
	originV5.CachegroupID = old.CachegroupID
	originV5.Coordinate = old.Coordinate
	originV5.CoordinateID = old.CoordinateID
	originV5.DeliveryService = *old.DeliveryService
	originV5.DeliveryServiceID = *old.DeliveryServiceID
	originV5.FQDN = *old.FQDN
	originV5.ID = *old.ID
	originV5.IP6Address = old.IP6Address
	originV5.IPAddress = old.IPAddress
	originV5.IsPrimary = *old.IsPrimary
	originV5.LastUpdated = r
	originV5.Name = *old.Name
	originV5.Port = old.Port
	originV5.Profile = old.Profile
	originV5.ProfileID = old.ProfileID
	originV5.Protocol = *old.Protocol
	originV5.Tenant = *old.Tenant
	originV5.TenantID = *old.TenantID

	return originV5
}
