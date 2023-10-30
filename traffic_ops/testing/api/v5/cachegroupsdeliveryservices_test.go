package v5

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
	"strconv"
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

func TestCacheGroupsDeliveryServices(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, ServerCapabilities, DeliveryServices, CacheGroupsDeliveryServices}, func() {

		methodTests := utils.TestCase[client.Session, client.RequestOptions, []int]{
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
						GetDeliveryServiceId(t, "ds3")(),
						GetDeliveryServiceId(t, "ds4")(),
						GetDeliveryServiceId(t, "DS5")(),
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
							resp, reqInf, err := testCase.ClientSession.SetCacheGroupDeliveryServices(testCase.EndpointID(), testCase.RequestBody, testCase.RequestOpts)
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
		opts := client.NewRequestOptions()
		for _, serverName := range cgDsResp.ServerNames {
			opts.QueryParameters.Set("hostName", string(serverName))
			resp, _, err := TOSession.GetServers(opts)
			assert.NoError(t, err, "Error: Getting server: %v - alerts: %+v", err, resp.Alerts)
			assert.Equal(t, len(resp.Response), 1, "Error: Getting servers: expected 1 got %v", len(resp.Response))

			serverDSes, _, err := TOSession.GetDeliveryServicesByServer(resp.Response[0].ID, client.RequestOptions{})
			assert.NoError(t, err, "Error: Getting Delivery Service Servers #%d: %v - alerts: %+v", resp.Response[0].ID, err, serverDSes.Alerts)
			for _, dsID := range cgDsResp.DeliveryServices {
				found := false
				for _, serverDS := range serverDSes.Response {
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
	dses, _, err := TOSession.GetDeliveryServices(client.RequestOptions{})
	assert.RequireNoError(t, err, "Cannot GET DeliveryServices: %v - %v", err, dses)

	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", "cachegroup3")
	clientCGs, _, err := TOSession.GetCacheGroups(opts)
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
	resp, _, err := TOSession.SetCacheGroupDeliveryServices(*clientCGs.Response[0].ID, dsIDs, client.RequestOptions{})
	assert.RequireNoError(t, err, "Setting cachegroup delivery services returned error: %v", err)
	assert.RequireGreaterOrEqual(t, len(resp.Response.ServerNames), 1, "Setting cachegroup delivery services returned success, but no servers set")
}

func setInactive(t *testing.T, dsID int) {
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("id", strconv.Itoa(dsID))
	resp, _, err := TOSession.GetDeliveryServices(opts)
	assert.RequireNoError(t, err, "Failed to fetch details for Delivery Service #%d: %v - alerts: %+v", dsID, err, resp.Alerts)
	assert.RequireEqual(t, len(resp.Response), 1, "Expected exactly one Delivery Service to exist with ID %d, found: %d", dsID, len(resp.Response))

	ds := resp.Response[0]
	if ds.Active == tc.DSActiveStateActive {
		ds.Active = tc.DSActiveStateInactive
		_, _, err = TOSession.UpdateDeliveryService(dsID, ds, client.RequestOptions{})
		assert.NoError(t, err, "Failed to set Delivery Service #%d to inactive: %v", dsID, err)
	}
}

func DeleteTestCachegroupsDeliveryServices(t *testing.T) {
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("limit", "1000000")
	dss, _, err := TOSession.GetDeliveryServiceServers(opts)
	assert.NoError(t, err, "Unexpected error retrieving server-to-Delivery-Service assignments: %v - alerts: %+v", err, dss.Alerts)

	for _, ds := range dss.Response {
		setInactive(t, *ds.DeliveryService)
		alerts, _, err := TOSession.DeleteDeliveryServiceServer(*ds.DeliveryService, *ds.Server, client.RequestOptions{})
		assert.NoError(t, err, "Error deleting delivery service servers: %v - alerts: %+v", err, alerts.Alerts)
	}

	dss, _, err = TOSession.GetDeliveryServiceServers(client.RequestOptions{})
	assert.NoError(t, err, "Unexpected error retrieving server-to-Delivery-Service assignments: %v - alerts: %+v", err, dss.Alerts)
	assert.Equal(t, len(dss.Response), 0, "Deleting delivery service servers: Expected empty subsequent get, actual %v", len(dss.Response))
}
