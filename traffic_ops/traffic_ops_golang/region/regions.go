package region

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
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/crudder"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
)

//we need a type alias to define functions on
type TORegion struct {
	api.APIInfoImpl `json:"-"`
	tc.Region
}

func (v *TORegion) GetLastUpdated() (*time.Time, bool, error) {
	return api.GetLastUpdated(v.APIInfo().Tx, v.ID, "region")
}

func (v *TORegion) SetLastUpdated(t tc.TimeNoMod) { v.LastUpdated = t }
func (v *TORegion) InsertQuery() string           { return insertQuery() }
func (v *TORegion) NewReadObj() interface{}       { return &tc.Region{} }
func (v *TORegion) SelectQuery() string           { return selectQuery() }
func (v *TORegion) ParamColumns() map[string]dbhelpers.WhereColumnInfo {
	return map[string]dbhelpers.WhereColumnInfo{
		"name":     dbhelpers.WhereColumnInfo{Column: "r.name"},
		"division": dbhelpers.WhereColumnInfo{Column: "r.division"},
		"id":       dbhelpers.WhereColumnInfo{Column: "r.id", Checker: api.IsInt},
	}
}
func (v *TORegion) UpdateQuery() string { return updateQuery() }

// DeleteQuery returns a query, including a WHERE clause.
func (v *TORegion) DeleteQuery() string { return deleteQuery() }

// DeleteQueryBase returns a query with no WHERE clause.
func (v *TORegion) DeleteQueryBase() string { return deleteQueryBase() }

func (region TORegion) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{Field: "id", Func: api.GetIntKey}}
}

//Implementation of the Identifier, Validator interface functions
func (region TORegion) GetKeys() (map[string]interface{}, bool) {
	return map[string]interface{}{"id": region.ID}, true
}

// DeleteKeyOptions returns a map containing the different fields a resource can be deleted by.
func (region TORegion) DeleteKeyOptions() map[string]dbhelpers.WhereColumnInfo {
	return map[string]dbhelpers.WhereColumnInfo{
		"id":   dbhelpers.WhereColumnInfo{Column: "r.id", Checker: api.IsInt},
		"name": dbhelpers.WhereColumnInfo{Column: "r.name"},
	}
}

func (region *TORegion) SetKeys(keys map[string]interface{}) {
	//this utilizes the non panicking type assertion, if the thrown away ok variable is false i will be the zero of the type, 0 here.
	if id, exists := keys["id"].(int); exists {
		region.ID = id
	}
	if name, exists := keys["name"].(string); exists {
		region.Name = name
	}
}

func (region *TORegion) GetAuditName() string {
	return region.Name
}

func (region *TORegion) GetType() string {
	return "region"
}

func (region *TORegion) Validate() error {
	if len(region.Name) < 1 {
		return errors.New(`region 'name' is required`)
	}
	if region.Division == 0 {
		return errors.New(`region 'division' is required`)
	}
	return nil
}

func (rg *TORegion) Read(h http.Header, useIMS bool) ([]interface{}, api.Errors, *time.Time) {
	api.DefaultSort(rg.APIInfo(), "name")
	return crudder.GenericRead(h, rg, useIMS)
}
func (v *TORegion) SelectMaxLastUpdatedQuery(where, orderBy, pagination, tableName string) string {
	return `SELECT max(t) from (
		SELECT max(r.last_updated) as t FROM region r
JOIN division d ON r.division = d.id ` + where + orderBy + pagination +
		` UNION ALL
	select max(last_updated) as t from last_deleted l where l.table_name='region') as res`
}

func (rg *TORegion) Update(h http.Header) api.Errors { return crudder.GenericUpdate(h, rg) }

func (rg *TORegion) Create() api.Errors {
	resultRows, err := rg.APIInfo().Tx.NamedQuery(rg.InsertQuery(), rg)
	if err != nil {
		return api.ParseDBError(err)
	}
	defer resultRows.Close()

	var divisionName string
	var id int
	var lastUpdated tc.TimeNoMod

	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err = resultRows.Scan(&id, &lastUpdated, &divisionName); err != nil {
			return api.NewSystemError(fmt.Errorf("could not scan after insert: %w)", err))
		}
	}

	if rowsAffected == 0 {
		return api.NewSystemError(errors.New("no region was inserted, nothing was returned"))
	} else if rowsAffected > 1 {
		return api.NewSystemError(errors.New("too many rows affected from region insert"))
	}

	rg.DivisionName = divisionName
	rg.ID = id
	rg.LastUpdated = lastUpdated
	return api.NewErrors()
}
func (rg *TORegion) Delete() api.Errors { return crudder.GenericDelete(rg) }

// OptionsDelete deletes a resource identified either as a route parameter or as a query string parameter.
func (rg *TORegion) OptionsDelete() api.Errors { return crudder.GenericOptionsDelete(rg) }

func selectQuery() string {
	return `SELECT
r.division,
d.name as divisionname,
r.id,
r.last_updated,
r.name
FROM region r
JOIN division d ON r.division = d.id`
}

func updateQuery() string {
	query := `UPDATE
region SET
division=:division,
name=:name
WHERE id=:id RETURNING last_updated`
	return query
}

func insertQuery() string {
	query := `INSERT INTO region (
division,
name) VALUES (
:division,
:name) RETURNING id,last_updated,
(SELECT d.name FROM division d WHERE id = region.division)`
	return query
}

func deleteQueryBase() string {
	query := `DELETE FROM region r`
	return query
}

func deleteQuery() string {
	query := deleteQueryBase() + ` WHERE id=:id`
	return query
}
