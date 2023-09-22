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

func CreateTestServers(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, server := range td.Servers {
		resp, _, err := cl.CreateServer(server, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Could not create server '%s': %v - alerts: %+v", *server.HostName, err, resp.Alerts)
	}
}

func DeleteTestServers(t *testing.T, cl *toclient.Session) {
	servers, _, err := cl.GetServers(toclient.RequestOptions{})
	assert.NoError(t, err, "Cannot get Servers: %v - alerts: %+v", err, servers.Alerts)

	for _, server := range servers.Response {
		delResp, _, err := cl.DeleteServer(*server.ID, toclient.RequestOptions{})
		assert.NoError(t, err, "Could not delete Server: %v - alerts: %+v", err, delResp.Alerts)
		// Retrieve Server to see if it got deleted
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.Itoa(*server.ID))
		getServer, _, err := cl.GetServers(opts)
		assert.RequireNotNil(t, server.HostName, "Expected server host name to not be nil.")
		assert.NoError(t, err, "Error deleting Server for '%s' : %v - alerts: %+v", *server.HostName, err, getServer.Alerts)
		assert.Equal(t, 0, len(getServer.Response), "Expected Server '%s' to be deleted", *server.HostName)
	}
}

func GetServerID(t *testing.T, cl *toclient.Session, hostName string) func() int {
	return func() int {
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("hostName", hostName)
		serversResp, _, err := cl.GetServers(opts)
		assert.RequireNoError(t, err, "Get Servers Request failed with error:", err)
		assert.RequireEqual(t, 1, len(serversResp.Response), "Expected response object length 1, but got %d", len(serversResp.Response))
		assert.RequireNotNil(t, serversResp.Response[0].ID, "Expected id to not be nil")
		return *serversResp.Response[0].ID
	}
}
