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
	"strings"
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

	if len(hardware) != 2 {
		testHelper.Error(t, "Should get back \"2\" Hardware, got: %d", len(hardware))
	} else {
		testHelper.Success(t, "Should get back \"2\" Hardware")
	}

	for _, h := range hardware {
		if !strings.Contains(h.HostName, "edge-den") {
			testHelper.Error(t, "Should get back \"edge-den\" in \"Hostname\", got: %s", h.HostName)
		} else {
			testHelper.Success(t, "Should get back \"edge-den\" in \"Hostname\"")
		}

		if !strings.Contains(h.Value, "DIS") {
			testHelper.Error(t, "Should get back \"DIS1\" or \"DIS2\" for \"Value\", got: %s", h.Value)
		} else {
			testHelper.Success(t, "Should get back \"DIS1\" or \"DIS2\" for \"Value\"")
		}

		if !strings.Contains(h.Description, "Phys") {
			testHelper.Error(t, "Should get back \"Phys Disk\" or \"Physical Disk\" for \"Description\", got: %s", h.Description)
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
