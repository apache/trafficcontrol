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

	tc "github.com/apache/incubator-trafficcontrol/lib/go-tc"
)

// Parameters gets an array of parameter structs for the profile given
// Deprecated: use GetParameters
func (to *Session) Parameters(profileName string) ([]tc.Parameter, error) {
	ps, _, err := to.GetParameters(profileName)
	return ps, err
}

func (to *Session) GetParameters(profileName string) ([]tc.Parameter, ReqInf, error) {
	url := fmt.Sprintf("/api/1.2/parameters/profile/%s.json", profileName)
	resp, remoteAddr, err := to.request("GET", url, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.ParametersResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}
