/*
   Copyright 2015 Comcast Cable Communications Management, LLC
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
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
)

type StatsSummaryResponse struct {
	Version  string         `json:"version"`
	Response []StatsSummary `json:"response"`
}

type StatsSummary struct {
	CdnName         string `json:"cdnName"`
	DeliveryService string `json:"deliveryServiceName"`
	StatName        string `json:"statName"`
	StatValue       string `json:"statValue"`
	SummaryTime     string `json:"summaryTime"`
}

type LastUpdated struct {
	Version  string `json:"version"`
	Response struct {
		SummaryTime string `json:"summaryTime"`
	} `json:"response"`
}

func (to *Session) SummaryStats(cdn string, deliveryService string, statName string) ([]StatsSummary, error) {
	var queryParams []string
	if len(cdn) > 0 {
		queryParams = append(queryParams, "cdnName="+cdn)
	}
	if len(deliveryService) > 0 {
		queryParams = append(queryParams, "deliveryServiceName="+deliveryService)
	}
	if len(statName) > 0 {
		queryParams = append(queryParams, "statName="+statName)
	}
	queryUrl := "/api/1.2/stats_summary.json"
	queryParamString := "?"
	if len(queryParams) > 0 {
		for i, param := range queryParams {
			if i == 0 {
				queryParamString += param
			} else {
				queryParamString += "&" + param
			}
		}
		queryUrl += queryParamString
	}
	body, err := to.getBytes(queryUrl)
	if err != nil {
		return nil, err
	}
	ssList, err := ssUnmarshall(body)
	return ssList.Response, err
}

func (to *Session) SummaryStatsLastUpdated(statName string) (string, error) {
	queryUrl := "/api/1.2/stats_summary.json?lastSummaryDate=true"
	if len(statName) > 0 {
		queryUrl += "?statName=" + statName
	}
	body, err := to.getBytes(queryUrl)
	if err != nil {
		return "", err
	}
	var data LastUpdated
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Printf("err is %v\n", err)
		return "", err
	}
	if len(data.Response.SummaryTime) > 0 {
		return data.Response.SummaryTime, nil
	} else {
		return "1970-01-01 00:00:00", nil
	}
}

func (to *Session) AddSummaryStats(statsSummary StatsSummary) (string, error) {
	body, err := json.Marshal(statsSummary)
	if err != nil {
		return "", err
	}
	response, err := to.postJson("/api/1.2/stats_summary/create", body)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	if response.StatusCode != 200 {
		err := errors.New("Response code = " + strconv.Itoa(response.StatusCode) + "and response body = " + string(body))
		return "", err
	}
	return string(body), nil
}

func ssUnmarshall(body []byte) (StatsSummaryResponse, error) {
	var data StatsSummaryResponse
	err := json.Unmarshal(body, &data)
	return data, err
}
