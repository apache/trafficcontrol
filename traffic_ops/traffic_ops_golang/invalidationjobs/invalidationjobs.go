package invalidationjobs

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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/tenant"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/util/ims"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/lib/pq"
)

// Deprecated, only to be used with versions below 4.0
type InvalidationJob struct {
	api.APIInfoImpl `json:"-"`
	tc.InvalidationJob
}

type InvalidationJobV4 struct {
	api.APIInfoImpl `json:"-"`
	tc.InvalidationJobV4
}

// Deprecated, only to be used with versions below 4.0
const insertQuery = `
INSERT INTO job (
	ttl_hr,
	asset_url,
	start_time,
	entered_time,
	job_user,
	job_deliveryservice,
	invalidation_type)
VALUES (
	$1,
	(
		SELECT o.protocol::text || '://' || o.fqdn || rtrim(concat(':', o.port::text), ':')
		FROM origin o
		WHERE o.deliveryservice = $2
		AND o.is_primary
	) || $3,
	$4,
	$5,
	$6,
	$7,
	$8
)
RETURNING
	asset_url,
	(
		SELECT deliveryservice.xml_id
		FROM deliveryservice
		WHERE deliveryservice.id=job_deliveryservice) AS deliveryservice,
	id,
	(
		SELECT tm_user.username
		FROM tm_user
		WHERE tm_user.id=job_user) AS createdBy,
	'PURGE' AS keyword,
	CONCAT('TTL:', ttl_hr, 'h') AS parameters,
	start_time
`

// Almost the same as insertQuery, but returns appropriate values for API 4.0+
const insertQueryV4 = `
INSERT INTO job (
	ttl_hr,
	asset_url,
	start_time,
	entered_time,
	job_user,
	job_deliveryservice,
	invalidation_type)
VALUES (
	$1,
	(
		SELECT o.protocol::text || '://' || o.fqdn || rtrim(concat(':', o.port::text), ':')
		FROM origin o
		WHERE o.deliveryservice = $2
		AND o.is_primary
	) || $3,
	$4,
	$5,
	$6,
	$7,
	$8
)
RETURNING
	id,
	asset_url,
	(
		SELECT tm_user.username
		FROM tm_user
		WHERE tm_user.id=job_user) AS createdBy,
	(
		SELECT deliveryservice.xml_id
		FROM deliveryservice
		WHERE deliveryservice.id=job_deliveryservice) AS deliveryServiceXML,
	ttl_hr as ttlHrs,
	invalidation_type as invalidationType,
	start_time as startTime
`

const queueUpdateOrRevalQuery = `
UPDATE public.server
SET %s = now()
WHERE server.status IN (
		SELECT status.id
		FROM status
		WHERE name IN ('ONLINE', 'REPORTED', 'ADMIN_DOWN')
		)
     AND server.profile IN (
		SELECT profile_parameter.profile
		FROM profile_parameter
		WHERE profile_parameter.parameter IN (
			SELECT parameter.id
			FROM parameter
			WHERE parameter.name='location'
			AND parameter.config_file='regex_revalidate.config'
			)
		)
     AND server.cdn_id  =  (
		SELECT deliveryservice.cdn_id
		FROM deliveryservice
		WHERE deliveryservice.%s=$1
		);
`

const updateQuery = `
UPDATE job
SET asset_url=$1,
    ttl_hr=$2,
    start_time=$3
WHERE job.id=$4
RETURNING asset_url,
	(
		SELECT tm_user.username
		FROM tm_user
		WHERE tm_user.id=job.job_user
	) AS created_by,
	(
		SELECT deliveryservice.xml_id
		FROM deliveryservice
		WHERE deliveryservice.id=job.job_deliveryservice
	) AS delivery_service,
	job.id,
	'PURGE' as keyword,
	CONCAT('TTL:', ttl_hr, 'h') AS parameters,
	start_time
`

// Almost the same as updateQuery, but returns appropriate values for API 4.0+
const updateQueryV4 = `
UPDATE job
SET asset_url=$1,
	ttl_hr=$2,
	start_time=$3,
	invalidation_type=$4
WHERE job.id=$5
RETURNING asset_url,
	(
		SELECT tm_user.username
		FROM tm_user
		WHERE tm_user.id=job.job_user
	) AS created_by,
	(
		SELECT deliveryservice.xml_id
		FROM deliveryservice
		WHERE deliveryservice.id=job.job_deliveryservice
	) AS delivery_service,
	job.id,
	ttl_hr,
	start_time,
	invalidation_type
`

// Deprecated, only to be used with versions below 4.0
const putInfoQuery = `
SELECT job.id AS id,
	tm_user.username AS createdBy,
	job.job_user AS createdByID,
	job.job_deliveryservice AS dsid,
	deliveryservice.xml_id AS dsxmlid,
	job.asset_url AS assetURL,
	CONCAT('TTL:', ttl_hr, 'h') AS parameters,
	job.start_time AS start_time,
	origin.protocol || '://' || origin.fqdn || rtrim(concat(':', origin.port), ':') AS OFQDN
FROM job
INNER JOIN origin ON origin.deliveryservice=job.job_deliveryservice AND origin.is_primary
INNER JOIN tm_user ON tm_user.id=job.job_user
INNER JOIN deliveryservice ON deliveryservice.id=job.job_deliveryservice
WHERE job.id=$1
`

// Almost the same as putInfoQuery, but returns appropriate values for API 4.0+
const putInfoQueryV4 = `
SELECT job.id AS id,
	tm_user.username AS createdBy,
	job.job_user AS createdByID,
	job.job_deliveryservice AS dsid,
	deliveryservice.xml_id AS dsxmlid,
	job.asset_url AS assetURL,
	job.ttl_hr AS ttlhrs,
	job.start_time AS start_time,
	job.invalidation_type as invalidationType,
	origin.protocol || '://' || origin.fqdn || rtrim(concat(':', origin.port), ':') AS OFQDN
FROM job
INNER JOIN origin ON origin.deliveryservice=job.job_deliveryservice AND origin.is_primary
INNER JOIN tm_user ON tm_user.id=job.job_user
INNER JOIN deliveryservice ON deliveryservice.id=job.job_deliveryservice
WHERE job.id=$1
`

// Deprecated, only to be used with versions below 4.0
const deleteQuery = `
DELETE
FROM job
WHERE job.id=$1
RETURNING job.asset_url,
	(
		SELECT tm_user.username
		FROM tm_user
		WHERE tm_user.id=job.job_user
	) AS created_by,
	(
		SELECT deliveryservice.xml_id
		FROM deliveryservice
		WHERE deliveryservice.id=job.job_deliveryservice
	) AS deliveryservice,
	job.id,
	'PURGE' as keyword,
	CONCAT('TTL:', ttl_hr, 'h') AS parameters,
	job.start_time
`

// Almost the same as deleteQuery, but returns appropriate values for API 4.0+
const deleteQueryV4 = `
DELETE
FROM job
WHERE job.id=$1
RETURNING 
	job.id,
	job.asset_url,
	(
		SELECT tm_user.username
		FROM tm_user
		WHERE tm_user.id=job.job_user
	) AS created_by,
	(
		SELECT deliveryservice.xml_id
		FROM deliveryservice
		WHERE deliveryservice.id=job.job_deliveryservice
	) AS deliveryservice,
	ttl_hr,
	job.invalidation_type,
	job.start_time
`

