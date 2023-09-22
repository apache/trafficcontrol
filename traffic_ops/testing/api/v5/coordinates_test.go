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

func TestCoordinates(t *testing.T) {
	WithObjs(t, []TCObj{Parameters, Coordinates}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.TestCase[client.Session, client.RequestOptions, tc.CoordinateV5]{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {tomorrow}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateCoordinateSort()),
				},
				"OK when VALID NAME parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"coordinate1"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						validateCoordinateFields(map[string]interface{}{"Name": "coordinate1"})),
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
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateCoordinateDescSort()),
				},
				"FIRST RESULT when LIMIT=1": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateCoordinatePagination("limit")),
				},
				"SECOND RESULT when LIMIT=1 OFFSET=1": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}, "offset": {"1"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateCoordinatePagination("offset")),
				},
				"SECOND RESULT when LIMIT=1 PAGE=2": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}, "page": {"2"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateCoordinatePagination("page")),
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
				"BAD REQUEST when INVALID NAME": {
					ClientSession: TOSession,
					RequestBody: tc.CoordinateV5{
						Latitude:  1.1,
						Longitude: 2.2,
						Name:      "",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID LATITUDE": {
					ClientSession: TOSession,
					RequestBody: tc.CoordinateV5{
						Latitude:  20000,
						Longitude: 2.2,
						Name:      "testlatitude",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID LONGITUDE": {
					ClientSession: TOSession,
					RequestBody: tc.CoordinateV5{
						Latitude:  1.1,
						Longitude: 20000,
						Name:      "testlongitude",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
			"PUT": {
				"OK when VALID request": {
					EndpointID:    GetCoordinateID(t, "coordinate2"),
					ClientSession: TOSession,
					RequestBody: tc.CoordinateV5{
						Latitude:  7.7,
						Longitude: 8.8,
						Name:      "coordinate2",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateCoordinateUpdateCreateFields("coordinate2", map[string]interface{}{"Latitude": 7.7, "Longitude": 8.8})),
				},
				"NOT FOUND when INVALID ID parameter": {
					EndpointID: func() int { return 111111 },
					RequestBody: tc.CoordinateV5{
						Latitude:  1.1,
						Longitude: 2.2,
						Name:      "coordinate1",
					},
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"PRECONDITION FAILED when updating with IMS & IUS Headers": {
					EndpointID:    GetCoordinateID(t, "coordinate1"),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}}},
					RequestBody: tc.CoordinateV5{
						Latitude:  1.1,
						Longitude: 2.2,
						Name:      "coordinate1",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"PRECONDITION FAILED when updating with IFMATCH ETAG Header": {
					EndpointID:    GetCoordinateID(t, "coordinate1"),
					ClientSession: TOSession,
					RequestBody: tc.CoordinateV5{
						Latitude:  1.1,
						Longitude: 2.2,
						Name:      "coordinate1",
					},
					RequestOpts:  client.RequestOptions{Header: http.Header{rfc.IfMatch: {rfc.ETag(currentTime)}}},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
			},
			"DELETE": {
				"NOT FOUND when INVALID ID parameter": {
					EndpointID:    func() int { return 12345 },
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "GET", "GET AFTER CHANGES":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetCoordinates(testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.CreateCoordinate(testCase.RequestBody, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, resp.Alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.UpdateCoordinate(testCase.EndpointID(), testCase.RequestBody, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, resp.Alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.DeleteCoordinate(testCase.EndpointID(), testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, resp.Alerts, err)
							}
						})
					}
				}
			})
		}
	})
}

func validateCoordinateFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Coordinate response to not be nil.")
		coordinateResp := resp.([]tc.CoordinateV5)
		for field, expected := range expectedResp {
			for _, coordinate := range coordinateResp {
				switch field {
				case "Name":
					assert.Equal(t, expected, coordinate.Name, "Expected Name to be %v, but got %s", expected, coordinate.Name)
				case "Latitude":
					assert.Equal(t, expected, coordinate.Latitude, "Expected Latitude to be %v, but got %f", expected, coordinate.Latitude)
				case "Longitude":
					assert.Equal(t, expected, coordinate.Longitude, "Expected Longitude to be %v, but got %f", expected, coordinate.Longitude)
				default:
					t.Errorf("Expected field: %v, does not exist in response", field)
				}
			}
		}
	}
}

func validateCoordinateUpdateCreateFields(name string, expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("name", name)
		coordinates, _, err := TOSession.GetCoordinates(opts)
		assert.RequireNoError(t, err, "Error getting Coordinate: %v - alerts: %+v", err, coordinates.Alerts)
		assert.RequireEqual(t, 1, len(coordinates.Response), "Expected one Coordinate returned Got: %d", len(coordinates.Response))
		validateCoordinateFields(expectedResp)(t, toclientlib.ReqInf{}, coordinates.Response, tc.Alerts{}, nil)
	}
}

