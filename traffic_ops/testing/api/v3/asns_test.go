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
	"encoding/json"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/assert"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/traffic_ops/v3-client"
)

func TestASN(t *testing.T) {
	WithObjs(t, []TCObj{Types, CacheGroups, ASN}, func() {

		tomorrow := time.Now().AddDate(0, 0, 1).Format(time.RFC1123)
		currentTime := time.Now().UTC().Add(-5 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)

		methodTests := map[string]map[string]struct {
			endpointId     func() int
			clientSession  *client.Session
			requestParams  url.Values
			requestHeaders http.Header
			requestBody    map[string]interface{}
			expectations   []utils.CkReqFunc
		}{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					clientSession: TOSession, requestHeaders: http.Header{rfc.IfModifiedSince: {tomorrow}},
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID request": {
					clientSession: TOSession, expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateSorted()),
				},
				"OK when VALID ASN PARAMETER": {
					clientSession: TOSession, requestParams: url.Values{"asn": {"9999"}},
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1)),
				},
			},
			"PUT": {
				"OK when VALID request": {
					clientSession: TOSession, endpointId: GetASNId(t, "8888"),
					requestBody: map[string]interface{}{
						"asn":            7777,
						"cachegroupName": "originCachegroup",
						"cachegroupId":   -1,
					},
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
			"GET AFTER CHANGES": {
				"OK when CHANGES made": {
					clientSession:  TOSession,
					requestHeaders: http.Header{rfc.IfModifiedSince: {currentTimeRFC}, rfc.IfUnmodifiedSince: {currentTimeRFC}},
					expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					asn := tc.ASN{}

					if testCase.requestBody != nil {
						if cgId, ok := testCase.requestBody["cachegroupId"]; ok {
							if cgId == -1 {
								if cgName, ok := testCase.requestBody["cachegroupName"]; ok {
									testCase.requestBody["cachegroupId"] = GetCacheGroupId(t, cgName.(string))()
								}
							}
						}
						dat, err := json.Marshal(testCase.requestBody)
						assert.NoError(t, err, "Error occurred when marshalling request body: %v", err)
						err = json.Unmarshal(dat, &asn)
						assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
					}

					switch method {
					case "GET", "GET AFTER CHANGES":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.clientSession.GetASNsWithHeader(&testCase.requestParams, testCase.requestHeaders)
							for _, check := range testCase.expectations {
								check(t, reqInf, resp, tc.Alerts{}, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.clientSession.UpdateASNByID(testCase.endpointId(), asn)
							for _, check := range testCase.expectations {
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
		asnResp := resp.([]tc.ASN)
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
		params := url.Values{"asn": {ASN}}

		resp, _, err := TOSession.GetASNsWithHeader(&params, http.Header{})
		assert.RequireNoError(t, err, "Get ASNs Request failed with error: %v", err)
		assert.RequireEqual(t, len(resp), 1, "Expected response object length 1, but got %d", len(resp))
		assert.RequireNotNil(t, &resp[0].ID, "Expected id to not be nil")

		return resp[0].ID
	}
}

func CreateTestASNs(t *testing.T) {
	resp, _, err := TOSession.GetCacheGroupNullableByNameWithHdr(*testData.CacheGroups[0].Name, http.Header{})
	assert.RequireNoError(t, err, "Unable to get cachgroup ID: %v - resp: %+v", err, resp)

	for _, asn := range testData.ASNs {
		asn.CachegroupID = *resp[0].ID
		resp, _, err := TOSession.CreateASN(asn)
		assert.NoError(t, err, "Could not create ASN: %v - resp: %+v", err, resp)
	}
}

func DeleteTestASNs(t *testing.T) {
	var header http.Header
	params := url.Values{}
	// Retrieve ASNs to delete
	resp, _, err := TOSession.GetASNsWithHeader(&params, header)
	assert.NoError(t, err, "Error trying to fetch ASNs for deletion: %v - resp: %+v", err, resp)
	for _, asn := range resp {
		_, _, err := TOSession.DeleteASNByASN(asn.ID)
		assert.NoError(t, err, "Cannot delete ASN by ASN number: '%v' %v", asn.ASN, err)

		// Retrieve the ASN to see if it got deleted
		params.Set("asn", strconv.Itoa(asn.ASN))
		asns, _, err := TOSession.GetASNsWithHeader(&params, header)
		assert.NoError(t, err, "Error deleting ASN: %s", err)
		assert.Equal(t, 0, len(asns), "Expected ASN: %v to be deleted", asn.ASN)
	}
}
