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
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc/v13"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/tovalidate"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type TOCacheGroup v13.CacheGroupNullable

//the refType is passed into the handlers where a copy of its type is used to decode the json.
var refType = TOCacheGroup{}

func GetRefType() *TOCacheGroup {
	return &refType
}

func (cachegroup TOCacheGroup) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{"id", api.GetIntKey}}
}

//Implementation of the Identifier, Validator interface functions
func (cachegroup TOCacheGroup) GetKeys() (map[string]interface{}, bool) {
	if cachegroup.ID == nil {
		return map[string]interface{}{"id": 0}, false
	}
	return map[string]interface{}{"id": *cachegroup.ID}, true
}

func (cachegroup *TOCacheGroup) SetKeys(keys map[string]interface{}) {
	i, _ := keys["id"].(int) //this utilizes the non panicking type assertion, if the thrown away ok variable is false i will be the zero of the type, 0 here.
	cachegroup.ID = &i
}

//Implementation of the Identifier, Validator interface functions
func (cachegroup TOCacheGroup) GetID() (int, bool) {
	if cachegroup.ID == nil {
		return 0, false
	}
	return *cachegroup.ID, true
}

func (cachegroup TOCacheGroup) GetAuditName() string {
	if cachegroup.Name != nil {
		return *cachegroup.Name
	}
	id, _ := cachegroup.GetID()
	return strconv.Itoa(id)
}

func (cachegroup TOCacheGroup) GetType() string {
	return "cachegroup"
}

func (cachegroup *TOCacheGroup) SetID(i int) {
	cachegroup.ID = &i
}

