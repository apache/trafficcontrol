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
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

// CreateDeliveryServiceServers associates the given servers with the given
// Delivery Services. If replace is true, it deletes any existing associations
// for the given Delivery Service.
func (to *Session) CreateDeliveryServiceServers(dsID int, serverIDs []int, replace bool) (*tc.DSServerIDs, toclientlib.ReqInf, error) {
	path := APIDeliveryServiceServer
	req := tc.DSServerIDs{
		DeliveryServiceID: util.IntPtr(dsID),
		ServerIDs:         serverIDs,
		Replace:           util.BoolPtr(replace),
	}
	resp := struct {
		Response tc.DSServerIDs `json:"response"`
	}{}
	reqInf, err := to.post(path, req, nil, &resp)
	if err != nil {
		return nil, reqInf, err
	}
	return &resp.Response, reqInf, nil
}

// DeleteDeliveryServiceServer removes the association between the Delivery
// Service identified by dsID and the server identified by serverID.
func (to *Session) DeleteDeliveryServiceServer(dsID int, serverID int) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d/%d", APIDeliveryServiceServer, dsID, serverID)
	resp := tc.Alerts{}
	reqInf, err := to.del(route, nil, &resp)
	return resp, reqInf, err
}

// AssignServersToDeliveryService assigns the given list of servers to the
// Delivery Service with the given xmlID.
func (to *Session) AssignServersToDeliveryService(servers []string, xmlID string) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf(APIDeliveryServicesServers, url.QueryEscape(xmlID))
	dss := tc.DeliveryServiceServers{ServerNames: servers, XmlId: xmlID}
	resp := tc.Alerts{}
	reqInf, err := to.post(route, dss, nil, &resp)
	return resp, reqInf, err
}

// GetDeliveryServiceServers returns associations between Delivery Services and
// servers.
func (to *Session) GetDeliveryServiceServers(urlQuery url.Values, h http.Header) (tc.DeliveryServiceServerResponse, toclientlib.ReqInf, error) {
	route := APIDeliveryServiceServer
	if qry := urlQuery.Encode(); qry != "" {
		route += `?` + qry
	}
	resp := tc.DeliveryServiceServerResponse{}
	reqInf, err := to.get(route, h, &resp)
	return resp, reqInf, err
}
