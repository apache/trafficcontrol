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
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
)

func TestServers(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Topologies, DeliveryServices, Servers}, func() {
		GetTestServersIMS(t)
		currentTime := time.Now().UTC().Add(-5 * time.Second)
		timestamp := currentTime.Format(time.RFC1123)
		var header http.Header
		header = make(map[string][]string)
		header.Set(rfc.IfModifiedSince, timestamp)
		UpdateTestServers(t)
		GetTestServersDetails(t)
		GetTestServers(t)
		GetTestServersIMSAfterChange(t, header)
		GetTestServersQueryParameters(t)
		CreateTestServerWithoutProfileId(t)
		UniqueIPProfileTestServers(t)
		UpdateTestServerStatus(t)
	})
}

func UpdateTestServerStatus(t *testing.T) {
	if len(testData.Servers) < 1 {
		t.Fatal("Need at least one server to test updating")
	}

	firstServer := testData.Servers[0]
	if firstServer.HostName == nil {
		t.Fatalf("First test server had nil hostname: %+v", firstServer)
	}

	hostName := *firstServer.HostName
	params := url.Values{}
	params.Add("hostName", hostName)

	// Retrieve the server by hostname so we can get the id for the Update
	resp, _, err := TOSession.GetServers(&params)
	if err != nil {
		t.Fatalf("cannot GET Server by hostname '%s': %v - %v", hostName, err, resp.Alerts)
	}
	if len(resp.Response) < 1 {
		t.Fatalf("Expected at least one server to exist by hostname '%s'", hostName)
	}
	if len(resp.Response) > 1 {
		t.Errorf("Expected exactly one server to exist by hostname '%s' - actual: %d", hostName, len(resp.Response))
		t.Logf("Testing will proceed with server: %+v", resp.Response[0])
	}
	remoteServer := resp.Response[0]
	if remoteServer.ID == nil {
		t.Fatalf("Got null ID for server '%s'", hostName)
	}
	id := fmt.Sprintf("%v", *resp.Response[0].ID)
	idParam := url.Values{}
	idParam.Add("id", id)
	originalStatusID := 0
	updatedStatusID := 0

	statuses, _, err := TOSession.GetStatusesWithHdr(nil)
	if err != nil {
		t.Fatalf("cannot get statuses: %v", err.Error())
	}
	for _, status := range statuses {
		if status.Name == "REPORTED" {
			originalStatusID = status.ID
		}
		if status.Name == "ONLINE" {
			updatedStatusID = status.ID
		}
	}
	// Keeping the status same, perform an update and make sure that statusLastUpdated didnt change
	remoteServer.StatusID = &originalStatusID

	alerts, _, err := TOSession.UpdateServerByID(*remoteServer.ID, remoteServer)
	if err != nil {
		t.Fatalf("cannot UPDATE Server by ID %d (hostname '%s'): %v - %v", *remoteServer.ID, hostName, err, alerts)
	}

	resp, _, err = TOSession.GetServers(&idParam)
	if err != nil {
		t.Errorf("cannot GET Server by ID: %v - %v", *remoteServer.HostName, err)
	}
	if len(resp.Response) < 1 {
		t.Fatalf("Expected at least one server to exist by hostname '%s'", hostName)
	}
	if len(resp.Response) > 1 {
		t.Errorf("Expected exactly one server to exist by hostname '%s' - actual: %d", hostName, len(resp.Response))
		t.Logf("Testing will proceed with server: %+v", resp.Response[0])
	}

	respServer := resp.Response[0]

	if *remoteServer.StatusLastUpdated != *respServer.StatusLastUpdated {
		t.Errorf("since status didnt change, no change in 'StatusLastUpdated' time was expected. Difference observer: old value: %v, new value: %v",
			remoteServer.StatusLastUpdated.String(), respServer.StatusLastUpdated.String())
	}

	// Changing the status, perform an update and make sure that statusLastUpdated changed
	remoteServer.StatusID = &updatedStatusID

	alerts, _, err = TOSession.UpdateServerByID(*remoteServer.ID, remoteServer)
	if err != nil {
		t.Fatalf("cannot UPDATE Server by ID %d (hostname '%s'): %v - %v", *remoteServer.ID, hostName, err, alerts)
	}

	resp, _, err = TOSession.GetServers(&idParam)
	if err != nil {
		t.Errorf("cannot GET Server by ID: %v - %v", *remoteServer.HostName, err)
	}
	if len(resp.Response) < 1 {
		t.Fatalf("Expected at least one server to exist by hostname '%s'", hostName)
	}
	if len(resp.Response) > 1 {
		t.Errorf("Expected exactly one server to exist by hostname '%s' - actual: %d", hostName, len(resp.Response))
		t.Logf("Testing will proceed with server: %+v", resp.Response[0])
	}

	respServer = resp.Response[0]

	if *remoteServer.StatusLastUpdated == *respServer.StatusLastUpdated {
		t.Errorf("since status was changed, expected to see a time difference between the old and new 'StatusLastUpdated' values, got the same value")
	}

	// Changing the status, perform an update and make sure that statusLastUpdated changed
	remoteServer.StatusID = &originalStatusID

	alerts, _, err = TOSession.UpdateServerByID(*remoteServer.ID, remoteServer)
	if err != nil {
		t.Fatalf("cannot UPDATE Server by ID %d (hostname '%s'): %v - %v", *remoteServer.ID, hostName, err, alerts)
	}

	resp, _, err = TOSession.GetServers(&idParam)
	if err != nil {
		t.Errorf("cannot GET Server by ID: %v - %v", *remoteServer.HostName, err)
	}
	if len(resp.Response) < 1 {
		t.Fatalf("Expected at least one server to exist by hostname '%s'", hostName)
	}
	if len(resp.Response) > 1 {
		t.Errorf("Expected exactly one server to exist by hostname '%s' - actual: %d", hostName, len(resp.Response))
		t.Logf("Testing will proceed with server: %+v", resp.Response[0])
	}

	respServer = resp.Response[0]

	if *remoteServer.StatusLastUpdated == *respServer.StatusLastUpdated {
		t.Errorf("since status was changed, expected to see a time difference between the old and new 'StatusLastUpdated' values, got the same value")
	}
}

