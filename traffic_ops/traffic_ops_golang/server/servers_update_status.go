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
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"

	"github.com/lib/pq"
)

func GetServerUpdateStatusHandler(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"host_name"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	serverUpdateStatuses, err, _ := getServerUpdateStatus(inf.Tx.Tx, inf.Params["host_name"])
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
		return
	}
	if inf.Version.LessThan(&api.Version{Major: 5}) {
		downgradedStatusesV4 := make([]tc.ServerUpdateStatusV40, len(serverUpdateStatuses))
		for i, status := range serverUpdateStatuses {
			downgradedStatusesV4[i] = status.Downgrade()
		}
		if inf.Version.LessThan(&api.Version{Major: 4}) {
			downgradedStatuses := make([]tc.ServerUpdateStatus, len(downgradedStatusesV4))
			for i, status := range downgradedStatusesV4 {
				downgradedStatuses[i] = status.Downgrade()
			}
			api.WriteRespRaw(w, r, downgradedStatuses)
		} else {
			api.WriteResp(w, r, downgradedStatusesV4)
		}
	} else {
		api.WriteResp(w, r, serverUpdateStatuses)
	}
}

func getServerUpdateStatus(tx *sql.Tx, hostName string) ([]tc.ServerUpdateStatusV5, error, error) {
	if serverUpdateStatusCacheIsInitialized() {
		return getServerUpdateStatusFromCache(hostName), nil, nil
	}

	updateStatuses := []tc.ServerUpdateStatusV5{}

	selectQuery := `
/* topology_ancestors finds the ancestor topology nodes of the topology node for
 * the cachegroup containing server $5.
 */
WITH RECURSIVE topology_ancestors AS (
/* This is the base case of the recursive CTE, the topology node for the
 * cachegroup containing server $5.
 */
	SELECT tcp.child parent, NULL cachegroup, s.id base_server_id
	FROM "server" s
	JOIN cachegroup c ON s.cachegroup = c.id
	JOIN topology_cachegroup tc ON c."name" = tc.cachegroup
	JOIN topology_cachegroup_parents tcp ON tc.id = tcp.child
	WHERE s.host_name = $5
UNION ALL
/* Find all direct topology parent nodes tc of a given topology ancestor ta. */
	SELECT tcp.parent, tc.cachegroup, ta.base_server_id
	FROM topology_ancestors ta, topology_cachegroup_parents tcp
	JOIN topology_cachegroup tc ON tcp.parent = tc.id
	JOIN cachegroup c ON tc.cachegroup = c."name"
	JOIN "type" t ON c."type" = t.id
	WHERE ta.parent = tcp.child
	AND t."name" LIKE ANY($4::TEXT[])
/* server_topology_ancestors is the set of every server whose cachegroup is an
 * ancestor topology node found by topology_ancestors.
 */
), server_topology_ancestors AS (
SELECT s.id,
	s.cachegroup,
	s.cdn_id,
	s.config_update_time > s.config_apply_time AS upd_pending,
	s.revalidate_update_time > s.revalidate_apply_time AS reval_pending,
	s.status,
	ta.base_server_id
	FROM server s
	JOIN cachegroup c ON s.cachegroup = c.id
	JOIN topology_ancestors ta ON c."name" = ta.cachegroup
	JOIN status ON status.id = s.status
	WHERE status.name = ANY($1::TEXT[])
), parentservers AS (
SELECT ps.id,
	ps.cachegroup,
	ps.cdn_id,
	ps.config_update_time > ps.config_apply_time AS upd_pending,
	ps.revalidate_update_time > ps.revalidate_apply_time AS reval_pending,
	ps.status
	FROM server ps
	LEFT JOIN status AS pstatus ON pstatus.id = ps.status
	LEFT JOIN type t ON ps."type" = t.id
	WHERE pstatus.name = ANY($1::TEXT[])
	AND t."name" LIKE ANY($4::TEXT[])
), use_reval_pending AS (
SELECT value::BOOLEAN
	FROM parameter
	WHERE name = $2
	AND config_file = $3
	UNION ALL SELECT FALSE FETCH FIRST 1 ROW ONLY
)
SELECT
	s.id,
	s.host_name,
	type.name AS type,
	s.revalidate_update_time > s.revalidate_apply_time AS server_reval_pending,
	use_reval_pending.value,
	s.config_update_time > s.config_apply_time AS server_upd_pending,
	status.name AS status,
	/* True if the cachegroup parent or any ancestor topology node has pending updates. */
	TRUE IN (
		SELECT sta.upd_pending FROM server_topology_ancestors sta
		WHERE sta.base_server_id = s.id
		AND sta.cdn_id = s.cdn_id
		UNION SELECT COALESCE(BOOL_OR(ps.upd_pending), FALSE)
	) AS parent_upd_pending,
	/* True if the cachegroup parent or any ancestor topology node has pending revalidation. */
	TRUE IN (
		SELECT sta.reval_pending FROM server_topology_ancestors sta
		WHERE sta.base_server_id = s.id
		AND sta.cdn_id = s.cdn_id
		UNION SELECT COALESCE(BOOL_OR(ps.reval_pending), FALSE)
	) AS parent_reval_pending,
	s.config_update_time,
	s.config_apply_time,
	s.config_update_failed,
	s.revalidate_update_time,
	s.revalidate_apply_time,
	s.revalidate_update_failed
	FROM use_reval_pending,
		 server s
LEFT JOIN status ON s.status = status.id
LEFT JOIN cachegroup cg ON s.cachegroup = cg.id
LEFT JOIN type ON type.id = s.type
LEFT JOIN parentservers ps ON ps.cachegroup = cg.parent_cachegroup_id
	AND ps.cdn_id = s.cdn_id
WHERE s.host_name = $5
GROUP BY s.id, s.host_name, type.name, server_reval_pending, use_reval_pending.value, server_upd_pending, status.name, config_update_time, config_apply_time, revalidate_update_time, revalidate_apply_time
ORDER BY s.id
`

	cacheStatusesToCheck := []tc.CacheStatus{tc.CacheStatusOnline, tc.CacheStatusReported, tc.CacheStatusAdminDown}
	cacheGroupTypes := []string{tc.EdgeTypePrefix + "%", tc.MidTypePrefix + "%"}
	rows, err := tx.Query(selectQuery, pq.Array(cacheStatusesToCheck), tc.UseRevalPendingParameterName, tc.GlobalConfigFileName, pq.Array(cacheGroupTypes), hostName)
	if err != nil {
		log.Errorf("could not execute query: %s\n", err)
		return nil, nil, fmt.Errorf("could not execute query: %w", err)
	}
	defer log.Close(rows, "getServerUpdateStatus(): unable to close db connection")

	for rows.Next() {
		var us tc.ServerUpdateStatusV5
		var serverType string
		if err := rows.Scan(&us.HostId, &us.HostName, &serverType, &us.RevalPending, &us.UseRevalPending, &us.UpdatePending, &us.Status, &us.ParentPending, &us.ParentRevalPending, &us.ConfigUpdateTime, &us.ConfigApplyTime, &us.ConfigUpdateFailed, &us.RevalidateUpdateTime, &us.RevalidateApplyTime, &us.RevalidateUpdateFailed); err != nil {
			return nil, nil, fmt.Errorf("could not scan server update status: %w", err)
		}
		updateStatuses = append(updateStatuses, us)
	}
	return updateStatuses, nil, nil
}

