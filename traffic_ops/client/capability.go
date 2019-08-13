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
import "fmt"
import "net"
import "net/http"

import "github.com/apache/trafficcontrol/lib/go-tc"

const API_v4_CAPABILITIES = "/api/1.4/capabilities"

// CreateCapability creates the passed capability.
func (to *Session) CreateCapability(c tc.Capability) (tc.Alerts, ReqInf, error) {
	var remoteAddr net.Addr
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}

	reqBody, err := json.Marshal(c)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}

	resp, remoteAddr, err := to.request(http.MethodPost, API_v4_CAPABILITIES, reqBody)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	reqInf.RemoteAddr = remoteAddr

	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, err
}

// ReplaceCapabilityByName replaces the capability named 'name' with the one passed as 'c'.
func (to *Session) ReplaceCapabilityByName(name string, c tc.Capability) (tc.Alerts, ReqInf, error) {
	var remoteAddr net.Addr
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}

	reqBody, err := json.Marshal(c)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}

	endpoint := fmt.Sprintf("%s?name=%s", API_v4_CAPABILITIES, name)

	var resp *http.Response
	resp, reqInf.RemoteAddr, err = to.request(http.MethodPut, endpoint, reqBody)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()

	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, err
}

// GetCapabilities retrieves all capabilities.
func (to *Session) GetCapabilities() ([]tc.Capability, ReqInf, error) {
	resp, remoteAddr, err := to.request(http.MethodGet, API_v4_CAPABILITIES, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.CapabilitiesResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	return data.Response, reqInf, err
}

// GetCapability retrieves only the capability named 'c'
func (to *Session) GetCapability(c string) (tc.Capability, ReqInf, error) {
	endpoint := fmt.Sprintf("%s?name=%s", API_v4_CAPABILITIES, c)
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
		return tc.Capability{}, reqInf, fmt.Errorf("Invalid response - no capability returned!")
	}

	return data.Response[0], reqInf, nil
}

// DeleteCapability deletes the capability named 'c'.
func (to *Session) DeleteCapability(c string) (alerts tc.Alerts, reqInf ReqInf, err error) {
	reqInf.CacheHitStatus = CacheHitStatusMiss
	endpoint := fmt.Sprintf("%s?name=%s", API_v4_CAPABILITIES, c)

	var resp *http.Response
	resp, reqInf.RemoteAddr, err = to.request(http.MethodDelete, endpoint, nil)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return
}
