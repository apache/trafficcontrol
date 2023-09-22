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
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/cache-config/testing/ort-tests/tcdata"
	"github.com/apache/trafficcontrol/v8/lib/go-atscfg"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
)

func TestTOUpdater(t *testing.T) {
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices, tcdata.Jobs}, func() {

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
			t.Fatalf("t3c-request failed: %v", err)
		}
		err = json.Unmarshal([]byte(output), &serverStatus)
		if err != nil {
			t.Fatalf("failed to parse t3c-request output: %v", err)
		}
		if serverStatus.RevalPending != true { // should be true since invalidation jobs were queued
			t.Fatal("expected RevalPending to be 'true'")
		}
		if serverStatus.UpdatePending != true {
			t.Fatal("expected UpdatePending to be 'true'")
		}

		// set config apply time to the config update time to signal the update was applied
		err = ExecTOUpdater(DefaultCacheHostName, serverStatus.ConfigUpdateTime, nil, util.BoolPtr(false), nil)
		if err != nil {
			t.Fatalf("t3c-update failed: %v", err)
		}
		// verify the update status is now 'false'
		output, err = runRequest(DefaultCacheHostName, CMDUpdateStatus)
		if err != nil {
			t.Fatalf("t3c-request failed: %v", err)
		}
		err = json.Unmarshal([]byte(output), &serverStatus)
		if err != nil {
			t.Fatalf("failed to parse t3c-request output: %v", err)
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

		// now change the reval stat and put server update status back
		// set config apply time to the config update time to signal the update was applied
		err = ExecTOUpdater(DefaultCacheHostName, nil, serverStatus.RevalidateUpdateTime, nil, util.BoolPtr(false))
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
		if serverStatus.RevalPending != false {
			t.Fatal("expected RevalPending to be 'false'")
		}
		if serverStatus.UpdatePending != false {
			t.Fatal("expected UpdatePending to be 'false'")
		}
		if serverStatus.RevalidateApplyTime != nil && serverStatus.RevalidateUpdateTime != nil {
			if !(*serverStatus.RevalidateApplyTime).Equal(*serverStatus.RevalidateUpdateTime) {
				t.Fatalf("failed to set reval apply time.\nExpc: %v\nRecv: %v", *serverStatus.RevalidateUpdateTime, *serverStatus.RevalidateApplyTime)
			}
		}

	})
}

func ExecTOUpdater(host string, configApplyTime, revalApplyTime *time.Time, configApplyBool, revalApplyBool *bool) error {
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

	// *** Compatability requirement until ATC (v7.0+) is deployed with the timestamp features
	if configApplyBool != nil {
		args = append(args, "--set-update-status="+strconv.FormatBool(*configApplyBool))
	}
	if revalApplyBool != nil {
		args = append(args, "--set-reval-status="+strconv.FormatBool(*revalApplyBool))
	}
	// ***

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
