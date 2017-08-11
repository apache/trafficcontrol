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

func TestParameters(t *testing.T) {
	resp := fixtures.Parameters()
	server := testHelper.ValidHTTPServer(resp)
	defer server.Close()

	to := client.NewSession("", "", server.URL, "", &http.Client{}, false)

	testHelper.Context(t, "Given the need to test a successful Traffic Ops request for Parameters")

	parameters, err := to.Parameters("test")
	if err != nil {
		testHelper.Error(t, "Should be able to make a request to Traffic Ops")
	} else {
		testHelper.Success(t, "Should be able to make a request to Traffic Ops")
	}

	if len(parameters) != 1 {
		testHelper.Error(t, "Should get back \"1\" Parameter, got: %d", len(parameters))
	} else {
		testHelper.Success(t, "Should get back \"1\" Parameter")
	}

	for _, param := range parameters {
		if param.Name != "location" {
			testHelper.Error(t, "Should get back \"location\" for \"Name\", got: %s", param.Name)
		} else {
			testHelper.Success(t, "Should get back \"location\" for \"Name\"")
		}

		if param.Value != "/foo/trafficserver/" {
			testHelper.Error(t, "Should get back \"/foo/trafficserver/\" for \"Value\", got: %s", param.Value)
		} else {
			testHelper.Success(t, "Should get back \"/foo/trafficserver/\" for \"Value\"")
		}

		if param.ConfigFile != "parent.config" {
			testHelper.Error(t, "Should get back \"parent.config\" for \"ConfigFile\", got: %s", param.ConfigFile)
		} else {
			testHelper.Success(t, "Should get back \"parent.config\" for \"ConfigFile\"")
		}
	}
}

func TestParametersUnauthorized(t *testing.T) {
	server := testHelper.InvalidHTTPServer(http.StatusUnauthorized)
	defer server.Close()

	to := client.NewSession("", "", server.URL, "", &http.Client{}, false)

	testHelper.Context(t, "Given the need to test a failed Traffic Ops request for Parameters")

	_, err := to.Parameters("test")
	if err == nil {
		testHelper.Error(t, "Should not be able to make a request to Traffic Ops")
	} else {
		testHelper.Success(t, "Should not be able to make a request to Traffic Ops")
	}
}
