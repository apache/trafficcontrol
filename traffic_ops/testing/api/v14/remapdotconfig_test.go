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
		t.Fatalf("GET delivery service servers: %v\n", err)
	} else if len(dsServers.Response) == 0 {
		t.Fatalf("GET delivery service servers: no servers found\n")
	} else if dsServers.Response[0].Server == nil {
		t.Fatalf("GET delivery service servers: returned nil server\n")
	} else if dsServers.Response[0].DeliveryService == nil {
		t.Fatalf("GET delivery service servers: returned nil ds\n")
	}
	serverID := *dsServers.Response[0].Server

	ds := (*tc.DeliveryService)(nil)
	for _, dsServer := range dsServers.Response {
		ds, _, err = TOSession.GetDeliveryService(strconv.Itoa(*dsServer.DeliveryService))
		if err != nil {
			t.Fatalf("Getting ds %+v: "+err.Error()+"\n", *dsServers.Response[0].DeliveryService)
		} else if ds == nil {
			t.Fatalf("Getting ds %+v: "+"got nil response"+"\n", *dsServers.Response[0].DeliveryService)
		} else if ds.OrgServerFQDN == "" {
			t.Fatalf("Getting ds %+v: "+"got empty ds.OrgServerFQDN"+"\n", *dsServers.Response[0].DeliveryService)
		}
		if ds.Type == tc.DSTypeAnyMap {
			continue
		}
		break
	}
	if ds == nil || ds.XMLID == "" {
		t.Fatalf("no Delivery Service found with assigned servers that isn't an ANY_MAP service, can't test remap.config")
	}

	originURI, err := url.Parse(ds.OrgServerFQDN)
	if err != nil {
		t.Fatalf("Getting ds %+v: "+" ds.OrgServerFQDN '%+v' failed to parse as a URL: %+v\n", *dsServers.Response[0].DeliveryService, ds.OrgServerFQDN, err)
	}
	originHost := originURI.Hostname()

	remapDotConfig, _, err := TOSession.GetATSServerConfig(serverID, "remap.config")
	if err != nil {
		t.Fatalf("Getting server %+v config remap.config: "+err.Error()+"\n", serverID)
	}

	if !strings.Contains(remapDotConfig, originHost) {
		t.Errorf("expected: remap.config to contain delivery service origin FQDN '%+v' host '%+v', actual: '''%+v'''", ds.OrgServerFQDN, originHost, remapDotConfig)
	}

	remapDotConfigLines := strings.Split(remapDotConfig, "\n")
	for i, line := range remapDotConfigLines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if line[0] == '#' {
			continue
		}
		if !strings.HasPrefix(line, "map") {
			t.Errorf("expected: remap.config line %v to start with 'map', actual: '%v'\n", i, line)
		}
	}
}
