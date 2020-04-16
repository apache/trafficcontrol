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
	"net/http"
	"strconv"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"

	"github.com/go-ozzo/ozzo-validation"
)

//we need a type alias to define functions on
type TOType struct {
	api.APIInfoImpl `json:"-"`
	tc.TypeNullable
}

func (v *TOType) SetLastUpdated(t tc.TimeNoMod) { v.LastUpdated = &t }
func (v *TOType) InsertQuery() string           { return insertQuery() }
func (v *TOType) NewReadObj() interface{}       { return &tc.TypeNullable{} }
func (v *TOType) SelectQuery() string           { return selectQuery() }
func (v *TOType) ParamColumns() map[string]dbhelpers.WhereColumnInfo {
	return map[string]dbhelpers.WhereColumnInfo{
		"name":       dbhelpers.WhereColumnInfo{"typ.name", nil},
		"id":         dbhelpers.WhereColumnInfo{"typ.id", api.IsInt},
		"useInTable": dbhelpers.WhereColumnInfo{"typ.use_in_table", nil},
	}
}
func (v *TOType) UpdateQuery() string { return updateQuery() }
func (v *TOType) DeleteQuery() string { return deleteQuery() }

func (typ TOType) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{"id", api.GetIntKey}}
}

//Implementation of the Identifier, Validator interface functions
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

func (typ *TOType) Validate() error {
	errs := validation.Errors{
		"name":         validation.Validate(typ.Name, validation.Required),
		"description":  validation.Validate(typ.Description, validation.Required),
		"use_in_table": validation.Validate(typ.UseInTable, validation.Required),
	}
	if errs != nil {
		return util.JoinErrs(tovalidate.ToErrors(errs))
	}
	return nil
}

func (tp *TOType) Read(h map[string][]string) ([]interface{}, error, error, int) {
	ims := h["If-Modified-Since"]
	if ims != nil && len(ims) != 0 {
		fmt.Println("SRIJEET!! INSIDE TOTYPE IMS is ", ims[0])
	}
	//modifiedSince, ok := web.GetHTTPDate(h, "If-Modified-Since")
	//if !ok {
	//	fmt.Println("SC!!!!!!! ERRRRRRRRRROOOOOOOOOORRRRRRRRRRR!")
	//}
	var modifiedSince time.Time
	if t, err := time.Parse(time.RFC1123, ims[0]); err == nil {
		modifiedSince = t
	} else {
		fmt.Println("ERROR!!!! ", err)
	}
	fmt.Println("REQUEST TIME ", modifiedSince)
	results, e1, e2, code := api.GenericRead(tp)
	if e1 != nil && e2 != nil {
		return results, e1, e2, code
	}
	var finalResults []interface{}
	for _,r := range results {
		sri := r.(*tc.TypeNullable)
		fmt.Println("Srijeet!! TYPE IS TO TYPE 222222")
		fmt.Println("LAST UPDATED!!! ", sri.LastUpdated)

		if sri.LastUpdated.After(modifiedSince) {
			fmt.Println("MODIFIEDDDDDDDDD!!! Send 200 with new value")
			finalResults = append(finalResults, r)
		} else {
			fmt.Println("NOT MODIFIEDDDDDDD!!! Send 304")
		}
	}
	if len(finalResults) == 0 {
		return finalResults, e1, e2, 304
	}
	return results, e1, e2, code
}

func (tp *TOType) Update() (error, error, int) {
	if !tp.AllowMutation(false) {
		return errors.New("can not update type"), nil, http.StatusBadRequest
	}
	return api.GenericUpdate(tp)
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
	apiInfo := tp.APIInfo()
	if apiInfo.Version.Major < 2 {
		return true
	}
	if !forCreation {
		userErr, sysErr, actualUseInTable := tp.loadUseInTable()
		if userErr != nil || sysErr != nil {
			fmt.Println("USER ERR ", userErr)
			fmt.Println("System ERR ", sysErr)
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
