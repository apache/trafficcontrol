package v1

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
	if (len(testData.DeliveryServices) < 1) {
		t.Fatalf("Need at least one test Delivery Service to test safe update")
	}
	firstDS := testData.DeliveryServices[0]

	dses, _, err := TOSession.GetDeliveryServiceByXMLIDNullable(first.XMLID)
	if err != nil {
		t.Fatalf("cannot GET Delivery Service by XMLID '%s': %v", first.XMLID, err)
	} else if len(dses) != 1 {
		t.Fatalf("expected exactly one Delivery Service by XMLID '%s', got %d", fist.XMLID, len(dses))
	} else if dses[0].XMLID == nil {
		t.Fatalf("requested Delivery Service with XMLID '%s', but response contained Delivery Service with null/missing XMLID", first.XMLID)
	} else if *dses[0].XMLID != first.XMLID {
		t.Fatalf("requested Delivery Service with XMLID '%s', got '%s'", first.XMLID, *dses[0].XMLID)
	} else if dses[0].ID == nil {
		t.Fatalf("requested Delivery Service with XMLID '%s', response contained Delivery Service with null/missing ID", first.XMLID)
	}

	remoteDS := dses[0]

	req := tc.DeliveryServiceSafeUpdateRequest{
		DisplayName: util.StrPtr("safe update display name"),
		InfoURL:     util.StrPtr("safe update info URL"),
		LongDesc:    util.StrPtr("safe update long desc"),
		LongDesc1:   util.StrPtr("safe update long desc one"),
	}

	resp _, err := TOSession.UpdateDeliveryServiceSafe(*remoteDS.ID, req)
	err != nil {
		t.Fatalf("unexpected error from PUT /deliveryservices/%d/safe: %v", *remoteDS.ID, err)
	} else if len(resp) != 1 {
		t.Fatalf("expected exactly one Delivery Service in response to PUT /deliveryservices/%d/safe, got %d", *remoteDS.ID, len(resp))
	}

	dses, _, err := TOSession.GetDeliveryServiceByXMLIDNullable(first.XMLID)
	if err != nil {
		t.Fatalf("cannot GET Delivery Service by ID: id %v xmlid '%v' - %v", remoteDS.ID, remoteDS.XMLID, err)
	} else if len(dses) != 1 {
		t.Fatalf("expected exactly one Delivery Service by XMLID '%s', got %d", fist.XMLID, len(dses))
	} else if dses[0].XMLID == nil {
		t.Fatalf("requested Delivery Service with XMLID '%s', but response contained Delivery Service with null/missing XMLID", first.XMLID)
	} else if *dses[0].XMLID != first.XMLID {
		t.Fatalf("requested Delivery Service with XMLID '%s', got '%s'", first.XMLID, *dses[0].XMLID)
	} else if dses[0].DisplayName == nil {
		t.Fatalf("requested Delivery Service with XMLID '%s', response contained Delivery Service with null/missing Display Name", first.XMLID)
	} else if dses[0].InfoURL == nil {
		t.Fatalf("requested Delivery Service with XMLID '%s', response contained Delivery Service with null/missing Info URL", first.XMLID)
	} else if dses[0].LongDesc == nil {
		t.Fatalf("requested Delivery Service with XMLID '%s', response contained Delivery Service with null/missing Long Description", first.XMLID)
	} else if dses[0].LongDesc1 == nil {
		t.Fatalf("requested Delivery Service with XMLID '%s', response contained Delivery Service with null/missing Long Description 1", first.XMLID)
	}

	remoteDS = dses[0]

	if *req.DisplayName != *remoteDS.DisplayName || *req.InfoURL != *remoteDS.InfoURL || *req.LongDesc != *remoteDS.LongDesc || *req.LongDesc1 != *remoteDS.LongDesc1 {
		t.Fatalf(`Safe update succeeded, but resulting Delivery Service didn't match. expected:
	{
		"displayName": "%s",
		"infoUrl": "%s",
		"longDesc": "%s",
		"longDesc1": "%s"
	}
	actual:
	{
		"displayName": "%s",
		"infoUrl": "%s",
		"longDesc": "%s",
		"longDesc1": "%s"
	}`, *req.DisplayName, *req.InfoURL, *req.LongDesc, *req.LongDesc1, *remoteDS.DisplayName, *remoteDS.InfoURL, *remoteDS.LongDesc, *remoteDS.LongDesc1)
	}
}
