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
	"database/sql"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/common/log"
)

const RoutePrefix = "api" // TODO config?

type Route struct {
	// Order matters! Do not reorder this! Routes() uses positional construction for readability.
	Version float64
	Method  string
	Path    string
	Handler RegexHandlerFunc
}

type ServerData struct {
	Config
	DB *sql.DB
}

type PathParams map[string]string

type RegexHandlerFunc func(w http.ResponseWriter, r *http.Request, params PathParams)

type CompiledRoute struct {
	Handler RegexHandlerFunc
	Regex   *regexp.Regexp
	Params  []string
}

func getSortedRouteVersions(rs []Route) []float64 {
	m := map[float64]struct{}{}
	for _, r := range rs {
		m[r.Version] = struct{}{}
	}
	versions := []float64{}
	for v, _ := range m {
		versions = append(versions, v)
	}
	sort.Float64s(versions)
	return versions
}

type PathHandler struct {
	Path    string
	Handler RegexHandlerFunc
}

// CreateRouteMap returns a map of methods to a slice of paths and handlers. Uses Semantic Versioning: routes are added to every subsequent minor version, but not subsequent major versions. For example, a 1.2 route is added to 1.3 but not 2.1. Also truncates '2.0' to '2', creating succinct major versions.
func CreateRouteMap(rs []Route) map[string][]PathHandler {
	// TODO strong types for method, path
	versions := getSortedRouteVersions(rs)
	m := map[string][]PathHandler{}
	for _, r := range rs {
		versionI := sort.SearchFloat64s(versions, r.Version)
		nextMajorVer := float64(int(r.Version) + 1)
		for _, version := range versions[versionI:] {
			if version >= nextMajorVer {
				break
			}
			vstr := strconv.FormatFloat(version, 'f', -1, 64)
			path := RoutePrefix + "/" + vstr + "/" + r.Path
			m[r.Method] = append(m[r.Method], PathHandler{Path: path, Handler: r.Handler})
			log.Infof("adding route %v %v\n", r.Method, path)
		}
	}
	return m
}

// CompiledRoutes takes a map of methods to paths and handlers, and returns a map of methods to CompiledRoutes.
func CompileRoutes(routes map[string][]PathHandler) map[string][]CompiledRoute {
	compiledRoutes := map[string][]CompiledRoute{}
	for method, mRoutes := range routes {
		for _, pathHandler := range mRoutes {
			route := pathHandler.Path
			handler := pathHandler.Handler
			var params []string
			for open := strings.Index(route, "{"); open > 0; open = strings.Index(route, "{") {
				close := strings.Index(route, "}")
				if close < 0 {
					panic("malformed route")
				}
				param := route[open+1 : close]

				params = append(params, param)
				route = route[:open] + `(.+)` + route[close+1:]
			}
			regex := regexp.MustCompile(route)
			compiledRoutes[method] = append(compiledRoutes[method], CompiledRoute{Handler: handler, Regex: regex, Params: params})
		}
	}
	return compiledRoutes
}

func Handler(routes map[string][]CompiledRoute, catchall http.Handler, w http.ResponseWriter, r *http.Request) {
	requested := r.URL.Path[1:]

	mRoutes, ok := routes[r.Method]
	if !ok {
		catchall.ServeHTTP(w, r)
		return
	}

	for _, compiledRoute := range mRoutes {
		match := compiledRoute.Regex.FindStringSubmatch(requested)
		if len(match) == 0 {
			continue
		}

		params := map[string]string{}
		for i, v := range compiledRoute.Params {
			params[v] = match[i+1]
		}
		compiledRoute.Handler(w, r, params)
		return
	}
	catchall.ServeHTTP(w, r)
}

func RegisterRoutes(d ServerData) error {
	routeSlice, catchall, err := Routes(d)
	if err != nil {
		return err
	}
	routes := CreateRouteMap(routeSlice)
	compiledRoutes := CompileRoutes(routes)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Handler(compiledRoutes, catchall, w, r)
	})
	return nil
}
