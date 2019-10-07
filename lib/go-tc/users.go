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

import "database/sql"
import "encoding/json"
import "errors"
import "fmt"

import "github.com/apache/trafficcontrol/lib/go-rfc"
import "github.com/apache/trafficcontrol/lib/go-util"

import "github.com/go-ozzo/ozzo-validation"
import "github.com/go-ozzo/ozzo-validation/is"

// UserCredentials contains Traffic Ops login credentials
type UserCredentials struct {
	Username string `json:"u"`
	Password string `json:"p"`
}

// UserToken represents a request payload containing a UUID token for authentication
type UserToken struct {
	Token string `json:"t"`
}

// UserV13 contains non-nullable TO user information
type UserV13 struct {
	Username         string    `json:"username"`
	PublicSSHKey     string    `json:"publicSshKey"`
	Role             int       `json:"role"`
	RoleName         string    `json:"rolename"`
	ID               int       `json:"id"`
	UID              int       `json:"uid"`
	GID              int       `json:"gid"`
	Company          string    `json:"company"`
	Email            string    `json:"email"`
	FullName         string    `json:"fullName"`
	NewUser          bool      `json:"newUser"`
	LastUpdated      string    `json:"lastUpdated"`
	AddressLine1     string    `json:"addressLine1"`
	AddressLine2     string    `json:"addressLine2"`
	City             string    `json:"city"`
	Country          string    `json:"country"`
	PhoneNumber      string    `json:"phoneNumber"`
	PostalCode       string    `json:"postalCode"`
	RegistrationSent TimeNoMod `json:"registrationSent"`
	StateOrProvince  string    `json:"stateOrProvince"`
	Tenant           string    `json:"tenant"`
	TenantID         int       `json:"tenantId"`
}

// commonUserFields is unexported, but its contents are still visible when it is embedded
// LastUpdated is a new field for some structs
type commonUserFields struct {
	AddressLine1    *string `json:"addressLine1" db:"address_line1"`
	AddressLine2    *string `json:"addressLine2" db:"address_line2"`
	City            *string `json:"city" db:"city"`
	Company         *string `json:"company" db:"company"`
	Country         *string `json:"country" db:"country"`
	Email           *string `json:"email" db:"email"`
	FullName        *string `json:"fullName" db:"full_name"`
	GID             *int    `json:"gid"`
	ID              *int    `json:"id" db:"id"`
	NewUser         *bool   `json:"newUser" db:"new_user"`
	PhoneNumber     *string `json:"phoneNumber" db:"phone_number"`
	PostalCode      *string `json:"postalCode" db:"postal_code"`
	PublicSSHKey    *string `json:"publicSshKey" db:"public_ssh_key"`
	Role            *int    `json:"role" db:"role"`
	StateOrProvince *string `json:"stateOrProvince" db:"state_or_province"`
	Tenant          *string `json:"tenant"`
	TenantID        *int    `json:"tenantId" db:"tenant_id"`
	Token           *string `json:"-" db:"token"`
	UID             *int    `json:"uid"`
	//Username        *string    `json:"username" db:"username"`  //not including major change due to naming incompatibility
	LastUpdated *TimeNoMod `json:"lastUpdated" db:"last_updated"`
}

// User fields in v14 have been updated to be nullable
type User struct {
	Username             *string    `json:"username" db:"username"`
	RegistrationSent     *TimeNoMod `json:"registrationSent" db:"registration_sent"`
	LocalPassword        *string    `json:"localPasswd,omitempty" db:"local_passwd"`
	ConfirmLocalPassword *string    `json:"confirmLocalPasswd,omitempty" db:"confirm_local_passwd"`
	RoleName             *string    `json:"roleName,omitempty" db:"-"`
	commonUserFields
}

// UserCurrent represents the profile for the authenticated user
type UserCurrent struct {
	UserName  *string `json:"username"`
	LocalUser *bool   `json:"localUser"`
	RoleName  *string `json:"roleName"`
	commonUserFields
}

// CurrentUserUpdateRequest differs from a regular User/UserCurrent in that many of its fields are
// *parsed* but not *unmarshaled*. This allows a handler to distinguish between "null" and
// "undefined" values.
type CurrentUserUpdateRequest struct {
	// User, for whatever reason, contains all of the actual data.
	User CurrentUserUpdateRequestUser `json:"user"`
}

