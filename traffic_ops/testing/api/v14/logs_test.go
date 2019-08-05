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
)

func TestLogs(t *testing.T) {
	WithObjs(t, []TCObj{}, func() {
		GetTestLogs(t)
		GetTestLogsByLimit(t)
	})
}

func GetTestLogs(t *testing.T) {
	_, _, err := TOSession.GetLogs()
	if err != nil {
		t.Fatalf("error getting logs: " + err.Error())
	}
}

func GetTestLogsByLimit(t *testing.T) {
	toLogs, _, err := TOSession.GetLogsByLimit(10)
	if err != nil {
		t.Fatalf("error getting logs: " + err.Error())
	}
	if len(toLogs) != 10 {
		t.Fatalf("GET logs by limit: incorrect number of logs returned\n")
	}
}
