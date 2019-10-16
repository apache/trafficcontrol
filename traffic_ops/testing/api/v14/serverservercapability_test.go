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

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
)

func TestServerServerCapabilities(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, ServerCapabilities, ServerServerCapabilities}, func() {
		GetTestServerServerCapabilities(t)
	})
}

func CreateTestServerServerCapabilities(t *testing.T) {

	// Valid POSTs

	// loop through server ServerCapabilities, assign FKs and create
	for _, ssc := range testData.ServerServerCapabilities {
		servResp, _, err := TOSession.GetServerByHostName(*ssc.Server)
		if err != nil {
			t.Fatalf("cannot GET Server by hostname: %v - %v\n", *ssc.Server, err)
		}
		if len(servResp) != 1 {
			t.Fatalf("cannot GET Server by hostname: %v. Response did not include record.\n", *ssc.Server)
		}
		server := servResp[0]
		ssc.ServerID = &server.ID
		resp, _, err := TOSession.CreateServerServerCapability(ssc)
		if err != nil {
			t.Errorf("could not POST the server capability %v to server %v: %v\n", *ssc.ServerCapability, *ssc.Server, err)
		}
		log.Debugln("Response: ", server.HostName, " ", resp)
	}

	// Invalid POSTs

	ssc := testData.ServerServerCapabilities[0]

	// Attempt to assign already assigned server capability
	_, _, err := TOSession.CreateServerServerCapability(ssc)
	if err == nil {
		t.Error("expected to receive error when assigning a already assigned server capability\n")
	}

	// Attempt to assign an server capability with no ID
	sscNilID := tc.ServerServerCapability{
		ServerCapability: ssc.ServerCapability,
	}
	_, _, err = TOSession.CreateServerServerCapability(sscNilID)
	if err == nil {
		t.Error("expected to receive error when assigning a server capability without a server ID\n")
	}

	// Attempt to assign an server capability with no server capability
	sscNilCapability := tc.ServerServerCapability{
		ServerID: ssc.ServerID,
	}
	_, _, err = TOSession.CreateServerServerCapability(sscNilCapability)
	if err == nil {
		t.Error("expected to receive error when assigning a server capability to a server without a server capability\n")
	}

	// Attempt to assign an server capability with invalid server capability
	sscInvalidCapability := tc.ServerServerCapability{
		ServerID:         ssc.ServerID,
		ServerCapability: util.StrPtr("bogus"),
	}
	_, _, err = TOSession.CreateServerServerCapability(sscInvalidCapability)
	if err == nil {
		t.Error("expected to receive error when assigning a non existent server capability to a server\n")
	}

	// Attempt to assign an server capability with invalid server capability
	sscInvalidID := tc.ServerServerCapability{
		ServerID:         util.IntPtr(-1),
		ServerCapability: ssc.ServerCapability,
	}
	_, _, err = TOSession.CreateServerServerCapability(sscInvalidID)
	if err == nil {
		t.Error("expected to receive error when assigning a server capability to a non existent server ID\n")
	}
}

func GetTestServerServerCapabilities(t *testing.T) {
	// Get All Server Capabilities
	sscs, _, err := TOSession.GetServerServerCapabilities(nil, nil, nil)
	if err != nil {
		t.Fatalf("cannot GET server capabilities assigned to servers: %v\n", err)
	}
	if sscs == nil {
		t.Fatal("returned server capabilities assigned to servers was nil\n")
	}
	if len(sscs) != len(testData.ServerServerCapabilities) {
		t.Errorf("expect %v server capabilities assigned to servers received %v \n", len(testData.ServerServerCapabilities), len(sscs))
	}

	checkResp := func(t *testing.T, sscs []tc.ServerServerCapability) {
		if sscs == nil {
			t.Fatal("returned server capabilities assigned to servers was nil\n")
		}
		if len(sscs) != 1 {
			t.Errorf("expect 1 server capabilities assigned to server received %v \n", len(sscs))
		}
	}

	for _, ssc := range sscs {
		// Get assigned Server Capabilities by server id
		sscs, _, err := TOSession.GetServerServerCapabilities(ssc.ServerID, nil, nil)
		if err != nil {
			t.Fatalf("cannot GET server capabilities assigned to servers by server ID %v: %v\n", *ssc.ServerID, err)
		}
		checkResp(t, sscs)
		// Get assigned Server Capabilities by host name
		sscs, _, err = TOSession.GetServerServerCapabilities(nil, ssc.Server, nil)
		if err != nil {
			t.Fatalf("cannot GET server capabilities assigned to servers by server host name %v: %v\n", *ssc.Server, err)
		}
		checkResp(t, sscs)

		// Get assigned Server Capabilities by server capability
		sscs, _, err = TOSession.GetServerServerCapabilities(nil, nil, ssc.ServerCapability)
		if err != nil {
			t.Fatalf("cannot GET server capabilities assigned to servers by server capability %v: %v\n", *ssc.ServerCapability, err)
		}
		checkResp(t, sscs)
	}
}

func DeleteTestServerServerCapabilities(t *testing.T) {
	// Get Server Capabilities to delete them
	sscs, _, err := TOSession.GetServerServerCapabilities(nil, nil, nil)
	if err != nil {
		t.Fatalf("cannot GET server capabilities assigned to servers: %v\n", err)
	}
	if sscs == nil {
		t.Fatal("returned server capabilities assigned to servers was nil\n")
	}

	// Delete them
	for _, ssc := range sscs {
		_, _, err := TOSession.DeleteServerServerCapability(*ssc.ServerID, *ssc.ServerCapability)
		if err != nil {
			t.Errorf("could not DELETE the server capability %v from server %v: %v\n", *ssc.ServerCapability, *ssc.Server, err)
		}
	}
}
