
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Started from https://raw.githubusercontent.com/jordan-wright/gophish

package routes

import (
	"github.com/apache/trafficcontrol/traffic_ops/experimental/server/api"
	"github.com/apache/trafficcontrol/traffic_ops/experimental/server/auth"
	"github.com/apache/trafficcontrol/traffic_ops/experimental/server/csconfig"
	output "github.com/apache/trafficcontrol/traffic_ops/experimental/server/output_format"

	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

const apiPath = "/api/2.0/"

// CreateRouter creates the routes for handling requests to the web interface.
// This function returns an http.Handler to be used in http.ListenAndServe().
func CreateRouter(db *sqlx.DB) http.Handler {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc(apiPath+"login", wrapHeaders(auth.GetLoginOptionsFunc(), []api.ApiMethod{api.OPTIONS, api.POST})).Methods("OPTIONS")
	router.HandleFunc(apiPath+"login", wrapHeaders(auth.GetLoginFunc(db), []api.ApiMethod{api.OPTIONS, api.POST})).Methods("POST")
	router.HandleFunc(apiPath+"{table}", auth.Use(optionsHandler, auth.DONTRequireLogin)).Methods("OPTIONS")
	router.HandleFunc(apiPath+"{table}/{id}", auth.Use(optionsHandler, auth.DONTRequireLogin)).Methods("OPTIONS")
	router.HandleFunc(apiPath+"config/cr/{cdn}/CRConfig.json", auth.Use(getHandleCRConfigFunc(db), auth.RequireLogin))
	router.HandleFunc(apiPath+"config/csconfig/hostname/{hostname}/port/{port}", auth.Use(getHandleCSConfigFunc(db), auth.RequireLogin))
	addApiHandlers(router, db)
	return auth.Use(router.ServeHTTP, auth.GetContext)
}

// wrapHeaders wraps an http.HandlerFunc to call setHeaders with the given methods.
func wrapHeaders(f http.HandlerFunc, methods api.ApiMethods) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setHeaders(w, methods)
		f(w, r)
	}
}

// setHeaders writes the universal headers needed by all routes,
// along with the given accepted HTTP Methods.
func setHeaders(w http.ResponseWriter, methods api.ApiMethods) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", methods.String())
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Authorization, X-Requested-With, Content-Type")
}

// optionsHandler handles HTTP OPTIONS requests, writing the
// appropriate options and an HTTP OK.
func optionsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	table := vars["table"]
	route := table
	if _, ok := vars["id"]; ok {
		route += "/{id}"
	}

	apiHandlers := api.ApiHandlers()
	if tableHandlers, ok := apiHandlers[route]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	} else {
		setHeaders(w, tableHandlers.Methods())
		w.WriteHeader(http.StatusOK)
	}
}

// wrapApiHandler takes an api.ApiHandlerFunc and returns a func with the
// signature expected by http.HandleFunc.
//
// The returned func sets the headers, reads the params, calls the
// handler with the params, encodes the handlers response, and writes it.
func wrapApiHandler(f api.ApiHandlerFunc, methods api.ApiMethods, db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setHeaders(w, methods)
		body, err := ioutil.ReadAll(r.Body)
		response, err := f(mux.Vars(r), body, db)
		if err != nil {
			log.Println(err)
		}
		jresponse := output.MakeApiResponse(response, nil, err)
		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		enc.Encode(jresponse)
	}
}

// addApiHandlers adds each API handler to the given router.
func addApiHandlers(router *mux.Router, db *sqlx.DB) {
	for route, funcs := range api.ApiHandlers() {
		for method, f := range funcs {
			router.HandleFunc(apiPath+route, auth.Use(wrapApiHandler(f, funcs.Methods(), db), auth.RequireLogin)).Methods(method.String())
		}
	}
}

func getCrconfigSnapshot(cdn string, db *sqlx.DB) (string, error) {
	queryStr := `select snapshot from crconfig_snapshots where cdn = $1 and created_at = (select max(created_at) created_at from crconfig_snapshots where cdn = $1);`
	rows, err := db.Query(queryStr, cdn)
	if err != nil {
		return "", fmt.Errorf("getCrconfigSnapshot query error: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return "", fmt.Errorf("No Snapshot Found")
	}

	var snapshot string
	if err := rows.Scan(&snapshot); err != nil {
		return "", fmt.Errorf("getCrconfigSnapshot row error: %v", err)
	}
	return snapshot, rows.Err()
}

// getHandleCRConfigFunc returns a func which handles requests to the CRConfig endpoint,
// returning the encoded CRConfig data for the requested CDN.
func getHandleCRConfigFunc(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setHeaders(w, []api.ApiMethod{api.GET})
		vars := mux.Vars(r)
		cdn := vars["cdn"]

		snapshot, err := getCrconfigSnapshot(cdn, db)
		if err != nil {
			log.Println(err)
			json.NewEncoder(w).Encode(struct {
				Error string `json:"error"`
			}{Error: err.Error()})
			return
		}
		w.Write([]byte(snapshot))
	}
}

// getHandleCSConfigFunc returns a func which handles requests to the CSConfig endpoint,
// returning the encoded CSConfig data for the requested host.
func getHandleCSConfigFunc(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setHeaders(w, []api.ApiMethod{api.GET})
		enc := json.NewEncoder(w)
		vars := mux.Vars(r)
		hostName := vars["hostname"]
		portStr := vars["port"]
		port, err := strconv.Atoi(portStr)
		if err != nil {
			enc.Encode(struct {
				Error string `json:"error"`
			}{Error: err.Error()})
			return
		}

		resp, err := csconfig.GetCSConfig(hostName, int64(port), db)
		if err != nil {
			enc.Encode(struct {
				Error string `json:"error"`
			}{Error: err.Error()})
		} else {
			enc.Encode(resp)
		}
	}
}
