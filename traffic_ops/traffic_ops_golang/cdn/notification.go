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

func CreateNotification(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id"}, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	reqObj := tc.CDNNotificationRequest{}
	if err := json.NewDecoder(r.Body).Decode(&reqObj); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("malformed JSON: "+err.Error()), nil)
		return
	}
	if err := create(inf.Tx.Tx, int64(inf.IntParams["id"]), inf.User.UserName, reqObj.Notification); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("CDN create notification: "+err.Error()))
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
	api.CreateChangeLogRawTx(api.ApiChange, "CDN: "+string(cdnName)+", ID: "+strconv.Itoa(inf.IntParams["id"])+", ACTION: CDN notification created, Notification: " +reqObj.Notification, inf.User, inf.Tx.Tx)
	api.WriteRespAlertObj(w, r, tc.SuccessLevel, "CDN notification was created", tc.CDNNotificationResponse{CDNID: int64(inf.IntParams["id"]), Username: "bob", Notification: reqObj.Notification})
}

func DeleteNotification(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id"}, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	if err := delete(inf.Tx.Tx, int64(inf.IntParams["id"])); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("CDN delete notification: "+err.Error()))
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
	api.CreateChangeLogRawTx(api.ApiChange, "CDN: "+string(cdnName)+", ID: "+strconv.Itoa(inf.IntParams["id"])+", ACTION: CDN notification deleted", inf.User, inf.Tx.Tx)
	api.WriteRespAlertObj(w, r, tc.SuccessLevel, "CDN notification was deleted", nil)
}

func create(tx *sql.Tx, cdnID int64, username string, notification string) error {
	if _, err := tx.Exec(`UPDATE cdn SET notification_created_by = $1, notification = $2 WHERE id = $3`, username, notification, cdnID); err != nil {
		return errors.New("creating cdn notification: " + err.Error())
	}
	return nil
}

func delete(tx *sql.Tx, cdnID int64) error {
	if _, err := tx.Exec(`UPDATE cdn SET notification_created_by = $1, notification = $2 WHERE id = $3`, nil, nil, cdnID); err != nil {
		return errors.New("deleting cdn notification: " + err.Error())
	}
	return nil
}
