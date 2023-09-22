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

// apiLogs is the API version-relative path to the /logs API endpoint.
const apiLogs = "/logs"

// GetLogs gets a list of logs.
func (to *Session) GetLogs(opts RequestOptions) (tc.LogsResponseV5, toclientlib.ReqInf, error) {
	var data tc.LogsResponseV5
	reqInf, err := to.get(apiLogs, opts, &data)
	return data, reqInf, err
}
