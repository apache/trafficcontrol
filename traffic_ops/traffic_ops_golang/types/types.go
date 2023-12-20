package types

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
	"strconv"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/util/ims"

	validation "github.com/go-ozzo/ozzo-validation"
)

// TOType is a needed type alias to define functions on.
type TOType struct {
	api.APIInfoImpl `json:"-"`
	tc.TypeNullable
}

func (v *TOType) GetLastUpdated() (*time.Time, bool, error) {
	return api.GetLastUpdated(v.APIInfo().Tx, *v.ID, "type")
}

func (v *TOType) SetLastUpdated(t tc.TimeNoMod) { v.LastUpdated = &t }
func (v *TOType) InsertQuery() string           { return insertQuery() }
func (v *TOType) NewReadObj() interface{}       { return &tc.TypeNullable{} }
func (v *TOType) SelectQuery() string           { return selectQuery() }
func (v *TOType) SelectMaxLastUpdatedQuery(where string, orderBy string, pagination string, tableName string) string {
	return selectMaxLastUpdatedQuery(where, orderBy, pagination, tableName)
}
func (v *TOType) ParamColumns() map[string]dbhelpers.WhereColumnInfo {
	return map[string]dbhelpers.WhereColumnInfo{
		"name":       dbhelpers.WhereColumnInfo{Column: "typ.name"},
		"id":         dbhelpers.WhereColumnInfo{Column: "typ.id", Checker: api.IsInt},
		"useInTable": dbhelpers.WhereColumnInfo{Column: "typ.use_in_table"},
	}
}
func (v *TOType) UpdateQuery() string { return updateQuery() }
func (v *TOType) DeleteQuery() string { return deleteQuery() }

func (typ TOType) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{Field: "id", Func: api.GetIntKey}}
}

// Implementation of the Identifier, Validator interface functions
func (typ TOType) GetKeys() (map[string]interface{}, bool) {
	if typ.ID == nil {
		return map[string]interface{}{"id": 0}, false
	}
	return map[string]interface{}{"id": *typ.ID}, true
}

func (typ *TOType) SetKeys(keys map[string]interface{}) {
	i, _ := keys["id"].(int) //this utilizes the non panicking type assertion, if the thrown away ok variable is false i will be the zero of the type, 0 here.
	typ.ID = &i
}

func (typ *TOType) GetAuditName() string {
	if typ.Name != nil {
		return *typ.Name
	}
	if typ.ID != nil {
		return strconv.Itoa(*typ.ID)
	}
	return "unknown"
}

func (typ *TOType) GetType() string {
	return "type"
}

func (typ *TOType) Validate() (error, error) {
	errs := validation.Errors{
		"name":         validation.Validate(typ.Name, validation.Required),
		"description":  validation.Validate(typ.Description, validation.Required),
		"use_in_table": validation.Validate(typ.UseInTable, validation.Required),
	}
	if errs != nil {
		return util.JoinErrs(tovalidate.ToErrors(errs)), nil
	}
	return nil, nil
}

func (tp *TOType) Read(h http.Header, useIMS bool) ([]interface{}, error, error, int, *time.Time) {
	api.DefaultSort(tp.APIInfo(), "name")
	return api.GenericRead(h, tp, useIMS)
}

func (tp *TOType) Update(h http.Header) (error, error, int) {
	if !tp.AllowMutation(false) {
		return errors.New("can not update type"), nil, http.StatusBadRequest
	}
	return api.GenericUpdate(h, tp)
}

func (tp *TOType) Delete() (error, error, int) {
	if !tp.AllowMutation(false) {
		return errors.New(fmt.Sprintf("can not delete type")), nil, http.StatusBadRequest
	}
	return api.GenericDelete(tp)
}

func (tp *TOType) Create() (error, error, int) {
	if !tp.AllowMutation(true) {
		return errors.New("can not create type"), nil, http.StatusBadRequest
	}
	return api.GenericCreate(tp)
}

func (tp *TOType) AllowMutation(forCreation bool) bool {
	if !forCreation {
		userErr, sysErr, actualUseInTable := tp.loadUseInTable()
		if userErr != nil || sysErr != nil {
			return false
		} else if actualUseInTable != "server" {
			return false
		}
	} else if *tp.UseInTable != "server" { // Only allow creating of types with UseInTable being "server"
		return false
	}
	return true
}

