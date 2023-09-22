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

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

const (
	// DEPRECATED: will be removed in the next major version. Be aware this may not be the URI being requested, for clients created with Login and ClientOps.ForceLatestAPI false.
	API_LOGS = apiBase + "/logs"

	APILogs = "/logs"
)

// GetLogsByQueryParams gets a list of logs filtered by query params.
func (to *Session) GetLogsByQueryParams(queryParams string) ([]tc.Log, toclientlib.ReqInf, error) {
	uri := APILogs + queryParams
	var data tc.LogsResponse
	reqInf, err := to.get(uri, nil, &data)
	return data.Response, reqInf, err
}

// GetLogs gets a list of logs.
func (to *Session) GetLogs() ([]tc.Log, toclientlib.ReqInf, error) {
	return to.GetLogsByQueryParams("")
}

// GetLogsByLimit gets a list of logs limited to a certain number of logs.
func (to *Session) GetLogsByLimit(limit int) ([]tc.Log, toclientlib.ReqInf, error) {
	return to.GetLogsByQueryParams(fmt.Sprintf("?limit=%d", limit))
}

// GetLogsByDays gets a list of logs limited to a certain number of days.
func (to *Session) GetLogsByDays(days int) ([]tc.Log, toclientlib.ReqInf, error) {
	return to.GetLogsByQueryParams(fmt.Sprintf("?days=%d", days))
}
