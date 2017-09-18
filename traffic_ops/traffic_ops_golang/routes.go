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
	"net/http"
	"net/http/httputil"
)

// Routes returns the routes, and a catchall route for when no route matches.
func Routes(d ServerData) ([]Route, http.Handler, error) {
	rd, err := routeData(d)
	if err != nil {
		return nil, nil, err
	}
	return []Route{
		{1.2, http.MethodGet, "cdns/{cdn}/configs/monitoring$", wrapHeaders(wrapAuth(monitoringHandler(d.DB), d.Insecure, d.TOSecret, rd.PrivLevelStmt, MonitoringPrivLevel))},
		{1.2, http.MethodGet, "cdns/{cdn}/configs/monitoring.json$", wrapHeaders(wrapAuth(monitoringHandler(d.DB), d.Insecure, d.TOSecret, rd.PrivLevelStmt, MonitoringPrivLevel))},
		{1.2, http.MethodGet, "servers$", wrapHeaders(wrapAuthWithData(serversHandler(d.DB), d.Insecure, d.TOSecret, rd.PrivLevelStmt, ServersPrivLevel))},
		{1.2, http.MethodGet, "servers.json$", wrapHeaders(wrapAuthWithData(serversHandler(d.DB), d.Insecure, d.TOSecret, rd.PrivLevelStmt, ServersPrivLevel))},
		{1.2, http.MethodGet, "cdns$", wrapHeaders(wrapAuthWithData(cdnsHandler(d.DB), d.Insecure, d.TOSecret, rd.PrivLevelStmt, CdnsPrivLevel))},
		{1.2, http.MethodGet, "cdns.json$", wrapHeaders(wrapAuthWithData(cdnsHandler(d.DB), d.Insecure, d.TOSecret, rd.PrivLevelStmt, CdnsPrivLevel))},
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

// getRootHandler returns the / handler for the service, which reverse-proxies the old Perl Traffic Ops
func rootHandler(d ServerData) http.Handler {
	// debug
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	rp := httputil.NewSingleHostReverseProxy(d.TOURL)
	rp.Transport = tr

	loggingProxyHandler := wrapAccessLog(d.TOSecret, rp)
	return loggingProxyHandler
}