func validateCoordinatePagination(paginationParam string) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		paginationResp := resp.([]tc.CoordinateV5)
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("orderby", "id")
		respBase, _, err := TOSession.GetCoordinates(opts)
		assert.RequireNoError(t, err, "Cannot get Coordinates: %v - alerts: %+v", err, respBase.Alerts)

		coordinate := respBase.Response
		assert.RequireGreaterOrEqual(t, len(coordinate), 2, "Need at least 2 Coordinates in Traffic Ops to test pagination support, found: %d", len(coordinate))
		switch paginationParam {
		case "limit:":
			assert.Exactly(t, coordinate[:1], paginationResp, "expected GET Coordinates with limit = 1 to return first result")
		case "offset":
			assert.Exactly(t, coordinate[1:2], paginationResp, "expected GET Coordinates with limit = 1, offset = 1 to return second result")
		case "page":
			assert.Exactly(t, coordinate[1:2], paginationResp, "expected GET Coordinates with limit = 1, page = 2 to return second result")
		}
	}
}

func validateCoordinateSort() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Coordinate response to not be nil.")
		var coordinateNames []string
		coordinateResp := resp.([]tc.CoordinateV5)
		for _, coordinate := range coordinateResp {
			coordinateNames = append(coordinateNames, coordinate.Name)
		}
		assert.Equal(t, true, sort.StringsAreSorted(coordinateNames), "List is not sorted by their names: %v", coordinateNames)
	}
}

func validateCoordinateDescSort() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Coordinate response to not be nil.")
		coordinateDescResp := resp.([]tc.CoordinateV5)
		var descSortedList []string
		var ascSortedList []string
		assert.RequireGreaterOrEqual(t, len(coordinateDescResp), 2, "Need at least 2 Coordinates in Traffic Ops to test desc sort, found: %d", len(coordinateDescResp))
		// Get Coordinates in the default ascending order for comparison.
		coordinateAscResp, _, err := TOSession.GetCoordinates(client.RequestOptions{})
		assert.RequireNoError(t, err, "Unexpected error getting Coordinates with default sort order: %v - alerts: %+v", err, coordinateAscResp.Alerts)
		// Verify the response match in length, i.e. equal amount of Coordinates.
		assert.RequireEqual(t, len(coordinateAscResp.Response), len(coordinateDescResp), "Expected descending order response length: %d, to match ascending order response length %d", len(coordinateAscResp.Response), len(coordinateDescResp))
		// Insert Coordinate names to the front of a new list, so they are now reversed to be in ascending order.
		for _, division := range coordinateDescResp {
			descSortedList = append([]string{division.Name}, descSortedList...)
		}
		// Insert Coordinate names by appending to a new list, so they stay in ascending order.
		for _, coordinate := range coordinateAscResp.Response {
			ascSortedList = append(ascSortedList, coordinate.Name)
		}
		assert.Exactly(t, ascSortedList, descSortedList, "Coordinate responses are not equal after reversal: %v - %v", ascSortedList, descSortedList)
	}
}

func GetCoordinateID(t *testing.T, coordinateName string) func() int {
	return func() int {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("name", coordinateName)
		coordinatesResp, _, err := TOSession.GetCoordinates(opts)
		assert.RequireNoError(t, err, "Get Coordinate Request failed with error:", err)
		assert.RequireEqual(t, 1, len(coordinatesResp.Response), "Expected response object length 1, but got %d", len(coordinatesResp.Response))
		id := coordinatesResp.Response[0].ID
		assert.RequireNotNil(t, id, "Traffic Ops responded with nil Coordinate ID")
		return *id
	}
}

func CreateTestCoordinates(t *testing.T) {
	for _, coordinate := range testData.Coordinates {
		resp, _, err := TOSession.CreateCoordinate(coordinate, client.RequestOptions{})
		assert.RequireNoError(t, err, "Could not create coordinate: %v - alerts: %+v", err, resp.Alerts)
	}
}

func DeleteTestCoordinates(t *testing.T) {
	coordinates, _, err := TOSession.GetCoordinates(client.RequestOptions{})
	assert.NoError(t, err, "Cannot get Coordinates: %v - alerts: %+v", err, coordinates.Alerts)
	for _, coordinate := range coordinates.Response {
		id := coordinate.ID
		assert.RequireNotNil(t, id, "Traffic Ops responded with nil Coordinate ID")

		alerts, _, err := TOSession.DeleteCoordinate(*id, client.RequestOptions{})
		assert.NoError(t, err, "Unexpected error deleting Coordinate '%s' (#%d): %v - alerts: %+v", coordinate.Name, coordinate.ID, err, alerts.Alerts)
		// Retrieve the Coordinate to see if it got deleted
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.Itoa(*id))
		getCoordinate, _, err := TOSession.GetCoordinates(opts)
		assert.NoError(t, err, "Error getting Coordinate '%s' after deletion: %v - alerts: %+v", coordinate.Name, err, getCoordinate.Alerts)
		assert.Equal(t, 0, len(getCoordinate.Response), "Expected Coordinate '%s' to be deleted, but it was found in Traffic Ops", coordinate.Name)
	}
}
