package cache

import (
	"net/http"
	"time"

	"github.com/apache/incubator-trafficcontrol/grove/web"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
)

type StatLogger struct {
	W                 http.ResponseWriter
	Conn              *web.InterceptConn
	Stats             Stats
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
	RemappingProducer *RemappingProducer
}

func NewStatLogger(w http.ResponseWriter, conn *web.InterceptConn, h *CacheHandler, r *http.Request, moneyTraceHdr string, clientIP string, reqTime time.Time, remappingProducer *RemappingProducer) StatLogger {
	return StatLogger{
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

func (l *StatLogger) Log(code int, bytesWritten uint64, successfullyRespondedToClient bool, successfullyGotFromOrigin bool, cacheHitStr string, originStatus int, originBytes uint64) {
	bytesSent := WriteStats(l.Stats, l.W, l.Conn, l.Host, l.RemoteAddr, code, bytesWritten)
	toFQDN := ""
	proxyStr := ""
	if l.RemappingProducer != nil {
		toFQDN = l.RemappingProducer.ToFQDN()
		proxyStr = l.RemappingProducer.ProxyStr()
	}
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
		cacheHitStr,
		proxyStr,
		"-", // TODO fix?
		l.UserAgent,
		l.MoneyTraceHdr,
	))
}
