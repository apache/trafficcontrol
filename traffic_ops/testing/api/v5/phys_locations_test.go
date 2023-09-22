package v5

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
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

func TestPhysLocations(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Parameters, Divisions, Regions, PhysLocations}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.TestCase[client.Session, client.RequestOptions, tc.PhysLocationV5]{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {tomorrow}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validatePhysicalLocationSort()),
				},
				"OK when VALID NAME parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"Denver"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validatePhysicalLocationFields(map[string]interface{}{"Name": "Denver"})),
				},
				"SORTED by ID when ORDERBY=ID parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validatePhysicalLocationIDSort()),
				},
				"FIRST RESULT when LIMIT=1": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validatePhysicalLocationPagination("limit")),
				},
				"SECOND RESULT when LIMIT=1 OFFSET=1": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}, "offset": {"1"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validatePhysicalLocationPagination("offset")),
				},
				"SECOND RESULT when LIMIT=1 PAGE=2": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}, "page": {"2"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validatePhysicalLocationPagination("page")),
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
			"POST": {
				"OK when VALID request": {
					ClientSession: TOSession,
					RequestBody: tc.PhysLocationV5{
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
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusCreated),
						validatePhysicalLocationUpdateCreateFields("testPhysicalLocation", map[string]interface{}{"Name": "testPhysicalLocation"})),
				},
				"BAD REQUEST when REGION ID does NOT MATCH REGION NAME": {
					ClientSession: TOSession,
					RequestBody: tc.PhysLocationV5{
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
					RequestBody: tc.PhysLocationV5{
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
				"OK when REGION ID doesn't match REGION NAME": {
					EndpointID:    GetPhysicalLocationID(t, "HotAtlanta"),
					ClientSession: TOSession,
					RequestBody: tc.PhysLocationV5{
						Address:    "1234 southern way",
						City:       "NewCity",
						Name:       "HotAtlanta",
						Phone:      "404-222-2222",
						RegionID:   GetRegionID(t, "region1")(),
						RegionName: "notRegion1",
						ShortName:  "atlanta",
						State:      "GA",
						Zip:        "30301",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validatePhysicalLocationUpdateCreateFields("HotAtlanta", map[string]interface{}{"City": "NewCity"})),
				},
				"PRECONDITION FAILED when updating with IMS & IUS Headers": {
					EndpointID:    GetPhysicalLocationID(t, "HotAtlanta"),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}}},
					RequestBody: tc.PhysLocationV5{
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
					RequestBody: tc.PhysLocationV5{
						Address:   "1234 southern way",
						City:      "Atlanta",
						RegionID:  GetRegionID(t, "region1")(),
						Name:      "HotAtlanta",
						ShortName: "atlanta",
						State:     "GA",
						Zip:       "30301",
					},
					RequestOpts:  client.RequestOptions{Header: http.Header{rfc.IfMatch: {rfc.ETag(currentTime)}}},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
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
					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetPhysLocations(testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.CreatePhysLocation(testCase.RequestBody, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.UpdatePhysLocation(testCase.EndpointID(), testCase.RequestBody, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.DeletePhysLocation(testCase.EndpointID(), testCase.RequestOpts)
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
		plResp := resp.([]tc.PhysLocationV5)
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
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("name", name)
		pl, _, err := TOSession.GetPhysLocations(opts)
		assert.RequireNoError(t, err, "Error getting Physical Location: %v - alerts: %+v", err, pl.Alerts)
		assert.RequireEqual(t, 1, len(pl.Response), "Expected one Physical Location returned Got: %d", len(pl.Response))
		validatePhysicalLocationFields(expectedResp)(t, toclientlib.ReqInf{}, pl.Response, tc.Alerts{}, nil)
	}
}

func validatePhysicalLocationPagination(paginationParam string) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		paginationResp := resp.([]tc.PhysLocationV5)

		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("orderby", "id")
		respBase, _, err := TOSession.GetPhysLocations(opts)
		assert.RequireNoError(t, err, "Cannot get Physical Locations: %v - alerts: %+v", err, respBase.Alerts)

		pl := respBase.Response
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
		physLocResp := resp.([]tc.PhysLocationV5)
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
		physLocResp := resp.([]tc.PhysLocationV5)
		for _, pl := range physLocResp {
			physLocIDs = append(physLocIDs, pl.ID)
		}
		assert.Equal(t, true, sort.IntsAreSorted(physLocIDs), "List is not sorted by their ids: %v", physLocIDs)
	}
}

func GetRegionID(t *testing.T, name string) func() int {
	return func() int {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("name", name)
		regions, _, err := TOSession.GetRegions(opts)
		assert.RequireNoError(t, err, "Get Regions Request failed with error:", err)
		assert.RequireEqual(t, 1, len(regions.Response), "Expected response object length 1, but got %d", len(regions.Response))
		return regions.Response[0].ID
	}
}

func GetPhysicalLocationID(t *testing.T, name string) func() int {
	return func() int {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("name", name)
		physicalLocations, _, err := TOSession.GetPhysLocations(opts)
		assert.RequireNoError(t, err, "Get PhysLocation Request failed with error:", err)
		assert.RequireEqual(t, 1, len(physicalLocations.Response), "Expected response object length 1, but got %d", len(physicalLocations.Response))
		return physicalLocations.Response[0].ID
	}
}

func CreateTestPhysLocations(t *testing.T) {
	for _, pl := range testData.PhysLocations {
		resp, _, err := TOSession.CreatePhysLocation(pl, client.RequestOptions{})
		assert.RequireNoError(t, err, "Could not create Physical Location '%s': %v - alerts: %+v", pl.Name, err, resp.Alerts)
	}
}

func DeleteTestPhysLocations(t *testing.T) {
	physicalLocations, _, err := TOSession.GetPhysLocations(client.RequestOptions{})
	assert.NoError(t, err, "Cannot get Physical Locations: %v - alerts: %+v", err, physicalLocations.Alerts)

	for _, pl := range physicalLocations.Response {
		alerts, _, err := TOSession.DeletePhysLocation(pl.ID, client.RequestOptions{})
		assert.NoError(t, err, "Unexpected error deleting Physical Location '%s' (#%d): %v - alerts: %+v", pl.Name, pl.ID, err, alerts.Alerts)
		// Retrieve the PhysLocation to see if it got deleted
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.Itoa(pl.ID))
		getPL, _, err := TOSession.GetPhysLocations(opts)
		assert.NoError(t, err, "Error getting Physical Location '%s' after deletion: %v - alerts: %+v", pl.Name, err, getPL.Alerts)
		assert.Equal(t, 0, len(getPL.Response), "Expected Physical Location '%s' to be deleted, but it was found in Traffic Ops", pl.Name)
	}
}
