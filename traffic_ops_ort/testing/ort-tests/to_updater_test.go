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
	"fmt"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops_ort/testing/ort-tests/tcdata"
	"os/exec"
	"strconv"
	"testing"
)

func TestTOUpdater(t *testing.T) {
	fmt.Println("------------- Starting TestTOUpdater tests ---------------")
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters, tcdata.Statuses,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices}, func() {

		// retrieve the current server status
		output, err := ExecTORequester("atlanta-edge-03", "update-status")
		if err != nil {
			t.Fatalf("ERROR: to_requester Exec failed: %v\n", err)
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
		err = ExecTOUpdater("atlanta-edge-03", false, true)
		if err != nil {
			t.Fatalf("ERROR: to_updater Exec failed: %v\n", err)
		}
		// verify the update status is now 'true'
		output, err = ExecTORequester("atlanta-edge-03", "update-status")
		if err != nil {
			t.Fatalf("ERROR: to_requester Exec failed: %v\n", err)
		}
		err = json.Unmarshal([]byte(output), &serverStatus)
		if err != nil {
			t.Fatalf("ERROR unmarshalling json output: " + err.Error())
		}
		if serverStatus.RevalPending != false {
			t.Fatal("ERROR unexpected result, expected RevalPending is 'false'")
		}
		if serverStatus.UpdatePending != true {
			t.Fatal("ERROR unexpected result, expected UpdatePending is 'true'")
		}

		// now change the reval stat and put server update status back
		err = ExecTOUpdater("atlanta-edge-03", true, false)
		if err != nil {
			t.Fatalf("ERROR: to_updater Exec failed: %v\n", err)
		}
		// verify the change
		output, err = ExecTORequester("atlanta-edge-03", "update-status")
		if err != nil {
			t.Fatalf("ERROR: to_requester Exec failed: %v\n", err)
		}
		err = json.Unmarshal([]byte(output), &serverStatus)
		if err != nil {
			t.Fatalf("ERROR unmarshalling json output: " + err.Error())
		}
		if serverStatus.RevalPending != true {
			t.Fatal("ERROR unexpected result, expected RevalPending is 'false'")
		}
		if serverStatus.UpdatePending != false {
			t.Fatal("ERROR unexpected result, expected UpdatePending is 'true'")
		}

	})
	fmt.Println("------------- End of TestTOUpdater tests ---------------")
}

func ExecTOUpdater(host string, reval_status bool, update_status bool) error {
	args := []string{
		"--traffic-ops-insecure=true",
		"--login-dispersion=0",
		"--traffic-ops-timeout-milliseconds=3000",
		"--traffic-ops-user=" + tcd.Config.TrafficOps.Users.Admin,
		"--traffic-ops-password=" + tcd.Config.TrafficOps.UserPassword,
		"--traffic-ops-url=" + tcd.Config.TrafficOps.URL,
		"--cache-host-name=" + host,
		"--log-location-error=test.log",
		"--log-location-info=test.log",
		"--log-location-debug=test.log",
		"--set-reval-status=" + strconv.FormatBool(reval_status),
		"--set-update-status=" + strconv.FormatBool(update_status),
	}
	cmd := exec.Command("/opt/ort/to_updater", args...)
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
