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

func CreateTestParameters(t *testing.T, cl *toclient.Session, td TrafficControl) {
	alerts, _, err := cl.CreateMultipleParameters(td.Parameters, toclient.RequestOptions{})
	assert.RequireNoError(t, err, "Could not create Parameters: %v - alerts: %+v", err, alerts)
}

func DeleteTestParameters(t *testing.T, cl *toclient.Session) {
	parameters, _, err := cl.GetParameters(toclient.RequestOptions{})
	assert.RequireNoError(t, err, "Cannot get Parameters: %v - alerts: %+v", err, parameters.Alerts)

	for _, parameter := range parameters.Response {
		alerts, _, err := cl.DeleteParameter(parameter.ID, toclient.RequestOptions{})
		assert.NoError(t, err, "Cannot delete Parameter #%d: %v - alerts: %+v", parameter.ID, err, alerts.Alerts)

		// Retrieve the Parameter to see if it got deleted
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.Itoa(parameter.ID))
		getParameters, _, err := cl.GetParameters(opts)
		assert.NoError(t, err, "Unexpected error fetching Parameter #%d after deletion: %v - alerts: %+v", parameter.ID, err, getParameters.Alerts)
		assert.Equal(t, 0, len(getParameters.Response), "Expected Parameter '%s' to be deleted, but it was found in Traffic Ops", parameter.Name)
	}
}
