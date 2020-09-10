package request

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
)

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