type serverUpdateStatuses struct {
	serverMap map[string][]tc.ServerUpdateStatusV5
	*sync.RWMutex
	initialized bool
	enabled     bool // note: enabled is only written to once at startup, before serving requests, so it doesn't need synchronized access
}

var serverUpdateStatusCache = serverUpdateStatuses{RWMutex: &sync.RWMutex{}}

func serverUpdateStatusCacheIsInitialized() bool {
	if serverUpdateStatusCache.enabled {
		serverUpdateStatusCache.RLock()
		defer serverUpdateStatusCache.RUnlock()
		return serverUpdateStatusCache.initialized
	}
	return false
}

func getServerUpdateStatusFromCache(hostname string) []tc.ServerUpdateStatusV5 {
	serverUpdateStatusCache.RLock()
	defer serverUpdateStatusCache.RUnlock()
	return serverUpdateStatusCache.serverMap[hostname]
}

var once = sync.Once{}

func InitServerUpdateStatusCache(interval time.Duration, db *sql.DB, timeout time.Duration) {
	once.Do(func() {
		if interval <= 0 {
			return
		}
		serverUpdateStatusCache.enabled = true
		refreshServerUpdateStatusCache(db, timeout)
		startServerUpdateStatusCacheRefresher(interval, db, timeout)
	})
}

func startServerUpdateStatusCacheRefresher(interval time.Duration, db *sql.DB, timeout time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			refreshServerUpdateStatusCache(db, timeout)
		}
	}()
}

