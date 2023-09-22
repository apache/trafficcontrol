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
	"errors"
	"net"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

// FederationResolversResponse represents a Traffic Ops API response to a
// GET request to its /federation_resolvers endpoint.
type FederationResolversResponse struct {
	Alerts
	Response []FederationResolver `json:"response"`
}

// FederationResolverResponse represents a Traffic Ops API response to a
// POST or DELETE request to its /federation_resolvers endpoint.
type FederationResolverResponse struct {
	Alerts
	Response FederationResolver `json:"response"`
}

// FederationResolver represents a resolver record for a CDN Federation.
type FederationResolver struct {
	ID          *uint      `json:"id" db:"id"`
	IPAddress   *string    `json:"ipAddress" db:"ip_address"`
	LastUpdated *TimeNoMod `json:"lastUpdated,omitempty" db:"last_updated"`
	Type        *string    `json:"type"`
	TypeID      *uint      `json:"typeId,omitempty" db:"type"`
}

// FederationResolverV5 - is an alias for the Federal Resolver struct response used for the latest minor version associated with APIv5.
type FederationResolverV5 = FederationResolverV50

// FederationResolverV50 - is used for RFC3339 format timestamp in FederationResolver which represents a resolver record for a CDN Federation for APIv50.
type FederationResolverV50 struct {
	ID          *uint      `json:"id" db:"id"`
	IPAddress   *string    `json:"ipAddress" db:"ip_address"`
	LastUpdated *time.Time `json:"lastUpdated,omitempty" db:"last_updated"`
	Type        *string    `json:"type"`
	TypeID      *uint      `json:"typeId,omitempty" db:"type"`
}

// FederationResolversResponseV5 - an alias for the Federation Resolver's struct response used for the latest minor version associated with APIv5.
type FederationResolversResponseV5 = FederationResolversResponseV50

// FederationResolversResponseV50 - GET request to its /federation_resolvers endpoint for APIv50.
type FederationResolversResponseV50 struct {
	Alerts
	Response []FederationResolverV5 `json:"response"`
}

// FederationResolverResponseV5 - represents struct response used for the latest minor version associated with APIv5.
type FederationResolverResponseV5 = FederationResolverResponseV50

// FederationResolverResponseV50 - POST request to its /federation_resolvers endpoint APIv50.
type FederationResolverResponseV50 struct {
	Alerts
	Response FederationResolverV5 `json:"response"`
}

// Validate implements the github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api.ParseValidator
// interface.
func (fr *FederationResolver) Validate(tx *sql.Tx) error {
	return validation.ValidateStruct(fr,
		validation.Field(&fr.IPAddress, validation.Required, validation.By(func(v interface{}) error {
			if v == nil {
				return nil // this is handled by 'required'
			}

			if ip := net.ParseIP(*v.(*string)); ip != nil {
				return nil
			}

			if _, _, err := net.ParseCIDR(*v.(*string)); err != nil {
				return errors.New("invalid network IP or CIDR-notation subnet")
			}
			return nil
		})),
		validation.Field(&fr.TypeID, validation.Required),
	)
}

// UpgradeToFederationResolverV5 upgrades an APIv4 Federal Resolver into an APIv5 Federal Resolver of
// the latest minor version.
func UpgradeToFederationResolverV5(fr FederationResolver) *FederationResolverV5 {
	upgraded := FederationResolverV5{
		ID:        fr.ID,
		IPAddress: fr.IPAddress,
		LastUpdated: func() *time.Time {
			lastUpdated := time.Unix(fr.LastUpdated.Unix(), 0)
			return &lastUpdated
		}(),
		Type: fr.Type,
	}

	return &upgraded
}
