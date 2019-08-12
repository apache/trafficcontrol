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
	"github.com/apache/trafficcontrol/lib/go-tc"
	toclient "github.com/apache/trafficcontrol/traffic_ops/client"
)

func TestServerCheck(t *testing.T) {
	WithObjs(t, []TCObj{Servers}, func() {
		UpdateServerCheck(t)
	})
}

const SessionUserName = "extension" // TODO make dynamic?

func UpdateServerCheck(t *testing.T) {
	var statusData tc.ServercheckNullable
	for _, server := range testData.Servers {

		statusData.ID = server.ID
		statusData.Name = "DSCP"
		statusData.Value = 1
		resp, _, err := TOSession.UpdateCheckStatus(statusData)
		if err != nil {
			t.Errorf("could not set servercheck status user: %v\n", err)
		}
		log.Debugln("Response: ", resp.Alerts)
	}
}
