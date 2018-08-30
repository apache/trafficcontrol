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

package v13

import (
	"encoding/json"
	"fmt"

	"github.com/apache/trafficcontrol/lib/go-tc/v13"
)

/* Internally, the CDNName is only used in the GET method. The CDNName
 * seems to primarily be used to differentiate between `/federations` and
 * `/cdns/:name/federations`. Although the behavior is odd, it is kept to
 * keep the same behavior from perl. */

func (to *Session) CreateCDNFederationByName(f v13.CDNFederation, CDNName string) (*v13.CreateCDNFederationResponse, ReqInf, error) {
	jsonReq, err := json.Marshal(f)
	if err != nil { //There is no remoteAddr for ReqInf at this point
		return nil, ReqInf{CacheHitStatus: CacheHitStatusMiss}, err
	}

	data := v13.CreateCDNFederationResponse{}
	url := fmt.Sprintf("%s/cdns/%s/federations", apiBase, CDNName)
	inf, err := makeReq(to, "POST", url, jsonReq, &data)
	return &data, inf, err
}

func (to *Session) GetCDNFederationsByName(CDNName string) (*v13.CDNFederationResponse, ReqInf, error) {
	data := v13.CDNFederationResponse{}
	url := fmt.Sprintf("%s/cdns/%s/federations", apiBase, CDNName)
	inf, err := get(to, url, &data)
	return &data, inf, err
}

func (to *Session) GetCDNFederationsByID(CDNName string, ID int) (*v13.CDNFederationResponse, ReqInf, error) {
	data := v13.CDNFederationResponse{}
	url := fmt.Sprintf("%s/cdns/%s/federations/%d", apiBase, CDNName, ID)
	inf, err := get(to, url, &data)
	return &data, inf, err
}

func (to *Session) UpdateCDNFederationsByID(f v13.CDNFederation, CDNName string, ID int) (*v13.UpdateCDNFederationResponse, ReqInf, error) {
	jsonReq, err := json.Marshal(f)
	if err != nil { //There is no remoteAddr for ReqInf at this point
		return nil, ReqInf{CacheHitStatus: CacheHitStatusMiss}, err
	}

	data := v13.UpdateCDNFederationResponse{}
	url := fmt.Sprintf("%s/cdns/%s/federations/%d", apiBase, CDNName, ID)
	inf, err := makeReq(to, "PUT", url, jsonReq, &data)
	return &data, inf, err
}

func (to *Session) DeleteCDNFederationByID(CDNName string, ID int) (*v13.DeleteCDNFederationResponse, ReqInf, error) {
	data := v13.DeleteCDNFederationResponse{}
	url := fmt.Sprintf("%s/cdns/%s/federations/%d", apiBase, CDNName, ID)
	inf, err := makeReq(to, "DELETE", url, nil, &data)
	return &data, inf, err
}
