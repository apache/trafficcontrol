// Package logs contains handlers and logic for the /logs and /logs/newcount
// API endpoints.
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
	"strconv"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/dbhelpers"
)

// These are the default values for query string parameters if not provided.
const (
	// For the 'limit' parameter.
	DefaultLogLimit = 1000
	// For the 'limit' parameter when 'days' is given and 'limit' is not.
	DefaultLogLimitForDays = 1000000
	// For the 'days' parameter.
	DefaultLogDays = 30
)

// Get is the handler for GET requests to /logs.
func Get(w http.ResponseWriter, r *http.Request) {
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
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, err, err)
		return
	}
	api.WriteRespWithSummary(w, r, logs, count)
}

// Get is the handler for GET requests to /logs V4.0.
func Getv40(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, []string{"days", "limit"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	days := DefaultLogDays
	if pDays, ok := inf.IntParams["days"]; ok {
		days = pDays
	}

	a := tc.Alerts{}
	setLastSeenCookie(w)
	logs, count, err := getLogV40(inf, days)

	if err != nil {
		a.AddNewAlert(tc.ErrorLevel, err.Error())
		api.WriteAlerts(w, r, http.StatusInternalServerError, a)
		return
	}
	var result interface{}
	var logsV5 []tc.LogV5
	if inf.Version.GreaterThanOrEqualTo(&api.Version{
		Major: 5,
		Minor: 0,
	}) {
		for _, l := range logs {
			logV5 := l.Upgrade()
			logsV5 = append(logsV5, logV5)
		}
		result = logsV5
	} else {
		result = logs
	}
	if a.HasAlerts() {
		api.WriteAlertsObj(w, r, 200, a, result)
	} else {
		api.WriteRespWithSummary(w, r, result, count)
	}
}

// GetNewCount is the handler for GET requests to /logs/newcount.
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

// LastSeenLogCookieName is the name of the HTTP cookie that stores the
// date/time at which the client last requested logs, so that unread logs can
// be returned to them on request.
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

const countQuery = `SELECT count(l.tm_user) FROM log as l`

func getLogV40(inf *api.Info, days int) ([]tc.Log, uint64, error) {
	var count = uint64(0)
	var whereCount string

	queryParamsToQueryCols := map[string]dbhelpers.WhereColumnInfo{
		"username": {Column: "u.username", Checker: nil},
	}
	where, _, pagination, queryValues, errs :=
		dbhelpers.BuildWhereAndOrderByAndPagination(inf.Params, queryParamsToQueryCols)

	if len(errs) > 0 {
		return nil, 0, util.JoinErrs(errs)
	}

	timeInterval := fmt.Sprintf("l.last_updated > now() - INTERVAL '%d' DAY", days)
	if where != "" {
		whereCount = ", tm_user as u\n" + where + " AND l.tm_user = u.id"
		where = where + " AND " + timeInterval
	} else {
		whereCount = where
		where = "\nWHERE " + timeInterval
	}
	queryCount := countQuery + whereCount
	rowCount, err := inf.Tx.NamedQuery(queryCount, queryValues)
	if err != nil {
		return nil, count, fmt.Errorf("querying log count for a given user: %w", err)
	}
	defer rowCount.Close()
	for rowCount.Next() {
		if err = rowCount.Scan(&count); err != nil {
			return nil, count, fmt.Errorf("scanning logs: %w", err)
		}
	}

	query := selectFromQuery + where + "\n ORDER BY last_updated DESC" + pagination
	rows, err := inf.Tx.NamedQuery(query, queryValues)
	if err != nil {
		return nil, count, fmt.Errorf("querying logs: %w", err)
	}
	defer rows.Close()
	ls := []tc.Log{}
	for rows.Next() {
		l := tc.Log{}
		if err = rows.Scan(&l.ID, &l.Level, &l.Message, &l.User, &l.TicketNum, &l.LastUpdated); err != nil {
			return nil, count, fmt.Errorf("scanning logs: %w", err)
		}
		ls = append(ls, l)
	}
	return ls, count, nil
}

func getLog(inf *api.Info, days int, limit int) ([]tc.Log, uint64, error) {
	var count = uint64(0)
	var whereCount string
	if _, ok := inf.Params["limit"]; !ok {
		if _, ok := inf.Params["days"]; !ok {
			inf.Params["limit"] = strconv.Itoa(DefaultLogLimit)
		}
	} else {
		inf.Params["limit"] = strconv.Itoa(limit)
	}

	queryParamsToQueryCols := map[string]dbhelpers.WhereColumnInfo{
		"username": {Column: "u.username", Checker: nil},
	}
	where, _, pagination, queryValues, errs :=
		dbhelpers.BuildWhereAndOrderByAndPagination(inf.Params, queryParamsToQueryCols)
	if len(errs) > 0 {
		return nil, 0, util.JoinErrs(errs)
	}

	timeInterval := fmt.Sprintf("l.last_updated > now() - INTERVAL '%v' DAY", days)
	if where != "" {
		whereCount = ", tm_user as u\n" + where + " AND l.tm_user = u.id"
		where = where + " AND " + timeInterval
	} else {
		whereCount = where
		where = "\nWHERE " + timeInterval
	}

	queryCount := countQuery + whereCount
	rowCount, err := inf.Tx.NamedQuery(queryCount, queryValues)
	if err != nil {
		return nil, count, errors.New("querying log count for a given user: " + err.Error())
	}
	defer rowCount.Close()
	for rowCount.Next() {
		if err = rowCount.Scan(&count); err != nil {
			return nil, count, errors.New("scanning logs: " + err.Error())
		}
	}

	query := selectFromQuery + where + "\n ORDER BY last_updated DESC" + pagination
	rows, err := inf.Tx.NamedQuery(query, queryValues)
	if err != nil {
		return nil, count, errors.New("querying logs: " + err.Error())
	}
	defer rows.Close()
	ls := []tc.Log{}
	for rows.Next() {
		l := tc.Log{}
		if err = rows.Scan(&l.ID, &l.Level, &l.Message, &l.User, &l.TicketNum, &l.LastUpdated); err != nil {
			return nil, count, errors.New("scanning logs: " + err.Error())
		}
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
