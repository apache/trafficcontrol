package http_server

import (
	"fmt"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/log"
	"github.com/hydrogen18/stoppableListener"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"time"
)

// Server is a re-runnable HTTP server. Server.Run() may be called repeatedly, and
// each time the previous running server will be stopped, and the server will be
// restarted with the new port address and data request channel.
type Server struct {
	mgrReqChan                 chan<- DataRequest
	stoppableListener          *stoppableListener.StoppableListener
	stoppableListenerWaitGroup sync.WaitGroup
}

// Endpoints returns a map of HTTP paths to functions.
// This is a function because Go doesn't have constant map literals.
func (s Server) endpoints() (map[string]http.HandlerFunc, error) {
	handleRoot, err := s.handleRootFunc()
	handleSortableJs, err := s.handleSortableFunc()
	if err != nil {
		return nil, fmt.Errorf("Error getting root endpoint: %v", err)
	}

	// note: with the trailing slash, any non-trailing slash requests will get a 301 redirect
	return map[string]http.HandlerFunc{
		"/publish/CacheStats/":       s.dataRequestFunc(CacheStats),
		"/publish/CrConfig":          s.dataRequestFunc(TRConfig),
		"/publish/CrStates":          s.handleCrStatesFunc(),
		"/publish/DsStats":           s.dataRequestFunc(DSStats),
		"/publish/EventLog":          s.dataRequestFunc(EventLog),
		"/publish/PeerStates":        s.dataRequestFunc(PeerStates),
		"/publish/StatSummary":       s.dataRequestFunc(StatSummary),
		"/publish/Stats":             s.dataRequestFunc(Stats),
		"/publish/ConfigDoc":         s.dataRequestFunc(ConfigDoc),
		"/api/cache-count":           s.dataRequestFunc(APICacheCount),
		"/api/cache-available-count": s.dataRequestFunc(APICacheAvailableCount),
		"/api/cache-down-count":      s.dataRequestFunc(APICacheDownCount),
		"/api/version":               s.dataRequestFunc(APIVersion),
		"/api/traffic-ops-uri":       s.dataRequestFunc(APITrafficOpsURI),
		"/api/cache-statuses":        s.dataRequestFunc(APICacheStates),
		"/api/bandwidth-kbps":        s.dataRequestFunc(APIBandwidthKbps),
		"/":             handleRoot,
		"/sorttable.js": handleSortableJs,
	}, nil
}

func (s Server) registerEndpoints(sm *http.ServeMux) error {
	endpoints, err := s.endpoints()
	if err != nil {
		return err
	}
	for path, f := range endpoints {
		sm.HandleFunc(path, f)
	}
	return nil
}

// Run runs a new HTTP service at the given addr, making data requests to the given c.
// Run may be called repeatedly, and each time, will shut down any existing service first.
// Run is NOT threadsafe, and MUST NOT be called concurrently by multiple goroutines.
func (s Server) Run(c chan<- DataRequest, addr string) error {
	// TODO make an object, which itself is not threadsafe, but which encapsulates all data so multiple
	//      objects can be created and Run.

	if s.stoppableListener != nil {
		log.Infof("Stopping Web Server\n")
		s.stoppableListener.Stop()
		s.stoppableListenerWaitGroup.Wait()
	}
	log.Infof("Starting Web Server\n")

	var err error
	var originalListener net.Listener
	if originalListener, err = net.Listen("tcp", addr); err != nil {
		return err
	}
	if s.stoppableListener, err = stoppableListener.New(originalListener); err != nil {
		return err
	}

	s.mgrReqChan = c

	sm := http.NewServeMux()
	err = s.registerEndpoints(sm)
	if err != nil {
		return err
	}
	server := &http.Server{
		Addr:           addr,
		Handler:        sm,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	s.stoppableListenerWaitGroup = sync.WaitGroup{}
	s.stoppableListenerWaitGroup.Add(1)
	go func() {
		defer s.stoppableListenerWaitGroup.Done()
		server.Serve(s.stoppableListener)
	}()

	log.Infof("Web server listening on %s", addr)
	return nil
}

type Type int

// TODO rename these, all caps isn't recommended Go style
const (
	TRConfig Type = (1 << iota)
	TRStateDerived
	TRStateSelf
	CacheStats
	DSStats
	EventLog
	PeerStates
	StatSummary
	Stats
	ConfigDoc
	APICacheCount
	APICacheAvailableCount
	APICacheDownCount
	APIVersion
	APITrafficOpsURI
	APICacheStates
	APIBandwidthKbps
)

type Format int

const (
	XML Format = (1 << iota)
	JSON
)

type DataRequest struct {
	Type
	Format
	Response   chan<- []byte
	Date       string
	Parameters map[string][]string
}

func writeResponse(w http.ResponseWriter, f Format, response <-chan []byte) {
	data := <-response
	if len(data) > 0 {
		w.Write(data)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

func (s Server) dataRequest(w http.ResponseWriter, req *http.Request, t Type, f Format) {
	//pp: "0=[my-ats-edge-cache-0], hc=[1]",
	//dateLayout := "Thu Oct 09 20:28:36 UTC 2014"
	dateLayout := "Mon Jan 02 15:04:05 MST 2006"
	response := make(chan []byte, 1) // must be buffered, so if this is killed, the writer doesn't block forever
	s.mgrReqChan <- DataRequest{
		Type:       t,
		Format:     f,
		Response:   response,
		Date:       time.Now().UTC().Format(dateLayout),
		Parameters: req.URL.Query(),
	}
	writeResponse(w, f, response)
}

func (s Server) handleRootFunc() (http.HandlerFunc, error) {
	index, err := ioutil.ReadFile("index.html")
	if err != nil {
		return nil, err
	}
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "%s", index)
	}, nil
}

func (s Server) handleSortableFunc() (http.HandlerFunc, error) {
	index, err := ioutil.ReadFile("sorttable.js")
	if err != nil {
		return nil, err
	}
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "%s", index)
	}, nil
}

func (s Server) handleCrStatesFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		t := TRStateDerived
		if req.URL.RawQuery == "raw" {
			t = TRStateSelf
		}
		s.dataRequest(w, req, t, JSON)
	}
}

func (s Server) dataRequestFunc(t Type) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.dataRequest(w, r, t, JSON)
	}
}
