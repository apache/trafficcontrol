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
	"net/http"
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/tenant"

	"github.com/jmoiron/sqlx"
)

// TODeliveryServiceRequest is the type alias to define functions on
type TODeliveryServiceRequest struct {
	api.APIInfoImpl `json:"-"`
	tc.DeliveryServiceRequestNullable
}

func (v *TODeliveryServiceRequest) SetLastUpdated(t tc.TimeNoMod) { v.LastUpdated = &t }
func (v *TODeliveryServiceRequest) InsertQuery() string           { return insertRequestQuery() }
func (v *TODeliveryServiceRequest) UpdateQuery() string           { return updateRequestQuery() }
func (v *TODeliveryServiceRequest) DeleteQuery() string {
	return `DELETE FROM deliveryservice_request WHERE id = :id`
}

func (req TODeliveryServiceRequest) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{"id", api.GetIntKey}}
}

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
func (req *TODeliveryServiceRequest) Read() ([]interface{}, error, error, int) {
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

	p := req.APIInfo().Params
	if _, ok := req.APIInfo().Params["orderby"]; !ok {
		// if orderby not provided, default to orderby xmlId.  Making a copy of parameters to not modify input arg
		p = make(map[string]string, len(req.APIInfo().Params))
		for k, v := range req.APIInfo().Params {
			p[k] = v
		}
		p["orderby"] = "xmlId"
	}

	where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(p, queryParamsToQueryCols)
	if len(errs) > 0 {
		return nil, util.JoinErrs(errs), nil, http.StatusBadRequest
	}
	tenantIDs, err := tenant.GetUserTenantIDListTx(req.APIInfo().Tx.Tx, req.APIInfo().User.TenantID)
	if err != nil {
		return nil, nil, errors.New("dsr getting tenant list: " + err.Error()), http.StatusInternalServerError
	}
	where, queryValues = dbhelpers.AddTenancyCheck(where, queryValues, "CAST(r.deliveryservice->>'tenantId' AS bigint)", tenantIDs)

	query := selectDeliveryServiceRequestsQuery() + where + orderBy + pagination
	log.Debugln("Query is ", query)

	rows, err := req.APIInfo().Tx.NamedQuery(query, queryValues)
	if err != nil {
		return nil, nil, errors.New("dsr querying: " + err.Error()), http.StatusInternalServerError
	}
	defer rows.Close()

	deliveryServiceRequests := []interface{}{}
	for rows.Next() {
		var s TODeliveryServiceRequest
		if err = rows.StructScan(&s); err != nil {
			return nil, nil, errors.New("dsr scanning: " + err.Error()), http.StatusInternalServerError
		}
		deliveryServiceRequests = append(deliveryServiceRequests, s)
	}

	return deliveryServiceRequests, nil, nil, http.StatusOK
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
	return tenant.IsResourceAuthorizedToUserTx(*ds.TenantID, user, req.APIInfo().Tx.Tx)
}

// Update implements the tc.Updater interface.
//all implementations of Updater should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a request with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
func (req *TODeliveryServiceRequest) Update() (error, error, int) {
	if req.ID == nil {
		return errors.New("missing id"), nil, http.StatusBadRequest
	}

	current := TODeliveryServiceRequest{}
	err := req.ReqInfo.Tx.QueryRowx(selectDeliveryServiceRequestsQuery()+`WHERE r.id=$1`, *req.ID).StructScan(&current)
	if err != nil {
		return nil, errors.New("dsr update querying: " + err.Error()), http.StatusInternalServerError
	}

	// Update can only change status between draft & submitted.  All other transitions must go thru
	// the PUT /api/<version>/deliveryservice_request/:id/status endpoint
	if current.Status == nil || req.Status == nil {
		return errors.New("Missing status for DeliveryServiceRequest"), nil, http.StatusBadRequest
	}

	if *current.Status != tc.RequestStatusDraft && *current.Status != tc.RequestStatusSubmitted {
		return fmt.Errorf("Cannot change DeliveryServiceRequest in '%s' status.", string(*current.Status)), nil, http.StatusBadRequest
	}

	if *req.Status != tc.RequestStatusDraft && *req.Status != tc.RequestStatusSubmitted {
		return fmt.Errorf("Cannot change DeliveryServiceRequest status from '%s' to '%s'", string(*current.Status), string(*req.Status)), nil, http.StatusBadRequest
	}

	userID := tc.IDNoMod(req.APIInfo().User.ID)
	req.LastEditedByID = &userID

	return api.GenericUpdate(req)
}

