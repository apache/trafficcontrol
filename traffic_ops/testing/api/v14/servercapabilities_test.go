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

	"github.com/apache/trafficcontrol/lib/go-log"
)

func TestServerCapabilities(t *testing.T) {
	WithObjs(t, []TCObj{ServerCapabilities}, func() {
		GetTestServerCapabilities(t)
	})
}

func CreateTestServerCapabilities(t *testing.T) {
	log.Debugln("---- CreateTestServerCapabilities ----")

	for _, sc := range testData.ServerCapabilities {
		resp, _, err := TOSession.CreateServerCapability(sc)
		if err != nil {
			t.Errorf("could not CREATE server capability: %v\n", err)
		}
		log.Debugln("Response: ", resp)
	}

}

func GetTestServerCapabilities(t *testing.T) {
	log.Debugln("---- GetTestServerCapabilities ----")

	for _, sc := range testData.ServerCapabilities {
		resp, _, err := TOSession.GetServerCapability(sc.Name)
		if err != nil {
			t.Errorf("cannot GET server capability: %v - %v\n", err, resp)

		}

		log.Debugln("Response: ", resp)
	}

	resp, _, err := TOSession.GetServerCapabilities()
	if err != nil {
		t.Errorf("cannot GET server capabilities: %v\n", err)
	}
	if len(resp) != len(testData.ServerCapabilities) {
		t.Errorf("expected to GET %d server capabilities, actual: %d", len(testData.ServerCapabilities), len(resp))
	}
}

func DeleteTestServerCapabilities(t *testing.T) {
	log.Debugln("---- DeleteTestServerCapabilities ----")

	for _, sc := range testData.ServerCapabilities {
		delResp, _, err := TOSession.DeleteServerCapability(sc.Name)
		if err != nil {
			t.Errorf("cannot DELETE server capability: %v - %v\n", err, delResp)
		}

		serverCapability, _, err := TOSession.GetServerCapability(sc.Name)
		if err != nil {
			t.Errorf("error deleting server capability: %s\n", err.Error())
		}
		if len(serverCapability) > 0 {
			t.Errorf("expected server capability: %s to be deleted\n", sc.Name)
		}
	}
}
