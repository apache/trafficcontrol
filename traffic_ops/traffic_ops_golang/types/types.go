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
	"errors"
	"fmt"
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"

	validation "github.com/go-ozzo/ozzo-validation"
)

// TOTypeV5 is a needed type alias to define functions on.
type TOTypeV5 struct {
	api.APIInfoImpl `json:"-"`
	tc.TypeNullableV5
}

func (v *TOTypeV5) GetLastUpdated() (*time.Time, bool, error) {
	return api.GetLastUpdated(v.APIInfo().Tx, *v.ID, "type")
}

func (v *TOTypeV5) SetLastUpdated(t tc.TimeNoMod) {
	newTime, err := util.ConvertTimeFormat(t.Time, time.RFC3339)
	if err != nil {
		log.Errorf("Unable to convert Type last update time: %s\n", t.Time)
		v.LastUpdated = &t.Time
	}
	v.LastUpdated = newTime
}
func (v *TOTypeV5) InsertQuery() string     { return insertQuery() }
func (v *TOTypeV5) NewReadObj() interface{} { return &tc.TypeNullableV5{} }
func (v *TOTypeV5) SelectQuery() string     { return selectQuery() }
func (v *TOTypeV5) SelectMaxLastUpdatedQuery(where string, orderBy string, pagination string, tableName string) string {
	return selectMaxLastUpdatedQuery(where, orderBy, pagination, tableName)
}
func (v *TOTypeV5) ParamColumns() map[string]dbhelpers.WhereColumnInfo {
	return map[string]dbhelpers.WhereColumnInfo{
		"name":       dbhelpers.WhereColumnInfo{Column: "typ.name"},
		"id":         dbhelpers.WhereColumnInfo{Column: "typ.id", Checker: api.IsInt},
		"useInTable": dbhelpers.WhereColumnInfo{Column: "typ.use_in_table"},
	}
}
func (v *TOTypeV5) UpdateQuery() string { return updateQuery() }
func (v *TOTypeV5) DeleteQuery() string { return deleteQuery() }

func (typ TOTypeV5) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{Field: "id", Func: api.GetIntKey}}
}

// Implementation of the Identifier, Validator interface functions
func (typ TOTypeV5) GetKeys() (map[string]interface{}, bool) {
	if typ.ID == nil {
		return map[string]interface{}{"id": 0}, false
	}
	return map[string]interface{}{"id": *typ.ID}, true
}

func (typ *TOTypeV5) SetKeys(keys map[string]interface{}) {
	i, _ := keys["id"].(int) //this utilizes the non panicking type assertion, if the thrown away ok variable is false i will be the zero of the type, 0 here.
	typ.ID = &i
}

func (typ *TOTypeV5) GetAuditName() string {
	if typ.Name != nil {
		return *typ.Name
	}
	if typ.ID != nil {
		return strconv.Itoa(*typ.ID)
	}
	return "unknown"
}

func (typ *TOTypeV5) GetType() string {
	return "type"
}

func (typ *TOTypeV5) Validate() (error, error) {
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

func (tp *TOTypeV5) Read(h http.Header, useIMS bool) ([]interface{}, error, error, int, *time.Time) {
	api.DefaultSort(tp.APIInfo(), "name")
	return api.GenericRead(h, tp, useIMS)
}

func (tp *TOTypeV5) Update(h http.Header) (error, error, int) {
	if !tp.AllowMutation(false) {
		return errors.New("can not update type"), nil, http.StatusBadRequest
	}
	return api.GenericUpdate(h, tp)
}

func (tp *TOTypeV5) Delete() (error, error, int) {
	if !tp.AllowMutation(false) {
		return errors.New(fmt.Sprintf("can not delete type")), nil, http.StatusBadRequest
	}
	return api.GenericDelete(tp)
}

func (tp *TOTypeV5) Create() (error, error, int) {
	if !tp.AllowMutation(true) {
		return errors.New("can not create type"), nil, http.StatusBadRequest
	}
	return api.GenericCreate(tp)
}

func (tp *TOTypeV5) AllowMutation(forCreation bool) bool {
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

func (tp *TOTypeV5) loadUseInTable() (error, error, string) {
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

// V5 version ends here ---xxx---

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

// [V5]

func Get(w http.ResponseWriter, r *http.Request) {
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

	typeList, maxTime, code, usrErr, syErr = GetTypes(tx, inf.Params, useIMS, r.Header)
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

// GetTypes [Version : V5] returns types.
func GetTypes(tx *sqlx.Tx, params map[string]string, useIMS bool, header http.Header) ([]tc.TypeV5, time.Time, int, error, error) {
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

	query := selectQuery + where + orderBy + pagination
	rows, err := tx.NamedQuery(query, queryValues)
	if err != nil {
		return nil, time.Time{}, http.StatusInternalServerError, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		typ := tc.TypeV5{}

		if err = rows.Scan(&typ.ID, &typ.LastUpdated, &typ.Name, &typ.Description, &typ.UseInTable); err != nil {
			return nil, time.Time{}, http.StatusInternalServerError, nil, err
		}
		typeList = append(typeList, typ)
	}

	return typeList, maxTime, http.StatusOK, nil, nil
}
