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

import "encoding/json"
import "errors"
import "net/http"
import "net/url"

import "github.com/apache/trafficcontrol/v6/lib/go-tc"

const API_CAPABILITIES = apiBase + "/capabilities"

// GetCapabilities retrieves all capabilities.
func (to *Session) GetCapabilities() ([]tc.Capability, ReqInf, error) {
	resp, remoteAddr, err := to.request(http.MethodGet, API_CAPABILITIES, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.CapabilitiesResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	return data.Response, reqInf, err
}

// GetCapability retrieves only the capability named 'c'.
func (to *Session) GetCapability(c string) (tc.Capability, ReqInf, error) {
	v := url.Values{}
	v.Add("name", c)
	endpoint := API_CAPABILITIES + "?" + v.Encode()
	resp, remoteAddr, err := to.request(http.MethodGet, endpoint, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Capability{}, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.CapabilitiesResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return tc.Capability{}, reqInf, err
	} else if data.Response == nil || len(data.Response) < 1 {
		return tc.Capability{}, reqInf, errors.New("Invalid response - no capability returned!")
	}

	return data.Response[0], reqInf, nil
}
