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
	"fmt"
	"net"
	"net/http"

	"github.com/apache/trafficcontrol/v6/lib/go-tc"
)

const API_TO_EXTENSION = apiBase + "/servercheck/extensions"

// CreateServerCheckExtension creates a servercheck extension.
func (to *Session) CreateServerCheckExtension(ServerCheckExtension tc.ServerCheckExtensionNullable) (tc.Alerts, ReqInf, error) {
	var remoteAddr net.Addr
	reqBody, err := json.Marshal(ServerCheckExtension)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	resp, remoteAddr, err := to.request(http.MethodPost, API_TO_EXTENSION, reqBody)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, err
}

// DeleteServerCheckExtension deletes a servercheck extension.
func (to *Session) DeleteServerCheckExtension(id int) (tc.Alerts, ReqInf, error) {
	URI := fmt.Sprintf("%s/%d", API_TO_EXTENSION, id)
	resp, remoteAddr, err := to.request(http.MethodDelete, URI, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, err
}

// GetServerCheckExtensions gets all servercheck extensions.
func (to *Session) GetServerCheckExtensions() (tc.ServerCheckExtensionResponse, ReqInf, error) {
	resp, remoteAddr, err := to.request(http.MethodGet, API_TO_EXTENSION, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.ServerCheckExtensionResponse{}, reqInf, err
	}
	defer resp.Body.Close()
	var toExtResp tc.ServerCheckExtensionResponse
	err = json.NewDecoder(resp.Body).Decode(&toExtResp)
	return toExtResp, reqInf, err
}
