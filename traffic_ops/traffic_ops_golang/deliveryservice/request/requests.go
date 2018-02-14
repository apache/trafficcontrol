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
	"errors"
	"fmt"
	"strconv"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc"

	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/tenant"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// TODeliveryServiceRequest provides a type alias to define functions on
type TODeliveryServiceRequest tc.DeliveryServiceRequestNullable

//the refType is passed into the handlers where a copy of its type is used to decode the json.
var refType = TODeliveryServiceRequest(tc.DeliveryServiceRequestNullable{})

// GetRefType is used to decode the JSON for deliveryservice requests
func GetRefType() *TODeliveryServiceRequest {
	return &refType
}

//Implementation of the Identifier, Validator interface functions

// GetID is part of the tc.Identifier interface
func (req *TODeliveryServiceRequest) GetID() int {
	return req.ID
}

// GetAuditName is part of the tc.Identifier interface
func (req *TODeliveryServiceRequest) GetAuditName() string {
	return strconv.Itoa(req.ID)
}

// GetType is part of the tc.Identifier interface
func (req *TODeliveryServiceRequest) GetType() string {
	return "deliveryservice_request"
}

// SetID is part of the tc.Identifier interface
func (req *TODeliveryServiceRequest) SetID(i int) {
	req.ID = i
}

// Read implements the api.Reader interface
func (req *TODeliveryServiceRequest) Read(db *sqlx.DB, parameters map[string]string, user auth.CurrentUser) ([]interface{}, []error, tc.ApiErrorType) {
	queryParamsToQueryCols := map[string]dbhelpers.WhereColumnInfo{
		"assignee":   dbhelpers.WhereColumnInfo{Column: "s.username"},
		"assigneeId": dbhelpers.WhereColumnInfo{Column: "r.assignee_id", Checker: api.IsInt},
		"author":     dbhelpers.WhereColumnInfo{Column: "a.username"},
		"authorId":   dbhelpers.WhereColumnInfo{Column: "r.author_id", Checker: api.IsInt},
		"changeType": dbhelpers.WhereColumnInfo{Column: "r.change_type"},
		"id":         dbhelpers.WhereColumnInfo{Column: "r.id", Checker: api.IsInt},
		"status":     dbhelpers.WhereColumnInfo{Column: "r.status"},
		"xmlId":      dbhelpers.WhereColumnInfo{Column: "r.deliveryservice->>'xmlId'"},
	}

	p := parameters
	if _, ok := parameters["orderby"]; !ok {
		// if orderby not provided, default to orderby xmlId.  Making a copy of parameters to not modify input arg
		p = make(map[string]string, len(parameters))
		for k, v := range parameters {
			p[k] = v
		}
		p["orderby"] = "xmlId"
	}

	where, orderBy, queryValues, errs := dbhelpers.BuildWhereAndOrderBy(p, queryParamsToQueryCols)
	if len(errs) > 0 {
		return nil, errs, tc.DataConflictError
	}

	query := selectDeliveryServiceRequestsQuery() + where + orderBy
	log.Debugln("Query is ", query)

	rows, err := db.NamedQuery(query, queryValues)
	if err != nil {
		log.Errorf("Error querying DeliveryServiceRequests: %v", err)
		return nil, []error{tc.DBError}, tc.SystemError
	}
	defer rows.Close()

	var deliveryServiceRequests []interface{}
	for rows.Next() {
		var s TODeliveryServiceRequest
		if err = rows.StructScan(&s); err != nil {
			log.Errorf("error parsing DeliveryServiceRequest rows: %v", err)
			return nil, []error{tc.DBError}, tc.SystemError
		}

		// TODO: combine tenancy with the query above so there's a single db call
		t, err := s.IsTenantAuthorized(user, db)
		if err != nil {
			log.Errorf("error checking tenancy: %v", err)
			return nil, []error{tc.DBError}, tc.SystemError
		}
		if t {
			deliveryServiceRequests = append(deliveryServiceRequests, s)
		}
	}

	return deliveryServiceRequests, []error{}, tc.NoError
}

