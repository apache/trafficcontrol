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
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/v8/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/v8/cache-config/testing/ort-tests/tcdata"
)

func TestT3cReload(t *testing.T) {
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices, tcdata.Jobs}, func() {

		t.Run("reload header rewrite", doTestT3cReloadHeaderRewrite)
		t.Run("reload anything in trafficserver dir", doTestT3cReloadAnythingInTrafficserverDir)
		t.Run("reload no change", doTestT3cReloadNoChange)
		t.Run("reval calls reload", doTestT3cRevalCallsReload)
		t.Run("reload state", doTestT3cReloadState)
	})
}

func doTestT3cReloadHeaderRewrite(t *testing.T) {
	if stdOut, exitCode := t3cUpdateReload(DefaultCacheHostName, "badass"); exitCode != 0 {
		t.Fatalf("t3c badass failed with exit code %d, output: %s", exitCode, stdOut)
	}

	// delete a file that we know should trigger a reload.
	fileNameToRemove := filepath.Join(TestConfigDir, "hdr_rw_first_ds-top.config")
	if err := os.Remove(fileNameToRemove); err != nil {
		t.Fatalf("failed to remove file '%s': %v", fileNameToRemove, err)
	}

	// set the update flag, so syncds will run
	err := tcd.QueueUpdatesForServer(DefaultCacheHostName, true)
	if err != nil {
		t.Fatalf("failed to set config update: %v", err)
	}

	stdOut, _ := t3cUpdateReload(DefaultCacheHostName, "syncds")
	// Ignore the exit code error for now, because the ORT Integration Test Framework doesn't currently start ATS.
	// TODO check err, after running ATS is added to the tests.
	// if err != nil {
	// 	t.Fatalf("t3c syncds failed: %v", err)
	// }

	if !strings.Contains(stdOut, `Running 'traffic_ctl config reload' now`) {
		t.Errorf("expected t3c to do a reload after adding a header rewrite file, actual: %s", stdOut)
	}
}

func doTestT3cReloadAnythingInTrafficserverDir(t *testing.T) {
	if stdOut, exitCode := t3cUpdateReload(DefaultCacheHostName, "badass"); exitCode != 0 {
		t.Fatalf("t3c badass failed with exit code %d, output: %s", exitCode, stdOut)
	}

	// delete a random file in etc/trafficserver which should trigger a reload
	fileNameToRemove := filepath.Join(TestConfigDir, "non-empty-file.config")
	if err := os.Remove(fileNameToRemove); err != nil {
		t.Fatalf("failed to remove file '%s': %v", fileNameToRemove, err)
	}

	// set the update flag, so syncds will run
	err := tcd.QueueUpdatesForServer(DefaultCacheHostName, true)
	if err != nil {
		t.Fatalf("failed to set config update: %v", err)
	}

	stdOut, _ := t3cUpdateReload(DefaultCacheHostName, "syncds")
	// Ignore the exit code error for now, because the ORT Integration Test Framework doesn't currently start ATS.
	// TODO check err, after running ATS is added to the tests.
	// if err != nil {
	// 	t.Fatalf("t3c syncds failed: %v", err)
	// }

	if !strings.Contains(stdOut, `Running 'traffic_ctl config reload' now`) {
		t.Errorf("expected t3c to do a reload after adding a etc/trafficserver file, actual: %s", stdOut)
	}
}

func doTestT3cReloadNoChange(t *testing.T) {
	if stdOut, exitCode := t3cUpdateReload(DefaultCacheHostName, "badass"); exitCode != 0 {
		t.Fatalf("t3c badass failed with exit code %d, output: %s", exitCode, stdOut)
	}

	// no change, should not trigger a reload

	// set the update flag, so syncds will run
	err := tcd.QueueUpdatesForServer(DefaultCacheHostName, true)
	if err != nil {
		t.Fatalf("failed to set config update: %v", err)
	}

	stdOut, _ := t3cUpdateReload(DefaultCacheHostName, "syncds")
	// Ignore the exit code error for now, because the ORT Integration Test Framework doesn't currently start ATS.
	// TODO check err, after running ATS is added to the tests.
	// if err != nil {
	// 	t.Fatalf("t3c syncds failed: %v", err)
	// }

	if strings.Contains(stdOut, `Running 'traffic_ctl config reload' now`) {
		t.Errorf("expected t3c to not reload after no change, actual: %s", stdOut)
	}
}

