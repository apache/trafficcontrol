package cachegroup

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
	"strings"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type TOCacheGroup struct {
	ReqInfo *api.APIInfo `json:"-"`
	tc.CacheGroupNullable
}

func GetTypeSingleton() api.CRUDFactory {
	return func(reqInfo *api.APIInfo) api.CRUDer {
		toReturn := TOCacheGroup{reqInfo, tc.CacheGroupNullable{}}
		return &toReturn
	}
}

func (cg TOCacheGroup) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{"id", api.GetIntKey}}
}

//Implementation of the Identifier, Validator interface functions
func (cg TOCacheGroup) GetKeys() (map[string]interface{}, bool) {
	if cg.ID == nil {
		return map[string]interface{}{"id": 0}, false
	}
	return map[string]interface{}{"id": *cg.ID}, true
}

func (cg *TOCacheGroup) SetKeys(keys map[string]interface{}) {
	i, _ := keys["id"].(int) //this utilizes the non panicking type assertion, if the thrown away ok variable is false i will be the zero of the type, 0 here.
	cg.ID = &i
}

//Implementation of the Identifier, Validator interface functions
func (cg TOCacheGroup) GetID() (int, bool) {
	if cg.ID == nil {
		return 0, false
	}
	return *cg.ID, true
}

func (cg TOCacheGroup) GetAuditName() string {
	if cg.Name != nil {
		return *cg.Name
	}
	id, _ := cg.GetID()
	return strconv.Itoa(id)
}

func (cg TOCacheGroup) GetType() string {
	return "cg"
}

func (cg *TOCacheGroup) SetID(i int) {
	cg.ID = &i
}

// Is the cachegroup being used?
func isUsed(tx *sqlx.Tx, ID int) (bool, error) {

	var usedByServer bool
	var usedByParent bool
	var usedBySecondaryParent bool
	var usedByASN bool

	query := `SELECT
    (SELECT id FROM server WHERE server.cachegroup = $1 LIMIT 1) IS NOT NULL,
    (SELECT id FROM cachegroup WHERE cachegroup.parent_cachegroup_id = $1 LIMIT 1) IS NOT NULL,
    (SELECT id FROM cachegroup WHERE cachegroup.secondary_parent_cachegroup_id = $1 LIMIT 1) IS NOT NULL,
    (SELECT id FROM asn WHERE cachegroup = $1 LIMIT 1) IS NOT NULL;`

	err := tx.QueryRow(query, ID).Scan(&usedByServer, &usedByParent, &usedBySecondaryParent, &usedByASN)
	if err != nil {
		log.Errorf("received error: %++v from query execution", err)
		return false, err
	}
	//Only return the immediate error
	if usedByServer {
		return true, errors.New("cachegroup is in use by one or more servers")
	}
	if usedByParent {
		return true, errors.New("cachegroup is in use as a parent cachegroup")
	}
	if usedBySecondaryParent {
		return true, errors.New("cachegroup is in use as a secondary parent cachegroup")
	}
	if usedByASN {
		return true, errors.New("cachegroup is in use in one or more ASNs")
	}

	return false, nil
}

