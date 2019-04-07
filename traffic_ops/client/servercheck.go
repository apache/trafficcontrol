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
	"encoding/json"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

const API_V13_SERVERCHECK = "/api/1.3/servercheck"

// Update a Server Check Status
func (to *Session) UpdateCheckStatus(status tc.ServercheckNullable) (*tc.ServercheckPostResponse, error) {
	uri := API_V13_SERVERCHECK
	jsonReq, err := json.Marshal(status)
	if err != nil {
		return nil, err
	}
	resp := tc.ServercheckPostResponse{}
	_, err = post(to, uri, jsonReq, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
