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
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

// UpdateServerStatus updates a server's status and returns the response.
func (to *Session) UpdateServerStatus(serverID int, req tc.ServerPutStatus) (*tc.Alerts, toclientlib.ReqInf, error) {
	path := fmt.Sprintf("/servers/%d/status", serverID)
	alerts := tc.Alerts{}
	reqInf, err := to.put(path, req, nil, &alerts)
	if err != nil {
		return nil, reqInf, err
	}
	return &alerts, reqInf, nil
}

var queueUpdateActions = map[bool]string{
	false: "dequeue",
	true:  "queue",
}

// SetServerQueueUpdate updates a server's status and returns the response.
func (to *Session) SetServerQueueUpdate(serverID int, queueUpdate bool) (tc.ServerQueueUpdateResponse, toclientlib.ReqInf, error) {
	req := tc.ServerQueueUpdateRequest{Action: queueUpdateActions[queueUpdate]}
	resp := tc.ServerQueueUpdateResponse{}
	path := fmt.Sprintf("/servers/%d/queue_update", serverID)
	reqInf, err := to.post(path, req, nil, &resp)
	return resp, reqInf, err
}

// SetUpdateServerStatuses updates a server's queue status and/or reval status.
// Either updateStatus or revalStatus may be nil, in which case that status isn't updated (but not both, because that wouldn't do anything).
func (to *Session) SetUpdateServerStatuses(serverName string, updateStatus *bool, revalStatus *bool) (toclientlib.ReqInf, error) {
	reqInf := toclientlib.ReqInf{CacheHitStatus: toclientlib.CacheHitStatusMiss}
	if updateStatus == nil && revalStatus == nil {
		return reqInf, errors.New("either updateStatus or revalStatus must be non-nil; nothing to do")
	}

	path := `/servers/` + serverName + `/update?`
	queryParams := []string{}
	if updateStatus != nil {
		queryParams = append(queryParams, `updated=`+strconv.FormatBool(*updateStatus))
	}
	if revalStatus != nil {
		queryParams = append(queryParams, `reval_updated=`+strconv.FormatBool(*revalStatus))
	}
	path += strings.Join(queryParams, `&`)
	alerts := tc.Alerts{}
	reqInf, err := to.post(path, nil, nil, &alerts)
	return reqInf, err
}
