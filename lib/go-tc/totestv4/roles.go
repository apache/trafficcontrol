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

func CreateTestRoles(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, role := range td.Roles {
		_, _, err := cl.CreateRole(role, toclient.RequestOptions{})
		assert.NoError(t, err, "No error expected, but got %v", err)
	}
}

func DeleteTestRoles(t *testing.T, cl *toclient.Session) {
	roles, _, err := cl.GetRoles(toclient.RequestOptions{})
	assert.NoError(t, err, "Cannot get Roles: %v - alerts: %+v", err, roles.Alerts)
	for _, role := range roles.Response {
		// Don't delete active roles created by test setup
		if role.Name == "admin" || role.Name == "disallowed" || role.Name == "operations" || role.Name == "portal" || role.Name == "read-only" || role.Name == "steering" || role.Name == "federation" {
			continue
		}
		_, _, err := cl.DeleteRole(role.Name, toclient.NewRequestOptions())
		assert.NoError(t, err, "Expected no error while deleting role %s, but got %v", role.Name, err)
		// Retrieve the Role to see if it got deleted
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("name", role.Name)
		getRole, _, err := cl.GetRoles(opts)
		assert.NoError(t, err, "Error getting Role '%s' after deletion: %v - alerts: %+v", role.Name, err, getRole.Alerts)
		assert.Equal(t, 0, len(getRole.Response), "Expected Role '%s' to be deleted, but it was found in Traffic Ops", role.Name)
	}
}
