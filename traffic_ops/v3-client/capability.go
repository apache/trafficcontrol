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

import "errors"
import "net/http"
import "net/url"

import "github.com/apache/trafficcontrol/v8/lib/go-tc"
import "github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"

// API_CAPABILITIES is Deprecated: will be removed in the next major version. Be aware this may not be the URI being requested, for clients created with Login and ClientOps.ForceLatestAPI false.
const API_CAPABILITIES = apiBase + "/capabilities"

const APICapabilities = "/capabilities"

func (to *Session) GetCapabilitiesWithHdr(header http.Header) ([]tc.Capability, toclientlib.ReqInf, error) {
	var data tc.CapabilitiesResponse
	reqInf, err := to.get(APICapabilities, header, &data)
	return data.Response, reqInf, err
}

// GetCapabilities retrieves all capabilities.
// Deprecated: GetCapabilities will be removed in 6.0. Use GetCapabilitiesWithHdr.
func (to *Session) GetCapabilities() ([]tc.Capability, toclientlib.ReqInf, error) {
	return to.GetCapabilitiesWithHdr(nil)
}

func (to *Session) GetCapabilityWithHdr(c string, header http.Header) (tc.Capability, toclientlib.ReqInf, error) {
	v := url.Values{}
	v.Add("name", c)
	endpoint := APICapabilities + "?" + v.Encode()
	var data tc.CapabilitiesResponse
	reqInf, err := to.get(endpoint, header, &data)
	if err != nil {
		return tc.Capability{}, reqInf, err
	} else if data.Response == nil || len(data.Response) < 1 {
		return tc.Capability{}, reqInf, errors.New("invalid response - no capability returned")
	}

	return data.Response[0], reqInf, nil
}

// GetCapability retrieves only the capability named 'c'.
// Deprecated: GetCapability will be removed in 6.0. Use GetCapabilityWithHdr.
func (to *Session) GetCapability(c string) (tc.Capability, toclientlib.ReqInf, error) {
	return to.GetCapabilityWithHdr(c, nil)
}
