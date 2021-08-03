package tc

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/lib/go-util"
	validation "github.com/go-ozzo/ozzo-validation"
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
	Alerts
}

// StatsSummary ...
type StatsSummary struct {
	CDNName         *string    `json:"cdnName"  db:"cdn_name"`
	DeliveryService *string    `json:"deliveryServiceName"  db:"deliveryservice_name"`
	StatName        *string    `json:"statName"  db:"stat_name"`
	StatValue       *float64   `json:"statValue"  db:"stat_value"`
	SummaryTime     time.Time  `json:"summaryTime"  db:"summary_time"`
	StatDate        *time.Time `json:"statDate"  db:"stat_date"`
}

func (ss StatsSummary) Validate(tx *sql.Tx) error {
	errs := tovalidate.ToErrors(validation.Errors{
		"statName":  validation.Validate(ss.StatName, validation.Required),
		"statValue": validation.Validate(ss.StatValue, validation.Required),
	})
	return util.JoinErrs(errs)
}

// UnmarshalJSON Customized Unmarshal to force date format on statDate
func (ss *StatsSummary) UnmarshalJSON(data []byte) error {
	type Alias StatsSummary
	resp := struct {
		SummaryTime string  `json:"summaryTime"`
		StatDate    *string `json:"statDate"`
		*Alias
	}{
		Alias: (*Alias)(ss),
	}
	err := json.Unmarshal(data, &resp)
	if err != nil {
		return err
	}
	if resp.StatDate != nil {
		statDate, err := parseTime(*resp.StatDate)
		if err != nil {
			return errors.New("invalid timestamp given for statDate")
		}
		ss.StatDate = &statDate
	}

	ss.SummaryTime, err = parseTime(resp.SummaryTime)
	if err != nil {
		return errors.New("invalid timestamp given for summaryTime")
	}
	return nil
}

func parseTime(ts string) (time.Time, error) {
	rt, err := time.Parse(time.RFC3339, ts)
	if err == nil {
		return rt, err
	}
	rt, err = time.Parse(TimeLayout, ts)
	if err == nil {
		return rt, err
	}
	return time.Parse(dateFormat, ts)
}

// MarshalJSON Customized Marshal to force date format on statDate
func (ss StatsSummary) MarshalJSON() ([]byte, error) {
	type Alias StatsSummary
	resp := struct {
		StatDate    *string `json:"statDate"`
		SummaryTime string  `json:"summaryTime"`
		Alias
	}{
		SummaryTime: ss.SummaryTime.Format(TimeLayout),
		Alias:       (Alias)(ss),
	}
	if ss.StatDate != nil {
		resp.StatDate = util.StrPtr(ss.StatDate.Format(dateFormat))
	}
	return json.Marshal(&resp)
}

// StatsSummaryLastUpdated ...
type StatsSummaryLastUpdated struct {
	SummaryTime *time.Time `json:"summaryTime"  db:"summary_time"`
}

// MarshalJSON is a customized marshal to force date format on SummaryTime.
func (ss StatsSummaryLastUpdated) MarshalJSON() ([]byte, error) {
	resp := struct {
		SummaryTime *string `json:"summaryTime"`
	}{}
	if ss.SummaryTime != nil {
		resp.SummaryTime = util.StrPtr(ss.SummaryTime.Format(TimeLayout))
	}
	return json.Marshal(&resp)
}

// UnmarshalJSON Customized Unmarshal to force timestamp format
func (ss *StatsSummaryLastUpdated) UnmarshalJSON(data []byte) error {
	resp := struct {
		SummaryTime *string `json:"summaryTime"`
	}{}
	err := json.Unmarshal(data, &resp)
	if err != nil {
		return err
	}
	if resp.SummaryTime != nil {
		var summaryTime time.Time
		summaryTime, err = time.Parse(time.RFC3339, *resp.SummaryTime)
		if err == nil {
			ss.SummaryTime = &summaryTime
			return nil
		}
		summaryTime, err = time.Parse(TimeLayout, *resp.SummaryTime)
		ss.SummaryTime = &summaryTime
		return err
	}
	return nil
}

type StatsSummaryLastUpdatedResponse struct {
	Version  string                  `json:"version"`
	Response StatsSummaryLastUpdated `json:"response"`
}

// StatsSummaryLastUpdatedAPIResponse is the type of a response from Traffic
// Ops to a request to its /stats_summary endpoint with the 'lastSummaryDate'
// query string parameter set to 'true'.
type StatsSummaryLastUpdatedAPIResponse struct {
	Response StatsSummaryLastUpdated `json:"response"`
	Alerts
}
