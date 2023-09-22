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
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/cache-config/testing/ort-tests/tcdata"
)

func TestLockfile(t *testing.T) {
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices}, func() {

		const fileName = `records.config`

		firstOut := ""
		firstCode := 0
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			t.Logf("TestLockFile first t3c starting %s", time.Now())
			firstOut, firstCode = t3cUpdateReload(DefaultCacheHostName, "badass")
			t.Logf("TestLockFile first t3c finished %s", time.Now())
		}()

		time.Sleep(time.Millisecond * 100) // sleep long enough to ensure the concurrent t3c starts
		t.Logf("TestLockFile second t3c starting %s", time.Now())
		out, code := t3cUpdateReload(DefaultCacheHostName, "badass")
		t.Logf("TestLockFile second t3c finished %s", time.Now())
		if code != 0 {
			t.Fatalf("second t3c apply badass failed with exit code %d, output: %s", code, out)
		}

		wg.Wait()
		if firstCode != 0 {
			t.Fatalf("first t3c apply badass failed with exit code %d, output: %s", firstCode, firstOut)
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
			t.Fatalf("t3c apply output expected to contain 'Trying to acquire app lock' and 'Acquired app lock', actual: %s", out)
		}

		acquireStart, err := parseLogLineTime(acquireStartLine)
		if acquireStart == nil || err != nil {
			t.Fatalf("t3c apply acquire line failed to parse time, line '%s': %v", acquireStartLine, err)
		}
		acquireEnd, err := parseLogLineTime(acquireEndLine)
		if acquireEnd == nil || err != nil {
			t.Fatalf("t3c apply acquire line failed to parse time, line '%s': %v", acquireEndLine, err)
		}

		minDiff := time.Second * 1 // checking the file lock should never take 1s, so that's enough to verify it was hit
		if diff := acquireEnd.Sub(*acquireStart); diff < minDiff {
			t.Fatalf("t3c apply expected time to acquire while another t3c is running to be at least %s, actual %s start line '%s' end line '%s'", minDiff, diff, acquireStartLine, acquireEndLine)
		}

	})
}

func parseLogLineTime(line string) (*time.Time, error) {
	fields := strings.Fields(line)
	if len(fields) < 3 {
		return nil, fmt.Errorf("expected at least 3 fields, got: %d", len(fields))
	}
	timeStr := fields[2]
	timeStr = strings.TrimSuffix(timeStr, ":")
	tm, err := time.Parse(time.RFC3339Nano, timeStr)
	if err != nil {
		return nil, fmt.Errorf("parsing '%s' as an RFC3339 timestamp: %w", timeStr, err)
	}
	return &tm, nil
}
