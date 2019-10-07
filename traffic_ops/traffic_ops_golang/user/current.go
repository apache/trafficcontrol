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

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"

	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/tenant"
)

const replaceCurrentQuery = `
UPDATE tm_user
SET address_line1=$1,
    address_line2=$2,
    city=$3,
    company=$4,
    confirm_local_passwd=$5,
    country=$6,
    email=$7,
    full_name=$8,
    gid=$9,
    local_passwd=$10,
    new_user=FALSE,
    phone_number=$11,
    postal_code=$12,
    public_ssh_key=$13,
    state_or_province=$14,
    tenant_id=$15,
    token=NULL,
    uid=$16,
    username=$17
WHERE id=$18
RETURNING address_line1,
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
          ) AS role_name,
          state_or_province,
          (
          	SELECT tenant.name
          	FROM tenant
          	WHERE tenant.id=tm_user.tenant_id
          ) AS tenant,
          tenant_id,
          uid,
          username
`

func Current(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	currentUser, err := getUser(inf.Tx.Tx, inf.User.ID)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting current user: "+err.Error()))
		return
	}
	api.WriteResp(w, r, currentUser)
}

func getUser(tx *sql.Tx, id int) (tc.UserCurrent, error) {
	q := `
SELECT
u.address_line1,
u.address_line2,
u.city,
u.company,
u.country,
u.email,
u.full_name,
u.id,
u.last_updated,
u.local_passwd,
u.new_user,
u.phone_number,
u.postal_code,
u.public_ssh_key,
u.role,
r.name as role_name,
u.state_or_province,
t.name as tenant,
u.tenant_id,
u.username
FROM tm_user as u
LEFT JOIN role as r ON r.id = u.role
LEFT JOIN tenant as t ON t.id = u.tenant_id
WHERE u.id=$1
`
	u := tc.UserCurrent{}
	localPassword := sql.NullString{}
	if err := tx.QueryRow(q, id).Scan(&u.AddressLine1, &u.AddressLine2, &u.City, &u.Company, &u.Country, &u.Email, &u.FullName, &u.ID, &u.LastUpdated, &localPassword, &u.NewUser, &u.PhoneNumber, &u.PostalCode, &u.PublicSSHKey, &u.Role, &u.RoleName, &u.StateOrProvince, &u.Tenant, &u.TenantID, &u.UserName); err != nil {
		return tc.UserCurrent{}, errors.New("querying current user: " + err.Error())
	}
	u.LocalUser = util.BoolPtr(localPassword.Valid)
	return u, nil
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
		userErr = fmt.Errorf("Couldn't parse request: %v", err)
		api.HandleErr(w, r, tx, errCode, userErr, nil)
		return
	}

	user, err := userRequest.User.ValidateAndUnmarshal()
	if err != nil {
		errCode = http.StatusBadRequest
		userErr = fmt.Errorf("Couldn't parse request: %v", err)
		api.HandleErr(w, r, tx, errCode, userErr, nil)
		return
	}

	// obfuscate passwords (ValidateAndUnmarshal checks for equality with ConfirmLocalPassword)
	// TODO: check for valid password via bad password list like Perl did? User creation doesn't...
	if user.LocalPassword != nil && *user.LocalPassword != "" {
		hashPass, err := auth.DerivePassword(*user.LocalPassword)
		if err != nil {
			sysErr = fmt.Errorf("Hashing new password: %v", err)
			errCode = http.StatusInternalServerError
			api.HandleErr(w, r, tx, errCode, nil, sysErr)
			return
		}

		user.LocalPassword = util.StrPtr(hashPass)
		user.ConfirmLocalPassword = util.StrPtr(hashPass)
	}

	if *user.ID != inf.User.ID {
		userErr = errors.New("You cannot change your user ID!")
		errCode = http.StatusBadRequest
		api.HandleErr(w, r, tx, errCode, userErr, nil)
		return
	}

	if *user.Role != inf.User.Role {
		userErr = errors.New("You cannot change your permissions role!")
		errCode = http.StatusBadRequest
		api.HandleErr(w, r, tx, errCode, userErr, nil)
		return
	}

	if ok, err := tenant.IsResourceAuthorizedToUserTx(*user.TenantID, inf.User, tx); err != nil {
		if err == sql.ErrNoRows {
			userErr = errors.New("No such tenant!")
			errCode = http.StatusConflict
		} else {
			sysErr = fmt.Errorf("Checking user %s permissions on tenant #%d: %v", inf.User.UserName, *user.TenantID, err)
			errCode = http.StatusInternalServerError
		}
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	} else if !ok {
		// unlike Perl, this endpoint will not disclose the existence of tenants over which the current
		// user has no permission - in keeping with the behavior of the '/tenants' endpoint.
		userErr = errors.New("No such tenant!")
		errCode = http.StatusConflict
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	if *user.Username != inf.User.UserName {

		if ok, err := dbhelpers.UsernameExists(*user.Username, tx); err != nil {
			sysErr = fmt.Errorf("Checking existence of user %s: %v", *user.Username, err)
			errCode = http.StatusInternalServerError
			api.HandleErr(w, r, tx, errCode, nil, sysErr)
			return
		} else if ok {
			// TODO users are tenanted, so theoretically I should be hiding the existence of the
			// conflicting user - but then how do I tell the client how to fix their request?
			userErr = fmt.Errorf("Username %s already exists!", *user.Username)
			errCode = http.StatusConflict
			api.HandleErr(w, r, tx, errCode, userErr, nil)
			return
		}
	}

	row := tx.QueryRow(replaceCurrentQuery,
		user.AddressLine1,
		user.AddressLine2,
		user.City,
		user.Company,
		user.ConfirmLocalPassword,
		user.Country,
		user.Email,
		user.FullName,
		user.GID,
		user.LocalPassword,
		user.PhoneNumber,
		user.PostalCode,
		user.PublicSSHKey,
		user.StateOrProvince,
		user.TenantID,
		user.UID,
		user.Username,
		inf.User.ID,
	)

	err = row.Scan(&user.AddressLine1,
		&user.AddressLine2,
		&user.City,
		&user.Company,
		&user.Country,
		&user.Email,
		&user.FullName,
		&user.GID,
		&user.ID,
		&user.LastUpdated,
		&user.NewUser,
		&user.PhoneNumber,
		&user.PostalCode,
		&user.PublicSSHKey,
		&user.Role,
		&user.RoleName,
		&user.StateOrProvince,
		&user.Tenant,
		&user.TenantID,
		&user.UID,
		&user.Username,
	)
	if err != nil {
		errCode = http.StatusInternalServerError
		sysErr = fmt.Errorf("Updating user: %v", err)
		api.HandleErr(w, r, tx, errCode, nil, sysErr)
		return
	}

	// Hide the password fields
	user.LocalPassword = nil
	user.ConfirmLocalPassword = nil

	resp := struct {
		tc.Alerts
		Response tc.User `json:"response"`
	}{
		tc.CreateAlerts(tc.SuccessLevel, "User profile was successfully updated"),
		user,
	}

	respBts, err := json.Marshal(resp)
	if err != nil {
		errCode = http.StatusInternalServerError
		sysErr = fmt.Errorf("Marshalling response: %v", err)
		api.HandleErr(w, r, tx, errCode, nil, sysErr)
		return
	}
	w.Header().Set(tc.ContentType, tc.ApplicationJson)
	w.Write(append(respBts, '\n'))
}
