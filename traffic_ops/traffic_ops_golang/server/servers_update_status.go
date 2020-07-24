package server

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
	"net/http"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/config"
)

func GetServerUpdateStatusHandler(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"host_name"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	serverUpdateStatus, err := getServerUpdateStatus(inf.Tx.Tx, inf.Config, inf.Params["host_name"])
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
		return
	}
	api.WriteRespRaw(w, r, serverUpdateStatus)
}

func getServerUpdateStatus(tx *sql.Tx, cfg *config.Config, hostName string) ([]tc.ServerUpdateStatus, error) {

	updateStatuses := []tc.ServerUpdateStatus{}
	var rows *sql.Rows
	var err error

	selectQuery := `
WITH RECURSIVE topology_ancestors AS (
	SELECT tcp.child parent, tc.cachegroup
	FROM "server" s
	JOIN cachegroup c ON s.cachegroup = c.id
	JOIN topology_cachegroup tc ON c."name" = tc.cachegroup
	JOIN topology_cachegroup_parents tcp ON tc.id = tcp.child
	WHERE s.host_name = $1
UNION ALL
	SELECT tcp.parent, tc.cachegroup
	FROM topology_ancestors ta, topology_cachegroup_parents tcp
	JOIN topology_cachegroup tc ON tcp.parent = tc.id
	WHERE ta.parent = tcp.child
), server_topology_ancestors AS (
SELECT s.id, s.cachegroup, s.cdn_id, s.upd_pending, s.reval_pending, s.status
FROM server s
JOIN cachegroup c ON s.cachegroup = c.id
JOIN topology_ancestors ta ON c."name" = ta.cachegroup
WHERE s.host_name != $1
), parentservers AS (
	SELECT ps.id, ps.cachegroup, ps.cdn_id, ps.upd_pending, ps.reval_pending, ps.status
		FROM server ps
	LEFT JOIN status AS pstatus ON pstatus.id = ps.status
	WHERE pstatus.name != 'OFFLINE'
), use_reval_pending AS (
	SELECT value::BOOLEAN
	FROM parameter
	WHERE name = 'use_reval_pending'
	AND config_file = 'global'
	UNION ALL SELECT FALSE FETCH FIRST 1 ROW ONLY
)
SELECT
	s.id,
	s.host_name,
	type.name AS type,
	(s.reval_pending::BOOLEAN) AS server_reval_pending,
	use_reval_pending.value,
	s.upd_pending,
	status.name AS status,
	COALESCE(BOOL_OR(ps.upd_pending),
		TRUE IN (
			SELECT sta.upd_pending FROM  server_topology_ancestors sta)
		) AS parent_upd_pending,
	COALESCE(BOOL_OR(ps.reval_pending),
		TRUE IN (
			SELECT sta.reval_pending FROM  server_topology_ancestors sta)
		) AS parent_reval_pending
	FROM use_reval_pending,
		 server s
LEFT JOIN status ON s.status = status.id
LEFT JOIN cachegroup cg ON s.cachegroup = cg.id
LEFT JOIN type ON type.id = s.type
LEFT JOIN parentservers ps ON ps.cachegroup = cg.parent_cachegroup_id
	AND ps.cdn_id = s.cdn_id
	AND type.name = 'EDGE'
WHERE s.host_name = $1
GROUP BY s.id, s.host_name, type.name, server_reval_pending, use_reval_pending.value, s.upd_pending, status.name
ORDER BY s.id
` // remove the type.name = 'EDGE' condition if other server types should have their parents processed

	rows, err = tx.Query(selectQuery, hostName)
	if err != nil {
		log.Error.Printf("could not execute query: %s\n", err)
		return nil, tc.DBError
	}
	defer log.Close(rows, "unable to close db connection")

	for rows.Next() {
		var serverUpdateStatus tc.ServerUpdateStatus
		var serverType string
		if err := rows.Scan(&serverUpdateStatus.HostId, &serverUpdateStatus.HostName, &serverType, &serverUpdateStatus.RevalPending, &serverUpdateStatus.UseRevalPending, &serverUpdateStatus.UpdatePending, &serverUpdateStatus.Status, &serverUpdateStatus.ParentPending, &serverUpdateStatus.ParentRevalPending); err != nil {
			log.Error.Printf("could not scan server update status: %s\n", err)
			return nil, tc.DBError
		}
		updateStatuses = append(updateStatuses, serverUpdateStatus)
	}
	return updateStatuses, nil
}

