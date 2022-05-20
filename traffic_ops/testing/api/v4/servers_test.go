package v4

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
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

func TestServers(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, DeliveryServices}, func() {
		GetTestServersIMS(t)
		currentTime := time.Now().UTC().Add(-5 * time.Second)
		time := currentTime.Format(time.RFC1123)
		var header http.Header
		header = make(map[string][]string)
		header.Set(rfc.IfModifiedSince, time)
		header.Set(rfc.IfUnmodifiedSince, time)
		UpdateTestServers(t)
		UpdateTestServersWithHeaders(t, header)
		GetTestServers(t)
		GetTestServersIMSAfterChange(t, header)
		GetTestServersQueryParameters(t)
		header = make(map[string][]string)
		etag := rfc.ETag(currentTime)
		header.Set(rfc.IfMatch, etag)
		UpdateTestServersWithHeaders(t, header)
		CreateTestBlankFields(t)
		CreateTestServerWithoutProfileID(t)
		UniqueIPProfileTestServers(t)
		UpdateTestServerStatus(t)
		LastServerInTopologyCacheGroup(t)
		GetServersForNonExistentDeliveryService(t)
		CUDServerWithLocks(t)
		GetTestPaginationSupportServers(t)
	})
}

func CUDServerWithLocks(t *testing.T) {
	resp, _, err := TOSession.GetTenants(client.RequestOptions{})
	if err != nil {
		t.Fatalf("could not GET tenants: %v", err)
	}
	if len(resp.Response) == 0 {
		t.Fatalf("didn't get any tenant in response")
	}

	// Create a new user with operations level privileges
	user1 := tc.UserV4{
		Username:             "lock_user1",
		RegistrationSent:     new(time.Time),
		LocalPassword:        util.StrPtr("test_pa$$word"),
		ConfirmLocalPassword: util.StrPtr("test_pa$$word"),
		Role:                 "operations",
	}
	user1.Email = util.StrPtr("lockuseremail@domain.com")
	user1.TenantID = resp.Response[0].ID
	user1.FullName = util.StrPtr("firstName LastName")
	_, _, err = TOSession.CreateUser(user1, client.RequestOptions{})
	if err != nil {
		t.Fatalf("could not create test user with username: %s", user1.Username)
	}
	defer ForceDeleteTestUsersByUsernames(t, []string{"lock_user1"})

	// Establish a session with the newly created non admin level user
	userSession, _, err := client.LoginWithAgent(Config.TrafficOps.URL, user1.Username, *user1.LocalPassword, true, "to-api-v4-client-tests", false, toReqTimeout)
	if err != nil {
		t.Fatalf("could not login with user lock_user1: %v", err)
	}
	if len(testData.Servers) == 0 {
		t.Fatalf("no servers to run the test on, quitting")
	}

	server := testData.Servers[0]
	server.HostName = util.StrPtr("cdn_locks_test_server")
	server.Interfaces = []tc.ServerInterfaceInfoV40{
		{
			ServerInterfaceInfo: tc.ServerInterfaceInfo{
				IPAddresses: []tc.ServerIPAddress{
					{
						Address:        "123.32.43.21",
						Gateway:        util.StrPtr("100.100.100.100"),
						ServiceAddress: true,
					},
				},
				MaxBandwidth: util.Uint64Ptr(2500),
				Monitor:      true,
				MTU:          util.Uint64Ptr(1500),
				Name:         "cdn_locks_interfaceName",
			},
			RouterHostName: "router1",
			RouterPortName: "9090",
		},
	}
	// Create a lock for this user
	_, _, err = userSession.CreateCDNLock(tc.CDNLock{
		CDN:     *server.CDNName,
		Message: util.StrPtr("test lock"),
		Soft:    util.BoolPtr(false),
	}, client.RequestOptions{})
	if err != nil {
		t.Fatalf("couldn't create cdn lock: %v", err)
	}
	// Try to create a new server on a CDN that another user has a hard lock on -> this should fail
	_, reqInf, err := TOSession.CreateServer(server, client.RequestOptions{})
	if err == nil {
		t.Error("expected an error while creating a new server for a CDN for which a hard lock is held by another user, but got nothing")
	}
	if reqInf.StatusCode != http.StatusForbidden {
		t.Errorf("expected a 403 forbidden status while creating a new server for a CDN for which a hard lock is held by another user, but got %d", reqInf.StatusCode)
	}

	// Try to create a new profile on a CDN that the same user has a hard lock on -> this should succeed
	_, reqInf, err = userSession.CreateServer(server, client.RequestOptions{})
	if err != nil {
		t.Errorf("expected no error while creating a new server for a CDN for which a hard lock is held by the same user, but got %v", err)
	}

	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("hostName", *server.HostName)
	servers, _, err := userSession.GetServers(opts)
	if err != nil {
		t.Fatalf("couldn't get server: %v", err)
	}
	if len(servers.Response) != 1 {
		t.Fatal("couldn't get exactly one server in the response, quitting")
	}
	serverID := servers.Response[0].ID
	// Try to update a server on a CDN that another user has a hard lock on -> this should fail
	servers.Response[0].DomainName = util.StrPtr("changed_domain_name")
	_, reqInf, err = TOSession.UpdateServer(*serverID, servers.Response[0], client.RequestOptions{})
	if err == nil {
		t.Error("expected an error while updating a server for a CDN for which a hard lock is held by another user, but got nothing")
	}
	if reqInf.StatusCode != http.StatusForbidden {
		t.Errorf("expected a 403 forbidden status while updating a server for a CDN for which a hard lock is held by another user, but got %d", reqInf.StatusCode)
	}

	// Try to update a server on a CDN that the same user has a hard lock on -> this should succeed
	_, reqInf, err = userSession.UpdateServer(*serverID, servers.Response[0], client.RequestOptions{})
	if err != nil {
		t.Errorf("expected no error while updating a server for a CDN for which a hard lock is held by the same user, but got %v", err)
	}

	// Try to delete a server on a CDN that another user has a hard lock on -> this should fail
	_, reqInf, err = TOSession.DeleteServer(*serverID, client.RequestOptions{})
	if err == nil {
		t.Error("expected an error while deleting a server for a CDN for which a hard lock is held by another user, but got nothing")
	}
	if reqInf.StatusCode != http.StatusForbidden {
		t.Errorf("expected a 403 forbidden status while deleting a server for a CDN for which a hard lock is held by another user, but got %d", reqInf.StatusCode)
	}

	// Try to delete a server on a CDN that the same user has a hard lock on -> this should succeed
	_, reqInf, err = userSession.DeleteServer(*serverID, client.RequestOptions{})
	if err != nil {
		t.Errorf("expected no error while deleting a server for a CDN for which a hard lock is held by the same user, but got %v", err)
	}

	// Delete the lock
	_, _, err = userSession.DeleteCDNLocks(client.RequestOptions{QueryParameters: url.Values{"cdn": []string{*server.CDNName}}})
	if err != nil {
		t.Errorf("expected no error while deleting other user's lock using admin endpoint, but got %v", err)
	}
}

