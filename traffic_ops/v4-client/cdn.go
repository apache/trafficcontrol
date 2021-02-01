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

	"github.com/apache/trafficcontrol/lib/go-tc"
)

const (
	API_CDNS = apiBase + "/cdns"
)

// CreateCDN creates a CDN.
func (to *Session) CreateCDN(cdn tc.CDN) (tc.Alerts, ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(API_CDNS, cdn, nil, &alerts)
	return alerts, reqInf, err
}

func (to *Session) UpdateCDNByIDWithHdr(id int, cdn tc.CDN, header http.Header) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s/%d", API_CDNS, id)
	var alerts tc.Alerts
	reqInf, err := to.put(route, cdn, header, &alerts)
	return alerts, reqInf, err
}

// UpdateCDNByID updates a CDN by ID.
// Deprecated: UpdateCDNByID will be removed in 6.0. Use UpdateCDNByIDWithHdr.
func (to *Session) UpdateCDNByID(id int, cdn tc.CDN) (tc.Alerts, ReqInf, error) {
	return to.UpdateCDNByIDWithHdr(id, cdn, nil)
}

func (to *Session) GetCDNsWithHdr(header http.Header) ([]tc.CDN, ReqInf, error) {
	var data tc.CDNsResponse
	reqInf, err := to.get(API_CDNS, header, &data)
	return data.Response, reqInf, err
}

// GetCDNs eturns a list of CDNs.
// Deprecated: GetCDNs will be removed in 6.0. Use GetCDNsWithHdr.
func (to *Session) GetCDNs() ([]tc.CDN, ReqInf, error) {
	return to.GetCDNsWithHdr(nil)
}

func (to *Session) GetCDNByIDWithHdr(id int, header http.Header) ([]tc.CDN, ReqInf, error) {
	route := fmt.Sprintf("%s?id=%v", API_CDNS, id)
	var data tc.CDNsResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GetCDNByID a CDN by the CDN ID.
// Deprecated: GetCDNByID will be removed in 6.0. Use GetCDNByIDWithHdr.
func (to *Session) GetCDNByID(id int) ([]tc.CDN, ReqInf, error) {
	return to.GetCDNByIDWithHdr(id, nil)
}

func (to *Session) GetCDNByNameWithHdr(name string, header http.Header) ([]tc.CDN, ReqInf, error) {
	route := fmt.Sprintf("%s?name=%s", API_CDNS, url.QueryEscape(name))
	var data tc.CDNsResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GetCDNByName gets a CDN by the CDN name.
// Deprecated: GetCDNByName will be removed in 6.0. Use GetCDNByNameWithHdr.
func (to *Session) GetCDNByName(name string) ([]tc.CDN, ReqInf, error) {
	return to.GetCDNByNameWithHdr(name, nil)
}

// DeleteCDNByID deletes a CDN by ID.
func (to *Session) DeleteCDNByID(id int) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s/%d", API_CDNS, id)
	var alerts tc.Alerts
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}

func (to *Session) GetCDNSSLKeysWithHdr(name string, header http.Header) ([]tc.CDNSSLKeys, ReqInf, error) {
	route := fmt.Sprintf("%s/name/%s/sslkeys", API_CDNS, name)
	var data tc.CDNSSLKeysResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// Deprecated: GetCDNSSLKeys will be removed in 6.0. Use GetCDNSSLKeysWithHdr.
func (to *Session) GetCDNSSLKeys(name string) ([]tc.CDNSSLKeys, ReqInf, error) {
	return to.GetCDNSSLKeysWithHdr(name, nil)
}
