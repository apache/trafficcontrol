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

func TestFederationResolvers(t *testing.T) {
	WithObjs(t, []TCObj{Types, FederationResolvers}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.V3TestCaseT[tc.FederationResolver]{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfModifiedSince: {tomorrow}},
					Expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1)),
				},
				"OK when VALID ID parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"id": {strconv.Itoa(GetFederationResolverID(t, "0.0.0.0/12")())}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateFederationResolverFields(map[string]interface{}{"ID": uint(GetFederationResolverID(t, "0.0.0.0/12")())})),
				},
				"OK when VALID IPADDRESS parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"ipAddress": {"1.2.3.4"}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateFederationResolverFields(map[string]interface{}{"IPAddress": "1.2.3.4"})),
				},
				"OK when VALID TYPE parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"type": {"RESOLVE4"}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateFederationResolversType(map[string]interface{}{"Type": "RESOLVE4"})),
				},
			},
			"POST": {
				"BAD REQUEST when MISSING IPADDRESS and TYPE FIELDS": {
					ClientSession: TOSession,
					RequestBody:   tc.FederationResolver{},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID IP ADDRESS": {
					ClientSession: TOSession,
					RequestBody: tc.FederationResolver{
						IPAddress: util.Ptr("not a valid IP address"),
						TypeID:    util.Ptr(uint(GetTypeId(t, "RESOLVE4"))),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
			"DELETE": {
				"NOT FOUND when INVALID ID": {
					EndpointID:    func() int { return 0 },
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							if name == "OK when VALID ID parameter" {
								id, err := strconv.Atoi(testCase.RequestParams["id"][0])
								assert.RequireNoError(t, err, "Error converting string to int")
								resp, reqInf, err := testCase.ClientSession.GetFederationResolverByIDWithHdr(uint(id), testCase.RequestHeaders)
								for _, check := range testCase.Expectations {
									check(t, reqInf, resp, tc.Alerts{}, err)
								}
							} else if name == "OK when VALID IPADDRESS parameter" {
								resp, reqInf, err := testCase.ClientSession.GetFederationResolverByIPAddressWithHdr(testCase.RequestParams["ipAddress"][0], testCase.RequestHeaders)
								for _, check := range testCase.Expectations {
									check(t, reqInf, resp, tc.Alerts{}, err)
								}
							} else if name == "OK when VALID TYPE parameter" {
								resp, reqInf, err := testCase.ClientSession.GetFederationResolversByTypeWithHdr(testCase.RequestParams["type"][0], testCase.RequestHeaders)
								for _, check := range testCase.Expectations {
									check(t, reqInf, resp, tc.Alerts{}, err)
								}
							} else {
								resp, reqInf, err := testCase.ClientSession.GetFederationResolversWithHdr(testCase.RequestHeaders)
								for _, check := range testCase.Expectations {
									check(t, reqInf, resp, tc.Alerts{}, err)
								}
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.CreateFederationResolver(testCase.RequestBody)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.DeleteFederationResolver(uint(testCase.EndpointID()))
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

func validateFederationResolverFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Federation Resolver response to not be nil.")
		fr := resp.(tc.FederationResolver)
		for field, expected := range expectedResp {
			switch field {
			case "ID":
				assert.RequireNotNil(t, fr.ID, "Expected ID to not be nil")
				assert.Equal(t, expected, *fr.ID, "Expected ID to be %v, but got %d", expected, *fr.ID)
			case "IPAddress":
				assert.RequireNotNil(t, fr.IPAddress, "Expected IPAddress to not be nil")
				assert.Equal(t, expected, *fr.IPAddress, "Expected IPAddress to be %v, but got %s", expected, *fr.IPAddress)
			default:
				t.Errorf("Expected field: %v, does not exist in response", field)
			}
		}
	}
}

func validateFederationResolversType(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Federation Resolver response to not be nil.")
		frResp := resp.([]tc.FederationResolver)
		for field, expected := range expectedResp {
			for _, fr := range frResp {
				switch field {
				case "Type":
					assert.RequireNotNil(t, fr.Type, "Expected Type to not be nil")
					assert.Equal(t, expected, *fr.Type, "Expected Type to be %v, but got %s", expected, *fr.Type)
				default:
					t.Errorf("Expected field: %v, does not exist in response", field)
				}
			}
		}
	}
}

func GetFederationResolverID(t *testing.T, ipAddress string) func() int {
	return func() int {
		federationResolver, _, err := TOSession.GetFederationResolverByIPAddressWithHdr(ipAddress, nil)
		assert.RequireNoError(t, err, "Get FederationResolvers Request failed with error:", err)
		assert.RequireNotNil(t, federationResolver.ID, "Expected Federation Resolver ID to not be nil")
		return int(*federationResolver.ID)
	}
}

func CreateTestFederationResolvers(t *testing.T) {
	for _, fr := range testData.FederationResolvers {
		fr.TypeID = util.UIntPtr(uint(GetTypeId(t, *fr.Type)))
		alerts, _, err := TOSession.CreateFederationResolver(fr)
		assert.RequireNoError(t, err, "Failed to create Federation Resolver %+v: %v - alerts: %+v", fr, err, alerts.Alerts)
	}
}

func DeleteTestFederationResolvers(t *testing.T) {
	frs, _, err := TOSession.GetFederationResolversWithHdr(nil)
	assert.RequireNoError(t, err, "Unexpected error getting Federation Resolvers: %v", err)
	for _, fr := range frs {
		alerts, _, err := TOSession.DeleteFederationResolver(*fr.ID)
		assert.NoError(t, err, "Failed to delete Federation Resolver %+v: %v - alerts: %+v", fr, err, alerts.Alerts)
		// Retrieve the Federation Resolver to see if it got deleted
		getFR, _, err := TOSession.GetFederationResolverByIDWithHdr(*fr.ID, nil)
		assert.NoError(t, err, "Error getting Federation Resolver '%d' after deletion: %v", *fr.ID, err)
		assert.Equal(t, (*uint)(nil), getFR.ID, "Expected Federation Resolver '%d' to be deleted, but it was found in Traffic Ops", *fr.ID)
	}
}
