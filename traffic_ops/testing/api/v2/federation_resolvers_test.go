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
	"testing"

	"github.com/apache/trafficcontrol/v6/lib/go-tc"
	"github.com/apache/trafficcontrol/v6/lib/go-util"
)

func TestFederationResolvers(t *testing.T) {
	WithObjs(t, []TCObj{Types, FederationResolvers}, func() {
		GetTestFederationResolvers(t)
	})
}

func GetTestFederationResolvers(t *testing.T) {
	var tdlen = len(testData.FederationResolvers)
	if tdlen < 1 {
		t.Fatal("no federation resolvers test data")
	}

	frs, _, err := TOSession.GetFederationResolvers()
	if err != nil {
		t.Errorf("Unexpected error getting Federation Resolvers: %v", err)
	}
	if len(frs) != tdlen {
		t.Fatalf("Wrong number of Federation Resolvers from GET, want %d got %d", tdlen, len(frs))
	}

	var testFr = frs[0]
	if testFr.ID == nil || testFr.Type == nil || testFr.IPAddress == nil {
		t.Fatalf("Malformed federation resolver: %+v", testFr)
	}

	_ = t.Run("Get Federation Resolvers by ID", getFRByIDTest(testFr))
	_ = t.Run("Get Federation Resolvers by IPAddress", getFRByIPTest(testFr))
	_ = t.Run("Get Federation Resolvers by Type", getFRByTypeTest(testFr))
}

func getFRByIDTest(testFr tc.FederationResolver) func(*testing.T) {
	return func(t *testing.T) {
		fr, _, err := TOSession.GetFederationResolverByID(*testFr.ID)
		if err != nil {
			t.Fatalf("Unexpected error getting Federation Resolver by ID %d: %v", *testFr.ID, err)
		}

		cmpr(testFr, fr, t)

	}
}

func getFRByIPTest(testFr tc.FederationResolver) func(*testing.T) {
	return func(t *testing.T) {
		fr, _, err := TOSession.GetFederationResolverByIPAddress(*testFr.IPAddress)
		if err != nil {
			t.Fatalf("Unexpected error getting Federation Resolver by IP %s: %v", *testFr.IPAddress, err)
		}

		cmpr(testFr, fr, t)

	}
}

func getFRByTypeTest(testFr tc.FederationResolver) func(*testing.T) {
	return func(t *testing.T) {
		frs, _, err := TOSession.GetFederationResolversByType(*testFr.Type)
		if err != nil {
			t.Fatalf("Unexpected error getting Federation Resolvers by Type %s: %v", *testFr.Type, err)
		}

		if len(frs) < 1 {
			t.Errorf("Expected at least one Federation Resolver by Type %s to exist, but none did", *testFr.Type)
		}

		for _, fr := range frs {
			if fr.ID == nil {
				t.Error("Retrieved Federation Resolver has nil ID")
			}
			if fr.IPAddress == nil {
				t.Error("Retrieved Federation Resolver has nil IPAddress")
			}
			if fr.Type == nil {
				t.Error("Retrieved Federation Resolver has nil Type")
			} else if *fr.Type != *testFr.Type {
				t.Errorf("Retrieved Federation Resolver has incorrect Type; want %s, got %s", *testFr.Type, *fr.Type)
			}
		}
	}
}

func cmpr(testFr, apiFr tc.FederationResolver, t *testing.T) {
	if apiFr.ID == nil {
		t.Error("Retrieved Federation Resolver has nil ID")
	} else if *apiFr.ID != *testFr.ID {
		t.Errorf("Retrieved Federation Resolver has incorrect ID; want %d, got %d", *testFr.ID, *apiFr.ID)
	}

	if apiFr.IPAddress == nil {
		t.Error("Retrieved Federation Resolver has nil IP address")
	} else if *apiFr.IPAddress != *testFr.IPAddress {
		t.Errorf("Retrieved Federation Resolver has incorrect IPAddress; want %s, got %s", *testFr.IPAddress, *apiFr.IPAddress)
	}

	if apiFr.Type == nil {
		t.Error("Retrieved Federation Resolver has nil Type")
	} else if *apiFr.Type != *testFr.Type {
		t.Errorf("Retrieved Federation Resolver has incorrect Type; want %s, got %s", *testFr.Type, *apiFr.Type)
	}
}

