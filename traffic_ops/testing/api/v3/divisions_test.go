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
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

func TestDivisions(t *testing.T) {
	WithObjs(t, []TCObj{Parameters, Divisions, Regions}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.V3TestCaseT[tc.Division]{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfModifiedSince: {tomorrow}},
					Expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateDivisionSort()),
				},
				"OK when VALID NAME parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"name": {"division1"}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						validateDivisionFields(map[string]interface{}{"Name": "division1"})),
				},
				"OK when CHANGES made": {
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfModifiedSince: {currentTimeRFC}},
					Expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
			"PUT": {
				"OK when VALID request": {
					EndpointID:    GetDivisionID(t, "cdn-div2"),
					ClientSession: TOSession,
					RequestBody: tc.Division{
						Name: "testdivision",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateDivisionUpdateCreateFields("testdivision", map[string]interface{}{"Name": "testdivision"})),
				},
				"PRECONDITION FAILED when updating with IMS & IUS Headers": {
					EndpointID:     GetDivisionID(t, "division1"),
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}},
					RequestBody: tc.Division{
						Name: "division1",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"PRECONDITION FAILED when updating with IFMATCH ETAG Header": {
					EndpointID:    GetDivisionID(t, "division1"),
					ClientSession: TOSession,
					RequestBody: tc.Division{
						Name: "division1",
					},
					RequestHeaders: http.Header{rfc.IfMatch: {rfc.ETag(currentTime)}},
					Expectations:   utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
			},
			"DELETE": {
				"BAD REQUEST when DIVISION in use by REGION": {
					EndpointID:    GetDivisionID(t, "division1"),
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
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
								resp, reqInf, err := testCase.ClientSession.GetDivisionByNameWithHdr(testCase.RequestParams["name"][0], testCase.RequestHeaders)
								for _, check := range testCase.Expectations {
									check(t, reqInf, resp, tc.Alerts{}, err)
								}
							} else {
								resp, reqInf, err := testCase.ClientSession.GetDivisionsWithHdr(testCase.RequestHeaders)
								for _, check := range testCase.Expectations {
									check(t, reqInf, resp, tc.Alerts{}, err)
								}
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.CreateDivision(testCase.RequestBody)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.UpdateDivisionByIDWithHdr(testCase.EndpointID(), testCase.RequestBody, testCase.RequestHeaders)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.DeleteDivisionByID(testCase.EndpointID())
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

func validateDivisionFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Division response to not be nil.")
		divisionResp := resp.([]tc.Division)
		for field, expected := range expectedResp {
			for _, division := range divisionResp {
				switch field {
				case "Name":
					assert.Equal(t, expected, division.Name, "Expected Name to be %v, but got %s", expected, division.Name)
				default:
					t.Errorf("Expected field: %v, does not exist in response", field)
				}
			}
		}
	}
}

func validateDivisionUpdateCreateFields(name string, expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		divisions, _, err := TOSession.GetDivisionByNameWithHdr(name, nil)
		assert.RequireNoError(t, err, "Error getting Division: %v - alerts: %+v", err, divisions)
		assert.RequireEqual(t, 1, len(divisions), "Expected one Division returned Got: %d", len(divisions))
		validateDivisionFields(expectedResp)(t, toclientlib.ReqInf{}, divisions, tc.Alerts{}, nil)
	}
}

func validateDivisionSort() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Division response to not be nil.")
		var divisionNames []string
		divisionResp := resp.([]tc.Division)
		for _, division := range divisionResp {
			divisionNames = append(divisionNames, division.Name)
		}
		assert.Equal(t, true, sort.StringsAreSorted(divisionNames), "List is not sorted by their names: %v", divisionNames)
	}
}

func GetDivisionID(t *testing.T, divisionName string) func() int {
	return func() int {
		divisionsResp, _, err := TOSession.GetDivisionByNameWithHdr(divisionName, nil)
		assert.RequireNoError(t, err, "Get Divisions Request failed with error:", err)
		assert.RequireEqual(t, 1, len(divisionsResp), "Expected response object length 1, but got %d", len(divisionsResp))
		return divisionsResp[0].ID
	}
}

func CreateTestDivisions(t *testing.T) {
	for _, division := range testData.Divisions {
		alerts, _, err := TOSession.CreateDivision(division)
		assert.RequireNoError(t, err, "Could not create Division '%s': %v - alerts: %+v", division.Name, err, alerts)
	}
}

func DeleteTestDivisions(t *testing.T) {
	divisions, _, err := TOSession.GetDivisionsWithHdr(nil)
	assert.NoError(t, err, "Cannot get Divisions: %v", err)
	for _, division := range divisions {
		alerts, _, err := TOSession.DeleteDivisionByID(division.ID)
		assert.NoError(t, err, "Unexpected error deleting Division '%s' (#%d): %v - alerts: %+v", division.Name, division.ID, err, alerts.Alerts)
		// Retrieve the Division to see if it got deleted
		getDivision, _, err := TOSession.GetDivisionByIDWithHdr(division.ID, nil)
		assert.NoError(t, err, "Error getting Division '%s' after deletion: %v", division.Name, err)
		assert.Equal(t, 0, len(getDivision), "Expected Division '%s' to be deleted, but it was found in Traffic Ops", division.Name)
	}
}
