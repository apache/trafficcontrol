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

func CreateTestServerCheckExtensions(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, ext := range td.ServerCheckExtensions {
		resp, _, err := cl.CreateServerCheckExtension(ext, toclient.RequestOptions{})
		assert.NoError(t, err, "Could not create Servercheck Extension: %v - alerts: %+v", err, resp.Alerts)
	}
}

func DeleteTestServerCheckExtensions(t *testing.T, cl *toclient.Session, td TrafficControl) {
	extensions, _, err := cl.GetServerCheckExtensions(toclient.RequestOptions{})
	assert.RequireNoError(t, err, "Could not get Servercheck Extensions: %v - alerts: %+v", err, extensions.Alerts)

	for _, extension := range extensions.Response {
		alerts, _, err := cl.DeleteServerCheckExtension(*extension.ID, toclient.RequestOptions{})
		assert.NoError(t, err, "Unexpected error deleting Servercheck Extension '%s' (#%d): %v - alerts: %+v", *extension.Name, *extension.ID, err, alerts.Alerts)
		// Retrieve the Server Extension to see if it got deleted
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.Itoa(*extension.ID))
		getExtension, _, err := cl.GetServerCheckExtensions(opts)
		assert.NoError(t, err, "Error getting Servercheck Extension '%s' after deletion: %v - alerts: %+v", *extension.Name, err, getExtension.Alerts)
		assert.Equal(t, 0, len(getExtension.Response), "Expected Servercheck Extension '%s' to be deleted, but it was found in Traffic Ops", *extension.Name)
	}
}
