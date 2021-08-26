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
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	client "github.com/apache/trafficcontrol/traffic_ops/v5-client"
)

var (
	toReqTimeout = time.Second * time.Duration(Config.Default.Session.TimeoutInSecs)
)

func TestServerCheckExtensions(t *testing.T) {
	WithObjs(t, []TCObj{ServerCheckExtensions}, func() {
		CreateTestInvalidServerCheckExtensions(t)
	})
}

func CreateTestServerCheckExtensions(t *testing.T) {
	SwitchSession(toReqTimeout, Config.TrafficOps.URL, Config.TrafficOps.Users.Admin, Config.TrafficOps.UserPassword, Config.TrafficOps.Users.Extension, Config.TrafficOps.UserPassword)

	for _, ext := range testData.ServerCheckExtensions {
		resp, _, err := TOSession.CreateServerCheckExtension(ext, client.RequestOptions{})
		if err != nil {
			t.Errorf("could not create Servercheck Extension: %v - alerts: %+v", err, resp.Alerts)
		}
	}

	SwitchSession(toReqTimeout, Config.TrafficOps.URL, Config.TrafficOps.Users.Extension, Config.TrafficOps.UserPassword, Config.TrafficOps.Users.Admin, Config.TrafficOps.UserPassword)
}

func CreateTestInvalidServerCheckExtensions(t *testing.T) {
	if len(testData.ServerCheckExtensions) < 1 {
		t.Fatal("Need at least one Servercheck Extension to test invalid Servercheck Extension creation")
	}
	// Fail Attempt to Create ServerCheckExtension as non extension user
	_, _, err := TOSession.CreateServerCheckExtension(testData.ServerCheckExtensions[0], client.RequestOptions{})
	if err == nil {
		t.Error("expected to receive error with non extension user")
	}

	SwitchSession(toReqTimeout, Config.TrafficOps.URL, Config.TrafficOps.Users.Admin, Config.TrafficOps.UserPassword, Config.TrafficOps.Users.Extension, Config.TrafficOps.UserPassword)

	// Attempt to create another valid ServerCheckExtension and it should fail as there is no open slots
	toExt := tc.ServerCheckExtensionNullable{
		Name:                 util.StrPtr("MEM_CHECKER"),
		Version:              util.StrPtr("3.0.3"),
		InfoURL:              util.StrPtr("-"),
		ScriptFile:           util.StrPtr("mem.py"),
		ServercheckShortName: util.StrPtr("MC"),
		Type:                 util.StrPtr("CHECK_EXTENSION_MEM"),
	}
	_, _, err = TOSession.CreateServerCheckExtension(toExt, client.RequestOptions{})
	if err == nil {
		t.Error("expected to receive error with no open slots left")
	}

	// Attempt to create a TO Extension with an invalid type
	toExt.Type = util.StrPtr("INVALID_TYPE")
	_, _, err = TOSession.CreateServerCheckExtension(toExt, client.RequestOptions{})
	if err == nil {
		t.Error("expected to receive error with invalid TO extension type")
	}
	SwitchSession(toReqTimeout, Config.TrafficOps.URL, Config.TrafficOps.Users.Extension, Config.TrafficOps.UserPassword, Config.TrafficOps.Users.Admin, Config.TrafficOps.UserPassword)

}

func DeleteTestServerCheckExtensions(t *testing.T) {
	SwitchSession(toReqTimeout, Config.TrafficOps.URL, Config.TrafficOps.Users.Admin, Config.TrafficOps.UserPassword, Config.TrafficOps.Users.Extension, Config.TrafficOps.UserPassword)

	extensions, _, err := TOSession.GetServerCheckExtensions(client.RequestOptions{})
	if err != nil {
		t.Fatalf("could not get Servercheck Extensions: %v - alerts: %+v", err, extensions.Alerts)
	}

	ids := []int{}
	for _, ext := range testData.ServerCheckExtensions {
		if ext.Name == nil {
			t.Errorf("Found Servercheck Extension in the testing data with null or undefined Name")
		}
		found := false
		for _, respTOExt := range extensions.Response {
			if respTOExt.Name == nil || respTOExt.ID == nil {
				t.Error("Traffic Ops returned a representation for a Servercheck Extension with null or undefined ID and/or name")
				continue
			}
			if *ext.Name == *respTOExt.Name {
				ids = append(ids, *respTOExt.ID)
				found = true
				continue
			}
		}
		if !found {
			t.Errorf("expected to find to_extension %v", *ext.Name)
		}
	}

	for _, id := range ids {
		resp, _, err := TOSession.DeleteServerCheckExtension(id, client.RequestOptions{})
		if err != nil {
			t.Errorf("cannot delete Servercheck Extension #%d: %v - alerts: %+v", id, err, resp.Alerts)
		}
	}
	extensions, _, err = TOSession.GetServerCheckExtensions(client.RequestOptions{})
	if err != nil {
		t.Fatalf("could not get to_extensions: %v", err)
	}

	for _, ext := range testData.ServerCheckExtensions {
		if ext.Name == nil {
			t.Errorf("Found Servercheck Extension in the testing data with null or undefined Name")
		}
		found := false
		for _, respTOExt := range extensions.Response {
			if respTOExt.Name == nil {
				t.Error("Traffic Ops returned a representation for a Servercheck Extension with null or undefined name")
				continue
			}
			if *ext.Name == *respTOExt.Name {
				found = true
				continue
			}
		}
		if found {
			t.Errorf("to_extension %v should have been deleted", *ext.Name)
		}
	}

	SwitchSession(toReqTimeout, Config.TrafficOps.URL, Config.TrafficOps.Users.Extension, Config.TrafficOps.UserPassword, Config.TrafficOps.Users.Admin, Config.TrafficOps.UserPassword)
}
