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
	"testing"

	"github.com/apache/incubator-trafficcontrol/traffic_ops/client"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/client/fixtures"
	"github.com/jheitz200/test_helper"
)

func TestCacheGroup(t *testing.T) {
	resp := fixtures.Cachegroups()
	server := testHelper.ValidHTTPServer(resp)
	defer server.Close()

	var httpClient http.Client
	to := client.Session{
		URL:       server.URL,
		UserAgent: &httpClient,
	}

	testHelper.Context(t, "Given the need to test a successful Traffic Ops request for CacheGroups")

	cacheGroups, err := to.CacheGroups()
	if err != nil {
		testHelper.Error(t, "Should be able to make a request to Traffic Ops")
	} else {
		testHelper.Success(t, "Should be able to make a request to Traffic Ops")
	}

	if len(cacheGroups) != 1 {
		testHelper.Error(t, "Should get back \"1\" CacheGroups, got: %d", len(cacheGroups))
	} else {
		testHelper.Success(t, "Should get back \"1\" CacheGroups")
	}

	for _, cacheGroup := range cacheGroups {
		if cacheGroup.Name == "" {
			testHelper.Error(t, "Should get back \"edge-philadelphia\" for \"Name\", got: %s", cacheGroup.Name)
		} else {
			testHelper.Success(t, "Should get back \"edge-philadelphia\" for \"Name\"")
		}

		if cacheGroup.Longitude != 5 {
			testHelper.Error(t, "Should get back \"5\" for \"Longitude\", got: %v", cacheGroup.Longitude)
		} else {
			testHelper.Success(t, "Should get back \"5\" for \"Longitude\"")
		}

		if cacheGroup.Latitude != 55 {
			testHelper.Error(t, "Should get back \"55\" for \"Latitude\", got: %v", cacheGroup.Latitude)
		} else {
			testHelper.Success(t, "Should get back \"55\" for \"Latitude\"")
		}

		if cacheGroup.ParentName != "mid-northeast" {
			testHelper.Error(t, "Should get back \"mid-northeast\" for \"ParentName\", got: %s", cacheGroup.ParentName)
		} else {
			testHelper.Success(t, "Should get back \"mid-northeast\" for \"ParentName\"")
		}
	}
}

func TestCacheGroupsUnauthorized(t *testing.T) {
	server := testHelper.InvalidHTTPServer(http.StatusUnauthorized)
	defer server.Close()

	var httpClient http.Client
	to := client.Session{
		URL:       server.URL,
		UserAgent: &httpClient,
	}

	testHelper.Context(t, "Given the need to test a failed Traffic Ops request for CacheGroups")

	_, err := to.CacheGroups()
	if err == nil {
		testHelper.Error(t, "Should not be able to make a request to Traffic Ops")
	} else {
		testHelper.Success(t, "Should not be able to make a request to Traffic Ops")
	}
}
