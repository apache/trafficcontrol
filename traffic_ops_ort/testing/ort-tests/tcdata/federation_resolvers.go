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
	"github.com/apache/trafficcontrol/lib/go-util"
)

func (r *TCData) CreateTestFederationResolvers(t *testing.T) {
	for _, fr := range r.TestData.FederationResolvers {
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

func (r *TCData) DeleteTestFederationResolvers(t *testing.T) {
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
