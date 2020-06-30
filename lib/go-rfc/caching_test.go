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

import "fmt"
import "net/http"
import "testing"

func makeHeader(cacheControlVal string) http.Header {
	return map[string][]string{"Cache-Control": []string{cacheControlVal}}
}

func ExampleParseCacheControl(t *testing.T) {
	hdrs := http.Header{}

	hdrs.Set(CacheControl, "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	fmt.Println(ParseCacheControl(hdrs))

	// Output: Cache-Control: no-store, no-cache, must-revalidate, post-check=0, pre-check=0
}

func TestParseCacheControl(t *testing.T) {
	hdrs := http.Header{}

	ccStr := "no-store, no-cache, must-revalidate, post-check=0, pre-check=0"
	hdrs.Set(CacheControl, ccStr)
	cc := ParseCacheControl(hdrs)
	if len(cc) != 5 {
		t.Errorf("Incorrect number of parameters parsed from '%s'; expected: 5, actual: %d", ccStr, len(cc))
	}
	if _, ok := cc["no-store"]; !ok {
		t.Errorf("Expected parsed map from '%s' to have 'no-store' parameter, but it didn't", ccStr)
	}
	if _, ok := cc["no-cache"]; !ok {
		t.Errorf("Expected parsed map from '%s' to have 'no-cache' parameter, but it didn't", ccStr)
	}
	if _, ok := cc["must-revalidate"]; !ok {
		t.Errorf("Expected parsed map from '%s' to have 'must-revalidate' parameter, but it didn't", ccStr)
	}
	if pc, ok := cc["post-check"]; !ok {
		t.Errorf("Expected parsed map from '%s' to have 'post-check' parameter, but it didn't", ccStr)
	} else if pc != "0" {
		t.Errorf("Invalid value for 'post-check' parsed from '%s'; expected: '0', actual: '%s'", ccStr, pc)
	}
	if pc, ok := cc["pre-check"]; !ok {
		t.Errorf("Expected parsed map from '%s' to have 'pre-check' parameter, but it didn't", ccStr)
	} else if pc != "0" {
		t.Errorf("Invalid value for 'pre-check' parsed from '%s'; expected: '0', actual: '%s'", ccStr, pc)
	}

	ccStr = "no-store, no-cache"
	hdrs.Set(CacheControl, ccStr)
	cc = ParseCacheControl(hdrs)
	if len(cc) != 2 {
		t.Errorf("Incorrect number of parameters parsed from '%s'; expected: 2, actual: %d", ccStr, len(cc))
	}
	if _, ok := cc["no-store"]; !ok {
		t.Errorf("Expected parsed map from '%s' to have 'no-store' parameter, but it didn't", ccStr)
	}
	if _, ok := cc["no-cache"]; !ok {
		t.Errorf("Expected parsed map from '%s' to have 'no-cache' parameter, but it didn't", ccStr)
	}

	ccStr = "no-cache"
	hdrs.Set(CacheControl, ccStr)
	cc = ParseCacheControl(hdrs)
	if len(cc) != 1 {
		t.Errorf("Incorrect number of parameters parsed from '%s'; expected: 1, actual: %d", ccStr, len(cc))
	}
	if _, ok := cc["no-cache"]; !ok {
		t.Errorf("Expected parsed map from '%s' to have 'no-cache' parameter, but it didn't", ccStr)
	}

	ccStr = ""
	hdrs.Set(CacheControl, ccStr)
	cc = ParseCacheControl(hdrs)
	if len(cc) != 0 {
		t.Errorf("Incorrect number of parameters parsed from '%s'; expected: 0, actual: %d", ccStr, len(cc))
	}

	ccStr = `foo="bar"`
	hdrs.Set(CacheControl, ccStr)
	cc = ParseCacheControl(hdrs)
	if len(cc) != 1 {
		t.Errorf("Incorrect number of parameters parsed from '%s'; expected: 1, actual: %d", ccStr, len(cc))
	}
	if foo, ok := cc["foo"]; !ok {
		t.Errorf("Expected parsed map from '%s' to have 'foo' parameter, but it didn't", ccStr)
	} else if foo != "bar" {
		t.Errorf("Invalid value for 'foo' parsed from '%s'; expected: 'bar', actual: '%s'", ccStr, foo)
	}

	ccStr = `foo="ba\"r"`
	hdrs.Set(CacheControl, ccStr)
	cc = ParseCacheControl(hdrs)
	if len(cc) != 1 {
		t.Errorf("Incorrect number of parameters parsed from '%s'; expected: 1, actual: %d", ccStr, len(cc))
	}
	if foo, ok := cc["foo"]; !ok {
		t.Errorf("Expected parsed map from '%s' to have 'foo' parameter, but it didn't", ccStr)
	} else if foo != `ba"r` {
		t.Errorf(`Invalid value for 'foo' parsed from '%s'; expected: 'ba"r', actual: '%s'`, ccStr, foo)
	}

	ccStr = `foo="ba\"r", baz=blee, aaaa="bb\"\"\"", cc="dd", ee="ff\"f", gg=hh", i="", j="k", l="m\\\\o\"`
	hdrs.Set(CacheControl, ccStr)
	cc = ParseCacheControl(hdrs)
	if len(cc) != 9 {
		t.Errorf("Incorrect number of parameters parsed from '%s'; expected: 9, actual: %d", ccStr, len(cc))
	}
	if foo, ok := cc["foo"]; !ok {
		t.Errorf("Expected parsed map from '%s' to have 'foo' parameter, but it didn't", ccStr)
	} else if foo != `ba"r` {
		t.Errorf(`Invalid value for 'foo' parsed from '%s'; expected: 'ba"r', actual: '%s'`, ccStr, foo)
	}
	if baz, ok := cc["baz"]; !ok {
		t.Errorf("Expected parsed map from '%s' to have 'baz' parameter, but it didn't", ccStr)
	} else if baz != "blee" {
		t.Errorf("Invalid value for 'baz' parsed from '%s'; expected: 'blee', actual: '%s'", ccStr, baz)
	}
	if aaaa, ok := cc["aaaa"]; !ok {
		t.Errorf("Expected parsed map from '%s' to have 'aaaa' parameter, but it didn't", ccStr)
	} else if aaaa != `bb"""` {
		t.Errorf(`Invalid value for 'aaaa' parsed from '%s'; expected: 'b"""', actual: '%s'`, ccStr, aaaa)
	}
	if ccParam, ok := cc["cc"]; !ok {
		t.Errorf("Expected parsed map from '%s' to have 'cc' parameter, but it didn't", ccStr)
	} else if ccParam != "dd" {
		t.Errorf("Invalid value for 'cc' parsed from '%s'; expected: 'dd', actual: '%s'", ccStr, ccParam)
	}
	if ee, ok := cc["ee"]; !ok {
		t.Errorf("Expected parsed map from '%s' to have 'ee' parameter, but it didn't", ccStr)
	} else if ee != `ff"f` {
		t.Errorf(`Invalid value for 'ee' parsed from '%s'; expected: 'ff"f', actual: '%s'`, ccStr, ee)
	}
	if gg, ok := cc["gg"]; !ok {
		t.Errorf("Expected parsed map from '%s' to have 'gg' parameter, but it didn't", ccStr)
	} else if gg != `hh"` {
		t.Errorf(`Invalid value for 'gg' parsed from '%s'; expected: 'hh"', actual: '%s'`, ccStr, gg)
	}
	if i, ok := cc["i"]; !ok {
		t.Errorf("Expected parsed map from '%s' to have 'i' parameter, but it didn't", ccStr)
	} else if i != "" {
		t.Errorf("Invalid value for 'i' parsed from '%s'; expected: '', actual: '%s'", ccStr, i)
	}
	if j, ok := cc["j"]; !ok {
		t.Errorf("Expected parsed map from '%s' to have 'j' parameter, but it didn't", ccStr)
	} else if j != "k" {
		t.Errorf("Invalid value for 'j' parsed from '%s'; expected: 'k', actual: '%s'", ccStr, j)
	}
	if l, ok := cc["l"]; !ok {
		t.Errorf("Expected parsed map from '%s' to have 'l' parameter, but it didn't", ccStr)
	} else if l != `m\\o\` {
		t.Errorf(`Invalid value for 'l' parsed from '%s'; expected: 'm\\o\', actual: '%s'`, ccStr, l)
	}

	ccStr = `foo="ba\"r", baz`
	hdrs.Set(CacheControl, ccStr)
	cc = ParseCacheControl(hdrs)
	if len(cc) != 2 {
		t.Errorf("Incorrect number of parameters parsed from '%s'; expected: 2, actual: %d", ccStr, len(cc))
	}
	if foo, ok := cc["foo"]; !ok {
		t.Errorf("Expected parsed map from '%s' to have 'foo' parameter, but it didn't", ccStr)
	} else if foo != `ba"r` {
		t.Errorf(`Invalid value for 'foo' parsed from '%s'; expected: 'ba"r', actual: '%s'`, ccStr, foo)
	}
	if _, ok := cc["baz"]; !ok {
		t.Errorf("Expected parsed map from '%s' to have 'baz' parameter, but it didn't", ccStr)
	}

	ccStr = "foo="
	hdrs.Set(CacheControl, ccStr)
	cc = ParseCacheControl(hdrs)
	if len(cc) != 1 {
		t.Errorf("Incorrect number of parameters parsed from '%s'; expected: 1, actual: %d", ccStr, len(cc))
	}
	if foo, ok := cc["foo"]; !ok {
		t.Errorf("Expected parsed map from '%s' to have 'foo' parameter, but it didn't", ccStr)
	} else if foo != "" {
		t.Errorf("Invalid value for 'foo' parsed from '%s'; expected: '', actual: '%s'", ccStr, foo)
	}
}

func BenchmarkParseCacheControl(b *testing.B) {
	var hdrs http.Header
	ccStr := `foo="ba\"r", baz=blee, aaaa="bb\"\"\"", cc="dd", ee="ff\"f", gg=hh", i="", j="k", l="m\\\\o\"`
	hdrs.Set(CacheControl, ccStr)
	for i := 0; i < b.N; i++ {
		ParseCacheControl(hdrs)
	}
}
