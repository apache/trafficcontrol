package role

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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/routing/middleware"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type TORole struct {
	api.APIInfoImpl `json:"-"`
	tc.Role
	LastUpdated    *tc.TimeNoMod   `json:"-"`
	PQCapabilities *pq.StringArray `json:"-" db:"capabilities"`
}

func (v *TORole) GetLastUpdated() (*time.Time, bool, error) {
	return api.GetLastUpdated(v.APIInfo().Tx, *v.ID, "role")
}

func (v *TORole) SelectMaxLastUpdatedQuery(where, orderBy, pagination, tableName string) string {
	return `SELECT max(t) from (
		SELECT max(last_updated) as t from ` + tableName + ` r ` + where + orderBy + pagination +
		` UNION ALL
	select max(last_updated) as t from last_deleted l where l.table_name='` + tableName + `') as res`
}

func (v *TORole) SetLastUpdated(t tc.TimeNoMod) { v.LastUpdated = &t }
func (v *TORole) InsertQuery() string           { return insertQuery() }
func (v *TORole) NewReadObj() interface{}       { return &TORole{} }
func (v *TORole) SelectQuery() string           { return selectQuery() }
func (v *TORole) ParamColumns() map[string]dbhelpers.WhereColumnInfo {
	return map[string]dbhelpers.WhereColumnInfo{
		"name":      dbhelpers.WhereColumnInfo{Column: "name"},
		"id":        dbhelpers.WhereColumnInfo{Column: "id", Checker: api.IsInt},
		"privLevel": dbhelpers.WhereColumnInfo{Column: "priv_level", Checker: api.IsInt}}
}
func (v *TORole) UpdateQuery() string { return updateQuery() }
func (v *TORole) DeleteQuery() string { return deleteQuery() }

func (role TORole) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{Field: "id", Func: api.GetIntKey}}
}

//Implementation of the Identifier, Validator interface functions
func (role TORole) GetKeys() (map[string]interface{}, bool) {
	if role.ID == nil {
		return map[string]interface{}{"id": 0}, false
	}
	return map[string]interface{}{"id": *role.ID}, true
}

func (role TORole) GetAuditName() string {
	if role.Name != nil {
		return *role.Name
	}
	if role.ID != nil {
		return strconv.Itoa(*role.ID)
	}
	return "0"
}

func (role TORole) GetType() string {
	return "role"
}

func (role *TORole) SetKeys(keys map[string]interface{}) {
	i, _ := keys["id"].(int) //this utilizes the non panicking type assertion, if the thrown away ok variable is false i will be the zero of the type, 0 here.
	role.ID = &i
}

// Validate fulfills the api.Validator interface
func (role TORole) Validate() error {
	errs := validation.Errors{
		"name":        validation.Validate(role.Name, validation.Required),
		"description": validation.Validate(role.Description, validation.Required),
		"privLevel":   validation.Validate(role.PrivLevel, validation.Required)}

	errsToReturn := tovalidate.ToErrors(errs)
	checkCaps := `SELECT cap FROM UNNEST($1::text[]) AS cap WHERE NOT cap =  ANY(ARRAY(SELECT c.name FROM capability AS c WHERE c.name = ANY($1)))`
	var badCaps []string
	if role.ReqInfo.Tx != nil {
		err := role.ReqInfo.Tx.Select(&badCaps, checkCaps, pq.Array(role.Capabilities))
		if err != nil {
			log.Errorf("got error from selecting bad capabilities: %v", err)
			return tc.DBError
		}
		if len(badCaps) > 0 {
			errsToReturn = append(errsToReturn, fmt.Errorf("can not add non-existent capabilities: %v", badCaps))
		}
	}
	return util.JoinErrs(errsToReturn)
}

