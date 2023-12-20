package asn

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

// ASNsPrivLevel ...
const ASNsPrivLevel = 10

// we need a type alias to define functions on
type TOASNV11 struct {
	api.APIInfoImpl `json:"-"`
	tc.ASNNullable
}

func (v *TOASNV11) GetLastUpdated() (*time.Time, bool, error) {
	return api.GetLastUpdated(v.APIInfo().Tx, *v.ID, "asn")
}

func (v *TOASNV11) SetLastUpdated(t tc.TimeNoMod) { v.LastUpdated = &t }
func (v *TOASNV11) InsertQuery() string           { return insertQuery() }
func (v *TOASNV11) NewReadObj() interface{}       { return &tc.ASNNullable{} }
func (v *TOASNV11) SelectQuery() string           { return selectQuery() }
func (v *TOASNV11) ParamColumns() map[string]dbhelpers.WhereColumnInfo {
	return map[string]dbhelpers.WhereColumnInfo{
		"asn":            dbhelpers.WhereColumnInfo{Column: "a.asn", Checker: api.IsInt},
		"cachegroup":     dbhelpers.WhereColumnInfo{Column: "c.id", Checker: api.IsInt},
		"id":             dbhelpers.WhereColumnInfo{Column: "a.id", Checker: api.IsInt},
		"cachegroupName": dbhelpers.WhereColumnInfo{Column: "c.name"},
	}
}
func (v *TOASNV11) UpdateQuery() string { return updateQuery() }
func (v *TOASNV11) DeleteQuery() string { return deleteQuery() }
func (asn TOASNV11) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{Field: "id", Func: api.GetIntKey}}
}

// func (asn TOASNV12) GetKeyFieldsInfo() []api.KeyFieldInfo { return TOASNV11(asn).GetKeyFieldsInfo() }

// Implementation of the Identifier, Validator interface functions

func (asn TOASNV11) GetKeys() (map[string]interface{}, bool) {
	if asn.ID == nil {
		return map[string]interface{}{"id": 0}, false
	}
	return map[string]interface{}{"id": *asn.ID}, true
}

func (asn *TOASNV11) SetKeys(keys map[string]interface{}) {
	i, _ := keys["id"].(int) //this utilizes the non panicking type assertion, if the thrown away ok variable is false i will be the zero of the type, 0 here.
	asn.ID = &i
}

func (asn TOASNV11) GetAuditName() string {
	if asn.ASN != nil {
		return strconv.Itoa(*asn.ASN)
	}
	if asn.ID != nil {
		return strconv.Itoa(*asn.ID)
	}
	return "unknown"
}

func (asn TOASNV11) GetType() string {
	return "asn"
}

func (asn TOASNV11) Validate() (error, error) {
	errs := validation.Errors{
		"asn":          validation.Validate(asn.ASN, validation.NotNil, validation.Min(0)),
		"cachegroupId": validation.Validate(asn.CachegroupID, validation.NotNil, validation.Min(0)),
	}
	return util.JoinErrs(tovalidate.ToErrors(errs)), nil
}

func (as *TOASNV11) Create() (error, error, int) {
	err := as.ASNExists(true)
	if err != nil {
		return err, nil, http.StatusBadRequest
	}
	return api.GenericCreate(as)
}
func (as *TOASNV11) Read(h http.Header, useIMS bool) ([]interface{}, error, error, int, *time.Time) {
	api.DefaultSort(as.APIInfo(), "asn")
	return api.GenericRead(h, as, useIMS)
}
func (v *TOASNV11) SelectMaxLastUpdatedQuery(where, orderBy, pagination, tableName string) string {
	return `SELECT max(t) from (
		SELECT max(a.last_updated) as t from asn a
JOIN
  cachegroup c ON a.cachegroup = c.id ` + where + orderBy + pagination +
		` UNION ALL
	select max(last_updated) as t from last_deleted l where l.table_name='asn') as res`
}

func (as *TOASNV11) Update(h http.Header) (error, error, int) {
	err := as.ASNExists(false)
	if err != nil {
		return err, nil, http.StatusBadRequest
	}
	return api.GenericUpdate(h, as)
}
func (as *TOASNV11) Delete() (error, error, int) { return api.GenericDelete(as) }

