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

	"github.com/Comcast/test_helper"
	"github.com/jheitz200/traffic_control/traffic_ops/client"
	"github.com/jheitz200/traffic_control/traffic_ops/client/fixtures"
)

func TestDeliveryServices(t *testing.T) {
	resp := fixtures.DeliveryServices()
	server := test.ValidHTTPServer(resp)
	defer server.Close()

	var httpClient http.Client
	to := client.Session{
		URL:       server.URL,
		UserAgent: &httpClient,
	}

	test.Context(t, "Given the need to test a successful Traffic Ops request for DeliveryServices")

	ds, err := to.DeliveryServices()
	if err != nil {
		test.Error(t, "Should be able to make a request to Traffic Ops")
	} else {
		test.Success(t, "Should be able to make a request to Traffic Ops")
	}

	if len(ds) != 1 {
		test.Error(t, "Should get back \"1\" DeliveryService, got: %d", len(ds))
	} else {
		test.Success(t, "Should get back \"1\" DeliveryService")
	}

	for _, s := range ds {
		if s.XMLID != "ds-test" {
			test.Error(t, "Should get back \"ds-test\" for \"XMLID\", got: %s", s.XMLID)
		} else {
			test.Success(t, "Should get back \"ds-test\" for \"XMLID\"")
		}

		if s.MissLong != "-99.123456" {
			test.Error(t, "Should get back \"-99.123456\" for \"MissLong\", got: %s", s.MissLong)
		} else {
			test.Success(t, "Should get back \"-99.123456\" for \"MissLong\"")
		}
	}
}

func TestDeliveryServicesUnauthorized(t *testing.T) {
	server := test.InvalidHTTPServer(http.StatusUnauthorized)
	defer server.Close()

	var httpClient http.Client
	to := client.Session{
		URL:       server.URL,
		UserAgent: &httpClient,
	}

	test.Context(t, "Given the need to test a failed Traffic Ops request for DeliveryServices")

	_, err := to.DeliveryServices()
	if err == nil {
		test.Error(t, "Should not be able to make a request to Traffic Ops")
	} else {
		test.Success(t, "Should not be able to make a request to Traffic Ops")
	}
}
