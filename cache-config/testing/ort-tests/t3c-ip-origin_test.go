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
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/v8/cache-config/testing/ort-tests/tcdata"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
)

func TestT3CIPOrigin(t *testing.T) {
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices}, func() {

		out, err := t3cUpdateWaitForParents(DefaultCacheHostName, "badass", util.StrPtr("false"))
		if err != nil {
			t.Fatalf("t3c badass failed: %v", err)
		}
		t.Logf("t3c badass output: %s", out)

		fileName := filepath.Join(TestConfigDir, "parent.config")
		parentConfigBts, err := ioutil.ReadFile(fileName)
		if err != nil {
			t.Fatalf("reading %s: %v", fileName, err)
		}
		parentConfig := string(parentConfigBts)
		t.Logf("parentConfig full contents: %s", parentConfig)
		if !strings.Contains(parentConfig, "dest_ip=192.0.2.1") {
			t.Errorf("expected parent.config to contain dest_ip for fixture IP origin, actual: %s", parentConfig)
		}
		if strings.Contains(parentConfig, "dest_domain=192.0.2.1") {
			t.Errorf("expected parent.config to not contain dest_domain for IP origin, actual: %s", parentConfig)
		}

	})
}
