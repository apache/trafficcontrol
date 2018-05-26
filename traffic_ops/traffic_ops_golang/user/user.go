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
	"net/http"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/tenant"

	"github.com/jmoiron/sqlx"
)

func Get(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, intParams, userErr, sysErr, errCode := api.AllParams(r, nil, []string{"tenant"})
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, errCode, userErr, sysErr)
			return
		}
		if tenantID, ok := intParams["tenant"]; ok {
			api.RespWriter(w, r)(getUsersByTenantID(db, tenantID))
			return
		}
		api.RespWriter(w, r)(getUsers(db))
	}
}

func GetID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, intParams, userErr, sysErr, errCode := api.AllParams(r, []string{"id"}, []string{"id"})
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, errCode, userErr, sysErr)
			return
		}
		user, ok, err := getUserByID(db, intParams["id"])
		if err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("getting user "+strconv.Itoa(intParams["id"])+": "+err.Error()))
			return
		}
		if !ok {
			api.HandleErr(w, r, http.StatusNotFound, nil, nil)
			return
		}
		api.WriteResp(w, r, []tc.APIUser{user})
	}
}

func Post(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := tc.APIUserPost{}
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			api.HandleErr(w, r, http.StatusBadRequest, errors.New("malformed JSON"), nil)
			return
		}
		if err := validatePost(u); err != nil {
			api.HandleErr(w, r, http.StatusBadRequest, errors.New("validation error: "+err.Error()), nil)
			return
		}
		api.RespWithAlertWriter(w, r, tc.SuccessLevel, "User creation was successful.")(createUser(db, u))
	}
}

func Put(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.GetCurrentUser(r.Context())
		if err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("getting current user: "+err.Error()))
			return
		}
		_, intParams, userErr, sysErr, errCode := api.AllParams(r, []string{"id"}, []string{"id"})
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, errCode, userErr, sysErr)
			return
		}
		existingTenantID, ok, err := getUserTenantIDByID(db.DB, intParams["id"])
		if err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("getting user "+strconv.Itoa(intParams["id"])+": "+err.Error()))
			return
		}
		if !ok {
			api.HandleErr(w, r, http.StatusNotFound, nil, nil)
			return
		}
		u := tc.APIUserPost{}
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			log.Errorln("user put: malformed JSON: " + err.Error())
			api.HandleErr(w, r, http.StatusBadRequest, errors.New("malformed JSON"), nil)
			return
		}
		userID := intParams["id"]
		u.ID = &userID
		if err := validatePut(u); err != nil {
			api.HandleErr(w, r, http.StatusBadRequest, errors.New("validation error: "+err.Error()), nil)
			return
		}
		authorized, err := isTenantAuthorized(user, db, intParams["id"], existingTenantID, u.TenantID)
		if err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("checking user tenancy for "+strconv.Itoa(intParams["id"])+": "+err.Error()))
			return
		}
		if !authorized {
			log.Errorln("user put: tenant unauthorized!")
			api.HandleErr(w, r, http.StatusUnauthorized, nil, nil)
			return
		}
		api.RespWithAlertWriter(w, r, tc.SuccessLevel, "User update successful.")(updateUser(db.DB, u))
	}
}

// isTenantAuthorized returns whether the current user is authorized to modify the tenant of the given user, and any new tenant. The tenantID may be null, if the existing tenant is not being changed.
func isTenantAuthorized(user *auth.CurrentUser, db *sqlx.DB, userID int, oldTenantID *int, newTenantID *int) (bool, error) {
	if oldTenantID != nil {
		authorized, err := tenant.IsResourceAuthorizedToUser(*oldTenantID, user, db)
		if err != nil {
			return false, errors.New("checking authorization for existing user ID: " + err.Error())
		}
		if !authorized {
			return false, nil
		}
	}
	if newTenantID != nil && (oldTenantID == nil || *newTenantID != *oldTenantID) {
		authorized, err := tenant.IsResourceAuthorizedToUser(*newTenantID, user, db)
		if err != nil {
			return false, errors.New("checking authorization for new user ID: " + err.Error())
		}
		if !authorized {
			return false, nil
		}
	}
	return true, nil
}

