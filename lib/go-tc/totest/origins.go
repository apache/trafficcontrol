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

	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	toclient "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

func CreateTestOrigins(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, origin := range td.Origins {
		resp, _, err := cl.CreateOrigin(origin, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Could not create Origins: %v - alerts: %+v", err, resp.Alerts)
	}
}

func DeleteTestOrigins(t *testing.T, cl *toclient.Session) {
	origins, _, err := cl.GetOrigins(toclient.RequestOptions{})
	assert.NoError(t, err, "Cannot get Origins : %v - alerts: %+v", err, origins.Alerts)

	for _, origin := range origins.Response {
		assert.RequireNotNil(t, origin.ID, "Expected origin ID to not be nil.")
		assert.RequireNotNil(t, origin.Name, "Expected origin ID to not be nil.")
		assert.RequireNotNil(t, origin.IsPrimary, "Expected origin ID to not be nil.")
		if !origin.IsPrimary {
			alerts, _, err := cl.DeleteOrigin(origin.ID, toclient.RequestOptions{})
			assert.NoError(t, err, "Unexpected error deleting Origin '%s' (#%d): %v - alerts: %+v", origin.Name, origin.ID, err, alerts.Alerts)
			// Retrieve the Origin to see if it got deleted
			opts := toclient.NewRequestOptions()
			opts.QueryParameters.Set("id", strconv.Itoa(origin.ID))
			getOrigin, _, err := cl.GetOrigins(opts)
			assert.NoError(t, err, "Error getting Origin '%s' after deletion: %v - alerts: %+v", origin.Name, err, getOrigin.Alerts)
			assert.Equal(t, 0, len(getOrigin.Response), "Expected Origin '%s' to be deleted, but it was found in Traffic Ops", origin.Name)
		}
	}
}
