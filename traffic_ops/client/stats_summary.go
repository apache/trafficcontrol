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
	"encoding/json"
	"fmt"
	"net"
	"net/url"

	tc "github.com/apache/trafficcontrol/lib/go-tc"
)

const (
	basePathSummaryStats = apiBase + "/stats_summary"
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

	url := basePathSummaryStats
	if len(param) > 0 {
		url = fmt.Sprintf("%s?%s", basePathSummaryStats, param.Encode())
	}
	reqInf, err := get(to, url, &resp)
	if err != nil {
		return resp, reqInf, err
	}
	return resp, reqInf, nil
}

// GetSummaryStatsLastUpdated time of the last summary for a given stat
func (to *Session) GetSummaryStatsLastUpdated(statName *string) (tc.StatsSummaryLastUpdatedResponse, ReqInf, error) {
	resp := tc.StatsSummaryLastUpdatedResponse{}

	param := url.Values{}
	param.Add("lastSummaryDate", "true")
	if statName != nil {
		param.Add("statName", *statName)
	}
	url := fmt.Sprintf("%s?%s", basePathSummaryStats, param.Encode())

	reqInf, err := get(to, url, &resp)
	if err != nil {
		return resp, reqInf, err
	}
	return resp, reqInf, nil
}

// CreateSummaryStats creates a stats summary
func (to *Session) CreateSummaryStats(statsSummary tc.StatsSummary) (tc.Alerts, ReqInf, error) {
	var alerts tc.Alerts
	var remoteAddr net.Addr
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	reqBody, err := json.Marshal(statsSummary)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	reqInf, err = post(to, fmt.Sprintf("%s/%s", basePathSummaryStats, "create"), reqBody, &alerts)
	return alerts, reqInf, err
}
