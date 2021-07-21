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
	"time"
)

// AsyncStatus represents an async job status.
type AsyncStatus struct {
	// Id is the integral, unique identifier for the asynchronous job status.
	Id int `json:"id,omitempty" db:"id"`
	// Status is the status of the asynchronous job. This will be PENDING, SUCCEEDED, or FAILED.
	Status string `json:"status,omitempty" db:"status"`
	// StartTime is the time the asynchronous job was started.
	StartTime time.Time `json:"start_time,omitempty" db:"start_time"`
	// EndTime is the time the asynchronous job completed. This will be null if it has not completed yet.
	EndTime *time.Time `json:"end_time,omitempty" db:"end_time"`
	// Message is the message about the job status.
	Message *string `json:"message,omitempty" db:"message"`
}

// AsyncStatusResponse represents the response from the GET /async_status/{id} API.
type AsyncStatusResponse struct {
	Response AsyncStatus `json:"response"`
	Alerts
}
