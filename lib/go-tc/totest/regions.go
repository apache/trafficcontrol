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

func CreateTestRegions(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, region := range td.Regions {
		resp, _, err := cl.CreateRegion(region, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Could not create Region '%s': %v - alerts: %+v", region.Name, err, resp.Alerts)
	}
}

func DeleteTestRegions(t *testing.T, cl *toclient.Session) {
	regions, _, err := cl.GetRegions(toclient.RequestOptions{})
	assert.NoError(t, err, "Cannot get Regions: %v - alerts: %+v", err, regions.Alerts)

	for _, region := range regions.Response {
		alerts, _, err := cl.DeleteRegion(region.Name, toclient.RequestOptions{})
		assert.NoError(t, err, "Unexpected error deleting Region '%s' (#%d): %v - alerts: %+v", region.Name, region.ID, err, alerts.Alerts)
		// Retrieve the Region to see if it got deleted
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.Itoa(region.ID))
		getRegion, _, err := cl.GetRegions(opts)
		assert.NoError(t, err, "Error getting Region '%s' after deletion: %v - alerts: %+v", region.Name, err, getRegion.Alerts)
		assert.Equal(t, 0, len(getRegion.Response), "Expected Region '%s' to be deleted, but it was found in Traffic Ops", region.Name)
	}
}
