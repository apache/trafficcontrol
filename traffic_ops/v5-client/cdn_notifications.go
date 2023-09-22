package client

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
	"net/url"
	"strconv"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

// apiCDNNotifications is the API version-relative path to the
// /cdn_notifications API endpoint.
const apiCDNNotifications = "/cdn_notifications"

// GetCDNNotifications returns a list of CDN Notifications.
func (to *Session) GetCDNNotifications(opts RequestOptions) (tc.CDNNotificationsResponse, toclientlib.ReqInf, error) {
	var data tc.CDNNotificationsResponse
	reqInf, err := to.get(apiCDNNotifications, opts, &data)
	return data, reqInf, err
}

// CreateCDNNotification creates a CDN notification.
func (to *Session) CreateCDNNotification(notification tc.CDNNotificationRequest, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(apiCDNNotifications, opts, notification, &alerts)
	return alerts, reqInf, err
}

// DeleteCDNNotification deletes a CDN Notification by notification ID.
func (to *Session) DeleteCDNNotification(id int, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	if opts.QueryParameters == nil {
		opts.QueryParameters = url.Values{}
	}
	opts.QueryParameters.Set("id", strconv.Itoa(id))
	reqInf, err := to.del(apiCDNNotifications, opts, &alerts)
	return alerts, reqInf, err
}