// checks if a cachegroup with the given ID is in use as a parent or secondary parent.
func isUsedByChildCache(db *sqlx.DB, ID int) (bool, error) {
	pQuery := "SELECT count(*) from cachegroup WHERE parent_cachegroup_id=$1"
	sQuery := "SELECT count(*) from cachegroup WHERE secondary_parent_cachegroup_id=$1"
	count := 0

	err := db.QueryRow(pQuery, ID).Scan(&count)
	if err != nil {
		log.Errorf("received error: %++v from query execution", err)
		return false, err
	}
	if count > 0 {
		return true, errors.New("cache is in use as a parent cache")
	}

	err = db.QueryRow(sQuery, ID).Scan(&count)
	if err != nil {
		log.Errorf("received error: %++v from query execution", err)
		return false, err
	}
	if count > 0 {
		return true, errors.New("cache is in use as a secondary parent cache")
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
func (cachegroup TOCacheGroup) Validate(db *sqlx.DB) []error {
	validName := validation.NewStringRule(IsValidCacheGroupName, "invalid characters found - Use alphanumeric . or - or _ .")
	validShortName := validation.NewStringRule(IsValidCacheGroupName, "invalid characters found - Use alphanumeric . or - or _ .")
	latitudeErr := "Must be a floating point number within the range +-90"
	longitudeErr := "Must be a floating point number within the range +-180"
	errs := validation.Errors{
		"name":                        validation.Validate(cachegroup.Name, validation.Required, validName),
		"shortName":                   validation.Validate(cachegroup.ShortName, validation.Required, validShortName),
		"latitude":                    validation.Validate(cachegroup.Latitude, validation.Min(-90.0).Error(latitudeErr), validation.Max(90.0).Error(latitudeErr)),
		"longitude":                   validation.Validate(cachegroup.Longitude, validation.Min(-180.0).Error(longitudeErr), validation.Max(180.0).Error(longitudeErr)),
		"parentCacheGroupID":          validation.Validate(cachegroup.ParentCachegroupID, validation.Min(1)),
		"secondaryParentCachegroupID": validation.Validate(cachegroup.SecondaryParentCachegroupID, validation.Min(1)),
	}
	return tovalidate.ToErrors(errs)
}

// looks up the parent_cachegroup_id and the secondary_cachegroup_id
// if the respective names are defined in the cachegroup struct.  A
// sucessful lookup sets the two ids on the struct.
//
// used by Create()
func getParentCachegroupIDs(db *sqlx.DB, cachegroup *TOCacheGroup) error {
	query := `SELECT id FROM cachegroup where name=$1`
	var parentID int
	var secondaryParentID int

	if cachegroup.ParentName != nil && *cachegroup.ParentName != "" {
		err := db.QueryRow(query, *cachegroup.ParentName).Scan(&parentID)
		if err != nil {
			log.Errorf("received error: %++v from query execution", err)
			return err
		}
		cachegroup.ParentCachegroupID = &parentID
	}
	// not using 'omitempty' on the CacheGroup struct so a '0' is really an empty field, so set the pointer to nil
	if cachegroup.ParentCachegroupID != nil && *cachegroup.ParentCachegroupID == 0 {
		cachegroup.ParentCachegroupID = nil
	}

	if cachegroup.SecondaryParentName != nil && *cachegroup.SecondaryParentName != "" {
		err := db.QueryRow(query, *cachegroup.SecondaryParentName).Scan(&secondaryParentID)
		if err != nil {
			log.Errorf("received error: %++v from query execution", err)
			return err
		}
		cachegroup.SecondaryParentCachegroupID = &secondaryParentID
	}
	// not using 'omitempty' on the CacheGroup struct so a '0' is really an empty field, so set the pointer to nil
	if cachegroup.SecondaryParentCachegroupID != nil && *cachegroup.SecondaryParentCachegroupID == 0 {
		cachegroup.SecondaryParentCachegroupID = nil
	}
	return nil
}

// looks up the parent and secondary cachegroup names by cachegroup ID.
//  the names are set on the struct.
//
// used by Read()
func getParentCacheGroupNames(db *sqlx.DB, cachegroup *TOCacheGroup) error {
	query1 := `SELECT name FROM cachegroup where id=$1`
	var primaryName string
	var secondaryName string

	// primary parent lookup
	if cachegroup.ParentCachegroupID != nil {
		err := db.QueryRow(query1, *cachegroup.ParentCachegroupID).Scan(&primaryName)
		if err != nil {
			log.Errorf("received error: %++v from query execution", err)
			return err
		}
		cachegroup.ParentName = &primaryName
	}

	// secondary parent lookup
	if cachegroup.SecondaryParentCachegroupID != nil {
		err := db.QueryRow(query1, *cachegroup.SecondaryParentCachegroupID).Scan(&secondaryName)
		if err != nil {
			log.Errorf("received error: %++v from query execution", err)
			return err
		}
		cachegroup.SecondaryParentName = &secondaryName
	}

	return nil
}

//The TOCacheGroup implementation of the Creator interface
//all implementations of Creator should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a cachegroup with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
//The insert sql returns the id and lastUpdated values of the newly inserted cachegroup and have
//to be added to the struct
func (cachegroup *TOCacheGroup) Create(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
	rollbackTransaction := true
	tx, err := db.Beginx()
	defer func() {
		if tx == nil || !rollbackTransaction {
			return
		}
		err := tx.Rollback()
		if err != nil {
			log.Errorln(errors.New("rolling back transaction: " + err.Error()))
		}
	}()

	if err != nil {
		log.Error.Printf("could not begin transaction: %v", err)
		return tc.DBError, tc.SystemError
	}

	err = getParentCachegroupIDs(db, cachegroup)
	if err != nil {
		log.Error.Printf("failure looking up parent cache groups %v", err)
		return tc.DBError, tc.SystemError
	}

	resultRows, err := tx.NamedQuery(insertQuery(), cachegroup)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			err, eType := dbhelpers.ParsePQUniqueConstraintError(pqErr)
			if eType == tc.DataConflictError {
				return errors.New("a cachegroup with " + err.Error()), eType
			}
			return err, eType
		} else {
			log.Errorf("received non pq error: %++v from create execution", err)
			return tc.DBError, tc.SystemError
		}
	}
	defer resultRows.Close()

	var id int
	var lastUpdated tc.TimeNoMod
	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err := resultRows.Scan(&id, &lastUpdated); err != nil {
			log.Error.Printf("could not scan id from insert: %s\n", err)
			return tc.DBError, tc.SystemError
		}
	}
	if rowsAffected == 0 {
		err = errors.New("no cachegroup was inserted, no id was returned")
		log.Errorln(err)
		return tc.DBError, tc.SystemError
	} else if rowsAffected > 1 {
		err = errors.New("too many ids returned from cachegroup insert")
		log.Errorln(err)
		return tc.DBError, tc.SystemError
	}
	cachegroup.SetID(id)
	cachegroup.LastUpdated = &lastUpdated
	err = tx.Commit()
	if err != nil {
		log.Errorln("Could not commit transaction: ", err)
		return tc.DBError, tc.SystemError
	}
	rollbackTransaction = false
	return nil, tc.NoError
}