func refreshServerUpdateStatusCache(db *sql.DB, timeout time.Duration) {
	newServerUpdateStatuses, err := getServerUpdateStatuses(db, timeout)
	if err != nil {
		log.Errorf("refreshing server update status cache: %s", err.Error())
		return
	}
	serverUpdateStatusCache.Lock()
	defer serverUpdateStatusCache.Unlock()
	serverUpdateStatusCache.serverMap = newServerUpdateStatuses
	serverUpdateStatusCache.initialized = true
	log.Infof("refreshed server update status cache (len = %d)", len(serverUpdateStatusCache.serverMap))
}

type serverInfo struct {
	id                 int
	hostName           string
	typeName           string
	cdnId              int
	status             string
	cachegroup         int
	configUpdateTime   *time.Time
	configApplyTime    *time.Time
	configUpdateFailed bool
	revalUpdateTime    *time.Time
	revalApplyTime     *time.Time
	revalUpdateFailed  bool
}

const getUseRevalPendingQuery = `
	SELECT value::BOOLEAN
	FROM parameter
	WHERE name = 'use_reval_pending' AND config_file = 'global'
	UNION ALL SELECT FALSE FETCH FIRST 1 ROW ONLY
`

const getServerInfoQuery = `
	SELECT
		s.id,
		s.host_name,
		t.name,
		s.cdn_id,
		st.name,
		s.cachegroup,
		s.config_update_time,
		s.config_apply_time,
		s.config_update_failed,
		s.revalidate_update_time,
		s.revalidate_apply_time,
		s.revalidate_update_failed
	FROM server s
	JOIN type t ON t.id = s.type
	JOIN status st ON st.id = s.status
`

const getCacheGroupsQuery = `
	SELECT
		c.id,
		c.parent_cachegroup_id,
		c.secondary_parent_cachegroup_id
	FROM cachegroup c
`

const getTopologyCacheGroupParentsQuery = `
	SELECT
		cg_child.id,
		ARRAY_AGG(DISTINCT cg_parent.id)
	FROM topology_cachegroup_parents tcp
	JOIN topology_cachegroup tc_child ON tc_child.id = tcp.child
	JOIN cachegroup cg_child ON cg_child.name = tc_child.cachegroup
	JOIN topology_cachegroup tc_parent ON tc_parent.id = tcp.parent
	JOIN cachegroup cg_parent ON cg_parent.name = tc_parent.cachegroup
	GROUP BY cg_child.id
`

