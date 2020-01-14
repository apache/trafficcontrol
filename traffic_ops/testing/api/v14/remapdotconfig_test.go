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
	"github.com/apache/trafficcontrol/lib/go-tc/tce"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

func TestRemapDotConfig(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, DeliveryServices}, func() {
		defer DeleteTestDeliveryServiceServersCreated(t)
		CreateTestDeliveryServiceServers(t)
		GetTestRemapDotConfig(t)
	})
}

func GetTestRemapDotConfig(t *testing.T) {
	dsServers, _, err := TOSession.GetDeliveryServiceServers()
	if err != nil {
		t.Fatalf("GET delivery service servers: %v", err)
	}
	if len(dsServers.Response) == 0 {
		t.Fatal("GET delivery service servers: no servers found")
	}

	var ds *tc.DeliveryService
	var serverID int
	for _, dsServer := range dsServers.Response {
		if dsServer.Server == nil {
			t.Error("Found DS-Server assignment with nil server")
			continue
		}
		if dsServer.DeliveryService == nil {
			t.Error("Found DS-Server assignment with nil Delivery Service")
			continue
		}

		serverID = *dsServer.Server

		ds, _, err = TOSession.GetDeliveryService(strconv.Itoa(*dsServer.DeliveryService))
		if err != nil {
			t.Errorf("Getting ds %+v: %v", *dsServer.DeliveryService, err)
			continue
		}
		if ds == nil {
			t.Errorf("Getting ds %+v: got nil response", *dsServer.DeliveryService)
			continue
		}
		if ds.OrgServerFQDN == "" {
			t.Errorf("Getting ds %+v: got empty ds.OrgServerFQDN", *dsServer.DeliveryService)
			continue
		}

		if ds.Type != tce.DSTypeAnyMap {
			break
		}
	}

	if ds == nil || ds.XMLID == "" {
		t.Fatal("no Delivery Service found with assigned servers that isn't an ANY_MAP service, can't test remap.config")
	}

	originURI, err := url.Parse(ds.OrgServerFQDN)
	if err != nil {
		t.Fatalf("Getting ds %+v: ds.OrgServerFQDN '%+v' failed to parse as a URL: %+v", ds.XMLID, ds.OrgServerFQDN, err)
	}
	originHost := originURI.Hostname()

	remapDotConfig, _, err := TOSession.GetATSServerConfig(serverID, "remap.config")
	if err != nil {
		t.Fatalf("Getting server %+v config remap.config: %v", serverID, err)
	}

	if !strings.Contains(remapDotConfig, originHost) {
		t.Errorf("expected: remap.config to contain delivery service origin FQDN '%v' host '%v', actual:\n'''\n%v\n'''", ds.OrgServerFQDN, originHost, remapDotConfig)
	}

	remapDotConfigLines := strings.Split(remapDotConfig, "\n")
	for i, line := range remapDotConfigLines {
		line = strings.TrimSpace(line)
		if len(line) > 0 && line[0] != '#' && !strings.HasPrefix(line, "map") {
			t.Errorf("expected: remap.config line %v to start with 'map', actual: '%v'", i, line)
		}
	}
}
