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
	"net"
	"net/url"
)

const APIServersStatus = apiBase + "/servers/status"

// GetServerStatusCounts gets the counts of each server status in Traffic Ops.
// If typeName is non-nil, only statuses of the given server type name will be counted.
func (to *Session) GetServerStatusCounts(typeName *string) (map[string]int, ReqInf, error) {
	var remoteAddr net.Addr
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	reqUrl := APIServersStatus
	if typeName != nil {
		reqUrl += fmt.Sprintf("?type=%s", url.QueryEscape(*typeName))
	}
	resp := struct {
		Response map[string]int `json:"response"`
	}{make(map[string]int)}

	reqInf, err := get(to, reqUrl, &resp)
	if err != nil {
		return nil, reqInf, err
	}
	return resp.Response, reqInf, nil
}
