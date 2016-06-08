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

// Endpoints returns a map of HTTP paths to functions.
// This is a function because Go doesn't have constant map literals.
func Endpoints() map[string]http.HandlerFunc {
	// note: with the trailing slash, any non-trailing slash requests will get a 301 redirect
	return map[string]http.HandlerFunc{
		"/": http.NotFound,
		"/publish/CacheStats/": handleCacheStats,
		"/publish/CrConfig":    handleCrConfig,
		"/publish/CrStates":    handleCrStates,
		"/publish/DsStats":     handleDsStats,
		"/publish/EventLog":    handleEventLog,
		"/publish/PeerStates":  handlePeerStates,
		"/publish/StatSummary": handleStatSummary,
		"/publish/Stats":       handleStats,
		"/publish/ConfigDoc":   handleConfigDoc,
	}
}

func RegisterEndpoints(sm *http.ServeMux) {
	for path, f := range Endpoints() {
		sm.HandleFunc(path, f)
	}
}

func Run(c chan DataRequest, addr string) {
	mgrReqChan = c

	sm := http.NewServeMux()
	RegisterEndpoints(sm)

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
