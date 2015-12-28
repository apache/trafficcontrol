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
	// auth "./auth"
	router "./router"
	db "./todb"

	// "encoding/json"
	// "fmt"
	"github.com/gorilla/handlers"
	"log"
	"net/http"
	"os"
	// "github.com/gorilla/mux"
)

func main() {
	db.InitializeDatabase(os.Args[1], os.Args[2], os.Args[3], os.Args[4])
	// router := mux.NewRouter().StrictSlash(true)
	// router.HandleFunc("/login", auth.LoginPage).Methods("GET")
	// router.HandleFunc("/login", auth.Login).Methods("POST")
	// router.HandleFunc("/hello/{name}", auth.Use(index, auth.RequireLogin)).Methods("GET")
	// router.HandleFunc("/api/2.0/raw/{table}.json", handleTable)
	// router.HandleFunc("/api/2.0/{cdn}/CRConfig.json", handleCRConfig)

	// log.Fatal(http.ListenAndServe(":8080", router))
	var Logger = log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)
	Logger.Printf("Starting server...")
	http.ListenAndServe(":8080", handlers.CombinedLoggingHandler(os.Stdout, router.CreateRouter()))
}
