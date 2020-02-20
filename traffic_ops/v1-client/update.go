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
	"net/http"
	"strconv"
)

// UpdateStatusClear is the value received by and posted to the Traffic Ops cache update endpoint, indicating no updates are pending for a cache.
const UpdateStatusClear = 0

// UpdateStatusPending is the value received by and posted to the Traffic Ops cache update endpoint, indicating an update is pending for a cache. That is, someone has clicked "Queue Updates" in Traffic Ops, and the cache has yet to run ORT, which will fetch updates and clear the update pending flag.
const UpdateStatusPending = 1

type UpdateResponse []Update

type Update struct {
	UpdatePending      bool   `json:"upd_pending"`
	ParentPending      bool   `json:"parent_pending"`
	RevalPending       bool   `json:"reval_pending"`
	ParentRevalPending bool   `json:"parent_reval_pending"`
	Status             string `json:"status"`
	HostID             int    `json:"host_id"`
	HostName           string `json:"host_name"`
}

func (to *Session) GetUpdate(serverName string) (Update, ReqInf, error) {
	url := "/update/" + serverName
	resp, remoteAddr, err := to.request(http.MethodGet, url, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return Update{}, reqInf, err
	}
	defer resp.Body.Close()
	var data UpdateResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return Update{}, reqInf, err
	}
	if len(data) < 1 {
		return Update{}, reqInf, errors.New("empty response")
	}
	return data[0], reqInf, nil
}

// SetUpdate sets the Update Pending and Reval Pending flags for the given cache server in Traffic Ops. This MUST only be called by an app (like ORT) which has fetched new config updates for the cache, generated or downloaded the config files onto the cache, and instructed the cache service to reload its config. This MUST NOT be called unless the cache is running the latest configuration in Traffic Ops, else the Traffic Ops cache update status will be wrong. If only the Reval or Update status has been updated, but not both, the old status should be queried from the update endpoint, and the original status for the unchanged value sent here.
func (to *Session) SetUpdate(serverName string, updatePending int, revalPending int) (ReqInf, error) {
	updateURL := "/update/" + serverName + "?" + "host_name=" + serverName + "&updated=" + strconv.Itoa(updatePending) + "&reval_updated=" + strconv.Itoa(revalPending)
	resp, remoteAddr, err := to.request(http.MethodPost, updateURL, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return reqInf, err
	}
	defer resp.Body.Close()
	// TODO return error if body is not success response
	return reqInf, nil
}
