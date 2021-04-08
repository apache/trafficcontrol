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

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

const (
	// APIStatsSummary is the full path to the /stats_summary API endpoint.
	APIStatsSummary = "/stats_summary"
)

// GetSummaryStats gets a list of Summary Stats with the ability to filter on
// CDN, Delivery Service, and/or stat name.
func (to *Session) GetSummaryStats(cdn, deliveryService, statName *string) (tc.StatsSummaryResponse, toclientlib.ReqInf, error) {
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

	route := APIStatsSummary
	if len(param) > 0 {
		route = fmt.Sprintf("%s?%s", APIStatsSummary, param.Encode())
	}
	reqInf, err := to.get(route, nil, &resp)
	return resp, reqInf, err
}

// GetSummaryStatsLastUpdated gets the time at which Stat Summaries were last
// updated.
// If 'statName' isn't nil, the response will be limited to the stat thereby
// named.
func (to *Session) GetSummaryStatsLastUpdated(statName *string) (tc.StatsSummaryLastUpdatedResponse, toclientlib.ReqInf, error) {
	resp := tc.StatsSummaryLastUpdatedResponse{}

	param := url.Values{}
	param.Add("lastSummaryDate", "true")
	if statName != nil {
		param.Add("statName", *statName)
	}
	route := fmt.Sprintf("%s?%s", APIStatsSummary, param.Encode())

	reqInf, err := to.get(route, nil, &resp)
	return resp, reqInf, err
}

// CreateSummaryStats creates the given Stats Summary.
func (to *Session) CreateSummaryStats(statsSummary tc.StatsSummary) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(APIStatsSummary, statsSummary, nil, &alerts)
	return alerts, reqInf, err
}
