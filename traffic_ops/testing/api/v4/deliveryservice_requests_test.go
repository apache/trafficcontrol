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
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"
)

func TestDeliveryServiceRequests(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Tenants, Users, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServerCapabilities, ServiceCategories, DeliveryServices, DeliveryServiceRequests}, func() {

		t.Run("update DSR crud", testUpdateDSR)

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
				"BAD REQUEST when using LONG DESCRIPTION 2 and 3 fields": {
					EndpointID:    GetDeliveryServiceRequestId(t, "test-ds1"),
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"changeType": "create",
						"requested": generateDeliveryService(t, map[string]interface{}{
							"longDesc1": "long desc 1",
							"longDesc2": "long desc 2",
							"tenantId":  GetTenantID(t, "tenant1")(),
							"xmlId":     "test-ds1",
						}),
						"status": "draft",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
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
						"requested": map[string]interface{}{
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
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {currentTimeRFC}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					dsReq := tc.DeliveryServiceRequestV4{}

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

func testUpdateDSR(t *testing.T) {
	resp, _, err := TOSession.GetDeliveryServices(client.RequestOptions{})
	if err != nil {
		t.Fatalf("failed to get Delivery Services: %#v - response: %+v", err, resp)
	}

	if len(resp.Response) < 1 {
		t.Fatal("need at least one Delivery Service to test updating a DS with a DSR")
	}

	ds := resp.Response[0]
	if ds.DisplayName == nil {
		t.Fatalf("Traffic Ops returned a DS with a nil Display Name: %+v", ds)
	}
	*ds.DisplayName += " - Update DSR test"

	dsr := tc.DeliveryServiceRequestV4{
		ChangeType: tc.DSRChangeTypeUpdate,
		Requested:  &ds,
		Status:     tc.RequestStatusDraft,
	}
	creationResp, _, err := TOSession.CreateDeliveryServiceRequest(dsr, client.RequestOptions{})
	if err != nil {
		t.Fatalf("failed to create an update DSR: %#v - response: %+v", err, creationResp)
	}

	id := creationResp.Response.ID
	if id == nil {
		t.Fatalf("Traffic Ops returned a created DSR without an ID: %+v", creationResp)
	}
	deleteResp, _, err := TOSession.DeleteDeliveryServiceRequest(*id, client.RequestOptions{})
	if err != nil {
		t.Errorf("failed to delete the created update DSR: %#v - response: %+v", err, deleteResp)
	}
}

func GetDeliveryServiceRequestId(t *testing.T, xmlId string) func() int {
	return func() int {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("xmlId", xmlId)
		resp, _, err := TOSession.GetDeliveryServiceRequests(opts)
		assert.RequireNoError(t, err, "Get Delivery Service Request '%s' failed with error: %v", xmlId, err)
		assert.RequireGreaterOrEqual(t, len(resp.Response), 1, "Expected delivery service requests response object length of atleast 1, but got %d", len(resp.Response))
		assert.RequireNotNil(t, resp.Response[0].ID, "Expected id to not be nil")
		return *resp.Response[0].ID
	}
}

func validateGetDSRequestFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		dsReqResp := resp.([]tc.DeliveryServiceRequestV40)
		for field, expected := range expectedResp {
			for _, ds := range dsReqResp {
				switch field {
				case "XMLID":
					assert.Equal(t, expected, *ds.Requested.XMLID, "Expected XMLID to be %v, but got %v", expected, *ds.Requested.XMLID)
				default:
					t.Errorf("Expected field: %v, does not exist in response", field)
				}
			}
		}
	}
}

func validatePutDSRequestFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		dsReqResp := resp.(tc.DeliveryServiceRequestV40)
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
