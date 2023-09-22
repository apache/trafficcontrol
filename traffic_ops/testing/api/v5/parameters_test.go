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
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

func TestParameters(t *testing.T) {
	WithObjs(t, []TCObj{Parameters}, func() {

		opsUserSession := utils.CreateV5Session(t, Config.TrafficOps.URL, "operations", Config.TrafficOps.UserPassword, Config.Default.Session.TimeoutInSecs)

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.V5TestCase{
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
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1)),
				},
				"OK when VALID CONFIGFILE parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"configFile": {"plugin.config"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateParametersFields(map[string]interface{}{"ConfigFile": "plugin.config"})),
				},
				"OK when VALID VALUE parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"value": {"90"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateParametersFields(map[string]interface{}{"Value": "90"})),
				},
				"OK when VALID NAME parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"tm.instance_name"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateParametersFields(map[string]interface{}{"Name": "tm.instance_name"})),
				},
				"VALUE HIDDEN when OPERATIONS USER views SECURE PARAMETER": {
					ClientSession: opsUserSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"testSecure"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateParametersFields(map[string]interface{}{"Secure": true, "Value": "********"})),
				},
				"EMPTY RESPONSE when NON-EXISTENT ID parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"id": {"10000"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when NON-EXISTENT NAME parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"doesntexist"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when NON-EXISTENT CONFIGFILE parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"configFile": {"doesntexist"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when NON-EXISTENT VALUE parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"value": {"doesntexist"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"FIRST RESULT when LIMIT=1": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateParametersPagination("limit")),
				},
				"SECOND RESULT when LIMIT=1 OFFSET=1": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}, "offset": {"1"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateParametersPagination("offset")),
				},
				"SECOND RESULT when LIMIT=1 PAGE=2": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}, "page": {"2"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateParametersPagination("page")),
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
				"OK when MULTIPLE PARAMETERS": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"parameters": []map[string]interface{}{
							{
								"configFile": "multiple.config1",
								"name":       "CONFIG1 multiple config",
								"secure":     false,
								"value":      "INT 1",
							},
							{
								"configFile": "multiple.config2",
								"name":       "CONFIG2 multiple config",
								"secure":     false,
								"value":      "INT 2",
							},
							{
								"configFile": "multiple.config3",
								"name":       "CONFIG3 multiple config",
								"secure":     false,
								"value":      "INT 3",
							},
						},
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusCreated)),
				},
				"BAD REQUEST when ALREADY EXISTS": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"configFile": "records.config",
						"name":       "CONFIG proxy.config.allocator.enable_reclaim",
						"secure":     false,
						"value":      "INT 0",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when MISSING NAME FIELD": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"configFile": "missingname.config",
						"name":       "",
						"secure":     false,
						"value":      "test missing name",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when MISSING CONFIGFILE FIELD": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"configFile": "",
						"name":       "missing config file",
						"secure":     false,
						"value":      "test missing config file",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
			"PUT": {
				"OK when VALID REQUEST": {
					EndpointID:    GetParameterID(t, "LogObject.Format", "logs_xml.config", "custom_ats_2"),
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"configFile": "updated.config",
						"name":       "updated name",
						"secure":     true,
						"value":      "updated value",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateParametersUpdateCreateFields("updated name",
							map[string]interface{}{"ConfigFile": "updated.config", "Name": "updated name", "Secure": true, "Value": "updated value"})),
				},
				"OK when MISSING VALUE FIELD": {
					EndpointID:    GetParameterID(t, "LogObject.Filename", "logs_xml.config", "custom_ats_2"),
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"configFile": "logs_new.config",
						"name":       "LogObject.Filename",
						"secure":     true,
						"value":      "",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateParametersUpdateCreateFields("LogObject.Filename",
							map[string]interface{}{"ConfigFile": "logs_new.config", "Secure": true, "Value": ""})),
				},
				"BAD REQUEST when MISSING NAME FIELD": {
					EndpointID:    GetParameterID(t, "astats_over_http.so", "plugin.config", ""),
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"configFile": "missingname.config",
						"name":       "",
						"secure":     false,
						"value":      "test missing name",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when MISSING CONFIGFILE FIELD": {
					EndpointID:    GetParameterID(t, "astats_over_http.so", "plugin.config", ""),
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"configFile": "",
						"name":       "missing config file",
						"secure":     false,
						"value":      "test missing config file",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"PRECONDITION FAILED when updating with IMS & IUS Headers": {
					EndpointID:    GetParameterID(t, "LogFormat.Name", "logs_xml.config", "custom_ats_2"),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}}},
					RequestBody: map[string]interface{}{
						"configFile": "logs_xml.config",
						"name":       "LogFormat.Name",
						"secure":     false,
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"PRECONDITION FAILED when updating with IFMATCH ETAG Header": {
					EndpointID:    GetParameterID(t, "LogFormat.Name", "logs_xml.config", "custom_ats_2"),
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"configFile": "logs_xml.config",
						"name":       "LogFormat.Name",
						"secure":     false,
					},
					RequestOpts:  client.RequestOptions{Header: http.Header{rfc.IfMatch: {rfc.ETag(currentTime)}}},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
			},
			"DELETE": {
				"BAD REQUEST when DOESNT EXIST": {
					EndpointID:    func() int { return 100000 },
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					parameter := tc.ParameterV5{}
					parameters := []tc.ParameterV5{}

					if testCase.RequestBody != nil {
						if params, ok := testCase.RequestBody["parameters"]; ok {
							dat, err := json.Marshal(params)
							assert.NoError(t, err, "Error occurred when marshalling request body: %v", err)
							err = json.Unmarshal(dat, &parameters)
							assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
						}
						dat, err := json.Marshal(testCase.RequestBody)
						assert.NoError(t, err, "Error occurred when marshalling request body: %v", err)
						err = json.Unmarshal(dat, &parameter)
						assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
					}

					switch method {
					case "GET", "GET AFTER CHANGES":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetParameters(testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							if len(parameters) == 0 {
								alerts, reqInf, err := testCase.ClientSession.CreateParameter(parameter, testCase.RequestOpts)
								for _, check := range testCase.Expectations {
									check(t, reqInf, nil, alerts, err)
								}
							} else {
								alerts, reqInf, err := testCase.ClientSession.CreateMultipleParameters(parameters, testCase.RequestOpts)
								for _, check := range testCase.Expectations {
									check(t, reqInf, nil, alerts, err)
								}
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.UpdateParameter(testCase.EndpointID(), parameter, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.DeleteParameter(testCase.EndpointID(), testCase.RequestOpts)
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
		parameterResp := resp.([]tc.ParameterV5)
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
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("name", name)
		parameters, _, err := TOSession.GetParameters(opts)
		assert.RequireNoError(t, err, "Error getting Parameter: %v - alerts: %+v", err, parameters.Alerts)
		assert.RequireEqual(t, 1, len(parameters.Response), "Expected one Parameter returned Got: %d", len(parameters.Response))
		validateParametersFields(expectedResp)(t, toclientlib.ReqInf{}, parameters.Response, tc.Alerts{}, nil)
	}
}

func validateParametersPagination(paginationParam string) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		paginationResp := resp.([]tc.ParameterV5)

		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("orderby", "id")
		respBase, _, err := TOSession.GetParameters(opts)
		assert.RequireNoError(t, err, "Cannot get Parameters: %v - alerts: %+v", err, respBase.Alerts)

		parameters := respBase.Response
		assert.RequireGreaterOrEqual(t, len(parameters), 3, "Need at least 3 Parameters in Traffic Ops to test pagination support, found: %d", len(parameters))
		switch paginationParam {
		case "limit:":
			assert.Exactly(t, parameters[:1], paginationResp, "expected GET Parameters with limit = 1 to return first result")
		case "offset":
			assert.Exactly(t, parameters[1:2], paginationResp, "expected GET Parameters with limit = 1, offset = 1 to return second result")
		case "page":
			assert.Exactly(t, parameters[1:2], paginationResp, "expected GET Parameters with limit = 1, page = 2 to return second result")
		}
	}
}

func GetParameterID(t *testing.T, name string, configFile string, value string) func() int {
	return func() int {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("name", name)
		opts.QueryParameters.Set("configFile", configFile)
		if value != "" {
			opts.QueryParameters.Set("value", value)
		}
		resp, _, err := TOSession.GetParameters(opts)
		assert.RequireNoError(t, err, "Get Parameters Request failed with error: %v", err)
		assert.RequireEqual(t, 1, len(resp.Response), "Expected response object length 1, but got %d", len(resp.Response))
		return resp.Response[0].ID
	}
}

func CreateTestParameters(t *testing.T) {
	alerts, _, err := TOSession.CreateMultipleParameters(testData.Parameters, client.RequestOptions{})
	assert.RequireNoError(t, err, "Could not create Parameters: %v - alerts: %+v", err, alerts)
}

func DeleteTestParameters(t *testing.T) {
	parameters, _, err := TOSession.GetParameters(client.RequestOptions{})
	assert.RequireNoError(t, err, "Cannot get Parameters: %v - alerts: %+v", err, parameters.Alerts)

	for _, parameter := range parameters.Response {
		alerts, _, err := TOSession.DeleteParameter(parameter.ID, client.RequestOptions{})
		assert.NoError(t, err, "Cannot delete Parameter #%d: %v - alerts: %+v", parameter.ID, err, alerts.Alerts)

		// Retrieve the Parameter to see if it got deleted
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.Itoa(parameter.ID))
		getParameters, _, err := TOSession.GetParameters(opts)
		assert.NoError(t, err, "Unexpected error fetching Parameter #%d after deletion: %v - alerts: %+v", parameter.ID, err, getParameters.Alerts)
		assert.Equal(t, 0, len(getParameters.Response), "Expected Parameter '%s' to be deleted, but it was found in Traffic Ops", parameter.Name)
	}
}
