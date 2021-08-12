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
	"path/filepath"
	"strconv"
	"testing"

	"github.com/apache/trafficcontrol/cache-config/testing/ort-tests/tcdata"
	"github.com/apache/trafficcontrol/lib/go-util"
)

func TestT3CNoOutgoingIP(t *testing.T) {
	t.Logf("------------- Starting TestT3CNoOutgoingIP ---------------")
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters, tcdata.Statuses,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices}, func() {

		{
			if err := t3cUpdateNoOutgoingIP("atlanta-edge-03", nil); err != nil {
				t.Fatalf("ERROR: t3c badass failed: %v\n", err)
			}

			recordsName := filepath.Join(test_config_dir, "records.config")
			recordsDotConfig, err := ioutil.ReadFile(recordsName)
			if err != nil {
				t.Fatalf("reading %v: %v\n", recordsName, err)
			}

			if !bytes.Contains(recordsDotConfig, []byte("outgoing_ip_to_bind")) {
				t.Errorf("expected t3c with no --no-outgoing-ip to add records.config outgoing_ip_to_bind, actual: '%v'\n", string(recordsDotConfig))
			}
		}

		{
			if err := t3cUpdateNoOutgoingIP("atlanta-edge-03", util.BoolPtr(false)); err != nil {
				t.Fatalf("ERROR: t3c badass failed: %v\n", err)
			}

			recordsName := filepath.Join(test_config_dir, "records.config")
			recordsDotConfig, err := ioutil.ReadFile(recordsName)
			if err != nil {
				t.Fatalf("reading %v: %v\n", recordsName, err)
			}

			if !bytes.Contains(recordsDotConfig, []byte("outgoing_ip_to_bind")) {
				t.Errorf("expected t3c with --no-outgoing-ip=false to add records.config outgoing_ip_to_bind, actual: '%v'\n", string(recordsDotConfig))
			}
		}

		{
			if err := t3cUpdateNoOutgoingIP("atlanta-edge-03", util.BoolPtr(true)); err != nil {
				t.Fatalf("ERROR: t3c badass failed: %v\n", err)
			}

			recordsName := filepath.Join(test_config_dir, "records.config")
			recordsDotConfig, err := ioutil.ReadFile(recordsName)
			if err != nil {
				t.Fatalf("reading %v: %v\n", recordsName, err)
			}

			if bytes.Contains(recordsDotConfig, []byte("outgoing_ip_to_bind")) {
				t.Errorf("expected t3c with --no-outgoing-ip=true to not add records.config outgoing_ip_to_bind, actual: '%v'\n", string(recordsDotConfig))
			}
		}
	})
	t.Logf("------------- End of TestT3CNoOutgoingIP ---------------")
}

func t3cUpdateNoOutgoingIP(host string, noOutgoingIP *bool) error {
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
		"-vv",
		"--run-mode=" + "badass",
		"--git=no",
	}
	if noOutgoingIP != nil {
		args = append(args, "--no-outgoing-ip="+strconv.FormatBool(*noOutgoingIP))
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
