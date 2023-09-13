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

	"github.com/apache/trafficcontrol/lib/go-util/assert"
	toclient "github.com/apache/trafficcontrol/traffic_ops/v5-client"
)

func CreateTestTopologies(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, topology := range td.Topologies {
		resp, _, err := cl.CreateTopology(topology, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Could not create Topology: %v - alerts: %+v", err, resp.Alerts)
	}
}

func DeleteTestTopologies(t *testing.T, cl *toclient.Session) {
	topologies, _, err := cl.GetTopologies(toclient.RequestOptions{})
	assert.NoError(t, err, "Cannot get Topologies: %v - alerts: %+v", err, topologies.Alerts)

	for _, topology := range topologies.Response {
		alerts, _, err := cl.DeleteTopology(topology.Name, toclient.RequestOptions{})
		assert.NoError(t, err, "Cannot delete Topology: %v - alerts: %+v", err, alerts.Alerts)
		// Retrieve the Topology to see if it got deleted
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("name", topology.Name)
		resp, _, err := cl.GetTopologies(opts)
		assert.NoError(t, err, "Unexpected error trying to fetch Topologies after deletion: %v - alerts: %+v", err, resp.Alerts)
		assert.Equal(t, 0, len(resp.Response), "Expected Topology '%s' to be deleted, but it was found in Traffic Ops", topology.Name)
	}
}
