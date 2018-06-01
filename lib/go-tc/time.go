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
	"database/sql/driver"
	"strings"
	"time"
)

// Time wraps standard time.Time to allow indication of invalid times
type Time struct {
	time.Time
	Valid bool
}

// TimeLayout is the format used in lastUpdated fields in Traffic Ops
const TimeLayout = "2006-01-02 15:04:05-07"

// Scan implements the database/sql Scanner interface.
func (t *Time) Scan(value interface{}) error {
	t.Time, t.Valid = value.(time.Time)
	return nil
}

// Value implements the database/sql/driver Valuer interface.
func (t Time) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time, nil
}

// MarshalJSON implements the json.Marshaller interface
func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.Time.Format(TimeLayout) + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaller interface
func (t *Time) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		t.Time = time.Time{}
		return
	}
	t.Time, err = time.Parse(TimeLayout, s)
	return
}

// TimeNoMod supported JSON marshalling, but suppresses JSON unmarshalling
type TimeNoMod Time

func NewTimeNoMod() *TimeNoMod {
	return &TimeNoMod{Time: time.Now()}
}

// Scan implements the database/sql Scanner interface.
func (t *TimeNoMod) Scan(value interface{}) error {
	t.Time, t.Valid = value.(time.Time)
	return nil
}

// Value implements the database/sql/driver Valuer interface.
func (t TimeNoMod) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time, nil
}

// MarshalJSON implements the json.Marshaller interface
func (t TimeNoMod) MarshalJSON() ([]byte, error) {
	return Time(t).MarshalJSON()
}

// UnmarshalJSON for TimeNoMod suppresses unmarshalling
func (t *TimeNoMod) UnmarshalJSON([]byte) (err error) {
	return nil
}
