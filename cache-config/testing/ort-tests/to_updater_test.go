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

	"github.com/apache/trafficcontrol/cache-config/testing/ort-tests/tcdata"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

func TestTOUpdater(t *testing.T) {
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters, tcdata.Statuses,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices}, func() {

		// retrieve the current server status
		output, err := runRequest(cacheHostName, "update-status")
		if err != nil {
			t.Fatalf("t3c-request failed: %v", err)
		}
		var serverStatus tc.ServerUpdateStatus
		err = json.Unmarshal([]byte(output), &serverStatus)
		if err != nil {
			t.Fatalf("failed to parse t3c-request output: %v", err)
		}
		if serverStatus.HostName != cacheHostName {
			t.Fatalf("expected server status host name to be '%s', actual: %s", cacheHostName, serverStatus.HostName)
		}
		if serverStatus.RevalPending != false {
			t.Fatal("expected RevalPending to be 'false'")
		}
		if serverStatus.UpdatePending != false {
			t.Fatal("expected UpdatePending to be 'false'")
		}

		// change the server update status
		err = ExecTOUpdater(cacheHostName, false, true)
		if err != nil {
			t.Fatalf("t3c-update failed: %v", err)
		}
		// verify the update status is now 'true'
		output, err = runRequest(cacheHostName, "update-status")
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

		// now change the reval stat and put server update status back
		err = ExecTOUpdater(cacheHostName, true, false)
		if err != nil {
			t.Fatalf("t3c-update failed: %v", err)
		}
		// verify the change
		output, err = runRequest(cacheHostName, "update-status")
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

	})
}

func ExecTOUpdater(host string, reval_status bool, update_status bool) error {
	args := []string{
		"update",
		"--traffic-ops-insecure=true",
		"--traffic-ops-timeout-milliseconds=3000",
		"--traffic-ops-user=" + tcd.Config.TrafficOps.Users.Admin,
		"--traffic-ops-password=" + tcd.Config.TrafficOps.UserPassword,
		"--traffic-ops-url=" + tcd.Config.TrafficOps.URL,
		"--cache-host-name=" + host,
		"-vv",
		"--set-reval-status=" + strconv.FormatBool(reval_status),
		"--set-update-status=" + strconv.FormatBool(update_status),
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
