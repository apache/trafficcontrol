package user

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
	"net/http"
	"strings"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"

	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/tenant"
)

const replacePasswordQuery = `
UPDATE tm_user
SET local_passwd=$1
WHERE id=$2
`

const replaceConfirmPasswordQuery = `
UPDATE tm_user
SET confirm_local_passwd=$1
WHERE id=$2
`

const replacePasswordV4Query = `
UPDATE tm_user
SET
	confirm_local_passwd=$1,
	local_passwd=$1
WHERE id=$2
`

const replaceCurrentQuery = `
UPDATE tm_user
SET
	address_line1=$1,
	address_line2=$2,
	city=$3,
	company=$4,
	country=$5,
	email=$6,
	full_name=$7,
	gid=$8,
	new_user=FALSE,
	phone_number=$9,
	postal_code=$10,
	public_ssh_key=$11,
	role=$12,
	state_or_province=$13,
	tenant_id=$14,
	token=NULL,
	uid=$15,
	username=$16
WHERE id=$17
RETURNING
	address_line1,
	address_line2,
	city,
	company,
	country,
	email,
	full_name,
	gid,
	id,
	last_updated,
	new_user,
	phone_number,
	postal_code,
	public_ssh_key,
	role,
	(
		SELECT role.name
		FROM role
		WHERE role.id=tm_user.role
	),
	state_or_province,
	(
		SELECT tenant.name
		FROM tenant
		WHERE tenant.id=tm_user.tenant_id
	),
	tenant_id,
	uid,
	username
`

const replaceCurrentV4Query = `
UPDATE tm_user
SET
	address_line1=$1,
	address_line2=$2,
	city=$3,
	company=$4,
	country=$5,
	email=$6,
	full_name=$7,
	gid=$8,
	new_user=FALSE,
	phone_number=$9,
	postal_code=$10,
	public_ssh_key=$11,
	role=(
		SELECT role.id
		FROM role
		WHERE name=$12
	),
	state_or_province=$13,
	tenant_id=$14,
	token=NULL,
	ucdn=$15,
	uid=$16,
	username=$17
WHERE id=$18
RETURNING
	address_line1,
	address_line2,
	(
		SELECT count(l.tm_user)
		FROM log AS l
		WHERE l.tm_user = tm_user.id
	),
	city,
	company,
	country,
	email,
	full_name,
	gid,
	id,
	last_authenticated,
	last_updated,
	new_user,
	phone_number,
	postal_code,
	public_ssh_key,
	registration_sent,
	(
		SELECT role.name
		FROM role
		WHERE role.id=tm_user.role
	),
	state_or_province,
	(
		SELECT tenant.name
		FROM tenant
		WHERE tenant.id=tm_user.tenant_id
	),
	tenant_id,
	ucdn,
	uid,
	username
`

func Current(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	if inf.Version.Major < 4 {
		cu, err := getLegacyUser(tx, inf.User.ID)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("getting legacy current user: %w", err))
			return
		}
		api.WriteResp(w, r, cu)
		return
	}
	currentUser, err := getUser(inf.Tx.Tx, inf.User.ID)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("getting current user: %w", err))
		return
	}
	api.WriteResp(w, r, currentUser)
}

func getUser(tx *sql.Tx, id int) (tc.UserV4, error) {
	q := `
SELECT
u.address_line1,
u.address_line2,
(
	SELECT count(l.tm_user)
	FROM log AS l
	WHERE l.tm_user = u.id
),
u.city,
u.company,
u.country,
u.email,
u.full_name,
u.gid,
u.id,
u.last_authenticated,
u.last_updated,
u.new_user,
u.phone_number,
u.postal_code,
u.public_ssh_key,
u.registration_sent,
r.name as "role",
u.state_or_province,
t.name as tenant,
u.tenant_id,
u.ucdn,
u.uid,
u.username
FROM tm_user as u
LEFT JOIN role as r ON r.id = u.role
INNER JOIN tenant as t ON t.id = u.tenant_id
WHERE u.id=$1
`
	var u tc.UserV4
	err := tx.QueryRow(q, id).Scan(
		&u.AddressLine1,
		&u.AddressLine2,
		&u.ChangeLogCount,
		&u.City,
		&u.Company,
		&u.Country,
		&u.Email,
		&u.FullName,
		&u.GID,
		&u.ID,
		&u.LastAuthenticated,
		&u.LastUpdated,
		&u.NewUser,
		&u.PhoneNumber,
		&u.PostalCode,
		&u.PublicSSHKey,
		&u.RegistrationSent,
		&u.Role,
		&u.StateOrProvince,
		&u.Tenant,
		&u.TenantID,
		&u.UCDN,
		&u.UID,
		&u.Username,
	)
	if err != nil {
		err = fmt.Errorf("querying current user: %w", err)
	}
	return u, err
}

