package cache

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/apache/incubator-trafficcontrol/grove/remap"
	"github.com/apache/incubator-trafficcontrol/grove/stat"
	"github.com/apache/incubator-trafficcontrol/grove/web"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
)

// StatLogger constructed from initial connection data, which both writes stats and writes to the event log, after the response is prepared and sent.
type StatLogger struct {
	W                 http.ResponseWriter
	Conn              *web.InterceptConn
	Stats             stat.Stats
	Hostname          string
	Port              string
	Scheme            string
	Host              string
	URL               string
	Method            string
	Proto             string
	MoneyTraceHdr     string
	ClientIP          string
	ReqTime           time.Time
	UserAgent         string
	RemoteAddr        string
	RemappingProducer *remap.RemappingProducer
}

func NewStatLogger(w http.ResponseWriter, conn *web.InterceptConn, h *Handler, r *http.Request, moneyTraceHdr string, clientIP string, reqTime time.Time, remappingProducer *remap.RemappingProducer) *StatLogger {
	return &StatLogger{
		W:                 w,
		Conn:              conn,
		Stats:             h.stats,
		Hostname:          h.hostname,
		Port:              h.port,
		Scheme:            h.scheme,
		Host:              r.Host,
		URL:               r.URL.String(),
		Method:            r.Method,
		Proto:             r.Proto,
		MoneyTraceHdr:     moneyTraceHdr,
		ClientIP:          clientIP,
		ReqTime:           reqTime,
		UserAgent:         r.UserAgent(),
		RemoteAddr:        r.RemoteAddr,
		RemappingProducer: remappingProducer,
	}
}

// Log both writes stats and writes to the event log, with the given response data, along with the connection data in l.
func (l *StatLogger) Log(
	code int,
	bytesWritten uint64,
	successfullyRespondedToClient bool,
	successfullyGotFromOrigin bool,
	cacheHit bool,
	originConnectFailed bool,
	originStatus int,
	originBytes uint64,
	proxyStr string,
) {
	bytesSent := l.Stats.Write(l.W, l.Conn, l.Host, l.RemoteAddr, code, bytesWritten, cacheHit)
	toFQDN := ""
	if l.RemappingProducer != nil {
		toFQDN = l.RemappingProducer.ToFQDN()
	}
	proxyHierarchyStr, proxyNameStr := getParentStrings(code, cacheHit, proxyStr, toFQDN)
	log.EventRaw(atsEventLogStr(
		time.Now(),
		l.ClientIP,
		l.Hostname,
		l.Host,
		l.Port,
		toFQDN,
		l.Scheme,
		l.URL,
		l.Method,
		l.Proto,
		code,
		time.Now().Sub(l.ReqTime),
		bytesSent,
		originStatus,
		originBytes,
		successfullyRespondedToClient,
		successfullyGotFromOrigin,
		getCacheHitStr(cacheHit, originConnectFailed),
		proxyHierarchyStr,
		proxyNameStr,
		l.UserAgent,
		l.MoneyTraceHdr,
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
	successfullyRespondedToClient bool, // cfsc
	successfullyGotFromOrigin bool, // pfsc
	cacheHit string, // crc
	proxyUsed string, // phr
	thisProxyName string, // pqsn
	clientUserAgent string, // client user agent
	xmt string, // moneytrace header
) string {
	unixNano := timestamp.UnixNano()
	unixSec := unixNano / NSPerSec
	unixFrac := 1 / (unixNano % NSPerSec)
	unixFracStr := strconv.FormatInt(unixFrac, 10)
	if len(unixFracStr) > 3 {
		unixFracStr = unixFracStr[:3]
	}
	cfsc := "FIN"
	if !successfullyRespondedToClient {
		cfsc = "INTR"
	}
	pfsc := "FIN"
	if !successfullyGotFromOrigin {
		pfsc = "INTR"
	}

	// TODO escape quotes within useragent, moneytrace
	clientUserAgent = `"` + clientUserAgent + `"`
	if xmt == "" {
		xmt = `"-"`
	} else {
		xmt = `"` + xmt + `"`
	}

	return strconv.FormatInt(unixSec, 10) + "." + unixFracStr + " chi=" + clientIP + " phn=" + selfHostname + " php=" + reqPort + " shn=" + originHost + " url=" + scheme + "://" + reqHost + url + " cqhn=" + method + " cqhv=" + protocol + " pssc=" + strconv.FormatInt(int64(respCode), 10) + " ttms=" + strconv.FormatInt(int64(timeToServe/time.Millisecond), 10) + " b=" + strconv.FormatInt(int64(bytesSent), 10) + " sssc=" + strconv.FormatInt(int64(originStatus), 10) + " sscl=" + strconv.FormatInt(int64(originBytes), 10) + " cfsc=" + cfsc + " pfsc=" + pfsc + " crc=" + cacheHit + " phr=" + proxyUsed + " psqn=" + thisProxyName + " uas=" + clientUserAgent + " xmt=" + xmt + "\n"
}