func (role *TORole) Create() (error, error, int) {
	if *role.PrivLevel > role.ReqInfo.User.PrivLevel {
		return errors.New("can not create a role with a higher priv level than your own"), nil, http.StatusBadRequest
	}

	userErr, sysErr, errCode := api.GenericCreate(role)
	if userErr != nil || sysErr != nil {
		return userErr, sysErr, errCode
	}

	//after we have role ID we can associate the capabilities:
	if role.Capabilities != nil && len(*role.Capabilities) > 0 {
		userErr, sysErr, errCode = role.createRoleCapabilityAssociations(role.ReqInfo.Tx)
		if userErr != nil || sysErr != nil {
			return userErr, sysErr, errCode
		}
	}
	return nil, nil, http.StatusOK
}

func (role *TORole) createRoleCapabilityAssociations(tx *sqlx.Tx) (error, error, int) {
	result, err := tx.Exec(associateCapabilities(), role.ID, pq.Array(role.Capabilities))
	if err != nil {
		return nil, errors.New("creating role capabilities: " + err.Error()), http.StatusInternalServerError
	}

	if rows, err := result.RowsAffected(); err != nil {
		log.Errorf("could not check result after inserting role_capability relations: %v", err)
	} else if expected := len(*role.Capabilities); int(rows) != expected {
		log.Errorf("wrong number of role_capability rows created: %d expected: %d", rows, expected)
	}
	return nil, nil, http.StatusOK
}

func (role *TORole) deleteRoleCapabilityAssociations(tx *sqlx.Tx) (error, error, int) {
	result, err := tx.Exec(deleteAssociatedCapabilities(), role.ID)
	if err != nil {
		return nil, errors.New("deleting role capabilities: " + err.Error()), http.StatusInternalServerError
	}

	if _, err = result.RowsAffected(); err != nil {
		log.Errorf("could not check result after inserting role_capability relations: %v", err)
	}
	// TODO verify expected row count shouldn't be checked?
	return nil, nil, http.StatusOK
}

func (role *TORole) Read(h http.Header, useIMS bool) ([]interface{}, error, error, int, *time.Time) {
	api.DefaultSort(role.APIInfo(), "name")
	vals, userErr, sysErr, errCode, maxTime := api.GenericRead(h, role, useIMS)
	if errCode == http.StatusNotModified {
		return []interface{}{}, nil, nil, errCode, maxTime
	}
	if userErr != nil || sysErr != nil {
		return nil, userErr, sysErr, errCode, maxTime
	}

	returnable := []interface{}{}
	for _, val := range vals {
		rl := val.(*TORole)
		caps := ([]string)(*rl.PQCapabilities)
		rl.Capabilities = &caps
		returnable = append(returnable, rl)
	}
	return returnable, nil, nil, http.StatusOK, maxTime
}

func (role *TORole) Update(h http.Header) (error, error, int) {
	if *role.PrivLevel > role.ReqInfo.User.PrivLevel {
		return errors.New("can not create a role with a higher priv level than your own"), nil, http.StatusForbidden
	}
	userErr, sysErr, errCode := api.GenericUpdate(h, role)
	if userErr != nil || sysErr != nil {
		return userErr, sysErr, errCode
	}

	// TODO cascade delete, to automatically do this in SQL?
	if role.Capabilities != nil && *role.Capabilities != nil {
		userErr, sysErr, errCode = role.deleteRoleCapabilityAssociations(role.ReqInfo.Tx)
		if userErr != nil || sysErr != nil {
			return userErr, sysErr, errCode
		}
		return role.createRoleCapabilityAssociations(role.ReqInfo.Tx)
	}
	return nil, nil, http.StatusOK
}

func (role *TORole) Delete() (error, error, int) {
	assignedUsers := 0
	if err := role.ReqInfo.Tx.Get(&assignedUsers, "SELECT COUNT(id) FROM tm_user WHERE role=$1", role.ID); err != nil {
		return nil, errors.New("role delete counting assigned users: " + err.Error()), http.StatusInternalServerError
	} else if assignedUsers != 0 {
		return fmt.Errorf("can not delete a role with %d assigned users", assignedUsers), nil, http.StatusBadRequest
	}

	userErr, sysErr, errCode := api.GenericDelete(role)
	if userErr != nil || sysErr != nil {
		return userErr, sysErr, errCode
	}
	return role.deleteRoleCapabilityAssociations(role.ReqInfo.Tx)
}

