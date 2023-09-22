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
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	totest "github.com/apache/trafficcontrol/v8/lib/go-tc/totestv4"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"
)

func TestASN(t *testing.T) {
	WithObjs(t, []TCObj{Types, CacheGroups, ASN}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.V4TestCase{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {tomorrow}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {currentTimeRFC}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateASNsSort()),
				},
				"OK when VALID ASN PARAMETER": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"asn": {"9999"}}},
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
						"cachegroupId":   totest.GetCacheGroupId(t, TOSession, "originCachegroup")(),
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
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("asn", asn)
		asnResp, _, err := TOSession.GetASNs(opts)
		assert.RequireNoError(t, err, "Error getting ASN: %v - alerts: %+v", err, asnResp.Alerts)
		assert.RequireEqual(t, 1, len(asnResp.Response), "Expected one ASN returned Got: %d", len(asnResp.Response))
		validateASNsFields(expectedResp)(t, toclientlib.ReqInf{}, asnResp.Response, tc.Alerts{}, nil)
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
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("asn", asn)
		resp, _, err := TOSession.GetASNs(opts)
		assert.RequireNoError(t, err, "Get ASNs Request failed with error: %v", err)
		assert.RequireEqual(t, len(resp.Response), 1, "Expected response object length 1, but got %d", len(resp.Response))
		return resp.Response[0].ID
	}
}
