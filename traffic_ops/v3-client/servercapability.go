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
	"net/http"
	"net/url"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

const (
	// API_SERVER_CAPABILITIES is Deprecated: will be removed in the next major version. Be aware this may not be the URI being requested, for clients created with Login and ClientOps.ForceLatestAPI false.
	API_SERVER_CAPABILITIES = apiBase + "/server_capabilities"

	APIServerCapabilities = "/server_capabilities"
)

// CreateServerCapability creates a server capability and returns the response.
func (to *Session) CreateServerCapability(sc tc.ServerCapability) (*tc.ServerCapabilityDetailResponse, toclientlib.ReqInf, error) {
	var scResp tc.ServerCapabilityDetailResponse
	reqInf, err := to.post(APIServerCapabilities, sc, nil, &scResp)
	if err != nil {
		return nil, reqInf, err
	}
	return &scResp, reqInf, nil
}

func (to *Session) GetServerCapabilitiesWithHdr(header http.Header) ([]tc.ServerCapability, toclientlib.ReqInf, error) {
	var data tc.ServerCapabilitiesResponse
	reqInf, err := to.get(APIServerCapabilities, header, &data)
	return data.Response, reqInf, err
}

// GetServerCapabilities returns all the server capabilities.
// Deprecated: GetServerCapabilities will be removed in 6.0. Use GetServerCapabilitiesWithHdr.
func (to *Session) GetServerCapabilities() ([]tc.ServerCapability, toclientlib.ReqInf, error) {
	return to.GetServerCapabilitiesWithHdr(nil)
}

func (to *Session) GetServerCapabilityWithHdr(name string, header http.Header) (*tc.ServerCapability, toclientlib.ReqInf, error) {
	reqUrl := fmt.Sprintf("%s?name=%s", APIServerCapabilities, url.QueryEscape(name))
	var data tc.ServerCapabilitiesResponse
	reqInf, err := to.get(reqUrl, header, &data)
	if err != nil {
		return nil, reqInf, err
	}
	if len(data.Response) == 1 {
		return &data.Response[0], reqInf, nil
	}
	return nil, reqInf, fmt.Errorf("expected one server capability in response, instead got: %+v", data.Response)
}

// GetServerCapability returns the given server capability by name.
// Deprecated: GetServerCapability will be removed in 6.0. Use GetServerCapabilityWithHdr.
func (to *Session) GetServerCapability(name string) (*tc.ServerCapability, toclientlib.ReqInf, error) {
	return to.GetServerCapabilityWithHdr(name, nil)
}

// DeleteServerCapability deletes the given server capability by name.
func (to *Session) DeleteServerCapability(name string) (tc.Alerts, toclientlib.ReqInf, error) {
	reqUrl := fmt.Sprintf("%s?name=%s", APIServerCapabilities, url.QueryEscape(name))
	var alerts tc.Alerts
	reqInf, err := to.del(reqUrl, nil, &alerts)
	return alerts, reqInf, err
}
