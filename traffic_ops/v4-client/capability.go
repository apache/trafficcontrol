package client

/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * 	http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import (
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

// apiCapabilities is the API version-relative path for the /capabilities API endpoint.
const apiCapabilities = "/capabilities"

// GetCapabilities retrieves capabilities.
func (to *Session) GetCapabilities(opts RequestOptions) (tc.CapabilitiesResponse, toclientlib.ReqInf, error) {
	var data tc.CapabilitiesResponse
	reqInf, err := to.get(apiCapabilities, opts, &data)
	return data, reqInf, err
}
