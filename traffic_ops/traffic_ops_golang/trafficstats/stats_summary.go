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

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/jmoiron/sqlx"
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

func getLastSummaryDate(w http.ResponseWriter, r *http.Request, inf *api.Info) {
	queryParamsToSQLCols := map[string]dbhelpers.WhereColumnInfo{
		"statName": dbhelpers.WhereColumnInfo{Column: "stat_name"},
	}
	where, _, _, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(inf.Params, queryParamsToSQLCols)
	if len(errs) > 0 {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, util.JoinErrs(errs))
		return
	}
	query := selectQuery() + where + " ORDER BY summary_time DESC"
	statsSummaries, err := queryStatsSummary(inf.Tx, inf.Version.Major, query, queryValues)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
		return
	}

	if inf.Version.Major >= 5 {
		resp := tc.StatsSummaryLastUpdatedV5{}
		if len(statsSummaries) >= 1 {
			if summary, ok := statsSummaries[0].(tc.StatsSummaryV5); ok {
				resp.SummaryTime = &summary.SummaryTime
			}
		}
		api.WriteResp(w, r, resp)

	} else {
		resp := tc.StatsSummaryLastUpdated{}
		if len(statsSummaries) >= 1 {
			if summary, ok := statsSummaries[0].(tc.StatsSummary); ok {
				resp.SummaryTime = &summary.SummaryTime
			}
		}
		api.WriteResp(w, r, resp)
	}

}

func getStatsSummary(w http.ResponseWriter, r *http.Request, inf *api.Info) {
	queryParamsToSQLCols := map[string]dbhelpers.WhereColumnInfo{
		"statName":            dbhelpers.WhereColumnInfo{Column: "stat_name"},
		"cdnName":             dbhelpers.WhereColumnInfo{Column: "cdn_name"},
		"deliveryServiceName": dbhelpers.WhereColumnInfo{Column: "deliveryservice_name"},
	}
	where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(inf.Params, queryParamsToSQLCols)
	if len(errs) > 0 {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, util.JoinErrs(errs))
		return
	}
	query := selectQuery() + where + orderBy + pagination
	queryStatsSummaries, err := queryStatsSummary(inf.Tx, inf.Version.Major, query, queryValues)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
		return
	}

	//api.WriteResp(w, r, statsSummaries)

	if inf.Version.Major >= 5 {
		statsSummariesV5 := make([]tc.StatsSummaryV5, len(queryStatsSummaries))
		for i, oldStat := range queryStatsSummaries {
			if summary, ok := oldStat.(tc.StatsSummaryV5); ok {
				newStat := tc.StatsSummaryV5{
					CDNName:         summary.CDNName,
					DeliveryService: summary.DeliveryService,
					StatName:        summary.StatName,
					StatValue:       summary.StatValue,
					SummaryTime:     summary.SummaryTime,
					StatDate:        summary.StatDate,
				}
				statsSummariesV5[i] = newStat
			}
		}
		api.WriteResp(w, r, statsSummariesV5)
	} else {
		statsSummaries := make([]tc.StatsSummary, len(queryStatsSummaries))
		for i, oldStat := range queryStatsSummaries {
			if summary, ok := oldStat.(tc.StatsSummary); ok {
				newStat := tc.StatsSummary{
					CDNName:         summary.CDNName,
					DeliveryService: summary.DeliveryService,
					StatName:        summary.StatName,
					StatValue:       summary.StatValue,
					SummaryTime:     summary.SummaryTime,
					StatDate:        summary.StatDate,
				}
				statsSummaries[i] = newStat
			}
		}
		api.WriteResp(w, r, statsSummaries)
	}

}

func queryStatsSummary(tx *sqlx.Tx, version uint64, q string, queryValues map[string]interface{}) ([]interface{}, error) {
	rows, err := tx.NamedQuery(q, queryValues)
	if err != nil {
		return nil, fmt.Errorf("querying stats summary: %v", err)
	}
	defer rows.Close()

	var returnStatsSummaries []interface{}

	if version >= 5 {
		var statsSummariesV5 []tc.StatsSummaryV5
		for rows.Next() {
			s := tc.StatsSummaryV5{}
			if err = rows.StructScan(&s); err != nil {
				return nil, fmt.Errorf("scanning stats summary: %v", err)
			}
			statsSummariesV5 = append(statsSummariesV5, s)
		}
		returnStatsSummaries = make([]interface{}, len(statsSummariesV5))
		for i, v := range statsSummariesV5 {
			returnStatsSummaries[i] = v
		}

	} else {
		var statsSummaries []tc.StatsSummary
		for rows.Next() {
			s := tc.StatsSummary{}
			if err = rows.StructScan(&s); err != nil {
				return nil, fmt.Errorf("scanning stats summary: %v", err)
			}
			statsSummaries = append(statsSummaries, s)
		}

		returnStatsSummaries = make([]interface{}, len(statsSummaries))
		for i, v := range statsSummaries {
			returnStatsSummaries[i] = v
		}
	}
	return returnStatsSummaries, nil

}

// CreateStatsSummary handler for creating stats summaries
func CreateStatsSummary(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{}, []string{})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	ss := tc.StatsSummary{}

	if err := api.Parse(r.Body, inf.Tx.Tx, &ss); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, err, nil)
		return
	}

	// CDN Name and Delivery service name are defaulted to "all" if not defined
	if ss.CDNName == nil || len(*ss.CDNName) == 0 {
		ss.CDNName = util.StrPtr("all")
	}

	if ss.DeliveryService == nil || len(*ss.DeliveryService) == 0 {
		ss.DeliveryService = util.StrPtr("all")
	}

	id := -1
	rows, err := inf.Tx.NamedQuery(insertQuery(), &ss)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("inserting stats summary: %v", err))
		return
	}
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("scanning created stats summary id: %v", err))
			return
		}
	}
	if id == -1 {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("sstats summary id: %v", err))
		return
	}

	successMsg := "Stats Summary was successfully created"
	api.WriteRespAlert(w, r, tc.SuccessLevel, successMsg)
}

func selectQuery() string {
	return `SELECT
cdn_name,
deliveryservice_name,
stat_name,
stat_value,
summary_time,
stat_date
FROM stats_summary`
}

func insertQuery() string {
	return `
INSERT INTO stats_summary (
	cdn_name,
	deliveryservice_name,
	stat_name,
	stat_value,
	summary_time,
	stat_date)
VALUES (
	:cdn_name,
	:deliveryservice_name,
	:stat_name,
	:stat_value,
	:summary_time,
	:stat_date) RETURNING id
`
}
