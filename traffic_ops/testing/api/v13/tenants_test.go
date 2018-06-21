package v13

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

	tc "github.com/apache/trafficcontrol/lib/go-tc"
)

func TestTenants(t *testing.T) {

	CreateTestTenants(t)
	UpdateTestTenants(t)
	GetTestTenants(t)
	//DeleteTestTenants(t)
}

func CreateTestTenants(t *testing.T) {
	for _, ten := range testData.Tenants {
		// testData does not define ParentID -- look up by name and fill in
		if ten.ParentID == 0 {
			parents, _, err := TOSession.GetTenantByName(ten.ParentName)
			if err != nil {
				t.Errorf("parent tenant %s: %++v", ten.ParentName, err)
				continue
			}
			ten.ParentID = parents[0].ID
		}
		resp, _, err := TOSession.CreateTenant(ten)
		t.Logf("response: %++v", resp)

		if err != nil {
			t.Errorf("could not CREATE tenant %s: %v\n", ten.Name, err)
		}
	}
}

func GetTestTenants(t *testing.T) {
	resp, _, err := TOSession.GetTenants()
	if err != nil {
		t.Errorf("cannot GET all tenants: %v - %v\n", err, resp)
		return
	}

	// expect root and badTenant (defined in todb.go) + all defined in testData.Tenants
	if len(resp) != 2+len(testData.Tenants) {
		t.Errorf("expected %d tenants,  got %d", 2+len(testData.Tenants), len(resp))
	}

	for _, ten := range testData.Tenants {
		resp, _, err := TOSession.GetTenantByName(ten.Name)
		if err != nil {
			t.Errorf("cannot GET Tenant by name: %v - %v\n", err, resp)
			continue
		}
		if len(resp) != 1 {
			t.Errorf("expected 1 tenant for %s,  got %d", ten.Name, len(resp))
			continue
		}
	}
}

func UpdateTestTenants(t *testing.T) {

	// Retrieve the Tenant by name so we can get the id for the Update
	name := "tenant2"
	parentName := "tenant1"
	resp, _, err := TOSession.GetTenantByName(name)
	if err != nil {
		t.Errorf("cannot GET Tenant by name: %s - %v\n", name, err)
	}
	modTenant := resp[0]
	resp, _, err = TOSession.GetTenantByName(parentName)
	if err != nil {
		t.Errorf("cannot GET Tenant by name: %s - %v\n", name, err)
	}
	newParent := resp[0]
	modTenant.ParentID = newParent.ID
	var alert tc.Alerts
	alert, _, err = TOSession.UpdateTenantByID(modTenant.ID, modTenant)
	if err != nil {
		t.Errorf("cannot UPDATE Tenant by id: %v - %v\n", err, alert)
	}

	// Retrieve the Tenant to check Tenant parent name got updated
	resp, _, err = TOSession.GetTenantByID(modTenant.ID)
	if err != nil {
		t.Errorf("cannot GET Tenant by name: %v - %v\n", name, err)
	}
	t.Logf("modified: %++v", resp)
	respTenant := resp[0]
	if respTenant.ParentName != parentName {
		t.Errorf("results do not match actual: %s, expected: %s\n", respTenant.ParentName, parentName)
	}

}

func DeleteTestTenants(t *testing.T) {

	for _, ten := range testData.Tenants {
		// Retrieve the Tenant by name so we can get the id for the Update
		resp, _, err := TOSession.GetTenantByName(ten.Name)
		if err != nil {
			t.Errorf("cannot GET Tenant by name: %v - %v\n", ten.Name, err)
		}
		respTenant := resp[0]

		delResp, _, err := TOSession.DeleteTenantByID(respTenant.ID)
		if err != nil {
			t.Errorf("cannot DELETE Tenant by name: %v - %v\n", err, delResp)
		}

		// Retrieve the Tenant to see if it got deleted
		Tenants, _, err := TOSession.GetTenantByName(ten.Name)
		if err != nil {
			t.Errorf("error deleting Tenant name: %s\n", err.Error())
		}
		if len(Tenants) > 0 {
			t.Errorf("expected Tenant name: %s to be deleted\n", ten.Name)
		}
	}
}
