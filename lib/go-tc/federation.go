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
	"fmt"
	"net"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-util"

	validation "github.com/go-ozzo/ozzo-validation"
)

// CDNFederationResponse represents a Traffic Ops API response to a request for one or more of a CDN's
// Federations.
type CDNFederationResponse struct {
	Response []CDNFederation `json:"response"`
	Alerts
}

// CreateCDNFederationResponse represents a Traffic Ops API response to a request to
// create a new Federation for a CDN.
type CreateCDNFederationResponse struct {
	Response CDNFederation `json:"response"`
	Alerts
}

// UpdateCDNFederationResponse represents a Traffic Ops API response to a request to replace a
// Federation of a CDN with the one provided in the request body.
type UpdateCDNFederationResponse struct {
	Response CDNFederation `json:"response"`
	Alerts
}

// DeleteCDNFederationResponse represents a Traffic Ops API response to a request to remove a
// Federation from a CDN.
type DeleteCDNFederationResponse struct {
	Alerts
}

// CDNFederation represents a Federation.
type CDNFederation struct {
	ID          *int       `json:"id" db:"id"`
	CName       *string    `json:"cname" db:"cname"`
	TTL         *int       `json:"ttl" db:"ttl"`
	Description *string    `json:"description" db:"description"`
	LastUpdated *TimeNoMod `json:"lastUpdated" db:"last_updated"`

	// omitempty only works with primitive types and pointers
	*DeliveryServiceIDs `json:"deliveryService,omitempty"`
}

// CDNFederationDeliveryService holds information about an assigned Delivery
// Service within a CDNFederationV5 structure.
type CDNFederationDeliveryService = struct {
	ID    int    `json:"id" db:"ds_id"`
	XMLID string `json:"xmlID" db:"xml_id"`
}

// CDNFederationV5 represents a Federation of some CDN as it appears in version
// 5 of the Traffic Ops API.
type CDNFederationV5 struct {
	ID          int       `json:"id" db:"id"`
	CName       string    `json:"cname" db:"cname"`
	TTL         int       `json:"ttl" db:"ttl"`
	Description *string   `json:"description" db:"description"`
	LastUpdated time.Time `json:"lastUpdated" db:"last_updated"`

	DeliveryService *CDNFederationDeliveryService `json:"deliveryService,omitempty"`
}

// CDNFederationsV5Response represents a Traffic Ops APIv5 response to a request
// for one or more of a CDN's Federations.
type CDNFederationsV5Response struct {
	Response []CDNFederationV5 `json:"response"`
	Alerts
}

// CDNFederationV5Response represents a Traffic Ops APIv5 response to a request
// for a single CDN's Federations.
type CDNFederationV5Response struct {
	Response CDNFederationV5 `json:"response"`
	Alerts
}

// DeliveryServiceIDs are pairs of identifiers for Delivery Services.
type DeliveryServiceIDs struct {
	DsId  *int    `json:"id,omitempty" db:"ds_id"`
	XmlId *string `json:"xmlId,omitempty" db:"xml_id"`
}

// FederationNullable represents a relationship between some Federation
// Mappings and a Delivery Service.
//
// This is not known to be used anywhere.
type FederationNullable struct {
	Mappings        []FederationMapping `json:"mappings"`
	DeliveryService *string             `json:"deliveryService"`
}

// A FederationMapping is a Federation, without any information about its
// resolver mappings or any relation to Delivery Services.
type FederationMapping struct {
	CName string `json:"cname"`
	TTL   int    `json:"ttl"`
}

// AllDeliveryServiceFederationsMapping is a structure that contains identifying information for a
// Delivery Service as well as any and all Federation Resolver mapping assigned to it (or all those
// getting assigned to it).
type AllDeliveryServiceFederationsMapping struct {
	Mappings        []FederationResolverMapping `json:"mappings"`
	DeliveryService DeliveryServiceName         `json:"deliveryService"`
}

// IsAllFederations implements the IAllFederation interface. Always returns true.
func (a AllDeliveryServiceFederationsMapping) IsAllFederations() bool { return true }

// FederationsResponse is the type of a response from Traffic Ops to
// requests to its /federations/all and /federations endpoints.
type FederationsResponse struct {
	Response []AllDeliveryServiceFederationsMapping `json:"response"`
	Alerts
}

// AllFederationCDN is the structure of a response from Traffic Ops to a GET
// request made to its /federations/all endpoint.
type AllFederationCDN struct {
	CDNName *CDNName `json:"cdnName"`
}

// IsAllFederations implements the IAllFederation interface. Always returns true.
func (a AllFederationCDN) IsAllFederations() bool { return true }

// A ResolverMapping is a set of Resolvers, ostensibly for a Federation.
type ResolverMapping struct {
	Resolve4 []string `json:"resolve4,omitempty"`
	Resolve6 []string `json:"resolve6,omitempty"`
}

