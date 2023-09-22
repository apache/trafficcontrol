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

package v3

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

func TestParameters(t *testing.T) {
	WithObjs(t, []TCObj{Parameters}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.V3TestCaseT[tc.Parameter]{
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
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1)),
				},
				"OK when VALID NAME parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"name": {"tm.instance_name"}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateParametersFields(map[string]interface{}{"Name": "tm.instance_name"})),
				},
			},
			"PUT": {
				"OK when VALID REQUEST": {
					EndpointID:    GetParameterID(t, "LogObject.Format", "logs_xml.config", "custom_ats_2"),
					ClientSession: TOSession,
					RequestBody: tc.Parameter{
						ConfigFile: "updated.config",
						Name:       "updated name",
						Secure:     true,
						Value:      "updated value",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateParametersUpdateCreateFields("updated name",
							map[string]interface{}{"ConfigFile": "updated.config", "Name": "updated name", "Secure": true, "Value": "updated value"})),
				},
				"PRECONDITION FAILED when updating with IMS & IUS Headers": {
					EndpointID:     GetParameterID(t, "LogFormat.Name", "logs_xml.config", "custom_ats_2"),
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}},
					RequestBody: tc.Parameter{
						ConfigFile: "logs_xml.config",
						Name:       "LogFormat.Name",
						Secure:     false,
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"PRECONDITION FAILED when updating with IFMATCH ETAG Header": {
					EndpointID:    GetParameterID(t, "LogFormat.Name", "logs_xml.config", "custom_ats_2"),
					ClientSession: TOSession,
					RequestBody: tc.Parameter{
						ConfigFile: "logs_xml.config",
						Name:       "LogFormat.Name",
						Secure:     false,
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
								resp, reqInf, err := testCase.ClientSession.GetParameterByNameWithHdr(testCase.RequestParams["name"][0], testCase.RequestHeaders)
								for _, check := range testCase.Expectations {
									check(t, reqInf, resp, tc.Alerts{}, err)
								}
							} else {
								resp, reqInf, err := testCase.ClientSession.GetParametersWithHdr(testCase.RequestHeaders)
								for _, check := range testCase.Expectations {
									check(t, reqInf, resp, tc.Alerts{}, err)
								}
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.CreateParameter(testCase.RequestBody)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.UpdateParameterByIDWithHdr(testCase.EndpointID(), testCase.RequestBody, testCase.RequestHeaders)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.DeleteParameterByID(testCase.EndpointID())
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

func validateParametersFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Parameters response to not be nil.")
		parameterResp := resp.([]tc.Parameter)
		for field, expected := range expectedResp {
			for _, parameter := range parameterResp {
				switch field {
				case "ConfigFile":
					assert.Equal(t, expected, parameter.ConfigFile, "Expected ConfigFile to be %v, but got %s", expected, parameter.ConfigFile)
				case "ID":
					assert.Equal(t, expected, parameter.ID, "Expected ID to be %v, but got %d", expected, parameter.ID)
				case "Name":
					assert.Equal(t, expected, parameter.Name, "Expected Name to be %v, but got %s", expected, parameter.Name)
				case "Secure":
					assert.Equal(t, expected, parameter.Secure, "Expected Secure to be %v, but got %v", expected, parameter.Secure)
				case "Value":
					assert.Equal(t, expected, parameter.Value, "Expected Value to be %v, but got %s", expected, parameter.Value)
				default:
					t.Fatalf("Expected field: %v, does not exist in response", field)
				}
			}
		}
	}
}

func validateParametersUpdateCreateFields(name string, expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		parameters, _, err := TOSession.GetParameterByNameWithHdr(name, nil)
		assert.RequireNoError(t, err, "Error getting Parameter: %v", err)
		assert.RequireEqual(t, 1, len(parameters), "Expected one Parameter returned Got: %d", len(parameters))
		validateParametersFields(expectedResp)(t, toclientlib.ReqInf{}, parameters, tc.Alerts{}, nil)
	}
}

func GetParameterID(t *testing.T, name string, configFile string, value string) func() int {
	return func() int {
		resp, _, err := TOSession.GetParameterByNameAndConfigFileAndValueWithHdr(name, configFile, value, nil)
		assert.RequireNoError(t, err, "Get Parameters Request failed with error: %v", err)
		assert.RequireEqual(t, 1, len(resp), "Expected response object length 1, but got %d", len(resp))
		return resp[0].ID
	}
}

func CreateTestParameters(t *testing.T) {
	alerts, _, err := TOSession.CreateMultipleParameters(testData.Parameters)
	assert.RequireNoError(t, err, "Could not create Parameters: %v - alerts: %+v", err, alerts)
}

func DeleteTestParameters(t *testing.T) {
	parameters, _, err := TOSession.GetParametersWithHdr(nil)
	assert.RequireNoError(t, err, "Cannot get Parameters: %v", err)

	for _, parameter := range parameters {
		alerts, _, err := TOSession.DeleteParameterByID(parameter.ID)
		assert.NoError(t, err, "Cannot delete Parameter #%d: %v - alerts: %+v", parameter.ID, err, alerts.Alerts)

		// Retrieve the Parameter to see if it got deleted
		getParameters, _, err := TOSession.GetParameterByIDWithHdr(parameter.ID, nil)
		assert.NoError(t, err, "Unexpected error fetching Parameter #%d after deletion: %v", parameter.ID, err)
		assert.Equal(t, 0, len(getParameters), "Expected Parameter '%s' to be deleted, but it was found in Traffic Ops", parameter.Name)
	}
}
