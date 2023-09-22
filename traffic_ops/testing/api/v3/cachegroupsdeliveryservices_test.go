package v3

/*

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

import (
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

func TestCacheGroupsDeliveryServices(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, DeliveryServices, CacheGroupsDeliveryServices}, func() {

		methodTests := utils.V3TestCaseT[[]int]{
			"POST": {
				"BAD REQUEST assigning TOPOLOGY-BASED DS to CACHEGROUP": {
					EndpointID:    GetCacheGroupId(t, "cachegroup3"),
					ClientSession: TOSession,
					RequestBody:   []int{GetDeliveryServiceId(t, "top-ds-in-cdn1")()},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"OK when valid request": {
					EndpointID:    GetCacheGroupId(t, "cachegroup3"),
					ClientSession: TOSession,
					RequestBody: []int{
						GetDeliveryServiceId(t, "ds1")(),
						GetDeliveryServiceId(t, "ds2")(),
						GetDeliveryServiceId(t, "ds3")(),
						GetDeliveryServiceId(t, "ds3")(),
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateCGDSServerAssignments()),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "POST":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.SetCachegroupDeliveryServices(testCase.EndpointID(), testCase.RequestBody)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					}
				}
			})
		}

	})
}

func validateCGDSServerAssignments() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		cgDsResp := resp.(tc.CacheGroupPostDSResp)
		params := url.Values{}
		for _, serverName := range cgDsResp.ServerNames {
			params.Set("hostName", string(serverName))
			resp, _, err := TOSession.GetServersWithHdr(&params, nil)
			assert.NoError(t, err, "Error: Getting server: %v - alerts: %+v", err, resp.Alerts)
			assert.Equal(t, len(resp.Response), 1, "Error: Getting servers: expected 1 got %v", len(resp.Response))

			serverDSes, _, err := TOSession.GetDeliveryServicesByServerV30WithHdr(*resp.Response[0].ID, nil)
			assert.NoError(t, err, "Error: Getting Delivery Service Servers #%d: %v", *resp.Response[0].ID, err)
			for _, dsID := range cgDsResp.DeliveryServices {
				found := false
				for _, serverDS := range serverDSes {
					if *serverDS.ID == dsID {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("POST succeeded, but didn't assign delivery service %v to server", dsID)
				}
			}
		}
	}
}

func CreateTestCachegroupsDeliveryServices(t *testing.T) {
	dses, _, err := TOSession.GetDeliveryServicesV30WithHdr(nil, nil)
	assert.RequireNoError(t, err, "Cannot GET DeliveryServices: %v - %v", err, dses)

	clientCGs, _, err := TOSession.GetCacheGroupNullableByNameWithHdr("cachegroup3", nil)
	assert.RequireNoError(t, err, "Cannot GET cachegroup: %v", err)
	assert.RequireEqual(t, len(clientCGs), 1, "Getting cachegroup expected 1, got %v", len(clientCGs))
	assert.RequireNotNil(t, clientCGs[0].ID, "Cachegroup has a nil ID")

	dsIDs := []int{}
	for _, ds := range dses {
		if *ds.CDNName == "cdn1" && ds.Topology == nil {
			dsIDs = append(dsIDs, *ds.ID)
		}
	}
	assert.RequireGreaterOrEqual(t, len(dsIDs), 1, "No Delivery Services found in CDN 'cdn1', cannot continue.")
	resp, _, err := TOSession.SetCachegroupDeliveryServices(*clientCGs[0].ID, dsIDs)
	assert.RequireNoError(t, err, "Setting cachegroup delivery services returned error: %v", err)
	assert.RequireGreaterOrEqual(t, len(resp.Response.ServerNames), 1, "Setting cachegroup delivery services returned success, but no servers set")
}

func setInactive(t *testing.T, dsID int) {
	opts := url.Values{}
	opts.Set("id", strconv.Itoa(dsID))
	resp, _, err := TOSession.GetDeliveryServicesV30WithHdr(nil, opts)
	assert.RequireNoError(t, err, "Failed to fetch details for Delivery Service #%d: %v", dsID, err)
	assert.RequireEqual(t, len(resp), 1, "Expected exactly one Delivery Service to exist with ID %d, found: %d", dsID, len(resp))

	ds := resp[0]
	if ds.Active == nil {
		t.Errorf("Deliver Service #%d had null or undefined 'active'", dsID)
		ds.Active = new(bool)
	}
	if *ds.Active {
		*ds.Active = false
		_, _, err = TOSession.UpdateDeliveryServiceV30WithHdr(dsID, ds, nil)
		assert.NoError(t, err, "Failed to set Delivery Service #%d to inactive: %v", dsID, err)
	}
}

func DeleteTestCachegroupsDeliveryServices(t *testing.T) {
	dss, _, err := TOSession.GetDeliveryServiceServersNWithHdr(1000000, nil)
	assert.NoError(t, err, "Unexpected error retrieving server-to-Delivery-Service assignments: %v - alerts: %+v", err, dss.Alerts)

	for _, ds := range dss.Response {
		setInactive(t, *ds.DeliveryService)
		alerts, _, err := TOSession.DeleteDeliveryServiceServer(*ds.DeliveryService, *ds.Server)
		assert.NoError(t, err, "Error deleting delivery service servers: %v - alerts: %+v", err, alerts.Alerts)
	}

	dss, _, err = TOSession.GetDeliveryServiceServersWithHdr(nil)
	assert.NoError(t, err, "Unexpected error retrieving server-to-Delivery-Service assignments: %v - alerts: %+v", err, dss.Alerts)
	assert.Equal(t, len(dss.Response), 0, "Deleting delivery service servers: Expected empty subsequent get, actual %v", len(dss.Response))
}