func isValidCacheGroupChar(r rune) bool {
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

// IsValidCacheGroupName returns true if the name contains only characters valid for a CacheGroup name
func IsValidCacheGroupName(str string) bool {
	i := strings.IndexFunc(str, func(r rune) bool { return !isValidCacheGroupChar(r) })
	return i == -1
}

func IsValidParentCachegroupID(id *int) bool {
	if id == nil || *id > 0 {
		return true
	}
	return false
}

// Validate fulfills the api.Validator interface
func (cg TOCacheGroup) Validate() error {
	if _, err := tc.ValidateTypeID(cg.ReqInfo.Tx.Tx, cg.TypeID, "cachegroup"); err != nil {
		return err
	}

	validName := validation.NewStringRule(IsValidCacheGroupName, "invalid characters found - Use alphanumeric . or - or _ .")
	validShortName := validation.NewStringRule(IsValidCacheGroupName, "invalid characters found - Use alphanumeric . or - or _ .")
	latitudeErr := "Must be a floating point number within the range +-90"
	longitudeErr := "Must be a floating point number within the range +-180"
	errs := validation.Errors{
		"name":                        validation.Validate(cg.Name, validation.Required, validName),
		"shortName":                   validation.Validate(cg.ShortName, validation.Required, validShortName),
		"latitude":                    validation.Validate(cg.Latitude, validation.Min(-90.0).Error(latitudeErr), validation.Max(90.0).Error(latitudeErr)),
		"longitude":                   validation.Validate(cg.Longitude, validation.Min(-180.0).Error(longitudeErr), validation.Max(180.0).Error(longitudeErr)),
		"parentCacheGroupID":          validation.Validate(cg.ParentCachegroupID, validation.Min(1)),
		"secondaryParentCachegroupID": validation.Validate(cg.SecondaryParentCachegroupID, validation.Min(1)),
		"localizationMethods":         validation.Validate(cg.LocalizationMethods, validation.By(tovalidate.IsPtrToSliceOfUniqueStringersICase("CZ", "DEEP_CZ", "GEO"))),
	}
	return util.JoinErrs(tovalidate.ToErrors(errs))
}

//The TOCacheGroup implementation of the Creator interface
//all implementations of Creator should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a cachegroup with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
//The insert sql returns the id and lastUpdated values of the newly inserted cachegroup and have
//to be added to the struct
func (cg *TOCacheGroup) Create() (error, error, int) {
	coordinateID, err := cg.createCoordinate()
	if err != nil {
		return nil, errors.New("cg create: creating coord:" + err.Error()), http.StatusInternalServerError
	}

	resultRows, err := cg.ReqInfo.Tx.Tx.Query(
		insertQuery(),
		cg.Name,
		cg.ShortName,
		coordinateID,
		cg.TypeID,
		cg.ParentCachegroupID,
		cg.SecondaryParentCachegroupID,
	)
	if err != nil {
		return api.ParseDBError(err)
	}
	defer resultRows.Close()

	var id int
	var lastUpdated tc.TimeNoMod
	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err := resultRows.Scan(&id, &lastUpdated); err != nil {
			return nil, errors.New("cg create scanning: " + err.Error()), http.StatusInternalServerError
		}
	}
	if rowsAffected == 0 {
		return nil, errors.New("cg create: no rows returned"), http.StatusInternalServerError
	} else if rowsAffected > 1 {
		return nil, errors.New("cg create: multiple rows returned"), http.StatusInternalServerError
	}
	cg.SetID(id)
	if err = cg.createLocalizationMethods(); err != nil {
		return nil, errors.New("cg create: creating localization methods: " + err.Error()), http.StatusInternalServerError
	}
	cg.LastUpdated = &lastUpdated
	return nil, nil, http.StatusOK
}

func (cg *TOCacheGroup) createLocalizationMethods() error {
	q := `DELETE FROM cachegroup_localization_method where cachegroup = $1`
	if _, err := cg.ReqInfo.Tx.Tx.Exec(q, *cg.ID); err != nil {
		return fmt.Errorf("unable to delete cachegroup_localization_methods for cachegroup %d: %s", *cg.ID, err.Error())
	}
	if cg.LocalizationMethods != nil {
		q = `INSERT INTO cachegroup_localization_method (cachegroup, method) VALUES ($1, $2)`
		for _, method := range *cg.LocalizationMethods {
			if _, err := cg.ReqInfo.Tx.Tx.Exec(q, *cg.ID, method.String()); err != nil {
				return fmt.Errorf("unable to insert cachegroup_localization_methods for cachegroup %d: %s", *cg.ID, err.Error())
			}
		}
	}
	return nil
}

