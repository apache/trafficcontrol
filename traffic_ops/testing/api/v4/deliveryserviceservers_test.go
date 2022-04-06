package v4

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
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/utils"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

//----------------------------------------------------------------------------------------------------------------
// get a bad status i.e. a status thats not online or reported
// PREREQ ASSIGNMENTS
// xmlId = test-ds-server-assignments
// edge server = test-ds-server-assignments
// org server = test-mso-org-01

//-----------SERVERS-------------------------------------------------------------------------------
// DELETE SERVER /server/id WHERE server ID is for edge server = test-ds-server-assignments
// Expect error removing edge server

// DELETE SERVER /server/id WHERE server ID is for origin server =  test-mso-org-01
// Expect error removing origin server

// PUT server/id/status to a "BAD" status // server = test-ds-server-as
//  Expect error updating state of server

// PUT server/id/status to a "BAD" status //  server = test-mso-org-01
//  Expect error updating state of server

// PUT server/id update status id to a "BAD status // server = test-ds-server-as
//  Expect error updating state of server

// PUT server/id  update status id to a "BAD status // server = test-mso-org-01
//  Expect error updating state of server
//----------------------------------------------------------------------------------------------------------------

// POST /deliveryserviceserver ds in cdn1 with cdn1 edges (PREREQS)
// POST /deliveryserviceserver ds in cdn2 with cdn2 edges (PREREQS)
// GET  /deliveryserviceserver with cdn1 parameters
// GET  /deliveryserviceserver with cdn2 parameters

