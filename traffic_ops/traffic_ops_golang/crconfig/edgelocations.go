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

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/lib/pq"
)

// getCachegroupFallbacks returns a map[primaryCacheGroupID][]fallbackCacheGroupName.
func getCachegroupFallbacks(tx *sql.Tx) (map[int][]string, error) {
	q := `
SELECT
  cachegroup_fallbacks.primary_cg,
  cachegroup.name
FROM
  cachegroup_fallbacks
  JOIN cachegroup on cachegroup_fallbacks.backup_cg = cachegroup.id
ORDER BY cachegroup_fallbacks.set_order
`
	rows, err := tx.Query(q)
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

func makeLocations(cdn string, tx *sql.Tx) (map[string]tc.CRConfigLatitudeLongitude, map[string]tc.CRConfigLatitudeLongitude, error) {
	edgeLocs := map[string]tc.CRConfigLatitudeLongitude{}
	routerLocs := map[string]tc.CRConfigLatitudeLongitude{}

	fallbacks, err := getCachegroupFallbacks(tx)
	if err != nil {
		return nil, nil, err
	}

	// TODO test whether it's faster to do a single query, joining lat/lon into servers
	q := `
select cg.name, cg.id, t.name as type, co.latitude, co.longitude, COALESCE(cg.fallback_to_closest, TRUE),
(SELECT array_agg(method::text) as localization_methods FROM cachegroup_localization_method WHERE cachegroup = cg.id)
from cachegroup as cg
left join coordinate as co on co.id = cg.coordinate
inner join server as s on s.cachegroup = cg.id
inner join type as t on t.id = s.type
inner join status as st ON st.id = s.status
where s.cdn_id = (select id from cdn where name = $1)
and (t.name like 'EDGE%' or t.name = 'CCR')
and (st.name = 'REPORTED' or st.name = 'ONLINE' or st.name = 'ADMIN_DOWN')
`
	// TODO pass edge type prefix, router type name
	rows, err := tx.Query(q, cdn)
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
