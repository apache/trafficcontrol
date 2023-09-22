// Package coordinate contains API handlers and associated logic for servicing
// the `/coordinates` API endpoint.
package coordinate

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

// TOCoordinate is a "CRUDer"-based API wrapper for Coordinate objects.
// Deprecated: All future Coordinate versions should use the non-"CRUDer"
// methodology.
type TOCoordinate struct {
	api.APIInfoImpl `json:"-"`
	tc.CoordinateNullable
}

// SetLastUpdated implements a "CRUDer" interface.
// Deprecated: All future Coordinate versions should use the non-"CRUDer"
// methodology.
func (coordinate *TOCoordinate) SetLastUpdated(t tc.TimeNoMod) { coordinate.LastUpdated = &t }

// InsertQuery implements a "CRUDer" interface.
// Deprecated: All future Coordinate versions should use the non-"CRUDer"
// methodology.
func (*TOCoordinate) InsertQuery() string { return insertQuery() }

// NewReadObj implements a "CRUDer" interface.
// Deprecated: All future Coordinate versions should use the non-"CRUDer"
// methodology.
func (*TOCoordinate) NewReadObj() interface{} { return &tc.CoordinateNullable{} }

// SelectQuery implements a "CRUDer" interface.
// Deprecated: All future Coordinate versions should use the non-"CRUDer"
// methodology.
func (*TOCoordinate) SelectQuery() string { return selectQuery() }

// ParamColumns implements a "CRUDer" interface.
// Deprecated: All future Coordinate versions should use the non-"CRUDer"
// methodology.
func (*TOCoordinate) ParamColumns() map[string]dbhelpers.WhereColumnInfo {
	return map[string]dbhelpers.WhereColumnInfo{
		"id":   {Column: "id", Checker: api.IsInt},
		"name": {Column: "name"},
	}
}

// GetLastUpdated implements a "CRUDer" interface.
// Deprecated: All future Coordinate versions should use the non-"CRUDer"
// methodology.
func (coordinate *TOCoordinate) GetLastUpdated() (*time.Time, bool, error) {
	return api.GetLastUpdated(coordinate.APIInfo().Tx, *coordinate.ID, "coordinate")
}

// UpdateQuery implements a "CRUD"er interface.
// Deprecated: All future Coordinate versions should use the non-"CRUDer"
// methodology.
func (*TOCoordinate) UpdateQuery() string { return updateQuery() }

// DeleteQuery implements a "CRUD"er interface.
// Deprecated: All future Coordinate versions should use the non-"CRUDer"
// methodology.
func (*TOCoordinate) DeleteQuery() string { return deleteQuery() }

// GetKeyFieldsInfo implements a "CRUD"er interface.
// Deprecated: All future Coordinate versions should use the non-"CRUDer"
// methodology.
func (coordinate TOCoordinate) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{Field: "id", Func: api.GetIntKey}}
}

// GetKeys implements the Identifier and Validator interfaces.
// Deprecated: All future Coordinate versions should use the non-"CRUDer"
// methodology.
func (coordinate TOCoordinate) GetKeys() (map[string]interface{}, bool) {
	if coordinate.ID == nil {
		return map[string]interface{}{"id": 0}, false
	}
	return map[string]interface{}{"id": *coordinate.ID}, true
}

// GetAuditName implements a "CRUD"er interface.
// Deprecated: All future Coordinate versions should use the non-"CRUDer"
// methodology.
func (coordinate TOCoordinate) GetAuditName() string {
	if coordinate.Name != nil {
		return *coordinate.Name
	}
	if coordinate.ID != nil {
		return strconv.Itoa(*coordinate.ID)
	}
	return "0"
}

// GetType implements a "CRUDer" interface.
// Deprecated: All future Coordinate versions should use the non-"CRUDer"
// methodology.
func (coordinate TOCoordinate) GetType() string {
	return "coordinate"
}

// SetKeys implements a "CRUDer" interface.
// Deprecated: All future Coordinate versions should use the non-"CRUDer"
// methodology.
func (coordinate *TOCoordinate) SetKeys(keys map[string]interface{}) {
	i, _ := keys["id"].(int)
	coordinate.ID = &i
}

func isValidCoordinateChar(r rune) bool {
	if r >= 'a' && r <= 'z' {
		return true
	}
	if r >= 'A' && r <= 'Z' {
		return true
	}
	if r >= '0' && r <= '9' {
		return true
	}
	if r == '.' || r == '-' || r == '_' {
		return true
	}
	return false
}

// IsValidCoordinateName returns true if the name contains only characters valid
// for a Coordinate name.
// Deprecated: All future Coordinate versions should use the non-"CRUDer"
// methodology.
func IsValidCoordinateName(str string) bool {
	i := strings.IndexFunc(str, func(r rune) bool { return !isValidCoordinateChar(r) })
	return i == -1
}

