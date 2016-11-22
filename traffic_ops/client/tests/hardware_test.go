/*

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

func TestHardware(t *testing.T) {
	resp := fixtures.Hardware()
	server := testHelper.ValidHTTPServer(resp)
	defer server.Close()

	var httpClient http.Client
	to := client.Session{
		URL:       server.URL,
		UserAgent: &httpClient,
	}

	testHelper.Context(t, "Given the need to test a successful Traffic Ops request for Hardware")

	hardware, err := to.Hardware(0)
	if err != nil {
		testHelper.Error(t, "Should be able to make a request to Traffic Ops")
	} else {
		testHelper.Success(t, "Should be able to make a request to Traffic Ops")
	}

	if len(hardware) != 1 {
		testHelper.Error(t, "Should get back \"1\" Hardware, got: %d", len(hardware))
	} else {
		testHelper.Success(t, "Should get back \"1\" Hardware")
	}

	for _, h := range hardware {
		if h.HostName != "odol-atsmid-cen-09" {
			testHelper.Error(t, "Should get back \"odol-atsmid-cen-09\" for \"Hostname\", got: %s", h.HostName)
		} else {
			testHelper.Success(t, "Should get back \"odol-atsmid-cen-09\" for \"Hostname\"")
		}

		if h.Value != "1.00" {
			testHelper.Error(t, "Should get back \"1.00\" for \"Value\", got: %s", h.Value)
		} else {
			testHelper.Success(t, "Should get back \"1.00\" for \"Value\"")
		}

		if h.Description != "BACKPLANE FIRMWARE" {
			testHelper.Error(t, "Should get back \"BACKPLANE FIRMWARE\" for \"Description\", got: %s", h.Description)
		} else {
			testHelper.Success(t, "Should get back \"BACKPLANE FIRMWARE\" for \"Description\"")
		}
	}
}

func TestHardwareUnauthorized(t *testing.T) {
	server := testHelper.InvalidHTTPServer(http.StatusUnauthorized)
	defer server.Close()

	var httpClient http.Client
	to := client.Session{
		URL:       server.URL,
		UserAgent: &httpClient,
	}

	testHelper.Context(t, "Given the need to test a failed Traffic Ops request for Hardware")

	_, err := to.Hardware(0)
	if err == nil {
		testHelper.Error(t, "Should not be able to make a request to Traffic Ops")
	} else {
		testHelper.Success(t, "Should not be able to make a request to Traffic Ops")
	}
}
