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
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

func TestCDNLocks(t *testing.T) {
	WithObjs(t, []TCObj{Types, CacheGroups, CDNs, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, Servers, ServerCapabilities, ServerServerCapabilitiesForTopologies, Topologies, Tenants, DeliveryServices, TopologyBasedDeliveryServiceRequiredCapabilities, Roles, Users}, func() {
		CRDCdnLocks(t)
		AdminCdnLocks(t)
		SnapshotWithLock(t)
		QueueUpdatesWithLock(t)
		QueueUpdatesFromTopologiesWithLock(t)
	})
}

func getCDNName(t *testing.T) string {
	cdnResp, _, err := TOSession.GetCDNs(client.RequestOptions{})
	if err != nil {
		t.Fatalf("couldn't get CDNs: %v", err)
	}
	if len(cdnResp.Response) < 1 {
		t.Fatalf("no valid CDNs in response")
	}
	return cdnResp.Response[0].Name
}

func getCDNNameAndServerID(t *testing.T) (string, int) {
	serverID := -1
	cdnResp, _, err := TOSession.GetCDNs(client.RequestOptions{})
	if err != nil {
		t.Fatalf("couldn't get CDNs: %v", err)
	}
	if len(cdnResp.Response) < 1 {
		t.Fatalf("no valid CDNs in response")
	}
	for _, cdn := range cdnResp.Response {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("cdn", strconv.Itoa(cdn.ID))
		serversResp, _, err := TOSession.GetServers(opts)
		if err != nil {
			t.Errorf("could not get servers for cdn %s: %v", cdn.Name, err)
		}
		if len(serversResp.Response) != 0 {
			serverID = *serversResp.Response[0].ID
			return cdn.Name, serverID
		}
	}
	return "", serverID
}

func getCDNDetailsAndTopologyName(t *testing.T) (int, string, string) {
	opts := client.NewRequestOptions()
	topologiesResp, _, err := TOSession.GetTopologies(client.RequestOptions{})
	if err != nil {
		t.Fatalf("couldn't get topologies, err: %v", err)
	}
	if len(topologiesResp.Response) == 0 {
		t.Fatal("no topologies returned")
	}
	for _, top := range topologiesResp.Response {
		for _, node := range top.Nodes {
			opts.QueryParameters.Set("name", node.Cachegroup)
			cacheGroupResp, _, err := TOSession.GetCacheGroups(opts)
			if err != nil {
				t.Errorf("error while GETting cachegroups: %v", err)
			}
			if len(cacheGroupResp.Response) != 0 && cacheGroupResp.Response[0].ID != nil {
				cacheGroupID := *cacheGroupResp.Response[0].ID
				opts.QueryParameters.Del("name")
				opts.QueryParameters.Set("cachegroup", strconv.Itoa(cacheGroupID))
				serversResp, _, err := TOSession.GetServers(opts)
				if err != nil {
					t.Errorf("couldn't get servers: %v", err)
				}
				if len(serversResp.Response) != 0 && serversResp.Response[0].CDNName != nil && serversResp.Response[0].CDNID != nil {
					return *serversResp.Response[0].CDNID, *serversResp.Response[0].CDNName, top.Name
				}
			}
		}
	}
	return -1, "", ""
}

