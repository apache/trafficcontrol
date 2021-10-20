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
	"time"

	"github.com/apache/trafficcontrol/v6/lib/go-util"
	client "github.com/apache/trafficcontrol/v6/traffic_ops/v2-client"
)

var SteeringUserSession *client.Session

func TestSteeringTargets(t *testing.T) {

	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, DeliveryServices, Users, SteeringTargets}, func() {
		GetTestSteeringTargets(t)
		UpdateTestSteeringTargets(t)
	})

}

// SetupSteeringTargets calls the CreateSteeringTargets test. It also sets the steering user session
// with the logged in steering user. SteeringUserSession is used by steering target test functions.
// Running this function depends on CreateTestUsers.
func SetupSteeringTargets(t *testing.T) {
	var err error
	toReqTimeout := time.Second * time.Duration(Config.Default.Session.TimeoutInSecs)
	SteeringUserSession, _, err = client.LoginWithAgent(TOSession.URL, "steering", "pa$$word", true, "to-api-v1-client-tests/steering", true, toReqTimeout)
	if err != nil {
		t.Fatalf("failed to get log in with steering user: %v", err.Error())
	}

	CreateTestSteeringTargets(t)
}

func CreateTestSteeringTargets(t *testing.T) {
	for _, st := range testData.SteeringTargets {
		if st.Type == nil {
			t.Fatal("creating steering target: test data missing type")
		}
		if st.DeliveryService == nil {
			t.Fatal("creating steering target: test data missing ds")
		}
		if st.Target == nil {
			t.Fatal("creating steering target: test data missing target")
		}

		{
			respTypes, _, err := SteeringUserSession.GetTypeByName(*st.Type)
			if err != nil {
				t.Fatalf("creating steering target: getting type: %v", err)
			} else if len(respTypes) < 1 {
				t.Fatal("creating steering target: getting type: not found")
			}
			st.TypeID = util.IntPtr(respTypes[0].ID)
		}
		{
			respDS, _, err := SteeringUserSession.GetDeliveryServiceByXMLIDNullable(string(*st.DeliveryService))
			if err != nil {
				t.Fatalf("creating steering target: getting ds: %v", err)
			} else if len(respDS) < 1 {
				t.Fatal("creating steering target: getting ds: not found")
			} else if respDS[0].ID == nil {
				t.Fatal("creating steering target: getting ds: nil ID returned")
			}
			dsID := uint64(*respDS[0].ID)
			st.DeliveryServiceID = &dsID
		}
		{
			respTarget, _, err := SteeringUserSession.GetDeliveryServiceByXMLIDNullable(string(*st.Target))
			if err != nil {
				t.Fatalf("creating steering target: getting target ds: %v", err)
			} else if len(respTarget) < 1 {
				t.Fatal("creating steering target: getting target ds: not found")
			} else if respTarget[0].ID == nil {
				t.Fatal("creating steering target: getting target ds: nil ID returned")
			}
			targetID := uint64(*respTarget[0].ID)
			st.TargetID = &targetID
		}

		resp, _, err := SteeringUserSession.CreateSteeringTarget(st)
		t.Log("Response: ", resp)
		if err != nil {
			t.Fatalf("creating steering target: %v", err)
		}
	}
}

