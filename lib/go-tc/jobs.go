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
	"encoding/json"
	"errors"
	"time"
)

// Job represents a content invalidation job as stored in the database.
//
// Deprecated: Use InvalidationJob instead, as it's more flexible.
type Job struct {
	Parameters      string `json:"parameters"`
	Keyword         string `json:"keyword"`
	AssetURL        string `json:"assetUrl"`
	CreatedBy       string `json:"createdBy"`
	StartTime       string `json:"startTime"`
	ID              int64  `json:"id"`
	DeliveryService string `json:"deliveryService"`
}

// JobRequest contains the data to create a job.
// Note this is a convenience struct for posting users; the actual JSON object is a JobRequestAPI
//
// Deprecated: Use InvalidationJobInput instead, as it's more flexible.
type JobRequest struct {
	TTL               time.Duration
	StartTime         time.Time
	DeliveryServiceID int64
	Regex             string
	Urgent            bool
}

// JobRequestTimeFormat is a Go reference time format, for use with time.Format, of the format
// used by the Perl version of the /jobs and /user/current/jobs endpoints.
//
// Deprecated: Since job inputs no longer strictly require this format, an RFC 3339 format or a
// timestamp in milliseconds should be used instead.
const JobRequestTimeFormat = `2006-01-02 15:04:05`

// JobTimeFormat is a Go reference time format, for use with time.Format.
//
// Deprecated: this format is the same as TimeLayout - which should be used instead.
const JobTimeFormat = `2006-01-02 15:04:05-07`

// MarshalJSON implements the encoding/json.Marshaler interface.
func (jr JobRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(JobRequestAPI{
		TTLSeconds: int64(jr.TTL / time.Second),
		StartTime:  jr.StartTime.Format(JobRequestTimeFormat),
		DSID:       jr.DeliveryServiceID,
		Regex:      jr.Regex,
		Urgent:     jr.Urgent,
	})
}

// UnmarshalJSON implements the encoding/json Unmarshaler interface.
func (jr *JobRequest) UnmarshalJSON(b []byte) error {
	jri := JobRequestAPI{}
	if err := json.Unmarshal(b, &jri); err != nil {
		return err
	}
	startTime, err := time.Parse(JobRequestTimeFormat, jri.StartTime)
	if err != nil {
		return errors.New("startTime '" + jri.StartTime + "' is not of the required format '" + JobRequestTimeFormat + "'")
	}
	*jr = JobRequest{
		TTL:               time.Duration(jri.TTLSeconds) * time.Second,
		StartTime:         startTime,
		DeliveryServiceID: jri.DSID,
		Regex:             jri.Regex,
		Urgent:            jri.Urgent,
	}
	return nil
}

// JobRequestAPI represents the JSON input accepted by the API for creating new
// content invalidation jobs.
//
// Deprecated: This structure is technically incorrect, and has been superseded by the
// InvalidationJobInput structure.
type JobRequestAPI struct {
	TTLSeconds int64  `json:"ttl"`
	StartTime  string `json:"startTime"`
	DSID       int64  `json:"dsId"`
	Regex      string `json:"regex"`
	Urgent     bool   `json:"urgent"`
}