// CurrentUserUpdateRequestUser holds all of the actual data in a request to update the current user.
type CurrentUserUpdateRequestUser struct {
	AddressLine1       json.RawMessage `json:"addressLine1"`
	AddressLine2       json.RawMessage `json:"addressLine2"`
	City               json.RawMessage `json:"city"`
	Company            json.RawMessage `json:"company"`
	ConfirmLocalPasswd *string         `json:"confirmLocalPasswd"`
	Country            json.RawMessage `json:"country"`
	Email              json.RawMessage `json:"email"`
	FullName           json.RawMessage `json:"fullName"`
	GID                json.RawMessage `json:"gid"`
	ID                 *uint64         `json:"id"`
	LocalPasswd        *string         `json:"localPasswd"`
	PhoneNumber        json.RawMessage `json:"phoneNumber"`
	PostalCode         json.RawMessage `json:"postalCode"`
	PublicSSHKey       json.RawMessage `json:"publicSshKey"`
	Role               *uint64         `json:"role"`
	StateOrProvince    json.RawMessage `json:"stateOrProvince"`
	TenantID           *uint64         `json:"tenantId"`
	UID                json.RawMessage `json:"uid"`
	Username           *string         `json:"username"`
}

// ValidateAndUnmarshal validates the request and returns a User into which the request's information
// has been unmarshalled. This allows many fields to be "null", but explicitly checks that they are
// present in the JSON payload.
func (u *CurrentUserUpdateRequestUser) ValidateAndUnmarshal() (User, error) {
	var user User
	errs := []error{}
	if u.AddressLine1 == nil {
		errs = append(errs, errors.New("addressLine1: required"))
	} else if err := json.Unmarshal(u.AddressLine1, &user.AddressLine1); err != nil {
		errs = append(errs, fmt.Errorf("addressLine1: %v", err))
	}

	if u.AddressLine2 == nil {
		errs = append(errs, errors.New("addressLine2: required"))
	} else if err := json.Unmarshal(u.AddressLine2, &user.AddressLine2); err != nil {
		errs = append(errs, fmt.Errorf("addressLine2: %v", err))
	}

	if u.City == nil {
		errs = append(errs, errors.New("city: required"))
	} else if err := json.Unmarshal(u.City, &user.City); err != nil {
		errs = append(errs, fmt.Errorf("city: %v", err))
	}

	if u.Company == nil {
		errs = append(errs, errors.New("company: required"))
	} else if err := json.Unmarshal(u.Company, &user.Company); err != nil {
		errs = append(errs, fmt.Errorf("company: %v", err))
	}

	user.ConfirmLocalPassword = u.ConfirmLocalPasswd
	user.LocalPassword = u.LocalPasswd

	if u.LocalPasswd != nil && *u.LocalPasswd != "" {
		if u.ConfirmLocalPasswd == nil || *u.ConfirmLocalPasswd == "" {
			errs = append(errs, errors.New("confirmLocalPasswd: required when changing password"))
		} else if *u.LocalPasswd != *u.ConfirmLocalPasswd {
			errs = append(errs, errors.New("localPasswd and confirmLocalPasswd do not match"))
		}
	}

	if u.Country == nil {
		errs = append(errs, errors.New("country: required"))
	} else if err := json.Unmarshal(u.Country, &user.Country); err != nil {
		errs = append(errs, fmt.Errorf("country: %v", err))
	}

	if u.Email == nil {
		errs = append(errs, errors.New("email: required"))
	} else if err := json.Unmarshal(u.Email, &user.Email); err != nil {
		errs = append(errs, fmt.Errorf("email: %v", err))
	}
	if user.Email != nil {
		if err := validation.Validate(*user.Email, is.Email); err != nil {
			errs = append(errs, err)
		}
	}

	if u.FullName == nil {
		errs = append(errs, errors.New("fullName: required"))
	} else if err := json.Unmarshal(u.FullName, &user.FullName); err != nil {
		errs = append(errs, fmt.Errorf("fullName: %v", err))
	}

	if u.GID == nil {
		errs = append(errs, errors.New("gid: required"))
	} else if err := json.Unmarshal(u.GID, &user.GID); err != nil {
		errs = append(errs, fmt.Errorf("gid: %v", err))
	}

	if u.ID == nil {
		errs = append(errs, errors.New("id: required (and can't be null)"))
	} else {
		id := int(*u.ID)
		user.ID = &id
	}

	if u.PhoneNumber == nil {
		errs = append(errs, errors.New("phoneNumber: required"))
	} else if err := json.Unmarshal(u.PhoneNumber, &user.PhoneNumber); err != nil {
		errs = append(errs, fmt.Errorf("phoneNumber: %v", err))
	}

	if u.PostalCode == nil {
		errs = append(errs, errors.New("postalCode: required"))
	} else if err := json.Unmarshal(u.PostalCode, &user.PostalCode); err != nil {
		errs = append(errs, fmt.Errorf("postalCode: %v", err))
	}

	if u.PublicSSHKey == nil {
		errs = append(errs, errors.New("publicSshKey: required"))
	} else if err := json.Unmarshal(u.PublicSSHKey, &user.PublicSSHKey); err != nil {
		errs = append(errs, fmt.Errorf("publicSshKey: %v", err))
	}

	if u.Role == nil {
		errs = append(errs, errors.New("role: required (and can't be null)"))
	} else {
		role := int(*u.Role)
		user.Role = &role
	}

	if u.StateOrProvince == nil {
		errs = append(errs, errors.New("stateOrProvince: required"))
	} else if err := json.Unmarshal(u.StateOrProvince, &user.StateOrProvince); err != nil {
		errs = append(errs, fmt.Errorf("stateOrProvince: %v", err))
	}

	if u.TenantID == nil {
		errs = append(errs, errors.New("tenantId: required"))
	} else {
		tenantID := int(*u.TenantID)
		user.TenantID = &tenantID
	}

	if u.UID == nil {
		errs = append(errs, errors.New("uid: required"))
	} else if err := json.Unmarshal(u.UID, &user.UID); err != nil {
		errs = append(errs, fmt.Errorf("uid: %v", err))
	}

	if u.Username == nil || *u.Username == "" {
		errs = append(errs, errors.New("username: required (and cannot be null or blank)"))
	} else {
		uname := *u.Username
		user.Username = &uname
	}

	var err error
	if len(errs) > 0 {
		err = util.JoinErrs(errs)
	}
	return user, err
}

