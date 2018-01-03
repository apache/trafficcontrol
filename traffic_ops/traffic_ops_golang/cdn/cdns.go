package cdn

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

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc"

	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

//we need a type alias to define functions on
type TOCDN tc.CDN

//the refType is passed into the handlers where a copy of its type is used to decode the json.
var refType = TOCDN(tc.CDN{})

func GetRefType() *TOCDN {
	return &refType
}

//Implementation of the Identifier, Validator interface functions
func (cdn *TOCDN) GetID() int {
	return cdn.ID
}

func (cdn *TOCDN) GetAuditName() string {
	return cdn.Name
}

func (cdn *TOCDN) GetType() string {
	return "cdn"
}

func (cdn *TOCDN) SetID(i int) {
	cdn.ID = i
}

func (cdn *TOCDN) Validate(db *sqlx.DB) []error {
	errs := []error{}
	if len(cdn.Name) < 1 {
		errs = append(errs, errors.New(`CDN 'name' is required.`))
	}
	if len(cdn.DomainName) < 1 {
		errs = append(errs, errors.New("Domain Name is required."))
	}
	return errs
}

//The TOCDN implementation of the Updater interface
//all implementations of Updater should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a cdn with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
func (cdn *TOCDN) Update(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
	tx, err := db.Beginx()
	defer func() {
		if tx == nil {
			return
		}
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	if err != nil {
		log.Error.Printf("could not begin transaction: %v", err)
		return tc.DBError, tc.SystemError
	}
	log.Debugf("about to run exec query: %s with cdn: %++v", updateCDNQuery(), cdn)
	resultRows, err := tx.NamedQuery(updateCDNQuery(), cdn)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			err, eType := dbhelpers.ParsePQUniqueConstraintError(err)
			if eType == tc.DataConflictError {
				return errors.New("a cdn with " + err.Error()), eType
			}
			return err, eType
		} else {
			log.Errorf("received error: %++v from update execution", err)
			return tc.DBError, tc.SystemError
		}
	}
	var lastUpdated tc.Time
	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err := resultRows.Scan(&lastUpdated); err != nil {
			log.Error.Printf("could not scan lastUpdated from insert: %s\n", err)
			return tc.DBError, tc.SystemError
		}
	}
	log.Debugf("lastUpdated: %++v", lastUpdated)
	cdn.LastUpdated = lastUpdated
	if rowsAffected != 1 {
		if rowsAffected < 1 {
			return errors.New("no cdn found with this id"), tc.DataMissingError
		} else {
			return fmt.Errorf("this update affected too many rows: %d", rowsAffected), tc.SystemError
		}
	}
	return nil, tc.NoError
}

//The TOCDN implementation of the Inserter interface
//all implementations of Inserter should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a cdn with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
//The insert sql returns the id and lastUpdated values of the newly inserted cdn and have
//to be added to the struct
func (cdn *TOCDN) Insert(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
	tx, err := db.Beginx()
	defer func() {
		if tx == nil {
			return
		}
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	if err != nil {
		log.Error.Printf("could not begin transaction: %v", err)
		return tc.DBError, tc.SystemError
	}
	resultRows, err := tx.NamedQuery(insertCDNQuery(), cdn)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			err, eType := dbhelpers.ParsePQUniqueConstraintError(err)
			return errors.New("a cdn with " + err.Error()), eType
		} else {
			log.Errorf("received non pq error: %++v from create execution", err)
			return tc.DBError, tc.SystemError
		}
	}
	var id int
	var lastUpdated tc.Time
	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err := resultRows.Scan(&id, &lastUpdated); err != nil {
			log.Error.Printf("could not scan id from insert: %s\n", err)
			return tc.DBError, tc.SystemError
		}
	}
	if rowsAffected == 0 {
		err = errors.New("no cdn was inserted, no id was returned")
		log.Errorln(err)
		return tc.DBError, tc.SystemError
	} else if rowsAffected > 1 {
		err = errors.New("too many ids returned from cdn insert")
		log.Errorln(err)
		return tc.DBError, tc.SystemError
	}
	cdn.SetID(id)
	cdn.LastUpdated = lastUpdated
	return nil, tc.NoError
}

//The CDN implementation of the Deleter interface
//all implementations of Deleter should use transactions and return the proper errorType
func (cdn *TOCDN) Delete(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
	tx, err := db.Beginx()
	defer func() {
		if tx == nil {
			return
		}
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	if err != nil {
		log.Error.Printf("could not begin transaction: %v", err)
		return tc.DBError, tc.SystemError
	}
	log.Debugf("about to run exec query: %s with cdn: %++v", deleteCDNQuery(), cdn)
	result, err := tx.NamedExec(deleteCDNQuery(), cdn)
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
			return errors.New("no cdn with that id found"), tc.DataMissingError
		} else {
			return fmt.Errorf("this create affected too many rows: %d", rowsAffected), tc.SystemError
		}
	}
	return nil, tc.NoError
}

func updateCDNQuery() string {
	query := `UPDATE
cdn SET
dnssec_enabled=:dnssec_enabled,
domain_name=:domain_name,
name=:name
WHERE id=:id RETURNING last_updated`
	return query
}

func insertCDNQuery() string {
	query := `INSERT INTO cdn (
dnssec_enabled,
domain_name,
name) VALUES (
:dnssec_enabled,
:domain_name,
:name) RETURNING id,last_updated`
	return query
}

func deleteCDNQuery() string {
	query := `DELETE FROM cdn
WHERE id=:id`
	return query
}
