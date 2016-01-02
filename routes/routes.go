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
	"../api"
	"../auth"
	"../crconfig"
	"../csconfig"
	output "../output_format"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// CreateAdminRouter creates the routes for handling requests to the web interface.
// This function returns an http.Handler to be used in http.ListenAndServe().
func CreateRouter() http.Handler {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/login", auth.LoginPage).Methods("GET")
	router.HandleFunc("/login", auth.Login).Methods("POST")
	router.HandleFunc("/logout", auth.Use(auth.Logout, auth.RequireLogin)).Methods("GET")

	router.HandleFunc("/api/2.0/{table}", auth.Use(apiHandler, auth.RequireLogin)).Methods("GET", "POST")
	router.HandleFunc("/api/2.0/{table}/{id}", auth.Use(apiHandler, auth.RequireLogin)).Methods("GET", "PUT", "DELETE")

	router.HandleFunc("/config/cr/{cdn}/CRConfig.json", auth.Use(handleCRConfig, auth.RequireLogin))
	router.HandleFunc("/config/csconfig/{hostname}", auth.Use(handleCSConfig, auth.RequireLogin))

	return auth.Use(router.ServeHTTP, auth.GetContext)
}

func handleCRConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cdn := vars["cdn"]
	resp, _ := crconfig.GetCRConfig(cdn)
	enc := json.NewEncoder(w)
	enc.Encode(resp)
}

func handleCSConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hostName := vars["hostname"]
	resp, _ := csconfig.GetCSConfig(hostName)
	enc := json.NewEncoder(w)
	enc.Encode(resp)
}

func setHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Authorization, X-Requested-With, Content-Type")
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	// maybe it is better to just pass (w http.ResponseWriter, r *http.Request) to the actions,
	// and have the action funcs write to w without returning?
	// TODO: handle admin/oper can CUD, rest can r
	// TODO: handle deliveryservice_tmuser for portal
	setHeaders(w)
	vars := mux.Vars(r)
	table := vars["table"]
	id := -1
	if vars["id"] != "" {
		num, err := strconv.Atoi(vars["id"])
		if err != nil {
			fmt.Println("error 323222")
		}
		id = num
	}
	body, err := ioutil.ReadAll(r.Body)
	response, err := api.Action(table, r.Method, id, body)
	if err != nil {
		fmt.Println("error 42 ", err)
	}
	jresponse := output.MakeApiResponse(response, nil, err)
	enc := json.NewEncoder(w)
	enc.Encode(jresponse)
}
