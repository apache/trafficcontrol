package client

import (
	"fmt"
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

// GetAPICapabilities retrieves all (or filtered) api_capability from Traffic Ops.
func (to *Session) GetAPICapabilities(capability string, order string) (tc.APICapabilityResponse, ReqInf, error) {
	var (
		vals = url.Values{}
		resp tc.APICapabilityResponse
	)

	// if capability != nil {
	vals.Set("capability", capability)
	// }

	// if order != nil {
	vals.Set("orderby", order)
	// }

	//path := fmt.Sprintf("%s/api_capabilities?%s", apiBase, vals.Encode())
	path := fmt.Sprintf("%s/api_capabilities", apiBase)
	fmt.Println("*** path:", path, "&&&")

	inf, err := get(to, path, &resp)

	return resp, inf, err
}