func (cg *TOCacheGroup) createCoordinate() (*int, error) {
	var coordinateID *int
	if cg.Latitude != nil && cg.Longitude != nil {
		q := `INSERT INTO coordinate (name, latitude, longitude) VALUES ($1, $2, $3) RETURNING id`
		if err := cg.ReqInfo.Tx.Tx.QueryRow(q, tc.CachegroupCoordinateNamePrefix+*cg.Name, *cg.Latitude, *cg.Longitude).Scan(&coordinateID); err != nil {
			return nil, fmt.Errorf("insert coordinate for cg '%s': %s", *cg.Name, err.Error())
		}
	}
	return coordinateID, nil
}

func (cg *TOCacheGroup) updateCoordinate() error {
	if cg.Latitude != nil && cg.Longitude != nil {
		q := `UPDATE coordinate SET name = $1, latitude = $2, longitude = $3 WHERE id = (SELECT coordinate FROM cachegroup WHERE id = $4)`
		result, err := cg.ReqInfo.Tx.Tx.Exec(q, tc.CachegroupCoordinateNamePrefix+*cg.Name, *cg.Latitude, *cg.Longitude, *cg.ID)
		if err != nil {
			return fmt.Errorf("update coordinate for cg '%s': %s", *cg.Name, err.Error())
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("update coordinate for cg '%s', getting rows affected: %s", *cg.Name, err.Error())
		}
		if rowsAffected == 0 {
			return fmt.Errorf("update coordinate for cg '%s', zero rows affected", *cg.Name)
		}
	}
	return nil
}

func (cg *TOCacheGroup) deleteCoordinate(coordinateID int) error {
	q := `UPDATE cachegroup SET coordinate = NULL WHERE id = $1`
	result, err := cg.ReqInfo.Tx.Tx.Exec(q, *cg.ID)
	if err != nil {
		return fmt.Errorf("updating cg %d coordinate to null: %s", *cg.ID, err.Error())
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("updating cg %d coordinate to null, getting rows affected: %s", *cg.ID, err.Error())
	}
	if rowsAffected == 0 {
		return fmt.Errorf("updating cg %d coordinate to null, zero rows affected", *cg.ID)
	}

	q = `DELETE FROM coordinate WHERE id = $1`
	result, err = cg.ReqInfo.Tx.Tx.Exec(q, coordinateID)
	if err != nil {
		return fmt.Errorf("delete coordinate %d for cg %d: %s", coordinateID, *cg.ID, err.Error())
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete coordinate %d for cg %d, getting rows affected: %s", coordinateID, *cg.ID, err.Error())
	}
	if rowsAffected == 0 {
		return fmt.Errorf("delete coordinate %d for cg %d, zero rows affected", coordinateID, *cg.ID)
	}
	return nil
}

func (cg *TOCacheGroup) Read() ([]interface{}, error, error, int) {
	// Query Parameters to Database Query column mappings
	// see the fields mapped in the SQL query
	queryParamsToQueryCols := map[string]dbhelpers.WhereColumnInfo{
		"id":        dbhelpers.WhereColumnInfo{"cachegroup.id", api.IsInt},
		"name":      dbhelpers.WhereColumnInfo{"cachegroup.name", nil},
		"shortName": dbhelpers.WhereColumnInfo{"short_name", nil},
		"type":      dbhelpers.WhereColumnInfo{"cachegroup.type", nil},
	}
	where, orderBy, queryValues, errs := dbhelpers.BuildWhereAndOrderBy(cg.ReqInfo.Params, queryParamsToQueryCols)
	if len(errs) > 0 {
		return nil, util.JoinErrs(errs), nil, http.StatusBadRequest
	}

	query := selectQuery() + where + orderBy
	log.Debugln("Query is ", query)

	rows, err := cg.ReqInfo.Tx.NamedQuery(query, queryValues)
	if err != nil {
		return nil, nil, errors.New("cg read: querying: " + err.Error()), http.StatusInternalServerError
	}
	defer rows.Close()

	cacheGroups := []interface{}{}
	for rows.Next() {
		var s TOCacheGroup
		lms := make([]tc.LocalizationMethod, 0)
		if err = rows.Scan(
			&s.ID,
			&s.Name,
			&s.ShortName,
			&s.Latitude,
			&s.Longitude,
			pq.Array(&lms),
			&s.ParentCachegroupID,
			&s.ParentName,
			&s.SecondaryParentCachegroupID,
			&s.SecondaryParentName,
			&s.Type,
			&s.TypeID,
			&s.LastUpdated,
		); err != nil {
			return nil, nil, errors.New("cg read: scanning: " + err.Error()), http.StatusInternalServerError
		}
		s.LocalizationMethods = &lms
		cacheGroups = append(cacheGroups, s)
	}

	return cacheGroups, nil, nil, http.StatusOK
}

