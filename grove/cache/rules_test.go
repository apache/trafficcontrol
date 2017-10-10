package cache

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/apache/incubator-trafficcontrol/grove/web"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
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

		if CanCache(reqHdr, respCode, respHdr, strictRFC) {
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

		if !CanCache(reqHdr, respCode, respHdr, strictRFC) {
			t.Errorf("CanCache returned false for no-store request and strict RFC disabled")
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

		if !CanCache(reqHdr, respCode, respHdr, strictRFC) {
			t.Errorf("CanCache returned false for no-store request and strict RFC disabled")
		}
	}

	log.Init(log.NopCloser(os.Stdout), log.NopCloser(os.Stdout), log.NopCloser(os.Stdout), log.NopCloser(os.Stdout), log.NopCloser(os.Stdout))

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

		if reuse := CanReuseStored(reqHdr, respHdr, reqCC, respCC, respReqHdrs, respReqTime, respRespTime, strictRFC); reuse != ReuseCan {
			t.Errorf("CanReuseStored for no-cache request and strict RFC disabled: expected ReuseCan, actual %v", reuse)
		}
	}
}

func TestCanReuseStored(t *testing.T) {
	// func CanReuseStored(reqHeaders http.Header, respHeaders http.Header, reqCacheControl web.CacheControl, respCacheControl web.CacheControl, respReqHeaders http.Header, respReqTime time.Time, respRespTime time.Time, strictRFC bool) Reuse {

}
