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
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"

	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/routing/middleware"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/tenant"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/util/ims"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

const selectQuery = `
SELECT
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
	r.status
FROM deliveryservice_request r
JOIN tm_user a ON r.author_id = a.id
LEFT OUTER JOIN tm_user s ON r.assignee_id = s.id
LEFT OUTER JOIN tm_user e ON r.last_edited_by_id = e.id
`

const insertQuery = `
INSERT INTO deliveryservice_request (
	assignee_id,
	author_id,
	change_type,
	last_edited_by_id,
	deliveryservice,
	status
) VALUES (
	$1,
	$2,
	$3,
	$2,
	$4,
	$5
)
RETURNING
	id,
	last_updated,
	created_at
`

const updateQuery = `
UPDATE deliveryservice_request
SET
	assignee_id = $1,
	change_type = $2,
	last_edited_by_id = $3,
	deliveryservice = $4,
	status = $5
WHERE id = $6
RETURNING
	last_updated,
	created_at
`

const deleteQuery = `
DELETE
FROM deliveryservice_request
WHERE id=$1
`

//TODO: figure out how to modify 'AddTenancyCheck' so this isn't necessary
const customTenancyCheck = `(
	CAST(r.deliveryservice->>'tenantId' AS BIGINT) = ANY(CAST(:accessibleTenants AS BIGINT[]))
)`

func selectMaxLastUpdatedQuery(where string) string {
	return `SELECT max(t) FROM (
		SELECT max(r.last_updated) as t FROM deliveryservice_request r
	JOIN tm_user a ON r.author_id = a.id
	LEFT OUTER JOIN tm_user s ON r.assignee_id = s.id
	LEFT OUTER JOIN tm_user e ON r.last_edited_by_id = e.id ` + where +
		` UNION ALL
	SELECT MAX(last_updated) AS t FROM last_deleted l WHERE l.table_name='deliveryservice_request') AS res`
}

// Get is the GET handler for /deliveryservice_requests.
func Get(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	// Middleware should've already handled this, so idk why this is a pointer at all tbh
	version := inf.Version
	if version == nil {
		middleware.NotImplementedHandler().ServeHTTP(w, r)
		return
	}

	queryParamsToQueryCols := map[string]dbhelpers.WhereColumnInfo{
		"assignee":   {Column: "s.username"},
		"assigneeId": {Column: "r.assignee_id", Checker: api.IsInt},
		"author":     {Column: "a.username"},
		"authorId":   {Column: "r.author_id", Checker: api.IsInt},
		"changeType": {Column: "r.change_type"},
		"createdAt":  {Column: "r.created_at"},
		"id":         {Column: "r.id", Checker: api.IsInt},
		"status":     {Column: "r.status"},
	}
	if _, ok := inf.Params["orderby"]; !ok {
		inf.Params["orderby"] = "xmlId"
	}

	where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(inf.Params, queryParamsToQueryCols)
	if len(errs) > 0 {
		api.HandleErr(w, r, tx, http.StatusBadRequest, util.JoinErrs(errs), nil)
		return
	}

	// TODO: add this functionality to the query builder in dbhelpers
	if xmlID, ok := inf.Params["xmlId"]; ok {
		queryValues["xmlId"] = xmlID
		if where != "" {
			where += " AND "
		} else {
			where = "WHERE "
		}
		where += "(r.deliveryservice->>'xmlId' = :xmlId)"
	}

	var maxTime *time.Time
	if inf.UseIMS() {
		maxTime = new(time.Time)
		var runSecond bool
		runSecond, *maxTime = ims.TryIfModifiedSinceQuery(inf.Tx, r.Header, queryValues, selectMaxLastUpdatedQuery(where))
		if !runSecond {
			log.Debugln("IMS HIT")
			api.WriteIMSHitResp(w, r, *maxTime)
			return
		}
		log.Debugln("IMS MISS")
	} else {
		log.Debugln("Non IMS request")
	}

	tenantIDs, err := tenant.GetUserTenantIDListTx(tx, inf.User.TenantID)
	if err != nil {
		sysErr = fmt.Errorf("dsr getting tenant list: %v", err)
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, sysErr)
		return
	}

	if where == "" {
		where = "WHERE "
	} else {
		where += " AND "
	}
	where += customTenancyCheck
	queryValues["accessibleTenants"] = pq.Array(tenantIDs)

	query := selectQuery + where + orderBy + pagination
	log.Debugln("Query is ", query)

	rows, err := inf.Tx.NamedQuery(query, queryValues)
	if err != nil {
		sysErr = fmt.Errorf("dsr querying: %v", err)
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, sysErr)
		return
	}
	defer rows.Close()

	dsrs := []tc.DeliveryServiceRequestV40{}
	for rows.Next() {
		var dsr tc.DeliveryServiceRequestV40
		if err = rows.StructScan(&dsr); err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("dsr scanning: %v", err))
			return
		}
		dsrs = append(dsrs, dsr)
	}

	if maxTime != nil {
		w.Header().Set(rfc.LastModified, maxTime.Format(rfc.LastModifiedFormat))
	}

	if version.Major >= 4 {
		api.WriteResp(w, r, dsrs)
	}

	downgraded := make([]tc.DeliveryServiceRequestNullable, 0, len(dsrs))
	for _, dsr := range dsrs {
		downgraded = append(downgraded, dsr.Downgrade())
	}

	api.WriteResp(w, r, downgraded)
}