type apiResponse struct {
	Alerts   []tc.Alert         `json:"alerts,omitempty"`
	Response tc.InvalidationJob `json:"response,omitempty"`
}

type apiResponseV4 struct {
	Alerts   []tc.Alert           `json:"alerts,omitempty"`
	Response tc.InvalidationJobV4 `json:"response,omitempty"`
}

func selectMaxLastUpdatedQuery(where string) string {
	return `SELECT max(t) from (
		SELECT max(job.last_updated) as t FROM job
	JOIN tm_user u ON job.job_user = u.id
	JOIN deliveryservice ds  ON job.job_deliveryservice = ds.id ` + where +
		` UNION ALL
	select max(last_updated) as t from last_deleted l where l.table_name='job') as res`
}

// Deprecated, only to be used with versions below 4.0
const readQuery = `
SELECT job.id,
	'PURGE' AS keyword,
	CONCAT('TTL::', ttl_hr, 'h') AS parameters,
	asset_url,
	start_time,
	u.username as createdBy,
	ds.xml_id as dsId
FROM job
JOIN tm_user u ON job.job_user = u.id
JOIN deliveryservice ds ON job.job_deliveryservice = ds.id
`

// Almost the same as readQuery, but returns appropriate values for API 4.0+
const readQueryV4 = `
SELECT job.id,
	asset_url,
	u.username as createdBy,
	ds.xml_id,
	ttl_hr,
	invalidation_type,
	start_time
FROM job
JOIN tm_user u ON job.job_user = u.id
JOIN deliveryservice ds ON job.job_deliveryservice = ds.id
`

// Used by GET requests to `/jobs`, simply returns a filtered list of
// content invalidation jobs according to the provided query parameters.
func (job *InvalidationJobV4) Read(h http.Header, useIMS bool) ([]interface{}, error, error, int, *time.Time) {
	var maxTime time.Time
	var runSecond bool
	queryParamsToSQLCols := map[string]dbhelpers.WhereColumnInfo{
		"id":               dbhelpers.WhereColumnInfo{Column: "job.id", Checker: api.IsInt},
		"assetUrl":         dbhelpers.WhereColumnInfo{Column: "asset_url"},
		"startTime":        dbhelpers.WhereColumnInfo{Column: "start_time"},
		"userId":           dbhelpers.WhereColumnInfo{Column: "job_user", Checker: api.IsInt},
		"createdBy":        dbhelpers.WhereColumnInfo{Column: `(SELECT tm_user.username FROM tm_user WHERE tm_user.id=job.job_user)`},
		"deliveryService":  dbhelpers.WhereColumnInfo{Column: `(SELECT deliveryservice.xml_id FROM deliveryservice WHERE deliveryservice.id=job.job_deliveryservice)`},
		"dsId":             dbhelpers.WhereColumnInfo{Column: "job.job_deliveryservice", Checker: api.IsInt},
		"invalidationType": dbhelpers.WhereColumnInfo{Column: "invalidation_type"},
	}

	where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(job.APIInfo().Params, queryParamsToSQLCols)
	if len(errs) > 0 {
		return nil, util.JoinErrs(errs), nil, http.StatusBadRequest, nil
	}

	accessibleTenants, err := tenant.GetUserTenantIDListTx(job.APIInfo().Tx.Tx, job.APIInfo().User.TenantID)
	if err != nil {
		return nil, nil, fmt.Errorf("getting accessible tenants for user - %v", err), http.StatusInternalServerError, nil
	}
	cdn := ""
	if cdnName, ok := job.APIInfo().Params["cdn"]; ok {
		queryValues["cdn"] = cdnName
		cdn = ` AND ds.cdn_id = (SELECT id FROM cdn WHERE name = :cdn) `
	}
	maxDays := ""
	if _, ok := job.APIInfo().Params["maxRevalDurationDays"]; ok {
		// jobs started within the last $maxRevalDurationDays days (defaulting to 90 days if the parameter doesn't exist)
		maxDays = ` AND job.start_time >= NOW() - CAST(
                                                       (SELECT COALESCE(
                                                                       (SELECT value
                                                                        FROM parameter
                                                                        WHERE name = 'maxRevalDurationDays'
                                                                        AND config_file = 'regex_revalidate.config'
                                                                        LIMIT 1),
                                                                       '90'))
                                                       || ' days' AS INTERVAL) `
	}
	if len(where) > 0 {
		where += " AND ds.tenant_id = ANY(:tenants) " + maxDays + cdn
	} else {
		where = dbhelpers.BaseWhere + " ds.tenant_id = ANY(:tenants) " + maxDays + cdn
	}
	queryValues["tenants"] = pq.Array(accessibleTenants)

	if useIMS {
		runSecond, maxTime = ims.TryIfModifiedSinceQuery(job.APIInfo().Tx, h, queryValues, selectMaxLastUpdatedQuery(where))
		if !runSecond {
			log.Debugln("IMS HIT")
			return []interface{}{}, nil, nil, http.StatusNotModified, &maxTime
		}
		log.Debugln("IMS MISS")
	} else {
		log.Debugln("Non IMS request")
	}

	query := readQueryV4 + where + orderBy + pagination
	log.Debugln("generated job query: " + query)
	log.Debugf("executing with values: %++v\n", queryValues)

	returnable := []interface{}{}
	rows, err := job.APIInfo().Tx.NamedQuery(query, queryValues)
	if err != nil {
		return nil, nil, fmt.Errorf("querying: %v", err), http.StatusInternalServerError, nil
	}
	defer rows.Close()

	for rows.Next() {
		job := tc.InvalidationJobV4{}
		if err := rows.Scan(&job.ID,
			&job.AssetURL,
			&job.CreatedBy,
			&job.DeliveryService,
			&job.TTLHours,
			&job.InvalidationType,
			&job.StartTime); err != nil {
			return nil, nil, fmt.Errorf("parsing db response: %v", err), http.StatusInternalServerError, nil
		}

		returnable = append(returnable, job)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("Parsing db responses: %v", err), http.StatusInternalServerError, nil
	}

	return returnable, nil, nil, http.StatusOK, &maxTime
}