func selectQuery() string {
	return `SELECT
id,
name,
description,
priv_level,
ARRAY(SELECT rc.cap_name FROM role_capability AS rc WHERE rc.role_id=id) AS capabilities
FROM role`
}

func updateQuery() string {
	return `UPDATE
role SET
name=:name,
description=:description
WHERE id=:id RETURNING last_updated`
}

func updateRoleQuery() string {
	return `UPDATE
role SET
description=$1
WHERE name=$2 RETURNING last_updated`
}

func deleteAssociatedCapabilities() string {
	return `DELETE FROM role_capability
WHERE role_id=$1`
}

func associateCapabilities() string {
	return `INSERT INTO role_capability (
role_id,
cap_name) WITH
	q1 AS ( SELECT * FROM (VALUES ($1::bigint)) AS role_id ),
	q2 AS (SELECT UNNEST($2::text[]))
	SELECT * FROM q1,q2`
}

func insertQuery() string {
	return `INSERT INTO role (
name,
description,
priv_level
) VALUES (
:name,
:description,
:priv_level
)
RETURNING id, last_updated`
}

func deleteQuery() string {
	return `DELETE FROM role WHERE id = :id`
}

func Update(w http.ResponseWriter, r *http.Request) {
	var roleID int
	var roleName string
	var roleDesc string
	var privLevel int
	var roleCapabilities *[]string
	var roleV50 tc.RoleV50
	var role tc.Role
	var ok bool
	var err error

	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	version := inf.Version
	if version == nil {
		middleware.NotImplementedHandler().ServeHTTP(w, r)
		return
	}

	if version.Major >= 5 {
		if err := json.NewDecoder(r.Body).Decode(&roleV50); err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, err, nil)
			return
		}
		if err := Validate(inf.Tx, role, roleV50, version); err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, err, nil)
			return
		}
		roleDesc = *roleV50.Description
		roleCapabilities = &roleV50.Permissions
		if roleName, ok = inf.Params["name"]; !ok {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, err, nil)
			return
		}
		roleID, err = dbhelpers.GetRoleIDFromName(inf.Tx.Tx, roleName)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("no ID exists for the supplied role name"), nil)
			return
		}

		existingLastUpdated, found, err := api.GetLastUpdated(inf.Tx, roleID, "role")
		if err == nil && found == false {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, errors.New("no role found with that ID"), nil)
			return
		}
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
			return
		}

		if !api.IsUnmodified(r.Header, *existingLastUpdated) {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusPreconditionFailed, api.ResourceModifiedError, nil)
			return
		}
		rows, err := inf.Tx.Tx.Query(updateRoleQuery(), roleDesc, roleName)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("updating role: "+err.Error()))
			return
		}
		defer rows.Close()
		if !rows.Next() {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, errors.New("no role found with this ID"), nil)
			return
		}
		lastUpdated := tc.TimeNoMod{}
		if err := rows.Scan(&lastUpdated); err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("scanning lastUpdated from role update: "+err.Error()))
			return
		}
		roleV50.LastUpdated = &lastUpdated
	} else {
		if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, err, nil)
			return
		}
		if err := Validate(inf.Tx, role, roleV50, version); err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("name and/or description and/or privLevel can not be empty"), nil)
			return
		}
		roleName = *role.Name
		roleDesc = *role.Description
		privLevel = *role.PrivLevel
		roleCapabilities = role.Capabilities
		if roleID, ok = inf.IntParams["id"]; !ok {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("must supply a role ID to delete"), nil)
			return
		}
		if privLevel > inf.User.PrivLevel {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusForbidden, errors.New("can not create a role with a higher priv level than your own"), nil)
			return
		}
		existingLastUpdated, found, err := api.GetLastUpdated(inf.Tx, roleID, "role")
		if err == nil && found == false {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, errors.New("no role found with that ID"), nil)
			return
		}
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
			return
		}

		if !api.IsUnmodified(r.Header, *existingLastUpdated) {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusPreconditionFailed, api.ResourceModifiedError, nil)
			return
		}

		rows, err := inf.Tx.Tx.Query(updateQuery(), role)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("updating role: "+err.Error()))
			return
		}
		rows.Close()
	}

	if roleCapabilities != nil && *roleCapabilities != nil {
		userErr, sysErr, errCode = deleteRoleCapabilityAssociations(inf.Tx, roleID)
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
			return
		}
		userErr, sysErr, errCode = createRoleCapabilityAssociations(inf.Tx, roleID, roleCapabilities)
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
			return
		}
	}
	alerts := tc.CreateAlerts(tc.SuccessLevel, "role was updated.")
	var roleResponse interface{}
	if version.Major >= 5 {
		var capabilities []string
		if roleCapabilities != nil {
			capabilities = *roleCapabilities
		}
		roleResponse = tc.RoleV50{
			Name:        util.StrPtr(roleName),
			Permissions: capabilities,
			Description: util.StrPtr(roleDesc),
		}
	} else {
		roleResponse = tc.Role{
			RoleV11: tc.RoleV11{
				ID:          util.IntPtr(roleID),
				Name:        util.StrPtr(roleName),
				Description: util.StrPtr(roleDesc),
				PrivLevel:   util.IntPtr(privLevel),
			},
			Capabilities: roleCapabilities,
		}
	}
	api.WriteAlertsObj(w, r, http.StatusOK, alerts, roleResponse)
}

