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
	"sort"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	tc "github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

func TestCDNs(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Parameters, Tenants, Users}, func() {

		readOnlyUserSession := utils.CreateV3Session(t, Config.TrafficOps.URL, "readonlyuser", "pa$$word", Config.Default.Session.TimeoutInSecs)

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.V3TestCaseT[tc.CDN]{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfModifiedSince: {tomorrow}},
					Expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1), validateCDNSort()),
				},
				"OK when VALID NAME parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"name": {"cdn1"}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						validateCDNFields(map[string]interface{}{"Name": "cdn1"})),
				},
				"OK when CHANGES made": {
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfModifiedSince: {currentTimeRFC}},
					Expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
			"POST": {
				"FORBIDDEN when READ ONLY USER": {
					ClientSession: readOnlyUserSession,
					RequestBody: tc.CDN{
						Name:          "readOnlyTest",
						DNSSECEnabled: false,
						DomainName:    "test.ro",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
			"PUT": {
				"OK when VALID request": {
					EndpointID:    GetCDNID(t, "cdn1"),
					ClientSession: TOSession,
					RequestBody: tc.CDN{
						DNSSECEnabled: false,
						DomainName:    "domain2",
						Name:          "cdn1",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateCDNUpdateFields("cdn1", map[string]interface{}{"DomainName": "domain2"})),
				},
				"PRECONDITION FAILED when updating with IF-UNMODIFIED-SINCE Headers": {
					EndpointID:     GetCDNID(t, "cdn1"),
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}},
					RequestBody: tc.CDN{
						DNSSECEnabled: false,
						DomainName:    "newDomain",
						Name:          "cdn1",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"PRECONDITION FAILED when updating with IFMATCH ETAG Header": {
					EndpointID:     GetCDNID(t, "cdn1"),
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfMatch: {rfc.ETag(currentTime)}},
					RequestBody: tc.CDN{
						DNSSECEnabled: false,
						DomainName:    "newDomain",
						Name:          "cdn1",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
			},
		}
		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							if name == "OK when VALID NAME parameter" {
								resp, reqInf, err := testCase.ClientSession.GetCDNByNameWithHdr(testCase.RequestParams["name"][0], testCase.RequestHeaders)
								for _, check := range testCase.Expectations {
									check(t, reqInf, resp, tc.Alerts{}, err)
								}
							} else {
								resp, reqInf, err := testCase.ClientSession.GetCDNsWithHdr(testCase.RequestHeaders)
								for _, check := range testCase.Expectations {
									check(t, reqInf, resp, tc.Alerts{}, err)
								}
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.CreateCDN(testCase.RequestBody)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.UpdateCDNByIDWithHdr(testCase.EndpointID(), testCase.RequestBody, testCase.RequestHeaders)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.DeleteCDNByID(testCase.EndpointID())
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

func validateCDNFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		cdnResp := resp.([]tc.CDN)
		for field, expected := range expectedResp {
			for _, cdn := range cdnResp {
				switch field {
				case "Name":
					assert.Equal(t, expected, cdn.Name, "Expected Name to be %v, but got %v", expected, cdn.Name)
				case "DomainName":
					assert.Equal(t, expected, cdn.DomainName, "Expected DomainName to be %v, but got %v", expected, cdn.DomainName)
				case "DNSSECEnabled":
					assert.Equal(t, expected, cdn.DNSSECEnabled, "Expected DNSSECEnabled to be %v, but got %v", expected, cdn.DNSSECEnabled)
				default:
					t.Errorf("Expected field: %v, does not exist in response", field)
				}
			}
		}
	}
}

func validateCDNUpdateFields(name string, expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		cdn, _, err := TOSession.GetCDNByNameWithHdr(name, nil)
		assert.NoError(t, err, "Error getting CDN: %v", err)
		assert.Equal(t, 1, len(cdn), "Expected one CDN returned Got: %d", len(cdn))
		validateCDNFields(expectedResp)(t, toclientlib.ReqInf{}, cdn, tc.Alerts{}, nil)
	}
}

func validateCDNSort() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected CDN response to not be nil.")
		var cdnNames []string
		cdnResp := resp.([]tc.CDN)
		for _, cdn := range cdnResp {
			cdnNames = append(cdnNames, cdn.Name)
		}
		assert.Equal(t, true, sort.StringsAreSorted(cdnNames), "List is not sorted by their names: %v", cdnNames)
	}
}

func GetCDNID(t *testing.T, cdnName string) func() int {
	return func() int {
		cdnsResp, _, err := TOSession.GetCDNByNameWithHdr(cdnName, http.Header{})
		assert.RequireNoError(t, err, "Get CDNs Request failed with error:", err)
		assert.RequireEqual(t, 1, len(cdnsResp), "Expected response object length 1, but got %d", len(cdnsResp))
		assert.RequireNotNil(t, cdnsResp[0].ID, "Expected id to not be nil")
		return cdnsResp[0].ID
	}
}

func CreateTestCDNs(t *testing.T) {
	for _, cdn := range testData.CDNs {
		resp, _, err := TOSession.CreateCDN(cdn)
		assert.NoError(t, err, "Could not create CDN: %v - alerts: %+v", err, resp.Alerts)
	}
}

func DeleteTestCDNs(t *testing.T) {
	resp, _, err := TOSession.GetCDNsWithHdr(http.Header{})
	assert.NoError(t, err, "Cannot get CDNs: %v", err)
	for _, cdn := range resp {
		delResp, _, err := TOSession.DeleteCDNByID(cdn.ID)
		assert.NoError(t, err, "Cannot delete CDN '%s' (#%d): %v - alerts: %+v", cdn.Name, cdn.ID, err, delResp.Alerts)

		// Retrieve the CDN to see if it got deleted
		cdns, _, err := TOSession.GetCDNByIDWithHdr(cdn.ID, http.Header{})
		assert.NoError(t, err, "Error deleting CDN '%s': %v", cdn.Name, err)
		assert.Equal(t, 0, len(cdns), "Expected CDN '%s' to be deleted", cdn.Name)
	}
}
