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
	"github.com/jheitz200/test_helper"
)

func TestLogin(t *testing.T) {
	resp := client.Result{
		Alerts: []client.Alert{
			client.Alert{
				Level: "success",
				Text:  "Successfully logged in.",
			},
		},
	}

	server := testHelper.ValidHTTPServer(resp)

	testHelper.Context(t, "Given the need to test a successful login to Traffic Ops")

	session, err := client.Login(server.URL, "test", "password", true)
	if err != nil {
		testHelper.Error(t, "Should be able to login")
	} else {
		testHelper.Success(t, "Should be able to login")
	}

	if session.UserName != "test" {
		testHelper.Error(t, "Should get back \"test\" for \"UserName\", got %s", session.UserName)
	} else {
		testHelper.Success(t, "Should get back \"test\" for \"UserName\"")
	}

	if session.Password != "password" {
		testHelper.Error(t, "Should get back \"password\" for \"Password\", got %s", session.Password)
	} else {
		testHelper.Success(t, "Should get back \"password\" for \"Password\"")
	}

	if session.URL != server.URL {
		testHelper.Error(t, "Should get back \"%s\" for \"URL\", got %s", server.URL, session.URL)
	} else {
		testHelper.Success(t, "Should get back \"%s\" for \"URL\"", server.URL)
	}
}

func TestLoginUnauthorized(t *testing.T) {
	server := testHelper.InvalidHTTPServer(http.StatusUnauthorized)
	defer server.Close()

	testHelper.Context(t, "Given the need to test an unsuccessful login to Traffic Ops")

	_, err := client.Login(server.URL, "test", "password", true)
	if err == nil {
		testHelper.Error(t, "Should not be able to login")
	} else {
		testHelper.Success(t, "Should not be able to login")
	}
}
