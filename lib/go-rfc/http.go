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

const ApplicationJSON = "application/json"
const Gzip = "gzip"
const ContentType = "Content-Type"
const ContentEncoding = "Content-Encoding"
const ContentTypeTextPlain = "text/plain"
const AcceptEncoding = "Accept-Encoding"

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
