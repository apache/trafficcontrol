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

	"github.com/apache/trafficcontrol/v6/lib/go-tc"
)

const (
	API_LOGS = apiBase + "/logs"
)

// GetLogsByQueryParams gets a list of logs filtered by query params.
func (to *Session) GetLogsByQueryParams(queryParams string) ([]tc.Log, ReqInf, error) {
	URI := API_LOGS + queryParams
	resp, remoteAddr, err := to.request(http.MethodGet, URI, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.LogsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// GetLogs gets a list of logs.
func (to *Session) GetLogs() ([]tc.Log, ReqInf, error) {
	return to.GetLogsByQueryParams("")
}

// GetLogsByLimit gets a list of logs limited to a certain number of logs.
func (to *Session) GetLogsByLimit(limit int) ([]tc.Log, ReqInf, error) {
	return to.GetLogsByQueryParams(fmt.Sprintf("?limit=%d", limit))
}

// GetLogsByDays gets a list of logs limited to a certain number of days.
func (to *Session) GetLogsByDays(days int) ([]tc.Log, ReqInf, error) {
	return to.GetLogsByQueryParams(fmt.Sprintf("?days=%d", days))
}
