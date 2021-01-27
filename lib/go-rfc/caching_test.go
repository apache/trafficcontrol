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
	"fmt"
	"net/http"
	"testing"
	"time"
)

func ExampleParseCacheControl() {
	hdrs := http.Header{}

	hdrs.Set(CacheControl, "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	cchm := ParseCacheControl(hdrs)
	fmt.Println(cchm.Has("no-store"))
	fmt.Println(cchm.Has("no-cache"))
	fmt.Println(cchm.Has("must-revalidate"))
	fmt.Println(cchm["post-check"])
	fmt.Println(cchm["pre-check"])

	// Output: true
	// true
	// true
	// 0
	// 0
}

func ExampleCacheControlMap_Has() {
	hdrs := http.Header{}
	hdrs.Set(CacheControl, "no-cache")

	ccm := ParseCacheControl(hdrs)
	if ccm.Has("no-cache") {
		fmt.Println("Has 'no-cache'")
	}
	if ccm.Has("no-store") {
		fmt.Println("Has 'no-store'")
	}
	// Output: Has 'no-cache'
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
	hdrs := http.Header{}
	ccStr := `foo="ba\"r", baz=blee, aaaa="bb\"\"\"", cc="dd", ee="ff\"f", gg=hh", i="", j="k", l="m\\\\o\"`
	hdrs.Set(CacheControl, ccStr)
	for i := 0; i < b.N; i++ {
		ParseCacheControl(hdrs)
	}
}

func TestCanCache(t *testing.T) {
	// tests RFC7234§5.2.1.5 compliance
	t.Run("client no-store with strict RFC", func(t *testing.T) {
		reqHdr := http.Header{
			"Cache-Control": {"no-store"},
		}
		respCode := http.StatusOK
		respHdr := http.Header{}
		strictRFC := true

		if CanCache(http.MethodGet, reqHdr, respCode, respHdr, strictRFC) {
			t.Errorf("CanCache returned true for no-store request and strict RFC")
		}
	})

	// tests RFC7234§5.2.1.5 violation to protect origins
	t.Run("client no-store without strict RFC", func(t *testing.T) {
		reqHdr := http.Header{
			"Cache-Control": {"no-store"},
		}
		respCode := 200
		respHdr := http.Header{}
		strictRFC := false

		if !CanCache(http.MethodGet, reqHdr, respCode, respHdr, strictRFC) {
			t.Errorf("CanCache returned false for no-store request and strict RFC disabled")
		}
	})

	// tests RFC7234§5.2.1.5 violation to protect origins
	t.Run("client no-cache no-store without strict RFC", func(t *testing.T) {
		reqHdr := http.Header{
			"Cache-Control": {"no-cache no-store"},
		}
		respCode := 200
		respHdr := http.Header{}
		strictRFC := false

		if !CanCache(http.MethodGet, reqHdr, respCode, respHdr, strictRFC) {
			t.Errorf("CanCache returned false for no-store request and strict RFC disabled")
		}
	})

	t.Run("parent no-cache with strict RFC", func(t *testing.T) {
		reqHdr := http.Header{}
		respCode := 200
		respHdr := http.Header{
			"Cache-Control": {"no-cache"},
		}
		strictRFC := false

		if CanCache(http.MethodGet, reqHdr, respCode, respHdr, strictRFC) {
			t.Errorf("CanCache returned true for no-cache request and strict RFC disabled")
		}
	})

	t.Run("parent no-store with strict RFC", func(t *testing.T) {
		reqHdr := http.Header{}
		respCode := 200
		respHdr := http.Header{
			"Cache-Control": {"no-store"},
		}
		strictRFC := false

		if CanCache(http.MethodGet, reqHdr, respCode, respHdr, strictRFC) {
			t.Errorf("CanCache returned true for no-cache request and strict RFC disabled")
		}
	})

	t.Run("parent no-cache without strict RFC", func(t *testing.T) {
		reqHdr := http.Header{}
		respCode := 200
		respHdr := http.Header{
			"Cache-Control": {"no-cache"},
		}
		strictRFC := false

		if CanCache(http.MethodGet, reqHdr, respCode, respHdr, strictRFC) {
			t.Errorf("CanCache returned true for no-cache request and strict RFC enabled")
		}
	})

	t.Run("parent no-store without strict RFC", func(t *testing.T) {
		reqHdr := http.Header{}
		respCode := 200
		respHdr := http.Header{
			"Cache-Control": {"no-store"},
		}
		strictRFC := true

		if CanCache(http.MethodGet, reqHdr, respCode, respHdr, strictRFC) {
			t.Errorf("CanCache returned true for no-cache request and strict RFC enabled")
		}
	})

	// tests RFC7234§5.2.1.3 and RFC7231§6.1 compliance
	t.Run("cache-able response codes", func(t *testing.T) {
		defaultCacheableCodes := []int{200, 203, 204, 206, 300, 301, 404, 405, 410, 414, 501} // RFC7231§6.1
		for _, code := range defaultCacheableCodes {
			reqHdr := http.Header{}
			respCode := code
			respHdr := http.Header{}
			strictRFC := true

			if !CanCache(http.MethodGet, reqHdr, respCode, respHdr, strictRFC) {
				t.Errorf("CanCache returned false for request with no cache control and default-cacheable response code %v", code)
			}
		}
	})

	// tests RFC7234§5.2.1.3 compliance
	t.Run("non-cache-able response codes without Cache-Control", func(t *testing.T) {
		nonDefaultCacheableCodes := []int{
			201, 202, 205, 207, 208, 226,
			302, 303, 304, 305, 306, 307, 308,
			400, 401, 402, 403, 406, 407, 408, 409, 411, 412, 413, 4015, 416, 417, 418, 421, 422, 423, 424, 428, 429, 431, 451,
			500, 502, 503, 504, 505, 506, 507, 508, 510, 511,
		}
		for _, code := range nonDefaultCacheableCodes {
			reqHdr := http.Header{}
			respCode := code
			respHdr := http.Header{}
			strictRFC := true

			if CanCache(http.MethodGet, reqHdr, respCode, respHdr, strictRFC) {
				t.Errorf("CanCache returned true for request with no cache control and non-default-cacheable response code %v", code)
			}
		}
	})

	// tests RFC7234§3 compliance
	t.Run("non-cache-able response codes with Cache-Control", func(t *testing.T) {
		nonDefaultCacheableCodes := []int{
			201, 202, 205, 207, 208, 226,
			302, 303, 304, 305, 306, 307, 308,
			400, 401, 402, 403, 406, 407, 408, 409, 411, 412, 413, 4015, 416, 417, 418, 421, 422, 423, 424, 428, 429, 431, 451,
			500, 502, 503, 504, 505, 506, 507, 508, 510, 511,
		}
		cacheableRespHdrs := []map[string][]string{
			{"Expires": {time.Now().Format(time.RFC1123)}},
			{"Cache-Control": {"max-age=42"}},
			{"Cache-Control": {"s-maxage=42"}},
		}

		for _, code := range nonDefaultCacheableCodes {
			for _, hdr := range cacheableRespHdrs {
				reqHdr := http.Header{}
				respCode := code
				respHdr := hdr
				strictRFC := true

				if !CanCache(http.MethodGet, reqHdr, respCode, respHdr, strictRFC) {
					t.Errorf("CanCache returned false for request with non-default-cacheable response code %v and cacheable header %v", respCode, respHdr)
				}
			}
		}
	})
}

// This benchmarks one of the longer checking processes: when the response code
// is not cache-able by default, but a Cache-Control header is present that
// makes it cache-able.
func BenchmarkCanCache(b *testing.B) {
	hdrs := http.Header{}
	hdrs.Set(CacheControl, "s-maxage=42")

	for i := 0; i < b.N; i++ {
		CanCache(http.MethodGet, http.Header{}, http.StatusCreated, hdrs, true)
	}
}

func TestCanReuseStored(t *testing.T) {

	// tests RFC7234§5.2.1.4 violation to protect origins
	t.Run("test client no-cache is ignored without strict RFC", func(t *testing.T) {
		reqHdr := http.Header{
			"Cache-Control": {"no-cache"},
		}
		respHdr := http.Header{}
		reqCC := CacheControlMap{
			"no-cache": "",
		}
		respCC := CacheControlMap{}
		respReqHdrs := http.Header{}
		respReqTime := time.Now()
		respRespTime := time.Now()

		strictRFC := false

		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != ReuseCan {
			t.Errorf("CanReuseStored for no-cache request and strict RFC disabled: expected ReuseCan, actual %v", reuse)
		}
	})

	// tests RFC7234§5.2.1.4 violation to protect origins
	t.Run("test client no-store is ignored without strict RFC", func(t *testing.T) {
		reqHdr := http.Header{
			"Cache-Control": {"no-store no-cache"},
		}
		respHdr := http.Header{}
		reqCC := CacheControlMap{
			"no-store": "",
			"no-cache": "",
		}
		respCC := CacheControlMap{}
		respReqHdrs := http.Header{}
		respReqTime := time.Now()
		respRespTime := time.Now()

		strictRFC := false

		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != ReuseCan {
			t.Errorf("CanReuseStored for no-cache request and strict RFC disabled: expected ReuseCan, actual %v", reuse)
		}
	})

	// tests RFC7234§5.2.1.4 violation to protect origins
	t.Run("test client no-cache and no-store is ignored without strict RFC", func(t *testing.T) {
		reqHdr := http.Header{
			"Cache-Control": {"no-store no-cache"},
		}
		respHdr := http.Header{}
		reqCC := CacheControlMap{
			"no-store": "",
			"no-cache": "",
		}
		respCC := CacheControlMap{}
		respReqHdrs := http.Header{}
		respReqTime := time.Now()
		respRespTime := time.Now()

		strictRFC := false

		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != ReuseCan {
			t.Errorf("CanReuseStored for no-cache request and strict RFC disabled: expected ReuseCan, actual %v", reuse)
		}
	})

	t.Run("test parent Expires in future is reused", func(t *testing.T) {
		now := time.Now()
		tenMinsBeforeExpires := now.Add(time.Minute * -10)
		expires := now.Format(time.RFC1123)
		reqHdr := http.Header{}
		respHdr := http.Header{
			"Expires": {expires},
		}
		reqCC := CacheControlMap{}
		respCC := CacheControlMap{}
		respReqHdrs := http.Header{}
		respReqTime := tenMinsBeforeExpires
		respRespTime := tenMinsBeforeExpires
		strictRFC := false
		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != ReuseCan {
			t.Errorf("CanReuseStored request after expires: expected ReuseCan, actual %v", reuse)
		}
	})

	t.Run("test parent Expires in past has revaldiate and can stale", func(t *testing.T) {
		now := time.Now()
		tenMinutesAgo := now.Add(time.Minute * -10)
		expires := tenMinutesAgo.Format(time.RFC1123)
		reqHdr := http.Header{}
		respHdr := http.Header{
			"Expires": {expires},
			"Date":    {expires},
		}
		reqCC := CacheControlMap{}
		respCC := CacheControlMap{}
		respReqHdrs := http.Header{}
		respReqTime := now
		respRespTime := now
		strictRFC := false
		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != ReuseMustRevalidateCanStale {
			t.Errorf("CanReuseStored request after expires: expected ReuseMustRevalidateCanStale, actual %v", reuse)
		}
	})

	t.Run("test parent Expires in past with must-revaldiate, cannot stale", func(t *testing.T) {
		now := time.Now()
		tenMinutesAgo := now.Add(time.Minute * -10)
		expires := tenMinutesAgo.Format(time.RFC1123)
		reqHdr := http.Header{}
		respHdr := http.Header{
			"Expires":       {expires},
			"Date":          {expires},
			"Cache-Control": {"must-revalidate"},
		}
		reqCC := CacheControlMap{}
		respCC := CacheControlMap{
			"must-revalidate": ""}
		respReqHdrs := http.Header{}
		respReqTime := now
		respRespTime := now
		strictRFC := false
		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != ReuseMustRevalidate {
			t.Errorf("CanReuseStored request after expires and response must-revalidate: expected ReuseMustRevalidate, actual %v", reuse)
		}
	})

	t.Run("test parent Expires in past proxy-revaldiate, and no-stale", func(t *testing.T) {
		now := time.Now()
		tenMinutesAgo := now.Add(time.Minute * -10)
		expires := tenMinutesAgo.Format(time.RFC1123)
		reqHdr := http.Header{}
		respHdr := http.Header{
			"Expires":       {expires},
			"Date":          {expires},
			"Cache-Control": {"must-revalidate"},
		}
		reqCC := CacheControlMap{}
		respCC := CacheControlMap{
			"must-revalidate": ""}
		respReqHdrs := http.Header{}
		respReqTime := now
		respRespTime := now
		strictRFC := false
		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != ReuseMustRevalidate {
			t.Errorf("CanReuseStored request after expires and response must-revalidate: expected ReuseMustRevalidate, actual %v", reuse)
		}
	})

	t.Run("test parent Expires in past with proxy-revaldiate returns MustRevalidateNoStale", func(t *testing.T) {
		now := time.Now()
		tenMinutesAgo := now.Add(time.Minute * -10)
		expires := tenMinutesAgo.Format(time.RFC1123)
		reqHdr := http.Header{}
		respHdr := http.Header{
			"Expires":       {expires},
			"Date":          {expires},
			"Cache-Control": {"proxy-revalidate"},
		}
		reqCC := CacheControlMap{}
		respCC := CacheControlMap{
			"proxy-revalidate": ""}
		respReqHdrs := http.Header{}
		respReqTime := now
		respRespTime := now
		strictRFC := false
		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != ReuseMustRevalidate {
			t.Errorf("CanReuseStored request after expires and response must-revalidate: expected ReuseMustRevalidate, actual %v", reuse)
		}
	})

	t.Run("test parent max-age in future returns CanReuse", func(t *testing.T) {
		now := time.Now()
		tenMinutesAgo := now.Add(time.Minute * -10)
		reqHdr := http.Header{}
		respHdr := http.Header{
			"Date":          {tenMinutesAgo.Format(time.RFC1123)},
			"Cache-Control": {"max-age=1200"}, // 20 minutes
		}
		reqCC := CacheControlMap{}
		respCC := CacheControlMap{
			"max-age": "1200"}
		respReqHdrs := http.Header{}
		respReqTime := now
		respRespTime := now
		strictRFC := false
		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != ReuseCan {
			t.Errorf("CanReuseStored request after expires and response must-revalidate: expected ReuseCan, actual %v", reuse)
		}
	})

	t.Run("test parent max-age in past returns MustRevalidate", func(t *testing.T) {
		now := time.Now()
		tenMinutesAgo := now.Add(time.Minute * -10)
		reqHdr := http.Header{}
		respHdr := http.Header{
			"Date":          {tenMinutesAgo.Format(time.RFC1123)},
			"Cache-Control": {"max-age=300"}, // 5 minutes
		}
		reqCC := CacheControlMap{}
		respCC := CacheControlMap{
			"max-age": "300"}
		respReqHdrs := http.Header{}
		respReqTime := tenMinutesAgo
		respRespTime := tenMinutesAgo
		strictRFC := false
		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != ReuseMustRevalidateCanStale {
			t.Errorf("CanReuseStored request after response max-age: expected ReuseMustRevalidateCanStale, actual %v", reuse)
		}
	})

	t.Run("test parent s-maxage in future returns CanReuse", func(t *testing.T) {
		now := time.Now()
		tenMinutesAgo := now.Add(time.Minute * -10)
		reqHdr := http.Header{}
		respHdr := http.Header{
			"Date":          {tenMinutesAgo.Format(time.RFC1123)},
			"Cache-Control": {"s-maxage=1200"}, // 20 minutes
		}
		reqCC := CacheControlMap{}
		respCC := CacheControlMap{
			"s-maxage": "1200"}
		respReqHdrs := http.Header{}
		respReqTime := now
		respRespTime := now
		strictRFC := false
		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != ReuseCan {
			t.Errorf("CanReuseStored request before s-maxage: expected ReuseCan, actual %v", reuse)
		}
	})

	t.Run("test parent s-age in past returns MustRevalidate", func(t *testing.T) {
		now := time.Now()
		tenMinutesAgo := now.Add(time.Minute * -10)
		reqHdr := http.Header{}
		respHdr := http.Header{
			"Date":          {tenMinutesAgo.Format(time.RFC1123)},
			"Cache-Control": {"s-maxage=300"}, // 5 minutes
		}
		reqCC := CacheControlMap{}
		respCC := CacheControlMap{
			"s-maxage": "300"}
		respReqHdrs := http.Header{}
		respReqTime := tenMinutesAgo
		respRespTime := tenMinutesAgo
		strictRFC := false
		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != ReuseMustRevalidateCanStale {
			t.Errorf("CanReuseStored request after response s-maxage: expected ReuseMustRevalidateCanStale, actual %v", reuse)
		}
	})

	t.Run("test parent future s-maxage overrides past max-age and returns CanReuse", func(t *testing.T) {
		now := time.Now()
		tenMinutesAgo := now.Add(time.Minute * -10)
		reqHdr := http.Header{}
		respHdr := http.Header{
			"Date":          {tenMinutesAgo.Format(time.RFC1123)},
			"Cache-Control": {"max-age=300,s-maxage=1200"},
		}
		reqCC := CacheControlMap{}
		respCC := CacheControlMap{
			"s-maxage": "1200",
			"max-age":  "300",
		}
		respReqHdrs := http.Header{}
		respReqTime := tenMinutesAgo
		respRespTime := tenMinutesAgo
		strictRFC := false
		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != ReuseCan {
			t.Errorf("CanReuseStored request before s-maxage but after max-age: expected ReuseCan, actual %v", reuse)
		}
	})

	t.Run("test parent past s-maxage overrides future max-age and returns MustReval", func(t *testing.T) {
		now := time.Now()
		tenMinutesAgo := now.Add(time.Minute * -10)
		reqHdr := http.Header{}
		respHdr := http.Header{
			"Date":          {tenMinutesAgo.Format(time.RFC1123)},
			"Cache-Control": {"max-age=1200,s-maxage=300"},
		}
		reqCC := CacheControlMap{}
		respCC := CacheControlMap{
			"s-maxage": "300",
			"max-age":  "1200",
		}
		respReqHdrs := http.Header{}
		respReqTime := tenMinutesAgo
		respRespTime := tenMinutesAgo
		strictRFC := false
		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != ReuseMustRevalidateCanStale {
			t.Errorf("CanReuseStored request before s-maxage but after max-age: expected ReuseMustRevalidateCanStale, actual %v", reuse)
		}
	})

	t.Run("test parent future max-age overrides past Expires and returns CanReuse", func(t *testing.T) {
		now := time.Now()
		twentyMinutesAgo := now.Add(time.Minute * -10)
		tenMinutesAgo := now.Add(time.Minute * -10)
		expires := twentyMinutesAgo.Format(time.RFC1123)
		reqHdr := http.Header{}
		respHdr := http.Header{
			"Date":          {tenMinutesAgo.Format(time.RFC1123)},
			"Expires":       {expires},
			"Cache-Control": {"max-age=1200"},
		}
		reqCC := CacheControlMap{}
		respCC := CacheControlMap{
			"max-age": "1200",
		}
		respReqHdrs := http.Header{}
		respReqTime := tenMinutesAgo
		respRespTime := tenMinutesAgo
		strictRFC := false
		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != ReuseCan {
			t.Errorf("CanReuseStored request before max-age but after Expires: expected ReuseCan, actual %v", reuse)
		}
	})

	t.Run("test parent past max-age overrides future Expires and returns MustRevalidate", func(t *testing.T) {
		now := time.Now()
		tenMinutesAgo := now.Add(time.Minute * -10)
		fiveMinutesHence := now.Add(time.Minute * 5)
		expires := fiveMinutesHence.Format(time.RFC1123)
		reqHdr := http.Header{}
		respHdr := http.Header{
			"Date":          {tenMinutesAgo.Format(time.RFC1123)},
			"Expires":       {expires},
			"Cache-Control": {"max-age=300"},
		}
		reqCC := CacheControlMap{}
		respCC := CacheControlMap{
			"max-age": "300",
		}
		respReqHdrs := http.Header{}
		respReqTime := tenMinutesAgo
		respRespTime := tenMinutesAgo
		strictRFC := false
		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != ReuseMustRevalidateCanStale {
			t.Errorf("CanReuseStored request after max-age but before Expires: expected ReuseMustRevalidateCanStale, actual %v", reuse)
		}
	})

	t.Run("test parent future s-maxage overrides past Expires and returns CanReuse", func(t *testing.T) {
		now := time.Now()
		twentyMinutesAgo := now.Add(time.Minute * -10)
		tenMinutesAgo := now.Add(time.Minute * -10)
		expires := twentyMinutesAgo.Format(time.RFC1123)
		reqHdr := http.Header{}
		respHdr := http.Header{
			"Date":          {tenMinutesAgo.Format(time.RFC1123)},
			"Expires":       {expires},
			"Cache-Control": {"s-maxage=1200"},
		}
		reqCC := CacheControlMap{}
		respCC := CacheControlMap{
			"s-maxage": "1200",
		}
		respReqHdrs := http.Header{}
		respReqTime := tenMinutesAgo
		respRespTime := tenMinutesAgo
		strictRFC := false
		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != ReuseCan {
			t.Errorf("CanReuseStored request before s-maxage but after Expires: expected ReuseCan, actual %v", reuse)
		}
	})

	t.Run("test parent past s-maxage overrides future Expires and returns MustRevalidate", func(t *testing.T) {
		now := time.Now()
		tenMinutesAgo := now.Add(time.Minute * -10)
		fiveMinutesHence := now.Add(time.Minute * 5)
		expires := fiveMinutesHence.Format(time.RFC1123)
		reqHdr := http.Header{}
		respHdr := http.Header{
			"Date":          {tenMinutesAgo.Format(time.RFC1123)},
			"Expires":       {expires},
			"Cache-Control": {"s-maxage=300"},
		}
		reqCC := CacheControlMap{}
		respCC := CacheControlMap{
			"s-maxage": "300",
		}
		respReqHdrs := http.Header{}
		respReqTime := tenMinutesAgo
		respRespTime := tenMinutesAgo
		strictRFC := false
		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != ReuseMustRevalidateCanStale {
			t.Errorf("CanReuseStored request after max-age but before Expires: expected ReuseMustRevalidateCanStale, actual %v", reuse)
		}
	})

	t.Run("test parent future s-maxage overrides past Expires and past max-age and returns CanReuse", func(t *testing.T) {
		now := time.Now()
		twentyMinutesAgo := now.Add(time.Minute * -10)
		tenMinutesAgo := now.Add(time.Minute * -10)
		expires := twentyMinutesAgo.Format(time.RFC1123)
		reqHdr := http.Header{}
		respHdr := http.Header{
			"Date":          {tenMinutesAgo.Format(time.RFC1123)},
			"Expires":       {expires},
			"Cache-Control": {"s-maxage=1200,max-age=300"},
		}
		reqCC := CacheControlMap{}
		respCC := CacheControlMap{
			"s-maxage": "1200",
			"max-age":  "300",
		}
		respReqHdrs := http.Header{}
		respReqTime := tenMinutesAgo
		respRespTime := tenMinutesAgo
		strictRFC := false
		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != ReuseCan {
			t.Errorf("CanReuseStored request before s-maxage but after Expires: expected ReuseCan, actual %v", reuse)
		}
	})

	t.Run("test parent past s-maxage overrides future Expires and future max-age and returns MustRevalidate", func(t *testing.T) {
		now := time.Now()
		tenMinutesAgo := now.Add(time.Minute * -10)
		fiveMinutesHence := now.Add(time.Minute * 5)
		expires := fiveMinutesHence.Format(time.RFC1123)
		reqHdr := http.Header{}
		respHdr := http.Header{
			"Date":          {tenMinutesAgo.Format(time.RFC1123)},
			"Expires":       {expires},
			"Cache-Control": {"s-maxage=300,max-age=1200"},
		}
		reqCC := CacheControlMap{}
		respCC := CacheControlMap{
			"s-maxage": "300",
			"max-age":  "1200",
		}
		respReqHdrs := http.Header{}
		respReqTime := tenMinutesAgo
		respRespTime := tenMinutesAgo
		strictRFC := false
		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != ReuseMustRevalidateCanStale {
			t.Errorf("CanReuseStored request after s-maxage but before Expires: expected ReuseMustRevalidateCanStale, actual %v", reuse)
		}
	})

	// tests RFC7234§5.2.1.3 compliance
	t.Run("test client min-fresh is obeyed with strict RFC", func(t *testing.T) {
		now := time.Now()
		tenMinutesAgo := now.Add(time.Minute * -10)
		reqHdr := http.Header{
			"Cache-Control": {"min-fresh=900"},
		}
		respHdr := http.Header{
			"Date":          {tenMinutesAgo.Format(time.RFC1123)},
			"Cache-Control": {"max-age=1200"},
		}
		reqCC := CacheControlMap{
			"min-fresh": "900",
		}
		respCC := CacheControlMap{
			"max-age": "1200",
		}
		respReqHdrs := http.Header{}
		respReqTime := tenMinutesAgo
		respRespTime := tenMinutesAgo
		strictRFC := true
		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != ReuseMustRevalidate {
			t.Errorf("CanReuseStored request with strictRFC min-fresh 300 with 600 remaining: expected ReuseMustRevalidate, actual %v", reuse)
		}
	})

	// tests RFC7234§5.2.1.3 violation to protect origins
	t.Run("test client min-fresh is ignored without strict RFC", func(t *testing.T) {
		now := time.Now()
		tenMinutesAgo := now.Add(time.Minute * -10)
		reqHdr := http.Header{
			"Cache-Control": {"min-fresh=900"},
		}
		respHdr := http.Header{
			"Date":          {tenMinutesAgo.Format(time.RFC1123)},
			"Cache-Control": {"max-age=1200"},
		}
		reqCC := CacheControlMap{
			"min-fresh": "900",
		}
		respCC := CacheControlMap{
			"max-age": "1200",
		}
		respReqHdrs := http.Header{}
		respReqTime := tenMinutesAgo
		respRespTime := tenMinutesAgo
		strictRFC := false
		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != ReuseCan {
			t.Errorf("CanReuseStored request with strictRFC min-fresh 1200 with 600 remaining: expected ReuseCan, actual %v", reuse)
		}
	})

	// tests RFC7234§5.2.1.3 violation to protect origins
	t.Run("test client min-fresh is ignored without strict RFC", func(t *testing.T) {
		now := time.Now()
		tenMinutesAgo := now.Add(time.Minute * -10)
		reqHdr := http.Header{
			"Cache-Control": {"min-fresh=900"},
		}
		respHdr := http.Header{
			"Date":          {tenMinutesAgo.Format(time.RFC1123)},
			"Cache-Control": {"max-age=1200"},
		}
		reqCC := CacheControlMap{
			"min-fresh": "900",
		}
		respCC := CacheControlMap{
			"max-age": "1200",
		}
		respReqHdrs := http.Header{}
		respReqTime := tenMinutesAgo
		respRespTime := tenMinutesAgo
		strictRFC := false
		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != ReuseCan {
			t.Errorf("CanReuseStored request with strictRFC min-fresh 1200 with 600 remaining: expected ReuseCan, actual %v", reuse)
		}
	})
}

func BenchmarkCanReuseStored(b *testing.B) {
	tenMinutesAgo := time.Now().Add(time.Minute * -10)
	reqHdr := http.Header{
		"Cache-Control": {"min-fresh=900"},
	}
	respHdr := http.Header{
		"Date":          {tenMinutesAgo.Format(time.RFC1123)},
		"Cache-Control": {"max-age=1200"},
	}
	reqCC := CacheControlMap{
		"min-fresh": "900",
	}
	respCC := CacheControlMap{
		"max-age": "1200",
	}
	respReqHdrs := http.Header{}
	respReqTime := tenMinutesAgo
	respRespTime := tenMinutesAgo
	strictRFC := true

	for i := 0; i < b.N; i++ {
		CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC)
	}
}
