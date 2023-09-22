package server

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
	"net/http"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/dbhelpers"
)

// QueueUpdateHandler implements an http handler that sets a server's
// config update time value.
func QueueUpdateHandler(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id"}, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	var reqObj tc.ServerQueueUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&reqObj); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, fmt.Errorf("malformed JSON: %v", err), nil)
		return
	}

	if reqObj.Action != "queue" && reqObj.Action != "dequeue" {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("action must be 'queue' or 'dequeue'"), nil)
		return
	}

	serverID := int64(inf.IntParams["id"])
	cdnName, err := dbhelpers.GetCDNNameFromServerID(inf.Tx.Tx, serverID)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
		return
	}
	if reqObj.Action == "queue" {
		userErr, sysErr, statusCode := dbhelpers.CheckIfCurrentUserHasCdnLock(inf.Tx.Tx, string(cdnName), inf.User.UserName)
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, inf.Tx.Tx, statusCode, userErr, sysErr)
			return
		}
	}

	if reqObj.Action == "queue" {
		err = dbhelpers.QueueUpdateForServer(inf.Tx.Tx, serverID)
	} else {
		err = dbhelpers.DequeueUpdateForServer(inf.Tx.Tx, serverID)
	}
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("queueing updates: %v", err))
		return
	}

	err = api.CreateChangeLogBuildMsg(
		api.ApiChange,
		api.Updated,
		inf.User,
		inf.Tx.Tx,
		"server",
		fmt.Sprint(serverID),
		map[string]interface{}{"id": serverID},
	)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("writing changelog: %v", err))
		return
	}

	api.WriteResp(w, r, tc.ServerQueueUpdate{
		ServerID: util.JSONIntStr(serverID),
		Action:   reqObj.Action,
	})
}
