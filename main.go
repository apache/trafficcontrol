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
	routes "./routes"
	db "./todb"

	"github.com/gorilla/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	db.InitializeDatabase(os.Args[1], os.Args[2], os.Args[3], os.Args[4])

	var Logger = log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)

	Logger.Printf("Starting server...")
	http.ListenAndServe(":8080", handlers.CombinedLoggingHandler(os.Stdout, routes.CreateRouter()))
}
