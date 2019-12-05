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
	"net/http"

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
