package tc

import (
	"encoding/json"
	"time"
)

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

const dateFormat = "2006-01-02"

// StatsSummaryResponse ...
type StatsSummaryResponse struct {
	Response []StatsSummary `json:"response"`
}

// StatsSummary ...
type StatsSummary struct {
	ID              int       `json:"-"  db:"id"`
	CDNName         string    `json:"cdnName"  db:"cdn_name"`
	DeliveryService string    `json:"deliveryServiceName"  db:"deliveryservice_name"`
	StatName        string    `json:"statName"  db:"stat_name"`
	StatValue       float64   `json:"statValue"  db:"stat_value"`
	SummaryTime     time.Time `json:"summaryTime"  db:"summary_time"`
	StatDate        time.Time `json:"statDate"  db:"stat_date"`
}

// UnmarshalJSON Customized Unmarshal to force date format on statDate
func (ss *StatsSummary) UnmarshalJSON(data []byte) error {
	type Alias StatsSummary
	resp := struct {
		StatDate string `json:"statDate"`
		*Alias
	}{
		Alias: (*Alias)(ss),
	}
	err := json.Unmarshal(data, &resp)
	if err != nil {
		return err
	}
	ss.StatDate, err = time.Parse(dateFormat, resp.StatDate)
	if err != nil {
		return err
	}
	return nil
}

// UnmarshalJSON Customized Marshal to force date format on statDate
func (ss StatsSummary) MarshalJSON() ([]byte, error) {
	type Alias StatsSummary
	return json.Marshal(&struct {
		StatDate string `json:"statDate"`
		Alias
	}{
		StatDate: ss.StatDate.Format(dateFormat),
		Alias:    (Alias)(ss),
	})
}

// StatsSummaryLastUpdated ...
type StatsSummaryLastUpdated struct {
	SummaryTime time.Time `json:"summaryTime"  db:"summary_time"`
}

// StatsSummaryLastUpdatedResponse ...
type StatsSummaryLastUpdatedResponse struct {
	Response StatsSummaryLastUpdated `json:"response"`
}
