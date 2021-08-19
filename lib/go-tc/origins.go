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
