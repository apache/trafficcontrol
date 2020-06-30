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
import "strings"

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

// CacheableRequestMethods is the list of all request methods which elicit
// cache-able responses.
var CacheableRequestMethods = map[string]struct{}{
	http.MethodGet: {},
}

// CacheControlMap is the parameters found in an HTTP Cache-Control header,
// each mapped to its specified value.
type CacheControlMap map[string]string

// String implements the Stringer interface by returning a textual representation
// of the CacheControlMap.
func (ccm CacheControlMap) String() string {
	s := "Cache-Control:"

	parts := make([]string, 0, len(ccm))
	for k, v := range ccm {
		if v != "" {
			parts = append(parts, k+"="+v)
		} else {
			parts = append(parts, k)
		}
	}

	if len(parts) > 0 {
		s += " " + strings.Join(parts, ", ")
	}

	return s
}

// Gets the position in the string at which a *quoted* Cache-Control parameter value
// ends - assuming it begins with the start of such a value.
//
// If the end can't be determined, returns -1.
func getQuotedValueEndPos(cacheControlStr string) int {
	if len(cacheControlStr) == 0 {
		return -1
	}
	if cacheControlStr[0] != '"' {
		return -1 // should never happen - log?
	}
	cacheControlStr = cacheControlStr[1:]

	skip := 0
	for {
		nextQuotePos := strings.Index(cacheControlStr[skip:], `"`) + skip
		if nextQuotePos == 0 || nextQuotePos == skip-1 { // -1 because we skip = nextQuotePos+1, to skip the actual quote
			return skip + 1 + 1 // +1 for the " we stripped at the beginning, +1 for quote itself
		}

		charBeforeQuote := cacheControlStr[nextQuotePos-1]
		if charBeforeQuote == '\\' {
			skip = nextQuotePos + 1
			continue
		}
		return nextQuotePos + 1 + 1 // +1 for the " we stripped at the beginning, +1 for the quote itself
	}
}

// Gets the position in the string at which a Cache-Control parameter value ends -
// assuming it begins with the start of such a value.
//
// If the end can't be determined, returns -1.
func getValueEndPos(cacheControlStr string) int {
	if len(cacheControlStr) == 0 {
		return -1
	}
	if cacheControlStr[0] != '"' {
		return strings.Index(cacheControlStr, `,`)
	}
	return getQuotedValueEndPos(cacheControlStr)
}

// Strips escape characters from the string.
//
// For example, `\w` becomes just `w` and `\\w` becomes `\w`.
func stripEscapes(s string) string {
	before := ""
	after := s
	for {
		i := strings.IndexAny(after, `\`)
		if i == -1 {
			return before + after
		}
		if len(after) <= i+1 {
			return before + after
		}
		if after[i+1] == '\\' {
			i++
		}
		if len(after) < i {
			return before + after
		}
		before += after[:i]
		after = after[i+1:]
	}
}

// ParseCacheControl parses the Cache-Control header from the headers object,
// and returns the parsed map of cache control directives.
//
// TODO verify Header/CacheControl are properly CanonicalCase/lowercase. Put cache-control text in constants?
func ParseCacheControl(h http.Header) CacheControlMap {
	c := CacheControlMap{}

	for _, cacheControlStr := range h[CacheControl] {
		for len(cacheControlStr) > 0 {
			nextSpaceOrEqualPos := strings.IndexAny(cacheControlStr, "=,")
			if nextSpaceOrEqualPos == -1 {
				c[cacheControlStr] = ""
				return c
			}

			key := cacheControlStr[:nextSpaceOrEqualPos] // TODO trim
			if cacheControlStr[nextSpaceOrEqualPos] == ' ' {
				cacheControlStr = cacheControlStr[nextSpaceOrEqualPos+1:]
				c[key] = ""
				continue
			}

			if len(cacheControlStr) < nextSpaceOrEqualPos+2 {
				c[key] = ""
				return c
			}
			cacheControlStr = cacheControlStr[nextSpaceOrEqualPos+1:]
			quoted := cacheControlStr[0] == '"'
			valueEndPos := getValueEndPos(cacheControlStr)
			if valueEndPos == -1 {
				c[key] = cacheControlStr
				return c
			}

			if len(cacheControlStr) < valueEndPos {
				value := cacheControlStr
				if quoted && len(value) > 1 {
					value = value[1 : len(value)-1]
					value = stripEscapes(value)
				}
				c[key] = value // TODO trim
				return c
			}
			value := cacheControlStr[:valueEndPos]
			if quoted && len(value) > 1 {
				value = value[1 : len(value)-1]
				value = stripEscapes(value)
			}
			c[key] = value // TODO trim

			if len(cacheControlStr) < valueEndPos+2 {
				return c
			}
			cacheControlStr = cacheControlStr[valueEndPos+2:]
		}
	}
	return c
}

// CanCache returns whether an object can be cached per RFC 7234, based on the
// request headers, response headers, and response code.
//
// If strictRFC is false, this ignores request headers denying cacheability such
// as `no-cache`, in order to protect origins.
// TODO add options to ignore/violate request cache-control (to protect origins)
// func CanCache(reqMethod string, reqHeaders http.Header, respCode int, respHeaders http.Header, strictRFC bool) bool {
// 	log.Debugf("CanCache start\n")
// 	if reqMethod != http.MethodGet {
// 		return false // for now, we only support GET as a cacheable method.
// 	}
// 	reqCacheControl := web.ParseCacheControl(reqHeaders)
// 	respCacheControl := web.ParseCacheControl(respHeaders)
// 	log.Debugf("CanCache reqCacheControl %+v respCacheControl %+v\n", reqCacheControl, respCacheControl)
// 	return canStoreResponse(respCode, respHeaders, reqCacheControl, respCacheControl, strictRFC) && canStoreAuthenticated(reqCacheControl, respCacheControl)
// }
