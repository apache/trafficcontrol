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

// Package client provides Go bindings to the Traffic Ops RPC API.
package client

import (
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

// GetExpiringCerts gets the exiring certs within the days if 'days' param is passed
// or the full list of all Delivery services and there expirations
func (to *Session) GetExpiringCerts(opts RequestOptions) (tc.SSLKeyExpirationGetResponse, toclientlib.ReqInf, error) {
	const sslKeyExpirations = "/sslkey_expirations"

	var data tc.SSLKeyExpirationGetResponse

	reqInf, err := to.get(sslKeyExpirations, opts, &data)
	return data, reqInf, err
}
