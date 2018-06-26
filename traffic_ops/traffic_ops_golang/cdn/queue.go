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
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/auth"
)

type QueueReq struct {
	Action string `json:"action"`
}

type QueueResp struct {
	Action string `json:"action"`
	CDNID  int64  `json:"cdnId"`
}

func Queue(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.GetCurrentUser(r.Context())
		if err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("getting user: "+err.Error()))
			return
		}
		params, intParams, userErr, sysErr, errCode := api.AllParams(r, []string{"id"}, []string{"id"})
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, errCode, userErr, sysErr)
			return
		}
		reqObj := QueueReq{}
		if err := json.NewDecoder(r.Body).Decode(&reqObj); err != nil {
			api.HandleErr(w, r, http.StatusBadRequest, errors.New("malformed JSON: "+err.Error()), nil)
			return
		}
		if reqObj.Action != "queue" && reqObj.Action != "dequeue" {
			api.HandleErr(w, r, http.StatusBadRequest, errors.New("action must be 'queue' or 'dequeue'"), nil)
			return
		}
		if err := queueUpdates(db, int64(intParams["id"]), reqObj.Action == "queue"); err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("CDN queueing updates: "+err.Error()))
			return
		}
		api.WriteResp(w, r, QueueResp{Action: reqObj.Action, CDNID: int64(intParams["id"])})
		api.CreateChangeLogRaw(api.ApiChange, "Server updates "+reqObj.Action+"d for cdn "+params["id"], user, db)
	}
}

func queueUpdates(db *sql.DB, cdnID int64, queue bool) error {
	if _, err := db.Exec(`UPDATE server SET upd_pending = $1 WHERE server.cdn_id = $2`, queue, cdnID); err != nil {
		return errors.New("querying queue updates: " + err.Error())
	}
	return nil
}
