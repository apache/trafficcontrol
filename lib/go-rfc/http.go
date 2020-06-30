package rfc

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
	"net/http"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-util"
)

const (
	ApplicationJSON        = "application/json"         // RFC4627§6
	Gzip                   = "gzip"                     // RFC7230§4.2.3
	ContentType            = "Content-Type"             // RFC7231§3.1.1.5
	ContentEncoding        = "Content-Encoding"         // RFC7231§3.1.2.2
	ContentTypeTextPlain   = "text/plain"               // RFC2046§4.1
	AcceptEncoding         = "Accept-Encoding"          // RFC7231§5.3.4
	ContentDisposition     = "Content-Disposition"      // RFC6266
	ApplicationOctetStream = "application/octet-stream" // RFC2046§4.5.2
	Vary                   = "Vary"                     // RFC7231§7.1.4
	IfModifiedSince        = "If-Modified-Since"        // RFC7232§3.3
	LastModified           = "Last-Modified"            // RFC7232§2.2
)

// AcceptsGzip returns whether r accepts gzip encoding, per RFC7231§5.3.4.
func AcceptsGzip(r *http.Request) bool {
	encodingHeaders := r.Header[AcceptEncoding] // headers are case-insensitive, but Go promises to Canonical-Case requests
	for _, encodingHeader := range encodingHeaders {
		encodingHeader = util.StripAllWhitespace(encodingHeader)
		encodings := strings.Split(encodingHeader, ",")
		for _, encoding := range encodings {
			if strings.ToLower(encoding) == Gzip { // encoding is case-insensitive, per the RFC7231§3.1.2.1.
				return true
			}
		}
	}
	return false
}
