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
	"strconv"
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	toclient "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

func CreateTestFederationResolvers(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, fr := range td.FederationResolvers {
		fr.TypeID = util.Ptr(uint(GetTypeId(t, cl, *fr.Type)))
		resp, _, err := cl.CreateFederationResolver(fr, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Failed to create Federation Resolver %+v: %v - alerts: %+v", fr, err, resp.Alerts)
	}
}

func DeleteTestFederationResolvers(t *testing.T, cl *toclient.Session) {
	frs, _, err := cl.GetFederationResolvers(toclient.RequestOptions{})
	assert.RequireNoError(t, err, "Unexpected error getting Federation Resolvers: %v - alerts: %+v", err, frs.Alerts)
	for _, fr := range frs.Response {
		alerts, _, err := cl.DeleteFederationResolver(*fr.ID, toclient.RequestOptions{})
		assert.NoError(t, err, "Failed to delete Federation Resolver %+v: %v - alerts: %+v", fr, err, alerts.Alerts)
		// Retrieve the Federation Resolver to see if it got deleted
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.Itoa(int(*fr.ID)))
		getFR, _, err := cl.GetFederationResolvers(opts)
		assert.NoError(t, err, "Error getting Federation Resolver '%d' after deletion: %v - alerts: %+v", *fr.ID, err, getFR.Alerts)
		assert.Equal(t, 0, len(getFR.Response), "Expected Federation Resolver '%d' to be deleted, but it was found in Traffic Ops", *fr.ID)
	}
}

func GetFederationResolverID(t *testing.T, cl *toclient.Session, ipAddress string) func() int {
	return func() int {
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("ipAddress", ipAddress)
		federationResolvers, _, err := cl.GetFederationResolvers(opts)
		assert.RequireNoError(t, err, "Get FederationResolvers Request failed with error:", err)
		assert.RequireEqual(t, 1, len(federationResolvers.Response), "Expected response object length 1, but got %d", len(federationResolvers.Response))
		assert.RequireNotNil(t, federationResolvers.Response[0].ID, "Expected Federation Resolver ID to not be nil")
		return int(*federationResolvers.Response[0].ID)
	}
}
