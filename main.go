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

	"fmt"
	"github.com/gorilla/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	db.InitializeDatabase(os.Args[1], os.Args[2], os.Args[3], os.Args[4])

	var Logger = log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)

	Logger.Printf("Starting server...")
	err := http.ListenAndServe(":8080", handlers.CombinedLoggingHandler(os.Stdout, routes.CreateRouter()))

	// for https. Make sure you have the server.pem and server.key file. To gen self signed:
	// openssl genrsa -out server.key 2048
	// openssl req -new -x509 -key server.key -out server.pem -days 3650
	// err := http.ListenAndServeTLS(":1443", "server.pem", "server.key", handlers.CombinedLoggingHandler(os.Stdout, routes.CreateRouter()))
	if err != nil {
		fmt.Println(err.Error())
	}
}
