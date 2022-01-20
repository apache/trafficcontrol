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
	"time"

	"github.com/apache/trafficcontrol/lib/go-util"
	client "github.com/apache/trafficcontrol/traffic_ops/v3-client"
)

var SteeringUserSession *client.Session

func (r *TCData) CreateTestSteeringTargets(t *testing.T) {
	for _, st := range r.TestData.SteeringTargets {
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
			respDS, _, err := SteeringUserSession.GetDeliveryServiceByXMLIDNullableWithHdr(string(*st.DeliveryService), nil)
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
			respTarget, _, err := SteeringUserSession.GetDeliveryServiceByXMLIDNullableWithHdr(string(*st.Target), nil)
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

func (r *TCData) DeleteTestSteeringTargets(t *testing.T) {
	dsIDs := []uint64{}
	for _, st := range r.TestData.SteeringTargets {
		if st.DeliveryService == nil {
			t.Fatal("deleting steering target: test data missing ds")
		}
		if st.Target == nil {
			t.Fatal("deleting steering target: test data missing target")
		}

		respDS, _, err := SteeringUserSession.GetDeliveryServiceByXMLIDNullableWithHdr(string(*st.DeliveryService), nil)
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

		respTarget, _, err := SteeringUserSession.GetDeliveryServiceByXMLIDNullableWithHdr(string(*st.Target), nil)
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
			t.Fatalf("deleting steering target: deleting: %v", err)
		}
	}

	for _, dsID := range dsIDs {
		sts, _, err := SteeringUserSession.GetSteeringTargets(int(dsID))
		if err != nil {
			t.Fatalf("deleting steering targets: getting steering target: %v", err)
		}
		if len(sts) != 0 {
			t.Fatalf("deleting steering targets: after delete, getting steering target: expected 0 actual %d", len(sts))
		}
	}
}

// SetupSteeringTargets calls the CreateSteeringTargets test. It also sets the steering user session
// with the logged in steering user. SteeringUserSession is used by steering target test functions.
// Running this function depends on CreateTestUsers.
func (r *TCData) SetupSteeringTargets(t *testing.T) {
	var err error
	toReqTimeout := time.Second * time.Duration(r.Config.Default.Session.TimeoutInSecs)
	SteeringUserSession, _, err = client.LoginWithAgent(TOSession.URL, "steering", "pa$$word", true, "to-api-v1-client-tests/steering", true, toReqTimeout)
	if err != nil {
		t.Fatalf("failed to get log in with steering user: %v", err)
	}

	r.CreateTestSteeringTargets(t)
}
