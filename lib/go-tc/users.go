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
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// Upgrade converts a User to a UserV4 (as seen in API versions 4.x).
func (u User) Upgrade() UserV4 {
	var ret UserV4
	ret.AddressLine1 = util.CopyIfNotNil(u.AddressLine1)
	ret.AddressLine2 = util.CopyIfNotNil(u.AddressLine2)
	ret.City = util.CopyIfNotNil(u.City)
	ret.Company = util.CopyIfNotNil(u.Company)
	ret.Country = util.CopyIfNotNil(u.Country)
	ret.Email = util.CopyIfNotNil(u.Email)
	ret.GID = util.CopyIfNotNil(u.GID)
	ret.ID = util.CopyIfNotNil(u.ID)
	ret.LocalPassword = util.CopyIfNotNil(u.LocalPassword)
	ret.PhoneNumber = util.CopyIfNotNil(u.PhoneNumber)
	ret.PostalCode = util.CopyIfNotNil(u.PostalCode)
	ret.PublicSSHKey = util.CopyIfNotNil(u.PublicSSHKey)
	ret.StateOrProvince = util.CopyIfNotNil(u.StateOrProvince)
	ret.Tenant = util.CopyIfNotNil(u.Tenant)
	ret.Token = util.CopyIfNotNil(u.Token)
	ret.UID = util.CopyIfNotNil(u.UID)
	ret.FullName = util.CopyIfNotNil(u.FullName)
	if u.LastUpdated != nil {
		ret.LastUpdated = u.LastUpdated.Time
	}
	if u.NewUser != nil {
		ret.NewUser = *u.NewUser
	}
	if u.RegistrationSent != nil {
		ret.RegistrationSent = new(time.Time)
		*ret.RegistrationSent = u.RegistrationSent.Time
	}
	if u.RoleName != nil {
		ret.Role = *u.RoleName
	}
	if u.TenantID != nil {
		ret.TenantID = *u.TenantID
	}
	if u.Username != nil {
		ret.Username = *u.Username
	}
	return ret
}

// Downgrade converts a UserV4 to a User (as seen in API versions < 4.0).
func (u UserV4) Downgrade() User {
	var ret User
	ret.FullName = new(string)
	ret.FullName = u.FullName
	ret.LastUpdated = TimeNoModFromTime(u.LastUpdated)
	ret.NewUser = new(bool)
	*ret.NewUser = u.NewUser
	ret.RoleName = new(string)
	*ret.RoleName = u.Role
	ret.Role = nil
	ret.TenantID = new(int)
	*ret.TenantID = u.TenantID
	ret.Username = new(string)
	*ret.Username = u.Username

	ret.AddressLine1 = util.CopyIfNotNil(u.AddressLine1)
	ret.AddressLine2 = util.CopyIfNotNil(u.AddressLine2)
	ret.City = util.CopyIfNotNil(u.City)
	ret.Company = util.CopyIfNotNil(u.Company)
	ret.ConfirmLocalPassword = util.CopyIfNotNil(u.LocalPassword)
	ret.Country = util.CopyIfNotNil(u.Country)
	ret.Email = util.CopyIfNotNil(u.Email)
	ret.GID = util.CopyIfNotNil(u.GID)
	ret.ID = util.CopyIfNotNil(u.ID)
	ret.LocalPassword = util.CopyIfNotNil(u.LocalPassword)
	ret.PhoneNumber = util.CopyIfNotNil(u.PhoneNumber)
	ret.PostalCode = util.CopyIfNotNil(u.PostalCode)
	ret.PublicSSHKey = util.CopyIfNotNil(u.PublicSSHKey)
	if u.RegistrationSent != nil {
		ret.RegistrationSent = TimeNoModFromTime(*u.RegistrationSent)
	}
	ret.StateOrProvince = util.CopyIfNotNil(u.StateOrProvince)
	ret.Tenant = util.CopyIfNotNil(u.Tenant)
	ret.Token = util.CopyIfNotNil(u.Token)
	ret.UID = util.CopyIfNotNil(u.UID)

	return ret
}

// UserCredentials contains Traffic Ops login credentials.
type UserCredentials struct {
	Username string `json:"u"`
	Password string `json:"p"`
}

// UserToken represents a request payload containing a UUID token for
// authentication.
type UserToken struct {
	Token string `json:"t"`
}

