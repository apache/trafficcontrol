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

	"github.com/apache/trafficcontrol/lib/go-tc"
)

func TestDeliveryServiceMatches(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, DeliveryServices}, func() {
		GetTestDeliveryServiceMatches(t)
	})
}

func GetTestDeliveryServiceMatches(t *testing.T) {
	dsMatches, _, err := TOSession.GetDeliveryServiceMatches()
	if err != nil {
		t.Errorf("cannot GET DeliveryService matches: %v\n", err)
	}

	dsMatchMap := map[tc.DeliveryServiceName][]string{}
	for _, ds := range dsMatches {
		dsMatchMap[ds.DSName] = ds.Patterns
	}

	for _, ds := range testData.DeliveryServices {
		if ds.Type == tc.DSTypeAnyMap || len(ds.MatchList) == 0 {
			continue // ANY_MAP DSes don't require matchLists
		}
		if _, ok := dsMatchMap[tc.DeliveryServiceName(ds.XMLID)]; !ok {
			t.Errorf("GET DeliveryService matches missing: %v\n", ds.XMLID)
		}
	}
}
