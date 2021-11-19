// Package orttest provides testing for the t3c utility(ies).
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
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/cache-config/testing/ort-tests/tcdata"
	"github.com/apache/trafficcontrol/cache-config/testing/ort-tests/util"
)

const trafficServerOwner = "ats"

var (
	base_line_dir   = "baseline-configs"
	test_config_dir = "/opt/trafficserver/etc/trafficserver"

	testFiles = [8]string{
		"astats.config",
		"hdr_rw_first_ds-top.config",
		"hosting.config",
		"parent.config",
		"records.config",
		"remap.config",
		"storage.config",
		"volume.config",
	}
)

func verifyPluginConfigs(t *testing.T) {
	err := runCheckRefs("/opt/trafficserver/etc/trafficserver/remap.config")
	if err != nil {
		t.Errorf("Plugin verification failed for remap.config: %v", err)
	}
	err = runCheckRefs("/opt/trafficserver/etc/trafficserver/plugin.config")
	if err != nil {
		t.Errorf("Plugin verification failed for plugin.config: %v", err)
	}

}

func syncDSTest(t *testing.T) {
	// remove the remap.config in preparation for running syncds
	remap := filepath.Join(test_config_dir, "remap.config")
	err := os.Remove(remap)
	if err != nil {
		t.Fatalf("unable to remove %s: %v", remap, err)
	}
	// prepare for running syncds.
	err = ExecTOUpdater(cacheHostName, false, true)
	if err != nil {
		t.Fatalf("queue updates failed: %v", err)
	}

	// remap.config is removed and atlanta-edge-03 should have
	// queue updates enabled.  run t3c to verify a new remap.config
	// is pulled down.
	err = runApply(cacheHostName, "syncds")
	if err != nil {
		t.Fatalf("t3c syncds failed: %v", err)
	}
	if !util.FileExists(remap) {
		t.Fatalf("syncds failed to pull down %s", remap)
	}

}

func TestT3cBadassAndSyncDs(t *testing.T) {
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters, tcdata.Statuses,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices}, func() {

		// traffic_ctl doesn't work because the test framework doesn't currently run ATS.
		// So, temporarily replace it with a no-op
		// TODO: remove this when running ATS is added to the test framework

		if err := os.Rename(`/opt/trafficserver/bin/traffic_ctl`, `/opt/trafficserver/bin/traffic_ctl.real`); err != nil {
			t.Fatal("temporarily moving traffic_ctl: " + err.Error())
		}

		fi, err := os.OpenFile(`/opt/trafficserver/bin/traffic_ctl`, os.O_RDWR|os.O_CREATE, 755)
		if err != nil {
			t.Fatal("creating temp no-op traffic_ctl file: " + err.Error())
		}
		if _, err := fi.WriteString(`#!/usr/bin/env bash` + "\n"); err != nil {
			fi.Close()
			t.Fatal("writing temp no-op traffic_ctl file: " + err.Error())
		}
		fi.Close()

		defer func() {
			if err := os.Rename(`/opt/trafficserver/bin/traffic_ctl.real`, `/opt/trafficserver/bin/traffic_ctl`); err != nil {
				t.Fatal("moving real traffic_ctl back: " + err.Error())
			}
		}()

		// run badass and check config files.
		if err := runApply(cacheHostName, "badass"); err != nil {
			t.Fatalf("t3c badass failed: %v", err)
		}

		// Use this for uid/gid file check
		atsUser, err := user.Lookup(trafficServerOwner)
		var atsUid string
		var atsGid string

		if err != nil {
			t.Logf("Unable to look up user: %s: %v", trafficServerOwner, err)
		} else {
			atsUid = atsUser.Uid
			atsGid = atsUser.Gid
		}

		for _, v := range testFiles {
			bfn := filepath.Join(base_line_dir, v)
			if !util.FileExists(bfn) {
				t.Fatalf("missing baseline config file, %s,  needed for tests", bfn)
			}
			tfn := filepath.Join(test_config_dir, v)
			if !util.FileExists(tfn) {
				t.Fatalf("missing the expected config file, %s", tfn)
			}

			diffStr, err := util.DiffFiles(bfn, tfn)
			if err != nil {
				t.Fatalf("diffing %s and %s: %v", tfn, bfn, err)
			} else if diffStr != "" {
				t.Errorf("%s and %s differ: %v", tfn, bfn, diffStr)
			} else {
				t.Logf("%s and %s diff clean", tfn, bfn)
			}

			fileInfo, err := os.Stat(tfn)
			if err != nil {
				t.Errorf("Error getting stats on %s: %v", tfn, err)
			} else {
				if statStruct, ok := fileInfo.Sys().(*syscall.Stat_t); ok {
					uid := strconv.Itoa(int(statStruct.Uid))
					if uid != atsUid {
						t.Errorf("Unexpected uid for file: %s: %s, expected %s", v, uid, atsUid)
					}
					gid := strconv.Itoa(int(statStruct.Gid))
					if gid != atsGid {
						t.Errorf("Unexpected gid for file: %s: %s, expected %s", v, gid, atsGid)
					}
				}
			}
		}

		time.Sleep(time.Second * 5)

		t.Run("Verify Plugin Configs", verifyPluginConfigs)
		t.Run("SyncDS Test", syncDSTest)
	})
}