// Used by GET requests to `/jobs`, simply returns a filtered list of
// content invalidation jobs according to the provided query parameters.
//
// Deprecated. To be used only with versions less than 4.0
func (job *InvalidationJob) Read(h http.Header, useIMS bool) ([]interface{}, error, error, int, *time.Time) {
	var maxTime time.Time
	var runSecond bool
	queryParamsToSQLCols := map[string]dbhelpers.WhereColumnInfo{
		"id":              dbhelpers.WhereColumnInfo{Column: "job.id", Checker: api.IsInt},
		"keyword":         dbhelpers.WhereColumnInfo{Column: "keyword"},
		"assetUrl":        dbhelpers.WhereColumnInfo{Column: "asset_url"},
		"startTime":       dbhelpers.WhereColumnInfo{Column: "start_time"},
		"userId":          dbhelpers.WhereColumnInfo{Column: "job_user", Checker: api.IsInt},
		"createdBy":       dbhelpers.WhereColumnInfo{Column: `(SELECT tm_user.username FROM tm_user WHERE tm_user.id=job.job_user)`},
		"deliveryService": dbhelpers.WhereColumnInfo{Column: `(SELECT deliveryservice.xml_id FROM deliveryservice WHERE deliveryservice.id=job.job_deliveryservice)`},
		"dsId":            dbhelpers.WhereColumnInfo{Column: "job.job_deliveryservice", Checker: api.IsInt},
	}

	where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(job.APIInfo().Params, queryParamsToSQLCols)
	if len(errs) > 0 {
		return nil, util.JoinErrs(errs), nil, http.StatusBadRequest, nil
	}

	accessibleTenants, err := tenant.GetUserTenantIDListTx(job.APIInfo().Tx.Tx, job.APIInfo().User.TenantID)
	if err != nil {
		return nil, nil, fmt.Errorf("getting accessible tenants for user - %v", err), http.StatusInternalServerError, nil
	}
	cdn := ""
	if cdnName, ok := job.APIInfo().Params["cdn"]; ok {
		queryValues["cdn"] = cdnName
		cdn = ` AND ds.cdn_id = (SELECT id FROM cdn WHERE name = :cdn) `
	}
	maxDays := ""
	if _, ok := job.APIInfo().Params["maxRevalDurationDays"]; ok {
		// jobs started within the last $maxRevalDurationDays days (defaulting to 90 days if the parameter doesn't exist)
		maxDays = ` AND job.start_time >= NOW() - CAST(
                                                       (SELECT COALESCE(
                                                                       (SELECT value
                                                                        FROM parameter
                                                                        WHERE name = 'maxRevalDurationDays'
                                                                        AND config_file = 'regex_revalidate.config'
                                                                        LIMIT 1),
                                                                       '90'))
                                                       || ' days' AS INTERVAL) `
	}
	if len(where) > 0 {
		where += " AND ds.tenant_id = ANY(:tenants) " + maxDays + cdn
	} else {
		where = dbhelpers.BaseWhere + " ds.tenant_id = ANY(:tenants) " + maxDays + cdn
	}
	queryValues["tenants"] = pq.Array(accessibleTenants)

	if useIMS {
		runSecond, maxTime = ims.TryIfModifiedSinceQuery(job.APIInfo().Tx, h, queryValues, selectMaxLastUpdatedQuery(where))
		if !runSecond {
			log.Debugln("IMS HIT")
			return []interface{}{}, nil, nil, http.StatusNotModified, &maxTime
		}
		log.Debugln("IMS MISS")
	} else {
		log.Debugln("Non IMS request")
	}

	query := readQuery + where + orderBy + pagination
	log.Debugln("generated job query: " + query)
	log.Debugf("executing with values: %++v\n", queryValues)

	returnable := []interface{}{}
	rows, err := job.APIInfo().Tx.NamedQuery(query, queryValues)
	if err != nil {
		return nil, nil, fmt.Errorf("querying: %v", err), http.StatusInternalServerError, nil
	}
	defer rows.Close()

	for rows.Next() {
		j := tc.InvalidationJob{}
		err := rows.Scan(&j.ID,
			&j.Keyword,
			&j.Parameters,
			&j.AssetURL,
			&j.StartTime,
			&j.CreatedBy,
			&j.DeliveryService)
		if err != nil {
			return nil, nil, fmt.Errorf("parsing db response: %v", err), http.StatusInternalServerError, nil
		}

		returnable = append(returnable, j)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("Parsing db responses: %v", err), http.StatusInternalServerError, nil
	}

	return returnable, nil, nil, http.StatusOK, &maxTime
}

// Used by POST requests to `/jobs`, creates a new content invalidation job
// from the provided request body.
func CreateV40(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	job := tc.InvalidationJobCreateV4{}
	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("Unable to parse Invalidation Job"), fmt.Errorf("parsing jobs/ POST: %v", err))
		return
	}

	// Check if request object is valid
	w.Header().Set(rfc.ContentType, rfc.ApplicationJSON)
	if err := validateJobCreateV4(job, inf.Tx.Tx); err != nil {
		response := tc.Alerts{
			Alerts: []tc.Alert{
				{
					Text:  err.Error(),
					Level: tc.ErrorLevel.String(),
				},
			},
		}

		resp, err := json.Marshal(response)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("encoding bad request response: %v", err))
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		api.WriteAndLogErr(w, r, append(resp, '\n'))
		return
	}

	// Check if authorized
	if ok, err := IsUserAuthorizedToModifyDSXMLID(inf, job.DeliveryService); err != nil {
		sysErr = fmt.Errorf("failed checking current user permissions for DS #%s: %v", job.DeliveryService, err)
		errCode = http.StatusInternalServerError
		api.HandleErr(w, r, inf.Tx.Tx, errCode, nil, sysErr)
		return
	} else if !ok {
		userErr = fmt.Errorf("failed to authorize based on tenancy")
		errCode = http.StatusNotFound
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, nil)
		return
	}

	// DS existence was already verified in the Validate() function
	dsid, exists, err := dbhelpers.GetDSIDFromXMLID(inf.Tx.Tx, job.DeliveryService)
	if err != nil {
		sysErr = fmt.Errorf("failed to match XML ID to int ID for Delivery Service %s: %v", job.DeliveryService, err)
		errCode = http.StatusInternalServerError
		api.HandleErr(w, r, inf.Tx.Tx, errCode, nil, sysErr)
		return
	}
	if !exists {
		userErr = fmt.Errorf("delivery service \"%v\" does not exist", job.DeliveryService)
		errCode = http.StatusNotFound
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, nil)
		return
	}

	row := inf.Tx.Tx.QueryRow(insertQueryV4,
		job.TTLHours,
		dsid, // Used in inner select for deliveryservice
		job.Regex,
		job.StartTime,
		time.Now(),
		inf.User.ID,
		dsid,
		job.InvalidationType) // Defaults for all api versions below 4.0

	result := tc.InvalidationJobV4{}
	err = row.Scan(
		&result.ID,
		&result.AssetURL,
		&result.CreatedBy,
		&result.DeliveryService,
		&result.TTLHours,
		&result.InvalidationType,
		&result.StartTime)
	if err != nil {
		userErr, sysErr, errCode = api.ParseDBError(err)
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}

	if err := setRevalFlags(uint(dsid), inf.Tx.Tx); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("setting reval flags: %v", err))
		return
	}

	conflicts := tc.ValidateJobUniqueness(inf.Tx.Tx, uint(dsid), result.StartTime, result.AssetURL, result.TTLHours)
	response := apiResponseV4{
		make([]tc.Alert, len(conflicts)+1),
		result,
	}
	for i, conflict := range conflicts {
		response.Alerts[i] = tc.Alert{
			Text:  conflict,
			Level: tc.WarnLevel.String(),
		}
	}
	response.Alerts[len(conflicts)] = tc.Alert{
		Text: fmt.Sprintf("Invalidation (%s) request created for %v, start:%v end %v",
			result.InvalidationType,
			result.AssetURL,
			result.StartTime,
			result.StartTime.Add(time.Hour*time.Duration(job.TTLHours))),
		Level: tc.SuccessLevel.String(),
	}
	resp, err := json.Marshal(response)

	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("Marshaling JSON: %v", err))
		return
	}

	if inf.Version == nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("nil API version"))
		return
	}

	w.Header().Set(http.CanonicalHeaderKey("location"), fmt.Sprintf("%s://%s/api/%d.%d/jobs?id=%d",
		inf.Config.URL.Scheme,
		r.Host,
		inf.Version.Major,
		inf.Version.Minor,
		result.ID))
	w.WriteHeader(http.StatusOK)
	api.WriteAndLogErr(w, r, append(resp, '\n'))

	duplicate := ""
	if len(conflicts) > 0 {
		duplicate = "(duplicate) "
	}
	changeLogMsg := fmt.Sprintf("%s content invalidation job %s- ID: %d DSXMLID: %s ASSET_URL: '%s' TTLHRs: %d INVALIDATION: %s",
		api.Created,
		duplicate,
		result.ID,
		result.DeliveryService,
		result.AssetURL,
		result.TTLHours,
		result.InvalidationType,
	)
	api.CreateChangeLogRawTx(api.ApiChange,
		changeLogMsg,
		inf.User,
		inf.Tx.Tx)
}

