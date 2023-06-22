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

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/jmoiron/sqlx"

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

// GetV5 [Version :V5] - GetV5 will retrieve a list of types
func GetV5(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	code := http.StatusOK
	useIMS := false
	config, e := api.GetConfig(r.Context())
	if e == nil && config != nil {
		useIMS = config.UseIMS
	} else {
		log.Warnf("Couldn't get config %v", e)
	}

	var maxTime time.Time
	var usrErr error
	var syErr error

	var typeList []tc.TypeV5

	tx := inf.Tx

	typeList, maxTime, code, usrErr, syErr = func(tx *sqlx.Tx, params map[string]string, useIMS bool, header http.Header) ([]tc.TypeV5, time.Time, int, error, error) {
		var runSecond bool
		var maxTime time.Time
		typeList := []tc.TypeV5{}

		selectQuery := `SELECT id, name, description, use_in_table, last_updated FROM type as typ`

		// Query Parameters to Database Query column mappings
		queryParamsToQueryCols := map[string]dbhelpers.WhereColumnInfo{
			"name":       {Column: "typ.name"},
			"id":         {Column: "typ.id", Checker: api.IsInt},
			"useInTable": {Column: "typ.use_in_table"},
		}
		if _, ok := params["orderby"]; !ok {
			params["orderby"] = "name"
		}
		where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(params, queryParamsToQueryCols)
		if len(errs) > 0 {
			return nil, time.Time{}, http.StatusBadRequest, util.JoinErrs(errs), nil
		}

		if useIMS {
			runSecond, maxTime = TryIfModifiedSinceQuery(header, tx, where, queryValues)
			if !runSecond {
				log.Debugln("IMS HIT")
				return typeList, maxTime, http.StatusNotModified, nil, nil
			}
			log.Debugln("IMS MISS")
		} else {
			log.Debugln("Non IMS request")
		}

		query := selectQuery + where + orderBy + pagination
		rows, err := tx.NamedQuery(query, queryValues)
		if err != nil {
			return nil, time.Time{}, http.StatusInternalServerError, nil, err
		}
		defer rows.Close()

		for rows.Next() {
			typ := tc.TypeV5{}

			if err = rows.Scan(&typ.ID, &typ.Name, &typ.Description, &typ.UseInTable, &typ.LastUpdated); err != nil {
				return nil, time.Time{}, http.StatusInternalServerError, nil, err
			}
			typeList = append(typeList, typ)
		}

		return typeList, maxTime, http.StatusOK, nil, nil
	}(tx, inf.Params, useIMS, r.Header)

	if code == http.StatusNotModified {
		w.WriteHeader(code)
		api.WriteResp(w, r, []tc.TypeV5{})
		return
	}

	if code == http.StatusBadRequest {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, usrErr, nil)
		return
	}

	if sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, syErr)
		return
	}

	if maxTime != (time.Time{}) && api.SetLastModifiedHeader(r, useIMS) {
		api.AddLastModifiedHdr(w, maxTime)
	}

	api.WriteResp(w, r, typeList)
}

// TryIfModifiedSinceQuery [Version : V5] function receives types and header from GetTypesV5 function and returns bool value if status is not modified.
func TryIfModifiedSinceQuery(header http.Header, tx *sqlx.Tx, where string, queryValues map[string]interface{}) (bool, time.Time) {
	var max time.Time
	var imsDate time.Time
	var ok bool
	imsDateHeader := []string{}
	runSecond := true
	dontRunSecond := false

	if header == nil {
		return runSecond, max
	}

	imsDateHeader = header[rfc.IfModifiedSince]
	if len(imsDateHeader) == 0 {
		return runSecond, max
	}

	if imsDate, ok = rfc.ParseHTTPDate(imsDateHeader[0]); !ok {
		log.Warnf("IMS request header date '%s' not parsable", imsDateHeader[0])
		return runSecond, max
	}

	imsQuery := `SELECT max(last_updated) as t from type typ`
	query := imsQuery + where
	rows, err := tx.NamedQuery(query, queryValues)

	if errors.Is(err, sql.ErrNoRows) {
		return dontRunSecond, max
	}

	if err != nil {
		log.Warnf("Couldn't get the max last updated time: %v", err)
		return runSecond, max
	}

	defer rows.Close()
	// This should only ever contain one row
	if rows.Next() {
		v := time.Time{}
		if err = rows.Scan(&v); err != nil {
			log.Warnf("Failed to parse the max time stamp into a struct %v", err)
			return runSecond, max
		}

		max = v
		// The request IMS time is later than the max of (lastUpdated, deleted_time)
		if imsDate.After(v) {
			return dontRunSecond, max
		}
	}
	return runSecond, max
}

// CreateType [Version : V5] - CreateTypeV5 function creates the type with the passed data.
func CreateTypeV5(w http.ResponseWriter, r *http.Request) {
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
	w.Header().Set("Location", fmt.Sprintf("/api/%d.%d/type?name=%s", inf.Version.Major, inf.Version.Minor, typ.Name))
	api.WriteAlertsObj(w, r, http.StatusCreated, alerts, typ)
	return
}

// UpdateType [Version : V5] - UpdateTypeV5 function updates name & description of the type passed.
func UpdateTypeV5(w http.ResponseWriter, r *http.Request) {
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

	requestedId := inf.Params["id"]
	// check if the entity was already updated
	userErr, sysErr, errCode = api.CheckIfUnModifiedByName(r.Header, inf.Tx, requestedId, "type")
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
			api.HandleErr(w, r, tx, http.StatusNotFound, fmt.Errorf("type with id: %s not found", requestedId), nil)
			return
		}
		usrErr, sysErr, code := api.ParseDBError(err)
		api.HandleErr(w, r, tx, code, usrErr, sysErr)
		return
	}
	alerts := tc.CreateAlerts(tc.SuccessLevel, "type was updated")
	api.WriteAlertsObj(w, r, http.StatusOK, alerts, typ)
	return
}

// DeleteType [Version : V5] - DeleteTypeV5 function deletes the type passed.
func DeleteTypeV5(w http.ResponseWriter, r *http.Request) {
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
	return
}

// readAndValidateJsonStructV5 [Version : V5] - readAndValidateJsonStructV5 function validates the JSON object passed.
func readAndValidateJsonStructV5(r *http.Request) (tc.TypeV5, error) {
	var typ tc.TypeV5
	if err := json.NewDecoder(r.Body).Decode(&typ); err != nil {
		userErr := fmt.Errorf("error decoding POST request body into TypeV5 struct %w", err)
		return typ, userErr
	}

	// validate JSON body
	rule := validation.NewStringRule(tovalidate.IsAlphanumericUnderscoreDash, "must consist of only alphanumeric, dash, or underscore characters")
	errs := tovalidate.ToErrors(validation.Errors{
		"name": validation.Validate(typ.Name, validation.Required, rule),
	})
	if len(errs) > 0 {
		userErr := util.JoinErrs(errs)
		return typ, userErr
	}
	return typ, nil
}
