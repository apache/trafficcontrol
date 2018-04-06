package profile

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

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc/v13"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/parameter"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/tovalidate"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

const (
	CDNQueryParam         = "cdn"
	NameQueryParam        = "name"
	IDQueryParam          = "id"
	DescriptionQueryParam = "description"
	TypeQueryParam        = "type"
)

//we need a type alias to define functions on
type TOProfile v13.ProfileNullable
type TOParameter v13.ParameterNullable

//the refType is passed into the handlers where a copy of its type is used to decode the json.
var refType = TOProfile{}

func GetRefType() *TOProfile {
	return &refType
}

func (prof TOProfile) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{IDQueryParam, api.GetIntKey}}
}

//Implementation of the Identifier, Validator interface functions
func (prof TOProfile) GetKeys() (map[string]interface{}, bool) {
	if prof.ID == nil {
		return map[string]interface{}{IDQueryParam: 0}, false
	}
	return map[string]interface{}{IDQueryParam: *prof.ID}, true
}

func (prof *TOProfile) SetKeys(keys map[string]interface{}) {
	i, _ := keys[IDQueryParam].(int) //this utilizes the non panicking type assertion, if the thrown away ok variable is false i will be the zero of the type, 0 here.
	prof.ID = &i
}

func (prof *TOProfile) GetAuditName() string {
	if prof.Name != nil {
		return *prof.Name
	}
	if prof.ID != nil {
		return strconv.Itoa(*prof.ID)
	}
	return "unknown"
}

func (prof *TOProfile) GetType() string {
	return "profile"
}

func (prof *TOProfile) Validate(db *sqlx.DB) []error {
	errs := validation.Errors{
		NameQueryParam:        validation.Validate(prof.Name, validation.Required),
		DescriptionQueryParam: validation.Validate(prof.Description, validation.Required),
		CDNQueryParam:         validation.Validate(prof.CDNID, validation.Required),
		TypeQueryParam:        validation.Validate(prof.Type, validation.Required),
	}
	if errs != nil {
		return tovalidate.ToErrors(errs)
	}
	return nil
}

func (prof *TOProfile) Read(db *sqlx.DB, parameters map[string]string, user auth.CurrentUser) ([]interface{}, []error, tc.ApiErrorType) {
	var rows *sqlx.Rows

	// Query Parameters to Database Query column mappings
	// see the fields mapped in the SQL query
	queryParamsToQueryCols := map[string]dbhelpers.WhereColumnInfo{
		NameQueryParam: dbhelpers.WhereColumnInfo{"prof.name", nil},
		IDQueryParam:   dbhelpers.WhereColumnInfo{"prof.id", api.IsInt},
	}
	where, orderBy, queryValues, errs := dbhelpers.BuildWhereAndOrderBy(parameters, queryParamsToQueryCols)
	if len(errs) > 0 {
		return nil, errs, tc.DataConflictError
	}

	query := selectProfilesQuery() + where + orderBy
	log.Debugln("Query is ", query)

	rows, err := db.NamedQuery(query, queryValues)
	if err != nil {
		log.Errorf("Error querying Profile: %v", err)
		return nil, []error{tc.DBError}, tc.SystemError
	}
	defer rows.Close()

	profiles := []interface{}{}
	for rows.Next() {
		var p v13.ProfileNullable
		if err = rows.StructScan(&p); err != nil {
			log.Errorf("error parsing Profile rows: %v", err)
			return nil, []error{tc.DBError}, tc.SystemError
		}

		// Attach Parameters if the 'id' parameter is sent
		if _, ok := parameters[IDQueryParam]; ok {
			params, err := ReadParameters(db, parameters, user, p)
			p.Parameters = params
			if len(errs) > 0 {
				log.Errorf("Error getting Parameters: %v", err)
				return nil, []error{tc.DBError}, tc.SystemError
			}
		}
		profiles = append(profiles, p)
	}

	return profiles, []error{}, tc.NoError

}

func selectProfilesQuery() string {

	query := `SELECT
prof.description,
prof.id,
prof.last_updated,
prof.name,
prof.routing_disabled,
prof.type,
c.id as cdn,
c.name as cdn_name
FROM profile prof
JOIN cdn c ON prof.cdn = c.id`

	return query
}

