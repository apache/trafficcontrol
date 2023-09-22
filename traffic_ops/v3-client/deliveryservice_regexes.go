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

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

const (
	// API_DS_REGEXES is Deprecated: will be removed in the next major version. Be aware this may not be the URI being requested, for clients created with Login and ClientOps.ForceLatestAPI false.
	// See: https://traffic-control-cdn.readthedocs.io/en/latest/api/v3/deliveryservices_id_regexes.html
	API_DS_REGEXES = apiBase + "/deliveryservices/%d/regexes"

	APIDSRegexes = "/deliveryservices/%d/regexes"
)

// GetDeliveryServiceRegexesByDSID gets DeliveryServiceRegexes by a DS id
// also accepts an optional map of query parameters
func (to *Session) GetDeliveryServiceRegexesByDSID(dsID int, params map[string]string) ([]tc.DeliveryServiceIDRegex, toclientlib.ReqInf, error) {
	response := struct {
		Response []tc.DeliveryServiceIDRegex `json:"response"`
	}{}
	reqInf, err := to.get(fmt.Sprintf(APIDSRegexes, dsID)+mapToQueryParameters(params), nil, &response)
	return response.Response, reqInf, err
}

// GetDeliveryServiceRegexes returns the "Regexes" (Regular Expressions) used by all (tenant-visible)
// Delivery Services.
// Deprecated: GetDeliveryServiceRegexes will be removed in 6.0. Use GetDeliveryServiceRegexesWithHdr.
func (to *Session) GetDeliveryServiceRegexes() ([]tc.DeliveryServiceRegexes, toclientlib.ReqInf, error) {
	return to.GetDeliveryServiceRegexesWithHdr(nil)
}

func (to *Session) GetDeliveryServiceRegexesWithHdr(header http.Header) ([]tc.DeliveryServiceRegexes, toclientlib.ReqInf, error) {
	var data tc.DeliveryServiceRegexResponse
	reqInf, err := to.get(APIDeliveryServicesRegexes, header, &data)
	return data.Response, reqInf, err
}

func (to *Session) PostDeliveryServiceRegexesByDSID(dsID int, regex tc.DeliveryServiceRegexPost) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	route := fmt.Sprintf(APIDSRegexes, dsID)
	reqInf, err := to.post(route, regex, nil, &alerts)
	return alerts, reqInf, err
}