func (cachegroup *TOCacheGroup) Read(db *sqlx.DB, parameters map[string]string, user auth.CurrentUser) ([]interface{}, []error, tc.ApiErrorType) {
	var rows *sqlx.Rows

	// Query Parameters to Database Query column mappings
	// see the fields mapped in the SQL query
	queryParamsToQueryCols := map[string]dbhelpers.WhereColumnInfo{
		"id":        dbhelpers.WhereColumnInfo{"cachegroup.id", api.IsInt},
		"name":      dbhelpers.WhereColumnInfo{"cachegroup.name", nil},
		"shortName": dbhelpers.WhereColumnInfo{"short_name", nil},
	}
	where, orderBy, queryValues, errs := dbhelpers.BuildWhereAndOrderBy(parameters, queryParamsToQueryCols)
	if len(errs) > 0 {
		return nil, errs, tc.DataConflictError
	}

	query := selectQuery() + where + orderBy
	log.Debugln("Query is ", query)

	rows, err := db.NamedQuery(query, queryValues)
	if err != nil {
		log.Errorf("Error querying CacheGroup: %v", err)
		return nil, []error{tc.DBError}, tc.SystemError
	}
	defer rows.Close()

	CacheGroups := []interface{}{}
	for rows.Next() {
		var s TOCacheGroup
		if err = rows.StructScan(&s); err != nil {
			log.Errorf("error parsing CacheGroup rows: %v", err)
			return nil, []error{tc.DBError}, tc.SystemError
		}
		getParentCacheGroupNames(db, &s)
		CacheGroups = append(CacheGroups, s)
	}

	return CacheGroups, []error{}, tc.NoError
}

//The TOCacheGroup implementation of the Updater interface
//all implementations of Updater should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a cachegroup with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
func (cachegroup *TOCacheGroup) Update(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
	rollbackTransaction := true
	tx, err := db.Beginx()
	defer func() {
		if tx == nil || !rollbackTransaction {
			return
		}
		err := tx.Rollback()
		if err != nil {
			log.Errorln(errors.New("rolling back transaction: " + err.Error()))
		}
	}()

	if err != nil {
		log.Error.Printf("could not begin transaction: %v", err)
		return tc.DBError, tc.SystemError
	}

	// fix up parent ids.
	err = getParentCachegroupIDs(db, cachegroup)
	if err != nil {
		log.Error.Printf("failure looking up parent cache groups %v", err)
		return tc.DBError, tc.SystemError
	}

	log.Debugf("about to run exec query: %s with cachegroup: %++v", updateQuery(), cachegroup)
	resultRows, err := tx.NamedQuery(updateQuery(), cachegroup)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			err, eType := dbhelpers.ParsePQUniqueConstraintError(pqErr)
			if eType == tc.DataConflictError {
				return errors.New("a cachegroup with " + err.Error()), eType
			}
			return err, eType
		} else {
			log.Errorf("received error: %++v from update execution", err)
			return tc.DBError, tc.SystemError
		}
	}
	defer resultRows.Close()

	var lastUpdated tc.TimeNoMod
	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err := resultRows.Scan(&lastUpdated); err != nil {
			log.Error.Printf("could not scan lastUpdated from insert: %s\n", err)
			return tc.DBError, tc.SystemError
		}
	}
	log.Debugf("lastUpdated: %++v", lastUpdated)
	cachegroup.LastUpdated = &lastUpdated
	if rowsAffected != 1 {
		if rowsAffected < 1 {
			return errors.New("no cachegroup found with this id"), tc.DataMissingError
		} else {
			return fmt.Errorf("this update affected too many rows: %d", rowsAffected), tc.SystemError
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Errorln("Could not commit transaction: ", err)
		return tc.DBError, tc.SystemError
	}
	rollbackTransaction = false
	return nil, tc.NoError
}

