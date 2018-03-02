package plugin

import (
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/apache/incubator-trafficcontrol/grove/web"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
)

func init() {
	AddPlugin(10000, Funcs{onRequest: callgc})
}

const CallGCEndpoint = "/_callgc"

func callgc(icfg interface{}, d OnRequestData) bool {
	if !strings.HasPrefix(d.R.URL.Path, CallGCEndpoint) {
		return false
	}
	reqTime := time.Now()

	log.Debugf("plugin onrequest http_callgc calling\n")

	w := d.W
	req := d.R

	// TODO access log? Stats byte count?
	ip, err := web.GetIP(req)
	if err != nil {
		code := http.StatusInternalServerError
		w.WriteHeader(code)
		w.Write([]byte(http.StatusText(code)))
		log.Errorln("statHandler ServeHTTP failed to get IP: " + ip.String())
		return true
	}
	if !d.StatRules.Allowed(ip) {
		code := http.StatusForbidden
		w.WriteHeader(code)
		w.Write([]byte(http.StatusText(code)))
		log.Debugln("statHandler.ServeHTTP IP " + ip.String() + " FORBIDDEN") // TODO event?
		return true
	}

	runtime.GC()

	respCode := http.StatusNoContent
	w.WriteHeader(respCode)

	clientIP, _ := web.GetClientIPPort(req)

	now := time.Now()
	// log, so we know if someone is hitting this endpoint when they shouldn't be. GC is expensive, this could become an accidental DDOS.
	log.EventRaw(atsEventLogStr(now, clientIP, d.Hostname, req.Host, d.Port, "-", d.Scheme, req.URL.String(), req.Method, req.Proto, respCode, now.Sub(reqTime), 0, 0, 0, true, true, getCacheHitStr(true, false), "-", "-", req.UserAgent(), req.Header.Get("X-Money-Trace")))

	return true
}
