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
	"github.com/apache/trafficcontrol/lib/go-tc"
	"testing"
)

func TestCDNNotifications(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, CDNNotifications}, func() {
		GetTestCDNotifications(t)
	})
}

func GetTestCDNotifications(t *testing.T) {
	for _, cdn := range testData.CDNs {
		resp, _, err := TOSession.GetCDNByName(cdn.Name)
		if err != nil {
			t.Errorf("cannot GET CDN by name: %v - %v", err, resp)
		}
		if len(resp) > 0 {
			respCDN := resp[0]
			expectedNotification := "test notification: " + respCDN.Name
			if respCDN.Notification != expectedNotification {
				t.Errorf("expected notification does not match actual: %s, expected: %s", respCDN.Notification, expectedNotification)
			}
		}
	}
}

func CreateTestCDNNotifications(t *testing.T) {
	for _, cdn := range testData.CDNs {
		// Retrieve the CDN by name so we can get the CDN ID to create the notification
		resp, _, err := TOSession.GetCDNByName(cdn.Name)
		if err != nil {
			t.Errorf("cannot GET CDN by name: %v - %v", cdn.Name, err)
		}
		if len(resp) > 0 {
			respCDN := resp[0]
			_, _, err := TOSession.CreateCDNNotification(respCDN, tc.CDNNotificationRequest{Notification: "test notification: " + respCDN.Name})
			if err != nil {
				t.Errorf("cannot CREATE CDN notification: '%s' %v", respCDN.Name, err)
			}
		}
	}
}

func DeleteTestCDNNotifications(t *testing.T) {
	for _, cdn := range testData.CDNs {
		// Retrieve the CDN by name so we can get the CDN ID to delete the notification
		resp, _, err := TOSession.GetCDNByName(cdn.Name)
		if err != nil {
			t.Errorf("cannot GET CDN by name: %v - %v", cdn.Name, err)
		}
		if len(resp) > 0 {
			respCDN := resp[0]
			_, _, err := TOSession.DeleteCDNNotification(respCDN)
			if err != nil {
				t.Errorf("cannot DELETE CDN notification: '%s' %v", respCDN.Name, err)
			}
		}
	}
}
