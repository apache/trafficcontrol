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

	tc "github.com/apache/trafficcontrol/lib/go-tc"
)

func TestTenants(t *testing.T) {
	CreateTestTenants(t)
	GetTestTenants(t)
	UpdateTestTenants(t)
	DeleteTestTenants(t)
}

func CreateTestTenants(t *testing.T) {
	for _, ten := range testData.Tenants {
		resp, err := TOSession.CreateTenant(&ten)

		if err != nil {
			t.Errorf("could not CREATE tenant %s: %v\n", ten.Name, err)
		}
		if resp.Response.Name != ten.Name {
			t.Errorf("expected tenant %+v; got %+v", ten, resp.Response)
		}
	}
}

func GetTestTenants(t *testing.T) {
	resp, _, err := TOSession.Tenants()
	if err != nil {
		t.Errorf("cannot GET all tenants: %v - %v\n", err, resp)
		return
	}
	foundTenants := make(map[string]tc.Tenant, len(resp))
	for _, ten := range resp {
		foundTenants[ten.Name] = ten
	}

	// expect root and badTenant (defined in todb.go) + all defined in testData.Tenants
	if len(resp) != 2+len(testData.Tenants) {
		t.Errorf("expected %d tenants,  got %d", 2+len(testData.Tenants), len(resp))
	}

	for _, ten := range testData.Tenants {
		if ft, ok := foundTenants[ten.Name]; ok {
			if ft.ParentName != ten.ParentName {
				t.Errorf("tenant %s: expected parent %s,  got %s", ten.Name, ten.ParentName, ft.ParentName)
			}
		} else {
			t.Errorf("expected tenant %s: not found", ten.Name)
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

	newParent, _, err := TOSession.TenantByName(parentName)
	if err != nil {
		t.Errorf("cannot GET Tenant by name: %s - %v\n", parentName, err)
	}
	modTenant.ParentID = newParent.ID

	_, err = TOSession.UpdateTenant(strconv.Itoa(modTenant.ID), modTenant)
	if err != nil {
		t.Errorf("cannot UPDATE Tenant by id: %v\n", err)
	}

	// Retrieve the Tenant to check Tenant parent name got updated
	respTenant, _, err := TOSession.Tenant(strconv.Itoa(modTenant.ID))
	if err != nil {
		t.Errorf("cannot GET Tenant by name: %v - %v\n", name, err)
	}
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
	expected := `Tenant '` + strconv.Itoa(tenant1.ID) + `' has child tenants. Please update these child tenants and retry.`
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
