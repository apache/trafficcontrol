package client

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
	"net/url"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

// apiStatsSummary is the full path to the /stats_summary API endpoint.
const apiStatsSummary = "/stats_summary"

// GetSummaryStats gets a list of Summary Stats with the ability to filter on
// CDN, Delivery Service, and/or stat name.
func (to *Session) GetSummaryStats(opts RequestOptions) (tc.StatsSummaryResponse, toclientlib.ReqInf, error) {
	var resp tc.StatsSummaryResponse
	reqInf, err := to.get(apiStatsSummary, opts, &resp)
	return resp, reqInf, err
}

// GetSummaryStatsLastUpdated gets the time at which Stat Summaries were last
// updated.
// If 'statName' isn't nil, the response will be limited to the stat thereby
// named.
func (to *Session) GetSummaryStatsLastUpdated(opts RequestOptions) (tc.StatsSummaryLastUpdatedAPIResponse, toclientlib.ReqInf, error) {
	if opts.QueryParameters == nil {
		opts.QueryParameters = url.Values{}
	}
	opts.QueryParameters.Set("lastSummaryDate", "true")

	var resp tc.StatsSummaryLastUpdatedAPIResponse
	reqInf, err := to.get(apiStatsSummary, opts, &resp)
	return resp, reqInf, err
}

// CreateSummaryStats creates the given Stats Summary.
func (to *Session) CreateSummaryStats(statsSummary tc.StatsSummary, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(apiStatsSummary, opts, statsSummary, &alerts)
	return alerts, reqInf, err
}
