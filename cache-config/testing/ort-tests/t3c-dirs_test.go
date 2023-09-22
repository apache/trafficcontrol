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
	"os"
	"testing"

	"github.com/apache/trafficcontrol/v8/cache-config/t3c-apply/util"
)

// TestDirs tests that the t3c rpm creates directories that t3c-apply requires to run.
// Right now, that's just /var/lib/trafficcontrol-cache-config, but in the future it may include others like /var/log/trafficcontrol-cache-config or /etc/trafficcontrol-cache-config.
func TestDirs(t *testing.T) {
	requiredDirs := []string{
		`/var/lib/trafficcontrol-cache-config`,
	}

	for _, dir := range requiredDirs {
		dirInf, err := os.Stat(dir)
		if os.IsNotExist(err) {
			t.Errorf("directory '%s' must exist for t3c-apply to run successfully, expected: exists, actual: not found", dir)
		} else if err != nil {
			t.Errorf("checking if directory '%s' exists, expected: no error, actual: %v", dir, err)
		} else if !dirInf.IsDir() {
			t.Errorf("directory '%s' must exist for t3c-apply to run successfully, expected: exists, actual: is a file, not a directory", dir)
		} else {
			t.Logf("successfully verified '%s' exists and is a directory", dir)
		}
	}
}

// Test_MkDir ensures that the directory creation in t3c-apply succeeds and fails as expected
func Test_MkDir(t *testing.T) {

	os.WriteFile("isAFile", []byte("nothing to see here"), 0755)

	mkDirTests := []struct {
		description    string
		path           string
		reportOnly     bool
		expectedResult bool
	}{
		{"Success - Create a test directory", "testDir", false, true},
		{"Fail - Create a test directory, but file exists", "isAFile", false, false},
		{"Success - Report Only even if dir doesn't exist", "ShouldntExist", true, true},
	}

	for _, test := range mkDirTests {
		result := util.MkDir(test.path, test.reportOnly)
		if result != test.expectedResult {
			t.Errorf("Directory creation test failed: %s", test.description)
		}
	}

}
