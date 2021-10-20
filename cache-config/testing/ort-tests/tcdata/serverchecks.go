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

	"github.com/apache/trafficcontrol/v6/lib/go-tc"
	"github.com/apache/trafficcontrol/v6/lib/go-util"
)

func (r *TCData) CreateTestServerChecks(t *testing.T) {
	toReqTimeout := time.Second * time.Duration(r.Config.Default.Session.TimeoutInSecs)
	r.SwitchSession(toReqTimeout, r.Config.TrafficOps.URL, r.Config.TrafficOps.Users.Admin, r.Config.TrafficOps.UserPassword, r.Config.TrafficOps.Users.Extension, r.Config.TrafficOps.UserPassword)

	for _, servercheck := range r.TestData.Serverchecks {
		resp, _, err := TOSession.InsertServerCheckStatus(servercheck)
		t.Logf("Response: %v host_name %v check %v", *servercheck.HostName, *servercheck.Name, resp)
		if err != nil {
			t.Errorf("could not CREATE servercheck: %v", err)
		}
	}
	r.SwitchSession(toReqTimeout, r.Config.TrafficOps.URL, r.Config.TrafficOps.Users.Extension, r.Config.TrafficOps.UserPassword, r.Config.TrafficOps.Users.Admin, r.Config.TrafficOps.UserPassword)
}

func (r *TCData) CreateTestInvalidServerChecks(t *testing.T) {
	toReqTimeout := time.Second * time.Duration(r.Config.Default.Session.TimeoutInSecs)

	_, _, err := TOSession.InsertServerCheckStatus(r.TestData.Serverchecks[0])
	if err == nil {
		t.Error("expected to receive error with non extension user")
	}

	r.SwitchSession(toReqTimeout, r.Config.TrafficOps.URL, r.Config.TrafficOps.Users.Admin, r.Config.TrafficOps.UserPassword, r.Config.TrafficOps.Users.Extension, r.Config.TrafficOps.UserPassword)

	invalidServerCheck := tc.ServercheckRequestNullable{
		Name:     util.StrPtr("BOGUS"),
		Value:    util.IntPtr(1),
		ID:       util.IntPtr(-1),
		HostName: util.StrPtr("bogus_hostname"),
	}

	// Attempt to create a ServerCheck with invalid server ID
	_, _, err = TOSession.InsertServerCheckStatus(invalidServerCheck)
	if err == nil {
		t.Error("expected to receive error with invalid id")
	}

	invalidServerCheck.ID = nil
	// Attempt to create a ServerCheck with invalid host name
	_, _, err = TOSession.InsertServerCheckStatus(invalidServerCheck)
	if err == nil {
		t.Error("expected to receive error with invalid host name")
	}

	// get valid name to get past host check
	invalidServerCheck.Name = r.TestData.Serverchecks[0].Name

	// Attempt to create a ServerCheck with invalid servercheck name
	_, _, err = TOSession.InsertServerCheckStatus(invalidServerCheck)
	if err == nil {
		t.Error("expected to receive error with invalid servercheck name")
	}
	r.SwitchSession(toReqTimeout, r.Config.TrafficOps.URL, r.Config.TrafficOps.Users.Extension, r.Config.TrafficOps.UserPassword, r.Config.TrafficOps.Users.Admin, r.Config.TrafficOps.UserPassword)
}

// Need to define no-op function as TCObj interface expects a delete function
// There is no delete path for serverchecks
func (r *TCData) DeleteTestServerChecks(t *testing.T) {
	return
}
