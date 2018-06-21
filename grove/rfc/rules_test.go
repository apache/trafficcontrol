package rfc

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

// Package rfc contains functions implementing RFC 7234, 2616, and other RFCs.
// When changing functions, be sure they still conform to the corresponding RFC.
// When adding symbols, document the RFC and section they correspond to.

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/grove/remapdata"
	"github.com/apache/trafficcontrol/grove/web"

	"github.com/apache/trafficcontrol/lib/go-log"
)

func TestRules(t *testing.T) {
	// test client no-store is obeyed with strict RFC - tests RFC7234§5.2.1.5 compliance
	{
		reqHdr := http.Header{
			"Cache-Control": {"no-store"},
		}
		respCode := 200
		respHdr := http.Header{}
		strictRFC := true

		if CanCache(http.MethodGet, reqHdr, respCode, respHdr, strictRFC) {
			t.Errorf("CanCache returned true for no-store request and strict RFC")
		}
	}

	// test client no-store is ignored without strict RFC - tests RFC7234§5.2.1.5 violation to protect origins
	{
		reqHdr := http.Header{
			"Cache-Control": {"no-store"},
		}
		respCode := 200
		respHdr := http.Header{}
		strictRFC := false

		if !CanCache(http.MethodGet, reqHdr, respCode, respHdr, strictRFC) {
			t.Errorf("CanCache returned false for no-store request and strict RFC disabled")
		}
	}

	// test client no-cache no-store is ignored without strict RFC - tests RFC7234§5.2.1.5 violation to protect origins
	{
		reqHdr := http.Header{
			"Cache-Control": {"no-cache no-store"},
		}
		respCode := 200
		respHdr := http.Header{}
		strictRFC := false

		if !CanCache(http.MethodGet, reqHdr, respCode, respHdr, strictRFC) {
			t.Errorf("CanCache returned false for no-store request and strict RFC disabled")
		}
	}

	// test client no-cache is ignored without strict RFC - tests RFC7234§5.2.1.4 violation to protect origins
	{
		reqHdr := http.Header{
			"Cache-Control": {"no-cache"},
		}
		respHdr := http.Header{}
		reqCC := web.CacheControl{
			"no-cache": "",
		}
		respCC := web.CacheControl{}
		respReqHdrs := http.Header{}
		respReqTime := time.Now()
		respRespTime := time.Now()

		strictRFC := false

		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != remapdata.ReuseCan {
			t.Errorf("CanReuseStored for no-cache request and strict RFC disabled: expected ReuseCan, actual %v", reuse)
		}
	}

	// test client no-store is ignored without strict RFC - tests RFC7234§5.2.1.4 violation to protect origins
	{
		reqHdr := http.Header{
			"Cache-Control": {"no-store no-cache"},
		}
		respHdr := http.Header{}
		reqCC := web.CacheControl{
			"no-store": "",
			"no-cache": "",
		}
		respCC := web.CacheControl{}
		respReqHdrs := http.Header{}
		respReqTime := time.Now()
		respRespTime := time.Now()

		strictRFC := false

		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != remapdata.ReuseCan {
			t.Errorf("CanReuseStored for no-cache request and strict RFC disabled: expected ReuseCan, actual %v", reuse)
		}
	}

	// test client no-cache and no-store is ignored without strict RFC - tests RFC7234§5.2.1.4 violation to protect origins
	{
		reqHdr := http.Header{
			"Cache-Control": {"no-store no-cache"},
		}
		respHdr := http.Header{}
		reqCC := web.CacheControl{
			"no-store": "",
			"no-cache": "",
		}
		respCC := web.CacheControl{}
		respReqHdrs := http.Header{}
		respReqTime := time.Now()
		respRespTime := time.Now()

		strictRFC := false

		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != remapdata.ReuseCan {
			t.Errorf("CanReuseStored for no-cache request and strict RFC disabled: expected ReuseCan, actual %v", reuse)
		}
	}

	// test parent no-cache is obeyed with strict RFC
	{
		reqHdr := http.Header{}
		respCode := 200
		respHdr := http.Header{
			"Cache-Control": {"no-cache"},
		}
		strictRFC := false

		if CanCache(http.MethodGet, reqHdr, respCode, respHdr, strictRFC) {
			t.Errorf("CanCache returned true for no-cache request and strict RFC disabled")
		}
	}

	// test parent no-store is obeyed with strict RFC
	{
		reqHdr := http.Header{}
		respCode := 200
		respHdr := http.Header{
			"Cache-Control": {"no-store"},
		}
		strictRFC := false

		if CanCache(http.MethodGet, reqHdr, respCode, respHdr, strictRFC) {
			t.Errorf("CanCache returned true for no-cache request and strict RFC disabled")
		}
	}

	// test parent no-cache is obeyed without strict RFC
	{
		reqHdr := http.Header{}
		respCode := 200
		respHdr := http.Header{
			"Cache-Control": {"no-cache"},
		}
		strictRFC := false

		if CanCache(http.MethodGet, reqHdr, respCode, respHdr, strictRFC) {
			t.Errorf("CanCache returned true for no-cache request and strict RFC enabled")
		}
	}

	// test parent no-store is obeyed without strict RFC
	{
		reqHdr := http.Header{}
		respCode := 200
		respHdr := http.Header{
			"Cache-Control": {"no-store"},
		}
		strictRFC := true

		if CanCache(http.MethodGet, reqHdr, respCode, respHdr, strictRFC) {
			t.Errorf("CanCache returned true for no-cache request and strict RFC enabled")
		}
	}

	// test parent Expires in future is reused
	{
		now := time.Now()
		tenMinsBeforeExpires := now.Add(time.Minute * -10)
		expires := now.Format(time.RFC1123)
		reqHdr := http.Header{}
		respHdr := http.Header{
			"Expires": {expires},
		}
		reqCC := web.CacheControl{}
		respCC := web.CacheControl{}
		respReqHdrs := http.Header{}
		respReqTime := tenMinsBeforeExpires
		respRespTime := tenMinsBeforeExpires
		strictRFC := false
		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != remapdata.ReuseCan {
			t.Errorf("CanReuseStored request after expires: expected ReuseCan, actual %v", reuse)
		}
	}

	// test parent Expires in past has revaldiate and can stale
	{
		now := time.Now()
		tenMinutesAgo := now.Add(time.Minute * -10)
		expires := tenMinutesAgo.Format(time.RFC1123)
		reqHdr := http.Header{}
		respHdr := http.Header{
			"Expires": {expires},
			"Date":    {expires},
		}
		reqCC := web.CacheControl{}
		respCC := web.CacheControl{}
		respReqHdrs := http.Header{}
		respReqTime := now
		respRespTime := now
		strictRFC := false
		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != remapdata.ReuseMustRevalidateCanStale {
			t.Errorf("CanReuseStored request after expires: expected ReuseMustRevalidateCanStale, actual %v", reuse)
		}
	}

	// test parent Expires in past with must-revaldiate, cannot stale
	{
		now := time.Now()
		tenMinutesAgo := now.Add(time.Minute * -10)
		expires := tenMinutesAgo.Format(time.RFC1123)
		reqHdr := http.Header{}
		respHdr := http.Header{
			"Expires":       {expires},
			"Date":          {expires},
			"Cache-Control": {"must-revalidate"},
		}
		reqCC := web.CacheControl{}
		respCC := web.CacheControl{
			"must-revalidate": ""}
		respReqHdrs := http.Header{}
		respReqTime := now
		respRespTime := now
		strictRFC := false
		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != remapdata.ReuseMustRevalidate {
			t.Errorf("CanReuseStored request after expires and response must-revalidate: expected ReuseMustRevalidate, actual %v", reuse)
		}
	}

	// test parent Expires in past proxy-revaldiate, and no-stale
	{
		now := time.Now()
		tenMinutesAgo := now.Add(time.Minute * -10)
		expires := tenMinutesAgo.Format(time.RFC1123)
		reqHdr := http.Header{}
		respHdr := http.Header{
			"Expires":       {expires},
			"Date":          {expires},
			"Cache-Control": {"must-revalidate"},
		}
		reqCC := web.CacheControl{}
		respCC := web.CacheControl{
			"must-revalidate": ""}
		respReqHdrs := http.Header{}
		respReqTime := now
		respRespTime := now
		strictRFC := false
		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != remapdata.ReuseMustRevalidate {
			t.Errorf("CanReuseStored request after expires and response must-revalidate: expected ReuseMustRevalidate, actual %v", reuse)
		}
	}

	// test parent Expires in past with proxy-revaldiate returns MustRevalidateNoStale
	{
		now := time.Now()
		tenMinutesAgo := now.Add(time.Minute * -10)
		expires := tenMinutesAgo.Format(time.RFC1123)
		reqHdr := http.Header{}
		respHdr := http.Header{
			"Expires":       {expires},
			"Date":          {expires},
			"Cache-Control": {"proxy-revalidate"},
		}
		reqCC := web.CacheControl{}
		respCC := web.CacheControl{
			"proxy-revalidate": ""}
		respReqHdrs := http.Header{}
		respReqTime := now
		respRespTime := now
		strictRFC := false
		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != remapdata.ReuseMustRevalidate {
			t.Errorf("CanReuseStored request after expires and response must-revalidate: expected ReuseMustRevalidate, actual %v", reuse)
		}
	}

	// test parent max-age in future returns CanReuse
	{
		now := time.Now()
		tenMinutesAgo := now.Add(time.Minute * -10)
		reqHdr := http.Header{}
		respHdr := http.Header{
			"Date":          {tenMinutesAgo.Format(time.RFC1123)},
			"Cache-Control": {"max-age=1200"}, // 20 minutes
		}
		reqCC := web.CacheControl{}
		respCC := web.CacheControl{
			"max-age": "1200"}
		respReqHdrs := http.Header{}
		respReqTime := now
		respRespTime := now
		strictRFC := false
		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != remapdata.ReuseCan {
			t.Errorf("CanReuseStored request after expires and response must-revalidate: expected ReuseCan, actual %v", reuse)
		}
	}

	// test parent max-age in past returns MustRevalidate
	{
		now := time.Now()
		tenMinutesAgo := now.Add(time.Minute * -10)
		reqHdr := http.Header{}
		respHdr := http.Header{
			"Date":          {tenMinutesAgo.Format(time.RFC1123)},
			"Cache-Control": {"max-age=300"}, // 5 minutes
		}
		reqCC := web.CacheControl{}
		respCC := web.CacheControl{
			"max-age": "300"}
		respReqHdrs := http.Header{}
		respReqTime := tenMinutesAgo
		respRespTime := tenMinutesAgo
		strictRFC := false
		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != remapdata.ReuseMustRevalidateCanStale {
			t.Errorf("CanReuseStored request after response max-age: expected ReuseMustRevalidateCanStale, actual %v", reuse)
		}
	}

	// test parent s-maxage in future returns CanReuse
	{
		now := time.Now()
		tenMinutesAgo := now.Add(time.Minute * -10)
		reqHdr := http.Header{}
		respHdr := http.Header{
			"Date":          {tenMinutesAgo.Format(time.RFC1123)},
			"Cache-Control": {"s-maxage=1200"}, // 20 minutes
		}
		reqCC := web.CacheControl{}
		respCC := web.CacheControl{
			"s-maxage": "1200"}
		respReqHdrs := http.Header{}
		respReqTime := now
		respRespTime := now
		strictRFC := false
		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != remapdata.ReuseCan {
			t.Errorf("CanReuseStored request before s-maxage: expected ReuseCan, actual %v", reuse)
		}
	}

	// test parent s-age in past returns MustRevalidate
	{
		now := time.Now()
		tenMinutesAgo := now.Add(time.Minute * -10)
		reqHdr := http.Header{}
		respHdr := http.Header{
			"Date":          {tenMinutesAgo.Format(time.RFC1123)},
			"Cache-Control": {"s-maxage=300"}, // 5 minutes
		}
		reqCC := web.CacheControl{}
		respCC := web.CacheControl{
			"s-maxage": "300"}
		respReqHdrs := http.Header{}
		respReqTime := tenMinutesAgo
		respRespTime := tenMinutesAgo
		strictRFC := false
		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != remapdata.ReuseMustRevalidateCanStale {
			t.Errorf("CanReuseStored request after response s-maxage: expected ReuseMustRevalidateCanStale, actual %v", reuse)
		}
	}

	// test parent future s-maxage overrides past max-age and  returns CanReuse
	{
		now := time.Now()
		tenMinutesAgo := now.Add(time.Minute * -10)
		reqHdr := http.Header{}
		respHdr := http.Header{
			"Date":          {tenMinutesAgo.Format(time.RFC1123)},
			"Cache-Control": {"max-age=300,s-maxage=1200"},
		}
		reqCC := web.CacheControl{}
		respCC := web.CacheControl{
			"s-maxage": "1200",
			"max-age":  "300",
		}
		respReqHdrs := http.Header{}
		respReqTime := tenMinutesAgo
		respRespTime := tenMinutesAgo
		strictRFC := false
		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != remapdata.ReuseCan {
			t.Errorf("CanReuseStored request before s-maxage but after max-age: expected ReuseCan, actual %v", reuse)
		}
	}

	// test parent past s-maxage overrides future max-age and returns MustReval
	{
		now := time.Now()
		tenMinutesAgo := now.Add(time.Minute * -10)
		reqHdr := http.Header{}
		respHdr := http.Header{
			"Date":          {tenMinutesAgo.Format(time.RFC1123)},
			"Cache-Control": {"max-age=1200,s-maxage=300"},
		}
		reqCC := web.CacheControl{}
		respCC := web.CacheControl{
			"s-maxage": "300",
			"max-age":  "1200",
		}
		respReqHdrs := http.Header{}
		respReqTime := tenMinutesAgo
		respRespTime := tenMinutesAgo
		strictRFC := false
		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != remapdata.ReuseMustRevalidateCanStale {
			t.Errorf("CanReuseStored request before s-maxage but after max-age: expected ReuseMustRevalidateCanStale, actual %v", reuse)
		}
	}

	// test parent future max-age overrides past Expires and returns CanReuse
	{
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
		reqCC := web.CacheControl{}
		respCC := web.CacheControl{
			"max-age": "1200",
		}
		respReqHdrs := http.Header{}
		respReqTime := tenMinutesAgo
		respRespTime := tenMinutesAgo
		strictRFC := false
		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != remapdata.ReuseCan {
			t.Errorf("CanReuseStored request before max-age but after Expires: expected ReuseCan, actual %v", reuse)
		}
	}

	// test parent past max-age overrides future Expires and returns MustRevalidate
	{
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
		reqCC := web.CacheControl{}
		respCC := web.CacheControl{
			"max-age": "300",
		}
		respReqHdrs := http.Header{}
		respReqTime := tenMinutesAgo
		respRespTime := tenMinutesAgo
		strictRFC := false
		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != remapdata.ReuseMustRevalidateCanStale {
			t.Errorf("CanReuseStored request after max-age but before Expires: expected ReuseMustRevalidateCanStale, actual %v", reuse)
		}
	}

	// test parent future s-maxage overrides past Expires and returns CanReuse
	{
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
		reqCC := web.CacheControl{}
		respCC := web.CacheControl{
			"s-maxage": "1200",
		}
		respReqHdrs := http.Header{}
		respReqTime := tenMinutesAgo
		respRespTime := tenMinutesAgo
		strictRFC := false
		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != remapdata.ReuseCan {
			t.Errorf("CanReuseStored request before s-maxage but after Expires: expected ReuseCan, actual %v", reuse)
		}
	}

	// test parent past s-maxage overrides future Expires and returns MustRevalidate
	{
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
		reqCC := web.CacheControl{}
		respCC := web.CacheControl{
			"s-maxage": "300",
		}
		respReqHdrs := http.Header{}
		respReqTime := tenMinutesAgo
		respRespTime := tenMinutesAgo
		strictRFC := false
		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != remapdata.ReuseMustRevalidateCanStale {
			t.Errorf("CanReuseStored request after max-age but before Expires: expected ReuseMustRevalidateCanStale, actual %v", reuse)
		}
	}

	// test parent future s-maxage overrides past Expires and past max-age and returns CanReuse
	{
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
		reqCC := web.CacheControl{}
		respCC := web.CacheControl{
			"s-maxage": "1200",
			"max-age":  "300",
		}
		respReqHdrs := http.Header{}
		respReqTime := tenMinutesAgo
		respRespTime := tenMinutesAgo
		strictRFC := false
		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != remapdata.ReuseCan {
			t.Errorf("CanReuseStored request before s-maxage but after Expires: expected ReuseCan, actual %v", reuse)
		}
	}

	// test parent past s-maxage overrides future Expires and future max-age and returns MustRevalidate
	{
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
		reqCC := web.CacheControl{}
		respCC := web.CacheControl{
			"s-maxage": "300",
			"max-age":  "1200",
		}
		respReqHdrs := http.Header{}
		respReqTime := tenMinutesAgo
		respRespTime := tenMinutesAgo
		strictRFC := false
		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != remapdata.ReuseMustRevalidateCanStale {
			t.Errorf("CanReuseStored request after s-maxage but before Expires: expected ReuseMustRevalidateCanStale, actual %v", reuse)
		}
	}

	// test client min-fresh is obeyed with strict RFC - tests RFC7234§5.2.1.3 compliance
	{
		now := time.Now()
		tenMinutesAgo := now.Add(time.Minute * -10)
		reqHdr := http.Header{
			"Cache-Control": {"min-fresh=900"},
		}
		respHdr := http.Header{
			"Date":          {tenMinutesAgo.Format(time.RFC1123)},
			"Cache-Control": {"max-age=1200"},
		}
		reqCC := web.CacheControl{
			"min-fresh": "900",
		}
		respCC := web.CacheControl{
			"max-age": "1200",
		}
		respReqHdrs := http.Header{}
		respReqTime := tenMinutesAgo
		respRespTime := tenMinutesAgo
		strictRFC := true
		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != remapdata.ReuseMustRevalidate {
			t.Errorf("CanReuseStored request with strictRFC min-fresh 300 with 600 remaining: expected ReuseMustRevalidate, actual %v", reuse)
		}
	}

	// test client min-fresh is ignored without strict RFC - tests RFC7234§5.2.1.3 violation to protect origins
	{
		now := time.Now()
		tenMinutesAgo := now.Add(time.Minute * -10)
		reqHdr := http.Header{
			"Cache-Control": {"min-fresh=900"},
		}
		respHdr := http.Header{
			"Date":          {tenMinutesAgo.Format(time.RFC1123)},
			"Cache-Control": {"max-age=1200"},
		}
		reqCC := web.CacheControl{
			"min-fresh": "900",
		}
		respCC := web.CacheControl{
			"max-age": "1200",
		}
		respReqHdrs := http.Header{}
		respReqTime := tenMinutesAgo
		respRespTime := tenMinutesAgo
		strictRFC := false
		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != remapdata.ReuseCan {
			t.Errorf("CanReuseStored request with strictRFC min-fresh 1200 with 600 remaining: expected ReuseCan, actual %v", reuse)
		}
	}

	// test client min-fresh is ignored without strict RFC - tests RFC7234§5.2.1.3 violation to protect origins
	{
		now := time.Now()
		tenMinutesAgo := now.Add(time.Minute * -10)
		reqHdr := http.Header{
			"Cache-Control": {"min-fresh=900"},
		}
		respHdr := http.Header{
			"Date":          {tenMinutesAgo.Format(time.RFC1123)},
			"Cache-Control": {"max-age=1200"},
		}
		reqCC := web.CacheControl{
			"min-fresh": "900",
		}
		respCC := web.CacheControl{
			"max-age": "1200",
		}
		respReqHdrs := http.Header{}
		respReqTime := tenMinutesAgo
		respRespTime := tenMinutesAgo
		strictRFC := false
		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != remapdata.ReuseCan {
			t.Errorf("CanReuseStored request with strictRFC min-fresh 1200 with 600 remaining: expected ReuseCan, actual %v", reuse)
		}
	}

	// test default-cacheable response is cached. Tests RFC7234§5.2.1.3 and RFC7231§6.1 compliance
	{
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
	}

	// test non-default-cacheable response with no Cache-Control is not cached. Tests RFC7234§5.2.1.3 compliance
	{
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
	}

	// test non-default-cacheable response with Cache-Control is cached. Tests RFC7234§3 compliance
	{
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
	}

	log.Init(log.NopCloser(os.Stdout), log.NopCloser(os.Stdout), log.NopCloser(os.Stdout), log.NopCloser(os.Stdout), log.NopCloser(os.Stdout))
}
