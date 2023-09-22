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
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"
)

func TestDivisions(t *testing.T) {
	WithObjs(t, []TCObj{Parameters, Divisions, Regions}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.TestCase[client.Session, client.RequestOptions, tc.Division]{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {tomorrow}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateDivisionSort()),
				},
				"OK when VALID NAME parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"division1"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						validateDivisionFields(map[string]interface{}{"Name": "division1"})),
				},
				"EMPTY RESPONSE when INVALID ID parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"id": {"10000"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when INVALID NAME parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"abcd"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"VALID when SORTORDER param is DESC": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"sortOrder": {"desc"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateDivisionDescSort()),
				},
				"FIRST RESULT when LIMIT=1": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateDivisionPagination("limit")),
				},
				"SECOND RESULT when LIMIT=1 OFFSET=1": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}, "offset": {"1"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateDivisionPagination("offset")),
				},
				"SECOND RESULT when LIMIT=1 PAGE=2": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}, "page": {"2"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateDivisionPagination("page")),
				},
				"BAD REQUEST when INVALID LIMIT parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"limit": {"-2"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID OFFSET parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"limit": {"1"}, "offset": {"0"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID PAGE parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"limit": {"1"}, "page": {"0"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"OK when CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {currentTimeRFC}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
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
					EndpointID:    GetDivisionID(t, "division1"),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}}},
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
					RequestOpts:  client.RequestOptions{Header: http.Header{rfc.IfMatch: {rfc.ETag(currentTime)}}},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
			},
			"DELETE": {
				"BAD REQUEST when DIVISION in use by REGION": {
					EndpointID:    GetDivisionID(t, "division1"),
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"NOT FOUND when INVALID ID parameter": {
					EndpointID:    func() int { return 111111 },
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
							resp, reqInf, err := testCase.ClientSession.GetDivisions(testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.CreateDivision(testCase.RequestBody, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.UpdateDivision(testCase.EndpointID(), testCase.RequestBody, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.DeleteDivision(testCase.EndpointID(), testCase.RequestOpts)
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
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("name", name)
		divisions, _, err := TOSession.GetDivisions(opts)
		assert.RequireNoError(t, err, "Error getting Division: %v - alerts: %+v", err, divisions.Alerts)
		assert.RequireEqual(t, 1, len(divisions.Response), "Expected one Division returned Got: %d", len(divisions.Response))
		validateDivisionFields(expectedResp)(t, toclientlib.ReqInf{}, divisions.Response, tc.Alerts{}, nil)
	}
}

func validateDivisionPagination(paginationParam string) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		paginationResp := resp.([]tc.Division)
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("orderby", "id")
		respBase, _, err := TOSession.GetDivisions(opts)
		assert.RequireNoError(t, err, "Cannot get Divisions: %v - alerts: %+v", err, respBase.Alerts)

		division := respBase.Response
		assert.RequireGreaterOrEqual(t, len(division), 2, "Need at least 2 Divisions in Traffic Ops to test pagination support, found: %d", len(division))
		switch paginationParam {
		case "limit:":
			assert.Exactly(t, division[:1], paginationResp, "expected GET Divisions with limit = 1 to return first result")
		case "offset":
			assert.Exactly(t, division[1:2], paginationResp, "expected GET Divisions with limit = 1, offset = 1 to return second result")
		case "page":
			assert.Exactly(t, division[1:2], paginationResp, "expected GET Divisions with limit = 1, page = 2 to return second result")
		}
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

func validateDivisionDescSort() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Division response to not be nil.")
		divisionDescResp := resp.([]tc.Division)
		var descSortedList []string
		var ascSortedList []string
		assert.RequireGreaterOrEqual(t, len(divisionDescResp), 2, "Need at least 2 Divisions in Traffic Ops to test desc sort, found: %d", len(divisionDescResp))
		// Get Divisions in the default ascending order for comparison.
		divisionAscResp, _, err := TOSession.GetDivisions(client.RequestOptions{})
		assert.RequireNoError(t, err, "Unexpected error getting Divisions with default sort order: %v - alerts: %+v", err, divisionAscResp.Alerts)
		// Verify the response match in length, i.e. equal amount of Divisions.
		assert.RequireEqual(t, len(divisionAscResp.Response), len(divisionDescResp), "Expected descending order response length: %v, to match ascending order response length %v", len(divisionAscResp.Response), len(divisionDescResp))
		// Insert Division names to the front of a new list, so they are now reversed to be in ascending order.
		for _, division := range divisionDescResp {
			descSortedList = append([]string{division.Name}, descSortedList...)
		}
		// Insert Division names by appending to a new list, so they stay in ascending order.
		for _, division := range divisionAscResp.Response {
			ascSortedList = append(ascSortedList, division.Name)
		}
		assert.Exactly(t, ascSortedList, descSortedList, "Division responses are not equal after reversal: %v - %v", ascSortedList, descSortedList)
	}
}

func GetDivisionID(t *testing.T, divisionName string) func() int {
	return func() int {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("name", divisionName)
		divisionsResp, _, err := TOSession.GetDivisions(opts)
		assert.RequireNoError(t, err, "Get Divisions Request failed with error:", err)
		assert.RequireEqual(t, 1, len(divisionsResp.Response), "Expected response object length 1, but got %d", len(divisionsResp.Response))
		return divisionsResp.Response[0].ID
	}
}
