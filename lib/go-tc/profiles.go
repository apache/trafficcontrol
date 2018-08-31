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

// ProfilesResponse ...
type ProfilesResponse struct {
	Response []Profile `json:"response"`
}

// A Single Profile Response for Update and Create to depict what changed
// swagger:response ProfileResponse
// in: body
type ProfileResponse struct {
	// in: body
	Response Profile `json:"response"`
}

// Profile ...
type Profile struct {
	ID              int                 `json:"id" db:"id"`
	LastUpdated     TimeNoMod           `json:"lastUpdated"`
	Name            string              `json:"name"`
	Parameter       string              `json:"param"`
	Description     string              `json:"description"`
	CDNName         string              `json:"cdnName"`
	CDNID           int                 `json:"cdn"`
	RoutingDisabled bool                `json:"routingDisabled"`
	Type            string              `json:"type"`
	Parameters      []ParameterNullable `json:"params,omitempty"`
}

type ProfileNullable struct {

	// Unique identifier for the Profile
	//
	ID *int `json:"id" db:"id"`

	// LastUpdated
	//
	LastUpdated *TimeNoMod `json:"lastUpdated" db:"last_updated"`

	// The Profile name
	//
	Name *string `json:"name" db:"name"`

	// The Profile Description
	//
	Description *string `json:"description" db:"description"`

	// The CDN name associated with the Profile
	//
	CDNName *string `json:"cdnName" db:"cdn_name"`

	// The CDN id associated with the Profile
	//
	CDNID *int `json:"cdn" db:"cdn"`

	// Enables
	//
	RoutingDisabled *bool `json:"routingDisabled" db:"routing_disabled"`

	// The Type name associated with the Profile
	//
	Type *string `json:"type" db:"type"`

	// Parameters associated to the profile
	//
	Parameters []ParameterNullable `json:"params,omitempty"`
}
type ProfileTrimmed struct {
	Name string `json:"name"`
}
