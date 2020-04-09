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
	"fmt"

	tc "github.com/apache/trafficcontrol/lib/go-tc"
)

// Hardware gets an array of Hardware
// Deprecated: use GetHardware
func (to *Session) Hardware(limit int) ([]tc.Hardware, error) {
	h, _, err := to.GetHardware(limit)
	return h, err
}

// GetHardware fetches an array of Hardware up to as many as 'limit' specifies.
// Deprecated: Hardware is deprecated and will not be exposed through the API in the future.
func (to *Session) GetHardware(limit int) ([]tc.Hardware, ReqInf, error) {
	url := "/api/1.2/hwinfo.json"
	if limit > 0 {
		url += fmt.Sprintf("?limit=%v", limit)
	}
	resp, remoteAddr, err := to.request("GET", url, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.HardwareResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}