// isTenantAuthorized ensures the user is authorized on the DSR's
// DeliveryService's Tenant, as appropriate to the change type.
func isTenantAuthorized(dsr tc.DeliveryServiceRequestV40, inf *api.APIInfo) (bool, error) {
	if dsr.DeliveryService != nil && (dsr.ChangeType == tc.DSRChangeTypeUpdate || dsr.ChangeType == tc.DSRChangeTypeCreate) {
		if dsr.DeliveryService.TenantID == nil {
			return false, errors.New("deliveryService.tenantID is nil")
		}
		ok, err := tenant.IsResourceAuthorizedToUserTx(*dsr.DeliveryService.TenantID, inf.User, inf.Tx.Tx)
		if !ok || err != nil {
			return ok, err
		}
	}

	return true, nil
}

// Warning: this assumes inf isn't nil, and neither is dsr, inf.Tx or inf.User or inf.Tx.Tx.
func insert(dsr *tc.DeliveryServiceRequestV40, inf *api.APIInfo) (int, error, error) {
	dsr.Author = inf.User.UserName
	dsr.LastEditedBy = inf.User.UserName

	dsr.ID = new(int)
	if err := inf.Tx.Tx.QueryRow(insertQuery, dsr.AssigneeID, inf.User.ID, dsr.ChangeType, dsr.DeliveryService, dsr.Status).Scan(dsr.ID, &dsr.LastUpdated, &dsr.CreatedAt); err != nil {
		userErr, sysErr, errCode := api.ParseDBError(err)
		return errCode, userErr, sysErr
	}
	return http.StatusOK, nil, nil
}

// dsrManipulationResult encodes the result of manipulating a DSR.
type dsrManipulationResult struct {
	// Action is the action performed to manipulate the DSR.
	Action string
	// Assignee is a pointer to the name of the user assigned to a DSR - or nil
	// if there isn't one.
	Assignee *string
	// ChangeType is the DSR's change type.
	ChangeType tc.DSRChangeType
	// Successful is whether or not the manipulation encountered no errors.
	Successful bool
	// XMLID is the XMLID of the Delivery Service affected by the DSR.
	XMLID string
}

// String constructs a changelog message for the result.
// Unsuccessful results do not have a changelog message.
func (d dsrManipulationResult) String() string {
	if !d.Successful {
		return ""
	}

	var builder strings.Builder
	builder.WriteString(d.Action)
	builder.Write([]byte(" Delivery Service Request of type "))
	builder.WriteString(d.ChangeType.String())
	builder.Write([]byte(" for Delivery Service '"))
	builder.WriteString(d.XMLID)
	builder.WriteRune('\'')

	if d.Assignee != nil {
		builder.Write([]byte(" (assigned to user "))
		builder.WriteString(*d.Assignee)
		builder.WriteRune(')')
	}

	return builder.String()
}

