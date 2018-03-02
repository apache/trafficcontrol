package plugin

import (
	"encoding/json"
	"net/http"
	"runtime"
	"strings"

	"github.com/apache/incubator-trafficcontrol/grove/web"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
)

func init() {
	AddPlugin(10000, Funcs{onRequest: memstats})
}

const MemStatsEndpoint = "/_memstats"

func memstats(icfg interface{}, d OnRequestData) bool {
	if !strings.HasPrefix(d.R.URL.Path, MemStatsEndpoint) {
		log.Debugf("plugin onrequest http_memstats returning, not in path '" + d.R.URL.Path + "'\n")
		return false
	}

	log.Debugf("plugin onrequest http_memstats calling\n")

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

	// TODO gzip
	// TODO cache for 1 second

	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)

	bytes, err := json.Marshal(stats)
	if err != nil {
		code := http.StatusInternalServerError
		w.WriteHeader(code)
		w.Write([]byte(http.StatusText(code)))
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
	return true
}
