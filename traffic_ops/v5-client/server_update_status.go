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

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
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

// SetUpdateServerStatuses updates a server's queue status and/or reval status.
// Either updateStatus or revalStatus may be nil, in which case that status
// isn't updated (but not both, because that wouldn't do anything).
func (to *Session) SetUpdateServerStatuses(serverName string, updateStatus *bool, revalStatus *bool, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	reqInf := toclientlib.ReqInf{CacheHitStatus: toclientlib.CacheHitStatusMiss}
	if updateStatus == nil && revalStatus == nil {
		return tc.Alerts{}, reqInf, errors.New("either updateStatus or revalStatus must be non-nil; nothing to do")
	}

	if opts.QueryParameters == nil {
		opts.QueryParameters = url.Values{}
	}

	if updateStatus != nil {
		if *updateStatus {
			opts.QueryParameters.Set("updated", "true")
		} else {
			opts.QueryParameters.Set("updated", "false")
		}
	}
	if revalStatus != nil {
		if *revalStatus {
			opts.QueryParameters.Set("reval_updated", "true")
		} else {
			opts.QueryParameters.Set("reval_updated", "false")
		}
	}
	var alerts tc.Alerts
	path := `/servers/` + url.PathEscape(serverName) + `/update`
	reqInf, err := to.post(path, opts, nil, &alerts)
	return alerts, reqInf, err
}
