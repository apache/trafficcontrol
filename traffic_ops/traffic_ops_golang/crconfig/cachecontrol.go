package crconfig

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
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/grove/web" // TODO move to lib/go-rfc
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	// "github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
)

// TODO HTTP RFC symbols to lib/go-rfc along with the Grove rfc package

const HdrETag = "ETag"
const HdrIfNoneMatch = "If-None-Match"
const HdrIfModifiedSince = "If-Modified-Since"
const HdrLastModified = "Last-Modified"

const ETagVersion = 0

// AddModifiedHdrs adds the ETag and Last-Modified headers to the response of w.
// This must be called before any body is written, of course, or the headers won't be.
// The lastModified should be the last time this resource was modified.
func AddModifiedHdrs(w http.ResponseWriter, lastModified time.Time) {
	w.Header().Set(HdrLastModified, FormatHTTPDate(lastModified))
	w.Header().Set(HdrETag, ETag(lastModified))
}

// GetModifiedHdr gets the modified time from the ETag or Last-Modified.
// Returns the modified time, and whether a time was found.
func GetModifiedHdr(r *http.Request) (time.Time, bool) {
	// Note we intentionally ignore errors. Logging client errs would be spammy and a potential attack vector, and we want to just return the real object to the client, not an error, if their headers were malformed.
	if eTag := r.Header.Get(HdrIfNoneMatch); eTag != "" {
		if t, err := ParseETag(eTag); err == nil {
			return t, true
		}
	}
	if lastModified := r.Header.Get(HdrIfModifiedSince); lastModified != "" {
		if t, ok := web.ParseHTTPDate(lastModified); ok {
			return t, true
		}
	}
	return time.Time{}, false
}

// FormatHTTPDate formats t as an RFC7231§7.1.1 HTTP-date.
func FormatHTTPDate(t time.Time) string {
	return t.Format(time.RFC1123)
}

// ETag takes the last time the CRConfig was modified (the snapshot time), and returns an ETag string. Note the string is the complete header value, including quotes. ETags must be quoted strings.
func ETag(t time.Time) string {
	return `"v` + strconv.Itoa(ETagVersion) + `-` + strconv.FormatInt(t.UnixNano(), 36) + `"`
}

// ParseETag takes a complete ETag header string, including the quotes (if the client correctly set them), and returns the last modified time encoded in the ETag.
func ParseETag(e string) (time.Time, error) {
	if len(e) < 2 || e[0] != '"' || e[len(e)-1] != '"' {
		return time.Time{}, errors.New("unquoted string, value must be quoted")
	}
	e = e[1 : len(e)-1] // strip quotes

	fmt.Println("DEBUG stripped quotes: X" + e + "X")

	prefix := `v` + strconv.Itoa(ETagVersion) + `-`
	if len(e) < len(prefix) || !strings.HasPrefix(e, prefix) {
		return time.Time{}, errors.New("malformed, no version prefix")
	}

	timeStr := e[len(prefix):]

	fmt.Println("DEBUG stripped version: X" + timeStr + "X")

	i, err := strconv.ParseInt(timeStr, 36, 64)
	if err != nil {
		return time.Time{}, errors.New("malformed")
	}

	t := time.Unix(0, i)

	const year = time.Hour * 24 * 365

	// sanity check - if the time isn't +/- 20 years, error. This catches overflows and near-zero errors
	if t.After(time.Now().Add(20*year)) || t.Before(time.Now().Add(-20*year)) {
		return time.Time{}, errors.New("malformed, out of range")
	}

	return t, nil
}

// IsNotModified returns whether the CRConfig is unmodified, and if unmodified, the latest unmodified time..
// Whether the CRConfig is modified is determined if the request contains a modified time header, and no CDN or Delivery Service has been snapshotted since that time.
// If false is returned (the CRConfig is modified), the time returned is undefined, an MUST NOT be used. If false is returned, the time IS NOT the snapshotted time.
func IsNotModified(r *http.Request, cdn tc.CDNName, tx *sql.Tx) (bool, time.Time) {
	lastModified, ok := GetModifiedHdr(r)
	if !ok {
		return false, time.Time{}
	}

	snapshotTime, ok, err := GetCRConfigSnapshotTime(tx, cdn)
	if err != nil {
		log.Errorln("Getting CRConfig Snapshot Time: " + err.Error())
		return false, time.Time{}
	}
	if !ok {
		return false, time.Time{}
	}
	if snapshotTime.After(lastModified) {
		return false, time.Time{}
	}
	return true, snapshotTime
}
