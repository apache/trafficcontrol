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
	"context"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"

	"fmt"

	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/jmoiron/sqlx"
)

const RoutePrefix = "api" // TODO config?

type Middleware func(handlerFunc http.HandlerFunc) http.HandlerFunc

type Route struct {
	// Order matters! Do not reorder this! Routes() uses positional construction for readability.
	Version           float64
	Method            string
	Path              string
	Handler           http.HandlerFunc
	RequiredPrivLevel int
	Authenticated     bool
	Middlewares       []Middleware
}

func getDefaultMiddleware() []Middleware {
	return []Middleware{wrapHeaders}
}

type ServerData struct {
	Config
	DB *sqlx.DB
}

type CompiledRoute struct {
	Handler http.HandlerFunc
	Regex   *regexp.Regexp
	Params  []string
}

func getSortedRouteVersions(rs []Route) []float64 {
	m := map[float64]struct{}{}
	for _, r := range rs {
		m[r.Version] = struct{}{}
	}
	versions := []float64{}
	for v := range m {
		versions = append(versions, v)
	}
	sort.Float64s(versions)
	return versions
}

type PathHandler struct {
	Path    string
	Handler http.HandlerFunc
}

// CreateRouteMap returns a map of methods to a slice of paths and handlers; wrapping the handlers in the appropriate middleware. Uses Semantic Versioning: routes are added to every subsequent minor version, but not subsequent major versions. For example, a 1.2 route is added to 1.3 but not 2.1. Also truncates '2.0' to '2', creating succinct major versions.
func CreateRouteMap(rs []Route, authBase AuthBase) map[string][]PathHandler {
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

			middlewares := r.Middlewares

			if middlewares == nil {
				middlewares = getDefaultMiddleware()
			}
			if r.Authenticated { //a privLevel of zero is an unauthenticated endpoint.
				authWrapper := authBase.GetWrapper(r.RequiredPrivLevel)
				middlewares = append([]Middleware{authWrapper}, middlewares...)
			}

			m[r.Method] = append(m[r.Method], PathHandler{Path: path, Handler: use(r.Handler, middlewares)})

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
				route = route[:open] + `([^/]+)` + route[close+1:]
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

		ctx := r.Context()

		params := api.PathParams{}
		for i, v := range compiledRoute.Params {
			params[v] = match[i+1]
		}

		ctx = context.WithValue(ctx, api.PathParamsKey, params)
		compiledRoute.Handler(w, r.WithContext(ctx))
		return
	}
	catchall.ServeHTTP(w, r)
}

func RegisterRoutes(d ServerData) error {
	routeSlice, catchall, err := Routes(d)
	if err != nil {
		return err
	}

	userInfoStmt, err := prepareUserInfoStmt(d.DB)
	if err != nil {
		return fmt.Errorf("Error preparing db priv level query: %s", err)
	}

	authBase := AuthBase{d.Insecure, d.Config.Secrets[0], userInfoStmt, nil} //we know d.Config.Secrets is a slice of at least one or start up would fail.
	routes := CreateRouteMap(routeSlice, authBase)
	compiledRoutes := CompileRoutes(routes)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Handler(compiledRoutes, catchall, w, r)
	})
	return nil
}

func prepareUserInfoStmt(db *sqlx.DB) (*sqlx.Stmt, error) {
	return db.Preparex("SELECT r.priv_level, u.id, u.username, COALESCE(u.tenant_id, -1) AS tenant_id FROM tm_user AS u JOIN role AS r ON u.role = r.id WHERE u.username = $1")
}

func use(h http.HandlerFunc, middlewares []Middleware) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- { //apply them in reverse order so they are used in a natural order.
		h = middlewares[i](h)
	}
	return h
}
