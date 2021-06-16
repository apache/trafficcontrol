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
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

const apiAPICapabilities = "/api_capabilities"

// GetAPICapabilities will retrieve API Capabilities. In the event that no capability parameter
// is supplied, it will return all existing. If a capability is supplied, it will return only
// those with an exact match. Order may be specified to change the default sort order.
func (to *Session) GetAPICapabilities(opts RequestOptions) (tc.APICapabilityResponse, toclientlib.ReqInf, error) {
	var resp tc.APICapabilityResponse
	reqInf, err := to.get(apiAPICapabilities, opts, &resp)
	return resp, reqInf, err
}