func CRDCdnLocks(t *testing.T) {
	cdn := getCDNName(t)
	// CREATE
	var cdnLock tc.CDNLock
	cdnLock.CDN = cdn
	cdnLock.UserName = TOSession.UserName
	cdnLock.Message = util.StrPtr("snapping cdn")
	cdnLock.Soft = util.BoolPtr(true)
	cdnLockResp, _, err := TOSession.CreateCDNLock(cdnLock, client.RequestOptions{})
	if err != nil {
		t.Fatalf("couldn't create cdn lock: %v", err)
	}
	if cdnLockResp.Response.UserName != cdnLock.UserName {
		t.Errorf("expected username %v, got %v", cdnLock.UserName, cdnLockResp.Response.UserName)
	}
	if cdnLockResp.Response.CDN != cdnLock.CDN {
		t.Errorf("expected cdn %v, got %v", cdnLock.CDN, cdnLockResp.Response.CDN)
	}
	if cdnLockResp.Response.Message == nil {
		t.Errorf("expected a valid message, but got nothing")
	}
	if cdnLockResp.Response.Message != nil && *cdnLockResp.Response.Message != *cdnLock.Message {
		t.Errorf("expected Message %v, got %v", *cdnLock.Message, *cdnLockResp.Response.Message)
	}
	if cdnLockResp.Response.Soft == nil {
		t.Errorf("expected a valid soft/hard setting, but got nothing")
	}
	if cdnLockResp.Response.Soft != nil && *cdnLockResp.Response.Soft != *cdnLock.Soft {
		t.Errorf("expected 'Soft' to be %v, got %v", *cdnLock.Soft, *cdnLockResp.Response.Soft)
	}

	// READ
	cdnLocksReadResp, _, err := TOSession.GetCDNLocks(client.RequestOptions{})
	if err != nil {
		t.Fatalf("could not get CDN Locks: %v", err)
	}
	if len(cdnLocksReadResp.Response) != 1 {
		t.Fatalf("expected to get back one CDN lock, but got %d instead", len(cdnLocksReadResp.Response))
	}
	if cdnLocksReadResp.Response[0].UserName != cdnLock.UserName {
		t.Errorf("expected username %v, got %v", cdnLock.UserName, cdnLocksReadResp.Response[0].UserName)
	}
	if cdnLocksReadResp.Response[0].CDN != cdnLock.CDN {
		t.Errorf("expected cdn %v, got %v", cdnLock.CDN, cdnLocksReadResp.Response[0].CDN)
	}
	if cdnLocksReadResp.Response[0].Message == nil {
		t.Errorf("expected a valid message, but got nothing")
	}
	if cdnLocksReadResp.Response[0].Message != nil && *cdnLocksReadResp.Response[0].Message != *cdnLock.Message {
		t.Errorf("expected Message %v, got %v", *cdnLock.Message, *cdnLocksReadResp.Response[0].Message)
	}
	if cdnLocksReadResp.Response[0].Soft == nil {
		t.Errorf("expected a valid soft/hard setting, but got nothing")
	}
	if cdnLocksReadResp.Response[0].Soft != nil && *cdnLocksReadResp.Response[0].Soft != *cdnLock.Soft {
		t.Errorf("expected 'Soft' to be %v, got %v", *cdnLock.Soft, *cdnLocksReadResp.Response[0].Soft)
	}

	// DELETE
	_, reqInf, err := TOSession.DeleteCDNLocks(client.RequestOptions{QueryParameters: url.Values{"cdn": []string{cdnLock.CDN}}})
	if err != nil {
		t.Fatalf("couldn't delete cdn lock, err: %v", err)
	}
	if reqInf.StatusCode != http.StatusOK {
		t.Errorf("expected status code of 200, but got %d instead", reqInf.StatusCode)
	}

}

