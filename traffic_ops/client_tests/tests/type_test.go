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

func TestTypes(t *testing.T) {
	resp := fixtures.Types()
	server := testHelper.ValidHTTPServer(resp)
	defer server.Close()

	to := client.NewSession("", "", server.URL, "", &http.Client{}, false)

	testHelper.Context(t, "Given the need to test a successful Traffic Ops request for Types")

	types, err := to.Types()
	if err != nil {
		testHelper.Error(t, "Should be able to make a request to Traffic Ops")
	} else {
		testHelper.Success(t, "Should be able to make a request to Traffic Ops")
	}

	for _, n := range types {
		if n.ID != 1 {
			testHelper.Error(t, "Should get back 1 for \"ID\", got %d", n.ID)
		} else {
			testHelper.Success(t, "Should get back 1 for \"ID\"")
		}

		if n.Name != "EDGE" {
			testHelper.Error(t, "Should get back \"EDGE\" for \"Name\", got %s", n.Name)
		} else {
			testHelper.Success(t, "Should get back \"EDGE\" for \"Name\"")
		}

		if n.Description != "edge cache" {
			testHelper.Error(t, "Should get back \"edge cache\" for \"Description\", got %s", n.Description)
		} else {
			testHelper.Success(t, "Should get back \"edge cache\" for \"Description\"")
		}

		if n.UseInTable != "server" {
			testHelper.Error(t, "Should get back \"server\" for \"UseInTable\", got %s", n.UseInTable)
		} else {
			testHelper.Success(t, "Should get back \"server\" for \"UseInTable\"")
		}
	}
}

func TestTypesUnauthorized(t *testing.T) {
	server := testHelper.InvalidHTTPServer(http.StatusUnauthorized)
	defer server.Close()

	to := client.NewSession("", "", server.URL, "", &http.Client{}, false)

	testHelper.Context(t, "Given the need to test a failed Traffic Ops request for Types")

	_, err := to.Types()
	if err == nil {
		testHelper.Error(t, "Should not be able to make a request to Traffic Ops")
	} else {
		testHelper.Success(t, "Should not be able to make a request to Traffic Ops")
	}
}
