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
	"testing"

	"github.com/apache/trafficcontrol/cache-config/testing/ort-tests/tcdata"
	"github.com/apache/trafficcontrol/cache-config/testing/ort-tests/util"
)

func TestApplyDiff(t *testing.T) {
	testName := "TestApplyDiff"
	t.Logf("------------- Starting " + testName + " ---------------")
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters, tcdata.Statuses,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices}, func() {

		hostName := "atlanta-edge-03"
		const fileName = `records.config`

		// badass to get initial config files

		if out, code := t3cUpdateReload(hostName, "badass"); code != 0 {
			t.Fatalf("t3c apply badass failed: output '''%v''' code %v\n", out, code)
		}

		filePath := test_config_dir + "/" + fileName

		if !util.FileExists(filePath) {
			t.Fatalf("missing config file '%v' needed to test", filePath)
		}

		// write a comment, which should diff as unchanged

		{
			f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				t.Fatalf("opening file '%v': %v", filePath, err)
			}
			_, err = f.Write([]byte(" #mycomment\n"))
			f.Close()
			if err != nil {
				t.Fatalf("writing comment to file '%v': %v", filePath, err)
			}
		}

		// queue and syncds to get changes

		if err := ExecTOUpdater(hostName, false, true); err != nil {
			t.Fatalf("updating queue status failed: %v\n", err)
		}
		out, code := t3cUpdateReload(hostName, "syncds")
		if code != 0 {
			t.Fatalf("t3c apply failed: output '''%v''' code %v\n", out, code)
		}

		// verify the file wasn't overwritten, as it would be if there were a diff

		{
			recordsDotConfig, err := ioutil.ReadFile(filePath)
			if err != nil {
				t.Fatalf("reading %v: %v\n", filePath, err)
			}
			if !bytes.Contains(recordsDotConfig, []byte("#mycomment")) {
				t.Fatalf("expected records.config to diff clean and not be replaced with comment difference, actual: '%v' t3c-apply output '''%v'''\n", string(recordsDotConfig), out)
			}
		}

		// write a non-comment, which should diff as changed

		{
			f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				t.Fatalf("opening file '%v': %v", filePath, err)
			}
			_, err = f.Write([]byte("\nmynocomment this line isn't a comment\n"))
			f.Close()
			if err != nil {
				t.Fatalf("writing line to file '%v': %v", filePath, err)
			}
		}

		// queue and syncds to get changes

		if err := ExecTOUpdater(hostName, false, true); err != nil {
			t.Fatalf("updating queue status failed: %v\n", err)
		}
		out, code = t3cUpdateReload(hostName, "syncds")
		if code != 0 {
			t.Fatalf("t3c apply failed: output '''%v''' code %v\n", out, code)
		}

		// verify the file was overwritten and our changes disappared, because there was a diff

		{
			recordsDotConfig, err := ioutil.ReadFile(filePath)
			if err != nil {
				t.Fatalf("reading %v: %v\n", filePath, err)
			}
			if bytes.Contains(recordsDotConfig, []byte("#mycomment")) {
				t.Fatalf("expected records.config to have a diff and be replaced with a non-comment difference, actual: '%v' t3c-apply output '''%v'''", string(recordsDotConfig), out)
			} else if bytes.Contains(recordsDotConfig, []byte("mynocomment")) {
				t.Fatalf("expected records.config to have a diff and be replaced with a non-comment difference, actual: '%v' t3c-apply output '''%v'''", string(recordsDotConfig), out)
			}
		}

		t.Logf(testName + " succeeded")

	})
	t.Logf("------------- End of " + testName + " ---------------")
}
