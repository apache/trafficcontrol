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
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

func TestServersIDDeliveryServices(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Tenants, Topologies, ServiceCategories, DeliveryServices, DeliveryServiceServerAssignments}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.V3TestCaseT[map[string]interface{}]{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					EndpointID:     GetServerID(t, "atlanta-edge-14"),
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfModifiedSince: {tomorrow}},
					Expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID request": {
					EndpointID:    GetServerID(t, "atlanta-edge-14"),
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
			"POST": {
				"OK when VALID request": {
					EndpointID:    GetServerID(t, "atlanta-edge-01"),
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"dsIds":   []int{GetDeliveryServiceId(t, "ds1")()},
						"replace": true,
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateServersDeliveryServicesPost(
							GetServerID(t, "atlanta-edge-01")(),
							[]int{
								GetDeliveryServiceId(t, "ds1")(),
								GetDeliveryServiceId(t, "ds-based-top-with-no-mids")(),
							},
							2)),
				},
				"OK when ASSIGNING EDGE to TOPOLOGY BASED DELIVERY SERVICE": {
					EndpointID:    GetServerID(t, "atlanta-edge-03"),
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"dsIds":   []int{GetDeliveryServiceId(t, "top-ds-in-cdn1")()},
						"replace": true,
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateServersDeliveryServicesPost(
							GetServerID(t, "atlanta-edge-03")(),
							[]int{
								GetDeliveryServiceId(t, "top-ds-in-cdn1")(),
							},
							1)),
				},
				"OK when ASSIGNING ORIGIN to TOPOLOGY BASED DELIVERY SERVICE": {
					EndpointID:    GetServerID(t, "denver-mso-org-01"),
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"dsIds":   []int{GetDeliveryServiceId(t, "ds-top")()},
						"replace": true,
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateServersDeliveryServicesPost(
							GetServerID(t, "denver-mso-org-01")(),
							[]int{
								GetDeliveryServiceId(t, "ds-top")(),
							},
							1)),
				},
				"CONFLICT when SERVER NOT IN SAME CDN as DELIVERY SERVICE": {
					EndpointID:    GetServerID(t, "cdn2-test-edge"),
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"dsIds":   []int{GetDeliveryServiceId(t, "ds1")()},
						"replace": true,
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusConflict)),
				},
				"BAD REQUEST when ORIGIN'S CACHEGROUP IS NOT A PART OF TOPOLOGY BASED DELIVERY SERVICE": {
					EndpointID:    GetServerID(t, "denver-mso-org-01"),
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"dsIds":   []int{GetDeliveryServiceId(t, "ds-top-req-cap")()},
						"replace": true,
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"CONFLICT when REMOVING ONLY EDGE SERVER ASSIGNMENT": {
					EndpointID:    GetServerID(t, "test-ds-server-assignments"),
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"dsIds":   []int{},
						"replace": true,
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusConflict)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {

					var dsIds []int
					var replace bool

					if testCase.RequestBody != nil {
						if val, ok := testCase.RequestBody["dsIds"]; ok {
							dsIds = val.([]int)
						}
						if val, ok := testCase.RequestBody["replace"]; ok {
							replace = val.(bool)
						}
					}

					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetServerIDDeliveryServicesWithHdr(testCase.EndpointID(), testCase.RequestHeaders)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp, tc.Alerts{}, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.AssignDeliveryServiceIDsToServerID(testCase.EndpointID(), dsIds, replace)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					}
				}
			})
		}
	})
}

func validateServersDeliveryServices(expectedDSID int) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Server Delivery Service response to not be nil.")
		var found bool
		deliveryServices := resp.([]tc.DeliveryServiceNullable)
		for _, ds := range deliveryServices {
			if ds.ID != nil && *ds.ID == expectedDSID {
				found = true
				break
			}
		}
		assert.Equal(t, true, found, "Expected to find Delivery Service ID: %d in response.")
	}
}

func validateServersDeliveryServicesPost(serverID int, expectedDSID []int, expectedDSCount int) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		serverDeliveryServices, _, err := TOSession.GetServerIDDeliveryServicesWithHdr(serverID, nil)
		assert.RequireNoError(t, err, "Error getting Server Delivery Services: %v", err)
		assert.RequireEqual(t, expectedDSCount, len(serverDeliveryServices), "Expected %d Delivery Service returned Got: %d", expectedDSCount, len(serverDeliveryServices))
		for i := 0; i < len(expectedDSID); i++ {
			validateServersDeliveryServices(expectedDSID[i])(t, toclientlib.ReqInf{}, serverDeliveryServices, tc.Alerts{}, nil)
		}
	}
}
