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

func TestCacheGroup(t *testing.T) {
	resp := fixtures.Cachegroups()
	server := validServer(resp)
	defer server.Close()

	var httpClient http.Client
	to := client.Session{
		URL:       server.URL,
		UserAgent: &httpClient,
	}

	Context(t, "Given the need to test a successful Traffic Ops request for CacheGroups")

	cacheGroups, err := to.CacheGroups()
	if err != nil {
		Error(t, "Should be able to make a request to Traffic Ops")
	} else {
		Success(t, "Should be able to make a request to Traffic Ops")
	}

	if len(cacheGroups) != 1 {
		Error(t, "Should get back \"1\" CacheGroups, got: %d", len(cacheGroups))
	} else {
		Success(t, "Should get back \"1\" CacheGroups")
	}

	for _, cacheGroup := range cacheGroups {
		if cacheGroup.Name == "" {
			Error(t, "Should get back \"edge-philadelphia\" for \"Name\", got: %s", cacheGroup.Name)
		} else {
			Success(t, "Should get back \"edge-philadelphia\" for \"Name\"")
		}

		if cacheGroup.Longitude != 5 {
			Error(t, "Should get back \"5\" for \"Longitude\", got: %v", cacheGroup.Longitude)
		} else {
			Success(t, "Should get back \"5\" for \"Longitude\"")
		}

		if cacheGroup.Latitude != 55 {
			Error(t, "Should get back \"55\" for \"Latitude\", got: %v", cacheGroup.Latitude)
		} else {
			Success(t, "Should get back \"55\" for \"Latitude\"")
		}

		if cacheGroup.ParentName != "mid-northeast" {
			Error(t, "Should get back \"mid-northeast\" for \"ParentName\", got: %s", cacheGroup.ParentName)
		} else {
			Success(t, "Should get back \"mid-northeast\" for \"ParentName\"")
		}
	}
}

func TestCacheGroupsUnauthorized(t *testing.T) {
	server := invalidCacheGroupServer(http.StatusUnauthorized)
	defer server.Close()

	var httpClient http.Client
	to := client.Session{
		URL:       server.URL,
		UserAgent: &httpClient,
	}

	Context(t, "Given the need to test a failed Traffic Ops request for CacheGroups")

	_, err := to.CacheGroups()
	if err == nil {
		Error(t, "Should not be able to make a request to Traffic Ops")
	} else {
		Success(t, "Should not be able to make a request to Traffic Ops")
	}
}

func invalidCacheGroupServer(statusCode int) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		w.Header().Set("Content-Type", "application/json")
	}))
	return server
}
