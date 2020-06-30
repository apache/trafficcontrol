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
	http.MethodGet:  {},
	http.MethodHead: {},
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
				c[strings.TrimSpace(cacheControlStr)] = ""
				return c
			}

			key := strings.TrimSpace(cacheControlStr[:nextSpaceOrEqualPos])
			if cacheControlStr[nextSpaceOrEqualPos] == ',' {
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

// Checks if the cache control allows responses to be cached.
func cacheControlAllows(respCode int, respHeaders http.Header, respCacheControl CacheControlMap) bool {
	if _, ok := respHeaders["Expires"]; ok {
		return true
	}
	if _, ok := respCacheControl["max-age"]; ok {
		return true
	}
	if _, ok := respCacheControl["s-maxage"]; ok {
		return true
	}
	// This used to be a stub function that just returns false, the original rationale for
	// why it was always false is shown here in the comment from that original function:
	// This MUST return false unless a specific Cache Control cache-extension token exists for an extension which allows. Which is to say, returning true here without a cache-extension token is in strict violation of RFC7234.
	// In practice, all returning true does is override whether a response code is default-cacheable. If we wanted to do that, it would be better to make codeDefaultCacheable take a strictRFC parameter.
	// if extensionAllows() {
	// 	return true
	// }
	if _, ok := CacheableResponseCodes[respCode]; ok {
		return true
	}
	// log.Debugf("CacheControlAllows false: no expires, no max-age, no s-max-age, no extension allows, code not default cacheable\n")
	return false
}

// canStoreResponse checks the constraints in RFC7234.
func canStoreResponse(respCode int, respHeaders http.Header, reqCC, respCC CacheControlMap, strictRFC bool) bool {
	if _, ok := reqCC["no-store"]; strictRFC && ok {
		// log.Debugf("CanStoreResponse false: request has no-store\n")
		return false
	}
	if _, ok := respCC["no-store"]; ok {
		// log.Debugf("CanStoreResponse false: response has no-store\n") // RFC7234§5.2.2.3
		return false
	}
	if _, ok := respCC["no-cache"]; ok {
		// log.Debugf("CanStoreResponse false: response has no-cache\n") // RFC7234§5.2.2.2
		return false
	}
	if _, ok := respCC["private"]; ok {
		// log.Debugf("CanStoreResponse false: has private\n")
		return false
	}
	if _, ok := respCC["authorization"]; ok {
		// log.Debugf("CanStoreResponse false: has authorization\n")
		return false
	}
	return cacheControlAllows(respCode, respHeaders, respCC)
}

// canStoreAuthenticated checks the constraints in RFC7234§3.2
// TODO: ensure RFC7234§3.2 requirements that max-age=0, must-revlaidate, s-maxage=0 are revalidated
func canStoreAuthenticated(reqCacheControl, respCacheControl CacheControlMap) bool {
	if _, ok := reqCacheControl["authorization"]; !ok {
		return true
	}
	if _, ok := respCacheControl["must-revalidate"]; ok {
		return true
	}
	if _, ok := respCacheControl["public"]; ok {
		return true
	}
	if _, ok := respCacheControl["s-maxage"]; ok {
		return true
	}
	// log.Debugf("CanStoreAuthenticated false: has authorization, and no must-revalidate/public/s-maxage\n")
	return false
}

// CanCache returns whether an object can be cached per RFC 7234, based on the
// request headers, response headers, and response code.
//
// If strictRFC is false, this ignores request headers denying cacheability such
// as `no-cache`, in order to protect origins.
// TODO add options to ignore/violate request Cache-Control (to protect origins)
func CanCache(reqMethod string, reqHeaders http.Header, respCode int, respHeaders http.Header, strictRFC bool) bool {
	// log.Debugf("CanCache start\n")
	if _, ok := CacheableRequestMethods[reqMethod]; !ok {
		return false // for now, we only support GET and HEAD as cacheable methods.
	}

	reqCacheControl := ParseCacheControl(reqHeaders)
	respCacheControl := ParseCacheControl(respHeaders)
	// log.Debugf("CanCache reqCacheControl %+v respCacheControl %+v\n", reqCacheControl, respCacheControl)
	return canStoreResponse(respCode, respHeaders, reqCacheControl, respCacheControl, strictRFC) && canStoreAuthenticated(reqCacheControl, respCacheControl)
}
