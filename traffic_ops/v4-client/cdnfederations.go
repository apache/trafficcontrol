package client

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

import (
	"fmt"
	"net/url"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

/* Internally, the CDNName is only used in the GET method. The CDNName
 * seems to primarily be used to differentiate between `/federations` and
 * `/cdns/:name/federations`. Although the behavior is odd, it is kept to
 * keep the same behavior from perl. */

// CreateCDNFederation creates the given Federation in the CDN with the given
// name.
func (to *Session) CreateCDNFederation(f tc.CDNFederation, cdnName string, opts RequestOptions) (tc.CreateCDNFederationResponse, toclientlib.ReqInf, error) {
	var data tc.CreateCDNFederationResponse
	route := "/cdns/" + url.PathEscape(cdnName) + "/federations"
	inf, err := to.post(route, opts, f, &data)
	return data, inf, err
}

// GetCDNFederationsByName retrieves all Federations in the CDN with the given
// name.
func (to *Session) GetCDNFederationsByName(cdnName string, opts RequestOptions) (tc.CDNFederationResponse, toclientlib.ReqInf, error) {
	var data tc.CDNFederationResponse
	route := "/cdns/" + url.PathEscape(cdnName) + "/federations"
	inf, err := to.get(route, opts, &data)
	return data, inf, err
}

// UpdateCDNFederation replaces the Federation with the given ID in the CDN
// with the given name with the provided Federation.
func (to *Session) UpdateCDNFederation(f tc.CDNFederation, cdnName string, id int, opts RequestOptions) (tc.UpdateCDNFederationResponse, toclientlib.ReqInf, error) {
	var data tc.UpdateCDNFederationResponse
	route := fmt.Sprintf("/cdns/%s/federations/%d", url.PathEscape(cdnName), id)
	inf, err := to.put(route, opts, f, &data)
	return data, inf, err
}

// DeleteCDNFederation deletes the Federation with the given ID in the CDN
// with the given name.
func (to *Session) DeleteCDNFederation(cdnName string, id int, opts RequestOptions) (tc.DeleteCDNFederationResponse, toclientlib.ReqInf, error) {
	var data tc.DeleteCDNFederationResponse
	route := fmt.Sprintf("/cdns/%s/federations/%d", url.PathEscape(cdnName), id)
	inf, err := to.del(route, opts, &data)
	return data, inf, err
}