func selectDeliveryServiceRequestsQuery() string {

	query := `SELECT
a.username AS author,
r.assignee_id,
r.author_id,
r.change_type,
r.created_at,
r.id,
r.last_edited_by_id,
r.last_updated,
r.deliveryservice,
r.status,
s.username AS assignee,
r.deliveryservice->>'xmlId' as xml_id

FROM deliveryservice_request r
JOIN tm_user a ON r.author_id = a.id
LEFT OUTER JOIN tm_user s ON r.assignee_id = s.id
`
	return query
}

// IsTenantAuthorized implements the Tenantable interface to ensure the user is authorized on the deliveryservice tenant
func (req *TODeliveryServiceRequest) IsTenantAuthorized(user auth.CurrentUser, db *sqlx.DB) (bool, error) {
	ds := req.DeliveryService
	if ds == nil {
		// No deliveryservice applied yet -- wide open
		return true, nil
	}
	if ds.TenantID == 0 {
		log.Debugf("tenantID is 0 -- THIS SHOULD NEVER HAPPEN!!")
		return false, errors.New("tenantID is 0")
	}
	return tenant.IsResourceAuthorizedToUser(ds.TenantID, user, db)
}

// Update implements the tc.Updater interface.
//all implementations of Updater should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a request with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
func (req *TODeliveryServiceRequest) Update(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
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
		log.Error.Println("could not begin transaction: ", err.Error())
		return err, tc.SystemError
	}

	// update lasteditedby field prior to making the query -- updated by current user
	req.LastEditedByID = tc.IDNoMod(user.ID)
	resultRows, err := tx.NamedQuery(updateRequestQuery(), req)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			err, eType := dbhelpers.ParsePQUniqueConstraintError(pqErr)
			if eType == tc.DataConflictError {
				return errors.New("a deliveryservice request with " + err.Error()), eType
			}
			return err, eType
		}
		log.Errorf("received error: %++v from update execution", err)
		return tc.DBError, tc.SystemError
	}
	defer resultRows.Close()

	// get LastUpdated field -- updated by trigger in the db
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
	req.LastUpdated = &lastUpdated

	if rowsAffected < 1 {
		return errors.New("no deliveryservice request found with this id"), tc.DataMissingError
	} else if rowsAffected > 1 {
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

// Insert implements the tc.Inserter interface
//all implementations of Inserter should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a request with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
//The insert sql returns the id and lastUpdated values of the newly inserted request and have
//to be added to the struct
func (req *TODeliveryServiceRequest) Insert(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
	if req == nil {
		return errors.New("nil deliveryservice_request"), tc.SystemError
	}
	if req.Status != "draft" && req.Status != "submitted" {
		return errors.New("invalid initial request status " + string(req.Status) + ".  Must be 'draft' or 'submitted'"), tc.DataConflictError
	}
	// first, ensure there's not an active request with this XMLID
	ds := req.DeliveryService
	if ds == nil {
		log.Debugln(" -- no ds")
		return errors.New("no delivery service associated with this request"), tc.DataMissingError
	}
	if ds.XMLID == nil {
		log.Debugln(" -- no XMLID")
		return errors.New("no xmlId associated with this request"), tc.DataMissingError
	}
	XMLID := *ds.XMLID
	active, err := isActiveRequest(db, XMLID)
	if err != nil {
		return err, tc.SystemError
	}
	if active {
		return errors.New("An active request exists for delivery service '" + XMLID), tc.DataConflictError
	}

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

	req.AuthorID = tc.IDNoMod(user.ID)
	req.LastEditedByID = tc.IDNoMod(user.ID)
	resultRows, err := tx.NamedQuery(insertRequestQuery(), req)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			err, eType := dbhelpers.ParsePQUniqueConstraintError(pqErr)
			if eType == tc.DataConflictError {
				return errors.New("a deliveryservice request with " + err.Error()), eType
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
		err = errors.New("no deliveryservice request inserted, no id was returned")
		log.Errorln(err)
		return tc.DBError, tc.SystemError
	} else if rowsAffected > 1 {
		err = errors.New("too many ids returned from deliveryservice request insert")
		log.Errorln(err)
		return tc.DBError, tc.SystemError
	}
	req.SetID(id)
	req.LastUpdated = &lastUpdated
	err = tx.Commit()
	if err != nil {
		log.Errorln("Could not commit transaction: ", err)
		return tc.DBError, tc.SystemError
	}
	rollbackTransaction = false
	return nil, tc.NoError
}

//The Request implementation of the Deleter interface
//all implementations of Deleter should use transactions and return the proper errorType

// Delete removes the request from the db
func (req *TODeliveryServiceRequest) Delete(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
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
		log.Error.Println("could not begin transaction: ", err.Error())
		return tc.DBError, tc.SystemError
	}
	log.Debugf("about to run exec query: %s with ds request: %++v", deleteRequestQuery(), req)
	result, err := tx.NamedExec(deleteRequestQuery(), req)
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
	err = tx.Commit()
	if err != nil {
		log.Errorln("Could not commit transaction: ", err)
		return tc.DBError, tc.SystemError
	}
	// success!
	rollbackTransaction = false
	log.Infoln("Delete succeeded")
	return nil, tc.NoError
}

