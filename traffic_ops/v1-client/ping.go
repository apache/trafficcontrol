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
	"net/http"
)

const (
	API_v13_PING = "/api/1.3/ping"
)

// Ping returns a static json object to show that traffic_ops is responsive
func (to *Session) Ping() (map[string]string, ReqInf, error) {
	resp, remoteAddr, err := to.request(http.MethodGet, API_v13_PING, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data map[string]string
	err = json.NewDecoder(resp.Body).Decode(&data)
	return data, reqInf, nil
}
