package v5

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
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	client "github.com/apache/trafficcontrol/traffic_ops/v5-client"
)

func TestServerChecks(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, ServerCheckExtensions, ServerChecks}, func() {
		CreateTestInvalidServerChecks(t)
		UpdateTestServerChecks(t)
		GetTestServerChecks(t)
		GetTestServerChecksWithName(t)
		GetTestServerChecksWithID(t)
	})
}

func CreateTestServerChecks(t *testing.T) {
	SwitchSession(toReqTimeout, Config.TrafficOps.URL, Config.TrafficOps.Users.Admin, Config.TrafficOps.UserPassword, Config.TrafficOps.Users.Extension, Config.TrafficOps.UserPassword)

	for _, servercheck := range testData.Serverchecks {
		resp, _, err := TOSession.InsertServerCheckStatus(servercheck, client.RequestOptions{})
		if err != nil {
			t.Errorf("could not insert Servercheck: %v - alerts: %+v", err, resp.Alerts)
		}
	}
	SwitchSession(toReqTimeout, Config.TrafficOps.URL, Config.TrafficOps.Users.Extension, Config.TrafficOps.UserPassword, Config.TrafficOps.Users.Admin, Config.TrafficOps.UserPassword)
}

func CreateTestInvalidServerChecks(t *testing.T) {
	if len(testData.Serverchecks) < 1 {
		t.Fatal("Need at least one Servercheck to test creating an invalid Servercheck")
	}
	toReqTimeout := time.Second * time.Duration(Config.Default.Session.TimeoutInSecs)

	_, _, err := TOSession.InsertServerCheckStatus(testData.Serverchecks[0], client.RequestOptions{})
	if err == nil {
		t.Error("expected to receive error with non extension user")
	}

	SwitchSession(toReqTimeout, Config.TrafficOps.URL, Config.TrafficOps.Users.Admin, Config.TrafficOps.UserPassword, Config.TrafficOps.Users.Extension, Config.TrafficOps.UserPassword)

	invalidServerCheck := tc.ServercheckRequestNullable{
		Name:     util.StrPtr("BOGUS"),
		Value:    util.IntPtr(1),
		ID:       util.IntPtr(-1),
		HostName: util.StrPtr("bogus_hostname"),
	}

	// Attempt to create a ServerCheck with invalid server ID
	_, _, err = TOSession.InsertServerCheckStatus(invalidServerCheck, client.RequestOptions{})
	if err == nil {
		t.Error("expected to receive error with invalid id")
	}

	invalidServerCheck.ID = nil
	// Attempt to create a ServerCheck with invalid host name
	_, _, err = TOSession.InsertServerCheckStatus(invalidServerCheck, client.RequestOptions{})
	if err == nil {
		t.Error("expected to receive error with invalid host name")
	}

	// get valid name to get past host check
	invalidServerCheck.Name = testData.Serverchecks[0].Name

	// Attempt to create a ServerCheck with invalid servercheck name
	_, _, err = TOSession.InsertServerCheckStatus(invalidServerCheck, client.RequestOptions{})
	if err == nil {
		t.Error("expected to receive error with invalid servercheck name")
	}
	SwitchSession(toReqTimeout, Config.TrafficOps.URL, Config.TrafficOps.Users.Extension, Config.TrafficOps.UserPassword, Config.TrafficOps.Users.Admin, Config.TrafficOps.UserPassword)
}

func UpdateTestServerChecks(t *testing.T) {
	SwitchSession(toReqTimeout, Config.TrafficOps.URL, Config.TrafficOps.Users.Admin, Config.TrafficOps.UserPassword, Config.TrafficOps.Users.Extension, Config.TrafficOps.UserPassword)
	for _, servercheck := range testData.Serverchecks {
		if servercheck.Value == nil {
			t.Error("Found servercheck in the testing data with null or undefined Value")
			continue
		}
		*servercheck.Value--
		resp, _, err := TOSession.InsertServerCheckStatus(servercheck, client.RequestOptions{})
		if err != nil {
			if servercheck.Name != nil {
				t.Logf("Servercheck Name: %s", *servercheck.Name)
			}
			if servercheck.HostName != nil {
				t.Logf("Servercheck Host Name: %s", *servercheck.HostName)
			}
			t.Errorf("could not update servercheck: %v - alerts: %+v", err, resp.Alerts)
		}
	}
	SwitchSession(toReqTimeout, Config.TrafficOps.URL, Config.TrafficOps.Users.Extension, Config.TrafficOps.UserPassword, Config.TrafficOps.Users.Admin, Config.TrafficOps.UserPassword)
}

