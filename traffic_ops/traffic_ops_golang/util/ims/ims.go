package ims

import (
	"database/sql"
	"github.com/apache/trafficcontrol/grove/web"
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/jmoiron/sqlx"
	"net/http"
	"time"
)

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

// LatestTimestamp to keep track of the max of "last updated" times in tables
type LatestTimestamp struct {
	LatestTime *tc.TimeNoMod `json:"latestTime" db:"max"`
}

// TryIfModifiedSinceQuery for components that DO NOT implement the CRUDER interface
func TryIfModifiedSinceQuery(tx *sqlx.Tx, h http.Header, queryValues map[string]interface{}, query string) (bool, time.Time) {
	var maxTime time.Time
	ims := []string{}
	runSecond := true
	if h == nil {
		return runSecond, maxTime
	}
	ims = h[rfc.IfModifiedSince]
	if ims == nil || len(ims) == 0 {
		return runSecond, maxTime
	}
	if l, ok := web.ParseHTTPDate(ims[0]); !ok {
		return runSecond, maxTime
	} else {
		rows, err := tx.NamedQuery(query, queryValues)
		if rows != nil {
			defer rows.Close()
		}
		if err != nil {
			log.Warnf("Couldn't get the max last updated time: %v", err)
			return runSecond, maxTime
		}
		if err == sql.ErrNoRows {
			runSecond = false
			return runSecond, maxTime
		}
		// This should only ever contain one row
		if rows.Next() {
			v := &LatestTimestamp{}
			if err = rows.StructScan(v); err != nil || v == nil {
				log.Warnf("Failed to parse the max time stamp into a struct %v", err)
				return runSecond, maxTime
			}
			maxTime = v.LatestTime.Time
			// The request IMS time is later than the max of (lastUpdated, deleted_time)
			if v.LatestTime != nil && l.After(v.LatestTime.Time) {
				runSecond = false
				return runSecond, maxTime
			}
		} else {
			runSecond = false
		}
	}
	return runSecond, maxTime
}
