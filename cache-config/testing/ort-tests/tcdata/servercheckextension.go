package tcdata

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
)

func (r *TCData) CreateTestServerCheckExtensions(t *testing.T) {
	toReqTimeout := time.Second * time.Duration(r.Config.Default.Session.TimeoutInSecs)
	r.SwitchSession(toReqTimeout, r.Config.TrafficOps.URL, r.Config.TrafficOps.Users.Admin, r.Config.TrafficOps.UserPassword, r.Config.TrafficOps.Users.Extension, r.Config.TrafficOps.UserPassword)

	for _, ext := range r.TestData.ServerCheckExtensions {
		resp, _, err := TOSession.CreateServerCheckExtension(ext)
		t.Logf("Response: %s %v", *ext.Name, resp)
		if err != nil {
			t.Errorf("could not create to_extension %s: %v", *ext.Name, err)
		}
	}

	r.SwitchSession(toReqTimeout, r.Config.TrafficOps.URL, r.Config.TrafficOps.Users.Extension, r.Config.TrafficOps.UserPassword, r.Config.TrafficOps.Users.Admin, r.Config.TrafficOps.UserPassword)
}

func (r *TCData) CreateTestInvalidServerCheckExtensions(t *testing.T) {
	toReqTimeout := time.Second * time.Duration(r.Config.Default.Session.TimeoutInSecs)
	// Fail Attempt to Create ServerCheckExtension as non extension user
	_, _, err := TOSession.CreateServerCheckExtension(r.TestData.ServerCheckExtensions[0])
	if err == nil {
		t.Error("expected to receive error with non extension user")
	}

	r.SwitchSession(toReqTimeout, r.Config.TrafficOps.URL, r.Config.TrafficOps.Users.Admin, r.Config.TrafficOps.UserPassword, r.Config.TrafficOps.Users.Extension, r.Config.TrafficOps.UserPassword)

	// Attempt to create another valid ServerCheckExtension and it should fail as there is no open slots
	toExt := tc.ServerCheckExtensionNullable{
		Name:                 util.StrPtr("MEM_CHECKER"),
		Version:              util.StrPtr("3.0.3"),
		InfoURL:              util.StrPtr("-"),
		ScriptFile:           util.StrPtr("mem.py"),
		ServercheckShortName: util.StrPtr("MC"),
		Type:                 util.StrPtr("CHECK_EXTENSION_MEM"),
	}
	_, _, err = TOSession.CreateServerCheckExtension(toExt)
	if err == nil {
		t.Error("expected to receive error with no open slots left")
	}

	// Attempt to create a TO Extension with an invalid type
	toExt.Type = util.StrPtr("INVALID_TYPE")
	_, _, err = TOSession.CreateServerCheckExtension(toExt)
	if err == nil {
		t.Error("expected to receive error with invalid TO extension type")
	}
	r.SwitchSession(toReqTimeout, r.Config.TrafficOps.URL, r.Config.TrafficOps.Users.Extension, r.Config.TrafficOps.UserPassword, r.Config.TrafficOps.Users.Admin, r.Config.TrafficOps.UserPassword)

}

func (r *TCData) DeleteTestServerCheckExtensions(t *testing.T) {
	toReqTimeout := time.Second * time.Duration(r.Config.Default.Session.TimeoutInSecs)
	r.SwitchSession(toReqTimeout, r.Config.TrafficOps.URL, r.Config.TrafficOps.Users.Admin, r.Config.TrafficOps.UserPassword, r.Config.TrafficOps.Users.Extension, r.Config.TrafficOps.UserPassword)

	extensions, _, err := TOSession.GetServerCheckExtensions()
	if err != nil {
		t.Fatalf("could not get to_extensions: %v", err)
	}

	ids := []int{}
	for _, ext := range r.TestData.ServerCheckExtensions {
		found := false
		for _, respTOExt := range extensions.Response {
			if *ext.Name == *respTOExt.Name {
				ids = append(ids, *respTOExt.ID)
				found = true
				continue
			}
		}
		if !found {
			t.Errorf("expected to find to_extension %s", *ext.Name)
		}
	}

	for _, id := range ids {
		resp, _, err := TOSession.DeleteServerCheckExtension(id)
		t.Logf("Response: %d %v", id, resp)
		if err != nil {
			t.Errorf("cannot delete to_extension: %d - %v", id, err)
		}
	}
	extensions, _, err = TOSession.GetServerCheckExtensions()
	if err != nil {
		t.Fatalf("could not get to_extensions: %v", err)
	}

	for _, ext := range r.TestData.ServerCheckExtensions {
		found := false
		for _, respTOExt := range extensions.Response {
			if *ext.Name == *respTOExt.Name {
				found = true
				continue
			}
		}
		if found {
			t.Errorf("to_extension %s should have been deleted", *ext.Name)
		}
	}

	r.SwitchSession(toReqTimeout, r.Config.TrafficOps.URL, r.Config.TrafficOps.Users.Extension, r.Config.TrafficOps.UserPassword, r.Config.TrafficOps.Users.Admin, r.Config.TrafficOps.UserPassword)
}
