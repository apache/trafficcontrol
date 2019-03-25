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
import "github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api";
import "github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/auth";
import "github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/tocookie";

const API_MIN_MINOR_VERSION = 1;
const API_MAX_MINOR_VERSION = 5;

// Will be used for various `lastUpdated` fields
var CURRENT_TIME *tc.TimeNoMod = tc.NewTimeNoMod();

// Current user fields
// (These _should_ be `const`, but you can't take the address of a `const` (for some reason))
var USERNAME = "admin";
var LOCAL_USER = true;
var USER_ID = 1;
var TENANT = "root";
var TENANT_ID = 1; // NOTE: This places a hard requirement on `tenant` implementation - `root` == `1`
var ROLE = "admin";
var ROLE_ID = 1; // NOTE: This places a hard requirement on `roles` implementation - `admin` == `1`
var NEW_USER = false;

var COMMON_USER_FIELDS = tc.CommonUserFields {
	AddressLine1: nil,
	AddressLine2: nil,
	City: nil,
	Company: nil,
	Country: nil,
	Email: nil,
	FullName: nil,
	GID: nil,
	ID: &USER_ID,
	NewUser: &NEW_USER,
	PhoneNumber: nil,
	PostalCode: nil,
	PublicSSHKey: nil,
	Role: &ROLE_ID,
	StateOrProvince: nil,
	Tenant: &TENANT,
	TenantID: &TENANT_ID,
	UID: nil,
	LastUpdated: CURRENT_TIME,
}

var CURRENT_USER = tc.UserCurrent{
	UserName: &USERNAME,
	LocalUser: &LOCAL_USER,
	RoleName: &ROLE,
	CommonUserFields: COMMON_USER_FIELDS,
}

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

func cuser(w http.ResponseWriter, r *http.Request) {
	api.WriteResp(w, r, CURRENT_USER);
}

func main() {
	server := &http.Server{
		Addr:              ":443",
		TLSConfig:         &tls.Config{InsecureSkipVerify: false},
	};

	for i := API_MIN_MINOR_VERSION; i <= API_MAX_MINOR_VERSION; i+=1 {
		v := "1." + strconv.Itoa(i);
		log.Printf("Loading API v%s\n", v);
		http.HandleFunc("/api/"+v+"/ping", ping);
		http.HandleFunc("/api/"+v+"/user/login", login);
		http.HandleFunc("/api/"+v+"/user/current", cuser);
	}

	log.Printf("Finished loading API routes at %s, server listening on port 443", CURRENT_TIME.String());

	if err := server.ListenAndServeTLS("./localhost.crt", "./localhost.key"); err != nil {
		log.Fatalf("Server crashed: %v\n", err);
	}
	os.Exit(0);
}
