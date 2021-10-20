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

// Validate implements the github.com/apache/trafficcontrol/v6/traffic_ops/traffic_ops_golang/api.ParseValidator
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