func GetServerUpdateStatusHandlerV2(w http.ResponseWriter, r *http.Request) {
	GetServerUpdateStatusHandlerV1(w, r)
}

func GetServerUpdateStatusHandlerV1(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"host_name"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	serverUpdateStatus, err := getServerUpdateStatusV2(inf.Tx.Tx, inf.Config, inf.Params["host_name"])
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
		return
	}
	api.WriteRespRaw(w, r, serverUpdateStatus)
}

// getServerUpdateStatusV2 supports /servers/all/update_status in addition to /servers/{host_name}/update_status
// This special case is believed to be used nowhere.
func getServerUpdateStatusV2(tx *sql.Tx, cfg *config.Config, hostName string) ([]tc.ServerUpdateStatus, error) {
	// language=SQL
	baseSelectStatement := `
WITH parentservers AS (
	SELECT ps.id, ps.cachegroup, ps.cdn_id, ps.upd_pending, ps.reval_pending
	FROM server ps
	LEFT JOIN status AS pstatus ON pstatus.id = ps.status
	WHERE pstatus.name != 'OFFLINE'
), use_reval_pending AS (
	SELECT value::BOOLEAN
	FROM parameter
	WHERE name = 'use_reval_pending'
	AND config_file = 'global'
	UNION ALL SELECT FALSE FETCH FIRST 1 ROW ONLY
)
SELECT
	s.id,
	s.host_name,
	type.name AS type,
	(s.reval_pending::BOOLEAN) AS server_reval_pending,
	use_reval_pending.value,
	s.upd_pending,
	status.name AS status,
	COALESCE(BOOL_OR(ps.upd_pending), FALSE) AS parent_upd_pending,
	COALESCE(BOOL_OR(ps.reval_pending), FALSE) AS parent_reval_pending
	FROM use_reval_pending,
		 server s
LEFT JOIN status ON s.status = status.id
LEFT JOIN cachegroup cg ON s.cachegroup = cg.id
LEFT JOIN type ON type.id = s.type
LEFT JOIN parentservers ps ON ps.cachegroup = cg.parent_cachegroup_id
	AND ps.cdn_id = s.cdn_id
	AND type.name = 'EDGE'
` // remove the EDGE reference if other server types should have their parents processed

	// language=SQL
	groupBy := `
GROUP BY s.id, s.host_name, type.name, server_reval_pending, use_reval_pending.value, s.upd_pending, status.name
ORDER BY s.id
`

	updateStatuses := []tc.ServerUpdateStatus{}
	var rows *sql.Rows
	var err error
	if hostName == "all" {
		rows, err = tx.Query(baseSelectStatement + groupBy)
		if err != nil {
			log.Error.Printf("could not execute select server update status query: %s\n", err)
			return nil, tc.DBError
		}
	} else {
		rows, err = tx.Query(baseSelectStatement+` WHERE s.host_name = $1`+groupBy, hostName)
		if err != nil {
			log.Error.Printf("could not execute select server update status by hostname query: %s\n", err)
			return nil, tc.DBError
		}
	}
	defer rows.Close()

	for rows.Next() {
		var serverUpdateStatus tc.ServerUpdateStatus
		var serverType string
		if err := rows.Scan(&serverUpdateStatus.HostId, &serverUpdateStatus.HostName, &serverType, &serverUpdateStatus.RevalPending, &serverUpdateStatus.UseRevalPending, &serverUpdateStatus.UpdatePending, &serverUpdateStatus.Status, &serverUpdateStatus.ParentPending, &serverUpdateStatus.ParentRevalPending); err != nil {
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
