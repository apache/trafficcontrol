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
	"os"
	"path/filepath"
	"testing"

	"github.com/apache/trafficcontrol/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/cache-config/testing/ort-tests/tcdata"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

func TestT3cUnsetsUpdateFlag(t *testing.T) {
	fmt.Println("------------- Starting TestT3cUnsetsUpdateFlag tests ---------------")
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters, tcdata.Statuses,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices}, func() {

		const cacheHostName = `atlanta-edge-03`
		const cmdUpdateStatus = `update-status`

		t.Logf("DEBUG TestT3cReload calling badass")
		if stdOut, exitCode := t3cUpdateUnsetFlag(cacheHostName, "badass"); exitCode != 0 {
			t.Fatalf("ERROR: t3c badass failed: code '%v' output '%v'\n", exitCode, stdOut)
		}

		// delete a file that we know should trigger a reload.
		fileNameToRemove := filepath.Join(test_config_dir, "hdr_rw_first_ds-top.config")
		if err := os.Remove(fileNameToRemove); err != nil {
			t.Fatalf("failed to remove file '" + fileNameToRemove + "': " + err.Error())
		}

		t.Logf("DEBUG TestT3cReload setting upate flag")
		// set the update flag, so syncds will run
		if err := ExecTOUpdater(cacheHostName, false, true); err != nil {
			t.Fatalf("t3c-update failed: %v\n", err)
		}

		{
			// verify update status is now true

			output, err := runRequest(cacheHostName, cmdUpdateStatus)
			if err != nil {
				t.Fatalf("ERROR: to_requester run failed: %v\n", err)
			}
			serverStatus := tc.ServerUpdateStatus{}
			if err = json.Unmarshal([]byte(output), &serverStatus); err != nil {
				t.Fatalf("ERROR unmarshalling json output: " + err.Error())
			}
			if serverStatus.HostName != cacheHostName {
				t.Fatalf("expected request update-status host '%v' actual %v", cacheHostName, serverStatus.HostName)
			} else if serverStatus.RevalPending {
				t.Fatal("expected RevalPending false after update")
			} else if !serverStatus.UpdatePending {
				t.Fatal("expected UpdatePending true after update")
			}
		}

		_, _ = t3cUpdateUnsetFlag(cacheHostName, "syncds")
		// Ignore the exit code error for now, because the ORT Integration Test Framework doesn't currently start ATS.
		// TODO check err, after running ATS is added to the tests.
		// if err != nil {
		// 	t.Fatalf("t3c syncds failed: %v\n", err)
		// }

		{
			// verify update status after syncds is now false

			output, err := runRequest(cacheHostName, cmdUpdateStatus)
			if err != nil {
				t.Fatalf("t3c-request failed: %v\n", err)
			}
			serverStatus := tc.ServerUpdateStatus{}
			if err = json.Unmarshal([]byte(output), &serverStatus); err != nil {
				t.Fatalf("unmarshalling request update-status json: " + err.Error())
			}
			if serverStatus.HostName != cacheHostName {
				t.Errorf("expected update-status host '%v' actual %v", cacheHostName, serverStatus.HostName)
			} else if serverStatus.RevalPending {
				t.Error("expected RevalPending false after syncds run")
			} else if serverStatus.UpdatePending {
				t.Error("expected UpdatePending false after syncds run")
			}
		}
	})
	fmt.Println("------------- End of TestT3cTOUpdates tests ---------------")
}

func t3cUpdateUnsetFlag(host string, runMode string) (string, int) {
	args := []string{
		"apply",
		"--traffic-ops-insecure=true",
		"--dispersion=0",
		"--login-dispersion=0",
		"--traffic-ops-timeout-milliseconds=3000",
		"--traffic-ops-user=" + tcd.Config.TrafficOps.Users.Admin,
		"--traffic-ops-password=" + tcd.Config.TrafficOps.UserPassword,
		"--traffic-ops-url=" + tcd.Config.TrafficOps.URL,
		"--cache-host-name=" + host,
		"--log-location-error=stdout",
		"--log-location-info=stdout",
		"--log-location-debug=test.log",
		"--omit-via-string-release=true",
		"--git=no",
		"--run-mode=" + runMode,
	}
	stdOut, _, exitCode := t3cutil.Do("t3c", args...) // should be no stderr, we told it to log to stdout
	return string(stdOut), exitCode
}
