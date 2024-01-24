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
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/util/ims"

	validation "github.com/go-ozzo/ozzo-validation"
)

// we need a type alias to define functions on
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

// Implementation of the Identifier, Validator interface functions
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

func (region *TORegion) Validate() (error, error) {
	if len(region.Name) < 1 {
		return errors.New(`region 'name' is required`), nil
	}
	if region.Division == 0 {
		return errors.New(`region 'division' is required`), nil
	}
	return nil, nil
}

func (rg *TORegion) Read(h http.Header, useIMS bool) ([]interface{}, error, error, int, *time.Time) {
	api.DefaultSort(rg.APIInfo(), "name")
	return api.GenericRead(h, rg, useIMS)
}
func (v *TORegion) SelectMaxLastUpdatedQuery(where, orderBy, pagination, tableName string) string {
	return `SELECT max(t) from (
		SELECT max(r.last_updated) as t FROM region r
JOIN division d ON r.division = d.id ` + where + orderBy + pagination +
		` UNION ALL
	select max(last_updated) as t from last_deleted l where l.table_name='region') as res`
}

func (rg *TORegion) Update(h http.Header) (error, error, int) { return api.GenericUpdate(h, rg) }

func (rg *TORegion) Create() (error, error, int) {
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
			return nil, fmt.Errorf("could not scan after insert: %w)", err), http.StatusInternalServerError
		}
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("no region was inserted, nothing was returned"), http.StatusInternalServerError
	} else if rowsAffected > 1 {
		return nil, fmt.Errorf("too many rows affected from region insert"), http.StatusInternalServerError
	}

	rg.DivisionName = divisionName
	rg.ID = id
	rg.LastUpdated = lastUpdated
	return nil, nil, http.StatusOK
}
func (rg *TORegion) Delete() (error, error, int) { return api.GenericDelete(rg) }

// OptionsDelete deletes a resource identified either as a route parameter or as a query string parameter.
func (rg *TORegion) OptionsDelete() (error, error, int) { return api.GenericOptionsDelete(rg) }

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

// Read is the handler for GET requests to Regions of APIv5
func Read(w http.ResponseWriter, r *http.Request) {
	var maxTime time.Time
	var runSecond bool

	inf, sysErr, userErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx
	if sysErr != nil || userErr != nil {
		api.HandleErr(w, r, tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	// Query Parameters to Database Query column mappings
	queryParamsToQueryCols := map[string]dbhelpers.WhereColumnInfo{
		"division": {Column: "r.division"},
		"id":       {Column: "r.id", Checker: api.IsInt},
		"name":     {Column: "r.name"},
	}
	if _, ok := inf.Params["orderby"]; !ok {
		inf.Params["orderby"] = "name"
	}
	where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(inf.Params, queryParamsToQueryCols)
	if len(errs) > 0 {
		api.HandleErr(w, r, tx.Tx, http.StatusBadRequest, util.JoinErrs(errs), nil)
		return
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

	query := `SELECT r.division, d.name as divisionname, r.id, r.last_updated, r.name FROM region r
				JOIN division d ON r.division = d.id` + where + orderBy + pagination
	rows, err := tx.NamedQuery(query, queryValues)
	if err != nil {
		api.HandleErr(w, r, tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("region get: error getting region(s): %w", err))
	}
	defer log.Close(rows, "unable to close DB connection")

	typ := tc.RegionV5{}
	regionList := []tc.RegionV5{}
	for rows.Next() {
		if err = rows.Scan(&typ.Division, &typ.DivisionName, &typ.ID, &typ.LastUpdated, &typ.Name); err != nil {
			api.HandleErr(w, r, tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("error getting region(s): %w", err))
		}
		regionList = append(regionList, typ)
	}

	api.WriteResp(w, r, regionList)
	return
}

// Create region with the passed data for APIv5.
func Create(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	tx := inf.Tx.Tx
	rg, readValErr := readAndValidateJsonStruct(r)
	if readValErr != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, readValErr, nil)
		return
	}

	var exists bool
	existErr := tx.QueryRow(`SELECT EXISTS(SELECT * from region where name = $1)`, rg.Name).Scan(&exists)

	if existErr != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("error: %w, when checking if region with name %s exists", existErr, rg.Name))
		return
	}
	if exists {
		api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("region name '%s' already exists", rg.Name), nil)
		return
	}

	query := `INSERT INTO region ( division, name) VALUES ( $1, $2 ) 
				RETURNING 
					id,last_updated,
					(SELECT d.name FROM division d WHERE d.id = region.division)`
	err := tx.QueryRow(query, rg.Division, rg.Name).Scan(&rg.ID, &rg.LastUpdated, &rg.DivisionName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, fmt.Errorf("error: %w in creating region with name: %s", err, rg.Name), nil)
			return
		}
		usrErr, sysErr, code := api.ParseDBError(err)
		api.HandleErr(w, r, tx, code, usrErr, sysErr)
		return
	}
	alerts := tc.CreateAlerts(tc.SuccessLevel, "region is created.")
	api.WriteAlertsObj(w, r, http.StatusCreated, alerts, rg)
	changeLogMsg := fmt.Sprintf("REGION: %s, ID:%d, ACTION: Created region", rg.Name, rg.ID)
	api.CreateChangeLogRawTx(api.Created, changeLogMsg, inf.User, tx)
	return
}

