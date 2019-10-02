package v14

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
			t.Errorf("could not create to_extension %v: %v\n", ext.Name, err)
		}
	}

	SwitchSession(toReqTimeout, Config.TrafficOps.URL, Config.TrafficOps.Users.Extension, Config.TrafficOps.UserPassword, Config.TrafficOps.Users.Admin, Config.TrafficOps.UserPassword)
}

func CreateTestInvalidTOExtensions(t *testing.T) {
	// Fail Attempt to Create ToExtension as non extension user
	_, _, err := TOSession.CreateTOExtension(testData.TOExtensions[0])
	if err == nil {
		t.Errorf("expected to receive error with non extension user\n")
	}
}

func DeleteTestTOExtensions(t *testing.T) {
	SwitchSession(toReqTimeout, Config.TrafficOps.URL, Config.TrafficOps.Users.Admin, Config.TrafficOps.UserPassword, Config.TrafficOps.Users.Extension, Config.TrafficOps.UserPassword)

	extensions, _, err := TOSession.GetTOExtensions()
	if err != nil {
		t.Errorf("could not get to_extensions: %v\n", err)
	}
	if len(extensions.Response) != len(testData.TOExtensions) {
		t.Errorf("%v to_extensions returned - expected %v\n", len(extensions.Response), len(testData.TOExtensions))
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
			t.Errorf("expected to find to_extension %v\n", *ext.Name)
		}
	}

	for _, id := range ids {
		resp, _, err := TOSession.DeleteTOExtension(id)
		log.Debugf("Response: %v %v", id, resp)
		if err != nil {
			t.Errorf("cannot delete to_extension: %v - %v\n", id, err)
		}
	}
	extensions, _, err = TOSession.GetTOExtensions()
	if err != nil {
		t.Errorf("could not get to_extensions: %v\n", err)
	}
	if len(extensions.Response) != 0 {
		t.Errorf("%v to_extensions returned - expected %v\n", len(extensions.Response), 0)
	}
	SwitchSession(toReqTimeout, Config.TrafficOps.URL, Config.TrafficOps.Users.Extension, Config.TrafficOps.UserPassword, Config.TrafficOps.Users.Admin, Config.TrafficOps.UserPassword)
}
