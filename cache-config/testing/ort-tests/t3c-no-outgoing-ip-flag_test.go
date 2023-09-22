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
	"strconv"
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/v8/cache-config/testing/ort-tests/tcdata"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
)

const outgoingIPToBind = "outgoing_ip_to_bind"

func testNoOutgoingIPAfterUpdate(t *testing.T, noOutgoingIP *bool) {
	if err := t3cUpdateNoOutgoingIP(DefaultCacheHostName, noOutgoingIP); err != nil {
		t.Fatalf("t3c badass failed: %v", err)
	}

	recordsDotConfig, err := ioutil.ReadFile(RecordsConfigFileName)
	if err != nil {
		t.Fatalf("reading %s: %v", RecordsConfigFileName, err)
	}
	contents := string(recordsDotConfig)

	// The default behavior when --no-outgoing-ip isn't given is equivalent to
	// passing --no-outgoing-ip=false
	if noOutgoingIP == nil || !*noOutgoingIP {
		if !strings.Contains(contents, outgoingIPToBind) {
			t.Errorf("expected t3c to add records.config outgoing_ip_to_bind, actual: %s", contents)
		}
	} else if strings.Contains(contents, outgoingIPToBind) {
		t.Errorf("expected t3c to not add records.config outgoing_ip_to_bind, actual: %s", contents)
	}

}

func TestT3CNoOutgoingIP(t *testing.T) {
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices}, func() {

		t.Run("not passing a no-outgoing-ip flag", func(t *testing.T) { testNoOutgoingIPAfterUpdate(t, nil) })
		t.Run("passing a no-outgoing-ip flag that's explicitly false", func(t *testing.T) { testNoOutgoingIPAfterUpdate(t, util.BoolPtr(false)) })
		t.Run("passing a no-outgoing-ip flag that's true", func(t *testing.T) { testNoOutgoingIPAfterUpdate(t, util.BoolPtr(true)) })
	})
}

func t3cUpdateNoOutgoingIP(host string, noOutgoingIP *bool) error {
	args := []string{
		"apply",
		"--no-confirm-service-action",
		"--traffic-ops-insecure=true",
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
