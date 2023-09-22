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

func TestPhysLocations(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Parameters, Divisions, Regions, PhysLocations}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.V3TestCaseT[tc.PhysLocation]{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfModifiedSince: {tomorrow}},
					Expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validatePhysicalLocationSort()),
				},
				"OK when VALID NAME parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"name": {"Denver"}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validatePhysicalLocationFields(map[string]interface{}{"Name": "Denver"})),
				},
				"SORTED by ID when ORDERBY=ID parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"orderby": {"id"}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validatePhysicalLocationIDSort()),
				},
				"FIRST RESULT when LIMIT=1": {
					ClientSession: TOSession,
					RequestParams: url.Values{"orderby": {"id"}, "limit": {"1"}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validatePhysicalLocationPagination("limit")),
				},
				"SECOND RESULT when LIMIT=1 OFFSET=1": {
					ClientSession: TOSession,
					RequestParams: url.Values{"orderby": {"id"}, "limit": {"1"}, "offset": {"1"}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validatePhysicalLocationPagination("offset")),
				},
				"SECOND RESULT when LIMIT=1 PAGE=2": {
					ClientSession: TOSession,
					RequestParams: url.Values{"orderby": {"id"}, "limit": {"1"}, "page": {"2"}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validatePhysicalLocationPagination("page")),
				},
				"BAD REQUEST when INVALID LIMIT parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"limit": {"-2"}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID OFFSET parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"limit": {"1"}, "offset": {"0"}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID PAGE parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"limit": {"1"}, "page": {"0"}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"OK when CHANGES made": {
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfModifiedSince: {currentTimeRFC}},
					Expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
			"POST": {
				"OK when VALID request": {
					ClientSession: TOSession,
					RequestBody: tc.PhysLocation{
						Address:    "100 blah lane",
						City:       "foo",
						Comments:   "comment",
						Email:      "bar@foobar.com",
						Name:       "testPhysicalLocation",
						Phone:      "111-222-3333",
						RegionName: "region1",
						RegionID:   GetRegionID(t, "region1")(),
						ShortName:  "testLocation1",
						State:      "CO",
						Zip:        "80602",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validatePhysicalLocationUpdateCreateFields("testPhysicalLocation", map[string]interface{}{"Name": "testPhysicalLocation"})),
				},
				"BAD REQUEST when REGION ID does NOT MATCH REGION NAME": {
					EndpointID:    GetPhysicalLocationID(t, "HotAtlanta"),
					ClientSession: TOSession,
					RequestBody: tc.PhysLocation{
						Address:    "1234 southern way",
						City:       "Atlanta",
						Name:       "HotAtlanta",
						Phone:      "404-222-2222",
						RegionName: "region1",
						RegionID:   GetRegionID(t, "cdn-region2")(),
						ShortName:  "atlanta",
						State:      "GA",
						Zip:        "30301",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
			"PUT": {
				"OK when VALID request": {
					EndpointID:    GetPhysicalLocationID(t, "HotAtlanta"),
					ClientSession: TOSession,
					RequestBody: tc.PhysLocation{
						Address:   "1234 southern way",
						City:      "NewCity",
						Name:      "HotAtlanta",
						Phone:     "404-222-2222",
						RegionID:  GetRegionID(t, "region1")(),
						ShortName: "atlanta",
						State:     "GA",
						Zip:       "30301",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validatePhysicalLocationUpdateCreateFields("HotAtlanta", map[string]interface{}{"City": "NewCity"})),
				},
				"PRECONDITION FAILED when updating with IMS & IUS Headers": {
					EndpointID:     GetPhysicalLocationID(t, "HotAtlanta"),
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}},
					RequestBody: tc.PhysLocation{
						Address:   "1234 southern way",
						City:      "Atlanta",
						RegionID:  GetRegionID(t, "region1")(),
						Name:      "HotAtlanta",
						ShortName: "atlanta",
						State:     "GA",
						Zip:       "30301",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"PRECONDITION FAILED when updating with IFMATCH ETAG Header": {
					EndpointID:    GetPhysicalLocationID(t, "HotAtlanta"),
					ClientSession: TOSession,
					RequestBody: tc.PhysLocation{
						Address:   "1234 southern way",
						City:      "Atlanta",
						RegionID:  GetRegionID(t, "region1")(),
						Name:      "HotAtlanta",
						ShortName: "atlanta",
						State:     "GA",
						Zip:       "30301",
					},
					RequestHeaders: http.Header{rfc.IfMatch: {rfc.ETag(currentTime)}},
					Expectations:   utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
			},
			"DELETE": {
				"OK when VALID request": {
					EndpointID:    GetPhysicalLocationID(t, "testDelete"),
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {

					params := make(map[string]string)
					if testCase.RequestParams != nil {
						for k, v := range testCase.RequestParams {
							params[k] = v[0]
						}
					}

					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetPhysLocationsWithHdr(params, testCase.RequestHeaders)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp, tc.Alerts{}, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.CreatePhysLocation(testCase.RequestBody)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.UpdatePhysLocationByIDWithHdr(testCase.EndpointID(), testCase.RequestBody, testCase.RequestHeaders)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.DeletePhysLocationByID(testCase.EndpointID())
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

func validatePhysicalLocationFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Physical Location response to not be nil.")
		plResp := resp.([]tc.PhysLocation)
		for field, expected := range expectedResp {
			for _, pl := range plResp {
				switch field {
				case "Name":
					assert.Equal(t, expected, pl.Name, "Expected Name to be %v, but got %s", expected, pl.Name)
				case "City":
					assert.Equal(t, expected, pl.City, "Expected City to be %v, but got %s", expected, pl.City)
				default:
					t.Errorf("Expected field: %v, does not exist in response", field)
				}
			}
		}
	}
}

func validatePhysicalLocationUpdateCreateFields(name string, expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		pl, _, err := TOSession.GetPhysLocationByNameWithHdr(name, nil)
		assert.RequireNoError(t, err, "Error getting Physical Location: %v", err)
		assert.RequireEqual(t, 1, len(pl), "Expected one Physical Location returned Got: %d", len(pl))
		validatePhysicalLocationFields(expectedResp)(t, toclientlib.ReqInf{}, pl, tc.Alerts{}, nil)
	}
}

func validatePhysicalLocationPagination(paginationParam string) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		paginationResp := resp.([]tc.PhysLocation)

		params := map[string]string{"orderby": "id"}
		respBase, _, err := TOSession.GetPhysLocationsWithHdr(params, nil)
		assert.RequireNoError(t, err, "Cannot get Physical Locations: %v", err)

		pl := respBase
		assert.RequireGreaterOrEqual(t, len(pl), 3, "Need at least 3 Physical Locations in Traffic Ops to test pagination support, found: %d", len(pl))
		switch paginationParam {
		case "limit:":
			assert.Exactly(t, pl[:1], paginationResp, "expected GET Physical Locations with limit = 1 to return first result")
		case "offset":
			assert.Exactly(t, pl[1:2], paginationResp, "expected GET Physical Locations with limit = 1, offset = 1 to return second result")
		case "page":
			assert.Exactly(t, pl[1:2], paginationResp, "expected GET Physical Locations with limit = 1, page = 2 to return second result")
		}
	}
}

func validatePhysicalLocationSort() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Physical Location response to not be nil.")
		var physLocNames []string
		physLocResp := resp.([]tc.PhysLocation)
		for _, pl := range physLocResp {
			physLocNames = append(physLocNames, pl.Name)
		}
		assert.Equal(t, true, sort.StringsAreSorted(physLocNames), "List is not sorted by their names: %v", physLocNames)
	}
}

