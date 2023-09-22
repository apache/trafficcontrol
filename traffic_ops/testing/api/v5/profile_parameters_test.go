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

package v5

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

func TestProfileParameters(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Profiles, ProfileParameters}, func() {

		// This is a one off test to check POST with an empty JSON body
		TestPostWithEmptyBody(t)
		currentTime := time.Now().UTC().Add(-15 * time.Second)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.V5TestCase{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {tomorrow}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID request": {
					ClientSession: TOSession,
					RequestOpts: client.RequestOptions{QueryParameters: url.Values{
						"profileId":   {strconv.Itoa(GetProfileID(t, "RASCAL1")())},
						"parameterId": {strconv.Itoa(GetParameterID(t, "peers.polling.interval", "rascal-config.txt", "60")())}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
			"POST": {
				"OK when MULTIPLE PARAMETERS": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"profileParameters": []map[string]interface{}{
							{
								"profileId":   GetProfileID(t, "MID1")(),
								"parameterId": GetParameterID(t, "CONFIG proxy.config.admin.user_id", "records.config", "STRING ats")(),
							},
							{
								"profileId":   GetProfileID(t, "MID2")(),
								"parameterId": GetParameterID(t, "CONFIG proxy.config.admin.user_id", "records.config", "STRING ats")(),
							},
						},
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusCreated)),
				},
				"BAD REQUEST when INVALID PROFILEID and PARAMETERID": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"profileId":   0,
						"parameterId": 0,
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when MISSING PROFILEID field": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"parameterId": GetParameterID(t, "health.threshold.queryTime", "rascal.properties", "1000")(),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when MISSING PARAMETERID field": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"profileId": GetProfileID(t, "EDGE2")(),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when EMPTY BODY": {
					ClientSession: TOSession,
					RequestBody:   map[string]interface{}{},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when ALREADY EXISTS": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"profileId":   GetProfileID(t, "EDGE1")(),
						"parameterId": GetParameterID(t, "health.threshold.availableBandwidthInKbps", "rascal.properties", ">1750000")(),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
			"DELETE": {
				"OK when VALID request": {
					EndpointID:    GetProfileID(t, "ATS_EDGE_TIER_CACHE"),
					ClientSession: TOSession,
					RequestOpts: client.RequestOptions{QueryParameters: url.Values{
						"parameterId": {strconv.Itoa(GetParameterID(t, "location", "set_dscp_37.config", "/etc/trafficserver/dscp")())},
					}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					profileParameter := tc.ProfileParameterCreationRequest{}
					profileParameters := []tc.ProfileParameterCreationRequest{}

					if testCase.RequestBody != nil {
						if profileParams, ok := testCase.RequestBody["profileParameters"]; ok {
							dat, err := json.Marshal(profileParams)
							assert.NoError(t, err, "Error occurred when marshalling request body: %v", err)
							err = json.Unmarshal(dat, &profileParameters)
							assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
						}
						dat, err := json.Marshal(testCase.RequestBody)
						assert.NoError(t, err, "Error occurred when marshalling request body: %v", err)
						err = json.Unmarshal(dat, &profileParameter)
						assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
					}

					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetProfileParameters(testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							if len(profileParameters) == 0 {
								alerts, reqInf, err := testCase.ClientSession.CreateProfileParameter(profileParameter, testCase.RequestOpts)
								for _, check := range testCase.Expectations {
									check(t, reqInf, nil, alerts, err)
								}
							} else {
								alerts, reqInf, err := testCase.ClientSession.CreateMultipleProfileParameters(profileParameters, testCase.RequestOpts)
								for _, check := range testCase.Expectations {
									check(t, reqInf, nil, alerts, err)
								}
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							parameterId, _ := strconv.Atoi(testCase.RequestOpts.QueryParameters["parameterId"][0])
							alerts, reqInf, err := testCase.ClientSession.DeleteProfileParameter(testCase.EndpointID(), parameterId, testCase.RequestOpts)
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

func TestPostWithEmptyBody(t *testing.T) {
	resp, err := TOSession.Client.Post(TOSession.URL+"/api/5.0/profileparameters", "application/json", nil)
	if err != nil {
		t.Fatalf("error sending post to create profile parameter with an empty body: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected to get a 400 error code, but received %d instead", resp.StatusCode)
	}
}

func TestProfileParameter(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Profiles}, func() {

		methodTests := utils.V5TestCase{
			"POST": {
				"OK when VALID request": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"profileId": GetProfileID(t, "ATS_EDGE_TIER_CACHE")(),
						"paramIds": []int64{
							int64(GetParameterID(t, "CONFIG proxy.config.allocator.enable_reclaim", "records.config", "INT 0")()),
							int64(GetParameterID(t, "CONFIG proxy.config.allocator.max_overage", "records.config", "INT 3")()),
						},
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					profileParameter := tc.PostProfileParam{}

					if testCase.RequestBody != nil {
						dat, err := json.Marshal(testCase.RequestBody)
						assert.NoError(t, err, "Error occurred when marshalling request body: %v", err)
						err = json.Unmarshal(dat, &profileParameter)
						assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
					}

					switch method {
					case "POST":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.CreateProfileWithMultipleParameters(profileParameter, testCase.RequestOpts)
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
func CreateTestProfileParameters(t *testing.T) {
	for _, profile := range testData.Profiles {
		profileID := GetProfileID(t, profile.Name)()

		for _, parameter := range profile.Parameters {
			assert.RequireNotNil(t, parameter.Name, "Expected parameter name to not be nil.")
			assert.RequireNotNil(t, parameter.Value, "Expected parameter value to not be nil.")
			assert.RequireNotNil(t, parameter.ConfigFile, "Expected parameter configFile to not be nil.")

			parameterOpts := client.NewRequestOptions()
			parameterOpts.QueryParameters.Set("name", *parameter.Name)
			parameterOpts.QueryParameters.Set("configFile", *parameter.ConfigFile)
			parameterOpts.QueryParameters.Set("value", *parameter.Value)
			getParameter, _, err := TOSession.GetParameters(parameterOpts)
			assert.RequireNoError(t, err, "Could not get Parameter %s: %v - alerts: %+v", *parameter.Name, err, getParameter.Alerts)
			if len(getParameter.Response) == 0 {
				alerts, _, err := TOSession.CreateParameter(tc.ParameterV5{Name: *parameter.Name, Value: *parameter.Value, ConfigFile: *parameter.ConfigFile}, client.RequestOptions{})
				assert.RequireNoError(t, err, "Could not create Parameter %s: %v - alerts: %+v", parameter.Name, err, alerts.Alerts)
				getParameter, _, err = TOSession.GetParameters(parameterOpts)
				assert.RequireNoError(t, err, "Could not get Parameter %s: %v - alerts: %+v", *parameter.Name, err, getParameter.Alerts)
				assert.RequireNotEqual(t, 0, len(getParameter.Response), "Could not get parameter %s: not found", *parameter.Name)
			}
			profileParameter := tc.ProfileParameterCreationRequest{ProfileID: profileID, ParameterID: getParameter.Response[0].ID}
			alerts, _, err := TOSession.CreateProfileParameter(profileParameter, client.RequestOptions{})
			assert.NoError(t, err, "Could not associate Parameter %s with Profile %s: %v - alerts: %+v", parameter.Name, profile.Name, err, alerts.Alerts)
		}
	}
}

func DeleteTestProfileParameters(t *testing.T) {
	profileParameters, _, err := TOSession.GetProfileParameters(client.RequestOptions{})
	assert.NoError(t, err, "Cannot get Profile Parameters: %v - alerts: %+v", err, profileParameters.Alerts)

	for _, profileParameter := range profileParameters.Response {
		alerts, _, err := TOSession.DeleteProfileParameter(GetProfileID(t, profileParameter.Profile)(), profileParameter.Parameter, client.RequestOptions{})
		assert.NoError(t, err, "Unexpected error deleting Profile Parameter: Profile: '%s' Parameter ID: (#%d): %v - alerts: %+v", profileParameter.Profile, profileParameter.Parameter, err, alerts.Alerts)
	}
	// Retrieve the Profile Parameters to see if it got deleted
	getProfileParameter, _, err := TOSession.GetProfileParameters(client.RequestOptions{})
	assert.NoError(t, err, "Error getting Profile Parameters after deletion: %v - alerts: %+v", err, getProfileParameter.Alerts)
	assert.Equal(t, 0, len(getProfileParameter.Response), "Expected Profile Parameters to be deleted, but %d were found in Traffic Ops", len(getProfileParameter.Response))
}
