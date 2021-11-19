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
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/cache-config/testing/ort-tests/tcdata"
	"github.com/apache/trafficcontrol/cache-config/testing/ort-tests/util"
)

var recordsConfigFileName = filepath.Join(test_config_dir, "records.config")

func writeComment(t *testing.T) {
	f, err := os.OpenFile(recordsConfigFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("opening file '%s': %v", recordsConfigFileName, err)
	}
	defer f.Close()
	_, err = f.Write([]byte(" #mycomment\n"))
	if err != nil {
		t.Fatalf("writing comment to file '%s': %v", recordsConfigFileName, err)
	}

}

func writeNonComment(t *testing.T) {
	f, err := os.OpenFile(recordsConfigFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("opening file '%s': %v", recordsConfigFileName, err)
	}
	_, err = f.Write([]byte("\nmynocomment this line isn't a comment\n"))
	f.Close()
	if err != nil {
		t.Fatalf("writing line to file '%s': %v", recordsConfigFileName, err)
	}
}

func verifyOverwritten(t *testing.T) {
	recordsDotConfig, err := ioutil.ReadFile(recordsConfigFileName)
	if err != nil {
		t.Fatalf("reading %s: %v", recordsConfigFileName, err)
	}
	content := string(recordsDotConfig)
	if strings.Contains(content, "#mycomment") {
		t.Fatalf("expected records.config to have a diff and be replaced with a non-comment difference, actual: %s", content)
	} else if strings.Contains(content, "mynocomment") {
		t.Fatalf("expected records.config to have a diff and be replaced with a non-comment difference, actual: %s", content)
	}
}

func TestApplyDiff(t *testing.T) {
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters, tcdata.Statuses,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices}, func() {

		// badass to get initial config files

		if out, code := t3cUpdateReload(cacheHostName, "badass"); code != 0 {
			t.Fatalf("t3c apply badass failed with exit code %d, output: %s", code, out)
		}

		if !util.FileExists(recordsConfigFileName) {
			t.Fatalf("missing config file '%s' needed to test", recordsConfigFileName)
		}

		t.Run("write a comment, which should diff as unchanged", writeComment)

		// queue and syncds to get changes

		if err := ExecTOUpdater(cacheHostName, false, true); err != nil {
			t.Fatalf("updating queue status failed: %v", err)
		}
		out, code := t3cUpdateReload(cacheHostName, "syncds")
		if code != 0 {
			t.Fatalf("t3c apply failed with exit code %d, output: %s", code, out)
		}

		// verify the file wasn't overwritten, as it would be if there were a diff

		recordsDotConfig, err := ioutil.ReadFile(recordsConfigFileName)
		if err != nil {
			t.Fatalf("reading %s: %v", recordsConfigFileName, err)
		}
		if !bytes.Contains(recordsDotConfig, []byte("#mycomment")) {
			t.Fatalf("expected records.config to diff clean and not be replaced with comment difference, actual: '%s' t3c-apply output: %s", string(recordsDotConfig), out)
		}

		t.Run("write a non-comment, which should diff as changed", writeNonComment)

		// queue and syncds to get changes

		if err := ExecTOUpdater(cacheHostName, false, true); err != nil {
			t.Fatalf("updating queue status failed: %v", err)
		}
		out, code = t3cUpdateReload(cacheHostName, "syncds")
		if code != 0 {
			t.Fatalf("t3c apply failed with exit code %d, output: %s", code, out)
		}
		t.Logf("t3c apply output: %s", out)

		t.Run("verify the file was overwritten and our changes disappared, because there was a diff", verifyOverwritten)
	})
}
