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

	"github.com/apache/trafficcontrol/v8/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/v8/cache-config/testing/ort-tests/tcdata"
)

func TestT3cApplyOSHostnameShort(t *testing.T) {
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices}, func() {
		t.Run("test t3c apply getting short hostname from the OS", doTestT3cApplyOSHostnameShort)
	})
}

func doTestT3cApplyOSHostnameShort(t *testing.T) {
	// verifies that when no hostname is passed,
	// t3c-apply will get the OS hostname,
	// and use the short name (which is what will be in Traffic Ops),
	// and not the full FQDN.

	startingHost, errCode := getHostName()
	if errCode != 0 {
		t.Fatalf("getting the hostname failed with exit code %d, message: %s", errCode, startingHost)
	}

	t.Logf("original host is %s", startingHost)

	if msg, code := setHostName(DefaultCacheHostName); code != 0 {
		t.Fatalf("setting the hostname failed with exit code %d, message: %s", code, msg)
	}

	defer func() {
		if msg, code := setHostName(startingHost); code != 0 {
			t.Fatalf("setting the hostname back to the original '%s' failed with exit code %d message %s", startingHost, code, msg)
		} else {
			t.Logf("set hostname back to original '%s'", startingHost)
		}
	}()

	// verify the host was really set
	if newHost, errCode := getHostName(); errCode != 0 {
		t.Fatalf("getting the hostname failed with exit code %d, message: %s", errCode, newHost)
	} else if newHost != DefaultCacheHostName {
		t.Fatalf("setting hostname claimed it succeeded, but was '%s' expected '%s'", newHost, DefaultCacheHostName)
	} else {
		t.Logf("set hostname to '%s'", newHost)
	}

	t.Logf("calling t3c-apply with no host flag, with a short hostname")
	if stdOut, exitCode := t3cApplyNoHost(); exitCode != 0 {
		t.Fatalf("t3c-apply with no hostname arg and and system hostname '%s' failed with exit code %d, output: %s", DefaultCacheHostName, exitCode, stdOut)
	}

	fqdnHostName := DefaultCacheHostName + ".fqdn.example.test"
	if msg, code := setHostName(fqdnHostName); code != 0 {
		t.Fatalf("setting the hostname failed with exit code %d, message: %s", code, msg)
	}

	// verify the host was really set
	if newHost, errCode := getHostName(); errCode != 0 {
		t.Fatalf("getting the hostname failed with exit code %d, message: %s", errCode, newHost)
	} else if newHost != fqdnHostName {
		t.Fatalf("setting hostname claimed it succeeded, but was '%s' expected '%s'", newHost, fqdnHostName)
	} else {
		t.Logf("set hostname to '%s'", newHost)
	}

	t.Logf("calling t3c-apply with no host flag, with a fqdn hostname")
	if stdOut, exitCode := t3cApplyNoHost(); exitCode != 0 {
		t.Fatalf("t3c-apply with no hostname arg and and system hostname '%s' failed with exit code %d, output: %s", DefaultCacheHostName, exitCode, stdOut)
	}
}

func setHostName(host string) (string, int) {
	stdOut, stdErr, exitCode := t3cutil.Do("hostname", host)
	return "out: " + string(stdOut) + " err: " + string(stdErr), exitCode
}

func getHostName() (string, int) {
	stdOut, stdErr, exitCode := t3cutil.Do("hostname")
	if exitCode == 0 {
		return strings.TrimSpace(string(stdOut)), 0
	}
	return "out: " + string(stdOut) + " err: " + string(stdErr), exitCode
}

func t3cApplyNoHost() (string, int) {
	args := []string{
		"apply",
		"--no-confirm-service-action",
		"--traffic-ops-insecure=true",
		"--traffic-ops-timeout-milliseconds=3000",
		"--traffic-ops-user=" + tcd.Config.TrafficOps.Users.Admin,
		"--traffic-ops-password=" + tcd.Config.TrafficOps.UserPassword,
		"--traffic-ops-url=" + tcd.Config.TrafficOps.URL,
		"-vv",
		"--omit-via-string-release=true",
		"--git=" + "yes",
		"--run-mode=" + "badass",
	}
	_, stdErr, exitCode := t3cutil.Do("t3c", args...)
	return string(stdErr), exitCode
}