func (tp *TOType) loadUseInTable() (error, error, string) {
	var useInTable string
	// ID is only nil on creation, should not call this method in that case
	if tp.ID != nil {
		query := `SELECT use_in_table from type where id=$1`
		err := tp.ReqInfo.Tx.Tx.QueryRow(query, tp.ID).Scan(&useInTable)
		if err == sql.ErrNoRows {
			if tp.UseInTable == nil {
				return nil, nil, ""
			}
			return nil, nil, *tp.UseInTable
		}
		if err != nil {
			return nil, err, ""
		}
	} else {
		return errors.New("no type with that key found"), nil, ""
	}

	return nil, nil, useInTable
}

func selectMaxLastUpdatedQuery(where, orderBy, pagination, tableName string) string {
	return `SELECT max(t) from (
		SELECT max(last_updated) as t from ` + tableName + ` typ ` + where + orderBy + pagination +
		` UNION ALL
	select max(last_updated) as t from last_deleted l where l.table_name='type') as res`
}

func selectQuery() string {
	return `SELECT
id,
name,
description,
use_in_table,
last_updated
FROM type typ`
}

func updateQuery() string {
	query := `UPDATE
type SET
name=:name,
description=:description,
use_in_table=:use_in_table
WHERE id=:id RETURNING last_updated`
	return query
}

func insertQuery() string {
	query := `INSERT INTO type (
name,
description,
use_in_table) VALUES (
:name,
:description,
:use_in_table) RETURNING id,last_updated`
	return query
}

func deleteQuery() string {
	query := `DELETE FROM type
WHERE id=:id`
	return query
}

// Read [V5] - gets a list of types for APIv5
func Read(w http.ResponseWriter, r *http.Request) {
	var runSecond bool
	var maxTime time.Time
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	// Query Parameters to Database Query column mappings
	queryParamsToQueryCols := map[string]dbhelpers.WhereColumnInfo{
		"name":       {Column: "typ.name"},
		"id":         {Column: "typ.id", Checker: api.IsInt},
		"useInTable": {Column: "typ.use_in_table"},
	}
	if _, ok := inf.Params["orderby"]; !ok {
		inf.Params["orderby"] = "name"
	}
	where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(inf.Params, queryParamsToQueryCols)
	if len(errs) > 0 {
		api.HandleErr(w, r, tx.Tx, http.StatusBadRequest, util.JoinErrs(errs), nil)
	}

	if inf.Config.UseIMS {
		runSecond, maxTime = ims.TryIfModifiedSinceQuery(tx, r.Header, queryValues, SelectMaxLastUpdatedQuery(where))
		if !runSecond {
			log.Debugln("IMS HIT")
			api.AddLastModifiedHdr(w, maxTime)
			w.WriteHeader(http.StatusNotModified)
			return
		}
		log.Debugln("IMS MISS")
	} else {
		log.Debugln("Non IMS request")
	}

	query := selectQuery() + where + orderBy + pagination
	rows, err := tx.NamedQuery(query, queryValues)
	if err != nil {
		api.HandleErr(w, r, tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("type get: error getting type(s): %w", err))
	}
	defer log.Close(rows, "unable to close DB connection")

	typ := tc.TypeV5{}
	typeList := []tc.TypeV5{}
	for rows.Next() {
		if err = rows.Scan(&typ.ID, &typ.Name, &typ.Description, &typ.UseInTable, &typ.LastUpdated); err != nil {
			api.HandleErr(w, r, tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("error getting type(s): %w", err))
		}
		typeList = append(typeList, typ)
	}

	api.WriteResp(w, r, typeList)
	return
}

// Create [V5] -  creates the type with the passed data fpr APIv5.
func Create(w http.ResponseWriter, r *http.Request) {
	typ := tc.TypeV5{}

	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	tx := inf.Tx.Tx

	typ, readValErr := readAndValidateJsonStructV5(r)
	if readValErr != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, readValErr, nil)
		return
	}

	if typ.UseInTable != "server" {
		api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("can not create type."), nil)
		return
	}

	// check if type already exists
	var exists bool
	err := tx.QueryRow(`SELECT EXISTS(SELECT * from type where name = $1)`, typ.Name).Scan(&exists)

	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("error: %w, when checking if type with name %s exists", err, typ.Name))
		return
	}
	if exists {
		api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("type name '%s' already exists.", typ.Name), nil)
		return
	}

	// create type
	query := `INSERT INTO type (name, description, use_in_table) VALUES ($1, $2, $3) RETURNING id,last_updated`
	err = tx.QueryRow(query, typ.Name, typ.Description, typ.UseInTable).Scan(&typ.ID, &typ.LastUpdated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, fmt.Errorf("error: %w in creating type with name: %s", err, typ.Name), nil)
			return
		}
		usrErr, sysErr, code := api.ParseDBError(err)
		api.HandleErr(w, r, tx, code, usrErr, sysErr)
		return
	}
	alerts := tc.CreateAlerts(tc.SuccessLevel, "type was created.")
	w.Header().Set(rfc.Location, fmt.Sprintf("/api/%s/type?name=%s", inf.Version, typ.Name))
	api.WriteAlertsObj(w, r, http.StatusCreated, alerts, typ)
	changeLogMsg := fmt.Sprintf("TYPE: %s, ID:%d, ACTION: Created type", typ.Name, typ.ID)
	api.CreateChangeLogRawTx(api.Created, changeLogMsg, inf.User, tx)
	return
}

