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

func CreateTestServiceCategories(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, serviceCategory := range td.ServiceCategories {
		resp, _, err := cl.CreateServiceCategory(serviceCategory, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Could not create Service Category: %v - alerts: %+v", err, resp.Alerts)
	}
}

func DeleteTestServiceCategories(t *testing.T, cl *toclient.Session) {
	serviceCategories, _, err := cl.GetServiceCategories(toclient.RequestOptions{})
	assert.NoError(t, err, "Cannot get Service Categories: %v - alerts: %+v", err, serviceCategories.Alerts)

	for _, serviceCategory := range serviceCategories.Response {
		alerts, _, err := cl.DeleteServiceCategory(serviceCategory.Name, toclient.RequestOptions{})
		assert.NoError(t, err, "Unexpected error deleting Service Category '%s': %v - alerts: %+v", serviceCategory.Name, err, alerts.Alerts)
		// Retrieve the Service Category to see if it got deleted
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("name", serviceCategory.Name)
		getServiceCategory, _, err := cl.GetServiceCategories(opts)
		assert.NoError(t, err, "Error getting Service Category '%s' after deletion: %v - alerts: %+v", serviceCategory.Name, err, getServiceCategory.Alerts)
		assert.Equal(t, 0, len(getServiceCategory.Response), "Expected Service Category '%s' to be deleted, but it was found in Traffic Ops", serviceCategory.Name)
	}
}
