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
	"strconv"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

const (
	// API_DELIVERY_SERVICES_REQUIRED_CAPABILITIES is Deprecated: will be removed in the next major version. Be aware this may not be the URI being requested, for clients created with Login and ClientOps.ForceLatestAPI false.
	API_DELIVERY_SERVICES_REQUIRED_CAPABILITIES = apiBase + "/deliveryservices_required_capabilities"

	APIDeliveryServicesRequiredCapabilities = "/deliveryservices_required_capabilities"
)

// CreateDeliveryServicesRequiredCapability assigns a Required Capability to a Delivery Service
func (to *Session) CreateDeliveryServicesRequiredCapability(capability tc.DeliveryServicesRequiredCapability) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(APIDeliveryServicesRequiredCapabilities, capability, nil, &alerts)
	return alerts, reqInf, err
}

// DeleteDeliveryServicesRequiredCapability unassigns a Required Capability from a Delivery Service
func (to *Session) DeleteDeliveryServicesRequiredCapability(deliveryserviceID int, capability string) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	param := url.Values{}
	param.Add("deliveryServiceID", strconv.Itoa(deliveryserviceID))
	param.Add("requiredCapability", capability)
	route := fmt.Sprintf("%s?%s", APIDeliveryServicesRequiredCapabilities, param.Encode())
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}

func (to *Session) GetDeliveryServicesRequiredCapabilitiesWithHdr(deliveryServiceID *int, xmlID, capability *string, header http.Header) ([]tc.DeliveryServicesRequiredCapability, toclientlib.ReqInf, error) {
	param := url.Values{}
	if deliveryServiceID != nil {
		param.Add("deliveryServiceID", strconv.Itoa(*deliveryServiceID))
	}
	if xmlID != nil {
		param.Add("xmlID", *xmlID)
	}
	if capability != nil {
		param.Add("requiredCapability", *capability)
	}

	route := APIDeliveryServicesRequiredCapabilities
	if len(param) > 0 {
		route = fmt.Sprintf("%s?%s", route, param.Encode())
	}

	resp := struct {
		Response []tc.DeliveryServicesRequiredCapability `json:"response"`
	}{}
	reqInf, err := to.get(route, header, &resp)
	return resp.Response, reqInf, err
}

// GetDeliveryServicesRequiredCapabilities retrieves a list of Required Capabilities that are assigned to a Delivery Service
// Callers can filter the results by delivery service id, xml id and/or required capability via the optional parameters
// Deprecated: GetDeliveryServicesRequiredCapabilities will be removed in 6.0. Use GetDeliveryServicesRequiredCapabilitiesWithHdr.
func (to *Session) GetDeliveryServicesRequiredCapabilities(deliveryServiceID *int, xmlID, capability *string) ([]tc.DeliveryServicesRequiredCapability, toclientlib.ReqInf, error) {
	return to.GetDeliveryServicesRequiredCapabilitiesWithHdr(deliveryServiceID, xmlID, capability, nil)
}
