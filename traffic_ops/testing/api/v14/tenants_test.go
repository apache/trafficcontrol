package v14

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
	"strconv"
	"strings"
	"testing"
)

func TestTenants(t *testing.T) {
	CreateTestTenants(t)
	GetTestTenants(t)
	UpdateTestTenants(t)
	DeleteTestTenants(t)
}

func CreateTestTenants(t *testing.T) {
	for _, ten := range testData.Tenants {
		// testData does not define ParentID -- look up by name and fill in
		if ten.ParentID == 0 {
			parent, _, err := TOSession.TenantByName(ten.ParentName)
			if err != nil || parent == nil {
				t.Errorf("parent tenant %s: %++v", ten.ParentName, err)
				continue
			}
			ten.ParentID = parent.ID
		}
		resp, err := TOSession.CreateTenant(&ten)
		t.Logf("response: %++v", resp)

		if err != nil {
			t.Errorf("could not CREATE tenant %s: %v\n", ten.Name, err)
		}
	}
}

func GetTestTenants(t *testing.T) {
	resp, _, err := TOSession.Tenants()
	if err != nil {
		t.Errorf("cannot GET all tenants: %v - %v\n", err, resp)
		return
	}

	t.Logf("resp: %++v\n", resp)
	// expect root and badTenant (defined in todb.go) + all defined in testData.Tenants
	if len(resp) != 2+len(testData.Tenants) {
		t.Errorf("expected %d tenants,  got %d", 2+len(testData.Tenants), len(resp))
	}

	for _, ten := range testData.Tenants {
		resp, _, err := TOSession.TenantByName(ten.Name)
		if err != nil {
			t.Errorf("cannot GET Tenant by name: %v - %v\n", err, resp)
			continue
		}
		if resp.Name != ten.Name {
			t.Errorf("expected tenant %s,  got %s", ten.Name, resp.Name)
			continue
		}
	}
}

func UpdateTestTenants(t *testing.T) {

	// Retrieve the Tenant by name so we can get the id for the Update
	name := "tenant2"
	parentName := "tenant1"
	modTenant, _, err := TOSession.TenantByName(name)
	if err != nil {
		t.Errorf("cannot GET Tenant by name: %s - %v\n", name, err)
	}

	t.Logf("modTenant is %++v", modTenant)
	newParent, _, err := TOSession.TenantByName(parentName)
	if err != nil {
		t.Errorf("cannot GET Tenant by name: %s - %v\n", parentName, err)
	}
	t.Logf("newParent is %++v", newParent)
	modTenant.ParentID = newParent.ID

	resp, err := TOSession.UpdateTenant(strconv.Itoa(modTenant.ID), modTenant)
	if err != nil {
		t.Errorf("cannot UPDATE Tenant by id: %v\n", err)
	}

	t.Logf("AFTER UPDATE modTenant is %++v", modTenant)
	// Retrieve the Tenant to check Tenant parent name got updated
	respTenant, _, err := TOSession.Tenant(strconv.Itoa(modTenant.ID))
	if err != nil {
		t.Errorf("cannot GET Tenant by name: %v - %v\n", name, err)
	}
	t.Logf("modified: %++v", resp)
	if respTenant.ParentName != parentName {
		t.Errorf("results do not match actual: %s, expected: %s\n", respTenant.ParentName, parentName)
	}

}

func DeleteTestTenants(t *testing.T) {

	t1 := "tenant1"
	tenant1, _, err := TOSession.TenantByName(t1)

	if err != nil {
		t.Errorf("cannot GET Tenant by name: %v - %v\n", t1, err)
	}

	_, err = TOSession.DeleteTenant(strconv.Itoa(tenant1.ID))
	if err == nil {
		t.Errorf("%s has child tenants -- should not be able to delete", t1)
	}
	expected := `Tenant '`+ strconv.Itoa(tenant1.ID) + `' has child tenants. Please update these child tenants and retry.`
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("expected error: %s;  got %s", expected, err.Error())
	}

	t2 := "tenant2"
	tenant2, _, err := TOSession.TenantByName(t2)
	_, err = TOSession.DeleteTenant(strconv.Itoa(tenant2.ID))
	if err != nil {
		t.Errorf("error deleting tenant %s: %v", t2, err)
	}

	// Now should be able to delete t1
	tenant1, _, err = TOSession.TenantByName(t1)
	_, err = TOSession.DeleteTenant(strconv.Itoa(tenant1.ID))
	if err != nil {
		t.Errorf("error deleting tenant %s: %v", t1, err)
	}

}
