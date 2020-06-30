package web

/*
   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

import (
	"net/http"
	"strings"
)

type CacheControl map[string]string

// ParseCacheControl parses the Cache-Control header from the headers object, and returns the parsed map of cache control directives.
// TODO verify Header/CacheControl are properly CanonicalCase/LowerCase. Put cache-control text in constants?
func ParseCacheControl(h http.Header) CacheControl {
	c := CacheControl{}

	cacheControlStrs := h["Cache-Control"]
	for _, cacheControlStr := range cacheControlStrs {
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

func getValueEndPos(cacheControlStr string) int {
	if len(cacheControlStr) == 0 {
		return -1
	}
	if cacheControlStr[0] != '"' {
		nextSpace := strings.Index(cacheControlStr, `,`)
		return nextSpace
	}
	return getQuotedValueEndPos(cacheControlStr)
}

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