func createV4(w http.ResponseWriter, r *http.Request, inf *api.APIInfo) (result dsrManipulationResult) {
	tx := inf.Tx.Tx
	var dsr tc.DeliveryServiceRequestV40
	if err := json.NewDecoder(r.Body).Decode(&dsr); err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("decoding: %v", err), nil)
		return
	}
	if userErr, sysErr := validateV4(dsr, tx); userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, userErr, sysErr)
		return
	}

	if dsr.Status != tc.RequestStatusDraft && dsr.Status != tc.RequestStatusSubmitted {
		userErr := fmt.Errorf("invalid initial request status '%s' - must be '%s' or '%s'", dsr.Status, tc.RequestStatusDraft, tc.RequestStatusSubmitted)
		api.HandleErr(w, r, tx, http.StatusBadRequest, userErr, nil)
		return
	}

	ok, err := isTenantAuthorized(dsr, inf)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
		return
	}
	if !ok {
		api.HandleErr(w, r, tx, http.StatusForbidden, errors.New("not authorized on this tenant"), nil)
		return
	}

	dsr.SetXMLID()
	if ok, err = dbhelpers.DSRExistsWithXMLID(dsr.XMLID, tx); err != nil {
		err = fmt.Errorf("checking for existence of DSR with xmlid '%s'", dsr.XMLID)
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
		return
	} else if ok {
		userErr := fmt.Errorf("an open Delivery Service Request for XMLID '%s' already exists", dsr.XMLID)
		api.HandleErr(w, r, tx, http.StatusBadRequest, userErr, nil)
		return
	}

	errCode, userErr, sysErr := insert(&dsr, inf)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/api/%d.%d/deliveryservice_requests/%d", inf.Version.Major, inf.Version.Minor, *dsr.ID))
	w.WriteHeader(http.StatusCreated)
	api.WriteRespAlertObj(w, r, tc.SuccessLevel, "Delivery Service request created", dsr)

	result.Successful = true
	result.Assignee = dsr.Assignee
	result.XMLID = dsr.XMLID
	result.ChangeType = dsr.ChangeType
	result.Action = api.Created
	return
}

func createLegacy(w http.ResponseWriter, r *http.Request, inf *api.APIInfo) (result dsrManipulationResult) {
	tx := inf.Tx.Tx
	var dsr tc.DeliveryServiceRequestNullable
	if err := json.NewDecoder(r.Body).Decode(&dsr); err != nil {
		userErr := fmt.Errorf("decoding: %v", err)
		api.HandleErr(w, r, tx, http.StatusBadRequest, userErr, nil)
		return
	}
	if err := validateLegacy(dsr, tx); err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
		return
	}

	upgraded := dsr.Upgrade()
	authorized, err := isTenantAuthorized(upgraded, inf)
	if err != nil {
		sysErr := fmt.Errorf("checking tenant authorized: %v", err)
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, sysErr)
		return
	}
	if !authorized {
		userErr := errors.New("not authorized on this tenant")
		api.HandleErr(w, r, tx, http.StatusForbidden, userErr, nil)
		return
	}

	if *dsr.Status != tc.RequestStatusDraft && *dsr.Status != tc.RequestStatusSubmitted {
		userErr := fmt.Errorf("invalid initial request status '%v'. Must be '%v' or '%v'", *dsr.Status, tc.RequestStatusDraft, tc.RequestStatusSubmitted)
		api.HandleErr(w, r, tx, http.StatusBadRequest, userErr, nil)
		return
	}

	// first, ensure there's not an active request with this XMLID
	ds := dsr.DeliveryService
	if ds == nil {
		userErr := errors.New("no delivery service associated with this request")
		api.HandleErr(w, r, tx, http.StatusBadRequest, userErr, nil)
		return
	}

	if ds.XMLID == nil {
		userErr := errors.New("no XMLID associated with this request")
		api.HandleErr(w, r, tx, http.StatusBadRequest, userErr, nil)
		return
	}
	XMLID := *ds.XMLID
	active, err := isActiveRequest(inf.Tx, XMLID)
	if err != nil {
		sysErr := fmt.Errorf("checking request active: %v", err)
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, sysErr)
		return
	}
	if active {
		userErr := fmt.Errorf("an active request exists for Delivery Service '%s'", XMLID)
		api.HandleErr(w, r, tx, http.StatusBadRequest, userErr, nil)
		return
	}

	errCode, userErr, sysErr := insert(&upgraded, inf)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	api.WriteRespAlertObj(w, r, tc.SuccessLevel, "Delivery Service request created", upgraded.Downgrade())

	result.Successful = true
	result.Assignee = dsr.Assignee
	result.XMLID = upgraded.XMLID
	result.ChangeType = upgraded.ChangeType
	result.Action = api.Created
	return
}

// Post is the handler for POST requests to /deliveryservice_requests.
func Post(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	// Middleware should've already handled this, so idk why this is a pointer at all tbh
	version := inf.Version
	if version == nil {
		middleware.NotImplementedHandler().ServeHTTP(w, r)
		return
	}
	if inf.User == nil {
		sysErr = errors.New("no user in API Info")
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, sysErr)
		return
	}

	var result dsrManipulationResult
	if version.Major >= 4 {
		result = createV4(w, r, inf)
	} else {
		result = createLegacy(w, r, inf)
	}

	if result.Successful {
		inf.CreateChangeLog(result.String())
	}
}

