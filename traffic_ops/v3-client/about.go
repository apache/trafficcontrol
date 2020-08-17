// Package client implements methods for interacting with the Traffic Ops API.
//
// Warning: Using the un-versioned import path ("client") is deprecated, and the
// ability to do so will be removed in ATC 6.0 - please use versioned client
// imports (e.g. "v3-client") instead
package client

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

import (
	"encoding/json"
	"net/http"
)

const (
	API_ABOUT = apiBase + "/about"
)

// GetAbout gets data about the TO instance.
func (to *Session) GetAbout() (map[string]string, ReqInf, error) {
	route := API_ABOUT
	resp, remoteAddr, err := to.request(http.MethodGet, route, nil, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data, reqInf, nil
}
