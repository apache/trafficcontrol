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
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/assert"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/utils"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

var (
	toReqTimeout = time.Second * time.Duration(Config.Default.Session.TimeoutInSecs)
)

func TestServerCheckExtensions(t *testing.T) {
	WithObjs(t, []TCObj{ServerCheckExtensions}, func() {

		extensionUser := utils.CreateV4Session(t, Config.TrafficOps.URL, Config.TrafficOps.Users.Extension, Config.TrafficOps.UserPassword, Config.Default.Session.TimeoutInSecs)

		methodTests := utils.V4TestCase{
			"POST": {
				"FORBIDDEN when NOT EXTENSION USER": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"name":                   "MEM_CHECKER",
						"version":                "3.0.3",
						"info_url":               "-",
						"script_file":            "mem.py",
						"isactive":               1,
						"servercheck_short_name": "MC",
						"type":                   "CHECK_EXTENSION_MEM",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
				"BAD REQUEST when NO OPEN SLOTS": {
					ClientSession: extensionUser,
					RequestBody: map[string]interface{}{
						"name":                   "MEM_CHECKER",
						"version":                "3.0.3",
						"info_url":               "-",
						"script_file":            "mem.py",
						"isactive":               1,
						"servercheck_short_name": "MC",
						"type":                   "CHECK_EXTENSION_NUM",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID TYPE": {
					ClientSession: extensionUser,
					RequestBody: map[string]interface{}{
						"name":                   "MEM_CHECKER",
						"version":                "3.0.3",
						"info_url":               "-",
						"script_file":            "mem.py",
						"isactive":               1,
						"servercheck_short_name": "MC",
						"type":                   "INVALID_TYPE",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
		}
		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					serverCheckExtension := tc.ServerCheckExtensionNullable{}

					if testCase.RequestBody != nil {
						dat, err := json.Marshal(testCase.RequestBody)
						assert.NoError(t, err, "Error occurred when marshalling request body: %v", err)
						err = json.Unmarshal(dat, &serverCheckExtension)
						assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
					}

					switch method {
					case "POST":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.CreateServerCheckExtension(serverCheckExtension, testCase.RequestOpts)
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

func CreateTestServerCheckExtensions(t *testing.T) {
	extensionUser := utils.CreateV4Session(t, Config.TrafficOps.URL, Config.TrafficOps.Users.Extension, Config.TrafficOps.UserPassword, Config.Default.Session.TimeoutInSecs)
	for _, ext := range testData.ServerCheckExtensions {
		resp, _, err := extensionUser.CreateServerCheckExtension(ext, client.RequestOptions{})
		assert.NoError(t, err, "Could not create Servercheck Extension: %v - alerts: %+v", err, resp.Alerts)
	}
}

func DeleteTestServerCheckExtensions(t *testing.T) {
	extensionUser := utils.CreateV4Session(t, Config.TrafficOps.URL, Config.TrafficOps.Users.Extension, Config.TrafficOps.UserPassword, Config.Default.Session.TimeoutInSecs)
	extensions, _, err := TOSession.GetServerCheckExtensions(client.RequestOptions{})
	assert.RequireNoError(t, err, "Could not get Servercheck Extensions: %v - alerts: %+v", err, extensions.Alerts)

	for _, extension := range extensions.Response {
		alerts, _, err := extensionUser.DeleteServerCheckExtension(*extension.ID, client.RequestOptions{})
		assert.NoError(t, err, "Unexpected error deleting Servercheck Extension '%s' (#%d): %v - alerts: %+v", *extension.Name, *extension.ID, err, alerts.Alerts)
		// Retrieve the Server Extension to see if it got deleted
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.Itoa(*extension.ID))
		getExtension, _, err := TOSession.GetServerCheckExtensions(opts)
		assert.NoError(t, err, "Error getting Servercheck Extension '%s' after deletion: %v - alerts: %+v", *extension.Name, err, getExtension.Alerts)
		assert.Equal(t, 0, len(getExtension.Response), "Expected Servercheck Extension '%s' to be deleted, but it was found in Traffic Ops", *extension.Name)
	}
}
