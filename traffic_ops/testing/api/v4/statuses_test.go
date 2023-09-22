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

func TestStatuses(t *testing.T) {
	WithObjs(t, []TCObj{Parameters, Statuses}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.TestCase[client.Session, client.RequestOptions, tc.Status]{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {tomorrow}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateStatusesSort()),
				},
				"OK when VALID NAME parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"CCR_IGNORE"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateStatusesFields(map[string]interface{}{"Name": "CCR_IGNORE"})),
				},
				"OK when CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {currentTimeRFC}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
			"PUT": {
				"OK when VALID request": {
					EndpointID:    GetStatusID(t, "TEST_NULL_DESCRIPTION"),
					ClientSession: TOSession,
					RequestBody: tc.Status{
						Description: "new description",
						Name:        "TEST_NULL_DESCRIPTION",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateStatusesUpdateCreateFields("TEST_NULL_DESCRIPTION", map[string]interface{}{"Description": "new description"})),
				},
				"PRECONDITION FAILED when updating with IMS & IUS Headers": {
					EndpointID:    GetStatusID(t, "TEST_NULL_DESCRIPTION"),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}}},
					RequestBody: tc.Status{
						Description: "new description",
						Name:        "TEST_NULL_DESCRIPTION",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"PRECONDITION FAILED when updating with IFMATCH ETAG Header": {
					EndpointID:    GetStatusID(t, "TEST_NULL_DESCRIPTION"),
					ClientSession: TOSession,
					RequestBody: tc.Status{
						Description: "new description",
						Name:        "TEST_NULL_DESCRIPTION",
					},
					RequestOpts:  client.RequestOptions{Header: http.Header{rfc.IfMatch: {rfc.ETag(currentTime)}}},
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
							resp, reqInf, err := testCase.ClientSession.GetStatuses(testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.UpdateStatus(testCase.EndpointID(), testCase.RequestBody, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.DeleteStatus(testCase.EndpointID(), testCase.RequestOpts)
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

func validateStatusesFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Status response to not be nil.")
		statusResp := resp.([]tc.Status)
		for field, expected := range expectedResp {
			for _, status := range statusResp {
				switch field {
				case "Description":
					assert.Equal(t, expected, status.Description, "Expected Description to be %v, but got %s", expected, status.Description)
				case "Name":
					assert.Equal(t, expected, status.Name, "Expected Name to be %v, but got %s", expected, status.Name)
				default:
					t.Errorf("Expected field: %v, does not exist in response", field)
				}
			}
		}
	}
}

func validateStatusesUpdateCreateFields(name string, expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("name", name)
		statuses, _, err := TOSession.GetStatuses(opts)
		assert.RequireNoError(t, err, "Error getting Statuses: %v - alerts: %+v", err, statuses.Alerts)
		assert.RequireEqual(t, 1, len(statuses.Response), "Expected one Status returned Got: %d", len(statuses.Response))
		validateStatusesFields(expectedResp)(t, toclientlib.ReqInf{}, statuses.Response, tc.Alerts{}, nil)
	}
}

func validateStatusesSort() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Status response to not be nil.")
		var statusNames []string
		statusResp := resp.([]tc.Status)
		for _, status := range statusResp {
			statusNames = append(statusNames, status.Name)
		}
		assert.Equal(t, true, sort.StringsAreSorted(statusNames), "List is not sorted by their names: %v", statusNames)
	}
}

func GetStatusID(t *testing.T, name string) func() int {
	return func() int {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("name", name)
		statusResp, _, err := TOSession.GetStatuses(opts)
		assert.RequireNoError(t, err, "Get Statuses Request failed with error:", err)
		assert.RequireEqual(t, 1, len(statusResp.Response), "Expected response object length 1, but got %d", len(statusResp.Response))
		return statusResp.Response[0].ID
	}
}

func CreateTestStatuses(t *testing.T) {
	for _, status := range testData.Statuses {
		resp, _, err := TOSession.CreateStatus(status, client.RequestOptions{})
		assert.RequireNoError(t, err, "Could not create Status: %v - alerts: %+v", err, resp.Alerts)
	}
}

func DeleteTestStatuses(t *testing.T) {
	opts := client.NewRequestOptions()
	for _, status := range testData.Statuses {
		assert.RequireNotNil(t, status.Name, "Cannot get test statuses: test data statuses must have names")
		// Retrieve the Status by name, so we can get the id for the Update
		opts.QueryParameters.Set("name", *status.Name)
		resp, _, err := TOSession.GetStatuses(opts)
		assert.RequireNoError(t, err, "Cannot get Statuses filtered by name '%s': %v - alerts: %+v", *status.Name, err, resp.Alerts)
		assert.RequireEqual(t, 1, len(resp.Response), "Expected 1 status returned. Got: %d", len(resp.Response))
		respStatus := resp.Response[0]

		delResp, _, err := TOSession.DeleteStatus(respStatus.ID, client.RequestOptions{})
		assert.NoError(t, err, "Cannot delete Status: %v - alerts: %+v", err, delResp.Alerts)

		// Retrieve the Status to see if it got deleted
		resp, _, err = TOSession.GetStatuses(opts)
		assert.NoError(t, err, "Unexpected error getting Statuses filtered by name after deletion: %v - alerts: %+v", err, resp.Alerts)
		assert.Equal(t, 0, len(resp.Response), "Expected Status '%s' to be deleted, but it was found in Traffic Ops", *status.Name)
	}
}
