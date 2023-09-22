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

func CreateTestCoordinates(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, coordinate := range td.Coordinates {
		resp, _, err := cl.CreateCoordinate(coordinate, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Could not create coordinate: %v - alerts: %+v", err, resp.Alerts)
	}
}

func DeleteTestCoordinates(t *testing.T, cl *toclient.Session) {
	coordinates, _, err := cl.GetCoordinates(toclient.RequestOptions{})
	assert.NoError(t, err, "Cannot get Coordinates: %v - alerts: %+v", err, coordinates.Alerts)
	for _, coordinate := range coordinates.Response {
		alerts, _, err := cl.DeleteCoordinate(coordinate.ID, toclient.RequestOptions{})
		assert.NoError(t, err, "Unexpected error deleting Coordinate '%s' (#%d): %v - alerts: %+v", coordinate.Name, coordinate.ID, err, alerts.Alerts)
		// Retrieve the Coordinate to see if it got deleted
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.Itoa(coordinate.ID))
		getCoordinate, _, err := cl.GetCoordinates(opts)
		assert.NoError(t, err, "Error getting Coordinate '%s' after deletion: %v - alerts: %+v", coordinate.Name, err, getCoordinate.Alerts)
		assert.Equal(t, 0, len(getCoordinate.Response), "Expected Coordinate '%s' to be deleted, but it was found in Traffic Ops", coordinate.Name)
	}
}
