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
	"strconv"

	tc "github.com/apache/trafficcontrol/lib/go-tc"
)

// SummaryStats ...
// Deprecated: use GetSummaryStats
func (to *Session) SummaryStats(cdn string, deliveryService string, statName string) ([]tc.StatsSummary, error) {
	ss, _, err := to.GetSummaryStats(cdn, deliveryService, statName)
	return ss, err
}

func (to *Session) GetSummaryStats(cdn string, deliveryService string, statName string) ([]tc.StatsSummary, ReqInf, error) {
	var queryParams []string
	if len(cdn) > 0 {
		queryParams = append(queryParams, fmt.Sprintf("cdnName=%s", cdn))
	}
	if len(deliveryService) > 0 {
		queryParams = append(queryParams, fmt.Sprintf("deliveryServiceName=%s", deliveryService))
	}
	if len(statName) > 0 {
		queryParams = append(queryParams, fmt.Sprintf("statName=%s", statName))
	}
	queryURL := "/api/1.2/stats_summary.json"
	queryParamString := "?"
	if len(queryParams) > 0 {
		for i, param := range queryParams {
			if i == 0 {
				queryParamString += param
			} else {
				queryParamString += fmt.Sprintf("&%s", param)
			}
		}
		queryURL += queryParamString
	}

	resp, remoteAddr, err := to.request("GET", queryURL, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.StatsSummaryResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// SummaryStatsLastUpdated ...
// Deprecated: use GetSummaryStatsLastUpdated
func (to *Session) SummaryStatsLastUpdated(statName string) (string, error) {
	s, _, err := to.GetSummaryStatsLastUpdated(statName)
	return s, err
}

func (to *Session) GetSummaryStatsLastUpdated(statName string) (string, ReqInf, error) {
	queryURL := "/api/1.2/stats_summary.json?lastSummaryDate=true"
	if len(statName) > 0 {
		queryURL += fmt.Sprintf("?statName=%s", statName)
	}

	resp, remoteAddr, err := to.request("GET", queryURL, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return "", reqInf, err
	}
	defer resp.Body.Close()

	var data tc.LastUpdated
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", reqInf, err
	}

	if len(data.Response.SummaryTime) > 0 {
		return data.Response.SummaryTime, reqInf, nil
	}
	t := "1970-01-01 00:00:00"
	return t, reqInf, nil
}

// AddSummaryStats ...
// Deprecated: use DoAddSummaryStats
func (to *Session) AddSummaryStats(statsSummary tc.StatsSummary) error {
	_, err := to.DoAddSummaryStats(statsSummary)
	return err
}
func (to *Session) DoAddSummaryStats(statsSummary tc.StatsSummary) (ReqInf, error) {
	reqBody, err := json.Marshal(statsSummary)
	if err != nil {
		return ReqInf{}, err
	}

	url := "/api/1.2/stats_summary/create"
	resp, remoteAddr, err := to.request("POST", url, reqBody)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return reqInf, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err := fmt.Errorf("Response code = %s and Status = %s", strconv.Itoa(resp.StatusCode), resp.Status)
		return reqInf, err
	}
	return reqInf, nil
}
