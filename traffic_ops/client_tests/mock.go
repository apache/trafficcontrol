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

// TODO: read this in properly from the VERSION file
const VERSION = "3.0.0";

// Will be used for various `lastUpdated` fields
var CURRENT_TIME *tc.TimeNoMod = tc.NewTimeNoMod();

// Static CDN Fields
var ALL_CDN = "ALL";
var ALL_CDN_ID = 1;
var ALL_CDN_DOMAINNAME = "-";
var ALL_CDN_DNSSEC_ENABLED = false;
var CDN = "Mock-CDN";
var CDN_ID = 2;
var CDN_DOMAIN_NAME = "mock.cdn.test";
var CDN_DNSSEC_ENABLED = false;

var CDNS = []tc.CDNNullable{
	tc.CDNNullable {
		DNSSECEnabled: &ALL_CDN_DNSSEC_ENABLED,
		DomainName: &ALL_CDN_DOMAINNAME,
		ID: &ALL_CDN_ID,
		LastUpdated: CURRENT_TIME,
		Name: &ALL_CDN,
	},
	tc.CDNNullable {
		DNSSECEnabled: &CDN_DNSSEC_ENABLED,
		DomainName: &CDN_DOMAIN_NAME,
		ID: &CDN_ID,
		Name: &CDN,
	},
};

// Static Cache Group fields
var EDGE_CACHEGROUP_ID = 1;
var EDGE_CACHEGROUP = "Edge";
var EDGE_CACHEGROUP_SHORT_NAME = "Edge";
var EDGE_CACHEGROUP_LATITUDE = 0.0;
var EDGE_CACHEGROUP_LONGITUDE = 0.0;
var EDGE_CACHEGROUP_PARENT_NAME= "Mid"; // NOTE: This places a hard requirement on the `cachegroups` implementation - must have a `MID_LOC` Cache Group named "Mid"
var EDGE_CACHEGROUP_PARENT_ID = 2; // NOTE: This places a hard requirement on the `cachegroups` implementation - must have a `MID_LOC` Cache Group identified by `2`
var EDGE_CACHEGROUP_FALLBACK_TO_CLOSEST = true;
var EDGE_CACHEGROUP_LOCALIZATION_METHODS = []tc.LocalizationMethod{
	tc.LocalizationMethodCZ,
	tc.LocalizationMethodDeepCZ,
	tc.LocalizationMethodGeo,
};
var EDGE_CACHEGROUP_TYPE = "EDGE_LOC";
var EDGE_CACHEGROUP_TYPE_ID = 1; // NOTE: This places a hard requirement on the `types` implementation - must have `EDGE_LOC` == 1

var MID_CACHEGROUP_ID = 2;
var MID_CACHEGROUP = "Mid";
var MID_CACHEGROUP_SHORT_NAME = "Mid";
var MID_CACHEGROUP_LATITUDE = 0.0;
var MID_CACHEGROUP_LONGITUDE = 0.0;
var MID_CACHEGROUP_FALLBACK_TO_CLOSEST = true;
var MID_CACHEGROUP_LOCALIZATION_METHODS = []tc.LocalizationMethod{
	tc.LocalizationMethodCZ,
	tc.LocalizationMethodDeepCZ,
	tc.LocalizationMethodGeo,
};
var MID_CACHEGROUP_TYPE = "MID_LOC";
var MID_CACHEGROUP_TYPE_ID = 2; // NOTE: This places a hard requirement on the `types` implementation - must have `MID_LOC` == 2


var CACHEGROUPS = []tc.CacheGroupNullable{
	tc.CacheGroupNullable{
		ID: &EDGE_CACHEGROUP_ID,
		Name: &EDGE_CACHEGROUP,
		ShortName: &EDGE_CACHEGROUP_SHORT_NAME,
		Latitude: &EDGE_CACHEGROUP_LATITUDE,
		Longitude: &EDGE_CACHEGROUP_LONGITUDE,
		ParentName: &EDGE_CACHEGROUP_PARENT_NAME,
		ParentCachegroupID: &EDGE_CACHEGROUP_PARENT_ID,
		SecondaryParentName: nil,
		SecondaryParentCachegroupID: nil,
		FallbackToClosest: &EDGE_CACHEGROUP_FALLBACK_TO_CLOSEST,
		LocalizationMethods: &EDGE_CACHEGROUP_LOCALIZATION_METHODS,
		Type: &EDGE_CACHEGROUP_TYPE,
		TypeID: &EDGE_CACHEGROUP_TYPE_ID,
		LastUpdated: CURRENT_TIME,
		Fallbacks: nil,
	},
	tc.CacheGroupNullable{
		ID: &MID_CACHEGROUP_ID,
		Name: &MID_CACHEGROUP,
		ShortName: &MID_CACHEGROUP_SHORT_NAME,
		Latitude: &MID_CACHEGROUP_LATITUDE,
		Longitude: &MID_CACHEGROUP_LONGITUDE,
		ParentName: nil,
		ParentCachegroupID: nil,
		SecondaryParentName: nil,
		SecondaryParentCachegroupID: nil,
		FallbackToClosest: &MID_CACHEGROUP_FALLBACK_TO_CLOSEST,
		LocalizationMethods: &MID_CACHEGROUP_LOCALIZATION_METHODS,
		Type: &MID_CACHEGROUP_TYPE,
		TypeID: &MID_CACHEGROUP_TYPE_ID,
		LastUpdated: CURRENT_TIME,
		Fallbacks: nil,
	},
};


// Static user fields
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
};

var CURRENT_USER = tc.UserCurrent{
	UserName: &USERNAME,
	LocalUser: &LOCAL_USER,
	RoleName: &ROLE,
	CommonUserFields: COMMON_USER_FIELDS,
};

// Just sets some common headers
func common(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json");
	w.Header().Set("Server", "Traffic Ops/"+VERSION+" (Mock)");
}

func ping(w http.ResponseWriter, r *http.Request) {
	common(w);
	w.Write([]byte("{\"ping\":\"pong\"}\n"));
}

func login(w http.ResponseWriter, r *http.Request) {
	form := auth.PasswordForm{};
	common(w);
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
	common(w);
	api.WriteResp(w, r, CURRENT_USER);
}

func getCDNS(w http.ResponseWriter, r *http.Request) {
	common(w);
	api.WriteResp(w, r, CDNS);
}

func getCacheGroups(w http.ResponseWriter, r *http.Request) {
	common(w);
	api.WriteResp(w, r, CACHEGROUPS);
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
		http.HandleFunc("/api/"+v+"/cdns", getCDNS);
		http.HandleFunc("/api/"+v+"/cachegroups", getCacheGroups);
	}

	log.Printf("Finished loading API routes at %s, server listening on port 443", CURRENT_TIME.String());

	if err := server.ListenAndServeTLS("./localhost.crt", "./localhost.key"); err != nil {
		log.Fatalf("Server crashed: %v\n", err);
	}
	os.Exit(0);
}
