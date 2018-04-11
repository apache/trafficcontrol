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

	"github.com/apache/incubator-trafficcontrol/lib/go-tc/v13"
)

// CacheGroups gets the CacheGroups in an array of CacheGroup structs
// (note CacheGroup used to be called location)
// Deprecated: use GetCacheGroups.
func (to *Session) CacheGroups() ([]v13.CacheGroup, error) {
	cgs, _, err := to.GetCacheGroups()
	return cgs, err
}

func (to *Session) GetCacheGroups() ([]v13.CacheGroup, ReqInf, error) {
	url := "/api/1.2/cachegroups.json"
	resp, remoteAddr, err := to.request("GET", url, nil) // TODO change to getBytesWithTTL, return CacheHitStatus
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data v13.CacheGroupsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}
