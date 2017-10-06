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
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
)

// Routes returns the routes, and a catchall route for when no route matches.
func Routes(d ServerData) ([]Route, http.Handler, error) {
	rd, err := routeData(d)
	if err != nil {
		return nil, nil, err
	}
	return []Route{
		{1.2, http.MethodGet, "cdns/{cdn}/configs/monitoring$", wrapHeaders(wrapAuth(monitoringHandler(d.DB), d.Insecure, d.Secrets[0], rd.PrivLevelStmt, MonitoringPrivLevel))},
		{1.2, http.MethodGet, "cdns/{cdn}/configs/monitoring.json$", wrapHeaders(wrapAuth(monitoringHandler(d.DB), d.Insecure, d.Secrets[0], rd.PrivLevelStmt, MonitoringPrivLevel))},
		{1.2, http.MethodGet, "regions-wip/{id}$", wrapHeaders(wrapAuthWithData(regionsHandler(d.DB), d.Insecure, d.Secrets[0], rd.PrivLevelStmt, ServersPrivLevel))},
		{1.2, http.MethodGet, "regions-wip.json$", wrapHeaders(wrapAuthWithData(regionsHandler(d.DB), d.Insecure, d.Secrets[0], rd.PrivLevelStmt, ServersPrivLevel))},
		{1.2, http.MethodGet, "regions-wip$", wrapHeaders(wrapAuthWithData(regionsHandler(d.DB), d.Insecure, d.Secrets[0], rd.PrivLevelStmt, ServersPrivLevel))},
		{1.2, http.MethodGet, "regions-wip.json$", wrapHeaders(wrapAuthWithData(serversHandler(d.DB), d.Insecure, d.Secrets[0], rd.PrivLevelStmt, ServersPrivLevel))},
		{1.2, http.MethodGet, "servers-wip.json$", wrapHeaders(wrapAuthWithData(serversHandler(d.DB), d.Insecure, d.Secrets[0], rd.PrivLevelStmt, ServersPrivLevel))},
		{1.2, http.MethodGet, "servers-wip$", wrapHeaders(wrapAuthWithData(serversHandler(d.DB), d.Insecure, d.Secrets[0], rd.PrivLevelStmt, ServersPrivLevel))},
		{1.2, http.MethodGet, "servers-wip.json$", wrapHeaders(wrapAuthWithData(serversHandler(d.DB), d.Insecure, d.Secrets[0], rd.PrivLevelStmt, ServersPrivLevel))},
		{1.2, http.MethodGet, "asns-wip$", wrapHeaders(wrapAuthWithData(ASNsHandler(d.DB), d.Insecure, d.Secrets[0], rd.PrivLevelStmt, ServersPrivLevel))},
		{1.2, http.MethodGet, "asns-wip.json$", wrapHeaders(wrapAuthWithData(ASNsHandler(d.DB), d.Insecure, d.Secrets[0], rd.PrivLevelStmt, ServersPrivLevel))},
		{1.2, http.MethodGet, "cdns-wip$", wrapHeaders(wrapAuthWithData(cdnsHandler(d.DB), d.Insecure, d.Secrets[0], rd.PrivLevelStmt, CDNsPrivLevel))},
		{1.2, http.MethodGet, "cdns-wip.json$", wrapHeaders(wrapAuthWithData(cdnsHandler(d.DB), d.Insecure, d.Secrets[0], rd.PrivLevelStmt, CDNsPrivLevel))},
		{1.2, http.MethodPost, "servers/{server}/deliveryservices$", wrapHeaders(wrapAuthWithData(assignDeliveryServicesToServerHandler(d.DB), d.Insecure, d.Secrets[0], rd.PrivLevelStmt, PrivLevelOperations))},
		{1.2, http.MethodGet, "divisions-wip$", wrapHeaders(wrapAuthWithData(divisionsHandler(d.DB), d.Insecure, d.Secrets[0], rd.PrivLevelStmt, CDNsPrivLevel))},
		{1.2, http.MethodGet, "divisions-wip.json$", wrapHeaders(wrapAuthWithData(divisionsHandler(d.DB), d.Insecure, d.Secrets[0], rd.PrivLevelStmt, DivisionsPrivLevel))},
		{1.2, http.MethodGet, "hwinfo-wip$", wrapHeaders(wrapAuthWithData(hwInfoHandler(d.DB), d.Insecure, d.Secrets[0], rd.PrivLevelStmt, HWInfoPrivLevel))},
		{1.2, http.MethodGet, "hwinfo-wip.json$", wrapHeaders(wrapAuthWithData(hwInfoHandler(d.DB), d.Insecure, d.Secrets[0], rd.PrivLevelStmt, HWInfoPrivLevel))},
		{1.2, http.MethodGet, "parameters-wip$", wrapHeaders(wrapAuthWithData(parametersHandler(d.DB), d.Insecure, d.Secrets[0], rd.PrivLevelStmt, ParametersPrivLevel))},
		{1.2, http.MethodGet, "parameters-wip.json$", wrapHeaders(wrapAuthWithData(parametersHandler(d.DB), d.Insecure, d.Secrets[0], rd.PrivLevelStmt, ParametersPrivLevel))},
		{1.2, http.MethodGet, "system/info-wip$", wrapHeaders(wrapAuthWithData(systemInfoHandler(d.DB), d.Insecure, d.Secrets[0], rd.PrivLevelStmt, SystemInfoPrivLevel))},
		{1.2, http.MethodGet, "system/info-wip.json$", wrapHeaders(wrapAuthWithData(systemInfoHandler(d.DB), d.Insecure, d.Secrets[0], rd.PrivLevelStmt, SystemInfoPrivLevel))},
	}, rootHandler(d), nil
}

type RouteData struct {
	PrivLevelStmt *sql.Stmt
}

func routeData(d ServerData) (RouteData, error) {
	rd := RouteData{}
	err := error(nil)

	if rd.PrivLevelStmt, err = preparePrivLevelStmt(d.DB); err != nil {
		return rd, fmt.Errorf("Error preparing db priv level query: ", err)
	}

	return rd, nil
}

// RootHandler returns the / handler for the service, which reverse-proxies the old Perl Traffic Ops
func rootHandler(d ServerData) http.Handler {
	// debug
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