// Update a Region for APIv5
func Update(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id"}, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	tx := inf.Tx.Tx
	rg, readValErr := readAndValidateJsonStruct(r)
	if readValErr != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, readValErr, nil)
		return
	}

	requestedRegionId := inf.IntParams["id"]
	// check if the entity was already updated
	userErr, sysErr, errCode = api.CheckIfUnModified(r.Header, inf.Tx, requestedRegionId, "region")
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	query := `UPDATE region SET division= $1, name= $2 WHERE id= $3 
				RETURNING 
					id,last_updated,
					(SELECT d.name FROM division d WHERE d.id = region.division)`

	err := tx.QueryRow(query, rg.Division, rg.Name, requestedRegionId).Scan(&rg.ID, &rg.LastUpdated, &rg.DivisionName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("region: %d not found", requestedRegionId), nil)
			return
		}
		usrErr, sysErr, code := api.ParseDBError(err)
		api.HandleErr(w, r, tx, code, usrErr, sysErr)
		return
	}

	alerts := tc.CreateAlerts(tc.SuccessLevel, "region was updated")
	api.WriteAlertsObj(w, r, http.StatusOK, alerts, rg)
	changeLogMsg := fmt.Sprintf("REGION: %s, ID:%d, ACTION: Updated region", rg.Name, rg.ID)
	api.CreateChangeLogRawTx(api.Updated, changeLogMsg, inf.User, tx)
	return
}

// Delete an Region for APIv5
func Delete(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	requestedRegionId := inf.Params["id"]
	requestedRegionName := inf.Params["name"]

	var exists bool
	var noExistsErrMsg string

	if requestedRegionId == "" && requestedRegionName == "" {
		api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("refusing to delete all resources of type Region"), nil)
		return
	} else if requestedRegionId != "" && requestedRegionName != "" { // checking if both id and name are passed
		existErr := tx.QueryRow(`SELECT EXISTS(SELECT * from region where id = $1 AND name = $2)`, requestedRegionId, requestedRegionName).Scan(&exists)
		noExistsErrMsg = fmt.Sprintf("no region exists by id: %s and name: %s", requestedRegionId, requestedRegionName)
		if existErr != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("error: %w, when checking if region with id %s exists", existErr, requestedRegionId))
			return
		}
	} else if requestedRegionId != "" {
		existErr := tx.QueryRow(`SELECT EXISTS(SELECT * from region where id = $1)`, requestedRegionId).Scan(&exists)
		noExistsErrMsg = fmt.Sprintf("no region exists by id: %s", requestedRegionId)
		if existErr != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("error: %w, when checking if region with id %s exists", existErr, requestedRegionId))
			return
		}
	} else if requestedRegionName != "" {
		existErr := tx.QueryRow(`SELECT EXISTS(SELECT * from region where name = $1)`, requestedRegionName).Scan(&exists)
		noExistsErrMsg = fmt.Sprintf("no region exists by name: %s", requestedRegionName)
		if existErr != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("error: %w, when checking if region with name %s exists", existErr, requestedRegionName))
			return
		}
		if exists {
			existErr := tx.QueryRow(`SELECT id from region where name = $1`, requestedRegionName).Scan(&requestedRegionId)
			if existErr != nil {
				api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("error: %w, when checking if region with name %s exists", existErr, requestedRegionName))
				return
			}
		}
	}

	if !exists {
		api.HandleErr(w, r, tx, http.StatusNotFound, fmt.Errorf(noExistsErrMsg), nil)
		return
	}

	res, err := tx.Exec("DELETE FROM region WHERE id = $1", requestedRegionId)

	if err != nil {
		userErr, sysErr, errCode := api.ParseDBError(err)
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("determining rows affected for delete region: %w", err))
		return
	}
	if rowsAffected == 0 {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, fmt.Errorf("no rows deleted for region"), nil)
		return
	}
	alerts := tc.CreateAlerts(tc.SuccessLevel, "region was deleted.")
	api.WriteAlerts(w, r, http.StatusOK, alerts)
	changeLogMsg := fmt.Sprintf("ID:%s, ACTION: Deleted region", requestedRegionId)
	api.CreateChangeLogRawTx(api.Deleted, changeLogMsg, inf.User, tx)
	return
}

// readAndValidateJsonStruct reads json body and validates json fields
func readAndValidateJsonStruct(r *http.Request) (tc.RegionV5, error) {
	var region tc.RegionV5
	if err := json.NewDecoder(r.Body).Decode(&region); err != nil {
		userErr := fmt.Errorf("error decoding POST request body into RegionV5 struct %w", err)
		return region, userErr
	}

	// validate JSON body
	errs := tovalidate.ToErrors(validation.Errors{
		"division": validation.Validate(region.Division, validation.NotNil, validation.Min(0)),
		"name":     validation.Validate(region.Name, validation.Required),
	})
	if len(errs) > 0 {
		userErr := util.JoinErrs(errs)
		return region, userErr
	}
	return region, nil
}

// SelectMaxLastUpdatedQuery used for TryIfModifiedSinceQuery()
func SelectMaxLastUpdatedQuery(where string) string {
	return `SELECT max(t) from (
		SELECT max(r.last_updated) as t FROM region r
JOIN division d ON r.division = d.id ` + where +
		` UNION ALL
	select max(last_updated) as t from last_deleted l where l.table_name='region') as res`
}
