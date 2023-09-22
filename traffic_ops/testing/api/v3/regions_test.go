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
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

func TestRegions(t *testing.T) {
	WithObjs(t, []TCObj{Parameters, Divisions, Regions}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.V3TestCaseT[tc.Region]{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfModifiedSince: {tomorrow}},
					Expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when CHANGES made": {
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfModifiedSince: {currentTimeRFC}},
					Expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateRegionsSort()),
				},
				"OK when VALID NAME parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"name": {"region1"}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						validateRegionsFields(map[string]interface{}{"Name": "region1"})),
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
					EndpointID:     GetRegionID(t, "region1"),
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}},
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
					RequestHeaders: http.Header{rfc.IfMatch: {rfc.ETag(currentTime)}},
					Expectations:   utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
			},
			"DELETE": {
				"OK when VALID request": {
					ClientSession: TOSession,
					RequestParams: url.Values{"name": {"test-deletion"}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"NOT FOUND when INVALID ID": {
					ClientSession: TOSession,
					RequestParams: url.Values{"id": {"99999999"}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"NOT FOUND when INVALID NAME": {
					ClientSession: TOSession,
					RequestParams: url.Values{"name": {"doesntexist"}},
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
							if name == "OK when VALID NAME parameter" {
								resp, reqInf, err := testCase.ClientSession.GetRegionByNameWithHdr(testCase.RequestParams["name"][0], testCase.RequestHeaders)
								for _, check := range testCase.Expectations {
									check(t, reqInf, resp, tc.Alerts{}, err)
								}
							} else {
								resp, reqInf, err := testCase.ClientSession.GetRegionsWithHdr(testCase.RequestHeaders)
								for _, check := range testCase.Expectations {
									check(t, reqInf, resp, tc.Alerts{}, err)
								}
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.CreateRegion(testCase.RequestBody)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.UpdateRegionByIDWithHdr(testCase.EndpointID(), testCase.RequestBody, testCase.RequestHeaders)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							var regionName *string
							var regionID *int
							if val, ok := testCase.RequestParams["name"]; ok {
								regionName = &val[0]
							}
							if val, ok := testCase.RequestParams["id"]; ok {
								id, _ := strconv.Atoi(val[0])
								regionID = &id
							}
							alerts, reqInf, err := testCase.ClientSession.DeleteRegion(regionID, regionName)
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
		region, _, err := TOSession.GetRegionByNameWithHdr(name, nil)
		assert.RequireNoError(t, err, "Error getting Region: %v", err)
		assert.RequireEqual(t, 1, len(region), "Expected one Region returned Got: %d", len(region))
		validateRegionsFields(expectedResp)(t, toclientlib.ReqInf{}, region, tc.Alerts{}, nil)
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

func CreateTestRegions(t *testing.T) {
	for _, region := range testData.Regions {
		resp, _, err := TOSession.CreateRegion(region)
		assert.RequireNoError(t, err, "Could not create Region '%s': %v - alerts: %+v", region.Name, err, resp.Alerts)
	}
}

func DeleteTestRegions(t *testing.T) {
	regions, _, err := TOSession.GetRegionsWithHdr(nil)
	assert.NoError(t, err, "Cannot get Regions: %v", err)

	for _, region := range regions {
		alerts, _, err := TOSession.DeleteRegion(nil, &region.Name)
		assert.NoError(t, err, "Unexpected error deleting Region '%s' (#%d): %v - alerts: %+v", region.Name, region.ID, err, alerts.Alerts)
		// Retrieve the Region to see if it got deleted
		getRegion, _, err := TOSession.GetRegionByIDWithHdr(region.ID, nil)
		assert.NoError(t, err, "Error getting Region '%s' after deletion: %v", region.Name, err)
		assert.Equal(t, 0, len(getRegion), "Expected Region '%s' to be deleted, but it was found in Traffic Ops", region.Name)
	}
}
