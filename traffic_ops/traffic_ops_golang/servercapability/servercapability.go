package servercapability

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
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/dbhelpers"

	validation "github.com/go-ozzo/ozzo-validation"
)

type TOServerCapability struct {
	api.APIInfoImpl `json:"-"`
	RequestedName   string `json:"-"`
	tc.ServerCapability
}

func (v *TOServerCapability) SetLastUpdated(t tc.TimeNoMod) { v.LastUpdated = &t }
func (v *TOServerCapability) NewReadObj() interface{}       { return &tc.ServerCapability{} }
func (v *TOServerCapability) InsertQuery() string {
	return `INSERT INTO server_capability (
  name
)
VALUES (
  :name
)
RETURNING last_updated
`
}

func (v *TOServerCapability) SelectQuery() string {
	return `SELECT
  name,
  last_updated
FROM
  server_capability sc
`
}

func (v *TOServerCapability) updateQuery() string {
	return `UPDATE server_capability sc SET
	name = $1
WHERE sc.name = $2
RETURNING sc.name, sc.last_updated
`
}

func (v *TOServerCapability) DeleteQuery() string {
	return `DELETE FROM server_capability 
WHERE name=:name
`
}

func (v *TOServerCapability) ParamColumns() map[string]dbhelpers.WhereColumnInfo {
	return map[string]dbhelpers.WhereColumnInfo{
		"name": {Column: "sc.name"},
	}
}

func (v TOServerCapability) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{Field: "name", Func: api.GetStringKey}}
}

// Implementation of the Identifier, Validator interface functions
func (v TOServerCapability) GetKeys() (map[string]interface{}, bool) {
	return map[string]interface{}{"name": v.Name}, true
}

func (v *TOServerCapability) SetKeys(keys map[string]interface{}) {
	v.RequestedName = v.Name
	v.Name, _ = keys["name"].(string)
}

func (v *TOServerCapability) GetAuditName() string {
	return v.Name
}

func (v *TOServerCapability) GetType() string {
	return "server capability"
}

func (v *TOServerCapability) Validate() (error, error) {
	rule := validation.NewStringRule(tovalidate.IsAlphanumericUnderscoreDash, "must consist of only alphanumeric, dash, or underscore characters")
	errs := validation.Errors{
		"name": validation.Validate(v.Name, validation.Required, rule),
	}
	return util.JoinErrs(tovalidate.ToErrors(errs)), nil
}

func (v *TOServerCapability) Read(h http.Header, useIMS bool) ([]interface{}, error, error, int, *time.Time) {
	api.DefaultSort(v.APIInfo(), "name")
	return api.GenericRead(h, v, useIMS)
}

func (v *TOServerCapability) Update(h http.Header) (error, error, int) {
	sc, userErr, sysErr, errCode, _ := v.Read(h, false)
	if userErr != nil || sysErr != nil {
		return userErr, sysErr, errCode
	}
	if len(sc) != 1 {
		return fmt.Errorf("cannot find exactly one server capability with the query string provided"), nil, http.StatusBadRequest
	}

	// check if the entity was already updated
	userErr, sysErr, errCode = api.CheckIfUnModifiedByName(h, v.ReqInfo.Tx, v.Name, "server_capability")
	if userErr != nil || sysErr != nil {
		return userErr, sysErr, errCode
	}

	// update server capability name
	rows, err := v.ReqInfo.Tx.Query(v.updateQuery(), v.RequestedName, v.Name)
	if err != nil {
		return nil, fmt.Errorf("server capability update: error setting the name for server capability %v: %v", v.Name, err.Error()), http.StatusInternalServerError
	}
	defer log.Close(rows, "unable to close DB connection")

	for rows.Next() {
		err = rows.Scan(&v.Name, &v.LastUpdated)
		if err != nil {
			return api.ParseDBError(err)
		}
	}
	return nil, nil, http.StatusOK
}

