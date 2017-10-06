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

func TestCDNs(t *testing.T) {
	resp := fixtures.CDNs()
	server := testHelper.ValidHTTPServer(resp)
	defer server.Close()

	to := client.NewSession("", "", server.URL, "", &http.Client{}, false)

	testHelper.Context(t, "Given the need to test a successful Traffic Ops request for CDNs")

	cdns, err := to.CDNs()
	if err != nil {
		testHelper.Error(t, "Should be able to make a request to Traffic Ops")
	} else {
		testHelper.Success(t, "Should be able to make a request to Traffic Ops")
	}

	for _, cdn := range cdns {
		if cdn.ID != 1 {
			testHelper.Error(t, "Should get back 1 for \"ID\", got %d", cdn.ID)
		} else {
			testHelper.Success(t, "Should get back 1 for \"ID\"")
		}
		if cdn.Name != "CDN-1" {
			testHelper.Error(t, "Should get back \"CDN-1\" for \"name\", got %s", cdn.Name)
		} else {
			testHelper.Success(t, "Should get back \"CDN-1\" for \"name\"")
		}

		if cdn.LastUpdated != "2016-03-22 17:00:30" {
			testHelper.Error(t, "Should get back \"2016-03-22 17:00:30\" for \"LastUpdated\", got %s", cdn.LastUpdated)
		} else {
			testHelper.Success(t, "Should get back \"2016-03-22 17:00:30\" for \"LastUpdated\"")
		}
	}
}

func TestCDNsUnauthorized(t *testing.T) {
	server := testHelper.InvalidHTTPServer(http.StatusUnauthorized)
	defer server.Close()

	to := client.NewSession("", "", server.URL, "", &http.Client{}, false)

	testHelper.Context(t, "Given the need to test a failed Traffic Ops request for CDNs")

	_, err := to.CDNs()
	if err == nil {
		testHelper.Error(t, "Should not be able to make a request to Traffic Ops")
	} else {
		testHelper.Success(t, "Should not be able to make a request to Traffic Ops")
	}
}

func TestCDNName(t *testing.T) {
	resp := fixtures.CDNs()
	server := testHelper.ValidHTTPServer(resp)
	defer server.Close()

	to := client.NewSession("", "", server.URL, "", &http.Client{}, false)

	cdn := "CDN-1"
	testHelper.Context(t, "Given the need to test a successful Traffic Ops request for CDN: \"%s\"", cdn)

	cdns, err := to.CDNName(cdn)
	if err != nil {
		testHelper.Error(t, "Should be able to make a request to Traffic Ops")
	} else {
		testHelper.Success(t, "Should be able to make a request to Traffic Ops")
	}

	for _, cdn := range cdns {
		if cdn.Name != "CDN-1" {
			testHelper.Error(t, "Should get back \"CDN-1\" for \"name\", got %s", cdn.Name)
		} else {
			testHelper.Success(t, "Should get back \"CDN-1\" for \"name\"")
		}

		if cdn.LastUpdated != "2016-03-22 17:00:30" {
			testHelper.Error(t, "Should get back \"2016-03-22 17:00:30\" for \"LastUpdated\", got %s", cdn.LastUpdated)
		} else {
			testHelper.Success(t, "Should get back \"2016-03-22 17:00:30\" for \"LastUpdated\"")
		}
	}
}

func TestCDNNameUnauthorized(t *testing.T) {
	server := testHelper.InvalidHTTPServer(http.StatusUnauthorized)
	defer server.Close()

	to := client.NewSession("", "", server.URL, "", &http.Client{}, false)

	cdn := "CDN-1"
	testHelper.Context(t, "Given the need to test a failed Traffic Ops request for CDN: \"%s\"", cdn)

	_, err := to.CDNName(cdn)
	if err == nil {
		testHelper.Error(t, "Should not be able to make a request to Traffic Ops")
	} else {
		testHelper.Success(t, "Should not be able to make a request to Traffic Ops")
	}
}
