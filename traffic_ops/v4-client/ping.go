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

// apiPing is the full path to the /ping API endpoint.
const apiPing = "/ping"

// PingResponse is the type of a response from Traffic Ops to a requestt made
// to its /ping API endpoint.
type PingResponse struct {
	Ping string `json:"ping"`
	tc.Alerts
}

// Ping returns a simple response to show that Traffic Ops is responsive.
func (to *Session) Ping(opts RequestOptions) (PingResponse, toclientlib.ReqInf, error) {
	var data PingResponse
	reqInf, err := to.get(apiPing, opts, &data)
	return data, reqInf, err
}
