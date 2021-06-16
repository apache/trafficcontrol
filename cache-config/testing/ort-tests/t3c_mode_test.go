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
	"fmt"
	"github.com/apache/trafficcontrol/cache-config/testing/ort-tests/tcdata"
	"github.com/apache/trafficcontrol/cache-config/testing/ort-tests/util"
	"os"
	"os/user"
	"strconv"
	"syscall"
	"testing"
	"time"
)

const TrafficServerOwner = "ats"

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

func TestT3cBadassAndSyncDs(t *testing.T) {
	fmt.Println("------------- Starting TestT3cBadassAndSyncDs ---------------")
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
		if err := runApply("atlanta-edge-03", "badass", 0); err != nil {
			t.Fatalf("ERROR: t3c badass failed: %v\n", err)
		}

		// Use this for uid/gid file check
		atsUser, err := user.Lookup(TrafficServerOwner)
		var atsUid string
		var atsGid string

		if err != nil {
			t.Logf("Unable to look up user: %s: %v", TrafficServerOwner, err)
		} else {
			atsUid = atsUser.Uid
			atsGid = atsUser.Gid
		}

		for _, v := range testFiles {
			bfn := base_line_dir + "/" + v
			if !util.FileExists(bfn) {
				t.Fatalf("ERROR: missing baseline config file, %s,  needed for tests", bfn)
			}
			tfn := test_config_dir + "/" + v
			if !util.FileExists(tfn) {
				t.Fatalf("ERROR: missing the expected config file, %s", tfn)
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

		fmt.Println("------------------------ Verify Plugin Configs ----------------")
		err = runCheckRefs("/opt/trafficserver/etc/trafficserver/remap.config")
		if err != nil {
			t.Errorf("Plugin verification failed for remap.config: " + err.Error())
		}
		err = runCheckRefs("/opt/trafficserver/etc/trafficserver/plugin.config")
		if err != nil {
			t.Errorf("Plugin verification failed for plugin.config: " + err.Error())
		}

		fmt.Println("----------------- End of Verify Plugin Configs ----------------")

		fmt.Println("------------------------ running SYNCDS Test ------------------")
		// remove the remap.config in preparation for running syncds
		remap := test_config_dir + "/remap.config"
		err = os.Remove(remap)
		if err != nil {
			t.Fatalf("ERROR: unable to remove %s\n", remap)
		}
		// prepare for running syncds.
		err = ExecTOUpdater("atlanta-edge-03", false, true)
		if err != nil {
			t.Fatalf("ERROR: queue updates failed: %v\n", err)
		}

		// remap.config is removed and atlanta-edge-03 should have
		// queue updates enabled.  run t3c to verify a new remap.config
		// is pulled down.
		err = runApply("atlanta-edge-03", "syncds", 0)
		if err != nil {
			t.Fatalf("ERROR: t3c syncds failed: %v\n", err)
		}
		if !util.FileExists(remap) {
			t.Fatalf("ERROR: syncds failed to pull down %s\n", remap)
		}
		fmt.Println("------------------------ end SYNCDS Test ------------------")

	})
	fmt.Println("------------- End of TestT3cBadassAndSyncDs ---------------")
}
