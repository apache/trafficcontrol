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
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/cache-config/testing/ort-tests/tcdata"
)

func TestT3cReload(t *testing.T) {
	t.Logf("------------- Starting TestT3cReload ---------------")
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters, tcdata.Statuses,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices}, func() {

		doTestT3cReloadHeaderRewrite(t)
		doTestT3cReloadAnythingInTrafficserverDir(t)
		doTestT3cReloadNoChange(t)
		doTestT3cRevalCallsReload(t)

	})
	t.Logf("------------- End of TestT3cReload ---------------")
}

func doTestT3cReloadHeaderRewrite(t *testing.T) {
	t.Logf("------------- Start doTestT3cReloadHeaderRewrite ---------------")

	cacheHostName := "atlanta-edge-03"

	t.Logf("DEBUG TestT3cReload calling badass")
	if stdOut, exitCode := t3cUpdateReload(cacheHostName, "badass"); exitCode != 0 {
		t.Fatalf("ERROR: t3c badass failed: code '%v' output '%v'\n", exitCode, stdOut)
	}

	t.Logf("DEBUG TestT3cReload deleting file")

	// delete a file that we know should trigger a reload.
	fileNameToRemove := filepath.Join(test_config_dir, "hdr_rw_first_ds-top.config")
	if err := os.Remove(fileNameToRemove); err != nil {
		t.Fatalf("failed to remove file '" + fileNameToRemove + "': " + err.Error())
	}

	t.Logf("DEBUG TestT3cReload setting upate flag")
	// set the update flag, so syncds will run
	if err := ExecTOUpdater("atlanta-edge-03", false, true); err != nil {
		t.Fatalf("t3c-update failed: %v\n", err)
	}

	t.Logf("DEBUG TestT3cReload calling syncds")
	stdOut, _ := t3cUpdateReload(cacheHostName, "syncds")
	// Ignore the exit code error for now, because the ORT Integration Test Framework doesn't currently start ATS.
	// TODO check err, after running ATS is added to the tests.
	// if err != nil {
	// 	t.Fatalf("t3c syncds failed: %v\n", err)
	// }

	t.Logf("DEBUG TestT3cReload looking for reload string")
	if !strings.Contains(stdOut, `Running 'traffic_ctl config reload' now`) {
		t.Errorf("expected t3c to do a reload after adding a header rewrite file, actual: '''%v'''\n", stdOut)
	}

	t.Logf("------------- End TestT3cReload doTestT3cReloadHeaderRewrite ---------------")
}

func doTestT3cReloadAnythingInTrafficserverDir(t *testing.T) {
	t.Logf("------------- Start doTestT3cReloadAnythingInTrafficserverDir ---------------")

	cacheHostName := "atlanta-edge-03"

	t.Logf("DEBUG doTestT3cReloadHeaderRewrite calling badass")
	if stdOut, exitCode := t3cUpdateReload(cacheHostName, "badass"); exitCode != 0 {
		t.Fatalf("ERROR: t3c badass failed: code '%v' output '%v'\n", exitCode, stdOut)
	}

	t.Logf("DEBUG doTestT3cReloadHeaderRewrite deleting file")

	// delete a random file in etc/trafficserver which should trigger a reload
	fileNameToRemove := filepath.Join(test_config_dir, "non-empty-file.config")
	if err := os.Remove(fileNameToRemove); err != nil {
		t.Fatalf("failed to remove file '" + fileNameToRemove + "': " + err.Error())
	}

	t.Logf("DEBUG doTestT3cReloadHeaderRewrite setting update flag")
	// set the update flag, so syncds will run
	if err := ExecTOUpdater("atlanta-edge-03", false, true); err != nil {
		t.Fatalf("t3c-update failed: %v\n", err)
	}

	t.Logf("DEBUG doTestT3cReloadHeaderRewrite calling syncds")
	stdOut, _ := t3cUpdateReload(cacheHostName, "syncds")
	// Ignore the exit code error for now, because the ORT Integration Test Framework doesn't currently start ATS.
	// TODO check err, after running ATS is added to the tests.
	// if err != nil {
	// 	t.Fatalf("t3c syncds failed: %v\n", err)
	// }

	t.Logf("DEBUG TestT3cReload looking for reload string")
	if !strings.Contains(stdOut, `Running 'traffic_ctl config reload' now`) {
		t.Errorf("expected t3c to do a reload after adding a etc/trafficserver file, actual: '''%v'''\n", stdOut)
	}

	t.Logf("------------- End TestT3cReload doTestT3cReloadAnythingInTrafficserverDir ---------------")
}

func doTestT3cReloadNoChange(t *testing.T) {
	t.Logf("------------- Start doTestT3cReloadNoChange ---------------")

	cacheHostName := "atlanta-edge-03"

	t.Logf("DEBUG doTestT3cReloadHeaderRewrite calling badass")
	if stdOut, exitCode := t3cUpdateReload(cacheHostName, "badass"); exitCode != 0 {
		t.Fatalf("ERROR: t3c badass failed: code '%v' output '%v'\n", exitCode, stdOut)
	}

	t.Logf("DEBUG doTestT3cReloadHeaderRewrite deleting file")

	// no change, should not trigger a reload

	t.Logf("DEBUG doTestT3cReloadHeaderRewrite setting update flag")
	// set the update flag, so syncds will run
	if err := ExecTOUpdater("atlanta-edge-03", false, true); err != nil {
		t.Fatalf("t3c-update failed: %v\n", err)
	}

	t.Logf("DEBUG doTestT3cReloadHeaderRewrite calling syncds")
	stdOut, _ := t3cUpdateReload(cacheHostName, "syncds")
	// Ignore the exit code error for now, because the ORT Integration Test Framework doesn't currently start ATS.
	// TODO check err, after running ATS is added to the tests.
	// if err != nil {
	// 	t.Fatalf("t3c syncds failed: %v\n", err)
	// }

	t.Logf("DEBUG TestT3cReload looking for reload string")
	if strings.Contains(stdOut, `Running 'traffic_ctl config reload' now`) {
		t.Errorf("expected t3c to not reload after no change, actual: '''%v'''\n", stdOut)
	}

	t.Logf("------------- End TestT3cReload doTestT3cReloadNoChange ---------------")
}

