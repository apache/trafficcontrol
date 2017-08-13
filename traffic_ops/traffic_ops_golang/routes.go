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
	"fmt"
	"net/http"
	"net/http/httputil"
)

// GetRoutes returns the map of regex routes, and a catchall route for when no regex matches.
func GetRoutes(d ServerData) (map[string]RegexHandlerFunc, http.Handler, error) {
	privLevelStmt, err := preparePrivLevelStmt(d.DB)
	if err != nil {
		return nil, nil, fmt.Errorf("Error preparing db priv level query: ", err)
	}

	return map[string]RegexHandlerFunc{
		"api/1.2/cdns/{cdn}/configs/monitoring":      wrapHeaders(wrapAuth(monitoringHandler(d.DB), d.Insecure, d.TOSecret, privLevelStmt, MonitoringPrivLevel)),
		"api/1.2/cdns/{cdn}/configs/monitoring.json": wrapHeaders(wrapAuth(monitoringHandler(d.DB), d.Insecure, d.TOSecret, privLevelStmt, MonitoringPrivLevel)),
	}, getRootHandler(d), nil
}

// getRootHandler returns the / handler for the service, which reverse-proxies the old Perl Traffic Ops
func getRootHandler(d ServerData) http.Handler {
	// debug
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	rp := httputil.NewSingleHostReverseProxy(d.TOURL)
	rp.Transport = tr

	loggingProxyHandler := wrapAccessLog(d.TOSecret, rp)
	return loggingProxyHandler
}
