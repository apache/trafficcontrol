package tc

import (
	"github.com/apache/trafficcontrol/v8/lib/go-util"
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

// LogsResponse is a list of Logs as a response.
type LogsResponse struct {
	Response []Log `json:"response"`
	Alerts
}

// LogsResponseV5 is a list of Logs as a response, for the latest minor version of api 5.x.
type LogsResponseV5 = LogsResponseV50

// LogsResponseV50 is a list of Logs as a response, for api version 5.0.
type LogsResponseV50 struct {
	Response []LogV50 `json:"response"`
	Alerts
}

// Log contains a change that has been made to the Traffic Control system.
type Log struct {
	ID          *int    `json:"id"`
	LastUpdated *Time   `json:"lastUpdated"`
	Level       *string `json:"level"`
	Message     *string `json:"message"`
	TicketNum   *int    `json:"ticketNum"`
	User        *string `json:"user"`
}

// Upgrade changes a Log structure into a LogV5 structure
func (l Log) Upgrade() LogV5 {
	logV5 := LogV5{
		ID:        l.ID,
		Level:     l.Level,
		Message:   l.Message,
		TicketNum: l.TicketNum,
		User:      l.User,
	}
	if l.LastUpdated != nil {
		t := &l.LastUpdated.Time
		if t != nil {
			t, _ := util.ConvertTimeFormat(*t, time.RFC3339)
			logV5.LastUpdated = t
		}
	}
	return logV5
}

// LogV50 contains a change that has been made to the Traffic Control system, for api version 5.0.
type LogV50 struct {
	ID          *int       `json:"id"`
	LastUpdated *time.Time `json:"lastUpdated"`
	Level       *string    `json:"level"`
	Message     *string    `json:"message"`
	TicketNum   *int       `json:"ticketNum"`
	User        *string    `json:"user"`
}

// LogV5 is the Log structure used by the latest 5.x API version
type LogV5 = LogV50

// NewLogCountResp is the response returned when the total number of new changes
// made to the Traffic Control system is requested. "New" means since the last
// time this information was requested.
type NewLogCountResp struct {
	NewLogCount uint64 `json:"newLogcount"`
}
