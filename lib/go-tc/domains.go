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

// DomainsResponse is a list of Domains as a response.
type DomainsResponse struct {
	Response []Domain `json:"response"`
	Alerts
}

// Domain contains information about a single Profile within a Domain.
type Domain struct {
	ProfileID int `json:"profileId" db:"profile_id"`

	// This property serves no known purpose; it is always -1.
	ParameterID int `json:"parameterId" db:"parameter_id"`

	ProfileName string `json:"profileName" db:"profile_name"`

	ProfileDescription string `json:"profileDescription" db:"profile_description"`

	// DomainName of the CDN
	DomainName string `json:"domainName" db:"domain_name"`
}

// DomainNullable is identical to a Domain but with reference properties that
// can have nil values, mostly used by the Traffic Ops API.
type DomainNullable struct {
	ProfileID          *int    `json:"profileId" db:"profile_id"`
	ParameterID        *int    `json:"parameterId" db:"parameter_id"`
	ProfileName        *string `json:"profileName" db:"profile_name"`
	ProfileDescription *string `json:"profileDescription" db:"profile_description"`
	DomainName         *string `json:"domainName" db:"domain_name"`
}