// Delete is the handler for DELETE requests to /deliveryservice_requests.
func Delete(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id"}, []string{"id"})
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	// Middleware should've already handled this, so idk why this is a pointer at all tbh
	version := inf.Version
	if version == nil {
		middleware.NotImplementedHandler().ServeHTTP(w, r)
		return
	}
	if inf.User == nil {
		sysErr = errors.New("no user in API Info")
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, sysErr)
		return
	}

	var dsr tc.DeliveryServiceRequestV40
	if err := inf.Tx.QueryRowx(selectQuery+"WHERE r.id=$1", inf.IntParams["id"]).StructScan(&dsr); err != nil {
		if err == sql.ErrNoRows {
			errCode = http.StatusNotFound
			userErr = fmt.Errorf("no such Delivery Service Request: #%d", inf.IntParams["id"])
			sysErr = nil
		} else {
			userErr, sysErr, errCode = api.ParseDBError(err)
		}
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	dsr.SetXMLID()

	authorized, err := isTenantAuthorized(dsr, inf)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
		return
	}
	if !authorized {
		api.HandleErr(w, r, tx, http.StatusForbidden, errors.New("not authorized on this tenant"), nil)
		return
	}

	result, err := tx.Exec(deleteQuery, inf.IntParams["id"])
	if err != nil {
		sysErr = fmt.Errorf("deleting DSR #%d: %v", inf.IntParams["id"], err)
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, sysErr)
		return
	}
	if affected, err := result.RowsAffected(); err != nil {
		sysErr = fmt.Errorf("checking affected rows: %v", err)
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, sysErr)
		return
	} else if affected != 1 {
		sysErr = fmt.Errorf("incorrect number of rows affected by delete: %d", affected)
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, sysErr)
		return
	}

	var resp interface{}
	if inf.Version.Major >= 4 {
		resp = dsr
	} else {
		resp = dsr.Downgrade()
	}

	api.WriteRespAlertObj(w, r, tc.SuccessLevel, fmt.Sprintf("Delivery Service Request #%d deleted", inf.IntParams["id"]), resp)

	res := dsrManipulationResult{
		Successful: true,
		XMLID:      dsr.XMLID,
		Action:     api.Deleted,
		Assignee:   dsr.Assignee,
		ChangeType: dsr.ChangeType,
	}
	inf.CreateChangeLog(res.String())
}

func putV40(w http.ResponseWriter, r *http.Request, inf *api.APIInfo) (result dsrManipulationResult) {
	tx := inf.Tx.Tx
	var dsr tc.DeliveryServiceRequestV40
	if err := json.NewDecoder(r.Body).Decode(&dsr); err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("decoding: %v", err), nil)
		return
	}
	if userErr, sysErr := validateV4(dsr, tx); userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, userErr, sysErr)
		return
	}

	if dsr.Status != tc.RequestStatusDraft && dsr.Status != tc.RequestStatusSubmitted {
		userErr := fmt.Errorf("cannot change DeliveryServiceRequest status to '%s'", dsr.Status)
		api.HandleErr(w, r, tx, http.StatusBadRequest, userErr, nil)
		return
	}

	authorized, err := isTenantAuthorized(dsr, inf)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
		return
	}
	if !authorized {
		api.HandleErr(w, r, tx, http.StatusForbidden, errors.New("not authorized on this tenant"), nil)
		return
	}

	dsr.LastEditedBy = inf.User.UserName
	dsr.LastEditedByID = new(int)
	*dsr.LastEditedByID = inf.User.ID

	args := []interface{}{
		dsr.AssigneeID,
		dsr.ChangeType,
		inf.User.ID,
		dsr.DeliveryService,
		dsr.Status,
		inf.IntParams["id"],
	}
	if err := tx.QueryRow(updateQuery, args...).Scan(&dsr.CreatedAt, &dsr.LastUpdated); err != nil {
		var userErr, sysErr error
		var errCode int
		if err == sql.ErrNoRows {
			userErr = fmt.Errorf("no such Delivery Service Request: #%d", inf.IntParams["id"])
			errCode = http.StatusNotFound
			sysErr = fmt.Errorf("running update query for Delivery Service Requests: %v", err)
		} else {
			userErr, sysErr, errCode = api.ParseDBError(err)
		}
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	dsr.SetXMLID()

	api.WriteRespAlertObj(w, r, tc.SuccessLevel, fmt.Sprintf("Delivery Service Request #%d updated", inf.IntParams["id"]), dsr)
	result.Successful = true
	result.Action = "Updated"
	result.Assignee = dsr.Assignee
	result.ChangeType = dsr.ChangeType
	result.XMLID = dsr.XMLID
	return
}

