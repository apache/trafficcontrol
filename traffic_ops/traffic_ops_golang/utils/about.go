package utils

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
	"encoding/json"
	"net/http"
	"runtime"
	"strings"
)

// about allows access to the version info identified at build time
type about struct {
	CommitHash           string `json:"commitHash,omitempty"`
	Commits              string `json:"commits,omitempty"`
	GoVersion            string `json:"goVersion,omitempty"`
	TrafficOpsName       string `json:"name,omitempty"`
	TrafficOpsRPMVersion string `json:"trafficOpsRPMVersion,omitempty"`
	TrafficOpsVersion    string `json:"trafficOpsVersion,omitempty"`
}

// About contains version info to be exposed by `api/.../about.json` endpoint
var About about

func splitRPMVersion(v string) (string, string, string, string) {

	if v == "" {
		return "UnknownVersion", "", "", ""
	}
	// RPM version is something like traffic_ops-2.3.0-8765.a0b1c3d4
	//  -- if not of that form, Name, Version, Commits, CommitHash may be missing
	s := strings.SplitN(v, "-", 3)
	if len(s) >= 3 {
		// 3rd field is commits.hash
		t := strings.SplitN(s[2], ".", 2)
		s = append(s[0:2], t...)
	}
	for cap(s) < 4 {
		s = append(s, "")
	}
	return s[0], s[1], s[2], s[3]
}

func init() {
	// name, version, commits, hash -- parts of rpm version string
	n, v, c, h := splitRPMVersion(About.TrafficOpsRPMVersion)
	About = about{
		GoVersion:            runtime.Version(),
		TrafficOpsRPMVersion: About.TrafficOpsRPMVersion,
		TrafficOpsName:       n,
		TrafficOpsVersion:    v,
		Commits:              c,
		CommitHash:           h,
	}

}

// AboutHandler returns info about running Traffic Ops
func AboutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(About)
	}
}
