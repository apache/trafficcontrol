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
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

func TestCDNNotifications(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, CDNNotifications}, func() {
		GetTestCDNotifications(t)
	})
}

// Note that this test will break if anyone adds a CDN notification to the test
// data that isn't exactly `test notification: {{CDN Name}}` (where {{CDN Name}}
// is the name of the associated CDN).
func GetTestCDNotifications(t *testing.T) {
	opts := client.NewRequestOptions()
	for _, cdn := range testData.CDNs {
		opts.QueryParameters.Set("cdn", cdn.Name)
		resp, _, err := TOSession.GetCDNNotifications(opts)
		if err != nil {
			t.Errorf("cannot get CDN Notification for CDN '%s': %v - alerts: %+v", cdn.Name, err, resp.Alerts)
		}
		if len(resp.Response) > 0 {
			respNotification := resp.Response[0]
			expectedNotification := "test notification: " + cdn.Name
			if respNotification.Notification != expectedNotification {
				t.Errorf("expected notification does not match actual: %s, expected: %s", respNotification.Notification, expectedNotification)
			}
		}
	}
}

func CreateTestCDNNotifications(t *testing.T) {
	var opts client.RequestOptions
	for _, cdn := range testData.CDNs {
		resp, _, err := TOSession.CreateCDNNotification(tc.CDNNotificationRequest{CDN: cdn.Name, Notification: "test notification: " + cdn.Name}, opts)
		if err != nil {
			t.Errorf("cannot create CDN Notification for CDN '%s': %v - alerts: %+v", cdn.Name, err, resp.Alerts)
		}
	}
}

func DeleteTestCDNNotifications(t *testing.T) {
	opts := client.NewRequestOptions()
	for _, cdn := range testData.CDNs {
		// Retrieve the notifications for a cdn
		resp, _, err := TOSession.GetCDNNotifications(opts)
		if err != nil {
			t.Errorf("cannot get notifications for CDN '%s': %v - alerts: %+v", cdn.Name, err, resp.Alerts)
		}
		if len(resp.Response) > 0 {
			respNotification := resp.Response[0]
			delResp, _, err := TOSession.DeleteCDNNotification(respNotification.ID, client.RequestOptions{})
			if err != nil {
				t.Errorf("cannot delete CDN notification #%d: %v - alerts: %+v", respNotification.ID, err, delResp.Alerts)
			}
		}
	}
}
