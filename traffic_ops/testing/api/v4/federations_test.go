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
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/assert"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

func TestFederations(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, DeliveryServices, CDNFederations, FederationDeliveryServices, FederationUsers}, func() {

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
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1)),
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
						validateFederationsUpdateFields(map[string]interface{}{
							"federations": []map[string]interface{}{
								{
									"DeliveryService": "ds1",
									"Mappings": map[string]interface{}{
										"resolve4": []string{"192.0.2.0/25", "192.0.2.128/25"},
										"resolve6": []string{"2001:db8::/33", "2001:db8:8000::/33"},
									},
								},
								{
									"DeliveryService": "ds2",
									"Mappings": map[string]interface{}{
										"resolve4": []string{"192.0.2.0/25", "192.0.2.128/25"},
										"resolve6": []string{"2001:db8::/33", "2001:db8:8000::/33"},
									},
								},
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
			"DELETE": {},
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
							resp, reqInf, err := testCase.ClientSession.AllFederations(testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.AddFederationResolverMappingsForCurrentUser(federation, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.ReplaceFederationResolverMappingsForCurrentUser(federation, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.DeleteFederationResolverMappingsForCurrentUser(testCase.RequestOpts)
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

func validateFederationFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Federation response to not be nil.")
		federationResp := resp.([]tc.AllDeliveryServiceFederationsMapping)
		for _, federation := range federationResp {
			if federation.DeliveryService.String() == expectedResp["DeliveryService"] {
				assert.RequireEqual(t, 1, len(federation.Mappings), "expected 1 mapping, got %d", len(federation.Mappings))
				//sort.Strings(federation.Mappings[0].Resolve4)
				//sort.Strings(federation.Mappings[0].Resolve6)
				assert.Exactly(t, expectedResp["Resolve4"], federation.Mappings[0].Resolve4, "checking federation resolver mappings, expected: %+v, actual: %+v", expectedResp["Resolve4"], federation.Mappings[0].Resolve4)
				assert.Exactly(t, expectedResp["Resolve6"], federation.Mappings[0].Resolve6, "checking federation resolver mappings, expected: %+v, actual: %+v", expectedResp["Resolve6"], federation.Mappings[0].Resolve6)
			}
		}
	}
}

func validateFederationsUpdateFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		opts := client.NewRequestOptions()
		federation, _, err := TOSession.Federations(opts)
		assert.RequireNoError(t, err, "Error getting Federations: %v - alerts: %+v", err, federation.Alerts)
		assert.RequireEqual(t, 1, len(federation.Response), "Expected one Federation returned Got: %d", len(federation.Response))
		for _, federationRequest := range expectedResp {
			dat, err := json.Marshal(federationRequest)
			assert.RequireNoError(t, err, "Error occurred when marshalling request body: %v", err)
			err = json.Unmarshal(dat, &federation)
			assert.RequireNoError(t, err, "Error occurred when unmarshalling request body: %v", err)
		}
		validateFederationFields(expectedResp)(t, toclientlib.ReqInf{}, federation.Response, tc.Alerts{}, nil)
	}
}

//
//func GetTestFederations(t *testing.T) {
//	if len(testData.Federations) == 0 {
//		t.Error("no federations test data")
//	}
//
//	feds, _, err := TOSession.AllFederations(client.RequestOptions{})
//	if err != nil {
//		t.Errorf("getting federations: %v - alerts: %+v", err, feds.Alerts)
//	}
//
//	if len(feds.Response) < 1 {
//		t.Errorf("expected atleast 1 federation, but got none")
//	}
//	fed := feds.Response[0]
//
//	if len(fed.Mappings) < 1 {
//		t.Fatal("federation mappings expected > 1, actual: 0")
//	}
//
//	mapping := fed.Mappings[0]
//	if mapping.CName == nil {
//		t.Fatal("federation mapping expected cname, actual: nil")
//	}
//	if mapping.TTL == nil {
//		t.Fatal("federation mapping expected ttl, actual: nil")
//	}
//
//	matched := false
//	for _, testFed := range testData.Federations {
//		if testFed.CName == nil {
//			t.Error("test federation missing cname!")
//			continue
//		}
//		if testFed.TTL == nil {
//			t.Error("test federation missing ttl!")
//			continue
//		}
//
//		if *mapping.CName != *testFed.CName {
//			continue
//		}
//		matched = true
//
//		if *mapping.TTL != *testFed.TTL {
//			t.Errorf("federation mapping ttl expected: %v, actual: %v", *testFed.TTL, *mapping.TTL)
//		}
//	}
//	if !matched {
//		t.Errorf("federation mapping expected to match test data, actual: cname %v not in test data", *mapping.CName)
//	}
//}
//
//func RemoveFederationResolversForCurrentUserTest(t *testing.T) {
//	if len(testData.Federations) < 1 {
//		t.Fatal("No test Federations, deleting resolvers cannot be tested!")
//	}
//
//	alerts, _, err := TOSession.DeleteFederationResolverMappingsForCurrentUser(client.RequestOptions{})
//	if err != nil {
//		t.Fatalf("Unexpected error deleting Federation Resolvers for current user: %v - alerts: %+v", err, alerts)
//	}
//	for _, a := range alerts.Alerts {
//		if a.Level == tc.ErrorLevel.String() {
//			t.Errorf("Unexpected error-level alert from deleting Federation Resolvers for current user: %s", a.Text)
//		}
//	}
//
//	// Now try deleting Federation Resolvers when there are none.
//	_, _, err = TOSession.DeleteFederationResolverMappingsForCurrentUser(client.RequestOptions{})
//	if err != nil {
//		t.Logf("Received expected error deleting Federation Resolvers for current user: %v", err)
//	} else {
//		t.Error("Expected an error deleting zero Federation Resolvers, but didn't get one.")
//	}
//}
