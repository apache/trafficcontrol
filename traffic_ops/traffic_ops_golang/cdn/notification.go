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
	"github.com/apache/trafficcontrol/lib/go-tc"
	"net/http"

	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
)

func CreateNotification(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"name"}, nil)
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
	ok, err := dbhelpers.CDNExists(inf.Params["name"], inf.Tx.Tx)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("CDN create notification: "+err.Error()))
		return
	} else if !ok {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, nil, nil)
		return
	}
	if err := create(inf.Tx.Tx, inf.Params["name"], inf.User.UserName, reqObj.Notification); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("CDN create notification: "+err.Error()))
		return
	}
	api.CreateChangeLogRawTx(api.ApiChange, "CDN: "+inf.Params["name"]+", ACTION: CDN notification created, Notification: "+reqObj.Notification, inf.User, inf.Tx.Tx)
	api.WriteRespAlertObj(w, r, tc.SuccessLevel, "CDN notification was created", tc.CDNNotificationResponse{Username: inf.User.UserName, Notification: reqObj.Notification})
}

func DeleteNotification(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"name"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	ok, err := dbhelpers.CDNExists(inf.Params["name"], inf.Tx.Tx)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("CDN delete notification: "+err.Error()))
		return
	} else if !ok {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, nil, nil)
		return
	}
	if err := delete(inf.Tx.Tx, inf.Params["name"]); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("CDN delete notification: "+err.Error()))
		return
	}
	api.CreateChangeLogRawTx(api.ApiChange, "CDN: "+inf.Params["name"]+", ACTION: CDN notification deleted", inf.User, inf.Tx.Tx)
	api.WriteRespAlertObj(w, r, tc.SuccessLevel, "CDN notification was deleted", nil)
}

func create(tx *sql.Tx, cdnName string, username string, notification string) error {
	if _, err := tx.Exec(`UPDATE cdn SET notification_created_by = $1, notification = $2 WHERE name = $3`, username, notification, cdnName); err != nil {
		return errors.New("creating cdn notification: " + err.Error())
	}
	return nil
}

func delete(tx *sql.Tx, cdnName string) error {
	if _, err := tx.Exec(`UPDATE cdn SET notification_created_by = $1, notification = $2 WHERE name = $3`, nil, nil, cdnName); err != nil {
		return errors.New("deleting cdn notification: " + err.Error())
	}
	return nil
}
