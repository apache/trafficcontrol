// Package capabilities contains logic and handlers for the /capabilities API
// endpoint.
//
// Deprecated: "Capabilities" (now called Permissions) are no longer handled
// this way, and this package should be removed once API versions that use it
// have been fully removed.
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
	"errors"
	"fmt"
	"net/http"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/dbhelpers"
)

const readQuery = `
SELECT description,
       last_updated,
       name
FROM capability
`

// Read handles GET requests to /capabilities.
//
// Deprecated: "Capabilities" (now called Permissions) are no longer handled
// this way, and this package should be removed once API versions that use it
// have been fully removed.
func Read(w http.ResponseWriter, r *http.Request) {
	inf, sysErr, userErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleDeprecatedErr(w, r, tx, errCode, userErr, sysErr, nil)
		return
	}
	defer inf.Close()

	cols := map[string]dbhelpers.WhereColumnInfo{
		"name": {Column: "capability.name"},
	}

	where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(inf.Params, cols)
	if len(errs) > 0 {
		errCode = http.StatusBadRequest
		userErr = util.JoinErrs(errs)
		api.HandleDeprecatedErr(w, r, tx, errCode, userErr, nil, nil)
		return
	}

	query := readQuery + where + orderBy + pagination
	rows, err := inf.Tx.NamedQuery(query, queryValues)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		errCode = http.StatusInternalServerError
		sysErr = fmt.Errorf("querying capabilities: %w", err)
		api.HandleDeprecatedErr(w, r, tx, errCode, nil, sysErr, nil)
		return
	}
	defer log.Close(rows, "closing Capabilities rows")

	caps := []tc.Capability{}
	for rows.Next() {
		var capability tc.Capability
		if err := rows.Scan(&capability.Description, &capability.LastUpdated, &capability.Name); err != nil {
			errCode = http.StatusInternalServerError
			sysErr = fmt.Errorf("parsing database response: %w", err)
			api.HandleDeprecatedErr(w, r, tx, errCode, nil, sysErr, nil)
			return
		}

		caps = append(caps, capability)
	}

	api.WriteRespAlertObj(w, r, tc.WarnLevel, "This endpoint is deprecated, and will be removed in the future", caps)
}
