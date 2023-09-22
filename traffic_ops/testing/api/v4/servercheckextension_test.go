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
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"
)

var (
	toReqTimeout = time.Second * time.Duration(Config.Default.Session.TimeoutInSecs)
)

func TestServerCheckExtensions(t *testing.T) {
	WithObjs(t, []TCObj{ServerCheckExtensions}, func() {

		t.Logf("TestServerCheckExtensions user '%v' pass '%v'\n", Config.TrafficOps.Users.Extension, Config.TrafficOps.UserPassword)
		extensionUser := utils.CreateV4Session(t, Config.TrafficOps.URL, Config.TrafficOps.Users.Extension, Config.TrafficOps.UserPassword, Config.Default.Session.TimeoutInSecs)

		methodTests := utils.TestCase[client.Session, client.RequestOptions, tc.ServerCheckExtensionNullable]{
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
							alerts, reqInf, err := testCase.ClientSession.CreateServerCheckExtension(testCase.RequestBody, testCase.RequestOpts)
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
