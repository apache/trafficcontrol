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

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

// APICapabilities is the API version-relative path for the /capabilities API endpoint.
const APICapabilities = "/capabilities"

// GetCapabilities retrieves all capabilities.
func (to *Session) GetCapabilities(header http.Header) ([]tc.Capability, toclientlib.ReqInf, error) {
	var data tc.CapabilitiesResponse
	reqInf, err := to.get(APICapabilities, header, &data)
	return data.Response, reqInf, err
}

// GetCapability retrieves only the capability named 'c'.
func (to *Session) GetCapability(c string, header http.Header) (tc.Capability, toclientlib.ReqInf, error) {
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
