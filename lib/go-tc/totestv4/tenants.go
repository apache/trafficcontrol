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

func CreateTestTenants(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, tenant := range td.Tenants {
		resp, _, err := cl.CreateTenant(tenant, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Could not create Tenant '%s': %v - alerts: %+v", tenant.Name, err, resp.Alerts)
	}
}

func DeleteTestTenants(t *testing.T, cl *toclient.Session) {
	opts := toclient.NewRequestOptions()
	opts.QueryParameters.Set("sortOrder", "desc")
	tenants, _, err := cl.GetTenants(opts)
	assert.NoError(t, err, "Cannot get Tenants: %v - alerts: %+v", err, tenants.Alerts)

	for _, tenant := range tenants.Response {
		if tenant.Name == "root" {
			continue
		}
		alerts, _, err := cl.DeleteTenant(tenant.ID, toclient.RequestOptions{})
		assert.NoError(t, err, "Unexpected error deleting Tenant '%s' (#%d): %v - alerts: %+v", tenant.Name, tenant.ID, err, alerts.Alerts)
		// Retrieve the Tenant to see if it got deleted
		opts.QueryParameters.Set("id", strconv.Itoa(tenant.ID))
		getTenants, _, err := cl.GetTenants(opts)
		assert.NoError(t, err, "Error getting Tenant '%s' after deletion: %v - alerts: %+v", tenant.Name, err, getTenants.Alerts)
		assert.Equal(t, 0, len(getTenants.Response), "Expected Tenant '%s' to be deleted, but it was found in Traffic Ops", tenant.Name)
	}
}