// commonUserFields contents are still visible when it is embedded
// LastUpdated is a new field for some structs.
type commonUserFields struct {
	AddressLine1    *string    `json:"addressLine1" db:"address_line1"`
	AddressLine2    *string    `json:"addressLine2" db:"address_line2"`
	City            *string    `json:"city" db:"city"`
	Company         *string    `json:"company" db:"company"`
	Country         *string    `json:"country" db:"country"`
	Email           *string    `json:"email" db:"email"`
	FullName        *string    `json:"fullName" db:"full_name"`
	GID             *int       `json:"gid"`
	ID              *int       `json:"id" db:"id"`
	NewUser         *bool      `json:"newUser" db:"new_user"`
	PhoneNumber     *string    `json:"phoneNumber" db:"phone_number"`
	PostalCode      *string    `json:"postalCode" db:"postal_code"`
	PublicSSHKey    *string    `json:"publicSshKey" db:"public_ssh_key"`
	Role            *int       `json:"role" db:"role"`
	StateOrProvince *string    `json:"stateOrProvince" db:"state_or_province"`
	Tenant          *string    `json:"tenant"`
	TenantID        *int       `json:"tenantId" db:"tenant_id"`
	Token           *string    `json:"-" db:"token"`
	UID             *int       `json:"uid"`
	LastUpdated     *TimeNoMod `json:"lastUpdated" db:"last_updated"`
}

// User represents a user of Traffic Ops.
type User struct {
	Username             *string    `json:"username" db:"username"`
	RegistrationSent     *TimeNoMod `json:"registrationSent" db:"registration_sent"`
	LocalPassword        *string    `json:"localPasswd,omitempty" db:"local_passwd"`
	ConfirmLocalPassword *string    `json:"confirmLocalPasswd,omitempty" db:"confirm_local_passwd"`
	// NOTE: RoleName db:"-" tag is required due to clashing with the DB query here:
	// https://github.com/apache/trafficcontrol/blob/3b5dd406bf1a0bb456c062b0f6a465ec0617d8ef/traffic_ops/traffic_ops_golang/user/user.go#L197
	// It's done that way in order to maintain "rolename" vs "roleName" JSON field capitalization for the different users APIs.
	RoleName *string `json:"roleName,omitempty" db:"role_name"`
	commonUserFields
}

// UserCurrent represents the profile for the authenticated user.
type UserCurrent struct {
	UserName  *string `json:"username"`
	LocalUser *bool   `json:"localUser"`
	RoleName  *string `json:"roleName"`
	commonUserFields
}

// ToLegacyCurrentUser will convert an APIv4 user to an APIv3 "current user"
// representation. A Role ID and "local user" value must be supplied, since the
// APIv4 User doesn't have them.
func (u UserV4) ToLegacyCurrentUser(roleID int, localUser bool) UserCurrent {
	var ret UserCurrent
	ret.FullName = new(string)
	*ret.FullName = *u.FullName
	ret.LastUpdated = TimeNoModFromTime(u.LastUpdated)
	ret.NewUser = new(bool)
	*ret.NewUser = u.NewUser
	ret.RoleName = new(string)
	*ret.RoleName = u.Role
	ret.Role = new(int)
	*ret.Role = roleID
	ret.TenantID = new(int)
	*ret.TenantID = u.TenantID
	ret.Tenant = u.Tenant
	ret.UserName = new(string)
	*ret.UserName = u.Username
	ret.LocalUser = new(bool)
	*ret.LocalUser = localUser
	ret.Token = util.CopyIfNotNil(u.Token)
	ret.AddressLine1 = util.CopyIfNotNil(u.AddressLine1)
	ret.AddressLine2 = util.CopyIfNotNil(u.AddressLine2)
	ret.City = util.CopyIfNotNil(u.City)
	ret.Company = util.CopyIfNotNil(u.Company)
	ret.Country = util.CopyIfNotNil(u.Country)
	ret.Email = util.CopyIfNotNil(u.Email)
	ret.GID = util.CopyIfNotNil(u.GID)
	ret.ID = util.CopyIfNotNil(u.ID)
	ret.UID = util.CopyIfNotNil(u.UID)
	ret.PhoneNumber = util.CopyIfNotNil(u.PhoneNumber)
	ret.PostalCode = util.CopyIfNotNil(u.PostalCode)
	ret.PublicSSHKey = util.CopyIfNotNil(u.PublicSSHKey)
	ret.StateOrProvince = util.CopyIfNotNil(u.StateOrProvince)
	ret.Tenant = util.CopyIfNotNil(u.Tenant)
	ret.Token = util.CopyIfNotNil(u.Token)
	ret.UID = util.CopyIfNotNil(u.UID)

	return ret
}

// UserV4 is an alias for the User struct used for the latest minor version associated with api major version 4.
type UserV4 UserV40

