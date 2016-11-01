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

func TestProfile(t *testing.T) {
	resp := fixtures.Profiles()
	server := testHelper.ValidHTTPServer(resp)
	defer server.Close()

	var httpClient http.Client
	to := client.Session{
		URL:       server.URL,
		UserAgent: &httpClient,
	}

	testHelper.Context(t, "Given the need to test a successful Traffic Ops request for Profiles")

	profiles, err := to.Profiles()
	if err != nil {
		testHelper.Error(t, "Should be able to make a request to Traffic Ops")
	} else {
		testHelper.Success(t, "Should be able to make a request to Traffic Ops")
	}

	if len(profiles) != 1 {
		testHelper.Error(t, "Should get back \"1\" Profile, got: %d", len(profiles))
	} else {
		testHelper.Success(t, "Should get back \"1\" Profile")
	}

	for _, p := range profiles {
		if p.Name != "TR_CDN2" {
			testHelper.Error(t, "Should get back \"TR_CDN2\" for \"Name\", got: %s", p.Name)
		} else {
			testHelper.Success(t, "Should get back \"TR_CDN2\" for \"Name\"")
		}

		if p.Description != "kabletown Content Router" {
			testHelper.Error(t, "Should get back \"kabletown Content Router\" for \"Description\", got: %s", p.Description)
		} else {
			testHelper.Success(t, "Should get back \"kabletown Content Router\" for \"Description\"")
		}
	}
}

func TestProfilesUnauthorized(t *testing.T) {
	server := testHelper.InvalidHTTPServer(http.StatusUnauthorized)
	defer server.Close()

	var httpClient http.Client
	to := client.Session{
		URL:       server.URL,
		UserAgent: &httpClient,
	}

	testHelper.Context(t, "Given the need to test a failed Traffic Ops request for Profiles")

	_, err := to.Profiles()
	if err == nil {
		testHelper.Error(t, "Should not be able to make a request to Traffic Ops")
	} else {
		testHelper.Success(t, "Should not be able to make a request to Traffic Ops")
	}
}