// Validate fulfills the api.Validator interface.
// Deprecated: All future Coordinate versions should use non-"CRUDer"
// validation.
func (coordinate TOCoordinate) Validate() (error, error) {
	validName := validation.NewStringRule(IsValidCoordinateName, "invalid characters found - Use alphanumeric . or - or _ .")
	latitudeErr := "Must be a floating point number within the range +-90"
	longitudeErr := "Must be a floating point number within the range +-180"
	errs := validation.Errors{
		"name":      validation.Validate(coordinate.Name, validation.Required, validName),
		"latitude":  validation.Validate(coordinate.Latitude, validation.Min(-90.0).Error(latitudeErr), validation.Max(90.0).Error(latitudeErr)),
		"longitude": validation.Validate(coordinate.Longitude, validation.Min(-180.0).Error(longitudeErr), validation.Max(180.0).Error(longitudeErr)),
	}
	return util.JoinErrs(tovalidate.ToErrors(errs)), nil
}

// Create implements a "CRUDer" interface.
// Deprecated: All future Coordinate versions should use the non-"CRUDer" Create
// function.
func (coord *TOCoordinate) Create() (error, error, int) { return api.GenericCreate(coord) }

// Read implements a "CRUDer" interface.
// Deprecated: All future Coordinate versions should use the non-"CRUDer" Read
// function.
func (coord *TOCoordinate) Read(h http.Header, useIMS bool) ([]interface{}, error, error, int, *time.Time) {
	api.DefaultSort(coord.APIInfo(), "name")
	return api.GenericRead(h, coord, useIMS)
}

func selectMaxLastUpdatedQuery(where, orderBy, pagination string) string {
	return `
SELECT max(t) FROM (
	SELECT max(last_updated) AS t
	FROM (
		SELECT *
		FROM coordinate c
		` + where + orderBy + pagination +
		`	) AS coords
	UNION ALL
	SELECT max(last_updated) AS t
	FROM last_deleted l
	WHERE l.table_name='coordinate'
) AS res`
}

// SelectMaxLastUpdatedQuery implements a "CRUDer" interface.
// Deprecated: All future Coordinate versions should use the non-"CRUDer"
// methodology.
func (*TOCoordinate) SelectMaxLastUpdatedQuery(where, orderBy, pagination, _ string) string {
	return selectMaxLastUpdatedQuery(where, orderBy, pagination)
}

// Update implements a "CRUDer" interface.
// Deprecated: All future Coordinate versions should use the non-"CRUDer" Update
// function.
func (coord *TOCoordinate) Update(h http.Header) (error, error, int) {
	return api.GenericUpdate(h, coord)
}

// Delete implements a "CRUDer" interface.
// Deprecated: All future Coordinate versions should use the non-"CRUDer"
// methodology.
func (coord *TOCoordinate) Delete() (error, error, int) { return api.GenericDelete(coord) }

const readQuery = `
SELECT
	id,
	latitude,
	longitude,
	last_updated,
	name
FROM coordinate c`

func selectQuery() string {
	return readQuery
}

const putQuery = `
UPDATE coordinate
SET
	latitude=$1,
	longitude=$2,
	name=$3
WHERE id=$4
RETURNING
	last_updated`

func updateQuery() string {
	query := `UPDATE
coordinate SET
latitude=:latitude,
longitude=:longitude,
name=:name
WHERE id=:id RETURNING last_updated`
	return query
}

const createQuery = `
INSERT INTO coordinate (
	latitude,
	longitude,
	name
) VALUES (
	$1,
	$2,
	$3
) RETURNING
	id,
	last_updated`

func insertQuery() string {
	query := `INSERT INTO coordinate (
latitude,
longitude,
name) VALUES (
:latitude,
:longitude,
:name) RETURNING id,last_updated`
	return query
}

const delQuery = `
DELETE FROM coordinate
WHERE id = $1
RETURNING
	latitude,
	longitude,
	name,
	last_updated
`

func deleteQuery() string {
	return `DELETE FROM coordinate WHERE id = :id`
}

// Read is the handler for GET requests made to the /coordinates API (in APIv5
// and later).
func Read(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	cols := map[string]dbhelpers.WhereColumnInfo{
		"id":   {Column: "c.id", Checker: api.IsInt},
		"name": {Column: "c.name", Checker: nil},
	}
	api.DefaultSort(inf, "name")

	where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(inf.Params, cols)
	if len(errs) > 0 {
		errCode = http.StatusBadRequest
		userErr = util.JoinErrs(errs)
		api.HandleErr(w, r, tx, errCode, userErr, nil)
		return
	}

	var maxTime time.Time
	if inf.UseIMS() {
		var runSecond bool
		runSecond, maxTime = ims.TryIfModifiedSinceQuery(inf.Tx, r.Header, queryValues, selectMaxLastUpdatedQuery(where, orderBy, pagination))
		if !runSecond {
			log.Debugln("IMS HIT")
			api.WriteNotModifiedResponse(maxTime, w, r)
			return
		}
		log.Debugln("IMS MISS")
	} else {
		log.Debugln("Non IMS request")
	}

	query := readQuery + where + orderBy + pagination
	rows, err := inf.Tx.NamedQuery(query, queryValues)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("querying coordinates: %w", err))
		return
	}
	defer log.Close(rows, "closing coordinate query rows")

	cs := []tc.CoordinateV5{}
	for rows.Next() {
		var c tc.CoordinateV5
		err := rows.Scan(&c.ID, &c.Latitude, &c.Longitude, &c.LastUpdated, &c.Name)
		if err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("scanning a coordinate: %w", err))
			return
		}
		cs = append(cs, c)
	}

	api.WriteResp(w, r, cs)
}