func getServerUpdateStatuses(db *sql.DB, timeout time.Duration) (map[string][]tc.ServerUpdateStatusV5, error) {
	dbCtx, dbClose := context.WithTimeout(context.Background(), timeout)
	defer dbClose()
	serversByID := make(map[int]serverInfo)
	updatePendingByCDNCachegroup := make(map[int]map[int]bool)
	revalPendingByCDNCachegroup := make(map[int]map[int]bool)
	tx, err := db.BeginTx(dbCtx, nil)
	if err != nil {
		return nil, fmt.Errorf("beginning server update status transaction: %w", err)
	}
	defer func() {
		if err := tx.Commit(); err != nil && err != sql.ErrTxDone {
			log.Errorln("committing server update status transaction: " + err.Error())
		}
	}()

	useRevalPending := false
	if err := tx.QueryRowContext(dbCtx, getUseRevalPendingQuery).Scan(&useRevalPending); err != nil {
		return nil, fmt.Errorf("querying use_reval_pending param: %w", err)
	}

	// get all servers and build map of update/revalPending by cachegroup+CDN
	serverRows, err := tx.QueryContext(dbCtx, getServerInfoQuery)
	if err != nil {
		return nil, fmt.Errorf("querying servers: %w", err)
	}
	defer log.Close(serverRows, "closing server rows")
	for serverRows.Next() {
		s := serverInfo{}
		if err := serverRows.Scan(&s.id, &s.hostName, &s.typeName, &s.cdnId, &s.status, &s.cachegroup, &s.configUpdateTime, &s.configApplyTime, &s.configUpdateFailed, &s.revalUpdateTime, &s.revalApplyTime, &s.revalUpdateFailed); err != nil {
			return nil, fmt.Errorf("scanning servers: %w", err)
		}
		serversByID[s.id] = s
		if _, ok := updatePendingByCDNCachegroup[s.cdnId]; !ok {
			updatePendingByCDNCachegroup[s.cdnId] = make(map[int]bool)
		}
		if _, ok := revalPendingByCDNCachegroup[s.cdnId]; !ok {
			revalPendingByCDNCachegroup[s.cdnId] = make(map[int]bool)
		}
		status := tc.CacheStatusFromString(s.status)
		if tc.IsValidCacheType(s.typeName) && (status == tc.CacheStatusOnline || status == tc.CacheStatusReported || status == tc.CacheStatusAdminDown) {
			if s.configUpdateTime.After(*s.configApplyTime) {
				updatePendingByCDNCachegroup[s.cdnId][s.cachegroup] = true
			}
			if s.revalUpdateTime.After(*s.revalApplyTime) {
				revalPendingByCDNCachegroup[s.cdnId][s.cachegroup] = true
			}
		}
	}
	if err := serverRows.Err(); err != nil {
		return nil, fmt.Errorf("iterating over server rows: %w", err)
	}

	// get all legacy cachegroup parents
	cacheGroupParents := make(map[int]map[int]struct{})
	cacheGroupRows, err := tx.QueryContext(dbCtx, getCacheGroupsQuery)
	if err != nil {
		return nil, fmt.Errorf("querying cachegroups: %w", err)
	}
	defer log.Close(cacheGroupRows, "closing cachegroup rows")
	for cacheGroupRows.Next() {
		id := 0
		parentID := new(int)
		secondaryParentID := new(int)
		if err := cacheGroupRows.Scan(&id, &parentID, &secondaryParentID); err != nil {
			return nil, fmt.Errorf("scanning cachegroups: %w", err)
		}
		cacheGroupParents[id] = make(map[int]struct{})
		if parentID != nil {
			cacheGroupParents[id][*parentID] = struct{}{}
		}
		if secondaryParentID != nil {
			cacheGroupParents[id][*secondaryParentID] = struct{}{}
		}
	}
	if err := cacheGroupRows.Err(); err != nil {
		return nil, fmt.Errorf("iterating over cachegroup rows: %w", err)
	}

	// get all topology-based cachegroup parents
	topologyCachegroupRows, err := tx.QueryContext(dbCtx, getTopologyCacheGroupParentsQuery)
	if err != nil {
		return nil, fmt.Errorf("querying topology cachegroups: %w", err)
	}
	defer log.Close(topologyCachegroupRows, "closing topology cachegroup rows")
	for topologyCachegroupRows.Next() {
		id := 0
		parents := []int32{}
		if err := topologyCachegroupRows.Scan(&id, pq.Array(&parents)); err != nil {
			return nil, fmt.Errorf("scanning topology cachegroup rows: %w", err)
		}
		for _, p := range parents {
			cacheGroupParents[id][int(p)] = struct{}{}
		}
	}
	if err = topologyCachegroupRows.Err(); err != nil {
		return nil, fmt.Errorf("iterating over topology cachegroup rows: %w", err)
	}

	serverUpdateStatuses := make(map[string][]tc.ServerUpdateStatusV5, len(serversByID))
	for serverID, server := range serversByID {
		updateStatus := tc.ServerUpdateStatusV5{
			HostName:               server.hostName,
			UpdatePending:          server.configUpdateTime.After(*server.configApplyTime),
			RevalPending:           server.revalUpdateTime.After(*server.revalApplyTime),
			UseRevalPending:        useRevalPending,
			HostId:                 serverID,
			Status:                 server.status,
			ParentPending:          getParentPending(cacheGroupParents[server.cachegroup], updatePendingByCDNCachegroup[server.cdnId]),
			ParentRevalPending:     getParentPending(cacheGroupParents[server.cachegroup], revalPendingByCDNCachegroup[server.cdnId]),
			ConfigUpdateTime:       server.configUpdateTime,
			ConfigApplyTime:        server.configApplyTime,
			ConfigUpdateFailed:     &server.configUpdateFailed,
			RevalidateUpdateTime:   server.revalUpdateTime,
			RevalidateApplyTime:    server.revalApplyTime,
			RevalidateUpdateFailed: &server.revalUpdateFailed,
		}
		serverUpdateStatuses[server.hostName] = append(serverUpdateStatuses[server.hostName], updateStatus)
	}
	return serverUpdateStatuses, nil
}

func getParentPending(parents map[int]struct{}, pendingByCacheGroup map[int]bool) bool {
	for parent := range parents {
		if pendingByCacheGroup[parent] {
			return true
		}
	}
	return false
}
