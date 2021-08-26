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
	"net/http"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	client "github.com/apache/trafficcontrol/traffic_ops/v5-client"
)

func TestDeliveryServicesEligible(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, DeliveryServices}, func() {
		GetTestDeliveryServicesEligibleIMS(t)
		GetTestDeliveryServicesEligible(t)
	})
}

func GetTestDeliveryServicesEligibleIMS(t *testing.T) {
	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)

	opts := client.NewRequestOptions()
	opts.Header.Set(rfc.IfModifiedSince, time)
	resp, reqInf, err := TOSession.GetDeliveryServices(opts)
	if err != nil {
		t.Fatalf("could not get eligible delivery services: %v - alerts: %+v", err, resp.Alerts)
	}
	if reqInf.StatusCode != http.StatusNotModified {
		t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
	}
}

func GetTestDeliveryServicesEligible(t *testing.T) {
	dses, _, err := TOSession.GetDeliveryServices(client.RequestOptions{})
	if err != nil {
		t.Errorf("cannot get Delivery Services: %v - alerts: %+v", err, dses.Alerts)
	}
	if len(dses.Response) == 0 {
		t.Fatal("GET DeliveryServices returned no delivery services, need at least 1 to test")
	}
	dsID := dses.Response[0].ID
	if dsID == nil {
		t.Fatal("Traffic Ops returned a representation of a Delivery Service that had null or undefined ID")
	}
	servers, _, err := TOSession.GetDeliveryServicesEligible(*dsID, client.RequestOptions{})
	if err != nil {
		t.Errorf("getting Delivery Services eligible: %v - alerts: %+v", err, servers.Alerts)
	}
	if len(servers.Response) == 0 {
		t.Error("getting delivery services eligible returned no servers")
	}
}