// Validate implements the
// github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api.ParseValidator
// interface.
func (r *ResolverMapping) Validate(tx *sql.Tx) error {
	errs := []error{}
	for _, res := range r.Resolve4 {
		ip := net.ParseIP(res)
		if ip != nil {
			if ip.To4() == nil {
				errs = append(errs, fmt.Errorf("[ %s ] is not a valid ip address.", res))
			}
			continue
		}

		if ip, _, err := net.ParseCIDR(res); err != nil || ip.To4() == nil {
			errs = append(errs, fmt.Errorf("[ %s ] is not a valid ip address.", res))
		}
	}

	for _, res := range r.Resolve6 {
		ip := net.ParseIP(res)
		if ip != nil {
			if ip.To16() == nil {
				errs = append(errs, fmt.Errorf("[ %s ] is not a valid ip address.", res))
			}
			continue
		}

		if ip, _, err := net.ParseCIDR(res); err != nil || ip.To16() == nil {
			errs = append(errs, fmt.Errorf("[ %s ] is not a valid ip address.", res))
		}
	}
	if len(errs) > 0 {
		return util.JoinErrs(errs)
	}
	return nil
}

// FederationResolverMapping is the set of all resolvers - both IPv4 and IPv6 - for a specific
// Federation.
type FederationResolverMapping struct {
	// TTL is the Time-to-Live of a DNS response to a request to resolve this Federation's CNAME
	TTL   *int    `json:"ttl"`
	CName *string `json:"cname"`
	ResolverMapping
}

// IAllFederation is an interface for the disparate objects returned by
// Federations-related Traffic Ops API endpoints.
// Adds additional safety, allowing functions to only return one of the valid
// object types for the endpoint.
type IAllFederation interface {
	IsAllFederations() bool
}

// FederationDSPost is the format of a request to associate a Federation with any number of Delivery Services.
type FederationDSPost struct {
	DSIDs []int `json:"dsIds"`
	// Replace indicates whether existing Federation-to-Delivery Service associations should be
	// replaced by the ones defined by this request, or otherwise merely augmented with them.
	Replace *bool `json:"replace"`
}

// Validate implements the
// github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api.ParseValidator
// interface.
func (f *FederationDSPost) Validate(tx *sql.Tx) error {
	return nil
}

// FederationUser represents Federation Users.
type FederationUser struct {
	Company  *string `json:"company" db:"company"`
	Email    *string `json:"email" db:"email"`
	FullName *string `json:"fullName" db:"full_name"`
	ID       *int    `json:"id" db:"id"`
	Role     *string `json:"role" db:"role_name"`
	Username *string `json:"username" db:"username"`
}

// FederationUsersResponse is the type of a response from Traffic Ops to a
// request made to its /federations/{{ID}}/users endpoint.
type FederationUsersResponse struct {
	Response []FederationUser `json:"response"`
	Alerts
}

// FederationUserPost represents POST body for assigning Users to Federations.
type FederationUserPost struct {
	IDs     []int `json:"userIds"`
	Replace *bool `json:"replace"`
}

// Validate implements the
// github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api.ParseValidator
// interface.
func (f *FederationUserPost) Validate(tx *sql.Tx) error {
	return validation.ValidateStruct(f,
		validation.Field(&f.IDs, validation.NotNil),
	)
}

// DeliveryServiceFederationResolverMapping structures represent resolvers to
// which a Delivery Service maps "federated" CDN traffic.
type DeliveryServiceFederationResolverMapping struct {
	DeliveryService string          `json:"deliveryService"`
	Mappings        ResolverMapping `json:"mappings"`
}

// Validate implements the
// github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api.ParseValidator
// interface.
func (d *DeliveryServiceFederationResolverMapping) Validate(tx *sql.Tx) error {
	return d.Mappings.Validate(tx)
}

// LegacyDeliveryServiceFederationResolverMappingRequest is the legacy format for a request to
// create (or modify) the Federation Resolver mappings of one or more Delivery Services. Use this
// for compatibility with API versions 1.3 and older.
type LegacyDeliveryServiceFederationResolverMappingRequest struct {
	Federations []DeliveryServiceFederationResolverMapping `json:"federations"`
}

// Validate implements the github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api.ParseValidator
// interface.
func (req *LegacyDeliveryServiceFederationResolverMappingRequest) Validate(tx *sql.Tx) error {
	if len(req.Federations) < 1 {
		return errors.New("federations: required")
	}
	r := DeliveryServiceFederationResolverMappingRequest(req.Federations)
	return (&r).Validate(tx)
}

// DeliveryServiceFederationResolverMappingRequest is the format of a request to create (or
// modify) the Federation Resolver mappings of one or more Delivery Services. Use this when working
// only with API versions 1.4 and newer.
type DeliveryServiceFederationResolverMappingRequest []DeliveryServiceFederationResolverMapping

// Validate implements the github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api.ParseValidator
// interface.
func (req *DeliveryServiceFederationResolverMappingRequest) Validate(tx *sql.Tx) error {
	if len(*req) < 1 {
		return errors.New("must specify at least one Delivery Service/Federation Resolver(s) mapping")
	}
	errs := []error{}
	for _, m := range *req {
		if err := m.Validate(tx); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return util.JoinErrs(errs)
	}
	return nil
}
