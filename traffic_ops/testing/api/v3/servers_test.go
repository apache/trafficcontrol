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
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

func TestServers(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, DeliveryServices, DeliveryServiceServerAssignments}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.V3TestCase{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfModifiedSince: {tomorrow}},
					Expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID HOSTNAME parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"hostName": {"atlanta-edge-01"}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						validateServerFields(map[string]interface{}{"HostName": "atlanta-edge-01"})),
				},
				"OK when VALID CACHEGROUP parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"cachegroup": {strconv.Itoa(GetCacheGroupId(t, "cachegroup1")())}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateServerFields(map[string]interface{}{"CachegroupID": GetCacheGroupId(t, "cachegroup1")()})),
				},
				"OK when VALID CACHEGROUPNAME parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"cachegroupName": {"topology-mid-cg-01"}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateServerFields(map[string]interface{}{"Cachegroup": "topology-mid-cg-01"})),
				},
				"OK when VALID CDN parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"cdn": {strconv.Itoa(GetCDNID(t, "cdn2")())}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateServerFields(map[string]interface{}{"CDNID": GetCDNID(t, "cdn2")()})),
				},
				"OK when VALID DSID parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"dsId": {strconv.Itoa(GetDeliveryServiceId(t, "test-ds-server-assignments")())}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateExpectedServers([]string{"test-ds-server-assignments"})),
				},
				"OK when VALID PARENTCACHEGROUP parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"parentCacheGroup": {strconv.Itoa(GetCacheGroupId(t, "parentCachegroup")())}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1)),
				},
				"OK when VALID PROFILEID parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"profileId": {strconv.Itoa(GetProfileID(t, "EDGE1")())}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1)),
				},
				"OK when VALID STATUS parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"status": {"REPORTED"}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateServerFields(map[string]interface{}{"Status": "REPORTED"})),
				},
				"OK when VALID TOPOLOGY parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"topology": {"mso-topology"}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateExpectedServers([]string{"denver-mso-org-01", "denver-mso-org-02", "edge1-cdn1-cg3", "edge2-cdn1-cg3",
							"atlanta-mid-01", "atlanta-mid-16", "atlanta-mid-17", "edgeInCachegroup3", "midInParentCachegroup",
							"midInSecondaryCachegroup", "midInSecondaryCachegroupInCDN1"})),
				},
				"OK when VALID TYPE parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"type": {"EDGE"}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateServerFields(map[string]interface{}{"Type": "EDGE"})),
				},
				"VALID SERVER LIST when using TOPOLOGY BASED DSID parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"dsId": {strconv.Itoa(GetDeliveryServiceId(t, "ds-top")())}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateExpectedServers([]string{"denver-mso-org-01"})),
				},
				"VALID SERVER TYPE when DS TOPOLOGY CONTAINS NO MIDS": {
					ClientSession: TOSession,
					RequestParams: url.Values{"dsId": {strconv.Itoa(GetDeliveryServiceId(t, "ds-based-top-with-no-mids")())}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1), validateServerTypeIsNotMid()),
				},
			},
			"POST": {
				"BAD REQUEST when BLANK PROFILEID": {
					ClientSession: TOSession,
					RequestBody:   generateServer(t, map[string]interface{}{"profileId": nil}),
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
			"PUT": {
				"OK when VALID request": {
					EndpointID:    GetServerID(t, "atlanta-edge-03"),
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"id":           GetServerID(t, "atlanta-edge-03")(),
						"cdnId":        GetCDNID(t, "cdn1")(),
						"cachegroupId": GetCacheGroupId(t, "cachegroup1")(),
						"domainName":   "updateddomainname",
						"hostName":     "atl-edge-01",
						"httpsPort":    8080,
						"interfaces": []map[string]interface{}{{
							"ipAddresses": []map[string]interface{}{
								{
									"address":        "2345:1234:12:2::4/64",
									"gateway":        "2345:1234:12:2::4",
									"serviceAddress": false,
								},
								{
									"address":        "127.0.0.13/30",
									"gateway":        "127.0.0.1",
									"serviceAddress": true,
								},
							},
							"monitor":        true,
							"mtu":            uint64(1280),
							"name":           "bond1",
							"routerHostName": "router5",
							"routerPort":     "9004",
						}},
						"physLocationId": GetPhysicalLocationID(t, "Denver")(),
						"profileId":      GetProfileID(t, "EDGE1")(),
						"rack":           "RR 119.03",
						"statusId":       GetStatusID(t, "REPORTED")(),
						"tcpPort":        8080,
						"typeId":         GetTypeId(t, "EDGE"),
						"updPending":     true,
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateServerFieldsForUpdate("atl-edge-01", map[string]interface{}{
							"CDNName": "cdn1", "Cachegroup": "cachegroup1", "DomainName": "updateddomainname", "HostName": "atl-edge-01",
							"HTTPSPort": 8080, "InterfaceName": "bond1", "MTU": uint64(1280), "PhysLocation": "Denver", "Rack": "RR 119.03",
							"TCPPort": 8080, "TypeID": GetTypeId(t, "EDGE"),
						})),
				},
				"BAD REQUEST when CHANGING XMPPID": {
					EndpointID:    GetServerID(t, "atlanta-edge-16"),
					ClientSession: TOSession,
					RequestBody: generateServer(t, map[string]interface{}{
						"id":     GetServerID(t, "atlanta-edge-16")(),
						"xmppId": "CHANGINGTHIS",
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"CONFLICT when UPDATING SERVER TYPE when ASSIGNED to DS": {
					EndpointID:    GetServerID(t, "test-ds-server-assignments"),
					ClientSession: TOSession,
					RequestBody: generateServer(t, map[string]interface{}{
						"id":           GetServerID(t, "test-ds-server-assignments")(),
						"cachegroupId": GetCacheGroupId(t, "cachegroup1")(),
						"typeId":       GetTypeId(t, "MID"),
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusConflict)),
				},
				"CONFLICT when UPDATING SERVER STATUS when its the ONLY EDGE SERVER ASSIGNED": {
					EndpointID:    GetServerID(t, "test-ds-server-assignments"),
					ClientSession: TOSession,
					RequestBody: generateServer(t, map[string]interface{}{
						"id":       GetServerID(t, "test-ds-server-assignments")(),
						"statusId": GetStatusID(t, "ADMIN_DOWN")(),
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusConflict)),
				},
				"BAD REQUEST when UPDATING CDN when LAST SERVER IN CACHEGROUP IN TOPOLOGY": {
					EndpointID:    GetServerID(t, "midInTopologyMidCg01"),
					ClientSession: TOSession,
					RequestBody: generateServer(t, map[string]interface{}{
						"id":           GetServerID(t, "midInTopologyMidCg01")(),
						"cdnId":        GetCDNID(t, "cdn1")(),
						"profileId":    GetProfileID(t, "MID1")(),
						"cachegroupId": GetCacheGroupId(t, "topology-mid-cg-01")(),
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when UPDATING CACHEGROUP when LAST SERVER IN CACHEGROUP IN TOPOLOGY": {
					EndpointID:    GetServerID(t, "midInTopologyMidCg01"),
					ClientSession: TOSession,
					RequestBody: generateServer(t, map[string]interface{}{
						"id":           GetServerID(t, "midInTopologyMidCg01")(),
						"hostName":     "midInTopologyMidCg01",
						"cdnId":        GetCDNID(t, "cdn2")(),
						"profileId":    GetProfileID(t, "CDN2_MID")(),
						"cachegroupId": GetCacheGroupId(t, "topology-mid-cg-02")(),
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when IPADDRESS EXISTS with SAME PROFILE": {
					EndpointID:    GetServerID(t, "atlanta-edge-16"),
					ClientSession: TOSession,
					RequestBody: generateServer(t, map[string]interface{}{
						"profileNames": []string{"EDGE1"},
						"interfaces": []map[string]interface{}{{
							"ipAddresses": []map[string]interface{}{{
								"address":        "127.0.0.11/22",
								"gateway":        "127.0.0.11",
								"serviceAddress": true,
							}},
							"name": "eth1",
						}},
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when BLANK HOSTNAME": {
					EndpointID:    GetServerID(t, "atlanta-edge-16"),
					ClientSession: TOSession,
					RequestBody:   generateServer(t, map[string]interface{}{"hostName": ""}),
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when BLANK DOMAINNAME": {
					EndpointID:    GetServerID(t, "atlanta-edge-16"),
					ClientSession: TOSession,
					RequestBody:   generateServer(t, map[string]interface{}{"domainName": ""}),
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"PRECONDITION FAILED when updating with IMS & IUS Headers": {
					EndpointID:     GetServerID(t, "atlanta-edge-01"),
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}},
					RequestBody: generateServer(t, map[string]interface{}{
						"id": GetServerID(t, "atlanta-edge-01")(),
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"PRECONDITION FAILED when updating with IFMATCH ETAG Header": {
					EndpointID:    GetServerID(t, "atlanta-edge-01"),
					ClientSession: TOSession,
					RequestBody: generateServer(t, map[string]interface{}{
						"id": GetServerID(t, "atlanta-edge-01")(),
					}),
					RequestHeaders: http.Header{rfc.IfMatch: {rfc.ETag(currentTime)}},
					Expectations:   utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
			},
			"DELETE": {
				"BAD REQUEST when LAST SERVER in CACHE GROUP": {
					EndpointID:    GetServerID(t, "midInTopologyMidCg01"),
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"CONFLICT when DELETING SERVER when its the ONLY EDGE SERVER ASSIGNED": {
					EndpointID:    GetServerID(t, "test-ds-server-assignments"),
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusConflict)),
				},
			},
			"GET AFTER CHANGES": {
				"OK when CHANGES made": {
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfModifiedSince: {currentTimeRFC}},
					Expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					server := tc.ServerV30{}

					if testCase.RequestBody != nil {
						dat, err := json.Marshal(testCase.RequestBody)
						assert.NoError(t, err, "Error occurred when marshalling request body: %v", err)
						err = json.Unmarshal(dat, &server)
						assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
					}

					switch method {
					case "GET", "GET AFTER CHANGES":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetServersWithHdr(&testCase.RequestParams, testCase.RequestHeaders)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.CreateServerWithHdr(server, testCase.RequestHeaders)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.UpdateServerByIDWithHdr(testCase.EndpointID(), server, testCase.RequestHeaders)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.DeleteServerByID(testCase.EndpointID())
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					}
				}
			})
		}
		t.Run("DS SERVER ASSIGNMENT REMOVED when DS UPDATED TO USE TOPOLOGY", func(t *testing.T) { UpdateDSGetServerDSID(t) })
		t.Run("STATUSLASTUPDATED ONLY CHANGES when STATUS CHANGES", func(t *testing.T) { UpdateTestServerStatusLastUpdated(t) })

	})
}

func validateServerFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected response to not be nil.")
		serverResp := resp.([]tc.ServerV30)
		for field, expected := range expectedResp {
			for _, server := range serverResp {
				switch field {
				case "CachegroupID":
					assert.RequireNotNil(t, server.CachegroupID, "Expected CachegroupID to not be nil")
					assert.Equal(t, expected, *server.CachegroupID, "Expected CachegroupID to be %d, but got %d", expected, *server.CachegroupID)
				case "Cachegroup":
					assert.RequireNotNil(t, server.Cachegroup, "Expected Cachegroup to not be nil")
					assert.Equal(t, expected, *server.Cachegroup, "Expected Cachegroup to be %s, but got %s", expected, *server.Cachegroup)
				case "CDNName":
					assert.RequireNotNil(t, server.CDNName, "Expected CDNName to not be nil")
					assert.Equal(t, expected, *server.CDNName, "Expected CDNName to be %s, but got %s", expected, *server.CDNName)
				case "CDNID":
					assert.RequireNotNil(t, server.CDNID, "Expected CDNID to not be nil")
					assert.Equal(t, expected, *server.CDNID, "Expected CDNID to be %d, but got %d", expected, *server.CDNID)
				case "DomainName":
					assert.RequireNotNil(t, server.DomainName, "Expected DomainName to not be nil")
					assert.Equal(t, expected, *server.DomainName, "Expected DomainName to be %s, but got %s", expected, *server.DomainName)
				case "HostName":
					assert.RequireNotNil(t, server.HostName, "Expected HostName to not be nil")
					assert.Equal(t, expected, *server.HostName, "Expected HostName to be %s, but got %s", expected, *server.HostName)
				case "HTTPSPort":
					assert.RequireNotNil(t, server.HTTPSPort, "Expected HTTPSPort to not be nil")
					assert.Equal(t, expected, *server.HTTPSPort, "Expected HTTPSPort to be %d, but got %d", expected, *server.HTTPSPort)
				case "InterfaceName":
					assert.RequireGreaterOrEqual(t, len(server.Interfaces), 1, "Expected Interfaces to have at least 1 interface")
					assert.Equal(t, expected, server.Interfaces[0].Name, "Expected InterfaceName to be %s, but got %s", expected, server.Interfaces[0].Name)
				case "MTU":
					assert.RequireGreaterOrEqual(t, len(server.Interfaces), 1, "Expected Interfaces to have at least 1 interface")
					assert.RequireNotNil(t, server.Interfaces[0].MTU, "Expected MTU to not be nil")
					assert.Equal(t, expected, *server.Interfaces[0].MTU, "Expected MTU to be %d, but got %d", expected, *server.Interfaces[0].MTU)
				case "PhysLocation":
					assert.RequireNotNil(t, server.PhysLocation, "Expected PhysLocation to not be nil")
					assert.Equal(t, expected, *server.PhysLocation, "Expected PhysLocation to be %s, but got %s", expected, *server.PhysLocation)
				case "ProfileID":
					assert.RequireNotNil(t, server.ProfileID, "Expected ProfileID to not be nil")
					assert.Exactly(t, expected, *server.ProfileID, "Expected ProfileID to be %v, but got %v", expected, server.ProfileID)
				case "Rack":
					assert.RequireNotNil(t, server.Rack, "Expected Rack to not be nil")
					assert.Equal(t, expected, *server.Rack, "Expected Rack to be %s, but got %s", expected, *server.Rack)
				case "Status":
					assert.RequireNotNil(t, server.Status, "Expected Status to not be nil")
					assert.Equal(t, expected, *server.Status, "Expected Status to be %s, but got %s", expected, *server.Status)
				case "TCPPort":
					assert.RequireNotNil(t, server.TCPPort, "Expected TCPPort to not be nil")
					assert.Equal(t, expected, *server.TCPPort, "Expected TCPPort to be %d, but got %d", expected, *server.TCPPort)
				case "Type":
					assert.Equal(t, expected, server.Type, "Expected Type to be %s, but got %s", expected, server.Type)
				case "TypeID":
					assert.RequireNotNil(t, server.TypeID, "Expected TypeID to not be nil")
					assert.Equal(t, expected, *server.TypeID, "Expected Type to be %d, but got %d", expected, *server.TypeID)
				default:
					t.Errorf("Expected field: %v, does not exist in response", field)
				}
			}
		}
	}
}

func validateServerFieldsForUpdate(hostName string, expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, _ interface{}, _ tc.Alerts, _ error) {
		params := url.Values{}
		params.Set("hostName", hostName)
		servers, _, err := TOSession.GetServersWithHdr(&params, nil)
		assert.NoError(t, err, "Error getting Server: %v - alerts: %+v", err, servers.Alerts)
		assert.Equal(t, 1, len(servers.Response), "Expected Server one server returned Got: %d", len(servers.Response))
		validateServerFields(expectedResp)(t, toclientlib.ReqInf{}, servers.Response, tc.Alerts{}, nil)
	}
}

func validateExpectedServers(expectedHostNames []string) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected response to not be nil.")
		serverResp := resp.([]tc.ServerV30)
		var notInResponse []string
		serverMap := make(map[string]struct{})
		for _, server := range serverResp {
			assert.RequireNotNil(t, server.HostName, "Expected server host name to not be nil.")
			serverMap[*server.HostName] = struct{}{}
		}
		for _, expected := range expectedHostNames {
			if _, exists := serverMap[expected]; !exists {
				notInResponse = append(notInResponse, expected)
			}
		}
		assert.Equal(t, len(notInResponse), 0, "%d servers missing from the response: %s", len(notInResponse), strings.Join(notInResponse, ", "))
	}
}

func validateServerTypeIsNotMid() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected response to not be nil.")
		serverResp := resp.([]tc.ServerV30)
		for _, server := range serverResp {
			assert.RequireNotNil(t, server.HostName, "Expected server host name to not be nil.")
			assert.NotEqual(t, server.Type, tc.CacheTypeMid.String(), "Expected to find no %s-typed servers but found server %s", tc.CacheTypeMid, *server.HostName)
		}
	}
}

