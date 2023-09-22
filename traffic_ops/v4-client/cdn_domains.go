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

// GetDomains gets all CDN Domains.
func (to *Session) GetDomains(opts RequestOptions) (tc.DomainsResponse, toclientlib.ReqInf, error) {
	var data tc.DomainsResponse
	inf, err := to.get("/cdns/domains", opts, &data)
	return data, inf, err
}
