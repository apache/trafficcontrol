package totestv4

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
	toclient "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"
)

func CreateTestProfiles(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, profile := range td.Profiles {
		resp, _, err := cl.CreateProfile(profile, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Could not create Profile '%s': %v - alerts: %+v", profile.Name, err, resp.Alerts)
	}
}

func DeleteTestProfiles(t *testing.T, cl *toclient.Session) {
	profiles, _, err := cl.GetProfiles(toclient.RequestOptions{})
	assert.NoError(t, err, "Cannot get Profiles: %v - alerts: %+v", err, profiles.Alerts)
	for _, profile := range profiles.Response {
		alerts, _, err := cl.DeleteProfile(profile.ID, toclient.RequestOptions{})
		assert.NoError(t, err, "Cannot delete Profile: %v - alerts: %+v", err, alerts.Alerts)
		// Retrieve the Profile to see if it got deleted
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.Itoa(profile.ID))
		getProfiles, _, err := cl.GetProfiles(opts)
		assert.NoError(t, err, "Error getting Profile '%s' after deletion: %v - alerts: %+v", profile.Name, err, getProfiles.Alerts)
		assert.Equal(t, 0, len(getProfiles.Response), "Expected Profile '%s' to be deleted, but it was found in Traffic Ops", profile.Name)
	}
}

func GetProfileID(t *testing.T, cl *toclient.Session, profileName string) func() int {
	return func() int {
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("name", profileName)
		resp, _, err := cl.GetProfiles(opts)
		assert.RequireNoError(t, err, "Get Profiles Request failed with error: %v", err)
		assert.RequireEqual(t, 1, len(resp.Response), "Expected response object length 1, but got %d", len(resp.Response))
		return resp.Response[0].ID
	}
}
