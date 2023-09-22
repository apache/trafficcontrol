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
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
)

func TestDeliveryServiceServers(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, ServerCapabilities, DeliveryServices, ServerServerCapabilities, DeliveryServicesRequiredCapabilities, DeliveryServiceServerAssignments}, func() {

		tomorrow := time.Now().UTC().AddDate(0, 0, 1).Format(time.RFC1123)

		dssTests := utils.V3TestCase{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession: TOSession, RequestHeaders: http.Header{rfc.IfModifiedSince: {tomorrow}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID request": {
					ClientSession: TOSession, Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
			"POST": {
				"OK when VALID request": {
					ClientSession: TOSession, RequestBody: map[string]interface{}{
						"dsId":    GetDeliveryServiceId(t, "ds3")(),
						"replace": true,
						"servers": []int{GetServerID(t, "atlanta-edge-01")(), GetServerID(t, "atlanta-edge-03")()},
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when ASSIGNING ORG SERVER IN CACHEGROUP of TOPOLOGY DS": {
					ClientSession: TOSession, RequestBody: map[string]interface{}{
						"dsId":    GetDeliveryServiceId(t, "ds-top")(),
						"replace": true,
						"servers": []int{GetServerID(t, "denver-mso-org-01")()},
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"BAD REQUEST when ASSIGNING ORG SERVER NOT IN CACHEGROUP of TOPOLOGY DS": {
					ClientSession: TOSession, RequestBody: map[string]interface{}{
						"dsId":    GetDeliveryServiceId(t, "ds-top-req-cap")(),
						"replace": true,
						"servers": []int{GetServerID(t, "denver-mso-org-01")()},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when ASSIGNING SERVERS to a TOPOLOGY DS": {
					ClientSession: TOSession, RequestBody: map[string]interface{}{
						"dsId":    GetDeliveryServiceId(t, "ds-top")(),
						"servers": []int{GetServerID(t, "atlanta-edge-01")(), GetServerID(t, "atlanta-edge-03")()},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"CONFLICT when REMOVING ONLY EDGE SERVER ASSIGNMENT": {
					ClientSession: TOSession, RequestBody: map[string]interface{}{
						"dsId":    GetDeliveryServiceId(t, "test-ds-server-assignments")(),
						"replace": true,
						"servers": []int{},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusConflict)),
				},
				"CONFLICT when REPLACING EDGE SERVER ASSIGNMENT with EDGE SERVER in BAD STATE": {
					ClientSession: TOSession, RequestBody: map[string]interface{}{
						"dsId":    GetDeliveryServiceId(t, "test-ds-server-assignments")(),
						"replace": true,
						"servers": []int{GetServerID(t, "admin-down-server")()},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusConflict)),
				},
				"OK when MAKING ASSIGNMENTS when DELIVERY SERVICE AND SERVER HAVE MATCHING CAPABILITIES": {
					ClientSession: TOSession, RequestBody: map[string]interface{}{
						"dsId":    GetDeliveryServiceId(t, "ds2")(),
						"replace": true,
						"servers": []int{GetServerID(t, "atlanta-org-2")()},
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when ASSIGNING a ORIGIN server to a DS with REQUIRED CAPABILITY": {
					ClientSession: TOSession, RequestBody: map[string]interface{}{
						"dsId":    GetDeliveryServiceId(t, "msods1")(),
						"replace": true,
						"servers": []int{GetServerID(t, "denver-mso-org-01")()},
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
			"SERVER STATUS PUT": {
				"BAD REQUEST when UPDATING SERVER STATUS when ONLY EDGE SERVER ASSIGNED": {
					EndpointID: GetServerID(t, "test-ds-server-assignments"), ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"status":        "ADMIN_DOWN",
						"offlineReason": "admin down",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusConflict)),
				},
			},
		}
		for method, testCases := range dssTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					var dsID int
					var replace bool
					var serverIDs []int
					status := tc.ServerPutStatus{}

					if testCase.RequestBody != nil {
						if val, ok := testCase.RequestBody["dsId"]; ok {
							dsID = val.(int)
						}
						if val, ok := testCase.RequestBody["replace"]; ok {
							replace = val.(bool)
						}
						if val, ok := testCase.RequestBody["servers"]; ok {
							serverIDs = val.([]int)
						}
						if _, ok := testCase.RequestBody["offlineReason"]; ok {
							dat, err := json.Marshal(testCase.RequestBody)
							assert.NoError(t, err, "Error occurred when marshalling request body: %v", err)
							err = json.Unmarshal(dat, &status)
							assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
						}
					}

					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetDeliveryServiceServersWithHdr(testCase.RequestHeaders)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.CreateDeliveryServiceServers(dsID, serverIDs, replace)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp, tc.Alerts{}, err)
							}
						})
					case "SERVER STATUS PUT":
						t.Run(name, func(t *testing.T) {
							_, reqInf, err := testCase.ClientSession.UpdateServerStatus(testCase.EndpointID(), status)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, tc.Alerts{}, err)
							}
						})
					}
				}
			})
		}
	})
}

