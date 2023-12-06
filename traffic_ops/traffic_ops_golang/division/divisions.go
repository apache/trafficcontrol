package division

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
	"strings"
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

// we need a type alias to define functions on
type TODivision struct {
	api.APIInfoImpl `json:"-"`
	tc.DivisionNullable
}

func (v *TODivision) GetLastUpdated() (*time.Time, bool, error) {
	return api.GetLastUpdated(v.APIInfo().Tx, *v.ID, "division")
}

func (v *TODivision) SetLastUpdated(t tc.TimeNoMod) { v.LastUpdated = &t }
func (v *TODivision) InsertQuery() string           { return insertQuery() }
func (v *TODivision) SelectMaxLastUpdatedQuery(where, orderBy, pagination, tableName string) string {
	return `SELECT max(t) from (
		SELECT max(last_updated) as t from ` + tableName + ` d ` + where + orderBy + pagination +
		` UNION ALL
	select max(last_updated) as t from last_deleted l where l.table_name='` + tableName + `') as res`
}

func (v *TODivision) NewReadObj() interface{} { return &tc.Division{} }
func (v *TODivision) SelectQuery() string     { return selectQuery() }
func (v *TODivision) ParamColumns() map[string]dbhelpers.WhereColumnInfo {
	return map[string]dbhelpers.WhereColumnInfo{
		"id":   dbhelpers.WhereColumnInfo{Column: "id", Checker: api.IsInt},
		"name": dbhelpers.WhereColumnInfo{Column: "name"},
	}
}
func (v *TODivision) UpdateQuery() string { return updateQuery() }
func (v *TODivision) DeleteQuery() string { return deleteQuery() }

func (division TODivision) GetAuditName() string {
	if division.Name != nil {
		return *division.Name
	}
	if division.ID != nil {
		return strconv.Itoa(*division.ID)
	}
	return "unknown"
}

func (division TODivision) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{Field: "id", Func: api.GetIntKey}}
}

// Implementation of the Identifier, Validator interface functions
func (division TODivision) GetKeys() (map[string]interface{}, bool) {
	if division.ID == nil {
		return map[string]interface{}{"id": 0}, false
	}
	return map[string]interface{}{"id": *division.ID}, true
}

func (division *TODivision) SetKeys(keys map[string]interface{}) {
	i, _ := keys["id"].(int) //this utilizes the non panicking type assertion, if the thrown away ok variable is false i will be the zero of the type, 0 here.
	division.ID = &i
}

func (division TODivision) GetType() string {
	return "division"
}

func (division TODivision) Validate() (error, error) {
	errs := validation.Errors{
		"name": validation.Validate(division.Name, validation.NotNil, validation.Required),
	}
	return util.JoinErrs(tovalidate.ToErrors(errs)), nil
}

func (dv *TODivision) Create() (error, error, int) { return api.GenericCreate(dv) }
func (dv *TODivision) Read(h http.Header, useIMS bool) ([]interface{}, error, error, int, *time.Time) {
	params := dv.APIInfo().Params
	// TODO move to router, and do for all endpoints
	if strings.HasSuffix(params["name"], ".json") {
		params["name"] = params["name"][:len(params["name"])-len(".json")]
	}
	api.DefaultSort(dv.APIInfo(), "name")
	return api.GenericRead(h, dv, useIMS)
}
func (dv *TODivision) Update(h http.Header) (error, error, int) { return api.GenericUpdate(h, dv) }
func (dv *TODivision) Delete() (error, error, int)              { return api.GenericDelete(dv) }

func insertQuery() string {
	return `INSERT INTO division (name) VALUES (:name) RETURNING id,last_updated`
}

func selectQuery() string {
	return `SELECT
id,
last_updated,
name
FROM division d`
}

func updateQuery() string {
	return `UPDATE
division SET
name=:name
WHERE id=:id RETURNING last_updated`
}

func deleteQuery() string {
	return `DELETE FROM division WHERE id=:id`
}

func GetDivisions(w http.ResponseWriter, r *http.Request) {
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
		"id":   {Column: "d.id", Checker: nil},
		"name": {Column: "d.name", Checker: nil},
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
		runSecond, maxTime = ims.TryIfModifiedSinceQuery(tx, r.Header, queryValues, selectMaxLastUpdatedQuery(where))
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

	selectQuery := "SELECT id, name, last_updated FROM division d"
	query := selectQuery + where + orderBy + pagination
	rows, err := tx.NamedQuery(query, queryValues)
	if err != nil {
		api.HandleErr(w, r, tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("Divisions read: error getting divison(s): %w", err))
		return
	}
	defer log.Close(rows, "unable to close DB connection")

	div := tc.DivisionV5{}
	divList := []tc.DivisionV5{}
	for rows.Next() {
		if err = rows.Scan(&div.ID, &div.Name, &div.LastUpdated); err != nil {
			api.HandleErr(w, r, tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("error getting division(s): %w", err))
			return
		}
		divList = append(divList, div)
	}

	api.WriteResp(w, r, divList)
	return
}

