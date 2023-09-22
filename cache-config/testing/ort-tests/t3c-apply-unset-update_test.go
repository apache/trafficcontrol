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
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/apache/trafficcontrol/v8/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/v8/cache-config/testing/ort-tests/tcdata"
	"github.com/apache/trafficcontrol/v8/lib/go-atscfg"
)

func verifyUpdateStatusIsFalse() error {
	output, err := runRequest(DefaultCacheHostName, CMDUpdateStatus)
	if err != nil {
		return fmt.Errorf("t3c-request failed: %w", err)
	}
	serverStatus := atscfg.ServerUpdateStatus{}
	if err = json.Unmarshal([]byte(output), &serverStatus); err != nil {
		return fmt.Errorf("failed to parse t3c-request output: %w", err)
	}
	if serverStatus.HostName != DefaultCacheHostName {
		return fmt.Errorf("expected update-status host '%s', actual: '%s'", DefaultCacheHostName, serverStatus.HostName)
	}
	if serverStatus.RevalPending {
		return fmt.Errorf("expected RevalPending false after syncds run")
	}
	if serverStatus.UpdatePending {
		return fmt.Errorf("expected UpdatePending false after syncds run")
	}
	return nil
}

func verifyUpdateStatusIsTrue() error {
	output, err := runRequest(DefaultCacheHostName, CMDUpdateStatus)
	if err != nil {
		return fmt.Errorf("update-status run failed: %w", err)
	}
	serverStatus := atscfg.ServerUpdateStatus{}
	if err = json.Unmarshal([]byte(output), &serverStatus); err != nil {
		return fmt.Errorf("failed to parse update-status output: %w", err)
	}
	if serverStatus.HostName != DefaultCacheHostName {
		return fmt.Errorf("expected request update-status host '%s', actual: '%s'", DefaultCacheHostName, serverStatus.HostName)
	}
	if serverStatus.RevalPending {
		return errors.New("expected RevalPending false after update")
	}
	if !serverStatus.UpdatePending {
		return errors.New("expected UpdatePending true after update")
	}

	return nil
}

func TestT3cUnsetsUpdateFlag(t *testing.T) {
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices}, func() {

		if stdOut, exitCode := t3cUpdateUnsetFlag(DefaultCacheHostName, "badass"); exitCode != 0 {
			t.Fatalf("t3c badass failed with code %d output: %s", exitCode, stdOut)
		}

		// delete a file that we know should trigger a reload.
		fileNameToRemove := filepath.Join(TestConfigDir, "hdr_rw_first_ds-top.config")
		if err := os.Remove(fileNameToRemove); err != nil {
			t.Fatalf("failed to remove file '%s': %v", fileNameToRemove, err)
		}

		// set the update flag, so syncds will run
		err := tcd.QueueUpdatesForServer(DefaultCacheHostName, true)
		if err != nil {
			t.Fatalf("failed to queue updates: %v", err)
		}

		if err := verifyUpdateStatusIsTrue(); err != nil {
			t.Errorf("verification that update status after syncds is true failed: %v", err)
		}

		// traffic_ctl doesn't work because the test framework doesn't currently run ATS.
		// So, temporarily replace it with a no-op, so t3c-apply gets far enough to un-set the update flag.
		// TODO: remove this when running ATS is added to the test framework

		if err := os.Rename(`/opt/trafficserver/bin/traffic_ctl`, `/opt/trafficserver/bin/traffic_ctl.real`); err != nil {
			t.Fatalf("temporarily moving traffic_ctl: %v", err)
		}

		fi, err := os.OpenFile(`/opt/trafficserver/bin/traffic_ctl`, os.O_RDWR|os.O_CREATE, 755)
		if err != nil {
			t.Fatalf("creating temp no-op traffic_ctl file: %v", err)
		}
		if _, err := fi.WriteString(`#!/usr/bin/env bash` + "\n"); err != nil {
			fi.Close()
			t.Fatalf("writing temp no-op traffic_ctl file: %v", err)
		}
		fi.Close()

		defer func() {
			if err := os.Rename(`/opt/trafficserver/bin/traffic_ctl.real`, `/opt/trafficserver/bin/traffic_ctl`); err != nil {
				t.Fatalf("moving real traffic_ctl back: %v", err)
			}
		}()

		stdOut, _ := t3cUpdateUnsetFlag(DefaultCacheHostName, "syncds")
		// Ignore the exit code error for now, because the ORT Integration Test Framework doesn't currently start ATS.
		// TODO check err, after running ATS is added to the tests.
		// if err != nil {
		// 	t.Fatalf("t3c syncds failed: %v", err)
		// }

		t.Logf("TestT3cTOUpdates t3cUpdateUnsetFlag stdout: %s", stdOut)
		if err := verifyUpdateStatusIsFalse(); err != nil {
			t.Errorf("verification that update status after syncds is false failed: %v", err)
		}
	})
}

func t3cUpdateUnsetFlag(host string, runMode string) (string, int) {
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
		"--git=no",
		"--run-mode=" + runMode,
	}
	stdOut, _, exitCode := t3cutil.Do("t3c", args...) // should be no stderr, we told it to log to stdout
	return string(stdOut), exitCode
}