func GetTestServerChecks(t *testing.T) {
	if len(testData.Serverchecks) < 1 {
		t.Fatal("Need at least one Servercheck to test creating an getting Serverchecks")
	}
	if testData.Serverchecks[0].HostName == nil {
		t.Fatal("Found a Servercheck in the testing data wih null or undefined Host Name")
	}
	hostname := *testData.Serverchecks[0].HostName

	// Get server checks
	serverChecksResp, _, err := TOSession.GetServersChecks(client.RequestOptions{})
	if err != nil {
		t.Fatalf("could not get Serverchecks: %v - alerts: %+v", err, serverChecksResp.Alerts)
	}
	found := false
	for _, sc := range serverChecksResp.Response {
		if sc.HostName == hostname {
			found = true

			if sc.Checks == nil {
				t.Errorf("server %s had no checks - expected it to have at least two", hostname)
				break
			}

			if ort, ok := sc.Checks["ORT"]; !ok {
				t.Error("no 'ORT' servercheck exists - expected it to exist")
			} else if ort == nil {
				t.Error("'null' returned for ORT value servercheck - expected pointer to 12")
			} else if *ort != 12 {
				t.Errorf("%d returned for ORT value servercheck - expected 12", *ort)
			}

			if ilo, ok := sc.Checks["ILO"]; !ok {
				t.Error("no 'ILO' servercheck exists - expected it to exist")
			} else if ilo == nil {
				t.Error("'null' returned for ILO value servercheck - expected pointer to 0")
			} else if *ilo != 0 {
				t.Errorf("%d returned for ILO value servercheck - expected 0", *ilo)
			}
			break
		}
	}
	if !found {
		t.Errorf("expected to find servercheck for host %s", hostname)
	}
}

func GetTestServerChecksWithName(t *testing.T) {
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("hostName", "atlanta-edge-01")

	// Get server checks
	scResp, _, err := TOSession.GetServersChecks(opts)
	if len(scResp.Response) == 0 {
		t.Fatal("no server checks in response, quitting")
	}
	if err != nil {
		t.Fatalf("could not get Serverchecks filtered by name 'atlanta-edge-01': %v - alerts: %+v", err, scResp.Alerts)
	}

	//Add unknown param key
	opts.QueryParameters.Add("foo", "car")
	// Get server checks
	resp, _, err := TOSession.GetServersChecks(opts)
	if len(resp.Response) == 0 {
		t.Fatal("no server checks in response, quitting")
	}
	if err != nil {
		t.Fatalf("could not get Serverchecks filtered by server Host Name '%s': %v - alerts: %+v", resp.Response[0].HostName, err, resp.Alerts)
	}

	if len(scResp.Response) != len(resp.Response) {
		t.Fatalf("expected: Both response lengths should be equal, got: first resp: %d - second resp: %d", len(scResp.Response), len(resp.Response))
	}
}

func GetTestServerChecksWithID(t *testing.T) {
	serverChecksResp, _, err := TOSession.GetServersChecks(client.RequestOptions{})
	if err != nil {
		t.Errorf("Unexpected error getting Serverchecks: %v - alerts: %+v", err, serverChecksResp.Alerts)
	}
	if len(serverChecksResp.Response) == 0 {
		t.Fatal("no server checks in response, quitting")
	}
	if serverChecksResp.Response[0].ID == 0 {
		t.Fatal("ID of the response server is nil, quitting")
	}
	id := serverChecksResp.Response[0].ID

	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("id", strconv.Itoa(id))

	// Get server checks
	scResp, _, err := TOSession.GetServersChecks(opts)
	if len(scResp.Response) == 0 {
		t.Fatal("no server checks in response, quitting")
	}
	if err != nil {
		t.Fatalf("could not get Serverchecks by ID %d: %v - alerts: %+v", scResp.Response[0].ID, err, scResp.Alerts)
	}

	//Add unknown param key
	opts.QueryParameters.Add("foo", "car")
	// Get server checks
	resp, _, err := TOSession.GetServersChecks(opts)
	if len(resp.Response) == 0 {
		t.Fatal("no server checks in response, quitting")
	}
	if err != nil {
		t.Fatalf("could not get Serverchecks filtered by ID %d with extraneous 'foo' parameter: %v - alerts: %+v", resp.Response[0].ID, err, resp.Alerts)
	}

	if len(scResp.Response) != len(resp.Response) {
		t.Fatalf("expected: Both response lengths should be equal, got: first resp:%v-second resp:%v", len(scResp.Response), len(resp.Response))
	}
}

// Need to define no-op function as TCObj interface expects a delete function
// There is no delete path for serverchecks
func DeleteTestServerChecks(*testing.T) {
	return
}
