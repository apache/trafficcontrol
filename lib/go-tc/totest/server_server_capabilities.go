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

func CreateTestServerServerCapabilities(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, ssc := range td.ServerServerCapabilities {
		assert.RequireNotNil(t, ssc.Server, "Expected Server to not be nil.")
		assert.RequireNotNil(t, ssc.ServerCapability, "Expected Server Capability to not be nil.")
		serverID := GetServerID(t, cl, *ssc.Server)()
		ssc.ServerID = &serverID
		resp, _, err := cl.CreateServerServerCapability(ssc, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Could not associate Capability '%s' with server '%s': %v - alerts: %+v", *ssc.ServerCapability, *ssc.Server, err, resp.Alerts)
	}
}

func DeleteTestServerServerCapabilities(t *testing.T, cl *toclient.Session) {
	sscs, _, err := cl.GetServerServerCapabilities(toclient.RequestOptions{})
	assert.RequireNoError(t, err, "Cannot get server server capabilities: %v - alerts: %+v", err, sscs.Alerts)
	for _, ssc := range sscs.Response {
		assert.RequireNotNil(t, ssc.Server, "Expected Server to not be nil.")
		assert.RequireNotNil(t, ssc.ServerCapability, "Expected Server Capability to not be nil.")
		alerts, _, err := cl.DeleteServerServerCapability(*ssc.ServerID, *ssc.ServerCapability, toclient.RequestOptions{})
		assert.NoError(t, err, "Could not remove Capability '%s' from server '%s' (#%d): %v - alerts: %+v", *ssc.ServerCapability, *ssc.Server, *ssc.ServerID, err, alerts.Alerts)
	}
}
