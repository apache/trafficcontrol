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

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/tenant"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

//we need a type alias to define functions on
type TODeliveryServiceRequest struct{
	ReqInfo *api.APIInfo `json:"-"`
	tc.DeliveryServiceRequestNullable
}

func GetTypeSingleton() func(reqInfo *api.APIInfo)api.CRUDer {
	return func(reqInfo *api.APIInfo)api.CRUDer {
		toReturn := TODeliveryServiceRequest{reqInfo, tc.DeliveryServiceRequestNullable{}}
		return &toReturn
	}
}

func (req TODeliveryServiceRequest) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{"id", api.GetIntKey}}
}

//Implementation of the Identifier, Validator interface functions
func (req TODeliveryServiceRequest) GetKeys() (map[string]interface{}, bool) {
	if req.ID == nil {
		return map[string]interface{}{"id": 0}, false
	}
	return map[string]interface{}{"id": *req.ID}, true
}

func (req *TODeliveryServiceRequest) SetKeys(keys map[string]interface{}) {
	i, _ := keys["id"].(int) //this utilizes the non panicking type assertion, if the thrown away ok variable is false i will be the zero of the type, 0 here.
	req.ID = &i
}

// GetAuditName is part of the tc.Identifier interface
func (req TODeliveryServiceRequest) GetAuditName() string {
	return req.getXMLID()
}

// GetType is part of the tc.Identifier interface
func (req TODeliveryServiceRequest) GetType() string {
	return "deliveryservice_request"
}

// Read implements the api.Reader interface
func (req *TODeliveryServiceRequest) Read(parameters map[string]string) ([]interface{}, []error, tc.ApiErrorType) {
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
	if tenant.IsTenancyEnabledTx(req.ReqInfo.Tx) {
		log.Debugln("Tenancy is enabled")
		tenantIDs, err := tenant.GetUserTenantIDListTx(req.ReqInfo.User, req.ReqInfo.Tx)
		if err != nil {
			log.Errorln("received error querying for user's tenants: " + err.Error())
			return nil, []error{tc.DBError}, tc.SystemError
		}
		where, queryValues = dbhelpers.AddTenancyCheck(where, queryValues, "r.deliveryservice->>'tenantId", tenantIDs)
	}

	query := selectDeliveryServiceRequestsQuery() + where + orderBy
	log.Debugln("Query is ", query)

	rows, err := req.ReqInfo.Tx.NamedQuery(query, queryValues)
	if err != nil {
		log.Errorf("Error querying DeliveryServiceRequests: %v", err)
		return nil, []error{tc.DBError}, tc.SystemError
	}
	defer rows.Close()

	deliveryServiceRequests := []interface{}{}
	for rows.Next() {
		var s TODeliveryServiceRequest
		if err = rows.StructScan(&s); err != nil {
			log.Errorf("error parsing DeliveryServiceRequest rows: %v", err)
			return nil, []error{tc.DBError}, tc.SystemError
		}
		deliveryServiceRequests = append(deliveryServiceRequests, s)
	}

	return deliveryServiceRequests, []error{}, tc.NoError
}

func selectDeliveryServiceRequestsQuery() string {

	query := `SELECT
a.username AS author,
e.username AS lastEditedBy,
s.username AS assignee,
r.assignee_id,
r.author_id,
r.change_type,
r.created_at,
r.id,
r.last_edited_by_id,
r.last_updated,
r.deliveryservice,
r.status,
r.deliveryservice->>'xmlId' as xml_id

FROM deliveryservice_request r
JOIN tm_user a ON r.author_id = a.id
LEFT OUTER JOIN tm_user s ON r.assignee_id = s.id
LEFT OUTER JOIN tm_user e ON r.last_edited_by_id = e.id
`
	return query
}

