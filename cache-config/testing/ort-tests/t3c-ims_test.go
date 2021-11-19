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
	"errors"
	"io/ioutil"
	"net/url"
	"path/filepath"
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/cache-config/testing/ort-tests/tcdata"
	testutil "github.com/apache/trafficcontrol/cache-config/testing/ort-tests/util"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
	toclient "github.com/apache/trafficcontrol/traffic_ops/v3-client"
)

func TestIMS(t *testing.T) {
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters, tcdata.Statuses,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices}, func() {

		t.Run("test IMS", doTestIMS)
		t.Run("test IMS when CDN changes", doTestIMSChangedCDN)

	})
}

func doTestIMS(t *testing.T) {
	if stdOut, exitCode := t3cApplyCache(cacheHostName, false); exitCode != 0 {
		t.Fatalf("t3c badass failed with exit code %d, output: %s", exitCode, stdOut)
	}

	if !testutil.FileExists(t3cutil.ApplyCachePath) {
		t.Fatalf("expected: cache '%s' to exist after badass, actual: doesn't exist", t3cutil.ApplyCachePath)
	}

	if stdOut, exitCode := t3cApplyCache(cacheHostName, false); exitCode != 0 {
		t.Fatalf("t3c badass failed with exit code %d, output: %s", exitCode, stdOut)
	} else if !strings.Contains(stdOut, "not modified, using old config") {
		t.Errorf("expected t3c second badass to have a successful IMS 304, actual: code %d output: %s", exitCode, stdOut)
	}

	if stdOut, exitCode := t3cApplyCache(cacheHostName, true); exitCode != 0 {
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

func checkRemapConfigHasInitialCDN(t *testing.T) {
	remapName := filepath.Join(test_config_dir, "remap.config")
	remapDotConfig, err := ioutil.ReadFile(remapName)
	if err != nil {
		t.Fatalf("reading %s: %v", remapName, err)
	}
	contents := string(remapDotConfig)

	if !strings.Contains(contents, cdn1Domain) {
		t.Errorf("expected remap.config to contain cdn1 domain '%s', actual: '%s'", cdn1Domain, contents)
	}

}

func changeServerCDN(t *testing.T) {
	cdn2Name := "cdn2"
	cdns, _, err := tcdata.TOSession.GetCDNByNameWithHdr(cdn2Name, nil)
	if err != nil {
		t.Fatalf("getting cdn: %v", err)
	}
	if len(cdns) != 1 {
		t.Fatalf("getting cdn: expected 1 cdn actual num %d cdns: %+v", len(cdns), cdns)
	}

	// have to change the profile at the same time, or TO will reject the change.
	profiles, _, err := tcdata.TOSession.GetProfileByNameWithHdr(cdn2ProfileName, nil)
	if err != nil {
		t.Fatalf("getting profile: %v", err)
	}
	if len(cdns) != 1 {
		t.Fatalf("getting profile: expected 1 cdn actual num %d objects: %+v", len(profiles), profiles)
	}

	cdn2ID := cdns[0].ID
	cdn2ProfileID := profiles[0].ID

	sv, _, err := GetServer(tcdata.TOSession, cacheHostName)
	if err != nil {
		t.Fatalf("getting server: %v", err)
	}

	sv.CDNID = &cdn2ID
	sv.CDNName = &cdn2Name
	sv.ProfileID = &cdn2ProfileID
	sv.Profile = util.StrPtr(cdn2ProfileName)

	_, _, err = tcdata.TOSession.UpdateServerByIDWithHdr(*sv.ID, *sv, nil)
	if err != nil {
		t.Fatalf("updating server: %v", err)
	}
}

func checkChangedCDN(t *testing.T) {
	remapName := filepath.Join(test_config_dir, "remap.config")
	remapDotConfig, err := ioutil.ReadFile(remapName)
	if err != nil {
		t.Fatalf("reading %s: %v", remapName, err)
	}
	contents := string(remapDotConfig)

	if !strings.Contains(contents, cdn2Domain) {
		t.Errorf("expected after changing server to cdn2 for remap.config to contain cdn2 domain '%s', actual: '%s'", cdn2Domain, contents)
	}

	if strings.Contains(contents, cdn1Domain) {
		t.Errorf("expected after changing server to cdn2 for remap.config to not contain cdn1 domain '%s', actual: '%s'", cdn1Domain, remapDotConfig)
	}
}

// doTestIMSChangedCDN tests that after caching, requests which use the CDN as a key don't use the invalid cache.
func doTestIMSChangedCDN(t *testing.T) {
	if stdOut, exitCode := t3cApplyCache(cacheHostName, false); exitCode != 0 {
		t.Fatalf("t3c badass failed with exit code %d, output: %s", exitCode, stdOut)
	}

	if !testutil.FileExists(t3cutil.ApplyCachePath) {
		t.Fatalf("expected: config data file '%s' to exist after badass, actual: doesn't exist", t3cutil.ApplyCachePath)
	}

	if stdOut, exitCode := t3cApplyCache(cacheHostName, false); exitCode != 0 {
		t.Fatalf("t3c badass failed with exit code %d, output: %s", exitCode, stdOut)
	} else if !strings.Contains(stdOut, "not modified, using old config") {
		t.Errorf("expected t3c second badass to have a successful IMS 304, actual: code %d output: %s", exitCode, stdOut)
	}

	t.Run("check that remap.config has the initial cdn1", checkRemapConfigHasInitialCDN)

	t.Run("change the server's CDN", changeServerCDN)

	// run t3c after changing the cdn

	stdOut, exitCode := t3cApplyCache(cacheHostName, false)
	if exitCode != 0 {
		t.Fatalf("t3c badass failed with exit code %d, output: %s", exitCode, stdOut)
	}

	t.Run("check that remap.config has the changed cdn2, and does not have the old cdn1", checkChangedCDN)
}

func t3cApplyCache(host string, noCache bool) (string, int) {
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
		"--git=" + "yes",
		"--run-mode=badass",
	}
	if noCache {
		args = append(args, `--no-cache=true`)
	}
	_, stdErr, exitCode := t3cutil.Do("t3c", args...) // should be no stdout, we told it to log to stderr
	return string(stdErr), exitCode
}

func GetServer(toClient *toclient.Session, hostName string) (*tc.ServerV30, toclientlib.ReqInf, error) {
	params := url.Values{}
	params.Add("hostName", hostName)
	resp, reqInf, err := toClient.GetServersWithHdr(&params, nil)
	if err != nil {
		return nil, reqInf, err
	}
	if len(resp.Response) == 0 {
		return nil, reqInf, errors.New("not found")
	}
	return &resp.Response[0], reqInf, nil
}
