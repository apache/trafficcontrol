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

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/lib/pq"
)

// getCachegroupFallbacks returns a map[primaryCacheGroupID][]fallbackCacheGroupName.
func getCachegroupFallbacks(tx *sql.Tx, cdn string, live bool) (map[int][]string, error) {
	qry := `
WITH cdn_name AS (
  SELECT $1::text as v
),
snapshot_time AS (
  SELECT time as v FROM snapshot sn where sn.cdn = (SELECT v from cdn_name)
)
SELECT primary_cg, backup_cg FROM (
SELECT DISTINCT ON (cgf.primary_cg, cg.name)
  cgf.primary_cg as primary_cg,
  cg.name as backup_cg,
  cgf.set_order,
	cgf.deleted
FROM
  cachegroup_fallbacks_snapshot cgf
  JOIN cachegroup_snapshot cg ON cgf.backup_cg = cg.id
`
	if !live {
		qry += `
WHERE
  cgf.last_updated <= (select v from snapshot_time)
  AND cg.last_updated <= (select v from snapshot_time)
`
	}
	qry += `
ORDER BY
  cgf.primary_cg,
  cg.name,
  cgf.last_updated,
  cg.last_updated
DESC
) v where deleted = false
ORDER BY set_order ASC
`
	rows, err := tx.Query(qry, cdn)
	if err != nil {
		return nil, errors.New("Error retrieving from cachegroup_fallbacks: " + err.Error())
	}
	defer rows.Close()

	fallbacks := map[int][]string{} // map[primaryCacheGroupID][]fallbackCacheGroupName
	for rows.Next() {
		primaryCGID := 0
		fallbackCG := ""
		if err := rows.Scan(&primaryCGID, &fallbackCG); err != nil {
			return nil, errors.New("scanning cachegroup_fallbacks: " + err.Error())
		}
		fallbacks[primaryCGID] = append(fallbacks[primaryCGID], fallbackCG)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.New("cachegroup_fallbacks rows: " + err.Error())
	}
	return fallbacks, nil
}

func makeLocations(tx *sql.Tx, cdn string, live bool) (map[string]tc.CRConfigLatitudeLongitude, map[string]tc.CRConfigLatitudeLongitude, error) {
	edgeLocs := map[string]tc.CRConfigLatitudeLongitude{}
	routerLocs := map[string]tc.CRConfigLatitudeLongitude{}

	fallbacks, err := getCachegroupFallbacks(tx, cdn, live)
	if err != nil {
		return nil, nil, err
	}

	// TODO test whether it's faster to do a single query, joining lat/lon into servers
	qry := `
WITH cdn_name AS (
  SELECT $1::text as v
),
snapshot_time AS (
  SELECT time as v FROM snapshot sn where sn.cdn = (SELECT v from cdn_name)
)
SELECT
  name,
  id,
  type,
  latitude,
  longitude,
  fallback_to_closest,
  localization_methods
FROM (
SELECT DISTINCT ON (cg.name)
  cg.name,
  cg.id,
  t.name as type,
  co.latitude,
  co.longitude,
  COALESCE(cg.fallback_to_closest, TRUE) as fallback_to_closest,
  (
    SELECT array_agg(method::text) as localization_methods FROM (
    SELECT DISTINCT ON (cgl.cachegroup, cgl.method)
      cgl.cachegroup,
      cgl.method,
      cgl.deleted
    FROM cachegroup_localization_method_snapshot cgl
    WHERE cachegroup = cg.id AND cgl.last_updated <= (select v from snapshot_time)
    ORDER BY
      cgl.cachegroup,
      cgl.method,
      cgl.last_updated
    DESC
    ) v WHERE deleted = false
  ),
	cg.deleted
FROM
  cachegroup_snapshot cg
  LEFT JOIN coordinate_snapshot co ON co.id = cg.coordinate
  JOIN server_snapshot s ON s.cachegroup = cg.id
  JOIN type_snapshot t ON t.id = s.type
  JOIN status_snapshot st ON st.id = s.status
WHERE
  s.cdn_id = (select id from cdn_snapshot c where c.name = (select v from cdn_name) and c.last_updated <= (select v from snapshot_time))
  AND (t.name like 'EDGE%' or t.name = 'CCR')
  AND (st.name = 'REPORTED' or st.name = 'ONLINE' or st.name = 'ADMIN_DOWN')
`
	if !live {
		qry += `
  AND cg.last_updated <= (select v from snapshot_time)
  AND co.last_updated <= (select v from snapshot_time)
  AND s.last_updated <= (select v from snapshot_time)
  AND t.last_updated <= (select v from snapshot_time)
  AND st.last_updated <= (select v from snapshot_time)
`
	}
	qry += `
ORDER BY
  cg.name,
  cg.last_updated,
  co.last_updated,
  s.last_updated,
  t.last_updated,
  st.last_updated
DESC
) v where deleted = false
`
	// TODO pass edge type prefix, router type name
	rows, err := tx.Query(qry, cdn)
	if err != nil {
		return nil, nil, errors.New("Error querying cachegroups: " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		cachegroup := ""
		primaryCacheID := 0
		ttype := ""
		fallbackToClosest := true
		latlon := tc.CRConfigLatitudeLongitude{}
		if err := rows.Scan(&cachegroup, &primaryCacheID, &ttype, &latlon.Lat, &latlon.Lon, &fallbackToClosest, pq.Array(&latlon.LocalizationMethods)); err != nil {
			return nil, nil, errors.New("Error scanning cachegroup: " + err.Error())
		}
		if len(latlon.LocalizationMethods) == 0 {
			// to keep current default behavior when localizationMethods is unset/empty, enable all current localization methods
			latlon.LocalizationMethods = []tc.LocalizationMethod{tc.LocalizationMethodGeo, tc.LocalizationMethodCZ, tc.LocalizationMethodDeepCZ}
		}
		if ttype == tc.RouterTypeName {
			routerLocs[cachegroup] = latlon
		} else {
			latlon.BackupLocations.FallbackToClosest = fallbackToClosest
			latlon.BackupLocations.List = fallbacks[primaryCacheID]
			edgeLocs[cachegroup] = latlon
		}
	}
	if err := rows.Err(); err != nil {
		return nil, nil, errors.New("Error iterating cachegroup rows: " + err.Error())
	}
	return edgeLocs, routerLocs, nil
}
