package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/apache/trafficcontrol/lib/go-tc"
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
func (to *Session) GetAPICapabilities(capability string, order string) (tc.APICapabilityResponse, ReqInf, error) {
	var (
		vals   = url.Values{}
		path   = fmt.Sprintf("%s/api_capabilities", apiBase)
		reqInf = ReqInf{CacheHitStatus: CacheHitStatusMiss}
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

	httpResp, remoteAddr, err := to.request(http.MethodGet, path, nil)
	reqInf.RemoteAddr = remoteAddr

	if err != nil {
		return tc.APICapabilityResponse{}, reqInf, err
	}
	defer httpResp.Body.Close()

	err = json.NewDecoder(httpResp.Body).Decode(&resp)

	return resp, reqInf, err
}
