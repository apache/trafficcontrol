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

package client

import (
	"fmt"
	"net/http"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

const (
	APICDNNotifications = "/cdn_notifications"
)

// GetCDNNotificationsWithHdr returns a list of CDN Notifications.
func (to *Session) GetCDNNotificationsWithHdr(cdnName string, header http.Header) ([]tc.CDNNotification, toclientlib.ReqInf, error) {
	var data tc.CDNNotificationsResponse
	route := fmt.Sprintf("%s?cdn=%s", APICDNNotifications, cdnName)
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// CreateCDNNotification creates a CDN notification.
func (to *Session) CreateCDNNotification(notification tc.CDNNotificationRequest) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(APICDNNotifications, notification, nil, &alerts)
	return alerts, reqInf, err
}

// DeleteCDNNotification deletes a CDN Notification by CDN name.
func (to *Session) DeleteCDNNotification(cdnName string) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?cdn=%s", APICDNNotifications, cdnName)
	var alerts tc.Alerts
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}
