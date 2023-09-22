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
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

const (
	// API_PING is Deprecated: will be removed in the next major version. Be aware this may not be the URI being requested, for clients created with Login and ClientOps.ForceLatestAPI false.
	API_PING = apiBase + "/ping"

	APIPing = "/ping"
)

// Ping returns a static json object to show that traffic_ops is responsive
func (to *Session) Ping() (map[string]string, toclientlib.ReqInf, error) {
	var data map[string]string
	reqInf, err := to.get(APIPing, nil, &data)
	return data, reqInf, err
}
