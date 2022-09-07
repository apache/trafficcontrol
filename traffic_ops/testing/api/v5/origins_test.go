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
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/assert"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/traffic_ops/v5-client"
)

func TestOrigins(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Coordinates, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Users, Topologies, ServiceCategories, DeliveryServices, Origins}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)

		tenant4UserSession := utils.CreateV5Session(t, Config.TrafficOps.URL, "tenant4user", "pa$$word", Config.Default.Session.TimeoutInSecs)

		methodTests := utils.V5TestCase{
			"GET": {
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1)),
				},
				"OK when VALID NAME parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"origin1"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						validateOriginsFields(map[string]interface{}{"Name": "origin1"})),
				},
				"OK when VALID DELIVERYSERVICE parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"deliveryservice": {strconv.Itoa(GetDeliveryServiceId(t, "ds1")())}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateOriginsFields(map[string]interface{}{"DeliveryServiceID": GetDeliveryServiceId(t, "ds1")()})),
				},
				"OK when VALID CACHEGROUP parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"cachegroup": {strconv.Itoa(GetCacheGroupId(t, "originCachegroup")())}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateOriginsFields(map[string]interface{}{"CachegroupID": GetCacheGroupId(t, "originCachegroup")()})),
				},
				"OK when VALID COORDINATE parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"coordinate": {strconv.Itoa(GetCoordinateID(t, "coordinate1")())}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateOriginsFields(map[string]interface{}{"Coordinate": "coordinate1"})),
				},
				"OK when VALID PROFILEID parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"profileId": {strconv.Itoa(GetProfileID(t, "ATS_EDGE_TIER_CACHE")())}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateOriginsFields(map[string]interface{}{"ProfileID": GetProfileID(t, "ATS_EDGE_TIER_CACHE")()})),
				},
				"OK when VALID PRIMARY parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"primary": {"true"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateOriginsFields(map[string]interface{}{"IsPrimary": true})),
				},
				"OK when VALID TENANT parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"tenant": {strconv.Itoa(GetTenantID(t, "tenant1")())}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateOriginsFields(map[string]interface{}{"TenantID": GetTenantID(t, "tenant1")()})),
				},
				"EMPTY RESPONSE when CHILD TENANT reads PARENT TENANT ORIGIN": {
					ClientSession: tenant4UserSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"tenant": {strconv.Itoa(GetTenantID(t, "tenant3")())}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when NAME parameter that DOESNT EXIST": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"doesntexist"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when DELIVERYSERVICE parameter that DOESNT EXIST": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"deliveryservice": {"1000000"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when CACHEGROUP parameter that DOESNT EXIST": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"cachegroup": {"1000000"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when COORDINATE parameter that DOESNT EXIST": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"coordinate": {"1000000"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when PROFILEID parameter that DOESNT EXIST": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"profileId": {"1000000"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"BAD REQUEST when INVALID PRIMARY parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"primary": {"1000000"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"EMPTY RESPONSE when TENANT parameter that DOESNT EXIST": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"tenant": {"1000000"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"FIRST RESULT when LIMIT=1": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateOriginsPagination("limit")),
				},
				"SECOND RESULT when LIMIT=1 OFFSET=1": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}, "offset": {"1"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateOriginsPagination("offset")),
				},
				"SECOND RESULT when LIMIT=1 PAGE=2": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}, "page": {"2"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateOriginsPagination("page")),
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
				"BAD REQUEST when ALREADY EXISTS": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"name":            "origin1",
						"cachegroup":      "originCachegroup",
						"Coordinate":      "coordinate1",
						"deliveryService": "ds1",
						"fqdn":            "origin1.example.com",
						"ipAddress":       "1.2.3.4",
						"ip6Address":      "dead:beef:cafe::42",
						"port":            1234,
						"Profile":         "ATS_EDGE_TIER_CACHE",
						"protocol":        "http",
						"tenantId":        GetTenantID(t, "tenant1")(),
						"isPrimary":       true,
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"FORBIDDEN when CHILD TENANT CREATES ORIGIN OUTSIDE TENANCY": {
					ClientSession: tenant4UserSession,
					RequestBody: map[string]interface{}{
						"name":              "originTenancyTest",
						"cachegroup":        "originCachegroup",
						"deliveryServiceId": GetDeliveryServiceId(t, "ds1")(),
						"fqdn":              "origintenancy.example.com",
						"protocol":          "http",
						"tenantId":          GetTenantID(t, "tenant3")(),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
				"NOT FOUND when CACHEGROUP DOESNT EXIST": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"name":              "testcg",
						"cachegroupId":      10000000,
						"deliveryServiceId": GetDeliveryServiceId(t, "ds1")(),
						"fqdn":              "test.cachegroupId.com",
						"protocol":          "http",
						"tenantId":          GetTenantID(t, "tenant1")(),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"NOT FOUND when PROFILEID DOESNT EXIST": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"name":              "testprofile",
						"deliveryServiceId": GetDeliveryServiceId(t, "ds1")(),
						"fqdn":              "test.profileId.com",
						"profileId":         1000000,
						"protocol":          "http",
						"tenantId":          GetTenantID(t, "tenant1")(),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"NOT FOUND when COORDINATE DOESNT EXIST": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"name":              "testcoordinate",
						"coordinateId":      10000000,
						"deliveryServiceId": GetDeliveryServiceId(t, "ds1")(),
						"fqdn":              "test.coordinate.com",
						"protocol":          "http",
						"tenantId":          GetTenantID(t, "tenant1")(),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"FORBIDDEN when INVALID TENANT": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"name":              "testtenant",
						"deliveryServiceId": GetDeliveryServiceId(t, "ds1")(),
						"fqdn":              "test.tenant.com",
						"protocol":          "http",
						"tenantId":          11111111,
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
				"BAD REQUEST when INVALID PROTOCOL": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"name":              "testprotocol",
						"deliveryServiceId": GetDeliveryServiceId(t, "ds1")(),
						"fqdn":              "test.protocol.com",
						"protocol":          "httttpppss",
						"tenantId":          GetTenantID(t, "tenant1")(),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID IPV4 ADDRESS": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"name":              "testip",
						"deliveryServiceId": GetDeliveryServiceId(t, "ds1")(),
						"fqdn":              "test.ip.com",
						"ipAddress":         "311.255.323.412",
						"protocol":          "http",
						"tenantId":          GetTenantID(t, "tenant1")(),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID IPV6 ADDRESS": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"name":              "testipv6",
						"deliveryServiceId": GetDeliveryServiceId(t, "ds1")(),
						"fqdn":              "origin1.example.com",
						"ip6Address":        "badipv6::addresss",
						"protocol":          "http",
						"tenantId":          GetTenantID(t, "tenant1")(),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
			"PUT": {
				"OK when VALID request": {
					EndpointId:    GetOriginID(t, "origin2"),
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"name":            "origin2",
						"cachegroup":      "multiOriginCachegroup",
						"Coordinate":      "coordinate2",
						"deliveryService": "ds3",
						"fqdn":            "originupdated.example.com",
						"ipAddress":       "1.2.3.4",
						"ip6Address":      "0000::1111",
						"port":            1234,
						"protocol":        "http",
						"tenantId":        GetTenantID(t, "tenant2")(),
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateOriginsUpdateCreateFields("origin2", map[string]interface{}{"Cachegroup": "multiOriginCachegroup", "Coordinate": "coordinate2", "DeliveryService": "ds3",
							"FQDN": "originupdated.example.com", "IPAddress": "1.2.3.4", "IP6Address": "0000::1111", "Port": 1234, "Protocol": "http", "Tenant": "tenant2"})),
				},
				"FORBIDDEN when CHILD TENANT updates PARENT TENANT ORIGIN": {
					EndpointId:    GetOriginID(t, "origin2"),
					ClientSession: tenant4UserSession,
					RequestBody: map[string]interface{}{
						"name":              "testtenancy",
						"deliveryServiceId": GetDeliveryServiceId(t, "ds1")(),
						"fqdn":              "testtenancy.example.com",
						"protocol":          "http",
						"tenantId":          GetTenantID(t, "tenant1")(),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
				"NOT FOUND when ORIGIN DOESNT EXIST": {
					EndpointId:    func() int { return 1111111 },
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"name":              "testid",
						"deliveryServiceId": GetDeliveryServiceId(t, "ds1")(),
						"fqdn":              "testid.example.com",
						"protocol":          "http",
						"tenantId":          GetTenantID(t, "tenant1")(),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"BAD REQUEST when DELIVERY SERVICE DOESNT EXIST": {
					EndpointId:    GetOriginID(t, "origin2"),
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"name":              "origin2",
						"deliveryServiceId": 11111111,
						"fqdn":              "origin2.example.com",
						"protocol":          "http",
						"tenantId":          GetTenantID(t, "tenant1")(),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"NOT FOUND when CACHEGROUP DOESNT EXIST": {
					EndpointId:    GetOriginID(t, "origin2"),
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"name":              "origin2",
						"cachegroupId":      1111111,
						"deliveryServiceId": GetDeliveryServiceId(t, "ds1")(),
						"fqdn":              "origin2.example.com",
						"protocol":          "http",
						"tenantId":          GetTenantID(t, "tenant1")(),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"NOT FOUND when PROFILEID DOESNT EXIST": {
					EndpointId:    GetOriginID(t, "origin2"),
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"name":              "origin2",
						"cachegroup":        "originCachegroup",
						"deliveryServiceId": GetDeliveryServiceId(t, "ds1")(),
						"fqdn":              "origin2.example.com",
						"profileId":         11111111,
						"protocol":          "http",
						"tenantId":          GetTenantID(t, "tenant1")(),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"NOT FOUND when COORDINATE DOESNT EXIST": {
					EndpointId:    GetOriginID(t, "origin2"),
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"name":              "origin2",
						"cachegroup":        "originCachegroup",
						"coordinateId":      1111111,
						"deliveryServiceId": GetDeliveryServiceId(t, "ds1")(),
						"fqdn":              "origin2.example.com",
						"protocol":          "http",
						"tenantId":          GetTenantID(t, "tenant1")(),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"FORBIDDEN when INVALID TENANT": {
					EndpointId:    GetOriginID(t, "origin2"),
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"name":              "origin1",
						"cachegroup":        "originCachegroup",
						"deliveryServiceId": GetDeliveryServiceId(t, "ds1")(),
						"fqdn":              "origin1.example.com",
						"protocol":          "http",
						"tenantId":          1111111,
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
				"BAD REQUEST when INVALID PROTOCOL": {
					EndpointId:    GetOriginID(t, "origin2"),
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"name":              "origin2",
						"cachegroup":        "originCachegroup",
						"deliveryServiceId": GetDeliveryServiceId(t, "ds1")(),
						"fqdn":              "origin2.example.com",
						"protocol":          "htttttpssss",
						"tenantId":          GetTenantID(t, "tenant1")(),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID IPV4 ADDRESS": {
					EndpointId:    GetOriginID(t, "origin2"),
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"name":              "origin2",
						"cachegroup":        "originCachegroup",
						"deliveryServiceId": GetDeliveryServiceId(t, "ds2")(),
						"fqdn":              "origin2.example.com",
						"ipAddress":         "300.254.123.1",
						"protocol":          "http",
						"tenantId":          GetTenantID(t, "tenant1")(),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID IPV6 ADDRESS": {
					EndpointId:    GetOriginID(t, "origin2"),
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"name":              "origin2",
						"cachegroup":        "originCachegroup",
						"deliveryServiceId": GetDeliveryServiceId(t, "ds2")(),
						"fqdn":              "origin2.example.com",
						"ip6Address":        "test::42",
						"protocol":          "http",
						"tenantId":          GetTenantID(t, "tenant1")(),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID PORT": {
					EndpointId:    GetOriginID(t, "origin2"),
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"name":              "origin2",
						"cachegroup":        "originCachegroup",
						"deliveryServiceId": GetDeliveryServiceId(t, "ds2")(),
						"fqdn":              "origin2.example.com",
						"port":              80000,
						"protocol":          "http",
						"tenantId":          GetTenantID(t, "tenant1")(),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"PRECONDITION FAILED when updating with IMS & IUS Headers": {
					EndpointId:    GetOriginID(t, "origin2"),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}}},
					RequestBody: map[string]interface{}{
						"name":            "origin2",
						"cachegroup":      "originCachegroup",
						"deliveryService": "ds2",
						"fqdn":            "origin2.example.com",
						"protocol":        "http",
						"tenantId":        GetTenantID(t, "tenant1")(),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"PRECONDITION FAILED when updating with IFMATCH ETAG Header": {
					EndpointId:    GetOriginID(t, "origin2"),
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"name":            "origin2",
						"cachegroup":      "originCachegroup",
						"deliveryService": "ds2",
						"fqdn":            "origin2.example.com",
						"protocol":        "http",
						"tenantId":        GetTenantID(t, "tenant1")(),
					},
					RequestOpts:  client.RequestOptions{Header: http.Header{rfc.IfMatch: {rfc.ETag(currentTime)}}},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
			},
			"DELETE": {
				"NOT FOUND when DOESNT EXIST": {
					EndpointId:    func() int { return 11111111 },
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"FORBIDDEN when CHILD TENANT deletes PARENT TENANT ORIGIN": {
					EndpointId:    GetOriginID(t, "origin2"),
					ClientSession: tenant4UserSession,
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					origin := tc.Origin{}

					if testCase.RequestBody != nil {
						dat, err := json.Marshal(testCase.RequestBody)
						assert.NoError(t, err, "Error occurred when marshalling request body: %v", err)
						err = json.Unmarshal(dat, &origin)
						assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
					}

					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetOrigins(testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.CreateOrigin(origin, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.UpdateOrigin(testCase.EndpointId(), origin, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.DeleteOrigin(testCase.EndpointId(), testCase.RequestOpts)
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
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("name", name)
		origin, _, err := TOSession.GetOrigins(opts)
		assert.RequireNoError(t, err, "Error getting Origin: %v - alerts: %+v", err, origin.Alerts)
		assert.RequireEqual(t, 1, len(origin.Response), "Expected one Origin returned Got: %d", len(origin.Response))
		validateOriginsFields(expectedResp)(t, toclientlib.ReqInf{}, origin.Response, tc.Alerts{}, nil)
	}
}

func validateOriginsPagination(paginationParam string) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		paginationResp := resp.([]tc.Origin)

		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("orderby", "id")
		respBase, _, err := TOSession.GetOrigins(opts)
		assert.RequireNoError(t, err, "Cannot get Origins: %v - alerts: %+v", err, respBase.Alerts)

		origin := respBase.Response
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
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("name", name)
		origins, _, err := TOSession.GetOrigins(opts)
		assert.RequireNoError(t, err, "Get Origins Request failed with error:", err)
		assert.RequireEqual(t, 1, len(origins.Response), "Expected response object length 1, but got %d", len(origins.Response))
		assert.RequireNotNil(t, origins.Response[0].ID, "Expected ID to not be nil.")
		return *origins.Response[0].ID
	}
}

func CreateTestOrigins(t *testing.T) {
	for _, origin := range testData.Origins {
		resp, _, err := TOSession.CreateOrigin(origin, client.RequestOptions{})
		assert.RequireNoError(t, err, "Could not create Origins: %v - alerts: %+v", err, resp.Alerts)
	}
}

func DeleteTestOrigins(t *testing.T) {
	origins, _, err := TOSession.GetOrigins(client.RequestOptions{})
	assert.NoError(t, err, "Cannot get Origins : %v - alerts: %+v", err, origins.Alerts)

	for _, origin := range origins.Response {
		assert.RequireNotNil(t, origin.ID, "Expected origin ID to not be nil.")
		assert.RequireNotNil(t, origin.Name, "Expected origin ID to not be nil.")
		assert.RequireNotNil(t, origin.IsPrimary, "Expected origin ID to not be nil.")
		if !*origin.IsPrimary {
			alerts, _, err := TOSession.DeleteOrigin(*origin.ID, client.RequestOptions{})
			assert.NoError(t, err, "Unexpected error deleting Origin '%s' (#%d): %v - alerts: %+v", *origin.Name, *origin.ID, err, alerts.Alerts)
			// Retrieve the Origin to see if it got deleted
			opts := client.NewRequestOptions()
			opts.QueryParameters.Set("id", strconv.Itoa(*origin.ID))
			getOrigin, _, err := TOSession.GetOrigins(opts)
			assert.NoError(t, err, "Error getting Origin '%s' after deletion: %v - alerts: %+v", *origin.Name, err, getOrigin.Alerts)
			assert.Equal(t, 0, len(getOrigin.Response), "Expected Origin '%s' to be deleted, but it was found in Traffic Ops", *origin.Name)
		}
	}
}

func alertsHaveError(alerts []tc.Alert, err string) bool {
	for _, alert := range alerts {
		if alert.Level == tc.ErrorLevel.String() && strings.Contains(alert.Text, err) {
			return true
		}
	}
	return false
}
