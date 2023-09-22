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
	"testing"

	"github.com/apache/trafficcontrol/v8/cache-config/testing/ort-tests/tcdata"
	"github.com/apache/trafficcontrol/v8/lib/go-atscfg"
)

func TestT3cTOUpdates(t *testing.T) {
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices, tcdata.Jobs}, func() {

		// retrieve the current server status
		output, err := runRequest(DefaultCacheHostName, CMDUpdateStatus)
		if err != nil {
			t.Fatalf("update-status run failed: %v", err)
		}
		var serverStatus atscfg.ServerUpdateStatus
		err = json.Unmarshal([]byte(output), &serverStatus)
		if err != nil {
			t.Fatalf("unmarshalling json output: %v", err)
		}
		if serverStatus.HostName != DefaultCacheHostName {
			t.Fatalf("incorrect server status hostname, expected '%s', got '%s'", DefaultCacheHostName, serverStatus.HostName)
		}
		if serverStatus.RevalPending != true { // should be true since invalidation jobs were queued
			t.Fatal("expected RevalPending to be 'true'")
		}
		if serverStatus.UpdatePending != false {
			t.Fatal("expected UpdatePending to be 'false'")
		}

		// change the server update status
		err = tcd.QueueUpdatesForServer(DefaultCacheHostName, true)
		if err != nil {
			t.Fatalf("failed to set config update: %v", err)
		}
		// verify the update status is now 'true'
		output, err = runRequest(DefaultCacheHostName, CMDUpdateStatus)
		if err != nil {
			t.Fatalf("update-status run failed: %v", err)
		}
		err = json.Unmarshal([]byte(output), &serverStatus)
		if err != nil {
			t.Fatalf("failed to parse update-status run output: %v", err)
		}
		if serverStatus.RevalPending != true { // should be true since invalidation jobs were queued
			t.Fatal("expected RevalPending to be 'true'")
		}
		if serverStatus.UpdatePending != true {
			t.Fatal("expected UpdatePending to be 'true'")
		}

		// run t3c syncds and verify only the queue update flag is reset to 'false'
		err = runApply(DefaultCacheHostName, "syncds")
		if err != nil {
			t.Fatalf("t3c syncds failed: %v", err)
		}
		output, err = runRequest(DefaultCacheHostName, CMDUpdateStatus)
		if err != nil {
			t.Fatalf("update-status run failed: %v", err)
		}
		err = json.Unmarshal([]byte(output), &serverStatus)
		if err != nil {
			t.Fatalf("failed to parse update-status run output: %v", err)
		}
		if serverStatus.RevalPending != true { // should be true since invalidation jobs were queued
			t.Fatal("expected RevalPending to be 'true'")
		}
		if serverStatus.UpdatePending != false {
			t.Fatal("expected UpdatePending to be 'false'")
		}
		if serverStatus.ConfigApplyTime != nil && serverStatus.ConfigUpdateTime != nil {
			if !(*serverStatus.ConfigApplyTime).Equal(*serverStatus.ConfigUpdateTime) {
				t.Fatalf("failed to set config apply time.\nExpc: %v\nRecv: %v", *serverStatus.ConfigUpdateTime, *serverStatus.ConfigApplyTime)
			}
		}

		// run t3c revalidate and verify only the queue update flag is still 'false'
		// and that the revalidate flag is now 'false'
		err = runApply(DefaultCacheHostName, "revalidate")
		if err != nil {
			t.Fatalf("t3c syncds failed: %v", err)
		}
		output, err = runRequest(DefaultCacheHostName, CMDUpdateStatus)
		if err != nil {
			t.Fatalf("update-status run failed: %v", err)
		}
		err = json.Unmarshal([]byte(output), &serverStatus)
		if err != nil {
			t.Fatalf("failed to parse update-status run output: %v", err)
		}
		if serverStatus.RevalPending != false {
			t.Error("expected RevalPending to be 'false'")
		}
		if serverStatus.UpdatePending != false {
			t.Error("expected UpdatePending to be 'false'")
		}
		if serverStatus.RevalidateApplyTime != nil && serverStatus.RevalidateUpdateTime != nil {
			if !(*serverStatus.RevalidateApplyTime).Equal(*serverStatus.RevalidateUpdateTime) {
				t.Fatalf("failed to set reval apply time.\nExpc: %v\nRecv: %v", *serverStatus.RevalidateUpdateTime, *serverStatus.RevalidateApplyTime)
			}
		}
	})
}
