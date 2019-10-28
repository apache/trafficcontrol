package atscdn

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
	"net/http"
	"strconv"
	"time"

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/ats"
)

const JobKeywordPurge = "PURGE"

func GetRegexRevalidateDotConfig(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"cdn-name-or-id"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	cdnName, userErr, sysErr, errCode := ats.GetCDNNameFromNameOrID(inf.Tx.Tx, inf.Params["cdn-name-or-id"])
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}

	toToolName, toURL, err := ats.GetToolNameAndURL(inf.Tx.Tx)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting tool name and url: "+err.Error()))
		return
	}

	params, err := ats.GetProfileParamsByName(inf.Tx.Tx, tc.GlobalProfileName, atscfg.RegexRevalidateFileName)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting profile params by name: "+err.Error()))
		return
	}

	maxDays, ok, err := getMaxDays(inf.Tx.Tx)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting profile params by name: "+err.Error()))
		return
	}
	if !ok {
		maxDays = atscfg.DefaultMaxRevalDurationDays
		log.Warnf("No maxRevalDurationDays regex_revalidate.config Parameter found, using default %v.\n", maxDays)
	}
	maxReval := time.Duration(maxDays) * time.Hour * 24

	jobs, err := getJobs(inf.Tx.Tx, cdnName, maxReval, atscfg.RegexRevalidateMinTTL)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting jobs: "+err.Error()))
		return
	}

	txt := atscfg.MakeRegexRevalidateDotConfig(tc.CDNName(cdnName), params, toToolName, toURL, jobs)
	w.Header().Set(tc.ContentType, tc.ContentTypeTextPlain)
	w.Write([]byte(txt))
}

// getJobs returns jobs which
//   - have a non-null deliveryservice
//   - have parameters of the form TTL:%dh
//   - have a start time later than (now + maxReval days). That is, we don't query jobs older than maxReval in the past.
//   - are "purge" jobs
//   - have a start_time+ttl > now. That is, jobs that haven't expired yet.
// The maxReval is used for both the max days, for which jobs older than that aren't selected, and for the maximum TTL.
func getJobs(tx *sql.Tx, cdnName string, maxReval time.Duration, minTTL time.Duration) ([]tc.Job, error) {
	qry := `
WITH
  cdn_name AS (select $1::text as v),
  max_days AS (select $2::integer as v)
SELECT
  j.parameters,
  j.keyword,
  j.asset_url,
  u.username,
  j.start_time,
  j.id,
  ds.xml_id
FROM
  job j
  JOIN deliveryservice ds ON j.job_deliveryservice = ds.id
  JOIN tm_user u ON j.job_user = u.id
WHERE
  j.parameters ~ 'TTL:(\d+)h'
  AND j.start_time > (NOW() - ((select v from max_days) * INTERVAL '1 day'))
  AND ds.cdn_id = (select id from cdn where name = (select v from cdn_name))
  AND j.job_deliveryservice IS NOT NULL
  AND j.keyword = '` + atscfg.JobKeywordPurge + `'
  AND (j.start_time + (CAST( (SELECT REGEXP_MATCHES(j.parameters, 'TTL:(\d+)h') FETCH FIRST 1 ROWS ONLY)[1] AS INTEGER) * INTERVAL '1 HOUR')) > NOW()
`

	maxRevalDays := maxReval / time.Hour / 24
	rows, err := tx.Query(qry, cdnName, maxRevalDays)
	if err != nil {
		return nil, errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	jobs := []tc.Job{}
	for rows.Next() {
		j := tc.Job{}
		startTime := time.Time{}
		if err := rows.Scan(&j.Parameters, &j.Keyword, &j.AssetURL, &j.CreatedBy, &startTime, &j.ID, &j.DeliveryService); err != nil {
			return nil, errors.New("scanning: " + err.Error())
		}
		j.StartTime = startTime.Format(tc.JobTimeFormat)
		jobs = append(jobs, j)
	}
	return jobs, nil
}

func getMaxDays(tx *sql.Tx) (int64, bool, error) {
	daysStr := ""
	if err := tx.QueryRow(`SELECT p.value FROM parameter p WHERE p.name = 'maxRevalDurationDays' AND p.config_file = 'regex_revalidate.config'`).Scan(&daysStr); err != nil {
		if err == sql.ErrNoRows {
			return 0, false, nil
		}
		return 0, false, errors.New("querying max reval duration days: " + err.Error())
	}
	days, err := strconv.ParseInt(daysStr, 10, 64)
	if err != nil {
		return 0, false, errors.New("querying max reval duration days: value '" + daysStr + "' is not an integer")
	}
	return days, true, nil
}
