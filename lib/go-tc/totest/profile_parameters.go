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

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	toclient "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

func CreateTestProfileParameters(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, profile := range td.Profiles {
		profileID := GetProfileID(t, cl, profile.Name)()

		for _, parameter := range profile.Parameters {
			assert.RequireNotNil(t, parameter.Name, "Expected parameter name to not be nil.")
			assert.RequireNotNil(t, parameter.Value, "Expected parameter value to not be nil.")
			assert.RequireNotNil(t, parameter.ConfigFile, "Expected parameter configFile to not be nil.")

			parameterOpts := toclient.NewRequestOptions()
			parameterOpts.QueryParameters.Set("name", *parameter.Name)
			parameterOpts.QueryParameters.Set("configFile", *parameter.ConfigFile)
			parameterOpts.QueryParameters.Set("value", *parameter.Value)
			getParameter, _, err := cl.GetParameters(parameterOpts)
			assert.RequireNoError(t, err, "Could not get Parameter %s: %v - alerts: %+v", *parameter.Name, err, getParameter.Alerts)
			if len(getParameter.Response) == 0 {
				alerts, _, err := cl.CreateParameter(tc.ParameterV5{Name: *parameter.Name, Value: *parameter.Value, ConfigFile: *parameter.ConfigFile}, toclient.RequestOptions{})
				assert.RequireNoError(t, err, "Could not create Parameter %s: %v - alerts: %+v", parameter.Name, err, alerts.Alerts)
				getParameter, _, err = cl.GetParameters(parameterOpts)
				assert.RequireNoError(t, err, "Could not get Parameter %s: %v - alerts: %+v", *parameter.Name, err, getParameter.Alerts)
				assert.RequireNotEqual(t, 0, len(getParameter.Response), "Could not get parameter %s: not found", *parameter.Name)
			}
			profileParameter := tc.ProfileParameterCreationRequest{ProfileID: profileID, ParameterID: getParameter.Response[0].ID}
			alerts, _, err := cl.CreateProfileParameter(profileParameter, toclient.RequestOptions{})
			assert.NoError(t, err, "Could not associate Parameter %s with Profile %s: %v - alerts: %+v", parameter.Name, profile.Name, err, alerts.Alerts)
		}
	}
}

func DeleteTestProfileParameters(t *testing.T, cl *toclient.Session) {
	profileParameters, _, err := cl.GetProfileParameters(toclient.RequestOptions{})
	assert.NoError(t, err, "Cannot get Profile Parameters: %v - alerts: %+v", err, profileParameters.Alerts)

	for _, profileParameter := range profileParameters.Response {
		alerts, _, err := cl.DeleteProfileParameter(GetProfileID(t, cl, profileParameter.Profile)(), profileParameter.Parameter, toclient.RequestOptions{})
		assert.NoError(t, err, "Unexpected error deleting Profile Parameter: Profile: '%s' Parameter ID: (#%d): %v - alerts: %+v", profileParameter.Profile, profileParameter.Parameter, err, alerts.Alerts)
	}
	// Retrieve the Profile Parameters to see if it got deleted
	getProfileParameter, _, err := cl.GetProfileParameters(toclient.RequestOptions{})
	assert.NoError(t, err, "Error getting Profile Parameters after deletion: %v - alerts: %+v", err, getProfileParameter.Alerts)
	assert.Equal(t, 0, len(getProfileParameter.Response), "Expected Profile Parameters to be deleted, but %d were found in Traffic Ops", len(getProfileParameter.Response))
}
