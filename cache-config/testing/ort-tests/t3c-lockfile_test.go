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
	"sync"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/cache-config/testing/ort-tests/tcdata"
)

func TestLockfile(t *testing.T) {
	testName := "TestLockfile"
	t.Logf("------------- Starting " + testName + " ---------------")
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters, tcdata.Statuses,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices}, func() {

		const hostName = `atlanta-edge-03`
		const fileName = `records.config`

		firstOut := ""
		firstCode := 0
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			t.Logf("TestLockFile first t3c starting %s", time.Now())
			firstOut, firstCode = t3cUpdateReload(hostName, "badass")
			t.Logf("TestLockFile first t3c finished %s", time.Now())
		}()

		time.Sleep(time.Millisecond * 100) // sleep long enough to ensure the concurrent t3c starts
		t.Logf("TestLockFile second t3c starting %s", time.Now())
		out, code := t3cUpdateReload(hostName, "badass")
		t.Logf("TestLockFile second t3c finished %s", time.Now())
		if code != 0 {
			t.Fatalf("second t3c apply badass failed: output '''%s''' code %d", out, code)
		}

		wg.Wait()
		if firstCode != 0 {
			t.Fatalf("first t3c apply badass failed: output '''%s''' code %d", firstOut, firstCode)
		}

		outStr := string(out)
		outLines := strings.Split(outStr, "\n")
		acquireStartLine := ""
		acquireEndLine := ""
		for _, line := range outLines {
			if strings.Contains(line, "Trying to acquire app lock") {
				acquireStartLine = line
			} else if strings.Contains(line, "Acquired app lock") {
				acquireEndLine = line
			}
			if acquireStartLine != "" && acquireEndLine != "" {
				break
			}
		}

		if acquireStartLine == "" || acquireEndLine == "" {
			t.Fatalf("t3c apply output expected to contain 'Trying to acquire app lock' and 'Acquired app lock', actual: '''%s'''", out)
		}

		acquireStart := parseLogLineTime(acquireStartLine)
		if acquireStart == nil {
			t.Fatalf("t3c apply acquire line failed to parse time, line '" + acquireStartLine + "'")
		}
		acquireEnd := parseLogLineTime(acquireEndLine)
		if acquireEnd == nil {
			t.Fatalf("t3c apply acquire line failed to parse time, line '" + acquireEndLine + "'")
		}

		minDiff := time.Second * 1 // checking the file lock should never take 1s, so that's enough to verify it was hit
		if diff := acquireEnd.Sub(*acquireStart); diff < minDiff {
			t.Fatalf("t3c apply expected time to acquire while another t3c is running to be at least %s, actual %s start line '%s' end line '%s'", minDiff, diff, acquireStartLine, acquireEndLine)
		}

		t.Logf(testName + " succeeded")
	})
	t.Logf("------------- End of " + testName + " ---------------")
}

func parseLogLineTime(line string) *time.Time {
	fields := strings.Fields(line)
	if len(fields) < 3 {
		return nil
	}
	timeStr := fields[2]
	timeStr = strings.TrimSuffix(timeStr, ":")
	tm, err := time.Parse(time.RFC3339Nano, timeStr)
	if err != nil {
		return nil
	}
	return &tm
}
