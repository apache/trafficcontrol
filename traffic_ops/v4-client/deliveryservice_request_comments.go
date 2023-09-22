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
	"net/url"
	"strconv"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

// apiDeliveryServiceRequestComments is the API version-relative route to
// the /deliveryservice_request_comments endpoint.
const apiDeliveryServiceRequestComments = "/deliveryservice_request_comments"

// CreateDeliveryServiceRequestComment creates the given Delivery Service
// Request comment.
func (to *Session) CreateDeliveryServiceRequestComment(comment tc.DeliveryServiceRequestComment, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(apiDeliveryServiceRequestComments, opts, comment, &alerts)
	return alerts, reqInf, err
}

// UpdateDeliveryServiceRequestComment replaces the Delivery Service Request
// comment identified by 'id' with the one provided.
func (to *Session) UpdateDeliveryServiceRequestComment(id int, comment tc.DeliveryServiceRequestComment, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	if opts.QueryParameters == nil {
		opts.QueryParameters = url.Values{}
	}
	opts.QueryParameters.Set("id", strconv.Itoa(id))
	var alerts tc.Alerts
	reqInf, err := to.put(apiDeliveryServiceRequestComments, opts, comment, &alerts)
	return alerts, reqInf, err
}

// GetDeliveryServiceRequestComments retrieves all comments on all Delivery
// Service Requests.
func (to *Session) GetDeliveryServiceRequestComments(opts RequestOptions) (tc.DeliveryServiceRequestCommentsResponse, toclientlib.ReqInf, error) {
	var data tc.DeliveryServiceRequestCommentsResponse
	reqInf, err := to.get(apiDeliveryServiceRequestComments, opts, &data)
	return data, reqInf, err
}

// DeleteDeliveryServiceRequestComment deletes the Delivery Service Request
// comment with the given ID.
func (to *Session) DeleteDeliveryServiceRequestComment(id int, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	if opts.QueryParameters == nil {
		opts.QueryParameters = url.Values{}
	}
	opts.QueryParameters.Set("id", strconv.Itoa(id))
	var alerts tc.Alerts
	reqInf, err := to.del(apiDeliveryServiceRequestComments, opts, &alerts)
	return alerts, reqInf, err
}
