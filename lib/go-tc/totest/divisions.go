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

func CreateTestDivisions(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, division := range td.Divisions {
		resp, _, err := cl.CreateDivision(division, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Could not create Division '%s': %v - alerts: %+v", division.Name, err, resp.Alerts)
	}
}

func DeleteTestDivisions(t *testing.T, cl *toclient.Session) {
	divisions, _, err := cl.GetDivisions(toclient.RequestOptions{})
	assert.NoError(t, err, "Cannot get Divisions: %v - alerts: %+v", err, divisions.Alerts)
	for _, division := range divisions.Response {
		alerts, _, err := cl.DeleteDivision(division.ID, toclient.RequestOptions{})
		assert.NoError(t, err, "Unexpected error deleting Division '%s' (#%d): %v - alerts: %+v", division.Name, division.ID, err, alerts.Alerts)
		// Retrieve the Division to see if it got deleted
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.Itoa(division.ID))
		getDivision, _, err := cl.GetDivisions(opts)
		assert.NoError(t, err, "Error getting Division '%s' after deletion: %v - alerts: %+v", division.Name, err, getDivision.Alerts)
		assert.Equal(t, 0, len(getDivision.Response), "Expected Division '%s' to be deleted, but it was found in Traffic Ops", division.Name)
	}
}