func doTestT3cRevalCallsReload(t *testing.T) {
	if stdOut, exitCode := t3cUpdateReload(DefaultCacheHostName, "badass"); exitCode != 0 {
		t.Fatalf("t3c badass failed with exit code %d, output: %s", exitCode, stdOut)
	}

	// delete a regex_revalidate.config to trigger a reval change and reload
	fileNameToRemove := filepath.Join(TestConfigDir, "regex_revalidate.config")
	if err := os.Remove(fileNameToRemove); err != nil {
		t.Fatalf("failed to remove file '%s': %v", fileNameToRemove, err)
	}

	// set the update flag, so reval will run
	// TODO this sets the config update, do we need to do reval instead?
	err := tcd.QueueUpdatesForServer(DefaultCacheHostName, true)
	if err != nil {
		t.Fatalf("failed to set config update: %v", err)
	}

	stdOut, _ := t3cUpdateReload(DefaultCacheHostName, "revalidate")
	// Ignore the exit code error for now, because the ORT Integration Test Framework doesn't currently start ATS.
	// TODO check err, after running ATS is added to the tests.
	// if err != nil {
	// 	t.Fatalf("t3c syncds failed: %v", err)
	// }

	if !strings.Contains(stdOut, `Running 'traffic_ctl config reload' now`) {
		t.Errorf("expected t3c to reload after reval change, actual: %s", stdOut)
	}
}

func doTestT3cReloadState(t *testing.T) {
	if stdOut, exitCode := t3cUpdateReload(DefaultCacheHostName, "badass"); exitCode != 0 {
		t.Fatalf("t3c badass failed with exit code %d, output: %s", exitCode, stdOut)
	}

	// delete header rewrite so we know should trigger a remap.config touch and reload.
	fileNameToRemove := filepath.Join(TestConfigDir, "hdr_rw_first_ds-top.config")
	if err := os.Remove(fileNameToRemove); err != nil {
		t.Fatalf("failed to remove file '%s': %v", fileNameToRemove, err)
	}

	// create plugin.config to trigger restart directive
	pluginConfigPath := filepath.Join(TestConfigDir, "plugin.config")
	contents := []byte("remap_stats.so")
	err := ioutil.WriteFile(pluginConfigPath, contents, 0666)
	if err != nil {
		t.Fatalf("Unable to create file %s", pluginConfigPath)
	}

	// set the update flag, so syncds will run
	// TODO this sets the config update, do we need to do reval instead?
	err = tcd.QueueUpdatesForServer(DefaultCacheHostName, true)
	if err != nil {
		t.Fatalf("failed to set config update: %v", err)
	}

	stdOut, _ := t3cUpdateReload(DefaultCacheHostName, "syncds")

	if !strings.Contains(stdOut, "Final state: remap.config: true reload: true restart: true ntpd: false sysctl: false") {
		t.Errorf("expected t3c Final reload state for remap.config, reload and restart, actual: %s", stdOut)
	}

	// remove plugin.config file for next test
	if err := os.Remove(pluginConfigPath); err != nil {
		t.Fatalf("failed to remove file '%s': %v", pluginConfigPath, err)
	}
}

func t3cUpdateReload(host string, runMode string) (string, int) {
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
		"--run-mode=" + runMode,
	}
	_, stdErr, exitCode := t3cutil.Do("t3c", args...)
	return string(stdErr), exitCode
}
