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
	"sort"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/assert"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

func TestFederations(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, DeliveryServices, CDNFederations, FederationDeliveryServices, FederationUsers}, func() {

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
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateAllFederationsFields([]map[string]interface{}{
							{
								"DeliveryService": "ds1",
								"Mappings": []map[string]interface{}{
									{
										"CName": "the.cname.com.",
										"TTL":   48,
									},
									{
										"CName": "booya.com.",
										"TTL":   34,
									},
									{
										"CName": "google.com.",
										"TTL":   30,
									},
								},
							},
						})),
				},
			},
			"POST": {
				"OK when VALID request": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"federations": []map[string]interface{}{
							{
								"deliveryService": "ds1",
								"mappings": map[string]interface{}{
									"resolve4": []string{"0.0.0.0"},
									"resolve6": []string{"::1"},
								},
							},
							{
								"deliveryService": "ds2",
								"mappings": map[string]interface{}{
									"resolve4": []string{"1.2.3.4/28"},
									"resolve6": []string{"1234::/110"},
								},
							},
						},
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"CONFLICT when INVALID DELIVERY SERVICE": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"federations": []map[string]interface{}{
							{
								"deliveryService": "aoeuhtns",
								"mappings": map[string]interface{}{
									"resolve4": []string{"1.2.3.4/28"},
									"resolve6": []string{"dead::beef", "f1d0::f00d/82"},
								},
							},
						},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusConflict)),
				},
			},
			"PUT": {
				"OK when VALID request": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"federations": []map[string]interface{}{
							{
								"deliveryService": "ds1",
								"mappings": map[string]interface{}{
									"resolve4": []string{"192.0.2.0/25", "192.0.2.128/25"},
									"resolve6": []string{"2001:db8::/33", "2001:db8:8000::/33"},
								},
							},
							{
								"deliveryService": "ds2",
								"mappings": map[string]interface{}{
									"resolve4": []string{"192.0.2.0/25", "192.0.2.128/25"},
									"resolve6": []string{"2001:db8::/33", "2001:db8:8000::/33"},
								},
							},
						},
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateFederationsUpdateFields([]map[string]interface{}{
							{
								"DeliveryService": "ds1",
								"Resolve4":        []string{"192.0.2.0/25", "192.0.2.128/25"},
								"Resolve6":        []string{"2001:db8::/33", "2001:db8:8000::/33"},
							},
							{
								"DeliveryService": "ds2",
								"Resolve4":        []string{"192.0.2.0/25", "192.0.2.128/25"},
								"Resolve6":        []string{"2001:db8::/33", "2001:db8:8000::/33"},
							},
						})),
				},
				"CONFLICT when INVALID DELIVERY SERVICE": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"federations": []map[string]interface{}{
							{
								"deliveryService": "aoeuhtns",
								"mappings": map[string]interface{}{
									"resolve4": []string{"1.2.3.4/28"},
									"resolve6": []string{"dead::beef", "f1d0::f00d/82"},
								},
							},
						},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusConflict)),
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
					federation := tc.DeliveryServiceFederationResolverMappingRequest{}

					if testCase.RequestBody != nil {
						for _, federationRequest := range testCase.RequestBody {
							dat, err := json.Marshal(federationRequest)
							assert.RequireNoError(t, err, "Error occurred when marshalling request body: %v", err)
							err = json.Unmarshal(dat, &federation)
							assert.RequireNoError(t, err, "Error occurred when unmarshalling request body: %v", err)
						}
					}

					switch method {
					case "GET", "GET AFTER CHANGES":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.AllFederationsWithHdr(testCase.RequestHeaders)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp, tc.Alerts{}, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.AddFederationResolverMappingsForCurrentUser(federation)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.ReplaceFederationResolverMappingsForCurrentUser(federation)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					}
				}
			})
		}

		t.Run("DELETE/OK when VALID request", func(t *testing.T) {
			alerts, reqInf, err := TOSession.DeleteFederationResolverMappingsForCurrentUser()
			for _, check := range utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)) {
				check(t, reqInf, nil, alerts, err)
			}
		})

		t.Run("DELETE/CONFLICT when NO FEDERATION RESOLVERS", func(t *testing.T) {
			alerts, reqInf, err := TOSession.DeleteFederationResolverMappingsForCurrentUser()
			for _, check := range utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusConflict)) {
				check(t, reqInf, nil, alerts, err)
			}
		})
	})
}

func validateFederationFields(expectedResp []map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Federation response to not be nil.")
		federationResp := resp.([]tc.AllDeliveryServiceFederationsMapping)
		for _, expectedFed := range expectedResp {
			for _, federation := range federationResp {
				if federation.DeliveryService.String() == expectedFed["DeliveryService"] {
					assert.RequireEqual(t, 1, len(federation.Mappings), "expected 1 mapping, got %d", len(federation.Mappings))
					sort.Strings(federation.Mappings[0].Resolve4)
					sort.Strings(federation.Mappings[0].Resolve6)
					sort.Strings(expectedFed["Resolve4"].([]string))
					sort.Strings(expectedFed["Resolve6"].([]string))
					assert.Exactly(t, expectedFed["Resolve4"], federation.Mappings[0].Resolve4, "checking federation resolver mappings, expected: %+v, actual: %+v", expectedFed["Resolve4"], federation.Mappings[0].Resolve4)
					assert.Exactly(t, expectedFed["Resolve6"], federation.Mappings[0].Resolve6, "checking federation resolver mappings, expected: %+v, actual: %+v", expectedFed["Resolve6"], federation.Mappings[0].Resolve6)
				}
			}
		}
	}
}

func validateFederationsUpdateFields(expectedResp []map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		federation, _, err := TOSession.FederationsWithHdr(nil)
		assert.RequireNoError(t, err, "Error getting Federations: %v", err)
		assert.RequireGreaterOrEqual(t, len(federation), 1, "Expected one Federation returned Got: %d", len(federation))
		validateFederationFields(expectedResp)(t, toclientlib.ReqInf{}, federation, tc.Alerts{}, nil)
	}
}

func validateAllFederationsFields(expectedResp []map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected All Federations response to not be nil.")
		allFederationResp := resp.([]tc.AllDeliveryServiceFederationsMapping)
		for _, expected := range expectedResp {
			for _, allFed := range allFederationResp {
				if allFed.DeliveryService.String() == expected["DeliveryService"] {
					for _, mapping := range allFed.Mappings {
						for _, expectedMapping := range expected["Mappings"].([]map[string]interface{}) {
							assert.RequireNotNil(t, mapping.CName, "Expected CName to not be nil.")
							if expectedMapping["CName"] == *mapping.CName {
								assert.RequireNotNil(t, mapping.TTL, "Expected TTL to not be nil.")
								assert.Equal(t, expectedMapping["TTL"], *mapping.TTL, "Expected TTL to be %v, but got %s", expected, allFed.Mappings[0].TTL)
							}
						}
					}
				}
			}
		}
	}
}