// A UserV40 is a representation of a Traffic Ops user as it appears in version
// 4.0 of Traffic Ops's API.
type UserV40 struct {
	AddressLine1   *string `json:"addressLine1" db:"address_line1"`
	AddressLine2   *string `json:"addressLine2" db:"address_line2"`
	ChangeLogCount int     `json:"changeLogCount" db:"change_log_count"`
	City           *string `json:"city" db:"city"`
	Company        *string `json:"company" db:"company"`
	Country        *string `json:"country" db:"country"`
	Email          *string `json:"email" db:"email"`
	FullName       *string `json:"fullName" db:"full_name"`
	// Deprecated: This has no known use, and will likely be removed in future
	// API versions.
	GID               *int       `json:"gid"`
	ID                *int       `json:"id" db:"id"`
	LastAuthenticated *time.Time `json:"lastAuthenticated" db:"last_authenticated"`
	LastUpdated       time.Time  `json:"lastUpdated" db:"last_updated"`
	LocalPassword     *string    `json:"localPasswd,omitempty" db:"local_passwd"`
	NewUser           bool       `json:"newUser" db:"new_user"`
	PhoneNumber       *string    `json:"phoneNumber" db:"phone_number"`
	PostalCode        *string    `json:"postalCode" db:"postal_code"`
	PublicSSHKey      *string    `json:"publicSshKey" db:"public_ssh_key"`
	RegistrationSent  *time.Time `json:"registrationSent" db:"registration_sent"`
	Role              string     `json:"role" db:"role"`
	StateOrProvince   *string    `json:"stateOrProvince" db:"state_or_province"`
	Tenant            *string    `json:"tenant"`
	TenantID          int        `json:"tenantId" db:"tenant_id"`
	Token             *string    `json:"-" db:"token"`
	UCDN              string     `json:"ucdn"`
	// Deprecated: This has no known use, and will likely be removed in future
	// API versions.
	UID      *int   `json:"uid"`
	Username string `json:"username" db:"username"`
}

// UsersResponseV4 is the type of a response from Traffic Ops to requests made
// to /users which return more than one user for the latest 4.x api version variant.
type UsersResponseV4 struct {
	Response []UserV4 `json:"response"`
	Alerts
}

// UserResponseV4 is the type of a response from Traffic Ops to requests made
// to /users which return one user for the latest 4.x api version variant.
type UserResponseV4 struct {
	Response UserV4 `json:"response"`
	Alerts
}

// CurrentUserUpdateRequest differs from a regular User/UserCurrent in that many of its fields are
// *parsed* but not *unmarshaled*. This allows a handler to distinguish between "null" and
// "undefined" values.
type CurrentUserUpdateRequest struct {
	// User, for whatever reason, contains all of the actual data.
	User *CurrentUserUpdateRequestUser `json:"user"`
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
	ID                 json.RawMessage `json:"id"`
	LocalPasswd        *string         `json:"localPasswd"`
	PhoneNumber        json.RawMessage `json:"phoneNumber"`
	PostalCode         json.RawMessage `json:"postalCode"`
	PublicSSHKey       json.RawMessage `json:"publicSshKey"`
	Role               json.RawMessage `json:"role"`
	StateOrProvince    json.RawMessage `json:"stateOrProvince"`
	TenantID           json.RawMessage `json:"tenantId"`
	UID                json.RawMessage `json:"uid"`
	Username           json.RawMessage `json:"username"`
}

// Upgrade converts an APIv3 and earlier "current user" to an APIv4 User.
// Fields not present in earlier API versions need to be passed explicitly.
func (u UserCurrent) Upgrade(registrationSent, lastAuthenticated *time.Time, ucdn string, changeLogCount int) UserV4 {
	var ret UserV4
	ret.AddressLine1 = util.CopyIfNotNil(u.AddressLine1)
	ret.AddressLine2 = util.CopyIfNotNil(u.AddressLine2)
	ret.ChangeLogCount = changeLogCount
	ret.City = util.CopyIfNotNil(u.City)
	ret.Company = util.CopyIfNotNil(u.Company)
	ret.Country = util.CopyIfNotNil(u.Country)
	ret.Email = util.CopyIfNotNil(u.Email)
	ret.GID = util.CopyIfNotNil(u.GID)
	ret.ID = util.CopyIfNotNil(u.ID)
	ret.LastAuthenticated = lastAuthenticated
	ret.PhoneNumber = util.CopyIfNotNil(u.PhoneNumber)
	ret.PostalCode = util.CopyIfNotNil(u.PostalCode)
	ret.PublicSSHKey = util.CopyIfNotNil(u.PublicSSHKey)
	ret.RegistrationSent = registrationSent
	ret.StateOrProvince = util.CopyIfNotNil(u.StateOrProvince)
	ret.Tenant = util.CopyIfNotNil(u.Tenant)
	ret.Token = util.CopyIfNotNil(u.Token)
	ret.UCDN = ucdn
	ret.UID = util.CopyIfNotNil(u.UID)
	ret.FullName = u.FullName
	if u.LastUpdated != nil {
		ret.LastUpdated = u.LastUpdated.Time
	}
	if u.NewUser != nil {
		ret.NewUser = *u.NewUser
	}

	if u.RoleName != nil {
		ret.Role = *u.RoleName
	}
	if u.TenantID != nil {
		ret.TenantID = *u.TenantID
	}
	if u.UserName != nil {
		ret.Username = *u.UserName
	}
	return ret
}

