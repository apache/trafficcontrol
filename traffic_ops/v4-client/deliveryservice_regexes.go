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
	// APIDSRegexes is the full API route to the
	// /deliveryservices/{{ID}}/regexes endpoint.
	APIDSRegexes = "/deliveryservices/%d/regexes"
)

// GetDeliveryServiceRegexesByDSID gets DeliveryServiceRegexes by a DS id
// also accepts an optional map of query parameters
func (to *Session) GetDeliveryServiceRegexesByDSID(dsID int, params url.Values) ([]tc.DeliveryServiceIDRegex, toclientlib.ReqInf, error) {
	response := struct {
		Response []tc.DeliveryServiceIDRegex `json:"response"`
	}{}
	route := fmt.Sprintf(APIDSRegexes, dsID)
	if len(params) > 0 {
		route += "?" + params.Encode()
	}
	reqInf, err := to.get(route, nil, &response)
	return response.Response, reqInf, err
}

// GetDeliveryServiceRegexes retrieves all Delivery Service Regexes in Traffic
// Ops.
func (to *Session) GetDeliveryServiceRegexes(header http.Header) ([]tc.DeliveryServiceRegexes, toclientlib.ReqInf, error) {
	var data tc.DeliveryServiceRegexResponse
	reqInf, err := to.get(APIDeliveryServicesRegexes, header, &data)
	return data.Response, reqInf, err
}

// PostDeliveryServiceRegexesByDSID adds the given Regex to the identified
// Delivery Service.
func (to *Session) PostDeliveryServiceRegexesByDSID(dsID int, regex tc.DeliveryServiceRegexPost) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	route := fmt.Sprintf(APIDSRegexes, dsID)
	reqInf, err := to.post(route, regex, nil, &alerts)
	return alerts, reqInf, err
}
