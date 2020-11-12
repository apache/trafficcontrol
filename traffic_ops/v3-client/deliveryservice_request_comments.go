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
	"net"
	"net/http"

	"github.com/apache/trafficcontrol/lib/go-tc"

	"fmt"
)

const (
	API_DELIVERY_SERVICE_REQUEST_COMMENTS = apiBase + "/deliveryservice_request_comments"
)

// Create a delivery service request comment
func (to *Session) CreateDeliveryServiceRequestComment(comment tc.DeliveryServiceRequestComment) (tc.Alerts, ReqInf, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(comment)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	resp, remoteAddr, err := to.request(http.MethodPost, API_DELIVERY_SERVICE_REQUEST_COMMENTS, reqBody, nil)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

func (to *Session) UpdateDeliveryServiceRequestCommentByIDWithHdr(id int, comment tc.DeliveryServiceRequestComment, header http.Header) (tc.Alerts, ReqInf, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(comment)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	route := fmt.Sprintf("%s?id=%d", API_DELIVERY_SERVICE_REQUEST_COMMENTS, id)
	resp, remoteAddr, err := to.request(http.MethodPut, route, reqBody, header)
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
	}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

// Update a delivery service request by ID
// Deprecated: UpdateDeliveryServiceRequestCommentByID will be removed in 6.0. Use UpdateDeliveryServiceRequestCommentByIDWithHdr.
func (to *Session) UpdateDeliveryServiceRequestCommentByID(id int, comment tc.DeliveryServiceRequestComment) (tc.Alerts, ReqInf, error) {
	return to.UpdateDeliveryServiceRequestCommentByIDWithHdr(id, comment, nil)
}

func (to *Session) GetDeliveryServiceRequestCommentsWithHdr(header http.Header) ([]tc.DeliveryServiceRequestComment, ReqInf, error) {
	resp, remoteAddr, err := to.request(http.MethodGet, API_DELIVERY_SERVICE_REQUEST_COMMENTS, nil, header)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		if reqInf.StatusCode == http.StatusNotModified {
			return []tc.DeliveryServiceRequestComment{}, reqInf, nil
		}
	}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.DeliveryServiceRequestCommentsResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	return data.Response, reqInf, nil
}

// Returns a list of delivery service request comments
// Deprecated: GetDeliveryServiceRequestComments will be removed in 6.0. Use GetDeliveryServiceRequestCommentsWithHdr.
func (to *Session) GetDeliveryServiceRequestComments() ([]tc.DeliveryServiceRequestComment, ReqInf, error) {
	return to.GetDeliveryServiceRequestCommentsWithHdr(nil)
}

func (to *Session) GetDeliveryServiceRequestCommentByIDWithHdr(id int, header http.Header) ([]tc.DeliveryServiceRequestComment, ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", API_DELIVERY_SERVICE_REQUEST_COMMENTS, id)
	resp, remoteAddr, err := to.request(http.MethodGet, route, nil, header)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		if reqInf.StatusCode == http.StatusNotModified {
			return []tc.DeliveryServiceRequestComment{}, reqInf, nil
		}
	}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.DeliveryServiceRequestCommentsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// GET a delivery service request comment by ID
// Deprecated: GetDeliveryServiceRequestCommentByID will be removed in 6.0. Use GetDeliveryServiceRequestCommentByIDWithHdr.
func (to *Session) GetDeliveryServiceRequestCommentByID(id int) ([]tc.DeliveryServiceRequestComment, ReqInf, error) {
	return to.GetDeliveryServiceRequestCommentByIDWithHdr(id, nil)
}

// DELETE a delivery service request comment by ID
func (to *Session) DeleteDeliveryServiceRequestCommentByID(id int) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", API_DELIVERY_SERVICE_REQUEST_COMMENTS, id)
	resp, remoteAddr, err := to.request(http.MethodDelete, route, nil, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}
