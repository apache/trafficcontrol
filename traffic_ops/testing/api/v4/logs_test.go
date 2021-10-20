package v4

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

	client "github.com/apache/trafficcontrol/v6/traffic_ops/v4-client"
)

func TestLogs(t *testing.T) {
	WithObjs(t, []TCObj{Roles, Tenants, Users}, func() { // Objs added to create logs when this test is run alone
		GetTestLogs(t)
		GetTestLogsByLimit(t)
		GetTestLogsByUsername(t)
	})
}

func GetTestLogs(t *testing.T) {
	resp, _, err := TOSession.GetLogs(client.RequestOptions{})
	if err != nil {
		t.Fatalf("error getting logs: %v - alerts: %+v", err, resp.Alerts)
	}
}

func GetTestLogsByLimit(t *testing.T) {
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("limit", "10")
	toLogs, _, err := TOSession.GetLogs(opts)
	if err != nil {
		t.Errorf("error getting logs: %v - alerts: %+v", err, toLogs.Alerts)
	}
	if len(toLogs.Response) != 10 {
		t.Fatalf("Get logs by limit: incorrect number of logs returned (%d)", len(toLogs.Response))
	}
}

func GetTestLogsByUsername(t *testing.T) {
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("username", "admin")
	toLogs, _, err := TOSession.GetLogs(opts)
	if err != nil {
		t.Errorf("error getting logs: %v - alerts: %+v", err, toLogs.Alerts)
	}
	if len(toLogs.Response) <= 0 {
		t.Fatalf("Get logs by username: incorrect number of logs returned (%d)", len(toLogs.Response))
	}
	for _, user := range toLogs.Response {
		if *user.User != TOSession.UserName {
			t.Errorf("incorrect username seen in logs, expected: `admin`, got: %v", *user.User)
		}
	}
}