func doTestT3cRevalCallsReload(t *testing.T) {
	t.Logf("------------- Start TestT3cReload doTestT3cRevalCallsReload ---------------")

	cacheHostName := "atlanta-edge-03"

	t.Logf("DEBUG doTestT3cRevalCallsReload calling badass")
	if stdOut, exitCode := t3cUpdateReload(cacheHostName, "badass"); exitCode != 0 {
		t.Fatalf("ERROR: t3c badass failed: code '%v' output '%v'\n", exitCode, stdOut)
	}

	t.Logf("DEBUG doTestT3cRevalCallsReload deleting file")

	// delete a regex_revalidate.config to trigger a reval change and reload
	fileNameToRemove := filepath.Join(test_config_dir, "regex_revalidate.config")
	if err := os.Remove(fileNameToRemove); err != nil {
		t.Fatalf("failed to remove file '" + fileNameToRemove + "': " + err.Error())
	}

	t.Logf("DEBUG doTestT3cRevalCallsReload setting reval flag")
	// set the update flag, so reval will run
	if err := ExecTOUpdater("atlanta-edge-03", true, false); err != nil {
		t.Fatalf("t3c-update failed: %v\n", err)
	}

	t.Logf("DEBUG doTestT3cReloadHeaderRewrite calling revalidate")
	stdOut, _ := t3cUpdateReload(cacheHostName, "revalidate")
	// Ignore the exit code error for now, because the ORT Integration Test Framework doesn't currently start ATS.
	// TODO check err, after running ATS is added to the tests.
	// if err != nil {
	// 	t.Fatalf("t3c syncds failed: %v\n", err)
	// }

	t.Logf("DEBUG TestT3cReload looking for reload string")
	if !strings.Contains(stdOut, `Running 'traffic_ctl config reload' now`) {
		t.Errorf("expected t3c to reload after reval change, actual: '''%v'''\n", stdOut)
	}

	t.Logf("------------- End TestT3cReload doTestT3cRevalCallsReload ---------------")
}

func doTestT3cReloadState(t *testing.T) {
	t.Logf("------------- Start doTestT3cReloadReloadState ---------------")

	cacheHostName := "atlanta-edge-03"

	t.Logf("DEBUG TestT3cReload calling badass")
	if stdOut, exitCode := t3cUpdateReload(cacheHostName, "badass"); exitCode != 0 {
		t.Fatalf("ERROR: t3c badass failed: code '%v' output '%v'\n", exitCode, stdOut)
	}

	t.Logf("DEBUG TestT3cReload deleting header rewrite")

	var fileNameToRemove string

	// delete header rewrite so we know should trigger a remap.config touch and reload.
	fileNameToRemove = filepath.Join(test_config_dir, "hdr_rw_first_ds-top.config")
	if err := os.Remove(fileNameToRemove); err != nil {
		t.Fatalf("failed to remove file '" + fileNameToRemove + "': " + err.Error())
	}

	// delete storage.config we know should trigger just a reload.
	fileNameToRemove = filepath.Join(test_config_dir, "storage.config")
	if err := os.Remove(fileNameToRemove); err != nil {
		t.Fatalf("failed to remove file '" + fileNameToRemove + "': " + err.Error())
	}

	t.Logf("DEBUG TestT3cReload setting upate flag")
	// set the update flag, so syncds will run
	if err := ExecTOUpdater("atlanta-edge-03", false, true); err != nil {
		t.Fatalf("t3c-update failed: %v\n", err)
	}

	t.Logf("DEBUG TestT3cReload calling syncds")
	stdOut, _ := t3cUpdateReload(cacheHostName, "syncds")
	// Ignore the exit code error for now, because the ORT Integration Test Framework doesn't currently start ATS.
	// TODO check err, after running ATS is added to the tests.
	// if err != nil {
	// 	t.Fatalf("t3c syncds failed: %v\n", err)
	// }

	t.Logf("DEBUG TestT3cReload looking for reload string")
	if !strings.Contains(stdOut, `Running 'traffic_ctl config reload' now`) {
		t.Errorf("expected t3c to do a reload after adding a header rewrite file, actual: '''%v'''\n", stdOut)
	}

	t.Logf("DEBUG TestT3cReload looking for remap.config reloading string")
	if !strings.Contains(stdOut, `updated the remap.config for reloading`) {
		t.Errorf("expected t3c to touch remap.config after adding a header rewrite file, actual: '''%v'''\n", stdOut)
	}

	t.Logf("------------- End TestT3cReload doTestT3cReloadState ---------------")
}

func t3cUpdateReload(host string, runMode string) (string, int) {
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
		"--log-location-error=stdout",
		"--log-location-info=stdout",
		"--log-location-debug=test.log",
		"--omit-via-string-release=true",
		"--git=" + "yes",
		"--run-mode=" + runMode,
	}
	stdOut, _, exitCode := t3cutil.Do("t3c", args...) // should be no stderr, we told it to log to stdout
	return string(stdOut), exitCode
}
