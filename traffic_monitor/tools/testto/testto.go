package main

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
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

func main() {
	port := flag.Int("port", 8000, "Port to serve on")
	flag.Parse()
	if *port < 0 || *port > 65535 {
		fmt.Println("port must be 0-65535")
		return
	}

	toDataThs := NewThs()
	toDataThs.Set(&FakeTOData{Servers: []tc.ServerV50{}})

	Serve(*port, toDataThs)
	fmt.Printf("Serving on %v\n", *port)

	for {
		// TODO handle sighup to die
		time.Sleep(time.Hour)
	}
}

type FakeTOData struct {
	Monitoring tc.TrafficMonitorConfig
	CRConfig   tc.CRConfig
	Servers    []tc.ServerV50
}

// TODO make timeouts configurable?

const readTimeout = time.Second * 10
const writeTimeout = time.Second * 10

func Serve(port int, fakeTOData Ths) *http.Server {
	// TODO add HTTPS
	server := &http.Server{
		Addr:           ":" + strconv.Itoa(port),
		Handler:        RouteHandler(fakeTOData),
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Println("Error serving on port " + strconv.Itoa(port) + ": " + err.Error())
		}
	}()
	return server
}

type Route struct {
	Regex   *regexp.Regexp
	Handler http.HandlerFunc
}

func GetRoutes(fakeTOData Ths) []Route {
	routes := []Route{}
	for route, makeHandler := range Routes {
		routeRegex := regexp.MustCompile(route)
		routes = append(routes, Route{Regex: routeRegex, Handler: makeHandler(fakeTOData)})
	}
	return routes
}

type MakeHandlerFunc func(fakeTOData Ths) http.HandlerFunc

var Routes = map[string]MakeHandlerFunc{
	`/api/(.*)/user/login/?(\.json)?$`:          loginHandler,
	`/api/(.*)/cdns/(.*)/configs/monitoring/?$`: monitoringHandler,
	`/api/(.*)/servers/?(\.json)?$`:             serversHandler,
	`/api/(.*)/cdns/(.*)/snapshot/?(\.json)?$`:  crConfigHandler,
	`/api/(.*)/ping`:                            pingHandler,
}

func RouteHandler(fakeTOData Ths) http.HandlerFunc {
	routes := GetRoutes(fakeTOData)
	return func(w http.ResponseWriter, r *http.Request) {
		for _, route := range routes {
			if route.Regex.MatchString(r.URL.Path) {
				route.Handler(w, r)
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
	}
}

func pingHandler(fakeTOData Ths) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set(rfc.ContentType, rfc.ApplicationJSON)
			w.Write(append([]byte(`{"ping":"pong"}`), '\n'))
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func monitoringHandler(fakeTOData Ths) http.HandlerFunc {
	return makeJSONGetPostHandler(fakeTOData, monitoringHandlerGet, monitoringHandlerPost)
}

func serversHandler(fakeTOData Ths) http.HandlerFunc {
	return makeJSONGetPostHandler(fakeTOData, serversHandlerGet, serversHandlerPost)
}

func crConfigHandler(fakeTOData Ths) http.HandlerFunc {
	return makeJSONGetPostHandler(fakeTOData, crConfigHandlerGet, crConfigHandlerPost)
}

func makeJSONGetPostHandler(
	fakeTOData Ths,
	makeGetHandler func(fakeTOData Ths) http.HandlerFunc,
	makePostHandler func(fakeTOData Ths) http.HandlerFunc,
) http.HandlerFunc {
	getHandler := makeGetHandler(fakeTOData)
	postHandler := makePostHandler(fakeTOData)
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			postHandler(w, r)
		case http.MethodGet:
			getHandler(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func monitoringHandlerPost(fakeTOData Ths) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		obj := (*FakeTOData)(fakeTOData.Get())
		postJSONObj(w, r, &obj.Monitoring, obj, fakeTOData)
	}
}

func monitoringHandlerGet(fakeTOData Ths) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSONObj(w, ((*FakeTOData)(fakeTOData.Get())).Monitoring)
	}
}

func serversHandlerGet(fakeTOData Ths) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSONObj(w, ((*FakeTOData)(fakeTOData.Get())).Servers)
	}
}

func serversHandlerPost(fakeTOData Ths) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		obj := (*FakeTOData)(fakeTOData.Get())
		postJSONObj(w, r, &obj.Servers, obj, fakeTOData)
	}
}

func crConfigHandlerGet(fakeTOData Ths) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		crConfig := ((*FakeTOData)(fakeTOData.Get())).CRConfig
		if crConfig.Stats.CDNName == nil && crConfig.Stats.TMHost == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write(append([]byte(http.StatusText(http.StatusNotFound)), '\n'))
			return
		}
		writeJSONObj(w, crConfig)
	}
}

func crConfigHandlerPost(fakeTOData Ths) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		obj := (*FakeTOData)(fakeTOData.Get())
		postJSONObj(w, r, &obj.CRConfig, obj, fakeTOData)
	}
}

func writeJSONObj(w http.ResponseWriter, obj interface{}) {
	bts, err := json.MarshalIndent(&obj, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "marshalling: ` + err.Error() + `"}`))
		return
	}
	// TODO write content type
	w.Write([]byte(`{"response":`))
	w.Write(bts)
	w.Write([]byte(`}`))
	w.Write([]byte("\n"))
}

func postJSONObj(w http.ResponseWriter, r *http.Request, obj interface{}, fakeTOData ThsT, fakeTODataThs Ths) {
	if err := json.NewDecoder(r.Body).Decode(obj); err != nil {
		log.Println("unmarshall:" + err.Error() + ", " + r.URL.String())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "unmarshalling posted body: ` + err.Error() + `"}`))
		return
	}
	fakeTODataThs.Set(fakeTOData)
	w.WriteHeader(http.StatusNoContent)
}

type ThsT *FakeTOData

type Ths struct {
	v *ThsT
	m *sync.RWMutex
}

func NewThs() Ths {
	v := ThsT(nil)
	return Ths{
		m: &sync.RWMutex{},
		v: &v,
	}
}

func (t Ths) Set(v ThsT) {
	t.m.Lock()
	defer t.m.Unlock()
	*t.v = v
}

func (t Ths) Get() ThsT {
	t.m.RLock()
	defer t.m.RUnlock()
	return *t.v
}

func loginHandler(Ths) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Set-Cookie", `mojolicious=fake; Path=/; Expires=Thu, 13 Dec 2018 21:21:33 GMT; HttpOnly`)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"alerts":[{"text": "Successfully logged in.","level": "success"}]}`))
	}
}
