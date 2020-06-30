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

import "net/http"

// CacheControl is the parameters found in an HTTP Cache-Control header,
// each mapped to its specified value.
type CacheControl map[string]string

// CacheableResponseCodes provides fast lookup of whether a HTTP response
// code is cache-able by default.
var CacheableResponseCodes = map[int]struct{}{
	http.StatusOK:                   {},
	http.StatusNonAuthoritativeInfo: {},
	http.StatusNoContent:            {},
	http.StatusPartialContent:       {},
	http.StatusMultipleChoices:      {},
	http.StatusMovedPermanently:     {},
	http.StatusNotFound:             {},
	http.StatusMethodNotAllowed:     {},
	http.StatusGone:                 {},
	http.StatusRequestURITooLong:    {},
	http.StatusNotImplemented:       {},
}
