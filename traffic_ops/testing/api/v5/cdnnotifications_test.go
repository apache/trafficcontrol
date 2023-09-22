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
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

func TestCDNNotifications(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, CDNNotifications}, func() {

		methodTests := utils.TestCase[client.Session, client.RequestOptions, struct{}]{
			"GET": {
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1)),
				},
				"OK when VALID CDN parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"cdn": {"cdn2"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						validateCDNNotificationFields(map[string]interface{}{"Notification": "test notification: cdn2"})),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetCDNNotifications(testCase.RequestOpts)
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

func validateCDNNotificationFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		notifications := resp.([]tc.CDNNotification)
		for field, expected := range expectedResp {
			for _, notification := range notifications {
				switch field {
				case "Notification":
					assert.Equal(t, expected, notification.Notification, "Expected Notification to be %v, but got %v", expected, notification.Notification)
				}
			}
		}
	}
}

func CreateTestCDNNotifications(t *testing.T) {
	var opts client.RequestOptions
	for _, cdn := range testData.CDNs {
		resp, _, err := TOSession.CreateCDNNotification(tc.CDNNotificationRequest{CDN: cdn.Name, Notification: "test notification: " + cdn.Name}, opts)
		assert.NoError(t, err, "Cannot create CDN Notification for CDN '%s': %v - alerts: %+v", cdn.Name, err, resp.Alerts)
	}
}

func DeleteTestCDNNotifications(t *testing.T) {
	resp, _, err := TOSession.GetCDNNotifications(client.RequestOptions{})
	assert.NoError(t, err, "Cannot get notifications for CDNs: %v - alerts: %+v", err, resp.Alerts)
	for _, notification := range resp.Response {
		delResp, _, err := TOSession.DeleteCDNNotification(notification.ID, client.RequestOptions{})
		assert.NoError(t, err, "Cannot delete CDN notification #%d: %v - alerts: %+v", notification.ID, err, delResp.Alerts)
		// Retrieve CDN Notification to see if it got deleted
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.Itoa(notification.ID))
		getNotification, _, err := TOSession.GetCDNNotifications(opts)
		assert.NoError(t, err, "Error deleting CDN Notification for '%s' : %v - alerts: %+v", notification.ID, err, getNotification.Alerts)
		assert.Equal(t, 0, len(getNotification.Response), "Expected CDN Notification '%s' to be deleted", notification.ID)
	}
}