func validatePhysicalLocationIDSort() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Physical Location response to not be nil.")
		var physLocIDs []int
		physLocResp := resp.([]tc.PhysLocation)
		for _, pl := range physLocResp {
			physLocIDs = append(physLocIDs, pl.ID)
		}
		assert.Equal(t, true, sort.IntsAreSorted(physLocIDs), "List is not sorted by their ids: %v", physLocIDs)
	}
}

func GetRegionID(t *testing.T, name string) func() int {
	return func() int {
		regions, _, err := TOSession.GetRegionByNameWithHdr(name, nil)
		assert.RequireNoError(t, err, "Get Regions Request failed with error:", err)
		assert.RequireEqual(t, 1, len(regions), "Expected response object length 1, but got %d", len(regions))
		return regions[0].ID
	}
}

func GetPhysicalLocationID(t *testing.T, name string) func() int {
	return func() int {
		physicalLocations, _, err := TOSession.GetPhysLocationByNameWithHdr(name, nil)
		assert.RequireNoError(t, err, "Get PhysLocation Request failed with error:", err)
		assert.RequireEqual(t, 1, len(physicalLocations), "Expected response object length 1, but got %d", len(physicalLocations))
		return physicalLocations[0].ID
	}
}

func CreateTestPhysLocations(t *testing.T) {
	for _, pl := range testData.PhysLocations {
		alerts, _, err := TOSession.CreatePhysLocation(pl)
		assert.RequireNoError(t, err, "Could not create Physical Location '%s': %v - alerts: %+v", pl.Name, err, alerts)
	}
}

func DeleteTestPhysLocations(t *testing.T) {
	physicalLocations, _, err := TOSession.GetPhysLocationsWithHdr(nil, nil)
	assert.NoError(t, err, "Cannot get Physical Locations: %v", err)

	for _, pl := range physicalLocations {
		alerts, _, err := TOSession.DeletePhysLocationByID(pl.ID)
		assert.NoError(t, err, "Unexpected error deleting Physical Location '%s' (#%d): %v - alerts: %+v", pl.Name, pl.ID, err, alerts.Alerts)
		// Retrieve the PhysLocation to see if it got deleted
		getPL, _, err := TOSession.GetPhysLocationByIDWithHdr(pl.ID, nil)
		assert.NoError(t, err, "Error getting Physical Location '%s' after deletion: %v", pl.Name, err)
		assert.Equal(t, 0, len(getPL), "Expected Physical Location '%s' to be deleted, but it was found in Traffic Ops", pl.Name)
	}
}
