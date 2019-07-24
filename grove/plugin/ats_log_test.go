package plugin

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
	"testing"
	"time"
)

func TestATSLogTimeFractionalSeconds(t *testing.T) {
	testTimes := []int64{
		1563936732547355432,
		1563937732000355432,
		1563936732000000000,
		1463136732999000000,
		1563916732009000000,
		1503936232090000000,
		1563936732099000000,
		1563936722900000000,
		1563236282909000000,
	}
	for _, testTime := range testTimes {
		timestamp := time.Unix(0, testTime)

		logStr := atsEventLogStr(timestamp, "", "", "", "", "", "", "", "", "", 0, 0, 0, 0, 0, false, false, "", "", "", "", "", 0)

		logFields := strings.Fields(logStr)
		if len(logFields) < 1 {
			t.Fatalf("atsEventLogStr expected >1 fields, actual %v", len(logFields))
		}

		timeField := logFields[0]

		// the time field should be the Unix timestamp in seconds, as a float with 3 decimal places.
		unixNano := timestamp.UnixNano()
		unixSec := float64(unixNano) / float64(NSPerSec)
		unixSecThreeDecimalPts := fmt.Sprintf("%.3f", unixSec)

		if timeField != unixSecThreeDecimalPts {
			t.Errorf("atsEventLogStr time expected '%v' actual '%v'", unixSecThreeDecimalPts, timeField)
		}
	}
}
