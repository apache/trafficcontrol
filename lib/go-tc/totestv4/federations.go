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

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	toclient "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"
)

// TODO fix/remove globals

var fedIDs = make(map[string]int)

// All prerequisite Federations are associated to this cdn and this xmlID
const FederationCDNName = "cdn1"

var fedXmlId = "ds1"

func GetFederationID(t *testing.T, cname string) func() int {
	return func() int {
		ID, ok := fedIDs[cname]
		assert.RequireEqual(t, true, ok, "Expected to find Federation CName: %s to have associated ID", cname)
		return ID
	}
}

func setFederationID(t *testing.T, cdnFederation tc.CDNFederation) {
	assert.RequireNotNil(t, cdnFederation.CName, "Federation CName was nil after posting.")
	assert.RequireNotNil(t, cdnFederation.ID, "Federation ID was nil after posting.")
	fedIDs[*cdnFederation.CName] = *cdnFederation.ID
}

func CreateTestCDNFederations(t *testing.T, cl *toclient.Session, dat TrafficControl) {
	for _, federation := range dat.Federations {
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("xmlId", *federation.DeliveryServiceIDs.XmlId)
		dsResp, _, err := cl.GetDeliveryServices(opts)
		assert.RequireNoError(t, err, "Could not get Delivery Service by XML ID: %v", err)
		assert.RequireEqual(t, 1, len(dsResp.Response), "Expected one Delivery Service, but got %d", len(dsResp.Response))
		assert.RequireNotNil(t, dsResp.Response[0].CDNName, "Expected Delivery Service CDN Name to not be nil.")

		resp, _, err := cl.CreateCDNFederation(federation, *dsResp.Response[0].CDNName, toclient.RequestOptions{})
		assert.NoError(t, err, "Could not create CDN Federations: %v - alerts: %+v", err, resp.Alerts)

		// Need to save the ids, otherwise the other tests won't be able to reference the federations
		setFederationID(t, resp.Response)
		assert.RequireNotNil(t, resp.Response.ID, "Federation ID was nil after posting.")
		assert.RequireNotNil(t, dsResp.Response[0].ID, "Delivery Service ID was nil.")
		_, _, err = cl.CreateFederationDeliveryServices(*resp.Response.ID, []int{*dsResp.Response[0].ID}, false, toclient.NewRequestOptions())
		assert.NoError(t, err, "Could not create Federation Delivery Service: %v", err)
	}
}

func DeleteTestCDNFederations(t *testing.T, cl *toclient.Session) {
	opts := toclient.NewRequestOptions()
	for _, id := range fedIDs {
		resp, _, err := cl.DeleteCDNFederation(FederationCDNName, id, opts)
		assert.NoError(t, err, "Cannot delete federation #%d: %v - alerts: %+v", id, err, resp.Alerts)

		opts.QueryParameters.Set("id", strconv.Itoa(id))
		data, _, err := cl.GetCDNFederationsByName(FederationCDNName, opts)
		assert.Equal(t, 0, len(data.Response), "expected federation to be deleted")
	}
	fedIDs = make(map[string]int) // reset the global variable for the next test
}
