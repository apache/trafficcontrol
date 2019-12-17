package api_capability

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/jmoiron/sqlx"
)

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

func GetAPICapabilitiesHandler(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	results, usrErr, sysErr := getAPICapabilities(inf.Tx, inf.Params)

	if usrErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, usrErr, nil)
		return
	}

	if sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, nil, sysErr)
		return
	}

	api.WriteResp(w, r, tc.APICapabilityResponse{Response: results})
	return
}

func getAPICapabilities(tx *sqlx.Tx, params map[string]string) ([]tc.APICapability, error, error) {
	selectQuery := `SELECT * FROM api_capability`
	queryParamsToQueryCols := map[string]dbhelpers.WhereColumnInfo{
		"id":          dbhelpers.WhereColumnInfo{"id", api.IsInt},
		"capability":  dbhelpers.WhereColumnInfo{"capability", nil},
		"httpMethod":  dbhelpers.WhereColumnInfo{"http_method", nil},
		"route":       dbhelpers.WhereColumnInfo{"route", nil},
		"lastUpdated": dbhelpers.WhereColumnInfo{"last_updated", nil},
	}

	where, orderBy, pagination, queryValues, errs :=
		dbhelpers.BuildWhereAndOrderByAndPagination(params, queryParamsToQueryCols)

	if len(errs) > 0 {
		return nil, errors.New(
			fmt.Sprintf(
				"query exception: could not build api_capbility query with params: %v, error: %s",
				params,
				errs[0].Error(),
			),
		), nil
	}

	query := selectQuery + where + orderBy + pagination
	rows, err := tx.NamedQuery(query, queryValues)
	if err != nil {
		return nil, errors.New(
			fmt.Sprintf(
				"db exception: could not execute api_capbility query with params: %v, error: %s",
				params,
				err.Error(),
			),
		), nil
	}
	defer rows.Close()

	var apiCaps []tc.APICapability
	for rows.Next() {
		var ac tc.APICapability
		err = rows.Scan(
			&ac.ID,
			&ac.HTTPMethod,
			&ac.Route,
			&ac.Capability,
			&ac.LastUpdated,
		)
		if err != nil {
			return nil, nil, errors.New(fmt.Sprintf("api capability read: scanning: %s", err.Error()))
		}
		apiCaps = append(apiCaps, ac)
	}

	return apiCaps, nil, nil
}