func AdminCdnLocks(t *testing.T) {
	resp, _, err := TOSession.GetTenants(client.RequestOptions{})
	if err != nil {
		t.Fatalf("could not GET tenants: %v", err)
	}
	if len(resp.Response) == 0 {
		t.Fatalf("didn't get any tenant in response")
	}

	// Create a new user with operations level privileges
	user1 := tc.User{
		Username:             util.StrPtr("lock_user1"),
		RegistrationSent:     tc.TimeNoModFromTime(time.Now()),
		LocalPassword:        util.StrPtr("test_pa$$word"),
		ConfirmLocalPassword: util.StrPtr("test_pa$$word"),
		RoleName:             util.StrPtr("operations"),
	}
	user1.Email = util.StrPtr("lockuseremail@domain.com")
	user1.TenantID = util.IntPtr(resp.Response[0].ID)
	user1.FullName = util.StrPtr("firstName LastName")
	_, _, err = TOSession.CreateUser(user1, client.RequestOptions{})
	if err != nil {
		t.Fatalf("could not create test user with username: %s", *user1.Username)
	}
	defer ForceDeleteTestUsersByUsernames(t, []string{"lock_user1"})

	// Create another new user with operations level privileges
	user2 := tc.User{
		Username:             util.StrPtr("lock_user2"),
		RegistrationSent:     tc.TimeNoModFromTime(time.Now()),
		LocalPassword:        util.StrPtr("test_pa$$word2"),
		ConfirmLocalPassword: util.StrPtr("test_pa$$word2"),
		RoleName:             util.StrPtr("operations"),
	}
	user2.Email = util.StrPtr("newlockuseremail@domain.com")
	user2.TenantID = util.IntPtr(resp.Response[0].ID)
	user2.FullName = util.StrPtr("firstName2 LastName2")
	_, _, err = TOSession.CreateUser(user2, client.RequestOptions{})
	if err != nil {
		fmt.Println(err)
		t.Fatalf("could not create test user with username: %s", *user2.Username)
	}
	defer ForceDeleteTestUsersByUsernames(t, []string{"lock_user2"})

	// Establish a session with the newly created non admin level user
	userSession, _, err := client.LoginWithAgent(Config.TrafficOps.URL, *user1.Username, *user1.LocalPassword, true, "to-api-v4-client-tests", false, toReqTimeout)
	if err != nil {
		t.Fatalf("could not login with user lock_user1: %v", err)
	}

	// Establish another session with the newly created non admin level user
	userSession2, _, err := client.LoginWithAgent(Config.TrafficOps.URL, *user2.Username, *user2.LocalPassword, true, "to-api-v4-client-tests", false, toReqTimeout)
	if err != nil {
		t.Fatalf("could not login with user lock_user1: %v", err)
	}

	cdn := getCDNName(t)
	// Create a lock for this user
	_, _, err = userSession.CreateCDNLock(tc.CDNLock{
		CDN:     cdn,
		Message: util.StrPtr("test lock"),
		Soft:    util.BoolPtr(true),
	}, client.RequestOptions{})
	if err != nil {
		t.Fatalf("couldn't create cdn lock: %v", err)
	}

	// Non admin user trying to delete another user's lock -> this should fail
	_, reqInf, err := userSession2.DeleteCDNLocks(client.RequestOptions{QueryParameters: url.Values{"cdn": []string{cdn}}})
	if err == nil {
		t.Fatalf("expected error when a non admin user tries to delete another user's lock, but got nothing")
	}
	if reqInf.StatusCode != http.StatusForbidden {
		t.Fatalf("expected a 403 status code, but got %d instead", reqInf.StatusCode)
	}

	// Now try to delete another user's lock by hitting the admin DELETE endpoint for cdn_locks -> this should pass
	_, reqInf, err = TOSession.DeleteCDNLocks(client.RequestOptions{QueryParameters: url.Values{"cdn": []string{cdn}}})
	if err != nil {
		t.Fatalf("expected no error while deleting other user's lock using admin endpoint, but got %v", err)
	}
	if reqInf.StatusCode != http.StatusOK {
		t.Fatalf("expected a 200 status code, but got %d instead", reqInf.StatusCode)
	}
}

func SnapshotWithLock(t *testing.T) {
	resp, _, err := TOSession.GetTenants(client.RequestOptions{})
	if err != nil {
		t.Fatalf("could not GET tenants: %v", err)
	}
	if len(resp.Response) == 0 {
		t.Fatalf("didn't get any tenant in response")
	}

	// Create a new user with operations level privileges
	user1 := tc.User{
		Username:             util.StrPtr("lock_user1"),
		RegistrationSent:     tc.TimeNoModFromTime(time.Now()),
		LocalPassword:        util.StrPtr("test_pa$$word"),
		ConfirmLocalPassword: util.StrPtr("test_pa$$word"),
		RoleName:             util.StrPtr("operations"),
	}
	user1.Email = util.StrPtr("lockuseremail@domain.com")
	user1.TenantID = util.IntPtr(resp.Response[0].ID)
	user1.FullName = util.StrPtr("firstName LastName")
	_, _, err = TOSession.CreateUser(user1, client.RequestOptions{})
	if err != nil {
		t.Fatalf("could not create test user with username: %s", *user1.Username)
	}
	defer ForceDeleteTestUsersByUsernames(t, []string{"lock_user1"})

	// Establish a session with the newly created non admin level user
	userSession, _, err := client.LoginWithAgent(Config.TrafficOps.URL, *user1.Username, *user1.LocalPassword, true, "to-api-v4-client-tests", false, toReqTimeout)
	if err != nil {
		t.Fatalf("could not login with user lock_user1: %v", err)
	}

	cdn := getCDNName(t)

	// Currently, no user has a lock on the "bar" CDN, so when "lock_user1", which does not have the lock on CDN "bar", tries to snap it, this should pass
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("cdn", cdn)
	_, _, err = userSession.SnapshotCRConfig(opts)
	if err != nil {
		t.Errorf("expected no error while snapping cdn %s by user %s, but got %v", cdn, *user1.Username, err)
	}

	// Create a lock for this user
	_, _, err = userSession.CreateCDNLock(tc.CDNLock{
		CDN:     cdn,
		Message: util.StrPtr("test lock"),
		Soft:    util.BoolPtr(true),
	}, client.RequestOptions{})
	if err != nil {
		t.Fatalf("couldn't create cdn lock: %v", err)
	}

	// "lock_user1", which has the lock on CDN "bar", tries to snap it -> this should pass
	_, _, err = userSession.SnapshotCRConfig(opts)
	if err != nil {
		t.Errorf("expected no error while snapping cdn %s by user %s, but got %v", cdn, *user1.Username, err)
	}

	// Admin user, which doesn't have the lock on the CDN "bar", is trying to snap it -> this should fail
	_, reqInf, err := TOSession.SnapshotCRConfig(opts)
	if err == nil {
		t.Errorf("expected error while snapping cdn %s by user admin, but got nothing", cdn)
	}
	if reqInf.StatusCode != http.StatusForbidden {
		t.Fatalf("expected a 403 status code, but got %d instead", reqInf.StatusCode)
	}

	// Delete the lock
	_, _, err = userSession.DeleteCDNLocks(client.RequestOptions{QueryParameters: url.Values{"cdn": []string{cdn}}})
	if err != nil {
		t.Fatalf("expected no error while deleting other user's lock using admin endpoint, but got %v", err)
	}
}

