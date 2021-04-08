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

/* Internally, the CDNName is only used in the GET method. The CDNName
 * seems to primarily be used to differentiate between `/federations` and
 * `/cdns/:name/federations`. Although the behavior is odd, it is kept to
 * keep the same behavior from perl. */

// CreateCDNFederationByName creates the given Federation in the CDN with the
// given name.
func (to *Session) CreateCDNFederation(f tc.CDNFederation, CDNName string) (tc.CreateCDNFederationResponse, toclientlib.ReqInf, error) {
	data := tc.CreateCDNFederationResponse{}
	route := "/cdns/" + url.QueryEscape(CDNName) + "/federations"
	inf, err := to.post(route, f, nil, &data)
	return data, inf, err
}

// GetCDNFederationsByName retrieves all Federations in the CDN with the given
// name.
func (to *Session) GetCDNFederationsByName(CDNName string, header http.Header) (tc.CDNFederationResponse, toclientlib.ReqInf, error) {
	data := tc.CDNFederationResponse{}
	route := "/cdns/" + url.QueryEscape(CDNName) + "/federations"
	inf, err := to.get(route, header, &data)
	return data, inf, err
}

// GetCDNFederationsByID retrieves the Federation in the CDN with the given
// name that has the given ID.
func (to *Session) GetCDNFederationsByID(CDNName string, ID int, header http.Header) (tc.CDNFederationResponse, toclientlib.ReqInf, error) {
	data := tc.CDNFederationResponse{}
	route := fmt.Sprintf("/cdns/%s/federations?id=%v", url.QueryEscape(CDNName), ID)
	inf, err := to.get(route, header, &data)
	return data, inf, err
}

// UpdateCDNFederation replaces the Federation with the given ID in the CDN
// with the given name with the provided Federation.
func (to *Session) UpdateCDNFederation(f tc.CDNFederation, CDNName string, ID int, h http.Header) (tc.UpdateCDNFederationResponse, toclientlib.ReqInf, error) {
	data := tc.UpdateCDNFederationResponse{}
	route := fmt.Sprintf("/cdns/%s/federations/%d", url.QueryEscape(CDNName), ID)
	inf, err := to.put(route, f, h, &data)
	return data, inf, err
}

// DeleteCDNFederationByID deletes the Federation with the given ID in the CDN
// with the given name.
func (to *Session) DeleteCDNFederation(CDNName string, ID int) (tc.DeleteCDNFederationResponse, toclientlib.ReqInf, error) {
	data := tc.DeleteCDNFederationResponse{}
	route := fmt.Sprintf("/cdns/%s/federations/%d", url.QueryEscape(CDNName), ID)
	inf, err := to.del(route, nil, &data)
	return data, inf, err
}
