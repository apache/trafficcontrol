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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

// UpdateServerStatus updates a server's status and returns the response.
func (to *Session) UpdateServerStatus(serverID int, req tc.ServerPutStatus) (*tc.Alerts, ReqInf, error) {
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, reqInf, err
	}

	path := fmt.Sprintf("%s/servers/%d/status", apiBase, serverID)
	resp, remoteAddr, err := to.request(http.MethodPut, path, reqBody)
	reqInf.RemoteAddr = remoteAddr
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	alerts := tc.Alerts{}
	if err = json.NewDecoder(resp.Body).Decode(&alerts); err != nil {
		return nil, reqInf, err
	}

	return &alerts, reqInf, nil
}

var queueUpdateActions = map[bool]string{
	false: "dequeue",
	true:  "queue",
}

// SetServerQueueUpdate updates a server's status and returns the response.
func (to *Session) SetServerQueueUpdate(serverID int, queueUpdate bool) (tc.ServerQueueUpdateResponse, ReqInf, error) {
	var (
		req    = tc.ServerQueueUpdateRequest{Action: queueUpdateActions[queueUpdate]}
		reqInf = ReqInf{CacheHitStatus: CacheHitStatusMiss}
		resp   tc.ServerQueueUpdateResponse
	)

	reqBody, err := json.Marshal(req)
	if err != nil {
		return resp, reqInf, err
	}

	path := fmt.Sprintf("%s/servers/%d/queue_update", apiBase, serverID)
	httpResp, remoteAddr, err := to.request(http.MethodPost, path, reqBody)
	reqInf.RemoteAddr = remoteAddr
	if err != nil {
		return resp, reqInf, err
	}
	defer httpResp.Body.Close()

	err = json.NewDecoder(httpResp.Body).Decode(&resp)

	return resp, reqInf, err
}

// SetUpdateServerStatuses updates a server's queue status and/or reval status.
// Either updateStatus or revalStatus may be nil, in which case that status isn't updated (but not both, because that wouldn't do anything).
//
// Deprecated: Prefer to use SetUpdateServerStatusTimes
func (to *Session) SetUpdateServerStatuses(serverName string, updateStatus *bool, revalStatus *bool) (ReqInf, error) {
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss}
	if updateStatus == nil && revalStatus == nil {
		return reqInf, errors.New("either updateStatus or revalStatus must be non-nil; nothing to do")
	}

	path := apiBase + `/servers/` + serverName + `/update?`
	queryParams := []string{}
	if updateStatus != nil {
		queryParams = append(queryParams, `updated=`+strconv.FormatBool(*updateStatus))
	}
	if revalStatus != nil {
		queryParams = append(queryParams, `reval_updated=`+strconv.FormatBool(*revalStatus))
	}
	path += strings.Join(queryParams, `&`)

	resp, remoteAddr, err := to.request(http.MethodPost, path, nil)
	reqInf.RemoteAddr = remoteAddr
	if err != nil {
		return reqInf, err
	}
	resp.Body.Close()
	return reqInf, nil
}