// isValid returns an error describing why c isn't a valid Coordinate, or nil if
// it's actually valid.
func isValid(c tc.CoordinateV5) error {
	validName := validation.NewStringRule(IsValidCoordinateName, "invalid characters found - Use alphanumeric . or - or _ .")
	latitudeErr := "Must be a floating point number within the range +-90"
	longitudeErr := "Must be a floating point number within the range +-180"
	errs := validation.Errors{
		"name":      validation.Validate(c.Name, validation.Required, validName),
		"latitude":  validation.Validate(c.Latitude, validation.Min(-90.0).Error(latitudeErr), validation.Max(90.0).Error(latitudeErr)),
		"longitude": validation.Validate(c.Longitude, validation.Min(-180.0).Error(longitudeErr), validation.Max(180.0).Error(longitudeErr)),
	}
	return util.JoinErrs(tovalidate.ToErrors(errs))
}

// Create is the handler for POST requests made to the /coordinates API (in
// APIv5 and later).
func Create(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	var c tc.CoordinateV5
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
		return
	}

	if err = isValid(c); err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
		return
	}

	if err = tx.QueryRow(createQuery, c.Latitude, c.Longitude, c.Name).Scan(&c.ID, &c.LastUpdated); err != nil {
		userErr, sysErr, errCode = api.ParseDBError(err)
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	w.Header().Set(rfc.Location, fmt.Sprintf("/api/%s/coordinates?id=%d", inf.Version, *c.ID))
	w.WriteHeader(http.StatusCreated)
	api.WriteRespAlertObj(w, r, tc.SuccessLevel, fmt.Sprintf("Coordinate '%s' (#%d) created", c.Name, *c.ID), c)

	changeLogMsg := fmt.Sprintf("USER: %s, COORDINATE: %s (#%d), ACTION: %s", inf.User.UserName, c.Name, *c.ID, api.Created)
	api.CreateChangeLogRawTx(api.ApiChange, changeLogMsg, inf.User, tx)
}

// Update is the handler for PUT requests made to the /coordinates API (in API
// v5 and later).
func Update(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id"}, []string{"id"})
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	var c tc.CoordinateV5
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
		return
	}

	id := inf.IntParams["id"]
	if c.ID != nil {
		if *c.ID != id {
			api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("ID mismatch; URI specifies %d but payload is for Coordinate #%d", id, *c.ID), nil)
			return
		}
	} else {
		c.ID = util.Ptr(id)
	}

	userErr, sysErr, statusCode := api.CheckIfUnModified(r.Header, inf.Tx, id, "coordinate")
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, statusCode, userErr, sysErr)
		return
	}

	if err = isValid(c); err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
		return
	}

	if err = tx.QueryRow(putQuery, c.Latitude, c.Longitude, c.Name, id).Scan(&c.LastUpdated); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			userErr = fmt.Errorf("no such Coordinate: #%d", id)
			errCode = http.StatusNotFound
			sysErr = nil
		} else {
			userErr, sysErr, errCode = api.ParseDBError(err)
		}
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	api.WriteRespAlertObj(w, r, tc.SuccessLevel, fmt.Sprintf("Coordinate '%s' (#%d) updated", c.Name, id), c)

	changeLogMsg := fmt.Sprintf("USER: %s, COORDINATE: %s (#%d), ACTION: %s", inf.User.UserName, c.Name, id, api.Updated)
	api.CreateChangeLogRawTx(api.ApiChange, changeLogMsg, inf.User, tx)
}

// Delete is the handler for PUT requests made to the /coordinates API (in API
// v5 and later).
func Delete(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id"}, []string{"id"})
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	id := inf.IntParams["id"]

	c := tc.CoordinateV5{
		ID: util.Ptr(id),
	}
	if err := tx.QueryRow(delQuery, id).Scan(&c.Latitude, &c.Longitude, &c.Name, &c.LastUpdated); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			userErr = fmt.Errorf("no such Coordinate: #%d", id)
			errCode = http.StatusNotFound
			sysErr = nil
		} else {
			userErr, sysErr, errCode = api.ParseDBError(err)
		}
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	api.WriteRespAlertObj(w, r, tc.SuccessLevel, fmt.Sprintf("Coordinate '%s' (#%d) deleted", c.Name, id), c)

	changeLogMsg := fmt.Sprintf("USER: %s, COORDINATE: %s (#%d), ACTION: %s", inf.User.UserName, c.Name, id, api.Deleted)
	api.CreateChangeLogRawTx(api.ApiChange, changeLogMsg, inf.User, tx)
}
