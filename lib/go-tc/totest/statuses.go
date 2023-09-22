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
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	toclient "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

func CreateTestStatuses(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, status := range td.Statuses {
		resp, _, err := cl.CreateStatus(status, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Could not create Status: %v - alerts: %+v", err, resp.Alerts)
	}
}

func DeleteTestStatuses(t *testing.T, cl *toclient.Session, td TrafficControl) {
	opts := toclient.NewRequestOptions()
	for _, status := range td.Statuses {
		assert.RequireNotNil(t, status.Name, "Cannot get test statuses: test data statuses must have names")
		// Retrieve the Status by name, so we can get the id for the Update
		opts.QueryParameters.Set("name", *status.Name)
		resp, _, err := cl.GetStatuses(opts)
		assert.RequireNoError(t, err, "Cannot get Statuses filtered by name '%s': %v - alerts: %+v", *status.Name, err, resp.Alerts)
		assert.RequireEqual(t, 1, len(resp.Response), "Expected 1 status returned. Got: %d", len(resp.Response))
		respStatus := resp.Response[0]

		delResp, _, err := cl.DeleteStatus(respStatus.ID, toclient.RequestOptions{})
		assert.NoError(t, err, "Cannot delete Status: %v - alerts: %+v", err, delResp.Alerts)

		// Retrieve the Status to see if it got deleted
		resp, _, err = cl.GetStatuses(opts)
		assert.NoError(t, err, "Unexpected error getting Statuses filtered by name after deletion: %v - alerts: %+v", err, resp.Alerts)
		assert.Equal(t, 0, len(resp.Response), "Expected Status '%s' to be deleted, but it was found in Traffic Ops", *status.Name)
	}
}