func generateServer(t *testing.T, requestServer map[string]interface{}) map[string]interface{} {
	// map for the most basic Server a user can create
	genericServer := map[string]interface{}{
		"cdnId":        GetCDNID(t, "cdn1")(),
		"cachegroupId": GetCacheGroupId(t, "cachegroup1")(),
		"domainName":   "localhost",
		"hostName":     "testserver",
		"interfaces": []map[string]interface{}{{
			"ipAddresses": []map[string]interface{}{{
				"address":        "127.0.0.1",
				"serviceAddress": true,
			}},
			"name": "eth0",
		}},
		"physLocationId": GetPhysicalLocationID(t, "Denver")(),
		"profileId":      GetProfileID(t, "EDGE1")(),
		"statusId":       GetStatusID(t, "REPORTED")(),
		"typeId":         GetTypeId(t, "EDGE"),
	}

	for k, v := range requestServer {
		genericServer[k] = v
	}
	return genericServer
}

func GetServerID(t *testing.T, hostName string) func() int {
	return func() int {
		params := url.Values{}
		params.Set("hostName", hostName)
		serversResp, _, err := TOSession.GetServersWithHdr(&params, nil)
		assert.RequireNoError(t, err, "Get Servers Request failed with error:", err)
		assert.RequireEqual(t, 1, len(serversResp.Response), "Expected response object length 1, but got %d", len(serversResp.Response))
		assert.RequireNotNil(t, serversResp.Response[0].ID, "Expected id to not be nil")
		return *serversResp.Response[0].ID
	}
}

