package http_server

import (
	"fmt"
	"net/http"
	"time"
)

var mgrReqChan chan DataRequest

func writeResponse(w http.ResponseWriter, f Format, dr DataRequest) {
	data := <-dr.C

	if len(data) > 0 {
		w.Write(data)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

func Run(c chan DataRequest, addr string) {
	mgrReqChan = c

	sm := http.NewServeMux()

	// note: with the trailing slash, any non-trailing slash requests will get a 301 redirect
	sm.HandleFunc("/", http.NotFound)
	sm.HandleFunc("/publish/CacheStats/", handleCacheStats)
	sm.HandleFunc("/publish/CrConfig", handleCrConfig)
	sm.HandleFunc("/publish/CrStates", handleCrStates)
	sm.HandleFunc("/publish/DsStats", handleDsStats)
	sm.HandleFunc("/publish/EventLog", handleEventLog)
	sm.HandleFunc("/publish/PeerStates", handlePeerStates)
	sm.HandleFunc("/publish/StatSummary", handleStatSummary)
	sm.HandleFunc("/publish/Stats", handleStats)

	s := &http.Server{
		Addr:           addr,
		Handler:        sm,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	s.ListenAndServe()
	fmt.Println("Web server listening on " + addr)
}
