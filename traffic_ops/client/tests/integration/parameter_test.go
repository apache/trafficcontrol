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

package integration

import (
	"encoding/json"
	"fmt"
	"testing"

	traffic_ops "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
)

func TestParameters(t *testing.T) {
	profile, err := GetProfile()
	if err != nil {
		t.Errorf("Could not get a profile, error was: %v\n", err)
	}

	uri := fmt.Sprintf("/api/1.2/parameters/profile/%s.json", profile.Name)
	resp, err := Request(*to, "GET", uri, nil)
	if err != nil {
		t.Errorf("Could not get %s reponse was: %v\n", uri, err)
		t.FailNow()
	}

	defer resp.Body.Close()
	var apiParamRes traffic_ops.ParamResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiParamRes); err != nil {
		t.Errorf("Could not decode parameter json.  Error is: %v\n", err)
		t.FailNow()
	}
	apiParams := apiParamRes.Response

	clientParams, err := to.Parameters(profile.Name)
	if err != nil {
		t.Errorf("Could not get parameters from client.  Error is: %v\n", err)
		t.FailNow()
	}

	if len(apiParams) != len(clientParams) {
		t.Errorf("Params Response Length -- expected %v, got %v\n", len(apiParams), len(clientParams))
	}

	for _, apiParam := range apiParams {
		match := false
		for _, clientParam := range clientParams {
			if apiParam.Name == clientParam.Name && apiParam.Value == clientParam.Value && apiParam.ConfigFile == clientParam.ConfigFile {
				match = true
			}
		}
		if !match {
			t.Errorf("Did not get a param matching %+v from the client\n", apiParam)
		}
	}
}
