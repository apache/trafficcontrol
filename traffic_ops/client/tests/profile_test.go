/*
   Copyright 2015 Comcast Cable Communications Management, LLC

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jheitz200/traffic_control/traffic_ops/client"
	"github.com/jheitz200/traffic_control/traffic_ops/client/fixtures"
)

func TestProfile(t *testing.T) {
	resp := fixtures.Profiles()
	server := validServer(resp)
	defer server.Close()

	var httpClient http.Client
	to := client.Session{
		URL:       server.URL,
		UserAgent: &httpClient,
	}

	Context(t, "Given the need to test a successful Traffic Ops request for Profiles")

	profiles, err := to.Profiles()
	if err != nil {
		Error(t, "Should be able to make a request to Traffic Ops")
	} else {
		Success(t, "Should be able to make a request to Traffic Ops")
	}

	if len(profiles) != 1 {
		Error(t, "Should get back \"1\" Profile, got: %d", len(profiles))
	} else {
		Success(t, "Should get back \"1\" Profile")
	}

	for _, p := range profiles {
		if p.Name != "TR_CDN2" {
			Error(t, "Should get back \"TR_CDN2\" for \"Name\", got: %s", p.Name)
		} else {
			Success(t, "Should get back \"TR_CDN2\" for \"Name\"")
		}

		if p.Description != "kabletown Content Router" {
			Error(t, "Should get back \"kabletown Content Router\" for \"Description\", got: %s", p.Description)
		} else {
			Success(t, "Should get back \"kabletown Content Router\" for \"Description\"")
		}
	}
}

func TestProfilesUnauthorized(t *testing.T) {
	server := invalidProfilesServer(http.StatusUnauthorized)
	defer server.Close()

	var httpClient http.Client
	to := client.Session{
		URL:       server.URL,
		UserAgent: &httpClient,
	}

	Context(t, "Given the need to test a failed Traffic Ops request for Profiles")

	_, err := to.Profiles()
	if err == nil {
		Error(t, "Should not be able to make a request to Traffic Ops")
	} else {
		Success(t, "Should not be able to make a request to Traffic Ops")
	}
}

func invalidProfilesServer(statusCode int) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		w.Header().Set("Content-Type", "application/json")
	}))
	return server
}
