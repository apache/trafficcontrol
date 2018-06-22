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
	"strconv"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
)

func TestPutDeliveryServiceSafe(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, DeliveryServices}, func() {
		UpdateTestDeliveryServiceSafe(t)
	})
}

func UpdateTestDeliveryServiceSafe(t *testing.T) {
	firstDS := testData.DeliveryServices[0]

	dses, _, err := TOSession.GetDeliveryServices()
	if err != nil {
		t.Fatalf("cannot GET Delivery Services: %v\n", err)
	}

	remoteDS := tc.DeliveryService{}
	found := false
	for _, ds := range dses {
		if ds.XMLID == firstDS.XMLID {
			found = true
			remoteDS = ds
			break
		}
	}
	if !found {
		t.Fatalf("GET Delivery Services missing: %v\n", firstDS.XMLID)
	}

	req := tc.DeliveryServiceSafeUpdate{
		DisplayName: util.StrPtr("safe update display name"),
		InfoURL:     util.StrPtr("safe update info URL"),
		LongDesc:    util.StrPtr("safe update long desc"),
		LongDesc1:   util.StrPtr("safe update long desc one"),
	}

	if _, err := TOSession.UpdateDeliveryServiceSafe(remoteDS.ID, &req); err != nil {
		t.Fatalf("cannot UPDATE DeliveryService safe: %v\n", err)
	}

	updatedDS, _, err := TOSession.GetDeliveryService(strconv.Itoa(remoteDS.ID))
	if err != nil {
		t.Fatalf("cannot GET Delivery Service by ID: id %v xmlid '%v' - %v\n", remoteDS.ID, remoteDS.XMLID, err)
	}
	if updatedDS == nil {
		t.Fatalf("GET Delivery Service by ID returned nil err, but nil DS: %v - nil\n", remoteDS.XMLID)
	}

	if *req.DisplayName != updatedDS.DisplayName || *req.InfoURL != updatedDS.InfoURL || *req.LongDesc != updatedDS.LongDesc || *req.LongDesc1 != updatedDS.LongDesc1 {
		t.Fatalf("ds safe update succeeded, but get delivery service didn't match. expected: {'%v' '%v' '%v' '%v'} actual: {'%v' '%v' '%v' '%v'}\n", *req.DisplayName, *req.InfoURL, *req.LongDesc, *req.LongDesc1, updatedDS.DisplayName, updatedDS.InfoURL, updatedDS.LongDesc, updatedDS.LongDesc1)
	}
}