// IsTenantAuthorized implements the Tenantable interface to ensure the user is authorized on the deliveryservice tenant
func (req TODeliveryServiceRequest) IsTenantAuthorized(user *auth.CurrentUser) (bool, error) {

	ds := req.DeliveryService
	if ds == nil {
		// No deliveryservice applied yet -- wide open
		return true, nil
	}
	if ds.TenantID == nil {
		log.Debugf("tenantID is nil")
		return false, errors.New("tenantID is nil")
	}
	return tenant.IsResourceAuthorizedToUserTx(*ds.TenantID, user, req.ReqInfo.Tx.Tx)
}

// Update implements the tc.Updater interface.
//all implementations of Updater should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a request with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
func (req *TODeliveryServiceRequest) Update() (error, tc.ApiErrorType) {
	var current TODeliveryServiceRequest
	if req.ID == nil {
		log.Errorf("error updating DeliveryServiceRequest: ID is nil")
		return errors.New("error updating DeliveryServiceRequest: ID is nil"), tc.DataMissingError
	}
	err := req.ReqInfo.Tx.QueryRowx(selectDeliveryServiceRequestsQuery() + `WHERE r.id=` + strconv.Itoa(*req.ID)).StructScan(&current)
	if err != nil {
		log.Errorf("Error querying DeliveryServiceRequests: %v", err)
		return err, tc.SystemError
	}

	// Update can only change status between draft & submitted.  All other transitions must go thru
	// the PUT /api/<version>/deliveryservice_request/:id/status endpoint
	if current.Status == nil || req.Status == nil {
		return errors.New("Missing status for DeliveryServiceRequest"), tc.DataMissingError
	}

	if *current.Status != tc.RequestStatusDraft && *current.Status != tc.RequestStatusSubmitted {
		return fmt.Errorf("Cannot change DeliveryServiceRequest in '%s' status.", string(*current.Status)),
			tc.DataConflictError
	}

	if *req.Status != tc.RequestStatusDraft && *req.Status != tc.RequestStatusSubmitted {
		return fmt.Errorf("Cannot change DeliveryServiceRequest status from '%s' to '%s'", string(*current.Status), string(*req.Status)),
			tc.DataConflictError
	}

	userID := tc.IDNoMod(req.ReqInfo.User.ID)
	req.LastEditedByID = &userID
	resultRows, err := req.ReqInfo.Tx.NamedQuery(updateRequestQuery(), req)
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

	return nil, tc.NoError
}