func GetServerCapability(w http.ResponseWriter, r *http.Request) {
	var sc tc.ServerCapabilityV4
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	// Query Parameters to Database Query column mappings
	queryParamsToQueryCols := map[string]dbhelpers.WhereColumnInfo{
		"name": {Column: "sc.name", Checker: nil},
	}
	if _, ok := inf.Params["orderby"]; !ok {
		inf.Params["orderby"] = "name"
	}
	where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(inf.Params, queryParamsToQueryCols)
	if len(errs) > 0 {
		api.HandleErr(w, r, tx.Tx, http.StatusBadRequest, util.JoinErrs(errs), nil)
		return
	}

	selectQuery := "SELECT name, description, last_updated FROM server_capability sc"
	query := selectQuery + where + orderBy + pagination
	rows, err := tx.NamedQuery(query, queryValues)
	if err != nil {
		api.HandleErr(w, r, tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("server capability read: error getting server capability(ies): %w", err))
		return
	}
	defer log.Close(rows, "unable to close DB connection")

	scList := []tc.ServerCapabilityV4{}
	for rows.Next() {
		if err = rows.Scan(&sc.Name, &sc.Description, &sc.LastUpdated); err != nil {
			api.HandleErr(w, r, tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("error getting server capability(ies): %w", err))
			return
		}
		scList = append(scList, sc)
	}

	api.WriteResp(w, r, scList)
	return
}

func UpdateServerCapability(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	tx := inf.Tx.Tx
	sc, readValErr := readAndValidateJsonStruct(r)
	if readValErr != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, readValErr, nil)
		return
	}

	requestedName := inf.Params["name"]
	// check if the entity was already updated
	userErr, sysErr, errCode = api.CheckIfUnModifiedByName(r.Header, inf.Tx, requestedName, "server_capability")
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	//update name and description of a capability
	query := `UPDATE server_capability sc SET
		name = $1,
		description = $2
	WHERE sc.name = $3
	RETURNING sc.name, sc.description, sc.last_updated`

	err := tx.QueryRow(query, sc.Name, sc.Description, requestedName).Scan(&sc.Name, &sc.Description, &sc.LastUpdated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("server capability with name: %s not found", sc.Name), nil)
			return
		}
		usrErr, sysErr, code := api.ParseDBError(err)
		api.HandleErr(w, r, tx, code, usrErr, sysErr)
		return
	}
	alerts := tc.CreateAlerts(tc.SuccessLevel, "server capability was updated")
	api.WriteAlertsObj(w, r, http.StatusOK, alerts, sc)
	changeLogMsg := fmt.Sprintf("CAPABILITY NAME:%s, ACTION: Updated serverCapability", sc.Name)
	api.CreateChangeLogRawTx(api.Updated, changeLogMsg, inf.User, tx)
	return
}

func CreateServerCapability(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	tx := inf.Tx.Tx

	sc, readValErr := readAndValidateJsonStruct(r)
	if readValErr != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, readValErr, nil)
		return
	}

	// check if capability already exists
	var count int
	err := tx.QueryRow("SELECT count(*) from server_capability where name = $1", sc.Name).Scan(&count)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("error: %w, when checking if server capability with name %s exists", err, sc.Name))
		return
	}
	if count == 1 {
		api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("server_capability name '%s' already exists.", sc.Name), nil)
		return
	}

	// create server capability
	query := `INSERT INTO server_capability (name, description) VALUES ($1, $2) RETURNING last_updated`
	err = tx.QueryRow(query, sc.Name, sc.Description).Scan(&sc.LastUpdated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, fmt.Errorf("error: %w in creating server capability with name: %s", err, sc.Name), nil)
			return
		}
		usrErr, sysErr, code := api.ParseDBError(err)
		api.HandleErr(w, r, tx, code, usrErr, sysErr)
		return
	}
	alerts := tc.CreateAlerts(tc.SuccessLevel, "server capability was created.")
	w.Header().Set(rfc.Location, fmt.Sprintf("/api/%s/server_capabilities?name=%s", inf.Version, sc.Name))
	api.WriteAlertsObj(w, r, http.StatusCreated, alerts, sc)
	changeLogMsg := fmt.Sprintf("CAPABILITY NAME:%s, ACTION: Created serverCapability", sc.Name)
	api.CreateChangeLogRawTx(api.Created, changeLogMsg, inf.User, tx)
	return
}