func UpdateTestServerStatusLastUpdated(t *testing.T) {
	const hostName = "atl-edge-01"

	params := url.Values{}
	params.Set("hostName", hostName)
	resp, _, err := TOSession.GetServersWithHdr(&params, nil)
	assert.RequireNoError(t, err, "Cannot get Server by hostname '%s': %v - alerts %+v", hostName, err, resp.Alerts)
	assert.RequireGreaterOrEqual(t, len(resp.Response), 1, "Expected at least one server to exist by hostname '%s'", hostName)
	assert.RequireNotNil(t, resp.Response[0].StatusLastUpdated, "Traffic Ops returned a representation for a server with null or undefined Status Last Updated time")
	originalServer := resp.Response[0]

	// Perform an update with no changes to status
	alerts, _, err := TOSession.UpdateServerByIDWithHdr(*originalServer.ID, originalServer, nil)
	assert.RequireNoError(t, err, "Cannot UPDATE Server by ID %d (hostname '%s'): %v - alerts: %+v", *originalServer.ID, hostName, err, alerts)

	resp, _, err = TOSession.GetServersWithHdr(&params, nil)
	assert.RequireNoError(t, err, "Cannot get Server by hostname '%s': %v - alerts %+v", hostName, err, resp.Alerts)
	assert.RequireGreaterOrEqual(t, len(resp.Response), 1, "Expected at least one server to exist by hostname '%s'", hostName)
	respServer := resp.Response[0]
	assert.RequireNotNil(t, respServer.StatusLastUpdated, "Traffic Ops returned a representation for a server with null or undefined Status Last Updated time")
	assert.Equal(t, *originalServer.StatusLastUpdated, *respServer.StatusLastUpdated, "Since status didnt change, no change in 'StatusLastUpdated' time was expected. "+
		"old value: %v, new value: %v", *originalServer.StatusLastUpdated, *respServer.StatusLastUpdated)

	// Changing the status, perform an update and make sure that statusLastUpdated changed
	newStatusID := GetStatusID(t, "ONLINE")()
	originalServer.StatusID = &newStatusID

	alerts, _, err = TOSession.UpdateServerByIDWithHdr(*originalServer.ID, originalServer, nil)
	assert.RequireNoError(t, err, "Cannot UPDATE Server by ID %d (hostname '%s'): %v - alerts: %+v", *originalServer.ID, hostName, err, alerts)

	resp, _, err = TOSession.GetServersWithHdr(&params, nil)
	assert.RequireNoError(t, err, "Cannot get Server by hostname '%s': %v - alerts %+v", hostName, err, resp.Alerts)
	assert.RequireGreaterOrEqual(t, len(resp.Response), 1, "Expected at least one server to exist by hostname '%s'", hostName)
	respServer = resp.Response[0]
	assert.RequireNotNil(t, respServer.StatusLastUpdated, "Traffic Ops returned a representation for a server with null or undefined Status Last Updated time")
	assert.NotEqual(t, *originalServer.StatusLastUpdated, *respServer.StatusLastUpdated, "Since status changed, expected 'StatusLastUpdated' to change. "+
		"old value: %v, new value: %v", *originalServer.StatusLastUpdated, *respServer.StatusLastUpdated)
}

