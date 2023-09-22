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

/* Internally, the CDNName is only used in the GET method. The CDNName
 * seems to primarily be used to differentiate between `/federations` and
 * `/cdns/:name/federations`. Although the behavior is odd, it is kept to
 * keep the same behavior from perl. */

func (to *Session) CreateCDNFederationByName(f tc.CDNFederation, CDNName string) (*tc.CreateCDNFederationResponse, toclientlib.ReqInf, error) {
	data := tc.CreateCDNFederationResponse{}
	route := fmt.Sprintf("/cdns/%s/federations", url.QueryEscape(CDNName))
	inf, err := to.post(route, f, nil, &data)
	return &data, inf, err
}

func (to *Session) GetCDNFederationsByNameWithHdr(CDNName string, header http.Header) (*tc.CDNFederationResponse, toclientlib.ReqInf, error) {
	data := tc.CDNFederationResponse{}
	route := fmt.Sprintf("/cdns/%s/federations", url.QueryEscape(CDNName))
	inf, err := to.get(route, header, &data)
	return &data, inf, err
}

// Deprecated: GetCDNFederationsByName will be removed in 6.0. Use GetCDNFederationsByNameWithHdr.
func (to *Session) GetCDNFederationsByName(CDNName string) (*tc.CDNFederationResponse, toclientlib.ReqInf, error) {
	return to.GetCDNFederationsByNameWithHdr(CDNName, nil)
}

func (to *Session) GetCDNFederationsByNameWithHdrReturnList(CDNName string, header http.Header) ([]tc.CDNFederation, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("/cdns/%s/federations", url.QueryEscape(CDNName))
	resp := struct {
		Response []tc.CDNFederation `json:"response"`
	}{}
	inf, err := to.get(route, header, &resp)
	return resp.Response, inf, err
}

func (to *Session) GetCDNFederationsByIDWithHdr(CDNName string, ID int, header http.Header) (*tc.CDNFederationResponse, toclientlib.ReqInf, error) {
	data := tc.CDNFederationResponse{}
	route := fmt.Sprintf("/cdns/%s/federations?id=%v", url.QueryEscape(CDNName), ID)
	inf, err := to.get(route, header, &data)
	return &data, inf, err
}

// Deprecated: GetCDNFederationsByID will be removed in 6.0. Use GetCDNFederationsByIDWithHdr.
func (to *Session) GetCDNFederationsByID(CDNName string, ID int) (*tc.CDNFederationResponse, toclientlib.ReqInf, error) {
	return to.GetCDNFederationsByIDWithHdr(CDNName, ID, nil)
}

func (to *Session) UpdateCDNFederationsByIDWithHdr(f tc.CDNFederation, CDNName string, ID int, h http.Header) (*tc.UpdateCDNFederationResponse, toclientlib.ReqInf, error) {
	data := tc.UpdateCDNFederationResponse{}
	route := fmt.Sprintf("/cdns/%s/federations/%d", url.QueryEscape(CDNName), ID)
	inf, err := to.put(route, f, h, &data)
	return &data, inf, err
}

// Deprecated: UpdateCDNFederationsByID will be removed in 6.0. Use UpdateCDNFederationsByIDWithHdr.
func (to *Session) UpdateCDNFederationsByID(f tc.CDNFederation, CDNName string, ID int) (*tc.UpdateCDNFederationResponse, toclientlib.ReqInf, error) {
	return to.UpdateCDNFederationsByIDWithHdr(f, CDNName, ID, nil)
}

func (to *Session) DeleteCDNFederationByID(CDNName string, ID int) (*tc.DeleteCDNFederationResponse, toclientlib.ReqInf, error) {
	data := tc.DeleteCDNFederationResponse{}
	route := fmt.Sprintf("/cdns/%s/federations/%d", url.QueryEscape(CDNName), ID)
	inf, err := to.del(route, nil, &data)
	return &data, inf, err
}
