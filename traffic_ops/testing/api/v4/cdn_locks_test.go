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
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/assert"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

func TestCDNLocks(t *testing.T) {
	WithObjs(t, []TCObj{Types, CacheGroups, CDNs, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, Servers, ServiceCategories, Topologies, Tenants, Roles, Users, DeliveryServices, CDNLocks}, func() {

		opsUserSession := utils.CreateV4Session(t, Config.TrafficOps.URL, "opsuser", "pa$$word", Config.Default.Session.TimeoutInSecs)
		opsUserWithLockSession := utils.CreateV4Session(t, Config.TrafficOps.URL, "opslockuser", "pa$$word", Config.Default.Session.TimeoutInSecs)

		methodTests := utils.V4TestCase{
			"GET": {
				"OK when VALID request": {
					ClientSession: TOSession, Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						utils.ResponseLengthGreaterOrEqual(1)),
				},
				"OK when VALID CDN parameter": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"cdn": {"cdn2"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						validateGetResponseFields(map[string]interface{}{"username": "opslockuser", "cdn": "cdn2", "message": "test lock for updates", "soft": false})),
				},
			},
			"POST": {
				"CREATED when VALID request": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"cdn":     "cdn3",
						"message": "snapping cdn",
						"soft":    true,
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusCreated),
						validateCreateResponseFields(map[string]interface{}{"username": "admin", "cdn": "cdn3", "message": "snapping cdn", "soft": true})),
				},
				"NOT CREATED when INVALID shared username": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"cdn":             "bar",
						"message":         "snapping cdn",
						"soft":            true,
						"sharedUserNames": []string{"adminuser2"},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"CREATED when VALID shared username": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"cdn":             "bar",
						"message":         "snapping cdn",
						"soft":            true,
						"sharedUserNames": []string{"adminuser"},
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusCreated)),
				},
			},
			"DELETE": {
				"OK when VALID request": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"cdn": {"cdn1"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"FORBIDDEN when NON-ADMIN USER DOESNT OWN LOCK": {
					ClientSession: opsUserSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"cdn": {"cdn4"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
				"OK when ADMIN USER DOESNT OWN LOCK": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"cdn": {"cdn4"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
			"SNAPSHOT": {
				"OK when USER OWNS LOCK": {
					ClientSession: opsUserWithLockSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"cdn": {"cdn2"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"FORBIDDEN when ADMIN USER DOESNT OWN LOCK": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"cdn": {"cdn2"}}},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
			"SERVERS QUEUE UPDATES": {
				"OK when USER OWNS LOCK": {
					EndpointId: GetServerID(t, "cdn2-test-edge"), ClientSession: opsUserWithLockSession,
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"FORBIDDEN when ADMIN USER DOESNT OWN LOCK": {
					EndpointId: GetServerID(t, "cdn2-test-edge"), ClientSession: TOSession,
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
			"TOPOLOGY QUEUE UPDATES": {
				"OK when USER OWNS LOCK": {
					ClientSession: opsUserWithLockSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"topology": {"top-for-ds-req"}}},
					RequestBody: map[string]interface{}{
						"action": "queue",
						"cdnId":  GetCDNID(t, "cdn2")(),
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"FORBIDDEN when ADMIN USER DOESNT OWN LOCK": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"topology": {"top-for-ds-req"}}},
					RequestBody: map[string]interface{}{
						"action": "queue",
						"cdnId":  GetCDNID(t, "cdn2")(),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
				"OK when ADMIN USER DOESNT OWN LOCK FOR DEQUEUE": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"topology": {"top-for-ds-req"}}},
					RequestBody: map[string]interface{}{
						"action": "dequeue",
						"cdnId":  GetCDNID(t, "cdn2")(),
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
			"CDN UPDATE": {
				"OK when USER OWNS LOCK": {
					EndpointId: GetCDNID(t, "cdn2"), ClientSession: opsUserWithLockSession,
					RequestBody: map[string]interface{}{
						"dnssecEnabled": false,
						"domainName":    "newdomain",
						"name":          "cdn2",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"FORBIDDEN when ADMIN USER DOESNT OWN LOCK": {
					EndpointId: GetCDNID(t, "cdn2"), ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"dnssecEnabled": false,
						"domainName":    "newdomaintest",
						"name":          "cdn2",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
			"CDN DELETE": {
				"OK when USER OWNS LOCK": {
					EndpointId: GetCDNID(t, "cdndelete"), ClientSession: opsUserWithLockSession,
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"FORBIDDEN when ADMIN USER DOESNT OWN LOCK": {
					EndpointId: GetCDNID(t, "cdn2"), ClientSession: TOSession,
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
			"CACHE GROUP UPDATE": {
				"OK when USER OWNS LOCK": {
					EndpointId: GetCacheGroupId(t, "cachegroup1"), ClientSession: opsUserWithLockSession,
					RequestBody: map[string]interface{}{
						"name":      "cachegroup1",
						"shortName": "newShortName",
						"typeName":  "EDGE_LOC",
						"typeId":    -1,
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"FORBIDDEN when ADMIN USER DOESNT OWN LOCK": {
					EndpointId: GetCacheGroupId(t, "cachegroup1"), ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"name":      "cachegroup1",
						"shortName": "newShortName",
						"typeName":  "EDGE_LOC",
						"typeId":    -1,
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
			"DELIVERY SERVICE POST": {
				"OK when USER OWNS LOCK": {
					ClientSession: opsUserWithLockSession, RequestBody: generateDeliveryService(t, map[string]interface{}{"xmlId": "testDSLock"}),
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusCreated)),
				},
				"FORBIDDEN when ADMIN USER DOESNT OWN LOCK": {
					ClientSession: TOSession, RequestBody: generateDeliveryService(t, map[string]interface{}{
						"xmlId": "testDSLock2", "cdnId": GetCDNID(t, "cdn2")()}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
			"DELIVERY SERVICE PUT": {
				"OK when USER OWNS LOCK": {
					EndpointId: GetDeliveryServiceId(t, "basic-ds-in-cdn2"), ClientSession: opsUserWithLockSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"xmlId": "basic-ds-in-cdn2", "cdnId": GetCDNID(t, "cdn2")(), "cdnName": "cdn2", "routingName": "cdn"}),
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"FORBIDDEN when ADMIN USER DOESNT OWN LOCK": {
					EndpointId: GetDeliveryServiceId(t, "basic-ds-in-cdn2"), ClientSession: TOSession,
					RequestBody:  generateDeliveryService(t, map[string]interface{}{}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
			"DELIVERY SERVICE DELETE": {
				"OK when USER OWNS LOCK": {
					EndpointId: GetDeliveryServiceId(t, "ds-forked-topology"), ClientSession: opsUserWithLockSession,
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"FORBIDDEN when ADMIN USER DOESNT OWN LOCK": {
					EndpointId: GetDeliveryServiceId(t, "top-ds-in-cdn2"), ClientSession: TOSession,
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
			"SERVER POST": {
				"OK when USER OWNS LOCK": {
					ClientSession: opsUserWithLockSession,
					RequestBody: generateServer(t, map[string]interface{}{
						"cdnId":        GetCDNID(t, "cdn2")(),
						"profileNames": []string{"EDGEInCDN2"},
					}),
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusCreated)),
				},
				"FORBIDDEN when ADMIN USER DOESNT OWN LOCK": {
					ClientSession: TOSession,
					RequestBody: generateServer(t, map[string]interface{}{
						"cdnId":        GetCDNID(t, "cdn2")(),
						"profileNames": []string{"EDGEInCDN2"},
						"interfaces": []map[string]interface{}{{
							"ipAddresses": []map[string]interface{}{{
								"address":        "127.0.0.2/30",
								"serviceAddress": true,
							}},
							"name": "eth0",
						}},
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
			"SERVER PUT": {
				"OK when USER OWNS LOCK": {
					EndpointId:    GetServerID(t, "edge1-cdn2"),
					ClientSession: opsUserWithLockSession,
					RequestBody: generateServer(t, generateServer(t, map[string]interface{}{
						"id":           GetServerID(t, "edge1-cdn2")(),
						"cdnId":        GetCDNID(t, "cdn2")(),
						"profileNames": []string{"EDGEInCDN2"},
						"interfaces": []map[string]interface{}{{
							"ipAddresses": []map[string]interface{}{{
								"address":        "0.0.0.1",
								"serviceAddress": true,
							}},
							"name": "eth0",
						}},
					})),
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"FORBIDDEN when ADMIN USER DOESNT OWN LOCK": {
					EndpointId:    GetServerID(t, "dtrc-edge-07"),
					ClientSession: TOSession,
					RequestBody: generateServer(t, generateServer(t, map[string]interface{}{
						"id":           GetServerID(t, "dtrc-edge-07")(),
						"cdnId":        GetCDNID(t, "cdn2")(),
						"cachegroupId": GetCacheGroupId(t, "dtrc2")(),
						"profileNames": []string{"CDN2_EDGE"},
						"interfaces": []map[string]interface{}{{
							"ipAddresses": []map[string]interface{}{{
								"address":        "192.0.2.11/24",
								"serviceAddress": true,
							}},
							"name": "eth0",
						}},
					})),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
			"SERVER DELETE": {
				"OK when USER OWNS LOCK": {
					EndpointId:    GetServerID(t, "atlanta-mid-17"),
					ClientSession: opsUserWithLockSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"FORBIDDEN when ADMIN USER DOESNT OWN LOCK": {
					EndpointId:    GetServerID(t, "denver-mso-org-02"),
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {

					topology := ""
					cdn := tc.CDN{}
					cdnLock := tc.CDNLock{}
					cacheGroup := tc.CacheGroupNullable{}
					ds := tc.DeliveryServiceV4{}
					server := tc.ServerV4{}
					topQueueUp := tc.TopologiesQueueUpdateRequest{}

					if testCase.RequestOpts.QueryParameters.Has("topology") {
						topology = testCase.RequestOpts.QueryParameters.Get("topology")
					}

					if testCase.RequestBody != nil {
						if _, ok := testCase.RequestBody["xmlId"]; ok {
							dat, err := json.Marshal(testCase.RequestBody)
							assert.NoError(t, err, "Error occurred when marshalling request body: %v", err)
							err = json.Unmarshal(dat, &ds)
							assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
						} else if _, ok := testCase.RequestBody["hostName"]; ok {
							dat, err := json.Marshal(testCase.RequestBody)
							assert.NoError(t, err, "Error occurred when marshalling request body: %v", err)
							err = json.Unmarshal(dat, &server)
							assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
						} else if getId, ok := testCase.RequestBody["cdnId"]; ok {
							testCase.RequestBody["cdnId"] = getId.(int)
							dat, err := json.Marshal(testCase.RequestBody)
							assert.NoError(t, err, "Error occurred when marshalling request body: %v", err)
							err = json.Unmarshal(dat, &topQueueUp)
							assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
						} else if typeName, ok := testCase.RequestBody["typeName"]; ok {
							testCase.RequestBody["typeId"] = GetTypeId(t, typeName.(string))
							dat, err := json.Marshal(testCase.RequestBody)
							assert.NoError(t, err, "Error occurred when marshalling request body: %v", err)
							err = json.Unmarshal(dat, &cacheGroup)
							assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
						} else if _, ok := testCase.RequestBody["dnssecEnabled"]; ok {
							dat, err := json.Marshal(testCase.RequestBody)
							assert.NoError(t, err, "Error occurred when marshalling request body: %v", err)
							err = json.Unmarshal(dat, &cdn)
							assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
						} else {
							dat, err := json.Marshal(testCase.RequestBody)
							assert.NoError(t, err, "Error occurred when marshalling request body: %v", err)
							err = json.Unmarshal(dat, &cdnLock)
							assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
						}
					}

					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetCDNLocks(testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.CreateCDNLock(cdnLock, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.DeleteCDNLocks(testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "SNAPSHOT":
						{
							t.Run(name, func(t *testing.T) {
								resp, reqInf, err := testCase.ClientSession.SnapshotCRConfig(testCase.RequestOpts)
								for _, check := range testCase.Expectations {
									check(t, reqInf, resp.Response, resp.Alerts, err)
								}
							})
						}
					case "SERVERS QUEUE UPDATES":
						{
							t.Run(name, func(t *testing.T) {
								resp, reqInf, err := testCase.ClientSession.SetServerQueueUpdate(testCase.EndpointId(), true, testCase.RequestOpts)
								for _, check := range testCase.Expectations {
									check(t, reqInf, resp.Response, resp.Alerts, err)
								}
							})
						}
					case "TOPOLOGY QUEUE UPDATES":
						{
							t.Run(name, func(t *testing.T) {
								resp, reqInf, err := testCase.ClientSession.TopologiesQueueUpdate(topology, topQueueUp, testCase.RequestOpts)
								for _, check := range testCase.Expectations {
									check(t, reqInf, nil, resp.Alerts, err)
								}
							})
						}
					case "CACHE GROUP UPDATE":
						{
							t.Run(name, func(t *testing.T) {
								resp, reqInf, err := testCase.ClientSession.UpdateCacheGroup(testCase.EndpointId(), cacheGroup, testCase.RequestOpts)
								for _, check := range testCase.Expectations {
									check(t, reqInf, nil, resp.Alerts, err)
								}
							})
						}
					case "CDN UPDATE":
						{
							t.Run(name, func(t *testing.T) {
								alerts, reqInf, err := testCase.ClientSession.UpdateCDN(testCase.EndpointId(), cdn, testCase.RequestOpts)
								for _, check := range testCase.Expectations {
									check(t, reqInf, nil, alerts, err)
								}
							})
						}
					case "CDN DELETE":
						{
							t.Run(name, func(t *testing.T) {
								alerts, reqInf, err := testCase.ClientSession.DeleteCDN(testCase.EndpointId(), testCase.RequestOpts)
								for _, check := range testCase.Expectations {
									check(t, reqInf, nil, alerts, err)
								}
							})
						}
					case "DELIVERY SERVICE POST":
						{
							t.Run(name, func(t *testing.T) {
								resp, reqInf, err := testCase.ClientSession.CreateDeliveryService(ds, testCase.RequestOpts)
								for _, check := range testCase.Expectations {
									check(t, reqInf, nil, resp.Alerts, err)
								}
							})
						}
					case "DELIVERY SERVICE PUT":
						{
							t.Run(name, func(t *testing.T) {
								resp, reqInf, err := testCase.ClientSession.UpdateDeliveryService(testCase.EndpointId(), ds, testCase.RequestOpts)
								for _, check := range testCase.Expectations {
									check(t, reqInf, nil, resp.Alerts, err)
								}
							})
						}
					case "DELIVERY SERVICE DELETE":
						{
							t.Run(name, func(t *testing.T) {
								resp, reqInf, err := testCase.ClientSession.DeleteDeliveryService(testCase.EndpointId(), testCase.RequestOpts)
								for _, check := range testCase.Expectations {
									check(t, reqInf, nil, resp.Alerts, err)
								}
							})
						}
					case "SERVER POST":
						{
							t.Run(name, func(t *testing.T) {
								alerts, reqInf, err := testCase.ClientSession.CreateServer(server, testCase.RequestOpts)
								for _, check := range testCase.Expectations {
									check(t, reqInf, nil, alerts, err)
								}
							})
						}
					case "SERVER PUT":
						{
							t.Run(name, func(t *testing.T) {
								alerts, reqInf, err := testCase.ClientSession.UpdateServer(testCase.EndpointId(), server, testCase.RequestOpts)
								for _, check := range testCase.Expectations {
									check(t, reqInf, nil, alerts, err)
								}
							})
						}
					case "SERVER DELETE":
						{
							t.Run(name, func(t *testing.T) {
								alerts, reqInf, err := testCase.ClientSession.DeleteServer(testCase.EndpointId(), testCase.RequestOpts)
								for _, check := range testCase.Expectations {
									check(t, reqInf, nil, alerts, err)
								}
							})
						}
					}
				}
			})
		}
	})
}

func validateGetResponseFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		cdnLockResp := resp.([]tc.CDNLock)
		assert.Equal(t, expectedResp["username"], cdnLockResp[0].UserName, "Expected username: %v Got: %v", expectedResp["username"], cdnLockResp[0].UserName)
		assert.Equal(t, expectedResp["cdn"], cdnLockResp[0].CDN, "Expected CDN: %v Got: %v", expectedResp["cdn"], cdnLockResp[0].CDN)
		assert.Equal(t, expectedResp["message"], *cdnLockResp[0].Message, "Expected Message %v Got: %v", expectedResp["message"], *cdnLockResp[0].Message)
		assert.Equal(t, expectedResp["soft"], *cdnLockResp[0].Soft, "Expected 'Soft' to be: %v Got: %v", expectedResp["soft"], *cdnLockResp[0].Soft)
	}
}

func validateCreateResponseFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		cdnLockResp := resp.(tc.CDNLock)
		assert.Equal(t, expectedResp["username"], cdnLockResp.UserName, "Expected username: %v Got: %v", expectedResp["username"], cdnLockResp.UserName)
		assert.Equal(t, expectedResp["cdn"], cdnLockResp.CDN, "Expected CDN: %v Got: %v", expectedResp["cdn"], cdnLockResp.CDN)
		assert.Equal(t, expectedResp["message"], *cdnLockResp.Message, "Expected Message %v Got: %v", expectedResp["message"], *cdnLockResp.Message)
		assert.Equal(t, expectedResp["soft"], *cdnLockResp.Soft, "Expected 'Soft' to be: %v Got: %v", expectedResp["soft"], *cdnLockResp.Soft)
	}
}

func CreateTestCDNLocks(t *testing.T) {
	for _, cl := range testData.CDNLocks {
		ClientSession := TOSession
		if cl.UserName != "" {
			for _, user := range testData.Users {
				if user.Username == cl.UserName {
					ClientSession = utils.CreateV4Session(t, Config.TrafficOps.URL, user.Username, *user.LocalPassword, Config.Default.Session.TimeoutInSecs)
				}
			}
		}
		resp, _, err := ClientSession.CreateCDNLock(cl, client.RequestOptions{})
		assert.NoError(t, err, "Could not create CDN Lock: %v - alerts: %+v", err, resp.Alerts)
	}
}

func DeleteTestCDNLocks(t *testing.T) {
	opts := client.NewRequestOptions()
	cdnlocks, _, err := TOSession.GetCDNLocks(opts)
	assert.NoError(t, err, "Error retrieving CDN Locks for deletion: %v - alerts: %+v", err, cdnlocks.Alerts)
	assert.GreaterOrEqual(t, len(cdnlocks.Response), 1, "Expected at least one CDN Lock for deletion")
	for _, cl := range cdnlocks.Response {
		opts.QueryParameters.Set("cdn", cl.CDN)
		resp, _, err := TOSession.DeleteCDNLocks(opts)
		assert.NoError(t, err, "Could not delete CDN Lock: %v - alerts: %+v", err, resp.Alerts)
		// Retrieve the CDN Lock to see if it got deleted
		cdnlock, _, err := TOSession.GetCDNLocks(opts)
		assert.NoError(t, err, "Error deleting CDN Lock for '%s' : %v - alerts: %+v", cl.CDN, err, cdnlock.Alerts)
		assert.Equal(t, 0, len(cdnlock.Response), "Expected CDN Lock for '%s' to be deleted", cl.CDN)
	}
}
