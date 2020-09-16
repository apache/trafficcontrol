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
	"errors"
	"fmt"
	"net/http"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/routing/middleware"
)

// GetStatus is the handler for GET requests to
// /deliveryservice_requests/{{ID}}/status.
func GetStatus(w http.ResponseWriter, r *http.Request) {
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

	// This should never happen because a route doesn't exist for it
	if version.Major < 3 {
		w.Header().Set("Allow", http.MethodPut)
		w.WriteHeader(http.StatusMethodNotAllowed)
		api.WriteRespAlert(w, r, tc.ErrorLevel, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	var dsr tc.DeliveryServiceRequestV30
	if err := inf.Tx.QueryRowx(selectQuery+"WHERE r.id=$1", inf.IntParams["id"]).StructScan(&dsr); err != nil {
		if err == sql.ErrNoRows {
			errCode = http.StatusNotFound
			userErr = fmt.Errorf("no such Delivery Service Request: %d", inf.IntParams["id"])
			sysErr = nil
		} else {
			errCode = http.StatusInternalServerError
			userErr = nil
			sysErr = fmt.Errorf("looking for DSR: %v", err)
		}
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	if dsr.ChangeType != tc.DSRChangeTypeDelete && dsr.IsOpen() && (dsr.Requested == nil || dsr.Requested.ID == nil) {
		sysErr = errors.New("retrieved open, non-delete, DSR that had nil Requested or Requested.ID")
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, sysErr)
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

	getOriginals([]int{*dsr.Requested.ID}, inf.Tx, map[int][]*tc.DeliveryServiceRequestV30{*dsr.Requested.ID: {&dsr}})

	api.WriteResp(w, r, dsr.Status)
}

type statusChangeRequest struct {
	Status tc.RequestStatus `json:"status"`
}

func (s *statusChangeRequest) Validate(*sql.Tx) error {
	return nil
}

const updateStatusQuery = `
UPDATE deliveryservice_request
SET status = $1
WHERE id = $2
RETURNING last_updated
`

func PutStatus(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id"}, []string{"id"})
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	var req statusChangeRequest
	if err := api.Parse(r.Body, tx, &req); err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
		return
	}

	// Middleware should've already handled this, so idk why this is a pointer at all tbh
	version := inf.Version
	if version == nil {
		middleware.NotImplementedHandler().ServeHTTP(w, r)
		return
	}

	var dsr tc.DeliveryServiceRequestV30
	if err := inf.Tx.QueryRowx(selectQuery+"WHERE r.id=$1", inf.IntParams["id"]).StructScan(&dsr); err != nil {
		if err == sql.ErrNoRows {
			errCode = http.StatusNotFound
			userErr = fmt.Errorf("no such Delivery Service Request: %d", inf.IntParams["id"])
			sysErr = nil
		} else {
			errCode = http.StatusInternalServerError
			userErr = nil
			sysErr = fmt.Errorf("looking for DSR: %v", err)
		}
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	if dsr.ChangeType != tc.DSRChangeTypeDelete && dsr.IsOpen() && (dsr.Requested == nil || dsr.Requested.ID == nil) {
		sysErr = errors.New("retrieved open, non-delete, DSR that had nil Requested or Requested.ID")
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, sysErr)
		return
	}
	dsr.SetXMLID()

	if err := dsr.Status.ValidTransition(req.Status); err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
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

	message := fmt.Sprintf("Changed status of '%s' Delivery Service Request from '%s' to '%s'", dsr.XMLID, dsr.Status, req.Status)
	dsr.Status = req.Status

	if err := tx.QueryRow(updateStatusQuery, req.Status, *dsr.ID).Scan(&dsr.LastUpdated); err != nil {
		userErr, sysErr, errCode = api.ParseDBError(err)
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	getOriginals([]int{*dsr.Requested.ID}, inf.Tx, map[int][]*tc.DeliveryServiceRequestV30{*dsr.Requested.ID: {&dsr}})

	var resp interface{}
	if inf.Version.Major >= 3 {
		resp = dsr
	} else {
		resp = dsr.Downgrade()
	}

	api.WriteRespAlertObj(w, r, tc.SuccessLevel, message, resp)
	message = fmt.Sprintf("Delivery Service Request: %d, ID: %d, ACTION: %s deliveryservice_request, keys: {id:%d }", *dsr.ID, *dsr.ID, message, *dsr.ID)
	inf.CreateChangeLog(message)
}
