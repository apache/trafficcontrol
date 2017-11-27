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

// Deprecated: use GetCDNs.
func (to *Session) CDNs() ([]tc.CDN, error) {
	cdns, _, err := to.GetCDNs()
	return cdns, err
}

func (to *Session) GetCDNs() ([]tc.CDN, ReqInf, error) {
	url := "/api/1.2/cdns.json"
	resp, remoteAddr, err := to.request(http.MethodGet, url, nil) // TODO change to getBytesWithTTL, which caches
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
// Deprecated: use GetCDNName
func (to *Session) CDNName(name string) ([]tc.CDN, error) {
	n, _, err := to.GetCDNName(name)
	return n, err
}

func (to *Session) GetCDNName(name string) ([]tc.CDN, ReqInf, error) {
	url := fmt.Sprintf("/api/1.2/cdns/name/%s.json", name)
	resp, remoteAddr, err := to.request(httpMethod.Get, url, nil) // TODO change to getBytesWithTTL, return CacheHitStatus
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
	url := fmt.Sprintf("/api/1.2/cdns/name/%s/sslkeys.json", name)
	resp, remoteAddr, err := to.request(httpMethod.Get, url, nil) // TODO change to getBytesWithTTL, return CacheHitStatus
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
