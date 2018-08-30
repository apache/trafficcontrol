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
	"time"
)

// UsersResponse ...
type UsersResponse struct {
	Response []User `json:"response"`
}

// User contains information about a given user in Traffic Ops.
type User struct {
	Username         string    `json:"username,omitempty"`
	PublicSSHKey     string    `json:"publicSshKey,omitempty"`
	Role             int       `json:"role,omitempty"`
	RoleName         string    `json:"rolename,omitempty"`
	ID               int       `json:"id,omitempty"`
	UID              int       `json:"uid,omitempty"`
	GID              int       `json:"gid,omitempty"`
	Company          string    `json:"company,omitempty"`
	Email            string    `json:"email,omitempty"`
	FullName         string    `json:"fullName,omitempty"`
	NewUser          bool      `json:"newUser,omitempty"`
	LastUpdated      string    `json:"lastUpdated,omitempty"`
	AddressLine1     string    `json:"addressLine1"`
	AddressLine2     string    `json:"addressLine2"`
	City             string    `json:"city"`
	Country          string    `json:"country"`
	PhoneNumber      string    `json:"phoneNumber"`
	PostalCode       string    `json:"postalCode"`
	RegistrationSent time.Time `json:"registrationSent"`
	StateOrProvince  string    `json:"stateOrProvince"`
	Tenant           string    `json:"tenant"`
	TenantID         int       `json:"tenantId"`
}

// Credentials contains Traffic Ops login credentials
type UserCredentials struct {
	Username string `json:"u"`
	Password string `json:"p"`
}

// TODO reconcile APIUser and User

type APIUser struct {
	AddressLine1     *string    `json:"addressLine1", db:"address_line1"`
	AddressLine2     *string    `json:"addressLine2" db:"address_line2"`
	City             *string    `json:"city" db:"city"`
	Company          *string    `json:"company,omitempty" db:"company"`
	Country          *string    `json:"country" db:"country"`
	Email            *string    `json:"email" db:"email"`
	FullName         *string    `json:"fullName" db:"full_name"`
	GID              *int       `json:"gid" db:"gid"`
	ID               *int       `json:"id" db:"id"`
	LastUpdated      *time.Time `json:"lastUpdated" db:"last_updated"`
	NewUser          *bool      `json:"newUser" db:"new_user"`
	PhoneNumber      *string    `json:"phoneNumber" db:"phone_number"`
	PostalCode       *string    `json:"postalCode" db:"postal_code"`
	PublicSSHKey     *string    `json:"publicSshKey" db:"public_ssh_key"`
	RegistrationSent *time.Time `json:"registrationSent" db:"registration_sent"`
	Role             *int       `json:"role" db:"role"`
	RoleName         *string    `json:"rolename"`
	StateOrProvince  *string    `json:"stateOrProvince" db:"state_or_province"`
	Tenant           *string    `json:"tenant"`
	TenantID         *int       `json:"tenantId" db:"tenant_id"`
	UID              *int       `json:"uid" db:"uid"`
	UserName         *string    `json:"username" db:"username"`
}

type APIUserPost struct {
	APIUser
	ConfirmLocalPassword *string `json:"confirmLocalPassword" db:"confirm_local_passwd"`
	LocalPassword        *string `json:"localPassword" db:"local_passwd"`
}

type APIUsersResponse struct {
	Response []APIUser `json:"response"`
}

type UserDeliveryServiceDeleteResponse struct {
	Alerts []Alert `json:"alerts"`
}

// UserCurrent contains all the user info about the current user, as returned by /api/1.x/user/current
type UserCurrent struct {
	AddressLine1    *string `json:"addressLine1"`
	AddressLine2    *string `json:"addressLine2"`
	City            *string `json:"city"`
	Company         *string `json:"company,omitempty"`
	Country         *string `json:"country"`
	Email           *string `json:"email,omitempty"`
	FullName        *string `json:"fullName,omitempty"`
	GID             *int    `json:"gid,omitempty"`
	ID              *int    `json:"id,omitempty"`
	LastUpdated     *string `json:"lastUpdated,omitempty"`
	LocalUser       *bool   `json:"localUser"`
	NewUser         *bool   `json:"newUser,omitempty"`
	PhoneNumber     *string `json:"phoneNumber"`
	PostalCode      *string `json:"postalCode"`
	PublicSSHKey    *string `json:"publicSshKey,omitempty"`
	Role            *int    `json:"role,omitempty"`
	RoleName        *string `json:"roleName,omitempty"`
	StateOrProvince *string `json:"stateOrProvince"`
	Tenant          *string `json:"tenant"`
	TenantID        *uint64 `json:"tenantId"`
	UID             *int    `json:"uid,omitempty"`
	UserName        *string `json:"username,omitempty"`
}

type UserCurrentResponse struct {
	Response UserCurrent `json:"response"`
}
