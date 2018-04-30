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
