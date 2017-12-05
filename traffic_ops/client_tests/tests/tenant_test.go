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
	"github.com/apache/incubator-trafficcontrol/traffic_ops/client_tests/fixtures"
	"github.com/jheitz200/test_helper"
)

func TestTenants(t *testing.T) {
	resp := fixtures.Tenants()
	server := testHelper.ValidHTTPServer(resp)
	defer server.Close()

	to := client.NewSession("", "", server.URL, "", &http.Client{}, false)

	testHelper.Context(t, "Given the need to test a successful Traffic Ops request for Tenants")

	t, err := to.Tenants()
	if err != nil {
		testHelper.Error(t, "Should be able to make a request to Traffic Ops")
	} else {
		testHelper.Success(t, "Should be able to make a request to Traffic Ops")
	}

	if len(t) != 1 {
		testHelper.Error(t, "Should get back \"1\" Tenant, got: %d", len(t))
	} else {
		testHelper.Success(t, "Should get back \"1\" Tenant")
	}

	for _, s := range t {
		if s.ID != 001 {
			testHelper.Error(t, "Should get back \"1\" for \"ID\", got: %s", s.ID)
		} else {
			testHelper.Success(t, "Should get back \"1\" for \"ID\"")
		}
	}
}

func TestTenantUnauthorized(t *testing.T) {
	server := testHelper.InvalidHTTPServer(http.StatusUnauthorized)
	defer server.Close()

	to := client.NewSession("", "", server.URL, "", &http.Client{}, false)

	testHelper.Context(t, "Given the need to test a failed Traffic Ops request for Tenants")

	_, err := to.Tenants()
	if err == nil {
		testHelper.Error(t, "Should not be able to make a request to Traffic Ops")
	} else {
		testHelper.Success(t, "Should not be able to make a request to Traffic Ops")
	}
}
