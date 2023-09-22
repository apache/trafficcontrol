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
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

func TestDeliveryServiceRequests(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Tenants, DeliveryServiceRequests}, func() {

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
				"OK when VALID XMLID parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"xmlId": {"test-ds1"}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateGetDSRequestFields(map[string]interface{}{"XMLID": "test-ds1"})),
				},
			},
			"PUT": {
				"OK when VALID request": {
					EndpointID:    GetDeliveryServiceRequestId(t, "test-ds1"),
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"changeType": "create",
						"deliveryService": generateDeliveryService(t, map[string]interface{}{
							"displayName": "NEW DISPLAY NAME",
							"tenantId":    GetTenantID(t, "tenant1")(),
							"xmlId":       "test-ds1",
						}),
						"status": "draft",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when UPDATING STATUS FROM DRAFT TO SUBMITTED": {
					EndpointID:    GetDeliveryServiceRequestId(t, "test-ds1"),
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"changeType": "create",
						"deliveryService": generateDeliveryService(t, map[string]interface{}{
							"tenantId": GetTenantID(t, "tenant1")(),
							"xmlId":    "test-ds1",
						}),
						"status": "submitted",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"PRECONDITION FAILED when updating with IF-UNMODIFIED-SINCE Header": {
					EndpointID:     GetDeliveryServiceRequestId(t, "test-ds1"),
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}},
					RequestBody:    map[string]interface{}{},
					Expectations:   utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"PRECONDITION FAILED when updating with IFMATCH ETAG Header": {
					EndpointID:     GetDeliveryServiceRequestId(t, "test-ds1"),
					ClientSession:  TOSession,
					RequestBody:    map[string]interface{}{},
					RequestHeaders: http.Header{rfc.IfMatch: {rfc.ETag(currentTime)}},
					Expectations:   utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
			},
			"POST": {
				"OK when VALID request": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"changeType": "create",
						"deliveryService": generateDeliveryService(t, map[string]interface{}{
							"ccrDnsTtl":          30,
							"deepCachingType":    "NEVER",
							"initialDispersion":  3,
							"ipv6RoutingEnabled": true,
							"longDesc":           "long desc",
							"orgServerFqdn":      "http://example.test",
							"profileName":        nil,
							"tenant":             "root",
							"xmlId":              "test-ds2",
						}),
						"status": "draft",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"BAD REQUEST when MISSING REQUIRED Delivery Service FIELDS": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"changeType": "create",
						"deliveryService": map[string]interface{}{
							"type":  "HTTP",
							"xmlId": "test-ds-fields",
						},
						"status": "draft",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"bad request when missing Delivery Service definition": {
					ClientSession: TOSession,
					RequestBody: map[string]any{
						"changeType": "create",
						"status":     "draft",
					},
				},
				"bad request when missing status": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"changeType": "create",
						"deliveryService": generateDeliveryService(t, map[string]interface{}{
							"ccrDnsTtl":          30,
							"deepCachingType":    "NEVER",
							"initialDispersion":  3,
							"ipv6RoutingEnabled": true,
							"longDesc":           "long desc",
							"orgServerFqdn":      "http://example.test",
							"profileName":        nil,
							"tenant":             "root",
							"xmlId":              "test-ds2",
						}),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"bad request when missing change type": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"deliveryService": generateDeliveryService(t, map[string]interface{}{
							"ccrDnsTtl":          30,
							"deepCachingType":    "NEVER",
							"initialDispersion":  3,
							"ipv6RoutingEnabled": true,
							"longDesc":           "long desc",
							"orgServerFqdn":      "http://example.test",
							"profileName":        nil,
							"tenant":             "root",
							"xmlId":              "test-ds2",
						}),
						"status": "draft",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when VALIDATION RULES ARE NOT FOLLOWED": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"changeType": "create",
						"deliveryService": map[string]interface{}{
							"ccrDnsTtl":            30,
							"deepCachingType":      "NEVER",
							"displayName":          strings.Repeat("X", 49),
							"dscp":                 0,
							"geoLimit":             0,
							"geoProvider":          1,
							"infoUrl":              "xxx",
							"initialDispersion":    1,
							"ipv6RoutingEnabled":   true,
							"logsEnabled":          true,
							"longDesc":             "long desc",
							"missLat":              0.0,
							"missLong":             0.0,
							"multiSiteOrigin":      false,
							"orgServerFqdn":        "http://example.test",
							"protocol":             0,
							"qstringIgnore":        0,
							"rangeRequestHandling": 0,
							"regionalGeoBlocking":  true,
							"routingName":          strings.Repeat("X", 1) + "." + strings.Repeat("X", 48),
							"tenant":               "tenant1",
							"type":                 "HTTP",
							"xmlId":                "X " + strings.Repeat("X", 46),
						},
						"status": "draft",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when NON-DRAFT": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"changeType": "create",
						"deliveryService": map[string]interface{}{
							"active":               false,
							"cdnName":              "cdn1",
							"displayName":          "Testing transitions",
							"dscp":                 3,
							"geoLimit":             1,
							"geoProvider":          1,
							"initialDispersion":    1,
							"ipv6RoutingEnabled":   true,
							"logsEnabled":          true,
							"missLat":              0.0,
							"missLong":             0.0,
							"multiSiteOrigin":      false,
							"orgServerFqdn":        "http://example.test",
							"protocol":             0,
							"qstringIgnore":        0,
							"rangeRequestHandling": 0,
							"regionalGeoBlocking":  true,
							"routingName":          "goodroute",
							"tenant":               "tenant1",
							"type":                 "HTTP",
							"xmlId":                "test-transitions",
						},
						"status": "pending",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when ALREADY EXISTS": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"changeType": "create",
						"deliveryService": map[string]interface{}{
							"active":               true,
							"cdnName":              "cdn1",
							"displayName":          "Good Kabletown CDN",
							"dscp":                 1,
							"geoLimit":             1,
							"geoProvider":          1,
							"initialDispersion":    1,
							"ipv6RoutingEnabled":   true,
							"logsEnabled":          true,
							"missLat":              0.0,
							"missLong":             0.0,
							"multiSiteOrigin":      false,
							"orgServerFqdn":        "http://example.test",
							"protocol":             0,
							"qstringIgnore":        0,
							"rangeRequestHandling": 0,
							"regionalGeoBlocking":  true,
							"routingName":          "goodroute",
							"tenant":               "tenant1",
							"type":                 "HTTP",
							"xmlId":                "test-ds1",
						},
						"status": "draft",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
			"DELETE": {
				"OK when VALID request": {
					EndpointID:    GetDeliveryServiceRequestId(t, "test-deletion"),
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
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
					dsReq := tc.DeliveryServiceRequest{}

					if testCase.RequestBody != nil {
						dat, err := json.Marshal(testCase.RequestBody)
						assert.NoError(t, err, "Error occurred when marshalling request body: %v", err)
						err = json.Unmarshal(dat, &dsReq)
						assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
					}

					switch method {
					case "GET", "GET AFTER CHANGES":
						t.Run(name, func(t *testing.T) {
							if name == "OK when VALID XMLID parameter" {
								resp, reqInf, err := testCase.ClientSession.GetDeliveryServiceRequestByXMLIDWithHdr(testCase.RequestParams["xmlId"][0], testCase.RequestHeaders)
								for _, check := range testCase.Expectations {
									check(t, reqInf, resp, tc.Alerts{}, err)
								}
							} else {
								resp, reqInf, err := testCase.ClientSession.GetDeliveryServiceRequestsWithHdr(testCase.RequestHeaders)
								for _, check := range testCase.Expectations {
									check(t, reqInf, resp, tc.Alerts{}, err)
								}
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.CreateDeliveryServiceRequest(dsReq)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.UpdateDeliveryServiceRequestByIDWithHdr(testCase.EndpointID(), dsReq, testCase.RequestHeaders)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.DeleteDeliveryServiceRequestByID(testCase.EndpointID())
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

func GetDeliveryServiceRequestId(t *testing.T, xmlId string) func() int {
	return func() int {
		resp, _, err := TOSession.GetDeliveryServiceRequestByXMLIDWithHdr(xmlId, http.Header{})
		assert.RequireNoError(t, err, "Get Delivery Service Requests failed with error: %v", err)
		assert.RequireGreaterOrEqual(t, len(resp), 1, "Expected delivery service requests response object length of atleast 1, but got %d", len(resp))
		assert.RequireNotNil(t, resp[0].ID, "Expected id to not be nil")
		return resp[0].ID
	}
}

func validateGetDSRequestFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		dsReqResp := resp.([]tc.DeliveryServiceRequest)
		for field, expected := range expectedResp {
			for _, ds := range dsReqResp {
				switch field {
				case "XMLID":
					assert.Equal(t, expected, ds.DeliveryService.XMLID, "Expected XMLID to be %v, but got %v", expected, ds.DeliveryService.XMLID)
				default:
					t.Errorf("Expected field: %v, does not exist in response", field)
				}
			}
		}
	}
}

func CreateTestDeliveryServiceRequests(t *testing.T) {
	for _, dsr := range testData.DeliveryServiceRequests {
		respDSR, _, err := TOSession.CreateDeliveryServiceRequest(dsr)
		assert.NoError(t, err, "Could not create Delivery Service Requests: %v - alerts: %+v", err, respDSR.Alerts)
	}
}

func DeleteTestDeliveryServiceRequests(t *testing.T) {
	resp, _, err := TOSession.GetDeliveryServiceRequestsWithHdr(http.Header{})
	assert.NoError(t, err, "Cannot get Delivery Service Requests: %v", err)
	for _, request := range resp {
		alert, _, err := TOSession.DeleteDeliveryServiceRequestByID(request.ID)
		assert.NoError(t, err, "Cannot delete Delivery Service Request #%d: %v - alerts: %+v", request.ID, err, alert.Alerts)

		// Retrieve the DeliveryServiceRequest to see if it got deleted
		dsr, _, err := TOSession.GetDeliveryServiceRequestByIDWithHdr(request.ID, http.Header{})
		assert.NoError(t, err, "Unexpected error fetching Delivery Service Request #%d after deletion: %v", request.ID, err)
		assert.Equal(t, len(dsr), 0, "Expected Delivery Service Request #%d to be deleted, but it was found in Traffic Ops", request.ID)
	}
}
