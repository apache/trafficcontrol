package totest

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
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	toclient "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

func CreateTestJobs(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, job := range td.InvalidationJobs {
		job.StartTime = time.Now().Add(time.Minute).UTC()
		resp, _, err := cl.CreateInvalidationJob(job, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Could not create job: %v - alerts: %+v", err, resp.Alerts)
	}
}

func DeleteTestJobs(t *testing.T, cl *toclient.Session) {
	jobs, _, err := cl.GetInvalidationJobs(toclient.RequestOptions{})
	assert.NoError(t, err, "Cannot get Jobs: %v - alerts: %+v", err, jobs.Alerts)

	for _, job := range jobs.Response {
		alerts, _, err := cl.DeleteInvalidationJob(job.ID, toclient.RequestOptions{})
		assert.NoError(t, err, "Unexpected error deleting Job with ID: (#%d): %v - alerts: %+v", job.ID, err, alerts.Alerts)
		// Retrieve the Job to see if it got deleted
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.Itoa(int(job.ID)))
		getJobs, _, err := cl.GetInvalidationJobs(opts)
		assert.NoError(t, err, "Error getting Job with ID: '%d' after deletion: %v - alerts: %+v", job.ID, err, getJobs.Alerts)
		assert.Equal(t, 0, len(getJobs.Response), "Expected Job to be deleted, but it was found in Traffic Ops")
	}
}
