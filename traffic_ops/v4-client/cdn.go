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
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

const (
	// APICDNs is the API version-relative path for the /cdns API endpoint.
	APICDNs = "/cdns"
)

// CreateCDN creates a CDN.
func (to *Session) CreateCDN(cdn tc.CDN) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(APICDNs, cdn, nil, &alerts)
	return alerts, reqInf, err
}

// UpdateCDNByID replaces the identified CDN with the provided CDN.
func (to *Session) UpdateCDN(id int, cdn tc.CDN, header http.Header) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d", APICDNs, id)
	var alerts tc.Alerts
	reqInf, err := to.put(route, cdn, header, &alerts)
	return alerts, reqInf, err
}

// GetCDNs retrieves all CDNs in Traffic Ops.
func (to *Session) GetCDNs(header http.Header) ([]tc.CDN, toclientlib.ReqInf, error) {
	var data tc.CDNsResponse
	reqInf, err := to.get(APICDNs, header, &data)
	return data.Response, reqInf, err
}

// GetCDNByID retrieves the CDN with the given ID.
func (to *Session) GetCDNByID(id int, header http.Header) ([]tc.CDN, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?id=%v", APICDNs, id)
	var data tc.CDNsResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GetCDNByName retrieves the CDN with the given Name.
func (to *Session) GetCDNByName(name string, header http.Header) ([]tc.CDN, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?name=%s", APICDNs, url.QueryEscape(name))
	var data tc.CDNsResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// DeleteCDNByID deletes the CDN with the given ID.
func (to *Session) DeleteCDN(id int) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d", APICDNs, id)
	var alerts tc.Alerts
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}

// GetCDNSSLKeys retrieves the SSL keys for the CDN with the given name.
func (to *Session) GetCDNSSLKeys(name string, header http.Header) ([]tc.CDNSSLKeys, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/name/%s/sslkeys", APICDNs, name)
	var data tc.CDNSSLKeysResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}
