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

func TestCoordinates(t *testing.T) {
	WithObjs(t, []TCObj{Parameters, Coordinates}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.V3TestCaseT[tc.Coordinate]{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfModifiedSince: {tomorrow}},
					Expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateCoordinateSort()),
				},
				"OK when VALID NAME parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"name": {"coordinate1"}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						validateCoordinateFields(map[string]interface{}{"Name": "coordinate1"})),
				},
				"OK when CHANGES made": {
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfModifiedSince: {currentTimeRFC}},
					Expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
			"PUT": {
				"OK when VALID request": {
					EndpointID:    GetCoordinateID(t, "coordinate2"),
					ClientSession: TOSession,
					RequestBody: tc.Coordinate{
						Latitude:  7.7,
						Longitude: 8.8,
						Name:      "coordinate2",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateCoordinateUpdateCreateFields("coordinate2", map[string]interface{}{"Latitude": 7.7, "Longitude": 8.8})),
				},
				"PRECONDITION FAILED when updating with IMS & IUS Headers": {
					EndpointID:     GetCoordinateID(t, "coordinate1"),
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}},
					RequestBody: tc.Coordinate{
						Latitude:  1.1,
						Longitude: 2.2,
						Name:      "coordinate1",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"PRECONDITION FAILED when updating with IFMATCH ETAG Header": {
					EndpointID:    GetCoordinateID(t, "coordinate1"),
					ClientSession: TOSession,
					RequestBody: tc.Coordinate{
						Latitude:  1.1,
						Longitude: 2.2,
						Name:      "coordinate1",
					},
					RequestHeaders: http.Header{rfc.IfMatch: {rfc.ETag(currentTime)}},
					Expectations:   utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
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
								resp, reqInf, err := testCase.ClientSession.GetCoordinateByNameWithHdr(testCase.RequestParams["name"][0], testCase.RequestHeaders)
								for _, check := range testCase.Expectations {
									check(t, reqInf, resp, tc.Alerts{}, err)
								}
							} else {
								resp, reqInf, err := testCase.ClientSession.GetCoordinatesWithHdr(testCase.RequestHeaders)
								for _, check := range testCase.Expectations {
									check(t, reqInf, resp, tc.Alerts{}, err)
								}
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.CreateCoordinate(testCase.RequestBody)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.UpdateCoordinateByIDWithHdr(testCase.EndpointID(), testCase.RequestBody, testCase.RequestHeaders)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.DeleteCoordinateByID(testCase.EndpointID())
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

func validateCoordinateFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Coordinate response to not be nil.")
		coordinateResp := resp.([]tc.Coordinate)
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
		coordinates, _, err := TOSession.GetCoordinateByNameWithHdr(name, nil)
		assert.RequireNoError(t, err, "Error: %v getting Coordinate: %s", err, name)
		assert.RequireEqual(t, 1, len(coordinates), "Expected one Coordinate returned Got: %d", len(coordinates))
		validateCoordinateFields(expectedResp)(t, toclientlib.ReqInf{}, coordinates, tc.Alerts{}, nil)
	}
}

func validateCoordinateSort() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Coordinate response to not be nil.")
		var coordinateNames []string
		coordinateResp := resp.([]tc.Coordinate)
		for _, coordinate := range coordinateResp {
			coordinateNames = append(coordinateNames, coordinate.Name)
		}
		assert.Equal(t, true, sort.StringsAreSorted(coordinateNames), "List is not sorted by their names: %v", coordinateNames)
	}
}

func GetCoordinateID(t *testing.T, coordinateName string) func() int {
	return func() int {
		coordinatesResp, _, err := TOSession.GetCoordinateByNameWithHdr(coordinateName, nil)
		assert.RequireNoError(t, err, "Get Coordinate Request failed with error:", err)
		assert.RequireEqual(t, 1, len(coordinatesResp), "Expected response object length 1, but got %d", len(coordinatesResp))
		return coordinatesResp[0].ID
	}
}

func CreateTestCoordinates(t *testing.T) {
	for _, coordinate := range testData.Coordinates {
		resp, _, err := TOSession.CreateCoordinate(coordinate)
		assert.RequireNoError(t, err, "Could not create coordinate: %v - alerts: %+v", err, resp.Alerts)
	}
}

func DeleteTestCoordinates(t *testing.T) {
	coordinates, _, err := TOSession.GetCoordinatesWithHdr(nil)
	assert.NoError(t, err, "Cannot get Coordinates: %v - alerts: %+v", err, coordinates)
	for _, coordinate := range coordinates {
		alerts, _, err := TOSession.DeleteCoordinateByID(coordinate.ID)
		assert.NoError(t, err, "Unexpected error deleting Coordinate '%s' (#%d): %v - alerts: %+v", coordinate.Name, coordinate.ID, err, alerts.Alerts)
		// Retrieve the Coordinate to see if it got deleted
		getCoordinate, _, err := TOSession.GetCoordinateByIDWithHdr(coordinate.ID, nil)
		assert.NoError(t, err, "Error getting Coordinate '%s' after deletion: %v", coordinate.Name, err)
		assert.Equal(t, 0, len(getCoordinate), "Expected Coordinate '%s' to be deleted, but it was found in Traffic Ops", coordinate.Name)
	}
}
