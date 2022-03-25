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
	"bytes"
	"encoding/json"
	"errors"
	"os/exec"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/cache-config/testing/ort-tests/tcdata"
	"github.com/apache/trafficcontrol/lib/go-atscfg"
)

func TestTOUpdater(t *testing.T) {
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters, tcdata.Statuses,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices}, func() {

		// retrieve the current server status
		output, err := runRequest(DefaultCacheHostName, CMDUpdateStatus)
		if err != nil {
			t.Fatalf("t3c-request failed: %v", err)
		}
		var serverStatus atscfg.ServerUpdateStatus
		err = json.Unmarshal([]byte(output), &serverStatus)
		if err != nil {
			t.Fatalf("failed to parse t3c-request output: %v", err)
		}
		if serverStatus.HostName != DefaultCacheHostName {
			t.Fatalf("expected server status host name to be '%s', actual: %s", DefaultCacheHostName, serverStatus.HostName)
		}
		if serverStatus.RevalPending != false {
			t.Fatal("expected RevalPending to be 'false'")
		}
		if serverStatus.UpdatePending != false {
			t.Fatal("expected UpdatePending to be 'false'")
		}

		// change the server update status
		// Send an apply time that is before an update time, signaling there is an update pending
		before := serverStatus.ConfigUpdateTime.Add(-time.Hour * 24)
		err = ExecTOUpdater(DefaultCacheHostName, &before, serverStatus.RevalidateUpdateTime)
		if err != nil {
			t.Fatalf("t3c-update failed: %v", err)
		}
		// verify the update status is now 'true'
		output, err = runRequest(DefaultCacheHostName, CMDUpdateStatus)
		if err != nil {
			t.Fatalf("t3c-request failed: %v", err)
		}
		err = json.Unmarshal([]byte(output), &serverStatus)
		if err != nil {
			t.Fatalf("failed to parse t3c-request output: %v", err)
		}
		if serverStatus.RevalPending != false {
			t.Fatal("expected RevalPending to be 'false'")
		}
		if serverStatus.UpdatePending != true {
			t.Fatal("expected UpdatePending to be 'true'")
		}
		if !serverStatus.ConfigApplyTime.Equal(before) {
			t.Fatalf("failed to set config apply time.\nSent: %v\nRecv: %v", before, serverStatus.ConfigApplyTime)
		}

		// now change the reval stat and put server update status back
		// Send an apply time that is before an update time, signaling there is an reval pending
		before = serverStatus.RevalidateUpdateTime.Add(-time.Hour * 24)
		err = ExecTOUpdater(DefaultCacheHostName, serverStatus.ConfigUpdateTime, &before)
		if err != nil {
			t.Fatalf("t3c-update failed: %v", err)
		}
		// verify the change
		output, err = runRequest(DefaultCacheHostName, CMDUpdateStatus)
		if err != nil {
			t.Fatalf("t3c-request failed: %v", err)
		}
		err = json.Unmarshal([]byte(output), &serverStatus)
		if err != nil {
			t.Fatalf("failed to parse t3c-request output: %v", err)
		}
		if serverStatus.RevalPending != true {
			t.Fatal("expected RevalPending to be 'false'")
		}
		if serverStatus.UpdatePending != false {
			t.Fatal("expected UpdatePending to be 'true'")
		}
		if !serverStatus.RevalidateApplyTime.Equal(before) {
			t.Fatalf("failed to set config apply time.\nSent: %v\nRecv: %v", before, serverStatus.RevalidateApplyTime)
		}

	})
}

func ExecTOUpdater(host string, configApplyTime, revalApplyTime *time.Time) error {
	args := []string{
		"update",
		"--traffic-ops-insecure=true",
		"--traffic-ops-timeout-milliseconds=3000",
		"--traffic-ops-user=" + tcd.Config.TrafficOps.Users.Admin,
		"--traffic-ops-password=" + tcd.Config.TrafficOps.UserPassword,
		"--traffic-ops-url=" + tcd.Config.TrafficOps.URL,
		"--cache-host-name=" + host,
		"-vv",
	}
	if configApplyTime != nil {
		args = append(args, "--set-config-apply-time="+(*configApplyTime).Format(time.RFC3339Nano))
	}
	if revalApplyTime != nil {
		args = append(args, "--set-reval-apply-time="+(*revalApplyTime).Format(time.RFC3339Nano))
	}
	cmd := exec.Command("t3c", args...)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err := cmd.Run()
	if err != nil {
		return errors.New(err.Error() + ": " + "stdout: " + out.String() + " stderr: " + errOut.String())
	}

	return nil
}
