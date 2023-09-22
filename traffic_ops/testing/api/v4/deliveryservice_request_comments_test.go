package v4

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
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"
)

func TestDeliveryServiceRequestComments(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Tenants, DeliveryServiceRequests, DeliveryServiceRequestComments}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.TestCase[client.Session, client.RequestOptions, tc.DeliveryServiceRequestComment]{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {tomorrow}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when VALID ID parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"id": {strconv.Itoa(GetDSRequestCommentId(t, "admin")())}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1)),
				},
				"VALIDATE SORT when DEFAULT is ASC ORDER": {
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateSortedDSRequestComments()),
				},
				"OK when CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {currentTimeRFC}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
			"PUT": {
				"OK when VALID request": {
					EndpointID:    GetDSRequestCommentId(t, "admin"),
					ClientSession: TOSession,
					RequestBody: tc.DeliveryServiceRequestComment{
						DeliveryServiceRequestID: GetDSRequestId(t, "test-ds1")(),
						Value:                    "updated comment",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"PRECONDITION FAILED when updating with IF-UNMODIFIED-SINCE Header": {
					EndpointID:    GetDSRequestCommentId(t, "admin"),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}}},
					RequestBody:   tc.DeliveryServiceRequestComment{},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"PRECONDITION FAILED when updating with IFMATCH ETAG Header": {
					EndpointID:    GetDSRequestCommentId(t, "admin"),
					ClientSession: TOSession,
					RequestBody:   tc.DeliveryServiceRequestComment{},
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfMatch: {rfc.ETag(currentTime)}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetDeliveryServiceRequestComments(testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.UpdateDeliveryServiceRequestComment(testCase.EndpointID(), testCase.RequestBody, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					}
				}
			})
		}
	})
}

func GetDSRequestCommentId(t *testing.T, author string) func() int {
	return func() int {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("author", author)

		resp, _, err := TOSession.GetDeliveryServiceRequestComments(opts)
		assert.RequireNoError(t, err, "Get Delivery Service Request Comments failed with error: %v", err)
		assert.RequireGreaterOrEqual(t, len(resp.Response), 1, "Expected delivery service request comments response object length of atleast 1, but got %d", len(resp.Response))
		assert.RequireNotNil(t, resp.Response[0].ID, "Expected id to not be nil")

		return resp.Response[0].ID
	}
}

func GetDSRequestId(t *testing.T, xmlId string) func() int {
	return func() int {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("xmlId", xmlId)

		resp, _, err := TOSession.GetDeliveryServiceRequests(opts)
		assert.RequireNoError(t, err, "Get Delivery Service Requests failed with error: %v", err)
		assert.RequireGreaterOrEqual(t, len(resp.Response), 1, "Expected delivery service requests response object length of atleast 1, but got %d", len(resp.Response))
		assert.RequireNotNil(t, resp.Response[0].ID, "Expected id to not be nil")

		return *resp.Response[0].ID
	}
}

func validateSortedDSRequestComments() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, err error) {
		var sortedList []string
		dsReqComments := resp.([]tc.DeliveryServiceRequestComment)

		for _, comment := range dsReqComments {
			sortedList = append(sortedList, comment.XMLID)
		}

		res := sort.SliceIsSorted(sortedList, func(p, q int) bool {
			return sortedList[p] < sortedList[q]
		})
		assert.Equal(t, res, true, "List is not sorted by their names: %v", sortedList)
	}
}
