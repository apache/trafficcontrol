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
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

func TestATSConfigMeta(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, DeliveryServices}, func() {
		GetTestATSConfigMeta(t)
	})
}

func GetTestATSConfigMeta(t *testing.T) {
	if len(testData.Servers) < 1 {
		t.Fatalf("cannot GET Server: no test data\n")
	}
	testServer := testData.Servers[0]

	serverList, _, err := TOSession.GetServerByHostName(testServer.HostName)
	if err != nil {
		t.Fatalf("cannot GET Server: %v\n", err)
	}
	if len(serverList) < 1 {
		t.Fatalf("cannot GET Server '" + testServer.HostName + "', returned no servers\n")
	}
	server := serverList[0]

	lst, _, err := TOSession.GetATSServerConfigList(server.ID)
	if err != nil {
		t.Fatalf("Getting server '" + server.HostName + "' config list: " + err.Error() + "\n")
	}

	expected := tc.ATSConfigMetaDataConfigFile{
		FileNameOnDisk: "hdr_rw_mid_anymap-ds.config",
		Location:       "/remap/config/location/parameter",
		APIURI:         "cdns/cdn1/configfiles/ats/hdr_rw_mid_anymap-ds.config", // expected suffix; config gen doesn't care about API version
		URL:            "",
		Scope:          "cdns",
	}

	actual := (*tc.ATSConfigMetaDataConfigFile)(nil)
	for _, cfg := range lst.ConfigFiles {
		if cfg.FileNameOnDisk == expected.FileNameOnDisk {
			actual = &cfg
			break
		}
	}
	if actual == nil {
		t.Fatalf("Getting server '"+server.HostName+"' config list: expected: %+v actual: not found\n", expected.FileNameOnDisk)
	}

	if expected.FileNameOnDisk != actual.FileNameOnDisk {
		t.Errorf("Getting server '"+server.HostName+"' config list: expected: %+v actual: %+v\n", expected, *actual)
	}
	if expected.Location != actual.Location {
		t.Errorf("Getting server '"+server.HostName+"' config list: expected: %+v actual: %+v\n", expected, *actual)
	}
	if !strings.HasSuffix(actual.APIURI, expected.APIURI) {
		t.Errorf("Getting server '"+server.HostName+"' config list: expected: %+v actual: %+v\n", expected, *actual)
	}
	if actual.Scope != expected.Scope {
		t.Errorf("Getting server '"+server.HostName+"' config list: expected: %+v actual: %+v\n", expected, *actual)
	}
}
