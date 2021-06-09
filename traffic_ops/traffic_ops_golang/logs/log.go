package logs

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
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
)

const DefaultLogLimit = 1000
const DefaultLogLimitForDays = 1000000
const DefaultLogDays = 30

func GetDeprecated(w http.ResponseWriter, r *http.Request) {
	get(w, r, api.CreateDeprecationAlerts(util.StrPtr("/logs")))
}

func Get(w http.ResponseWriter, r *http.Request) {
	get(w, r, tc.Alerts{})
}

func get(w http.ResponseWriter, r *http.Request, a tc.Alerts) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, []string{"days", "limit"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	limit := DefaultLogLimit
	days := DefaultLogDays
	if pDays, ok := inf.IntParams["days"]; ok {
		days = pDays
		limit = DefaultLogLimitForDays
	}
	if pLimit, ok := inf.IntParams["limit"]; ok {
		limit = pLimit
	}

	setLastSeenCookie(w)
	logs, count, err := getLog(inf, days, limit)
	if err != nil {
		a.AddNewAlert(tc.ErrorLevel, err.Error())
		api.WriteAlerts(w, r, http.StatusInternalServerError, a)
		return
	}
	if a.HasAlerts() {
		api.WriteAlertsObj(w, r, 200, a, logs)
	} else {
		api.WriteRespWithSummary(w, r, logs, count)
	}
}

func GetNewCount(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, []string{"days", "limit"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	lastSeen, ok := getLastSeenCookie(r)
	if !ok {
		setLastSeenCookie(w) // only set the cookie if it didn't exist; emulates old Perl behavior
		api.WriteResp(w, r, tc.NewLogCountResp{NewLogCount: 0})
		return
	}
	newCount, err := getLogCountSince(inf.Tx.Tx, lastSeen)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting log new count: "+err.Error()))
		return
	}
	api.WriteResp(w, r, tc.NewLogCountResp{NewLogCount: newCount})
}

const LastSeenLogCookieName = "last_seen_log"

func setLastSeenCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:    LastSeenLogCookieName,
		Value:   time.Now().Format(time.RFC3339Nano),
		Expires: time.Now().Add(time.Hour * 24 * 7),
		Path:    "/",
	})
}

// getLastSeenCookie returns the time of the last log seen cookie, or whether the cookie didn't exist or there was an error parsing it. Errors are not logged or returned, only that the cookie could not be retrieved.
func getLastSeenCookie(r *http.Request) (time.Time, bool) {
	cookie, err := r.Cookie(LastSeenLogCookieName)
	if err != nil {
		return time.Time{}, false
	}
	lastSeen, err := time.Parse(time.RFC3339Nano, cookie.Value)
	if err != nil {
		return time.Time{}, false
	}
	return lastSeen, true
}

const selectFromQuery = `
SELECT l.id, l.level, l.message, u.username as user, l.ticketnum, l.last_updated
FROM "log" as l JOIN tm_user as u ON l.tm_user = u.id`

func getLog(inf *api.APIInfo, days int, limit int) ([]tc.Log, uint64, error) {
	var count = uint64(0)

	queryParamsToQueryCols := map[string]dbhelpers.WhereColumnInfo{
		"username": dbhelpers.WhereColumnInfo{Column: "u.username", Checker: nil},
	}
	where, orderBy, pagination, queryValues, errs :=
		dbhelpers.BuildWhereAndOrderByAndPagination(inf.Params, queryParamsToQueryCols)
	if len(errs) > 0 {
		return nil, 0, util.JoinErrs(errs)
	}

	timeInterval := fmt.Sprintf("l.last_updated > now() - INTERVAL '%v' DAY", days)
	if _, ok := inf.Params["username"]; ok {
		where = where + " AND " + timeInterval
	} else {
		where = "\nWHERE " + timeInterval
	}
	if orderBy == "" {
		orderBy = orderBy + "\nORDER BY l.last_updated DESC"
	}
	if pagination == "" {
		pagination = pagination + fmt.Sprintf("\nLIMIT %v", limit)
	}
	query := selectFromQuery + where + orderBy + pagination

	rows, err := inf.Tx.NamedQuery(query, queryValues)
	if err != nil {
		return nil, count, errors.New("querying logs: " + err.Error())
	}
	ls := []tc.Log{}
	for rows.Next() {
		l := tc.Log{}
		if err = rows.Scan(&l.ID, &l.Level, &l.Message, &l.User, &l.TicketNum, &l.LastUpdated); err != nil {
			return nil, count, errors.New("scanning logs: " + err.Error())
		}
		count += 1
		ls = append(ls, l)
	}
	return ls, count, nil
}

func getLogCountSince(tx *sql.Tx, since time.Time) (uint64, error) {
	count := uint64(0)
	if err := tx.QueryRow(`SELECT count(*) from log where last_updated > $1`, since).Scan(&count); err != nil {
		return 0, errors.New("querying log last seen count: " + err.Error())
	}
	return count, nil
}
