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
	"errors"
	"fmt"
	"strconv"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc/common"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/tovalidate"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// ASNsPrivLevel ...
const ASNsPrivLevel = 10

//we need a type alias to define functions on
type TOASN tc.ASNNullable

//the refType is passed into the handlers where a copy of its type is used to decode the json.
var refType = TOASN(tc.ASNNullable{})

func GetRefType() *TOASN {
	return &refType
}

//Implementation of the Identifier, Validator interface functions
func (asn TOASN) GetID() (int, bool) {
	if asn.ID == nil {
		return 0, false
	}
	return *asn.ID, true
}

func (asn TOASN) GetAuditName() string {
	if asn.ASN == nil {
		id, _ := asn.GetID()
		return strconv.Itoa(id)
	}
	return strconv.Itoa(*asn.ASN)
}

func (asn TOASN) GetType() string {
	return "asn"
}

func (asn *TOASN) SetID(i int) {
	asn.ID = &i
}

func (asn TOASN) Validate(db *sqlx.DB) []error {
	errs := validation.Errors{
		"asn":          validation.Validate(asn.ASN, validation.NotNil, validation.Min(0)),
		"cachegroupId": validation.Validate(asn.CachegroupID, validation.NotNil, validation.Min(0)),
	}
	return tovalidate.ToErrors(errs)
}

//The TOASN implementation of the Creator interface
//all implementations of Creator should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a asn with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
//The insert sql returns the id and lastUpdated values of the newly inserted asn and have
//to be added to the struct
func (asn *TOASN) Create(db *sqlx.DB, user auth.CurrentUser) (error, common.ApiErrorType) {
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
		return common.DBError, common.SystemError
	}
	resultRows, err := tx.NamedQuery(insertQuery(), asn)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			err, eType := dbhelpers.ParsePQUniqueConstraintError(pqErr)
			if eType == common.DataConflictError {
				return errors.New("a asn with " + err.Error()), eType
			}
			return err, eType
		} else {
			log.Errorf("received non pq error: %++v from create execution", err)
			return common.DBError, common.SystemError
		}
	}
	defer resultRows.Close()

	var id int
	var lastUpdated common.TimeNoMod
	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err := resultRows.Scan(&id, &lastUpdated); err != nil {
			log.Error.Printf("could not scan id from insert: %s\n", err)
			return common.DBError, common.SystemError
		}
	}
	if rowsAffected == 0 {
		err = errors.New("no asn was inserted, no id was returned")
		log.Errorln(err)
		return common.DBError, common.SystemError
	} else if rowsAffected > 1 {
		err = errors.New("too many ids returned from asn insert")
		log.Errorln(err)
		return common.DBError, common.SystemError
	}
	asn.SetID(id)
	asn.LastUpdated = &lastUpdated
	err = tx.Commit()
	if err != nil {
		log.Errorln("Could not commit transaction: ", err)
		return common.DBError, common.SystemError
	}
	rollbackTransaction = false
	return nil, common.NoError
}

func (asn *TOASN) Read(db *sqlx.DB, parameters map[string]string, user auth.CurrentUser) ([]interface{}, []error, common.ApiErrorType) {
	var rows *sqlx.Rows

	// Query Parameters to Database Query column mappings
	// see the fields mapped in the SQL query
	queryParamsToQueryCols := map[string]dbhelpers.WhereColumnInfo{
		"asn":            dbhelpers.WhereColumnInfo{"a.asn", nil},
		"cachegroup":     dbhelpers.WhereColumnInfo{"c.id", nil},
		"id":             dbhelpers.WhereColumnInfo{"a.id", api.IsInt},
		"cachegroupName": dbhelpers.WhereColumnInfo{"c.name", nil},
	}
	where, orderBy, queryValues, errs := dbhelpers.BuildWhereAndOrderBy(parameters, queryParamsToQueryCols)
	if len(errs) > 0 {
		return nil, errs, common.DataConflictError
	}

	query := selectQuery() + where + orderBy
	log.Debugln("Query is ", query)

	rows, err := db.NamedQuery(query, queryValues)
	if err != nil {
		log.Errorf("Error querying ASNs: %v", err)
		return nil, []error{err}, common.SystemError
	}
	defer rows.Close()

	ASNs := []interface{}{}
	for rows.Next() {
		var s TOASN
		if err = rows.StructScan(&s); err != nil {
			log.Errorf("error parsing ASN rows: %v", err)
			return nil, []error{err}, common.SystemError
		}
		ASNs = append(ASNs, s)
	}

	return ASNs, []error{}, common.NoError
}