func getLegacyUser(tx *sql.Tx, id int) (tc.UserCurrent, error) {
	q := `
SELECT
u.address_line1,
u.address_line2,
u.city,
u.company,
u.country,
u.email,
u.full_name,
u.gid,
u.id,
u.last_updated,
u.local_passwd IS NOT NULL,
u.new_user,
u.phone_number,
u.postal_code,
u.public_ssh_key,
u.role as "role",
r.name as role_name,
u.state_or_province,
t.name as tenant,
u.tenant_id,
u.uid,
u.username
FROM tm_user as u
LEFT JOIN role as r ON r.id = u.role
INNER JOIN tenant as t ON t.id = u.tenant_id
WHERE u.id=$1
`
	var u tc.UserCurrent
	err := tx.QueryRow(q, id).Scan(
		&u.AddressLine1,
		&u.AddressLine2,
		&u.City,
		&u.Company,
		&u.Country,
		&u.Email,
		&u.FullName,
		&u.GID,
		&u.ID,
		&u.LastUpdated,
		&u.LocalUser,
		&u.NewUser,
		&u.PhoneNumber,
		&u.PostalCode,
		&u.PublicSSHKey,
		&u.Role,
		&u.RoleName,
		&u.StateOrProvince,
		&u.Tenant,
		&u.TenantID,
		&u.UID,
		&u.UserName,
	)
	if err != nil {
		err = fmt.Errorf("querying legacy current user: %w", err)
	}
	return u, err
}

func ReplaceCurrent(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	var userRequest tc.CurrentUserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		errCode = http.StatusBadRequest
		userErr = fmt.Errorf("couldn't parse request: %w", err)
		api.HandleErr(w, r, tx, errCode, userErr, nil)
		return
	}

	if userRequest.User == nil {
		errCode = http.StatusBadRequest
		userErr = fmt.Errorf("missing required 'user' object")
		api.HandleErr(w, r, tx, errCode, userErr, nil)
		return
	}

	user, exists, err := dbhelpers.GetUserByID(inf.User.ID, tx)
	if err != nil {
		sysErr = fmt.Errorf("getting user by ID %d: %w", inf.User.ID, err)
		errCode = http.StatusInternalServerError
		api.HandleErr(w, r, tx, errCode, nil, sysErr)
		return
	}
	if !exists {
		sysErr = fmt.Errorf("current user (#%d) doesn't exist... ??", inf.User.ID)
		errCode = http.StatusInternalServerError
		api.HandleErr(w, r, tx, errCode, nil, sysErr)
		return
	}

	if err := userRequest.User.UnmarshalAndValidate(&user); err != nil {
		errCode = http.StatusBadRequest
		userErr = fmt.Errorf("couldn't parse request: %w", err)
		api.HandleErr(w, r, tx, errCode, userErr, nil)
		return
	}

	changePasswd := false
	changeConfirmPasswd := false

	// obfuscate passwords (UnmarshalAndValidate checks for equality with ConfirmLocalPassword)
	// TODO: check for valid password via bad password list like Perl did? User creation doesn't...
	if user.LocalPassword != nil && *user.LocalPassword != "" {
		if ok, err := auth.IsGoodPassword(*user.LocalPassword); !ok {
			errCode = http.StatusBadRequest
			if err != nil {
				userErr = err
			} else {
				userErr = errors.New("unacceptable password")
			}
			api.HandleErr(w, r, tx, errCode, userErr, nil)
			return
		}

		hashPass, err := auth.DerivePassword(*user.LocalPassword)
		if err != nil {
			sysErr = fmt.Errorf("hashing new password: %w", err)
			errCode = http.StatusInternalServerError
			api.HandleErr(w, r, tx, errCode, nil, sysErr)
			return
		}
		changePasswd = true
		user.LocalPassword = util.StrPtr(hashPass)
	}

	// Perl did this although it serves no known purpose
	if user.ConfirmLocalPassword != nil && *user.ConfirmLocalPassword != "" {
		hashPass, err := auth.DerivePassword(*user.ConfirmLocalPassword)
		if err != nil {
			sysErr = fmt.Errorf("hashing new 'confirm' password: %w", err)
			errCode = http.StatusInternalServerError
			api.HandleErr(w, r, tx, errCode, nil, sysErr)
			return
		}
		user.ConfirmLocalPassword = util.StrPtr(hashPass)
		changeConfirmPasswd = true
	}

	if *user.Role != inf.User.Role {
		privLevel, exists, err := dbhelpers.GetPrivLevelFromRoleID(tx, *user.Role)
		if err != nil {
			sysErr = fmt.Errorf("getting privLevel for Role #%d: %w", *user.Role, err)
			errCode = http.StatusInternalServerError
			api.HandleErr(w, r, tx, errCode, nil, sysErr)
			return
		}
		if !exists {
			userErr = fmt.Errorf("role: no such role: %d", *user.Role)
			errCode = http.StatusNotFound
			api.HandleErr(w, r, tx, errCode, userErr, nil)
			return
		}
		if privLevel > inf.User.PrivLevel {
			userErr = errors.New("role: cannot have greater permissions than user's current role")
			errCode = http.StatusForbidden
			api.HandleErr(w, r, tx, errCode, userErr, nil)
			return
		}
	}

	if ok, err := tenant.IsResourceAuthorizedToUserTx(*user.TenantID, inf.User, tx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			userErr = errors.New("no such tenant")
			errCode = http.StatusNotFound
		} else {
			sysErr = fmt.Errorf("checking user %s permissions on tenant #%d: %w", inf.User.UserName, *user.TenantID, err)
			errCode = http.StatusInternalServerError
		}
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	} else if !ok {
		// unlike Perl, this endpoint will not disclose the existence of tenants over which the current
		// user has no permission - in keeping with the behavior of the '/tenants' endpoint.
		userErr = errors.New("no such tenant")
		errCode = http.StatusNotFound
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	if *user.Username != inf.User.UserName {
		if ok, err := dbhelpers.UsernameExists(*user.Username, tx); err != nil {
			sysErr = fmt.Errorf("checking existence of user %s: %w", *user.Username, err)
			errCode = http.StatusInternalServerError
			api.HandleErr(w, r, tx, errCode, nil, sysErr)
			return
		} else if ok {
			// TODO users are tenanted, so theoretically I should be hiding the existence of the
			// conflicting user - but then how do I tell the client how to fix their request?
			userErr = fmt.Errorf("username %s already exists", *user.Username)
			errCode = http.StatusConflict
			api.HandleErr(w, r, tx, errCode, userErr, nil)
			return
		}
	}

	if err = updateLegacyUser(&user, tx, changePasswd, changeConfirmPasswd); err != nil {
		userErr, sysErr, statusCode := api.ParseDBError(err)
		sysErr = fmt.Errorf("updating legacy user: %w", err)
		api.HandleErr(w, r, tx, statusCode, userErr, sysErr)
		return
	}

	api.WriteRespAlertObj(w, r, tc.SuccessLevel, "User profile was successfully updated", user)
}

