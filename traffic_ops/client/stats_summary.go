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
	"time"
)

type CacheGroupResponse struct {
	Version  string         `json:"version"`
	Response []StatsSummary `json:"response"`
}

type StatsSummary struct {
	CdnName          string     `json:"cdnName"`
	DeliveryService  string     `json:"dsName"`
	StatName         string     `json:"statName"`
	StatValue        float64    `json:"statValue"`
	SummaryTimestamp *time.Time `json:"summaryTimestamp"`
}

func (to *Session) SummaryStats(cdn string, deliveryService string, statName string, startDate *time.Time, endDate *time.Time, limit uint64) ([]CacheGroup, error) {
	body, err := to.getBytes("/api/1.2/summary_stats.json")
	if err != nil {
		return nil, err
	}
	cgList, err := cgUnmarshall(body)
	return cgList.Response, err
}

func (to *Session) SummaryStatsLastSummary(statName string) ([]CacheGroup, error) {
	body, err := to.getBytes("/api/1.1/cachegroups.json")
	if err != nil {
		return nil, err
	}
	cgList, err := cgUnmarshall(body)
	return cgList.Response, err
}

func cgUnmarshall(body []byte) (CacheGroupResponse, error) {
	var data CacheGroupResponse
	err := json.Unmarshal(body, &data)
	return data, err
}
