// Package cdn_lock contains the CRD methods which aid in locking and unlocking CDNs.
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
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"

	"github.com/lib/pq"
)

const readQuery = `SELECT username, cdn, message, soft, shared_usernames, last_updated FROM cdn_lock`
const insertQuery = `INSERT INTO cdn_lock (username, cdn, message, soft, shared_usernames) VALUES ($1, $2, $3, $4, $5) RETURNING username, cdn, message, soft, shared_usernames, last_updated`
const deleteQuery = `DELETE FROM cdn_lock WHERE cdn=$1 AND username=$2 RETURNING username, cdn, message, soft, shared_usernames, last_updated`
const deleteAdminQuery = `DELETE FROM cdn_lock WHERE cdn=$1 RETURNING username, cdn, message, soft, shared_usernames, last_updated`
const checkSharedUsersValidityQuery = `select count(*) from tm_user u join role r on r.id = u.role join role_capability rc on rc.role_id = r.id where u.tenant_id in (WITH RECURSIVE
user_tenant_id as (select (select tenant_id from tm_user where username=$1)::bigint as v),
user_tenant_parents AS (
  SELECT active, parent_id FROM tenant WHERE id = (select v from user_tenant_id)
  UNION
  SELECT t.active, t.parent_id FROM TENANT t JOIN user_tenant_parents ON user_tenant_parents.parent_id = t.id
),
user_tenant_active AS (
  SELECT bool_and(active) as v FROM user_tenant_parents
),
user_tenant_children AS (
  SELECT id, name, active, parent_id
  FROM tenant WHERE id = (SELECT v FROM user_tenant_id) AND (SELECT v FROM user_tenant_active)
  UNION
  SELECT t.id, t.name, t.active, t.parent_id
  FROM tenant t JOIN user_tenant_children ON user_tenant_children.id = t.parent_id
)
SELECT id as tenant_list FROM user_tenant_children) and u.username = ANY($2) and rc.cap_name='CDN-LOCK:CREATE'`

// Read is the handler for GET requests to /cdn_locks.
func Read(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	cols := map[string]dbhelpers.WhereColumnInfo{
		"cdn":      {Column: "cdn_lock.cdn", Checker: nil},
		"username": {Column: "cdn_lock.username", Checker: nil},
	}

	where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(inf.Params, cols)
	if len(errs) > 0 {
		errCode = http.StatusBadRequest
		userErr = util.JoinErrs(errs)
		api.HandleErr(w, r, tx, errCode, userErr, nil)
		return
	}

	cdnLock := []tc.CDNLock{}
	query := readQuery + where + orderBy + pagination
	rows, err := inf.Tx.NamedQuery(query, queryValues)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("querying cdn locks: "+err.Error()))
		return
	}
	defer rows.Close()

	for rows.Next() {
		var cLock tc.CDNLock
		if err = rows.Scan(&cLock.UserName, &cLock.CDN, &cLock.Message, &cLock.Soft, pq.Array(&cLock.SharedUserNames), &cLock.LastUpdated); err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("scanning cdn locks: "+err.Error()))
			return
		}
		cdnLock = append(cdnLock, cLock)
	}

	api.WriteResp(w, r, cdnLock)
}

