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
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"
)

func TestRegions(t *testing.T) {
	WithObjs(t, []TCObj{Parameters, Divisions, Regions}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.TestCase[client.Session, client.RequestOptions, tc.Region]{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {tomorrow}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {currentTimeRFC}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateRegionsSort()),
				},
				"OK when VALID NAME parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"region1"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						validateRegionsFields(map[string]interface{}{"Name": "region1"})),
				},
				"OK when VALID DIVISION parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"division": {strconv.Itoa(GetDivisionID(t, "division1")())}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateRegionsFields(map[string]interface{}{"DivisionName": "division1"})),
				},
				"EMPTY RESPONSE when REGION NAME DOESNT EXIST": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"doesntexist"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when REGION ID DOESNT EXIST": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"id": {"9999999"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when DIVISION DOESNT EXIST": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"division": {"9999999"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"VALID when SORTORDER param is DESC": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"sortOrder": {"desc"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateRegionsDescSort()),
				},
				"FIRST RESULT when LIMIT=1": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateRegionsPagination("limit")),
				},
				"SECOND RESULT when LIMIT=1 OFFSET=1": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}, "offset": {"1"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateRegionsPagination("offset")),
				},
				"SECOND RESULT when LIMIT=1 PAGE=2": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}, "page": {"2"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateRegionsPagination("page")),
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
			},
			"POST": {
				"NOT FOUND when DIVISION DOESNT EXIST": {
					ClientSession: TOSession,
					RequestBody: tc.Region{
						Name:         "invalidDivision",
						Division:     99999999,
						DivisionName: "doesntexist",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
			},
			"PUT": {
				"OK when VALID request": {
					EndpointID:    GetRegionID(t, "cdn-region2"),
					ClientSession: TOSession,
					RequestBody: tc.Region{
						Name:         "newName",
						Division:     GetDivisionID(t, "cdn-div2")(),
						DivisionName: "cdn-div2",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateRegionsUpdateCreateFields("newName", map[string]interface{}{"Name": "newName"})),
				},
				"PRECONDITION FAILED when updating with IMS & IUS Headers": {
					EndpointID:    GetRegionID(t, "region1"),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}}},
					RequestBody: tc.Region{
						Name:         "newName",
						Division:     GetDivisionID(t, "division1")(),
						DivisionName: "division1",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"PRECONDITION FAILED when updating with IFMATCH ETAG Header": {
					EndpointID:    GetRegionID(t, "region1"),
					ClientSession: TOSession,
					RequestBody: tc.Region{
						Name:         "newName",
						Division:     GetDivisionID(t, "division1")(),
						DivisionName: "division1",
					},
					RequestOpts:  client.RequestOptions{Header: http.Header{rfc.IfMatch: {rfc.ETag(currentTime)}}},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
			},
			"DELETE": {
				"OK when VALID request": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"test-deletion"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"NOT FOUND when INVALID ID": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"id": {"99999999"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"NOT FOUND when INVALID NAME": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"doesntexist"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					regionName := ""
					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetRegions(testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.CreateRegion(testCase.RequestBody, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.UpdateRegion(testCase.EndpointID(), testCase.RequestBody, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							if val, ok := testCase.RequestOpts.QueryParameters["name"]; ok {
								regionName = val[0]
							}
							alerts, reqInf, err := testCase.ClientSession.DeleteRegion(regionName, testCase.RequestOpts)
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

func validateRegionsFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Regions response to not be nil.")
		regionResp := resp.([]tc.Region)
		for field, expected := range expectedResp {
			for _, region := range regionResp {
				switch field {
				case "Division":
					assert.Equal(t, expected, region.Division, "Expected Division to be %v, but got %d", expected, region.Division)
				case "DivisionName":
					assert.Equal(t, expected, region.DivisionName, "Expected DivisionName to be %v, but got %s", expected, region.DivisionName)
				case "ID":
					assert.Equal(t, expected, region.ID, "Expected ID to be %v, but got %d", expected, region.ID)
				case "Name":
					assert.Equal(t, expected, region.Name, "Expected Name to be %v, but got %s", expected, region.Name)
				default:
					t.Errorf("Expected field: %v, does not exist in response", field)
				}
			}
		}
	}
}

func validateRegionsUpdateCreateFields(name string, expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("name", name)
		region, _, err := TOSession.GetRegions(opts)
		assert.RequireNoError(t, err, "Error getting Region: %v - alerts: %+v", err, region.Alerts)
		assert.RequireEqual(t, 1, len(region.Response), "Expected one Region returned Got: %d", len(region.Response))
		validateRegionsFields(expectedResp)(t, toclientlib.ReqInf{}, region.Response, tc.Alerts{}, nil)
	}
}

func validateRegionsSort() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Regions response to not be nil.")
		var regionNames []string
		regionResp := resp.([]tc.Region)
		for _, region := range regionResp {
			regionNames = append(regionNames, region.Name)
		}
		assert.Equal(t, true, sort.StringsAreSorted(regionNames), "List is not sorted by their names: %v", regionNames)
	}
}

func validateRegionsDescSort() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Regions response to not be nil.")
		regionDescResp := resp.([]tc.Region)
		var descSortedList []string
		var ascSortedList []string
		assert.RequireGreaterOrEqual(t, len(regionDescResp), 2, "Need at least 2 Regions in Traffic Ops to test desc sort, found: %d", len(regionDescResp))
		// Get Regions in the default ascending order for comparison.
		regionAscResp, _, err := TOSession.GetRegions(client.RequestOptions{})
		assert.RequireNoError(t, err, "Unexpected error getting Regions with default sort order: %v - alerts: %+v", err, regionAscResp.Alerts)
		// Verify the response match in length, i.e. equal amount of Regions.
		assert.RequireEqual(t, len(regionAscResp.Response), len(regionDescResp), "Expected descending order response length: %v, to match ascending order response length %v", len(regionAscResp.Response), len(regionDescResp))
		// Insert Region names to the front of a new list, so they are now reversed to be in ascending order.
		for _, region := range regionDescResp {
			descSortedList = append([]string{region.Name}, descSortedList...)
		}
		// Insert Region names by appending to a new list, so they stay in ascending order.
		for _, region := range regionAscResp.Response {
			ascSortedList = append(ascSortedList, region.Name)
		}
		assert.Exactly(t, ascSortedList, descSortedList, "Region responses are not equal after reversal: %v - %v", ascSortedList, descSortedList)
	}
}

func validateRegionsPagination(paginationParam string) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		paginationResp := resp.([]tc.Region)

		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("orderby", "id")
		respBase, _, err := TOSession.GetRegions(opts)
		assert.RequireNoError(t, err, "Cannot get Regions: %v - alerts: %+v", err, respBase.Alerts)

		region := respBase.Response
		assert.RequireGreaterOrEqual(t, len(region), 2, "Need at least 2 Regions in Traffic Ops to test pagination support, found: %d", len(region))
		switch paginationParam {
		case "limit:":
			assert.Exactly(t, region[:1], paginationResp, "expected GET Regions with limit = 1 to return first result")
		case "offset":
			assert.Exactly(t, region[1:2], paginationResp, "expected GET Regions with limit = 1, offset = 1 to return second result")
		case "page":
			assert.Exactly(t, region[1:2], paginationResp, "expected GET Regions with limit = 1, page = 2 to return second result")
		}
	}
}