// Used by POST requests to `/jobs`, creates a new content invalidation job
// from the provided request body.
//
// Deprecated. To be used only with versions less than 4.0
func Create(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	job := tc.InvalidationJobInput{}
	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("Unable to parse Invalidation Job"), fmt.Errorf("parsing jobs/ POST: %v", err))
		return
	}

	w.Header().Set(rfc.ContentType, rfc.ApplicationJSON)
	if err := job.Validate(inf.Tx.Tx); err != nil {
		response := tc.Alerts{
			Alerts: []tc.Alert{
				tc.Alert{
					Text:  err.Error(),
					Level: tc.ErrorLevel.String(),
				},
			},
		}

		resp, err := json.Marshal(response)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("Encoding bad request response: %v", err))
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		api.WriteAndLogErr(w, r, append(resp, '\n'))
		return
	}

	// Validate() would have already checked for deliveryservice existence and
	// parsed the ttl, so if either of these throws an error now, something
	// weird has happened
	dsid, err := job.DSID(nil)
	if err != nil {
		sysErr = fmt.Errorf("retrieving parsed DSID: %v", err)
		errCode = http.StatusInternalServerError
		api.HandleErr(w, r, inf.Tx.Tx, errCode, nil, sysErr)
		return
	}
	var ttl uint
	if ttl, err = job.TTLHours(); err != nil {
		sysErr = fmt.Errorf("retrieving parsed TTL: %v", err)
		errCode = http.StatusInternalServerError
		api.HandleErr(w, r, inf.Tx.Tx, errCode, nil, sysErr)
		return
	}

	if ok, err := IsUserAuthorizedToModifyDSID(inf, dsid); err != nil {
		sysErr = fmt.Errorf("Checking current user permissions for DS #%d: %v", dsid, err)
		errCode = http.StatusInternalServerError
		api.HandleErr(w, r, inf.Tx.Tx, errCode, nil, sysErr)
		return
	} else if !ok {
		userErr = fmt.Errorf("No such Delivery Service!")
		errCode = http.StatusNotFound
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, nil)
		return
	}

	_, cdnName, _, err := dbhelpers.GetDSNameAndCDNFromID(inf.Tx.Tx, int(dsid))
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting delivery service and CDN name from ID: "+err.Error()))
		return
	}
	userErr, sysErr, statusCode := dbhelpers.CheckIfCurrentUserCanModifyCDN(inf.Tx.Tx, string(cdnName), inf.User.UserName)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, statusCode, userErr, sysErr)
		return
	}
	row := inf.Tx.Tx.QueryRow(insertQuery,
		ttl,
		dsid, // Used in inner select for deliveryservice
		*job.Regex,
		(*job.StartTime).Time,
		time.Now(),
		inf.User.ID,
		dsid,
		tc.REFRESH) // Defaults for all api versions below 4.0

	result := tc.InvalidationJob{}
	err = row.Scan(&result.AssetURL,
		&result.DeliveryService,
		&result.ID,
		&result.CreatedBy,
		&result.Keyword,
		&result.Parameters,
		&result.StartTime)
	if err != nil {
		userErr, sysErr, errCode = api.ParseDBError(err)
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}

	if err := setRevalFlags(dsid, inf.Tx.Tx); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("setting reval flags: %v", err))
		return
	}

	conflicts := tc.ValidateJobUniqueness(inf.Tx.Tx, dsid, job.StartTime.Time, *result.AssetURL, ttl)
	response := apiResponse{
		make([]tc.Alert, len(conflicts)+1),
		result,
	}
	for i, conflict := range conflicts {
		response.Alerts[i] = tc.Alert{
			Text:  conflict,
			Level: tc.WarnLevel.String(),
		}
	}
	response.Alerts[len(conflicts)] = tc.Alert{
		Text: fmt.Sprintf("Invalidation request created for %v, start:%v end %v", *result.AssetURL, job.StartTime.Time,
			job.StartTime.Add(time.Hour*time.Duration(ttl))),
		Level: tc.SuccessLevel.String(),
	}
	resp, err := json.Marshal(response)

	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("Marshaling JSON: %v", err))
		return
	}

	if inf.Version == nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("nil API version"))
		return
	}

	w.Header().Set(http.CanonicalHeaderKey("location"), fmt.Sprintf("%s://%s/api/%d.%d/jobs?id=%d", inf.Config.URL.Scheme, r.Host, inf.Version.Major, inf.Version.Minor, *result.ID))
	w.WriteHeader(http.StatusOK)
	api.WriteAndLogErr(w, r, append(resp, '\n'))

	duplicate := ""
	if len(conflicts) > 0 {
		duplicate = "(duplicate) "
	}
	api.CreateChangeLogRawTx(api.ApiChange, api.Created+" content invalidation job "+duplicate+"- ID: "+
		strconv.FormatUint(*result.ID, 10)+" DS: "+*result.DeliveryService+" URL: '"+*result.AssetURL+
		"' Params: '"+*result.Parameters+"'", inf.User, inf.Tx.Tx)
}