// Create is the handler for POST requests to /cdn_locks.
func Create(w http.ResponseWriter, r *http.Request) {
	soft := "soft"
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	tx := inf.Tx.Tx
	var cdnLock tc.CDNLock
	if err := json.NewDecoder(r.Body).Decode(&cdnLock); err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
		return
	}
	if cdnLock.Soft == nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, errors.New("field 'soft' must be present"), nil)
		return
	}
	if !*cdnLock.Soft {
		soft = "hard"
	}
	if cdnLock.CDN == "" {
		api.HandleErr(w, r, tx, http.StatusBadRequest, errors.New("field 'cdn' must be present"), nil)
		return
	}
	cdnLock.UserName = inf.User.UserName
	if cdnLock.SharedUserNames != nil && len(cdnLock.SharedUserNames) > 0 {
		errCode, userErr, sysErr := checkSharedUserNamesValidity(tx, cdnLock)
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, tx, errCode, userErr, sysErr)
			return
		}
	}
	resultRows, err := inf.Tx.Query(insertQuery, cdnLock.UserName, cdnLock.CDN, cdnLock.Message, cdnLock.Soft, pq.Array(cdnLock.SharedUserNames))
	if err != nil {
		userErr, sysErr, errCode := api.ParseDBError(err)
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer resultRows.Close()

	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err := resultRows.Scan(&cdnLock.UserName, &cdnLock.CDN, &cdnLock.Message, &cdnLock.Soft, pq.Array(cdnLock.SharedUserNames), &cdnLock.LastUpdated); err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("cdn lock create: scanning locks: "+err.Error()))
			return
		}
	}
	if rowsAffected == 0 {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("cdn lock create: lock couldn't be acquired"))
		return
	}
	alerts := tc.CreateAlerts(tc.SuccessLevel, fmt.Sprintf("%s CDN lock acquired!", soft))
	api.WriteAlertsObj(w, r, http.StatusCreated, alerts, cdnLock)

	changeLogMsg := fmt.Sprintf("USER: %s, CDN: %s, ACTION: %s lock acquired", inf.User.UserName, cdnLock.CDN, soft)
	api.CreateChangeLogRawTx(api.ApiChange, changeLogMsg, inf.User, tx)
}

func checkSharedUserNamesValidity(tx *sql.Tx, lock tc.CDNLock) (int, error, error) {
	count := 0
	if err := tx.QueryRow(checkSharedUsersValidityQuery, lock.UserName, pq.Array(lock.SharedUserNames)).Scan(&count); err != nil {
		return http.StatusInternalServerError, nil, err
	}
	if count != len(lock.SharedUserNames) {
		return http.StatusBadRequest, errors.New("shared users must exist and be in the same/child tenants of the current user"), nil
	}
	return http.StatusOK, nil, nil
}

// Delete is the handler for DELETE requests to /cdn_locks.
func Delete(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"cdn"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	cdn := inf.Params["cdn"]
	tx := inf.Tx.Tx
	var result tc.CDNLock
	var err error
	adminPerms := inf.Config.RoleBasedPermissions && inf.User.Can("CDN-LOCK:DELETE-OTHERS")

	if adminPerms || inf.User.PrivLevel == auth.PrivLevelAdmin {
		err = inf.Tx.Tx.QueryRow(deleteAdminQuery, cdn).Scan(&result.UserName, &result.CDN, &result.Message, &result.Soft, pq.Array(&result.SharedUserNames), &result.LastUpdated)
	} else {
		err = inf.Tx.Tx.QueryRow(deleteQuery, cdn, inf.User.UserName).Scan(&result.UserName, &result.CDN, &result.Message, &result.Soft, pq.Array(&result.SharedUserNames), &result.LastUpdated)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if inf.User.PrivLevel != auth.PrivLevelAdmin {
				api.HandleErr(w, r, tx, http.StatusForbidden, fmt.Errorf("deleting cdn lock with cdn name %s: operation forbidden", cdn), nil)
				return
			}
			api.HandleErr(w, r, tx, http.StatusNotFound, fmt.Errorf("deleting cdn lock with cdn name %s: lock not found", cdn), nil)
			return
		}
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("deleting cdn lock with cdn name %s : %w", cdn, err))
		return
	}
	alerts := tc.CreateAlerts(tc.SuccessLevel, "cdn lock deleted")
	api.WriteAlertsObj(w, r, http.StatusOK, alerts, result)
	changeLogMsg := fmt.Sprintf("USER: %s, CDN: %s, ACTION: Lock Released", result.UserName, cdn)
	api.CreateChangeLogRawTx(api.ApiChange, changeLogMsg, inf.User, tx)
}
