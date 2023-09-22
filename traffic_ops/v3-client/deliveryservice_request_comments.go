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
	// API_DELIVERY_SERVICE_REQUEST_COMMENTS is Deprecated: will be removed in the next major version. Be aware this may not be the URI being requested, for clients created with Login and ClientOps.ForceLatestAPI false.
	API_DELIVERY_SERVICE_REQUEST_COMMENTS = apiBase + "/deliveryservice_request_comments"

	APIDeliveryServiceRequestComments = "/deliveryservice_request_comments"
)

// Create a delivery service request comment
func (to *Session) CreateDeliveryServiceRequestComment(comment tc.DeliveryServiceRequestComment) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(APIDeliveryServiceRequestComments, comment, nil, &alerts)
	return alerts, reqInf, err
}

func (to *Session) UpdateDeliveryServiceRequestCommentByIDWithHdr(id int, comment tc.DeliveryServiceRequestComment, header http.Header) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", APIDeliveryServiceRequestComments, id)
	var alerts tc.Alerts
	reqInf, err := to.put(route, comment, header, &alerts)
	return alerts, reqInf, err
}

// Update a delivery service request by ID
// Deprecated: UpdateDeliveryServiceRequestCommentByID will be removed in 6.0. Use UpdateDeliveryServiceRequestCommentByIDWithHdr.
func (to *Session) UpdateDeliveryServiceRequestCommentByID(id int, comment tc.DeliveryServiceRequestComment) (tc.Alerts, toclientlib.ReqInf, error) {
	return to.UpdateDeliveryServiceRequestCommentByIDWithHdr(id, comment, nil)
}

func (to *Session) GetDeliveryServiceRequestCommentsWithHdr(header http.Header) ([]tc.DeliveryServiceRequestComment, toclientlib.ReqInf, error) {
	var data tc.DeliveryServiceRequestCommentsResponse
	reqInf, err := to.get(APIDeliveryServiceRequestComments, header, &data)
	return data.Response, reqInf, err
}

// Returns a list of delivery service request comments
// Deprecated: GetDeliveryServiceRequestComments will be removed in 6.0. Use GetDeliveryServiceRequestCommentsWithHdr.
func (to *Session) GetDeliveryServiceRequestComments() ([]tc.DeliveryServiceRequestComment, toclientlib.ReqInf, error) {
	return to.GetDeliveryServiceRequestCommentsWithHdr(nil)
}

func (to *Session) GetDeliveryServiceRequestCommentByIDWithHdr(id int, header http.Header) ([]tc.DeliveryServiceRequestComment, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", APIDeliveryServiceRequestComments, id)
	var data tc.DeliveryServiceRequestCommentsResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GET a delivery service request comment by ID
// Deprecated: GetDeliveryServiceRequestCommentByID will be removed in 6.0. Use GetDeliveryServiceRequestCommentByIDWithHdr.
func (to *Session) GetDeliveryServiceRequestCommentByID(id int) ([]tc.DeliveryServiceRequestComment, toclientlib.ReqInf, error) {
	return to.GetDeliveryServiceRequestCommentByIDWithHdr(id, nil)
}

// DELETE a delivery service request comment by ID
func (to *Session) DeleteDeliveryServiceRequestCommentByID(id int) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", APIDeliveryServiceRequestComments, id)
	var alerts tc.Alerts
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}
