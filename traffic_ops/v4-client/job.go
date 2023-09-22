package client

/*

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

import (
	"net/url"
	"strconv"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

// apiJobs is the API version-relative path to the /jobs API route.
const apiJobs = "/jobs"

// CreateInvalidationJob creates the passed Content Invalidation Job.
func (to *Session) CreateInvalidationJob(job tc.InvalidationJobCreateV4, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(apiJobs, opts, job, &alerts)
	return alerts, reqInf, err
}

// DeleteInvalidationJob deletes the Content Invalidation Job identified by
// 'jobID'.
func (to *Session) DeleteInvalidationJob(jobID uint64, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	if opts.QueryParameters == nil {
		opts.QueryParameters = url.Values{}
	}
	opts.QueryParameters.Set("id", strconv.FormatUint(jobID, 10))
	var alerts tc.Alerts
	reqInf, err := to.del(apiJobs, opts, &alerts)
	return alerts, reqInf, err
}

// UpdateInvalidationJob updates the passed Content Invalidation Job (it is
// expected to have an ID).
func (to *Session) UpdateInvalidationJob(job tc.InvalidationJobV4, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	if opts.QueryParameters == nil {
		opts.QueryParameters = url.Values{}
	}
	opts.QueryParameters.Set("id", strconv.FormatUint(job.ID, 10))
	reqInf, err := to.put(apiJobs, opts, job, &alerts)
	return alerts, reqInf, err
}

// GetInvalidationJobs returns a list of Content Invalidation Jobs visible to
// your Tenant.
func (to *Session) GetInvalidationJobs(opts RequestOptions) (tc.InvalidationJobsResponseV4, toclientlib.ReqInf, error) {
	var data tc.InvalidationJobsResponseV4
	reqInf, err := to.get(apiJobs, opts, &data)
	return data, reqInf, err
}
