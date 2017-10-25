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

func (to *Session) CDNs() ([]tc.CDN, error) {
	url := "/api/1.2/cdns.json"
	resp, err := to.request("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data tc.CDNsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data.Response, nil
}

// CDNName gets an array of CDNs
func (to *Session) CDNName(name string) ([]tc.CDN, error) {
	url := fmt.Sprintf("/api/1.2/cdns/name/%s.json", name)
	resp, err := to.request("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data tc.CDNsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data.Response, nil
}

func (to *Session) CDNSSLKeys(name string) ([]tc.CDNSSLKeys, error) {
	url := fmt.Sprintf("/api/1.2/cdns/name/%s/sslkeys.json", name)
	resp, err := to.request("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data tc.CDNSSLKeysResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data.Response, nil
}