func ReadParameters(db *sqlx.DB, parameters map[string]string, user auth.CurrentUser, profile v13.ProfileNullable) ([]v13.ParameterNullable, []error) {

	var rows *sqlx.Rows
	privLevel := user.PrivLevel
	queryValues := make(map[string]interface{})
	queryValues["profile_id"] = *profile.ID

	query := selectParametersQuery()
	rows, err := db.NamedQuery(query, queryValues)
	if err != nil {
		log.Errorf("Error querying Parameter: %v", err)
		return nil, []error{tc.DBError}
	}
	defer rows.Close()

	var params []v13.ParameterNullable
	for rows.Next() {
		var param v13.ParameterNullable

		if err = rows.StructScan(&param); err != nil {
			log.Errorf("error parsing parameter rows: %v", err)
			return nil, []error{tc.DBError}
		}
		var isSecure bool
		if param.Secure != nil {
			isSecure = *param.Secure
		}
		if isSecure && (privLevel < auth.PrivLevelAdmin) {
			param.Value = &parameter.HiddenField
		}
		params = append(params, param)
	}
	return params, []error{}
}

func selectParametersQuery() string {

	query := `SELECT
p.id,
p.name,
p.config_file,
p.value,
p.secure
FROM parameter p
JOIN profile_parameter pp ON pp.parameter = p.id 
WHERE pp.profile = :profile_id`

	return query
}

//The TOProfile implementation of the Updater interface
//all implementations of Updater should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a profile with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
func (prof *TOProfile) Update(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
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
	log.Debugf("about to run exec query: %s with profile: %++v", updateQuery(), prof)
	resultRows, err := tx.NamedQuery(updateQuery(), prof)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			err, eType := dbhelpers.ParsePQUniqueConstraintError(pqErr)
			if eType == tc.DataConflictError {
				return errors.New("a profile with " + err.Error()), eType
			}
			return err, eType
		}
		log.Errorf("received error: %++v from update execution", err)
		return tc.DBError, tc.SystemError
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
	prof.LastUpdated = &lastUpdated
	if rowsAffected != 1 {
		if rowsAffected < 1 {
			return errors.New("no profile found with this id"), tc.DataMissingError
		}
		return fmt.Errorf("this update affected too many rows: %d", rowsAffected), tc.SystemError
	}
	err = tx.Commit()
	if err != nil {
		log.Errorln("Could not commit transaction: ", err)
		return tc.DBError, tc.SystemError
	}
	rollbackTransaction = false
	return nil, tc.NoError
}

//The TOProfile implementation of the Creator interface
//all implementations of Creator should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a profile with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
//The insert sql returns the id and lastUpdated values of the newly inserted profile and have
//to be added to the struct
func (prof *TOProfile) Create(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
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
	resultRows, err := tx.NamedQuery(insertQuery(), prof)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			err, eType := dbhelpers.ParsePQUniqueConstraintError(pqErr)
			if eType == tc.DataConflictError {
				return errors.New("a profile with " + err.Error()), eType
			}
			return err, eType
		}
		log.Errorf("received non pq error: %++v from create execution", err)
		return tc.DBError, tc.SystemError
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
		err = errors.New("no profile was inserted, no id was returned")
		log.Errorln(err)
		return tc.DBError, tc.SystemError
	}
	if rowsAffected > 1 {
		err = errors.New("too many ids returned from profile insert")
		log.Errorln(err)
		return tc.DBError, tc.SystemError
	}

	prof.SetKeys(map[string]interface{}{IDQueryParam: id})
	prof.LastUpdated = &lastUpdated
	err = tx.Commit()
	if err != nil {
		log.Errorln("Could not commit transaction: ", err)
		return tc.DBError, tc.SystemError
	}
	rollbackTransaction = false
	return nil, tc.NoError
}

//The Profile implementation of the Deleter interface
//all implementations of Deleter should use transactions and return the proper errorType
func (prof *TOProfile) Delete(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
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
	log.Debugf("about to run exec query: %s with profile: %++v", deleteQuery(), prof)
	result, err := tx.NamedExec(deleteQuery(), prof)
	if err != nil {
		log.Errorf("received error: %++v from delete execution", err)
		return tc.DBError, tc.SystemError
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return tc.DBError, tc.SystemError
	}
	if rowsAffected < 1 {
		return errors.New("no profile with that id found"), tc.DataMissingError
	}
	if rowsAffected > 1 {
		return fmt.Errorf("this create affected too many rows: %d", rowsAffected), tc.SystemError
	}

	err = tx.Commit()
	if err != nil {
		log.Errorln("Could not commit transaction: ", err)
		return tc.DBError, tc.SystemError
	}
	rollbackTransaction = false
	return nil, tc.NoError
}

func updateQuery() string {
	query := `UPDATE
profile SET
cdn=:cdn,
description=:description,
name=:name,
routing_disabled=:routing_disabled,
type=:type
WHERE id=:id RETURNING last_updated`
	return query
}

func insertQuery() string {
	query := `INSERT INTO profile (
cdn,
description,
name,
routing_disabled,
type) VALUES (
:cdn,
:description,
:name,
:routing_disabled,
:type) RETURNING id,last_updated`
	return query
}

func deleteQuery() string {
	query := `DELETE FROM profile
WHERE id=:id`
	return query
}
