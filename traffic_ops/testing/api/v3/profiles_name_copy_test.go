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
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
)

func TestProfilesNameCopy(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Profiles, ProfileParameters}, func() {

		methodTests := utils.V3TestCaseT[tc.ProfileCopy]{
			"POST": {
				"OK when VALID request": {
					ClientSession: TOSession,
					RequestBody: tc.ProfileCopy{
						ExistingName: "EDGE1",
						Name:         "edge1-copy",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"BAD REQUEST when NEW PROFILE NAME has SPACES": {
					ClientSession: TOSession,
					RequestBody: tc.ProfileCopy{
						ExistingName: "EDGE1",
						Name:         "Profile Has Spaces",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when PROFILE NAME ALREADY EXISTS": {
					ClientSession: TOSession,
					RequestBody: tc.ProfileCopy{
						ExistingName: "EDGE1",
						Name:         "EDGE2",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"NOT FOUND when PROFILE to COPY FROM DOESNT EXIST": {
					ClientSession: TOSession,
					RequestBody: tc.ProfileCopy{
						ExistingName: "DOESNTEXIST",
						Name:         "profileCopyFail",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "POST":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.CopyProfile(testCase.RequestBody)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					}
				}
			})
		}
	})
}
