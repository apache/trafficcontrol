package toreq

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
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
	toclient "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

func serversToLatest(svs tc.ServersV4Response) ([]atscfg.Server, error) {
	return atscfg.ToServers(svs.Response), nil
}

func serverToLatest(oldSv *tc.ServerV40) (*atscfg.Server, error) {
	asv := atscfg.Server(*oldSv)
	return &asv, nil
}

func dsesToLatest(dses []tc.DeliveryServiceV40) []atscfg.DeliveryService {
	return atscfg.V40ToDeliveryServices(dses)
}

func jobsToLatest(jobs []tc.InvalidationJobV4) []atscfg.InvalidationJob {
	return atscfg.ToInvalidationJobs(jobs)
}

func serverUpdateStatusesToLatest(statuses []tc.ServerUpdateStatusV40) []atscfg.ServerUpdateStatus {
	return atscfg.ToServerUpdateStatuses(statuses)
}

// GetJobsCompat gets jobs from any Traffic Ops built from the ATC `master` branch, and converts the different formats to the latest.
// This makes t3c work with old or new Traffic Ops deployed from `master`,
// though it doesn't make a version of t3c older than this work with a new TO,
// which isn't logically possible from the client.
func (cl *TOClient) GetJobsCompat(opts toclient.RequestOptions) (tc.InvalidationJobsResponseV4, toclientlib.ReqInf, error) {
	path := "/jobs"

	objs := struct {
		Response []InvalidationJobV4PlusLegacy `json:"response"`
		tc.Alerts
	}{}

	if len(opts.QueryParameters) > 0 {
		path += "?" + opts.QueryParameters.Encode()
	}
	reqInf, err := cl.c.TOClient.Req(http.MethodGet, path, nil, opts.Header, &objs)
	if err != nil {
		return tc.InvalidationJobsResponseV4{}, reqInf, errors.New("request: " + err.Error())
	}

	resp := tc.InvalidationJobsResponseV4{Alerts: objs.Alerts}
	for _, job := range objs.Response {
		newJob, err := InvalidationJobV4FromLegacy(job) // (InvalidationJobV4, error) {
		if err != nil {
			return tc.InvalidationJobsResponseV4{}, reqInf, errors.New("converting job from possible legacy format: " + err.Error())
		}
		resp.Response = append(resp.Response, newJob)
	}
	return resp, reqInf, nil
}

// InvalidationJobV4ForLegacy is a type alias to prevent MarshalJSON recursion.
type InvalidationJobV4ForLegacy tc.InvalidationJobV4

// InvalidationJobV4PlusLegacy has the data to deserialize both the latest and older versions that Traffic Ops could return.
type InvalidationJobV4PlusLegacy struct {
	// StartTime overrides the StartTime in InvalidationJobV4 in order to unmarshal any string format.
	//
	// A json.Unmarshal will place a 'startTime' value in this field,
	// rather than the anonymous embedded InvalidationJobV4ForLegacy (tc.InvalidationJobV4).
	//
	// InvalidationJobV4FromLegacy will then parse multiple time formats that different Traffic Ops servers may return,
	// and put the parsed time in tc.InvalidationJobV4.StartTime.
	StartTime *string `json:"startTime"`
	InvalidationJobV4ForLegacy
	InvalidationJobV4Legacy
}

type InvalidationJobV4Legacy struct {
	Keyword    *string `json:"keyword"`
	Parameters *string `json:"parameters"`
}

func InvalidationJobV4FromLegacy(job InvalidationJobV4PlusLegacy) (tc.InvalidationJobV4, error) {
	if job.StartTime != nil {
		err := error(nil)
		job.InvalidationJobV4ForLegacy.StartTime, err = time.Parse(atscfg.JobV4TimeFormat, *job.StartTime)
		if err != nil {
			job.InvalidationJobV4ForLegacy.StartTime, err = time.Parse(atscfg.JobLegacyTimeFormat, *job.StartTime)
			if err != nil {
				return tc.InvalidationJobV4{}, errors.New("malformed startTime")
			}
		}
	}

	if job.TTLHours == 0 && job.Parameters != nil {
		params := *job.Parameters
		params = strings.TrimSpace(params)
		params = strings.ToLower(params)
		params = strings.Replace(params, " ", "", -1)

		paramPrefix := strings.ToLower(atscfg.JobLegacyParamPrefix)
		paramSuffix := strings.ToLower(atscfg.JobLegacyParamSuffix)
		if !strings.HasPrefix(params, paramPrefix) || !strings.HasSuffix(params, paramSuffix) {
			return tc.InvalidationJobV4{}, errors.New("legacy job.Parameters was not nil, but unexpected format '" + params + "'")
		}

		hoursStr := params[len(paramPrefix) : len(params)-len(paramSuffix)]
		hours, err := strconv.Atoi(hoursStr)
		if err != nil {
			return tc.InvalidationJobV4{}, errors.New("legacy job.Parameters was not nil, but hours not an integer: '" + params + "'")
		}
		job.TTLHours = uint(hours)
	}

	if job.InvalidationType == "" && job.Parameters != nil {
		job.InvalidationType = tc.REFRESH
	}
	if strings.HasSuffix(job.AssetURL, atscfg.JobLegacyRefetchSuffix) {
		job.InvalidationType = tc.REFETCH
	}
	job.AssetURL = strings.TrimSuffix(job.AssetURL, atscfg.JobLegacyRefreshSuffix)
	job.AssetURL = strings.TrimSuffix(job.AssetURL, atscfg.JobLegacyRefetchSuffix)

	return tc.InvalidationJobV4(job.InvalidationJobV4ForLegacy), nil
}

// SetServerUpdateStatusCompat is a bridge to send both styles of query parameters to the
// TO endpoint /servers/{hostname-or-id}/update. The current (old) is to send a bool
// value, however this has resulted in an accidental race condition. The attempt to fix
// this is to send a timestamp representing when the config or revalidation changes
// have been applied.
//
// To ensure T3C is compatible with both the current releases and future releases
// this function will send both "styles". Once both T3C and TO have been deployed
// with the timestamp only V4 TO API endpoint, this function can be removed and the
// V4 client function `SetUpdateServerStatusTimes` may be used instead (as intended).
// *** Compatability requirement until TO (v6.3+ eta 04/22/22) is deployed with the timestamp features
func (cl *TOClient) SetServerUpdateStatusCompat(serverName string, configApplyTime, revalApplyTime *time.Time, configApplyBool, revalApplyBool *bool, opts toclient.RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	reqInf := toclientlib.ReqInf{CacheHitStatus: toclientlib.CacheHitStatusMiss}
	var alerts tc.Alerts

	if opts.QueryParameters == nil {
		opts.QueryParameters = url.Values{}
	}

	if configApplyTime != nil {
		opts.QueryParameters.Set("config_apply_time", configApplyTime.Format(time.RFC3339Nano))
	}

	if revalApplyTime != nil {
		opts.QueryParameters.Set("revalidate_apply_time", revalApplyTime.Format(time.RFC3339Nano))
	}

	if configApplyBool != nil {
		opts.QueryParameters.Set("updated", "false")
	}
	if revalApplyBool != nil {
		opts.QueryParameters.Set("reval_updated", "false")
	}

	path := `/servers/` + url.PathEscape(serverName) + `/update`
	if len(opts.QueryParameters) > 0 {
		path += "?" + opts.QueryParameters.Encode()
	}
	reqInf, err := cl.c.TOClient.Req(http.MethodPost, path, nil, opts.Header, &alerts)
	return alerts, reqInf, err
}
