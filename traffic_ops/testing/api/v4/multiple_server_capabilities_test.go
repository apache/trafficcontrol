package v4

import (
	"net/http"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/utils"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

func TestMultipleServerCapabilities(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, ServerCapabilities, ServerServerCapabilities}, func() {

		methodTests := utils.TestCase[client.Session, client.RequestOptions, tc.MultipleServerCapabilities]{
			"PUT": {
				"OK when VALID REQUEST": {
					ClientSession: TOSession,
					RequestBody: tc.MultipleServerCapabilities{
						ServerID:           GetServerID(t, "dtrc-mid-04")(),
						ServerCapabilities: []string{"disk", "blah"},
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "PUT":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.AssignMultipleServerCapability(testCase.RequestBody, testCase.RequestOpts, testCase.RequestBody.ServerID)
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