func validatePost(u tc.APIUserPost) error {
	errs := []string{}
	if u.ConfirmLocalPassword == nil || *u.ConfirmLocalPassword == "" {
		errs = append(errs, "confirmLocalPassword must be set")
	}
	if u.Email == nil || *u.Email == "" {
		errs = append(errs, "email must be set")
	}
	if u.FullName == nil || *u.FullName == "" {
		errs = append(errs, "fullName must be set")
	}
	if u.LocalPassword == nil || *u.LocalPassword == "" {
		errs = append(errs, "localPassword must be set")
	}
	if u.Role == nil {
		errs = append(errs, "role must be set")
	}
	if u.UserName == nil || *u.UserName == "" {
		errs = append(errs, "username must be set")
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}
	return nil
}

func validatePut(u tc.APIUserPost) error {
	errs := []string{}
	if u.Email == nil || *u.Email == "" {
		errs = append(errs, "email must be set")
	}
	if u.FullName == nil || *u.FullName == "" {
		errs = append(errs, "fullName must be set")
	}
	if u.Role == nil {
		errs = append(errs, "role must be set")
	}
	if u.UserName == nil || *u.UserName == "" {
		errs = append(errs, "username must be set")
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}
	return nil
}

func getUsers(db *sql.DB) ([]tc.APIUser, error) {
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
u.new_user,
u.phone_number,
u.postal_code,
u.public_ssh_key,
u.registration_sent,
u.role,
r.name,
u.state_or_province,
t.name,
u.tenant_id,
u.uid,
u.username
FROM tm_user as u
JOIN tenant as t on t.id = u.tenant_id
JOIN role as r on r.id = u.role
`
	rows, err := db.Query(q)
	if err != nil {
		return nil, errors.New("querying users: " + err.Error())
	}
	defer rows.Close()
	users := []tc.APIUser{}
	for rows.Next() {
		u := tc.APIUser{}
		if err := rows.Scan(&u.AddressLine1, &u.AddressLine2, &u.City, &u.Company, &u.Country, &u.Email, &u.FullName, &u.GID, &u.ID, &u.LastUpdated, &u.NewUser, &u.PhoneNumber, &u.PostalCode, &u.PublicSSHKey, &u.RegistrationSent, &u.Role, &u.RoleName, &u.StateOrProvince, &u.Tenant, &u.TenantID, &u.UID, &u.UserName); err != nil {
			return nil, errors.New("scanning users: " + err.Error())
		}
		users = append(users, u)
	}
	return users, nil
}

func getUsersByTenantID(db *sql.DB, tenantID int) ([]tc.APIUser, error) {
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
u.new_user,
u.phone_number,
u.postal_code,
u.public_ssh_key,
u.registration_sent,
u.role,
r.name,
u.state_or_province,
t.name,
u.tenant_id,
u.uid,
u.username
FROM tm_user as u
JOIN tenant as t on t.id = u.tenant_id
JOIN role as r on r.id = u.role
WHERE u.tenant_id = $1
`
	rows, err := db.Query(q, tenantID)
	if err != nil {
		return nil, errors.New("querying users: " + err.Error())
	}
	defer rows.Close()
	users := []tc.APIUser{}
	for rows.Next() {
		u := tc.APIUser{}
		if err := rows.Scan(&u.AddressLine1, &u.AddressLine2, &u.City, &u.Company, &u.Country, &u.Email, &u.FullName, &u.GID, &u.ID, &u.LastUpdated, &u.NewUser, &u.PhoneNumber, &u.PostalCode, &u.PublicSSHKey, &u.RegistrationSent, &u.Role, &u.RoleName, &u.StateOrProvince, &u.Tenant, &u.TenantID, &u.UID, &u.UserName); err != nil {
			return nil, errors.New("scanning users: " + err.Error())
		}
		users = append(users, u)
	}
	return users, nil
}

