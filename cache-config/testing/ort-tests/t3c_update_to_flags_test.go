package orttest

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
	"fmt"
	"testing"

	"github.com/apache/trafficcontrol/cache-config/testing/ort-tests/tcdata"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

func TestT3cTOUpdates(t *testing.T) {
	fmt.Println("------------- Starting TestT3cTOUpdates tests ---------------")
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters, tcdata.Statuses,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices}, func() {

		// retrieve the current server status
		output, err := runRequest("atlanta-edge-03", "update-status")
		if err != nil {
			t.Fatalf("ERROR: to_requester run failed: %v\n", err)
		}
		var serverStatus tc.ServerUpdateStatus
		err = json.Unmarshal([]byte(output), &serverStatus)
		if err != nil {
			t.Fatalf("ERROR unmarshalling json output: " + err.Error())
		}
		if serverStatus.HostName != "atlanta-edge-03" {
			t.Fatal("ERROR unexpected result, expected 'atlanta-edge-03'")
		}
		if serverStatus.RevalPending != false {
			t.Fatal("ERROR unexpected result, expected RevalPending is 'false'")
		}
		if serverStatus.UpdatePending != false {
			t.Fatal("ERROR unexpected result, expected UpdatePending is 'false'")
		}

		// change the server update status
		err = ExecTOUpdater("atlanta-edge-03", true, true)
		if err != nil {
			t.Fatalf("ERROR: to_updater run failed: %v\n", err)
		}
		// verify the update status is now 'true'
		output, err = runRequest("atlanta-edge-03", "update-status")
		if err != nil {
			t.Fatalf("ERROR: to_requester run failed: %v\n", err)
		}
		err = json.Unmarshal([]byte(output), &serverStatus)
		if err != nil {
			t.Fatalf("ERROR unmarshalling json output: " + err.Error())
		}
		if serverStatus.RevalPending != true {
			t.Fatal("ERROR unexpected result, expected RevalPending is 'true'")
		}
		if serverStatus.UpdatePending != true {
			t.Fatal("ERROR unexpected result, expected UpdatePending is 'true'")
		}

		// run t3c syncds and verify only the queue update flag is reset to 'false'
		err = runApply("atlanta-edge-03", "syncds")
		if err != nil {
			t.Fatalf("ERROR: t3c syncds failed: %v\n", err)
		}
		output, err = runRequest("atlanta-edge-03", "update-status")
		if err != nil {
			t.Fatalf("ERROR: to_requester run failed: %v\n", err)
		}
		err = json.Unmarshal([]byte(output), &serverStatus)
		if err != nil {
			t.Fatalf("ERROR unmarshalling json output: " + err.Error())
		}
		if serverStatus.RevalPending != true {
			t.Fatal("ERROR unexpected result, expected RevalPending is 'true'")
		}
		if serverStatus.UpdatePending != false {
			t.Fatal("ERROR unexpected result, expected UpdatePending is 'false'")
		}

		// run t3c revalidate and verify only the queue update flag is still 'false'
		// and that the revalidate flag is now 'false'
		err = runApply("atlanta-edge-03", "revalidate")
		if err != nil {
			t.Fatalf("ERROR: t3c syncds failed: %v\n", err)
		}
		output, err = runRequest("atlanta-edge-03", "update-status")
		if err != nil {
			t.Fatalf("ERROR: to_requester run failed: %v\n", err)
		}
		err = json.Unmarshal([]byte(output), &serverStatus)
		if err != nil {
			t.Fatalf("ERROR unmarshalling json output: " + err.Error())
		}
		if serverStatus.RevalPending != false {
			t.Fatal("ERROR unexpected result, expected RevalPending is 'false'")
		}
		if serverStatus.UpdatePending != false {
			t.Fatal("ERROR unexpected result, expected UpdatePending is 'false'")
		}
	})
	fmt.Println("------------- End of TestT3cTOUpdates tests ---------------")
}
