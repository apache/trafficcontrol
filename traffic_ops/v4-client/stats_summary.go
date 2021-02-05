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
	"net/url"

	tc "github.com/apache/trafficcontrol/lib/go-tc"
)

const (
	API_STATS_SUMMARY = apiBase + "/stats_summary"
)

// GetSummaryStats gets a list of summary stats with the ability to filter on cdn,deliveryService and/or stat
func (to *Session) GetSummaryStats(cdn, deliveryService, statName *string) (tc.StatsSummaryResponse, ReqInf, error) {
	resp := tc.StatsSummaryResponse{}

	param := url.Values{}
	if cdn != nil {
		param.Add("cdnName", *cdn)
	}
	if deliveryService != nil {
		param.Add("deliveryServiceName", *deliveryService)
	}
	if statName != nil {
		param.Add("statName", *statName)
	}

	route := API_STATS_SUMMARY
	if len(param) > 0 {
		route = fmt.Sprintf("%s?%s", API_STATS_SUMMARY, param.Encode())
	}
	reqInf, err := to.get(route, nil, &resp)
	return resp, reqInf, err
}

// GetSummaryStatsLastUpdated time of the last summary for a given stat
func (to *Session) GetSummaryStatsLastUpdated(statName *string) (tc.StatsSummaryLastUpdatedResponse, ReqInf, error) {
	resp := tc.StatsSummaryLastUpdatedResponse{}

	param := url.Values{}
	param.Add("lastSummaryDate", "true")
	if statName != nil {
		param.Add("statName", *statName)
	}
	route := fmt.Sprintf("%s?%s", API_STATS_SUMMARY, param.Encode())

	reqInf, err := to.get(route, nil, &resp)
	return resp, reqInf, err
}

// CreateSummaryStats creates a stats summary
func (to *Session) CreateSummaryStats(statsSummary tc.StatsSummary) (tc.Alerts, ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(API_STATS_SUMMARY, statsSummary, nil, &alerts)
	return alerts, reqInf, err
}
