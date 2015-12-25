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

package router

import (
	auth "../auth"
	db "../todb"

	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// CreateAdminRouter creates the routes for handling requests to the web interface.
// This function returns an http.Handler to be used in http.ListenAndServe().
func CreateRouter() http.Handler {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/login", auth.LoginPage).Methods("GET")
	router.HandleFunc("/logout", auth.Use(auth.Logout, auth.RequireLogin)).Methods("GET")
	router.HandleFunc("/login", auth.Login).Methods("POST")
	router.HandleFunc("/hello/{name}", auth.Use(hello, auth.RequireLogin)).Methods("GET")

	router.HandleFunc("/api/2.0/raw/{table}.json", auth.Use(handleTable, auth.RequireLogin))
	router.HandleFunc("/api/2.0/{cdn}/CRConfig.json", auth.Use(handleCRConfig, auth.RequireLogin))
	return auth.Use(router.ServeHTTP, auth.GetContext)
}

func handleCRConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cdn := vars["cdn"]
	resp, _ := db.GetCRConfig(cdn)
	enc := json.NewEncoder(w)
	enc.Encode(resp)
}

func handleTable(w http.ResponseWriter, r *http.Request) {
	log.Println("Responding to /api request")
	log.Println(r.UserAgent())

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Authorization, X-Requested-With, Content-Type")

	vars := mux.Vars(r)
	table := vars["table"]

	rows, _ := db.GetTable(table)
	// fmt.Print(rows)
	enc := json.NewEncoder(w)
	enc.Encode(rows)
	// w.WriteHeader(http.StatusOK)
	// fmt.Fprintln(w, "table:", table)
}

func hello(w http.ResponseWriter, r *http.Request) {
	log.Println("Responding to /hello request")
	log.Println(r.UserAgent())

	vars := mux.Vars(r)
	name := vars["name"]

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Hello:", name)
}