func CreateDivision(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	tx := inf.Tx.Tx

	div, readValErr := readAndValidateJsonStruct(r)
	if readValErr != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, readValErr, nil)
		return
	}

	// check if division already exists
	var count int
	err := tx.QueryRow("SELECT count(*) from division where name = $1", div.Name).Scan(&count)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("error: %w, when checking if division with name %s exists", err, div.Name))
		return
	}
	if count == 1 {
		api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("division name '%s' already exists", div.Name), nil)
		return
	}

	// create division
	query := `INSERT INTO division (name) VALUES ($1) RETURNING id, last_updated`
	err = tx.QueryRow(query, div.Name).Scan(&div.ID, &div.LastUpdated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, fmt.Errorf("error: %w in creating division with name: %s", err, div.Name), nil)
			return
		}
		usrErr, sysErr, code := api.ParseDBError(err)
		api.HandleErr(w, r, tx, code, usrErr, sysErr)
		return
	}
	alerts := tc.CreateAlerts(tc.SuccessLevel, "division was created.")
	w.Header().Set(rfc.Location, fmt.Sprintf("/api/%s/divisons?name=%s", inf.Version, div.Name))
	api.WriteAlertsObj(w, r, http.StatusCreated, alerts, div)
	changeLogMsg := fmt.Sprintf("DIVISION: %s, ID:%d, ACTION: Created division", div.Name, div.ID)
	api.CreateChangeLogRawTx(api.Created, changeLogMsg, inf.User, tx)
	return
}

func UpdateDivision(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	tx := inf.Tx.Tx
	div, readValErr := readAndValidateJsonStruct(r)
	if readValErr != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, readValErr, nil)
		return
	}

	requestedID := inf.Params["id"]

	intRequestId, convErr := strconv.Atoi(requestedID)
	if convErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, fmt.Errorf("division update error: %w, while converting from string to int", convErr), nil)
	}
	// check if the entity was already updated
	userErr, sysErr, errCode = api.CheckIfUnModified(r.Header, inf.Tx, intRequestId, "division")
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	//update name and description of a division
	query := `UPDATE division div SET
		name = $2
	WHERE div.id = $1
	RETURNING div.id, div.last_updated`

	err := tx.QueryRow(query, requestedID, div.Name).Scan(&div.ID, &div.LastUpdated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("division with ID: %v not found", div.ID), nil)
			return
		}
		usrErr, sysErr, code := api.ParseDBError(err)
		api.HandleErr(w, r, tx, code, usrErr, sysErr)
		return
	}
	alerts := tc.CreateAlerts(tc.SuccessLevel, "division was updated")
	api.WriteAlertsObj(w, r, http.StatusOK, alerts, div)
	changeLogMsg := fmt.Sprintf("DIVISION: %s, ID:%d, ACTION: Updated division", div.Name, div.ID)
	api.CreateChangeLogRawTx(api.Updated, changeLogMsg, inf.User, tx)
	return
}

func DeleteDivision(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	id := inf.Params["id"]
	exists, err := dbhelpers.DivisionExists(tx, id)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
		return
	}
	if !exists {
		if id != "" {
			api.HandleErr(w, r, tx, http.StatusNotFound, fmt.Errorf("no divisions exists by id: %s", id), nil)
			return
		} else {
			api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("no divisions exists for empty id: %s", id), nil)
			return
		}
	}

	assignedRegions := 0
	if err := inf.Tx.Get(&assignedRegions, "SELECT count(id) FROM region reg WHERE reg.division=$1", id); err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("Divisions delete, counting assigned Regions: %w", err))
		return
	} else if assignedRegions != 0 {
		api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("can not delete a division with %d assigned region", assignedRegions), nil)
		return
	}

	res, err := tx.Exec("DELETE FROM division AS div WHERE div.id=$1", id)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("determining rows affected for delete division: %w", err))
		return
	}
	if rowsAffected == 0 {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, fmt.Errorf("no rows deleted for division"), nil)
		return
	}
	alerts := tc.CreateAlerts(tc.SuccessLevel, "division was deleted.")
	api.WriteAlertsObj(w, r, http.StatusOK, alerts, inf.Params)
	changeLogMsg := fmt.Sprintf("ID:%s, ACTION: Deleted division", id)
	api.CreateChangeLogRawTx(api.Deleted, changeLogMsg, inf.User, tx)
	return
}

// selectMaxLastUpdatedQuery used for TryIfModifiedSinceQuery()
func selectMaxLastUpdatedQuery(where string) string {
	return `SELECT max(t) from (
		SELECT max(a.last_updated) as t from division a` + where +
		` UNION ALL
	select max(last_updated) as t from last_deleted l where l.table_name='division') as res`
}

func readAndValidateJsonStruct(r *http.Request) (tc.DivisionV5, error) {
	var div tc.DivisionV5
	if err := json.NewDecoder(r.Body).Decode(&div); err != nil {
		userErr := fmt.Errorf("error decoding POST request body into DivisionV5 struct %w", err)
		return div, userErr
	}

	// validate JSON body
	rule := validation.NewStringRule(tovalidate.IsAlphanumericUnderscoreDash, "must consist of only alphanumeric, dash, or underscore characters")
	errs := tovalidate.ToErrors(validation.Errors{
		"name": validation.Validate(div.Name, validation.Required, rule),
	})
	if len(errs) > 0 {
		userErr := util.JoinErrs(errs)
		return div, userErr
	}
	return div, nil
}
