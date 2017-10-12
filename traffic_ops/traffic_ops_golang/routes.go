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

var Authenticated = true
var NoAuth = false

// Routes returns the routes, and a catchall route for when no route matches.
func Routes(d ServerData) ([]Route, http.Handler, error) {

	routes := []Route{
		//ASNs
		{1.2, http.MethodGet, `asns-wip(\.json)?$`, ASNsHandler(d.DB), ASNsPrivLevel, Authenticated, nil},
		//CDNs
		{1.2, http.MethodGet, `cdns-wip(\.json)?$`, cdnsHandler(d.DB), CDNsPrivLevel, Authenticated, nil},
		{1.2, http.MethodGet, `cdns/{name}/configs/monitoring(\.json)?$`, monitoringHandler(d.DB), MonitoringPrivLevel, Authenticated, nil},
		// Delivery services
		{1.3, http.MethodGet, "deliveryservices/{xml-id}/urisignkeys$", urisignkeysHandler(d.DB, d.Config), PrivLevelAdmin, Authenticated, nil},
		{1.3, http.MethodPost, "deliveryservices/{xml-id}/urisignkeys$", assignDeliveryServiceUriKeysKeysHandler(d.DB, d.Config), PrivLevelAdmin, Authenticated, nil},
		//Divisions
		{1.2, http.MethodGet, `divisions-wip(\.json)?$`, divisionsHandler(d.DB), DivisionsPrivLevel, Authenticated, nil},
		//HwInfo
		{1.2, http.MethodGet, `hwinfo-wip(\.json)?$`, hwInfoHandler(d.DB), HWInfoPrivLevel, Authenticated, nil},
		//Parameters
		{1.2, http.MethodGet, `parameters-wip(\.json)?$`, parametersHandler(d.DB), ParametersPrivLevel, Authenticated, nil},
		//Regions
		{1.2, http.MethodGet, `regions-wip(\.json)?$`, regionsHandler(d.DB), RegionsPrivLevel, Authenticated, nil},
		{1.2, http.MethodGet, "regions-wip/{id}$", regionsHandler(d.DB), RegionsPrivLevel, Authenticated, nil},
		//Servers
		{1.2, http.MethodGet, `servers-wip(\.json)?$`, serversHandler(d.DB), ServersPrivLevel, Authenticated, nil},
		{1.2, http.MethodGet, "servers-wip/{id}$", serversHandler(d.DB), ServersPrivLevel, Authenticated, nil},
		{1.2, http.MethodPost, "servers/{id}/deliveryservices$", assignDeliveryServicesToServerHandler(d.DB), PrivLevelOperations, Authenticated, nil},

		//Statuses
		{1.2, http.MethodGet, `statuses-wip(\.json)?$`, statusesHandler(d.DB), StatusesPrivLevel, Authenticated, nil},
		{1.2, http.MethodGet, "statuses-wip/{id}$", statusesHandler(d.DB), StatusesPrivLevel, Authenticated, nil},
		//System
		{1.2, http.MethodGet, `system/info-wip(\.json)?$`, systemInfoHandler(d.DB), SystemInfoPrivLevel, Authenticated, nil},
	}
	return routes, rootHandler(d), nil
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

	rp.ErrorLog = log.Error //if we don't provide a logger to the reverse proxy it logs to stdout/err and is lost when ran by a script.
	log.Debugf("our reverseProxy: %++v\n", rp)
	log.Debugf("our reverseProxy's transport: %++v\n", tr)
	loggingProxyHandler := wrapAccessLog(d.Secrets[0], rp)
	return loggingProxyHandler
}
