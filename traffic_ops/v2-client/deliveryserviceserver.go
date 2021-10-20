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
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/v6/lib/go-tc"
	"github.com/apache/trafficcontrol/v6/lib/go-util"
)

// CreateDeliveryServiceServers associates the given servers with the given delivery services. If replace is true, it deletes any existing associations for the given delivery service.
func (to *Session) CreateDeliveryServiceServers(dsID int, serverIDs []int, replace bool) (*tc.DSServerIDs, error) {
	path := apiBase + `/deliveryserviceserver`
	req := tc.DSServerIDs{
		DeliveryServiceID: util.IntPtr(dsID),
		ServerIDs:         serverIDs,
		Replace:           util.BoolPtr(replace),
	}
	jsonReq, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}
	resp := struct {
		Response tc.DSServerIDs `json:"response"`
	}{}
	if _, err := post(to, path, jsonReq, &resp); err != nil {
		return nil, err
	}
	return &resp.Response, nil
}

func (to *Session) DeleteDeliveryServiceServer(dsID int, serverID int) (tc.Alerts, ReqInf, error) {
	route := apiBase + `/deliveryserviceserver/` + strconv.Itoa(dsID) + "/" + strconv.Itoa(serverID)
	reqResp, remoteAddr, err := to.request(http.MethodDelete, route, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, errors.New("requesting from Traffic Ops: " + err.Error())
	}
	defer reqResp.Body.Close()
	resp := tc.Alerts{}
	if err = json.NewDecoder(reqResp.Body).Decode(&resp); err != nil {
		return tc.Alerts{}, reqInf, errors.New("decoding response: " + err.Error())
	}
	return resp, reqInf, nil
}

// GetDeliveryServiceServers gets all delivery service servers, with the default API limit.
func (to *Session) GetDeliveryServiceServers() (tc.DeliveryServiceServerResponse, ReqInf, error) {
	return to.getDeliveryServiceServers(url.Values{})
}

// GetDeliveryServiceServersN gets all delivery service servers, with a limit of n.
func (to *Session) GetDeliveryServiceServersN(n int) (tc.DeliveryServiceServerResponse, ReqInf, error) {
	return to.getDeliveryServiceServers(url.Values{"limit": []string{strconv.Itoa(n)}})
}

// GetDeliveryServiceServersWithLimits gets all delivery service servers, allowing specifying the limit of mappings to return, the delivery services to return, and the servers to return.
// The limit may be 0, in which case the default limit will be applied. The deliveryServiceIDs and serverIDs may be nil or empty, in which case all delivery services and/or servers will be returned.
func (to *Session) GetDeliveryServiceServersWithLimits(limit int, deliveryServiceIDs []int, serverIDs []int) (tc.DeliveryServiceServerResponse, ReqInf, error) {
	vals := url.Values{}
	if limit != 0 {
		vals.Set("limit", strconv.Itoa(limit))
	}

	if len(deliveryServiceIDs) != 0 {
		dsIDStrs := []string{}
		for _, dsID := range deliveryServiceIDs {
			dsIDStrs = append(dsIDStrs, strconv.Itoa(dsID))
		}
		vals.Set("deliveryserviceids", strings.Join(dsIDStrs, ","))
	}

	if len(serverIDs) != 0 {
		serverIDStrs := []string{}
		for _, serverID := range serverIDs {
			serverIDStrs = append(serverIDStrs, strconv.Itoa(serverID))
		}
		vals.Set("serverids", strings.Join(serverIDStrs, ","))
	}

	return to.getDeliveryServiceServers(vals)
}

func (to *Session) getDeliveryServiceServers(urlQuery url.Values) (tc.DeliveryServiceServerResponse, ReqInf, error) {
	route := apiBase + `/deliveryserviceserver`
	if qry := urlQuery.Encode(); qry != "" {
		route += `?` + qry
	}
	reqResp, remoteAddr, err := to.request(http.MethodGet, route, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.DeliveryServiceServerResponse{}, reqInf, errors.New("requesting from Traffic Ops: " + err.Error())
	}
	defer reqResp.Body.Close()
	resp := tc.DeliveryServiceServerResponse{}
	if err = json.NewDecoder(reqResp.Body).Decode(&resp); err != nil {
		return tc.DeliveryServiceServerResponse{}, reqInf, errors.New("decoding response: " + err.Error())
	}
	return resp, reqInf, nil
}