func UpdateTestSteeringTargets(t *testing.T) {
	if len(testData.SteeringTargets) < 1 {
		t.Fatal("updating steering target: no steering target test data")
	}
	st := testData.SteeringTargets[0]
	if st.DeliveryService == nil {
		t.Fatal("updating steering target: test data missing ds")
	}
	if st.Target == nil {
		t.Fatal("updating steering target: test data missing target")
	}

	respDS, _, err := SteeringUserSession.GetDeliveryServiceByXMLIDNullable(string(*st.DeliveryService))
	if err != nil {
		t.Fatalf("updating steering target: getting ds: %v", err)
	}
	if len(respDS) < 1 {
		t.Fatal("updating steering target: getting ds: not found")
	}
	if respDS[0].ID == nil {
		t.Fatal("updating steering target: getting ds: nil id returned")
	}
	dsID := *respDS[0].ID

	sts, _, err := SteeringUserSession.GetSteeringTargets(dsID)
	if err != nil {
		t.Fatalf("updating steering targets: getting steering target: %v", err)
	}
	if len(sts) < 1 {
		t.Fatal("updating steering targets: getting steering target: got 0")
	}
	st = sts[0]

	expected := util.JSONIntStr(-12345)
	if st.Value != nil && *st.Value == expected {
		expected++
	}
	st.Value = &expected

	_, _, err = SteeringUserSession.UpdateSteeringTarget(st)
	if err != nil {
		t.Fatalf("updating steering targets: updating: %+v", err)
	}

	sts, _, err = SteeringUserSession.GetSteeringTargets(dsID)
	if err != nil {
		t.Fatalf("updating steering targets: getting updated steering target: %v", err)
	}
	if len(sts) < 1 {
		t.Fatal("updating steering targets: getting updated steering target: got 0")
	}
	actual := sts[0]

	if actual.DeliveryServiceID == nil {
		t.Fatalf("steering target update: ds id expected %v actual %v", dsID, nil)
	} else if *actual.DeliveryServiceID != uint64(dsID) {
		t.Fatalf("steering target update: ds id expected %v actual %v", dsID, *actual.DeliveryServiceID)
	}
	if actual.TargetID == nil {
		t.Fatalf("steering target update: ds id expected %v actual %v", dsID, nil)
	} else if *actual.TargetID != *st.TargetID {
		t.Fatalf("steering target update: ds id expected %v actual %v", *st.TargetID, *actual.TargetID)
	}
	if actual.TypeID == nil {
		t.Fatalf("steering target update: ds id expected %v actual %v", *st.TypeID, nil)
	} else if *actual.TypeID != *st.TypeID {
		t.Fatalf("steering target update: ds id expected %v actual %v", *st.TypeID, *actual.TypeID)
	}
	if actual.DeliveryService == nil {
		t.Fatalf("steering target update: ds expected %v actual %v", *st.DeliveryService, nil)
	} else if *st.DeliveryService != *actual.DeliveryService {
		t.Fatalf("steering target update: ds name expected %v actual %v", *st.DeliveryService, *actual.DeliveryService)
	}
	if actual.Target == nil {
		t.Fatalf("steering target update: target expected %v actual %v", *st.Target, nil)
	} else if *st.Target != *actual.Target {
		t.Fatalf("steering target update: target expected %v actual %v", *st.Target, *actual.Target)
	}
	if actual.Type == nil {
		t.Fatalf("steering target update: type expected %v actual %v", *st.Type, nil)
	} else if *st.Type != *actual.Type {
		t.Fatalf("steering target update: type expected %v actual %v", *st.Type, *actual.Type)
	}
	if actual.Value == nil {
		t.Fatalf("steering target update: ds expected %v actual %v", *st.Value, nil)
	} else if *st.Value != *actual.Value {
		t.Fatalf("steering target update: value expected %v actual %v", *st.Value, actual.Value)
	}
}

