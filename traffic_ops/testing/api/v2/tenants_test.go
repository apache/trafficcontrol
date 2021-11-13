package v2

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
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
	client "github.com/apache/trafficcontrol/traffic_ops/v2-client"
)

func TestTenants(t *testing.T) {
	WithObjs(t, []TCObj{Parameters, Tenants}, func() {
		GetTestTenants(t)
		UpdateTestTenants(t)
	})
}

func CreateTestTenants(t *testing.T) {
	for _, ten := range testData.Tenants {
		resp, err := TOSession.CreateTenant(&ten)

		if err != nil {
			t.Errorf("could not CREATE tenant %s: %v", ten.Name, err)
			continue
		}
		if resp == nil {
			t.Errorf("Traffic Ops returned null or undefined Tenant in creation response")
			continue
		}
		if resp.Response.Name != ten.Name {
			t.Errorf("expected tenant %+v; got %+v", ten, resp.Response)
		}
	}
}

func GetTestTenants(t *testing.T) {
	resp, _, err := TOSession.Tenants()
	if err != nil {
		t.Errorf("cannot GET all tenants: %v - %v", err, resp)
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
		t.Errorf("cannot GET Tenant by name: %s - %v", name, err)
	}

	newParent, _, err := TOSession.TenantByName(parentName)
	if err != nil {
		t.Errorf("cannot GET Tenant by name: %s - %v", parentName, err)
	}
	modTenant.ParentID = newParent.ID

	_, err = TOSession.UpdateTenant(strconv.Itoa(modTenant.ID), modTenant)
	if err != nil {
		t.Errorf("cannot UPDATE Tenant by id: %v", err)
	}

	// Retrieve the Tenant to check Tenant parent name got updated
	respTenant, _, err := TOSession.Tenant(strconv.Itoa(modTenant.ID))
	if err != nil {
		t.Errorf("cannot GET Tenant by name: %v - %v", name, err)
	}
	if respTenant.ParentName != parentName {
		t.Errorf("results do not match actual: %s, expected: %s", respTenant.ParentName, parentName)
	}

}

