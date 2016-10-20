package http_server

import (
	"fmt"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/common/log"
	"github.com/hydrogen18/stoppableListener"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)

func GetCommonAPIData(params url.Values, time time.Time) CommonAPIData {
	return CommonAPIData{
		QueryParams: ParametersStr(params),
		DateStr:     DateStr(time.Now()),
	}
}

type CommonAPIData struct {
	QueryParams string `json:"pp"`
	DateStr     string `json:"date"`
}

// Server is a re-runnable HTTP server. Server.Run() may be called repeatedly, and
// each time the previous running server will be stopped, and the server will be
// restarted with the new port address and data request channel.
type Server struct {
	getData                    GetDataFunc
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
		"/publish/CacheStats/":          s.dataRequestFunc(CacheStats),
		"/publish/CacheStats":           s.dataRequestFunc(CacheStats),
		"/publish/CrConfig/":            s.dataRequestFunc(TRConfig),
		"/publish/CrConfig":             s.dataRequestFunc(TRConfig),
		"/publish/CrStates/":            s.handleCrStatesFunc(),
		"/publish/CrStates":             s.handleCrStatesFunc(),
		"/publish/DsStats/":             s.dataRequestFunc(DSStats),
		"/publish/DsStats":              s.dataRequestFunc(DSStats),
		"/publish/EventLog/":            s.dataRequestFunc(EventLog),
		"/publish/EventLog":             s.dataRequestFunc(EventLog),
		"/publish/PeerStates/":          s.dataRequestFunc(PeerStates),
		"/publish/PeerStates":           s.dataRequestFunc(PeerStates),
		"/publish/StatSummary/":         s.dataRequestFunc(StatSummary),
		"/publish/StatSummary":          s.dataRequestFunc(StatSummary),
		"/publish/Stats/":               s.dataRequestFunc(Stats),
		"/publish/Stats":                s.dataRequestFunc(Stats),
		"/publish/ConfigDoc/":           s.dataRequestFunc(ConfigDoc),
		"/publish/ConfigDoc":            s.dataRequestFunc(ConfigDoc),
		"/api/cache-count/":             s.dataRequestFunc(APICacheCount),
		"/api/cache-count":              s.dataRequestFunc(APICacheCount),
		"/api/cache-available-count/":   s.dataRequestFunc(APICacheAvailableCount),
		"/api/cache-available-count":    s.dataRequestFunc(APICacheAvailableCount),
		"/api/cache-down-count/":        s.dataRequestFunc(APICacheDownCount),
		"/api/cache-down-count":         s.dataRequestFunc(APICacheDownCount),
		"/api/version/":                 s.dataRequestFunc(APIVersion),
		"/api/version":                  s.dataRequestFunc(APIVersion),
		"/api/traffic-ops-uri/":         s.dataRequestFunc(APITrafficOpsURI),
		"/api/traffic-ops-uri":          s.dataRequestFunc(APITrafficOpsURI),
		"/api/cache-statuses/":          s.dataRequestFunc(APICacheStates),
		"/api/cache-statuses":           s.dataRequestFunc(APICacheStates),
		"/api/bandwidth-kbps/":          s.dataRequestFunc(APIBandwidthKbps),
		"/api/bandwidth-kbps":           s.dataRequestFunc(APIBandwidthKbps),
		"/api/bandwidth-capacity-kbps/": s.dataRequestFunc(APIBandwidthCapacityKbps),
		"/api/bandwidth-capacity-kbps":  s.dataRequestFunc(APIBandwidthCapacityKbps),
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
func (s Server) Run(f GetDataFunc, addr string, readTimeout time.Duration, writeTimeout time.Duration) error {
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

	s.getData = f

	sm := http.NewServeMux()
	err = s.registerEndpoints(sm)
	if err != nil {
		return err
	}
	server := &http.Server{
		Addr:           addr,
		Handler:        sm,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
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
	APIBandwidthCapacityKbps
)

func (t Type) String() string {
	switch t {
	case TRConfig:
		return "TRConfig"
	case TRStateDerived:
		return "TRStateDerived"
	case TRStateSelf:
		return "TRStateSelf"
	case CacheStats:
		return "CacheStats"
	case DSStats:
		return "DSStats"
	case EventLog:
		return "EventLog"
	case PeerStates:
		return "PeerStates"
	case StatSummary:
		return "StatSummary"
	case Stats:
		return "Stats"
	case ConfigDoc:
		return "ConfigDoc"
	case APICacheCount:
		return "APICacheCount"
	case APICacheAvailableCount:
		return "APICacheAvailableCount"
	case APICacheDownCount:
		return "APICacheDownCount"
	case APIVersion:
		return "APIVersion"
	case APITrafficOpsURI:
		return "APITrafficOpsURI"
	case APICacheStates:
		return "APICacheStates"
	case APIBandwidthKbps:
		return "APIBandwidthKbps"
	case APIBandwidthCapacityKbps:
		return "APIBandwidthCapacityKbps"
	default:
		return "Invalid"
	}
}

type Format int

const (
	XML Format = (1 << iota)
	JSON
)

type DataRequest struct {
	Type
	Format
	Date       string
	Parameters map[string][]string
}

type GetDataFunc func(DataRequest) ([]byte, int)

// ParametersStr takes the URL query parameters, and returns a string as used by the Traffic Monitor 1.0 endpoints "pp" key.
func ParametersStr(params url.Values) string {
	fmt.Println("debug4 ParametersStr 0")
	pp := ""
	for param, vals := range params {
		for _, val := range vals {
			pp += param + "=[" + val + "], "
		}
	}
	if len(pp) > 2 {
		pp = pp[:len(pp)-2]
	}
	return pp
}

// DateStr returns the given time in the format expected by Traffic Monitor 1.0 API users
func DateStr(t time.Time) string {
	return t.UTC().Format("Mon Jan 02 15:04:05 UTC 2006")
}

func (s Server) dataRequest(w http.ResponseWriter, req *http.Request, t Type, f Format) {
	//pp: "0=[my-ats-edge-cache-0], hc=[1]",
	//dateLayout := "Thu Oct 09 20:28:36 UTC 2014"
	dateLayout := "Mon Jan 02 15:04:05 MST 2006"
	data, responseCode := s.getData(DataRequest{
		Type:       t,
		Format:     f,
		Date:       time.Now().UTC().Format(dateLayout),
		Parameters: req.URL.Query(),
	})
	if len(data) > 0 {
		w.WriteHeader(responseCode)
		w.Write(data)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}
}

func (s Server) handleRootFunc() (http.HandlerFunc, error) {
	return handleFile("index.html")
}

func (s Server) handleSortableFunc() (http.HandlerFunc, error) {
	return handleFile("sorttable.js")
}

func (s Server) handleFile(name string) (http.HandlerFunc, error) {
	index, err := ioutil.ReadFile(file)
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
