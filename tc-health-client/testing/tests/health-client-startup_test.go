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

package hctest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/cache-config/testing/ort-tests/tcdata"
	"github.com/apache/trafficcontrol/v8/tc-health-client/testing/tests/hcutil"
	"github.com/apache/trafficcontrol/v8/tc-health-client/tmagent"
)

func startHealthClient() {
	outbuf, errbuf, result := hcutil.Do("systemctl", "start", "tc-health-client", "-vvv")
	if result != 0 {
		fmt.Fprintf(os.Stdout, "Error starting the health-client: %s\n", string(errbuf))
	} else {
		fmt.Fprintf(os.Stdout, "the health-client was succesfully started: %s\n", string(outbuf))
	}
}

func TestHealthClientStartup(t *testing.T) {
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
	}, func() {

		// initialize variables
		cfg := tmagent.ParentInfo{}
		pollStateFile := "/var/log/trafficcontrol/poll-state.json"
		atlantaMid := "atlanta-mid-16.ga.atlanta.kabletown.net"
		dtrcMid := "dtrc-mid-02.kabletown.net"
		rascal := "rascal01.kabletown.net"

		waitTime, err := time.ParseDuration("5s")
		if err != nil {
			fmt.Fprintf(os.Stdout, "failed to parse a value for waitTime")
			os.Exit(1)
		}

		// mark down some parents using ATS traffic_ctl
		_, errbuf, result := hcutil.Do("/opt/trafficserver/bin/traffic_ctl", "host", "down", "--reason", "active", atlantaMid)
		if result != 0 {
			fmt.Fprintf(os.Stdout, "%s\n", string(errbuf))
			t.Fatalf("unable to mark down parent '%s'\n", atlantaMid)
		} else {
			fmt.Fprintf(os.Stdout, "marked down '%s'\n", atlantaMid)
		}
		_, errbuf, result = hcutil.Do("/opt/trafficserver/bin/traffic_ctl", "host", "down", "--reason", "active", dtrcMid)
		if result != 0 {
			fmt.Fprintf(os.Stdout, "%s\n", string(errbuf))
			t.Fatalf("unable to mark down parent '%s'\n", dtrcMid)
		} else {
			fmt.Fprintf(os.Stdout, "marked down '%s'\n", dtrcMid)
		}

		// startup the health-client
		fmt.Fprintf(os.Stdout, "Starting the tc-health-client\n")
		go startHealthClient()

		// wait for the health client to write it's poll state
		time.Sleep(waitTime)

		fmt.Fprintf(os.Stdout, "Running tests\n")

		// read the health-client poll-state file.
		cfg = tmagent.ParentInfo{}
		content, err := ioutil.ReadFile(pollStateFile)
		if err != nil {
			t.Fatalf("could not read the %s file: %s\n", pollStateFile, err.Error())
		}
		err = json.Unmarshal(content, &cfg)
		if err != nil {
			t.Fatalf("could not unmarshal %s: %s\n", pollStateFile, err.Error())
		}

		// we marked down mids, now test that the health client read the ATS
		// Host Status and see's that they are down.
		p, ok := cfg.LoadParentStatus(atlantaMid)
		if !ok {
			t.Fatalf("Expected %s to be in parents but it's not", atlantaMid)
		}
		if p.ActiveReason != false {
			t.Fatalf("Expected %s to be marked down but it's not", atlantaMid)
		} else {
			fmt.Fprintf(os.Stdout, "%s is available: %v\n", atlantaMid, p.ActiveReason)
		}
		p, ok = cfg.LoadParentStatus(dtrcMid)
		if !ok {
			t.Fatalf("Expected %s to be in parents but it's not", dtrcMid)
		}
		if p.ActiveReason != false {
			t.Fatalf("Expected %s to be marked down but it's not", dtrcMid)
		} else {
			fmt.Fprintf(os.Stdout, "%s is available: %v\n", dtrcMid, p.ActiveReason)
		}

		// verify that the health-client was able to poll and get an available
		// traffic monitor from TrafficOps
		_, ok = cfg.TOData.Get().Monitors[rascal]
		if !ok {
			t.Fatalf("Expected %s to be available but it's not", rascal)
		} else {
			fmt.Fprintf(os.Stdout, "%s is available: true\n", rascal)
		}

		fmt.Fprintf(os.Stdout, "Stopping the tc-health-client\n")
		_, errbuf, result = hcutil.Do("/usr/bin/systemctl", "stop", "tc-health-client", "-vvv")
		if result != 0 {
			fmt.Fprintf(os.Stdout, "%s\n", string(errbuf))
			t.Fatalf("unable to stop the tc-health-client\n")
		}
	})
}
