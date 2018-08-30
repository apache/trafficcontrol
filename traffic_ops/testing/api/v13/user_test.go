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
	"encoding/json"
	"github.com/apache/trafficcontrol/lib/go-log"
	"testing"
)

func TestUsers(t *testing.T) {
	CreateTestCDNs(t)
	CreateTestTypes(t)
	CreateTestProfiles(t)
	CreateTestStatuses(t)
	CreateTestDivisions(t)
	CreateTestRegions(t)
	CreateTestPhysLocations(t)
	CreateTestCacheGroups(t)
	CreateTestDeliveryServices(t)

	GetTestUsers(t)
	GetTestUserCurrent(t)

	DeleteTestDeliveryServices(t)
	DeleteTestCacheGroups(t)
	DeleteTestPhysLocations(t)
	DeleteTestRegions(t)
	DeleteTestDivisions(t)
	DeleteTestStatuses(t)
	DeleteTestProfiles(t)
	DeleteTestTypes(t)
	DeleteTestCDNs(t)

}

const SessionUserName = "admin" // TODO make dynamic?

func GetTestUsers(t *testing.T) {

	resp, _, err := TOSession.GetUsers()
	if err != nil {
		t.Fatalf("cannot GET users: %v\n", err)
	}
	if len(resp) == 0 {
		t.Fatalf("no users, must have at least 1 user to test\n")
	}

	log.Debugf("Response: %s\n")
	for _, user := range resp {
		bytes, _ := json.Marshal(user)
		log.Debugf("%s\n", bytes)
	}
}

func GetTestUserCurrent(t *testing.T) {
	log.Debugln("GetTestUserCurrent")
	user, _, err := TOSession.GetUserCurrent()
	if err != nil {
		t.Fatalf("cannot GET current user: %v\n", err)
	}
	if user.UserName == nil {
		t.Fatalf("current user expected: %v actual: %v\n", SessionUserName, nil)
	}
	if *user.UserName != SessionUserName {
		t.Fatalf("current user expected: %v actual: %v\n", SessionUserName, *user.UserName)
	}
}
