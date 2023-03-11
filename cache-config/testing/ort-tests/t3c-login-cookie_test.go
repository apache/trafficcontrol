package orttest

import (
	"github.com/apache/trafficcontrol/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/cache-config/testing/ort-tests/tcdata"
	"github.com/apache/trafficcontrol/cache-config/testing/ort-tests/util"
	"os"
	"strings"
	"testing"
)

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

const CookieDir = `/var/lib/trafficcontrol-cache-config/`

func TestT3cCookie(t *testing.T) {
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters, tcdata.Statuses,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices, tcdata.InvalidationJobs}, func() {

		t.Run("Use cookie for admin", doTestT3cCookieAdmin)
		t.Run("use cookie for any other", doTestT3cCookieOther)
	})
}

func doTestT3cCookieAdmin(t *testing.T) {
	cookieFile := CookieDir + tcd.Config.TrafficOps.Users.Admin + ".cookie"
	var cookieTimeStamp *int64
	if util.FileExists(cookieFile) {
		fileInfo, err := os.Stat(cookieFile)
		if err == nil {
			t := fileInfo.ModTime().UTC().Unix()
			cookieTimeStamp = &t
		}
	}
	stdOut, exitCode := t3cLoginCookie(DefaultCacheHostName, tcd.Config.TrafficOps.Users.Admin)
	if exitCode != 0 {
		t.Fatalf("t3c update failed with exitcode %d output: %s", exitCode, stdOut)
	}
	if cookieTimeStamp != nil && strings.Contains(stdOut, "with Cookie") {
		if newCookie, err := os.Stat(cookieFile); err == nil {
			if newCookie.ModTime().UTC().Unix() == *cookieTimeStamp {
				t.Error("cookie file timestamps match, expected new cookie file to be written after login")
			}
		}
	} else if cookieTimeStamp == nil && strings.Contains(stdOut, "Error retrieving cookie") {
		if !util.FileExists(cookieFile) {
			t.Error("Cookie file didn't exist, expected new cookie to be written.")
		}
	}
}

func doTestT3cCookieOther(t *testing.T) {
	stdOut, exitCode := t3cLoginCookie(DefaultCacheHostName, tcd.Config.TrafficOps.Users.Operations)
	if exitCode != 0 {
		t.Fatalf("t3c update failed with exitcode %d output: %s", exitCode, stdOut)
	}
	if !strings.Contains(stdOut, "Error retrieving cookie") {
		t.Error("found existing cookie, expected error retrieving cookie")
	}
}

func t3cLoginCookie(host string, user string) (string, int) {
	args := []string{
		"update",
		"--traffic-ops-insecure=true",
		"--traffic-ops-timeout-milliseconds=3000",
		"--traffic-ops-user=" + user,
		"--traffic-ops-password=" + tcd.Config.TrafficOps.UserPassword,
		"--traffic-ops-url=" + tcd.Config.TrafficOps.URL,
		"--cache-host-name=" + host,
		"-vv",
		"--set-reval-status=false",
		"--set-update-status=false",
	}
	_, stdErr, exitCode := t3cutil.Do("t3c", args...)
	return string(stdErr), exitCode
}
