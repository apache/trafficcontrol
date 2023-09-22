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
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

func TestDeliveryServices(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServerCapabilities, DeliveryServices, ServerServerCapabilities, DeliveryServicesRequiredCapabilities, DeliveryServiceServerAssignments}, func() {

		tomorrow := time.Now().AddDate(0, 0, 1).Format(time.RFC1123)
		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)

		tenant1UserSession := utils.CreateV3Session(t, Config.TrafficOps.URL, "tenant1user", "pa$$word", Config.Default.Session.TimeoutInSecs)
		tenant2UserSession := utils.CreateV3Session(t, Config.TrafficOps.URL, "tenant2user", "pa$$word", Config.Default.Session.TimeoutInSecs)
		tenant3UserSession := utils.CreateV3Session(t, Config.TrafficOps.URL, "tenant3user", "pa$$word", Config.Default.Session.TimeoutInSecs)
		tenant4UserSession := utils.CreateV3Session(t, Config.TrafficOps.URL, "tenant4user", "pa$$word", Config.Default.Session.TimeoutInSecs)

		methodTests := utils.V3TestCase{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession: TOSession, RequestHeaders: http.Header{rfc.IfModifiedSince: {tomorrow}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID request": {
					ClientSession: TOSession, Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when ACTIVE=TRUE": {
					ClientSession: TOSession, RequestParams: url.Values{"active": {"true"}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateGetDSExpectedFields(map[string]interface{}{"Active": true})),
				},
				"OK when ACTIVE=FALSE": {
					ClientSession: TOSession, RequestParams: url.Values{"active": {"false"}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateGetDSExpectedFields(map[string]interface{}{"Active": false})),
				},
				"OK when VALID ACCESSIBLETO parameter": {
					ClientSession: TOSession, RequestParams: url.Values{"accessibleTo": {"1"}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1)),
				},
				"OK when PARENT TENANT reads DS of INACTIVE CHILD TENANT": {
					ClientSession: tenant1UserSession,
					RequestParams: url.Values{"xmlId": {"ds2"}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1)),
				},
				"EMPTY RESPONSE when DS BELONGS to TENANT but PARENT TENANT is INACTIVE": {
					ClientSession: tenant3UserSession,
					RequestParams: url.Values{"xmlId": {"ds3"}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when INACTIVE TENANT reads DS of SAME TENANCY": {
					ClientSession: tenant2UserSession,
					RequestParams: url.Values{"xmlId": {"ds2"}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when TENANT reads DS OUTSIDE TENANCY": {
					ClientSession: tenant4UserSession,
					RequestParams: url.Values{"xmlId": {"ds3"}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when CHILD TENANT reads DS of PARENT TENANT": {
					ClientSession: tenant3UserSession,
					RequestParams: url.Values{"xmlId": {"ds2"}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
			},
			"POST": {
				"BAD REQUEST when XMLID left EMPTY": {
					ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"xmlId": "",
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when XMLID is NIL": {
					ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"xmlId": nil,
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when TOPOLOGY DOESNT EXIST": {
					ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"topology": "topology-doesnt-exist",
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when creating DS with TENANCY NOT THE SAME AS CURRENT TENANT": {
					ClientSession: tenant4UserSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"tenantId": GetTenantID(t, "tenant3")(),
						"xmlId":    "test-tenancy",
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
			"PUT": {
				"OK when VALID request": {
					EndpointID: GetDeliveryServiceId(t, "ds1"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"longDesc":              "changed long desc",
						"maxDNSAnswers":         500,
						"maxOriginConnections":  5,
						"matchList":             nil,
						"maxRequestHeaderBytes": 120000,
						"xmlId":                 "ds1",
					}),
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateUpdateDSExpectedFields(map[string]interface{}{"LongDesc": "changed long desc",
							"MaxDNSAnswers": 500, "MaxOriginConnections": 5, "MaxRequestHeaderBytes": 120000,
						})),
				},
				"OK when UPDATING MINOR VERSION FIELDS": {
					EndpointID: GetDeliveryServiceId(t, "ds-test-minor-versions"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"consistentHashQueryParams": []string{"d", "e", "f"},
						"consistentHashRegex":       "foo",
						"deepCachingType":           "NEVER",
						"fqPacingRate":              41,
						"maxOriginConnections":      500,
						"routingName":               "cdn",
						"signingAlgorithm":          "uri_signing",
						"tenantId":                  GetTenantID(t, "tenant1")(),
						"trRequestHeaders":          "X-ooF\nX-raB",
						"trResponseHeaders":         "Access-Control-Max-Age: 600\nContent-Type: text/html; charset=utf-8",
						"xmlId":                     "ds-test-minor-versions",
					}),
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateUpdateDSExpectedFields(map[string]interface{}{"ConsistentHashQueryParams": []string{"d", "e", "f"},
							"ConsistentHashRegex": "foo", "DeepCachingType": tc.DeepCachingTypeNever, "FQPacingRate": 41, "MaxOriginConnections": 500,
							"SigningAlgorithm": "uri_signing", "Tenant": "tenant1", "TRRequestHeaders": "X-ooF\nX-raB",
							"TRResponseHeaders": "Access-Control-Max-Age: 600\nContent-Type: text/html; charset=utf-8",
						})),
				},
				"BAD REQUEST when INVALID REMAP TEXT": {
					EndpointID: GetDeliveryServiceId(t, "ds1"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"remapText": "@plugin=tslua.so @pparam=/opt/trafficserver/etc/trafficserver/remapPlugin1.lua\nline2",
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when MISSING SLICE PLUGIN SIZE": {
					EndpointID: GetDeliveryServiceId(t, "ds1"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"rangeRequestHandling": 3,
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when SLICE PLUGIN SIZE SET with INVALID RANGE REQUEST SETTING": {
					EndpointID: GetDeliveryServiceId(t, "ds1"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"rangeRequestHandling": 1,
						"rangeSliceBlockSize":  262144,
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when SLICE PLUGIN SIZE TOO SMALL": {
					EndpointID: GetDeliveryServiceId(t, "ds1"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"rangeRequestHandling": 3,
						"rangeSliceBlockSize":  0,
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when SLICE PLUGIN SIZE TOO LARGE": {
					EndpointID: GetDeliveryServiceId(t, "ds1"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"rangeRequestHandling": 3,
						"rangeSliceBlockSize":  40000000,
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when CHANGING TOPOLOGY of DS with ORG SERVERS ASSIGNED": {
					EndpointID: GetDeliveryServiceId(t, "ds-top"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"topology": "another-topology",
						"xmlId":    "ds-top",
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when ADDING TOPOLOGY to CLIENT STEERING DS": {
					EndpointID: GetDeliveryServiceId(t, "ds-client-steering"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"topology": "mso-topology",
						"xmlId":    "ds-client-steering",
						"typeId":   GetTypeId(t, "CLIENT_STEERING"),
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when TOPOLOGY DOESNT EXIST": {
					EndpointID: GetDeliveryServiceId(t, "ds1"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"topology": "",
						"xmlId":    "ds1",
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when ADDING TOPOLOGY to DS with DS REQUIRED CAPABILITY": {
					EndpointID: GetDeliveryServiceId(t, "ds1"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"topology": "top-for-ds-req",
						"xmlId":    "ds1",
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when ADDING TOPOLOGY to DS when NO CACHES in SAME CDN as DS": {
					EndpointID: GetDeliveryServiceId(t, "top-ds-in-cdn2"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"cdnId":    GetCDNID(t, "cdn2")(),
						"topology": "top-with-caches-in-cdn1",
						"xmlId":    "top-ds-in-cdn2",
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"OK when REMOVING TOPOLOGY": {
					EndpointID: GetDeliveryServiceId(t, "ds-based-top-with-no-mids"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"topology": nil,
						"xmlId":    "ds-based-top-with-no-mids",
					}),
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when DS with TOPOLOGY updates HEADER REWRITE FIELDS": {
					EndpointID: GetDeliveryServiceId(t, "ds-top"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"firstHeaderRewrite": "foo",
						"innerHeaderRewrite": "bar",
						"lastHeaderRewrite":  "baz",
						"topology":           "mso-topology",
						"xmlId":              "ds-top",
					}),
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"BAD REQUEST when DS with NO TOPOLOGY updates HEADER REWRITE FIELDS": {
					EndpointID: GetDeliveryServiceId(t, "ds1"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"firstHeaderRewrite": "foo",
						"innerHeaderRewrite": "bar",
						"lastHeaderRewrite":  "baz",
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when UPDATING DS OUTSIDE TENANCY": {
					EndpointID: GetDeliveryServiceId(t, "ds3"), ClientSession: tenant4UserSession,
					RequestBody:  generateDeliveryService(t, map[string]interface{}{"xmlId": "ds3"}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
				"PRECONDITION FAILED when updating with IMS & IUS Headers": {
					EndpointID: GetDeliveryServiceId(t, "ds1"), ClientSession: TOSession,
					RequestHeaders: http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}},
					RequestBody:    generateDeliveryService(t, map[string]interface{}{"xmlId": "ds1"}),
					Expectations:   utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"PRECONDITION FAILED when updating with IFMATCH ETAG Header": {
					EndpointID: GetDeliveryServiceId(t, "ds1"), ClientSession: TOSession,
					RequestBody:    generateDeliveryService(t, map[string]interface{}{"xmlId": "ds1"}),
					RequestHeaders: http.Header{rfc.IfMatch: {rfc.ETag(currentTime)}},
					Expectations:   utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
			},
			"DELETE": {
				"ERROR when DELETING DS OUTSIDE TENANCY": {
					EndpointID: GetDeliveryServiceId(t, "ds3"), ClientSession: tenant4UserSession,
					Expectations: utils.CkRequest(utils.HasError()),
				},
			},
			"GET AFTER CHANGES": {
				"OK when CHANGES made": {
					ClientSession: TOSession, RequestHeaders: http.Header{rfc.IfModifiedSince: {currentTimeRFC}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
			"DELIVERY SERVICES CAPACITY": {
				"OK when VALID request": {
					EndpointID: GetDeliveryServiceId(t, "ds1"), ClientSession: TOSession,
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					ds := tc.DeliveryServiceNullableV30{}

					if val, ok := testCase.RequestParams["accessibleTo"]; ok {
						if _, err := strconv.Atoi(val[0]); err != nil {
							testCase.RequestParams.Set("accessibleTo", strconv.Itoa(GetTenantID(t, val[0])()))
						}
					}
					if val, ok := testCase.RequestParams["cdn"]; ok {
						if _, err := strconv.Atoi(val[0]); err != nil {
							testCase.RequestParams.Set("cdn", strconv.Itoa(GetCDNID(t, val[0])()))
						}
					}
					if val, ok := testCase.RequestParams["profile"]; ok {
						if _, err := strconv.Atoi(val[0]); err != nil {
							testCase.RequestParams.Set("profile", strconv.Itoa(GetProfileID(t, val[0])()))
						}
					}
					if val, ok := testCase.RequestParams["type"]; ok {
						if _, err := strconv.Atoi(val[0]); err != nil {
							testCase.RequestParams.Set("type", strconv.Itoa(GetTypeId(t, val[0])))
						}
					}
					if val, ok := testCase.RequestParams["tenant"]; ok {
						if _, err := strconv.Atoi(val[0]); err != nil {
							testCase.RequestParams.Set("tenant", strconv.Itoa(GetTenantID(t, val[0])()))
						}
					}

					if testCase.RequestBody != nil {
						dat, err := json.Marshal(testCase.RequestBody)
						assert.NoError(t, err, "Error occurred when marshalling request body: %v", err)
						err = json.Unmarshal(dat, &ds)
						assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
					}

					switch method {
					case "GET", "GET AFTER CHANGES":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetDeliveryServicesV30WithHdr(testCase.RequestHeaders, testCase.RequestParams)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp, tc.Alerts{}, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.CreateDeliveryServiceV30(ds)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp, tc.Alerts{}, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.UpdateDeliveryServiceV30WithHdr(testCase.EndpointID(), ds, testCase.RequestHeaders)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp, tc.Alerts{}, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							resp, err := testCase.ClientSession.DeleteDeliveryService(strconv.Itoa(testCase.EndpointID()))
							for _, check := range testCase.Expectations {
								if resp != nil {
									check(t, toclientlib.ReqInf{}, nil, resp.Alerts, err)
								}
							}
						})
					case "DELIVERY SERVICES CAPACITY":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetDeliveryServiceCapacityWithHdr(strconv.Itoa(testCase.EndpointID()), testCase.RequestHeaders)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp, tc.Alerts{}, err)
							}
						})
					}
				}
			})
		}
	})
}

func validateGetDSExpectedFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		dsResp := resp.([]tc.DeliveryServiceNullableV30)
		for field, expected := range expectedResp {
			for _, ds := range dsResp {
				switch field {
				case "Active":
					assert.Equal(t, expected, *ds.Active, "Expected active to be %v, but got %v", expected, *ds.Active)
				default:
					t.Errorf("Expected field: %v, does not exist in response", field)
				}
			}
		}
	}
}

func validateUpdateDSExpectedFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		ds := resp.(tc.DeliveryServiceNullableV30)
		for field, expected := range expectedResp {
			switch field {
			case "DeepCachingType":
				assert.Equal(t, expected, *ds.DeepCachingType, "Expected deepCachingType to be %v, but got %v", expected, *ds.DeepCachingType)
			case "ConsistentHashRegex":
				assert.Equal(t, expected, *ds.ConsistentHashRegex, "Expected ConsistentHashRegex to be %v, but got %v", expected, *ds.ConsistentHashRegex)
			case "ConsistentHashQueryParams":
				assert.Exactly(t, expected, ds.ConsistentHashQueryParams, "Expected ConsistentHashQueryParams to be %v, but got %v", expected, ds.ConsistentHashQueryParams)
			case "FQPacingRate":
				assert.Equal(t, expected, *ds.FQPacingRate, "Expected FQPacingRate to be %v, but got %v", expected, *ds.FQPacingRate)
			case "LogsEnabled":
				assert.Equal(t, expected, *ds.LogsEnabled, "Expected LogsEnabled to be %v, but got %v", expected, *ds.LogsEnabled)
			case "LongDesc":
				assert.Equal(t, expected, *ds.LongDesc, "Expected LongDesc to be %v, but got %v", expected, *ds.LongDesc)
			case "MaxDNSAnswers":
				assert.Equal(t, expected, *ds.MaxDNSAnswers, "Expected LogsEnabled to be %v, but got %v", expected, *ds.MaxDNSAnswers)
			case "MaxOriginConnections":
				assert.Equal(t, expected, *ds.MaxOriginConnections, "Expected MaxOriginConnections to be %v, but got %v", expected, *ds.MaxOriginConnections)
			case "MaxRequestHeaderBytes":
				assert.Equal(t, expected, *ds.MaxRequestHeaderBytes, "Expected MaxRequestHeaderBytes to be %v, but got %v", expected, *ds.MaxRequestHeaderBytes)
			case "SigningAlgorithm":
				assert.Equal(t, expected, *ds.SigningAlgorithm, "Expected SigningAlgorithm to be %v, but got %v", expected, *ds.SigningAlgorithm)
			case "Tenant":
				assert.Equal(t, expected, *ds.Tenant, "Expected Tenant to be %v, but got %v", expected, *ds.Tenant)
			case "Topology":
				assert.Equal(t, expected, *ds.Topology, "Expected Topology to be %v, but got %v", expected, *ds.Topology)
			case "TRRequestHeaders":
				assert.Equal(t, expected, *ds.TRRequestHeaders, "Expected TRRequestHeaders to be %v, but got %v", expected, *ds.TRRequestHeaders)
			case "TRResponseHeaders":
				assert.Equal(t, expected, *ds.TRResponseHeaders, "Expected TRResponseHeaders to be %v, but got %v", expected, *ds.TRResponseHeaders)
			case "Type":
				assert.Equal(t, expected, *ds.Type, "Expected Type to be %v, but got %v", expected, *ds.Type)
			case "XMLID":
				assert.Equal(t, expected, *ds.XMLID, "Expected XMLID to be %v, but got %v", expected, *ds.XMLID)
			default:
				t.Errorf("Expected field: %v, does not exist in response", field)
			}
		}
	}
}

func GetDeliveryServiceId(t *testing.T, xmlId string) func() int {
	return func() int {
		resp, _, err := TOSession.GetDeliveryServiceByXMLIDNullableWithHdr(xmlId, http.Header{})
		assert.RequireNoError(t, err, "Get Delivery Service Request failed with error: %v", err)
		assert.RequireEqual(t, len(resp), 1, "Expected response object length 1, but got %d", len(resp))
		assert.RequireNotNil(t, resp[0].ID, "Expected id to not be nil")
		return *resp[0].ID
	}
}

func generateDeliveryService(t *testing.T, requestDS map[string]interface{}) map[string]interface{} {
	// map for the most basic HTTP Delivery Service a user can create
	genericHTTPDS := map[string]interface{}{
		"active":               true,
		"cdnName":              "cdn1",
		"cdnId":                GetCDNID(t, "cdn1")(),
		"displayName":          "test ds",
		"dscp":                 0,
		"geoLimit":             0,
		"geoProvider":          0,
		"initialDispersion":    1,
		"ipv6RoutingEnabled":   false,
		"logsEnabled":          false,
		"missLat":              0.0,
		"missLong":             0.0,
		"multiSiteOrigin":      false,
		"orgServerFqdn":        "http://ds.test",
		"protocol":             0,
		"qstringIgnore":        0,
		"rangeRequestHandling": 0,
		"regionalGeoBlocking":  false,
		"routingName":          "ccr-ds1",
		"tenant":               "tenant1",
		"type":                 tc.DSTypeHTTP,
		"typeId":               GetTypeId(t, "HTTP"),
		"xmlId":                "testds",
	}
	for k, v := range requestDS {
		genericHTTPDS[k] = v
	}
	return genericHTTPDS
}

func CreateTestDeliveryServices(t *testing.T) {
	for _, ds := range testData.DeliveryServices {
		_, _, err := TOSession.CreateDeliveryServiceV30(ds)
		assert.NoError(t, err, "Could not create Delivery Service '%s': %v", *ds.XMLID, err)
	}
}

func DeleteTestDeliveryServices(t *testing.T) {
	dses, _, err := TOSession.GetDeliveryServicesV30WithHdr(nil, nil)
	assert.NoError(t, err, "Cannot get Delivery Services: %v", err)

	for _, ds := range dses {
		delResp, err := TOSession.DeleteDeliveryService(strconv.Itoa(*ds.ID))
		assert.NoError(t, err, "Could not delete Delivery Service: %v - alerts: %+v", err, delResp.Alerts)
		// Retrieve Delivery Service to see if it got deleted
		params := url.Values{}
		params.Set("id", strconv.Itoa(*ds.ID))
		getDS, _, err := TOSession.GetDeliveryServicesV30WithHdr(http.Header{}, params)
		assert.NoError(t, err, "Error deleting Delivery Service for '%s' : %v", *ds.XMLID, err)
		assert.Equal(t, 0, len(getDS), "Expected Delivery Service '%s' to be deleted", *ds.XMLID)
	}
}
