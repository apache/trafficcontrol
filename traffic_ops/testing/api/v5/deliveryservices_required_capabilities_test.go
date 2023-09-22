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
	"fmt"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
	"testing"
)

//func TestDeliveryServicesRequiredCapabilities(t *testing.T) {
//	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, ServerCapabilities, Topologies, ServiceCategories, DeliveryServices, DeliveryServiceServerAssignments, ServerServerCapabilities, DeliveryServicesRequiredCapabilities}, func() {
//
//		currentTime := time.Now().UTC().Add(-15 * time.Second)
//		currentTimeRFC := currentTime.Format(time.RFC1123)
//		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)
//
//		methodTests := utils.TestCase[client.Session, client.RequestOptions, tc.DeliveryServicesRequiredCapability]{
//			"GET": {
//				"NOT MODIFIED when NO CHANGES made": {
//					ClientSession: TOSession,
//					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {tomorrow}}},
//					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
//				},
//				"OK when VALID request": {
//					ClientSession: TOSession,
//					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
//				},
//				"OK when VALID DELIVERYSERVICEID parameter": {
//					ClientSession: TOSession,
//					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"deliveryServiceId": {strconv.Itoa(GetDeliveryServiceId(t, "ds1")())}}},
//					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
//						validateDSRCExpectedFields(map[string]interface{}{"DeliveryServiceId": "ds1"})),
//				},
//				"OK when VALID XMLID parameter": {
//					ClientSession: TOSession,
//					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"xmlID": {"ds2"}}},
//					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
//						validateDSRCExpectedFields(map[string]interface{}{"XMLID": "ds2"})),
//				},
//				"OK when VALID REQUIREDCAPABILITY parameter": {
//					ClientSession: TOSession,
//					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"requiredCapability": {"bar"}}},
//					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
//						validateDSRCExpectedFields(map[string]interface{}{"RequiredCapability": "bar"})),
//				},
//				"FIRST RESULT when LIMIT=1": {
//					ClientSession: TOSession,
//					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"requiredCapability"}, "limit": {"1"}}},
//					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateDSRCPagination("limit")),
//				},
//				"SECOND RESULT when LIMIT=1 OFFSET=1": {
//					ClientSession: TOSession,
//					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"requiredCapability"}, "limit": {"1"}, "offset": {"1"}}},
//					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateDSRCPagination("offset")),
//				},
//				"SECOND RESULT when LIMIT=1 PAGE=2": {
//					ClientSession: TOSession,
//					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"requiredCapability"}, "limit": {"1"}, "page": {"2"}}},
//					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateDSRCPagination("page")),
//				},
//				"BAD REQUEST when INVALID LIMIT parameter": {
//					ClientSession: TOSession,
//					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"limit": {"-2"}}},
//					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
//				},
//				"BAD REQUEST when INVALID OFFSET parameter": {
//					ClientSession: TOSession,
//					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"limit": {"1"}, "offset": {"0"}}},
//					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
//				},
//				"BAD REQUEST when INVALID PAGE parameter": {
//					ClientSession: TOSession,
//					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"limit": {"1"}, "page": {"0"}}},
//					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
//				},
//				"OK when CHANGES made": {
//					ClientSession: TOSession,
//					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {currentTimeRFC}}},
//					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
//				},
//			},
//			"POST": {
//				"BAD REQUEST when REASSIGNING REQUIRED CAPABILITY to DELIVERY SERVICE": {
//					ClientSession: TOSession,
//					RequestBody: tc.DeliveryServicesRequiredCapability{
//						DeliveryServiceID:  util.IntPtr(GetDeliveryServiceId(t, "ds1")()),
//						RequiredCapability: util.StrPtr("foo"),
//					},
//					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
//				},
//				"BAD REQUEST when SERVERS DONT have CAPABILITY": {
//					ClientSession: TOSession,
//					RequestBody: tc.DeliveryServicesRequiredCapability{
//						DeliveryServiceID:  util.IntPtr(GetDeliveryServiceId(t, "test-ds-server-assignments")()),
//						RequiredCapability: util.StrPtr("disk"),
//					},
//					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
//				},
//				"BAD REQUEST when DELIVERY SERVICE HAS TOPOLOGY where SERVERS DONT have CAPABILITY": {
//					ClientSession: TOSession,
//					RequestBody: tc.DeliveryServicesRequiredCapability{
//						DeliveryServiceID:  util.IntPtr(GetDeliveryServiceId(t, "ds-top-req-cap")()),
//						RequiredCapability: util.StrPtr("bar"),
//					},
//					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
//				},
//				"BAD REQUEST when DELIVERY SERVICE ID EMPTY": {
//					ClientSession: TOSession,
//					RequestBody: tc.DeliveryServicesRequiredCapability{
//						RequiredCapability: util.StrPtr("bar"),
//					},
//					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
//				},
//				"BAD REQUEST when REQUIRED CAPABILITY EMPTY": {
//					ClientSession: TOSession,
//					RequestBody: tc.DeliveryServicesRequiredCapability{
//						DeliveryServiceID: util.IntPtr(GetDeliveryServiceId(t, "ds1")()),
//					},
//					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
//				},
//				"NOT FOUND when NON-EXISTENT REQUIRED CAPABILITY": {
//					ClientSession: TOSession,
//					RequestBody: tc.DeliveryServicesRequiredCapability{
//						DeliveryServiceID:  util.IntPtr(GetDeliveryServiceId(t, "ds1")()),
//						RequiredCapability: util.StrPtr("bogus"),
//					},
//					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
//				},
//				"NOT FOUND when NON-EXISTENT DELIVERY SERVICE ID": {
//					ClientSession: TOSession,
//					RequestBody: tc.DeliveryServicesRequiredCapability{
//						DeliveryServiceID:  util.IntPtr(-1),
//						RequiredCapability: util.StrPtr("foo"),
//					},
//					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
//				},
//				"BAD REQUEST when INVALID DELIVERY SERVICE TYPE": {
//					ClientSession: TOSession,
//					RequestBody: tc.DeliveryServicesRequiredCapability{
//						DeliveryServiceID:  util.IntPtr(GetDeliveryServiceId(t, "anymap-ds")()),
//						RequiredCapability: util.StrPtr("foo"),
//					},
//					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
//				},
//			},
//			"DELETE": {
//				"OK when VALID request": {
//					ClientSession: TOSession,
//					RequestBody: tc.DeliveryServicesRequiredCapability{
//						DeliveryServiceID:  util.IntPtr(GetDeliveryServiceId(t, "ds-top-req-cap")()),
//						RequiredCapability: util.StrPtr("ram"),
//					},
//					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
//				},
//				"NOT FOUND when NON-EXISTENT DELIVERYSERVICEID parameter": {
//					ClientSession: TOSession,
//					RequestBody: tc.DeliveryServicesRequiredCapability{
//						DeliveryServiceID:  util.IntPtr(-1),
//						RequiredCapability: util.StrPtr("foo"),
//					},
//					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
//				},
//				"NOT FOUND when NON-EXISTENT REQUIREDCAPABILITY parameter": {
//					ClientSession: TOSession,
//					RequestBody: tc.DeliveryServicesRequiredCapability{
//						DeliveryServiceID:  util.IntPtr(GetDeliveryServiceId(t, "ds1")()),
//						RequiredCapability: util.StrPtr("bogus"),
//					},
//					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
//				},
//			},
//		}
//
//		for method, testCases := range methodTests {
//			t.Run(method, func(t *testing.T) {
//				for name, testCase := range testCases {
//					switch method {
//					case "GET":
//						t.Run(name, func(t *testing.T) {
//							resp, reqInf, err := testCase.ClientSession.GetDeliveryServicesRequiredCapabilities(testCase.RequestOpts)
//							for _, check := range testCase.Expectations {
//								check(t, reqInf, resp.Response, resp.Alerts, err)
//							}
//						})
//					case "POST":
//						t.Run(name, func(t *testing.T) {
//							alerts, reqInf, err := testCase.ClientSession.CreateDeliveryServicesRequiredCapability(testCase.RequestBody, testCase.RequestOpts)
//							for _, check := range testCase.Expectations {
//								check(t, reqInf, nil, alerts, err)
//							}
//						})
//					case "DELETE":
//						t.Run(name, func(t *testing.T) {
//							alerts, reqInf, err := testCase.ClientSession.DeleteDeliveryServicesRequiredCapability(*testCase.RequestBody.DeliveryServiceID, *testCase.RequestBody.RequiredCapability, testCase.RequestOpts)
//							for _, check := range testCase.Expectations {
//								check(t, reqInf, nil, alerts, err)
//							}
//						})
//					}
//				}
//			})
//		}
//	})
//
//}

func validateDSRCExpectedFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		dsrcResp := resp.([]tc.DeliveryServicesRequiredCapability)
		for field, expected := range expectedResp {
			for _, dsrc := range dsrcResp {
				switch field {
				case "DeliveryServiceID":
					assert.Equal(t, expected, *dsrc.DeliveryServiceID, "Expected deliveryServiceId to be %v, but got %v", expected, dsrc.DeliveryServiceID)
				case "XMLID":
					assert.Equal(t, expected, *dsrc.XMLID, "Expected xmlID to be %v, but got %v", expected, dsrc.XMLID)
				case "RequiredCapability":
					assert.Equal(t, expected, *dsrc.RequiredCapability, "Expected requiredCapability to be %v, but got %v", expected, dsrc.RequiredCapability)
				}
			}
		}
	}
}

func validateDSRCPagination(paginationParam string) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		paginationResp := resp.([]tc.DeliveryServicesRequiredCapability)

		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("orderby", "requiredCapability")
		respBase, _, err := TOSession.GetDeliveryServicesRequiredCapabilities(opts)
		assert.RequireNoError(t, err, "Cannot get Delivery Services Required Capabilities: %v - alerts: %+v", err, respBase.Alerts)

		dsrc := respBase.Response
		assert.RequireGreaterOrEqual(t, len(dsrc), 3, "Need at least 3 Delivery Services Required Capabilities in Traffic Ops to test pagination support, found: %d", len(dsrc))
		switch paginationParam {
		case "limit:":
			assert.Exactly(t, dsrc[:1], paginationResp, "Expected GET deliveryservices_required_capabilities with limit = 1 to return first result")
		case "offset":
			assert.Exactly(t, dsrc[1:2], paginationResp, "Expected GET deliveryservices_required_capabilities with limit = 1, offset = 1 to return second result")
		case "page":
			assert.Exactly(t, dsrc[1:2], paginationResp, "Expected GET deliveryservices_required_capabilities with limit = 1, page = 2 to return second result")
		}
	}
}

