package capabilities

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
	"fmt"
	"net/http"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
)

const readQuery = `
SELECT description,
       last_updated,
       name
FROM capability
`

func Read(w http.ResponseWriter, r *http.Request) {
	inf, sysErr, userErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	cols := map[string]dbhelpers.WhereColumnInfo{
		"name": dbhelpers.WhereColumnInfo{Column: "capability.name"},
	}

	where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(inf.Params, cols)
	if len(errs) > 0 {
		errCode = http.StatusBadRequest
		userErr = util.JoinErrs(errs)
		api.HandleErr(w, r, tx, errCode, userErr, nil)
		return
	}

	query := readQuery + where + orderBy + pagination
	rows, err := inf.Tx.NamedQuery(query, queryValues)
	if err != nil && err != sql.ErrNoRows {
		errCode = http.StatusInternalServerError
		sysErr = fmt.Errorf("querying capabilities: %v", err)
		api.HandleErr(w, r, tx, errCode, nil, sysErr)
		return
	}
	defer rows.Close()

	caps := []tc.Capability{}
	for rows.Next() {
		cap := tc.Capability{}
		if err := rows.Scan(&cap.Description, &cap.LastUpdated, &cap.Name); err != nil {
			errCode = http.StatusInternalServerError
			sysErr = fmt.Errorf("Parsing database response: %v", err)
			api.HandleErr(w, r, tx, errCode, nil, sysErr)
			return
		}

		caps = append(caps, cap)
	}

	api.WriteResp(w, r, caps)
}
