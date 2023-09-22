package topology

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

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/lib/pq"
)

// MakeTopologies makes the topologies data for the crconfig and tmconfig snapshots.
func MakeTopologies(tx *sql.Tx) (map[string]tc.CRConfigTopology, error) {
	query := `
SELECT
	t.name,
	(SELECT ARRAY_AGG(tc.cachegroup ORDER BY tc.cachegroup)
		FROM topology_cachegroup tc
		JOIN cachegroup c ON c.name = tc.cachegroup
		JOIN type ON type.id = c.type
		WHERE t.name = tc.topology
		AND type.name = $1
		) AS nodes
FROM topology t
ORDER BY t.name
`
	var rows *sql.Rows
	var err error
	if rows, err = tx.Query(query, tc.CacheGroupEdgeTypeName); err != nil {
		return nil, errors.New("querying topologies: " + err.Error())
	}
	defer log.Close(rows, "unable to close DB connection")

	topologies := map[string]tc.CRConfigTopology{}
	for rows.Next() {
		topology := tc.CRConfigTopology{}
		var name string
		if err = rows.Scan(
			&name,
			pq.Array(&topology.Nodes),
		); err != nil {
			return nil, errors.New("scanning topology: " + err.Error())
		}
		topologies[name] = topology
	}
	return topologies, nil
}
