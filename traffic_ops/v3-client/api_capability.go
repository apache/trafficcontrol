package client

import (
	"fmt"
	"net/url"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

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

// GetAPICapabilities will retrieve API Capabilities. In the event that no capability parameter
// is supplied, it will return all existing. If a capability is supplied, it will return only
// those with an exact match. Order may be specified to change the default sort order.
func (to *Session) GetAPICapabilities(capability string, order string) (tc.APICapabilityResponse, toclientlib.ReqInf, error) {
	var (
		vals   = url.Values{}
		path   = "/api_capabilities"
		reqInf = toclientlib.ReqInf{CacheHitStatus: toclientlib.CacheHitStatusMiss}
		resp   tc.APICapabilityResponse
	)

	if capability != "" {
		vals.Set("capability", capability)
	}

	if order != "" {
		vals.Set("orderby", order)
	}

	if len(vals) > 0 {
		path = fmt.Sprintf("%s?%s", path, vals.Encode())
	}

	reqInf, err := to.get(path, nil, &resp)

	return resp, reqInf, err
}
