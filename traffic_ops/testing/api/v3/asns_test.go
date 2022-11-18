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
)

func TestASN(t *testing.T) {
	WithObjs(t, []TCObj{Types, CacheGroups, ASN}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.V3TestCase{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfModifiedSince: {tomorrow}},
					Expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when CHANGES made": {
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfModifiedSince: {currentTimeRFC}},
					Expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateASNsSort()),
				},
				"OK when VALID ASN PARAMETER": {
					ClientSession: TOSession,
					RequestParams: url.Values{"asn": {"9999"}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						validateASNsFields(map[string]interface{}{"ASN": 9999, "Cachegroup": "multiOriginCachegroup"})),
				},
			},
			"PUT": {
				"OK when VALID request": {
					ClientSession: TOSession, EndpointID: GetASNID(t, "8888"),
					RequestBody: map[string]interface{}{
						"asn":            7777,
						"cachegroupName": "originCachegroup",
						"cachegroupId":   GetCacheGroupId(t, "originCachegroup")(),
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateASNsUpdateCreateFields("7777", map[string]interface{}{"ASN": 7777})),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					asn := tc.ASN{}

					if testCase.RequestBody != nil {
						dat, err := json.Marshal(testCase.RequestBody)
						assert.NoError(t, err, "Error occurred when marshalling request body: %v", err)
						err = json.Unmarshal(dat, &asn)
						assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
					}

					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetASNsWithHeader(&testCase.RequestParams, testCase.RequestHeaders)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp, tc.Alerts{}, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.UpdateASNByID(testCase.EndpointID(), asn)
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

func validateASNsFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected ASN response to not be nil.")
		asnResp := resp.([]tc.ASN)
		for field, expected := range expectedResp {
			for _, asn := range asnResp {
				switch field {
				case "ASN":
					assert.Equal(t, expected, asn.ASN, "Expected ASN to be %v, but got %d", expected, asn.ASN)
				case "Cachegroup":
					assert.Equal(t, expected, asn.Cachegroup, "Expected Cachegroup to be %v, but got %s", expected, asn.Cachegroup)
				case "CachegroupID":
					assert.Equal(t, expected, asn.CachegroupID, "Expected CachegroupID to be %v, but got %d", expected, asn.CachegroupID)
				case "ID":
					assert.Equal(t, expected, asn.ID, "Expected ID to be %v, but got %d", expected, asn.ID)
				default:
					t.Errorf("Expected field: %v, does not exist in response", field)
				}
			}
		}
	}
}

func validateASNsUpdateCreateFields(asn string, expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		params := url.Values{"asn": {asn}}
		asnResp, _, err := TOSession.GetASNsWithHeader(&params, nil)
		assert.RequireNoError(t, err, "Error getting ASN: %v", err)
		assert.RequireEqual(t, 1, len(asnResp), "Expected one ASN returned Got: %d", len(asnResp))
		validateASNsFields(expectedResp)(t, toclientlib.ReqInf{}, asnResp, tc.Alerts{}, nil)
	}
}

func validateASNsSort() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected ASN response to not be nil.")
		var asns []int
		asnResp := resp.([]tc.ASN)
		for _, asn := range asnResp {
			asns = append(asns, asn.ASN)
		}
		assert.Equal(t, true, sort.IntsAreSorted(asns), "List is not sorted by their ids: %v", asns)
	}
}

func GetASNID(t *testing.T, asn string) func() int {
	return func() int {
		params := url.Values{"asn": {asn}}
		asnResp, _, err := TOSession.GetASNsWithHeader(&params, nil)
		assert.RequireNoError(t, err, "Get ASNs Request failed with error: %v", err)
		assert.RequireEqual(t, len(asnResp), 1, "Expected response object length 1, but got %d", len(asnResp))
		return asnResp[0].ID
	}
}

func CreateTestASNs(t *testing.T) {
	for _, asn := range testData.ASNs {
		asn.CachegroupID = GetCacheGroupId(t, asn.Cachegroup)()
		resp, _, err := TOSession.CreateASN(asn)
		assert.RequireNoError(t, err, "Could not create ASN: %v - alerts: %+v", err, resp)
	}
}

func DeleteTestASNs(t *testing.T) {
	params := url.Values{}
	asns, _, err := TOSession.GetASNsWithHeader(&params, nil)
	assert.NoError(t, err, "Error trying to fetch ASNs for deletion: %v", err)

	for _, asn := range asns {
		alerts, _, err := TOSession.DeleteASNByASN(asn.ID)
		assert.NoError(t, err, "Cannot delete ASN %d: %v - alerts: %+v", asn.ASN, err, alerts)
		// Retrieve the ASN to see if it got deleted
		params.Set("asn", strconv.Itoa(asn.ASN))
		getAsns, _, err := TOSession.GetASNsWithHeader(&params, nil)
		assert.NoError(t, err, "Error trying to fetch ASN after deletion: %v", err)
		assert.Equal(t, 0, len(getAsns), "Expected ASN %d to be deleted, but it was found in Traffic Ops", asn.ASN)
	}
}
