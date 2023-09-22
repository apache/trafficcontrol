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
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

var latestTime time.Time

func TestStatsSummary(t *testing.T) {

	CreateTestStatsSummaries(t)

	methodTests := utils.TestCase[client.Session, client.RequestOptions, tc.StatsSummaryV5]{
		"GET": {
			"OK when VALID request": {
				ClientSession: TOSession,
				Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
			},
			"OK when VALID STATNAME parameter": {
				ClientSession: TOSession,
				RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"statName": {"daily_bytesserved"}}},
				Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
					validateStatsSummaryFields(map[string]interface{}{"StatName": "daily_bytesserved"})),
			},
			"OK when VALID CDNNAME parameter": {
				ClientSession: TOSession,
				RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"cdnName": {"cdn1"}}},
				Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(2),
					validateStatsSummaryFields(map[string]interface{}{"CDNName": "cdn1"})),
			},
			"OK when VALID DELIVERYSERVICENAME parameter": {
				ClientSession: TOSession,
				RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"deliveryServiceName": {"all"}}},
				Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(2),
					validateStatsSummaryFields(map[string]interface{}{"DeliveryService": "all"})),
			},
			"OK when VALID LASTSUMMARYDATE parameter": {
				ClientSession: TOSession,
				RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"statName": {"daily_bytesserved"}}},
				Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateStatsSummaryLastUpdatedField(latestTime)),
			},
			"EMPTY RESPONSE when NON-EXISTENT STATNAME": {
				ClientSession: TOSession,
				RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"statName": {"bogus"}}},
				Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
			},
			"EMPTY RESPONSE when NON-EXISTENT DELIVERYSERVICENAME": {
				ClientSession: TOSession,
				RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"deliveryServiceName": {"bogus"}}},
				Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
			},
			"EMPTY RESPONSE when NON-EXISTENT CDNNAME": {
				ClientSession: TOSession,
				RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"cdnName": {"bogus"}}},
				Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
			},
		},
	}

	for method, testCases := range methodTests {
		t.Run(method, func(t *testing.T) {
			for name, testCase := range testCases {
				switch method {
				case "GET":
					t.Run(name, func(t *testing.T) {
						if name == "OK when VALID LASTSUMMARYDATE parameter" {
							resp, reqInf, err := testCase.ClientSession.GetSummaryStatsLastUpdated(testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						} else {
							resp, reqInf, err := testCase.ClientSession.GetSummaryStats(testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						}
					})
				case "POST":
					t.Run(name, func(t *testing.T) {
						alerts, reqInf, err := testCase.ClientSession.CreateSummaryStats(testCase.RequestBody, testCase.RequestOpts)
						for _, check := range testCase.Expectations {
							check(t, reqInf, nil, alerts, err)
						}
					})
				}
			}
		})
	}
}

func validateStatsSummaryFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Stats Summary response to not be nil.")
		statsSummaryResp := resp.([]tc.StatsSummaryV5)
		for field, expected := range expectedResp {
			for _, statsSummary := range statsSummaryResp {
				switch field {
				case "CDNName":
					assert.RequireNotNil(t, statsSummary.CDNName, "Expected CDNName to not be nil.")
					assert.Equal(t, expected, *statsSummary.CDNName, "Expected CDNName to be %v, but got %s", expected, *statsSummary.CDNName)
				case "DeliveryService":
					assert.RequireNotNil(t, statsSummary.DeliveryService, "Expected DeliveryService to not be nil.")
					assert.Equal(t, expected, *statsSummary.DeliveryService, "Expected DeliveryService to be %v, but got %s", expected, *statsSummary.DeliveryService)
				case "StatName":
					assert.RequireNotNil(t, statsSummary.StatName, "Expected StatName to not be nil.")
					assert.Equal(t, expected, *statsSummary.StatName, "Expected StatName to be %v, but got %s", expected, *statsSummary.StatName)
				case "SummaryTime":
					assert.Equal(t, true, expected.(time.Time).Equal(statsSummary.SummaryTime), "Expected SummaryTime to be %v, but got %v", expected, statsSummary.SummaryTime)
				default:
					t.Errorf("Expected field: %v, does not exist in response", field)
				}
			}
		}
	}
}

func validateStatsSummaryLastUpdatedField(expectedTime time.Time) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected StatsSummaryLastUpdated response to not be nil.")
		statsSummaryLastUpdated := resp.(tc.StatsSummaryLastUpdatedV5)
		assert.RequireNotNil(t, statsSummaryLastUpdated.SummaryTime, "Expected SummaryTime to not be nil.")
		assert.Equal(t, true, expectedTime.Equal(*statsSummaryLastUpdated.SummaryTime), "Expected SummaryTime to be %v, but got %v", expectedTime, *statsSummaryLastUpdated.SummaryTime)

	}
}

// Note that these stats summaries are never cleaned up, and will be left in
// the TODB after the tests complete
func CreateTestStatsSummaries(t *testing.T) {
	for _, ss := range testData.StatsSummaries {
		latestTime = time.Now().Truncate(time.Second)
		ss.SummaryTime = latestTime
		alerts, _, err := TOSession.CreateSummaryStats(ss, client.RequestOptions{})
		assert.RequireNoError(t, err, "Creating Stats Summary for stat '%s': %v - alerts: %+v", *ss.StatName, err, alerts.Alerts)
	}
}