func (asn TOASNV11) ASNExists(create bool) error {
	if asn.APIInfo() == nil || asn.APIInfo().Tx == nil {
		return errors.New("couldn't perform check to see if asn number exists already")
	}
	if asn.ASN == nil || asn.CachegroupID == nil {
		return errors.New("no asn or cachegroup ID specified")
	}
	query := `SELECT id from asn where asn=$1`
	rows, err := asn.APIInfo().Tx.Query(query, *asn.ASN)
	if err != nil {
		return errors.New("selecting asns: " + err.Error())
	}
	defer rows.Close()
	if rows.Next() {
		if create {
			return errors.New("an asn with the specified number already exists")
		}
		var v int
		id := *asn.ID
		err = rows.Scan(&v)
		if err != nil {
			return errors.New("couldn't check if this number exists")
		}
		if v != id {
			return errors.New("another asn exists for this number")
		}
	}
	return nil
}

func selectQuery() string {
	return `
SELECT
 a.id,
 a.asn,
 a.last_updated,
 a.cachegroup AS cachegroup_id,
 c.name AS cachegroup
FROM
  asn a
JOIN
  cachegroup c ON a.cachegroup = c.id
`
}

func insertQuery() string {
	return `
INSERT INTO
  asn (asn, cachegroup)
VALUES
  (:asn, :cachegroup_id)
RETURNING id, last_updated
`
}

func updateQuery() string {
	return `
UPDATE
  asn
SET
  asn        = :asn,
  cachegroup = :cachegroup_id
WHERE
  id = :id
RETURNING
  last_updated
`
}

func deleteQuery() string {
	return `DELETE FROM asn WHERE id=:id`
}

// Read gets list of ASNs for APIv5
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
		"asn": {Column: "a.asn", Checker: api.IsInt},
		"id":  {Column: "a.id", Checker: api.IsInt},
	}
	if _, ok := inf.Params["orderby"]; !ok {
		inf.Params["orderby"] = "asn"
	}
	where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(inf.Params, queryParamsToQueryCols)
	if len(errs) > 0 {
		api.HandleErr(w, r, tx.Tx, http.StatusBadRequest, util.JoinErrs(errs), nil)
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

	query := selectQuery() + where + orderBy + pagination
	rows, err := tx.NamedQuery(query, queryValues)
	if err != nil {
		api.HandleErr(w, r, tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("asn read: error getting asn(s): %w", err))
	}
	defer log.Close(rows, "unable to close DB connection")

	asn := tc.ASNV5{}
	asnList := []tc.ASNV5{}
	for rows.Next() {
		if err = rows.Scan(&asn.ID, &asn.ASN, &asn.LastUpdated, &asn.CachegroupID, &asn.Cachegroup); err != nil {
			api.HandleErr(w, r, tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("error getting asn(s): %w", err))
		}
		asnList = append(asnList, asn)
	}

	api.WriteResp(w, r, asnList)
	return
}

// Create an ASN for APIv5
func Create(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	tx := inf.Tx.Tx

	asn, readValErr := readAndValidateJsonStruct(r)
	if readValErr != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, readValErr, nil)
		return
	}

	// check if asn already exists
	var count int
	err := tx.QueryRow("SELECT count(*) from asn where asn=$1", asn.ASN).Scan(&count)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("error: %w, when checking if asn '%d' exists", err, asn.ASN))
		return
	}
	if count == 1 {
		api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("asn:'%d' already exists", asn.ASN), nil)
		return
	}

	// create asn
	query := `INSERT INTO asn (asn, cachegroup) VALUES ($1, $2) RETURNING id, last_updated, (select name FROM cachegroup where id = $2)`
	err = tx.QueryRow(query, asn.ASN, asn.CachegroupID).Scan(&asn.ID, &asn.LastUpdated, &asn.Cachegroup)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, fmt.Errorf("error: %w in creating asn:%d", err, asn.ASN), nil)
			return
		}
		usrErr, sysErr, code := api.ParseDBError(err)
		api.HandleErr(w, r, tx, code, usrErr, sysErr)
		return
	}

	alerts := tc.CreateAlerts(tc.SuccessLevel, "asn was created.")
	w.Header().Set(rfc.Location, fmt.Sprintf("/api/%s/asns?id=%d", inf.Version, asn.ID))
	api.WriteAlertsObj(w, r, http.StatusCreated, alerts, asn)
	changeLogMsg := fmt.Sprintf("ASN: %d, ID:%d, ACTION: Created asn", asn.ASN, asn.ID)
	api.CreateChangeLogRawTx(api.Created, changeLogMsg, inf.User, tx)
	return
}

