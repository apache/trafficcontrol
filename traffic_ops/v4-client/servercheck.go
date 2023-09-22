package client

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

import (
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

// apiServercheck is the API version-relative path to the /servercheck API endpoint.
const apiServercheck = "/servercheck"

// InsertServerCheckStatus will insert/update the Servercheck value based on if
// it already exists or not.
func (to *Session) InsertServerCheckStatus(status tc.ServercheckRequestNullable, opts RequestOptions) (tc.ServercheckPostResponse, toclientlib.ReqInf, error) {
	var resp tc.ServercheckPostResponse
	reqInf, err := to.post(apiServercheck, opts, status, &resp)
	return resp, reqInf, err
}

// GetServersChecks fetches check and meta information about servers from
// /servercheck.
func (to *Session) GetServersChecks(opts RequestOptions) (tc.ServercheckAPIResponse, toclientlib.ReqInf, error) {
	var response tc.ServercheckAPIResponse
	reqInf, err := to.get(apiServercheck, opts, &response)
	return response, reqInf, err
}
