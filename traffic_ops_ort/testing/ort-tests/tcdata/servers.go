package tcdata

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
	"testing"

	"github.com/apache/trafficcontrol/lib/go-util"
)

func (r *TCData) CreateTestServers(t *testing.T) {
	// loop through servers, assign FKs and create
	for _, server := range r.TestData.Servers {
		if server.HostName == nil {
			t.Errorf("found server with nil hostname: %+v", server)
			continue
		}
		resp, _, err := TOSession.CreateServerWithHdr(server, nil)
		t.Log("Response: ", *server.HostName, " ", resp)
		if err != nil {
			t.Errorf("could not CREATE servers: %v", err)
		}
	}
}

func (r *TCData) CreateTestBlankFields(t *testing.T) {
	serverResp, _, err := TOSession.GetServersWithHdr(nil, nil)
	if err != nil {
		t.Fatalf("couldnt get servers: %v", err)
	}
	if len(serverResp.Response) < 1 {
		t.Fatal("expected at least one server")
	}
	server := serverResp.Response[0]
	originalHost := server.HostName

	server.HostName = util.StrPtr("")
	_, _, err = TOSession.UpdateServerByIDWithHdr(*server.ID, server, nil)
	if err == nil {
		t.Fatal("should not be able to update server with blank HostName")
	}

	server.HostName = originalHost
	server.DomainName = util.StrPtr("")
	_, _, err = TOSession.UpdateServerByIDWithHdr(*server.ID, server, nil)
	if err == nil {
		t.Fatal("should not be able to update server with blank DomainName")
	}
}

func (r *TCData) CreateTestServerWithoutProfileId(t *testing.T) {
	params := url.Values{}
	servers := r.TestData.Servers[19]
	params.Set("hostName", *servers.HostName)

	resp, _, err := TOSession.GetServersWithHdr(&params, nil)
	if err != nil {
		t.Fatalf("cannot GET Server by name '%s': %v - %v", *servers.HostName, err, resp.Alerts)
	}

	server := resp.Response[0]
	originalProfile := *server.Profile
	delResp, _, err := TOSession.DeleteServerByID(*server.ID)
	if err != nil {
		t.Fatalf("cannot DELETE Server by ID %d: %v - %v", *server.ID, err, delResp)
	}

	*server.Profile = ""
	server.ProfileID = nil
	response, reqInfo, errs := TOSession.CreateServerWithHdr(server, nil)
	t.Log("Response: ", *server.HostName, " ", response)
	if reqInfo.StatusCode != 400 {
		t.Fatalf("Expected status code: %v but got: %v", "400", reqInfo.StatusCode)
	}

	//Reverting it back for further tests
	*server.Profile = originalProfile
	response, _, errs = TOSession.CreateServerWithHdr(server, nil)
	t.Log("Response: ", *server.HostName, " ", response)
	if errs != nil {
		t.Fatalf("could not CREATE servers: %v", errs)
	}
}

func (r *TCData) DeleteTestServers(t *testing.T) {
	params := url.Values{}

	for _, server := range r.TestData.Servers {
		if server.HostName == nil {
			t.Errorf("found server with nil hostname: %+v", server)
			continue
		}

		params.Set("hostName", *server.HostName)

		resp, _, err := TOSession.GetServersWithHdr(&params, nil)
		if err != nil {
			t.Errorf("cannot GET Server by hostname '%s': %v - %v", *server.HostName, err, resp.Alerts)
			continue
		}
		if len(resp.Response) > 0 {
			if len(resp.Response) > 1 {
				t.Errorf("Expected exactly one server by hostname '%s' - actual: %d", *server.HostName, len(resp.Response))
				t.Logf("Testing will proceed with server: %+v", resp.Response[0])
			}
			respServer := resp.Response[0]

			if respServer.ID == nil {
				t.Errorf("Server '%s' had nil ID", *server.HostName)
				continue
			}

			delResp, _, err := TOSession.DeleteServerByID(*respServer.ID)
			if err != nil {
				t.Errorf("cannot DELETE Server by ID %d: %v - %v", *respServer.ID, err, delResp)
				continue
			}

			// Retrieve the Server to see if it got deleted
			resp, _, err := TOSession.GetServersWithHdr(&params, nil)
			if err != nil {
				t.Errorf("error deleting Server hostname '%s': %v - %v", *server.HostName, err, resp.Alerts)
			}
			if len(resp.Response) > 0 {
				t.Errorf("expected Server hostname: %s to be deleted", *server.HostName)
			}
		}
	}
}
