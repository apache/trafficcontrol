package plugin

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
	"strconv"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/v8/grove/web"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
)

const NSPerSec = 1000000000

func init() {
	AddPlugin(20000, Funcs{afterRespond: atsLog})
}

func atsLog(icfg interface{}, d AfterRespondData) {
	now := time.Now()
	bytesSent := web.TryGetBytesWritten(d.W, d.Conn, d.BytesWritten)

	proxyHierarchyStr, proxyNameStr := getParentStrings(d.RespCode, d.CacheHit, d.ProxyStr, d.ToFQDN)

	log.EventRaw(atsEventLogStr(
		now,
		d.ClientIP,
		d.Hostname,
		d.Req.Host,
		d.Port,
		d.ToFQDN,
		d.Scheme,
		d.Req.URL.String(),
		d.Req.Method,
		d.Req.Proto,
		d.RespCode,
		now.Sub(d.ReqTime),
		bytesSent,
		d.OriginCode,
		d.OriginBytes,
		d.RespSuccess,
		d.OriginReqSuccess,
		getCacheHitStr(d.CacheHit, d.OriginConnectFailed),
		proxyHierarchyStr,
		proxyNameStr,
		d.Req.UserAgent(),
		d.Req.Header.Get("X-Money-Trace"),
		d.RequestID,
	))
}

// getParentStrings returns the phr and pqsn ATS log event strings (in that order).
// This covers almost all occurences that we currently see from ATS.
func getParentStrings(code int, hit bool, proxyStr string, toFQDN string) (string, string) {
	// the most common case (hopefully), do this first
	if hit {
		return "NONE", "-"
	}
	if code >= 200 {
		if proxyStr != "" {
			return "PARENT_HIT", strings.Split(proxyStr, ":")[0]
		}
		return "DIRECT", toFQDN
	}
	return "EMPTY", "-"
}

// getCacheHitStr returns the event log string for whether the request was a cache hit. For a request not in the cache, pass `ReuseCannot` to indicate a cache miss.
func getCacheHitStr(hit bool, originConnectFailed bool) string {
	if originConnectFailed {
		return "ERR_CONNECT_FAIL"
	}
	if hit {
		return "TCP_HIT"
	}
	return "TCP_MISS"
}

func atsEventLogStr(
	timestamp time.Time, // (prefix)
	clientIP string, // chi
	selfHostname string,
	reqHost string, // phn
	reqPort string, // php
	originHost string, // shn
	scheme string, // url
	url string, // url
	method string, // cqhm
	protocol string, // cqhv
	respCode int, // pssc
	timeToServe time.Duration, // ttms
	bytesSent uint64, // b
	originStatus int, // sssc
	originBytes uint64, // sscl
	clientRespSuccess bool, // cfsc
	originReqSuccess bool, // pfsc
	cacheHit string, // crc
	proxyUsed string, // phr
	thisProxyName string, // pqsn
	clientUserAgent string, // client user agent
	xmt string, // moneytrace header
	requestID uint64, // Grove tracing ID - not part of real ATS log format
) string {
	unixNano := timestamp.UnixNano()
	unixSec := unixNano / NSPerSec
	unixFrac := (unixNano / (NSPerSec / 1000)) - (unixSec * 1000) // gives fractional seconds to three decimal points, like the ATS logs.
	unixFracStr := strconv.FormatInt(unixFrac, 10)
	for len(unixFracStr) < 3 {
		unixFracStr = "0" + unixFracStr // leading zeros, so e.g. a fraction of '42' becomes '1234.042' not '1234.42'
	}
	cfsc := "FIN"
	if !clientRespSuccess {
		cfsc = "INTR"
	}
	pfsc := "FIN"
	if !originReqSuccess {
		pfsc = "INTR"
	}

	// TODO escape quotes within useragent, moneytrace
	clientUserAgent = `"` + clientUserAgent + `"`
	if xmt == "" {
		xmt = `"-"`
	} else {
		xmt = `"` + xmt + `"`
	}

	// 	1505408269.011 chi=2001:beef:cafe:f::2 phn=cdn-ec-nyc-001-01.nyc.kabletown.net php=80 shn=disc-org.kabletown.net url=http://edge.disc.kabletown.net/250001/3306/lb.xml cqhm=GET cqhv=HTTP/1.1 pssc=200 ttms=0 b=1778 sssc=000 sscl=0 cfsc=FIN pfsc=FIN crc=TCP_MEM_HIT phr=NONE pqsn=- uas="Go-http-client/1.1" xmt="-"
	return strconv.FormatInt(unixSec, 10) + "." + unixFracStr + " chi=" + clientIP + " phn=" + selfHostname + " php=" + reqPort + " shn=" + originHost + " url=" + scheme + "://" + reqHost + url + " cqhn=" + method + " cqhv=" + protocol + " pssc=" + strconv.FormatInt(int64(respCode), 10) + " ttms=" + strconv.FormatInt(int64(timeToServe/time.Millisecond), 10) + " b=" + strconv.FormatInt(int64(bytesSent), 10) + " sssc=" + strconv.FormatInt(int64(originStatus), 10) + " sscl=" + strconv.FormatInt(int64(originBytes), 10) + " cfsc=" + cfsc + " pfsc=" + pfsc + " crc=" + cacheHit + " phr=" + proxyUsed + " pqsn=" + thisProxyName + " uas=" + clientUserAgent + " xmt=" + xmt + " reqid=" + strconv.FormatUint(requestID, 10) + "\n"
}
