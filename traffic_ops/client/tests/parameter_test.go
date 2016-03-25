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

func TestParameters(t *testing.T) {
	resp := fixtures.Parameters()
	server := validServer(resp)
	defer server.Close()

	var httpClient http.Client
	to := client.Session{
		URL:       server.URL,
		UserAgent: &httpClient,
	}

	Context(t, "Given the need to test a successful Traffic Ops request for Parameters")

	parameters, err := to.Parameters("test")
	if err != nil {
		Error(t, "Should be able to make a request to Traffic Ops")
	} else {
		Success(t, "Should be able to make a request to Traffic Ops")
	}

	if len(parameters) != 1 {
		Error(t, "Should get back \"1\" Parameter, got: %d", len(parameters))
	} else {
		Success(t, "Should get back \"1\" Parameter")
	}

	for _, param := range parameters {
		if param.Name != "location" {
			Error(t, "Should get back \"location\" for \"Name\", got: %s", param.Name)
		} else {
			Success(t, "Should get back \"location\" for \"Name\"")
		}

		if param.Value != "/foo/trafficserver/" {
			Error(t, "Should get back \"/foo/trafficserver/\" for \"Value\", got: %s", param.Value)
		} else {
			Success(t, "Should get back \"/foo/trafficserver/\" for \"Value\"")
		}

		if param.ConfigFile != "parent.config" {
			Error(t, "Should get back \"parent.config\" for \"ConfigFile\", got: %s", param.ConfigFile)
		} else {
			Success(t, "Should get back \"parent.config\" for \"ConfigFile\"")
		}
	}
}

func TestParametersUnauthorized(t *testing.T) {
	server := invalidParametersServer(http.StatusUnauthorized)
	defer server.Close()

	var httpClient http.Client
	to := client.Session{
		URL:       server.URL,
		UserAgent: &httpClient,
	}

	Context(t, "Given the need to test a failed Traffic Ops request for Parameters")

	_, err := to.Parameters("test")
	if err == nil {
		Error(t, "Should not be able to make a request to Traffic Ops")
	} else {
		Success(t, "Should not be able to make a request to Traffic Ops")
	}
}

func invalidParametersServer(statusCode int) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		w.Header().Set("Content-Type", "application/json")
	}))
	return server
}
