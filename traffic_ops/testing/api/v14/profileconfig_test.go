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
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

func TestProfileDotConfig(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, DeliveryServices}, func() {
		defer DeleteTestDeliveryServiceServersCreated(t)
		CreateTestDeliveryServiceServers(t)
		GetTestProfileDotConfig(t)
	})
}

func GetTestProfileDotConfig(t *testing.T) {
	dsServers, _, err := TOSession.GetDeliveryServiceServers()
	if err != nil {
		t.Fatalf("GET delivery service servers: %v\n", err)
	} else if len(dsServers.Response) == 0 {
		t.Fatalf("GET delivery service servers: no servers found\n")
	} else if dsServers.Response[0].Server == nil {
		t.Fatalf("GET delivery service servers: returned nil server\n")
	} else if dsServers.Response[0].DeliveryService == nil {
		t.Fatalf("GET delivery service servers: returned nil ds\n")
	}
	serverID := *dsServers.Response[0].Server

	ds, _, err := TOSession.GetDeliveryService(strconv.Itoa(*dsServers.Response[0].DeliveryService))
	if err != nil {
		t.Fatalf("Getting ds %+v: "+err.Error()+"\n", *dsServers.Response[0].DeliveryService)
	} else if ds == nil {
		t.Fatalf("Getting ds %+v: "+"got nil response"+"\n", *dsServers.Response[0].DeliveryService)
	} else if ds.OrgServerFQDN == "" {
		t.Fatalf("Getting ds %+v: "+"got empty ds.OrgServerFQDN"+"\n", *dsServers.Response[0].DeliveryService)
	}

	if _, err := url.Parse(ds.OrgServerFQDN); err != nil {
		t.Fatalf("Getting ds %+v: "+" ds.OrgServerFQDN '%+v' failed to parse as a URL: %+v\n", *dsServers.Response[0].DeliveryService, ds.OrgServerFQDN, err)
	}

	servers, _, err := TOSession.GetServers()
	if err != nil {
		t.Errorf("cannot GET Servers: %v\n", err)
	}

	server := tc.Server{ID: -1}
	for _, potentialServer := range servers {
		if potentialServer.Type != string(tc.CacheTypeEdge) {
			continue
		}
		server = potentialServer
	}
	if server.ID == -1 {
		t.Errorf("GET Servers returned no edge servers, must have at least 1 to test")
	}

	profileDotConfig, _, err := TOSession.GetATSProfileConfig(server.ProfileID, "cache.config")
	if err != nil {
		t.Fatalf("Getting server %+v config parent.config: "+err.Error()+"\n", serverID)
	}

	if !strings.Contains(profileDotConfig, server.Profile) {
		t.Errorf("expected: profile cache.config to contain profile name '%+v', actual: '''%+v'''", server.Profile, profileDotConfig)
	}
}