func UpdateDSGetServerDSID(t *testing.T) {
	const hostName = "atlanta-edge-14"
	const xmlId = "ds3"
	var topology = "mso-topology"
	var firstHeaderRewrite = "first header rewrite"
	var innerHeaderRewrite = "inner header rewrite"
	var lastHeaderRewrite = "last header rewrite"

	params := url.Values{}
	params.Set("dsId", strconv.Itoa(GetDeliveryServiceId(t, xmlId)()))
	servers, _, err := TOSession.GetServersWithHdr(&params, nil)
	assert.RequireNoError(t, err, "Failed to get Servers: %v - alerts: %+v", err, servers.Alerts)
	assert.RequireGreaterOrEqual(t, len(servers.Response), 1, "Failed to get at least one Server")
	assert.RequireEqual(t, hostName, *servers.Response[0].HostName, "Expected delivery service assignment between xmlId: %v and server: %v. Got server: %v", xmlId, hostName, servers.Response[0].HostName)

	dses, _, err := TOSession.GetDeliveryServiceByXMLIDNullableWithHdr(xmlId, nil)
	assert.RequireNoError(t, err, "Failed to get Delivery Services: %v", err)
	assert.RequireEqual(t, len(dses), 1, "Failed to get at least one Delivery Service")
	ds := dses[0]

	ds.Topology = &topology
	ds.FirstHeaderRewrite = &firstHeaderRewrite
	ds.InnerHeaderRewrite = &innerHeaderRewrite
	ds.LastHeaderRewrite = &lastHeaderRewrite
	ds.EdgeHeaderRewrite = nil
	ds.MidHeaderRewrite = nil

	_, _, err = TOSession.UpdateDeliveryServiceV30WithHdr(*ds.ID, ds, nil)
	assert.RequireNoError(t, err, "Unable to add topology-related fields to deliveryservice %s: %v", xmlId, err)

	params.Set("dsId", strconv.Itoa(*ds.ID))
	servers, _, err = TOSession.GetServersWithHdr(&params, nil)
	assert.RequireNoError(t, err, "Failed to get servers by Topology-based Delivery Service ID with xmlId %s: %v - alerts: %+v", xmlId, err, servers.Alerts)
	assert.RequireGreaterOrEqual(t, len(servers.Response), 1, "Expected at least one server")
	for _, server := range servers.Response {
		assert.NotEqual(t, hostName, *server.HostName, "Server: %v was not expected to be returned.")
	}
}

