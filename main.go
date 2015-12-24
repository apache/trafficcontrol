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

package main

import (
	db "./todb"

	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {

	db.InitializeDatabase(os.Args[1], os.Args[2], os.Args[3])
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/hello/{name}", index).Methods("GET")
	router.HandleFunc("/api/2.0/raw/{table}.json", handleTable)
	router.HandleFunc("/api/2.0/{cdn}/CRConfig.json", handleCRConfig)

	log.Fatal(http.ListenAndServe(":8080", router))
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

func index(w http.ResponseWriter, r *http.Request) {
	log.Println("Responding to /hello request")
	log.Println(r.UserAgent())

	vars := mux.Vars(r)
	name := vars["name"]

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Hello:", name)
}
