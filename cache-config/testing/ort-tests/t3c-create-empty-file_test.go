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
	"io/ioutil"
	"os/exec"
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/cache-config/testing/ort-tests/tcdata"
	"github.com/apache/trafficcontrol/cache-config/testing/ort-tests/util"
)

func TestT3cCreateEmptyFile(t *testing.T) {
	// t3c must create semantically blank files. Failing to do so will cause other config files that reference them to fail.
	t.Logf("------------- Starting TestT3cCreateEmptyFile ---------------")
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters, tcdata.Statuses,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices}, func() {

		err := t3cUpdateCreateEmptyFile("atlanta-edge-03", "badass")
		if err != nil {
			t.Fatalf("ERROR: t3c badass failed: %v\n", err)
		}

		const emptyFileName = `empty-file.config`

		filePath := test_config_dir + "/" + emptyFileName

		if !util.FileExists(filePath) {
			t.Fatalf("ERROR: missing empty config file, %s,  empty files must still be created", filePath)
		}

		emptyFile, err := ioutil.ReadFile(filePath)
		if err != nil {
			t.Fatalf("reading file '%v': %v\n", filePath, err)
		}
		emptyFile = bytes.TrimSpace(emptyFile)

		lines := strings.Split(string(emptyFile), "\n")
		if len(lines) > 0 && !strings.HasPrefix(lines[0], `#`) {
			t.Errorf("expected file '%v' to be empty except for comment, actual '%v'\n", filePath, string(emptyFile))
		}
		if len(lines) > 1 {
			t.Errorf("expected file '%v' to be empty for testing, actual '%v'\n", filePath, string(emptyFile))
		}
	})
	t.Logf("------------- End of TestT3cCreateEmptyFile ---------------")
}

func t3cUpdateCreateEmptyFile(host string, run_mode string) error {
	args := []string{
		"apply",
		"--traffic-ops-insecure=true",
		"--traffic-ops-timeout-milliseconds=3000",
		"--traffic-ops-user=" + tcd.Config.TrafficOps.Users.Admin,
		"--traffic-ops-password=" + tcd.Config.TrafficOps.UserPassword,
		"--traffic-ops-url=" + tcd.Config.TrafficOps.URL,
		"--cache-host-name=" + host,
		"-vv",
		"--omit-via-string-release=true",
		"--run-mode=" + run_mode,
		"--git=no",
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
