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
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
)

func TestServerChecks(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, TOExtensions, ServerChecks}, func() {
		CreateTestInvalidServerChecks(t)
		UpdateTestServerChecks(t)
		GetTestServerChecks(t)
	})
}

func CreateTestServerChecks(t *testing.T) {
	SwitchSession(toReqTimeout, Config.TrafficOps.URL, Config.TrafficOps.Users.Admin, Config.TrafficOps.UserPassword, Config.TrafficOps.Users.Extension, Config.TrafficOps.UserPassword)

	for _, servercheck := range testData.Serverchecks {
		resp, _, err := TOSession.InsertServerCheckStatus(servercheck)
		log.Debugf("Response: %v host_name %v check %v", *servercheck.HostName, *servercheck.Name, resp)
		if err != nil {
			t.Errorf("could not CREATE servercheck: %v\n", err)
		}
	}
	SwitchSession(toReqTimeout, Config.TrafficOps.URL, Config.TrafficOps.Users.Extension, Config.TrafficOps.UserPassword, Config.TrafficOps.Users.Admin, Config.TrafficOps.UserPassword)
}

func CreateTestInvalidServerChecks(t *testing.T) {
	toReqTimeout := time.Second * time.Duration(Config.Default.Session.TimeoutInSecs)

	_, _, err := TOSession.InsertServerCheckStatus(testData.Serverchecks[0])
	if err == nil {
		t.Errorf("expected to receive error with non extension user\n")
	}

	SwitchSession(toReqTimeout, Config.TrafficOps.URL, Config.TrafficOps.Users.Admin, Config.TrafficOps.UserPassword, Config.TrafficOps.Users.Extension, Config.TrafficOps.UserPassword)

	invalidServerCheck := tc.ServercheckRequestNullable{
		Name:     util.StrPtr("BOGUS"),
		Value:    util.IntPtr(1),
		ID:       util.IntPtr(-1),
		HostName: util.StrPtr("bogus_hostname"),
	}

	// Attempt to create a ServerCheck with invalid server ID
	_, _, err = TOSession.InsertServerCheckStatus(invalidServerCheck)
	if err == nil {
		t.Errorf("expected to receive error with invalid id\n")
	}

	invalidServerCheck.ID = nil
	// Attempt to create a ServerCheck with invalid host name
	_, _, err = TOSession.InsertServerCheckStatus(invalidServerCheck)
	if err == nil {
		t.Errorf("expected to receive error with invalid host name\n")
	}

	// get valid name to get past host check
	invalidServerCheck.Name = testData.Serverchecks[0].Name

	// Attempt to create a ServerCheck with invalid servercheck name
	_, _, err = TOSession.InsertServerCheckStatus(invalidServerCheck)
	if err == nil {
		t.Errorf("expected to receive error with invalid servercheck name\n")
	}
	SwitchSession(toReqTimeout, Config.TrafficOps.URL, Config.TrafficOps.Users.Extension, Config.TrafficOps.UserPassword, Config.TrafficOps.Users.Admin, Config.TrafficOps.UserPassword)
}

func UpdateTestServerChecks(t *testing.T) {
	SwitchSession(toReqTimeout, Config.TrafficOps.URL, Config.TrafficOps.Users.Admin, Config.TrafficOps.UserPassword, Config.TrafficOps.Users.Extension, Config.TrafficOps.UserPassword)
	for _, servercheck := range testData.Serverchecks {
		*servercheck.Value--
		resp, _, err := TOSession.InsertServerCheckStatus(servercheck)
		log.Debugf("Response: %v host_name %v check %v", *servercheck.HostName, *servercheck.Name, resp)
		if err != nil {
			t.Errorf("could not update servercheck: %v\n", err)
		}
	}
	SwitchSession(toReqTimeout, Config.TrafficOps.URL, Config.TrafficOps.Users.Extension, Config.TrafficOps.UserPassword, Config.TrafficOps.Users.Admin, Config.TrafficOps.UserPassword)
}

func GetTestServerChecks(t *testing.T) {
	hostname := testData.Serverchecks[0].HostName
	// Get server checks
	serverChecksResp, _, err := TOSession.GetServerChecks()
	if err != nil {
		t.Fatalf("could not GET serverchecks: %v\n", err)
	}
	found := false
	for _, sc := range serverChecksResp.Response {
		if sc.HostName == *hostname {
			found = true
			if sc.Checks.ORT != 12 {
				t.Errorf("%v returned for ORT value servercheck - expected 12\n", sc.Checks.ORT)
			}

			if sc.Checks.ILO != 0 {
				t.Errorf("%v returned for ILO value servercheck - expected 1\n", sc.Checks.ILO)
			}
			break
		}
	}
	if !found {
		t.Errorf("expected to find servercheck for host %v\n", hostname)
	}
}

// Need to define no-op function as TCObj interface expects a delete function
// There is no delete path for serverchecks
func DeleteTestServerChecks(t *testing.T) {
	return
}
