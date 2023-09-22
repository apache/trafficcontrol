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
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/v8/cache-config/testing/ort-tests/tcdata"
	"github.com/apache/trafficcontrol/v8/cache-config/testing/ort-tests/util"
)

func TestApplyDiff(t *testing.T) {
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices}, func() {

		// badass to get initial config files
		if out, code := t3cUpdateReload(DefaultCacheHostName, "badass"); code != 0 {
			t.Fatalf("t3c apply badass failed with exit code %d, output: %s", code, out)
		}

		if !util.FileExists(RecordsConfigFileName) {
			t.Fatalf("missing config file '%s' needed to test", RecordsConfigFileName)
		}

		t.Run("verify comment is unchanged", func(t *testing.T) {
			f, err := os.OpenFile(RecordsConfigFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				t.Fatalf("opening file '%s': %v", RecordsConfigFileName, err)
			}
			defer f.Close()
			_, err = f.Write([]byte(" #mycomment\n"))
			if err != nil {
				t.Fatalf("writing comment to file '%s': %v", RecordsConfigFileName, err)
			}
			// queue and syncds to get changes

			err = tcd.QueueUpdatesForServer(DefaultCacheHostName, true)
			if err != nil {
				t.Fatalf("failed to queue updates: %v", err)
			}
			out, code := t3cUpdateReload(DefaultCacheHostName, "syncds")
			if code != 0 {
				t.Fatalf("t3c apply failed with exit code %d, output: %s", code, out)
			}

			// verify the file wasn't overwritten, as it would be if there were a diff

			recordsDotConfig, err := ioutil.ReadFile(RecordsConfigFileName)
			if err != nil {
				t.Fatalf("reading %s: %v", RecordsConfigFileName, err)
			}
			if !bytes.Contains(recordsDotConfig, []byte("#mycomment")) {
				t.Fatalf("expected records.config to diff clean and not be replaced with comment difference, actual: '%s' t3c-apply output: %s", string(recordsDotConfig), out)
			}
		})

		t.Run("verify non-comment is overwritten", func(t *testing.T) {
			f, err := os.OpenFile(RecordsConfigFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				t.Fatalf("opening file '%s': %v", RecordsConfigFileName, err)
			}
			_, err = f.Write([]byte("\nmynocomment this line isn't a comment\n"))
			f.Close()
			if err != nil {
				t.Fatalf("writing line to file '%s': %v", RecordsConfigFileName, err)
			}

			// queue and syncds to get changes

			err = tcd.QueueUpdatesForServer(DefaultCacheHostName, true)
			if err != nil {
				t.Fatalf("failed to queue updates: %v", err)
			}
			out, code := t3cUpdateReload(DefaultCacheHostName, "syncds")
			if code != 0 {
				t.Fatalf("t3c apply failed with exit code %d, output: %s", code, out)
			}
			t.Logf("t3c apply output: %s", out)

			recordsDotConfig, err := ioutil.ReadFile(RecordsConfigFileName)
			if err != nil {
				t.Fatalf("reading %s: %v", RecordsConfigFileName, err)
			}
			content := string(recordsDotConfig)
			if strings.Contains(content, "#mycomment") {
				t.Fatalf("expected records.config to have a diff and be replaced with a non-comment difference, actual: %s", content)
			} else if strings.Contains(content, "mynocomment") {
				t.Fatalf("expected records.config to have a diff and be replaced with a non-comment difference, actual: %s", content)
			}
		})
	})
}
