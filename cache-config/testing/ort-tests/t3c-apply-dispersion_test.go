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
	"github.com/apache/trafficcontrol/cache-config/testing/ort-tests/tcdata"
	"testing"
	"time"
)

func TestDispersion(t *testing.T) {
	fmt.Println("------------- Starting TestDispersion ---------------")
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters, tcdata.Statuses,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices}, func() {

		const cacheHostName = `atlanta-edge-03`

		// get config once, because it'll take longer the first time

		if err := runApply(cacheHostName, "badass", 0); err != nil {
			t.Fatalf("ERROR: t3c badass failed: %v\n", err)
		}

		// get the average time to run with dispersion=0

		numAvgRuns := 5
		runSpans := make([]time.Duration, numAvgRuns)
		for i := 0; i < numAvgRuns; i++ {
			if err := ExecTOUpdater(cacheHostName, false, true); err != nil {
				t.Fatalf("t3c-update failed: %v\n", err)
			}
			start := time.Now()
			if err := runApply(cacheHostName, "syncds", 0); err != nil {
				t.Fatalf("ERROR: t3c badass failed: %v\n", err)
			}
			runSpans[i] = time.Since(start)
		}

		avgNoDispersionRunSpan := time.Duration(0)
		for i := 0; i < numAvgRuns; i++ {
			n := time.Duration(i + 1)
			avgNoDispersionRunSpan -= avgNoDispersionRunSpan / n
			avgNoDispersionRunSpan += runSpans[i] / n
		}

		t.Logf("t3c no-dispersion run times: %+v average %v\n", runSpans, avgNoDispersionRunSpan)

		// set the dispersion time to double the average runtime, or at least 2s

		dispersionValue := avgNoDispersionRunSpan * 4 // dispersion is randomly between 0 and d, so the amortized time here will be double, not 4x
		minDispersion := time.Second * 4
		if dispersionValue < minDispersion {
			dispersionValue = minDispersion
		}

		t.Logf("set dispersion to %v (average was %v)\n", dispersionValue, avgNoDispersionRunSpan)

		dispRunSpans := make([]time.Duration, numAvgRuns)
		for i := 0; i < numAvgRuns; i++ {
			if err := ExecTOUpdater(cacheHostName, false, true); err != nil {
				t.Fatalf("t3c-update failed: %v\n", err)
			}
			start := time.Now()
			if err := runApply(cacheHostName, "syncds", dispersionValue); err != nil {
				t.Fatalf("ERROR: t3c badass failed: %v\n", err)
			}
			dispRunSpans[i] = time.Since(start)
		}

		avgDispersionRunSpan := time.Duration(0)
		for i := 0; i < numAvgRuns; i++ {
			n := time.Duration(i + 1)
			avgDispersionRunSpan -= avgDispersionRunSpan / n
			avgDispersionRunSpan += dispRunSpans[i] / n
		}

		t.Logf("t3c dispersion run times: %+v average %v\n", dispRunSpans, avgDispersionRunSpan)

		if expected := avgNoDispersionRunSpan.Seconds() * 1.5; avgDispersionRunSpan.Seconds() < expected {
			t.Errorf("expected dispersion flag %v to average at least %v seconds, actual: %v dispersion runs averaged %v no dispersion averaged %v", dispersionValue, expected, numAvgRuns, avgDispersionRunSpan, avgNoDispersionRunSpan)
		} else {
			t.Logf("success: dispersion %v runs %v averaged %v dispersion=0 averaged %v, expectation met %v > %v", dispersionValue, numAvgRuns, avgDispersionRunSpan, avgNoDispersionRunSpan, avgDispersionRunSpan.Seconds(), expected)
		}

	})
	fmt.Println("------------- End of TestDispersion ---------------")
}
