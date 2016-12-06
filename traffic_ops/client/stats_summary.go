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
)

// StatsSummaryResponse ...
type StatsSummaryResponse struct {
	Response []StatsSummary `json:"response"`
}

// StatsSummary ...
type StatsSummary struct {
	CDNName         string `json:"cdnName"`
	DeliveryService string `json:"deliveryServiceName"`
	StatName        string `json:"statName"`
	StatValue       string `json:"statValue"`
	SummaryTime     string `json:"summaryTime"`
	StatDate        string `json:"statDate"`
}

// LastUpdated ...
type LastUpdated struct {
	Version  string `json:"version"`
	Response struct {
		SummaryTime string `json:"summaryTime"`
	} `json:"response"`
}

// SummaryStats ...
func (to *Session) SummaryStats(cdn string, deliveryService string, statName string) ([]StatsSummary, error) {
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

	resp, err := to.request("GET", queryURL, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data StatsSummaryResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data.Response, nil
}

// SummaryStatsLastUpdated ...
func (to *Session) SummaryStatsLastUpdated(statName string) (string, error) {
	queryURL := "/api/1.2/stats_summary.json?lastSummaryDate=true"
	if len(statName) > 0 {
		queryURL += fmt.Sprintf("?statName=%s", statName)
	}

	resp, err := to.request("GET", queryURL, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data LastUpdated
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	if len(data.Response.SummaryTime) > 0 {
		return data.Response.SummaryTime, nil
	}
	t := "1970-01-01 00:00:00"
	return t, nil
}

// AddSummaryStats ...
func (to *Session) AddSummaryStats(statsSummary StatsSummary) error {
	reqBody, err := json.Marshal(statsSummary)
	if err != nil {
		return err
	}

	url := "/api/1.2/stats_summary/create"
	resp, err := to.request("POST", url, reqBody)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		err := fmt.Errorf("Response code = %s and Status = %s", strconv.Itoa(resp.StatusCode), resp.Status)
		return err
	}
	return nil
}
