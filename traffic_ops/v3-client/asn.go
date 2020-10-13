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
	API_ASNS = apiBase + "/asns"
)

// CreateASN creates a ASN
func (to *Session) CreateASN(entity tc.ASN) (tc.Alerts, ReqInf, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(entity)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	resp, remoteAddr, err := to.request(http.MethodPost, API_ASNS, reqBody, nil)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

// UpdateASNByID updates a ASN by ID
func (to *Session) UpdateASNByID(id int, entity tc.ASN) (tc.Alerts, ReqInf, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(entity)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	route := fmt.Sprintf("%s?id=%d", API_ASNS, id)
	resp, remoteAddr, err := to.request(http.MethodPut, route, reqBody, nil)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

// GetASNsWithHeader Returns a list of ASNs matching query params
func (to *Session) GetASNsWithHeader(params *url.Values, header http.Header) ([]tc.ASN, ReqInf, error) {
	route := fmt.Sprintf("%s?%s", API_ASNS, params.Encode())
	resp, remoteAddr, err := to.request(http.MethodGet, route, nil, header)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}

	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		if reqInf.StatusCode == http.StatusNotModified {
			return []tc.ASN{}, reqInf, nil
		}
	}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.ASNsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return data.Response, reqInf, err
	}
	return data.Response, reqInf, nil
}

// DeleteASNByASN deletes an ASN by asn number
func (to *Session) DeleteASNByASN(asn int) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", API_ASNS, asn)
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
