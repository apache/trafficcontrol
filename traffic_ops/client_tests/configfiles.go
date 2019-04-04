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
import "net/http"

var EDGE_PROFILE_HEADER = "# DO NOT EDIT - Generated for " + EDGE_PROFILE_NAME + " by Traffic Ops (https://localhost:443/) on " + CURRENT_TIME.UTC().Format("Mon Jan 2 15:04:05 MST 2006")
var EDGE_SERVER_HEADER = "# DO NOT EDIT - Generated for " + EDGE_SERVER_HOSTNAME + " by Traffic Ops (https://localhost:443/) on " + CURRENT_TIME.Format("Mon Jan 2 15:04:05 MST 2006")
var MID_PROFILE_HEADER = "# DO NOT EDIT - Generated for " + MID_PROFILE_NAME + " by Traffic Ops (https://localhost:443/) on " + CURRENT_TIME.Format("Mon Jan 2 15:04:05 MST 2006")
var MID_SERVER_HEADER = "# DO NOT EDIT - Generated for " + MID_SERVER_HOSTNAME + " by Traffic Ops (https://localhost:443/) on " + CURRENT_TIME.Format("Mon Jan 2 15:04:05 MST 2006")

var ASTATS_CONFIG = `
allow_ip=0.0.0.0/0
allow_ip6=::/0
path=_astats
record_types=122
`

func edgeAstatsConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "Traffic Ops/" + VERSION + " (Mock)")
	if (r.Method == http.MethodGet) {
		w.Write([]byte(EDGE_PROFILE_HEADER + ASTATS_CONFIG))
	} else {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func midAstatsConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "Traffic Ops/" + VERSION + " (Mock)")
	if (r.Method == http.MethodGet) {
		w.Write([]byte(MID_PROFILE_HEADER + ASTATS_CONFIG))
	} else {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func edgeCacheConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "Traffic Ops/" + VERSION + " (Mock)")
	if (r.Method == http.MethodGet) {
		w.Write([]byte(EDGE_PROFILE_HEADER))
	} else {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func midCacheConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "Traffic Ops/" + VERSION + " (Mock)")
	if (r.Method == http.MethodGet) {
		w.Write([]byte(MID_PROFILE_HEADER))
	} else {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func chkConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "Traffic Ops/" + VERSION + " (Mock)")
	if (r.Method == http.MethodGet) {
		w.Write([]byte(`[{"value":"0:off\t1:off\t2:on\t3:on\t4:on\t5:on\t6:off","name":"trafficserver"}]`))
	} else {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusNotImplemented)
	}
}

var HOSTING_CONFIG = `
hostname=*   volume=1
`

func edgeHostingConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "Traffic Ops/" + VERSION + " (Mock)")
	if (r.Method == http.MethodGet) {
		w.Write([]byte(EDGE_SERVER_HEADER + HOSTING_CONFIG))
	} else {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func midHostingConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "Traffic Ops/" + VERSION + " (Mock)")
	if (r.Method == http.MethodGet) {
		w.Write([]byte(MID_SERVER_HEADER + HOSTING_CONFIG))
	} else {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func edgeIp_allowConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "Traffic Ops/" + VERSION + " (Mock)")
}
