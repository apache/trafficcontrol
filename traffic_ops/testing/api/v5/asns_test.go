package v5

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

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

func TestASN(t *testing.T) {
	WithObjs(t, []TCObj{Types, CacheGroups, ASN}, func() {
		tomorrow := time.Now().AddDate(0, 0, 1).Format(time.RFC1123)
		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)

		methodTests := utils.V5TestCase{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {tomorrow}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateSorted()),
				},
				"OK when VALID ASN PARAMETER": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"asn": {"9999"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1)),
				},
			},
			"PUT": {
				"OK when VALID request": {
					ClientSession: TOSession,
					EndpointID:    GetASNId(t, "8888"),
					RequestBody: map[string]interface{}{
						"asn":            7777,
						"cachegroupName": "originCachegroup",
						"cachegroupId":   -1,
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"PRECONDITION FAILED when updating with IMS & IUS Headers": {
					ClientSession: TOSession,
					EndpointID:    GetASNId(t, "9999"),
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}}},
					RequestBody: map[string]interface{}{
						"asn":            8888,
						"cachegroupName": "originCachegroup",
						"cachegroupId":   -1,
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"BAD REQUEST when ASN is not unique": {
					ClientSession: TOSession,
					EndpointID:    GetASNId(t, "5555"),
					RequestBody: map[string]interface{}{
						"asn":            9999,
						"cachegroupName": "originCachegroup",
						"cachegroupId":   -1,
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
			"GET AFTER CHANGES": {
				"OK when CHANGES made": {
					ClientSession: TOSession,
					RequestOpts: client.RequestOptions{
						Header: http.Header{
							rfc.IfModifiedSince: {currentTimeRFC}, rfc.IfUnmodifiedSince: {currentTimeRFC},
						},
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					asn := tc.ASNV5{}

					if testCase.RequestBody != nil {
						if cgId, ok := testCase.RequestBody["cachegroupId"]; ok {
							if cgId == -1 {
								if cgName, ok := testCase.RequestBody["cachegroupName"]; ok {
									testCase.RequestBody["cachegroupId"] = GetCacheGroupId(t, cgName.(string))()
								}
							}
						}
						dat, err := json.Marshal(testCase.RequestBody)
						assert.NoError(t, err, "Error occurred when marshalling request body: %v", err)
						err = json.Unmarshal(dat, &asn)
						assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
					}

					switch method {
					case "GET", "GET AFTER CHANGES":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetASNs(testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.UpdateASN(testCase.EndpointID(), asn, testCase.RequestOpts)
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

func validateSorted() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		asnResp := resp.([]tc.ASNV5)
		var sortedList []string
		assert.RequireGreaterOrEqual(t, len(asnResp), 2, "Need at least 2 ASNs in Traffic Ops to test sorted, found: %d", len(asnResp))

		for _, asn := range asnResp {
			sortedList = append(sortedList, strconv.Itoa(asn.ASN))
		}

		res := sort.SliceIsSorted(sortedList, func(p, q int) bool {
			return sortedList[p] < sortedList[q]
		})
		assert.Equal(t, res, true, "List is not sorted by their names: %v", sortedList)
	}
}

func GetASNId(t *testing.T, ASN string) func() int {
	return func() int {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("asn", ASN)

		resp, _, err := TOSession.GetASNs(opts)
		assert.RequireNoError(t, err, "Get ASNs Request failed with error: %v", err)
		assert.RequireEqual(t, len(resp.Response), 1, "Expected response object length 1, but got %d", len(resp.Response))
		assert.RequireNotNil(t, &resp.Response[0].ID, "Expected id to not be nil")

		return resp.Response[0].ID
	}
}

func CreateTestASNs(t *testing.T) {
	assert.RequireGreaterOrEqual(t, len(testData.CacheGroups), 1, "Need at least one Cache Group to test creating ASNs")

	cg := testData.CacheGroups[0]
	assert.RequireNotNil(t, cg.Name, "Cache Group found in the test data with null or undefined name")

	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", *cg.Name)
	resp, _, err := TOSession.GetCacheGroups(opts)
	assert.RequireNoError(t, err, "Unable to get cachgroup ID: %v - alerts: %+v", err, resp.Alerts)
	assert.RequireEqual(t, 1, len(resp.Response), "Expected exactly one Cache Group with Name '%s', got: %d", *cg.Name, len(resp.Response))
	assert.RequireNotNil(t, resp.Response[0].ID, "Cache Group '%s' had no ID in Traffic Ops response", *cg.Name)

	id := *resp.Response[0].ID
	for _, asn := range testData.ASNs {
		asn.CachegroupID = id
		resp, _, err := TOSession.CreateASN(asn, client.RequestOptions{})
		assert.NoError(t, err, "Could not create ASN: %v - alerts: %+v", err, resp)
	}
}

func DeleteTestASNs(t *testing.T) {
	opts := client.NewRequestOptions()
	// Retrieve the ASNs to delete
	asns, _, err := TOSession.GetASNs(opts)
	assert.NoError(t, err, "Error trying to fetch ASNs for deletion: %v - alerts: %+v", err, asns.Alerts)
	for _, asn := range asns.Response {
		alerts, _, err := TOSession.DeleteASN(asn.ID, client.RequestOptions{})
		assert.NoError(t, err, "Cannot delete ASN %d: %v - alerts: %+v", asn.ASN, err, alerts)

		// Retrieve the ASN to see if it got deleted
		opts.QueryParameters.Set("asn", strconv.Itoa(asn.ASN))
		asns, _, err := TOSession.GetASNs(opts)
		assert.NoError(t, err, "Error trying to fetch ASN after deletion: %v - alerts: %+v", err, asns.Alerts)
		assert.Equal(t, 0, len(asns.Response), "Expected ASN %d to be deleted, but it was found in Traffic Ops's response", asn.ASN)
	}
}
