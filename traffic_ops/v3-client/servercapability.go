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
	"net"
	"net/http"
	"net/url"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

const (
	API_SERVER_CAPABILITIES = apiBase + "/server_capabilities"
)

// CreateServerCapability creates a server capability and returns the response.
func (to *Session) CreateServerCapability(sc tc.ServerCapability) (*tc.ServerCapabilityDetailResponse, ReqInf, error) {
	var remoteAddr net.Addr
	reqBody, err := json.Marshal(sc)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	resp, remoteAddr, err := to.request(http.MethodPost, API_SERVER_CAPABILITIES, reqBody, nil)
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()
	var scResp tc.ServerCapabilityDetailResponse
	if err = json.NewDecoder(resp.Body).Decode(&scResp); err != nil {
		return nil, reqInf, err
	}
	return &scResp, reqInf, nil
}

func (to *Session) GetServerCapabilitiesWithHdr(header http.Header) ([]tc.ServerCapability, ReqInf, error) {
	resp, remoteAddr, err := to.request(http.MethodGet, API_SERVER_CAPABILITIES, nil, header)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		if reqInf.StatusCode == http.StatusNotModified {
			return []tc.ServerCapability{}, reqInf, nil
		}
	}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.ServerCapabilitiesResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// GetServerCapabilities returns all the server capabilities.
// Deprecated: GetServerCapabilities will be removed in 6.0. Use GetServerCapabilitiesWithHdr.
func (to *Session) GetServerCapabilities() ([]tc.ServerCapability, ReqInf, error) {
	return to.GetServerCapabilitiesWithHdr(nil)
}

func (to *Session) GetServerCapabilityWithHdr(name string, header http.Header) (*tc.ServerCapability, ReqInf, error) {
	reqUrl := fmt.Sprintf("%s?name=%s", API_SERVER_CAPABILITIES, url.QueryEscape(name))
	resp, remoteAddr, err := to.request(http.MethodGet, reqUrl, nil, header)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		if reqInf.StatusCode == http.StatusNotModified {
			return nil, reqInf, nil
		}
	}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.ServerCapabilitiesResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	if len(data.Response) == 1 {
		return &data.Response[0], reqInf, nil
	}
	return nil, reqInf, fmt.Errorf("expected one server capability in response, instead got: %+v", data.Response)
}

// GetServerCapability returns the given server capability by name.
// Deprecated: GetServerCapability will be removed in 6.0. Use GetServerCapabilityWithHdr.
func (to *Session) GetServerCapability(name string) (*tc.ServerCapability, ReqInf, error) {
	return to.GetServerCapabilityWithHdr(name, nil)
}

// DeleteServerCapability deletes the given server capability by name.
func (to *Session) DeleteServerCapability(name string) (tc.Alerts, ReqInf, error) {
	reqUrl := fmt.Sprintf("%s?name=%s", API_SERVER_CAPABILITIES, url.QueryEscape(name))
	resp, remoteAddr, err := to.request(http.MethodDelete, reqUrl, nil, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	if err = json.NewDecoder(resp.Body).Decode(&alerts); err != nil {
		return tc.Alerts{}, reqInf, err
	}
	return alerts, reqInf, nil
}