func validateV4(user tc.UserV4, inf *api.Info) (error, error) {
	validateErrs := validation.Errors{
		"email":    validation.Validate(user.Email, validation.Required, is.Email),
		"fullName": validation.Validate(user.FullName, validation.Required),
		"role":     validation.Validate(user.Role, validation.Required),
		"username": validation.Validate(user.Username, validation.Required),
		"tenantID": validation.Validate(user.TenantID, validation.Required),
	}

	// Password is not required for update
	if user.LocalPassword != nil {
		ok, err := auth.IsGoodLoginPair(user.Username, *user.LocalPassword)
		if err != nil {
			return err, nil
		}
		if !ok {
			return errors.New("unacceptable password"), nil
		}
	}

	if err := tovalidate.ToError(validateErrs); err != nil {
		return err, nil
	}

	caps, err := dbhelpers.GetCapabilitiesFromRoleName(inf.Tx.Tx, user.Role)
	if err != nil {
		return nil, fmt.Errorf("getting capabilities for user's requested Role (%s): %w", user.Role, err)
	}

	missing := inf.User.MissingPermissions(caps...)
	if len(missing) > 0 {
		return nil, fmt.Errorf("cannot request more than assigned permissions, current user needs %s permissions", strings.Join(missing, ","))
	}

	if user.Username != inf.User.UserName {
		if ok, err := dbhelpers.UsernameExists(user.Username, inf.Tx.Tx); err != nil {
			return nil, fmt.Errorf("checking existence of user %s: %w", user.Username, err)
		} else if ok {
			return fmt.Errorf("username %s already exists", user.Username), nil
		}
	}

	return nil, nil
}