func GetTestServersIMSAfterChange(t *testing.T, header http.Header) {
	params := url.Values{}
	for _, server := range testData.Servers {
		if server.HostName == nil {
			t.Errorf("found server with nil hostname: %+v", server)
			continue
		}
		params.Set("hostName", *server.HostName)
		_, reqInf, err := TOSession.GetServersWithHdr(&params, header)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusOK {
			t.Fatalf("Expected 200 status code, got %v", reqInf.StatusCode)
		}
	}
	currentTime := time.Now().UTC()
	currentTime = currentTime.Add(1 * time.Second)
	timeStr := currentTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, timeStr)
	for _, server := range testData.Servers {
		if server.HostName == nil {
			t.Errorf("found server with nil hostname: %+v", server)
			continue
		}
		params.Set("hostName", *server.HostName)
		_, reqInf, err := TOSession.GetServersWithHdr(&params, header)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func GetTestServersIMS(t *testing.T) {
	var header http.Header
	header = make(map[string][]string)
	futureTime := time.Now().AddDate(0, 0, 1)
	timestamp := futureTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, timestamp)
	params := url.Values{}
	for _, server := range testData.Servers {
		if server.HostName == nil {
			t.Errorf("found server with nil hostname: %+v", server)
			continue
		}
		params.Set("hostName", *server.HostName)
		_, reqInf, err := TOSession.GetServersWithHdr(&params, header)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func CreateTestServers(t *testing.T) {
	// loop through servers, assign FKs and create
	for _, server := range testData.Servers {
		if server.HostName == nil {
			t.Errorf("found server with nil hostname: %+v", server)
			continue
		}
		resp, _, err := TOSession.CreateServer(server)
		t.Log("Response: ", *server.HostName, " ", resp)
		if err != nil {
			t.Errorf("could not CREATE servers: %v", err)
		}
	}
}

func CreateTestServerWithoutProfileId(t *testing.T) {
	params := url.Values{}
	servers := testData.Servers[19]
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
	response, reqInfo, errs := TOSession.CreateServer(server)
	t.Log("Response: ", *server.HostName, " ", response)
	if reqInfo.StatusCode != 400 {
		t.Fatalf("Expected status code: %v but got: %v", "400", reqInfo.StatusCode)
	}

	//Reverting it back for further tests
	*server.Profile = originalProfile
	response, _, errs = TOSession.CreateServer(server)
	t.Log("Response: ", *server.HostName, " ", response)
	if errs != nil {
		t.Fatalf("could not CREATE servers: %v", errs)
	}
}

func GetTestServers(t *testing.T) {
	params := url.Values{}
	for _, server := range testData.Servers {
		if server.HostName == nil {
			t.Errorf("found server with nil hostname: %+v", server)
			continue
		}
		params.Set("hostName", *server.HostName)
		resp, _, err := TOSession.GetServersWithHdr(&params, nil)
		if err != nil {
			t.Errorf("cannot GET Server by name '%s': %v - %v", *server.HostName, err, resp.Alerts)
		} else if resp.Summary.Count != 1 {
			t.Errorf("incorrect server count, expected: 1, actual: %d", resp.Summary.Count)

		}
	}
}

func GetTestServersDetails(t *testing.T) {

	for _, server := range testData.Servers {
		if server.HostName == nil {
			t.Errorf("found server with nil hostname: %+v", server)
			continue
		}
		resp, _, err := TOSession.GetServerDetailsByHostNameWithHdr(*server.HostName, nil)
		if err != nil {
			t.Errorf("cannot GET Server Details by name: %v - %v", err, resp)
		}
	}
}

func GetTestServersQueryParameters(t *testing.T) {
	dses, _, err := TOSession.GetDeliveryServicesNullableWithHdr(nil)
	if err != nil {
		t.Fatalf("Failed to get Delivery Services: %v", err)
	}
	if len(dses) < 1 {
		t.Fatal("Failed to get at least one Delivery Service")
	}

	ds := dses[0]
	if ds.ID == nil {
		t.Fatal("Got Delivery Service with nil ID")
	}

	params := url.Values{}
	params.Add("dsId", strconv.Itoa(*ds.ID))
	_, _, err = TOSession.GetServersWithHdr(&params, nil)
	if err != nil {
		t.Fatalf("Failed to get server by Delivery Service ID: %v", err)
	}

	currentTime := time.Now().UTC().Add(5 * time.Second)
	timestamp := currentTime.Format(time.RFC1123)
	var header http.Header
	header = make(map[string][]string)
	header.Set(rfc.IfModifiedSince, timestamp)
	_, reqInf, _ := TOSession.GetServersWithHdr(&params, header)
	if reqInf.StatusCode != http.StatusNotModified {
		t.Errorf("Expected a status code of 304, got %v", reqInf.StatusCode)
	}

	foundTopDs := false
	const topDsXmlId = "ds-top"
	for _, ds = range dses {
		if ds.XMLID == nil || *ds.XMLID != topDsXmlId {
			continue
		}
		foundTopDs = true
		break
	}
	if !foundTopDs {
		t.Fatalf("unable to find deliveryservice %s", topDsXmlId)
	}
	params.Set("dsId", strconv.Itoa(*ds.ID))
	expectedHostnames := map[string]bool{
		"edge1-cdn1-cg3": true,
		"edge2-cdn1-cg3": true,
		"atlanta-mid-16": true,
	}
	response, _, err := TOSession.GetServersWithHdr(&params, nil)
	if err != nil {
		t.Fatalf("Failed to get servers by Topology-based Delivery Service ID with xmlId %s", err)
	}
	if len(response.Response) == 0 {
		t.Fatalf("Did not find any servers for Topology-based Delivery Service with xmlId %s", err)
	}
	for _, server := range response.Response {
		if _, exists := expectedHostnames[*server.HostName]; !exists {
			t.Fatalf("expected hostnames %v, actual %s actual: ", expectedHostnames, *server.HostName)
		}
	}

	params.Del("dsId")

	resp, _, err := TOSession.GetServersWithHdr(nil, nil)
	if err != nil {
		t.Fatalf("Failed to get servers: %v", err)
	}

	if len(resp.Response) < 1 {
		t.Fatalf("Failed to get at least one server")
	}

	s := resp.Response[0]

	params.Add("type", s.Type)
	if _, _, err := TOSession.GetServersWithHdr(&params, nil); err != nil {
		t.Errorf("Error getting servers by type: %v", err)
	}
	params.Del("type")

	if s.CachegroupID == nil {
		t.Error("Found server with no Cache Group ID")
	} else {
		params.Add("cachegroup", strconv.Itoa(*s.CachegroupID))
		if _, _, err := TOSession.GetServersWithHdr(&params, nil); err != nil {
			t.Errorf("Error getting servers by Cache Group ID: %v", err)
		}
		params.Del("cachegroup")
	}

	if s.Status == nil {
		t.Error("Found server with no status")
	} else {
		params.Add("status", *s.Status)
		if _, _, err := TOSession.GetServersWithHdr(&params, nil); err != nil {
			t.Errorf("Error getting servers by status: %v", err)
		}
		params.Del("status")
	}

	if s.ProfileID == nil {
		t.Error("Found server with no Profile ID")
	} else {
		params.Add("profileId", strconv.Itoa(*s.ProfileID))
		if _, _, err := TOSession.GetServersWithHdr(&params, nil); err != nil {
			t.Errorf("Error getting servers by Profile ID: %v", err)
		}
		params.Del("profileId")
	}

	cgs, _, err := TOSession.GetCacheGroupsNullableWithHdr(nil)
	if err != nil {
		t.Fatalf("Failed to get Cache Groups: %v", err)
	}
	if len(cgs) < 1 {
		t.Fatal("Failed to get at least one Cache Group")
	}
	if cgs[0].ID == nil {
		t.Fatal("Cache Group found with no ID")
	}

	params.Add("parentCacheGroup", strconv.Itoa(*cgs[0].ID))
	if _, _, err = TOSession.GetServersWithHdr(&params, nil); err != nil {
		t.Errorf("Error getting servers by parentCacheGroup: %v", err)
	}
	params.Del("parentCacheGroup")
}

func UniqueIPProfileTestServers(t *testing.T) {
	serversResp, _, err := TOSession.GetServersWithHdr(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(serversResp.Response) < 1 {
		t.Fatal("expected more than 0 servers")
	}
	xmppId := "unique"
	server := serversResp.Response[0]
	_, _, err = TOSession.CreateServer(tc.ServerNullable{
		CommonServerProperties: tc.CommonServerProperties{
			Cachegroup: server.Cachegroup,
			CDNName:    server.CDNName,
			DomainName: util.StrPtr("mydomain"),
			FQDN:       util.StrPtr("myfqdn"),
			FqdnTime:   time.Time{},
			HostName:   util.StrPtr("myhostname"),
			HTTPSPort:  util.IntPtr(443),
			LastUpdated: &tc.TimeNoMod{
				Time:  time.Time{},
				Valid: false,
			},
			PhysLocation: server.PhysLocation,
			Profile:      server.Profile,
			StatusID:     server.StatusID,
			Type:         server.Type,
			UpdPending:   util.BoolPtr(false),
			XMPPID:       &xmppId,
		},
		Interfaces: server.Interfaces,
	})

	if err == nil {
		t.Error("expected an error when updating a server with an ipaddress that already exists on another server with the same profile")
		// Cleanup, don't want to break other tests
		pathParams := url.Values{}
		pathParams.Add("xmppid", xmppId)
		server, _, err := TOSession.GetServersWithHdr(&pathParams, nil)
		if err != nil {
			t.Fatal(err)
		}
		_, _, err = TOSession.DeleteServerByID(*server.Response[0].ID)
		if err != nil {
			t.Fatalf("unable to delete server: %v", err)
		}
	}

	var changed bool
	for i, interf := range server.Interfaces {
		if interf.Monitor {
			for j, ip := range interf.IPAddresses {
				if ip.ServiceAddress {
					server.Interfaces[i].IPAddresses[j].Address = "127.0.0.1/24"
					changed = true
				}
			}
		}
	}
	if !changed {
		t.Fatal("did not find ip address to update")
	}
	_, _, err = TOSession.UpdateServerByID(*server.ID, server)
	if err != nil {
		t.Fatalf("expected update to pass: %s", err)
	}
}

func UpdateTestServers(t *testing.T) {
	if len(testData.Servers) < 1 {
		t.Fatal("Need at least one server to test updating")
	}

	firstServer := testData.Servers[0]
	if firstServer.HostName == nil {
		t.Fatalf("First test server had nil hostname: %+v", firstServer)
	}

	hostName := *firstServer.HostName
	params := url.Values{}
	params.Add("hostName", hostName)

	// Retrieve the server by hostname so we can get the id for the Update
	resp, _, err := TOSession.GetServersWithHdr(&params, nil)
	if err != nil {
		t.Fatalf("cannot GET Server by hostname '%s': %v - %v", hostName, err, resp.Alerts)
	}
	if len(resp.Response) < 1 {
		t.Fatalf("Expected at least one server to exist by hostname '%s'", hostName)
	}
	if len(resp.Response) > 1 {
		t.Errorf("Expected exactly one server to exist by hostname '%s' - actual: %d", hostName, len(resp.Response))
		t.Logf("Testing will proceed with server: %+v", resp.Response[0])
	}

	remoteServer := resp.Response[0]
	if remoteServer.ID == nil {
		t.Fatalf("Got null ID for server '%s'", hostName)
	}

	originalHostname := *resp.Response[0].HostName
	originalXMPIDD := *resp.Response[0].XMPPID
	// Creating idParam to get server when hostname changes.
	id := fmt.Sprintf("%v", *resp.Response[0].ID)
	idParam := url.Values{}
	idParam.Add("id", id)

	infs := remoteServer.Interfaces
	if len(infs) < 1 {
		t.Fatalf("Expected server '%s' to have at least one network interface", hostName)
	}
	inf := infs[0]

	updatedServerInterface := "bond1"
	updatedServerRack := "RR 119.03"
	updatedHostName := "atl-edge-01"
	updatedXMPPID := "change-it"

	// update rack, interfaceName and hostName values on server
	inf.Name = updatedServerInterface
	infs[0] = inf
	remoteServer.Interfaces = infs
	remoteServer.Rack = &updatedServerRack
	remoteServer.HostName = &updatedHostName

	alerts, _, err := TOSession.UpdateServerByID(*remoteServer.ID, remoteServer)
	if err != nil {
		t.Fatalf("cannot UPDATE Server by ID %d (hostname '%s'): %v - %v", *remoteServer.ID, hostName, err, alerts)
	}

	// Retrieve the server to check rack, interfaceName, hostName values were updated
	resp, _, err = TOSession.GetServersWithHdr(&idParam, nil)
	if err != nil {
		t.Errorf("cannot GET Server by ID: %v - %v", *remoteServer.HostName, err)
	}
	if len(resp.Response) < 1 {
		t.Fatalf("Expected at least one server to exist by hostname '%s'", hostName)
	}
	if len(resp.Response) > 1 {
		t.Errorf("Expected exactly one server to exist by hostname '%s' - actual: %d", hostName, len(resp.Response))
		t.Logf("Testing will proceed with server: %+v", resp.Response[0])
	}

	respServer := resp.Response[0]
	infs = respServer.Interfaces
	found := false
	for _, inf = range infs {
		if inf.Name == updatedServerInterface {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected server '%s' to have an interface named '%s' after update", hostName, updatedServerInterface)
		t.Logf("Actual interfaces: %+v", infs)
	}

	if respServer.Rack == nil {
		t.Errorf("results do not match actual: null, expected: '%s'", updatedServerRack)
	} else if *respServer.Rack != updatedServerRack {
		t.Errorf("results do not match actual: '%s', expected: '%s'", *respServer.Rack, updatedServerRack)
	}

	if remoteServer.TypeID == nil {
		t.Fatalf("Cannot test server type change update; server '%s' had nil type ID", hostName)
	}

	//Check change in hostname with no change to xmppid
	if originalHostname == *respServer.HostName && originalXMPIDD == *respServer.XMPPID {
		t.Errorf("HostName didn't change. Expected: #{updatedHostName}, actual: #{originalHostname}")
	}

	//Check to verify XMPPID never gets updated
	remoteServer.XMPPID = &updatedXMPPID
	al, reqInf, err := TOSession.UpdateServerByID(*remoteServer.ID, remoteServer)
	if err != nil && reqInf.StatusCode != http.StatusBadRequest {
		t.Logf("error making sure that XMPPID does not get updated, %d (hostname '%s'): %v - %v", *remoteServer.ID, hostName, err, al)
	}

	//Change back hostname and xmppid to its original name for other tests to pass
	remoteServer.HostName = &originalHostname
	remoteServer.XMPPID = &originalXMPIDD
	alert, _, err := TOSession.UpdateServerByID(*remoteServer.ID, remoteServer)
	if err != nil {
		t.Fatalf("cannot UPDATE Server by ID %d (hostname '%s'): %v - %v", *remoteServer.ID, hostName, err, alert)
	}
	resp, _, err = TOSession.GetServersWithHdr(&params, nil)
	if err != nil {
		t.Errorf("cannot GET Server by hostName: %v - %v", originalHostname, err)
	}

	// Assign server to DS and then attempt to update to a different type
	dses, _, err := TOSession.GetDeliveryServicesNullableWithHdr(nil)
	if err != nil {
		t.Fatalf("cannot GET DeliveryServices: %v", err)
	}
	if len(dses) < 1 {
		t.Fatal("GET DeliveryServices returned no dses, must have at least 1 to test invalid type server update")
	}

	serverTypes, _, err := TOSession.GetTypesWithHdr(nil, "server")
	if err != nil {
		t.Fatalf("cannot GET Server Types: %v", err)
	}
	if len(serverTypes) < 2 {
		t.Fatal("GET Server Types returned less then 2 types, must have at least 2 to test invalid type server update")
	}
	for _, t := range serverTypes {
		if t.ID != *remoteServer.TypeID {
			remoteServer.TypeID = &t.ID
			break
		}
	}

	// Assign server to DS
	_, _, err = TOSession.CreateDeliveryServiceServers(*dses[0].ID, []int{*remoteServer.ID}, true)
	if err != nil {
		t.Fatalf("POST delivery service servers: %v", err)
	}

	// Attempt Update - should fail
	alerts, _, err = TOSession.UpdateServerByID(*remoteServer.ID, remoteServer)
	if err == nil {
		t.Errorf("expected error when updating Server Type of a server assigned to DSes")
	} else {
		t.Logf("type change update alerts: %+v", alerts)
	}
}

func DeleteTestServers(t *testing.T) {
	params := url.Values{}

	for _, server := range testData.Servers {
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