func CreateTestServers(t *testing.T) {
	for _, server := range testData.Servers {
		resp, _, err := TOSession.CreateServerWithHdr(server, nil)
		assert.RequireNoError(t, err, "Could not create server '%s': %v - alerts: %+v", *server.HostName, err, resp.Alerts)
	}
}

func DeleteTestServers(t *testing.T) {
	servers, _, err := TOSession.GetServersWithHdr(nil, nil)
	assert.NoError(t, err, "Cannot get Servers: %v - alerts: %+v", err, servers.Alerts)

	for _, server := range servers.Response {
		delResp, _, err := TOSession.DeleteServerByID(*server.ID)
		assert.NoError(t, err, "Could not delete Server: %v - alerts: %+v", err, delResp.Alerts)
		// Retrieve Server to see if it got deleted
		params := url.Values{}
		params.Set("id", strconv.Itoa(*server.ID))
		getServer, _, err := TOSession.GetServersWithHdr(&params, nil)
		assert.RequireNotNil(t, server.HostName, "Expected server host name to not be nil.")
		assert.NoError(t, err, "Error deleting Server for '%s' : %v - alerts: %+v", *server.HostName, err, getServer.Alerts)
		assert.Equal(t, 0, len(getServer.Response), "Expected Server '%s' to be deleted", *server.HostName)
	}
}