func CreateTestDeliveryServicesRequiredCapabilities(t *testing.T) {
	// Assign all required capability to delivery services listed in `tc-fixtures.json`.
	for _, dsrc := range testData.DeliveryServicesRequiredCapabilities {
		dsId := GetDeliveryServiceId(t, *dsrc.XMLID)()
		dsrc = tc.DeliveryServicesRequiredCapability{
			DeliveryServiceID:  &dsId,
			RequiredCapability: dsrc.RequiredCapability,
		}
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("id", fmt.Sprint(dsId))
		resp, _, err := TOSession.GetDeliveryServices(opts)
		assert.NoError(t, err, "Error getting delivery services: %v - alerts: %v", err, resp.Alerts)
		assert.Equal(t, len(resp.Response), 1, "Expected response to have exactly 1 delivery service, but got %d", len(resp.Response))
		ds := resp.Response[0]
		ds.RequiredCapabilities = []string{*dsrc.RequiredCapability}
		_, _, err = TOSession.UpdateDeliveryService(dsId, ds, client.NewRequestOptions())
		assert.NoError(t, err, "Unexpected error creating a Delivery Service/Required Capability relationship: %v - alerts: %+v", err, resp.Alerts)
	}
}

func DeleteTestDeliveryServicesRequiredCapabilities(t *testing.T) {
	resp, _, err := TOSession.GetDeliveryServices(client.RequestOptions{})
	assert.NoError(t, err, "Error getting delivery services: %v - alerts: %v", err, resp.Alerts)
	for _, r := range resp.Response {
		r.RequiredCapabilities = []string{}
		response, _, err := TOSession.UpdateDeliveryService(*r.ID, r, client.RequestOptions{})
		assert.NoError(t, err, "Error removing Delivery Service/ Required Capability relationship: %v - alerts: %v", err, response.Alerts)
	}
}
