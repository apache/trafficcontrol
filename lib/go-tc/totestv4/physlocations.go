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

func CreateTestPhysLocations(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, pl := range td.PhysLocations {
		resp, _, err := cl.CreatePhysLocation(pl, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Could not create Physical Location '%s': %v - alerts: %+v", pl.Name, err, resp.Alerts)
	}
}

func DeleteTestPhysLocations(t *testing.T, cl *toclient.Session) {
	physicalLocations, _, err := cl.GetPhysLocations(toclient.RequestOptions{})
	assert.NoError(t, err, "Cannot get Physical Locations: %v - alerts: %+v", err, physicalLocations.Alerts)

	for _, pl := range physicalLocations.Response {
		alerts, _, err := cl.DeletePhysLocation(pl.ID, toclient.RequestOptions{})
		assert.NoError(t, err, "Unexpected error deleting Physical Location '%s' (#%d): %v - alerts: %+v", pl.Name, pl.ID, err, alerts.Alerts)
		// Retrieve the PhysLocation to see if it got deleted
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.Itoa(pl.ID))
		getPL, _, err := cl.GetPhysLocations(opts)
		assert.NoError(t, err, "Error getting Physical Location '%s' after deletion: %v - alerts: %+v", pl.Name, err, getPL.Alerts)
		assert.Equal(t, 0, len(getPL.Response), "Expected Physical Location '%s' to be deleted, but it was found in Traffic Ops", pl.Name)
	}
}
