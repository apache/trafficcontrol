package v13

import (
	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
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

type OriginsResponse struct {
	Response []OriginNullable `json:"response"`
}

// OriginDetailResponse is the JSON object returned for a single origin
type OriginDetailResponse struct {
	Response OriginNullable `json:"response"`
	tc.Alerts
}

type Origin struct {
	Cachegroup        string       `json:"cachegroup" db:"cachegroup"`
	CachegroupID      int          `json:"cachegroupId" db:"cachegroup_id"`
	Coordinate        string       `json:"coordinate" db:"coordinate"`
	CoordinateID      int          `json:"coordinateId" db:"coordinate_id"`
	DeliveryService   string       `json:"deliveryService" db:"deliveryservice"`
	DeliveryServiceID int          `json:"deliveryServiceId" db:"deliveryservice_id"`
	FQDN              string       `json:"fqdn" db:"fqdn"`
	ID                int          `json:"id" db:"id"`
	IP6Address        string       `json:"ip6Address" db:"ip6_address"`
	IPAddress         string       `json:"ipAddress" db:"ip_address"`
	IsPrimary         bool         `json:"isPrimary" db:"is_primary"`
	LastUpdated       tc.TimeNoMod `json:"lastUpdated" db:"last_updated"`
	Name              string       `json:"name" db:"name"`
	Port              int          `json:"port" db:"port"`
	Profile           string       `json:"profile" db:"profile"`
	ProfileID         int          `json:"profileId" db:"profile_id"`
	Protocol          string       `json:"protocol" db:"protocol"`
	Tenant            string       `json:"tenant" db:"tenant"`
	TenantID          int          `json:"tenantId" db:"tenant_id"`
}

type OriginNullable struct {
	Cachegroup        *string       `json:"cachegroup" db:"cachegroup"`
	CachegroupID      *int          `json:"cachegroupId" db:"cachegroup_id"`
	Coordinate        *string       `json:"coordinate" db:"coordinate"`
	CoordinateID      *int          `json:"coordinateId" db:"coordinate_id"`
	DeliveryService   *string       `json:"deliveryService" db:"deliveryservice"`
	DeliveryServiceID *int          `json:"deliveryServiceId" db:"deliveryservice_id"`
	FQDN              *string       `json:"fqdn" db:"fqdn"`
	ID                *int          `json:"id" db:"id"`
	IP6Address        *string       `json:"ip6Address" db:"ip6_address"`
	IPAddress         *string       `json:"ipAddress" db:"ip_address"`
	IsPrimary         *bool         `json:"isPrimary" db:"is_primary"`
	LastUpdated       *tc.TimeNoMod `json:"lastUpdated" db:"last_updated"`
	Name              *string       `json:"name" db:"name"`
	Port              *int          `json:"port" db:"port"`
	Profile           *string       `json:"profile" db:"profile"`
	ProfileID         *int          `json:"profileId" db:"profile_id"`
	Protocol          *string       `json:"protocol" db:"protocol"`
	Tenant            *string       `json:"tenant" db:"tenant"`
	TenantID          *int          `json:"tenantId" db:"tenant_id"`
}

type OriginBuilder struct {
	originBuild OriginNullable
}

func NewOriginBuilder() *OriginBuilder {
	return &OriginBuilder{}
}

func (o *OriginBuilder) Build() OriginNullable {
	return o.originBuild
}

func (o *OriginBuilder) CachegroupID(id int) *OriginBuilder {
	o.originBuild.CachegroupID = &id
	return o
}

func (o *OriginBuilder) CoordinateID(id int) *OriginBuilder {
	o.originBuild.CoordinateID = &id
	return o
}

func (o *OriginBuilder) DeliveryServiceID(id int) *OriginBuilder {
	o.originBuild.DeliveryServiceID = &id
	return o
}

func (o *OriginBuilder) FQDN(fqdn string) *OriginBuilder {
	o.originBuild.FQDN = &fqdn
	return o
}

func (o *OriginBuilder) IP6Address(ip string) *OriginBuilder {
	o.originBuild.IP6Address = &ip
	return o
}

func (o *OriginBuilder) IPAddress(ip string) *OriginBuilder {
	o.originBuild.IPAddress = &ip
	return o
}

func (o *OriginBuilder) IsPrimary(isPrimary bool) *OriginBuilder {
	o.originBuild.IsPrimary = &isPrimary
	return o
}

func (o *OriginBuilder) Name(name string) *OriginBuilder {
	o.originBuild.Name = &name
	return o
}

func (o *OriginBuilder) Port(port int) *OriginBuilder {
	o.originBuild.Port = &port
	return o
}

func (o *OriginBuilder) ProfileID(id int) *OriginBuilder {
	o.originBuild.ProfileID = &id
	return o
}

func (o *OriginBuilder) Protocol(proto string) *OriginBuilder {
	o.originBuild.Protocol = &proto
	return o
}

func (o *OriginBuilder) TenantID(id int) *OriginBuilder {
	o.originBuild.TenantID = &id
	return o
}
