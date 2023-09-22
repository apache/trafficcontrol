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

const (
	// apiVaultPing is the partial path (excluding the /api/<version> prefix) to the /vault/ping API endpoint.
	apiVaultPing = "/vault/ping"
)

// TrafficVaultPing returns a response indicating whether or not Traffic Vault is responsive.
func (to *Session) TrafficVaultPing(opts RequestOptions) (tc.TrafficVaultPingResponse, toclientlib.ReqInf, error) {
	var data tc.TrafficVaultPingResponse
	reqInf, err := to.get(apiVaultPing, opts, &data)
	return data, reqInf, err
}