// Used by PUT requests to `/jobs`, replaces an existing content invalidation job
// with the provided request body.
func UpdateV40(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id"}, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	var oFQDN string
	var dsid uint
	var uid uint
	job := tc.InvalidationJobV4{}
	row := inf.Tx.Tx.QueryRow(putInfoQueryV4, inf.Params["id"])
	err := row.Scan(&job.ID,
		&job.CreatedBy,
		&uid,
		&dsid,
		&job.DeliveryService,
		&job.AssetURL,
		&job.TTLHours,
		&job.StartTime,
		&job.InvalidationType,
		&oFQDN)
	if err != nil {
		if err == sql.ErrNoRows {
			userErr = fmt.Errorf("No job by id '%s'!", inf.Params["id"])
			errCode = http.StatusNotFound
		} else {
			sysErr = fmt.Errorf("fetching job update info: %v", err)
			errCode = http.StatusInternalServerError
		}
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}

	if ok, err := IsUserAuthorizedToModifyDSID(inf, dsid); err != nil {
		sysErr = fmt.Errorf("Checking user permissions on DS #%d: %v", dsid, err)
		errCode = http.StatusInternalServerError
		api.HandleErr(w, r, inf.Tx.Tx, errCode, nil, sysErr)
		return
	} else if !ok {
		userErr = errors.New("No such Delivery Service!")
		errCode = http.StatusNotFound
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, nil)
		return
	}

	if ok, err := IsUserAuthorizedToModifyJobsMadeByUsername(inf, job.CreatedBy); err != nil {
		sysErr = fmt.Errorf("Checking user permissions against user %s: %v", job.CreatedBy, err)
		errCode = http.StatusInternalServerError
		api.HandleErr(w, r, inf.Tx.Tx, errCode, nil, sysErr)
		return
	} else if !ok {
		userErr = fmt.Errorf("No job by id '%s'!", inf.Params["id"])
		errCode = http.StatusNotFound
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, nil)
		return
	}

	input := tc.InvalidationJobV4{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		userErr = fmt.Errorf("Unable to parse input: %v", err)
		sysErr = fmt.Errorf("parsing input to PUT jobs?id=%s: %v", inf.Params["id"], err)
		errCode = http.StatusBadRequest
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}

	if err := validateInvalidationJobV4(input); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, err, nil)
		return
	}

	if !strings.HasPrefix(input.AssetURL, oFQDN) {
		userErr = fmt.Errorf("Cannot set asset URL that does not start with Delivery Service origin URL: %s", oFQDN)
		errCode = http.StatusBadRequest
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, nil)
		return
	}

	if job.StartTime.Before(time.Now()) {
		userErr = errors.New("Cannot modify a job that has already started!")
		errCode = http.StatusMethodNotAllowed
		w.Header().Set(http.CanonicalHeaderKey("allow"), "GET,HEAD,DELETE")
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, nil)
		return
	}

	if job.DeliveryService != input.DeliveryService {
		userErr = errors.New("Cannot change 'deliveryService' of existing invalidation job!")
		errCode = http.StatusConflict
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, nil)
		return
	}

	if job.CreatedBy != input.CreatedBy {
		userErr = errors.New("Cannot change 'createdBy' of existing invalidation jobs!")
		errCode = http.StatusConflict
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, nil)
		return
	}

	if job.ID != input.ID {
		userErr = errors.New("Cannot change an invalidation job 'id'!")
		errCode = http.StatusConflict
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, nil)
		return
	}

	if job.InvalidationType != input.InvalidationType {
		if input.InvalidationType == tc.REFETCH && !refetchAllowed(inf.Tx.Tx) {
			userErr = errors.New("Invalidation Type REFRESH is disallowed")
			errCode = http.StatusBadRequest
			api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, nil)
			return
		}
	}

	row = inf.Tx.Tx.QueryRow(updateQueryV4,
		input.AssetURL,
		input.TTLHours,
		input.StartTime,
		input.InvalidationType,
		job.ID)
	err = row.Scan(&job.AssetURL,
		&job.CreatedBy,
		&job.DeliveryService,
		&job.ID,
		&job.TTLHours,
		&job.StartTime,
		&job.InvalidationType)
	if err != nil {
		sysErr = fmt.Errorf("Updating a job: %v", err)
		errCode = http.StatusInternalServerError
		api.HandleErr(w, r, inf.Tx.Tx, errCode, nil, sysErr)
		return
	}

	if err = setRevalFlags(job.DeliveryService, inf.Tx.Tx); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("Setting reval flags: %v", err))
		return
	}

	conflicts := tc.ValidateJobUniqueness(inf.Tx.Tx, dsid, input.StartTime, input.AssetURL, input.TTLHours)
	response := apiResponseV4{
		make([]tc.Alert, len(conflicts)+1),
		job,
	}
	for i, conflict := range conflicts {
		response.Alerts[i] = tc.Alert{
			Text:  conflict,
			Level: tc.WarnLevel.String(),
		}
	}
	response.Alerts[len(conflicts)] = tc.Alert{
		Text: fmt.Sprintf("Invalidation request created for %s, start: %v end: %v invalidation type: %v",
			job.AssetURL,
			job.StartTime,
			job.StartTime.Add(time.Hour*time.Duration(job.TTLHours)),
			job.InvalidationType),
		Level: tc.SuccessLevel.String(),
	}

	resp, err := json.Marshal(response)
	if err != nil {
		sysErr = fmt.Errorf("encoding response: %v", err)
		errCode = http.StatusInternalServerError
		api.HandleErr(w, r, inf.Tx.Tx, errCode, nil, sysErr)
		return
	}

	w.Header().Set(http.CanonicalHeaderKey("content-type"), rfc.ApplicationJSON)
	api.WriteAndLogErr(w, r, append(resp, '\n'))

	changeLogMsg := fmt.Sprintf("%s content invalidation job - ID: %d DSXMLID: %s ASSET_URL: '%s' TTLHRs: %d INVALIDATION: %s",
		api.Updated,
		input.ID,
		input.DeliveryService,
		input.AssetURL,
		input.TTLHours,
		input.InvalidationType,
	)
	api.CreateChangeLogRawTx(api.ApiChange,
		changeLogMsg,
		inf.User,
		inf.Tx.Tx)
}

