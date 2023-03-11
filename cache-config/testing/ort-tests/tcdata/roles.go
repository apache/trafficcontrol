package tcdata

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

import (
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

const (
	roleGood         = 0
	roleInvalidCap   = 1
	roleNeedCap      = 2
	roleBadPrivLevel = 3
)

func (r *TCData) CreateTestRoles(t *testing.T) {
	for _, role := range r.TestData.Roles {
		var alerts tc.Alerts
		alerts, _, status, err := TOSession.CreateRole(role)
		t.Log("Status Code: ", status)
		t.Log("Response: ", alerts)
		if err != nil {
			t.Errorf("creating Role failed: %v", err)
		}
	}
}

func (r *TCData) DeleteTestRoles(t *testing.T) {

	role := r.TestData.Roles[roleGood]
	// Retrieve the Role by name so we can get the id
	resp, _, status, err := TOSession.GetRoleByNameWithHdr(*role.Name, nil)
	t.Log("Status Code: ", status)
	if err != nil {
		t.Errorf("cannot GET Role by name: %s - %v", *role.Name, err)
	}
	respRole := resp[0]

	delResp, _, status, err := TOSession.DeleteRoleByID(*respRole.ID)
	t.Log("Status Code: ", status)

	if err != nil {
		t.Errorf("cannot DELETE Role by role: %v - %v", err, delResp)
	}

	// Retrieve the Role to see if it got deleted
	roleResp, _, status, err := TOSession.GetRoleByNameWithHdr(*role.Name, nil)
	t.Log("Status Code: ", status)

	if err != nil {
		t.Errorf("error deleting Role role: %v", err)
	}
	if len(roleResp) > 0 {
		t.Errorf("expected Role : %s to be deleted", *role.Name)
	}
}