// Update an ASN for APIv5
func Update(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id"}, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	tx := inf.Tx.Tx
	asn, readValErr := readAndValidateJsonStruct(r)
	if readValErr != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, readValErr, nil)
		return
	}

	requestedAsnId := inf.IntParams["id"]
	// check if the entity was already updated
	userErr, sysErr, errCode = api.CheckIfUnModified(r.Header, inf.Tx, requestedAsnId, "asn")
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	// check if asn already exists
	var id int
	err := tx.QueryRow("SELECT id from asn where asn=$1", asn.ASN).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("error: %w, when checking if asn '%d' exists", err, asn.ASN))
		return
	}
	if id != 0 && id != requestedAsnId {
		api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("asn:'%d' already exists", asn.ASN), nil)
		return
	}

	//update asn and cachegroup of an asn
	query := `UPDATE asn SET
		asn = $1,
		cachegroup = $2
	WHERE id = $3
	RETURNING id, last_updated, (select name FROM cachegroup where id = $2)`

	err = tx.QueryRow(query, asn.ASN, asn.CachegroupID, requestedAsnId).Scan(&asn.ID, &asn.LastUpdated, &asn.Cachegroup)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("asn: %d not found", asn.ASN), nil)
			return
		}
		usrErr, sysErr, code := api.ParseDBError(err)
		api.HandleErr(w, r, tx, code, usrErr, sysErr)
		return
	}

	alerts := tc.CreateAlerts(tc.SuccessLevel, "asn was updated")
	api.WriteAlertsObj(w, r, http.StatusOK, alerts, asn)
	changeLogMsg := fmt.Sprintf("ASN: %d, ID:%d, ACTION: Updated asn", asn.ASN, asn.ID)
	api.CreateChangeLogRawTx(api.Updated, changeLogMsg, inf.User, tx)
	return
}

// Delete an ASN for APIv5
func Delete(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	id := inf.Params["id"]
	exists, err := dbhelpers.ASNExists(tx, id)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
		return
	}
	if !exists {
		if id != "" {
			api.HandleErr(w, r, tx, http.StatusNotFound, fmt.Errorf("no asn exists by id: %s", id), nil)
			return
		} else {
			api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("no asn exists for empty id"), nil)
			return
		}
	}

	res, err := tx.Exec("DELETE FROM asn WHERE id=$1", id)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("determining rows affected for delete asn: %w", err))
		return
	}
	if rowsAffected == 0 {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, fmt.Errorf("no rows deleted for asn"), nil)
		return
	}
	alerts := tc.CreateAlerts(tc.SuccessLevel, "asn was deleted.")
	api.WriteAlerts(w, r, http.StatusOK, alerts)
	changeLogMsg := fmt.Sprintf("ID: %s, ACTION: Deleted asn", id)
	api.CreateChangeLogRawTx(api.Deleted, changeLogMsg, inf.User, tx)
	return
}

// readAndValidateJsonStruct reads json body and validates json fields
func readAndValidateJsonStruct(r *http.Request) (tc.ASNV5, error) {
	var asn tc.ASNV5
	if err := json.NewDecoder(r.Body).Decode(&asn); err != nil {
		userErr := fmt.Errorf("error decoding POST request body into ASNV5 struct %w", err)
		return asn, userErr
	}

	// validate JSON body
	errs := tovalidate.ToErrors(validation.Errors{
		"asn":          validation.Validate(asn.ASN, validation.NotNil, validation.Min(0)),
		"cachegroupId": validation.Validate(asn.CachegroupID, validation.NotNil, validation.Min(0)),
	})
	if len(errs) > 0 {
		userErr := util.JoinErrs(errs)
		return asn, userErr
	}
	return asn, nil
}

// selectMaxLastUpdatedQuery used for TryIfModifiedSinceQuery()
func selectMaxLastUpdatedQuery(where string) string {
	return `SELECT max(t) from (
		SELECT max(a.last_updated) as t from asn a
		JOIN cachegroup c ON a.cachegroup = c.id ` + where +
		` UNION ALL
	select max(last_updated) as t from last_deleted l where l.table_name='asn') as res`
}
