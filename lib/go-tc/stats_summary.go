package tc

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

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	validation "github.com/go-ozzo/ozzo-validation"
)

const dateFormat = "2006-01-02"

// StatsSummaryResponse is the structure of a response from Traffic Ops to
// GET requests made to its /stats_summary API endpoint.
type StatsSummaryResponse struct {
	Response []StatsSummary `json:"response"`
	Alerts
}

// StatsSummary is a summary of some kind of statistic for a CDN and/or
// Delivery Service.
type StatsSummary struct {
	CDNName         *string `json:"cdnName"  db:"cdn_name"`
	DeliveryService *string `json:"deliveryServiceName"  db:"deliveryservice_name"`
	// The name of the stat, which can be whatever the TO API client wants.
	StatName *string `json:"statName"  db:"stat_name"`
	// The value of the stat - this cannot actually be nil in a valid
	// StatsSummary.
	StatValue   *float64   `json:"statValue"  db:"stat_value"`
	SummaryTime time.Time  `json:"summaryTime"  db:"summary_time"`
	StatDate    *time.Time `json:"statDate"  db:"stat_date"`
}

// Validate implements the
// github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api.ParseValidator
// interface.
func (ss StatsSummary) Validate(tx *sql.Tx) error {
	errs := tovalidate.ToErrors(validation.Errors{
		"statName":  validation.Validate(ss.StatName, validation.Required),
		"statValue": validation.Validate(ss.StatValue, validation.Required),
	})
	return util.JoinErrs(errs)
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface with a
// customized decoding to force the date format on StatDate.
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

// MarshalJSON implements the encoding/json.Marshaler interface with a
// customized encoding to force the date format on StatDate.
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

// StatsSummaryLastUpdated is the type of the `response` property of a response
// from Traffic Ops to a GET request made to its /stats_summary endpoint when
// the 'lastSummaryDate' query string parameter is passed as 'true'.
type StatsSummaryLastUpdated struct {
	SummaryTime *time.Time `json:"summaryTime"  db:"summary_time"`
}

// MarshalJSON implements the encoding/json.Marshaler interface with a
// customized encoding to force the date format on SummaryTime.
func (ss StatsSummaryLastUpdated) MarshalJSON() ([]byte, error) {
	resp := struct {
		SummaryTime *string `json:"summaryTime"`
	}{}
	if ss.SummaryTime != nil {
		resp.SummaryTime = util.StrPtr(ss.SummaryTime.Format(TimeLayout))
	}
	return json.Marshal(&resp)
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface with a
// customized decoding to force the SummaryTime format.
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

// StatsSummaryLastUpdatedResponse is the type of a response from Traffic Ops
// to a GET request made to its /stats_summary endpoint when the
// 'lastSummaryDate' query string parameter is passed as 'true'.
//
// Deprecated: This structure includes an unknown field and drops Alerts
// returned by the API - use StatsSummaryLastUpdatedAPIResponse instead.
type StatsSummaryLastUpdatedResponse struct {
	// This field has unknown purpose and meaning - do not depend on its value
	// for anything.
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

// StatsSummaryV5 is an alias for the latest minor version for the major version 5.
type StatsSummaryV5 StatsSummaryV50

// StatsSummaryV50 is a summary of some kind of statistic for a CDN and/or
// Delivery Service.
type StatsSummaryV50 struct {
	CDNName         *string    `json:"cdnName"  db:"cdn_name"`
	DeliveryService *string    `json:"deliveryServiceName"  db:"deliveryservice_name"`
	StatName        *string    `json:"statName"  db:"stat_name"`
	StatValue       *float64   `json:"statValue"  db:"stat_value"`
	SummaryTime     time.Time  `json:"summaryTime"  db:"summary_time"`
	StatDate        *time.Time `json:"statDate"  db:"stat_date"`
}

// Validate implements the
// github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api.ParseValidator
// interface.
func (ss StatsSummaryV5) Validate(tx *sql.Tx) error {
	errs := tovalidate.ToErrors(validation.Errors{
		"statName":  validation.Validate(ss.StatName, validation.Required),
		"statValue": validation.Validate(ss.StatValue, validation.Required),
	})
	return util.JoinErrs(errs)
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface with a
// customized decoding to force the date format on StatDate.
func (ss *StatsSummaryV5) UnmarshalJSON(data []byte) error {
	type Alias StatsSummaryV5
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
		statDate, err := parseTimeV5(*resp.StatDate)
		if err != nil {
			return fmt.Errorf("invalid timestamp given for statDate: %v", err)
		}
		ss.StatDate = &statDate
	}

	ss.SummaryTime, err = parseTimeV5(resp.SummaryTime)
	if err != nil {
		return fmt.Errorf("invalid timestamp given for summaryTime: %v", err)
	}
	return nil
}

func parseTimeV5(ts string) (time.Time, error) {
	rt, err := time.Parse(time.RFC3339, ts)
	if err == nil {
		return rt, err
	}
	return time.Parse(dateFormat, ts)
}

// MarshalJSON implements the encoding/json.Marshaler interface with a
// customized encoding to force the date format on StatDate.
func (ss StatsSummaryV5) MarshalJSON() ([]byte, error) {
	type Alias StatsSummaryV5
	resp := struct {
		StatDate    *string `json:"statDate"`
		SummaryTime string  `json:"summaryTime"`
		Alias
	}{
		SummaryTime: ss.SummaryTime.Format(time.RFC3339),
		Alias:       (Alias)(ss),
	}
	if ss.StatDate != nil {
		resp.StatDate = util.Ptr(ss.StatDate.Format(dateFormat))
	}
	return json.Marshal(&resp)
}

// StatsSummaryResponseV5 is an alias for the latest minor version for the major version 5.
type StatsSummaryResponseV5 StatsSummaryResponseV50

// StatsSummaryResponseV50 is the structure of a response from Traffic Ops to
// GET requests made to its /stats_summary V5 API endpoint.
type StatsSummaryResponseV50 struct {
	Response []StatsSummaryV5 `json:"response"`
	Alerts
}

// StatsSummaryLastUpdatedV5 is an alias for the latest minor version for the major version 5.
type StatsSummaryLastUpdatedV5 StatsSummaryLastUpdatedV50

// StatsSummaryLastUpdatedV50 is the type of the `response` property of a response
// from Traffic Ops to a GET request made to its /stats_summary endpoint when
// the 'lastSummaryDate' query string parameter is passed as 'true'.
type StatsSummaryLastUpdatedV50 struct {
	SummaryTime *time.Time `json:"summaryTime"  db:"summary_time"`
}

// MarshalJSON implements the encoding/json.Marshaler interface with a
// customized encoding to force the date format on SummaryTime.
func (ss StatsSummaryLastUpdatedV5) MarshalJSON() ([]byte, error) {
	resp := struct {
		SummaryTime *string `json:"summaryTime"`
	}{}
	if ss.SummaryTime != nil {
		resp.SummaryTime = util.Ptr(ss.SummaryTime.Format(time.RFC3339))
	}
	return json.Marshal(&resp)
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface with a
// customized decoding to force the SummaryTime format.
func (ss *StatsSummaryLastUpdatedV5) UnmarshalJSON(data []byte) error {
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
		return err
	}
	return nil
}

// StatsSummaryLastUpdatedAPIResponseV5 is an alias for the latest minor version for the major version 5.
type StatsSummaryLastUpdatedAPIResponseV5 StatsSummaryLastUpdatedAPIResponseV50

// StatsSummaryLastUpdatedAPIResponseV50 is the type of a response from Traffic
// Ops to a request to its /stats_summary endpoint with the 'lastSummaryDate'
// query string parameter set to 'true'.
type StatsSummaryLastUpdatedAPIResponseV50 struct {
	Response StatsSummaryLastUpdatedV5 `json:"response"`
	Alerts
}
