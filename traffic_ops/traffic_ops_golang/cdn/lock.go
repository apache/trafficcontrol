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
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-tc"

	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
)

func Lock(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id"}, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	reqObj := tc.CDNLockRequest{}
	if err := json.NewDecoder(r.Body).Decode(&reqObj); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("malformed JSON: "+err.Error()), nil)
		return
	}
	if err := lock(inf.Tx.Tx, int64(inf.IntParams["id"]), inf.User.UserName, reqObj.Reason); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("CDN lock: "+err.Error()))
		return
	}

	cdnName, ok, err := dbhelpers.GetCDNNameFromID(inf.Tx.Tx, int64(inf.IntParams["id"]))
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting cdn name from ID '"+inf.Params["id"]+"': "+err.Error()))
		return
	} else if !ok {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, nil, nil)
		return
	}
	api.CreateChangeLogRawTx(api.ApiChange, "CDN: "+string(cdnName)+", ID: "+strconv.Itoa(inf.IntParams["id"])+", ACTION: CDN locked, REASON: " +reqObj.Reason, inf.User, inf.Tx.Tx)
	api.WriteResp(w, r, tc.CDNLockResponse{CDNID: int64(inf.IntParams["id"]), Username: "bob", Reason: reqObj.Reason})
}

func Unlock(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id"}, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	if err := unlock(inf.Tx.Tx, int64(inf.IntParams["id"])); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("CDN unlock: "+err.Error()))
		return
	}

	cdnName, ok, err := dbhelpers.GetCDNNameFromID(inf.Tx.Tx, int64(inf.IntParams["id"]))
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting cdn name from ID '"+inf.Params["id"]+"': "+err.Error()))
		return
	} else if !ok {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, nil, nil)
		return
	}
	api.CreateChangeLogRawTx(api.ApiChange, "CDN: "+string(cdnName)+", ID: "+strconv.Itoa(inf.IntParams["id"])+", ACTION: CDN unlocked", inf.User, inf.Tx.Tx)
	api.WriteResp(w, r, tc.CDNUnlockResponse{CDNID: int64(inf.IntParams["id"])})
}

func lock(tx *sql.Tx, cdnID int64, username string, reason string) error {
	if _, err := tx.Exec(`UPDATE cdn SET locked_by = $1, locked_reason = $2 WHERE id = $3`, username, reason, cdnID); err != nil {
		return errors.New("locking cdn: " + err.Error())
	}
	return nil
}

func unlock(tx *sql.Tx, cdnID int64) error {
	if _, err := tx.Exec(`UPDATE cdn SET locked_by = $1, locked_reason = $2 WHERE id = $3`, nil, nil, cdnID); err != nil {
		return errors.New("unlocking cdn: " + err.Error())
	}
	return nil
}