// Used by PUT requests to `/jobs`, replaces an existing content invalidation job
// with the provided request body.
//
// Deprecated. To be used only with versions less than 4.0
func Update(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	var oFQDN string
	var dsid uint
	var uid uint
	job := tc.InvalidationJob{}
	row := inf.Tx.Tx.QueryRow(putInfoQuery, inf.Params["id"])
	err := row.Scan(&job.ID,
		&job.CreatedBy,
		&uid,
		&dsid,
		&job.DeliveryService,
		&job.AssetURL,
		&job.Parameters,
		&job.StartTime,
		&oFQDN)
	if err != nil {
		if err == sql.ErrNoRows {
			userErr = fmt.Errorf("No job by id '%s'!", inf.Params["id"])
			errCode = http.StatusNotFound
		} else {
			sysErr = fmt.Errorf("fetching job update info: %v", err)
			errCode = http.StatusInternalServerError
		}
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}

	if ok, err := IsUserAuthorizedToModifyDSID(inf, dsid); err != nil {
		sysErr = fmt.Errorf("Checking user permissions on DS #%d: %v", dsid, err)
		errCode = http.StatusInternalServerError
		api.HandleErr(w, r, inf.Tx.Tx, errCode, nil, sysErr)
		return
	} else if !ok {
		userErr = errors.New("No such Delivery Service!")
		errCode = http.StatusNotFound
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, nil)
		return
	}

	if ok, err := IsUserAuthorizedToModifyJobsMadeByUsername(inf, *job.CreatedBy); err != nil {
		sysErr = fmt.Errorf("Checking user permissions against user %s: %v", *job.CreatedBy, err)
		errCode = http.StatusInternalServerError
		api.HandleErr(w, r, inf.Tx.Tx, errCode, nil, sysErr)
		return
	} else if !ok {
		userErr = fmt.Errorf("No job by id '%s'!", inf.Params["id"])
		errCode = http.StatusNotFound
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, nil)
		return
	}

	input := tc.InvalidationJob{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		userErr = fmt.Errorf("Unable to parse input: %v", err)
		sysErr = fmt.Errorf("parsing input to PUT jobs?id=%s: %v", inf.Params["id"], err)
		errCode = http.StatusBadRequest
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}

	if err := input.Validate(); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, err, nil)
		return
	}

	if !strings.HasPrefix(*input.AssetURL, oFQDN) {
		userErr = fmt.Errorf("Cannot set asset URL that does not start with Delivery Service origin URL: %s", oFQDN)
		errCode = http.StatusBadRequest
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, nil)
		return
	}

	if job.StartTime.Before(time.Now()) {
		userErr = errors.New("Cannot modify a job that has already started!")
		errCode = http.StatusMethodNotAllowed
		w.Header().Set(http.CanonicalHeaderKey("allow"), "GET,HEAD,DELETE")
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, nil)
		return
	}

	if *job.DeliveryService != *input.DeliveryService {
		userErr = errors.New("Cannot change 'deliveryService' of existing invalidation job!")
		errCode = http.StatusConflict
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, nil)
		return
	}

	if *job.CreatedBy != *input.CreatedBy {
		userErr = errors.New("Cannot change 'createdBy' of existing invalidation jobs!")
		errCode = http.StatusConflict
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, nil)
		return
	}

	if *job.ID != *input.ID {
		userErr = errors.New("Cannot change an invalidation job 'id'!")
		errCode = http.StatusConflict
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, nil)
		return
	}

	_, cdnName, _, err := dbhelpers.GetDSNameAndCDNFromID(inf.Tx.Tx, int(dsid))
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting delivery service and CDN name from ID: "+err.Error()))
		return
	}
	userErr, sysErr, statusCode := dbhelpers.CheckIfCurrentUserCanModifyCDN(inf.Tx.Tx, string(cdnName), inf.User.UserName)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, statusCode, userErr, sysErr)
		return
	}

	row = inf.Tx.Tx.QueryRow(updateQuery,
		input.AssetURL,
		strings.TrimSuffix(strings.TrimPrefix(*input.Parameters, "TTL:"), "h"), // Strip TTL: and h from 'TTL:##h'
		input.StartTime.Time,
		*job.ID)
	err = row.Scan(&job.AssetURL,
		&job.CreatedBy,
		&job.DeliveryService,
		&job.ID,
		&job.Keyword,
		&job.Parameters,
		&job.StartTime)
	if err != nil {
		sysErr = fmt.Errorf("Updating a job: %v", err)
		errCode = http.StatusInternalServerError
		api.HandleErr(w, r, inf.Tx.Tx, errCode, nil, sysErr)
		return
	}

	if err = setRevalFlags(*job.DeliveryService, inf.Tx.Tx); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("Setting reval flags: %v", err))
		return
	}

	ttlHours := input.TTLHours()
	conflicts := tc.ValidateJobUniqueness(inf.Tx.Tx, dsid, input.StartTime.Time, *input.AssetURL, ttlHours)
	response := apiResponse{
		make([]tc.Alert, len(conflicts)+1),
		job,
	}
	for i, conflict := range conflicts {
		response.Alerts[i] = tc.Alert{
			Text:  conflict,
			Level: tc.WarnLevel.String(),
		}
	}
	response.Alerts[len(conflicts)] = tc.Alert{
		Text: fmt.Sprintf("Invalidation request created for %v, start:%v end %v", *job.AssetURL, job.StartTime.Time,
			job.StartTime.Add(time.Hour*time.Duration(ttlHours))),
		Level: tc.SuccessLevel.String(),
	}

	resp, err := json.Marshal(response)
	if err != nil {
		sysErr = fmt.Errorf("encoding response: %v", err)
		errCode = http.StatusInternalServerError
		api.HandleErr(w, r, inf.Tx.Tx, errCode, nil, sysErr)
		return
	}

	w.Header().Set(http.CanonicalHeaderKey("content-type"), rfc.ApplicationJSON)
	api.WriteAndLogErr(w, r, append(resp, '\n'))

	api.CreateChangeLogRawTx(api.ApiChange, api.Updated+" content invalidation job - ID: "+strconv.FormatUint(*job.ID, 10)+" DS: "+*job.DeliveryService+" URL: '"+*job.AssetURL+"' Params: '"+*job.Parameters+"'", inf.User, inf.Tx.Tx)
}

// Used by DELETE requests to `/jobs`, deletes an existing content invalidation job
func DeleteV40(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id"}, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	var dsid uint
	var createdBy uint
	row := inf.Tx.Tx.QueryRow(`SELECT job_deliveryservice, job_user FROM job WHERE id=$1`, inf.Params["id"])
	if err := row.Scan(&dsid, &createdBy); err != nil {
		if err == sql.ErrNoRows {
			userErr = fmt.Errorf("No job by id '%s'!", inf.Params["id"])
			errCode = http.StatusNotFound
		} else {
			sysErr = fmt.Errorf("Getting info for job #%s: %v", inf.Params["id"], err)
			errCode = http.StatusInternalServerError
		}
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}

	if ok, err := IsUserAuthorizedToModifyDSID(inf, dsid); err != nil {
		sysErr = fmt.Errorf("Checking user permissions on DS #%d: %v", dsid, err)
		errCode = http.StatusInternalServerError
		api.HandleErr(w, r, inf.Tx.Tx, errCode, nil, sysErr)
		return
	} else if !ok {
		userErr = errors.New("No such Delivery Service!")
		errCode = http.StatusNotFound
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, nil)
		return
	}

	if ok, err := IsUserAuthorizedToModifyJobsMadeByUserID(inf, createdBy); err != nil {
		sysErr = fmt.Errorf("Checking user permissions against user %v: %v", createdBy, err)
		errCode = http.StatusInternalServerError
		api.HandleErr(w, r, inf.Tx.Tx, errCode, nil, sysErr)
		return
	} else if !ok {
		userErr = fmt.Errorf("No job by id '%s'!", inf.Params["id"])
		errCode = http.StatusNotFound
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, nil)
		return
	}

	result := tc.InvalidationJobV4{}
	row = inf.Tx.Tx.QueryRow(deleteQueryV4, inf.Params["id"])
	err := row.Scan(
		&result.ID,
		&result.AssetURL,
		&result.CreatedBy,
		&result.DeliveryService,
		&result.TTLHours,
		&result.InvalidationType,
		&result.StartTime)
	if err != nil {
		sysErr = fmt.Errorf("deleting job #%s: %v", inf.Params["id"], err)
		errCode = http.StatusInternalServerError
		api.HandleErr(w, r, inf.Tx.Tx, errCode, nil, sysErr)
		return
	}

	if err = setRevalFlags(dsid, inf.Tx.Tx); err != nil {
		sysErr = fmt.Errorf("setting reval_pending after deleting job #%s: %v", inf.Params["id"], err)
		errCode = http.StatusInternalServerError
		api.HandleErr(w, r, inf.Tx.Tx, errCode, nil, sysErr)
		return
	}

	response := apiResponseV4{[]tc.Alert{
		{Text: "Content invalidation job was deleted", Level: tc.SuccessLevel.String()},
	},
		result,
	}
	resp, err := json.Marshal(response)
	if err != nil {
		sysErr = fmt.Errorf("encoding response: %v", err)
		errCode = http.StatusInternalServerError
		api.HandleErr(w, r, inf.Tx.Tx, errCode, nil, sysErr)
		return
	}

	w.Header().Set(http.CanonicalHeaderKey("content-type"), rfc.ApplicationJSON)
	api.WriteAndLogErr(w, r, append(resp, '\n'))

	changeLogMsg := fmt.Sprintf("%s content invalidation job - ID: %d DSXMLID: %s ASSET_URL: '%s' TTLHRs: %d INVALIDATION: %s",
		api.Deleted,
		result.ID,
		result.DeliveryService,
		result.AssetURL,
		result.TTLHours,
		result.InvalidationType,
	)
	api.CreateChangeLogRawTx(api.ApiChange,
		changeLogMsg,
		inf.User,
		inf.Tx.Tx)
}

