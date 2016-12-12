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

func TestTypes(t *testing.T) {

	uri := fmt.Sprintf("/api/1.2/types.json")
	resp, err := Request(*to, "GET", uri, nil)
	if err != nil {
		t.Errorf("Could not get %s reponse was: %v\n", uri, err)
		t.FailNow()
	}

	defer resp.Body.Close()
	var apiTypeRes traffic_ops.TypeResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiTypeRes); err != nil {
		t.Errorf("Could not decode type json.  Error is: %v\n", err)
		t.FailNow()
	}
	apiTypes := apiTypeRes.Response

	clientTypes, err := to.Types()
	if err != nil {
		t.Errorf("Could not get types from client.  Error is: %v\n", err)
		t.FailNow()
	}

	if len(apiTypes) != len(clientTypes) {
		t.Errorf("Types Response Length -- expected %v, got %v\n", len(apiTypes), len(clientTypes))
	}

	for _, apiType := range apiTypes {
		match := false
		for _, clientType := range clientTypes {
			if apiType.ID == clientType.ID {
				match = true
				if apiType.Description != clientType.Description {
					t.Errorf("Description -- Expected %v, got %v\n", apiType.Description, clientType.Description)
				}
				if apiType.Name != clientType.Name {
					t.Errorf("Name -- Expected %v, got %v\n", apiType.Name, clientType.Name)
				}
				if apiType.UseInTable != clientType.UseInTable {
					t.Errorf("UseInTable -- Expected %v, got %v\n", apiType.UseInTable, clientType.UseInTable)
				}
			}
		}
		if !match {
			t.Errorf("Did not get a type matching %v\n", apiType.Name)
		}
	}
}
