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

	"github.com/apache/trafficcontrol/v8/cache-config/testing/ort-tests/tcdata"
	"github.com/apache/trafficcontrol/v8/cache-config/testing/ort-tests/util"
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

func verifyRemapConfigPlaced(t *testing.T) {
	// remove the remap.config in preparation for running syncds
	remap := filepath.Join(TestConfigDir, "remap.config")
	err := os.Remove(remap)
	if err != nil {
		t.Fatalf("unable to remove %s: %v", remap, err)
	}
	// prepare for running syncds.
	err = tcd.QueueUpdatesForServer(DefaultCacheHostName, true)
	if err != nil {
		t.Fatalf("failed to set config update: %v", err)
	}

	// remap.config is removed and atlanta-edge-03 should have
	// queue updates enabled.  run t3c to verify a new remap.config
	// is pulled down.
	err = runApply(DefaultCacheHostName, "syncds")
	if err != nil {
		t.Fatalf("t3c syncds failed: %v", err)
	}
	if !util.FileExists(remap) {
		t.Fatalf("syncds failed to pull down %s", remap)
	}
}

// given a filename checks that the baseline and generated files diff clean and
// that the generated files have the given user and owner group IDs
func checkDiff(fName, atsUid, atsGid string, t *testing.T) {
	bfn := filepath.Join(BaselineConfigDir, fName)
	if !util.FileExists(bfn) {
		t.Fatalf("missing baseline config file, %s,  needed for tests", bfn)
	}
	tfn := filepath.Join(TestConfigDir, fName)
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
				t.Errorf("Unexpected uid for file: %s: %s, expected %s", fName, uid, atsUid)
			}
			gid := strconv.Itoa(int(statStruct.Gid))
			if gid != atsGid {
				t.Errorf("Unexpected gid for file: %s: %s, expected %s", fName, gid, atsGid)
			}
		}
	}
}

func TestT3cBadassAndSyncDs(t *testing.T) {
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices}, func() {

		// traffic_ctl doesn't work because the test framework doesn't currently run ATS.
		// So, temporarily replace it with a no-op
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

		// run badass and check config files.
		if err := runApply(DefaultCacheHostName, "badass"); err != nil {
			t.Fatalf("t3c badass failed: %v", err)
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

		for _, v := range TestFiles {
			t.Run("check diff of "+v+" between baseline and badass-generated", func(t *testing.T) { checkDiff(v, atsUid, atsGid, t) })
		}

		time.Sleep(time.Second * 5)

		t.Run("Verify Plugin Configs", verifyPluginConfigs)
		t.Run("Verify remap.config placed as expected after SYNCDS run", verifyRemapConfigPlaced)
	})
}
