package ats

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
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
)

const DefaultMaxRevalDurationDays = 90
const JobKeywordPurge = "PURGE"

func GetRegexRevalidateDotConfig(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"cdn-name-or-id"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	cdnName, userErr, sysErr, errCode := getCDNNameFromNameOrID(inf.Tx.Tx, inf.Params["cdn-name-or-id"])
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}

	regexRevalTxt, err := getRegexRevalidate(inf.Tx.Tx, cdnName)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting regex_revalidate.config text: "+err.Error()))
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(regexRevalTxt))
}

func getRegexRevalidate(tx *sql.Tx, cdnName string) (string, error) {
	maxDays, ok, err := getMaxDays(tx)
	if err != nil {
		return "", errors.New("getting max reval duration days from Parameter: " + err.Error())
	}
	if !ok {
		maxDays = DefaultMaxRevalDurationDays
		log.Warnf("No maxRevalDurationDays regex_revalidate.config Parameter found, using default %v.\n", maxDays)
	}
	maxReval := time.Duration(maxDays) * time.Hour * 24
	minTTL := time.Hour * 1

	jobs, err := getJobs(tx, cdnName, maxReval, minTTL)
	if err != nil {
		return "", errors.New("getting jobs: " + err.Error())
	}

	text, err := headerComment(tx, "CDN "+cdnName)
	if err != nil {
		return "", errors.New("getting header comment: " + err.Error())
	}
	for _, job := range jobs {
		text += job.AssetURL + " " + strconv.FormatInt(job.PurgeEnd.Unix(), 10) + "\n"
	}

	return text, nil
}

type Job struct {
	AssetURL string
	PurgeEnd time.Time
}

type Jobs []Job

func (jb Jobs) Len() int      { return len(jb) }
func (jb Jobs) Swap(i, j int) { jb[i], jb[j] = jb[j], jb[i] }
func (jb Jobs) Less(i, j int) bool {
	if jb[i].AssetURL == jb[j].AssetURL {
		return jb[i].PurgeEnd.Before(jb[j].PurgeEnd)
	}
	return strings.Compare(jb[i].AssetURL, jb[j].AssetURL) < 0
}

// getJobs returns jobs which
//   - have a non-null deliveryservice
//   - have parameters of the form TTL:%dh
//   - have a start time later than (now + maxReval days). That is, we don't query jobs older than maxReval in the past.
//   - are "purge" jobs
//   - have a start_time+ttl > now. That is, jobs that haven't expired yet.
// The maxReval is used for both the max days, for which jobs older than that aren't selected, and for the maximum TTL.
func getJobs(tx *sql.Tx, cdnName string, maxReval time.Duration, minTTL time.Duration) ([]Job, error) {
	qry := `
WITH
  cdn_name AS (select $1::text as v),
  max_days AS (select $2::integer as v)
SELECT
  j.asset_url,
  CAST((SELECT REGEXP_MATCHES(j.parameters, 'TTL:(\d+)h') FETCH FIRST 1 ROWS ONLY)[1] AS INTEGER) as ttl,
  j.start_time
FROM
  job j
  JOIN deliveryservice ds ON j.job_deliveryservice = ds.id
WHERE
  j.parameters ~ 'TTL:(\d+)h'
  AND j.start_time > (NOW() - ((select v from max_days) * INTERVAL '1 day'))
  AND ds.cdn_id = (select id from cdn where name = (select v from cdn_name))
  AND j.job_deliveryservice IS NOT NULL
  AND j.keyword = '` + JobKeywordPurge + `'
  AND (j.start_time + (CAST( (SELECT REGEXP_MATCHES(j.parameters, 'TTL:(\d+)h') FETCH FIRST 1 ROWS ONLY)[1] AS INTEGER) * INTERVAL '1 HOUR')) > NOW()
`
	maxRevalDays := maxReval / time.Hour / 24
	rows, err := tx.Query(qry, cdnName, maxRevalDays)
	if err != nil {
		return nil, errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	jobMap := map[string]time.Time{}
	for rows.Next() {
		assetURL := ""
		ttlHours := 0
		startTime := time.Time{}
		if err := rows.Scan(&assetURL, &ttlHours, &startTime); err != nil {
			return nil, errors.New("scanning: " + err.Error())
		}

		ttl := time.Duration(ttlHours) * time.Hour
		if ttl > maxReval {
			ttl = maxReval
		} else if ttl < minTTL {
			ttl = minTTL
		}

		purgeEnd := startTime.Add(ttl)

		if existingPurgeEnd, ok := jobMap[assetURL]; !ok || purgeEnd.After(existingPurgeEnd) {
			jobMap[assetURL] = purgeEnd
		}
	}

	jobs := []Job{}
	for assetURL, purgeEnd := range jobMap {
		jobs = append(jobs, Job{AssetURL: assetURL, PurgeEnd: purgeEnd})
	}
	sort.Sort(Jobs(jobs))
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

func headerComment(tx *sql.Tx, name string) (string, error) {
	nameVersionStr, err := GetNameVersionString(tx)
	if err != nil {
		return "", errors.New("getting name version string: " + err.Error())
	}
	return "# DO NOT EDIT - Generated for " + name + " by " + nameVersionStr + " on " + time.Now().Format(HeaderCommentDateFormat) + "\n", nil
}
