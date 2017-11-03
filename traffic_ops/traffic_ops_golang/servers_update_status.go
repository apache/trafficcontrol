package main

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
	"fmt"
	"net/http"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/jmoiron/sqlx"
)

func getServerUpdateStatusHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErr := tc.GetHandleErrorFunc(w, r)

		// p PathParams, username string, privLevel int
		ctx := r.Context()
		pathParams, err := getPathParams(ctx)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}
		hostName := pathParams["host_name"]

		serverUpdateStatus, err := getServerUpdateStatus(hostName, db)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}

		respBts, err := json.Marshal(serverUpdateStatus)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}

		w.Header().Set(tc.ContentType, tc.ApplicationJson)
		fmt.Fprintf(w, "%s", respBts)
	}
}

func getServerUpdateStatus(hostName string, db *sqlx.DB) ([]tc.ServerUpdateStatus, error) {
	baseSelectStatement :=
		`WITH parentservers AS (SELECT ps.id, ps.cachegroup, ps.cdn_id, ps.upd_pending, ps.reval_pending FROM server ps
         LEFT JOIN status AS pstatus ON pstatus.id = ps.status
         WHERE pstatus.name != 'OFFLINE' ),
         use_reval AS (SELECT value::boolean FROM parameter WHERE name = 'use_reval_pending' AND config_file = 'global')
         SELECT s.id, s.host_name, type.name AS type, (s.reval_pending::boolean AND use_reval.value) as combined_reval_pending, s.upd_pending, status.name AS status, COALESCE(bool_or(ps.upd_pending), FALSE) AS parent_upd_pending, COALESCE(bool_or(ps.reval_pending), FALSE) AS parent_reval_pending FROM use_reval, server s
         LEFT JOIN status ON s.status = status.id
         LEFT JOIN cachegroup cg ON s.cachegroup = cg.id
         LEFT JOIN type ON type.id = s.type
         LEFT JOIN parentservers ps ON ps.cachegroup = cg.parent_cachegroup_id AND ps.cdn_id = s.cdn_id AND type.name = 'EDGE'` //remove the EDGE reference if other server types should have their parents processed

	groupBy := ` GROUP BY s.id, s.host_name, type.name, combined_reval_pending, s.upd_pending, status.name ORDER BY s.id;`

	updateStatuses := []tc.ServerUpdateStatus{}
	var rows *sql.Rows
	if hostName == "all" {
		selectAll, err := db.Prepare(baseSelectStatement + groupBy)
		if err != nil {
			log.Error.Printf("could not prepare select server update status query: %s\n", err)
			return nil, tc.DBError
		}
		rows, err = selectAll.Query()
		if err != nil {
			log.Error.Printf("could not execute select server update status query: %s\n", err)
			return nil, tc.DBError
		}
	} else {
		selectStatement, err := db.Prepare(baseSelectStatement + ` WHERE s.host_name = $1` + groupBy)
		if err != nil {
			log.Error.Printf("could not prepare select server update status by hostname query: %s\n", err)
			return nil, tc.DBError
		}
		rows, err = selectStatement.Query(hostName)
		if err != nil {
			log.Error.Printf("could not execute select server update status by hostname query: %s\n", err)
			return nil, tc.DBError
		}
	}
        defer rows.Close()
	for rows.Next() {
		var serverUpdateStatus tc.ServerUpdateStatus
		var serverType string
		if err := rows.Scan(&serverUpdateStatus.HostId, &serverUpdateStatus.HostName, &serverType, &serverUpdateStatus.RevalPending, &serverUpdateStatus.UpdatePending, &serverUpdateStatus.Status, &serverUpdateStatus.ParentPending, &serverUpdateStatus.ParentRevalPending); err != nil {
			log.Error.Printf("could not scan server update status: %s\n", err)
			return nil, tc.DBError
		}
		if hostName == "all" { //if we want to return the parent data for servers when all is used remove this block
			serverUpdateStatus.ParentRevalPending = false
			serverUpdateStatus.ParentPending = false
		}
		updateStatuses = append(updateStatuses, serverUpdateStatus)
	}
	return updateStatuses, nil
}
