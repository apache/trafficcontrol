package rfc

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
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

const (
	IfModifiedSince   = "If-Modified-Since" // RFC7232§3.3
	LastModified      = "Last-Modified"     // RFC7232§2.2
	ETagHeader        = "ETag"
	IfMatch           = "If-Match"
	IfUnmodifiedSince = "If-Unmodified-Since"
	Date              = "Date"
	ETagVersion       = 1
)

// ETag takes the last time the object was modified, and returns an ETag string. Note the string is the complete header value, including quotes. ETags must be quoted strings.
func ETag(t time.Time) string {
	return `"v` + strconv.Itoa(ETagVersion) + `-` + strconv.FormatInt(t.UnixNano(), 36) + `"`
}

// ParseETag takes a complete ETag header string, including the quotes (if the client correctly set them), and returns the last modified time encoded in the ETag.
func ParseETag(e string) (time.Time, error) {
	if len(e) < 2 || e[0] != '"' || e[len(e)-1] != '"' {
		return time.Time{}, errors.New("unquoted string, value must be quoted")
	}
	e, err := strconv.Unquote(e) // strip quotes

	if err != nil {
		return time.Time{}, err
	}
	prefix := `v` + strconv.Itoa(ETagVersion) + `-`
	if len(e) < len(prefix) || !strings.HasPrefix(e, prefix) {
		return time.Time{}, errors.New("malformed, no version prefix")
	}

	timeStr := e[len(prefix):]

	i, err := strconv.ParseInt(timeStr, 36, 64)
	if err != nil {
		return time.Time{}, err
	}

	t := time.Unix(0, i)

	const year = time.Hour * 24 * 365

	// sanity check - if the time isn't +/- 20 years, error. This catches overflows and near-zero errors
	if t.After(time.Now().Add(20*year)) || t.Before(time.Now().Add(-20*year)) {
		return time.Time{}, errors.New("malformed, out of range")
	}

	return t, nil
}

// GetUnmodifiedTime gets the latest time out of the Etags (if present), or the If-Unmodified-Since (if present)
func GetUnmodifiedTime(h http.Header) (time.Time, bool) {
	if h == nil {
		return time.Time{}, false
	}
	if im := h.Get(IfMatch); im != "" {
		if et, ok := ParseETags(strings.Split(im, ",")); ok {
			return et, true
		}
	}
	if ius := h.Get(IfUnmodifiedSince); ius != "" {
		if tm, ok := ParseHTTPDate(ius); ok {
			return tm, true
		}
	}
	return time.Time{}, false
}

// ParseETags the latest time of any valid ETag, and whether a valid ETag was found.
func ParseETags(eTags []string) (time.Time, bool) {
	latestTime := time.Time{}
	for _, tag := range eTags {
		tag = strings.TrimSpace(tag)
		et, err := ParseETag(tag)
		// errors are recoverable, keep going through the list of etags
		if err != nil {
			continue
		}
		if et.After(latestTime) {
			latestTime = et
		}
	}
	return latestTime, latestTime != time.Time{}
}
