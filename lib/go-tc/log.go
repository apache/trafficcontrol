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

// LogsResponse is a list of Logs as a response.
type LogsResponse struct {
	Response []Log `json:"response"`
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

// NewLogCountResp is the response returned when the total number of new changes
// made to the Traffic Control system is requested. "New" means since the last
// time this information was requested.
type NewLogCountResp struct {
	NewLogCount uint64 `json:"newLogcount"`
}
