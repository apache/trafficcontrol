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
	"io/ioutil"
	"net"
	"net/http"

	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
)

const (
	API_DS_REQUESTS = "/api/1.3/deliveryservice_requests"
)

// Create a Delivery Service Request
func (to *Session) CreateDeliveryServiceRequest(dsr tc.DeliveryServiceRequest) (tc.Alerts, ReqInf, error) {

	var alerts tc.Alerts
	var remoteAddr net.Addr
	reqBody, err := json.Marshal(dsr)
	fmt.Printf("reqBody ---> %v\n", string(reqBody))
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return alerts, reqInf, err
	}
	resp, remoteAddr, err := to.rawRequest(http.MethodPost, API_DS_REQUESTS, reqBody)
	if err == nil {
		body, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			return alerts, reqInf, readErr
		}
		if err := json.Unmarshal(body, &alerts); err == nil {
			return alerts, reqInf, err
		}
	}

	defer resp.Body.Close()
	return alerts, reqInf, nil
}

// GET a DeliveryServiceRequest by the DeliveryServiceRequest XMLID
func (to *Session) GetDeliveryServiceRequestByXMLID(XMLID string) ([]tc.DeliveryServiceRequest, ReqInf, error) {

	route := fmt.Sprintf("%s?xmlId=%s", API_DS_REQUESTS, XMLID)
	fmt.Printf("route ---> %v\n", route)

	resp, remoteAddr, err := to.request(http.MethodGet, route, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	data := struct {
		Response []tc.DeliveryServiceRequest `json:"response"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// GET a DeliveryServiceRequest by the DeliveryServiceRequest id
func (to *Session) GetDeliveryServiceRequestByID(id int) ([]tc.DeliveryServiceRequest, ReqInf, error) {
	route := fmt.Sprintf("%s/%d", API_DS_REQUESTS, id)
	resp, remoteAddr, err := to.request(http.MethodGet, route, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	data := struct {
		Response []tc.DeliveryServiceRequest `json:"response"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// Update a DeliveryServiceRequest by ID
func (to *Session) UpdateDeliveryServiceRequestByID(id int, dsr tc.DeliveryServiceRequest) (tc.Alerts, ReqInf, error) {

	route := fmt.Sprintf("%s/%d", API_DS_REQUESTS, id)
	fmt.Printf("route ---> %v\n", route)

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(dsr)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	fmt.Printf("reqBody ---> %v\n", string(reqBody))
	resp, remoteAddr, err := to.request(http.MethodPut, route, reqBody)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

/*

// Returns a list of DeliveryServiceRequests
func (to *Session) GetDeliveryServiceRequests() ([]tc.DeliveryServiceRequest, ReqInf, error) {
	resp, remoteAddr, err := to.request(http.MethodGet, API_DS_REQUESTS, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.DeliveryServiceRequestsResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	return data.Response, reqInf, nil
}

// GET a DeliveryServiceRequest by the DeliveryServiceRequest assignee
func (to *Session) GetDeliveryServiceRequestByAssignee(assignee string) ([]tc.DeliveryServiceRequest, ReqInf, error) {
	url := fmt.Sprintf("%s/assignee/%s", API_DS_REQUESTS, assignee)
	resp, remoteAddr, err := to.request(http.MethodGet, url, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.DeliveryServiceRequestsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// DELETE a DeliveryServiceRequest by DeliveryServiceRequest assignee
func (to *Session) DeleteDeliveryServiceRequestByAssignee(assignee string) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s/assignee/%s", API_DS_REQUESTS, assignee)
	resp, remoteAddr, err := to.request(http.MethodDelete, route, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}
*/
