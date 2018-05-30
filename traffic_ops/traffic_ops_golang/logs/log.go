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
	"net/http"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
)

const DefaultLogLimit = 1000
const DefaultLogLimitForDays = 1000000
const DefaultLogDays = 30

func Get(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		api.RespWriter(w, r, inf.Tx.Tx)(getLog(inf.Tx.Tx, days, limit))
	}
}

func getLog(tx *sql.Tx, days int, limit int) ([]tc.Log, error) {
	rows, err := tx.Query(`
SELECT l.id, l.level, l.message, u.username as user, l.ticketnum, l.last_updated
FROM "log" as l JOIN tm_user as u ON l.tm_user = u.id
WHERE l.last_updated > now() - ($1 || ' DAY')::INTERVAL
ORDER BY l.last_updated DESC
LIMIT $2
`, days, limit)
	if err != nil {
		return nil, errors.New("querying logs: " + err.Error())
	}
	ls := []tc.Log{}
	for rows.Next() {
		l := tc.Log{}
		if err = rows.Scan(&l.ID, &l.Level, &l.Message, &l.User, &l.TicketNum, &l.LastUpdated); err != nil {
			return nil, errors.New("scanning logs: " + err.Error())
		}
		ls = append(ls, l)
	}
	return ls, nil
}