//The TOCacheGroup implementation of the Updater interface
//all implementations of Updater should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a cachegroup with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
func (cg *TOCacheGroup) Update() (error, error, int) {
	coordinateID, err, errType := cg.handleCoordinateUpdate()
	if err != nil {
		return api.TypeErrToAPIErr(err, errType)
	}

	resultRows, err := cg.ReqInfo.Tx.Tx.Query(
		updateQuery(),
		cg.Name,
		cg.ShortName,
		coordinateID,
		cg.ParentCachegroupID,
		cg.SecondaryParentCachegroupID,
		cg.TypeID,
		cg.ID,
	)
	if err != nil {
		return api.ParseDBError(err)
	}
	defer resultRows.Close()

	var lastUpdated tc.TimeNoMod
	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err := resultRows.Scan(&lastUpdated); err != nil {
			return nil, errors.New("cg update: scanning: " + err.Error()), http.StatusInternalServerError
		}
	}
	log.Debugf("lastUpdated: %++v", lastUpdated)
	cg.LastUpdated = &lastUpdated
	if rowsAffected != 1 {
		if rowsAffected < 1 {
			return nil, nil, http.StatusNotFound
		} else {
			return nil, errors.New("cg update: affected multiple rows"), http.StatusInternalServerError
		}
	}
	if err = cg.createLocalizationMethods(); err != nil {
		return nil, errors.New("cg update: creating localization methods: " + err.Error()), http.StatusInternalServerError
	}
	return nil, nil, http.StatusOK
}

func (cg *TOCacheGroup) handleCoordinateUpdate() (*int, error, tc.ApiErrorType) {
	coordinateID, err := cg.getCoordinateID()
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no cg with id %d found", *cg.ID), tc.DataMissingError
		}
		log.Errorf("updating cg %d got error when querying coordinate: %s\n", *cg.ID, err)
		return nil, tc.DBError, tc.SystemError
	}
	if coordinateID == nil && cg.Latitude != nil && cg.Longitude != nil {
		newCoordinateID, err := cg.createCoordinate()
		if err != nil {
			log.Errorf("updating cg %d: %s\n", *cg.ID, err)
			return nil, tc.DBError, tc.SystemError
		}
		coordinateID = newCoordinateID
	} else if coordinateID != nil && (cg.Latitude == nil || cg.Longitude == nil) {
		if err = cg.deleteCoordinate(*coordinateID); err != nil {
			log.Errorf("updating cg %d: %s\n", *cg.ID, err)
			return nil, tc.DBError, tc.SystemError
		}
		coordinateID = nil
	} else {
		if err = cg.updateCoordinate(); err != nil {
			log.Errorf("updating cg %d: %s\n", *cg.ID, err)
			return nil, tc.DBError, tc.SystemError
		}
	}
	return coordinateID, nil, tc.NoError
}

