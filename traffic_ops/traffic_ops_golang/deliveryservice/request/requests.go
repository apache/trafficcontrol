package request

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
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc"

	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// TODeliveryServiceRequest provides a type alias to define functions on
type TODeliveryServiceRequest tc.DeliveryServiceRequest

//the refType is passed into the handlers where a copy of its type is used to decode the json.
var refType = TODeliveryServiceRequest(tc.DeliveryServiceRequest{})

// GetRefType is used to decode the JSON for deliveryservice requests
func GetRefType() *TODeliveryServiceRequest {
	return &refType
}

//Implementation of the Identifier, Validator interface functions

// GetID ...
func (request *TODeliveryServiceRequest) GetID() int {
	return request.ID
}

// GetAuditName ...
func (request *TODeliveryServiceRequest) GetAuditName() string {
	return strconv.Itoa(request.ID)
}

// GetType ...
func (request *TODeliveryServiceRequest) GetType() string {
	return "deliveryservice_request"
}

// SetID ...
func (request *TODeliveryServiceRequest) SetID(i int) {
	request.ID = i
}

// uniqQuery returns the query needed for a unique new request
func uniqXMLIDQuery(ds tc.DeliveryService) string {
	// no two active ds requests can exist for the same xmlid
	q := `SELECT * FROM deliveryservice_request r 
WHERE r.request->>'xml_id' = '` + ds.XMLID + `'
AND r.status IN ('draft', 'submitted', 'pending')`
	return q
}

func validRequest(db *sqlx.DB, request *TODeliveryServiceRequest) error {
	// get the xmlid from the ds
	var ds tc.DeliveryService
	err := json.Unmarshal([]byte(request.Request), &ds)
	if err != nil {
		return err
	}
	r, err := db.Exec(uniqXMLIDQuery(ds))
	if err != nil {
		return err
	}
	n, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("multiple requests for the same xmlId")
	}
	return nil
}

//The TODeliveryServiceRequest implementation of the Updater interface
//all implementations of Updater should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a request with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned

// Update ...
func (request *TODeliveryServiceRequest) Update(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
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
		log.Error.Println("could not begin transaction: ", err.Error())
		return tc.DBError, tc.SystemError
	}
	log.Debugf("about to run exec query: %s with request: %++v", updateRequestQuery(), request)
	resultRows, err := tx.NamedQuery(updateRequestQuery(), request)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			err, eType := dbhelpers.ParsePQUniqueConstraintError(err)
			if eType == tc.DataConflictError {
				return errors.New("a request with " + err.Error()), eType
			}
			return err, eType
		}
		log.Errorf("received error from update execution: %s", err.Error())
		return tc.DBError, tc.SystemError
	}
	var lastUpdated tc.Time
	var rowsAffected int
	for resultRows.Next() {
		rowsAffected++
		if err := resultRows.Scan(&lastUpdated); err != nil {
			log.Error.Println("could not scan lastUpdated from insert: ", err.Error())
			return tc.DBError, tc.SystemError
		}
	}
	log.Debugln("lastUpdated: ", lastUpdated)
	request.LastUpdated = lastUpdated
	if rowsAffected < 1 {
		return errors.New("no request found with this id"), tc.DataMissingError
	}
	if rowsAffected > 1 {
		return fmt.Errorf("this update affected too many rows: %d", rowsAffected), tc.SystemError
	}

	err = validRequest(db, request)
	if err != nil {
		return err, tc.SystemError
	}

	return nil, tc.NoError
}

//The TODeliveryServiceRequest implementation of the Inserter interface
//all implementations of Inserter should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a request with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
//The insert sql returns the id and lastUpdated values of the newly inserted request and have
//to be added to the struct

// Insert ...
func (request *TODeliveryServiceRequest) Insert(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
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
		log.Error.Println("could not begin transaction: ", err.Error())
		return tc.DBError, tc.SystemError
	}
	request.AuthorID = user.ID
	ir := insertRequestQuery()
	resultRows, err := tx.NamedQuery(ir, request)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			err, eType := dbhelpers.ParsePQUniqueConstraintError(err)
			return err, eType
		}
		log.Errorln("received non pq error from create execution: ", err.Error())
		return tc.DBError, tc.SystemError
	}
	var id int
	var lastUpdated tc.Time
	var rowsAffected int
	for resultRows.Next() {
		rowsAffected++
		if err := resultRows.Scan(&id, &lastUpdated); err != nil {
			log.Error.Println("could not scan id from insert: ", err.Error())
			return tc.DBError, tc.SystemError
		}
	}
	if rowsAffected == 0 {
		err = errors.New("no request was inserted, no id was returned")
		log.Errorln(err)
		return tc.DBError, tc.SystemError
	} else if rowsAffected > 1 {
		err = errors.New("too many ids returned from request insert")
		log.Errorln(err)
		return tc.DBError, tc.SystemError
	}
	request.SetID(id)
	request.LastUpdated = lastUpdated
	err = validRequest(db, request)
	if err != nil {
		return err, tc.SystemError
	}
	return nil, tc.NoError
}

//The Request implementation of the Deleter interface
//all implementations of Deleter should use transactions and return the proper errorType

// Delete ...
func (request *TODeliveryServiceRequest) Delete(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
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
		log.Error.Println("could not begin transaction: ", err.Error())
		return tc.DBError, tc.SystemError
	}
	log.Debugf("about to run exec query: %s with request: %++v", deleteRequestQuery(), request)
	result, err := tx.NamedExec(deleteRequestQuery(), request)
	if err != nil {
		log.Errorln("received error from delete execution: ", err.Error())
		return tc.DBError, tc.SystemError
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return tc.DBError, tc.SystemError
	}
	if rowsAffected < 1 {
		return errors.New("no request with that id found"), tc.DataMissingError
	}
	if rowsAffected > 1 {
		return fmt.Errorf("this create affected too many rows: %d", rowsAffected), tc.SystemError
	}
	return nil, tc.NoError
}

func updateRequestQuery() string {
	query := `UPDATE
deliveryservice_request SET
assignee_id=:assignee_id,
change_type=:change_type,
request=:request,
status=:status
WHERE id=:id RETURNING last_updated`
	return query
}

func insertRequestQuery() string {
	query := `INSERT INTO deliveryservice_request (
assignee_id,
author_id,
change_type,
request,
status
) VALUES (
:assignee_id,
:author_id,
:change_type,
:request,
:status
) RETURNING id,last_updated`
	return query
}

func deleteRequestQuery() string {
	query := `DELETE FROM deliveryservice_request
WHERE id=:id`
	return query
}
