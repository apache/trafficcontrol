package crconfig

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
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/lib/pq"
)

// getCachegroupFallbacks returns a map[primaryCacheGroupID][]fallbackCacheGroupName, the last updated time of any fallback, and any error.
func getCachegroupFallbacks(tx *sql.Tx) (map[int][]string, time.Time, error) {
	q := `
SELECT
  cachegroup_fallbacks.primary_cg,
  cachegroup.name,
  GREATEST(cachegroup.last_updated, cachegroup_fallbacks.last_updated) as last_updated
FROM
  cachegroup_fallbacks
  JOIN cachegroup on cachegroup_fallbacks.backup_cg = cachegroup.id
ORDER BY cachegroup_fallbacks.set_order
`
	rows, err := tx.Query(q)
	if err != nil {
		return nil, time.Time{}, errors.New("Error retrieving from cachegroup_fallbacks: " + err.Error())
	}
	defer rows.Close()

	lastUpdated := time.Time{}
	fallbacks := map[int][]string{} // map[primaryCacheGroupID][]fallbackCacheGroupName
	for rows.Next() {
		primaryCGID := 0
		fallbackCG := ""
		fbLastUpdated := time.Time{}
		if err := rows.Scan(&primaryCGID, &fallbackCG, &fbLastUpdated); err != nil {
			return nil, time.Time{}, errors.New("scanning cachegroup_fallbacks: " + err.Error())
		}
		fallbacks[primaryCGID] = append(fallbacks[primaryCGID], fallbackCG)
		if fbLastUpdated.After(lastUpdated) {
			lastUpdated = fbLastUpdated
		}
	}
	if err := rows.Err(); err != nil {
		return nil, time.Time{}, errors.New("cachegroup_fallbacks rows: " + err.Error())
	}
	return fallbacks, lastUpdated, nil
}

// makeLocations returns the map of edge locations, router locations, the last updated time of any location, and any error
func makeLocations(cdn string, tx *sql.Tx) (map[string]tc.CRConfigLatitudeLongitude, map[string]tc.CRConfigLatitudeLongitude, time.Time, error) {
	edgeLocs := map[string]tc.CRConfigLatitudeLongitude{}
	routerLocs := map[string]tc.CRConfigLatitudeLongitude{}

	fallbacks, lastUpdated, err := getCachegroupFallbacks(tx)
	if err != nil {
		return nil, nil, time.Time{}, err
	}

	// TODO test whether it's faster to do a single query, joining lat/lon into servers
	qry := `
SELECT
  cg.name,
  cg.id,
  t.name as type,
  co.latitude,
  co.longitude,
  COALESCE(cg.fallback_to_closest, TRUE),
  (SELECT array_agg(method::text) as localization_methods FROM cachegroup_localization_method WHERE cachegroup = cg.id),
  GREATEST(cg.last_updated, co.last_updated) as last_updated
FROM
  cachegroup cg
  LEFT JOIN coordinate co on co.id = cg.coordinate
  JOIN server s on s.cachegroup = cg.id
  JOIN type t on t.id = s.type
  JOIN status st ON st.id = s.status
  JOIN cdn ON cdn.id = s.cdn_id
WHERE
  cdn.name = $1
  AND (t.name like 'EDGE%' or t.name = 'CCR')
  AND (st.name = 'REPORTED' or st.name = 'ONLINE' or st.name = 'ADMIN_DOWN')
`
	// TODO pass edge type prefix, router type name
	rows, err := tx.Query(qry, cdn)
	if err != nil {
		return nil, nil, time.Time{}, errors.New("Error querying cachegroups: " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		cachegroup := ""
		primaryCacheID := 0
		ttype := ""
		fallbackToClosest := true
		latlon := tc.CRConfigLatitudeLongitude{}
		cgLastUpdated := time.Time{}
		if err := rows.Scan(&cachegroup, &primaryCacheID, &ttype, &latlon.Lat, &latlon.Lon, &fallbackToClosest, pq.Array(&latlon.LocalizationMethods), &cgLastUpdated); err != nil {
			return nil, nil, time.Time{}, errors.New("Error scanning cachegroup: " + err.Error())
		}
		if len(latlon.LocalizationMethods) == 0 {
			// to keep current default behavior when localizationMethods is unset/empty, enable all current localization methods
			latlon.LocalizationMethods = []tc.LocalizationMethod{tc.LocalizationMethodGeo, tc.LocalizationMethodCZ, tc.LocalizationMethodDeepCZ}
		}
		if ttype == RouterTypeName {
			routerLocs[cachegroup] = latlon
		} else {
			latlon.BackupLocations.FallbackToClosest = fallbackToClosest
			latlon.BackupLocations.List = fallbacks[primaryCacheID]
			edgeLocs[cachegroup] = latlon
		}
		if cgLastUpdated.After(lastUpdated) {
			lastUpdated = cgLastUpdated
		}
	}
	if err := rows.Err(); err != nil {
		return nil, nil, time.Time{}, errors.New("Error iterating cachegroup rows: " + err.Error())
	}
	return edgeLocs, routerLocs, lastUpdated, nil
}
