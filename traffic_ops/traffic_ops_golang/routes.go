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
	"net"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
)

// Routes returns the routes, and a catchall route for when no route matches.
func Routes(d ServerData) ([]Route, http.Handler, error) {
	routes := []Route{}
	//ASNs
	routes = append(routes, Route{Version: 1.2, Method: http.MethodGet, Path: `asns-wip(\.json)?$`, Handler: ASNsHandler(d.DB), RequiredPrivLevel: ServersPrivLevel})
	//CDNs
	routes = append(routes, Route{Version: 1.2, Method: http.MethodGet, Path: `cdns-wip(\.json)?$`, Handler: cdnsHandler(d.DB), RequiredPrivLevel: CdnsPrivLevel})
	routes = append(routes, Route{Version: 1.2, Method: http.MethodGet, Path: `cdns/{cdn}/configs/monitoring(\.json)?$`, Handler: monitoringHandler(d.DB), RequiredPrivLevel: MonitoringPrivLevel})
	//Divisions
	routes = append(routes, Route{Version: 1.2, Method: http.MethodGet, Path: `divisions-wip(\.json)?$`, Handler: divisionsHandler(d.DB), RequiredPrivLevel: DivisionsPrivLevel})
	//HwInfo
	routes = append(routes, Route{Version: 1.2, Method: http.MethodGet, Path: `hwinfo-wip(\.json)?$`, Handler: hwInfoHandler(d.DB), RequiredPrivLevel: HWInfoPrivLevel})
	//Parameters
	routes = append(routes, Route{Version: 1.2, Method: http.MethodGet, Path: `parameters-wip(\.json)?$`, Handler: parametersHandler(d.DB), RequiredPrivLevel: ParametersPrivLevel})
	//Regions
	routes = append(routes, Route{Version: 1.2, Method: http.MethodGet, Path: `regions-wip(\.json)?$`, Handler: regionsHandler(d.DB), RequiredPrivLevel: ServersPrivLevel})
	routes = append(routes, Route{Version: 1.2, Method: http.MethodGet, Path: "regions-wip/{id}$", Handler: regionsHandler(d.DB), RequiredPrivLevel: ServersPrivLevel})
	//Servers
	routes = append(routes, Route{Version: 1.2, Method: http.MethodGet, Path: `servers-wip(\.json)?$`, Handler: serversHandler(d.DB), RequiredPrivLevel: ServersPrivLevel})
	routes = append(routes, Route{Version: 1.2, Method: http.MethodGet, Path: "servers-wip/{id}$", Handler: serversHandler(d.DB), RequiredPrivLevel: ServersPrivLevel})
	routes = append(routes, Route{Version: 1.2, Method: http.MethodPost, Path: "servers/{server}/deliveryservices$", Handler: assignDeliveryServicesToServerHandler(d.DB), RequiredPrivLevel: PrivLevelOperations})
	//System
	routes = append(routes, Route{Version: 1.2, Method: http.MethodGet, Path: `system/info-wip(\.json)?$`, Handler: systemInfoHandler(d.DB), RequiredPrivLevel: SystemInfoPrivLevel})
	return routes, rootHandler(d), nil
}

// RootHandler returns the / handler for the service, which reverse-proxies the old Perl Traffic Ops
func rootHandler(d ServerData) http.Handler {
	tr := &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
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

	rp.ErrorLog = log.Error //if we don't provide a logger to the reverse proxy it logs to stdout/err and is lost when ran by a script.
	log.Debugf("our reverseProxy: %++v\n", rp)
	log.Debugf("our reverseProxy's transport: %++v\n", tr)
	loggingProxyHandler := wrapAccessLog(d.Secrets[0], rp)
	return loggingProxyHandler
}
