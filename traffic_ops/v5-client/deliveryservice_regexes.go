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

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

// apiDSRegexes is the full API route to the
// /deliveryservices/{{ID}}/regexes endpoint.
const apiDSRegexes = "/deliveryservices/%d/regexes"

// GetDeliveryServiceRegexesByDSID gets DeliveryServiceRegexes by a DS id
// also accepts an optional map of query parameters.
func (to *Session) GetDeliveryServiceRegexesByDSID(dsID int, opts RequestOptions) (tc.DeliveryServiceIDRegexResponse, toclientlib.ReqInf, error) {
	var response tc.DeliveryServiceIDRegexResponse
	route := fmt.Sprintf(apiDSRegexes, dsID)
	reqInf, err := to.get(route, opts, &response)
	return response, reqInf, err
}

// GetDeliveryServiceRegexes retrieves all Delivery Service Regexes in Traffic
// Ops.
func (to *Session) GetDeliveryServiceRegexes(opts RequestOptions) (tc.DeliveryServiceRegexResponse, toclientlib.ReqInf, error) {
	var data tc.DeliveryServiceRegexResponse
	reqInf, err := to.get(apiDeliveryServicesRegexes, opts, &data)
	return data, reqInf, err
}

// PostDeliveryServiceRegexesByDSID adds the given Regex to the identified
// Delivery Service.
func (to *Session) PostDeliveryServiceRegexesByDSID(dsID int, regex tc.DeliveryServiceRegexPost, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	route := fmt.Sprintf(apiDSRegexes, dsID)
	reqInf, err := to.post(route, opts, regex, &alerts)
	return alerts, reqInf, err
}