func deleteRoleQuery() string {
	return `DELETE FROM role WHERE id = $1`
}

func readQuery() string {
	return `SELECT
id,
name,
description,
priv_level,
last_updated,
ARRAY(SELECT rc.cap_name FROM role_capability AS rc WHERE rc.role_id=id) AS permissions
FROM role`
}

func createQuery() string {
	return `INSERT INTO role (
name,
description,
priv_level
) VALUES (
$1,
$2,
$3
)
RETURNING id, last_updated`
}

func createRoleCapabilityAssociations(tx *sqlx.Tx, roleID int, permissions *[]string) (error, error, int) {
	result, err := tx.Exec(associateCapabilities(), roleID, pq.Array(permissions))
	if err != nil {
		return nil, errors.New("creating role capabilities: " + err.Error()), http.StatusInternalServerError
	}

	if rows, err := result.RowsAffected(); err != nil {
		log.Errorf("could not check result after inserting role_capability relations: %v", err)
	} else if expected := len(*permissions); int(rows) != expected {
		log.Errorf("wrong number of role_capability rows created: %d expected: %d", rows, expected)
	}
	return nil, nil, http.StatusOK
}

func Validate(tx *sqlx.Tx, role tc.Role, roleV50 tc.RoleV50, version *api.Version) error {
	var capabilities *[]string
	errs := make(map[string]error)
	if version.Major >= 5 {
		errs = validation.Errors{
			"name":        validation.Validate(roleV50.Name, validation.Required),
			"description": validation.Validate(roleV50.Description, validation.Required),
		}
		capabilities = &roleV50.Permissions
	} else {
		errs = validation.Errors{
			"name":        validation.Validate(role.Name, validation.Required),
			"description": validation.Validate(role.Description, validation.Required),
			"privLevel":   validation.Validate(role.PrivLevel, validation.Required),
		}
		capabilities = role.Capabilities
	}

	errsToReturn := tovalidate.ToErrors(errs)
	checkCaps := `SELECT cap FROM UNNEST($1::text[]) AS cap WHERE NOT cap =  ANY(ARRAY(SELECT c.name FROM capability AS c WHERE c.name = ANY($1)))`
	var badCaps []string
	if tx != nil {
		err := tx.Select(&badCaps, checkCaps, pq.Array(capabilities))
		if err != nil {
			log.Errorf("got error from selecting bad capabilities: %v", err)
			return err
		}
		if len(badCaps) > 0 {
			errsToReturn = append(errsToReturn, fmt.Errorf("can not add non-existent capabilities: %v", badCaps))
		}
	}
	return util.JoinErrs(errsToReturn)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	var roleName string
	var ok bool
	var roleID int
	var err error
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	version := inf.Version
	if version == nil {
		middleware.NotImplementedHandler().ServeHTTP(w, r)
		return
	}

	if version.Major >= 5 {
		if roleName, ok = inf.Params["name"]; !ok {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("must supply a role name to delete"), nil)
			return
		}
		roleID, err = dbhelpers.GetRoleIDFromName(inf.Tx.Tx, roleName)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("no ID exists for the supplied role name"), nil)
			return
		}
	} else {
		if roleID, ok = inf.IntParams["id"]; !ok {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("must supply a role ID to delete"), nil)
			return
		}
	}

	assignedUsers := 0
	if err := inf.Tx.Get(&assignedUsers, "SELECT COUNT(id) FROM tm_user WHERE role=$1", roleID); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("role delete counting assigned users: "+err.Error()))
		return
	} else if assignedUsers != 0 {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, fmt.Errorf("can not delete a role with %d assigned users", assignedUsers), nil)
		return
	}

	rows, err := inf.Tx.Tx.Query(deleteRoleQuery(), roleID)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("deleting role: "+err.Error()))
		return
	}
	rows.Close()
	userErr, sysErr, errCode = deleteRoleCapabilityAssociations(inf.Tx, roleID)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	alerts := tc.CreateAlerts(tc.SuccessLevel, "role was deleted.")
	api.WriteAlerts(w, r, http.StatusOK, alerts)
}

