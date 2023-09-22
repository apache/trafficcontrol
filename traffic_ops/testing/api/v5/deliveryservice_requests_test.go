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

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

// this resets the IDs of things attached to a DS, which needs to be done
// because the WithObjs flow destroys and recreates those object IDs
// non-deterministically with each test - BUT, the client method permanently
// alters the DSR structures by adding these referential IDs. Older clients
// got away with it by not making 'DeliveryService' a pointer, but to add
// original/requested fields you need to sometimes allow each to be nil, so
// this is a problem that needs to be solved at some point.
// A better solution _might_ be to reload all the test fixtures every time
// to wipe any and all referential modifications made to any test data, but
// for now that's overkill.
func resetDS(ds *tc.DeliveryServiceV5) {
	if ds == nil {
		return
	}
	ds.CDNID = -1
	ds.ID = nil
	ds.ProfileID = nil
	ds.TenantID = -1
	ds.TypeID = -1
}

func TestDeliveryServiceRequests(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Tenants, DeliveryServiceRequests}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.V5TestCase{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {tomorrow}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID XMLID parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"xmlId": {"test-ds1"}}},
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
						"requested": generateDeliveryService(t, map[string]interface{}{
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
						"requested": generateDeliveryService(t, map[string]interface{}{
							"tenantId": GetTenantID(t, "tenant1")(),
							"xmlId":    "test-ds1",
						}),
						"status": "submitted",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validatePutDSRequestFields(map[string]interface{}{"STATUS": tc.RequestStatusSubmitted})),
				},
				"PRECONDITION FAILED when updating with IF-UNMODIFIED-SINCE Header": {
					EndpointID:    GetDeliveryServiceRequestId(t, "test-ds1"),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}}},
					RequestBody:   map[string]interface{}{},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"PRECONDITION FAILED when updating with IFMATCH ETAG Header": {
					EndpointID:    GetDeliveryServiceRequestId(t, "test-ds1"),
					ClientSession: TOSession,
					RequestBody:   map[string]interface{}{},
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfMatch: {rfc.ETag(currentTime)}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
			},
			"POST": {
				"CREATED when VALID request": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"changeType": "create",
						"requested": generateDeliveryService(t, map[string]interface{}{
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
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusCreated)),
				},
				"BAD REQUEST when MISSING REQUIRED FIELDS": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"changeType": "create",
						"requested": map[string]interface{}{
							"type":  "HTTP",
							"xmlId": "test-ds-fields",
						},
						"status": "draft",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when VALIDATION RULES ARE NOT FOLLOWED": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"changeType": "create",
						"requested": map[string]interface{}{
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
						"requested": map[string]interface{}{
							"active":               "INACTIVE",
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
						"requested": map[string]interface{}{
							"active":               "ACTIVE",
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
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {currentTimeRFC}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					var dsReq tc.DeliveryServiceRequestV5

					if testCase.RequestBody != nil {
						dat, err := json.Marshal(testCase.RequestBody)
						assert.NoError(t, err, "Error occurred when marshalling request body: %v", err)
						err = json.Unmarshal(dat, &dsReq)
						assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
					}

					switch method {
					case "GET", "GET AFTER CHANGES":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetDeliveryServiceRequests(testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.CreateDeliveryServiceRequest(dsReq, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.UpdateDeliveryServiceRequest(testCase.EndpointID(), dsReq, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.DeleteDeliveryServiceRequest(testCase.EndpointID(), testCase.RequestOpts)
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

func GetDeliveryServiceRequestId(t *testing.T, xmlId string) func() int {
	return func() int {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("xmlId", xmlId)
		resp, _, err := TOSession.GetDeliveryServiceRequests(opts)
		assert.RequireNoError(t, err, "Get Delivery Service Requests failed with error: %v", err)
		assert.RequireGreaterOrEqual(t, len(resp.Response), 1, "Expected delivery service requests response object length of atleast 1, but got %d", len(resp.Response))
		assert.RequireNotNil(t, resp.Response[0].ID, "Expected id to not be nil")
		return *resp.Response[0].ID
	}
}

func validateGetDSRequestFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		dsReqResp := resp.([]tc.DeliveryServiceRequestV5)
		for field, expected := range expectedResp {
			for _, ds := range dsReqResp {
				switch field {
				case "XMLID":
					assert.RequireNotNil(t, ds.Requested, "expected 'requested' DS in DSR to not be null/undefined")
					assert.Equal(t, expected, ds.Requested.XMLID, "Expected XMLID to be %v, but got %v", expected, ds.Requested.XMLID)
				default:
					t.Errorf("Expected field: %v, does not exist in response", field)
				}
			}
		}
	}
}

func validatePutDSRequestFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		dsReqResp := resp.(tc.DeliveryServiceRequestV5)
		for field, expected := range expectedResp {
			switch field {
			case "STATUS":
				assert.Equal(t, expected, dsReqResp.Status, "Expected status to be %v, but got %v", expected, dsReqResp.Status)
			default:
				t.Errorf("Expected field: %v, does not exist in response", field)
			}
		}
	}
}

func CreateTestDeliveryServiceRequests(t *testing.T) {
	for _, dsr := range testData.DeliveryServiceRequests {
		resetDS(dsr.Original)
		resetDS(dsr.Requested)
		_, _, err := TOSession.CreateDeliveryServiceRequest(dsr, client.RequestOptions{})
		assert.NoError(t, err, "Could not create Delivery Service Requests: %v", err)
	}
}

func DeleteTestDeliveryServiceRequests(t *testing.T) {
	resp, _, err := TOSession.GetDeliveryServiceRequests(client.RequestOptions{})
	assert.NoError(t, err, "Cannot get Delivery Service Requests: %v - alerts: %+v", err, resp.Alerts)
	for _, request := range resp.Response {
		alert, _, err := TOSession.DeleteDeliveryServiceRequest(*request.ID, client.RequestOptions{})
		assert.NoError(t, err, "Cannot delete Delivery Service Request #%d: %v - alerts: %+v", request.ID, err, alert.Alerts)

		// Retrieve the DeliveryServiceRequest to see if it got deleted
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.Itoa(*request.ID))
		dsr, _, err := TOSession.GetDeliveryServiceRequests(opts)
		assert.NoError(t, err, "Unexpected error fetching Delivery Service Request #%d after deletion: %v - alerts: %+v", *request.ID, err, dsr.Alerts)
		assert.Equal(t, len(dsr.Response), 0, "Expected Delivery Service Request #%d to be deleted, but it was found in Traffic Ops", *request.ID)
	}
}
