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
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	totest "github.com/apache/trafficcontrol/v8/lib/go-tc/totestv4"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"
)

func TestServers(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, DeliveryServices, DeliveryServiceServerAssignments}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)
		methodTests := utils.V4TestCase{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {tomorrow}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID HOSTNAME parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"hostName": {"atlanta-edge-01"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						validateServerFields(map[string]interface{}{"HostName": "atlanta-edge-01"})),
				},
				"OK when VALID CACHEGROUP parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"cachegroup": {strconv.Itoa(totest.GetCacheGroupId(t, TOSession, "cachegroup1")())}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateServerFields(map[string]interface{}{"CachegroupID": totest.GetCacheGroupId(t, TOSession, "cachegroup1")()})),
				},
				"OK when VALID CACHEGROUPNAME parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"cachegroupName": {"topology-mid-cg-01"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateServerFields(map[string]interface{}{"Cachegroup": "topology-mid-cg-01"})),
				},
				"OK when VALID CDN parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"cdn": {strconv.Itoa(totest.GetCDNID(t, TOSession, "cdn2")())}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateServerFields(map[string]interface{}{"CDNID": totest.GetCDNID(t, TOSession, "cdn2")()})),
				},
				"OK when VALID DSID parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"dsId": {strconv.Itoa(totest.GetDeliveryServiceId(t, TOSession, "test-ds-server-assignments")())}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateExpectedServers([]string{"test-ds-server-assignments", "test-mso-org-01"})),
				},
				"OK when VALID PARENTCACHEGROUP parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"parentCacheGroup": {strconv.Itoa(totest.GetCacheGroupId(t, TOSession, "parentCachegroup")())}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1)),
				},
				"OK when VALID PROFILENAME parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"profileName": {"EDGE1"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1)),
				},
				"OK when VALID STATUS parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"status": {"REPORTED"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateServerFields(map[string]interface{}{"Status": "REPORTED"})),
				},
				"OK when VALID TOPOLOGY parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"topology": {"mso-topology"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateExpectedServers([]string{"denver-mso-org-01", "denver-mso-org-02", "edge1-cdn1-cg3", "edge2-cdn1-cg3",
							"atlanta-mid-01", "atlanta-mid-16", "atlanta-mid-17", "edgeInCachegroup3", "midInParentCachegroup",
							"midInSecondaryCachegroup", "midInSecondaryCachegroupInCDN1", "test-mso-org-01"})),
				},
				"OK when VALID TYPE parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"type": {"EDGE"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateServerFields(map[string]interface{}{"Type": "EDGE"})),
				},
				"VALID SERVER LIST when using TOPOLOGY BASED DSID parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"dsId": {strconv.Itoa(totest.GetDeliveryServiceId(t, TOSession, "ds-top")())}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateExpectedServers([]string{"denver-mso-org-01"})),
				},
				"VALID SERVER TYPE when DS TOPOLOGY CONTAINS NO MIDS": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"dsId": {strconv.Itoa(totest.GetDeliveryServiceId(t, TOSession, "ds-based-top-with-no-mids")())}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1), validateServerTypeIsNotMid()),
				},
				"EMPTY RESPONSE when INVALID DSID parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"dsId": {"999999"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"FIRST RESULT when LIMIT=1": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateServerPagination("limit")),
				},
				"SECOND RESULT when LIMIT=1 OFFSET=1": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}, "offset": {"1"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateServerPagination("offset")),
				},
				"SECOND RESULT when LIMIT=1 PAGE=2": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}, "page": {"2"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateServerPagination("page")),
				},
				"BAD REQUEST when INVALID LIMIT parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"limit": {"-2"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID OFFSET parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"limit": {"1"}, "offset": {"0"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID PAGE parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"limit": {"1"}, "page": {"0"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
			"POST": {
				"BAD REQUEST when BLANK PROFILENAMES": {
					ClientSession: TOSession,
					RequestBody:   generateServer(t, map[string]interface{}{"profileNames": []string{""}}),
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
			"PUT": {
				"OK when VALID request": {
					EndpointID:    totest.GetServerID(t, TOSession, "atlanta-edge-03"),
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"id":           totest.GetServerID(t, TOSession, "atlanta-edge-03")(),
						"cdnId":        totest.GetCDNID(t, TOSession, "cdn1")(),
						"cachegroupId": totest.GetCacheGroupId(t, TOSession, "cachegroup1")(),
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
						"profileNames":   []string{"EDGE1"},
						"rack":           "RR 119.03",
						"statusId":       GetStatusID(t, "REPORTED")(),
						"tcpPort":        8080,
						"typeId":         totest.GetTypeId(t, TOSession, "EDGE"),
						"updPending":     true,
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateServerFieldsForUpdate("atl-edge-01", map[string]interface{}{
							"CDNName": "cdn1", "Cachegroup": "cachegroup1", "DomainName": "updateddomainname", "HostName": "atl-edge-01",
							"HTTPSPort": 8080, "InterfaceName": "bond1", "MTU": uint64(1280), "PhysLocation": "Denver", "Rack": "RR 119.03",
							"TCPPort": 8080, "TypeID": totest.GetTypeId(t, TOSession, "EDGE"),
						})),
				},
				"BAD REQUEST when CHANGING XMPPID": {
					EndpointID:    totest.GetServerID(t, TOSession, "atlanta-edge-16"),
					ClientSession: TOSession,
					RequestBody: generateServer(t, map[string]interface{}{
						"id":     totest.GetServerID(t, TOSession, "atlanta-edge-16")(),
						"xmppId": "CHANGINGTHIS",
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"CONFLICT when UPDATING SERVER TYPE when ASSIGNED to DS": {
					EndpointID:    totest.GetServerID(t, TOSession, "test-ds-server-assignments"),
					ClientSession: TOSession,
					RequestBody: generateServer(t, map[string]interface{}{
						"id":           totest.GetServerID(t, TOSession, "test-ds-server-assignments")(),
						"cachegroupId": totest.GetCacheGroupId(t, TOSession, "cachegroup1")(),
						"typeId":       totest.GetTypeId(t, TOSession, "MID"),
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusConflict)),
				},
				"CONFLICT when UPDATING SERVER STATUS when its the ONLY EDGE SERVER ASSIGNED": {
					EndpointID:    totest.GetServerID(t, TOSession, "test-ds-server-assignments"),
					ClientSession: TOSession,
					RequestBody: generateServer(t, map[string]interface{}{
						"id":       totest.GetServerID(t, TOSession, "test-ds-server-assignments")(),
						"statusId": GetStatusID(t, "ADMIN_DOWN")(),
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusConflict)),
				},
				"CONFLICT when UPDATING SERVER STATUS when its the ONLY ORG SERVER ASSIGNED": {
					EndpointID:    totest.GetServerID(t, TOSession, "test-mso-org-01"),
					ClientSession: TOSession,
					RequestBody: generateServer(t, map[string]interface{}{
						"id":       totest.GetServerID(t, TOSession, "test-mso-org-01")(),
						"statusId": GetStatusID(t, "ADMIN_DOWN")(),
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusConflict)),
				},
				"BAD REQUEST when UPDATING CDN when LAST SERVER IN CACHEGROUP IN TOPOLOGY": {
					EndpointID:    totest.GetServerID(t, TOSession, "midInTopologyMidCg01"),
					ClientSession: TOSession,
					RequestBody: generateServer(t, map[string]interface{}{
						"id":           totest.GetServerID(t, TOSession, "midInTopologyMidCg01")(),
						"cdnId":        totest.GetCDNID(t, TOSession, "cdn1")(),
						"profileNames": []string{"MID1"},
						"cachegroupId": totest.GetCacheGroupId(t, TOSession, "topology-mid-cg-01")(),
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when UPDATING CACHEGROUP when LAST SERVER IN CACHEGROUP IN TOPOLOGY": {
					EndpointID:    totest.GetServerID(t, TOSession, "midInTopologyMidCg01"),
					ClientSession: TOSession,
					RequestBody: generateServer(t, map[string]interface{}{
						"id":           totest.GetServerID(t, TOSession, "midInTopologyMidCg01")(),
						"hostName":     "midInTopologyMidCg01",
						"cdnId":        totest.GetCDNID(t, TOSession, "cdn2")(),
						"profileNames": []string{"CDN2_MID"},
						"cachegroupId": totest.GetCacheGroupId(t, TOSession, "topology-mid-cg-02")(),
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when IPADDRESS EXISTS with SAME PROFILE": {
					EndpointID:    totest.GetServerID(t, TOSession, "atlanta-edge-16"),
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
					EndpointID:    totest.GetServerID(t, TOSession, "atlanta-edge-16"),
					ClientSession: TOSession,
					RequestBody:   generateServer(t, map[string]interface{}{"hostName": ""}),
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when BLANK DOMAINNAME": {
					EndpointID:    totest.GetServerID(t, TOSession, "atlanta-edge-16"),
					ClientSession: TOSession,
					RequestBody:   generateServer(t, map[string]interface{}{"domainName": ""}),
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"PRECONDITION FAILED when updating with IMS & IUS Headers": {
					EndpointID:    totest.GetServerID(t, TOSession, "atlanta-edge-01"),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}}},
					RequestBody: generateServer(t, map[string]interface{}{
						"id": totest.GetServerID(t, TOSession, "atlanta-edge-01")(),
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"PRECONDITION FAILED when updating with IFMATCH ETAG Header": {
					EndpointID:    totest.GetServerID(t, TOSession, "atlanta-edge-01"),
					ClientSession: TOSession,
					RequestBody: generateServer(t, map[string]interface{}{
						"id": totest.GetServerID(t, TOSession, "atlanta-edge-01")(),
					}),
					RequestOpts:  client.RequestOptions{Header: http.Header{rfc.IfMatch: {rfc.ETag(currentTime)}}},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
			},
			"DELETE": {
				"BAD REQUEST when LAST SERVER in CACHE GROUP": {
					EndpointID:    totest.GetServerID(t, TOSession, "midInTopologyMidCg01"),
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"CONFLICT when DELETING SERVER when its the ONLY EDGE SERVER ASSIGNED": {
					EndpointID:    totest.GetServerID(t, TOSession, "test-ds-server-assignments"),
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusConflict)),
				},
			},
			"GET AFTER CHANGES": {
				"OK when CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {currentTimeRFC}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					server := tc.ServerV4{}

					if testCase.RequestBody != nil {
						dat, err := json.Marshal(testCase.RequestBody)
						assert.NoError(t, err, "Error occurred when marshalling request body: %v", err)
						err = json.Unmarshal(dat, &server)
						assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
					}

					switch method {
					case "GET", "GET AFTER CHANGES":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetServers(testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.CreateServer(server, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.UpdateServer(testCase.EndpointID(), server, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.DeleteServer(testCase.EndpointID(), testCase.RequestOpts)
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
		serverResp := resp.([]tc.ServerV40)
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
				case "ProfileNames":
					assert.Exactly(t, expected, server.ProfileNames, "Expected ProfileNames to be %v, but got %v", expected, server.ProfileNames)
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

func validateServerFieldsForUpdate(hostname string, expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, _ interface{}, _ tc.Alerts, _ error) {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("hostName", hostname)
		servers, _, err := TOSession.GetServers(opts)
		assert.NoError(t, err, "Error getting Server: %v - alerts: %+v", err, servers.Alerts)
		assert.Equal(t, 1, len(servers.Response), "Expected Server one server returned Got: %d", len(servers.Response))
		validateServerFields(expectedResp)(t, toclientlib.ReqInf{}, servers.Response, tc.Alerts{}, nil)
	}
}

func validateExpectedServers(expectedHostnames []string) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected response to not be nil.")
		serverResp := resp.([]tc.ServerV40)
		var notInResponse []string
		serverMap := make(map[string]struct{})
		for _, server := range serverResp {
			assert.RequireNotNil(t, server.HostName, "Expected server host name to not be nil.")
			serverMap[*server.HostName] = struct{}{}
		}
		for _, expected := range expectedHostnames {
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
		serverResp := resp.([]tc.ServerV40)
		for _, server := range serverResp {
			assert.RequireNotNil(t, server.HostName, "Expected server host name to not be nil.")
			assert.NotEqual(t, server.Type, tc.CacheTypeMid.String(), "Expected to find no %s-typed servers but found server %s", tc.CacheTypeMid, *server.HostName)
		}
	}
}

func validateServerPagination(paginationParam string) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected response to not be nil.")
		paginationResp := resp.([]tc.ServerV40)
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("orderby", "id")
		respBase, _, err := TOSession.GetServers(opts)
		assert.RequireNoError(t, err, "Cannot get Servers: %v - alerts: %+v", err, respBase.Alerts)

		ds := respBase.Response
		assert.RequireGreaterOrEqual(t, len(ds), 3, "Need at least 3 Servers in Traffic Ops to test pagination support, found: %d", len(ds))
		switch paginationParam {
		case "limit:":
			assert.Exactly(t, ds[:1], paginationResp, "expected GET Servers with limit = 1 to return first result")
		case "offset":
			assert.Exactly(t, ds[1:2], paginationResp, "expected GET Servers with limit = 1, offset = 1 to return second result")
		case "page":
			assert.Exactly(t, ds[1:2], paginationResp, "expected GET Servers with limit = 1, page = 2 to return second result")
		}
	}
}

func generateServer(t *testing.T, requestServer map[string]interface{}) map[string]interface{} {
	// map for the most basic Server a user can create
	genericServer := map[string]interface{}{
		"cdnId":        totest.GetCDNID(t, TOSession, "cdn1")(),
		"cachegroupId": totest.GetCacheGroupId(t, TOSession, "cachegroup1")(),
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
		"profileNames":   []string{"EDGE1"},
		"statusId":       GetStatusID(t, "REPORTED")(),
		"typeId":         totest.GetTypeId(t, TOSession, "EDGE"),
	}

	for k, v := range requestServer {
		genericServer[k] = v
	}
	return genericServer
}

func UpdateTestServerStatusLastUpdated(t *testing.T) {
	const hostName = "atl-edge-01"

	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("hostName", hostName)
	resp, _, err := TOSession.GetServers(opts)
	assert.RequireNoError(t, err, "Cannot get Server by hostname '%s': %v - alerts %+v", hostName, err, resp.Alerts)
	assert.RequireGreaterOrEqual(t, len(resp.Response), 1, "Expected at least one server to exist by hostname '%s'", hostName)
	assert.RequireNotNil(t, resp.Response[0].StatusLastUpdated, "Traffic Ops returned a representation for a server with null or undefined Status Last Updated time")
	originalServer := resp.Response[0]

	// Perform an update with no changes to status
	alerts, _, err := TOSession.UpdateServer(*originalServer.ID, originalServer, client.RequestOptions{})
	assert.RequireNoError(t, err, "Cannot UPDATE Server by ID %d (hostname '%s'): %v - alerts: %+v", *originalServer.ID, hostName, err, alerts)

	resp, _, err = TOSession.GetServers(opts)
	assert.RequireNoError(t, err, "Cannot get Server by hostname '%s': %v - alerts %+v", hostName, err, resp.Alerts)
	assert.RequireGreaterOrEqual(t, len(resp.Response), 1, "Expected at least one server to exist by hostname '%s'", hostName)
	respServer := resp.Response[0]
	assert.RequireNotNil(t, respServer.StatusLastUpdated, "Traffic Ops returned a representation for a server with null or undefined Status Last Updated time")
	assert.Equal(t, *originalServer.StatusLastUpdated, *respServer.StatusLastUpdated, "Since status didnt change, no change in 'StatusLastUpdated' time was expected. "+
		"old value: %v, new value: %v", *originalServer.StatusLastUpdated, *respServer.StatusLastUpdated)

	// Changing the status, perform an update and make sure that statusLastUpdated changed
	newStatusID := GetStatusID(t, "ONLINE")()
	originalServer.StatusID = &newStatusID

	alerts, _, err = TOSession.UpdateServer(*originalServer.ID, originalServer, client.RequestOptions{})
	assert.RequireNoError(t, err, "Cannot UPDATE Server by ID %d (hostname '%s'): %v - alerts: %+v", *originalServer.ID, hostName, err, alerts)

	resp, _, err = TOSession.GetServers(opts)
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

	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("dsId", strconv.Itoa(totest.GetDeliveryServiceId(t, TOSession, xmlId)()))
	servers, _, err := TOSession.GetServers(opts)
	assert.RequireNoError(t, err, "Failed to get Servers: %v - alerts: %+v", err, servers.Alerts)
	assert.RequireGreaterOrEqual(t, len(servers.Response), 1, "Failed to get at least one Server")
	assert.RequireEqual(t, hostName, *servers.Response[0].HostName, "Expected delivery service assignment between xmlId: %v and server: %v. Got server: %v", xmlId, hostName, servers.Response[0].HostName)

	opts.QueryParameters.Set("xmlId", xmlId)
	dses, _, err := TOSession.GetDeliveryServices(opts)
	assert.RequireNoError(t, err, "Failed to get Delivery Services: %v - alerts: %+v", err, dses.Alerts)
	assert.RequireEqual(t, len(dses.Response), 1, "Failed to get at least one Delivery Service")
	ds := dses.Response[0]

	ds.Topology = &topology
	ds.FirstHeaderRewrite = &firstHeaderRewrite
	ds.InnerHeaderRewrite = &innerHeaderRewrite
	ds.LastHeaderRewrite = &lastHeaderRewrite
	ds.EdgeHeaderRewrite = nil
	ds.MidHeaderRewrite = nil

	updResp, _, err := TOSession.UpdateDeliveryService(*ds.ID, ds, client.RequestOptions{})
	assert.RequireNoError(t, err, "Unable to add topology-related fields to deliveryservice %s: %v - alerts: %+v", xmlId, err, updResp.Alerts)

	opts.QueryParameters.Set("dsId", strconv.Itoa(*ds.ID))
	servers, _, err = TOSession.GetServers(opts)
	assert.RequireNoError(t, err, "Failed to get servers by Topology-based Delivery Service ID with xmlId %s: %v - alerts: %+v", xmlId, err, servers.Alerts)
	assert.RequireGreaterOrEqual(t, len(servers.Response), 1, "Expected at least one server")
	for _, server := range servers.Response {
		assert.NotEqual(t, hostName, *server.HostName, "Server: %v was not expected to be returned.")
	}
}
