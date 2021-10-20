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
	"net"

	"github.com/apache/trafficcontrol/v6/lib/go-tc"
)

const API_SERVERCHECK = apiBase + "/servercheck"

// InsertServerCheckStatus Will insert/update the servercheck value based on if it already exists or not.
func (to *Session) InsertServerCheckStatus(status tc.ServercheckRequestNullable) (*tc.ServercheckPostResponse, ReqInf, error) {
	uri := API_SERVERCHECK
	var remoteAddr net.Addr
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	jsonReq, err := json.Marshal(status)
	if err != nil {
		return nil, reqInf, err
	}
	resp := tc.ServercheckPostResponse{}
	reqInf, err = post(to, uri, jsonReq, &resp)
	if err != nil {
		return nil, reqInf, err
	}
	return &resp, reqInf, nil
}

// GetServersChecks fetches check and meta information about servers from /servercheck.
func (to *Session) GetServersChecks() ([]tc.GenericServerCheck, tc.Alerts, ReqInf, error) {
	var response struct {
		tc.Alerts
		Response []tc.GenericServerCheck `json:"response"`
	}
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss}
	reqInf, err := get(to, API_SERVERCHECK, &response)
	return response.Response, response.Alerts, reqInf, err
}