func deleteRoleCapabilityAssociations(tx *sqlx.Tx, roleID int) (error, error, int) {
	result, err := tx.Exec(deleteAssociatedCapabilities(), roleID)
	if err != nil {
		return nil, errors.New("deleting role capabilities: " + err.Error()), http.StatusInternalServerError
	}
	if _, err = result.RowsAffected(); err != nil {
		log.Errorf("could not check result after inserting role_capability relations: %v", err)
	}
	// TODO verify expected row count shouldn't be checked?
	return nil, nil, http.StatusOK
}

func Create(w http.ResponseWriter, r *http.Request) {
	var roleID int
	var roleName string
	var roleDesc string
	var privLevel int
	var roleCapabilities *[]string
	var roleV50 tc.RoleV50
	var role tc.Role

	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	version := inf.Version
	if version == nil {
		middleware.NotImplementedHandler().ServeHTTP(w, r)
		return
	}

	if version.Major >= 5 {
		if err := json.NewDecoder(r.Body).Decode(&roleV50); err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, err, nil)
			return
		}
		if err := Validate(inf.Tx, role, roleV50, version); err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, err, nil)
			return
		}
		roleName = *roleV50.Name
		roleDesc = *roleV50.Description
		privLevel = inf.User.PrivLevel
		roleCapabilities = &roleV50.Permissions
	} else {
		if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, err, nil)
			return
		}
		if err := Validate(inf.Tx, role, roleV50, version); err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, err, nil)
			return
		}
		roleName = *role.Name
		roleDesc = *role.Description
		privLevel = *role.PrivLevel
		roleCapabilities = role.Capabilities
	}

	rows, err := inf.Tx.Tx.Query(createQuery(), roleName, roleDesc, privLevel)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("creating role: "+err.Error()))
		return
	}
	defer rows.Close()

	for rows.Next() {
		var throwaway interface{}
		if err := rows.Scan(&roleID, &throwaway); err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("role create: scanning role ID: "+err.Error()))
			return
		}
	}

	if roleCapabilities != nil && len(*roleCapabilities) > 0 {
		userErr, sysErr, errCode = createRoleCapabilityAssociations(inf.Tx, roleID, roleCapabilities)
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
			return
		}
	}
	alerts := tc.CreateAlerts(tc.SuccessLevel, "role was created.")
	var roleResponse interface{}
	if version.Major >= 5 {
		var capabilities []string
		if roleCapabilities != nil {
			capabilities = *roleCapabilities
		}
		roleResponse = tc.RoleV50{
			Name:        util.StrPtr(roleName),
			Permissions: capabilities,
			Description: util.StrPtr(roleDesc),
		}
	} else {
		roleResponse = tc.Role{
			RoleV11: tc.RoleV11{
				ID:          util.IntPtr(roleID),
				Name:        util.StrPtr(roleName),
				Description: util.StrPtr(roleDesc),
				PrivLevel:   util.IntPtr(privLevel),
			},
			Capabilities: roleCapabilities,
		}
	}
	api.WriteAlertsObj(w, r, http.StatusOK, alerts, roleResponse)
}