func GetTestSteeringTargets(t *testing.T) {
	if len(testData.SteeringTargets) < 1 {
		t.Fatal("updating steering target: no steering target test data")
	}
	st := testData.SteeringTargets[0]
	if st.DeliveryService == nil {
		t.Fatal("updating steering target: test data missing ds")
	}

	respDS, _, err := SteeringUserSession.GetDeliveryServiceByXMLIDNullable(string(*st.DeliveryService))
	if err != nil {
		t.Fatalf("creating steering target: getting ds: %v", err)
	} else if len(respDS) < 1 {
		t.Fatal("steering target get: getting ds: not found")
	} else if respDS[0].ID == nil {
		t.Fatal("steering target get: getting ds: nil id returned")
	}
	dsID := *respDS[0].ID

	sts, _, err := SteeringUserSession.GetSteeringTargets(dsID)
	if err != nil {
		t.Fatalf("steering target get: getting steering target: %v", err)
	}

	if len(sts) != len(testData.SteeringTargets) {
		t.Fatalf("steering target get: expected %v actual %v", len(testData.SteeringTargets), len(sts))
	}

	expected := testData.SteeringTargets[0]
	actual := sts[0]

	if actual.DeliveryServiceID == nil {
		t.Fatalf("steering target get: ds id expected %v actual %v", dsID, nil)
	} else if *actual.DeliveryServiceID != uint64(dsID) {
		t.Fatalf("steering target get: ds id expected %v actual %v", dsID, *actual.DeliveryServiceID)
	}
	if actual.DeliveryService == nil {
		t.Fatalf("steering target get: ds expected %v actual %v", expected.DeliveryService, nil)
	} else if *expected.DeliveryService != *actual.DeliveryService {
		t.Fatalf("steering target get: ds name expected %v actual %v", expected.DeliveryService, actual.DeliveryService)
	}
	if actual.Target == nil {
		t.Fatalf("steering target get: target expected %v actual %v", expected.Target, nil)
	} else if *expected.Target != *actual.Target {
		t.Fatalf("steering target get: target expected %v actual %v", expected.Target, actual.Target)
	}
	if actual.Type == nil {
		t.Fatalf("steering target get: type expected %v actual %v", expected.Type, nil)
	} else if *expected.Type != *actual.Type {
		t.Fatalf("steering target get: type expected %v actual %v", expected.Type, actual.Type)
	}
	if actual.Value == nil {
		t.Fatalf("steering target get: ds expected %v actual %v", expected.Value, nil)
	} else if *expected.Value != *actual.Value {
		t.Fatalf("steering target get: value expected %v actual %v", *expected.Value, *actual.Value)
	}
}

func DeleteTestSteeringTargets(t *testing.T) {
	dsIDs := []uint64{}
	if SteeringUserSession == nil {
		t.Fatal("Steering user session is nil")
	}
	for _, st := range testData.SteeringTargets {
		if st.DeliveryService == nil {
			t.Fatal("deleting steering target: test data missing ds")
		}
		if st.Target == nil {
			t.Fatal("deleting steering target: test data missing target")
		}

		respDS, _, err := SteeringUserSession.GetDeliveryServiceByXMLIDNullable(string(*st.DeliveryService))
		if err != nil {
			t.Fatalf("deleting steering target: getting ds: %v", err)
		} else if len(respDS) < 1 {
			t.Fatal("deleting steering target: getting ds: not found")
		} else if respDS[0].ID == nil {
			t.Fatal("deleting steering target: getting ds: nil ID returned")
		}
		dsID := uint64(*respDS[0].ID)
		st.DeliveryServiceID = &dsID

		dsIDs = append(dsIDs, dsID)

		respTarget, _, err := SteeringUserSession.GetDeliveryServiceByXMLIDNullable(string(*st.Target))
		if err != nil {
			t.Fatalf("deleting steering target: getting target ds: %v", err)
		} else if len(respTarget) < 1 {
			t.Fatal("deleting steering target: getting target ds: not found")
		} else if respTarget[0].ID == nil {
			t.Fatal("deleting steering target: getting target ds: not found")
		}
		targetID := uint64(*respTarget[0].ID)
		st.TargetID = &targetID

		_, _, err = SteeringUserSession.DeleteSteeringTarget(int(*st.DeliveryServiceID), int(*st.TargetID))
		if err != nil {
			t.Fatalf("deleting steering target: deleting: %+v", err)
		}
	}

	for _, dsID := range dsIDs {
		sts, _, err := SteeringUserSession.GetSteeringTargets(int(dsID))
		if err != nil {
			t.Fatalf("deleting steering targets: getting steering target: %v", err)
		}
		if len(sts) != 0 {
			t.Fatalf("deleting steering targets: after delete, getting steering target: expected 0 actual %+v", len(sts))
		}
	}
}
