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
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/cache-config/testing/ort-tests/tcdata"
	testutil "github.com/apache/trafficcontrol/cache-config/testing/ort-tests/util"
)

func TestIMS(t *testing.T) {
	t.Logf("------------- Starting TestIMS ---------------")
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters, tcdata.Statuses,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices}, func() {

		doTestIMS(t)

	})
	t.Logf("------------- End of TestIMS ---------------")
}

func doTestIMS(t *testing.T) {
	t.Logf("------------- Start doTestIMS ---------------")

	cacheHostName := "atlanta-edge-03"

	t.Logf("doTestIMS calling badass with cache")
	if stdOut, exitCode := t3cApplyCache(cacheHostName, false); exitCode != 0 {
		t.Fatalf("ERROR: t3c badass failed: code '%v' output '%v'\n", exitCode, stdOut)
	}

	if !testutil.FileExists(t3cutil.ApplyCachePath) {
		t.Fatalf("expected: cache '%v' to exist after badass, actual: doesn't exist", t3cutil.ApplyCachePath)
	}

	if stdOut, exitCode := t3cApplyCache(cacheHostName, false); exitCode != 0 {
		t.Fatalf("ERROR: t3c badass failed: code '%v' output '%v'\n", exitCode, stdOut)
	} else if !strings.Contains(stdOut, "not modified, using old config") {
		t.Errorf("ERROR: expected t3c second badass to have a successful IMS 304, actual: code '%v' output '%v'\n", exitCode, stdOut)
	}

	if stdOut, exitCode := t3cApplyCache(cacheHostName, true); exitCode != 0 {
		t.Fatalf("ERROR: t3c badass failed: code '%v' output '%v'\n", exitCode, stdOut)
	} else if strings.Contains(stdOut, "not modified, using old config") {
		t.Errorf("ERROR: expected t3c second badass with --no-cache to not use the cache, actual: code '%v' output '%v'\n", exitCode, stdOut)
	}

	t.Logf("------------- End doTestIMS ---------------")
}

func t3cApplyCache(host string, noCache bool) (string, int) {
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
		"--omit-via-string-release=true",
		"--git=" + "yes",
		"--run-mode=badass",
	}
	if noCache {
		args = append(args, `--no-cache=true`)
	}
	_, stdErr, exitCode := t3cutil.Do("t3c", args...) // should be no stdout, we told it to log to stderr
	return string(stdErr), exitCode
}
