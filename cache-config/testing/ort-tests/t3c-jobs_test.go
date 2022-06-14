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
	"time"

	"github.com/apache/trafficcontrol/cache-config/testing/ort-tests/tcdata"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
)

func TestT3CJobs(t *testing.T) {
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices}, func() {

		now := tc.Time{Time: time.Now().Add(time.Minute)}
		dsi := (interface{})("ds1")
		ttli := (interface{})(72.0)
		job := tc.InvalidationJobInput{
			DeliveryService: &dsi,
			Regex:           util.StrPtr(`/refetch-test\.png##REFETCH##`),
			StartTime:       &now,
			TTL:             &ttli,
		}
		if _, _, err := tcdata.TOSession.CreateInvalidationJob(job); err != nil {
			t.Fatalf("create refetch job failed: %v", err)
		}
		job.Regex = util.StrPtr(`/refresh-test\.png`)
		if _, _, err := tcdata.TOSession.CreateInvalidationJob(job); err != nil {
			t.Fatalf("create refresh job failed: %v", err)
		}

		out, err := t3cUpdateWaitForParents(DefaultCacheHostName, "badass", util.StrPtr("false"))
		if err != nil {
			t.Fatalf("t3c badass failed: %v", err)
		}
		t.Logf("t3c badass output: %s", out)

		fileName := filepath.Join(TestConfigDir, "regex_revalidate.config")
		revalFileBts, err := ioutil.ReadFile(fileName)
		if err != nil {
			t.Fatalf("reading %s: %v", fileName, err)
		}
		revalFile := string(revalFileBts)
		t.Logf("revalFile full contents: %s", revalFile)
		lines := strings.Split(revalFile, "\n")
		sawRefresh := false
		sawRefetch := false
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue // blank lines are ok
			}
			if strings.Contains(line, "refresh-test") {
				sawRefresh = true
				if strings.HasSuffix(line, "STALE") || strings.HasSuffix(line, "MISS") {
					t.Errorf("expected regex_revalidate.config refresh-test line to contain no type, actual: %s", line)
				}
			}
			if strings.Contains(line, "refetch-test") {
				sawRefetch = true
				if !strings.HasSuffix(line, "MISS") {
					t.Errorf("expected regex_revalidate.config refetch-test line to contain 'MISS', actual: %s", line)
				}
			}
		}
		if !sawRefresh {
			t.Error("expected regex_revalidate.config to contain refresh-test")
		}
		if !sawRefetch {
			t.Error("expected regex_revalidate.config to contain refetch-test")
		}
	})
}