func putLegacy(w http.ResponseWriter, r *http.Request, inf *api.APIInfo) (result dsrManipulationResult) {
	tx := inf.Tx.Tx
	var dsr tc.DeliveryServiceRequestNullable
	if err := json.NewDecoder(r.Body).Decode(&dsr); err != nil {
		userErr := fmt.Errorf("decoding: %v", err)
		api.HandleErr(w, r, tx, http.StatusBadRequest, userErr, nil)
		return
	}
	if err := validateLegacy(dsr, tx); err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
		return
	}

	if *dsr.Status != tc.RequestStatusDraft && *dsr.Status != tc.RequestStatusSubmitted {
		userErr := fmt.Errorf("cannot change DeliveryServiceRequest status to '%s'", dsr.Status)
		api.HandleErr(w, r, tx, http.StatusBadRequest, userErr, nil)
		return
	}

	dsr.LastEditedBy = new(string)
	*dsr.LastEditedBy = inf.User.UserName
	dsr.LastEditedByID = new(tc.IDNoMod)
	*dsr.LastEditedByID = tc.IDNoMod(inf.User.ID)

	upgraded := dsr.Upgrade()

	authorized, err := isTenantAuthorized(upgraded, inf)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
		return
	}
	if !authorized {
		api.HandleErr(w, r, tx, http.StatusForbidden, errors.New("not authorized on this tenant"), nil)
		return
	}

	args := []interface{}{
		upgraded.AssigneeID,
		upgraded.ChangeType,
		inf.User.ID,
		upgraded.DeliveryService,
		upgraded.Status,
		inf.IntParams["id"],
	}
	if err := tx.QueryRow(updateQuery, args...).Scan(&dsr.CreatedAt, &dsr.LastUpdated); err != nil {
		var errCode int
		var userErr, sysErr error
		if err == sql.ErrNoRows {
			errCode = http.StatusNotFound
			userErr = fmt.Errorf("no such Delivery Service Request: #%d", inf.IntParams["id"])
			sysErr = nil
		} else {
			userErr, sysErr, errCode = api.ParseDBError(err)
		}
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	upgraded.SetXMLID()

	api.WriteRespAlertObj(w, r, tc.SuccessLevel, fmt.Sprintf("Delivery Service Request #%d updated", inf.IntParams["id"]), dsr)
	result.Action = api.Updated
	result.Assignee = dsr.Assignee
	result.ChangeType = upgraded.ChangeType
	result.Successful = true
	result.XMLID = upgraded.XMLID
	return
}

// Put is the handler for PUT requests to /deliveryservice_requests.
func Put(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id"}, []string{"id"})
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	// Middleware should've already handled this, so idk why this is a pointer at all tbh
	version := inf.Version
	if version == nil {
		middleware.NotImplementedHandler().ServeHTTP(w, r)
		return
	}
	if inf.User == nil {
		sysErr = errors.New("no user in API Info")
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, sysErr)
		return
	}

	id := inf.IntParams["id"]

	errCode, userErr, sysErr = inf.CheckPrecondition(selectMaxLastUpdatedQuery("WHERE r.id = $1"), id)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	var current tc.DeliveryServiceRequestV40
	if err := inf.Tx.QueryRowx(selectQuery+"WHERE r.id=$1", id).StructScan(&current); err != nil {
		if err == sql.ErrNoRows {
			errCode = http.StatusNotFound
			userErr = fmt.Errorf("no such Delivery Service Request: #%d", inf.IntParams["id"])
			sysErr = nil
		} else {
			userErr, sysErr, errCode = api.ParseDBError(err)
		}
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	authorized, err := isTenantAuthorized(current, inf)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
		return
	}
	if !authorized {
		api.HandleErr(w, r, tx, http.StatusForbidden, errors.New("not authorized on this tenant"), nil)
		return
	}

	if !current.IsOpen() {
		userErr = fmt.Errorf("cannot change DeliveryServiceRequest in '%s' status", current.Status)
		api.HandleErr(w, r, tx, http.StatusBadRequest, userErr, nil)
		return
	}

	var result dsrManipulationResult
	if inf.Version.Major >= 4 {
		result = putV40(w, r, inf)
	} else {
		result = putLegacy(w, r, inf)
	}

	if result.Successful {
		inf.CreateChangeLog(result.String())
	}
}

// isActiveRequest returns true if a request using this XMLID is currently in an active state.
func isActiveRequest(tx *sqlx.Tx, xmlID string) (bool, error) {
	qry := `SELECT EXISTS(SELECT 1 FROM deliveryservice_request WHERE deliveryservice->>'xmlId' = $1 AND status IN ('draft', 'submitted', 'pending'))`
	active := false
	if err := tx.QueryRow(qry, xmlID).Scan(&active); err != nil {
		return false, err
	}
	return active, nil
}
