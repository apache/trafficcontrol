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
	"net/url"
	"path/filepath"
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/cache-config/testing/ort-tests/tcdata"
	testutil "github.com/apache/trafficcontrol/cache-config/testing/ort-tests/util"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
	toclient "github.com/apache/trafficcontrol/traffic_ops/v3-client"
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
		doTestIMSChangedCDN(t)

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

// doTestIMSChangedCDN tests that after caching, requests which use the CDN as a key don't use the invalid cache.
func doTestIMSChangedCDN(t *testing.T) {
	t.Logf("------------- Start doTestIMSChangedCDN ---------------")
	defer func() { t.Logf("------------- End doTestIMSChangedCDN ---------------") }()

	cacheHostName := "atlanta-edge-03"

	t.Logf("doTestIMSChangedCDN calling badass with cache")
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

	cdn1Domain := "test.cdn1.net"
	cdn2Domain := "test.cdn2.net"
	cdn2ProfileName := "ATS_EDGE_TIER_CACHE_CDN2"

	{
		// check that remap.config has the initial cdn1

		remapName := filepath.Join(test_config_dir, "remap.config")
		remapDotConfig, err := ioutil.ReadFile(remapName)
		if err != nil {
			t.Fatalf("reading %v: %v\n", remapName, err)
		}

		if !bytes.Contains(remapDotConfig, []byte(cdn1Domain)) {
			t.Errorf("expected remap.config to contain cdn1 domain '%v', actual: '%v'\n", cdn1Domain, string(remapDotConfig))
		}
	}

	{
		// change the server's CDN

		cdn2Name := "cdn2"
		cdns, _, err := tcdata.TOSession.GetCDNByNameWithHdr(cdn2Name, nil)
		if err != nil {
			t.Fatalf("getting cdn: " + err.Error())
		} else if len(cdns) != 1 {
			t.Fatalf("getting cdn: expected 1 cdn actual num %v cdns %+v", len(cdns), cdns)
		}

		// have to change the profile at the same time, or TO will reject the change.
		profiles, _, err := tcdata.TOSession.GetProfileByNameWithHdr(cdn2ProfileName, nil)
		if err != nil {
			t.Fatalf("getting profile: " + err.Error())
		} else if len(cdns) != 1 {
			t.Fatalf("getting profile: expected 1 cdn actual num %v objects %+v", len(profiles), profiles)
		}

		cdn2ID := cdns[0].ID
		cdn2ProfileID := profiles[0].ID

		sv, _, err := GetServer(tcdata.TOSession, cacheHostName)
		if err != nil {
			t.Fatalf("getting server: " + err.Error())
		}

		sv.CDNID = &cdn2ID
		sv.CDNName = &cdn2Name
		sv.ProfileID = &cdn2ProfileID
		sv.Profile = &cdn2ProfileName

		_, _, err = tcdata.TOSession.UpdateServerByIDWithHdr(*sv.ID, *sv, nil)
		if err != nil {
			t.Fatalf("updating server: " + err.Error())
		}
	}

	// run t3c after changing the cdn

	stdOut, exitCode := t3cApplyCache(cacheHostName, false)
	if exitCode != 0 {
		t.Fatalf("ERROR: t3c badass failed: code '%v' output '%v'\n", exitCode, stdOut)
	}

	{
		// check that remap.config has the changed cdn2, and does not have the old cdn1

		remapName := filepath.Join(test_config_dir, "remap.config")
		remapDotConfig, err := ioutil.ReadFile(remapName)
		if err != nil {
			t.Fatalf("reading %v: %v\n", remapName, err)
		}

		if !bytes.Contains(remapDotConfig, []byte(cdn2Domain)) {
			t.Errorf("expected after changing server to cdn2 for remap.config to contain cdn2 domain '%v', actual: '%v'\n", cdn2Domain, string(remapDotConfig))
		}

		if bytes.Contains(remapDotConfig, []byte(cdn1Domain)) {
			t.Errorf("expected after changing server to cdn2 for remap.config to not contain cdn1 domain '%v', actual: '%v'\n", cdn1Domain, string(remapDotConfig))
		}
	}

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