func QueueUpdatesWithLock(t *testing.T) {
	resp, _, err := TOSession.GetTenants(client.RequestOptions{})
	if err != nil {
		t.Fatalf("could not GET tenants: %v", err)
	}
	if len(resp.Response) == 0 {
		t.Fatalf("didn't get any tenant in response")
	}

	// Create a new user with operations level privileges
	user1 := tc.User{
		Username:             util.StrPtr("lock_user1"),
		RegistrationSent:     tc.TimeNoModFromTime(time.Now()),
		LocalPassword:        util.StrPtr("test_pa$$word"),
		ConfirmLocalPassword: util.StrPtr("test_pa$$word"),
		RoleName:             util.StrPtr("operations"),
	}
	user1.Email = util.StrPtr("lockuseremail@domain.com")
	user1.TenantID = util.IntPtr(resp.Response[0].ID)
	user1.FullName = util.StrPtr("firstName LastName")
	_, _, err = TOSession.CreateUser(user1, client.RequestOptions{})
	if err != nil {
		t.Fatalf("could not create test user with username: %s", *user1.Username)
	}
	defer ForceDeleteTestUsersByUsernames(t, []string{"lock_user1"})

	// Establish a session with the newly created non admin level user
	userSession, _, err := client.LoginWithAgent(Config.TrafficOps.URL, *user1.Username, *user1.LocalPassword, true, "to-api-v4-client-tests", false, toReqTimeout)
	if err != nil {
		t.Fatalf("could not login with user lock_user1: %v", err)
	}

	cdn, serverID := getCDNNameAndServerID(t)
	if serverID == -1 {
		t.Fatalf("Could not get any valid servers to queue updates on")
	}

	// Currently, no user has a lock on the "bar" CDN, so when "lock_user1", which does not have the lock on CDN "bar", tries to queue updates on a server in the same CDN, this should pass
	_, _, err = userSession.SetServerQueueUpdate(serverID, true, client.RequestOptions{})
	if err != nil {
		t.Errorf("expected no error while queueing updates for a server in cdn %s by user %s, but got %v", cdn, *user1.Username, err)
	}

	// Create a lock for this user
	_, _, err = userSession.CreateCDNLock(tc.CDNLock{
		CDN:     cdn,
		Message: util.StrPtr("test lock"),
		Soft:    util.BoolPtr(true),
	}, client.RequestOptions{})
	if err != nil {
		t.Fatalf("couldn't create cdn lock: %v", err)
	}

	// "lock_user1", which has the lock on CDN "bar", tries to queue updates on a server in it -> this should pass
	_, _, err = userSession.SetServerQueueUpdate(serverID, true, client.RequestOptions{})
	if err != nil {
		t.Errorf("expected no error while queueing updates for a server in cdn %s by user %s, but got %v", cdn, *user1.Username, err)
	}

	// Admin user, which doesn't have the lock on the CDN "bar", is trying to queue updates on a server in it -> this should fail
	_, reqInf, err := TOSession.SetServerQueueUpdate(serverID, true, client.RequestOptions{})
	if err == nil {
		t.Errorf("expected error while queueing updates on a server in cdn %s by user admin, but got nothing", cdn)
	}
	if reqInf.StatusCode != http.StatusForbidden {
		t.Fatalf("expected a 403 status code, but got %d instead", reqInf.StatusCode)
	}

	// Delete the lock
	_, _, err = userSession.DeleteCDNLocks(client.RequestOptions{QueryParameters: url.Values{"cdn": []string{cdn}}})
	if err != nil {
		t.Fatalf("expected no error while deleting other user's lock using admin endpoint, but got %v", err)
	}
}