// isActiveRequest returns true if a request using this XMLID is currently in an active state
func isActiveRequest(db *sqlx.DB, XMLID string) (bool, error) {
	q := `SELECT EXISTS(SELECT 1 FROM deliveryservice_request
WHERE deliveryservice->>'xml_id' = '` + XMLID + `'
AND status IN ('draft', 'submitted', 'pending'))`
	row := db.QueryRow(q)
	var active bool
	err := row.Scan(&active)
	if err != nil {
		log.Debugln("ERROR: ", err, ";  QUERY:", q)
	}
	return active, err
}

func updateRequestQuery() string {
	query := `UPDATE
deliveryservice_request
SET change_type=:change_type,
last_edited_by_id=:last_edited_by_id,
deliveryservice=:deliveryservice,
WHERE id=:id RETURNING last_updated`
	return query
}

func insertRequestQuery() string {
	query := `INSERT INTO deliveryservice_request (
assignee_id,
author_id,
change_type,
last_edited_by_id,
deliveryservice,
status
) VALUES (
:assignee_id,
:author_id,
:change_type,
:last_edited_by_id,
:deliveryservice,
:status
) RETURNING id,last_updated`
	return query
}

func deleteRequestQuery() string {
	query := `DELETE FROM deliveryservice_request
WHERE id=:id`
	return query
}

////////////////////////////////////////////////////////////////
// Assignment change
type deliveryServiceRequestAssignment TODeliveryServiceRequest

// GetAssignRefType is used to decode the JSON for deliveryservice_request assignment
func GetAssignRefType() *deliveryServiceRequestAssignment {
	return &deliveryServiceRequestAssignment{}
}

// Update assignee only
func (req *deliveryServiceRequestAssignment) Update(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
	// req represents the state the deliveryservice_request is to transition to
	// we want to limit what changes here -- only assignee can change

	// get original
	var current deliveryServiceRequestAssignment
	err := db.QueryRowx(selectDeliveryServiceRequestsQuery() + " WHERE r.id=" + strconv.Itoa(req.ID)).StructScan(&current)
	if err != nil {
		log.Errorf("Error querying DeliveryServiceRequests: %v", err)
		return err, tc.SystemError
	}

	if current.AssigneeID == req.AssigneeID {
		log.Infof("assignee unchanged")
		return nil, tc.NoError
	}

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
		log.Error.Println("could not begin transaction: ", err.Error())
		return err, tc.SystemError
	}

	// LastEditedBy field should not change with status update
	v := "null"
	if req.AssigneeID != nil {
		v = strconv.Itoa(*req.AssigneeID)
	}
	query := fmt.Sprintf(`UPDATE deliveryservice_request SET assignee_id = %s WHERE id=%d`, v, current.ID)
	_, err = tx.Exec(query)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			err, eType := dbhelpers.ParsePQUniqueConstraintError(pqErr)
			if eType == tc.DataConflictError {
				return errors.New("a deliveryservice request with " + err.Error()), eType
			}
			return err, eType
		}
		log.Errorf("received error: %++v from update execution", err)
		return tc.DBError, tc.SystemError
	}

	err = tx.Commit()
	if err != nil {
		log.Errorln("Could not commit transaction: ", err)
		return tc.DBError, tc.SystemError
	}
	rollbackTransaction = false
	return nil, tc.NoError
}

