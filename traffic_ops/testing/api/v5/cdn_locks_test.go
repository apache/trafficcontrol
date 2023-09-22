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

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

func TestCDNLocks(t *testing.T) {
	WithObjs(t, []TCObj{Types, CacheGroups, CDNs, Parameters, Profiles, ProfileParameters, Statuses, Divisions, Regions, PhysLocations, Servers, ServiceCategories, Topologies, Tenants, Roles, Users, ServerCapabilities, DeliveryServices, StaticDNSEntries, CDNLocks}, func() {

		now := time.Now().Round(time.Microsecond)
		opsUserSession := utils.CreateV5Session(t, Config.TrafficOps.URL, "opsuser", "pa$$word", Config.Default.Session.TimeoutInSecs)
		opsUserWithLockSession := utils.CreateV5Session(t, Config.TrafficOps.URL, "opslockuser", "pa$$word", Config.Default.Session.TimeoutInSecs)

		methodTests := utils.V5TestCase{
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
					EndpointID: GetServerID(t, "cdn2-test-edge"), ClientSession: opsUserWithLockSession,
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"FORBIDDEN when ADMIN USER DOESNT OWN LOCK": {
					EndpointID: GetServerID(t, "cdn2-test-edge"), ClientSession: TOSession,
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
			"SERVERS HOSTNAME UPDATE": {
				"CONFIG_APPLY_TIME is SET EVEN when CDN LOCKED": {
					ClientSession: opsUserWithLockSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"hostName": {"cdn2-test-edge"}}},
					RequestBody: map[string]interface{}{
						"config_apply_time":    util.Ptr(now),
						"config_update_failed": util.Ptr(true),
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateServerApplyTimes("cdn2-test-edge", map[string]interface{}{"ConfigApplyTime": now, "ConfigUpdateFailed": true})),
				},
				"REVALIDATE_APPLY_TIME is SET EVEN when CDN LOCKED": {
					ClientSession: opsUserWithLockSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"hostName": {"cdn2-test-edge"}}},
					RequestBody: map[string]interface{}{
						"revalidate_apply_time":    util.Ptr(now),
						"revalidate_update_failed": util.Ptr(true),
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateServerApplyTimes("cdn2-test-edge", map[string]interface{}{"RevalApplyTime": now, "RevalUpdateFailed": true})),
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
					EndpointID: GetCDNID(t, "cdn2"), ClientSession: opsUserWithLockSession,
					RequestBody: map[string]interface{}{
						"dnssecEnabled": false,
						"domainName":    "newdomain",
						"name":          "cdn2",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"FORBIDDEN when ADMIN USER DOESNT OWN LOCK": {
					EndpointID: GetCDNID(t, "cdn2"), ClientSession: TOSession,
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
					EndpointID: GetCDNID(t, "cdndelete"), ClientSession: opsUserWithLockSession,
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"FORBIDDEN when ADMIN USER DOESNT OWN LOCK": {
					EndpointID: GetCDNID(t, "cdn2"), ClientSession: TOSession,
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
			"CACHE GROUP UPDATE": {
				"OK when USER OWNS LOCK": {
					EndpointID: GetCacheGroupId(t, "cachegroup1"), ClientSession: opsUserWithLockSession,
					RequestBody: map[string]interface{}{
						"name":      "cachegroup1",
						"shortName": "newShortName",
						"typeName":  "EDGE_LOC",
						"typeId":    GetTypeId(t, "EDGE_LOC"),
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"FORBIDDEN when ADMIN USER DOESNT OWN LOCK": {
					EndpointID: GetCacheGroupId(t, "cachegroup1"), ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"name":      "cachegroup1",
						"shortName": "newShortName",
						"typeName":  "EDGE_LOC",
						"typeId":    GetTypeId(t, "EDGE_LOC"),
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
					EndpointID: GetDeliveryServiceId(t, "basic-ds-in-cdn2"), ClientSession: opsUserWithLockSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"xmlId": "basic-ds-in-cdn2", "cdnId": GetCDNID(t, "cdn2")(), "cdnName": "cdn2", "routingName": "cdn"}),
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"FORBIDDEN when ADMIN USER DOESNT OWN LOCK": {
					EndpointID: GetDeliveryServiceId(t, "basic-ds-in-cdn2"), ClientSession: TOSession,
					RequestBody:  generateDeliveryService(t, map[string]interface{}{}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
			"DELIVERY SERVICE DELETE": {
				"OK when USER OWNS LOCK": {
					EndpointID: GetDeliveryServiceId(t, "ds-forked-topology"), ClientSession: opsUserWithLockSession,
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"FORBIDDEN when ADMIN USER DOESNT OWN LOCK": {
					EndpointID: GetDeliveryServiceId(t, "top-ds-in-cdn2"), ClientSession: TOSession,
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
			"PROFILE POST": {
				"OK when USER OWNS LOCK": {
					ClientSession: opsUserWithLockSession,
					RequestBody: map[string]interface{}{
						"cdn":              GetCDNID(t, "cdn2")(),
						"cdnName":          "cdn2",
						"description":      "test cdn locks description",
						"name":             "TestLocks",
						"routing_disabled": false,
						"type":             "ATS_PROFILE",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusCreated)),
				},
				"FORBIDDEN when ADMIN USER DOESNT OWN LOCK": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"cdn":              GetCDNID(t, "cdn2")(),
						"cdnName":          "cdn2",
						"description":      "test cdn locks description",
						"name":             "TestLocksForbidden",
						"routing_disabled": false,
						"type":             "ATS_PROFILE",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
			"PROFILE PUT": {
				"OK when USER OWNS LOCK": {
					EndpointID:    GetProfileID(t, "CDN2_EDGE"),
					ClientSession: opsUserWithLockSession,
					RequestBody: map[string]interface{}{
						"cdn":              GetCDNID(t, "cdn2")(),
						"cdnName":          "cdn2",
						"description":      "cdn2 edge description updated when user owns lock",
						"name":             "CDN2_EDGE",
						"routing_disabled": false,
						"type":             "ATS_PROFILE",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"FORBIDDEN when ADMIN USER DOESNT OWN LOCK": {
					EndpointID:    GetProfileID(t, "EDGEInCDN2"),
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"cdn":              GetCDNID(t, "cdn2")(),
						"cdnName":          "cdn2",
						"description":      "should fail",
						"name":             "EDGEInCDN2",
						"routing_disabled": false,
						"type":             "ATS_PROFILE",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
			"PROFILE DELETE": {
				"OK when USER OWNS LOCK": {
					EndpointID:    GetProfileID(t, "CCR2"),
					ClientSession: opsUserWithLockSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"FORBIDDEN when ADMIN USER DOESNT OWN LOCK": {
					EndpointID:    GetProfileID(t, "MID2"),
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
			"PROFILE PARAMETER POST": {
				"OK when USER OWNS LOCK": {
					ClientSession: opsUserWithLockSession,
					RequestBody: map[string]interface{}{
						"profileId":   GetProfileID(t, "EDGEInCDN2")(),
						"parameterId": GetParameterID(t, "CONFIG proxy.config.admin.user_id", "records.config", "STRING ats")(),
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusCreated)),
				},
				"FORBIDDEN when ADMIN USER DOESNT OWN LOCK": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"profileId":   GetProfileID(t, "EDGEInCDN2")(),
						"parameterId": GetParameterID(t, "CONFIG proxy.config.admin.user_id", "records.config", "STRING ats")(),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
			"PROFILE PARAMETER DELETE": {
				"OK when USER OWNS LOCK": {
					EndpointID:    GetProfileID(t, "OKwhenUserOwnLocks"),
					ClientSession: opsUserWithLockSession,
					RequestOpts: client.RequestOptions{QueryParameters: url.Values{
						"parameterId": {strconv.Itoa(GetParameterID(t, "test.cdnlock.delete", "rascal.properties", "25.0")())},
					}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"FORBIDDEN when ADMIN USER DOESNT OWN LOCK": {
					EndpointID:    GetProfileID(t, "FORBIDDENwhenDoesntOwnLock"),
					ClientSession: TOSession,
					RequestOpts: client.RequestOptions{QueryParameters: url.Values{
						"parameterId": {strconv.Itoa(GetParameterID(t, "test.cdnlock.forbidden.delete", "rascal.properties", "25.0")())},
					}},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
			"SERVER POST": {
				"OK when USER OWNS LOCK": {
					ClientSession: opsUserWithLockSession,
					RequestBody: generateServer(t, map[string]interface{}{
						"cdnID":    GetCDNID(t, "cdn2")(),
						"profiles": []string{"EDGEInCDN2"},
					}),
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusCreated)),
				},
				"FORBIDDEN when ADMIN USER DOESNT OWN LOCK": {
					ClientSession: TOSession,
					RequestBody: generateServer(t, map[string]interface{}{
						"cdnID":    GetCDNID(t, "cdn2")(),
						"profiles": []string{"EDGEInCDN2"},
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
					EndpointID:    GetServerID(t, "edge1-cdn2"),
					ClientSession: opsUserWithLockSession,
					RequestBody: generateServer(t, map[string]interface{}{
						"id":       GetServerID(t, "edge1-cdn2")(),
						"cdnID":    GetCDNID(t, "cdn2")(),
						"profiles": []string{"EDGEInCDN2"},
						"interfaces": []map[string]interface{}{{
							"ipAddresses": []map[string]interface{}{{
								"address":        "0.0.0.1",
								"serviceAddress": true,
							}},
							"name": "eth0",
						}},
					}),
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"FORBIDDEN when ADMIN USER DOESNT OWN LOCK": {
					EndpointID:    GetServerID(t, "dtrc-edge-07"),
					ClientSession: TOSession,
					RequestBody: generateServer(t, map[string]interface{}{
						"id":           GetServerID(t, "dtrc-edge-07")(),
						"cdnID":        GetCDNID(t, "cdn2")(),
						"cachegroupId": GetCacheGroupId(t, "dtrc2")(),
						"profiles":     []string{"CDN2_EDGE"},
						"interfaces": []map[string]interface{}{{
							"ipAddresses": []map[string]interface{}{{
								"address":        "192.0.2.11/24",
								"serviceAddress": true,
							}},
							"name": "eth0",
						}},
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
			"SERVER DELETE": {
				"OK when USER OWNS LOCK": {
					EndpointID:    GetServerID(t, "atlanta-mid-17"),
					ClientSession: opsUserWithLockSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"FORBIDDEN when ADMIN USER DOESNT OWN LOCK": {
					EndpointID:    GetServerID(t, "denver-mso-org-02"),
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
			"STATIC DNS ENTRIES POST": {
				"OK when USER OWNS LOCK": {
					ClientSession: opsUserWithLockSession,
					RequestBody: map[string]interface{}{
						"address":         "192.168.0.1",
						"cachegroup":      "cachegroup1",
						"deliveryservice": "basic-ds-in-cdn2",
						"host":            "cdn_locks_test_host",
						"type":            "A_RECORD",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"FORBIDDEN when ADMIN USER DOESNT OWN LOCK": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"address":         "192.168.0.1",
						"cachegroup":      "cachegroup1",
						"deliveryservice": "basic-ds-in-cdn2",
						"host":            "cdn_locks_test_host",
						"type":            "A_RECORD",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
			"STATIC DNS ENTRIES PUT": {
				"OK when USER OWNS LOCK": {
					EndpointID:    GetStaticDNSEntryID(t, "host2"),
					ClientSession: opsUserWithLockSession,
					RequestBody: map[string]interface{}{
						"address":         "192.168.0.2",
						"cachegroup":      "cachegroup2",
						"deliveryservice": "basic-ds-in-cdn2",
						"host":            "host2",
						"type":            "A_RECORD",
						"ttl":             int64(0),
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"FORBIDDEN when ADMIN USER DOESNT OWN LOCK": {
					EndpointID:    GetStaticDNSEntryID(t, "cdnlock-test-delete-host"),
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"address":         "192.168.0.2",
						"cachegroup":      "cachegroup2",
						"deliveryservice": "basic-ds-in-cdn2",
						"host":            "host2",
						"type":            "A_RECORD",
						"ttl":             int64(0),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
			"STATIC DNS ENTRIES DELETE": {
				"OK when USER OWNS LOCK": {
					EndpointID:    GetStaticDNSEntryID(t, "host3"),
					ClientSession: opsUserWithLockSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"FORBIDDEN when ADMIN USER DOESNT OWN LOCK": {
					EndpointID:    GetStaticDNSEntryID(t, "cdnlock-negtest-delete-host"),
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					var dat []byte
					var err error

					if testCase.RequestBody != nil {
						dat, err = json.Marshal(testCase.RequestBody)
						assert.NoError(t, err, "Error occurred when marshalling request body: %v", err)
					}

					cases := map[string]func(*testing.T){
						"GET": func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetCDNLocks(testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						},
						"POST": func(t *testing.T) {
							cdnLock := tc.CDNLock{}
							err = json.Unmarshal(dat, &cdnLock)
							assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
							resp, reqInf, err := testCase.ClientSession.CreateCDNLock(cdnLock, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						},
						"DELETE": func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.DeleteCDNLocks(testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						},
						"SNAPSHOT": func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.SnapshotCRConfig(testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						},
						"SERVERS QUEUE UPDATES": func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.SetServerQueueUpdate(testCase.EndpointID(), true, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						},
						"SERVERS HOSTNAME UPDATE": func(t *testing.T) {
							var hostName string
							var configApplyTime *time.Time
							var revalApplyTime *time.Time
							var revalUpdateFailed *bool
							var configUpdateFailed *bool

							if hostNameParam, ok := testCase.RequestOpts.QueryParameters["hostName"]; ok {
								hostName = hostNameParam[0]
							}
							if configApplyTimeVal, ok := testCase.RequestBody["config_apply_time"]; ok {
								configApplyTime = configApplyTimeVal.(*time.Time)
							}
							if revalApplyTimeVal, ok := testCase.RequestBody["revalidate_apply_time"]; ok {
								revalApplyTime = revalApplyTimeVal.(*time.Time)
							}
							if val, ok := testCase.RequestBody["config_update_failed"]; ok {
								configUpdateFailed = val.(*bool)
							}
							if val, ok := testCase.RequestBody["revalidate_update_failed"]; ok {
								revalUpdateFailed = val.(*bool)
							}
							alerts, reqInf, err := testCase.ClientSession.SetUpdateServerStatusTimes(hostName, configApplyTime, revalApplyTime, configUpdateFailed, revalUpdateFailed, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						},
						"TOPOLOGY QUEUE UPDATES": func(t *testing.T) {
							topology := testCase.RequestOpts.QueryParameters.Get("topology")
							topQueueUp := tc.TopologiesQueueUpdateRequest{}
							err = json.Unmarshal(dat, &topQueueUp)
							assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
							resp, reqInf, err := testCase.ClientSession.TopologiesQueueUpdate(topology, topQueueUp, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, resp.Alerts, err)
							}
						},
						"CACHE GROUP UPDATE": func(t *testing.T) {
							cacheGroup := tc.CacheGroupNullableV5{}
							err = json.Unmarshal(dat, &cacheGroup)
							assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
							resp, reqInf, err := testCase.ClientSession.UpdateCacheGroup(testCase.EndpointID(), cacheGroup, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, resp.Alerts, err)
							}
						},
						"CDN UPDATE": func(t *testing.T) {
							cdn := tc.CDNV5{}
							err = json.Unmarshal(dat, &cdn)
							assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
							alerts, reqInf, err := testCase.ClientSession.UpdateCDN(testCase.EndpointID(), cdn, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						},
						"CDN DELETE": func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.DeleteCDN(testCase.EndpointID(), testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						},
						"DELIVERY SERVICE POST": func(t *testing.T) {
							var ds tc.DeliveryServiceV5
							err = json.Unmarshal(dat, &ds)
							assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
							resp, reqInf, err := testCase.ClientSession.CreateDeliveryService(ds, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, resp.Alerts, err)
							}
						},
						"DELIVERY SERVICE PUT": func(t *testing.T) {
							var ds tc.DeliveryServiceV5
							err = json.Unmarshal(dat, &ds)
							assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
							resp, reqInf, err := testCase.ClientSession.UpdateDeliveryService(testCase.EndpointID(), ds, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, resp.Alerts, err)
							}
						},
						"DELIVERY SERVICE DELETE": func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.DeleteDeliveryService(testCase.EndpointID(), testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, resp.Alerts, err)
							}
						},
						"PROFILE POST": func(t *testing.T) {
							profile := tc.ProfileV5{}
							err = json.Unmarshal(dat, &profile)
							assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
							alerts, reqInf, err := testCase.ClientSession.CreateProfile(profile, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						},
						"PROFILE PUT": func(t *testing.T) {
							profile := tc.ProfileV5{}
							err = json.Unmarshal(dat, &profile)
							assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
							alerts, reqInf, err := testCase.ClientSession.UpdateProfile(testCase.EndpointID(), profile, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						},
						"PROFILE DELETE": func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.DeleteProfile(testCase.EndpointID(), testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						},
						"PROFILE PARAMETER POST": func(t *testing.T) {
							profileParameter := tc.ProfileParameterCreationRequest{}
							err = json.Unmarshal(dat, &profileParameter)
							assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
							alerts, reqInf, err := testCase.ClientSession.CreateProfileParameter(profileParameter, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						},
						"PROFILE PARAMETER DELETE": func(t *testing.T) {
							parameterId, _ := strconv.Atoi(testCase.RequestOpts.QueryParameters["parameterId"][0])
							alerts, reqInf, err := testCase.ClientSession.DeleteProfileParameter(testCase.EndpointID(), parameterId, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						},
						"SERVER POST": func(t *testing.T) {
							var server tc.ServerV5
							err = json.Unmarshal(dat, &server)
							assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
							alerts, reqInf, err := testCase.ClientSession.CreateServer(server, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						},
						"SERVER PUT": func(t *testing.T) {
							var server tc.ServerV5
							err = json.Unmarshal(dat, &server)
							assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
							alerts, reqInf, err := testCase.ClientSession.UpdateServer(testCase.EndpointID(), server, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						},
						"SERVER DELETE": func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.DeleteServer(testCase.EndpointID(), testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						},
						"STATIC DNS ENTRIES POST": func(t *testing.T) {
							staticDNSEntry := tc.StaticDNSEntryV5{}
							err = json.Unmarshal(dat, &staticDNSEntry)
							assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
							staticDNSEntry.TTL = util.Ptr(int64(0))
							alerts, reqInf, err := testCase.ClientSession.CreateStaticDNSEntry(staticDNSEntry, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						},
						"STATIC DNS ENTRIES PUT": func(t *testing.T) {
							staticDNSEntry := tc.StaticDNSEntryV5{}
							err = json.Unmarshal(dat, &staticDNSEntry)
							assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
							staticDNSEntry.TTL = util.Ptr(int64(0))
							alerts, reqInf, err := testCase.ClientSession.UpdateStaticDNSEntry(testCase.EndpointID(), staticDNSEntry, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						},
						"STATIC DNS ENTRIES DELETE": func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.DeleteStaticDNSEntry(testCase.EndpointID(), testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						},
					}

					if _, ok := cases[method]; ok {
						t.Run(name, cases[method])
					} else {
						t.Errorf("Test Case: %s not found. Test: %s failed to run.", method, name)
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
					ClientSession = utils.CreateV5Session(t, Config.TrafficOps.URL, user.Username, *user.LocalPassword, Config.Default.Session.TimeoutInSecs)
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
