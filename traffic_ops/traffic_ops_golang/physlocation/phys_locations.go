package physlocation

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

// we need a type alias to define functions on
type TOPhysLocation struct {
	api.APIInfoImpl `json:"-"`
	tc.PhysLocationNullable
}

func (v *TOPhysLocation) GetLastUpdated() (*time.Time, bool, error) {
	return api.GetLastUpdated(v.APIInfo().Tx, *v.ID, "phys_location")
}

func (v *TOPhysLocation) SetLastUpdated(t tc.TimeNoMod) { v.LastUpdated = &t }
func (v *TOPhysLocation) InsertQuery() string           { return insertQuery() }
func (v *TOPhysLocation) NewReadObj() interface{}       { return &tc.PhysLocationNullable{} }
func (v *TOPhysLocation) SelectQuery() string           { return selectQuery() }
func (v *TOPhysLocation) ParamColumns() map[string]dbhelpers.WhereColumnInfo {
	return map[string]dbhelpers.WhereColumnInfo{
		"name":   dbhelpers.WhereColumnInfo{Column: "pl.name"},
		"id":     dbhelpers.WhereColumnInfo{Column: "pl.id", Checker: api.IsInt},
		"region": dbhelpers.WhereColumnInfo{Column: "pl.region", Checker: api.IsInt},
	}
}
func (v *TOPhysLocation) UpdateQuery() string { return updateQuery() }
func (v *TOPhysLocation) DeleteQuery() string { return deleteQuery() }

func (pl TOPhysLocation) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{Field: "id", Func: api.GetIntKey}}
}

// Implementation of the Identifier, Validator interface functions
func (pl TOPhysLocation) GetKeys() (map[string]interface{}, bool) {
	if pl.ID == nil {
		return map[string]interface{}{"id": 0}, false
	}
	return map[string]interface{}{"id": *pl.ID}, true
}

func (pl *TOPhysLocation) SetKeys(keys map[string]interface{}) {
	i, _ := keys["id"].(int) //this utilizes the non panicking type assertion, if the thrown away ok variable is false i will be the zero of the type, 0 here.
	pl.ID = &i
}

func (pl *TOPhysLocation) GetAuditName() string {
	if pl.Name != nil {
		return *pl.Name
	}
	if pl.ID != nil {
		return strconv.Itoa(*pl.ID)
	}
	return "unknown"
}

func (pl *TOPhysLocation) GetType() string {
	return "physLocation"
}

func (pl *TOPhysLocation) Validate() (error, error) {
	errs := validation.Errors{
		"address":   validation.Validate(pl.Address, validation.Required),
		"city":      validation.Validate(pl.City, validation.Required),
		"name":      validation.Validate(pl.Name, validation.Required),
		"regionId":  validation.Validate(pl.RegionID, validation.Required, validation.Min(0)),
		"shortName": validation.Validate(pl.ShortName, validation.Required),
		"state":     validation.Validate(pl.State, validation.Required),
		"zip":       validation.Validate(pl.Zip, validation.Required),
	}
	if errs != nil {
		return util.JoinErrs(tovalidate.ToErrors(errs)), nil
	}
	return nil, nil
}

func (pl *TOPhysLocation) Read(h http.Header, useIMS bool) ([]interface{}, error, error, int, *time.Time) {
	api.DefaultSort(pl.APIInfo(), "name")
	return api.GenericRead(h, pl, useIMS)
}
func (v *TOPhysLocation) SelectMaxLastUpdatedQuery(where, orderBy, pagination, tableName string) string {
	return `SELECT max(t) from (
		SELECT max(pl.last_updated) as t FROM phys_location pl
JOIN region r ON pl.region = r.id ` + where + orderBy + pagination +
		` UNION ALL
	select max(last_updated) as t from last_deleted l where l.table_name='phys_location') as res`
}

// MatchRegionNameAndID checks to see if the supplied region name and ID in the phys_location body correspond to each other.
func (pl *TOPhysLocation) MatchRegionNameAndID() (error, error, int) {
	if pl.RegionName != nil {
		regionName, ok, err := dbhelpers.GetRegionNameFromID(pl.APIInfo().Tx.Tx, *pl.RegionID)
		if err != nil {
			return nil, fmt.Errorf("error fetching name from region ID: %w", err), http.StatusInternalServerError
		} else if !ok {
			return errors.New("no such region"), nil, http.StatusNotFound
		}
		if regionName != *pl.RegionName {
			return errors.New("region name and ID do not match"), nil, http.StatusBadRequest
		}
	}
	return nil, nil, http.StatusOK
}

func (pl *TOPhysLocation) Update(h http.Header) (error, error, int) { return api.GenericUpdate(h, pl) }
func (pl *TOPhysLocation) Create() (error, error, int) {
	if userErr, sysErr, statusCode := pl.MatchRegionNameAndID(); userErr != nil || sysErr != nil {
		return userErr, sysErr, statusCode
	}
	return api.GenericCreate(pl)
}
func (pl *TOPhysLocation) Delete() (error, error, int) { return api.GenericDelete(pl) }

