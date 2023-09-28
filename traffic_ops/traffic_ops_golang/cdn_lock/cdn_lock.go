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
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/dbhelpers"

	"github.com/lib/pq"
)

const readQuery = `SELECT c.username, c.cdn, c.message, c.soft, c.last_updated, ARRAY_REMOVE(ARRAY_AGG(DISTINCT(u.username)), null) AS shared_usernames FROM cdn_lock_user u FULL JOIN cdn_lock c ON c.username = u.owner AND c.cdn = u.cdn`
const insertQueryWithoutSharedUserNames = `INSERT INTO cdn_lock (username, cdn, message, soft) VALUES ($1, $2, $3, $4) RETURNING username, cdn, message, soft, last_updated`

const insertQueryWithSharedUserNames = `WITH first_insert AS (
INSERT INTO cdn_lock (username, cdn, message, soft)
VALUES($1, $2, $3, $4)
RETURNING *
),
second_insert AS (
INSERT INTO cdn_lock_user (owner, cdn, username)
VALUES($5, $6, UNNEST($7::TEXT[]))
RETURNING owner, username, cdn)
SELECT f.username, f.cdn, f.message, f.soft, ARRAY_AGG(s.username) AS shared_usernames, f.last_updated
FROM first_insert f
JOIN second_insert s
ON s.owner = f.username
AND s.cdn = f.cdn
GROUP BY f.username,
f.cdn,
f.message,
f.soft,
f.last_updated`

const deleteQuery = `DELETE FROM cdn_lock WHERE cdn=$1 AND username=$2 RETURNING username, cdn, message, soft, (SELECT ARRAY_AGG(u.username) AS shared_usernames FROM cdn_lock_user u JOIN cdn_lock c ON c.username = u.owner AND c.cdn = u.cdn WHERE u.cdn=$1 AND u.owner=$2), last_updated`

const deleteAdminQuery = `DELETE FROM cdn_lock WHERE cdn=$1 RETURNING username, cdn, message, soft, (SELECT ARRAY_AGG(u.username) AS shared_usernames FROM cdn_lock_user u JOIN cdn_lock c ON c.username = u.owner AND c.cdn = u.cdn WHERE u.cdn=$1), last_updated`

const checkSharedUsersValidityQuery = `select count(*) from tm_user u join role r on r.id = u.role join role_capability rc on rc.role_id = r.id where u.username = ANY($1) and (rc.cap_name='ALL' or rc.cap_name='CDN-LOCK:CREATE')`

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
		"cdn":      {Column: "c.cdn", Checker: nil},
		"username": {Column: "c.username", Checker: nil},
	}

	where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(inf.Params, cols)
	if len(errs) > 0 {
		errCode = http.StatusBadRequest
		userErr = util.JoinErrs(errs)
		api.HandleErr(w, r, tx, errCode, userErr, nil)
		return
	}

	cdnLock := []tc.CDNLock{}
	query := readQuery + where + orderBy + pagination + " GROUP BY c.cdn"
	rows, err := inf.Tx.NamedQuery(query, queryValues)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("querying cdn locks: "+err.Error()))
		return
	}
	defer rows.Close()

	for rows.Next() {
		var cLock tc.CDNLock
		if err = rows.Scan(&cLock.UserName, &cLock.CDN, &cLock.Message, &cLock.Soft, &cLock.LastUpdated, pq.Array(&cLock.SharedUserNames)); err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("scanning cdn locks: "+err.Error()))
			return
		}
		if inf.Version != nil && inf.Version.GreaterThanOrEqualTo(&api.Version{Major: 5, Minor: 0}) {
			t, err := util.ConvertTimeFormat(cLock.LastUpdated, time.RFC3339)
			if err != nil {
				api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("converting time formats: "+err.Error()))
				return
			}
			cLock.LastUpdated = *t
		}
		cdnLock = append(cdnLock, cLock)
	}

	api.WriteResp(w, r, cdnLock)
}

