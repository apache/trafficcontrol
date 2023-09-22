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
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

func TestMultipleServerCapabilities(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, ServerCapabilities, ServerServerCapabilities}, func() {
		var multipleSCs []string
		var multipleServerIDs []int64

		methodTests := utils.V5TestCase{
			"POST": {
				"OK When Assigned Multiple Server Capabilities": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"serverCapabilities": append(multipleSCs, "disk", "blah"),
						"serverIds":          append(multipleServerIDs, int64(GetServerID(t, "dtrc-mid-04")())),
						"pageType":           "server",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateSSC("dtrc-mid-04", "server")),
				},
				"OK When Assigned Multiple Servers Per Capability": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"serverCapabilities": append(multipleSCs, "ram"),
						"serverIds":          append(multipleServerIDs, int64(GetServerID(t, "dtrc-mid-04")()), int64(GetServerID(t, "dtrc-edge-08")())),
						"pageType":           "sc",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateSSC("ram", "sc")),
				},
				"BAD REQUEST When Assigning Many:Many Relation": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"serverCapabilities": append(multipleSCs, "ram", "disk"),
						"serverIds":          append(multipleServerIDs, int64(GetServerID(t, "dtrc-mid-04")()), int64(GetServerID(t, "dtrc-edge-08")())),
						"pageType":           "sc",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST When Missing JSON field": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"serverCapabilities": append(multipleSCs, "ram", "disk"),
						"serverIds":          append(multipleServerIDs, int64(GetServerID(t, "dtrc-mid-04")()), int64(GetServerID(t, "dtrc-edge-08")())),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
			"DELETE": {
				"OK When Delete Multiple Assigned Servers Per Capability": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"serverCapabilities": append(multipleSCs, "asdf"),
						"serverIds":          append(multipleServerIDs, int64(GetServerID(t, "dtrc-mid-04")()), int64(GetServerID(t, "dtrc-edge-08")())),
						"pageType":           "sc",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK When Delete Multiple Assigned Server Capabilities": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"serverCapabilities": append(multipleSCs, "disk", "blah"),
						"serverIds":          append(multipleServerIDs, int64(GetServerID(t, "dtrc-mid-04")())),
						"pageType":           "server",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					mssc := tc.MultipleServersCapabilities{}

					if testCase.RequestBody != nil {
						dat, err := json.Marshal(testCase.RequestBody)
						assert.NoError(t, err, "Error occurred when marshalling request body: %v", err)
						err = json.Unmarshal(dat, &mssc)
						assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
					}

					switch method {
					case "POST":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.AssignMultipleServersCapabilities(mssc, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})

					case "DELETE":
						alerts, reqInf, err := testCase.ClientSession.DeleteMultipleServersCapabilities(mssc, testCase.RequestOpts)
						for _, check := range testCase.Expectations {
							check(t, reqInf, nil, alerts, err)
						}
					}
				}
			})
		}
	})
}

func validateSSC(name, pageType string) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		opts := client.NewRequestOptions()
		switch pageType {
		case "server":
			opts.QueryParameters.Set("serverId", strconv.Itoa(GetServerID(t, name)()))
		case "sc":
			opts.QueryParameters.Set("serverCapability", name)
		}
		ssc, _, err := TOSession.GetServerServerCapabilities(opts)
		assert.RequireGreaterOrEqual(t, len(ssc.Response), 1, "Expected one or more association with:%s, Got:%d", name, len(ssc.Response))
		assert.RequireNoError(t, err, "Cannot get response: %v - alerts: %+v", err, ssc.Alerts)
	}
}
