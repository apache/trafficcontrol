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
	"net/url"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

func TestOrigins(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Coordinates, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Users, Topologies, ServiceCategories, DeliveryServices, Origins}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)

		tenant4UserSession := utils.CreateV3Session(t, Config.TrafficOps.URL, "tenant4user", "pa$$word", Config.Default.Session.TimeoutInSecs)

		methodTests := utils.V3TestCaseT[tc.Origin]{
			"GET": {
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1)),
				},
				"OK when VALID NAME parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"name": {"origin1"}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						validateOriginsFields(map[string]interface{}{"Name": "origin1"})),
				},
				"EMPTY RESPONSE when CHILD TENANT reads PARENT TENANT ORIGIN": {
					ClientSession: tenant4UserSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
			},
			"GET QUERY PARAMS": {
				"FIRST RESULT when LIMIT=1": {
					ClientSession: TOSession,
					RequestParams: url.Values{"orderby": {"id"}, "limit": {"1"}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateOriginsPagination("limit")),
				},
				"SECOND RESULT when LIMIT=1 OFFSET=1": {
					ClientSession: TOSession,
					RequestParams: url.Values{"orderby": {"id"}, "limit": {"1"}, "offset": {"1"}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateOriginsPagination("offset")),
				},
				"SECOND RESULT when LIMIT=1 PAGE=2": {
					ClientSession: TOSession,
					RequestParams: url.Values{"orderby": {"id"}, "limit": {"1"}, "page": {"2"}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateOriginsPagination("page")),
				},
				"BAD REQUEST when INVALID LIMIT parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"limit": {"-2"}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID OFFSET parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"limit": {"1"}, "offset": {"0"}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID PAGE parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"limit": {"1"}, "page": {"0"}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
			"PUT": {
				"OK when VALID request": {
					EndpointID:    GetOriginID(t, "origin2"),
					ClientSession: TOSession,
					RequestBody: tc.Origin{
						Name:            util.Ptr("origin2"),
						Cachegroup:      util.Ptr("multiOriginCachegroup"),
						Coordinate:      util.Ptr("coordinate2"),
						DeliveryService: util.Ptr("ds3"),
						FQDN:            util.Ptr("originupdated.example.com"),
						IPAddress:       util.Ptr("1.2.3.4"),
						IP6Address:      util.Ptr("0000::1111"),
						Port:            util.Ptr(1234),
						Protocol:        util.Ptr("http"),
						TenantID:        util.Ptr(GetTenantID(t, "tenant2")()),
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateOriginsUpdateCreateFields("origin2", map[string]interface{}{"Cachegroup": "multiOriginCachegroup", "Coordinate": "coordinate2", "DeliveryService": "ds3",
							"FQDN": "originupdated.example.com", "IPAddress": "1.2.3.4", "IP6Address": "0000::1111", "Port": 1234, "Protocol": "http", "Tenant": "tenant2"})),
				},
				"FORBIDDEN when CHILD TENANT updates PARENT TENANT ORIGIN": {
					EndpointID:    GetOriginID(t, "origin2"),
					ClientSession: tenant4UserSession,
					RequestBody: tc.Origin{
						Name:              util.Ptr("testtenancy"),
						DeliveryServiceID: util.Ptr(GetDeliveryServiceId(t, "ds1")()),
						FQDN:              util.Ptr("testtenancy.example.com"),
						Protocol:          util.Ptr("http"),
						TenantID:          util.Ptr(GetTenantID(t, "tenant1")()),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
				"NOT FOUND when ORIGIN DOESNT EXIST": {
					EndpointID:    func() int { return 1111111 },
					ClientSession: TOSession,
					RequestBody: tc.Origin{
						Name:              util.Ptr("testid"),
						DeliveryServiceID: util.Ptr(GetDeliveryServiceId(t, "ds1")()),
						FQDN:              util.Ptr("testid.example.com"),
						Protocol:          util.Ptr("http"),
						TenantID:          util.Ptr(GetTenantID(t, "tenant1")()),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"PRECONDITION FAILED when updating with IMS & IUS Headers": {
					EndpointID:     GetOriginID(t, "origin2"),
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}},
					RequestBody: tc.Origin{
						Name:            util.Ptr("origin2"),
						Cachegroup:      util.Ptr("originCachegroup"),
						DeliveryService: util.Ptr("ds2"),
						FQDN:            util.Ptr("origin2.example.com"),
						Protocol:        util.Ptr("http"),
						TenantID:        util.Ptr(GetTenantID(t, "tenant1")()),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"PRECONDITION FAILED when updating with IFMATCH ETAG Header": {
					EndpointID:    GetOriginID(t, "origin2"),
					ClientSession: TOSession,
					RequestBody: tc.Origin{
						Name:            util.Ptr("origin2"),
						Cachegroup:      util.Ptr("originCachegroup"),
						DeliveryService: util.Ptr("ds2"),
						FQDN:            util.Ptr("origin2.example.com"),
						Protocol:        util.Ptr("http"),
						TenantID:        util.Ptr(GetTenantID(t, "tenant1")()),
					},
					RequestHeaders: http.Header{rfc.IfMatch: {rfc.ETag(currentTime)}},
					Expectations:   utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
			},
			"DELETE": {
				"NOT FOUND when DOESNT EXIST": {
					EndpointID:    func() int { return 11111111 },
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"FORBIDDEN when CHILD TENANT deletes PARENT TENANT ORIGIN": {
					EndpointID:    GetOriginID(t, "origin2"),
					ClientSession: tenant4UserSession,
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							if name == "OK when VALID NAME parameter" {
								resp, reqInf, err := testCase.ClientSession.GetOriginByName(testCase.RequestParams["name"][0])
								for _, check := range testCase.Expectations {
									check(t, reqInf, resp, tc.Alerts{}, err)
								}
							} else {
								resp, reqInf, err := testCase.ClientSession.GetOrigins()
								for _, check := range testCase.Expectations {
									check(t, reqInf, resp, tc.Alerts{}, err)
								}
							}
						})
					case "GET QUERY PARAMS":
						t.Run(name, func(t *testing.T) {
							queryParams := "?" + testCase.RequestParams.Encode()
							resp, reqInf, err := testCase.ClientSession.GetOriginsByQueryParams(queryParams)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp, tc.Alerts{}, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.CreateOrigin(testCase.RequestBody)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.UpdateOriginByIDWithHdr(testCase.EndpointID(), testCase.RequestBody, testCase.RequestHeaders)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.DeleteOriginByID(testCase.EndpointID())
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

func validateOriginsFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Origin response to not be nil.")
		originResp := resp.([]tc.Origin)
		for field, expected := range expectedResp {
			for _, origin := range originResp {
				switch field {
				case "Cachegroup":
					assert.RequireNotNil(t, origin.Cachegroup, "Expected Cachegroup to not be nil.")
					assert.Equal(t, expected, *origin.Cachegroup, "Expected Cachegroup to be %v, but got %s", expected, *origin.Cachegroup)
				case "CachegroupID":
					assert.RequireNotNil(t, origin.CachegroupID, "Expected CachegroupID to not be nil.")
					assert.Equal(t, expected, *origin.CachegroupID, "Expected CachegroupID to be %v, but got %d", expected, *origin.Cachegroup)
				case "Coordinate":
					assert.RequireNotNil(t, origin.Coordinate, "Expected Coordinate to not be nil.")
					assert.Equal(t, expected, *origin.Coordinate, "Expected Coordinate to be %v, but got %s", expected, *origin.Coordinate)
				case "CoordinateID":
					assert.RequireNotNil(t, origin.CoordinateID, "Expected CoordinateID to not be nil.")
					assert.Equal(t, expected, *origin.CoordinateID, "Expected CoordinateID to be %v, but got %d", expected, *origin.CoordinateID)
				case "DeliveryService":
					assert.RequireNotNil(t, origin.DeliveryService, "Expected DeliveryService to not be nil.")
					assert.Equal(t, expected, *origin.DeliveryService, "Expected DeliveryService to be %v, but got %s", expected, *origin.DeliveryService)
				case "DeliveryServiceID":
					assert.RequireNotNil(t, origin.DeliveryServiceID, "Expected DeliveryServiceID to not be nil.")
					assert.Equal(t, expected, *origin.DeliveryServiceID, "Expected DeliveryServiceID to be %v, but got %d", expected, *origin.DeliveryServiceID)
				case "FQDN":
					assert.RequireNotNil(t, origin.FQDN, "Expected FQDN to not be nil.")
					assert.Equal(t, expected, *origin.FQDN, "Expected FQDN to be %v, but got %s", expected, *origin.FQDN)
				case "ID":
					assert.RequireNotNil(t, origin.ID, "Expected ID to not be nil.")
					assert.Equal(t, expected, *origin.ID, "Expected ID to be %v, but got %d", expected, *origin.ID)
				case "IPAddress":
					assert.RequireNotNil(t, origin.IPAddress, "Expected IPAddress to not be nil.")
					assert.Equal(t, expected, *origin.IPAddress, "Expected IPAddress to be %v, but got %s", expected, *origin.IPAddress)
				case "IP6Address":
					assert.RequireNotNil(t, origin.IP6Address, "Expected IP6Address to not be nil.")
					assert.Equal(t, expected, *origin.IP6Address, "Expected IP6Address to be %v, but got %s", expected, *origin.IP6Address)
				case "IsPrimary":
					assert.RequireNotNil(t, origin.IsPrimary, "Expected IsPrimary to not be nil.")
					assert.Equal(t, expected, *origin.IsPrimary, "Expected IsPrimary to be %v, but got %v", expected, *origin.IsPrimary)
				case "Name":
					assert.RequireNotNil(t, origin.Name, "Expected Name to not be nil.")
					assert.Equal(t, expected, *origin.Name, "Expected Name to be %v, but got %s", expected, *origin.Name)
				case "Port":
					assert.RequireNotNil(t, origin.Port, "Expected Port to not be nil.")
					assert.Equal(t, expected, *origin.Port, "Expected Port to be %v, but got %d", expected, *origin.Port)
				case "Profile":
					assert.RequireNotNil(t, origin.Profile, "Expected Profile to not be nil.")
					assert.Equal(t, expected, *origin.Profile, "Expected Profile to be %v, but got %s", expected, *origin.Profile)
				case "ProfileID":
					assert.RequireNotNil(t, origin.ProfileID, "Expected ProfileID to not be nil.")
					assert.Equal(t, expected, *origin.ProfileID, "Expected ProfileID to be %v, but got %d", expected, *origin.ProfileID)
				case "Protocol":
					assert.RequireNotNil(t, origin.Protocol, "Expected Protocol to not be nil.")
					assert.Equal(t, expected, *origin.Protocol, "Expected Tenant to be %v, but got %s", expected, *origin.Protocol)
				case "Tenant":
					assert.RequireNotNil(t, origin.Tenant, "Expected Tenant to not be nil.")
					assert.Equal(t, expected, *origin.Tenant, "Expected Tenant to be %v, but got %s", expected, *origin.Tenant)
				case "TenantID":
					assert.RequireNotNil(t, origin.TenantID, "Expected TenantID to not be nil.")
					assert.Equal(t, expected, *origin.TenantID, "Expected TenantID to be %v, but got %d", expected, *origin.TenantID)
				default:
					t.Errorf("Expected field: %v, does not exist in response", field)
				}
			}
		}
	}
}

func validateOriginsUpdateCreateFields(name string, expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		origin, _, err := TOSession.GetOriginByName(name)
		assert.RequireNoError(t, err, "Error getting Origin: %v", err)
		assert.RequireEqual(t, 1, len(origin), "Expected one Origin returned Got: %d", len(origin))
		validateOriginsFields(expectedResp)(t, toclientlib.ReqInf{}, origin, tc.Alerts{}, nil)
	}
}

func validateOriginsPagination(paginationParam string) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		paginationResp := resp.([]tc.Origin)
		respBase, _, err := TOSession.GetOriginsByQueryParams("?orderby=id")
		assert.RequireNoError(t, err, "Cannot get Origins: %v", err)

		origin := respBase
		assert.RequireGreaterOrEqual(t, len(origin), 3, "Need at least 3 Origins in Traffic Ops to test pagination support, found: %d", len(origin))
		switch paginationParam {
		case "limit:":
			assert.Exactly(t, origin[:1], paginationResp, "expected GET Origins with limit = 1 to return first result")
		case "offset":
			assert.Exactly(t, origin[1:2], paginationResp, "expected GET Origins with limit = 1, offset = 1 to return second result")
		case "page":
			assert.Exactly(t, origin[1:2], paginationResp, "expected GET Origins with limit = 1, page = 2 to return second result")
		}
	}
}

func GetOriginID(t *testing.T, name string) func() int {
	return func() int {
		origins, _, err := TOSession.GetOriginByName(name)
		assert.RequireNoError(t, err, "Get Origins Request failed with error:", err)
		assert.RequireEqual(t, 1, len(origins), "Expected response object length 1, but got %d", len(origins))
		assert.RequireNotNil(t, origins[0].ID, "Expected ID to not be nil.")
		return *origins[0].ID
	}
}

func CreateTestOrigins(t *testing.T) {
	for _, origin := range testData.Origins {
		resp, _, err := TOSession.CreateOrigin(origin)
		assert.RequireNoError(t, err, "Could not create Origins: %v - alerts: %+v", err, resp.Alerts)
	}
}

func DeleteTestOrigins(t *testing.T) {
	origins, _, err := TOSession.GetOrigins()
	assert.NoError(t, err, "Cannot get Origins : %v", err)

	for _, origin := range origins {
		assert.RequireNotNil(t, origin.ID, "Expected origin ID to not be nil.")
		assert.RequireNotNil(t, origin.Name, "Expected origin ID to not be nil.")
		assert.RequireNotNil(t, origin.IsPrimary, "Expected origin ID to not be nil.")
		if !*origin.IsPrimary {
			alerts, _, err := TOSession.DeleteOriginByID(*origin.ID)
			assert.NoError(t, err, "Unexpected error deleting Origin '%s' (#%d): %v - alerts: %+v", *origin.Name, *origin.ID, err, alerts.Alerts)
			// Retrieve the Origin to see if it got deleted
			getOrigin, _, err := TOSession.GetOriginByID(*origin.ID)
			assert.NoError(t, err, "Error getting Origin '%s' after deletion: %v", *origin.Name, err)
			assert.Equal(t, 0, len(getOrigin), "Expected Origin '%s' to be deleted, but it was found in Traffic Ops", *origin.Name)
		}
	}
}