func LastServerInTopologyCacheGroup(t *testing.T) {
	const cacheGroupName = "topology-mid-cg-01"
	const moveToCacheGroup = "topology-mid-cg-02"
	const topologyName = "forked-topology"
	const cdnName = "cdn2"
	const expectedLength = 1
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", cdnName)
	cdns, _, err := TOSession.GetCDNs(opts)
	if err != nil {
		t.Fatalf("unable to get CDN '%s': %v - alerts: %+v", cdnName, err, cdns.Alerts)
	}
	if len(cdns.Response) < 1 {
		t.Fatalf("Expected exactly one CDN to exist with name '%s', found: %d", cdnName, len(cdns.Response))
	}
	cdnID := cdns.Response[0].ID

	serverOpts := client.NewRequestOptions()
	serverOpts.QueryParameters.Add("cachegroupName", cacheGroupName)
	serverOpts.QueryParameters.Add("topology", topologyName)
	serverOpts.QueryParameters.Add("cdn", strconv.Itoa(cdnID))
	servers, _, err := TOSession.GetServers(serverOpts)
	if err != nil {
		t.Fatalf("getting server from CDN '%s', from Cache Group '%s', and in Topology '%s': %v - alerts: %+v", cdnName, cacheGroupName, topologyName, err, servers.Alerts)
	}
	if len(servers.Response) != expectedLength {
		t.Fatalf("expected to get %d server from cdn %s from cachegroup %s in topology %s, got %d servers", expectedLength, cdnName, cacheGroupName, topologyName, len(servers.Response))
	}
	server := servers.Response[0]
	if server.ID == nil || server.CDNID == nil || len(server.ProfileNames) == 0 || server.CachegroupID == nil || server.HostName == nil {
		t.Fatal("Traffic Ops returned a representation for a server with null or undefined ID and/or CDN ID and/or Profile Names and/or Cache Group ID and/or Host Name")
	}

	_, reqInf, err := TOSession.DeleteServer(*server.ID, client.RequestOptions{})
	if err == nil {
		t.Fatalf("expected an error deleting server with id %d, received no error", *server.ID)
	}
	if reqInf.StatusCode < http.StatusBadRequest || reqInf.StatusCode >= http.StatusInternalServerError {
		t.Fatalf("expected a 400-level error deleting server with id %d, got status code %d: %s", *server.ID, reqInf.StatusCode, err.Error())
	}

	// attempt to move it to another CDN while it's the last server in the cachegroup in its CDN
	opts.QueryParameters.Set("name", "cdn1")
	cdns, _, err = TOSession.GetCDNs(opts)
	if err != nil {
		t.Fatalf("unable to get CDN 'cdn1': %v - alerts: %+v", err, cdns.Alerts)
	}
	if len(cdns.Response) < 1 {
		t.Fatalf("Expected exactly one CDN to exist with name 'cdn1', found: %d", len(cdns.Response))
	}
	newCDNID := cdns.Response[0].ID
	oldCDNID := *server.CDNID
	server.CDNID = &newCDNID
	opts.QueryParameters.Set("name", "MID1")
	profiles, _, err := TOSession.GetProfiles(opts)
	if err != nil {
		t.Errorf("unable to get Profile 'MID1': %v - alerts: %+v", err, profiles.Alerts)
	}
	if len(profiles.Response) != 1 {
		t.Fatalf("Expected exactly one Profile to exist with name 'MID1', found: %d", len(profiles.Response))
	}
	newProfileID := profiles.Response[0].ID
	oldProfileName := server.ProfileNames[0]

	opts.QueryParameters.Set("id", strconv.Itoa(newProfileID))
	nps, _, err := TOSession.GetProfiles(opts)
	if err != nil {
		t.Fatalf("failed to query profiles: %v", err)
	}
	if len(nps.Response) != 1 {
		t.Fatalf("Expected exactly one Profile to exist, found: %d", len(profiles.Response))
	}
	server.ProfileNames = []string{nps.Response[0].Name}
	opts.QueryParameters.Del("id")

	_, _, err = TOSession.UpdateServer(*server.ID, server, client.RequestOptions{})
	if err == nil {
		t.Fatalf("changing the CDN of the last server (%s) in a CDN in a cachegroup used by a topology assigned to a delivery service(s) in that CDN - expected: error, actual: nil", *server.HostName)
	}
	server.CDNID = &oldCDNID
	server.ProfileNames = []string{oldProfileName}

	opts.QueryParameters.Set("name", moveToCacheGroup)
	cgs, _, err := TOSession.GetCacheGroups(opts)
	if err != nil {
		t.Fatalf("getting cachegroup with hostname %s: %v - alerts: %+v", moveToCacheGroup, err, cgs.Alerts)
	}
	if len(cgs.Response) != expectedLength {
		t.Fatalf("expected %d cachegroup with hostname %s, received %d cachegroups", expectedLength, moveToCacheGroup, len(cgs.Response))
	}
	if cgs.Response[0].ID == nil {
		t.Fatalf("Traffic Ops responded with Cache Group '%s' that had null or undefined ID", moveToCacheGroup)
	}

	alerts, _, err := TOSession.UpdateServer(*server.ID, server, client.RequestOptions{})
	if err != nil {
		t.Fatalf("error updating server with hostname %s without moving it to a different Cache Group: %v - alerts: %+v", *server.HostName, err, alerts.Alerts)
	}

	*server.CachegroupID = *cgs.Response[0].ID
	alerts, _, err = TOSession.UpdateServer(*server.ID, server, client.RequestOptions{})
	if err == nil {
		t.Fatalf("expected an error moving server with id %s to a different cachegroup, received no error", *server.HostName)
	}
	if reqInf.StatusCode < http.StatusBadRequest || reqInf.StatusCode >= http.StatusInternalServerError {
		t.Fatalf("expected a 400-level error moving server with id %d to a different cachegroup, got status code %d: %v - alerts: %+v", *server.ID, reqInf.StatusCode, err, alerts.Alerts)
	}
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
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("hostName", hostName)

	// Retrieve the server by hostname so we can get the id for the Update
	resp, _, err := TOSession.GetServers(opts)
	if err != nil {
		t.Fatalf("cannot get Server by hostname '%s': %v - alerts %+v", hostName, err, resp.Alerts)
	}
	if len(resp.Response) < 1 {
		t.Fatalf("Expected at least one server to exist by hostname '%s'", hostName)
	}
	if len(resp.Response) > 1 {
		t.Errorf("Expected exactly one server to exist by hostname '%s' - actual: %d", hostName, len(resp.Response))
		t.Logf("Testing will proceed with server: %+v", resp.Response[0])
	}
	remoteServer := resp.Response[0]
	if remoteServer.ID == nil || remoteServer.HostName == nil || remoteServer.StatusLastUpdated == nil {
		t.Fatalf("Traffic Ops returned a representation for server '%s' with null or undefined ID and/or Host Name and/or Status Last Updated time", hostName)
	}
	id := fmt.Sprintf("%v", *remoteServer.ID)
	originalStatusID := 0
	updatedStatusID := 0

	statuses, _, err := TOSession.GetStatuses(client.RequestOptions{})
	if err != nil {
		t.Fatalf("cannot get Statuses: %v - alerts: %+v", err, statuses.Alerts)
	}
	for _, status := range statuses.Response {
		if status.Name == "REPORTED" {
			originalStatusID = status.ID
		}
		if status.Name == "ONLINE" {
			updatedStatusID = status.ID
		}
	}
	// Keeping the status same, perform an update and make sure that statusLastUpdated didnt change
	remoteServer.StatusID = &originalStatusID

	alerts, _, err := TOSession.UpdateServer(*remoteServer.ID, remoteServer, client.RequestOptions{})
	if err != nil {
		t.Fatalf("cannot UPDATE Server by ID %d (hostname '%s'): %v - alerts: %+v", *remoteServer.ID, hostName, err, alerts)
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("id", id)
	resp, _, err = TOSession.GetServers(opts)
	if err != nil {
		t.Errorf("cannot get Server #%s by ID: %v - alerts %+v", id, err, resp.Alerts)
	}
	if len(resp.Response) < 1 {
		t.Fatalf("Expected at least one server to exist by hostname '%s'", hostName)
	}
	if len(resp.Response) > 1 {
		t.Errorf("Expected exactly one server to exist by hostname '%s' - actual: %d", hostName, len(resp.Response))
		t.Logf("Testing will proceed with server: %+v", resp.Response[0])
	}

	respServer := resp.Response[0]
	if respServer.StatusLastUpdated == nil {
		t.Fatal("Traffic Ops returned a representation for a server with null or undefined Status Last Updated time")
	}

	if !remoteServer.StatusLastUpdated.Equal(*respServer.StatusLastUpdated) {
		t.Errorf("since status didnt change, no change in 'StatusLastUpdated' time was expected. Difference observer: old value: %v, new value: %v",
			remoteServer.StatusLastUpdated.String(), respServer.StatusLastUpdated.String())
	}

	// Changing the status, perform an update and make sure that statusLastUpdated changed
	remoteServer.StatusID = &updatedStatusID

	alerts, _, err = TOSession.UpdateServer(*remoteServer.ID, remoteServer, client.RequestOptions{})
	if err != nil {
		t.Fatalf("cannot update Server #%d (hostname '%s'): %v - alerts %+v", *remoteServer.ID, hostName, err, alerts.Alerts)
	}

	resp, _, err = TOSession.GetServers(opts)
	if err != nil {
		t.Errorf("cannot get Server by ID: %v - %v", *remoteServer.HostName, err)
	}
	if len(resp.Response) < 1 {
		t.Fatalf("Expected at least one server to exist by hostname '%s'", hostName)
	}
	if len(resp.Response) > 1 {
		t.Errorf("Expected exactly one server to exist by hostname '%s' - actual: %d", hostName, len(resp.Response))
		t.Logf("Testing will proceed with server: %+v", resp.Response[0])
	}

	respServer = resp.Response[0]
	if respServer.StatusLastUpdated == nil {
		t.Fatal("Traffic Ops returned a representation for a server with null or undefined Status Last Updated time")
	}

	if *remoteServer.StatusLastUpdated == *respServer.StatusLastUpdated {
		t.Errorf("since status was changed, expected to see a time difference between the old and new 'StatusLastUpdated' values, got the same value")
	}

	// Changing the status, perform an update and make sure that statusLastUpdated changed
	remoteServer.StatusID = &originalStatusID

	alerts, _, err = TOSession.UpdateServer(*remoteServer.ID, remoteServer, client.RequestOptions{})
	if err != nil {
		t.Fatalf("cannot update Server by ID %d (hostname '%s'): %v - alerts: %+v", *remoteServer.ID, hostName, err, alerts)
	}

	resp, _, err = TOSession.GetServers(opts)
	if err != nil {
		t.Errorf("cannot get Server by ID %d: %v - alerts: %+v", *remoteServer.ID, err, resp.Alerts)
	}
	if len(resp.Response) < 1 {
		t.Fatalf("Expected at least one server to exist by hostname '%s'", hostName)
	}
	if len(resp.Response) > 1 {
		t.Errorf("Expected exactly one server to exist by hostname '%s' - actual: %d", hostName, len(resp.Response))
		t.Logf("Testing will proceed with server: %+v", resp.Response[0])
	}

	respServer = resp.Response[0]
	if respServer.StatusLastUpdated == nil {
		t.Fatal("Traffic Ops returned a representation for a server with null or undefined Status Last Updated time")
	}

	if *remoteServer.StatusLastUpdated == *respServer.StatusLastUpdated {
		t.Errorf("since status was changed, expected to see a time difference between the old and new 'StatusLastUpdated' values, got the same value")
	}
}

func UpdateTestServersWithHeaders(t *testing.T, header http.Header) {
	if len(testData.Servers) < 1 {
		t.Fatal("Need at least one server to test updating")
	}

	firstServer := testData.Servers[0]
	if firstServer.HostName == nil {
		t.Fatalf("First test server had nil hostname: %+v", firstServer)
	}

	hostName := *firstServer.HostName
	opts := client.NewRequestOptions()
	opts.QueryParameters.Add("hostName", hostName)
	opts.Header = header

	// Retrieve the server by hostname so we can get the id for the Update
	resp, _, err := TOSession.GetServers(opts)
	if err != nil {
		t.Fatalf("cannot get Server by hostname '%s': %v - alerts: %+v", hostName, err, resp.Alerts)
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
		t.Fatalf("Got null or undefined ID for server '%s'", hostName)
	}

	infs := remoteServer.Interfaces
	if len(infs) < 1 {
		t.Fatalf("Expected server '%s' to have at least one network interface", hostName)
	}
	inf := infs[0]

	updatedServerInterface := "bond1"
	updatedServerRack := "RR 119.03"
	updatedHostName := "atl-edge-01"

	// update rack, interfaceName and hostName values on server
	inf.Name = updatedServerInterface
	infs[0] = inf
	remoteServer.Interfaces = infs
	remoteServer.Rack = &updatedServerRack
	remoteServer.HostName = &updatedHostName

	opts.QueryParameters = nil
	_, reqInf, err := TOSession.UpdateServer(*remoteServer.ID, remoteServer, opts)
	if err == nil {
		t.Errorf("Expected error about precondition failed, but got none")
	}
	if reqInf.StatusCode != http.StatusPreconditionFailed {
		t.Errorf("Expected status code 412, got %v", reqInf.StatusCode)
	}
}

func GetTestServersIMSAfterChange(t *testing.T, header http.Header) {
	opts := client.NewRequestOptions()
	opts.Header = header
	for _, server := range testData.Servers {
		if server.HostName == nil {
			t.Errorf("found server with nil hostname: %+v", server)
			continue
		}
		opts.QueryParameters.Set("hostName", *server.HostName)
		resp, reqInf, err := TOSession.GetServers(opts)
		if err != nil {
			t.Fatalf("Expected no error, but got: %v - alerts: %+v", err, resp.Alerts)
		}
		if reqInf.StatusCode != http.StatusOK {
			t.Fatalf("Expected 200 status code, got %v", reqInf.StatusCode)
		}
	}
	currentTime := time.Now().UTC()
	currentTime = currentTime.Add(1 * time.Second)
	timeStr := currentTime.Format(time.RFC1123)

	opts.Header.Set(rfc.IfModifiedSince, timeStr)
	for _, server := range testData.Servers {
		if server.HostName == nil {
			t.Errorf("found server with nil hostname: %+v", server)
			continue
		}
		opts.QueryParameters.Set("hostName", *server.HostName)
		resp, reqInf, err := TOSession.GetServers(opts)
		if err != nil {
			t.Fatalf("Expected no error, but got: %v - alerts: %+v", err, resp.Alerts)
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func GetTestServersIMS(t *testing.T) {
	futureTime := time.Now().AddDate(0, 0, 1)
	timestamp := futureTime.Format(time.RFC1123)

	opts := client.NewRequestOptions()
	opts.Header.Set(rfc.IfModifiedSince, timestamp)
	for _, server := range testData.Servers {
		if server.HostName == nil {
			t.Errorf("found server with nil hostname: %+v", server)
			continue
		}
		opts.QueryParameters.Set("hostName", *server.HostName)
		resp, reqInf, err := TOSession.GetServers(opts)
		if err != nil {
			t.Fatalf("Expected no error, but got: %v - alerts: %+v", err, resp.Alerts)
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
		resp, _, err := TOSession.CreateServer(server, client.RequestOptions{})
		if err != nil {
			t.Errorf("could not create server '%s': %v - alerts: %+v", *server.HostName, err, resp.Alerts)
		}
	}
}

func CreateTestBlankFields(t *testing.T) {
	serverResp, _, err := TOSession.GetServers(client.RequestOptions{})
	if err != nil {
		t.Fatalf("couldnt get servers: %v - alerts: %+v", err, serverResp.Alerts)
	}
	if len(serverResp.Response) < 1 {
		t.Fatal("expected at least one server")
	}
	server := serverResp.Response[0]
	if server.ID == nil {
		t.Fatal("Traffic Ops returned a representation for a servver with null or undefined ID")
	}
	originalHost := server.HostName

	server.HostName = util.StrPtr("")
	_, _, err = TOSession.UpdateServer(*server.ID, server, client.RequestOptions{})
	if err == nil {
		t.Error("should not be able to update server with blank HostName")
	}

	server.HostName = originalHost
	server.DomainName = util.StrPtr("")
	_, _, err = TOSession.UpdateServer(*server.ID, server, client.RequestOptions{})
	if err == nil {
		t.Error("should not be able to update server with blank DomainName")
	}
}

// This test will break if the structure of the test data servers collection
// is changed at all.
func CreateTestServerWithoutProfileID(t *testing.T) {
	if len(testData.Servers) < 20 {
		t.Fatal("Need at least 20 servers to test creating a server without a Profile")
	}
	testServer := testData.Servers[19]
	if testServer.HostName == nil {
		t.Fatal("Found a server in the test data with null or undefined Host Name")
	}

	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("hostName", *testServer.HostName)

	resp, _, err := TOSession.GetServers(opts)
	if err != nil {
		t.Fatalf("cannot get Server by Host Name '%s': %v - alerts: %+v", *testServer.HostName, err, resp.Alerts)
	}

	server := resp.Response[0]
	if len(server.ProfileNames) == 0 || server.ID == nil || server.HostName == nil {
		t.Fatal("Traffic Ops returned a representation of a server with null or undefined ID and/or Profile and/or Host Name")
	}
	originalProfile := server.ProfileNames
	delResp, _, err := TOSession.DeleteServer(*server.ID, client.RequestOptions{})
	if err != nil {
		t.Fatalf("cannot delete Server by ID %d: %v - %v", *server.ID, err, delResp)
	}

	server.ProfileNames = []string{""}
	_, reqInfo, _ := TOSession.CreateServer(server, client.RequestOptions{})
	if reqInfo.StatusCode != 400 {
		t.Fatalf("Expected status code: %v but got: %v", "400", reqInfo.StatusCode)
	}

	//Reverting it back for further tests
	server.ProfileNames = originalProfile
	response, _, err := TOSession.CreateServer(server, client.RequestOptions{})
	if err != nil {
		t.Fatalf("could not create server: %v - alerts: %+v", err, response.Alerts)
	}
}

func GetTestServers(t *testing.T) {
	opts := client.NewRequestOptions()
	for _, server := range testData.Servers {
		if server.HostName == nil {
			t.Errorf("found server with nil hostname: %+v", server)
			continue
		}
		opts.QueryParameters.Set("hostName", *server.HostName)
		resp, _, err := TOSession.GetServers(opts)
		if err != nil {
			t.Errorf("cannot get Server by Host Name '%s': %v - alerts: %+v", *server.HostName, err, resp.Alerts)
		} else if resp.Summary.Count != 1 {
			t.Errorf("incorrect server count, expected: 1, actual: %d", resp.Summary.Count)
		}
	}
}

func GetTestServersQueryParameters(t *testing.T) {
	dses, _, err := TOSession.GetDeliveryServices(client.RequestOptions{QueryParameters: url.Values{"xmlId": []string{"ds1"}}})
	if err != nil {
		t.Fatalf("Failed to get Delivery Services: %v - alerts: %+v", err, dses.Alerts)
	}
	if len(dses.Response) < 1 {
		t.Fatal("Failed to get at least one Delivery Service")
	}

	ds := dses.Response[0]
	if ds.ID == nil {
		t.Fatal("Traffic Ops returned a representation of a Delivery Service with null or undefined ID")
	}

	AssignTestDeliveryService(t)
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("dsId", strconv.Itoa(*ds.ID))
	servers, _, err := TOSession.GetServers(opts)
	if err != nil {
		t.Fatalf("Failed to get server by Delivery Service ID: %v - alerts: %+v", err, servers.Alerts)
	}
	if len(servers.Response) != 3 {
		t.Fatalf("expected to get 3 servers for Delivery Service: %d, actual: %d", *ds.ID, len(servers.Response))
	}

	currentTime := time.Now().UTC().Add(5 * time.Second)
	timestamp := currentTime.Format(time.RFC1123)

	opts.Header.Set(rfc.IfModifiedSince, timestamp)
	_, reqInf, _ := TOSession.GetServers(opts)
	if reqInf.StatusCode != http.StatusNotModified {
		t.Errorf("Expected a status code of 304, got %v", reqInf.StatusCode)
	}

	dses, _, err = TOSession.GetDeliveryServices(client.RequestOptions{})
	if err != nil {
		t.Fatalf("Failed to get Delivery Services: %v - alerts: %+v", err, dses.Alerts)
	}

	foundTopDs := false
	const (
		topDSXmlID = "ds-top"
		topology   = "mso-topology"
	)
	for _, ds = range dses.Response {
		if ds.XMLID == nil || ds.ID == nil {
			t.Error("Traffic Ops returned a representation of a Delivery Service that had a null or undefined XMLID and/or ID")
			continue
		}
		if *ds.XMLID != topDSXmlID {
			continue
		}
		if ds.Topology == nil || ds.FirstHeaderRewrite == nil || ds.InnerHeaderRewrite == nil || ds.LastHeaderRewrite == nil {
			t.Errorf("Traffic Ops returned a representation of Delivery Service '%s' that had a null or undefined Topology and/or First Header Rewrite text and/or Inner Header Rewrite text and/or Last Header Rewrite text", topDSXmlID)
			continue
		}
		foundTopDs = true
		break
	}
	if !foundTopDs {
		t.Fatalf("unable to find deliveryservice %s", topDSXmlID)
	}

	/* Create a deliveryservice server assignment that should not show up in the
	 * client.GetServers( response because ds-top is topology-based
	 */
	const otherServerHostname = "topology-edge-02"
	serverResponse, _, err := TOSession.GetServers(client.RequestOptions{QueryParameters: url.Values{"hostName": []string{otherServerHostname}}})
	if err != nil {
		t.Fatalf("getting server by Host Name %s: %v - alerts: %+v", otherServerHostname, err, serverResponse.Alerts)
	}
	if len(serverResponse.Response) != 1 {
		t.Fatalf("unable to find server with hostname %s", otherServerHostname)
	}
	otherServer := serverResponse.Response[0]
	if otherServer.ID == nil || otherServer.HostName == nil {
		t.Fatal("Traffic Ops returned a representation of a Server that had a null or undefined ID and/or Host Name")
	}

	dsTopologyField, dsFirstHeaderRewriteField, innerHeaderRewriteField, lastHeaderRewriteField := *ds.Topology, *ds.FirstHeaderRewrite, *ds.InnerHeaderRewrite, *ds.LastHeaderRewrite
	ds.Topology, ds.FirstHeaderRewrite, ds.InnerHeaderRewrite, ds.LastHeaderRewrite = nil, nil, nil, nil
	updResp, _, err := TOSession.UpdateDeliveryService(*ds.ID, ds, client.RequestOptions{})
	if err != nil {
		t.Fatalf("unable to temporary remove topology-related fields from deliveryservice '%s': %v - alerts: %+v", topDSXmlID, err, updResp.Alerts)
	}
	if len(updResp.Response) != 1 {
		t.Fatalf("Expected updating a Delivery Service to update exactly one Delivery Service, but Traffic Ops indicates that %d were updated", len(updResp.Response))
	}
	ds = updResp.Response[0]
	if ds.ID == nil {
		t.Fatal("Traffic Ops returned a representation of a Delivery Service that had null or undefined ID")
	}
	assignResp, _, err := TOSession.CreateDeliveryServiceServers(*ds.ID, []int{*otherServer.ID}, false, client.RequestOptions{})
	if err != nil {
		t.Fatalf("unable to assign server '%s' to Delivery Service '%s': %v - alerts: %+v", *otherServer.HostName, topDSXmlID, err, assignResp.Alerts)
	}
	ds.Topology, ds.FirstHeaderRewrite, ds.InnerHeaderRewrite, ds.LastHeaderRewrite = &dsTopologyField, &dsFirstHeaderRewriteField, &innerHeaderRewriteField, &lastHeaderRewriteField
	updResp, _, err = TOSession.UpdateDeliveryService(*ds.ID, ds, client.RequestOptions{})
	if err != nil {
		t.Fatalf("unable to re-add topology-related fields to deliveryservice %s: %v - alerts: %+v", topDSXmlID, err, updResp.Alerts)
	}

	opts.Header = nil
	opts.QueryParameters.Set("dsId", strconv.Itoa(*ds.ID))
	expectedHostnames := map[string]bool{
		"edge1-cdn1-cg3":                 false,
		"edge2-cdn1-cg3":                 false,
		"atlanta-mid-01":                 false,
		"atlanta-mid-16":                 false,
		"edgeInCachegroup3":              false,
		"midInSecondaryCachegroupInCDN1": false,
	}
	response, _, err := TOSession.GetServers(opts)
	if err != nil {
		t.Fatalf("Failed to get servers by Topology-based Delivery Service ID with xmlId %s: %v - alerts: %+v", topDSXmlID, err, response.Alerts)
	}
	if len(response.Response) == 0 {
		t.Fatalf("Did not find any servers for Topology-based Delivery Service with xmlId %s", topDSXmlID)
	}
	for _, server := range response.Response {
		if server.HostName == nil {
			t.Fatal("Traffic Ops responded with a representation for a server with null or undefined Host Name")
		}
		if _, exists := expectedHostnames[*server.HostName]; !exists {
			t.Fatalf("expected hostnames %v, actual %s", expectedHostnames, *server.HostName)
		}
		expectedHostnames[*server.HostName] = true
	}
	var notInResponse []string
	for hostName, inResponse := range expectedHostnames {
		if !inResponse {
			notInResponse = append(notInResponse, hostName)
		}
	}
	if len(notInResponse) != 0 {
		t.Fatalf("%d servers missing from the response: %s", len(notInResponse), strings.Join(notInResponse, ", "))
	}
	const originHostname = "denver-mso-org-01"
	if resp, _, err := TOSession.AssignServersToDeliveryService([]string{originHostname}, topDSXmlID, client.RequestOptions{}); err != nil {
		t.Fatalf("assigning origin server '%s' to Delivery Service '%s': %v - alerts: %+v", originHostname, topDSXmlID, err, resp.Alerts)
	}
	response, _, err = TOSession.GetServers(opts)
	if err != nil {
		t.Fatalf("Failed to get servers by Topology-based Delivery Service ID with xmlId %s: %v - alerts: %+v", topDSXmlID, err, response.Alerts)
	}
	if len(response.Response) == 0 {
		t.Fatalf("Did not find any servers for Topology-based Delivery Service with xmlId %s", topDSXmlID)
	}
	containsOrigin := false
	for _, server := range response.Response {
		if server.HostName == nil || *server.HostName != originHostname {
			continue
		}
		containsOrigin = true
		break
	}
	if !containsOrigin {
		t.Fatalf("did not find origin server %s when querying servers by dsId after assigning %s to delivery service %s", originHostname, originHostname, topDSXmlID)
	}

	const topDsWithNoMids = "ds-based-top-with-no-mids"
	dses, _, err = TOSession.GetDeliveryServices(client.RequestOptions{QueryParameters: url.Values{"xmlId": []string{topDsWithNoMids}}})
	if err != nil {
		t.Fatalf("Failed to get Delivery Services: %v - alerts: %+v", err, dses.Alerts)
	}
	if len(dses.Response) < 1 {
		t.Fatal("Failed to get at least one Delivery Service")
	}

	ds = dses.Response[0]
	if ds.ID == nil {
		t.Fatal("Got Delivery Service with nil ID")
	}
	opts.QueryParameters.Set("dsId", strconv.Itoa(*ds.ID))

	response, _, err = TOSession.GetServers(opts)
	if err != nil {
		t.Fatalf("Failed to get servers by Topology-based Delivery Service ID with xmlId %s: %s", topDsWithNoMids, err)
	}
	if len(response.Response) == 0 {
		t.Fatalf("Did not find any servers for Topology-based Delivery Service with xmlId %s: %s", topDsWithNoMids, err)
	}
	for _, server := range response.Response {
		if server.HostName == nil {
			t.Fatal("Traffic Ops returned a server with null or undefined Host Name")
		}
		if server.Type == tc.CacheTypeMid.String() {
			t.Fatalf("Expected to find no %s-typed servers when querying servers by the ID for Delivery Service with XMLID %s but found %s-typed server %s", tc.CacheTypeMid, topDsWithNoMids, tc.CacheTypeMid, *server.HostName)
		}
	}

	opts.QueryParameters.Del("dsId")
	opts.QueryParameters.Set("topology", topology)
	expectedHostnames = map[string]bool{
		originHostname:                   false,
		"denver-mso-org-02":              false,
		"edge1-cdn1-cg3":                 false,
		"edge2-cdn1-cg3":                 false,
		"atlanta-mid-01":                 false,
		"atlanta-mid-16":                 false,
		"atlanta-mid-17":                 false,
		"edgeInCachegroup3":              false,
		"midInParentCachegroup":          false,
		"midInSecondaryCachegroup":       false,
		"midInSecondaryCachegroupInCDN1": false,
		"test-mso-org-01":                false,
	}
	response, _, err = TOSession.GetServers(opts)
	if err != nil {
		t.Fatalf("Failed to get servers belonging to Cache Groups in Topology %s: %v - alerts: %+v", topology, err, response.Alerts)
	}
	if len(response.Response) == 0 {
		t.Fatalf("Did not find any servers belonging to Cache Groups in Topology %s:", topology)
	}
	for _, server := range response.Response {
		if server.HostName == nil {
			t.Fatal("Traffic Ops returned a server with null or undefined Host Name")
		}
		if _, exists := expectedHostnames[*server.HostName]; !exists {
			t.Fatalf("expected hostnames %v, actual %s", expectedHostnames, *server.HostName)
		}
		expectedHostnames[*server.HostName] = true
	}
	notInResponse = []string{}
	for hostName, inResponse := range expectedHostnames {
		if !inResponse {
			notInResponse = append(notInResponse, hostName)
		}
	}
	if len(notInResponse) != 0 {
		t.Fatalf("%d servers missing from the response: %s", len(notInResponse), strings.Join(notInResponse, ", "))
	}
	opts.QueryParameters.Del("topology")

	resp, _, err := TOSession.GetServers(client.RequestOptions{})
	if err != nil {
		t.Fatalf("Failed to get servers: %v - alerts: %+v", err, resp.Alerts)
	}

	if len(resp.Response) < 1 {
		t.Fatal("Failed to get at least one server")
	}

	s := resp.Response[0]

	opts.QueryParameters.Set("type", s.Type)
	if resp, _, err := TOSession.GetServers(opts); err != nil {
		t.Errorf("Error getting servers by Type: %v - alerts: %+v", err, resp.Alerts)
	}
	opts.QueryParameters.Del("type")

	if s.CachegroupID == nil {
		t.Error("Found server with no Cache Group ID")
	} else {
		opts.QueryParameters.Add("cachegroup", strconv.Itoa(*s.CachegroupID))
		if resp, _, err := TOSession.GetServers(opts); err != nil {
			t.Errorf("Error getting servers by Cache Group ID: %v - alerts: %+v", err, resp.Alerts)
		}
		opts.QueryParameters.Del("cachegroup")
	}

	if s.Status == nil {
		t.Error("Found server with no status")
	} else {
		opts.QueryParameters.Add("status", *s.Status)
		if resp, _, err := TOSession.GetServers(opts); err != nil {
			t.Errorf("Error getting servers by status: %v - alerts: %+v", err, resp.Alerts)
		}
		opts.QueryParameters.Del("status")
	}

	opts.QueryParameters.Add("name", s.ProfileNames[0])
	pr, _, err := TOSession.GetProfiles(opts)
	if err != nil {
		t.Fatalf("failed to query profile: %v", err)
	}
	if len(pr.Response) != 1 {
		t.Error("Found server with no Profile ID")
	} else {
		profileName := pr.Response[0].Name
		opts.QueryParameters.Add("profileName", profileName)
		if resp, _, err := TOSession.GetServers(opts); err != nil {
			t.Errorf("Error getting servers by Profile ID: %v - alerts: %+v", err, resp.Alerts)
		}
		opts.QueryParameters.Del("profileName")
	}

	cgs, _, err := TOSession.GetCacheGroups(client.RequestOptions{})
	if err != nil {
		t.Fatalf("Failed to get Cache Groups: %v", err)
	}
	if len(cgs.Response) < 1 {
		t.Fatal("Failed to get at least one Cache Group")
	}
	if cgs.Response[0].ID == nil {
		t.Fatal("Cache Group found with no ID")
	}

	opts.QueryParameters.Add("parentCacheGroup", strconv.Itoa(*cgs.Response[0].ID))
	if resp, _, err = TOSession.GetServers(opts); err != nil {
		t.Errorf("Error getting servers by parent Cache Group: %v - alerts: %+v", err, resp.Alerts)
	}
	opts.QueryParameters.Del("parentCacheGroup")
}

func UniqueIPProfileTestServers(t *testing.T) {
	serversResp, _, err := TOSession.GetServers(client.RequestOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(serversResp.Response) < 1 {
		t.Fatal("expected more than 0 servers")
	}
	xmppID := "unique"
	var server tc.ServerV40
	for _, s := range serversResp.Response {
		if len(s.Interfaces) >= 1 && s.Interfaces[0].Monitor {
			server = s
			break
		}
	}
	_, _, err = TOSession.CreateServer(tc.ServerV40{
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
		ProfileNames: server.ProfileNames,
		StatusID:     server.StatusID,
		Type:         server.Type,
		UpdPending:   util.BoolPtr(false),
		XMPPID:       &xmppID,
		Interfaces:   server.Interfaces,
	}, client.RequestOptions{})

	if err == nil {
		t.Error("expected an error when updating a server with an ipaddress that already exists on another server with the same profile")
		// Cleanup, don't want to break other tests
		opts := client.NewRequestOptions()
		opts.QueryParameters.Add("xmppid", xmppID)
		server, _, err := TOSession.GetServers(opts)
		if err != nil {
			t.Fatalf("Unexpected error getting servers filtered by XMPPID '%s': %v - alerts: %+v", xmppID, err, server.Alerts)
		}
		if len(server.Response) < 1 {
			t.Fatalf("Expected at least one server to exist with XMPPID '%s'", xmppID)
		}
		alerts, _, err := TOSession.DeleteServer(*server.Response[0].ID, client.RequestOptions{})
		if err != nil {
			t.Fatalf("unable to delete server: %v - alerts: %+v", err, alerts.Alerts)
		}
	}

	changed := false
	for i, interf := range server.Interfaces {
		if interf.Monitor {
			for j, ip := range interf.IPAddresses {
				if ip.ServiceAddress {
					server.Interfaces[i].IPAddresses[j].Address = "127.0.0.5/24"
					changed = true
				}
			}
		}
	}
	if !changed {
		t.Fatal("did not find ip address to update")
	}
	alerts, _, err := TOSession.UpdateServer(*server.ID, server, client.RequestOptions{})
	if err != nil {
		t.Fatalf("expected update to pass: %v - alerts: %+v", err, alerts.Alerts)
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
	opts := client.NewRequestOptions()
	opts.QueryParameters.Add("hostName", hostName)

	// Retrieve the server by hostname so we can get the id for the Update
	resp, _, err := TOSession.GetServers(opts)
	if err != nil {
		t.Fatalf("cannot get Server by hostname '%s': %v - alerts: %+v", hostName, err, resp.Alerts)
	}
	if len(resp.Response) < 1 {
		t.Fatalf("Expected at least one server to exist by hostname '%s'", hostName)
	}
	if len(resp.Response) > 1 {
		t.Errorf("Expected exactly one server to exist by hostname '%s' - actual: %d", hostName, len(resp.Response))
		t.Logf("Testing will proceed with server: %+v", resp.Response[0])
	}

	remoteServer := resp.Response[0]
	if remoteServer.ID == nil || remoteServer.HostName == nil || remoteServer.XMPPID == nil {
		t.Fatalf("Traffic Ops returned a representation for server '%s' with null or undefined ID and/or Host Name and/or XMPPID", hostName)
	}

	originalHostname := *remoteServer.HostName
	originalXMPIDD := *remoteServer.XMPPID
	originalDomainName := remoteServer.DomainName
	originalHttpsPort := remoteServer.HTTPSPort
	originalTcpPort := remoteServer.TCPPort
	originalupdpending := remoteServer.UpdPending
	originalCgid := remoteServer.CachegroupID
	originalPhysLocId := remoteServer.PhysLocationID
	originalTypeId := remoteServer.TypeID

	//Get cachegroup id
	if len(testData.CacheGroups) < 1 {
		t.Fatal("Need at least one Cache Group to test updating servers")
	}
	firstCG := testData.CacheGroups[0]
	if firstCG.Name == nil {
		t.Fatal("Found Cache Group with null or undefined name in testing data")
	}
	opts = client.NewRequestOptions()
	opts.QueryParameters.Set("name", *firstCG.Name)
	cachegroupResp, _, err := TOSession.GetCacheGroups(opts)
	if err != nil {
		t.Fatalf("cannot get Cache Group '%s': %v - alerts: %+v", *firstCG.Name, err, cachegroupResp.Alerts)
	}
	if len(cachegroupResp.Response) != 1 {
		t.Fatalf("Expected exactly one Cache Group to exist with name '%s', but got: %d", *firstCG.Name, len(cachegroupResp.Response))
	}
	cg := cachegroupResp.Response[0]
	if cg.ID == nil {
		t.Fatalf("Traffic Ops returned Cache Group '%s' with null or undefined ID", *cg.Name)
	}

	// Retrieve the PhysLocation ID by name
	if len(testData.PhysLocations) < 1 {
		t.Fatal("Need at least one Physical Location to test updating servers")
	}
	firstPhysLocation := testData.PhysLocations[0]
	if firstPhysLocation.Name == "" {
		t.Fatal("Found Physical location with null or undefined name in testing data")
	}
	opts = client.NewRequestOptions()
	opts.QueryParameters.Set("name", firstPhysLocation.Name)
	physicalLocResp, _, err := TOSession.GetPhysLocations(opts)
	if err != nil {
		t.Errorf("cannot get Physical Location by name '%s': %v - alerts: %+v", firstPhysLocation.Name, err, physicalLocResp.Alerts)
	}
	if len(physicalLocResp.Response) != 1 {
		t.Fatalf("Expected exactly one Physical Location to exist with name '%s', found: %d", firstPhysLocation.Name, len(physicalLocResp.Response))
	}
	phylocation := physicalLocResp.Response[0]

	// Retrieve the type ID by useInTable
	opts = client.NewRequestOptions()
	opts.QueryParameters.Set("useInTable", "server")
	typeResp, _, err := TOSession.GetTypes(opts)
	if err != nil {
		t.Errorf("cannot get Types by useInTable '%s': %v - alerts: %+v", "server", err, typeResp.Alerts)
	}
	if len(typeResp.Response) < 1 {
		t.Fatalf("Expected atleast one Types to exist with useInTable '%s', found: %d", "server", len(typeResp.Response))
	}
	types := typeResp.Response[0]

	// Creating idParam to get server when hostname changes.
	id := fmt.Sprintf("%v", *remoteServer.ID)
	idOpts := client.NewRequestOptions()
	idOpts.QueryParameters.Add("id", id)

	infs := remoteServer.Interfaces
	if len(infs) < 1 {
		t.Fatalf("Expected server '%s' to have at least one network interface", hostName)
	}
	inf := infs[0]
	if remoteServer.Interfaces[0].MTU == nil {
		t.Fatalf("got null value for interface MTU related to server %s", hostName)
	}
	originalMTU := *remoteServer.Interfaces[0].MTU

	updatedServerInterface := "bond1"
	updatedServerRack := "RR 119.03"
	updatedHostName := "atl-edge-01"
	updatedXMPPID := "change-it"
	updatedMTU := uint64(1280)
	updateDomainName := "updateddomainname"
	updateHttpsPort := 8080
	updatedTcpPort := 8080
	updatedPending := true

	// update rack, interfaceName and hostName values on server
	inf.Name = updatedServerInterface
	infs[0] = inf
	remoteServer.Interfaces = infs
	remoteServer.Rack = &updatedServerRack
	remoteServer.HostName = &updatedHostName
	remoteServer.Interfaces[0].MTU = &updatedMTU
	remoteServer.CachegroupID = cg.ID
	remoteServer.DomainName = &updateDomainName
	remoteServer.HTTPSPort = &updateHttpsPort
	remoteServer.PhysLocationID = &phylocation.ID
	remoteServer.TCPPort = &updatedTcpPort
	remoteServer.TypeID = &types.ID
	remoteServer.UpdPending = &updatedPending

	alerts, _, err := TOSession.UpdateServer(*remoteServer.ID, remoteServer, client.RequestOptions{})
	if err != nil {
		t.Fatalf("cannot update Server by ID %d (hostname '%s'): %v - alerts: %+v", *remoteServer.ID, hostName, err, alerts.Alerts)
	}

	// Retrieve the server to check rack, interfaceName, hostName and MTU values were updated
	resp, _, err = TOSession.GetServers(idOpts)
	if err != nil {
		t.Errorf("cannot get Server: %v - alerts: %+v", err, resp.Alerts)
	}
	if len(resp.Response) < 1 {
		t.Fatalf("Expected at least one server to exist by hostname '%s'", hostName)
	}
	if len(resp.Response) > 1 {
		t.Errorf("Expected exactly one server to exist by hostname '%s' - actual: %d", hostName, len(resp.Response))
		t.Logf("Testing will proceed with server: %+v", resp.Response[0])
	}

	respServer := resp.Response[0]
	if respServer.HostName == nil || respServer.XMPPID == nil {
		t.Fatal("Traffic Ops returned a representation for a server with null or undefined Host Name and/or XMPPID")
	}
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

	// Check to verify mtu changed
	if len(respServer.Interfaces) >= 1 {
		if respServer.Interfaces[0].MTU != nil {
			if originalMTU == *respServer.Interfaces[0].MTU {
				t.Errorf("MTU value didn't update. Expected: %v, actual: %v", updatedMTU, originalMTU)
			}
		}
	}

	//Check change in hostname with no change to xmppid
	if originalHostname == *respServer.HostName && originalXMPIDD == *respServer.XMPPID {
		t.Errorf("HostName didn't change. Expected: %s, actual: %s", updatedHostName, originalHostname)
	}

	if *respServer.CachegroupID != *remoteServer.CachegroupID {
		t.Errorf("Cachegroup ID is not updated while updating the servers")
	}
	if *respServer.DomainName != *remoteServer.DomainName {
		t.Errorf("DomainName is not updated while updating the servers")
	}
	if *respServer.HTTPSPort != *remoteServer.HTTPSPort {
		t.Errorf("Https Port ID is not updated while updating the servers")
	}
	if *respServer.PhysLocationID != *remoteServer.PhysLocationID {
		t.Errorf("Physical Location ID is not updated while updating the servers")
	}
	if *respServer.TCPPort != *remoteServer.TCPPort {
		t.Errorf("TCP Port is not updated while updating the servers")
	}
	if *respServer.TypeID != *remoteServer.TypeID {
		t.Errorf("Type ID is not updated while updating the servers")
	}

	//Check to verify XMPPID never gets updated
	remoteServer.XMPPID = &updatedXMPPID
	al, reqInf, err := TOSession.UpdateServer(*remoteServer.ID, remoteServer, client.RequestOptions{})
	if err != nil && reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("error making sure that XMPPID does not get updated, %d (hostname '%s'): %v - %v", *remoteServer.ID, hostName, err, al.Alerts)
	}

	//Change back hostname, xmppid, mtu to its original name for other tests to pass
	remoteServer.HostName = &originalHostname
	remoteServer.XMPPID = &originalXMPIDD
	remoteServer.Interfaces[0].MTU = &originalMTU
	remoteServer.CachegroupID = originalCgid
	remoteServer.DomainName = originalDomainName
	remoteServer.HTTPSPort = originalHttpsPort
	remoteServer.PhysLocationID = originalPhysLocId
	remoteServer.TCPPort = originalTcpPort
	remoteServer.TypeID = originalTypeId
	remoteServer.UpdPending = originalupdpending

	alert, _, err := TOSession.UpdateServer(*remoteServer.ID, remoteServer, client.RequestOptions{})
	if err != nil {
		t.Fatalf("cannot UPDATE Server by ID %d (hostname '%s'): %v - %v", *remoteServer.ID, hostName, err, alert)
	}
	resp, _, err = TOSession.GetServers(opts)
	if err != nil {
		t.Errorf("cannot get Server by Host Name '%s': %v - alerts: %+v", originalHostname, err, resp.Alerts)
	}

	// Assign server to DS and then attempt to update to a different type
	dses, _, err := TOSession.GetDeliveryServices(client.RequestOptions{})
	if err != nil {
		t.Fatalf("cannot get Delivery Services: %v - alerts: %+v", err, dses.Alerts)
	}
	if len(dses.Response) < 1 {
		t.Fatal("GET DeliveryServices returned no dses, must have at least 1 to test invalid type server update")
	}
	ds := dses.Response[0]
	if ds.ID == nil {
		t.Fatal("Traffic Ops returned a representation of a Delivery Servvice with a null or undefined ID")
	}

	typeOpts := client.NewRequestOptions()
	typeOpts.QueryParameters.Set("useInTable", "server")
	serverTypes, _, err := TOSession.GetTypes(typeOpts)
	if err != nil {
		t.Fatalf("cannot get Server Types: %v - alerts: %+v", err, serverTypes.Alerts)
	}
	if len(serverTypes.Response) < 2 {
		t.Fatal("GET Server Types returned less then 2 types, must have at least 2 to test invalid type server update")
	}
	for _, t := range serverTypes.Response {
		if t.ID != *remoteServer.TypeID {
			remoteServer.TypeID = &t.ID
			break
		}
	}

	// Assign server to DS
	assignResp, _, err := TOSession.CreateDeliveryServiceServers(*ds.ID, []int{*remoteServer.ID}, true, client.RequestOptions{})
	if err != nil {
		t.Fatalf("Unexpected error creating server-to-Delivery-Service assignments: %v - alerts: %+v", err, assignResp.Alerts)
	}

	// Attempt Update - should fail
	alerts, _, err = TOSession.UpdateServer(*remoteServer.ID, remoteServer, client.RequestOptions{})
	if err == nil {
		t.Errorf("expected error when updating Server Type of a server assigned to DSes")
	} else {
		t.Logf("got expected error when updating Server Type of a server assigned to DSes - type change update alerts: %+v, err: %v", alerts, err)
	}
}

func DeleteTestServers(t *testing.T) {
	opts := client.NewRequestOptions()

	for _, server := range testData.Servers {
		if server.HostName == nil {
			t.Errorf("found server with nil hostname: %+v", server)
			continue
		}

		opts.QueryParameters.Set("hostName", *server.HostName)

		resp, _, err := TOSession.GetServers(opts)
		if err != nil {
			t.Errorf("cannot get Server by Host Name '%s': %v - alerts: %+v", *server.HostName, err, resp.Alerts)
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

			delResp, _, err := TOSession.DeleteServer(*respServer.ID, client.RequestOptions{})
			if err != nil {
				t.Errorf("cannot delete Server by ID %d: %v - alerts: %+v", *respServer.ID, err, delResp.Alerts)
				continue
			}

			// Retrieve the Server to see if it got deleted
			resp, _, err := TOSession.GetServers(opts)
			if err != nil {
				t.Errorf("error filtering Servers by hostname '%s' after supposed deletion: %v - alerts: %+v", *server.HostName, err, resp.Alerts)
			}
			if len(resp.Response) > 0 {
				t.Errorf("expected Server hostname: %s to be deleted", *server.HostName)
			}
		}
	}
}

func GetServersForNonExistentDeliveryService(t *testing.T) {
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("dsId", "999999")
	resp, reqInf, err := TOSession.GetServers(opts)
	if err != nil {
		t.Errorf("error getting the servers for DS with ID %d: %v - alerts: %+v", 999999, err, resp.Alerts)
	}
	if reqInf.StatusCode != http.StatusOK {
		t.Errorf("expected status code of 200, but got %d", reqInf.StatusCode)
	}
	if len(resp.Response) != 0 {
		t.Errorf("expected an empty list of servers associated with a non existent DS, but got %d servers", len(resp.Response))
	}
}

func GetTestPaginationSupportServers(t *testing.T) {
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("orderby", "id")
	resp, _, err := TOSession.GetServers(opts)
	if err != nil {
		t.Fatalf("cannot Get Servers: %v - alerts: %+v", err, resp.Alerts)
	}
	server := resp.Response
	if len(server) < 3 {
		t.Fatalf("Need at least 3 server in Traffic Ops to test pagination support, found: %d", len(server))
	}

	opts.QueryParameters.Set("limit", "1")
	serversWithLimit, _, err := TOSession.GetServers(opts)
	if err != nil {
		t.Fatalf("cannot Get server with Limit: %v - alerts: %+v", err, serversWithLimit.Alerts)
	}
	if !reflect.DeepEqual(server[:1], serversWithLimit.Response) {
		t.Error("expected GET server with limit = 1 to return first result")
	}

	opts.QueryParameters.Set("offset", "1")
	serversWithOffset, _, err := TOSession.GetServers(opts)
	if err != nil {
		t.Fatalf("cannot Get server with Limit and Offset: %v - alerts: %+v", err, serversWithOffset.Alerts)
	}
	if !reflect.DeepEqual(server[1:2], serversWithOffset.Response) {
		t.Error("expected GET server with limit = 1, offset = 1 to return second result")
	}

	opts.QueryParameters.Del("offset")
	opts.QueryParameters.Set("page", "2")
	serversWithPage, _, err := TOSession.GetServers(opts)
	if err != nil {
		t.Fatalf("cannot Get server with Limit and Page: %v - alerts: %+v", err, serversWithPage.Alerts)
	}
	if !reflect.DeepEqual(server[1:2], serversWithPage.Response) {
		t.Error("expected GET server with limit = 1, page = 2 to return second result")
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("limit", "-2")
	resp, _, err = TOSession.GetServers(opts)
	if err == nil {
		t.Error("expected GET server to return an error when limit is not bigger than -1")
	} else if !alertsHaveError(resp.Alerts.Alerts, "must be bigger than -1") {
		t.Errorf("expected GET server to return an error for limit is not bigger than -1, actual error: %v - alerts: %+v", err, resp.Alerts)
	}

	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("offset", "0")
	resp, _, err = TOSession.GetServers(opts)
	if err == nil {
		t.Error("expected GET server to return an error when offset is not a positive integer")
	} else if !alertsHaveError(resp.Alerts.Alerts, "must be a positive integer") {
		t.Errorf("expected GET server to return an error for offset is not a positive integer, actual error: %v - alerts: %+v", err, resp.Alerts)
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("page", "0")
	resp, _, err = TOSession.GetServers(opts)
	if err == nil {
		t.Error("expected GET server to return an error when page is not a positive integer")
	} else if !alertsHaveError(resp.Alerts.Alerts, "must be a positive integer") {
		t.Errorf("expected GET server to return an error for page is not a positive integer, actual error: %v - alerts: %+v", err, resp.Alerts)
	}
}