// ------------------- Response structs -------------------- //
//  Response structs should only be used in the client       //
//  The client's use of these will eventually be deprecated  //
// --------------------------------------------------------- //

type UsersResponseV13 struct {
	Response []UserV13 `json:"response"`
}

type UsersResponse struct {
	Response []User `json:"response"`
}

type CreateUserResponse struct {
	Response User `json:"response"`
	Alerts
}

type UpdateUserResponse struct {
	Response User `json:"response"`
	Alerts
}

type DeleteUserResponse struct {
	Alerts
}

type UserCurrentResponse struct {
	Response UserCurrent `json:"response"`
}

type UserDeliveryServiceDeleteResponse struct {
	Alerts
}

type UserPasswordResetRequest struct {
	Email rfc.EmailAddress `json:"email"`
}

// UserRegistrationRequest is the request submitted by operators when they want to register a new
// user.
type UserRegistrationRequest struct {
	Email rfc.EmailAddress `json:"email"`
	// Role - despite being named "Role" - is actually merely the *ID* of a Role to give the new user.
	Role     uint `json:"role"`
	TenantID uint `json:"tenantId"`
}

// Validate implements the
// github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api.ParseValidator interface.
func (urr *UserRegistrationRequest) Validate(tx *sql.Tx) error {
	var errs = []error{}
	if urr.Role == 0 {
		errs = append(errs, errors.New("role: required and cannot be zero."))
	}

	if urr.TenantID == 0 {
		errs = append(errs, errors.New("tenantId: required and cannot be zero."))
	}

	// This can only happen if an email isn't present in the request; the JSON parse handles actually
	// invalid email addresses.
	if urr.Email.Address.Address == "" {
		errs = append(errs, errors.New("email: required"))
	}

	return util.JoinErrs(errs)
}
