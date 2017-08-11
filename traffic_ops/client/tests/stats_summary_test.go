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

func TestStatsSummary(t *testing.T) {
	resp := fixtures.StatsSummary()
	server := testHelper.ValidHTTPServer(resp)
	defer server.Close()

	to := client.NewSession("", "", server.URL, "", &http.Client{}, false)

	testHelper.Context(t, "Given the need to test a successful Traffic Ops request for Stats Summary")

	stats, err := to.SummaryStats("test-cdn", "test-ds1", "test-stat")
	if err != nil {
		testHelper.Error(t, "Should be able to make a request to Traffic Ops")
	} else {
		testHelper.Success(t, "Should be able to make a request to Traffic Ops")
	}

	if len(stats) != 1 {
		testHelper.Error(t, "Should get back \"1\" Parameter, got: %d", len(stats))
	} else {
		testHelper.Success(t, "Should get back \"1\" Parameter")
	}

	for _, s := range stats {
		if s.StatName != "test-stat" {
			testHelper.Error(t, "Should get back \"test-stat\" for \"StatName\", got: %s", s.StatName)
		} else {
			testHelper.Success(t, "Should get back \"test-stat\" for \"StatName\"")
		}

		if s.DeliveryService != "test-ds1" {
			testHelper.Error(t, "Should get back \"test-ds1\" for \"DeliveryService\", got: %s", s.DeliveryService)
		} else {
			testHelper.Success(t, "Should get back \"test-ds1\" for \"DeliveryService\"")
		}

		if s.StatValue != "3.1415" {
			testHelper.Error(t, "Should get back \"3.1415\" for \"StatValue\", got: %s", s.StatValue)
		} else {
			testHelper.Success(t, "Should get back \"3.1415\" for \"StatValue\"")
		}

		if s.CDNName != "test-cdn" {
			testHelper.Error(t, "Should get back \"test-cdn\" for \"CDNName\", got: %s", s.CDNName)
		} else {
			testHelper.Success(t, "Should get back \"test-cdn\" for \"CDNName\"")
		}
	}
}

func TestStatsSummaryUnauthorized(t *testing.T) {
	server := testHelper.InvalidHTTPServer(http.StatusUnauthorized)
	defer server.Close()

	to := client.NewSession("", "", server.URL, "", &http.Client{}, false)

	testHelper.Context(t, "Given the need to test a failed Traffic Ops request for Stats Summary")

	_, err := to.SummaryStats("test-cdn", "test-ds1", "test-stat")
	if err == nil {
		testHelper.Error(t, "Should not be able to make a request to Traffic Ops")
	} else {
		testHelper.Success(t, "Should not be able to make a request to Traffic Ops")
	}
}
