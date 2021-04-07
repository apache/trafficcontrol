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

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

const (
	// APIDeliveryServiceRequestComments is the API version-relative route to
	// the /deliveryservice_request_comments endpoint.
	APIDeliveryServiceRequestComments = "/deliveryservice_request_comments"
)

// CreateDeliveryServiceRequestComment creates the given Delivery Service
// Request comment.
func (to *Session) CreateDeliveryServiceRequestComment(comment tc.DeliveryServiceRequestComment) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(APIDeliveryServiceRequestComments, comment, nil, &alerts)
	return alerts, reqInf, err
}

// UpdateDeliveryServiceRequestComment replaces the Delivery Service Request
// comment identified by 'id' with the one provided.
func (to *Session) UpdateDeliveryServiceRequestComment(id int, comment tc.DeliveryServiceRequestComment, header http.Header) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", APIDeliveryServiceRequestComments, id)
	var alerts tc.Alerts
	reqInf, err := to.put(route, comment, header, &alerts)
	return alerts, reqInf, err
}

// GetDeliveryServiceRequestComments retrieves all comments on all Delivery
// Service Requests.
func (to *Session) GetDeliveryServiceRequestComments(header http.Header) ([]tc.DeliveryServiceRequestComment, toclientlib.ReqInf, error) {
	var data tc.DeliveryServiceRequestCommentsResponse
	reqInf, err := to.get(APIDeliveryServiceRequestComments, header, &data)
	return data.Response, reqInf, err
}

// GetDeliveryServiceRequestComment retrieves the Delivery Service Request
// comment with the given ID.
func (to *Session) GetDeliveryServiceRequestComment(id int, header http.Header) ([]tc.DeliveryServiceRequestComment, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", APIDeliveryServiceRequestComments, id)
	var data tc.DeliveryServiceRequestCommentsResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// DeleteDeliveryServiceRequestComment deletes the Delivery Service Request
// comment with the given ID.
func (to *Session) DeleteDeliveryServiceRequestComment(id int) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", APIDeliveryServiceRequestComments, id)
	var alerts tc.Alerts
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}
