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
	"encoding/json"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	tc "github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/assert"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

func TestCDNs(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Parameters}, func() {
		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.V4TestCase{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {tomorrow}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID request": {
					ClientSession: TOSession, Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						utils.ResponseLengthGreaterOrEqual(1), validateCDNSort()),
				},
				"OK when VALID NAME parameter": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"name": {"cdn1"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						validateCDNFields(map[string]interface{}{"Name": "cdn1"})),
				},
				"OK when VALID DOMAINNAME parameter": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"domainName": {"test.cdn2.net"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						validateCDNFields(map[string]interface{}{"DomainName": "test.cdn2.net"})),
				},
				"OK when VALID DNSSECENABLED parameter": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"dnssecEnabled": {"false"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateCDNFields(map[string]interface{}{"DNSSECEnabled": false})),
				},
				"VALID when SORTORDER param is DESC": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"sortOrder": {"desc"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateCDNDescSort()),
				},
				"FIRST RESULT when LIMIT=1": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateCDNPagination("limit")),
				},
				"SECOND RESULT when LIMIT=1 OFFSET=1": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}, "offset": {"1"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateCDNPagination("offset")),
				},
				"SECOND RESULT when LIMIT=1 PAGE=2": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}, "page": {"2"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateCDNPagination("page")),
				},
				"BAD REQUEST when INVALID LIMIT parameter": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"limit": {"-2"}}},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID OFFSET parameter": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"limit": {"1"}, "offset": {"0"}}},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID PAGE parameter": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"limit": {"1"}, "page": {"0"}}},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
			"POST": {
				"BAD REQUEST when CDN ALREADY EXISTS": {
					ClientSession: TOSession, RequestBody: map[string]interface{}{
						"name":          "cdn3",
						"dnssecEnabled": false,
						"domainName":    "test.cdn3.net",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when EMPTY NAME": {
					ClientSession: TOSession, RequestBody: map[string]interface{}{
						"name":          "",
						"dnssecEnabled": false,
						"domainName":    "test.noname.net",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when EMPTY DOMAIN NAME": {
					ClientSession: TOSession, RequestBody: map[string]interface{}{
						"name":          "nodomain",
						"dnssecEnabled": false,
						"domainName":    "",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
			"PUT": {
				"OK when VALID request": {
					EndpointId: GetCDNID(t, "cdn1"), ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"dnssecEnabled": false,
						"domainName":    "newDomain",
						"name":          "cdn1",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"PRECONDITION FAILED when updating with IF-UNMODIFIED-SINCE Headers": {
					EndpointId: GetCDNID(t, "cdn1"), ClientSession: TOSession,
					RequestOpts: client.RequestOptions{Header: http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}}},
					RequestBody: map[string]interface{}{
						"dnssecEnabled": false,
						"domainName":    "newDomain",
						"name":          "cdn1",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"PRECONDITION FAILED when updating with IFMATCH ETAG Header": {
					EndpointId: GetCDNID(t, "cdn1"), ClientSession: TOSession,
					RequestOpts: client.RequestOptions{Header: http.Header{rfc.IfMatch: {rfc.ETag(currentTime)}}},
					RequestBody: map[string]interface{}{
						"dnssecEnabled": false,
						"domainName":    "newDomain",
						"name":          "cdn1",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
			},
			"DELETE": {
				"NOT FOUND when INVALID ID parameter": {
					EndpointId: func() int { return 111111 }, ClientSession: TOSession,
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
			},
			"GET AFTER CHANGES": {
				"OK when CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {currentTimeRFC}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
		}
		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					cdn := tc.CDN{}

					if testCase.RequestBody != nil {
						dat, err := json.Marshal(testCase.RequestBody)
						assert.NoError(t, err, "Error occurred when marshalling request body: %v", err)
						err = json.Unmarshal(dat, &cdn)
						assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
					}

					switch method {
					case "GET", "GET AFTER CHANGES":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetCDNs(testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.CreateCDN(cdn, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.UpdateCDN(testCase.EndpointId(), cdn, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.DeleteCDN(testCase.EndpointId(), testCase.RequestOpts)
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

func validateCDNFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		cdnResp := resp.([]tc.CDN)
		for field, expected := range expectedResp {
			for _, cdn := range cdnResp {
				switch field {
				case "Name":
					assert.Equal(t, expected, cdn.Name, "Expected Name to be %v, but got %v", expected, cdn.Name)
				case "DomainName":
					assert.Equal(t, expected, cdn.DomainName, "Expected DomainName to be %v, but got %v", expected, cdn.DomainName)
				case "DNSSECEnabled":
					assert.Equal(t, expected, cdn.DNSSECEnabled, "Expected DNSSECEnabled to be %v, but got %v", expected, cdn.DNSSECEnabled)
				default:
					t.Errorf("Expected field: %v, does not exist in response", field)
				}
			}
		}
	}
}

func validateCDNPagination(paginationParam string) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		paginationResp := resp.([]tc.CDN)

		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("orderby", "id")
		respBase, _, err := TOSession.GetCDNs(opts)
		assert.RequireNoError(t, err, "Cannot get CDNs: %v - alerts: %+v", err, respBase.Alerts)

		cachegroup := respBase.Response
		assert.RequireGreaterOrEqual(t, len(cachegroup), 3, "Need at least 3 CDNs in Traffic Ops to test pagination support, found: %d", len(cachegroup))
		switch paginationParam {
		case "limit:":
			assert.Exactly(t, cachegroup[:1], paginationResp, "Expected GET CDNs with limit = 1 to return first result")
		case "offset":
			assert.Exactly(t, cachegroup[1:2], paginationResp, "Expected GET CDNs with limit = 1, offset = 1 to return second result")
		case "page":
			assert.Exactly(t, cachegroup[1:2], paginationResp, "Expected GET CDNs with limit = 1, page = 2 to return second result")
		}
	}
}

func validateCDNSort() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		var sortedList []string
		cdnResp := resp.([]tc.CDN)

		for _, cdn := range cdnResp {
			sortedList = append(sortedList, cdn.Name)
		}

		res := sort.SliceIsSorted(sortedList, func(p, q int) bool {
			return sortedList[p] < sortedList[q]
		})
		assert.Equal(t, res, true, "List is not sorted by their names: %v", sortedList)
	}
}

func validateCDNDescSort() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		cdnDescResp := resp.([]tc.CDN)
		var descSortedList []string
		var ascSortedList []string
		assert.GreaterOrEqual(t, len(cdnDescResp), 2, "Need at least 2 CDNs in Traffic Ops to test desc sort, found: %d", len(cdnDescResp))
		// Get CDNs in the default ascending order for comparison.
		cdnAscResp, _, err := TOSession.GetCDNs(client.RequestOptions{})
		assert.NoError(t, err, "Unexpected error getting CDNs with default sort order: %v - alerts: %+v", err, cdnAscResp.Alerts)
		assert.GreaterOrEqual(t, len(cdnAscResp.Response), 2, "Need at least 2 CDNs in Traffic Ops to test sort, found %d", len(cdnAscResp.Response))
		// Verify the response match in length, i.e. equal amount of CDNs.
		assert.Equal(t, len(cdnAscResp.Response), len(cdnDescResp), "Expected descending order response length: %v, to match ascending order response length %v", len(cdnAscResp.Response), len(cdnDescResp))
		// Insert CDN names to the front of a new list, so they are now reversed to be in ascending order.
		for _, cdn := range cdnDescResp {
			descSortedList = append([]string{cdn.Name}, descSortedList...)
		}
		// Insert CDN names by appending to a new list, so they stay in ascending order.
		for _, cdn := range cdnAscResp.Response {
			ascSortedList = append(ascSortedList, cdn.Name)
		}
		assert.Exactly(t, ascSortedList, descSortedList, "CDN responses are not equal after reversal: %v - %v", ascSortedList, descSortedList)
	}
}

func GetCDNID(t *testing.T, cdnName string) func() int {
	return func() int {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("name", cdnName)
		cdnsResp, _, err := TOSession.GetCDNs(opts)
		assert.NoError(t, err, "Get CDNs Request failed with error:", err)
		assert.Equal(t, 1, len(cdnsResp.Response), "Expected response object length 1, but got %d", len(cdnsResp.Response))
		assert.NotNil(t, cdnsResp.Response[0].ID, "Expected id to not be nil")
		return cdnsResp.Response[0].ID
	}
}

func CreateTestCDNs(t *testing.T) {
	for _, cdn := range testData.CDNs {
		resp, _, err := TOSession.CreateCDN(cdn, client.RequestOptions{})
		assert.NoError(t, err, "Could not create CDN: %v - alerts: %+v", err, resp.Alerts)
	}
}

func DeleteTestCDNs(t *testing.T) {
	resp, _, err := TOSession.GetCDNs(client.RequestOptions{})
	assert.NoError(t, err, "Cannot get CDNs: %v - alerts: %+v", err, resp.Alerts)
	for _, cdn := range resp.Response {
		delResp, _, err := TOSession.DeleteCDN(cdn.ID, client.RequestOptions{})
		assert.NoError(t, err, "Cannot delete CDN '%s' (#%d): %v - alerts: %+v", cdn.Name, cdn.ID, err, delResp.Alerts)

		// Retrieve the CDN to see if it got deleted
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.Itoa(cdn.ID))
		cdns, _, err := TOSession.GetCDNs(opts)
		assert.NoError(t, err, "Error deleting CDN '%s': %v - alerts: %+v", cdn.Name, err, cdns.Alerts)
		assert.Equal(t, 0, len(cdns.Response), "Expected CDN '%s' to be deleted", cdn.Name)
	}
}
