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

func TestUsers(t *testing.T) {

	uri := fmt.Sprintf("/api/1.2/users.json")
	resp, err := Request(*to, "GET", uri, nil)
	if err != nil {
		t.Errorf("Could not get %s reponse was: %v\n", uri, err)
		t.FailNow()
	}

	defer resp.Body.Close()
	var apiUserRes traffic_ops.UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiUserRes); err != nil {
		t.Errorf("Could not decode user json.  Error is: %v\n", err)
		t.FailNow()
	}
	apiUsers := apiUserRes.Response

	clientUsers, err := to.Users()
	if err != nil {
		t.Errorf("Could not get users from client.  Error is: %v\n", err)
		t.FailNow()
	}

	if len(apiUsers) != len(clientUsers) {
		t.Errorf("Users Response Length -- expected %v, got %v\n", len(apiUsers), len(clientUsers))
	}

	for _, apiUser := range apiUsers {
		match := false
		for _, clientUser := range clientUsers {
			if apiUser.ID == clientUser.ID {
				match = true
				if apiUser != clientUser {
					t.Errorf("apiUser and clientUser do not match! Expected %+v, got %+v\n", apiUser, clientUser)
				}
			}
		}
		if !match {
			t.Errorf("Did not get a user matching %v\n", apiUser.Email)
		}
	}
}
