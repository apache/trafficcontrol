package web

import (
	"net/http"
	"testing"
)

func makeHeader(cacheControlVal string) http.Header {
	return map[string][]string{"Cache-Control": []string{cacheControlVal}}
}

func TestParseCacheControl(t *testing.T) {
	testCacheControls := []string{
		"no-store, no-cache, must-revalidate, post-check=0, pre-check=0",
		"no-store, no-cache",
		"no-cache",
		"",
		`foo="bar"`,
		`foo="ba\"r"`,
		`foo="ba\"r", baz=blee, aaaa="bb\"\"\"", cc="dd", ee="ff\"f", gg=hh", i="", j="k", l="m\\\\o\"`,
		`foo="ba\"r", baz`,
		`foo=`,
	}

	for _, ccStr := range testCacheControls {
		ParseCacheControl(makeHeader(ccStr))
		// TODO actually test
		// fmt.Printf("parsed: %+v\n", cc)
	}
}
