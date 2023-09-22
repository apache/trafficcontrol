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

func CreateTestCachegroupsDeliveryServices(t *testing.T, cl *toclient.Session) {
	dses, _, err := cl.GetDeliveryServices(toclient.RequestOptions{})
	assert.RequireNoError(t, err, "Cannot GET DeliveryServices: %v - %v", err, dses)

	opts := toclient.NewRequestOptions()
	opts.QueryParameters.Set("name", "cachegroup3")
	clientCGs, _, err := cl.GetCacheGroups(opts)
	assert.RequireNoError(t, err, "Cannot GET cachegroup: %v", err)
	assert.RequireEqual(t, len(clientCGs.Response), 1, "Getting cachegroup expected 1, got %v", len(clientCGs.Response))
	assert.RequireNotNil(t, clientCGs.Response[0].ID, "Cachegroup has a nil ID")

	dsIDs := []int{}
	for _, ds := range dses.Response {
		if *ds.CDNName == "cdn1" && ds.Topology == nil {
			dsIDs = append(dsIDs, *ds.ID)
		}
	}
	assert.RequireGreaterOrEqual(t, len(dsIDs), 1, "No Delivery Services found in CDN 'cdn1', cannot continue.")
	resp, _, err := cl.SetCacheGroupDeliveryServices(*clientCGs.Response[0].ID, dsIDs, toclient.RequestOptions{})
	assert.RequireNoError(t, err, "Setting cachegroup delivery services returned error: %v", err)
	assert.RequireGreaterOrEqual(t, len(resp.Response.ServerNames), 1, "Setting cachegroup delivery services returned success, but no servers set")
}

func DeleteTestCachegroupsDeliveryServices(t *testing.T, cl *toclient.Session) {
	opts := toclient.NewRequestOptions()
	opts.QueryParameters.Set("limit", "1000000")
	dss, _, err := cl.GetDeliveryServiceServers(opts)
	assert.NoError(t, err, "Unexpected error retrieving server-to-Delivery-Service assignments: %v - alerts: %+v", err, dss.Alerts)

	for _, ds := range dss.Response {
		setInactive(t, cl, *ds.DeliveryService)
		alerts, _, err := cl.DeleteDeliveryServiceServer(*ds.DeliveryService, *ds.Server, toclient.RequestOptions{})
		assert.NoError(t, err, "Error deleting delivery service servers: %v - alerts: %+v", err, alerts.Alerts)
	}

	dss, _, err = cl.GetDeliveryServiceServers(toclient.RequestOptions{})
	assert.NoError(t, err, "Unexpected error retrieving server-to-Delivery-Service assignments: %v - alerts: %+v", err, dss.Alerts)
	assert.Equal(t, len(dss.Response), 0, "Deleting delivery service servers: Expected empty subsequent get, actual %v", len(dss.Response))
}