// Update [V5] - updates name & description of the type passed for APIv5.
func Update(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id"}, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	tx := inf.Tx.Tx

	typ, readValErr := readAndValidateJsonStructV5(r)
	if readValErr != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, readValErr, nil)
		return
	}

	requestedId := inf.IntParams["id"]
	// check if the entity was already updated
	userErr, sysErr, errCode = api.CheckIfUnModified(r.Header, inf.Tx, requestedId, "type")
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	if typ.UseInTable != "server" {
		api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("can not update type."), nil)
		return
	}

	//update type query
	query := `UPDATE type typ SET name= $1, description= $2, use_in_table= $3 WHERE typ.id=$4 RETURNING typ.id, typ.last_updated`

	err := tx.QueryRow(query, typ.Name, typ.Description, typ.UseInTable, requestedId).Scan(&typ.ID, &typ.LastUpdated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			api.HandleErr(w, r, tx, http.StatusNotFound, fmt.Errorf("type with id: %d not found", requestedId), nil)
			return
		}
		usrErr, sysErr, code := api.ParseDBError(err)
		api.HandleErr(w, r, tx, code, usrErr, sysErr)
		return
	}
	alerts := tc.CreateAlerts(tc.SuccessLevel, "type was updated")
	api.WriteAlertsObj(w, r, http.StatusOK, alerts, typ)
	changeLogMsg := fmt.Sprintf("TYPE: %s, ID:%d, ACTION: Updated type", typ.Name, typ.ID)
	api.CreateChangeLogRawTx(api.Updated, changeLogMsg, inf.User, tx)
	return
}

// Delete [V5] - deletes the type passed for APIv5.
func Delete(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	id := inf.Params["id"]
	// check if type already exists
	var exists bool
	err := tx.QueryRow(`SELECT EXISTS(SELECT * from type where id = $1)`, id).Scan(&exists)

	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
		return
	}
	if !exists {
		api.HandleErr(w, r, tx, http.StatusNotFound, fmt.Errorf("can not delete type"), nil)
		return
	}

	res, err := tx.Exec("DELETE FROM type AS typ WHERE typ.id=$1", id)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("determining rows affected for delete type: %w", err))
		return
	}
	if rowsAffected == 0 {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("no rows deleted for type"))
		return
	}

	alertMessage := fmt.Sprintf("type was deleted.")
	alerts := tc.CreateAlerts(tc.SuccessLevel, alertMessage)
	api.WriteAlerts(w, r, http.StatusOK, alerts)
	changeLogMsg := fmt.Sprintf("ID:%s, ACTION: Deleted type", id)
	api.CreateChangeLogRawTx(api.Deleted, changeLogMsg, inf.User, tx)
	return
}

// readAndValidateJsonStructV5 [V5] - validates the JSON object passed.
func readAndValidateJsonStructV5(r *http.Request) (tc.TypeV5, error) {
	var typ tc.TypeV5
	if err := json.NewDecoder(r.Body).Decode(&typ); err != nil {
		userErr := fmt.Errorf("error decoding POST request body into TypeV5 struct %w", err)
		return typ, userErr
	}

	// validate JSON body
	rule := validation.NewStringRule(tovalidate.IsAlphanumericUnderscoreDash, "must consist of only alphanumeric, dash, or underscore characters")
	errs := tovalidate.ToErrors(validation.Errors{
		"name":         validation.Validate(typ.Name, validation.Required, rule),
		"description":  validation.Validate(typ.Description, validation.Required),
		"use_in_table": validation.Validate(typ.UseInTable, validation.Required),
	})
	if len(errs) > 0 {
		userErr := util.JoinErrs(errs)
		return typ, userErr
	}
	return typ, nil
}

// SelectMaxLastUpdatedQuery used for TryIfModifiedSinceQuery()
func SelectMaxLastUpdatedQuery(where string) string {
	return `SELECT max(t) from (
		SELECT max(last_updated) as t from type typ ` + where +
		` UNION ALL
	select max(last_updated) as t from last_deleted l where l.table_name='type') as res`
}
