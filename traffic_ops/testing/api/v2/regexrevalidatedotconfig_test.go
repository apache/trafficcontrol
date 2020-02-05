package v2

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
	"testing"
)

func TestRegexRevalidateDotConfig(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, DeliveryServices}, func() {
		CreateTestJobs(t)
		GetTestRegexRevalidateDotConfig(t)
	})
}

func GetTestRegexRevalidateDotConfig(t *testing.T) {
	cdnName := "cdn1"

	cfg, _, err := TOSession.GetATSCDNConfigByName(cdnName, "regex_revalidate.config")
	if err != nil {
		t.Fatalf("Getting cdn '" + cdnName + "' config regex_revalidate.config: " + err.Error() + "\n")
	}

	for _, testJob := range testData.Jobs {
		if !strings.Contains(cfg, testJob.Request.Regex) {
			t.Errorf("regex_revalidate.config '''%+v''' expected: contains '%+v' actual: missing", cfg, testJob.Request.Regex)
		}
	}
}