// Used by DELETE requests to `/jobs`, deletes an existing content invalidation job
//
// Deprecated. To be used only with versions less than 4.0
func Delete(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id"}, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	var dsid uint
	var createdBy uint
	row := inf.Tx.Tx.QueryRow(`SELECT job_deliveryservice, job_user FROM job WHERE id=$1`, inf.Params["id"])
	if err := row.Scan(&dsid, &createdBy); err != nil {
		if err == sql.ErrNoRows {
			userErr = fmt.Errorf("No job by id '%s'!", inf.Params["id"])
			errCode = http.StatusNotFound
		} else {
			sysErr = fmt.Errorf("Getting info for job #%s: %v", inf.Params["id"], err)
			errCode = http.StatusInternalServerError
		}
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}

	if ok, err := IsUserAuthorizedToModifyDSID(inf, dsid); err != nil {
		sysErr = fmt.Errorf("Checking user permissions on DS #%d: %v", dsid, err)
		errCode = http.StatusInternalServerError
		api.HandleErr(w, r, inf.Tx.Tx, errCode, nil, sysErr)
		return
	} else if !ok {
		userErr = errors.New("No such Delivery Service!")
		errCode = http.StatusNotFound
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, nil)
		return
	}

	if ok, err := IsUserAuthorizedToModifyJobsMadeByUserID(inf, createdBy); err != nil {
		sysErr = fmt.Errorf("Checking user permissions against user %v: %v", createdBy, err)
		errCode = http.StatusInternalServerError
		api.HandleErr(w, r, inf.Tx.Tx, errCode, nil, sysErr)
		return
	} else if !ok {
		userErr = fmt.Errorf("No job by id '%s'!", inf.Params["id"])
		errCode = http.StatusNotFound
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, nil)
		return
	}

	_, cdnName, _, err := dbhelpers.GetDSNameAndCDNFromID(inf.Tx.Tx, int(dsid))
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting delivery service and CDN name from ID: "+err.Error()))
		return
	}
	userErr, sysErr, statusCode := dbhelpers.CheckIfCurrentUserCanModifyCDN(inf.Tx.Tx, string(cdnName), inf.User.UserName)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, statusCode, userErr, sysErr)
		return
	}

	result := tc.InvalidationJob{}
	row = inf.Tx.Tx.QueryRow(deleteQuery, inf.Params["id"])
	err = row.Scan(&result.AssetURL,
		&result.CreatedBy,
		&result.DeliveryService,
		&result.ID,
		&result.Keyword,
		&result.Parameters,
		&result.StartTime)
	if err != nil {
		sysErr = fmt.Errorf("deleting job #%s: %v", inf.Params["id"], err)
		errCode = http.StatusInternalServerError
		api.HandleErr(w, r, inf.Tx.Tx, errCode, nil, sysErr)
		return
	}

	if err = setRevalFlags(dsid, inf.Tx.Tx); err != nil {
		sysErr = fmt.Errorf("setting reval_pending after deleting job #%s: %v", inf.Params["id"], err)
		errCode = http.StatusInternalServerError
		api.HandleErr(w, r, inf.Tx.Tx, errCode, nil, sysErr)
		return
	}

	response := apiResponse{[]tc.Alert{tc.Alert{Text: "Content invalidation job was deleted", Level: tc.SuccessLevel.String()}}, result}
	resp, err := json.Marshal(response)
	if err != nil {
		sysErr = fmt.Errorf("encoding response: %v", err)
		errCode = http.StatusInternalServerError
		api.HandleErr(w, r, inf.Tx.Tx, errCode, nil, sysErr)
		return
	}

	w.Header().Set(http.CanonicalHeaderKey("content-type"), rfc.ApplicationJSON)
	api.WriteAndLogErr(w, r, append(resp, '\n'))

	api.CreateChangeLogRawTx(api.ApiChange, api.Deleted+" content invalidation job - ID: "+strconv.FormatUint(*result.ID, 10)+" DS: "+*result.DeliveryService+" URL: '"+*result.AssetURL+"' Params: '"+*result.Parameters+"'", inf.User, inf.Tx.Tx)
}

// Validates the fields submitted for an InvalidationJobCreateV40. These errors
// are ultimately returned to the user
func validateJobCreateV4(job tc.InvalidationJobCreateV4, tx *sql.Tx) error {
	errs := []string{}
	err := validation.ValidateStruct(&job,
		validation.Field(&job.DeliveryService, validation.Required),
		validation.Field(&job.Regex, validation.Required, validation.NewStringRule(func(s string) bool {
			return strings.HasPrefix(s, `\/`) || strings.HasPrefix(s, "/")
		}, `must start with '/' (or '\/')`)),
		validation.Field(&job.StartTime, validation.Required),
		validation.Field(&job.TTLHours, validation.Required),
		validation.Field(&job.InvalidationType, validation.Required, validation.NewStringRule(func(s string) bool {
			return s == tc.REFRESH || s == tc.REFETCH
		}, fmt.Sprintf("must be either %s or %s (case sensitive)", tc.REFRESH, tc.REFETCH))),
	)
	if err != nil {
		errs = append(errs, err.Error())
	}

	if _, _, err := dbhelpers.GetDSIDFromXMLID(tx, job.DeliveryService); err != nil {
		errs = append(errs, "Delivery Service is invalid: "+err.Error())
	}

	if _, err := regexp.Compile(job.Regex); err != nil {
		errs = append(errs, "regex: is not a valid Regular Expression: "+err.Error())
	}

	if job.StartTime.Before(time.Now()) {
		errs = append(errs, "startTime: must be in the future")
	}

	if valid, err := validateTLLHours(job.TTLHours, tx); !valid {
		if err != nil {
			errs = append(errs, "TTL is invalid: "+err.Error())
		} else {
			errs = append(errs, "TTL is invalid")
		}
	}

	if job.InvalidationType == tc.REFETCH && !refetchAllowed(tx) {
		errs = append(errs, "InvalidationType is invalid")
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ", "))
	}

	return nil
}

// validateInvalidationJobV4 checks that the InvalidationJob is valid, by ensuring all of its fields are well-defined.
// This returns an error describing any and all problematic fields encountered during validation.
func validateInvalidationJobV4(job tc.InvalidationJobV4) error {
	errs := []string{}
	err := validation.ValidateStruct(&job,
		validation.Field(&job.DeliveryService, validation.Required),
		validation.Field(&job.AssetURL, validation.Required, is.URL),
		validation.Field(&job.CreatedBy, validation.Required),
		validation.Field(&job.ID, validation.Required),
		validation.Field(&job.TTLHours, validation.Required),
		validation.Field(&job.StartTime, validation.Required),
		validation.Field(&job.InvalidationType, validation.Required, validation.NewStringRule(func(s string) bool {
			return s == tc.REFRESH || s == tc.REFETCH
		}, fmt.Sprintf("must be either %s or %s (case sensitive)", tc.REFRESH, tc.REFETCH))),
	)

	if err != nil {
		errs = append(errs, err.Error())
	}

	if job.StartTime.After(time.Now().Add(time.Hour * 48)) {
		errs = append(errs, "startTime: must be within two days from now")
	}

	if job.StartTime.Before(time.Now()) {
		errs = append(errs, "startTime: cannot be in the past")
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ", "))
	}

	return nil
}

