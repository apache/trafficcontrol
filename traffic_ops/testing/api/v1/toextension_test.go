package v1

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

	"github.com/apache/trafficcontrol/lib/go-log"
)

var (
	toReqTimeout = time.Second * time.Duration(Config.Default.Session.TimeoutInSecs)
)

func TestTOExtensions(t *testing.T) {
	WithObjs(t, []TCObj{TOExtensions}, func() {
		CreateTestInvalidTOExtensions(t)
	})
}

func CreateTestTOExtensions(t *testing.T) {
	SwitchSession(toReqTimeout, Config.TrafficOps.URL, Config.TrafficOps.Users.Admin, Config.TrafficOps.UserPassword, Config.TrafficOps.Users.Extension, Config.TrafficOps.UserPassword)

	for _, ext := range testData.TOExtensions {
		resp, _, err := TOSession.CreateTOExtension(ext)
		log.Debugf("Response: %v %v", *ext.Name, resp)
		if err != nil {
			t.Errorf("could not create to_extension %v: %v", ext.Name, err)
		}
	}

	SwitchSession(toReqTimeout, Config.TrafficOps.URL, Config.TrafficOps.Users.Extension, Config.TrafficOps.UserPassword, Config.TrafficOps.Users.Admin, Config.TrafficOps.UserPassword)
}

func CreateTestInvalidTOExtensions(t *testing.T) {
	// Fail Attempt to Create ToExtension as non extension user
	_, _, err := TOSession.CreateTOExtension(testData.TOExtensions[0])
	if err == nil {
		t.Error("expected to receive error with non extension user")
	}

	SwitchSession(toReqTimeout, Config.TrafficOps.URL, Config.TrafficOps.Users.Admin, Config.TrafficOps.UserPassword, Config.TrafficOps.Users.Extension, Config.TrafficOps.UserPassword)

	// Attempt to create another valid TOExtension and it should fail as there is no open slots
	toExt := tc.TOExtensionNullable{
		Name:                 util.StrPtr("MEM_CHECKER"),
		Version:              util.StrPtr("3.0.3"),
		InfoURL:              util.StrPtr("-"),
		ScriptFile:           util.StrPtr("mem.py"),
		ServercheckShortName: util.StrPtr("MC"),
		Type:                 util.StrPtr("CHECK_EXTENSION_MEM"),
	}
	_, _, err = TOSession.CreateTOExtension(toExt)
	if err == nil {
		t.Error("expected to receive error with no open slots left")
	}

	// Attempt to create a TO Extension with an invalid type
	toExt.Type = util.StrPtr("INVALID_TYPE")
	_, _, err = TOSession.CreateTOExtension(toExt)
	if err == nil {
		t.Error("expected to receive error with invalid TO extension type")
	}
	SwitchSession(toReqTimeout, Config.TrafficOps.URL, Config.TrafficOps.Users.Extension, Config.TrafficOps.UserPassword, Config.TrafficOps.Users.Admin, Config.TrafficOps.UserPassword)

}

func DeleteTestTOExtensions(t *testing.T) {
	SwitchSession(toReqTimeout, Config.TrafficOps.URL, Config.TrafficOps.Users.Admin, Config.TrafficOps.UserPassword, Config.TrafficOps.Users.Extension, Config.TrafficOps.UserPassword)

	extensions, _, err := TOSession.GetTOExtensions()
	if err != nil {
		t.Fatalf("could not get to_extensions: %v", err)
	}

	ids := []int{}
	for _, ext := range testData.TOExtensions {
		found := false
		for _, respTOExt := range extensions.Response {
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
		resp, _, err := TOSession.DeleteTOExtension(id)
		log.Debugf("Response: %v %v", id, resp)
		if err != nil {
			t.Errorf("cannot delete to_extension: %v - %v", id, err)
		}
	}
	extensions, _, err = TOSession.GetTOExtensions()
	if err != nil {
		t.Fatalf("could not get to_extensions: %v", err)
	}

	for _, ext := range testData.TOExtensions {
		found := false
		for _, respTOExt := range extensions.Response {
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
