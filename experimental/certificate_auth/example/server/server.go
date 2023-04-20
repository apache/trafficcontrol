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
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
)

func main() {
	handler := http.NewServeMux()
	handler.HandleFunc("/", HelloHandler)

	tlsConfig := &tls.Config{
		ClientAuth: tls.RequestClientCert,
	}

	server := http.Server{
		Addr:      "server.local:8443",
		Handler:   handler,
		TLSConfig: tlsConfig,
	}

	if err := server.ListenAndServeTLS("../certs/server.crt.pem", "../certs/server.key.pem"); err != nil {
		log.Fatalf("error listening to port: %v", err)
	}
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {

	if r.TLS.PeerCertificates != nil {
		clientCert := r.TLS.PeerCertificates[0]
		fmt.Println("Client cert subject: ", clientCert.Subject)
	}

	fmt.Println("Hello")
}