func DeleteServerCapability(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	name := inf.Params["name"]
	exists, err := dbhelpers.GetSCInfo(tx, name)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
		return
	}
	if !exists {
		if name != "" {
			api.HandleErr(w, r, tx, http.StatusNotFound, fmt.Errorf("no server capability exists by name: %s", name), nil)
			return
		} else {
			api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("no server capability exists for empty name: %s", name), nil)
			return
		}
	}

	assignedServer := 0
	if err := inf.Tx.Get(&assignedServer, "SELECT count(server) FROM server_server_capability ssc WHERE ssc.server_capability=$1", name); err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("server capability delete, counting assigned servers: %w", err))
		return
	} else if assignedServer != 0 {
		api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("can not delete a server capability with %d assigned servers", assignedServer), nil)
		return
	}

	res, err := tx.Exec("DELETE FROM server_capability AS sc WHERE sc.name=$1", name)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("determining rows affected for delete server capability: %w", err))
		return
	}
	if rowsAffected == 0 {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, fmt.Errorf("no rows deleted for server capability"), nil)
		return
	}
	alerts := tc.CreateAlerts(tc.SuccessLevel, "server capability was deleted.")
	api.WriteAlertsObj(w, r, http.StatusOK, alerts, name)
	changeLogMsg := fmt.Sprintf("CAPABILITY NAME:%s, ACTION: Deleted serverCapability", name)
	api.CreateChangeLogRawTx(api.Deleted, changeLogMsg, inf.User, tx)
	return
}

func readAndValidateJsonStruct(r *http.Request) (tc.ServerCapabilityV41, error) {
	var sc tc.ServerCapabilityV41
	if err := json.NewDecoder(r.Body).Decode(&sc); err != nil {
		userErr := fmt.Errorf("error decoding POST request body into ServerCapabilityV41 struct %w", err)
		return sc, userErr
	}

	// validate JSON body
	rule := validation.NewStringRule(tovalidate.IsAlphanumericUnderscoreDash, "must consist of only alphanumeric, dash, or underscore characters")
	errs := tovalidate.ToErrors(validation.Errors{
		"name": validation.Validate(sc.Name, validation.Required, rule),
	})
	if len(errs) > 0 {
		userErr := util.JoinErrs(errs)
		return sc, userErr
	}
	return sc, nil
}

func (v *TOServerCapability) SelectMaxLastUpdatedQuery(where, orderBy, pagination, tableName string) string {
	return `SELECT max(t) from (
		SELECT max(sc.last_updated) as t from server_capability sc ` + where + orderBy + pagination +
		` UNION ALL
	select max(l.last_updated) as t from last_deleted l where l.table_name='server_capability') as res`
}

func (v *TOServerCapability) Create() (error, error, int) { return api.GenericCreateNameBasedID(v) }
func (v *TOServerCapability) Delete() (error, error, int) { return api.GenericDelete(v) }

func GetServerCapabilityV5(w http.ResponseWriter, r *http.Request) {
	var sc tc.ServerCapabilityV5
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	// Query Parameters to Database Query column mappings
	queryParamsToQueryCols := map[string]dbhelpers.WhereColumnInfo{
		"name": {Column: "sc.name", Checker: nil},
	}
	if _, ok := inf.Params["orderby"]; !ok {
		inf.Params["orderby"] = "name"
	}
	where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(inf.Params, queryParamsToQueryCols)
	if len(errs) > 0 {
		api.HandleErr(w, r, tx.Tx, http.StatusBadRequest, util.JoinErrs(errs), nil)
		return
	}

	selectQuery := "SELECT name, description, last_updated FROM server_capability sc"
	query := selectQuery + where + orderBy + pagination
	rows, err := tx.NamedQuery(query, queryValues)
	if err != nil {
		api.HandleErr(w, r, tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("server capability read: error getting server capability(ies): %w", err))
		return
	}
	defer log.Close(rows, "unable to close DB connection")

	scList := []tc.ServerCapabilityV5{}
	for rows.Next() {
		if err = rows.Scan(&sc.Name, &sc.Description, &sc.LastUpdated); err != nil {
			api.HandleErr(w, r, tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("error getting server capability(ies): %w", err))
			return
		}
		scList = append(scList, sc)
	}

	api.WriteResp(w, r, scList)
	return
}