//The CacheGroup implementation of the Deleter interface
//all implementations of Deleter should use transactions and return the proper errorType
func (cachegroup *TOCacheGroup) Delete(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
	rollbackTransaction := true
	tx, err := db.Beginx()
	defer func() {
		if tx == nil || !rollbackTransaction {
			return
		}
		err := tx.Rollback()
		if err != nil {
			log.Errorln(errors.New("rolling back transaction: " + err.Error()))
		}
	}()

	if err != nil {
		log.Error.Printf("could not begin transaction: %v", err)
		return tc.DBError, tc.SystemError
	}

	inUse, err := isUsedByChildCache(db, *cachegroup.ID)
	log.Debugf("inUse: %d, err: %v", inUse, err)
	if inUse == false && err != nil {
		return tc.DBError, tc.SystemError
	}
	if inUse == true {
		return err, tc.DataConflictError
	}

	log.Debugf("about to run exec query: %s with cachegroup: %++v", deleteQuery(), cachegroup)
	result, err := tx.NamedExec(deleteQuery(), cachegroup)
	if err != nil {
		log.Errorf("received error: %++v from delete execution", err)
		return tc.DBError, tc.SystemError
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return tc.DBError, tc.SystemError
	}
	if rowsAffected != 1 {
		if rowsAffected < 1 {
			return errors.New("no cachegroup with that id found"), tc.DataMissingError
		} else {
			return fmt.Errorf("this create affected too many rows: %d", rowsAffected), tc.SystemError
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Errorln("Could not commit transaction: ", err)
		return tc.DBError, tc.SystemError
	}
	rollbackTransaction = false
	return nil, tc.NoError
}

// insert query
func insertQuery() string {
	// to disambiguate struct scans, the named
	// parameter 'type_id' is an alias to cachegroup.type
	//see also the v13.CacheGroupNullable struct 'db' metadata
	query := `INSERT INTO cachegroup (
name,
short_name,
latitude,
longitude,
type,
parent_cachegroup_id,
secondary_parent_cachegroup_id
) VALUES(
:name,
:short_name,
:latitude,
:longitude,
:type_id,
:parent_cachegroup_id,
:secondary_parent_cachegroup_id
) RETURNING id,last_updated`
	return query
}

// select query
func selectQuery() string {
	// the 'type_name' and 'type_id' aliases on the 'type.name'
	// and cachegroup.type' fields are needed
	// to disambiguate the struct scan, see also the
	// v13.CacheGroupNullable struct 'db' metadata
	query := `SELECT
cachegroup.id,
cachegroup.name,
cachegroup.short_name,
cachegroup.latitude,
cachegroup.longitude,
cachegroup.parent_cachegroup_id,
cachegroup.secondary_parent_cachegroup_id,
type.name AS type_name,
cachegroup.type AS type_id,
cachegroup.last_updated
FROM cachegroup
INNER JOIN type ON cachegroup.type = type.id`
	return query
}

// update query
func updateQuery() string {
	// to disambiguate struct scans, the named
	// parameter 'type_id' is an alias to cachegroup.type
	//see also the v13.CacheGroupNullable struct 'db' metadata
	query := `UPDATE
cachegroup SET
name=:name,
short_name=:short_name,
latitude=:latitude,
longitude=:longitude,
parent_cachegroup_id=:parent_cachegroup_id,
secondary_parent_cachegroup_id=:secondary_parent_cachegroup_id,
type=:type_id WHERE id=:id RETURNING last_updated`
	return query
}

//delete query
func deleteQuery() string {
	query := `DELETE FROM cachegroup WHERE id=:id`
	return query
}
