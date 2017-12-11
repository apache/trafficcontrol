package fakesrvr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor/tools/testcaches/fakesrvrdata"
)

// TODO config?
const readTimeout = time.Second * 10
const writeTimeout = time.Second * 10

func reqIsApplicationSystem(r *http.Request) bool {
	return r.URL.Query().Get("application") == "system"
}

func astatsHandler(fakeSrvrDataThs fakesrvrdata.Ths) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		srvr := (*fakesrvrdata.FakeServerData)(fakeSrvrDataThs.Get())
		// TODO cast to System, if query string `application=system`
		b := []byte{}
		err := error(nil)
		if reqIsApplicationSystem(r) {
			system := srvr.GetSystem()
			b, err = json.MarshalIndent(&system, "", "  ") // TODO debug, change to Marshal
		} else {
			b, err = json.MarshalIndent(&srvr, "", "  ") // TODO debug, change to Marshal
		}
		if err != nil {
			w.Write([]byte(`{"error": "marshalling: ` + err.Error() + `"}`)) // TODO escape error for JSON
		}
		w.Write(b)
	}
}

func Serve(port int, fakeSrvrData fakesrvrdata.Ths) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/_astats", astatsHandler(fakeSrvrData))
	server := &http.Server{
		Addr:           ":" + strconv.Itoa(port),
		Handler:        mux,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			// TODO pass the error somewhere, somehow?
			fmt.Println("Error serving on port " + strconv.Itoa(port) + ": " + err.Error())
		}
	}()
	return server
}
