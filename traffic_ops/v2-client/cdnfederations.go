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

	"github.com/apache/trafficcontrol/v6/lib/go-tc"
)

/* Internally, the CDNName is only used in the GET method. The CDNName
 * seems to primarily be used to differentiate between `/federations` and
 * `/cdns/:name/federations`. Although the behavior is odd, it is kept to
 * keep the same behavior from perl. */

func (to *Session) CreateCDNFederationByName(f tc.CDNFederation, CDNName string) (*tc.CreateCDNFederationResponse, ReqInf, error) {
	jsonReq, err := json.Marshal(f)
	if err != nil { //There is no remoteAddr for ReqInf at this point
		return nil, ReqInf{CacheHitStatus: CacheHitStatusMiss}, err
	}

	data := tc.CreateCDNFederationResponse{}
	url := fmt.Sprintf("%s/cdns/%s/federations", apiBase, CDNName)
	inf, err := makeReq(to, "POST", url, jsonReq, &data)
	return &data, inf, err
}

func (to *Session) GetCDNFederationsByName(CDNName string) (*tc.CDNFederationResponse, ReqInf, error) {
	data := tc.CDNFederationResponse{}
	url := fmt.Sprintf("%s/cdns/%s/federations", apiBase, CDNName)
	inf, err := get(to, url, &data)
	return &data, inf, err
}

func (to *Session) GetCDNFederationsByID(CDNName string, ID int) (*tc.CDNFederationResponse, ReqInf, error) {
	data := tc.CDNFederationResponse{}
	url := fmt.Sprintf("%s/cdns/%s/federations?id=%v", apiBase, CDNName, ID)
	inf, err := get(to, url, &data)
	return &data, inf, err
}

func (to *Session) UpdateCDNFederationsByID(f tc.CDNFederation, CDNName string, ID int) (*tc.UpdateCDNFederationResponse, ReqInf, error) {
	jsonReq, err := json.Marshal(f)
	if err != nil { //There is no remoteAddr for ReqInf at this point
		return nil, ReqInf{CacheHitStatus: CacheHitStatusMiss}, err
	}

	data := tc.UpdateCDNFederationResponse{}
	url := fmt.Sprintf("%s/cdns/%s/federations/%d", apiBase, CDNName, ID)
	inf, err := makeReq(to, "PUT", url, jsonReq, &data)
	return &data, inf, err
}

func (to *Session) DeleteCDNFederationByID(CDNName string, ID int) (*tc.DeleteCDNFederationResponse, ReqInf, error) {
	data := tc.DeleteCDNFederationResponse{}
	url := fmt.Sprintf("%s/cdns/%s/federations/%d", apiBase, CDNName, ID)
	inf, err := makeReq(to, "DELETE", url, nil, &data)
	return &data, inf, err
}