// Create is the handler for POST requests to /cdn_locks.
func Create(w http.ResponseWriter, r *http.Request) {
	var err error
	var shared bool
	var resultRows *sql.Rows
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
	if len(cdnLock.SharedUserNames) == 0 {
		resultRows, err = inf.Tx.Query(insertQueryWithoutSharedUserNames, cdnLock.UserName, cdnLock.CDN, cdnLock.Message, cdnLock.Soft)
	} else {
		shared = true
		for _, sharedUser := range cdnLock.SharedUserNames {
			if sharedUser == cdnLock.UserName {
				api.HandleErr(w, r, tx, http.StatusBadRequest, errors.New("shared username cannot be the same as the one creating the lock"), nil)
				return
			}
		}
		resultRows, err = inf.Tx.Query(insertQueryWithSharedUserNames, cdnLock.UserName, cdnLock.CDN, cdnLock.Message, cdnLock.Soft, cdnLock.UserName, cdnLock.CDN, pq.Array(cdnLock.SharedUserNames))
	}
	if err != nil {
		userErr, sysErr, errCode := api.ParseDBError(err)
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer resultRows.Close()

	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if shared {
			if err := resultRows.Scan(&cdnLock.UserName, &cdnLock.CDN, &cdnLock.Message, &cdnLock.Soft, pq.Array(cdnLock.SharedUserNames), &cdnLock.LastUpdated); err != nil {
				api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("cdn lock create: scanning locks: "+err.Error()))
				return
			}
		} else {
			if err := resultRows.Scan(&cdnLock.UserName, &cdnLock.CDN, &cdnLock.Message, &cdnLock.Soft, &cdnLock.LastUpdated); err != nil {
				api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("cdn lock create: scanning locks: "+err.Error()))
				return
			}
		}
	}
	if rowsAffected == 0 {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("cdn lock create: lock couldn't be acquired"))
		return
	}
	if inf.Version != nil && inf.Version.Major >= 5 && inf.Version.Minor >= 0 {
		t, err := util.ConvertTimeFormat(cdnLock.LastUpdated, time.RFC3339)
		if err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("converting time formats: "+err.Error()))
			return
		}
		cdnLock.LastUpdated = *t
	}
	alerts := tc.CreateAlerts(tc.SuccessLevel, fmt.Sprintf("%s CDN lock acquired!", soft))
	api.WriteAlertsObj(w, r, http.StatusCreated, alerts, cdnLock)

	changeLogMsg := fmt.Sprintf("USER: %s, CDN: %s, ACTION: %s lock acquired", inf.User.UserName, cdnLock.CDN, soft)
	api.CreateChangeLogRawTx(api.ApiChange, changeLogMsg, inf.User, tx)
}

func checkSharedUserNamesValidity(tx *sql.Tx, lock tc.CDNLock) (int, error, error) {
	count := 0
	if err := tx.QueryRow(checkSharedUsersValidityQuery, pq.Array(lock.SharedUserNames)).Scan(&count); err != nil {
		return http.StatusInternalServerError, nil, err
	}
	if count != len(lock.SharedUserNames) {
		return http.StatusBadRequest, errors.New("shared users must exist and have the correct permissions to create a lock"), nil
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
	var adminPerms bool

	if (inf.Version.GreaterThanOrEqualTo(&api.Version{Major: 4}) && inf.Config.RoleBasedPermissions) || inf.Version.GreaterThanOrEqualTo(&api.Version{Major: 5}) {
		adminPerms = inf.User.Can(tc.PermCDNLocksDeleteOthers)
	} else {
		adminPerms = inf.User.PrivLevel == auth.PrivLevelAdmin
	}
	if adminPerms {
		err = inf.Tx.Tx.QueryRow(deleteAdminQuery, cdn).Scan(&result.UserName, &result.CDN, &result.Message, &result.Soft, pq.Array(&result.SharedUserNames), &result.LastUpdated)
	} else {
		err = inf.Tx.Tx.QueryRow(deleteQuery, cdn, inf.User.UserName).Scan(&result.UserName, &result.CDN, &result.Message, &result.Soft, pq.Array(&result.SharedUserNames), &result.LastUpdated)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if !adminPerms {
				api.HandleErr(w, r, tx, http.StatusForbidden, fmt.Errorf("deleting cdn lock with cdn name %s: operation forbidden", cdn), nil)
				return
			}
			api.HandleErr(w, r, tx, http.StatusNotFound, fmt.Errorf("deleting cdn lock with cdn name %s: lock not found", cdn), nil)
			return
		}
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("deleting cdn lock with cdn name %s : %w", cdn, err))
		return
	}
	if inf.Version != nil && inf.Version.Major >= 5 && inf.Version.Minor >= 0 {
		t, err := util.ConvertTimeFormat(result.LastUpdated, time.RFC3339)
		if err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("converting time formats: "+err.Error()))
			return
		}
		result.LastUpdated = *t
	}
	alerts := tc.CreateAlerts(tc.SuccessLevel, "cdn lock deleted")
	api.WriteAlertsObj(w, r, http.StatusOK, alerts, result)
	changeLogMsg := fmt.Sprintf("USER: %s, CDN: %s, ACTION: Lock Released", result.UserName, cdn)
	api.CreateChangeLogRawTx(api.ApiChange, changeLogMsg, inf.User, tx)
}
