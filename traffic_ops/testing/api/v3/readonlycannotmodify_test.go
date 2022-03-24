package v3

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
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
	toclient "github.com/apache/trafficcontrol/traffic_ops/v3-client"
)

func TestReadOnlyCannotModify(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, DeliveryServices, Users}, func() {
		CreateTestCDNWithReadOnlyUser(t)
	})
}

func CreateTestCDNWithReadOnlyUser(t *testing.T) {
	if len(testData.CDNs) == 0 {
		t.Fatal("Can't test readonly user creating a cdns: test data has no cdns")
	}

	toReqTimeout := time.Second * time.Duration(Config.Default.Session.TimeoutInSecs)
	readonlyTOClient, _, err := toclient.LoginWithAgent(TOSession.URL, "readonlyuser", "pa$$word", true, "to-api-v3-client-tests/readonlyuser", true, toReqTimeout)
	if err != nil {
		t.Fatalf("failed to get log in with readonlyuser: " + err.Error())
	}

	cdn := tc.CDN{
		Name:       "cdn-test-readonly-create-failure",
		DomainName: "createfailure.invalid",
	}

	alerts, _, err := readonlyTOClient.CreateCDN(cdn)

	if err == nil {
		t.Error("readonlyuser creating cdn error expected: not nil, actual: nil error")
	}

	if !strings.Contains(strings.ToLower(err.Error()), "forbidden") {
		t.Errorf("readonlyuser creating cdn error expected: contains 'forbidden', actual: '" + err.Error() + "'")
	}

	for _, alert := range alerts.Alerts {
		if alert.Level == tc.SuccessLevel.String() {
			t.Errorf("readonlyuser creating cdn, alerts expected: no success alert, actual: got success alert '" + alert.Text + "'")
		}
	}

	cdns, _, _ := TOSession.GetCDNByName(cdn.Name)
	if len(cdns) > 0 {
		t.Errorf("readonlyuser getting created cdn, len(cdns) expected: 0, actual: %+v %+v", len(cdns), cdns)
	}
}
