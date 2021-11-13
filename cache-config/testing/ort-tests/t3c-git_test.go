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
	"os/exec"
	"strconv"
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/cache-config/testing/ort-tests/tcdata"
	"github.com/apache/trafficcontrol/cache-config/testing/ort-tests/util"
)

func TestT3cGit(t *testing.T) {
	t.Logf("------------- Starting TestT3cGit ---------------")
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters, tcdata.Statuses,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices}, func() {

		if err := util.RMGit(test_config_dir); err != nil {
			t.Fatalf("removing existing git directory: %v", err)
		}

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

			diffStr, err := util.DiffFiles(bfn, tfn)
			if err != nil {
				t.Fatalf("diffing %s and %s: %v", tfn, bfn, err)
			} else if diffStr != "" {
				t.Errorf("%s and %s differ: %v", tfn, bfn, diffStr)
			} else {
				t.Logf("%s and %s diff clean", tfn, bfn)
			}
		}

		gitLog, err := gitLogOneline(test_config_dir)
		if err != nil {
			t.Fatalf("getting git log: %v", err)
		}

		numCommits, err := gitNumCommits(test_config_dir)
		if err != nil {
			t.Errorf("ERROR: checking number of git commits: %v\n", err)
		} else if numCommits != 3 {
			// expecting 3 commits: initial commit, startup commit of preexisting files, and the post-run commit.
			t.Errorf("ERROR: git commits expected >=3 actual %v\n%v", numCommits, gitLog)

			for i := 0; i < numCommits; i++ {
				showTxt, err := gitShow(i, test_config_dir)
				if err != nil {
					t.Errorf("git show: %v\n", err)
				} else {
					t.Logf("git HEAD~%v: %v\n\n\n", i, showTxt)
				}
			}
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

func gitShow(n int, dir string) (string, error) {
	cmd := exec.Command("git", "show", "HEAD~"+strconv.Itoa(n))
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git error: in dir '%v' returned err %v msg '%v'", dir, err, string(output))
	}
	return strings.TrimSpace(string(output)), nil
}

func gitLogOneline(dir string) (string, error) {
	cmd := exec.Command("git", "log", "--pretty=oneline")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git error: in dir '%v' returned err %v msg '%v'", dir, err, string(output))
	}
	return strings.TrimSpace(string(output)), nil
}

func t3cUpdateGit(host string, run_mode string) error {
	args := []string{
		"apply",
		"--traffic-ops-insecure=true",
		"--traffic-ops-timeout-milliseconds=3000",
		"--traffic-ops-user=" + tcd.Config.TrafficOps.Users.Admin,
		"--traffic-ops-password=" + tcd.Config.TrafficOps.UserPassword,
		"--traffic-ops-url=" + tcd.Config.TrafficOps.URL,
		"--cache-host-name=" + host,
		"--verbose", // first verbose option to enable warnings
		"--verbose", // second verbose option to enable info
		"--omit-via-string-release=true",
		"--run-mode=" + run_mode,
		"--git=" + "yes",
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