// validateTLLHours ensures the supplied TTL hours is within acceptable limits
func validateTLLHours(ttlHours uint32, tx *sql.Tx) (bool, error) {
	var maxDays uint
	err := tx.QueryRow(`SELECT value FROM parameter WHERE name='maxRevalDurationDays' AND config_file='regex_revalidate.config'`).Scan(&maxDays)
	maxHours := maxDays * 24
	if err != nil {
		log.Errorf("error querying \"maxRevalDurationDays\" parameter: %v", err)
		return false, nil // sent to the user, hide server error
	}
	if err == nil && uint(ttlHours) > maxHours {
		return false, fmt.Errorf("cannot exceed %s", strconv.FormatUint(uint64(maxHours), 10))
	}
	return true, nil
}

// refetchAllowed checks whether Refetch is allowed and enabled in the parameter table
func refetchAllowed(tx *sql.Tx) bool {
	refetchEnabled := false
	err := tx.QueryRow(`SELECT 'true' = lower(trim(p.value)) FROM "parameter" p WHERE p.name=$1 AND p.config_file=$2`,
		tc.RefetchEnabled, tc.GlobalConfigFileName).Scan(&refetchEnabled)
	if err != nil {
		log.Errorf("error querying \"refetch_enabled\" from parameter: %v", err)
		return refetchEnabled // sent to the user, hide server error
	}
	return refetchEnabled
}

// API versions below 4.0 allowed for either the Delivery Service ID (uint) OR Delivery Service XML-ID (string).
// This can be refactored once api versions below 4.0 are removed to take a Delivery Service XML-ID (string), rather
// than an empty interface {}.
func setRevalFlags(d interface{}, tx *sql.Tx) error {
	var useReval string
	row := tx.QueryRow(`SELECT value FROM parameter WHERE name=$1 AND config_file=$2`, tc.UseRevalPendingParameterName, tc.GlobalConfigFileName)
	if err := row.Scan(&useReval); err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		useReval = "0"
	}

	column := "revalidate_update_time"
	if useReval == "0" {
		column = "config_update_time"
	}

	var q string
	switch t := d.(type) {
	case uint:
		q = fmt.Sprintf(queueUpdateOrRevalQuery, column, "id")
	case string:
		q = fmt.Sprintf(queueUpdateOrRevalQuery, column, "xml_id")
	default:
		return fmt.Errorf("invalid type passed to 'setRevalFlags': %v", t)
	}

	row = tx.QueryRow(q, d)
	if err := row.Scan(); err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil
}

// Checks if the current user's (identified in the APIInfo) tenant has permissions to
// edit a Delivery Service. `ds` is expected to be the integral, unique identifer of the
// Delivery Service in question.
//
// This returns, in order, a boolean that indicates whether or not the current user
// has the required tenancy to modify the indicated Delivery Service, and an error
// indicating what, if anything, went wrong during the check.
// returned errors is not nil, otherwise its value is undefined.
//
// Note: If no such delivery service exists, the return values shall indicate that the
// user isn't authorized.
func IsUserAuthorizedToModifyDSID(inf *api.APIInfo, ds uint) (bool, error) {
	var t uint
	row := inf.Tx.Tx.QueryRow(`SELECT tenant_id FROM deliveryservice WHERE id=$1`, ds)
	if err := row.Scan(&t); err != nil {
		if err == sql.ErrNoRows {
			return false, nil //I do this to conceal the existence of DSes for which the user has no permission to see
		}
		return false, err
	}

	return tenant.IsResourceAuthorizedToUserTx(int(t), inf.User, inf.Tx.Tx)
}

// Checks if the current user's (identified in the APIInfo) tenant has permissions to
// edit a Delivery Service. `ds` is expected to be the "xml_id" of the
// Delivery Service in question.
//
// This returns, in order, a boolean that indicates whether or not the current user
// has the required tenancy to modify the indicated Delivery Service, and an error
// indicating what, if anything, went wrong during the check.
// returned errors is not nil, otherwise its value is undefined.
//
// Note: If no such delivery service exists, the return values shall indicate that the
// user isn't authorized.
func IsUserAuthorizedToModifyDSXMLID(inf *api.APIInfo, ds string) (bool, error) {
	var t uint
	row := inf.Tx.Tx.QueryRow(`SELECT tenant_id FROM deliveryservice WHERE xml_id=$1`, ds)
	if err := row.Scan(&t); err != nil {
		if err == sql.ErrNoRows {
			return false, nil //I do this to conceal the existence of DSes for which the user has no permission to see
		}
		return false, err
	}

	return tenant.IsResourceAuthorizedToUserTx(int(t), inf.User, inf.Tx.Tx)
}

// Checks if the current user's (identified in the APIInfo) tenant has permissions to
// edit on par with the user identified by `u`. `u` is expected to be the integral,
// unique identifer of the user in question (not the current, requesting user).
//
// This returns, in order, a boolean that indicates whether or not the current user
// has the required tenancy to modify the indicated Delivery Service, and an error
// indicating what, if anything, went wrong during the check.
// returned errors is not nil, otherwise its value is undefined.
//
// Note: If no such delivery service exists, the return values shall indicate that the
// user isn't authorized.
func IsUserAuthorizedToModifyJobsMadeByUserID(inf *api.APIInfo, u uint) (bool, error) {
	var t uint
	row := inf.Tx.Tx.QueryRow(`SELECT tenant_id FROM tm_user WHERE id=$1`, u)
	if err := row.Scan(&t); err != nil {
		if err == sql.ErrNoRows {
			return false, nil //I do this to conceal the existence of DSes for which the user has no permission to see
		}
		return false, err
	}

	return tenant.IsResourceAuthorizedToUserTx(int(t), inf.User, inf.Tx.Tx)
}

// Checks if the current user's (identified in the APIInfo) tenant has permissions to
// edit on par with the user identified by `u`. `u` is expected to be the username of
// the user in question (not the current, requesting user).
//
// This returns, in order, a boolean that indicates whether or not the current user
// has the required tenancy to modify the indicated Delivery Service, and an error
// indicating what, if anything, went wrong during the check.
// returned errors is not nil, otherwise its value is undefined.
//
// Note: If no such delivery service exists, the return values shall indicate that the
// user isn't authorized.
func IsUserAuthorizedToModifyJobsMadeByUsername(inf *api.APIInfo, u string) (bool, error) {
	var t uint
	row := inf.Tx.Tx.QueryRow(`SELECT tenant_id FROM tm_user WHERE username=$1`, u)
	if err := row.Scan(&t); err != nil {
		if err == sql.ErrNoRows {
			return false, nil //I do this to conceal the existence of DSes for which the user has no permission to see
		}
		return false, err
	}

	return tenant.IsResourceAuthorizedToUserTx(int(t), inf.User, inf.Tx.Tx)
}
