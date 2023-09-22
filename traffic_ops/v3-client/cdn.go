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
	// API_CDNS is Deprecated: will be removed in the next major version. Be aware this may not be the URI being requested, for clients created with Login and ClientOps.ForceLatestAPI false.
	API_CDNS = apiBase + "/cdns"

	APICDNs = "/cdns"
)

// CreateCDN creates a CDN.
func (to *Session) CreateCDN(cdn tc.CDN) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(APICDNs, cdn, nil, &alerts)
	return alerts, reqInf, err
}

func (to *Session) UpdateCDNByIDWithHdr(id int, cdn tc.CDN, header http.Header) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d", APICDNs, id)
	var alerts tc.Alerts
	reqInf, err := to.put(route, cdn, header, &alerts)
	return alerts, reqInf, err
}

// UpdateCDNByID updates a CDN by ID.
// Deprecated: UpdateCDNByID will be removed in 6.0. Use UpdateCDNByIDWithHdr.
func (to *Session) UpdateCDNByID(id int, cdn tc.CDN) (tc.Alerts, toclientlib.ReqInf, error) {
	return to.UpdateCDNByIDWithHdr(id, cdn, nil)
}

func (to *Session) GetCDNsWithHdr(header http.Header) ([]tc.CDN, toclientlib.ReqInf, error) {
	var data tc.CDNsResponse
	reqInf, err := to.get(APICDNs, header, &data)
	return data.Response, reqInf, err
}

// GetCDNs eturns a list of CDNs.
// Deprecated: GetCDNs will be removed in 6.0. Use GetCDNsWithHdr.
func (to *Session) GetCDNs() ([]tc.CDN, toclientlib.ReqInf, error) {
	return to.GetCDNsWithHdr(nil)
}

func (to *Session) GetCDNByIDWithHdr(id int, header http.Header) ([]tc.CDN, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?id=%v", APICDNs, id)
	var data tc.CDNsResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GetCDNByID a CDN by the CDN ID.
// Deprecated: GetCDNByID will be removed in 6.0. Use GetCDNByIDWithHdr.
func (to *Session) GetCDNByID(id int) ([]tc.CDN, toclientlib.ReqInf, error) {
	return to.GetCDNByIDWithHdr(id, nil)
}

func (to *Session) GetCDNByNameWithHdr(name string, header http.Header) ([]tc.CDN, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?name=%s", APICDNs, url.QueryEscape(name))
	var data tc.CDNsResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GetCDNByName gets a CDN by the CDN name.
// Deprecated: GetCDNByName will be removed in 6.0. Use GetCDNByNameWithHdr.
func (to *Session) GetCDNByName(name string) ([]tc.CDN, toclientlib.ReqInf, error) {
	return to.GetCDNByNameWithHdr(name, nil)
}

// DeleteCDNByID deletes a CDN by ID.
func (to *Session) DeleteCDNByID(id int) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d", APICDNs, id)
	var alerts tc.Alerts
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}

func (to *Session) GetCDNSSLKeysWithHdr(name string, header http.Header) ([]tc.CDNSSLKeys, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/name/%s/sslkeys", APICDNs, name)
	var data tc.CDNSSLKeysResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// Deprecated: GetCDNSSLKeys will be removed in 6.0. Use GetCDNSSLKeysWithHdr.
func (to *Session) GetCDNSSLKeys(name string) ([]tc.CDNSSLKeys, toclientlib.ReqInf, error) {
	return to.GetCDNSSLKeysWithHdr(name, nil)
}

// QueueUpdatesForCDN set the "updPending" field of a list of servers identified by
// 'cdnID' and any other query params (type or profile) to the value of 'queueUpdate'
func (to *Session) QueueUpdatesForCDN(cdnID int, cdnQueueUpdate tc.CDNQueueUpdateRequest) (tc.CDNQueueUpdateResponse, toclientlib.ReqInf, error) {
	var resp tc.CDNQueueUpdateResponse
	path := fmt.Sprintf("/cdns/%d/queue_update", cdnID)
	reqInf, err := to.post(path, cdnQueueUpdate, nil, &resp)
	return resp, reqInf, err
}
