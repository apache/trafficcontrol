package v3

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
)

func TestSteering(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, DeliveryServices, Users, SteeringTargets}, func() {
		GetTestSteering(t)
	})
}

func GetTestSteering(t *testing.T) {
	if len(testData.SteeringTargets) < 1 {
		t.Error("get steering: no steering target test data")
	}
	st := testData.SteeringTargets[0]
	if st.DeliveryService == nil {
		t.Error("get steering: test data missing ds")
	}

	steerings, _, err := TOSession.Steering()
	if err != nil {
		t.Errorf("steering get: getting steering: %v", err)
	}

	if len(steerings) != len(testData.SteeringTargets) {
		t.Errorf("steering get: expected %v actual %v", len(testData.SteeringTargets), len(steerings))
	}

	if steerings[0].ClientSteering {
		t.Errorf("steering get: ClientSteering expected %v actual %v", false, true)
	}
	if len(steerings[0].Targets) != 1 {
		t.Errorf("steering get: Targets expected %v actual %v", 1, len(steerings[0].Targets))
	}
	if steerings[0].Targets[0].Order != 0 {
		t.Errorf("steering get: Targets Order expected %v actual %v", 0, steerings[0].Targets[0].Order)
	}
	if testData.SteeringTargets[0].Value != nil && steerings[0].Targets[0].Weight != int32(*testData.SteeringTargets[0].Value) {
		t.Errorf("steering get: Targets Order expected %v actual %v", testData.SteeringTargets[0].Value, steerings[0].Targets[0].Weight)
	}
	if steerings[0].Targets[0].GeoOrder != nil {
		t.Errorf("steering get: Targets Order expected %v actual %+v", nil, *steerings[0].Targets[0].GeoOrder)
	}
	if steerings[0].Targets[0].Longitude != nil {
		t.Errorf("steering get: Targets Order expected %v actual %+v", nil, *steerings[0].Targets[0].Longitude)
	}
	if steerings[0].Targets[0].Latitude != nil {
		t.Errorf("steering get: Targets Order expected %v actual %+v", nil, *steerings[0].Targets[0].Latitude)
	}
}