func QueueUpdatesFromTopologiesWithLock(t *testing.T) {
	resp, _, err := TOSession.GetTenants(client.RequestOptions{})
	if err != nil {
		t.Fatalf("could not GET tenants: %v", err)
	}
	if len(resp.Response) == 0 {
		t.Fatalf("didn't get any tenant in response")
	}

	// Create a new user with operations level privileges
	user1 := tc.User{
		Username:             util.StrPtr("lock_user1"),
		RegistrationSent:     tc.TimeNoModFromTime(time.Now()),
		LocalPassword:        util.StrPtr("test_pa$$word"),
		ConfirmLocalPassword: util.StrPtr("test_pa$$word"),
		RoleName:             util.StrPtr("operations"),
	}
	user1.Email = util.StrPtr("lockuseremail@domain.com")
	user1.TenantID = util.IntPtr(resp.Response[0].ID)
	user1.FullName = util.StrPtr("firstName LastName")
	_, _, err = TOSession.CreateUser(user1, client.RequestOptions{})
	if err != nil {
		t.Fatalf("could not create test user with username: %s", *user1.Username)
	}
	defer ForceDeleteTestUsersByUsernames(t, []string{"lock_user1"})

	// Establish a session with the newly created non admin level user
	userSession, _, err := client.LoginWithAgent(Config.TrafficOps.URL, *user1.Username, *user1.LocalPassword, true, "to-api-v4-client-tests", false, toReqTimeout)
	if err != nil {
		t.Fatalf("could not login with user lock_user1: %v", err)
	}

	cdnID, cdnName, topology := getCDNDetailsAndTopologyName(t)
	if topology == "" || cdnName == "" || cdnID == -1 {
		t.Fatalf("Could not get any valid topologies/ cdns to queue updates on")
	}

	// Currently, no user has a lock on the "bar" CDN, so when "lock_user1", which does not have the lock on CDN "bar", tries to queue updates on a topology in the same CDN, this should pass
	_, _, err = userSession.TopologiesQueueUpdate(topology, tc.TopologiesQueueUpdateRequest{Action: "queue", CDNID: int64(cdnID)}, client.RequestOptions{})
	if err != nil {
		t.Errorf("expected no error while queueing updates for a topology server in cdn %s by user %s, but got %v", cdnName, *user1.Username, err)
	}

	// Create a lock for this user
	_, _, err = userSession.CreateCDNLock(tc.CDNLock{
		CDN:     cdnName,
		Message: util.StrPtr("test lock"),
		Soft:    util.BoolPtr(true),
	}, client.RequestOptions{})
	if err != nil {
		t.Fatalf("couldn't create cdn lock: %v", err)
	}

	// "lock_user1", which has the lock on CDN "bar", tries to queue updates on a topology in it -> this should pass
	_, _, err = userSession.TopologiesQueueUpdate(topology, tc.TopologiesQueueUpdateRequest{Action: "queue", CDNID: int64(cdnID)}, client.RequestOptions{})
	if err != nil {
		t.Errorf("expected no error while queueing updates for a topology server in cdn %s by user %s, but got %v", cdnName, *user1.Username, err)
	}

	// Admin user, which doesn't have the lock on the CDN "bar", is trying to queue updates on a topology in it -> this should fail
	_, reqInf, err := TOSession.TopologiesQueueUpdate(topology, tc.TopologiesQueueUpdateRequest{Action: "queue", CDNID: int64(cdnID)}, client.RequestOptions{})
	if err == nil {
		t.Errorf("expected error while queueing updates on topology servers on cdn %s by user admin, but got nothing", cdnName)
	}
	if reqInf.StatusCode != http.StatusForbidden {
		t.Fatalf("expected a 403 status code, but got %d instead", reqInf.StatusCode)
	}

	// Delete the lock
	_, _, err = userSession.DeleteCDNLocks(client.RequestOptions{QueryParameters: url.Values{"cdn": []string{cdnName}}})
	if err != nil {
		t.Fatalf("expected no error while deleting lock, but got %v", err)
	}
}