func Get(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	version := inf.Version
	if version == nil {
		middleware.NotImplementedHandler().ServeHTTP(w, r)
		return
	}

	params := make(map[string]dbhelpers.WhereColumnInfo, 0)
	if version.Major >= 5 {
		params["name"] = dbhelpers.WhereColumnInfo{Column: "name"}
	} else {
		params["name"] = dbhelpers.WhereColumnInfo{Column: "name"}
		params["id"] = dbhelpers.WhereColumnInfo{Column: "id", Checker: api.IsInt}
		params["privLevel"] = dbhelpers.WhereColumnInfo{Column: "priv_level", Checker: api.IsInt}
	}

	if _, ok := inf.Params["orderby"]; !ok {
		inf.Params["orderby"] = "name"
	}
	where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(inf.Params, params)
	if len(errs) != 0 {
		api.HandleErr(w, r, tx, http.StatusBadRequest, util.JoinErrs(errs), nil)
		return
	}
	if version.Major >= 5 {
		if perm, ok := inf.Params["can"]; ok {
			queryValues["can"] = perm
			where = dbhelpers.AppendWhere(where, "permissions @> :can")
		}
	}

	query := readQuery() + where + orderBy + pagination

	rows, err := inf.Tx.NamedQuery(query, queryValues)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("querying Roles: %w", err))
		return
	}
	defer log.Close(rows, "reading in Roles from the database")

	var response struct {
		Response interface{}
	}

	var roleV50 tc.RoleV50
	rolesV50 := []tc.RoleV50{}

	var role tc.Role
	roles := []tc.Role{}

	if version.Major >= 5 {
		for rows.Next() {
			throwAway := new(interface{})
			if err = rows.Scan(throwAway, &roleV50.Name, &roleV50.Description, throwAway, &roleV50.LastUpdated, pq.Array(&roleV50.Permissions)); err != nil {
				api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("scanning RoleV50 row: %w", err))
				return
			}
			rolesV50 = append(rolesV50, roleV50)
		}
		response.Response = rolesV50
	} else {
		for rows.Next() {
			throwAway := new(interface{})
			var capabilities []string
			if err = rows.Scan(&role.ID, &role.Name, &role.Description, &role.PrivLevel, throwAway, pq.Array(&capabilities)); err != nil {
				api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("scanning RoleV11 row: %w", err))
				return
			}
			role.Capabilities = &capabilities
			roles = append(roles, role)
		}
		response.Response = roles
	}

	if inf.UseIMS() {
		maxTime, err := getMaxLastUpdated(where, queryValues, inf.Tx)
		if err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
			return
		}
		w.Header().Add(rfc.LastModified, maxTime.Format(rfc.LastModifiedFormat))
	}
	api.WriteResp(w, r, response)
}

func getMaxLastUpdated(where string, queryValues map[string]interface{}, tx *sqlx.Tx) (time.Time, error) {
	query := selectMaxLastUpdatedQuery(where)
	var t time.Time
	rows, err := tx.NamedQuery(query, queryValues)
	if err != nil {
		return t, fmt.Errorf("query for max user last updated time: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&t); err != nil {
			return t, fmt.Errorf("scanning user max last updated time: %w", err)
		}
	}
	return t, nil
}

func selectMaxLastUpdatedQuery(where string) string {
	return `
		SELECT max(t)
		FROM (
			SELECT max(r.last_updated) AS t
			FROM role r ` + where + `
			UNION ALL
			SELECT max(l.last_updated)
			FROM last_deleted l
			WHERE l.tablename=role OR l.tablename=role_capability
			UNION ALL
			SELECT max(rc.last_updated)
			FROM role_capability rc
			WHERE rc.role_id = r.id
		) AS res`
}