func (cg *TOCacheGroup) getCoordinateID() (*int, error) {
	q := `SELECT coordinate FROM cachegroup WHERE id = $1`
	var coordinateID *int
	if err := cg.ReqInfo.Tx.Tx.QueryRow(q, *cg.ID).Scan(&coordinateID); err != nil {
		return nil, err
	}
	return coordinateID, nil
}

//The CacheGroup implementation of the Deleter interface
//all implementations of Deleter should use transactions and return the proper errorType
func (cg *TOCacheGroup) Delete() (error, error, int) {
	inUse, err := isUsed(cg.ReqInfo.Tx, *cg.ID)
	if err != nil {
		return nil, errors.New("cg delete: checking use: " + err.Error()), http.StatusInternalServerError
	}
	if inUse == true {
		return errors.New("cannot delete cachegroup in use"), nil, http.StatusInternalServerError
	}

	coordinateID, err := cg.getCoordinateID()
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, http.StatusNotFound
		}
		return nil, errors.New("cg delete: getting coord: " + err.Error()), http.StatusInternalServerError
	}

	if coordinateID != nil {
		if err = cg.deleteCoordinate(*coordinateID); err != nil {
			return nil, errors.New("cg delete: deleting coord: " + err.Error()), http.StatusInternalServerError
		}
	}

	result, err := cg.ReqInfo.Tx.NamedExec(deleteQuery(), cg)
	if err != nil {
		return nil, errors.New("cg delete querying: " + err.Error()), http.StatusInternalServerError
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, errors.New("cg delete getting rows affected: " + err.Error()), http.StatusInternalServerError
	}
	if rowsAffected != 1 {
		if rowsAffected < 1 {
			return nil, nil, http.StatusNotFound
		} else {
			return nil, errors.New("cg delete affected multiple rows"), http.StatusInternalServerError
		}
	}

	return nil, nil, http.StatusOK
}

// insert query
func insertQuery() string {
	query := `INSERT INTO cachegroup (
name,
short_name,
coordinate,
type,
parent_cachegroup_id,
secondary_parent_cachegroup_id
) VALUES($1,$2,$3,$4,$5,$6)
RETURNING id,last_updated`
	return query
}

// select query
func selectQuery() string {
	// the 'type_name' and 'type_id' aliases on the 'type.name'
	// and cachegroup.type' fields are needed
	// to disambiguate the struct scan, see also the
	// tc.CacheGroupNullable struct 'db' metadata
	query := `SELECT
cachegroup.id,
cachegroup.name,
cachegroup.short_name,
coordinate.latitude,
coordinate.longitude,
(SELECT array_agg(CAST(method as text)) AS localization_methods FROM cachegroup_localization_method clm WHERE clm.cachegroup = cachegroup.id),
cachegroup.parent_cachegroup_id,
cgp.name AS parent_cachegroup_name,
cachegroup.secondary_parent_cachegroup_id,
cgs.name AS secondary_parent_cachegroup_name,
type.name AS type_name,
cachegroup.type AS type_id,
cachegroup.last_updated
FROM cachegroup
LEFT JOIN coordinate ON coordinate.id = cachegroup.coordinate
INNER JOIN type ON cachegroup.type = type.id
LEFT JOIN cachegroup AS cgp ON cachegroup.parent_cachegroup_id = cgp.id
LEFT JOIN cachegroup AS cgs ON cachegroup.secondary_parent_cachegroup_id = cgs.id`
	return query
}

// update query
func updateQuery() string {
	// to disambiguate struct scans, the named
	// parameter 'type_id' is an alias to cachegroup.type
	//see also the tc.CacheGroupNullable struct 'db' metadata
	query := `UPDATE
cachegroup SET
name=$1,
short_name=$2,
coordinate=$3,
parent_cachegroup_id=$4,
secondary_parent_cachegroup_id=$5,
type=$6 WHERE id=$7 RETURNING last_updated`
	return query
}

//delete query
func deleteQuery() string {
	query := `DELETE FROM cachegroup WHERE id=:id`
	return query
}
