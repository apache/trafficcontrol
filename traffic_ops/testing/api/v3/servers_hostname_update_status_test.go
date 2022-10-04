package v4

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/apache/trafficcontrol/traffic_ops/testing/api/utils"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

func TestServersHostnameUpdateStatus(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers}, func() {

		methodTests := utils.V4TestCase{
			"GET": {
				"OK when VALID request": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"hostName": {"atlanta-edge-01"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when NON-UNIQUE HOSTNAME": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"hostName": {"non-unique-hostname"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {

					switch method {
					case "GET":
						var hostName string
						if hostNameParam, ok := testCase.RequestOpts.QueryParameters["hostName"]; ok {
							hostName = hostNameParam[0]
						}
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetServerUpdateStatus(hostName, testCase.RequestOpts)
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
