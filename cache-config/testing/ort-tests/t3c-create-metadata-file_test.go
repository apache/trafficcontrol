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
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/v8/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/v8/cache-config/testing/ort-tests/tcdata"
	"github.com/apache/trafficcontrol/v8/cache-config/testing/ort-tests/util"
)

func TestT3cCreateMetaDataFile(t *testing.T) {
	// t3c should create a metadata file
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices}, func() {

		err := t3cUpdateCreateEmptyFile(DefaultCacheHostName, "badass")
		if err != nil {
			t.Fatalf("t3c badass failed: %v", err)
		}

		const metaDataFileName = `t3c-apply-metadata.json`

		filePath := filepath.Join(TestConfigDir, metaDataFileName)

		if !util.FileExists(filePath) {
			t.Fatalf("missing metadata file '%s'", filePath)
		}

		mdFileBts, err := ioutil.ReadFile(filePath)
		if err != nil {
			t.Fatalf("reading file '%s': %v", filePath, err)
		}

		// Test that the file is a valid JSON object.
		// Other than that, we don't want to assert any particular data.
		mdObj := map[string]interface{}{}
		if err := json.Unmarshal(mdFileBts, &mdObj); err != nil {
			t.Errorf("expected metadata file '%s' to be a json object, actual: %s", filePath, err)
		}

		metaDataStr := string(mdFileBts)
		if !strings.Contains(metaDataStr, `"actions":`) {
			t.Errorf("expected metadata file '%s' to contain action log, actual: %s", filePath, metaDataStr)
		}
		if !strings.Contains(metaDataStr, `"`+string(t3cutil.ActionLogActionApplyStart)+`"`) {
			t.Errorf("expected metadata file '%s' to contain action log apply start message '%v', actual: %s", filePath, t3cutil.ActionLogActionApplyStart, metaDataStr)
		}

	})
}
