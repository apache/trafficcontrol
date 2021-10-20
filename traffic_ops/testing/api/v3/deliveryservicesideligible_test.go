package v3

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
	"net/http"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v6/lib/go-rfc"
)

func TestDeliveryServicesEligible(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, DeliveryServices}, func() {
		GetTestDeliveryServicesEligibleIMS(t)
		GetTestDeliveryServicesEligible(t)
	})
}

func GetTestDeliveryServicesEligibleIMS(t *testing.T) {
	var header http.Header
	header = make(map[string][]string)
	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, time)
	_, reqInf, err := TOSession.GetDeliveryServicesNullableWithHdr(header)
	if err != nil {
		t.Fatalf("could not GET eligible delivery services: %v", err)
	}
	if reqInf.StatusCode != http.StatusNotModified {
		t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
	}
}

func GetTestDeliveryServicesEligible(t *testing.T) {
	dses, _, err := TOSession.GetDeliveryServicesNullable()
	if err != nil {
		t.Errorf("cannot GET DeliveryServices: %v", err)
	}
	if len(dses) == 0 {
		t.Error("GET DeliveryServices returned no delivery services, need at least 1 to test")
	}
	dsID := dses[0].ID
	servers, _, err := TOSession.GetDeliveryServicesEligible(*dsID)
	if err != nil {
		t.Errorf("getting delivery services eligible: %v", err)
	}
	if len(servers) == 0 {
		t.Error("getting delivery services eligible returned no servers")
	}
}
