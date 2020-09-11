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
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/routing/middleware"
)

// GetAssignment is the handler for GET requests to
// /deliveryservice_requests/{{ID}}/assign.
func GetAssignment(w http.ResponseWriter, r *http.Request) {
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

	api.WriteResp(w, r, dsr.Assignee)
}

func GetAssignmentSingleton() api.Updater {
	return &deliveryServiceRequestAssignment{}
}

type deliveryServiceRequestAssignment struct {
	api.APIInfoImpl `json:"-"`
	tc.DeliveryServiceRequestV15
}

func (req *deliveryServiceRequestAssignment) GetAuditName() string {
	if req != nil && req.ID != nil {
		return strconv.Itoa(*req.ID)
	}
	return "UNKNOWN"
}

func (req *deliveryServiceRequestAssignment) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{Field: "id", Func: api.GetIntKey}}
}

func (req *deliveryServiceRequestAssignment) GetKeys() (map[string]interface{}, bool) {
	keys := map[string]interface{}{"id": 0}
	success := false
	if req.ID != nil {
		keys["id"] = *req.ID
		success = true
	}
	return keys, success
}

func (req *deliveryServiceRequestAssignment) SetKeys(keys map[string]interface{}) {
	i, _ := keys["id"].(int)
	req.ID = &i
}

func (*deliveryServiceRequestAssignment) GetType() string {
	return "deliveryservice_request"
}

// Update assignee only
func (req *deliveryServiceRequestAssignment) Update() (error, error, int) {
	// req represents the state the deliveryservice_request is to transition to
	// we want to limit what changes here -- only assignee can change
	if req.ID == nil {
		return errors.New("missing id"), nil, http.StatusBadRequest
	}

	var current tc.DeliveryServiceRequestV30
	err := req.ReqInfo.Tx.QueryRowx(selectQuery+`WHERE r.id = $1`, *req.ID).StructScan(&current)
	if err != nil {
		return nil, errors.New("dsr assignment querying existing: " + err.Error()), http.StatusInternalServerError
	}

	// unchanged (maybe both nil)
	if current.AssigneeID == req.AssigneeID {
		log.Infof("dsr assignment update: assignee unchanged")
		return nil, nil, http.StatusOK
	}

	// LastEditedBy field should not change with status update
	if _, err = req.APIInfo().Tx.Tx.Exec(`UPDATE deliveryservice_request SET assignee_id = $1 WHERE id = $2`, req.AssigneeID, *current.ID); err != nil {
		return api.ParseDBError(err)
	}

	// Only assigneeID changes -- nothing else
	assigneeID := req.AssigneeID
	req.DeliveryServiceRequestV15 = current.Downgrade()
	req.AssigneeID = assigneeID

	if err = req.APIInfo().Tx.QueryRowx(selectQuery+` WHERE r.id = $1`, *req.ID).StructScan(&current); err != nil {
		return nil, errors.New("dsr assignment querying: " + err.Error()), http.StatusInternalServerError
	}

	req.DeliveryServiceRequestV15 = current.Downgrade()

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
	XMLID := "UNKNOWN"
	if req.XMLID != nil {
		XMLID = *req.XMLID
	} else if req.DeliveryService != nil && req.DeliveryService.XMLID != nil {
		XMLID = *req.DeliveryService.XMLID
	}
	message := fmt.Sprintf("Changed assignee of '%s' Delivery Service Request to '%s'", XMLID, a)

	return message, nil
}
