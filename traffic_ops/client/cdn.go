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
	"net/http"

	tc "github.com/apache/incubator-trafficcontrol/lib/go-tc"
)

var (
	uri = "/api/1.2/cdns"
)

// Deprecated: use GetCDNs.
func (to *Session) CDNs() ([]tc.CDN, error) {
	cdns, _, err := to.GetCDNs()
	return cdns, err
}

func (to *Session) GetCDNs() ([]tc.CDN, ReqInf, error) {
	resp, remoteAddr, err := to.request(http.MethodGet, uri, nil) // TODO change to getBytesWithTTL, which caches
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
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

// CDNName gets an array of CDNs
// Deprecated: use GetCDNByName
func (to *Session) CDNName(name string) ([]tc.CDN, error) {
	n, _, err := to.GetCDNByName(name)
	return n, err
}

// CDNName gets an array of CDNs
// Deprecated: use GetCDNByName
func (to *Session) GetCDNName(name string) ([]tc.CDN, error) {
	n, _, err := to.GetCDNByName(name)
	return n, err
}

// GetCDNByName gets an array of CDNs
func (to *Session) GetCDNByName(name string) ([]tc.CDN, ReqInf, error) {
	url := fmt.Sprintf("%s/name/%s", uri, name)
	resp, remoteAddr, err := to.request(http.MethodGet, url, nil) // TODO change to getBytesWithTTL, return CacheHitStatus
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
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

// Deprecated: use GetCDNSSLKeys
func (to *Session) CDNSSLKeys(name string) ([]tc.CDNSSLKeys, error) {
	ks, _, err := to.GetCDNSSLKeys(name)
	return ks, err
}

func (to *Session) GetCDNSSLKeys(name string) ([]tc.CDNSSLKeys, ReqInf, error) {
	url := fmt.Sprintf("%s/name/%s/sslkeys", uri, name)
	resp, remoteAddr, err := to.request(http.MethodGet, url, nil) // TODO change to getBytesWithTTL, return CacheHitStatus
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
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

func (to *Session) DeleteCDNByName(name string) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s/name/%s", name)
	resp, remoteAddr, err := to.request(http.MethodDelete, route, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	if err := json.NewDecoder(resp.Body).Decode(&alerts); err != nil {
		return tc.Alerts{}, reqInf, err
	}
	return alerts, reqInf, nil
}
