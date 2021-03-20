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
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/traffic_ops_ort/testing/ort-tests/tcdata"
	"github.com/apache/trafficcontrol/traffic_ops_ort/testing/ort-tests/util"
)

func TestT3cGit(t *testing.T) {
	t.Logf("------------- Starting TestT3cGit ---------------")
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters, tcdata.Statuses,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices}, func() {

		// run badass and check config files.
		err := t3cUpdateGit("atlanta-edge-03", "badass")
		if err != nil {
			t.Fatalf("ERROR: t3c badass failed: %v\n", err)
		}
		for _, v := range testFiles {
			bfn := base_line_dir + "/" + v
			if !util.FileExists(bfn) {
				t.Fatalf("ERROR: missing baseline config file, %s,  needed for tests", bfn)
			}
			tfn := test_config_dir + "/" + v
			if !util.FileExists(tfn) {
				t.Fatalf("ERROR: missing the expected config file, %s", tfn)
			}

			result, err := util.DiffFiles(bfn, tfn)
			if err != nil || !result {
				t.Fatalf("ERROR: the contents of '%s' does not match those in %s",
					tfn, bfn)
			}
		}

		time.Sleep(time.Second * 5)

		t.Logf("------------------------ running SYNCDS Test ------------------")
		// remove the remap.config in preparation for running syncds
		remap := test_config_dir + "/remap.config"
		err = os.Remove(remap)
		if err != nil {
			t.Fatalf("ERROR: unable to remove %s\n", remap)
		}
		// prepare for running syncds.
		err = setQueueUpdateStatus("atlanta-edge-03", "true")
		if err != nil {
			t.Fatalf("ERROR: queue updates failed: %v\n", err)
		}

		// remap.config is removed and atlanta-edge-03 should have
		// queue updates enabled.  run t3c to verify a new remap.config
		// is pulled down.
		err = t3cUpdateGit("atlanta-edge-03", "syncds")
		if err != nil {
			t.Fatalf("ERROR: t3c syncds failed: %v\n", err)
		}
		if !util.FileExists(remap) {
			t.Fatalf("ERROR: syncds failed to pull down %s\n", remap)
		}
		t.Logf("------------------------ end SYNCDS Test ------------------")

		numCommits, err := gitNumCommits(test_config_dir)
		if err != nil {
			t.Errorf("ERROR: checking number of git commits: %v\n", err)
		} else if numCommits != 3 { // expecting 3 commits: initial commit, startup commit of preexisting files, and the post-run commit
			t.Errorf("ERROR: git commits expected %v actual %v\n", 3, numCommits)
		}

	})
	t.Logf("------------- End of TestT3cGit ---------------")
}

func gitNumCommits(dir string) (int, error) {
	cmd := exec.Command("git", "rev-list", "--count", "HEAD")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("git error: in dir '%v' returned err %v msg '%v'", dir, err, string(output))
	}
	numChanges, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		return 0, fmt.Errorf("git error: in dir '%v' expected number, but got '%v'", dir, string(output))
	}
	return numChanges, nil
}

func t3cUpdateGit(host string, run_mode string) error {
	args := []string{
		"--traffic-ops-insecure=true",
		"--dispersion=0",
		"--login-dispersion=0",
		"--traffic-ops-timeout-milliseconds=3000",
		"--traffic-ops-user=" + tcd.Config.TrafficOps.Users.Admin,
		"--traffic-ops-password=" + tcd.Config.TrafficOps.UserPassword,
		"--traffic-ops-url=" + tcd.Config.TrafficOps.URL,
		"--cache-host-name=" + host,
		"--log-location-error=test.log",
		"--log-location-info=test.log",
		"--log-location-debug=test.log",
		"--run-mode=" + run_mode,
		"--git=" + "yes",
	}
	cmd := exec.Command("/opt/ort/t3c", args...)
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