// Creator implements the tc.Creator interface
//all implementations of Creator should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a request with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
//The insert sql returns the id and lastUpdated values of the newly inserted request and have
//to be added to the struct
func (req *TODeliveryServiceRequest) Create() (error, error, int) {
	// TODO move to Validate()
	if req.Status == nil {
		return errors.New("missing status"), nil, http.StatusBadRequest
	}
	if *req.Status != tc.RequestStatusDraft && *req.Status != tc.RequestStatusSubmitted {
		return fmt.Errorf("invalid initial request status '%v'.  Must be '%v' or '%v'",
			*req.Status, tc.RequestStatusDraft, tc.RequestStatusSubmitted), nil, http.StatusBadRequest
	}
	// first, ensure there's not an active request with this XMLID
	ds := req.DeliveryService
	if ds == nil {
		return errors.New("no delivery service associated with this request"), nil, http.StatusBadRequest
	}
	if ds.XMLID == nil {
		return errors.New("no xmlId associated with this request"), nil, http.StatusBadRequest
	}
	XMLID := *ds.XMLID
	active, err := isActiveRequest(req.APIInfo().Tx, XMLID)
	if err != nil {
		return errors.New("checking request active: " + err.Error()), nil, http.StatusInternalServerError
	}
	if active {
		return errors.New(`An active request exists for delivery service '` + XMLID + `'`), nil, http.StatusBadRequest
	}

	userID := tc.IDNoMod(req.APIInfo().User.ID)
	req.AuthorID = &userID
	req.LastEditedByID = &userID

	return api.GenericCreate(req)
}