func TestDeliveryServiceServers(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, DeliveryServices}, func() {

		tomorrow := time.Now().UTC().AddDate(0, 0, 1).Format(time.RFC1123)

		// Tests for the /deliveryservicesserver route
		dssTests := utils.V4TestCase{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {tomorrow}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID request": {
					ClientSession: TOSession, Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				// GET deliveryservices using CDN Param // length should match // cdn name should match
				"OK when VALID CDN parameter": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"cdn": {"cdn1"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
			"POST": {
				"OK when VALID request": {
					ClientSession: TOSession, RequestBody: map[string]interface{}{
						"dsId":    GetDeliveryServiceId(t, "ds1")(),
						"replace": true,
						"servers":  []int{getServerID(t, "atlanta-edge-01")(), getServerID(t, "influxdb02")(),
							getServerID(t, "atlanta-router-01")(), getServerID(t, "atlanta-edge-03")(),
							getServerID(t, "atlanta-edge-14")(), getServerID(t, "atlanta-edge-15")(),
							getServerID(t, "edge1-cdn1-cg3")(), getServerID(t, "edge2-cdn1-cg3")(),
							getServerID(t, "dtrc-edge-01")(), getServerID(t, "dtrc-edge-02")(),
							getServerID(t, "dtrc-edge-03")(), getServerID(t, "dtrc-edge-04")(),
							getServerID(t, "dtrc-edge-05")(), getServerID(t, "dtrc-edge-06")(),
							getServerID(t, "dtrc-edge-07")(), getServerID(t, "dtrc-edge-08")()},
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when ASSIGNING ORG SERVER IN CACHEGROUP of TOPOLOGY DS": {
					ClientSession: TOSession, RequestBody: map[string]interface{}{
						"dsId":    GetDeliveryServiceId(t, "ds-top")(),
						"replace": true,
						"servers": []int{getServerID(t, "atlanta-edge-01")(), getServerID(t, "denver-mso-org-01")()},
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"BAD REQUEST when ASSIGNING ORG SERVER NOT IN CACHEGROUP of TOPOLOGY DS": {
					ClientSession: TOSession, RequestBody: map[string]interface{}{
						"dsId":    GetDeliveryServiceId(t, "ds-top-req-cap")(),
						"replace": true,
						"servers": []int{getServerID(t, "denver-mso-org-01")()},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when ASSIGNING SERVERS to a TOPOLOGY DS": {
					ClientSession: TOSession, RequestBody: map[string]interface{}{
						"dsId":    "ds-top",
						"servers": []int{getServerID(t, "atlanta-edge-01")(), getServerID(t, "atlanta-edge-03")()},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when REMOVING ONLY EDGE SERVER ASSIGNMENT": {
					ClientSession: TOSession, RequestBody: map[string]interface{}{
						"dsId":    "test-ds-server-assignments",
						"replace": true,
						"servers": []int{getServerID(t, "test-mso-org-01")()},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when REMOVING ONLY ORIGIN SERVER ASSIGNMENT": {
					ClientSession: TOSession, RequestBody: map[string]interface{}{
						"dsId":    "test-ds-server-assignments",
						"replace": true,
						"servers": []int{getServerID(t, "test-ds-server-assignments")()},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when REPLACING EDGE SERVER ASSIGNMENT with EDGE SERVER in BAD STATE": {
					ClientSession: TOSession, RequestBody: map[string]interface{}{
						"dsId":    "test-ds-server-assignments",
						"replace": true,
						"servers": []int{getServerID(t, "admin-down-server")()},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				// server="atlanta-edge-01" xmlId="ds-test-minor-versions" // use testdata serverserv
				// Create DSRC for ds-test-minor-versions with serverservercap at index 1
				// Create SSC for  atlanta-edge-01 with serverservercap at index 1
				// Create DSS between ds and server expect no failure
				// SAME TEST BELOW
				// server="atlanta-mid-01" xmlId="ds3" // use testdata serverservercap at index 1
				"OK when MAKING ASSIGNMENTS when DELIVERY SERVICE AND SERVER HAVE MATCHING CAPABILITIES": {
					ClientSession: TOSession, RequestBody: map[string]interface{}{
						"dsId":    GetDeliveryServiceId(t, "ds-top")(),
						"replace": true,
						"servers": []int{getServerID(t, "atlanta-edge-01")(), getServerID(t, "denver-mso-org-01")()},
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				// NO ERROR WHEN ASSIGNING A DSRC to a DS with an existing DSS assignment <--- DSRC TEST
				"OK when ASSIGNING a ORIGIN server to a DS with REQUIRED CAPABILITY": {
					ClientSession: TOSession, RequestBody: map[string]interface{}{
						"dsId":    GetDeliveryServiceId(t, "msods1")(),
						"replace": true,
						"servers": []int{getServerID(t, "denver-mso-org-01")()},
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
		}

		// Tests for the /deliveryservices/{xmlId}/servers route
		dsXMLIDServersTests := utils.V4TestCase{
			"GET": {
				// should match expected assignment // GetServersByDeliveryService // PREREQ have assignment made
				"OK when VALID request": {
					ClientSession: TOSession, Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						utils.ResponseHasLength(1)),
				},
			},
			"POST": {
				// Using ds-top // server must belong in same cdn and be an edge cache // Multiple server assignments
				"BAD REQUEST when ASSIGNING SERVERS to a TOPOLOGY DS": {
					ClientSession: TOSession, RequestBody: map[string]interface{}{
						"xmlId":       "ds-top",
						"serverNames": []string{"atlanta-edge-01", "atlanta-edge-03"},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when ASSIGNING ORG SERVER NOT IN CACHEGROUP of TOPOLOGY DS": {
					ClientSession: TOSession, RequestBody: map[string]interface{}{
						"dsId":    GetDeliveryServiceId(t, "ds-top-req-cap")(),
						"replace": true,
						"servers": []string{"denver-mso-org-01"},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"OK when ASSIGNING ORG SERVER IN CACHEGROUP of TOPOLOGY DS": {
					ClientSession: TOSession, RequestBody: map[string]interface{}{
						"dsId":    GetDeliveryServiceId(t, "ds-top")(),
						"replace": true,
						"servers": []string{"denver-mso-org-01"},
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
		}

		// Tests for the /servers/%d/deliveryservices route
		serversIDDSTests := utils.V4TestCase{
			"POST": {
				"BAD REQUEST when REMOVING ONLY EDGE SERVER ASSIGNMENT": {
					ClientSession: TOSession, RequestBody: map[string]interface{}{
						"xmlId":       "test-ds-server-assignments",
						"serverNames": []string{"denver-mso-org-01"},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when REMOVING ONLY ORIGIN SERVER ASSIGNMENT": {
					ClientSession: TOSession, RequestBody: map[string]interface{}{
						"xmlId":       "test-ds-server-assignments",
						"serverNames": []string{"test-ds-server-assignments"},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
		}

		// Tests for the /deliveryserviceserver/{{DSID}}/{{serverID}} route
		dssDSIDServerIDTests := utils.V4TestCase{
			"DELETE": {
				// PREREQ have assignment made // Endpoint needs two ids
				"OK when VALID REQUEST": {
					ClientSession: TOSession, Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				// WHERE server ID is for edge server = test-ds-server-assignments xmlid = test-ds-server-assignments
				"BAD REQUEST when REMOVING ONLY EDGE SERVER ASSIGNMENT": {
					ClientSession: TOSession, Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				// DELETE /deliveryserviceserver/{{DSID}}/{{serverID}} WHERE server ID is for origin server =  test-mso-org-01
				"BAD REQUEST when REMOVING ONLY ORIGIN SERVER ASSIGNMENT": {
					ClientSession: TOSession, Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
		}

		for method, testCases := range dssTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					var dsID int
					var replace bool
					var serverIDs []int

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
					}
				}
			})
		}

		for method, testCases := range dsXMLIDServersTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {

					if testCase.RequestBody != nil {

					}

					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetServersByDeliveryService(id, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						//AssignServersToDeliveryService
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.AssignServersToDeliveryService(),,, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					}
				}
			})
		}

		for method, testCases := range dssDSIDServerIDTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {

					if testCase.RequestBody != nil {

					}

					switch method {
					case "POST":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.AssignDeliveryServiceIDsToServerID(testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					}
				}
			})
		}

		for method, testCases := range dssDSIDServerIDTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {

					if testCase.RequestBody != nil {

					}

					switch method {
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.DeleteDeliveryServiceServer(testCase.RequestOpts)
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