// UnmarshalAndValidate validates the request and returns a User into which the request's information
// has been unmarshalled.
func (u *CurrentUserUpdateRequestUser) UnmarshalAndValidate(user *User) error {
	errs := []error{}
	if u.AddressLine1 != nil {
		if err := json.Unmarshal(u.AddressLine1, &user.AddressLine1); err != nil {
			errs = append(errs, fmt.Errorf("addressLine1: %w", err))
		}
	}

	if u.AddressLine2 != nil {
		if err := json.Unmarshal(u.AddressLine2, &user.AddressLine2); err != nil {
			errs = append(errs, fmt.Errorf("addressLine2: %w", err))
		}
	}

	if u.City != nil {
		if err := json.Unmarshal(u.City, &user.City); err != nil {
			errs = append(errs, fmt.Errorf("city: %w", err))
		}
	}

	if u.Company != nil {
		if err := json.Unmarshal(u.Company, &user.Company); err != nil {
			errs = append(errs, fmt.Errorf("company: %w", err))
		}
	}

	user.ConfirmLocalPassword = u.ConfirmLocalPasswd
	user.LocalPassword = u.LocalPasswd

	if u.Country != nil {
		if err := json.Unmarshal(u.Country, &user.Country); err != nil {
			errs = append(errs, fmt.Errorf("country: %w", err))
		}
	}

	if u.Email != nil {
		if err := json.Unmarshal(u.Email, &user.Email); err != nil {
			errs = append(errs, fmt.Errorf("email: %w", err))
		} else if user.Email == nil || *user.Email == "" {
			errs = append(errs, errors.New("email: cannot be null or an empty string"))
		} else if err = validation.Validate(*user.Email, is.Email); err != nil {
			errs = append(errs, err)
		}
	}

	if u.FullName != nil {
		if err := json.Unmarshal(u.FullName, &user.FullName); err != nil {
			errs = append(errs, fmt.Errorf("fullName: %w", err))
		} else if user.FullName == nil || *user.FullName == "" {
			// Perl enforced this
			errs = append(errs, errors.New("fullName: cannot be set to 'null' or empty string"))
		}
	}

	if u.GID != nil {
		if err := json.Unmarshal(u.GID, &user.GID); err != nil {
			errs = append(errs, fmt.Errorf("gid: %w", err))
		}
	}

	if u.ID != nil {
		var uid int
		if err := json.Unmarshal(u.ID, &uid); err != nil {
			errs = append(errs, fmt.Errorf("id: %w", err))
		} else if user.ID != nil && *user.ID != uid {
			errs = append(errs, errors.New("id: cannot change user id"))
		} else {
			user.ID = &uid
		}
	}

	if u.PhoneNumber != nil {
		if err := json.Unmarshal(u.PhoneNumber, &user.PhoneNumber); err != nil {
			errs = append(errs, fmt.Errorf("phoneNumber: %w", err))
		}
	}

	if u.PostalCode != nil {
		if err := json.Unmarshal(u.PostalCode, &user.PostalCode); err != nil {
			errs = append(errs, fmt.Errorf("postalCode: %w", err))
		}
	}

	if u.PublicSSHKey != nil {
		if err := json.Unmarshal(u.PublicSSHKey, &user.PublicSSHKey); err != nil {
			errs = append(errs, fmt.Errorf("publicSshKey: %w", err))
		}
	}

	if u.Role != nil {
		if err := json.Unmarshal(u.Role, &user.Role); err != nil {
			errs = append(errs, fmt.Errorf("role: %w", err))
		} else if user.Role == nil {
			errs = append(errs, errors.New("role: cannot be null"))
		}
	}

	if u.StateOrProvince != nil {
		if err := json.Unmarshal(u.StateOrProvince, &user.StateOrProvince); err != nil {
			errs = append(errs, fmt.Errorf("stateOrProvince: %w", err))
		}
	}

	if u.TenantID != nil {
		if err := json.Unmarshal(u.TenantID, &user.TenantID); err != nil {
			errs = append(errs, fmt.Errorf("tenantID: %w", err))
		} else if user.TenantID == nil {
			errs = append(errs, errors.New("tenantID: cannot be null"))
		}
	}

	if u.UID != nil {
		if err := json.Unmarshal(u.UID, &user.UID); err != nil {
			errs = append(errs, fmt.Errorf("uid: %w", err))
		}
	}

	if u.Username != nil {
		if err := json.Unmarshal(u.Username, &user.Username); err != nil {
			errs = append(errs, fmt.Errorf("username: %w", err))
		} else if user.Username == nil || *user.Username == "" {
			errs = append(errs, errors.New("username: cannot be null or empty string"))
		}
	}

	return util.JoinErrs(errs)
}

