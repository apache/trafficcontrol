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
	"strconv"
	"strings"
	"testing"
)

func (r *TCData) CreateTestTenants(t *testing.T) {
	for _, ten := range r.TestData.Tenants {
		resp, err := TOSession.CreateTenant(&ten)

		if err != nil {
			t.Errorf("could not CREATE tenant %s: %v", ten.Name, err)
		}
		if resp.Response.Name != ten.Name {
			t.Errorf("expected tenant %+v; got %+v", ten, resp.Response)
		}
	}
}

func (r *TCData) DeleteTestTenants(t *testing.T) {

	t1 := "tenant1"
	tenant1, _, err := TOSession.TenantByName(t1)

	if err != nil {
		t.Errorf("cannot GET Tenant by name: %v - %v", t1, err)
	}
	expectedChildDeleteErrMsg := `Tenant '` + strconv.Itoa(tenant1.ID) + `' has child tenants. Please update these child tenants and retry.`
	if _, err := TOSession.DeleteTenant(strconv.Itoa(tenant1.ID)); err == nil {
		t.Fatalf("%s has child tenants -- should not be able to delete", t1)
	} else if !strings.Contains(err.Error(), expectedChildDeleteErrMsg) {
		t.Errorf("expected error: %s;  got %s", expectedChildDeleteErrMsg, err.Error())
	}

	deletedTenants := map[string]struct{}{}
	for {
		initLenDeleted := len(deletedTenants)
		for _, tn := range r.TestData.Tenants {
			if _, ok := deletedTenants[tn.Name]; ok {
				continue
			}

			hasParent := false
			for _, otherTenant := range r.TestData.Tenants {
				if _, ok := deletedTenants[otherTenant.Name]; ok {
					continue
				}
				if otherTenant.ParentName == tn.Name {
					hasParent = true
					break
				}
			}
			if hasParent {
				continue
			}

			toTenant, _, err := TOSession.TenantByName(tn.Name)
			if err != nil {
				t.Fatalf("getting tenant %s: %v", tn.Name, err)
			}
			if _, err = TOSession.DeleteTenant(strconv.Itoa(toTenant.ID)); err != nil {
				t.Fatalf("deleting tenant %s: %v", toTenant.Name, err)
			}
			deletedTenants[tn.Name] = struct{}{}

		}
		if len(deletedTenants) == len(r.TestData.Tenants) {
			break
		}
		if len(deletedTenants) == initLenDeleted {
			t.Fatal("could not delete tenants: not tenant without an existing child found (cycle?)")
		}
	}
}
