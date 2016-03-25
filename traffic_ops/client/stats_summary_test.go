/*
   Copyright 2015 Comcast Cable Communications Management, LLC

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

package client

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestStatsSummary(t *testing.T) {
	fmt.Println("Running StatsSummary Tests")
	text, err := os.Open("testdata/stats_summary.json")
	if err != nil {
		t.Skip("Skipping stats_summary test, no stats_summary.json found.")
	}

	var data StatsSummaryResponse
	if err := json.NewDecoder(text).Decode(&data); err != nil {
		t.Fatal(err)
	}

	for _, summaryStat := range data.Response {
		statName := summaryStat.StatName
		if len(statName) == 0 {
			t.Fatal("statSummary result does not contain 'StatName'")
		}
		statValue := summaryStat.StatValue
		if len(statValue) == 0 {
			t.Fatal("statSummary result does not contain 'StatValue'")
		}
		if len(summaryStat.CdnName) == 0 {
			t.Errorf("CdnName is null for summaryStat with name %v and value %v", statName, statValue)
		}
		if len(summaryStat.DeliveryService) == 0 {
			t.Errorf("DeliverySerice is null for summaryStat with name %v and value %v", statName, statValue)
		}
		if len(summaryStat.SummaryTime) == 0 {
			t.Errorf("SummaryTime is null for summaryStat with name %v and value %v", statName, statValue)
		}
	}
}
