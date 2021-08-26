package v5

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
	client "github.com/apache/trafficcontrol/traffic_ops/v5-client"
)

var (
	testStatsSummaries []tc.StatsSummary
	latestTime         time.Time
)

func TestStatsSummary(t *testing.T) {
	testStatsSummaries = []tc.StatsSummary{}
	latestTime = time.Now().Truncate(time.Second).UTC()
	CreateTestStatsSummaries(t)
	GetTestStatsSummaries(t)
	GetTestStatsSummariesLastUpdated(t)
}

// Note that these stats summaries are never cleaned up, and will be left in
// the TODB after the tests complete
func CreateTestStatsSummaries(t *testing.T) {
	tmpTime := latestTime
	for _, ss := range testData.StatsSummaries {
		ss.SummaryTime = tmpTime
		alerts, _, err := TOSession.CreateSummaryStats(ss, client.RequestOptions{})
		if err != nil {
			t.Errorf("creating Stats Summary for stat '%s': %v - alerts: %+v", *ss.StatName, err, alerts.Alerts)
		}

		tmpTime = tmpTime.AddDate(0, 0, -1)

		testStatsSummaries = append(testStatsSummaries, ss)
	}
}

func GetTestStatsSummaries(t *testing.T) {
	var testCases = []struct {
		description            string
		stat                   string
		cdn                    string
		ds                     string
		expectedStatsSummaries []tc.StatsSummary
	}{
		{
			description:            "get all summary stats",
			expectedStatsSummaries: testStatsSummaries,
		},
		{
			description: "non-existant stat name",
			stat:        "bogus",
		},
		{
			description: "non-existant ds name",
			ds:          "bogus",
		},
		{
			description: "non-existant cdn name",
			cdn:         "bogus",
		},
		{
			description: "get stats summary by stat name",
			stat:        "daily_bytesserved",
			expectedStatsSummaries: func() []tc.StatsSummary {
				statsSummaries := []tc.StatsSummary{}
				for _, ss := range testStatsSummaries {
					if ss.StatName == nil {
						t.Error("testing stats summaries collection contains a Stats Summary with nil StatName")
						continue
					}
					if *ss.StatName == "daily_bytesserved" {
						statsSummaries = append(statsSummaries, ss)
					}
				}
				return statsSummaries
			}(),
		},
		{
			description:            "get stats summary by cdn name",
			cdn:                    "cdn1",
			expectedStatsSummaries: testStatsSummaries,
		},
		{
			description:            "get stats summary by ds name",
			ds:                     "all",
			expectedStatsSummaries: testStatsSummaries,
		},
	}

	for _, tc := range testCases {
		opts := client.NewRequestOptions()
		t.Run(tc.description, func(t *testing.T) {
			if tc.cdn != "" {
				opts.QueryParameters.Set("cdnName", tc.cdn)
			}
			if tc.ds != "" {
				opts.QueryParameters.Set("deliveryServiceName", tc.ds)
			}
			if tc.stat != "" {
				opts.QueryParameters.Set("statName", tc.stat)
			}
			tsr, _, err := TOSession.GetSummaryStats(opts)
			if err != nil {
				t.Fatalf("Unexpected error getting Stats Summary: %v - alerts: %+v", err, tsr.Alerts)
			}
			if len(tc.expectedStatsSummaries) == 0 && len(tsr.Response) != 0 {
				t.Fatalf("expected to recieve no stats summaries but received %d", len(tsr.Response))
			}
			for _, ess := range tc.expectedStatsSummaries {
				if ess.StatName == nil {
					t.Error("testing stats summaries collection contains a Stats Summary with nil StatName")
					continue
				}
				found := false
				for _, ss := range tsr.Response {
					if ss.StatName == nil {
						t.Error("Traffic Ops returned a representation for a Stats Summary with null or undefined name")
						continue
					}
					if *ess.StatName == *ss.StatName && ess.SummaryTime.Equal(ss.SummaryTime) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected to find stat '%s' in stats summary response", *ess.StatName)
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
		nullTimeStamp     bool
	}
	testCases := []testCase{
		{
			description:       "latest updated timestamp",
			stat:              nil,
			errExpected:       false,
			expectedTimestamp: latestTime,
		},
		{
			description:   "non-existant stat name",
			stat:          util.StrPtr("bogus"),
			errExpected:   false,
			nullTimeStamp: true,
		},
	}
	for _, ss := range testStatsSummaries {
		if ss.StatName == nil {
			t.Error("testing stats summaries collection contains a Stats Summary with nil StatName")
			continue
		}
		testCases = append(testCases, testCase{
			description:       fmt.Sprintf("latest updated timestamp for - %v", *ss.StatName),
			stat:              ss.StatName,
			errExpected:       false,
			expectedTimestamp: ss.SummaryTime,
		})
	}

	for _, tc := range testCases {
		opts := client.NewRequestOptions()
		t.Run(tc.description, func(t *testing.T) {
			if tc.stat != nil {
				opts.QueryParameters.Set("statName", *tc.stat)
			}
			tsr, _, err := TOSession.GetSummaryStatsLastUpdated(opts)
			if tc.errExpected {
				if err == nil {
					t.Fatalf("expected to get error on getting stats_summary latest updated timestamp but received nil")
				}
			} else if err != nil {
				t.Fatalf("received unexpected error getting Stats Summary latest updated timestamp: %v - alerts: %+v", err, tsr.Alerts)
			} else if tc.nullTimeStamp {
				if tsr.Response.SummaryTime != nil {
					t.Fatalf("expected to get null on latest timestamp but instead got %v", *tsr.Response.SummaryTime)
				}
			} else if !tsr.Response.SummaryTime.Equal(tc.expectedTimestamp) {
				t.Fatalf("received latest timestamp %v does not match up to expected timestamp %v", tsr.Response.SummaryTime, tc.expectedTimestamp)
			}
		})
	}
}