func selectQuery() string {
	return `
SELECT
pl.address,
pl.city,
pl.comments,
pl.email,
pl.id,
pl.last_updated,
pl.name,
pl.phone,
pl.poc,
r.id as region,
r.name as region_name,
pl.short_name,
pl.state,
pl.zip
FROM phys_location pl
JOIN region r ON pl.region = r.id
`
}

func updateQuery() string {
	query := `UPDATE
phys_location SET
address=:address,
city=:city,
comments=:comments,
email=:email,
name=:name,
phone=:phone,
poc=:poc,
region=:region,
short_name=:short_name,
state=:state,
zip=:zip
WHERE id=:id RETURNING last_updated`
	return query
}

func insertQuery() string {
	query := `INSERT INTO phys_location (
address,
city,
comments,
email,
name,
phone,
poc,
region,
short_name,
state,
zip) VALUES (
:address,
:city,
:comments,
:email,
:name,
:phone,
:poc,
:region,
:short_name,
:state,
:zip) RETURNING id,last_updated`
	return query
}

func deleteQuery() string {
	query := `DELETE FROM phys_location
WHERE id=:id`
	return query
}

func GetPhysLocation(w http.ResponseWriter, r *http.Request) {
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
		"name":   {Column: "pl.name", Checker: nil},
		"id":     {Column: "pl.id", Checker: api.IsInt},
		"region": {Column: "pl.region", Checker: api.IsInt},
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

	query := selectQuery() + where + orderBy + pagination
	rows, err := tx.NamedQuery(query, queryValues)
	if err != nil {
		api.HandleErr(w, r, tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("Phy_Location read: error getting Phy_Location(s): %w", err))
		return
	}
	defer log.Close(rows, "unable to close DB connection")

	physLocation := tc.PhysLocationNullableV5{}
	physLocationList := []tc.PhysLocationNullableV5{}
	for rows.Next() {
		if err = rows.Scan(&physLocation.Address, &physLocation.City, &physLocation.Comments, &physLocation.Email, &physLocation.ID, &physLocation.LastUpdated, &physLocation.Name, &physLocation.Phone, &physLocation.POC, &physLocation.RegionID, &physLocation.RegionName, &physLocation.ShortName, &physLocation.State, &physLocation.Zip); err != nil {
			api.HandleErr(w, r, tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("error getting physLocation(s): %w", err))
			return
		}
		physLocationList = append(physLocationList, physLocation)
	}

	api.WriteResp(w, r, physLocationList)
	return
}

func CreatePhysLocation(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	tx := inf.Tx.Tx

	physLocation, readValErr := readAndValidateJsonStruct(r)
	if readValErr != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, readValErr, nil)
		return
	}

	// checks to see if the supplied region name and ID in the phys_location body correspond to each other.
	if physLocation.RegionName != "" {
		regionName, ok, err := dbhelpers.GetRegionNameFromID(tx, physLocation.RegionID)
		if err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, fmt.Errorf("error fetching name from region ID: %w", err), nil)
			return
		} else if !ok {
			api.HandleErr(w, r, tx, http.StatusNotFound, errors.New("no such region"), nil)
			return
		}
		if regionName != physLocation.RegionName {
			api.HandleErr(w, r, tx, http.StatusBadRequest, errors.New("region name and ID do not match"), nil)
			return
		}
	}

	// check if phys_location already exists
	var count int
	err := tx.QueryRow("SELECT count(*) from phys_location where name = $1", physLocation.Name).Scan(&count)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("error: %w, when checking if physLocation with name %s exists", err, physLocation.Name))
		return
	}
	if count == 1 {
		api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("physLocation name '%s' already exists", physLocation.Name), nil)
		return
	}

	// create phys_location
	query := `
	INSERT INTO phys_location (
		address, 
		city, 
		comments, 
		email, 
		name, 
		phone, 
		poc, 
		region, 
		short_name, 
		state, 
		zip
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
	) RETURNING id, last_updated
`
	err = tx.QueryRow(
		query,
		physLocation.Address,
		physLocation.City,
		physLocation.Comments,
		physLocation.Email,
		physLocation.Name,
		physLocation.Phone,
		physLocation.POC,
		physLocation.RegionID,
		physLocation.ShortName,
		physLocation.State,
		physLocation.Zip,
	).Scan(
		&physLocation.ID,
		&physLocation.LastUpdated,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, fmt.Errorf("error: %w in physLocation  with name: %s", err, physLocation.Name), nil)
			return
		}
		usrErr, sysErr, code := api.ParseDBError(err)
		api.HandleErr(w, r, tx, code, usrErr, sysErr)
		return
	}
	alerts := tc.CreateAlerts(tc.SuccessLevel, "physLocation was created.")
	w.Header().Set(rfc.Location, fmt.Sprintf("/api/%s/phys_locations?name=%s", inf.Version, physLocation.Name))
	api.WriteAlertsObj(w, r, http.StatusCreated, alerts, physLocation)
	changeLogMsg := fmt.Sprintf("PHYSLOCATION: %s, ID:%d, ACTION: Created physLocation", physLocation.Name, physLocation.ID)
	api.CreateChangeLogRawTx(api.Created, changeLogMsg, inf.User, tx)
	return
}