func (req *deliveryServiceRequestAssignment) GetID() int {
	return (*TODeliveryServiceRequest)(req).GetID()
}

func (req *deliveryServiceRequestAssignment) GetType() string {
	return (*TODeliveryServiceRequest)(req).GetType()
}

func (req *deliveryServiceRequestAssignment) GetAuditName() string {
	return (*TODeliveryServiceRequest)(req).GetAuditName()
}

func (req *deliveryServiceRequestAssignment) Validate(db *sqlx.DB) []error {
	return nil
}

////////////////////////////////////////////////////////////////
// Status change

// deliveryServiceRequestStatus implements interfaces needed to update the request status only
type deliveryServiceRequestStatus TODeliveryServiceRequest

// GetStatusRefType is used to decode the JSON for deliveryservice_request status change
func GetStatusRefType() *deliveryServiceRequestStatus {
	return &deliveryServiceRequestStatus{}
}

// Update status only
func (req *deliveryServiceRequestStatus) Update(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
	// req represents the state the deliveryservice_request is to transition to
	// we want to limit what changes here -- only status can change,  and only according to the established rules
	// for status transition

	// get original
	var current deliveryServiceRequestStatus
	err := db.QueryRowx(selectDeliveryServiceRequestsQuery() + " WHERE r.id=" + strconv.Itoa(req.ID)).StructScan(&current)
	if err != nil {
		log.Errorf("Error querying DeliveryServiceRequests: %v", err)
		return err, tc.SystemError
	}

	st := req.Status
	if err = current.Status.ValidTransition(st); err != nil {
		return err, tc.DataConflictError
	}

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
		log.Error.Println("could not begin transaction: ", err.Error())
		return err, tc.SystemError
	}

	// LastEditedBy field should not change with status update
	query := fmt.Sprintf(`UPDATE deliveryservice_request SET status = '%s' WHERE id=%d`, req.Status, current.ID)
	_, err = tx.Exec(query)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			err, eType := dbhelpers.ParsePQUniqueConstraintError(pqErr)
			if eType == tc.DataConflictError {
				return errors.New("a deliveryservice request with " + err.Error()), eType
			}
			return err, eType
		}
		log.Errorf("received error: %++v from update execution", err)
		return tc.DBError, tc.SystemError
	}

	err = tx.Commit()
	if err != nil {
		log.Errorln("Could not commit transaction: ", err)
		return tc.DBError, tc.SystemError
	}
	rollbackTransaction = false
	return nil, tc.NoError
}

// GetID from base type
func (req *deliveryServiceRequestStatus) GetID() int {
	return (*TODeliveryServiceRequest)(req).GetID()
}

// GetType from base type
func (req *deliveryServiceRequestStatus) GetType() string {
	return (*TODeliveryServiceRequest)(req).GetType()
}

// GetAuditName from base type
func (req *deliveryServiceRequestStatus) GetAuditName() string {
	return (*TODeliveryServiceRequest)(req).GetAuditName()
}

// Validate is not needed when only Status is updated
func (req *deliveryServiceRequestStatus) Validate(db *sqlx.DB) []error {
	return nil
}
