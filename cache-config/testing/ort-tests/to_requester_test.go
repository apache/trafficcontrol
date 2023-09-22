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
	"encoding/json"
	"errors"
	"os/exec"
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/v8/cache-config/testing/ort-tests/tcdata"
	"github.com/apache/trafficcontrol/v8/lib/go-atscfg"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

type Package struct {
	Name    *string `json:"name"`
	Version *string `json:"version"`
}

func TestTORequester(t *testing.T) {
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices}, func() {

		// chkconfig test
		output, err := ExecTORequester(DefaultCacheHostName, "chkconfig")
		if err != nil {
			t.Fatalf("t3c-request exec failed: %v", err)
		}
		var chkConfig []map[string]interface{}
		err = json.Unmarshal([]byte(output), &chkConfig)
		if err != nil {
			t.Fatalf("failed to parse t3c-request output: %v", err)
		}
		if len(chkConfig) < 1 {
			t.Fatal("expected at least one chkconfig entry, got zero")
		}
		firstName := chkConfig[0]["name"]
		if firstName != "trafficserver" {
			t.Fatalf("expected the name of the first chkconfig entry to be 'trafficserver', actual: %s", firstName)
		}

		// get system-info test
		output, err = ExecTORequester(DefaultCacheHostName, "system-info")
		if err != nil {
			t.Fatalf("t3c-request exec failed: %v", err)
		}
		var sysInfo map[string]interface{}
		err = json.Unmarshal([]byte(output), &sysInfo)
		if err != nil {
			t.Fatalf("failed to parse t3c-request output: %v", err)
		}
		instanceName := sysInfo["tm.instance_name"]
		if instanceName != "Traffic Ops ORT Tests" {
			t.Fatalf("expected 'tm.instance_name' to be 'Traffic Ops ORT Tests', actual: %s", instanceName)
		}

		// statuses test
		output, err = ExecTORequester(DefaultCacheHostName, "statuses")
		if err != nil {
			t.Fatalf("t3c-request exec failed: %v", err)
		}
		// should parse json to an array of 'tc.Status'
		var statuses []tc.Status
		err = json.Unmarshal([]byte(output), &statuses)
		if err != nil {
			t.Fatalf("failed to parse t3c-request output: %v", err)
		}

		// packages test
		output, err = ExecTORequester(DefaultCacheHostName, "packages")
		if err != nil {
			t.Fatalf("t3c-request exec failed: %v", err)
		}
		// should parse to an array of 'Package'
		var packages []Package
		err = json.Unmarshal([]byte(output), &packages)
		if err != nil {
			t.Fatalf("failed to parse t3c-request output: %v", err)
		}
		if len(packages) < 1 {
			t.Fatal("expected at least one package, got zero")
		}
		if packages[0].Name == nil {
			t.Fatal("null or undefined name for the first package in t3c-request output")
		}
		pkgName := *packages[0].Name
		if pkgName != "trafficserver" {
			t.Fatalf("expected first package to be named 'trafficserver', actual: %s", pkgName)
		}

		// update-status test
		output, err = ExecTORequester(DefaultCacheHostName, CMDUpdateStatus)
		if err != nil {
			t.Fatalf("t3c-request exec failed: %v", err)
		}
		var serverStatus atscfg.ServerUpdateStatus
		err = json.Unmarshal([]byte(output), &serverStatus)
		if err != nil {
			t.Fatalf("failed to parse t3c-request output: %v", err)
		}
		if serverStatus.HostName != DefaultCacheHostName {
			t.Fatalf("expected server status hosname to be '%s', actual: %s", DefaultCacheHostName, serverStatus.HostName)
		}

	})
}

func ExecTORequester(host string, data_req string) (string, error) {
	args := []string{
		"--traffic-ops-insecure=true",
		"--traffic-ops-timeout-milliseconds=3000",
		"--traffic-ops-user=" + tcd.Config.TrafficOps.Users.Admin,
		"--traffic-ops-password=" + tcd.Config.TrafficOps.UserPassword,
		"--traffic-ops-url=" + tcd.Config.TrafficOps.URL,
		"--cache-host-name=" + host,
		"-vv",
		"--get-data=" + data_req,
	}
	cmd := exec.Command("/usr/bin/t3c-request", args...)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err := cmd.Run()
	if err != nil {
		return "", errors.New(err.Error() + ": " + "stdout: " + out.String() + " stderr: " + errOut.String())
	}

	// capture the last line of JSON in the 'Stdout' buffer 'out'
	output := strings.Split(strings.TrimSpace(strings.Replace(out.String(), "\r\n", "\n", -1)), "\n")
	lastLine := output[len(output)-1]

	return lastLine, nil
}
