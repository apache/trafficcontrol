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
	WithObjs(t, []TCObj{Types, CacheGroups, CDNs, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, Servers, Topologies, Tenants, Roles, Users, CDNLocks}, func() {

		opsUserSession := createSession(t, "opsuser", "pa$$word")
		opsUserWithLockSession := createSession(t, "opslockuser", "pa$$word")

		methodTests := map[string]map[string]struct {
			endpointId    func() int
			clientSession *client.Session
			requestOpts   client.RequestOptions
			requestBody   map[string]interface{}
			expectations  []utils.CkReqFunc
		}{
			"GET": {
				"OK when VALID request": {
					clientSession: TOSession, expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						utils.ResponseLengthGreaterOrEqual(1)),
				},
			},
			"POST": {
				"CREATED when VALID request": {
					clientSession: TOSession,
					requestBody: map[string]interface{}{
						"cdn":     "cdn3",
						"message": "snapping cdn",
						"soft":    true,
					},
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusCreated),
						validateResponseFields(map[string]interface{}{"username": "admin", "cdn": "cdn3", "message": "snapping cdn", "soft": true})),
				},
			},
			"DELETE": {
				"OK when VALID request": {
					clientSession: TOSession, requestOpts: client.RequestOptions{QueryParameters: url.Values{"cdn": {"cdn1"}}},
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"FORBIDDEN when NON-ADMIN USER DOESNT OWN LOCK": {
					clientSession: opsUserSession,
					requestOpts:   client.RequestOptions{QueryParameters: url.Values{"cdn": {"cdn4"}}},
					expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
				"OK when ADMIN USER DOESNT OWN LOCK": {
					clientSession: TOSession, requestOpts: client.RequestOptions{QueryParameters: url.Values{"cdn": {"cdn4"}}},
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
			"SNAPSHOT": {
				"OK when USER OWNS LOCK": {
					clientSession: opsUserWithLockSession,
					requestOpts:   client.RequestOptions{QueryParameters: url.Values{"cdn": {"cdn2"}}},
					expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"FORBIDDEN when ADMIN USER DOESNT OWN LOCK": {
					clientSession: TOSession, requestOpts: client.RequestOptions{QueryParameters: url.Values{"cdn": {"cdn2"}}},
					expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
			"SERVERS QUEUE UPDATES": {
				"OK when USER OWNS LOCK": {
					endpointId: getServerID(t, "cdn2-test-edge"), clientSession: opsUserWithLockSession,
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"FORBIDDEN when ADMIN USER DOESNT OWN LOCK": {
					endpointId: getServerID(t, "cdn2-test-edge"), clientSession: TOSession,
					expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
			"TOPOLOGY QUEUE UPDATES": {
				"OK when USER OWNS LOCK": {
					clientSession: opsUserWithLockSession,
					requestOpts:   client.RequestOptions{QueryParameters: url.Values{"topology": {"top-for-ds-req"}}},
					requestBody: map[string]interface{}{
						"action": "queue",
						"cdnId":  getCDNID(t, "cdn2"),
					},
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"FORBIDDEN when ADMIN USER DOESNT OWN LOCK": {
					clientSession: TOSession,
					requestOpts:   client.RequestOptions{QueryParameters: url.Values{"topology": {"top-for-ds-req"}}},
					requestBody: map[string]interface{}{
						"action": "queue",
						"cdnId":  getCDNID(t, "cdn2"),
					},
					expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
			"CACHE GROUP UPDATE": {
				"OK when USER OWNS LOCK": {
					endpointId: GetCacheGroupId(t, "cachegroup1"), clientSession: opsUserWithLockSession,
					requestBody: map[string]interface{}{
						"name":      "cachegroup1",
						"shortName": "newShortName",
						"typeName":  "EDGE_LOC",
						"typeId":    -1,
					},
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"FORBIDDEN when ADMIN USER DOESNT OWN LOCK": {
					endpointId: GetCacheGroupId(t, "cachegroup1"), clientSession: TOSession,
					requestBody: map[string]interface{}{
						"name":      "cachegroup1",
						"shortName": "newShortName",
						"typeName":  "EDGE_LOC",
						"typeId":    -1,
					},
					expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {

					topology := ""
					cdnLock := tc.CDNLock{}
					cacheGroup := tc.CacheGroupNullable{}
					topQueueUp := tc.TopologiesQueueUpdateRequest{}

					if testCase.requestOpts.QueryParameters.Has("topology") {
						topology = testCase.requestOpts.QueryParameters.Get("topology")
					}

					if testCase.requestBody != nil {
						if getId, ok := testCase.requestBody["cdnId"]; ok {
							testCase.requestBody["cdnId"] = getId.(func() int)()
							dat, err := json.Marshal(testCase.requestBody)
							assert.NoError(t, err, "Error occurred when marshalling request body: %v", err)
							err = json.Unmarshal(dat, &topQueueUp)
							assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
						} else if typeName, ok := testCase.requestBody["typeName"]; ok {
							testCase.requestBody["typeId"] = GetTypeId(t, typeName.(string))
							dat, err := json.Marshal(testCase.requestBody)
							assert.NoError(t, err, "Error occurred when marshalling request body: %v", err)
							err = json.Unmarshal(dat, &cacheGroup)
							assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
						} else {
							dat, err := json.Marshal(testCase.requestBody)
							assert.NoError(t, err, "Error occurred when marshalling request body: %v", err)
							err = json.Unmarshal(dat, &cdnLock)
							assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
						}
					}

					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.clientSession.GetCDNLocks(testCase.requestOpts)
							for _, check := range testCase.expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.clientSession.CreateCDNLock(cdnLock, testCase.requestOpts)
							for _, check := range testCase.expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.clientSession.DeleteCDNLocks(testCase.requestOpts)
							for _, check := range testCase.expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "SNAPSHOT":
						{
							t.Run(name, func(t *testing.T) {
								resp, reqInf, err := testCase.clientSession.SnapshotCRConfig(testCase.requestOpts)
								for _, check := range testCase.expectations {
									check(t, reqInf, resp.Response, resp.Alerts, err)
								}
							})
						}
					case "SERVERS QUEUE UPDATES":
						{
							t.Run(name, func(t *testing.T) {
								resp, reqInf, err := testCase.clientSession.SetServerQueueUpdate(testCase.endpointId(), true, testCase.requestOpts)
								for _, check := range testCase.expectations {
									check(t, reqInf, resp.Response, resp.Alerts, err)
								}
							})
						}
					case "TOPOLOGY QUEUE UPDATES":
						{
							t.Run(name, func(t *testing.T) {
								resp, reqInf, err := testCase.clientSession.TopologiesQueueUpdate(topology, topQueueUp, testCase.requestOpts)
								for _, check := range testCase.expectations {
									check(t, reqInf, nil, resp.Alerts, err)
								}
							})
						}
					case "CACHE GROUP UPDATE":
						{
							t.Run(name, func(t *testing.T) {
								resp, reqInf, err := testCase.clientSession.UpdateCacheGroup(testCase.endpointId(), cacheGroup, testCase.requestOpts)
								for _, check := range testCase.expectations {
									check(t, reqInf, nil, resp.Alerts, err)
								}
							})
						}
					}
				}
			})
		}
	})
}

// createSession creates a session using the passed in username and password.
func createSession(t *testing.T, username string, password string) *client.Session {
	userSession, _, err := client.LoginWithAgent(Config.TrafficOps.URL, username, password, true, "to-api-v4-client-tests", false, toReqTimeout)
	assert.RequireNoError(t, err, "Could not login with user %v: %v", username, err)
	return userSession
}

func validateResponseFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		cdnLockResp := resp.(tc.CDNLock)
		assert.Equal(t, expectedResp["username"], cdnLockResp.UserName, "Expected username: %v Got: %v", expectedResp["username"], cdnLockResp.UserName)
		assert.Equal(t, expectedResp["cdn"], cdnLockResp.CDN, "Expected CDN: %v Got: %v", expectedResp["cdn"], cdnLockResp.CDN)
		assert.Equal(t, expectedResp["message"], *cdnLockResp.Message, "Expected Message %v Got: %v", expectedResp["message"], *cdnLockResp.Message)
		assert.Equal(t, expectedResp["soft"], *cdnLockResp.Soft, "Expected 'Soft' to be: %v Got: %v", expectedResp["soft"], *cdnLockResp.Soft)
	}
}

func getCDNID(t *testing.T, cdnName string) func() int {
	return func() int {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("name", cdnName)
		cdnsResp, _, err := TOSession.GetCDNs(opts)
		assert.NoError(t, err, "Get CDNs Request failed with error:", err)
		assert.Equal(t, 1, len(cdnsResp.Response), "Expected response object length 1, but got %d", len(cdnsResp.Response))
		assert.NotNil(t, cdnsResp.Response[0].ID, "Expected id to not be nil")
		return cdnsResp.Response[0].ID
	}
}

func getServerID(t *testing.T, hostName string) func() int {
	return func() int {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("hostName", hostName)
		serversResp, _, err := TOSession.GetServers(opts)
		assert.NoError(t, err, "Get Servers Request failed with error:", err)
		assert.Equal(t, 1, len(serversResp.Response), "Expected response object length 1, but got %d", len(serversResp.Response))
		assert.NotNil(t, serversResp.Response[0].ID, "Expected id to not be nil")
		return *serversResp.Response[0].ID
	}
}

func CreateTestCDNLocks(t *testing.T) {
	for _, cl := range testData.CDNLocks {
		clientSession := TOSession
		if cl.UserName != "" {
			for _, user := range testData.Users {
				if user.Username == cl.UserName {
					clientSession = createSession(t, user.Username, *user.LocalPassword)
				}
			}
		}
		resp, _, err := clientSession.CreateCDNLock(cl, client.RequestOptions{})
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
