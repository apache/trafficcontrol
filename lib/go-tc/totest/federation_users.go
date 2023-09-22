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

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	toclient "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

func CreateTestFederationUsers(t *testing.T, cl *toclient.Session) {
	// Prerequisite Federation Users
	federationUsers := map[string]tc.FederationUserPost{
		"the.cname.com.": {
			IDs:     []int{GetUserID(t, cl, "admin")(), GetUserID(t, cl, "adminuser")(), GetUserID(t, cl, "disalloweduser")(), GetUserID(t, cl, "readonlyuser")()},
			Replace: util.Ptr(false),
		},
		"booya.com.": {
			IDs:     []int{GetUserID(t, cl, "adminuser")()},
			Replace: util.Ptr(false),
		},
	}

	for cname, federationUser := range federationUsers {
		fedID := GetFederationID(t, cname)()
		resp, _, err := cl.CreateFederationUsers(fedID, federationUser.IDs, *federationUser.Replace, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Assigning users %v to federation %d: %v - alerts: %+v", federationUser.IDs, fedID, err, resp.Alerts)
	}
}

func DeleteTestFederationUsers(t *testing.T, cl *toclient.Session) {
	for _, fedID := range fedIDs {
		fedUsers, _, err := cl.GetFederationUsers(fedID, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Error getting users for federation %d: %v - alerts: %+v", fedID, err, fedUsers.Alerts)
		for _, fedUser := range fedUsers.Response {
			if fedUser.ID == nil {
				t.Error("Traffic Ops returned a representation of a relationship between a user and a Federation that had null or undefined ID")
				continue
			}
			alerts, _, err := cl.DeleteFederationUser(fedID, *fedUser.ID, toclient.RequestOptions{})
			assert.NoError(t, err, "Error deleting user #%d from federation #%d: %v - alerts: %+v", *fedUser.ID, fedID, err, alerts.Alerts)
		}
		fedUsers, _, err = cl.GetFederationUsers(fedID, toclient.RequestOptions{})
		assert.NoError(t, err, "Error getting users for federation %d: %v - alerts: %+v", fedID, err, fedUsers.Alerts)
		assert.Equal(t, 0, len(fedUsers.Response), "Federation users expected 0, actual: %+v", len(fedUsers.Response))
	}
}
