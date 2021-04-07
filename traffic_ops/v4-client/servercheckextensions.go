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
	"fmt"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

// APITOExtension is the API version-relative path to the /servercheck/extensions API
// endpoint.
const APITOExtension = "/servercheck/extensions"

// CreateServerCheckExtension creates the given Servercheck Extension.
func (to *Session) CreateServerCheckExtension(ServerCheckExtension tc.ServerCheckExtensionNullable) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(APITOExtension, ServerCheckExtension, nil, &alerts)
	return alerts, reqInf, err
}

// DeleteServerCheckExtension deletes the Servercheck Extension identified by
// 'id'.
func (to *Session) DeleteServerCheckExtension(id int) (tc.Alerts, toclientlib.ReqInf, error) {
	URI := fmt.Sprintf("%s/%d", APITOExtension, id)
	var alerts tc.Alerts
	reqInf, err := to.del(URI, nil, &alerts)
	return alerts, reqInf, err
}

// GetServerCheckExtensions gets all Servercheck Extensions in Traffic Ops.
func (to *Session) GetServerCheckExtensions() (tc.ServerCheckExtensionResponse, toclientlib.ReqInf, error) {
	var toExtResp tc.ServerCheckExtensionResponse
	reqInf, err := to.get(APITOExtension, nil, &toExtResp)
	return toExtResp, reqInf, err
}