func (req *TODeliveryServiceRequest) Delete() (error, error, int) {
	if req.ID == nil {
		return errors.New("missing id"), nil, http.StatusBadRequest
	}

	st := tc.RequestStatus(0)
	if err := req.APIInfo().Tx.Tx.QueryRow(`SELECT status FROM deliveryservice_request WHERE id=$1`, *req.ID).Scan(&st); err != nil {
		return nil, errors.New("dsr delete querying status: " + err.Error()), http.StatusBadRequest
	}
	if st == tc.RequestStatusComplete || st == tc.RequestStatusPending || st == tc.RequestStatusRejected {
		return errors.New("cannot delete a deliveryservice_request with state " + string(st)), nil, http.StatusBadRequest
	}

	return api.GenericDelete(req)
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
func isActiveRequest(tx *sqlx.Tx, xmlID string) (bool, error) {
	qry := `SELECT EXISTS(SELECT 1 FROM deliveryservice_request WHERE deliveryservice->>'xmlId' = $1 AND status IN ('draft', 'submitted', 'pending'))`
	active := false
	if err := tx.QueryRow(qry, xmlID).Scan(&active); err != nil {
		return false, err
	}
	return active, nil
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

func GetAssignmentSingleton() api.CRUDer {
	return &deliveryServiceRequestAssignment{}
}

type deliveryServiceRequestAssignment struct {
	TODeliveryServiceRequest
}

// Update assignee only
func (req *deliveryServiceRequestAssignment) Update() (error, error, int) {
	// req represents the state the deliveryservice_request is to transition to
	// we want to limit what changes here -- only assignee can change
	if req.ID == nil {
		return errors.New("missing id"), nil, http.StatusBadRequest
	}

	current := TODeliveryServiceRequest{}
	err := req.ReqInfo.Tx.QueryRowx(selectDeliveryServiceRequestsQuery()+`WHERE r.id = $1`, *req.ID).StructScan(&current)
	if err != nil {
		return nil, errors.New("dsr assignment querying existing: " + err.Error()), http.StatusInternalServerError
	}

	// unchanged (maybe both nil)
	if current.AssigneeID == req.AssigneeID {
		log.Infof("dsr assignment update: assignee unchanged")
		return nil, nil, http.StatusOK
	}

	// Only assigneeID changes -- nothing else
	assigneeID := req.AssigneeID
	req.DeliveryServiceRequestNullable = current.DeliveryServiceRequestNullable
	req.AssigneeID = assigneeID

	// LastEditedBy field should not change with status update
	v := "null"
	if req.AssigneeID != nil {
		v = strconv.Itoa(*req.AssigneeID)
	}

	if _, err = req.APIInfo().Tx.Tx.Exec(`UPDATE deliveryservice_request SET assignee_id = $1 WHERE id = $2`, v, *req.ID); err != nil {
		return api.ParseDBError(err)
	}

	if err = req.APIInfo().Tx.QueryRowx(selectDeliveryServiceRequestsQuery()+` WHERE r.id = $1`, *req.ID).StructScan(req); err != nil {
		return nil, errors.New("dsr assignment querying: " + err.Error()), http.StatusInternalServerError
	}

	return nil, nil, http.StatusOK
}

func (req deliveryServiceRequestAssignment) Validate() error {
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

func GetStatusSingleton() api.CRUDer {
	return &deliveryServiceRequestStatus{}
}

// deliveryServiceRequestStatus implements interfaces needed to update the request status only
type deliveryServiceRequestStatus struct {
	TODeliveryServiceRequest
}

func (req *deliveryServiceRequestStatus) Update() (error, error, int) {
	// req represents the state the deliveryservice_request is to transition to
	// we want to limit what changes here -- only status can change,  and only according to the established rules
	// for status transition
	if req.ID == nil {
		return errors.New("missing id"), nil, http.StatusBadRequest
	}

	current := TODeliveryServiceRequest{}
	err := req.APIInfo().Tx.QueryRowx(selectDeliveryServiceRequestsQuery()+` WHERE r.id = $1`, *req.ID).StructScan(&current)
	if err != nil {
		return nil, errors.New("dsr status querying existing: " + err.Error()), http.StatusInternalServerError
	}

	if err = current.Status.ValidTransition(*req.Status); err != nil {
		return err, nil, http.StatusBadRequest // TODO verify err is secure to send to user
	}

	// keep everything else the same -- only update status
	st := req.Status
	req.DeliveryServiceRequestNullable = current.DeliveryServiceRequestNullable
	req.Status = st

	// LastEditedBy field should not change with status update

	if _, err = req.APIInfo().Tx.Tx.Exec(`UPDATE deliveryservice_request SET status = $1 WHERE id = $2`, *req.Status, *req.ID); err != nil {
		return api.ParseDBError(err)
	}

	if err = req.APIInfo().Tx.QueryRowx(selectDeliveryServiceRequestsQuery()+` WHERE r.id = $1`, *req.ID).StructScan(req); err != nil {
		return nil, errors.New("dsr status update querying: " + err.Error()), http.StatusInternalServerError
	}

	return nil, nil, http.StatusOK
}

// Validate is not needed when only Status is updated
func (req deliveryServiceRequestStatus) Validate() error {
	return nil
}

// ChangeLogMessage implements the api.ChangeLogger interface for a custom log message
func (req deliveryServiceRequestStatus) ChangeLogMessage(action string) (string, error) {
	message := `Changed status of ‘` + req.getXMLID() + `’ ` + req.GetType() + ` to '` + string(*req.Status) + `'`
	return message, nil
}
