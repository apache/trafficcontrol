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
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

// UpdateServerStatus updates the Status of the server identified by
// 'serverID'.
func (to *Session) UpdateServerStatus(serverID int, req tc.ServerPutStatus, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	path := fmt.Sprintf("servers/%d/status", serverID)
	var alerts tc.Alerts
	reqInf, err := to.put(path, opts, req, &alerts)
	return alerts, reqInf, err
}

var queueUpdateActions = map[bool]string{
	false: "dequeue",
	true:  "queue",
}

// SetServerQueueUpdate set the "updPending" field of th eserver identified by
// 'serverID' to the value of 'queueUpdate - and properly queues updates on
// parents/children as necessary.
func (to *Session) SetServerQueueUpdate(serverID int, queueUpdate bool, opts RequestOptions) (tc.ServerQueueUpdateResponse, toclientlib.ReqInf, error) {
	req := tc.ServerQueueUpdateRequest{Action: queueUpdateActions[queueUpdate]}
	var resp tc.ServerQueueUpdateResponse
	path := fmt.Sprintf("/servers/%d/queue_update", serverID)
	reqInf, err := to.post(path, opts, req, &resp)
	return resp, reqInf, err
}

// SetUpdateServerStatusTimes updates a server's config queue status and/or reval status.
// Each argument individually is optional, however at least one argument must not be nil.
func (to *Session) SetUpdateServerStatusTimes(serverName string, configApplyTime, revalApplyTime *time.Time, configUpdateFailed, revalUpdateFailed *bool, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	reqInf := toclientlib.ReqInf{CacheHitStatus: toclientlib.CacheHitStatusMiss}
	var alerts tc.Alerts

	if configApplyTime == nil && revalApplyTime == nil && configUpdateFailed == nil && revalUpdateFailed == nil {
		return alerts, reqInf, errors.New("one must be non-nil (configApplyTime, configUpdateFailed, revalApplyTime); nothing to do")
	}

	if opts.QueryParameters == nil {
		opts.QueryParameters = url.Values{}
	}

	if configApplyTime != nil {
		cat := configApplyTime.Format(time.RFC3339Nano)
		opts.QueryParameters.Set("config_apply_time", cat)
	}
	if configUpdateFailed != nil {
		opts.QueryParameters.Set("config_update_failed", strconv.FormatBool(*configUpdateFailed))
	}
	if revalApplyTime != nil {
		rat := revalApplyTime.Format(time.RFC3339Nano)
		opts.QueryParameters.Set("revalidate_apply_time", rat)
	}
	if revalUpdateFailed != nil {
		opts.QueryParameters.Set("revalidate_update_failed", strconv.FormatBool(*revalUpdateFailed))
	}

	path := `/servers/` + url.PathEscape(serverName) + `/update`
	reqInf, err := to.post(path, opts, nil, &alerts)
	return alerts, reqInf, err
}
