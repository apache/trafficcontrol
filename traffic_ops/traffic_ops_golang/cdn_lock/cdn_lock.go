package cdn_lock

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

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
)

const readQuery = `SELECT username, cdn, message, soft, last_updated FROM cdn_lock`
const insertQuery = `INSERT INTO cdn_lock (username, cdn, message, soft) VALUES (:username, :cdn, :message, :soft) RETURNING username, cdn, message, soft, last_updated`
const deleteQuery = `DELETE FROM cdn_lock WHERE cdn=$1 AND username=$2 RETURNING username, cdn, message, soft, last_updated`
const deleteAdminQuery = `DELETE FROM cdn_lock WHERE cdn=$1 RETURNING username, cdn, message, soft, last_updated`

func Read(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	cols := map[string]dbhelpers.WhereColumnInfo{
		"cdn":      {"cdn_lock.cdn", nil},
		"username": {"cdn_lock.username", nil},
	}

	where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(inf.Params, cols)
	if len(errs) > 0 {
		errCode = http.StatusBadRequest
		userErr = util.JoinErrs(errs)
		api.HandleErr(w, r, tx, errCode, userErr, nil)
		return
	}

	var cdnLock []tc.CdnLock
	query := readQuery + where + orderBy + pagination
	rows, err := inf.Tx.NamedQuery(query, queryValues)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("querying cdn locks: "+err.Error()))
		return
	}
	defer rows.Close()

	for rows.Next() {
		var cLock tc.CdnLock
		if err = rows.Scan(&cLock.UserName, &cLock.Cdn, &cLock.Message, &cLock.Soft, &cLock.LastUpdated); err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("scanning cdn locks: "+err.Error()))
			return
		}
		cdnLock = append(cdnLock, cLock)
	}

	api.WriteResp(w, r, cdnLock)
}

func Create(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	tx := inf.Tx.Tx

	var cdnLock tc.CdnLock
	if err := json.NewDecoder(r.Body).Decode(&cdnLock); err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
		return
	}
	// by default, always create soft (or shared) locks
	if cdnLock.Soft == nil {
		cdnLock.Soft = util.BoolPtr(true)
	}

	c, err := api.GetConfig(r.Context())
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
		return
	}
	u, userErr, sysErr, errCode := api.GetUserFromReq(w, r, c.Secrets[0])
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	cdnLock.UserName = u.UserName
	resultRows, err := inf.Tx.NamedQuery(insertQuery, cdnLock)
	if err != nil {
		userErr, sysErr, errCode := api.ParseDBError(err)
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer resultRows.Close()

	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err := resultRows.Scan(&cdnLock.UserName, &cdnLock.Cdn, &cdnLock.Message, &cdnLock.Soft, &cdnLock.LastUpdated); err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("cdn lock create: scanning locks: "+err.Error()))
			return
		}
	}
	if rowsAffected == 0 {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("cdn lock create: lock couldn't be acquired"))
		return
	}

	alerts := tc.CreateAlerts(tc.SuccessLevel, "CDN lock acquired!")
	api.WriteAlertsObj(w, r, http.StatusCreated, alerts, cdnLock)

	changeLogMsg := fmt.Sprintf("USER: %s, CDN: %s, ACTION: Lock Acquired", u.UserName, cdnLock.Cdn)
	api.CreateChangeLogRawTx(api.ApiChange, changeLogMsg, inf.User, tx)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"cdn"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	cdn := inf.Params["cdn"]

	c, err := api.GetConfig(r.Context())
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
		return
	}
	u, userErr, sysErr, errCode := api.GetUserFromReq(w, r, c.Secrets[0])
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}

	tx := inf.Tx.Tx
	var result tc.CdnLock

	err = inf.Tx.Tx.QueryRow(deleteQuery, cdn, u.UserName).Scan(&result.UserName, &result.Cdn, &result.Message, &result.Soft, &result.LastUpdated)
	if err != nil {
		if err == sql.ErrNoRows {
			api.HandleErr(w, r, tx, http.StatusNotFound, errors.New(fmt.Sprintf("deleting cdn lock with cdn name %s: lock not found", cdn)), nil)
			return
		}
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New(fmt.Sprintf("deleting cdn lock with cdn name %s : %v", cdn, err.Error())))
		return
	}

	alerts := tc.CreateAlerts(tc.SuccessLevel, "cdn lock deleted")
	api.WriteAlertsObj(w, r, http.StatusOK, alerts, result)

	changeLogMsg := fmt.Sprintf("USER: %s, CDN: %s, ACTION: Lock Released", u.UserName, cdn)
	api.CreateChangeLogRawTx(api.ApiChange, changeLogMsg, inf.User, tx)
}

func AdminDelete(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"cdn"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	cdn := inf.Params["cdn"]
	tx := inf.Tx.Tx
	var result tc.CdnLock
	err := inf.Tx.Tx.QueryRow(deleteAdminQuery, cdn).Scan(&result.UserName, &result.Cdn, &result.Message, &result.Soft, &result.LastUpdated)
	if err != nil {
		if err == sql.ErrNoRows {
			api.HandleErr(w, r, tx, http.StatusNotFound, errors.New(fmt.Sprintf("deleting cdn lock with cdn name %s: lock not found", cdn)), nil)
			return
		}
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New(fmt.Sprintf("deleting cdn lock with cdn name %s : %v", cdn, err.Error())))
		return
	}

	alerts := tc.CreateAlerts(tc.SuccessLevel, "cdn lock deleted by admin")
	api.WriteAlertsObj(w, r, http.StatusOK, alerts, result)

	changeLogMsg := fmt.Sprintf("USER: %s, CDN: %s, ACTION: Lock Released By admin", result.UserName, cdn)
	api.CreateChangeLogRawTx(api.ApiChange, changeLogMsg, inf.User, tx)
}
