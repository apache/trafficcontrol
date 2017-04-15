package grove

import (
	"fmt"
	"net/http"
	"strings"
)

type CacheControl map[string]string

// TODO verify Header/CacheControl are properly CanonicalCase/LowerCase. Put cache-control text in constants?
func ParseCacheControl(h http.Header) CacheControl {
	c := CacheControl{}

	cacheControlStrs := h["Cache-Control"]
	for _, cacheControlStr := range cacheControlStrs {
		for len(cacheControlStr) > 0 {
			fmt.Printf("looping cacheControlStr ''%v''\n", cacheControlStr)
			nextSpaceOrEqualPos := strings.IndexAny(cacheControlStr, "=,")
			fmt.Printf("nextSpaceOrEqualPos %v\n", nextSpaceOrEqualPos)
			if nextSpaceOrEqualPos == -1 {
				fmt.Printf("setting0 c[''%v'']=\"\"\n", cacheControlStr)
				c[cacheControlStr] = ""
				return c
			}

			key := cacheControlStr[:nextSpaceOrEqualPos] // TODO trim
			fmt.Printf("key %v\n", cacheControlStr[:nextSpaceOrEqualPos])
			fmt.Printf("nextSpaceOrEqualPos %v\n", nextSpaceOrEqualPos)
			fmt.Printf("cacheControlStr[nextSpaceOrEqualPos] ''%v''\n", string(cacheControlStr[nextSpaceOrEqualPos]))
			if cacheControlStr[nextSpaceOrEqualPos] == ' ' {
				cacheControlStr = cacheControlStr[nextSpaceOrEqualPos+1:]
				fmt.Printf("cacheControlStr0 ''%v''\n", cacheControlStr)
				fmt.Printf("setting1 c[''%v'']=\"\"\n", key)
				c[key] = ""
				continue
			}

			if len(cacheControlStr) < nextSpaceOrEqualPos+2 {
				c[key] = ""
				return c
			}
			cacheControlStr = cacheControlStr[nextSpaceOrEqualPos+1:]
			fmt.Printf("cacheControlStr1 ''%v''\n", cacheControlStr)
			quoted := cacheControlStr[0] == '"'
			valueEndPos := getValueEndPos(cacheControlStr)
			fmt.Printf("valueEndPos %v\n", valueEndPos)
			if valueEndPos == -1 {
				fmt.Printf("setting2 c[''%v'']=''%v''\n", key, cacheControlStr)
				c[key] = cacheControlStr
				return c
			}
			fmt.Printf("cacheControlStr2 '%v' valueEndPos %v\n", cacheControlStr, valueEndPos)

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
			fmt.Printf("cacheControlStr2 value ''%v''\n", value)
			if quoted && len(value) > 1 {
				value = value[1 : len(value)-1]
				value = stripEscapes(value)
			}
			fmt.Printf("setting3 c[''%v'']=''%v''\n", key, value)
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
	fmt.Printf("getValueEndPos cacheControlStr ''%v''\n", cacheControlStr)
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
	fmt.Printf("getQuotedValueEndPos cacheControlStr ''%v''\n", cacheControlStr)

	if len(cacheControlStr) == 0 {
		fmt.Printf("getQuotedValueEndPos len(cacheControlStr)==0 !!\n")
		return -1
	}
	if cacheControlStr[0] != '"' {
		fmt.Printf("getQuotedValueEndPos cacheControlStr[0]==\"!!\n")
		return -1 // should never happen - log?
	}
	cacheControlStr = cacheControlStr[1:]

	skip := 0
	for {
		fmt.Printf("cacheControlStr ''%v''\n", cacheControlStr)
		fmt.Printf("skip %v\n", skip)
		fmt.Printf("cacheControlStr[skip:] ''%v''\n", cacheControlStr[skip:])
		nextQuotePos := strings.Index(cacheControlStr[skip:], `"`) + skip
		fmt.Printf("nextQuotePos %v\n", nextQuotePos)
		if nextQuotePos == 0 || nextQuotePos == skip-1 { // -1 because we skip = nextQuotePos+1, to skip the actual quote
			fmt.Printf("getQuotedValueEndPos returning skip %v\n", skip)
			return skip + 1 + 1 // +1 for the " we stripped at the beginning, +1 for quote itself
		}

		charBeforeQuote := cacheControlStr[nextQuotePos-1]
		fmt.Printf("getQuotedValueEndPos charBeforeQuote %v\n", string(charBeforeQuote))
		if charBeforeQuote == '\\' {
			skip = nextQuotePos + 1
			continue
		}
		fmt.Printf("getQuotedValueEndPos returning nextQuotePos %v + 1\n", nextQuotePos)
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
		fmt.Printf("stripEscapes before: ''%s'' after: ''%s''\n", before, after)
	}
}
