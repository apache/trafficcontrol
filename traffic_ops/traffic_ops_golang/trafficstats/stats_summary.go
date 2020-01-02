package trafficstats

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
	"fmt"
	"net/http"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/jmoiron/sqlx"

	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
)

// GetStatsSummary handler for getting stats summaries
func GetStatsSummary(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{}, []string{})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	lastSummaryDateStr := inf.Params["lastSummaryDate"]
	if len(lastSummaryDateStr) != 0 { // Perl only checked for existence of query param
		getLastSummaryDate(w, r, inf)
		return
	}

	getStatsSummary(w, r, inf)
	return
}

func getLastSummaryDate(w http.ResponseWriter, r *http.Request, inf *api.APIInfo) {
	queryParamsToSQLCols := map[string]dbhelpers.WhereColumnInfo{
		"statName": dbhelpers.WhereColumnInfo{"stat_name", nil},
	}
	where, _, _, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(inf.Params, queryParamsToSQLCols)
	if len(errs) > 0 {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, util.JoinErrs(errs))
		return
	}
	query := selectQuery() + where + " ORDER BY summary_time DESC"
	statsSummaries, err := queryStatsSummary(inf.Tx, query, queryValues)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
		return
	}
	resp := tc.StatsSummaryLastUpdated{}
	if len(statsSummaries) >= 1 {
		resp.SummaryTime = &statsSummaries[0].SummaryTime
	}
	api.WriteResp(w, r, resp)
}

func getStatsSummary(w http.ResponseWriter, r *http.Request, inf *api.APIInfo) {
	queryParamsToSQLCols := map[string]dbhelpers.WhereColumnInfo{
		"statName":            dbhelpers.WhereColumnInfo{"stat_name", nil},
		"cdnName":             dbhelpers.WhereColumnInfo{"cdn_name", nil},
		"deliveryServiceName": dbhelpers.WhereColumnInfo{"deliveryservice_name", nil},
	}
	where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(inf.Params, queryParamsToSQLCols)
	if len(errs) > 0 {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, util.JoinErrs(errs))
		return
	}
	query := selectQuery() + where + orderBy + pagination
	statsSummaries, err := queryStatsSummary(inf.Tx, query, queryValues)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
		return
	}

	api.WriteResp(w, r, statsSummaries)
}

func queryStatsSummary(tx *sqlx.Tx, q string, queryValues map[string]interface{}) ([]tc.StatsSummary, error) {
	rows, err := tx.NamedQuery(q, queryValues)
	if err != nil {
		return nil, fmt.Errorf("querying stats summary: %v", err)
	}
	defer rows.Close()

	statsSummaries := []tc.StatsSummary{}
	for rows.Next() {
		s := tc.StatsSummary{}
		if err = rows.StructScan(&s); err != nil {
			return nil, fmt.Errorf("scanning stats summary: %v", err)
		}
		statsSummaries = append(statsSummaries, s)
	}
	return statsSummaries, nil
}

func selectQuery() string {
	return `SELECT
id,
cdn_name,
deliveryservice_name,
stat_name,
stat_value,
summary_time,
stat_date
FROM stats_summary`
}