func CreateServerCapabilityV5(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	tx := inf.Tx.Tx

	sc, readValErr := readAndValidateJsonStructV5(r)
	if readValErr != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, readValErr, nil)
		return
	}

	// check if capability already exists
	var count int
	err := tx.QueryRow("SELECT count(*) from server_capability where name = $1", sc.Name).Scan(&count)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("error: %w, when checking if server capability with name %s exists", err, sc.Name))
		return
	}
	if count == 1 {
		api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("server_capability name '%s' already exists.", sc.Name), nil)
		return
	}

	// create server capability
	query := `INSERT INTO server_capability (name, description) VALUES ($1, $2) RETURNING last_updated`
	err = tx.QueryRow(query, sc.Name, sc.Description).Scan(&sc.LastUpdated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, fmt.Errorf("error: %w in creating server capability with name: %s", err, sc.Name), nil)
			return
		}
		usrErr, sysErr, code := api.ParseDBError(err)
		api.HandleErr(w, r, tx, code, usrErr, sysErr)
		return
	}
	alerts := tc.CreateAlerts(tc.SuccessLevel, "server capability was created.")
	w.Header().Set(rfc.Location, fmt.Sprintf("/api/%s/server_capabilities?name=%s", inf.Version, sc.Name))
	api.WriteAlertsObj(w, r, http.StatusCreated, alerts, sc)
	return
}

func UpdateServerCapabilityV5(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	tx := inf.Tx.Tx
	sc, readValErr := readAndValidateJsonStructV5(r)
	if readValErr != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, readValErr, nil)
		return
	}

	requestedName := inf.Params["name"]
	// check if the entity was already updated
	userErr, sysErr, errCode = api.CheckIfUnModifiedByName(r.Header, inf.Tx, requestedName, "server_capability")
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	//update name and description of a capability
	query := `UPDATE server_capability sc SET
		name = $1,
		description = $2
	WHERE sc.name = $3
	RETURNING sc.name, sc.description, sc.last_updated`

	err := tx.QueryRow(query, sc.Name, sc.Description, requestedName).Scan(&sc.Name, &sc.Description, &sc.LastUpdated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("server capability with name: %s not found", sc.Name), nil)
			return
		}
		usrErr, sysErr, code := api.ParseDBError(err)
		api.HandleErr(w, r, tx, code, usrErr, sysErr)
		return
	}
	alerts := tc.CreateAlerts(tc.SuccessLevel, "server capability was updated")
	api.WriteAlertsObj(w, r, http.StatusOK, alerts, sc)
	return
}

func readAndValidateJsonStructV5(r *http.Request) (tc.ServerCapabilityV5, error) {
	var sc tc.ServerCapabilityV5
	if err := json.NewDecoder(r.Body).Decode(&sc); err != nil {
		userErr := fmt.Errorf("error decoding POST request body into ServerCapabilityV5 struct %w", err)
		return sc, userErr
	}

	// validate JSON body
	rule := validation.NewStringRule(tovalidate.IsAlphanumericUnderscoreDash, "must consist of only alphanumeric, dash, or underscore characters")
	errs := tovalidate.ToErrors(validation.Errors{
		"name": validation.Validate(sc.Name, validation.Required, rule),
	})
	if len(errs) > 0 {
		userErr := util.JoinErrs(errs)
		return sc, userErr
	}
	return sc, nil
}
