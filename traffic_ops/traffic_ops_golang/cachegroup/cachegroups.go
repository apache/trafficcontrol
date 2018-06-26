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

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/lib/go-tc/v13"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type TOCacheGroup struct{
	ReqInfo *api.APIInfo `json:"-"`
	v13.CacheGroupNullable
	}

func GetTypeSingleton() func(reqInfo *api.APIInfo) api.CRUDer {
	return func(reqInfo *api.APIInfo) api.CRUDer {
		toReturn := TOCacheGroup{reqInfo, v13.CacheGroupNullable{}}
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

// checks if a cachegroup with the given ID is in use as a parent or secondary parent.
func isUsedByChildCache(tx *sqlx.Tx, ID int) (bool, error) {
	pQuery := "SELECT count(*) from cachegroup WHERE parent_cachegroup_id=$1"
	sQuery := "SELECT count(*) from cachegroup WHERE secondary_parent_cachegroup_id=$1"
	count := 0

	err := tx.QueryRow(pQuery, ID).Scan(&count)
	if err != nil {
		log.Errorf("received error: %++v from query execution", err)
		return false, err
	}
	if count > 0 {
		return true, errors.New("cache is in use as a parent cache")
	}

	err = tx.QueryRow(sQuery, ID).Scan(&count)
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
func (cg TOCacheGroup) Validate() []error {
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
	}
	return tovalidate.ToErrors(errs)
}

// looks up the parent_cachegroup_id and the secondary_cachegroup_id
// if the respective names are defined in the cachegroup struct.  A
// sucessful lookup sets the two ids on the struct.
//
// used by Create()
func getParentCachegroupIDs(tx *sqlx.Tx, cachegroup *TOCacheGroup) error {
	query := `SELECT id FROM cachegroup where name=$1`
	var parentID int
	var secondaryParentID int

	if cachegroup.ParentName != nil && *cachegroup.ParentName != "" {
		err := tx.QueryRow(query, *cachegroup.ParentName).Scan(&parentID)
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
		err := tx.QueryRow(query, *cachegroup.SecondaryParentName).Scan(&secondaryParentID)
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
func getParentCacheGroupNames(tx *sqlx.Tx, cachegroup *TOCacheGroup) error {
	query1 := `SELECT name FROM cachegroup where id=$1`
	var primaryName string
	var secondaryName string

	// primary parent lookup
	if cachegroup.ParentCachegroupID != nil {
		err := tx.QueryRow(query1, *cachegroup.ParentCachegroupID).Scan(&primaryName)
		if err != nil {
			log.Errorf("received error: %++v from query execution", err)
			return err
		}
		cachegroup.ParentName = &primaryName
	}

	// secondary parent lookup
	if cachegroup.SecondaryParentCachegroupID != nil {
		err := tx.QueryRow(query1, *cachegroup.SecondaryParentCachegroupID).Scan(&secondaryName)
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
func (cg *TOCacheGroup) Create() (error, tc.ApiErrorType) {
	err := getParentCachegroupIDs(cg.ReqInfo.Tx, cg)
	if err != nil {
		log.Error.Printf("failure looking up parent cache groups %v", err)
		return tc.DBError, tc.SystemError
	}

	resultRows, err := cg.ReqInfo.Tx.NamedQuery(insertQuery(), cg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			err, eType := dbhelpers.ParsePQUniqueConstraintError(pqErr)
			if eType == tc.DataConflictError {
				return errors.New("a cg with " + err.Error()), eType
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
		err = errors.New("no cg was inserted, no id was returned")
		log.Errorln(err)
		return tc.DBError, tc.SystemError
	} else if rowsAffected > 1 {
		err = errors.New("too many ids returned from cg insert")
		log.Errorln(err)
		return tc.DBError, tc.SystemError
	}
	cg.SetID(id)
	cg.LastUpdated = &lastUpdated
	return nil, tc.NoError
}

func (cg *TOCacheGroup) Read(parameters map[string]string) ([]interface{}, []error, tc.ApiErrorType) {
	// Query Parameters to Database Query column mappings
	// see the fields mapped in the SQL query
	queryParamsToQueryCols := map[string]dbhelpers.WhereColumnInfo{
		"id":        dbhelpers.WhereColumnInfo{"cg.id", api.IsInt},
		"name":      dbhelpers.WhereColumnInfo{"cg.name", nil},
		"shortName": dbhelpers.WhereColumnInfo{"short_name", nil},
		"type":      dbhelpers.WhereColumnInfo{"cg.type", nil},
	}
	where, orderBy, queryValues, errs := dbhelpers.BuildWhereAndOrderBy(parameters, queryParamsToQueryCols)
	if len(errs) > 0 {
		return nil, errs, tc.DataConflictError
	}

	query := selectQuery() + where + orderBy
	log.Debugln("Query is ", query)

	rows, err := cg.ReqInfo.Tx.NamedQuery(query, queryValues)
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
		getParentCacheGroupNames(cg.ReqInfo.Tx, &s)
		CacheGroups = append(CacheGroups, s)
	}

	return CacheGroups, []error{}, tc.NoError
}

//The TOCacheGroup implementation of the Updater interface
//all implementations of Updater should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a cachegroup with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
func (cg *TOCacheGroup) Update() (error, tc.ApiErrorType) {
	// fix up parent ids.
	err := getParentCachegroupIDs(cg.ReqInfo.Tx, cg)
	if err != nil {
		log.Error.Printf("failure looking up parent cache groups %v", err)
		return tc.DBError, tc.SystemError
	}

	log.Debugf("about to run exec query: %s with cg: %++v", updateQuery(), cg)
	resultRows, err := cg.ReqInfo.Tx.NamedQuery(updateQuery(), cg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			err, eType := dbhelpers.ParsePQUniqueConstraintError(pqErr)
			if eType == tc.DataConflictError {
				return errors.New("a cg with " + err.Error()), eType
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
	cg.LastUpdated = &lastUpdated
	if rowsAffected != 1 {
		if rowsAffected < 1 {
			return errors.New("no cg found with this id"), tc.DataMissingError
		} else {
			return fmt.Errorf("this update affected too many rows: %d", rowsAffected), tc.SystemError
		}
	}
	return nil, tc.NoError
}

//The CacheGroup implementation of the Deleter interface
//all implementations of Deleter should use transactions and return the proper errorType
func (cg *TOCacheGroup) Delete() (error, tc.ApiErrorType) {
	inUse, err := isUsedByChildCache(cg.ReqInfo.Tx, *cg.ID)
	log.Debugf("inUse: %d, err: %v", inUse, err)
	if inUse == false && err != nil {
		return tc.DBError, tc.SystemError
	}
	if inUse == true {
		return err, tc.DataConflictError
	}

	log.Debugf("about to run exec query: %s with cg: %++v", deleteQuery(), cg)
	result, err := cg.ReqInfo.Tx.NamedExec(deleteQuery(), cg)
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
			return errors.New("no cg with that id found"), tc.DataMissingError
		} else {
			return fmt.Errorf("this create affected too many rows: %d", rowsAffected), tc.SystemError
		}
	}

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