func UpdatePhysLocation(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	tx := inf.Tx.Tx
	physLocation, readValErr := readAndValidateJsonStruct(r)
	if readValErr != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, readValErr, nil)
		return
	}

	requestedID := inf.Params["id"]

	intRequestId, convErr := strconv.Atoi(requestedID)
	if convErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, fmt.Errorf("physLocation update error: %w, while converting from string to int", convErr), nil)
	}
	// check if the entity was already updated
	userErr, sysErr, errCode = api.CheckIfUnModified(r.Header, inf.Tx, intRequestId, "phys_location")
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	//update name and description of a phys_location
	query := `
	UPDATE phys_location pl
	SET
		address = $1,
		city = $2,
		comments = $3,
		email = $4,
		name = $5,
		phone = $6,
		poc = $7,
		region = $8,
		short_name = $9,
		state = $10,
		zip = $11
	WHERE
		pl.id = $12
	RETURNING
		pl.id, pl.last_updated
`

	err := tx.QueryRow(
		query,
		physLocation.Address,
		physLocation.City,
		physLocation.Comments,
		physLocation.Email,
		physLocation.Name,
		physLocation.Phone,
		physLocation.POC,
		physLocation.RegionID,
		physLocation.ShortName,
		physLocation.State,
		physLocation.Zip,
		requestedID,
	).Scan(
		&physLocation.ID,
		&physLocation.LastUpdated,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("physLocation with ID: %v not found", physLocation.ID), nil)
			return
		}
		usrErr, sysErr, code := api.ParseDBError(err)
		api.HandleErr(w, r, tx, code, usrErr, sysErr)
		return
	}
	alerts := tc.CreateAlerts(tc.SuccessLevel, "physLocation was updated")
	api.WriteAlertsObj(w, r, http.StatusOK, alerts, physLocation)
	changeLogMsg := fmt.Sprintf("PHYSLOCATION: %s, ID:%d, ACTION: Updated physLocation", physLocation.Name, physLocation.ID)
	api.CreateChangeLogRawTx(api.Updated, changeLogMsg, inf.User, tx)
	return
}

func DeletePhysLocation(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	id := inf.Params["id"]
	if id == "" {
		api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("couldn't delete Phys_Location. Invalid ID. Id Cannot be empty for Delete Operation"), nil)
		return
	}
	exists, err := dbhelpers.PhysLocationExists(tx, id)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
		return
	}
	if !exists {
		api.HandleErr(w, r, tx, http.StatusNotFound, fmt.Errorf("no PhysLocation exists by id: %s", id), nil)
		return
	}

	assignedServer := 0
	if err := inf.Tx.Get(&assignedServer, "SELECT count(id) FROM server sv WHERE sv.phys_location=$1", id); err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("phys_location delete error, could not count assigned servers: %w", err))
		return
	} else if assignedServer != 0 {
		api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("can not delete a phys_location with %d assigned servers", assignedServer), nil)
		return
	}

	res, err := tx.Exec("DELETE FROM phys_location AS pl WHERE pl.id=$1", id)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("determining rows affected for delete phys_location: %w", err))
		return
	}
	if rowsAffected == 0 {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, fmt.Errorf("no rows deleted for phys_location"), nil)
		return
	}
	alerts := tc.CreateAlerts(tc.SuccessLevel, "phys_location was deleted.")
	api.WriteAlerts(w, r, http.StatusOK, alerts)
	changeLogMsg := fmt.Sprintf("ID:%s, ACTION: Deleted physLocation", id)
	api.CreateChangeLogRawTx(api.Deleted, changeLogMsg, inf.User, tx)
	return
}

func readAndValidateJsonStruct(r *http.Request) (tc.PhysLocationV5, error) {
	var physLocation tc.PhysLocationV5
	if err := json.NewDecoder(r.Body).Decode(&physLocation); err != nil {
		userErr := fmt.Errorf("error decoding POST request body into PhysLocationV5 struct %w", err)
		return physLocation, userErr
	}

	// validate JSON body
	errs := tovalidate.ToErrors(validation.Errors{
		"address":   validation.Validate(physLocation.Address, validation.Required),
		"city":      validation.Validate(physLocation.City, validation.Required),
		"name":      validation.Validate(physLocation.Name, validation.Required),
		"regionId":  validation.Validate(physLocation.RegionID, validation.Required, validation.Min(0)),
		"shortName": validation.Validate(physLocation.ShortName, validation.Required),
		"state":     validation.Validate(physLocation.State, validation.Required),
		"zip":       validation.Validate(physLocation.Zip, validation.Required)})
	if len(errs) > 0 {
		userErr := util.JoinErrs(errs)
		return physLocation, userErr
	}
	return physLocation, nil
}

// selectMaxLastUpdatedQuery used for TryIfModifiedSinceQuery()
func selectMaxLastUpdatedQuery(where string) string {
	return `SELECT max(t) from (
		SELECT max(a.last_updated) as t from phys_location a` + where +
		` UNION ALL
	select max(last_updated) as t from last_deleted l where l.table_name='phys_location') as res`
}
