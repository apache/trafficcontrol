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
	"net"
	"net/url"
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

const (
	v14DeliveryServiceServerCapabilities = apiBase + "/deliveryservice_server_capabilities"
)

// CreateDeliveryServiceServerCapability assigns a Server Capability to a Delivery Service
func (to *Session) CreateDeliveryServiceServerCapability(capability tc.DeliveryServiceServerCapability) (tc.Alerts, ReqInf, error) {
	var alerts tc.Alerts
	var remoteAddr net.Addr
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	reqBody, err := json.Marshal(capability)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	reqInf, err = post(to, v14DeliveryServiceServerCapabilities, reqBody, &alerts)
	return alerts, reqInf, err
}

// DeleteDeliveryServiceServerCapability unassigns a Server Capability from a Delivery Service
func (to *Session) DeleteDeliveryServiceServerCapability(deliveryserviceID int, serverCapability string) (tc.Alerts, ReqInf, error) {
	var alerts tc.Alerts
	endpoint := fmt.Sprintf("%v?deliveryServiceID=%v&serverCapability=%v", v14DeliveryServiceServerCapabilities, deliveryserviceID, serverCapability)
	reqInf, err := del(to, endpoint, &alerts)
	return alerts, reqInf, err
}

// GetDeliveryServiceServerCapabilities retrieves a list of Server Capabilities that are assigned to a Delivery Service
// Callers can filter the results by delivery service id, xml id and/or server capability via the optional parameters
func (to *Session) GetDeliveryServiceServerCapabilities(deliveryServiceID *int, xmlID, serverCapability *string) ([]tc.DeliveryServiceServerCapability, ReqInf, error) {
	v := url.Values{}
	if deliveryServiceID != nil {
		v.Add("deliveryServiceID", strconv.Itoa(*deliveryServiceID))
	}
	if xmlID != nil {
		v.Add("xmlID", *xmlID)
	}
	if serverCapability != nil {
		v.Add("serverCapability", *serverCapability)
	}
	queryURL := v14DeliveryServiceServerCapabilities
	if qStr := v.Encode(); len(qStr) > 0 {
		queryURL = fmt.Sprintf("%s?%s", queryURL, qStr)
	}

	resp := struct {
		Response []tc.DeliveryServiceServerCapability `json:"response"`
	}{}
	reqInf, err := get(to, queryURL, &resp)
	if err != nil {
		return nil, reqInf, err
	}
	return resp.Response, reqInf, nil
}
