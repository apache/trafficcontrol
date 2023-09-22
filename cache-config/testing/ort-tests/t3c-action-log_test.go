package orttest

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

	"github.com/apache/trafficcontrol/v8/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/v8/cache-config/testing/ort-tests/tcdata"
)

func TestT3cActionLog(t *testing.T) {
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices}, func() {

		stdErr, exitCode := t3cUpdateReload(DefaultCacheHostName, "badass")
		if exitCode != 0 {
			t.Fatalf("t3c badass failed, code: %v", exitCode)
		}

		startActionMsg := `ACTION='` + t3cutil.ActionLogActionApplyStart + `'`
		if !strings.Contains(stdErr, string(startActionMsg)) {
			t.Errorf("expected log to contain action log message '%s', actual: %s", startActionMsg, string(stdErr))
		}

	})
}
