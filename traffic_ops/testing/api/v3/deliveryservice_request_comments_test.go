package v3

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
	"sort"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

func TestDeliveryServiceRequestComments(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Tenants, DeliveryServiceRequests, DeliveryServiceRequestComments}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.V3TestCaseT[tc.DeliveryServiceRequestComment]{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfModifiedSince: {tomorrow}},
					Expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when VALID ID parameter": {
					EndpointID:    GetDSRequestCommentId(t),
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1)),
				},
				"VALIDATE SORT when DEFAULT is ASC ORDER": {
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateSortedDSRequestComments()),
				},
				"OK when CHANGES made": {
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfModifiedSince: {currentTimeRFC}},
					Expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
			"PUT": {
				"OK when VALID request": {
					EndpointID:    GetDSRequestCommentId(t),
					ClientSession: TOSession,
					RequestBody: tc.DeliveryServiceRequestComment{
						DeliveryServiceRequestID: GetDSRequestId(t, "test-ds1")(),
						Value:                    "updated comment",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"PRECONDITION FAILED when updating with IF-UNMODIFIED-SINCE Header": {
					EndpointID:     GetDSRequestCommentId(t),
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}},
					RequestBody:    tc.DeliveryServiceRequestComment{},
					Expectations:   utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"PRECONDITION FAILED when updating with IFMATCH ETAG Header": {
					EndpointID:     GetDSRequestCommentId(t),
					ClientSession:  TOSession,
					RequestBody:    tc.DeliveryServiceRequestComment{},
					RequestHeaders: http.Header{rfc.IfMatch: {rfc.ETag(currentTime)}},
					Expectations:   utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							if name == "OK when VALID ID parameter" {
								resp, reqInf, err := testCase.ClientSession.GetDeliveryServiceRequestCommentByIDWithHdr(testCase.EndpointID(), testCase.RequestHeaders)
								for _, check := range testCase.Expectations {
									check(t, reqInf, resp, tc.Alerts{}, err)
								}
							} else {
								resp, reqInf, err := testCase.ClientSession.GetDeliveryServiceRequestCommentsWithHdr(testCase.RequestHeaders)
								for _, check := range testCase.Expectations {
									check(t, reqInf, resp, tc.Alerts{}, err)
								}
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.UpdateDeliveryServiceRequestCommentByIDWithHdr(testCase.EndpointID(), testCase.RequestBody, testCase.RequestHeaders)
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

func GetDSRequestCommentId(t *testing.T) func() int {
	return func() int {
		resp, _, err := TOSession.GetDeliveryServiceRequestCommentsWithHdr(http.Header{})
		assert.RequireNoError(t, err, "Get Delivery Service Request Comments failed with error: %v", err)
		assert.RequireGreaterOrEqual(t, len(resp), 1, "Expected delivery service request comments response object length of atleast 1, but got %d", len(resp))
		assert.RequireNotNil(t, resp[0].ID, "Expected id to not be nil")

		return resp[0].ID
	}
}

func GetDSRequestId(t *testing.T, xmlId string) func() int {
	return func() int {
		resp, _, err := TOSession.GetDeliveryServiceRequestByXMLIDWithHdr(xmlId, http.Header{})
		assert.RequireNoError(t, err, "Get Delivery Service Requests failed with error: %v", err)
		assert.RequireGreaterOrEqual(t, len(resp), 1, "Expected delivery service requests response object length of atleast 1, but got %d", len(resp))
		assert.RequireNotNil(t, resp[0].ID, "Expected id to not be nil")

		return resp[0].ID
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

func CreateTestDeliveryServiceRequestComments(t *testing.T) {
	for _, comment := range testData.DeliveryServiceRequestComments {
		resp, _, err := TOSession.GetDeliveryServiceRequestByXMLIDWithHdr(comment.XMLID, http.Header{})
		assert.NoError(t, err, "Cannot get Delivery Service Request by XMLID '%s': %v", comment.XMLID, err)
		assert.Equal(t, len(resp), 1, "Found %d Delivery Service request by XMLID '%s, expected exactly one", len(resp), comment.XMLID)
		assert.NotNil(t, resp[0].ID, "Got Delivery Service Request with xml_id '%s' that had a null ID", comment.XMLID)

		comment.DeliveryServiceRequestID = resp[0].ID
		alerts, _, err := TOSession.CreateDeliveryServiceRequestComment(comment)
		assert.NoError(t, err, "Could not create Delivery Service Request Comment: %v - alerts: %+v", err, alerts.Alerts)
	}
}

func DeleteTestDeliveryServiceRequestComments(t *testing.T) {
	resp, _, err := TOSession.GetDeliveryServiceRequestCommentsWithHdr(http.Header{})
	assert.NoError(t, err, "Unexpected error getting Delivery Service Request Comments: %v", err)

	for _, comment := range resp {
		alerts, _, err := TOSession.DeleteDeliveryServiceRequestCommentByID(comment.ID)
		assert.NoError(t, err, "Cannot delete Delivery Service Request Comment #%d: %v - alerts: %+v", comment.ID, err, alerts.Alerts)

		// Retrieve the delivery service request comment to see if it got deleted
		resp, _, err := TOSession.GetDeliveryServiceRequestCommentByIDWithHdr(comment.ID, http.Header{})
		assert.NoError(t, err, "Unexpected error fetching Delivery Service Request Comment %d after deletion: %v", comment.ID, err)
		assert.Equal(t, len(resp), 0, "Expected Delivery Service Request Comment #%d to be deleted, but it was found in Traffic Ops", comment.ID)
	}
}
