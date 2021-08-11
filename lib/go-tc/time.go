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
	"math"
	"strconv"
	"strings"
	"time"
)

// Time wraps standard time.Time to allow indication of invalid times.
type Time struct {
	time.Time
	Valid bool
}

// TimeLayout is the format used in lastUpdated fields in Traffic Ops.
//
// Deprecated: New code should use RFC3339 (with optional nanosecond
// precision).
const TimeLayout = "2006-01-02 15:04:05-07"

// Do not ever use this. It only exists for compatibility with Perl.
const legacyLayout = "2006-01-02 15:04:05"

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

// MarshalJSON implements the json.Marshaller interface.
//
// Time structures marshal in the format defined by TimeLayout.
func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.Time.Format(TimeLayout) + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaller interface
//
// Time structures accept both RFC3339-compliant date/time strings as well as
// the format defined by TimeLayout and Unix(-ish) timestamps. Timestamps are
// expected to be integer numbers that represent milliseconds since Jan 1,
// 1970 00:00:00.000 UTC.
func (t *Time) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		t.Time = time.Time{}
		return
	}

	// timestamp support
	var i int64
	i, err = strconv.ParseInt(s, 10, 64)
	if err == nil {
		seconds := float64(i) / 1000.0
		t.Time = time.Unix(int64(seconds), int64(math.Copysign(float64(int64(math.Abs(seconds)))-math.Abs(seconds), seconds)*1000000))
		return
	}

	t.Time, err = time.Parse(time.RFC3339, s)
	if err == nil {
		return
	}
	t.Time, err = time.Parse(TimeLayout, s)
	if err == nil {
		return
	}

	// legacy
	t.Time, err = time.Parse(legacyLayout, s)
	return
}

// TimeNoMod supports JSON marshalling, but suppresses JSON unmarshalling.
//
// Deprecated: As TimeNoMod only supports the deprecated TimeLayout format for
// date/times, new code should just use time.Time instead (which defaults to
// the appropriate RFC3339), or Time to support Unix timestamps and backwards
// compatibility with this format.
type TimeNoMod Time

// NewTimeNoMod returns the address of a TimeNoMod with a Time value of the
// current time.
//
// Deprecated: TimeNoMod is deprecated.
func NewTimeNoMod() *TimeNoMod {
	return &TimeNoMod{Time: time.Now()}
}

// TimeNoModFromTime returns a reference to a TimeNoMod with the given Time
// value.
//
// Deprecated: TimeNoMod is deprecated.
func TimeNoModFromTime(t time.Time) *TimeNoMod {
	return &TimeNoMod{Time: t}
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

// MarshalJSON implements the json.Marshaller interface.
func (t TimeNoMod) MarshalJSON() ([]byte, error) {
	return Time(t).MarshalJSON()
}

// UnmarshalJSON for TimeNoMod suppresses unmarshalling.
func (t *TimeNoMod) UnmarshalJSON([]byte) (err error) {
	return nil
}

// TimeStamp holds the current time with nanosecond precision. It is unused.
//
// Deprecated: This is unused, so it will be removed in the future.
type TimeStamp time.Time