// Creator implements the tc.Creator interface
//all implementations of Creator should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a request with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
//The insert sql returns the id and lastUpdated values of the newly inserted request and have
//to be added to the struct
func (req *TODeliveryServiceRequest) Create() (error, tc.ApiErrorType) {
	if req == nil {
		return errors.New("nil deliveryservice_request"), tc.SystemError
	}
	if req.Status == nil {
		return errors.New("nil deliveryservice_request status"), tc.SystemError
	}
	if *req.Status != tc.RequestStatusDraft && *req.Status != tc.RequestStatusSubmitted {
		return fmt.Errorf("invalid initial request status '%v'.  Must be '%v' or '%v'",
				*req.Status, tc.RequestStatusDraft, tc.RequestStatusSubmitted),
			tc.DataConflictError
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
	active, err := isActiveRequest(req.ReqInfo.Tx, XMLID)
	if err != nil {
		return err, tc.SystemError
	}
	if active {
		return errors.New(`An active request exists for delivery service '` + XMLID + `'`), tc.DataConflictError
	}

	userID := tc.IDNoMod(req.ReqInfo.User.ID)
	req.AuthorID = &userID
	req.LastEditedByID = &userID
	resultRows, err := req.ReqInfo.Tx.NamedQuery(insertRequestQuery(), req)
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
	req.SetKeys(map[string]interface{}{"id": id})
	req.LastUpdated = &lastUpdated

	return nil, tc.NoError
}

//The Request implementation of the Deleter interface
//all implementations of Deleter should use transactions and return the proper errorType

// Delete removes the request from the db
func (req *TODeliveryServiceRequest) Delete() (error, tc.ApiErrorType) {
	var st tc.RequestStatus
	log.Debugln("DELETING REQUEST WITH ID ", strconv.Itoa(*req.ID))
	if req.ID == nil {
		return errors.New("cannot delete deliveryservice_request -- ID is nil"), tc.SystemError
	}

	err := req.ReqInfo.Tx.QueryRow(`SELECT status FROM deliveryservice_request WHERE id=` + strconv.Itoa(*req.ID)).Scan(&st)
	if err != nil {
		return err, tc.SystemError
	}

	if st == tc.RequestStatusComplete || st == tc.RequestStatusPending || st == tc.RequestStatusRejected {
		return fmt.Errorf("cannot delete a deliveryservice_request with state %s", string(st)), tc.DataConflictError
	}

	query := `DELETE FROM deliveryservice_request WHERE id=` + strconv.Itoa(*req.ID)
	log.Debugf("about to run exec query: %s", query)

	result, err := req.ReqInfo.Tx.Exec(query)
	if err != nil {
		log.Errorln("received error from delete execution: ", err.Error())
		return tc.DBError, tc.SystemError
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Errorln("error getting rows affected: ", err.Error())
		return tc.DBError, tc.SystemError
	}
	if rowsAffected < 1 {
		log.Errorln("no request with that id found")
		return errors.New("no request with that id found"), tc.DataMissingError
	}
	if rowsAffected > 1 {
		log.Errorln("the delete affected too many rows")
		return fmt.Errorf("this delete affected too many rows: %d", rowsAffected), tc.SystemError
	}
	// success!

	return nil, tc.NoError
}

func (req TODeliveryServiceRequest) getXMLID() string {
	if req.DeliveryService == nil || req.DeliveryService.XMLID == nil {

		if req.ID != nil {
			return strconv.Itoa(*req.ID)
		}
		return "0"
	}
	return *req.DeliveryService.XMLID
}

// ChangeLogMessage implements the api.ChangeLogger interface for a custom log message
func (req TODeliveryServiceRequest) ChangeLogMessage(action string) (string, error) {
	changeType := "unknown change type"
	if req.ChangeType != nil {
		changeType = *req.ChangeType
	}
	// use ID in case don't have access to XMLID (e.g. on DELETE)
	message := action + ` ` + req.GetType() + ` of type '` + changeType + `' for deliveryservice '` + req.getXMLID() + `'`
	return message, nil
}

// isActiveRequest returns true if a request using this XMLID is currently in an active state
func isActiveRequest(tx *sqlx.Tx, XMLID string) (bool, error) {
	q := `SELECT EXISTS(SELECT 1 FROM deliveryservice_request
WHERE deliveryservice->>'xmlId' = '` + XMLID + `'
AND status IN ('draft', 'submitted', 'pending'))`
	row := tx.QueryRow(q)
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
status=:status
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

func GetAssignmentTypeSingleton() func(reqInfo *api.APIInfo)api.CRUDer {
	return func(reqInfo *api.APIInfo)api.CRUDer {
		toReturn := deliveryServiceRequestAssignment{TODeliveryServiceRequest{reqInfo, tc.DeliveryServiceRequestNullable{}}}
		return &toReturn
	}
}


type deliveryServiceRequestAssignment struct {
	TODeliveryServiceRequest
}

// Update assignee only
func (req *deliveryServiceRequestAssignment) Update() (error, tc.ApiErrorType) {
	// req represents the state the deliveryservice_request is to transition to
	// we want to limit what changes here -- only assignee can change
	if req.ID == nil {
		return errors.New("cannot update DeliveryServiceRequestAssignment -- ID is nil"), tc.SystemError
	}

	// get original
	var current TODeliveryServiceRequest
	err := req.ReqInfo.Tx.QueryRowx(selectDeliveryServiceRequestsQuery() + `WHERE r.id=` + strconv.Itoa(*req.ID)).StructScan(&current)
	if err != nil {
		log.Errorf("Error querying DeliveryServiceRequests: %v", err)
		return err, tc.SystemError
	}

	// unchanged (maybe both nil)
	if current.AssigneeID == req.AssigneeID {
		log.Infof("assignee unchanged")
		return nil, tc.NoError
	}

	// Only assigneeID changes -- nothing else
	assigneeID := req.AssigneeID
	*req = deliveryServiceRequestAssignment{current}
	req.AssigneeID = assigneeID

	// LastEditedBy field should not change with status update
	v := "null"
	if req.AssigneeID != nil {
		v = strconv.Itoa(*req.AssigneeID)
	}
	query := `UPDATE deliveryservice_request SET assignee_id = ` + v + ` WHERE id=` + strconv.Itoa(*req.ID)
	_, err = req.ReqInfo.Tx.Exec(query)
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

	// update req with current info
	err = req.ReqInfo.Tx.QueryRowx(selectDeliveryServiceRequestsQuery() + ` WHERE r.id=` + strconv.Itoa(*req.ID)).StructScan(req)
	if err != nil {
		log.Errorf("Error querying DeliveryServiceRequests: %v", err)
		return err, tc.SystemError
	}

	return nil, tc.NoError
}

func (req deliveryServiceRequestAssignment) Validate() []error {
	return nil
}

// ChangeLogMessage implements the api.ChangeLogger interface for a custom log message
func (req deliveryServiceRequestAssignment) ChangeLogMessage(action string) (string, error) {
	a := "NONE"
	if req.Assignee != nil {
		a = *req.Assignee
	}
	message := `Changed assignee of ‘` + req.getXMLID() + `’ ` + req.GetType() + ` to '` + a + `'`

	return message, nil
}

////////////////////////////////////////////////////////////////
// Status change

// deliveryServiceRequestStatus implements interfaces needed to update the request status only
type deliveryServiceRequestStatus struct {
	TODeliveryServiceRequest
}

func GetStatusTypeSingleton() func(reqInfo *api.APIInfo)api.CRUDer {
	return func(reqInfo *api.APIInfo)api.CRUDer {
		toReturn := deliveryServiceRequestStatus{TODeliveryServiceRequest{reqInfo, tc.DeliveryServiceRequestNullable{}}}
		return &toReturn
	}
}
// Update status only
func (req *deliveryServiceRequestStatus) Update() (error, tc.ApiErrorType) {
	// req represents the state the deliveryservice_request is to transition to
	// we want to limit what changes here -- only status can change,  and only according to the established rules
	// for status transition

	// get original
	var current TODeliveryServiceRequest
	if req.ID == nil {
		log.Errorf("error updating DeliveryServiceRequestStatus: ID is nil")
	}
	err := req.ReqInfo.Tx.QueryRowx(selectDeliveryServiceRequestsQuery() + ` WHERE r.id=` + strconv.Itoa(*req.ID)).StructScan(&current)
	if err != nil {
		log.Errorf("Error querying DeliveryServiceRequests: %v", err)
		return err, tc.SystemError
	}

	if err = current.Status.ValidTransition(*req.Status); err != nil {
		return err, tc.DataConflictError
	}

	// keep everything else the same -- only update status
	st := req.Status
	*req = deliveryServiceRequestStatus{current}
	req.Status = st

	// LastEditedBy field should not change with status update
	query := `UPDATE deliveryservice_request SET status = '` + string(*req.Status) + `' WHERE id=` + strconv.Itoa(*req.ID)
	_, err = req.ReqInfo.Tx.Exec(query)
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

	// update req with current info
	err = req.ReqInfo.Tx.QueryRowx(selectDeliveryServiceRequestsQuery() + ` WHERE r.id=` + strconv.Itoa(*req.ID)).StructScan(req)
	if err != nil {
		log.Errorf("Error querying DeliveryServiceRequests: %v", err)
		return err, tc.SystemError
	}

	return nil, tc.NoError
}

// Validate is not needed when only Status is updated
func (req deliveryServiceRequestStatus) Validate() []error {
	return nil
}

// ChangeLogMessage implements the api.ChangeLogger interface for a custom log message
func (req deliveryServiceRequestStatus) ChangeLogMessage(action string) (string, error) {
	message := `Changed status of ‘` + req.getXMLID() + `’ ` + req.GetType() + ` to '` + string(*req.Status) + `'`
	return message, nil
}
