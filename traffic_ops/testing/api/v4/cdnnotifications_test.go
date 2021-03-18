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
)

func TestCDNNotifications(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, CDNNotifications}, func() {
		GetTestCDNotifications(t)
	})
}

func GetTestCDNotifications(t *testing.T) {
	for _, cdn := range testData.CDNs {
		resp, _, err := TOSession.GetCDNNotifications(cdn.Name, nil)
		if err != nil {
			t.Errorf("cannot GET cdn notification for cdn: %v - %v", err, resp)
		}
		if len(resp) > 0 {
			respNotification := resp[0]
			expectedNotification := "test notification: " + cdn.Name
			if respNotification.Notification != expectedNotification {
				t.Errorf("expected notification does not match actual: %s, expected: %s", respNotification.Notification, expectedNotification)
			}
		}
	}
}

func CreateTestCDNNotifications(t *testing.T) {
	for _, cdn := range testData.CDNs {
		_, _, err := TOSession.CreateCDNNotification(tc.CDNNotificationRequest{CDN: cdn.Name, Notification: "test notification: " + cdn.Name})
		if err != nil {
			t.Errorf("cannot CREATE CDN notification: '%s' %v", cdn.Name, err)
		}
	}
}

func DeleteTestCDNNotifications(t *testing.T) {
	for _, cdn := range testData.CDNs {
		// Retrieve the notifications for a cdn
		resp, _, err := TOSession.GetCDNNotifications(cdn.Name, nil)
		if err != nil {
			t.Errorf("cannot GET notifications for a CDN: %v - %v", cdn.Name, err)
		}
		if len(resp) > 0 {
			respNotification := resp[0]
			_, _, err := TOSession.DeleteCDNNotification(respNotification.ID)
			if err != nil {
				t.Errorf("cannot DELETE CDN notification by ID: '%d' %v", respNotification.ID, err)
			}
		}
	}
}