// UsersResponse can hold a Traffic Ops API response to a request to get a list of users.
type UsersResponse struct {
	Response []User `json:"response"`
	Alerts
}

// UserResponse can hold a Traffic Ops API response to a request to get a user.
type UserResponse struct {
	Response User `json:"response"`
	Alerts
}

// CreateUserResponse can hold a Traffic Ops API response to a POST request to create a user.
type CreateUserResponse struct {
	Response User `json:"response"`
	Alerts
}

// CreateUserResponseV4 can hold a Traffic Ops API response to a POST request to create a user in api v4.
type CreateUserResponseV4 struct {
	Response UserV4 `json:"response"`
	Alerts
}

// UpdateUserResponse can hold a Traffic Ops API response to a PUT request to update a user.
type UpdateUserResponse struct {
	Response User `json:"response"`
	Alerts
}

// UpdateUserResponseV4 can hold a Traffic Ops API response to a PUT request to update a user for the latest 4.x api version variant.
type UpdateUserResponseV4 struct {
	Response UserV4 `json:"response"`
	Alerts
}

// DeleteUserResponse can theoretically hold a Traffic Ops API response to a
// DELETE request to update a user. It is unused.
type DeleteUserResponse struct {
	Alerts
}

// UserCurrentResponse can hold a Traffic Ops API response to a request to get
// or update the current user.
type UserCurrentResponse struct {
	Response UserCurrent `json:"response"`
	Alerts
}

// UserCurrentResponseV4 is the latest 4.x Traffic Ops API version variant of UserResponse.
type UserCurrentResponseV4 struct {
	Response UserV4 `json:"response"`
	Alerts
}

// UserDeliveryServiceDeleteResponse can hold a Traffic Ops API response to
// a request to remove a delivery service from a user.
type UserDeliveryServiceDeleteResponse struct {
	Alerts
}

// UserPasswordResetRequest can hold  Traffic Ops API request to reset a user's password.
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

// UserRegistrationRequestV4 is the alias for the UserRegistrationRequest for the latest 4.x api version variant.
type UserRegistrationRequestV4 UserRegistrationRequestV40

// UserRegistrationRequestV40 is the request submitted by operators when they want to register a new
// user in api V4.
type UserRegistrationRequestV40 struct {
	Email    rfc.EmailAddress `json:"email"`
	Role     string           `json:"role"`
	TenantID uint             `json:"tenantId"`
}

// Validate implements the
// github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api.ParseValidator
// interface.
func (urr *UserRegistrationRequestV4) Validate(tx *sql.Tx) error {
	var errs = []error{}
	if urr.Role == "" {
		errs = append(errs, errors.New("role: required and cannot be empty"))
	}

	if urr.TenantID == 0 {
		errs = append(errs, errors.New("tenantId: required and cannot be zero"))
	}

	// This can only happen if an email isn't present in the request; the JSON parse handles actually
	// invalid email addresses.
	if urr.Email.Address.Address == "" {
		errs = append(errs, errors.New("email: required"))
	}

	return util.JoinErrs(errs)
}

// Validate implements the
// github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api.ParseValidator
// interface.
func (urr *UserRegistrationRequest) Validate(tx *sql.Tx) error {
	var errs = []error{}
	if urr.Role == 0 {
		errs = append(errs, errors.New("role: required and cannot be zero"))
	}

	if urr.TenantID == 0 {
		errs = append(errs, errors.New("tenantId: required and cannot be zero"))
	}

	// This can only happen if an email isn't present in the request; the JSON parse handles actually
	// invalid email addresses.
	if urr.Email.Address.Address == "" {
		errs = append(errs, errors.New("email: required"))
	}

	return util.JoinErrs(errs)
}
