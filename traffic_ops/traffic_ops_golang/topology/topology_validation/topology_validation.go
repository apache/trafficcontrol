// The topology_validation package is for topology validation functions that are used outside of the topology
// package, in order to prevent import cycles.

package topology_validation

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
	"errors"
	"fmt"
	"strings"

	"github.com/apache/trafficcontrol/v8/lib/go-log"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// CheckForEmptyCacheGroups checks if the cachegroups are empty (altogether) or empty in any of the given CDN IDs.
// If cachegroupsInTopology is true, it will only check cachegroups that are used in a topology. Any server IDs in
// excludeServerIds will not be counted.
func CheckForEmptyCacheGroups(tx *sqlx.Tx, cacheGroupIds []int, CDNIDs []int, cachegroupsInTopology bool, excludeServerIds []int) error {
	if excludeServerIds == nil {
		excludeServerIds = []int{}
	}
	var (
		baseError   = errors.New("unable to check for cachegroups with no servers")
		systemError = "checking for cachegroups with no servers: %s"
		query       = selectEmptyCacheGroupsQuery(cachegroupsInTopology)
		parameters  = map[string]interface{}{
			"cachegroup_ids":     pq.Array(cacheGroupIds),
			"exclude_server_ids": pq.Array(excludeServerIds),
		}
	)

	rows, err := tx.NamedQuery(query, parameters)
	if err != nil {
		log.Errorf(systemError, err.Error())
		return baseError
	}

	var (
		serverCountByCDN int
		cacheGroup       string
		cdnID            *int
	)
	cgServerCountsByCDN := make(map[int]map[string]int)
	cgServerCounts := make(map[string]int)
	topologySetByCachegroup := make(map[string]map[string]struct{})
	defer log.Close(rows, "unable to close DB connection when checking for cachegroups with no servers")
	for rows.Next() {
		var scanTo = []interface{}{&cacheGroup, &cdnID, &serverCountByCDN}
		var topologiesForRow []string
		if cachegroupsInTopology {
			scanTo = append(scanTo, pq.Array(&topologiesForRow))
		}
		if err := rows.Scan(scanTo...); err != nil {
			log.Errorf(systemError, err.Error())
			return baseError
		}
		if cdnID != nil {
			if _, ok := cgServerCountsByCDN[*cdnID]; !ok {
				cgServerCountsByCDN[*cdnID] = make(map[string]int)
			}
			cgServerCountsByCDN[*cdnID][cacheGroup] = serverCountByCDN
		}
		cgServerCounts[cacheGroup] += serverCountByCDN

		if cachegroupsInTopology {
			if _, ok := topologySetByCachegroup[cacheGroup]; !ok {
				topologySetByCachegroup[cacheGroup] = make(map[string]struct{})
			}
			for _, topology := range topologiesForRow {
				topologySetByCachegroup[cacheGroup][topology] = struct{}{}
			}
		}
	}
	topologiesByCachegroup := make(map[string][]string, len(topologySetByCachegroup))
	for cg, topologySet := range topologySetByCachegroup {
		for topology := range topologySet {
			topologiesByCachegroup[cg] = append(topologiesByCachegroup[cg], topology)
		}
	}
	emptyCachegroups := []string{}
	for cg, count := range cgServerCounts {
		if count == 0 {
			messageEntry := cg
			if cachegroupsInTopology {
				messageEntry += " (in topologies: " + strings.Join(topologiesByCachegroup[cg], ", ") + ")"
			}
			emptyCachegroups = append(emptyCachegroups, messageEntry)
		}
	}

	if len(emptyCachegroups) > 0 {
		errMessage := "cachegroups with no servers in them: " + strings.Join(emptyCachegroups, ", ")
		return errors.New(errMessage)
	}

	errMessage := []string{}
	for _, cdnID := range CDNIDs {
		if _, ok := cgServerCountsByCDN[cdnID]; !ok {
			return fmt.Errorf("topology is assigned to delivery service on CDN %d, but that CDN has no servers", cdnID)
		}
		emptyCachegroupsByCDN := []string{}
		for cg, serverCount := range cgServerCountsByCDN[cdnID] {
			if serverCount == 0 {
				emptyCachegroupsByCDN = append(emptyCachegroupsByCDN, cg)
			}
		}
		// check that this CDN has a count for all given cachegroups
		for cg := range cgServerCounts {
			if _, ok := cgServerCountsByCDN[cdnID][cg]; !ok {
				emptyCachegroupsByCDN = append(emptyCachegroupsByCDN, cg)
			}
		}
		if len(emptyCachegroupsByCDN) > 0 {
			errMessage = append(errMessage, fmt.Sprintf("topology is assigned to delivery service(s) on CDN %d, but the following cachegroups have no servers in CDN %d: %s", cdnID, cdnID, strings.Join(emptyCachegroupsByCDN, ", ")))
		}
	}
	if len(errMessage) > 0 {
		return errors.New(strings.Join(errMessage, "; "))
	}
	return nil
}

func selectEmptyCacheGroupsQuery(cachegroupsInTopology bool) string {
	var joinTopologyCachegroups string
	var topologyNames string
	if cachegroupsInTopology {
		// language=SQL
		topologyNames = `
		, ARRAY_AGG(tc.topology)
`
		// language=SQL
		joinTopologyCachegroups = `
		JOIN topology_cachegroup tc ON c."name" = tc.cachegroup
`
	}
	// language=SQL
	query := fmt.Sprintf(`
		SELECT
			c."name",
			s.cdn_id,
			COUNT(*) FILTER (
			    WHERE s.id IS NOT NULL
			    AND NOT(s."id" = ANY(CAST(:exclude_server_ids AS INT[])))
			) AS server_count %s
		FROM cachegroup c
		%s
		LEFT JOIN "server" s ON c.id = s.cachegroup
		WHERE c."id" = ANY(CAST(:cachegroup_ids AS BIGINT[]))
		GROUP BY c."name", s.cdn_id
	`, topologyNames, joinTopologyCachegroups)
	return query
}