func DeleteTestTenants(t *testing.T) {

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
		for _, tn := range testData.Tenants {
			if _, ok := deletedTenants[tn.Name]; ok {
				continue
			}

			hasParent := false
			for _, otherTenant := range testData.Tenants {
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
		if len(deletedTenants) == len(testData.Tenants) {
			break
		}
		if len(deletedTenants) == initLenDeleted {
			t.Fatal("could not delete tenants: not tenant without an existing child found (cycle?)")
		}
	}
}

func TestTenantsActive(t *testing.T) {
	CreateTestCDNs(t)
	CreateTestTypes(t)
	CreateTestTenants(t)
	CreateTestParameters(t)
	CreateTestProfiles(t)
	CreateTestStatuses(t)
	CreateTestDivisions(t)
	CreateTestRegions(t)
	CreateTestPhysLocations(t)
	CreateTestCacheGroups(t)
	CreateTestServers(t)
	CreateTestDeliveryServices(t)
	CreateTestUsers(t)

	UpdateTestTenantsActive(t)

	ForceDeleteTestUsers(t)
	DeleteTestDeliveryServices(t)
	DeleteTestServers(t)
	DeleteTestCacheGroups(t)
	DeleteTestPhysLocations(t)
	DeleteTestRegions(t)
	DeleteTestDivisions(t)
	DeleteTestStatuses(t)
	DeleteTestProfiles(t)
	DeleteTestParameters(t)
	DeleteTestTenants(t)
	DeleteTestTypes(t)
	DeleteTestCDNs(t)
}

func ExtractXMLID(ds *tc.DeliveryServiceNullable) string {
	if ds.XMLID != nil {
		return *ds.XMLID
	}
	return "nil"
}

func UpdateTestTenantsActive(t *testing.T) {
	originalTenants, _, err := TOSession.Tenants()
	if err != nil {
		t.Fatalf("getting tenants error expected: nil, actual: %+v", err)
	}

	setTenantActive(t, "tenant1", true)
	setTenantActive(t, "tenant2", true)
	setTenantActive(t, "tenant3", false)

	// ds3 has tenant3. Even though tenant3 is inactive, we should still be able to get it, because our user is tenant1, which is active.
	dses, _, err := TOSession.GetDeliveryServiceByXMLIDNullable("ds3")
	if err != nil {
		t.Fatal("failed to get delivery service, when the DS's tenant was inactive (even though our user's tenant was active)")
	} else if len(dses) != 1 {
		t.Error("admin user getting delivery service ds3 with tenant3, expected: ds, actual: empty")
	}

	setTenantActive(t, "tenant1", true)
	setTenantActive(t, "tenant2", false)
	setTenantActive(t, "tenant3", true)

	// ds3 has tenant3. Even though tenant3's parent, tenant2, is inactive, we should still be able to get it, because our user is tenant1, which is active.
	_, _, err = TOSession.GetDeliveryServiceByXMLIDNullable("ds3")
	if err != nil {
		t.Fatal("failed to get delivery service, when a parent tenant was inactive (even though our user's tenant was active)")
	}

	toReqTimeout := time.Second * time.Duration(Config.Default.Session.TimeoutInSecs)
	tenant3Session, _, err := client.LoginWithAgent(TOSession.URL, "tenant3user", "pa$$word", true, "to-api-v1-client-tests/tenant3user", true, toReqTimeout)
	if err != nil {
		t.Fatalf("failed to get log in with tenant3user: " + err.Error())
	}

	tenant4Session, _, err := client.LoginWithAgent(TOSession.URL, "tenant4user", "pa$$word", true, "to-api-v1-client-tests/tenant4user", true, toReqTimeout)
	if err != nil {
		t.Fatalf("failed to get log in with tenant4user: " + err.Error())
	}

	// tenant3user with tenant3 has no access to ds3 with tenant3 when parent tenant2 is inactive
	dses, _, err = tenant3Session.GetDeliveryServiceByXMLIDNullable("ds3")
	for _, ds := range dses {
		t.Errorf("tenant3user got delivery service %+v with tenant3 but tenant3 parent tenant2 is inactive, expected: no ds", ExtractXMLID(&ds))
	}

	setTenantActive(t, "tenant1", true)
	setTenantActive(t, "tenant2", true)
	setTenantActive(t, "tenant3", false)

	// tenant3user with tenant3 has no access to ds3 with tenant3 when tenant3 is inactive
	dses, _, err = tenant3Session.GetDeliveryServiceByXMLIDNullable("ds3")
	for _, ds := range dses {
		t.Errorf("tenant3user got delivery service %+v with tenant3 but tenant3 is inactive, expected: no ds", ExtractXMLID(&ds))
	}

	setTenantActive(t, "tenant1", true)
	setTenantActive(t, "tenant2", true)
	setTenantActive(t, "tenant3", true)

	// tenant3user with tenant3 has access to ds3 with tenant3
	dses, _, err = tenant3Session.GetDeliveryServiceByXMLIDNullable("ds3")
	if err != nil {
		t.Errorf("tenant3user getting delivery service ds3 error expected: nil, actual: %+v", err)
	} else if len(dses) == 0 {
		t.Error("tenant3user getting delivery service ds3 with tenant3, expected: ds, actual: empty")
	}

	// 1. ds2 has tenant2.
	// 2. tenant3user has tenant3.
	// 3. tenant2 is not a child of tenant3 (tenant3 is a child of tenant2)
	// 4. Therefore, tenant3user should not have access to ds2
	dses, _, _ = tenant3Session.GetDeliveryServiceByXMLIDNullable("ds2")
	for _, ds := range dses {
		t.Errorf("tenant3user got delivery service %+v with tenant2, expected: no ds", ExtractXMLID(&ds))
	}

	// 1. ds1 has tenant1.
	// 2. tenant4user has tenant4.
	// 3. tenant1 is not a child of tenant4 (tenant4 is unrelated to tenant1)
	// 4. Therefore, tenant4user should not have access to ds1
	dses, _, _ = tenant4Session.GetDeliveryServiceByXMLIDNullable("ds1")
	for _, ds := range dses {
		t.Errorf("tenant4user got delivery service %+v with tenant1, expected: no ds", ExtractXMLID(&ds))
	}

	setTenantActive(t, "tenant3", false)
	dses, _, _ = tenant3Session.GetDeliveryServiceByXMLIDNullable("ds3")
	for _, ds := range dses {
		t.Errorf("tenant3user was inactive, but got delivery service %+v with tenant3, expected: no ds", ExtractXMLID(&ds))
	}

	for _, tn := range originalTenants {
		if tn.Name == "root" {
			continue
		}
		if _, err := TOSession.UpdateTenant(strconv.Itoa(tn.ID), &tn); err != nil {
			t.Fatalf("restoring original tenants: " + err.Error())
		}
	}
}

func setTenantActive(t *testing.T, name string, active bool) {
	tn, _, err := TOSession.TenantByName(name)
	if err != nil {
		t.Fatalf("cannot GET Tenant by name: %s - %v", name, err)
	}
	tn.Active = active
	_, err = TOSession.UpdateTenant(strconv.Itoa(tn.ID), tn)
	if err != nil {
		t.Fatalf("cannot UPDATE Tenant by id: %v", err)
	}
}
