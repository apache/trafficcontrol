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
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/v8/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/v8/cache-config/t3cutil/toreq"
	"github.com/apache/trafficcontrol/v8/cache-config/testing/ort-tests/tcdata"
	testutil "github.com/apache/trafficcontrol/v8/cache-config/testing/ort-tests/util"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	toclient "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

func TestIMS(t *testing.T) {
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices}, func() {

		t.Run("test IMS", doTestIMS)
		t.Run("test IMS when CDN changes", doTestIMSChangedCDN)

	})
}

func doTestIMS(t *testing.T) {
	if stdOut, exitCode := t3cApplyCache(DefaultCacheHostName, false); exitCode != 0 {
		t.Fatalf("t3c badass failed with exit code %d, output: %s", exitCode, stdOut)
	}

	if !testutil.FileExists(t3cutil.ApplyCachePath) {
		t.Fatalf("expected: cache '%s' to exist after badass, actual: doesn't exist", t3cutil.ApplyCachePath)
	}

	if stdOut, exitCode := t3cApplyCache(DefaultCacheHostName, false); exitCode != 0 {
		t.Fatalf("t3c badass failed with exit code %d, output: %s", exitCode, stdOut)
	} else if !strings.Contains(stdOut, "not modified, using old config") {
		t.Errorf("expected t3c second badass to have a successful IMS 304, actual: code %d output: %s", exitCode, stdOut)
	}

	if stdOut, exitCode := t3cApplyCache(DefaultCacheHostName, true); exitCode != 0 {
		t.Fatalf("t3c badass failed with exit code %d, output: %s", exitCode, stdOut)
	} else if strings.Contains(stdOut, "not modified, using old config") {
		t.Errorf("expected t3c second badass with --no-cache to not use the cache, actual: code %d output: %s", exitCode, stdOut)
	}
}

const (
	cdn1Domain      = "test.cdn1.net"
	cdn2Domain      = "test.cdn2.net"
	cdn2ProfileName = "ATS_EDGE_TIER_CACHE_CDN2"
)

// doTestIMSChangedCDN tests that after caching, requests which use the CDN as a key don't use the invalid cache.
func doTestIMSChangedCDN(t *testing.T) {
	if stdOut, exitCode := t3cApplyCache(DefaultCacheHostName, false); exitCode != 0 {
		t.Fatalf("t3c badass failed with exit code %d, output: %s", exitCode, stdOut)
	}

	if !testutil.FileExists(t3cutil.ApplyCachePath) {
		t.Fatalf("expected: config data file '%s' to exist after badass, actual: doesn't exist", t3cutil.ApplyCachePath)
	}

	if stdOut, exitCode := t3cApplyCache(DefaultCacheHostName, false); exitCode != 0 {
		t.Fatalf("t3c badass failed with exit code %d, output: %s", exitCode, stdOut)
	} else if !strings.Contains(stdOut, "not modified, using old config") {
		t.Errorf("expected t3c second badass to have a successful IMS 304, actual: code %d output: %s", exitCode, stdOut)
	}

	// check that remap.config has the initial cdn1
	remapName := filepath.Join(TestConfigDir, "remap.config")
	remapDotConfig, err := ioutil.ReadFile(remapName)
	if err != nil {
		t.Fatalf("reading %s: %v", remapName, err)
	}
	contents := string(remapDotConfig)

	if !strings.Contains(contents, cdn1Domain) {
		t.Errorf("expected remap.config to contain cdn1 domain '%s', actual: '%s'", cdn1Domain, contents)
	}

	cdn2Name := "cdn2"
	cdn, _, err := toreq.GetCDNByName(tcdata.TOSession, tc.CDNName(cdn2Name), nil)
	if err != nil {
		t.Fatalf("getting cdn: %v", err)
	}

	// have to change the profile at the same time, or TO will reject the change.

	cdn2ID := cdn.ID
	sv, _, err := toreq.GetServerByHostName(tcdata.TOSession, DefaultCacheHostName)
	if err != nil {
		t.Fatalf("getting server: %v", err)
	}

	sv.CDNID = cdn2ID
	sv.CDN = cdn2Name
	sv.Profiles = []string{cdn2ProfileName}

	_, _, err = tcdata.TOSession.UpdateServer(sv.ID, *sv, toclient.RequestOptions{})
	if err != nil {
		t.Fatalf("updating server: %v", err)
	}

	// run t3c after changing the cdn

	stdOut, exitCode := t3cApplyCache(DefaultCacheHostName, false)
	if exitCode != 0 {
		t.Fatalf("t3c badass failed with exit code %d, output: %s", exitCode, stdOut)
	}

	remapDotConfig, err = ioutil.ReadFile(remapName)
	if err != nil {
		t.Fatalf("reading %s: %v", remapName, err)
	}
	contents = string(remapDotConfig)

	if !strings.Contains(contents, cdn2Domain) {
		t.Errorf("expected after changing server to cdn2 for remap.config to contain cdn2 domain '%s', actual: '%s'", cdn2Domain, contents)
	}

	if strings.Contains(contents, cdn1Domain) {
		t.Errorf("expected after changing server to cdn2 for remap.config to not contain cdn1 domain '%s', actual: '%s'", cdn1Domain, remapDotConfig)
	}
}

func t3cApplyCache(host string, noCache bool) (string, int) {
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
