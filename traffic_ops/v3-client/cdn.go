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
	API_CDNS = apiBase + "/cdns"
)

// CreateCDN creates a CDN.
func (to *Session) CreateCDN(cdn tc.CDN) (tc.Alerts, ReqInf, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(cdn)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	resp, remoteAddr, err := to.request(http.MethodPost, API_CDNS, reqBody, nil)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

// UpdateCDNByID updates a CDN by ID.
func (to *Session) UpdateCDNByID(id int, cdn tc.CDN) (tc.Alerts, ReqInf, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(cdn)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	route := fmt.Sprintf("%s/%d", API_CDNS, id)
	resp, remoteAddr, err := to.request(http.MethodPut, route, reqBody, nil)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

func (to *Session) GetCDNsWithHdr(header http.Header) ([]tc.CDN, ReqInf, error) {
	resp, remoteAddr, err := to.request(http.MethodGet, API_CDNS, nil, header)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		if reqInf.StatusCode == http.StatusNotModified {
			return []tc.CDN{}, reqInf, nil
		}
	}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.CDNsResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	return data.Response, reqInf, nil
}

// GetCDNs eturns a list of CDNs.
// Deprecated: GetCDNs will be removed in 6.0. Use GetCDNsWithHdr.
func (to *Session) GetCDNs() ([]tc.CDN, ReqInf, error) {
	return to.GetCDNsWithHdr(nil)
}

func (to *Session) GetCDNByIDWithHdr(id int, header http.Header) ([]tc.CDN, ReqInf, error) {
	route := fmt.Sprintf("%s?id=%v", API_CDNS, id)
	resp, remoteAddr, err := to.request(http.MethodGet, route, nil, header)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		if reqInf.StatusCode == http.StatusNotModified {
			return []tc.CDN{}, reqInf, nil
		}
	}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.CDNsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// GetCDNByID a CDN by the CDN ID.
// Deprecated: GetCDNByID will be removed in 6.0. Use GetCDNByIDWithHdr.
func (to *Session) GetCDNByID(id int) ([]tc.CDN, ReqInf, error) {
	return to.GetCDNByIDWithHdr(id, nil)
}

func (to *Session) GetCDNByNameWithHdr(name string, header http.Header) ([]tc.CDN, ReqInf, error) {
	url := fmt.Sprintf("%s?name=%s", API_CDNS, url.QueryEscape(name))
	resp, remoteAddr, err := to.request(http.MethodGet, url, nil, header)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		if reqInf.StatusCode == http.StatusNotModified {
			return []tc.CDN{}, reqInf, nil
		}
	}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.CDNsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// GetCDNByName gets a CDN by the CDN name.
// Deprecated: GetCDNByName will be removed in 6.0. Use GetCDNByNameWithHdr.
func (to *Session) GetCDNByName(name string) ([]tc.CDN, ReqInf, error) {
	return to.GetCDNByNameWithHdr(name, nil)
}

// DeleteCDNByID deletes a CDN by ID.
func (to *Session) DeleteCDNByID(id int) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s/%d", API_CDNS, id)
	resp, remoteAddr, err := to.request(http.MethodDelete, route, nil, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

func (to *Session) GetCDNSSLKeysWithHdr(name string, header http.Header) ([]tc.CDNSSLKeys, ReqInf, error) {
	url := fmt.Sprintf("%s/name/%s/sslkeys", API_CDNS, name)
	resp, remoteAddr, err := to.request(http.MethodGet, url, nil, header)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		if reqInf.StatusCode == http.StatusNotModified {
			return []tc.CDNSSLKeys{}, reqInf, nil
		}
	}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.CDNSSLKeysResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// Deprecated: GetCDNSSLKeys will be removed in 6.0. Use GetCDNSSLKeysWithHdr.
func (to *Session) GetCDNSSLKeys(name string) ([]tc.CDNSSLKeys, ReqInf, error) {
	return to.GetCDNSSLKeysWithHdr(name, nil)
}