func TestDeliveryServiceXMLIDServers(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, DeliveryServices, DeliveryServiceServerAssignments}, func() {
		dsXMLIDServersTests := utils.V3TestCase{
			"POST": {
				"BAD REQUEST when ASSIGNING SERVERS to a TOPOLOGY DS": {
					ClientSession: TOSession, RequestBody: map[string]interface{}{
						"xmlID":       "ds-top",
						"serverNames": []string{"atlanta-edge-01", "atlanta-edge-03"},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when ASSIGNING ORG SERVER NOT IN CACHEGROUP of TOPOLOGY DS": {
					ClientSession: TOSession, RequestBody: map[string]interface{}{
						"xmlID":       "ds-top-req-cap",
						"serverNames": []string{"denver-mso-org-01"},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
		}
		for method, testCases := range dsXMLIDServersTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					var xmlID string
					var servers []string

					if testCase.RequestBody != nil {
						if val, ok := testCase.RequestBody["xmlID"]; ok {
							xmlID = val.(string)
						}
						if val, ok := testCase.RequestBody["serverNames"]; ok {
							servers = val.([]string)
						}
					}

					switch method {
					case "POST":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.AssignServersToDeliveryService(servers, xmlID)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, resp, err)
							}
						})
					}
				}
			})
		}
	})
}

func TestDeliveryServicesIDServers(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, DeliveryServices, DeliveryServiceServerAssignments}, func() {
		dsIDServersTests := utils.V3TestCase{
			"GET": {
				"OK when VALID request": {
					EndpointID: GetDeliveryServiceId(t, "test-ds-server-assignments"), ClientSession: TOSession,
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						utils.ResponseHasLength(1)),
				},
			},
		}
		for method, testCases := range dsIDServersTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetServersByDeliveryService(testCase.EndpointID())
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

func TestDeliveryServicesDSIDServerID(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, DeliveryServices, DeliveryServiceServerAssignments}, func() {
		dssDSIDServerIDTests := utils.V3TestCase{
			"DELETE": {
				"OK when VALID REQUEST": {
					ClientSession: TOSession, RequestBody: map[string]interface{}{
						"server": GetServerID(t, "denver-mso-org-01")(),
						"dsId":   GetDeliveryServiceId(t, "ds-top")(),
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"BAD REQUEST when REMOVING ONLY EDGE SERVER ASSIGNMENT": {
					ClientSession: TOSession, RequestBody: map[string]interface{}{
						"server": GetServerID(t, "test-ds-server-assignments")(),
						"dsId":   GetDeliveryServiceId(t, "test-ds-server-assignments")(),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusConflict)),
				},
			},
		}
		for method, testCases := range dssDSIDServerIDTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					var dsID int
					var serverId int

					if testCase.RequestBody != nil {
						if val, ok := testCase.RequestBody["server"]; ok {
							serverId = val.(int)
						}
						if val, ok := testCase.RequestBody["dsId"]; ok {
							dsID = val.(int)
						}
					}

					switch method {
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.DeleteDeliveryServiceServer(dsID, serverId)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, resp, err)
							}
						})
					}
				}
			})
		}
	})
}

func DeleteTestDeliveryServiceServers(t *testing.T) {
	dsServers, _, err := TOSession.GetDeliveryServiceServersWithHdr(nil)
	assert.NoError(t, err, "Unexpected error retrieving server-to-Delivery-Service assignments: %v - alerts: %+v", err, dsServers.Alerts)

	for _, dss := range dsServers.Response {
		// Retrieve Delivery Service in order to update its active field to false
		params := url.Values{"id": {strconv.Itoa(*dss.DeliveryService)}}
		getDS, _, err := TOSession.GetDeliveryServicesV30WithHdr(nil, params)
		assert.NoError(t, err, "Error retrieving Delivery Service: %v - resp: %+v", err, getDS)
		assert.Equal(t, 1, len(getDS), "Expected 1 Delivery Service.")
		// Update active to false in order to remove the server assignment
		active := false
		getDS[0].Active = &active
		updResp, _, err := TOSession.UpdateDeliveryServiceV30WithHdr(*dss.DeliveryService, getDS[0], nil)
		assert.NoError(t, err, "Error updating Delivery Service: %v - resp: %+v", err, updResp)
		assert.Equal(t, false, *updResp.Active, "Expected Delivery Service to be Inactive.")

		alerts, _, err := TOSession.DeleteDeliveryServiceServer(*dss.DeliveryService, *dss.Server)
		assert.NoError(t, err, "Unexpected error removing server-to-Delivery-Service assignments: %v - alerts: %+v", err, alerts.Alerts)
	}
	dsServers, _, err = TOSession.GetDeliveryServiceServersWithHdr(nil)
	assert.NoError(t, err, "Unexpected error retrieving server-to-Delivery-Service assignments: %v - alerts: %+v", err, dsServers.Alerts)
	assert.Equal(t, dsServers.Size, 0, "Expected all Delivery Service Server assignments to be deleted.")
}

func CreateTestDeliveryServiceServerAssignments(t *testing.T) {
	for _, dss := range testData.DeliveryServiceServerAssignments {
		resp, _, err := TOSession.AssignServersToDeliveryService(dss.ServerNames, dss.XmlId)
		assert.NoError(t, err, "Could not create Delivery Service Server Assignments: %v - alerts: %+v", err, resp.Alerts)
	}
}
