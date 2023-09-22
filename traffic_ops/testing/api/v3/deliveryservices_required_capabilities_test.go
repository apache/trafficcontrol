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
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

func TestDeliveryServicesRequiredCapabilities(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, ServerCapabilities, Topologies, ServiceCategories, DeliveryServices, DeliveryServiceServerAssignments, ServerServerCapabilities, DeliveryServicesRequiredCapabilities}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.V3TestCaseT[tc.DeliveryServicesRequiredCapability]{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfModifiedSince: {tomorrow}},
					Expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when VALID DELIVERYSERVICEID parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"deliveryServiceId": {"ds1"}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateDSRCExpectedFields(map[string]interface{}{"DeliveryServiceId": "ds1"})),
				},
				"OK when VALID XMLID parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"xmlID": {"ds2"}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateDSRCExpectedFields(map[string]interface{}{"XMLID": "ds2"})),
				},
				"OK when VALID REQUIREDCAPABILITY parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"requiredCapability": {"bar"}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateDSRCExpectedFields(map[string]interface{}{"RequiredCapability": "bar"})),
				},
				"OK when CHANGES made": {
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfModifiedSince: {currentTimeRFC}},
					Expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
			"POST": {
				"BAD REQUEST when REASSIGNING REQUIRED CAPABILITY to DELIVERY SERVICE": {
					ClientSession: TOSession,
					RequestBody: tc.DeliveryServicesRequiredCapability{
						DeliveryServiceID:  util.Ptr(GetDeliveryServiceId(t, "ds1")()),
						RequiredCapability: util.Ptr("foo"),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when SERVERS DONT have CAPABILITY": {
					ClientSession: TOSession,
					RequestBody: tc.DeliveryServicesRequiredCapability{
						DeliveryServiceID:  util.Ptr(GetDeliveryServiceId(t, "test-ds-server-assignments")()),
						RequiredCapability: util.Ptr("disk"),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when DELIVERY SERVICE HAS TOPOLOGY where SERVERS DONT have CAPABILITY": {
					ClientSession: TOSession,
					RequestBody: tc.DeliveryServicesRequiredCapability{
						DeliveryServiceID:  util.Ptr(GetDeliveryServiceId(t, "ds-top-req-cap")()),
						RequiredCapability: util.Ptr("bar"),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when DELIVERY SERVICE ID EMPTY": {
					ClientSession: TOSession,
					RequestBody: tc.DeliveryServicesRequiredCapability{
						RequiredCapability: util.Ptr("bar"),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when REQUIRED CAPABILITY EMPTY": {
					ClientSession: TOSession,
					RequestBody: tc.DeliveryServicesRequiredCapability{
						DeliveryServiceID: util.Ptr(GetDeliveryServiceId(t, "ds1")()),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"NOT FOUND when NON-EXISTENT REQUIRED CAPABILITY": {
					ClientSession: TOSession,
					RequestBody: tc.DeliveryServicesRequiredCapability{
						DeliveryServiceID:  util.Ptr(GetDeliveryServiceId(t, "ds1")()),
						RequiredCapability: util.Ptr("bogus"),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"NOT FOUND when NON-EXISTENT DELIVERY SERVICE ID": {
					ClientSession: TOSession,
					RequestBody: tc.DeliveryServicesRequiredCapability{
						DeliveryServiceID:  util.Ptr(-1),
						RequiredCapability: util.Ptr("foo"),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"BAD REQUEST when INVALID DELIVERY SERVICE TYPE": {
					ClientSession: TOSession,
					RequestBody: tc.DeliveryServicesRequiredCapability{
						DeliveryServiceID:  util.Ptr(GetDeliveryServiceId(t, "anymap-ds")()),
						RequiredCapability: util.Ptr("foo"),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
			"DELETE": {
				"OK when VALID request": {
					EndpointID:    GetDeliveryServiceId(t, "ds-top-req-cap"),
					ClientSession: TOSession,
					RequestBody: tc.DeliveryServicesRequiredCapability{
						RequiredCapability: util.Ptr("ram"),
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"NOT FOUND when NON-EXISTENT DELIVERYSERVICEID parameter": {
					EndpointID:    func() int { return -1 },
					ClientSession: TOSession,
					RequestBody: tc.DeliveryServicesRequiredCapability{
						RequiredCapability: util.Ptr("foo"),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"NOT FOUND when NON-EXISTENT REQUIREDCAPABILITY parameter": {
					EndpointID:    GetDeliveryServiceId(t, "ds1"),
					ClientSession: TOSession,
					RequestBody: tc.DeliveryServicesRequiredCapability{
						RequiredCapability: util.Ptr("bogus"),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					var deliveryServiceId *int
					var xmlId *string
					var capability *string

					if val, ok := testCase.RequestParams["deliveryServiceId"]; ok {
						if _, err := strconv.Atoi(val[0]); err != nil {
							dsId := GetDeliveryServiceId(t, val[0])()
							deliveryServiceId = &dsId
						}
					}
					if val, ok := testCase.RequestParams["xmlID"]; ok {
						xmlId = &val[0]
					}
					if val, ok := testCase.RequestParams["requiredCapability"]; ok {
						capability = &val[0]
					}

					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetDeliveryServicesRequiredCapabilitiesWithHdr(deliveryServiceId, xmlId, capability, testCase.RequestHeaders)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp, tc.Alerts{}, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.CreateDeliveryServicesRequiredCapability(testCase.RequestBody)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.DeleteDeliveryServicesRequiredCapability(testCase.EndpointID(), *testCase.RequestBody.RequiredCapability)
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

func CreateTestDeliveryServicesRequiredCapabilities(t *testing.T) {
	// Assign all required capability to delivery services listed in `tc-fixtures.json`.
	for _, dsrc := range testData.DeliveryServicesRequiredCapabilities {
		dsId := GetDeliveryServiceId(t, *dsrc.XMLID)()
		dsrc = tc.DeliveryServicesRequiredCapability{
			DeliveryServiceID:  &dsId,
			RequiredCapability: dsrc.RequiredCapability,
		}
		resp, _, err := TOSession.CreateDeliveryServicesRequiredCapability(dsrc)
		assert.NoError(t, err, "Unexpected error creating a Delivery Service/Required Capability relationship: %v - alerts: %+v", err, resp.Alerts)
	}
}

func DeleteTestDeliveryServicesRequiredCapabilities(t *testing.T) {
	// Get Required Capabilities to delete them
	dsrcs, _, err := TOSession.GetDeliveryServicesRequiredCapabilitiesWithHdr(nil, nil, nil, nil)
	assert.NoError(t, err, "Error getting Delivery Service/Required Capability relationships: %v - resp: %+v", err, dsrcs)

	for _, dsrc := range dsrcs {
		alerts, _, err := TOSession.DeleteDeliveryServicesRequiredCapability(*dsrc.DeliveryServiceID, *dsrc.RequiredCapability)
		assert.NoError(t, err, "Error deleting a relationship between a Delivery Service and a Capability: %v - alerts: %+v", err, alerts.Alerts)
	}
}
