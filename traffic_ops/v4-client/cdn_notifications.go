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
	"net/url"
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

const (
	APICDNNotifications = "/cdn_notifications"
)

// GetCDNNotifications returns a list of CDN Notifications.
func (to *Session) GetCDNNotifications(cdnName string, header http.Header) ([]tc.CDNNotification, toclientlib.ReqInf, error) {
	var data tc.CDNNotificationsResponse
	params := url.Values{}
	params.Add("cdn", cdnName)
	route := fmt.Sprintf("%s?%s", APICDNNotifications, params.Encode())
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// CreateCDNNotification creates a CDN notification.
func (to *Session) CreateCDNNotification(notification tc.CDNNotificationRequest) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(APICDNNotifications, notification, nil, &alerts)
	return alerts, reqInf, err
}

// DeleteCDNNotification deletes a CDN Notification by notification ID.
func (to *Session) DeleteCDNNotification(id int) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	params := url.Values{}
	params.Add("id", strconv.Itoa(id))
	route := fmt.Sprintf("%s?%s", APICDNNotifications, params.Encode())
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}