// ReplaceCurrentV4 replaces the current user with the definition in the user's
// request (assuming it meets validation constraints).
func ReplaceCurrentV4(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	var user tc.UserV4
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		errCode = http.StatusBadRequest
		userErr = fmt.Errorf("couldn't parse request: %w", err)
		api.HandleErr(w, r, tx, errCode, userErr, nil)
		return
	}
	// Token must never be updated this way
	user.Token = nil

	user.ID = new(int)
	*user.ID = inf.User.ID

	userErr, sysErr = validateV4(user, inf)
	if userErr != nil || sysErr != nil {
		errCode = http.StatusBadRequest
		if sysErr != nil {
			errCode = http.StatusInternalServerError
		}
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	roleID, ok, err := dbhelpers.GetRoleIDFromName(tx, user.Role)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
		return
	} else if !ok {
		api.HandleErr(w, r, tx, http.StatusNotFound, errors.New("no such role"), nil)
		return
	}
	if inf.User.Role != roleID {
		api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("users cannot update their own role"), nil)
		return
	}

	changePasswd := false

	// obfuscate password
	if user.LocalPassword != nil {
		hashPass, err := auth.DerivePassword(*user.LocalPassword)
		if err != nil {
			sysErr = fmt.Errorf("hashing new password for user %s (#%d): %w", inf.User.UserName, inf.User.ID, err)
			errCode = http.StatusInternalServerError
			api.HandleErr(w, r, tx, errCode, nil, sysErr)
			return
		}
		changePasswd = true
		*user.LocalPassword = hashPass
	}

	if ok, err := tenant.IsResourceAuthorizedToUserTx(user.TenantID, inf.User, tx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			userErr = fmt.Errorf("no such tenant: #%d", user.TenantID)
			errCode = http.StatusNotFound
		} else {
			sysErr = fmt.Errorf("checking user %s permissions on tenant #%d: %w", inf.User.UserName, user.TenantID, err)
			errCode = http.StatusInternalServerError
		}
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	} else if !ok {
		userErr = fmt.Errorf("no such tenant: #%d", user.TenantID)
		errCode = http.StatusNotFound
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	if err := updateUser(&user, tx, changePasswd); err != nil {
		userErr, sysErr, statusCode := api.ParseDBError(err)
		sysErr = fmt.Errorf("updating user: %w", err)
		api.HandleErr(w, r, tx, statusCode, userErr, sysErr)
		return
	}

	api.WriteRespAlertObj(w, r, tc.SuccessLevel, "User profile was successfully updated", user)
}

func updateLegacyUser(u *tc.User, tx *sql.Tx, changePassword bool, changeConfirmPasswd bool) error {
	row := tx.QueryRow(replaceCurrentQuery,
		u.AddressLine1,
		u.AddressLine2,
		u.City,
		u.Company,
		u.Country,
		u.Email,
		u.FullName,
		u.GID,
		u.PhoneNumber,
		u.PostalCode,
		u.PublicSSHKey,
		u.Role,
		u.StateOrProvince,
		u.TenantID,
		u.UID,
		u.Username,
		u.ID,
	)
	err := row.Scan(
		&u.AddressLine1,
		&u.AddressLine2,
		&u.City,
		&u.Company,
		&u.Country,
		&u.Email,
		&u.FullName,
		&u.GID,
		&u.ID,
		&u.LastUpdated,
		&u.NewUser,
		&u.PhoneNumber,
		&u.PostalCode,
		&u.PublicSSHKey,
		&u.Role,
		&u.RoleName,
		&u.StateOrProvince,
		&u.Tenant,
		&u.TenantID,
		&u.UID,
		&u.Username,
	)
	if err != nil {
		return err
	}

	if changePassword {
		_, err = tx.Exec(replacePasswordQuery, u.LocalPassword, u.ID)
		if err != nil {
			return fmt.Errorf("resetting password: %w", err)
		}
	}

	if changeConfirmPasswd {
		_, err = tx.Exec(replaceConfirmPasswordQuery, u.ConfirmLocalPassword, u.ID)
		if err != nil {
			return fmt.Errorf("resetting confirm password: %w", err)
		}
	}

	u.LocalPassword = nil
	u.ConfirmLocalPassword = nil
	return nil
}

func updateUser(u *tc.UserV4, tx *sql.Tx, changePassword bool) error {
	row := tx.QueryRow(replaceCurrentV4Query,
		u.AddressLine1,
		u.AddressLine2,
		u.City,
		u.Company,
		u.Country,
		u.Email,
		u.FullName,
		u.GID,
		u.PhoneNumber,
		u.PostalCode,
		u.PublicSSHKey,
		u.Role,
		u.StateOrProvince,
		u.TenantID,
		u.UCDN,
		u.UID,
		u.Username,
		u.ID,
	)

	err := row.Scan(
		&u.AddressLine1,
		&u.AddressLine2,
		&u.ChangeLogCount,
		&u.City,
		&u.Company,
		&u.Country,
		&u.Email,
		&u.FullName,
		&u.GID,
		&u.ID,
		&u.LastAuthenticated,
		&u.LastUpdated,
		&u.NewUser,
		&u.PhoneNumber,
		&u.PostalCode,
		&u.PublicSSHKey,
		&u.RegistrationSent,
		&u.Role,
		&u.StateOrProvince,
		&u.Tenant,
		&u.TenantID,
		&u.UCDN,
		&u.UID,
		&u.Username,
	)
	if err != nil {
		return err
	}

	if changePassword {
		_, err = tx.Exec(replacePasswordQuery, u.LocalPassword, u.ID)
		if err != nil {
			return fmt.Errorf("resetting password: %w", err)
		}
	}

	u.LocalPassword = nil
	return nil
}