func getUserByID(db *sql.DB, id int) (tc.APIUser, bool, error) {
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
u.new_user,
u.phone_number,
u.postal_code,
u.public_ssh_key,
u.registration_sent,
u.role,
r.name,
u.state_or_province,
t.name,
u.tenant_id,
u.uid,
u.username
FROM tm_user as u
JOIN tenant as t on t.id = u.tenant_id
JOIN role as r on r.id = u.role
WHERE u.id = $1
`
	u := tc.APIUser{}
	if err := db.QueryRow(q, id).Scan(&u.AddressLine1, &u.AddressLine2, &u.City, &u.Company, &u.Country, &u.Email, &u.FullName, &u.GID, &u.ID, &u.LastUpdated, &u.NewUser, &u.PhoneNumber, &u.PostalCode, &u.PublicSSHKey, &u.RegistrationSent, &u.Role, &u.RoleName, &u.StateOrProvince, &u.Tenant, &u.TenantID, &u.UID, &u.UserName); err != nil {
		if err == sql.ErrNoRows {
			return tc.APIUser{}, false, nil
		}
		return tc.APIUser{}, true, errors.New("querying user: " + err.Error())
	}
	return u, true, nil
}

func createUser(db *sql.DB, u tc.APIUserPost) (tc.APIUser, error) {
	cols, params, vals, err := dbhelpers.BuildInsertColumns(u)
	if err != nil {
		return tc.APIUser{}, errors.New("building insert query: " + err.Error())
	}
	q := `INSERT INTO tm_user (` + cols + `) VALUES (` + params + `) RETURNING id, last_updated, (select r2.name from role as r2 where id = tm_user.role)`
	if err := db.QueryRow(q, vals...).Scan(&u.ID, &u.LastUpdated, &u.RoleName); err != nil {
		return tc.APIUser{}, errors.New("inserting user: " + err.Error())
	}
	return u.APIUser, nil
}

func updateUser(db *sql.DB, u tc.APIUserPost) (tc.APIUser, error) {
	q := `
UPDATE tm_user SET
address_line1=$1,
address_line2=$2,
city=$3,
company=$4,
country=$5,
email=$6,
full_name=$7,
gid=$8,
new_user=$9,
phone_number=$10,
postal_code=$11,
public_ssh_key=$12,
registration_sent=$13,
role=$14,
state_or_province=$15,
tenant_id=$16,
uid=$17,
username=$18
`
	nextParam := `$19`
	if u.LocalPassword != nil {
		q += `,local_passwd=$19`
		nextParam = `$20`
	}
	if u.ConfirmLocalPassword != nil {
		if u.LocalPassword == nil {
			q += `,confirm_local_passwd=$19`
			nextParam = `$20`
		} else {
			q += `,confirm_local_passwd=$20`
			nextParam = `$21`
		}
	}

	q += `
WHERE id=` + nextParam + `
RETURNING
last_updated,
(select t2.name from tenant as t2 where id = tm_user.tenant_id),
(select r2.name from role as r2 where id = tm_user.role)
`
	vals := []interface{}{u.AddressLine1, u.AddressLine2, u.City, u.Company, u.Country, u.Email, u.FullName, u.GID, u.NewUser, u.PhoneNumber, u.PostalCode, u.PublicSSHKey, u.RegistrationSent, u.Role, u.StateOrProvince, u.TenantID, u.UID, u.UserName}
	if u.LocalPassword != nil {
		vals = append(vals, u.LocalPassword)
	}
	if u.ConfirmLocalPassword != nil {
		vals = append(vals, u.ConfirmLocalPassword)
	}
	vals = append(vals, u.ID)

	if err := db.QueryRow(q, vals...).Scan(&u.LastUpdated, &u.Tenant, &u.RoleName); err != nil {
		return tc.APIUser{}, errors.New("updating user: " + err.Error())
	}
	return u.APIUser, nil
}
