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
)

// TestDirs tests that the t3c rpm creates directories that t3c-apply requires to run.
// Right now, that's just /var/lib/trafficcontrol-cache-config, but in the future it may include others like /var/log/trafficcontrol-cache-config or /etc/trafficcontrol-cache-config.
func TestDirs(t *testing.T) {
	t.Logf("------------- Starting TestDirs ---------------")

	requiredDirs := []string{
		`/var/lib/trafficcontrol-cache-config`,
	}

	for _, dir := range requiredDirs {
		dirInf, err := os.Stat(dir)
		if os.IsNotExist(err) {
			t.Errorf("directory '%v' must exist for t3c-apply to run successfully, expected: exists, actual: not found", dir)
		} else if err != nil {
			t.Errorf("checking if directory '%v' exists, expected: no error, actual: %v", dir, err)
		} else if !dirInf.IsDir() {
			t.Errorf("directory '%v' must exist for t3c-apply to run successfully, expected: exists, actual: is a file, not a directory", dir)
		} else {
			t.Logf("successfully verified '%v' exists and is a directory", dir)
		}
	}

	t.Logf("------------- End of TestDirs ---------------")
}
