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

type Time struct {
	time.Time
	Valid bool
}

// TimeLayout is the format used in lastUpdated fields in Traffic Ops
const TimeLayout = "2006-01-02 15:04:05-07"

// Scan implements the database/sql Scanner interface.
func (jt *Time) Scan(value interface{}) error {
	jt.Time, jt.Valid = value.(time.Time)
	return nil
}

// Value implements the database/sql/driver Valuer interface.
func (jt Time) Value() (driver.Value, error) {
	if !jt.Valid {
		return nil, nil
	}
	return jt.Time, nil
}

// MarshalJSON formats the Time field as Traffic Control expects it.
func (t *Time) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + t.Time.Format(TimeLayout) + `"`), nil
}

// UnmarshalJSON reads time from JSON into Time var
func (t *Time) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		t.Time = time.Time{}
		return
	}
	t.Time, err = time.Parse(TimeLayout, s)
	return
}
