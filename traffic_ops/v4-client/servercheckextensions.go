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
)

const API_TO_EXTENSION = apiBase + "/servercheck/extensions"

// CreateServerCheckExtension creates a servercheck extension.
func (to *Session) CreateServerCheckExtension(ServerCheckExtension tc.ServerCheckExtensionNullable) (tc.Alerts, ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(API_TO_EXTENSION, ServerCheckExtension, nil, &alerts)
	return alerts, reqInf, err
}

// DeleteServerCheckExtension deletes a servercheck extension.
func (to *Session) DeleteServerCheckExtension(id int) (tc.Alerts, ReqInf, error) {
	URI := fmt.Sprintf("%s/%d", API_TO_EXTENSION, id)
	var alerts tc.Alerts
	reqInf, err := to.del(URI, nil, &alerts)
	return alerts, reqInf, err
}

// GetServerCheckExtensions gets all servercheck extensions.
func (to *Session) GetServerCheckExtensions() (tc.ServerCheckExtensionResponse, ReqInf, error) {
	var toExtResp tc.ServerCheckExtensionResponse
	reqInf, err := to.get(API_TO_EXTENSION, nil, &toExtResp)
	return toExtResp, reqInf, err
}
