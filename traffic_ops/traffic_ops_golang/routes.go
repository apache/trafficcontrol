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
	"crypto/tls"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"time"

	tclog "github.com/apache/incubator-trafficcontrol/lib/go-log"
)

var Authenticated = true
var NoAuth = false

func handlerToFunc(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}
}

// Routes returns the routes, and a catchall route for when no route matches.
func Routes(d ServerData) ([]Route, http.Handler, error) {
	proxyHandler := rootHandler(d)

	routes := []Route{
		//ASNs
		{1.2, http.MethodGet, `asns(\.json)?$`, ASNsHandler(d.DB), ASNsPrivLevel, Authenticated, nil},
		//CDNs
		{1.2, http.MethodGet, `cdns(\.json)?$`, cdnsHandler(d.DB), CDNsPrivLevel, Authenticated, nil},
		{1.2, http.MethodGet, `cdns/{name}/configs/monitoring(\.json)?$`, monitoringHandler(d.DB), MonitoringPrivLevel, Authenticated, nil},
		// Delivery services
		{1.3, http.MethodGet, "deliveryservices/{xml-id}/urisignkeys$", getUrisignkeysHandler(d.DB, d.Config), PrivLevelAdmin, Authenticated, nil},
		{1.3, http.MethodPost, "deliveryservices/{xml-id}/urisignkeys$", assignDeliveryServiceUriKeysHandler(d.DB, d.Config), PrivLevelAdmin, Authenticated, nil},
		//Divisions
		{1.2, http.MethodGet, `divisions(\.json)?$`, divisionsHandler(d.DB), DivisionsPrivLevel, Authenticated, nil},
		//HwInfo
		{1.2, http.MethodGet, `hwinfo-wip(\.json)?$`, hwInfoHandler(d.DB), HWInfoPrivLevel, Authenticated, nil},
		//Parameters
		{1.2, http.MethodGet, `parameters(\.json)?$`, parametersHandler(d.DB), ParametersPrivLevel, Authenticated, nil},
		//Regions
		{1.2, http.MethodGet, `regions(\.json)?$`, regionsHandler(d.DB), RegionsPrivLevel, Authenticated, nil},
		{1.2, http.MethodGet, "regions/{id}$", regionsHandler(d.DB), RegionsPrivLevel, Authenticated, nil},
		//Servers
		// explicitly passed to legacy system until fully implemented.  Auth handled by legacy system.
		{1.2, http.MethodGet, "servers/checks$", handlerToFunc(proxyHandler), 0, NoAuth, []Middleware{}},
		{1.2, http.MethodGet, "servers/details$", handlerToFunc(proxyHandler), 0, NoAuth, []Middleware{}},
		{1.2, http.MethodGet, "servers/status$", handlerToFunc(proxyHandler), 0, NoAuth, []Middleware{}},
		{1.2, http.MethodGet, "servers/totals$", handlerToFunc(proxyHandler), 0, NoAuth, []Middleware{}},

		{1.2, http.MethodGet, `servers(\.json)?$`, serversHandler(d.DB), ServersPrivLevel, Authenticated, nil},
		{1.2, http.MethodGet, "servers/{id}$", serversHandler(d.DB), ServersPrivLevel, Authenticated, nil},
		{1.2, http.MethodPost, "servers/{id}/deliveryservices$", assignDeliveryServicesToServerHandler(d.DB), PrivLevelOperations, Authenticated, nil},
		{1.2, http.MethodGet, "servers/{host_name}/update_status$", getServerUpdateStatusHandler(d.DB), PrivLevelReadOnly, Authenticated, nil},

		//Statuses
		{1.2, http.MethodGet, `statuses(\.json)?$`, statusesHandler(d.DB), StatusesPrivLevel, Authenticated, nil},
		{1.2, http.MethodGet, "statuses/{id}$", statusesHandler(d.DB), StatusesPrivLevel, Authenticated, nil},
		//System
		{1.2, http.MethodGet, `system/info(\.json)?$`, systemInfoHandler(d.DB), SystemInfoPrivLevel, Authenticated, nil},
	}
	return routes, proxyHandler, nil
}

// RootHandler returns the / handler for the service, which reverse-proxies the old Perl Traffic Ops
func rootHandler(d ServerData) http.Handler {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout:   time.Duration(d.Config.ProxyTimeout) * time.Second,
			KeepAlive: time.Duration(d.Config.ProxyKeepAlive) * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   time.Duration(d.Config.ProxyTLSTimeout) * time.Second,
		ResponseHeaderTimeout: time.Duration(d.Config.ProxyReadHeaderTimeout) * time.Second,
		//Other knobs we can turn: ExpectContinueTimeout,IdleConnTimeout
	}
	rp := httputil.NewSingleHostReverseProxy(d.URL)
	rp.Transport = tr

	var logger interface{}
	logger, err := tclog.GetLogWriter(d.Config.ErrorLog())
	if err != nil {
		tclog.Errorln("could not create error log writer for proxy: ", err)
	}
	rp.ErrorLog = log.New(logger.(io.Writer), "proxy error: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC) //if we don't provide a logger to the reverse proxy it logs to stdout/err and is lost when ran by a script.
	tclog.Debugf("our reverseProxy: %++v\n", rp)
	tclog.Debugf("our reverseProxy's transport: %++v\n", tr)
	loggingProxyHandler := wrapAccessLog(d.Secrets[0], rp)

	managerHandler := CreateThrottledHandler(loggingProxyHandler, d.BackendMaxConnections["mojolicious"])
	return managerHandler
}

//CreateThrottledHandler takes a handler, and a max and uses a channel to insure the handler is used concurrently by only max number of routines
func CreateThrottledHandler(handler http.Handler, maxConcurrentCalls int) ThrottledHandler {
	return ThrottledHandler{handler, make(chan struct{}, maxConcurrentCalls)}
}

type ThrottledHandler struct {
	Handler http.Handler
	ReqChan chan struct{}
}

func (m ThrottledHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.ReqChan <- struct{}{}
	defer func() { <-m.ReqChan }()
	m.Handler.ServeHTTP(w, r)
}
