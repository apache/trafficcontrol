package http_server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

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
)

type Format int

const (
	XML Format = (1 << iota)
	JSON
)

type DataRequest struct {
	T          Type
	F          Format
	C          chan []byte
	Date       string
	Parameters map[string][]string
}

func dataRequest(w http.ResponseWriter, req *http.Request, t Type, f Format) {
	//pp: "0=[my-ats-edge-cache-0], hc=[1]",
	//dateLayout := "Thu Oct 09 20:28:36 UTC 2014"
	dateLayout := "Mon Jan 02 15:04:05 MST 2006"
	time := time.Now()
	p := map[string][]string{}

	for key, v := range req.URL.Query() {
		for _, value := range v {
			p[key] = append(p[key], value)
		}
	}

	dr := DataRequest{
		T:          t,
		F:          f,
		C:          make(chan []byte, 1), // must be buffered, so if this is killed, the writer doesn't block forever
		Date:       time.UTC().Format(dateLayout),
		Parameters: p,
	}

	mgrReqChan <- dr
	writeResponse(w, f, dr)
}

func handleRootFunc() (http.HandlerFunc, error) {
	index, err := ioutil.ReadFile("index.html")
	if err != nil {
		return nil, err
	}
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "%s", index)
	}, nil
}

func handleCrStates(w http.ResponseWriter, req *http.Request) {
	t := TRStateDerived
	if req.URL.RawQuery == "raw" {
		t = TRStateSelf
	}
	dataRequest(w, req, t, JSON)
}

func DataRequestFunc(t Type) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dataRequest(w, r, t, JSON)
	}
}
