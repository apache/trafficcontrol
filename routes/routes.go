// Copyright 2015 Comcast Cable Communications Management, LLC

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
	"github.com/Comcast/traffic_control/traffic_ops/goto2/api"
	"github.com/Comcast/traffic_control/traffic_ops/goto2/auth"
	"github.com/Comcast/traffic_control/traffic_ops/goto2/crconfig"
	"github.com/Comcast/traffic_control/traffic_ops/goto2/csconfig"
	output "github.com/Comcast/traffic_control/traffic_ops/goto2/output_format"

	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

const apiPath = "/api/2.0/"

// CreateRouter creates the routes for handling requests to the web interface.
// This function returns an http.Handler to be used in http.ListenAndServe().
func CreateRouter() http.Handler {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/login", auth.Login).Methods("POST")
	router.HandleFunc(apiPath+"{table}", auth.Use(optionsHandler, auth.DONTRequireLogin)).Methods("OPTIONS")
	router.HandleFunc(apiPath+"{table}/{id}", auth.Use(optionsHandler, auth.DONTRequireLogin)).Methods("OPTIONS")
	router.HandleFunc("/config/cr/{cdn}/CRConfig.json", auth.Use(handleCRConfig, auth.RequireLogin))
	router.HandleFunc("/config/csconfig/{hostname}", auth.Use(handleCSConfig, auth.RequireLogin))
	addApiHandlers(router)
	return auth.Use(router.ServeHTTP, auth.GetContext)
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
func wrapApiHandler(f api.ApiHandlerFunc, methods api.ApiMethods) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		setHeaders(w, methods)
		body, err := ioutil.ReadAll(r.Body)
		response, err := f(mux.Vars(r), body)
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
func addApiHandlers(router *mux.Router) {
	for route, funcs := range api.ApiHandlers() {
		for method, f := range funcs {
			router.HandleFunc(apiPath+route, auth.Use(wrapApiHandler(f, funcs.Methods()), auth.RequireLogin)).Methods(method.String())
		}
	}
}

// handleCRConfig handles requests to the CRConfig endpoint,
// returning the encoded CRConfig data for the requested CDN.
func handleCRConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cdn := vars["cdn"]
	resp, _ := crconfig.GetCRConfig(cdn)
	enc := json.NewEncoder(w)
	enc.Encode(resp)
}

// handleCSConfig handles requests to the CSConfig endpoint,
// returning the encoded CSConfig data for the requested host.
func handleCSConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hostName := vars["hostname"]
	resp, _ := csconfig.GetCSConfig(hostName)
	enc := json.NewEncoder(w)
	enc.Encode(resp)
}
