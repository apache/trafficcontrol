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

	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
)

func makeLocations(cdn string, db *sql.DB) (map[string]tc.CRConfigLatitudeLongitude, map[string]tc.CRConfigLatitudeLongitude, error) {
	edgeLocs := map[string]tc.CRConfigLatitudeLongitude{}
	routerLocs := map[string]tc.CRConfigLatitudeLongitude{}

	// TODO test whether it's faster to do a single query, joining lat/lon into servers
	q := `
select cg.name, t.name as type, cg.latitude, cg.longitude from cachegroup as cg
inner join server as s on s.cachegroup = cg.id
inner join type as t on t.id = s.type
inner join status as st ON st.id = s.status
where s.cdn_id = (select id from cdn where name = $1)
and (t.name like 'EDGE%' or t.name = 'CCR')
and (st.name = 'REPORTED' or st.name = 'ONLINE' or st.name = 'ADMIN_DOWN')
`
	// TODO pass edge type prefix, router type name
	rows, err := db.Query(q, cdn)
	if err != nil {
		return nil, nil, errors.New("Error querying cachegroups: " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		cachegroup := ""
		ttype := ""
		latlon := tc.CRConfigLatitudeLongitude{}
		if err := rows.Scan(&cachegroup, &ttype, &latlon.Lat, &latlon.Lon); err != nil {
			return nil, nil, errors.New("Error scanning cachegroup: " + err.Error())
		}
		if ttype == RouterTypeName {
			routerLocs[cachegroup] = latlon
		} else {
			edgeLocs[cachegroup] = latlon
		}
	}
	if err := rows.Err(); err != nil {
		return nil, nil, errors.New("Error iterating cachegroup rows: " + err.Error())
	}
	return edgeLocs, routerLocs, nil
}
