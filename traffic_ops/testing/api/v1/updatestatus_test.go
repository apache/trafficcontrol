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
)

func TestUpdateStatus(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers}, func() {
		GetTestUpdateStatus(t)
	})
}

func GetTestUpdateStatus(t *testing.T) {
	if len(testData.Servers) < 1 {
		t.Fatal("cannot GET Server: no test data")
	}
	testServer := testData.Servers[0]

	if _, _, err := TOSession.GetServerUpdateStatus(testServer.HostName); err != nil {
		t.Errorf("GetServerUpdateStatus error expected: nil, actual: %v", err)
	}
}
