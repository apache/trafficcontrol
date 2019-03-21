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

import "crypto/tls";
import "encoding/json";
import "log";
import "net/http";
import "os";
import "strconv";
import "time";

import "github.com/apache/trafficcontrol/lib/go-tc";
import "github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/auth";
import "github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/tocookie";

const API_MIN_MINOR_VERSION = 1;
const API_MAX_MINOR_VERSION = 5;

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("{\"ping\":\"pong\"}\n"));
}

func login(w http.ResponseWriter, r *http.Request) {
	form := auth.PasswordForm{};
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		errBytes, JsonErr := json.Marshal(tc.CreateErrorAlerts(err));
		if JsonErr != nil {
			w.WriteHeader(http.StatusInternalServerError);
			log.Fatalf("Failed to create an alerts structure from '%v': %s\n", err, JsonErr.Error());
		}
		w.WriteHeader(http.StatusBadRequest);
		w.Write(errBytes);
		return;
	}

	expiry := time.Now().Add(time.Hour * 6);
	cookie := tocookie.New(form.Username, expiry, "foo");
	httpCookie := http.Cookie{Name: "mojolicious", Value: cookie, Path: "/", Expires: expiry, HttpOnly: true};
	http.SetCookie(w, &httpCookie);
	resp := struct {
		tc.Alerts
	}{tc.CreateAlerts(tc.SuccessLevel, "Successfully logged in.")};
	respBts, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError);
		log.Fatalf("Couldn't marshal /login response: %s", err.Error());
	}
	w.Write(respBts);
}

func main() {
	server := &http.Server{
		Addr:              ":443",
		TLSConfig:         &tls.Config{InsecureSkipVerify: false},
	};
	log.Println("Server listening on port 443");

	for i := API_MIN_MINOR_VERSION; i <= API_MAX_MINOR_VERSION; i+=1 {
		v := "1." + strconv.Itoa(i);
		http.HandleFunc("/api/"+v+"/ping", ping);
		http.HandleFunc("/api/"+v+"/user/login", login);
	}

	if err := server.ListenAndServeTLS("./localhost.crt", "./localhost.key"); err != nil {
		log.Fatalf("Server crashed: %v\n", err);
	}
	os.Exit(0);
}
