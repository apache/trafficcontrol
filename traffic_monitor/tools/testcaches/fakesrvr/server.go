package fakesrvr

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"encoding/json"
	"fmt"
	"html"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/tools/testcaches/fakesrvrdata"
)

// TODO config?
const readTimeout = time.Second * 10
const writeTimeout = time.Second * 10

func reqIsApplicationSystem(r *http.Request) bool {
	return r.URL.Query().Get("application") == "system"
}

func astatsHandler(fakeSrvrDataThs fakesrvrdata.Ths) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(rfc.ContentType, rfc.ApplicationJSON)
		srvr := (*fakesrvrdata.FakeServerData)(fakeSrvrDataThs.Get())

		delayMSPtr := (*fakesrvrdata.MinMaxUint64)(atomic.LoadPointer(fakeSrvrDataThs.DelayMS))
		minDelayMS := delayMSPtr.Min
		maxDelayMS := delayMSPtr.Max

		if maxDelayMS != 0 {
			delayMS := minDelayMS
			if minDelayMS != maxDelayMS {
				delayMS += uint64(rand.Int63n(int64(maxDelayMS - minDelayMS)))
			}
			delay := time.Duration(delayMS) * time.Millisecond
			time.Sleep(delay)
		}

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
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "marshalling: ` + err.Error() + `"}`)) // TODO escape error for JSON
			return
		}
		w.Write(b)
	}
}

func cmdHandler(fakeSrvrDataThs fakesrvrdata.Ths) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := html.EscapeString(r.URL.Path)
		path = strings.ToLower(path)
		path = strings.TrimLeft(path, "/cmd")
		for cmd, cmdF := range cmds {
			if strings.HasPrefix(path, cmd) {
				cmdF(w, r, fakeSrvrDataThs)
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("command '" + path + "' not found\n"))
	}
}

func Serve(port int, fakeSrvrData fakesrvrdata.Ths) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/_astats", astatsHandler(fakeSrvrData))
	mux.HandleFunc("/cmd", cmdHandler(fakeSrvrData))
	mux.HandleFunc("/cmd/", cmdHandler(fakeSrvrData))
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
