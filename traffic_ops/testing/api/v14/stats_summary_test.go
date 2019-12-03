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
	"fmt"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
)

var (
	testStatsSummaries []tc.StatsSummary
	latestTime         time.Time
)

const dateFormat = "2006-01-02"

func TestStatsSummary(t *testing.T) {
	WithObjs(t, []TCObj{}, func() {
		testStatsSummaries = []tc.StatsSummary{}
		latestTime = time.Now().Truncate(time.Second).UTC()
		CreateTestStatsSummaries(t)
		GetTestStatsSummaries(t)
		GetTestStatsSummariesLastUpdated(t)
	})
}

func CreateTestStatsSummaries(t *testing.T) {
	tmpTime := latestTime
	for _, ss := range testData.StatsSummaries {
		ss.SummaryTime = tmpTime
		timeStamp, err := time.Parse(dateFormat, tmpTime.Format(dateFormat))
		if err != nil {
			t.Errorf("error getting timestamp at beginning of the day: %v", err)
		}

		ss.StatDate = timeStamp

		_, _, err = TOSession.CreateSummaryStats(ss)
		if err != nil {
			t.Errorf("creating stats_summary %v: %v", ss.StatName, err)
		}

		tmpTime = tmpTime.AddDate(0, 0, -1)

		testStatsSummaries = append(testStatsSummaries, ss)
	}
}

func GetTestStatsSummaries(t *testing.T) {
	var testCases = []struct {
		description            string
		stat                   *string
		cdn                    *string
		ds                     *string
		expectedStatsSummaries []tc.StatsSummary
	}{
		{
			description:            "get all summary stats",
			expectedStatsSummaries: testStatsSummaries,
		},
		{
			description: "non-existant stat name",
			stat:        util.StrPtr("bogus"),
		},
		{
			description: "non-existant ds name",
			ds:          util.StrPtr("bogus"),
		},
		{
			description: "non-existant cdn name",
			cdn:         util.StrPtr("bogus"),
		},
		{
			description: "get stats summary by stat name",
			stat:        util.StrPtr("daily_bytesserved"),
			expectedStatsSummaries: func() []tc.StatsSummary {
				statsSummaries := []tc.StatsSummary{}
				for _, ss := range testStatsSummaries {
					if ss.StatName == "daily_bytesserved" {
						statsSummaries = append(statsSummaries, ss)
					}
				}
				return statsSummaries
			}(),
		},
		{
			description:            "get stats summary by cdn name",
			cdn:                    util.StrPtr("cdn1"),
			expectedStatsSummaries: testStatsSummaries,
		},
		{
			description:            "get stats summary by ds name",
			ds:                     util.StrPtr("all"),
			expectedStatsSummaries: testStatsSummaries,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			tsr, _, err := TOSession.GetSummaryStats(tc.cdn, tc.ds, tc.stat)
			if err != nil {
				t.Fatalf("received unexpectant error %v on GET to stats_summary", err)
			}
			if len(tc.expectedStatsSummaries) == 0 && len(tsr.Response) != 0 {
				t.Fatalf("expected to recieve no stats summaries but received %v", len(tsr.Response))
			}
			for _, ess := range tc.expectedStatsSummaries {
				found := false
				for _, ss := range tsr.Response {
					if ess == ss {
						found = true
						break
					}
				}
				if !found {
					t.Fatalf("expected to find stat %v in stats summary response", ess.StatName)
				}
			}
		})
	}
}

func GetTestStatsSummariesLastUpdated(t *testing.T) {
	type testCase struct {
		description       string
		stat              *string
		errExpected       bool
		expectedTimestamp time.Time
	}
	testCases := []testCase{
		testCase{
			description:       "latest updated timestamp",
			stat:              nil,
			errExpected:       false,
			expectedTimestamp: latestTime,
		},
		testCase{
			description: "non-existant stat name",
			stat:        util.StrPtr("bogus"),
			errExpected: true,
		},
	}
	for _, ss := range testStatsSummaries {
		testCases = append(testCases, testCase{
			description:       fmt.Sprintf("latest updated timestamp for - %v", ss.StatName),
			stat:              util.StrPtr(ss.StatName),
			errExpected:       false,
			expectedTimestamp: ss.SummaryTime,
		})
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			tsr, _, err := TOSession.GetSummaryStatsLastUpdated(tc.stat)
			if tc.errExpected && err == nil {
				t.Fatalf("expected to get error on getting stats_summary latest updated timestamp but recieved nil")
			}

			if !tc.errExpected && err != nil {
				t.Fatalf("recieved unexpected error getting stats_summary latest updated timestamp: %v", err)
			}

			if !tc.errExpected && !tsr.Response.SummaryTime.Equal(tc.expectedTimestamp) {
				t.Fatalf("recieved latest timestamp %v does not match up to expected timestamp %v", tsr.Response.SummaryTime, tc.expectedTimestamp)
			}
		})
	}

}
