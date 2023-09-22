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

	"github.com/apache/trafficcontrol/v8/cache-config/testing/ort-tests/tcdata"
)

func TestT3cApplyFailMsg(t *testing.T) {
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices}, func() {
		t.Run("test t3c apply logging a run failure message", doTestT3cApplyFailMsg)
	})
}

func doTestT3cApplyFailMsg(t *testing.T) {
	// verifies that when no hostname is passed,
	// t3c-apply will get the OS hostname,
	// and use the short name (which is what will be in Traffic Ops),
	// and not the full FQDN.

	stdErr, exitCode := t3cUpdateReload("nonexistent-host-to-cause-failure", "badass")
	if exitCode == 0 {
		t.Fatalf("t3c-apply with nonexistent host expected failure, actual code %v stderr: %s", exitCode, stdErr)
	}

	errStr := strings.TrimSpace(string(stdErr))

	if !strings.HasSuffix(errStr, "CRITICAL FAILURE, ABORTING") {
		t.Fatalf("t3c-apply failure expected to end with critical failure message, actual code %v stderr: %s", exitCode, stdErr)
	}
}
