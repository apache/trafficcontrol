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
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

func TestDeliveryServiceServers(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, ServerCapabilities, ServerServerCapabilities, DeliveryServices, DeliveryServiceServerAssignments}, func() {

		tomorrow := time.Now().UTC().AddDate(0, 0, 1).Format(time.RFC1123)

		dssTests := utils.V5TestCase{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {tomorrow}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID request": {
					ClientSession: TOSession, Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when VALID CDN parameter": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"cdn": {"cdn1"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1)),
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
						"servers": []int{GetServerID(t, "test-mso-org-01")()},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusConflict)),
				},
				"CONFLICT when REMOVING ONLY ORIGIN SERVER ASSIGNMENT": {
					ClientSession: TOSession, RequestBody: map[string]interface{}{
						"dsId":    GetDeliveryServiceId(t, "test-ds-server-assignments")(),
						"replace": true,
						"servers": []int{GetServerID(t, "test-ds-server-assignments")()},
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
				"BAD REQUEST when UPDATING SERVER STATUS when ONLY ORIGIN SERVER ASSIGNED": {
					EndpointID: GetServerID(t, "test-mso-org-01"), ClientSession: TOSession,
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
							resp, reqInf, err := testCase.ClientSession.GetDeliveryServiceServers(testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.CreateDeliveryServiceServers(dsID, serverIDs, replace, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "SERVER STATUS PUT":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.UpdateServerStatus(testCase.EndpointID(), status, testCase.RequestOpts)
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

func TestDeliveryServiceXMLIDServers(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, ServerCapabilities, ServerServerCapabilities, DeliveryServices, DeliveryServiceServerAssignments}, func() {
		dsXMLIDServersTests := utils.V5TestCase{
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
				"OK when ASSIGNING ORG SERVER IN CACHEGROUP of TOPOLOGY DS": {
					ClientSession: TOSession, RequestBody: map[string]interface{}{
						"xmlID":       "ds-top",
						"serverNames": []string{"test-mso-org-01"},
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
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
							resp, reqInf, err := testCase.ClientSession.AssignServersToDeliveryService(servers, xmlID, testCase.RequestOpts)
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
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, ServerCapabilities, ServerServerCapabilities, DeliveryServices, DeliveryServiceServerAssignments}, func() {
		dsIDServersTests := utils.V5TestCase{
			"GET": {
				"OK when VALID request": {
					EndpointID: GetDeliveryServiceId(t, "test-ds-server-assignments"), ClientSession: TOSession,
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						utils.ResponseHasLength(2)),
				},
			},
		}
		for method, testCases := range dsIDServersTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetServersByDeliveryService(testCase.EndpointID(), testCase.RequestOpts)
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
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, ServerCapabilities, ServerServerCapabilities, DeliveryServices, DeliveryServiceServerAssignments}, func() {
		dssDSIDServerIDTests := utils.V5TestCase{
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
				"BAD REQUEST when REMOVING ONLY ORIGIN SERVER ASSIGNMENT": {
					ClientSession: TOSession, RequestBody: map[string]interface{}{
						"server": GetServerID(t, "test-mso-org-01")(),
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
							resp, reqInf, err := testCase.ClientSession.DeleteDeliveryServiceServer(dsID, serverId, testCase.RequestOpts)
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
	dsServers, _, err := TOSession.GetDeliveryServiceServers(client.RequestOptions{})
	assert.NoError(t, err, "Unexpected error retrieving server-to-Delivery-Service assignments: %v - alerts: %+v", err, dsServers.Alerts)

	for _, dss := range dsServers.Response {
		// Retrieve Delivery Service in order to update its active field to false
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.Itoa(*dss.DeliveryService))
		getDS, _, err := TOSession.GetDeliveryServices(opts)
		assert.NoError(t, err, "Error retrieving Delivery Service: %v - alerts: %+v", err, getDS.Alerts)
		assert.Equal(t, 1, len(getDS.Response), "Expected 1 Delivery Service.")
		// Update active to false in order to remove the server assignment
		getDS.Response[0].Active = tc.DSActiveStateInactive
		updResp, _, err := TOSession.UpdateDeliveryService(*dss.DeliveryService, getDS.Response[0], client.RequestOptions{})
		assert.NoError(t, err, "Error updating Delivery Service: %v - alerts: %+v", err, updResp.Alerts)
		assert.Equal(t, tc.DSActiveStateInactive, updResp.Response.Active, "Expected Delivery Service to be Inactive.")

		alerts, _, err := TOSession.DeleteDeliveryServiceServer(*dss.DeliveryService, *dss.Server, client.RequestOptions{})
		assert.NoError(t, err, "Unexpected error removing server-to-Delivery-Service assignments: %v - alerts: %+v", err, alerts.Alerts)
	}
	dsServers, _, err = TOSession.GetDeliveryServiceServers(client.RequestOptions{})
	assert.NoError(t, err, "Unexpected error retrieving server-to-Delivery-Service assignments: %v - alerts: %+v", err, dsServers.Alerts)
	assert.Equal(t, dsServers.Size, 0, "Expected all Delivery Service Server assignments to be deleted.")
}

func CreateTestDeliveryServiceServerAssignments(t *testing.T) {
	for _, dss := range testData.DeliveryServiceServerAssignments {
		resp, _, err := TOSession.AssignServersToDeliveryService(dss.ServerNames, dss.XmlId, client.RequestOptions{})
		assert.NoError(t, err, "Could not create Delivery Service Server Assignments: %v - alerts: %+v", err, resp.Alerts)
	}
}
