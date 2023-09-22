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
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
)

var (
	toReqTimeout = time.Second * time.Duration(Config.Default.Session.TimeoutInSecs)
)

func TestServerCheckExtensions(t *testing.T) {
	WithObjs(t, []TCObj{ServerCheckExtensions}, func() {

		extensionUser := utils.CreateV3Session(t, Config.TrafficOps.URL, Config.TrafficOps.Users.Extension, Config.TrafficOps.UserPassword, Config.Default.Session.TimeoutInSecs)

		methodTests := utils.V3TestCaseT[tc.ServerCheckExtensionNullable]{
			"POST": {
				"FORBIDDEN when NOT EXTENSION USER": {
					ClientSession: TOSession,
					RequestBody: tc.ServerCheckExtensionNullable{
						Name:                 util.Ptr("MEM_CHECKER"),
						Version:              util.Ptr("3.0.3"),
						InfoURL:              util.Ptr("-"),
						ScriptFile:           util.Ptr("mem.py"),
						IsActive:             util.Ptr(1),
						ServercheckShortName: util.Ptr("MC"),
						Type:                 util.Ptr("CHECK_EXTENSION_MEM"),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
				"BAD REQUEST when NO OPEN SLOTS": {
					ClientSession: extensionUser,
					RequestBody: tc.ServerCheckExtensionNullable{
						Name:                 util.Ptr("MEM_CHECKER"),
						Version:              util.Ptr("3.0.3"),
						InfoURL:              util.Ptr("-"),
						ScriptFile:           util.Ptr("mem.py"),
						IsActive:             util.Ptr(1),
						ServercheckShortName: util.Ptr("MC"),
						Type:                 util.Ptr("CHECK_EXTENSION_NUM"),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID TYPE": {
					ClientSession: extensionUser,
					RequestBody: tc.ServerCheckExtensionNullable{
						Name:                 util.Ptr("MEM_CHECKER"),
						Version:              util.Ptr("3.0.3"),
						InfoURL:              util.Ptr("-"),
						ScriptFile:           util.Ptr("mem.py"),
						IsActive:             util.Ptr(1),
						ServercheckShortName: util.Ptr("MC"),
						Type:                 util.Ptr("INVALID_TYPE"),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
		}
		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "POST":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.CreateServerCheckExtension(testCase.RequestBody)
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
	extensionUser := utils.CreateV3Session(t, Config.TrafficOps.URL, Config.TrafficOps.Users.Extension, Config.TrafficOps.UserPassword, Config.Default.Session.TimeoutInSecs)
	for _, ext := range testData.ServerCheckExtensions {
		resp, _, err := extensionUser.CreateServerCheckExtension(ext)
		assert.NoError(t, err, "Could not create Servercheck Extension: %v - alerts: %+v", err, resp.Alerts)
	}
}

func DeleteTestServerCheckExtensions(t *testing.T) {
	extensionUser := utils.CreateV3Session(t, Config.TrafficOps.URL, Config.TrafficOps.Users.Extension, Config.TrafficOps.UserPassword, Config.Default.Session.TimeoutInSecs)
	extensions, _, err := TOSession.GetServerCheckExtensions()
	assert.RequireNoError(t, err, "Could not get Servercheck Extensions: %v - alerts: %+v", err, extensions.Alerts)

	for _, extension := range extensions.Response {
		alerts, _, err := extensionUser.DeleteServerCheckExtension(*extension.ID)
		assert.NoError(t, err, "Unexpected error deleting Servercheck Extension '%s' (#%d): %v - alerts: %+v", *extension.Name, *extension.ID, err, alerts.Alerts)
		// Retrieve the Server Extension to see if it got deleted
		getExtensions, _, err := TOSession.GetServerCheckExtensions()
		assert.NoError(t, err, "Error getting Servercheck Extensions after deletion: %v - alerts: %+v", err, getExtensions.Alerts)
		for _, getExtension := range getExtensions.Response {
			if *getExtension.ID == *extension.ID {
				t.Errorf("Expected Servercheck Extension '%s' to be deleted, but it was found in Traffic Ops", *extension.Name)
			}
		}
	}
}