func selectQuery() string {
	query := `SELECT
a.id,
a.asn,
a.last_updated,
a.cachegroup AS cachegroup_id,
c.name AS cachegroup

FROM asn a JOIN cachegroup c ON a.cachegroup = c.id`
	return query
}

//The TOASN implementation of the Updater interface
//all implementations of Updater should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a asn with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
func (asn *TOASN) Update(db *sqlx.DB, user auth.CurrentUser) (error, common.ApiErrorType) {
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
		return common.DBError, common.SystemError
	}
	log.Debugf("about to run exec query: %s with asn: %++v", updateQuery(), asn)
	resultRows, err := tx.NamedQuery(updateQuery(), asn)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			err, eType := dbhelpers.ParsePQUniqueConstraintError(pqErr)
			if eType == common.DataConflictError {
				return errors.New("a asn with " + err.Error()), eType
			}
			return err, eType
		} else {
			log.Errorf("received error: %++v from update execution", err)
			return common.DBError, common.SystemError
		}
	}
	defer resultRows.Close()

	var lastUpdated common.TimeNoMod
	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err := resultRows.Scan(&lastUpdated); err != nil {
			log.Error.Printf("could not scan lastUpdated from insert: %s\n", err)
			return common.DBError, common.SystemError
		}
	}
	log.Debugf("lastUpdated: %++v", lastUpdated)
	asn.LastUpdated = &lastUpdated
	if rowsAffected != 1 {
		if rowsAffected < 1 {
			return errors.New("no asn found with this id"), common.DataMissingError
		} else {
			return fmt.Errorf("this update affected too many rows: %d", rowsAffected), common.SystemError
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Errorln("Could not commit transaction: ", err)
		return common.DBError, common.SystemError
	}
	rollbackTransaction = false
	return nil, common.NoError
}

//The ASN implementation of the Deleter interface
//all implementations of Deleter should use transactions and return the proper errorType
func (asn *TOASN) Delete(db *sqlx.DB, user auth.CurrentUser) (error, common.ApiErrorType) {
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
		return common.DBError, common.SystemError
	}
	log.Debugf("about to run exec query: %s with asn: %++v", deleteQuery(), asn)
	result, err := tx.NamedExec(deleteQuery(), asn)
	if err != nil {
		log.Errorf("received error: %++v from delete execution", err)
		return common.DBError, common.SystemError
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return common.DBError, common.SystemError
	}
	if rowsAffected != 1 {
		if rowsAffected < 1 {
			return errors.New("no asn with that id found"), common.DataMissingError
		} else {
			return fmt.Errorf("this create affected too many rows: %d", rowsAffected), common.SystemError
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Errorln("Could not commit transaction: ", err)
		return common.DBError, common.SystemError
	}
	rollbackTransaction = false
	return nil, common.NoError
}

func insertQuery() string {
	query := `INSERT INTO asn (
asn,
cachegroup) 
VALUES (
:asn,
:cachegroup_id
)
RETURNING id,last_updated`
	return query
}

func updateQuery() string {
	query := `UPDATE
asn SET
asn=:asn,
cachegroup=:cachegroup_id
WHERE id=:id RETURNING last_updated`
	return query
}

func deleteQuery() string {
	query := `DELETE FROM asn
WHERE id=:id`
	return query
}
