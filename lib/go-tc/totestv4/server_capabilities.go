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
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	toclient "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"
)

func CreateTestServerCapabilities(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, sc := range td.ServerCapabilities {
		resp, _, err := cl.CreateServerCapabilityV41(sc, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Unexpected error creating Server Capability '%s': %v - alerts: %+v", sc.Name, err, resp.Alerts)
	}
}

func DeleteTestServerCapabilities(t *testing.T, cl *toclient.Session) {
	serverCapabilities, _, err := cl.GetServerCapabilities(toclient.RequestOptions{})
	assert.NoError(t, err, "Cannot get Server Capabilities: %v - alerts: %+v", err, serverCapabilities.Alerts)

	for _, serverCapability := range serverCapabilities.Response {
		alerts, _, err := cl.DeleteServerCapability(serverCapability.Name, toclient.RequestOptions{})
		assert.NoError(t, err, "Unexpected error deleting Server Capability '%s': %v - alerts: %+v", serverCapability.Name, err, alerts.Alerts)
		// Retrieve the Server Capability to see if it got deleted
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("name", serverCapability.Name)
		getServerCapability, _, err := cl.GetServerCapabilities(opts)
		assert.NoError(t, err, "Error getting Server Capability '%s' after deletion: %v - alerts: %+v", serverCapability.Name, err, getServerCapability.Alerts)
		assert.Equal(t, 0, len(getServerCapability.Response), "Expected Server Capability '%s' to be deleted, but it was found in Traffic Ops", serverCapability.Name)
	}
}
