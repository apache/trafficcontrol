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
	"os/exec"
	"testing"

	"github.com/apache/trafficcontrol/v8/cache-config/testing/ort-tests/tcdata"
)

func TestT3cTail(t *testing.T) {
	// t3c must create semantically blank files. Failing to do so will cause other config files that reference them to fail.
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices}, func() {

		t.Run("confirm default flags run t3c-tail ATS restart confirmation", func(t *testing.T) {
			stdErr, err := t3cUpdateWithTail(DefaultCacheHostName, "badass", true)
			if err != nil {
				t.Fatalf("t3c badass failed: %v", err)
			}

			if !bytes.Contains(stdErr, []byte(`confirming ATS restart succeeded`)) {
				t.Errorf("expected t3c log to have confirmed ATS restart (t3c-tail), actual: %v", string(stdErr))
			}

			if bytes.Contains(stdErr, []byte(`skipping ATS restart success confirmation`)) {
				t.Errorf("expected t3c log to not have skipped ATS restart (t3c-tail), actual: %v", string(stdErr))
			}

		})

		t.Run("confirm t3c-apply --no-confirm-service-action does not t3c-tail for ATS restart confirmation", func(t *testing.T) {
			stdErr, err := t3cUpdateWithTail(DefaultCacheHostName, "badass", false)
			if err != nil {
				t.Fatalf("t3c badass failed: %v", err)
			}

			if bytes.Contains(stdErr, []byte(`confirming ATS restart succeeded`)) {
				t.Errorf("expected t3c --no-confirm-service-action flag log to not confirm ATS restart (t3c-tail), actual: %v", string(stdErr))
			}

			if !bytes.Contains(stdErr, []byte(`skipping ATS restart success confirmation`)) {
				t.Errorf("expected t3c --no-confirm-service-action flag to log not confirming ATS restart (t3c-tail), actual: %v", string(stdErr))
			}

		})

	})
}

func t3cUpdateWithTail(host string, runMode string, withTail bool) ([]byte, error) {
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
		"--run-mode=" + runMode,
		"--git=no",
	}

	if !withTail {
		args = append(args, "--no-confirm-service-action")
	}

	cmd := exec.Command("t3c", args...)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err := cmd.Run()
	if err != nil {
		return nil, errors.New(err.Error() + ": " + "stdout: " + out.String() + " stderr: " + errOut.String())
	}
	return errOut.Bytes(), nil
}
