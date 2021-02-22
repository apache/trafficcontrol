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
	API_ASNS = apiBase + "/asns"
)

// CreateASN creates a ASN
func (to *Session) CreateASN(entity tc.ASN) (tc.Alerts, ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(API_ASNS, entity, nil, &alerts)
	return alerts, reqInf, err
}

// UpdateASNByID updates a ASN by ID
func (to *Session) UpdateASNByID(id int, entity tc.ASN) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", API_ASNS, id)
	var alerts tc.Alerts
	reqInf, err := to.put(route, entity, nil, &alerts)
	return alerts, reqInf, err
}

// GetASNsWithHeader Returns a list of ASNs matching query params
func (to *Session) GetASNsWithHeader(params *url.Values, header http.Header) ([]tc.ASN, ReqInf, error) {
	route := fmt.Sprintf("%s?%s", API_ASNS, params.Encode())
	var data tc.ASNsResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// DeleteASNByASN deletes an ASN by asn number
func (to *Session) DeleteASNByASN(asn int) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", API_ASNS, asn)
	var alerts tc.Alerts
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}
