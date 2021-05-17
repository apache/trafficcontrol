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
)

func (r *TCData) CreateTestServerCapabilities(t *testing.T) {

	for _, sc := range r.TestData.ServerCapabilities {
		resp, _, err := TOSession.CreateServerCapability(sc)
		if err != nil {
			t.Errorf("could not CREATE server capability: %v", err)
		}
		t.Log("Response: ", resp)
	}

}

func (r *TCData) DeleteTestServerCapabilities(t *testing.T) {

	for _, sc := range r.TestData.ServerCapabilities {
		delResp, _, err := TOSession.DeleteServerCapability(sc.Name)
		if err != nil {
			t.Errorf("cannot DELETE server capability: %v - %v", err, delResp)
		}

		serverCapability, _, err := TOSession.GetServerCapability(sc.Name)
		if err == nil {
			t.Errorf("expected error trying to GET deleted server capability: %s, actual: nil", sc.Name)
		}
		if serverCapability != nil {
			t.Errorf("expected nil trying to GET deleted server capability: %s, actual: non-nil", sc.Name)
		}
	}
}