func CreateTestFederationResolvers(t *testing.T) {
	for _, fr := range testData.FederationResolvers {
		if fr.Type == nil {
			t.Fatal("testData Federation Resolver has nil Type")
		}

		tid, _, err := TOSession.GetTypeByName(*fr.Type)
		if err != nil {
			t.Fatalf("Couldn't get an ID for type %s", *fr.Type)
		}
		if len(tid) != 1 {
			t.Fatalf("Expected exactly one Type by name %s, got %d", *fr.Type, len(tid))
		}

		fr.TypeID = util.UIntPtr(uint(tid[0].ID))

		alerts, _, err := TOSession.CreateFederationResolver(fr)
		if err != nil {
			t.Fatalf("failed to create Federation resolver %+v: %v\n\talerts: %+v", fr, err, alerts)
		}
		for _, a := range alerts.Alerts {
			if a.Level != tc.SuccessLevel.String() {
				t.Errorf("Unexpected %s creating a federation resolver: %s", a.Level, a.Text)
			} else {
				t.Logf("Received expected success creating federation resolver: %s", a.Text)
			}
		}
	}

	var invalidFR tc.FederationResolver
	alerts, _, err := TOSession.CreateFederationResolver(invalidFR)
	if err == nil {
		t.Error("Expected an error creating a bad Federation Resolver, but didn't get one")
	}
	for _, a := range alerts.Alerts {
		if a.Level == tc.SuccessLevel.String() {
			t.Errorf("Unexpected success creating a bad Federation Resolver: %s", a.Text)
		} else {
			t.Logf("Received expected %s creating federation resolver: %s", a.Level, a.Text)
		}
	}

	invalidFR.TypeID = util.UIntPtr(1)
	invalidFR.IPAddress = util.StrPtr("not a valid IP address")
	alerts, _, err = TOSession.CreateFederationResolver(invalidFR)
	if err == nil {
		t.Error("Expected an error creating a bad Federation Resolver, but didn't get one")
	}
	for _, a := range alerts.Alerts {
		if a.Level == tc.SuccessLevel.String() {
			t.Errorf("Unexpected success creating a bad Federation Resolver: %s", a.Text)
		} else {
			t.Logf("Received expected %s creating a bad federation resolver: %s", a.Level, a.Text)
		}
	}
}

func DeleteTestFederationResolvers(t *testing.T) {
	frs, _, err := TOSession.GetFederationResolvers()
	if err != nil {
		t.Errorf("Unexpected error getting Federation Resolvers: %v", err)
	}
	if len(frs) < 1 {
		t.Fatal("Found no Federation Resolvers to delete")
	}
	for _, fr := range frs {
		if fr.ID == nil {
			t.Fatalf("Malformed Federation Resolver: %+v", fr)
		}
		alerts, _, err := TOSession.DeleteFederationResolver(*fr.ID)
		if err != nil {
			t.Fatalf("failed to delete Federation Resolver %+v: %v\n\talerts: %+v", fr, err, alerts)
		}
		for _, a := range alerts.Alerts {
			if a.Level != tc.SuccessLevel.String() {
				t.Errorf("Unexpected %s deleting a federation resolver: %s", a.Level, a.Text)
			} else {
				t.Logf("Received expected success deleting federation resolver: %s", a.Text)
			}
		}
	}

	alerts, _, err := TOSession.DeleteFederationResolver(0)
	if err == nil {
		t.Error("Expected an error deleting a non-existent Federation Resolver, but didn't get one")
	}
	for _, a := range alerts.Alerts {
		if a.Level == tc.SuccessLevel.String() {
			t.Errorf("Unexpected success deleting a non-existent Federation Resolver: %s", a.Text)
		} else {
			t.Logf("Received expected %s deleting a non-existent federation resolver: %s", a.Level, a.Text)
		}
	}

}
