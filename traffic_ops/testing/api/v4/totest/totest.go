// package totest is a utility lib for applications using the Traffic Ops API to create integration tests.

package totest

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
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/assert"
	toclient "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

// TrafficControl - maps to the tc-fixtures.json file
type TrafficControl struct {
	ASNs                                              []tc.ASN                                `json:"asns"`
	CDNs                                              []tc.CDN                                `json:"cdns"`
	CDNLocks                                          []tc.CDNLock                            `json:"cdnlocks"`
	CacheGroups                                       []tc.CacheGroupNullable                 `json:"cachegroups"`
	Capabilities                                      []tc.Capability                         `json:"capability"`
	Coordinates                                       []tc.Coordinate                         `json:"coordinates"`
	DeliveryServicesRegexes                           []tc.DeliveryServiceRegexesTest         `json:"deliveryServicesRegexes"`
	DeliveryServiceRequests                           []tc.DeliveryServiceRequestV40          `json:"deliveryServiceRequests"`
	DeliveryServiceRequestComments                    []tc.DeliveryServiceRequestComment      `json:"deliveryServiceRequestComments"`
	DeliveryServices                                  []tc.DeliveryServiceV4                  `json:"deliveryservices"`
	DeliveryServicesRequiredCapabilities              []tc.DeliveryServicesRequiredCapability `json:"deliveryservicesRequiredCapabilities"`
	DeliveryServiceServerAssignments                  []tc.DeliveryServiceServers             `json:"deliveryServiceServerAssignments"`
	TopologyBasedDeliveryServicesRequiredCapabilities []tc.DeliveryServicesRequiredCapability `json:"topologyBasedDeliveryServicesRequiredCapabilities"`
	Divisions                                         []tc.Division                           `json:"divisions"`
	Federations                                       []tc.CDNFederation                      `json:"federations"`
	FederationResolvers                               []tc.FederationResolver                 `json:"federation_resolvers"`
	Jobs                                              []tc.InvalidationJobCreateV4            `json:"jobs"`
	Origins                                           []tc.Origin                             `json:"origins"`
	Profiles                                          []tc.Profile                            `json:"profiles"`
	Parameters                                        []tc.Parameter                          `json:"parameters"`
	ProfileParameters                                 []tc.ProfileParameter                   `json:"profileParameters"`
	PhysLocations                                     []tc.PhysLocation                       `json:"physLocations"`
	Regions                                           []tc.Region                             `json:"regions"`
	Roles                                             []tc.RoleV4                             `json:"roles"`
	Servers                                           []tc.ServerV40                          `json:"servers"`
	ServerServerCapabilities                          []tc.ServerServerCapability             `json:"serverServerCapabilities"`
	ServerCapabilities                                []tc.ServerCapability                   `json:"serverCapabilities"`
	ServiceCategories                                 []tc.ServiceCategory                    `json:"serviceCategories"`
	Statuses                                          []tc.StatusNullable                     `json:"statuses"`
	StaticDNSEntries                                  []tc.StaticDNSEntry                     `json:"staticdnsentries"`
	StatsSummaries                                    []tc.StatsSummary                       `json:"statsSummaries"`
	Tenants                                           []tc.Tenant                             `json:"tenants"`
	ServerCheckExtensions                             []tc.ServerCheckExtensionNullable       `json:"servercheck_extensions"`
	Topologies                                        []tc.Topology                           `json:"topologies"`
	Types                                             []tc.Type                               `json:"types"`
	SteeringTargets                                   []tc.SteeringTargetNullable             `json:"steeringTargets"`
	Serverchecks                                      []tc.ServercheckRequestNullable         `json:"serverchecks"`
	Users                                             []tc.UserV4                             `json:"users"`
	InvalidationJobs                                  []tc.InvalidationJobCreateV4            `json:"invalidationJobs"`
	InvalidationJobsRefetch                           []tc.InvalidationJobCreateV4            `json:"invalidationJobsRefetch"`
}

func GetCacheGroupId(t *testing.T, cl *toclient.Session, cacheGroupName string) func() int {
	return func() int {
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("name", cacheGroupName)

		resp, _, err := cl.GetCacheGroups(opts)
		assert.RequireNoError(t, err, "Get Cache Groups Request failed with error: %v", err)
		assert.RequireEqual(t, len(resp.Response), 1, "Expected response object length 1, but got %d", len(resp.Response))
		assert.RequireNotNil(t, resp.Response[0].ID, "Expected id to not be nil")

		return *resp.Response[0].ID
	}
}

func CreateTestASNs(t *testing.T, cl *toclient.Session, dat TrafficControl) {
	for _, asn := range dat.ASNs {
		asn.CachegroupID = GetCacheGroupId(t, cl, asn.Cachegroup)()
		resp, _, err := cl.CreateASN(asn, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Could not create ASN: %v - alerts: %+v", err, resp)
	}
}

func DeleteTestASNs(t *testing.T, cl *toclient.Session) {
	asns, _, err := cl.GetASNs(toclient.RequestOptions{})
	assert.NoError(t, err, "Error trying to fetch ASNs for deletion: %v - alerts: %+v", err, asns.Alerts)

	for _, asn := range asns.Response {
		alerts, _, err := cl.DeleteASN(asn.ID, toclient.RequestOptions{})
		assert.NoError(t, err, "Cannot delete ASN %d: %v - alerts: %+v", asn.ASN, err, alerts)
		// Retrieve the ASN to see if it got deleted
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("asn", strconv.Itoa(asn.ASN))
		asns, _, err := cl.GetASNs(opts)
		assert.NoError(t, err, "Error trying to fetch ASN after deletion: %v - alerts: %+v", err, asns.Alerts)
		assert.Equal(t, 0, len(asns.Response), "Expected ASN %d to be deleted, but it was found in Traffic Ops", asn.ASN)
	}
}

func CreateTestCacheGroups(t *testing.T, cl *toclient.Session, dat TrafficControl) {
	for _, cg := range dat.CacheGroups {

		resp, _, err := cl.CreateCacheGroup(cg, toclient.RequestOptions{})
		if err != nil {
			t.Errorf("could not create Cache Group: %v - alerts: %+v", err, resp.Alerts)
			continue
		}

		// Testing 'join' fields during create
		if cg.ParentName != nil && resp.Response.ParentName == nil {
			t.Error("Parent cachegroup is null in response when it should have a value")
		}
		if cg.SecondaryParentName != nil && resp.Response.SecondaryParentName == nil {
			t.Error("Secondary parent cachegroup is null in response when it should have a value")
		}
		if cg.Type != nil && resp.Response.Type == nil {
			t.Error("Type is null in response when it should have a value")
		}
		assert.NotNil(t, resp.Response.LocalizationMethods, "Localization methods are null")
		assert.NotNil(t, resp.Response.Fallbacks, "Fallbacks are null")
	}
}

func DeleteTestCacheGroups(t *testing.T, cl *toclient.Session, dat TrafficControl) {
	var parentlessCacheGroups []tc.CacheGroupNullable
	opts := toclient.NewRequestOptions()

	// delete the edge caches.
	for _, cg := range dat.CacheGroups {
		if cg.Name == nil {
			t.Error("Found a Cache Group with null or undefined name")
			continue
		}

		// Retrieve the CacheGroup by name so we can get the id for Deletion
		opts.QueryParameters.Set("name", *cg.Name)
		resp, _, err := cl.GetCacheGroups(opts)
		assert.NoError(t, err, "Cannot GET CacheGroup by name '%s': %v - alerts: %+v", *cg.Name, err, resp.Alerts)

		if len(resp.Response) < 1 {
			t.Errorf("Could not find test data Cache Group '%s' in Traffic Ops", *cg.Name)
			continue
		}
		cg = resp.Response[0]

		// Cachegroups that are parents (usually mids but sometimes edges)
		// need to be deleted only after the children cachegroups are deleted.
		if cg.ParentCachegroupID == nil && cg.SecondaryParentCachegroupID == nil {
			parentlessCacheGroups = append(parentlessCacheGroups, cg)
			continue
		}

		if cg.ID == nil {
			t.Error("Traffic Ops returned a Cache Group with null or undefined ID")
			continue
		}

		alerts, _, err := cl.DeleteCacheGroup(*cg.ID, toclient.RequestOptions{})
		assert.NoError(t, err, "Cannot delete Cache Group: %v - alerts: %+v", err, alerts)

		// Retrieve the CacheGroup to see if it got deleted
		opts.QueryParameters.Set("name", *cg.Name)
		cgs, _, err := cl.GetCacheGroups(opts)
		assert.NoError(t, err, "Error deleting Cache Group by name: %v - alerts: %+v", err, cgs.Alerts)
		assert.Equal(t, 0, len(cgs.Response), "Expected CacheGroup name: %s to be deleted", *cg.Name)
	}

	opts = toclient.NewRequestOptions()
	// now delete the parentless cachegroups
	for _, cg := range parentlessCacheGroups {
		// nil check for cg.Name occurs prior to insertion into parentlessCacheGroups
		opts.QueryParameters.Set("name", *cg.Name)
		// Retrieve the CacheGroup by name so we can get the id for Deletion
		resp, _, err := cl.GetCacheGroups(opts)
		assert.NoError(t, err, "Cannot get Cache Group by name '%s': %v - alerts: %+v", *cg.Name, err, resp.Alerts)

		if len(resp.Response) < 1 {
			t.Errorf("Cache Group '%s' somehow stopped existing since the last time we ask Traffic Ops about it", *cg.Name)
			continue
		}

		respCG := resp.Response[0]
		if respCG.ID == nil {
			t.Errorf("Traffic Ops returned Cache Group '%s' with null or undefined ID", *cg.Name)
			continue
		}
		delResp, _, err := cl.DeleteCacheGroup(*respCG.ID, toclient.RequestOptions{})
		assert.NoError(t, err, "Cannot delete Cache Group '%s': %v - alerts: %+v", *respCG.Name, err, delResp.Alerts)

		// Retrieve the CacheGroup to see if it got deleted
		opts.QueryParameters.Set("name", *cg.Name)
		cgs, _, err := cl.GetCacheGroups(opts)
		assert.NoError(t, err, "Error attempting to fetch Cache Group '%s' after deletion: %v - alerts: %+v", *cg.Name, err, cgs.Alerts)
		assert.Equal(t, 0, len(cgs.Response), "Expected Cache Group '%s' to be deleted", *cg.Name)
	}
}

func CreateTestCachegroupsDeliveryServices(t *testing.T, cl *toclient.Session) {
	dses, _, err := cl.GetDeliveryServices(toclient.RequestOptions{})
	assert.RequireNoError(t, err, "Cannot GET DeliveryServices: %v - %v", err, dses)

	opts := toclient.NewRequestOptions()
	opts.QueryParameters.Set("name", "cachegroup3")
	clientCGs, _, err := cl.GetCacheGroups(opts)
	assert.RequireNoError(t, err, "Cannot GET cachegroup: %v", err)
	assert.RequireEqual(t, len(clientCGs.Response), 1, "Getting cachegroup expected 1, got %v", len(clientCGs.Response))
	assert.RequireNotNil(t, clientCGs.Response[0].ID, "Cachegroup has a nil ID")

	dsIDs := []int{}
	for _, ds := range dses.Response {
		if *ds.CDNName == "cdn1" && ds.Topology == nil {
			dsIDs = append(dsIDs, *ds.ID)
		}
	}
	assert.RequireGreaterOrEqual(t, len(dsIDs), 1, "No Delivery Services found in CDN 'cdn1', cannot continue.")
	resp, _, err := cl.SetCacheGroupDeliveryServices(*clientCGs.Response[0].ID, dsIDs, toclient.RequestOptions{})
	assert.RequireNoError(t, err, "Setting cachegroup delivery services returned error: %v", err)
	assert.RequireGreaterOrEqual(t, len(resp.Response.ServerNames), 1, "Setting cachegroup delivery services returned success, but no servers set")
}

func setInactive(t *testing.T, cl *toclient.Session, dsID int) {
	opts := toclient.NewRequestOptions()
	opts.QueryParameters.Set("id", strconv.Itoa(dsID))
	resp, _, err := cl.GetDeliveryServices(opts)
	assert.RequireNoError(t, err, "Failed to fetch details for Delivery Service #%d: %v - alerts: %+v", dsID, err, resp.Alerts)
	assert.RequireEqual(t, len(resp.Response), 1, "Expected exactly one Delivery Service to exist with ID %d, found: %d", dsID, len(resp.Response))

	ds := resp.Response[0]
	if ds.Active == nil {
		t.Errorf("Deliver Service #%d had null or undefined 'active'", dsID)
		ds.Active = new(bool)
	}
	if *ds.Active {
		*ds.Active = false
		_, _, err = cl.UpdateDeliveryService(dsID, ds, toclient.RequestOptions{})
		assert.NoError(t, err, "Failed to set Delivery Service #%d to inactive: %v", dsID, err)
	}
}

func DeleteTestCachegroupsDeliveryServices(t *testing.T, cl *toclient.Session) {
	opts := toclient.NewRequestOptions()
	opts.QueryParameters.Set("limit", "1000000")
	dss, _, err := cl.GetDeliveryServiceServers(opts)
	assert.NoError(t, err, "Unexpected error retrieving server-to-Delivery-Service assignments: %v - alerts: %+v", err, dss.Alerts)

	for _, ds := range dss.Response {
		setInactive(t, cl, *ds.DeliveryService)
		alerts, _, err := cl.DeleteDeliveryServiceServer(*ds.DeliveryService, *ds.Server, toclient.RequestOptions{})
		assert.NoError(t, err, "Error deleting delivery service servers: %v - alerts: %+v", err, alerts.Alerts)
	}

	dss, _, err = cl.GetDeliveryServiceServers(toclient.RequestOptions{})
	assert.NoError(t, err, "Unexpected error retrieving server-to-Delivery-Service assignments: %v - alerts: %+v", err, dss.Alerts)
	assert.Equal(t, len(dss.Response), 0, "Deleting delivery service servers: Expected empty subsequent get, actual %v", len(dss.Response))
}

// TODO fix/remove globals

var fedIDs = make(map[string]int)

// All prerequisite Federations are associated to this cdn and this xmlID
const FederationCDNName = "cdn1"

var fedXmlId = "ds1"

func GetFederationID(t *testing.T, cname string) func() int {
	return func() int {
		ID, ok := fedIDs[cname]
		assert.RequireEqual(t, true, ok, "Expected to find Federation CName: %s to have associated ID", cname)
		return ID
	}
}

func setFederationID(t *testing.T, cdnFederation tc.CDNFederation) {
	assert.RequireNotNil(t, cdnFederation.CName, "Federation CName was nil after posting.")
	assert.RequireNotNil(t, cdnFederation.ID, "Federation ID was nil after posting.")
	fedIDs[*cdnFederation.CName] = *cdnFederation.ID
}

func CreateTestCDNFederations(t *testing.T, cl *toclient.Session, dat TrafficControl) {
	for _, federation := range dat.Federations {
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("xmlId", *federation.DeliveryServiceIDs.XmlId)
		dsResp, _, err := cl.GetDeliveryServices(opts)
		assert.RequireNoError(t, err, "Could not get Delivery Service by XML ID: %v", err)
		assert.RequireEqual(t, 1, len(dsResp.Response), "Expected one Delivery Service, but got %d", len(dsResp.Response))
		assert.RequireNotNil(t, dsResp.Response[0].CDNName, "Expected Delivery Service CDN Name to not be nil.")

		resp, _, err := cl.CreateCDNFederation(federation, *dsResp.Response[0].CDNName, toclient.RequestOptions{})
		assert.NoError(t, err, "Could not create CDN Federations: %v - alerts: %+v", err, resp.Alerts)

		// Need to save the ids, otherwise the other tests won't be able to reference the federations
		setFederationID(t, resp.Response)
		assert.RequireNotNil(t, resp.Response.ID, "Federation ID was nil after posting.")
		assert.RequireNotNil(t, dsResp.Response[0].ID, "Delivery Service ID was nil.")
		_, _, err = cl.CreateFederationDeliveryServices(*resp.Response.ID, []int{*dsResp.Response[0].ID}, false, toclient.NewRequestOptions())
		assert.NoError(t, err, "Could not create Federation Delivery Service: %v", err)
	}
}

func DeleteTestCDNFederations(t *testing.T, cl *toclient.Session) {
	opts := toclient.NewRequestOptions()
	for _, id := range fedIDs {
		resp, _, err := cl.DeleteCDNFederation(FederationCDNName, id, opts)
		assert.NoError(t, err, "Cannot delete federation #%d: %v - alerts: %+v", id, err, resp.Alerts)

		opts.QueryParameters.Set("id", strconv.Itoa(id))
		data, _, err := cl.GetCDNFederationsByName(FederationCDNName, opts)
		assert.Equal(t, 0, len(data.Response), "expected federation to be deleted")
	}
	fedIDs = make(map[string]int) // reset the global variable for the next test
}

func GetUserID(t *testing.T, cl *toclient.Session, username string) func() int {
	return func() int {
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("username", username)
		users, _, err := cl.GetUsers(opts)
		assert.RequireNoError(t, err, "Get Users Request failed with error:", err)
		assert.RequireEqual(t, 1, len(users.Response), "Expected response object length 1, but got %d", len(users.Response))
		assert.RequireNotNil(t, users.Response[0].ID, "Expected ID to not be nil.")
		return *users.Response[0].ID
	}
}

func CreateTestFederationUsers(t *testing.T, cl *toclient.Session) {
	// Prerequisite Federation Users
	federationUsers := map[string]tc.FederationUserPost{
		"the.cname.com.": {
			IDs:     []int{GetUserID(t, cl, "admin")(), GetUserID(t, cl, "adminuser")(), GetUserID(t, cl, "disalloweduser")(), GetUserID(t, cl, "readonlyuser")()},
			Replace: util.BoolPtr(false),
		},
		"booya.com.": {
			IDs:     []int{GetUserID(t, cl, "adminuser")()},
			Replace: util.BoolPtr(false),
		},
	}

	for cname, federationUser := range federationUsers {
		fedID := GetFederationID(t, cname)()
		resp, _, err := cl.CreateFederationUsers(fedID, federationUser.IDs, *federationUser.Replace, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Assigning users %v to federation %d: %v - alerts: %+v", federationUser.IDs, fedID, err, resp.Alerts)
	}
}

func DeleteTestFederationUsers(t *testing.T, cl *toclient.Session) {
	for _, fedID := range fedIDs {
		fedUsers, _, err := cl.GetFederationUsers(fedID, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Error getting users for federation %d: %v - alerts: %+v", fedID, err, fedUsers.Alerts)
		for _, fedUser := range fedUsers.Response {
			if fedUser.ID == nil {
				t.Error("Traffic Ops returned a representation of a relationship between a user and a Federation that had null or undefined ID")
				continue
			}
			alerts, _, err := cl.DeleteFederationUser(fedID, *fedUser.ID, toclient.RequestOptions{})
			assert.NoError(t, err, "Error deleting user #%d from federation #%d: %v - alerts: %+v", *fedUser.ID, fedID, err, alerts.Alerts)
		}
		fedUsers, _, err = cl.GetFederationUsers(fedID, toclient.RequestOptions{})
		assert.NoError(t, err, "Error getting users for federation %d: %v - alerts: %+v", fedID, err, fedUsers.Alerts)
		assert.Equal(t, 0, len(fedUsers.Response), "Federation users expected 0, actual: %+v", len(fedUsers.Response))
	}
}

func GetCDNID(t *testing.T, cl *toclient.Session, cdnName string) func() int {
	return func() int {
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("name", cdnName)
		cdnsResp, _, err := cl.GetCDNs(opts)
		assert.RequireNoError(t, err, "Get CDNs Request failed with error:", err)
		assert.RequireEqual(t, 1, len(cdnsResp.Response), "Expected response object length 1, but got %d", len(cdnsResp.Response))
		assert.RequireNotNil(t, cdnsResp.Response[0].ID, "Expected id to not be nil")
		return cdnsResp.Response[0].ID
	}
}

func CreateTestCDNs(t *testing.T, cl *toclient.Session, dat TrafficControl) {
	for _, cdn := range dat.CDNs {
		resp, _, err := cl.CreateCDN(cdn, toclient.RequestOptions{})
		assert.NoError(t, err, "Could not create CDN: %v - alerts: %+v", err, resp.Alerts)
	}
}

func DeleteTestCDNs(t *testing.T, cl *toclient.Session) {
	resp, _, err := cl.GetCDNs(toclient.RequestOptions{})
	assert.NoError(t, err, "Cannot get CDNs: %v - alerts: %+v", err, resp.Alerts)
	for _, cdn := range resp.Response {
		delResp, _, err := cl.DeleteCDN(cdn.ID, toclient.RequestOptions{})
		assert.NoError(t, err, "Cannot delete CDN '%s' (#%d): %v - alerts: %+v", cdn.Name, cdn.ID, err, delResp.Alerts)

		// Retrieve the CDN to see if it got deleted
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.Itoa(cdn.ID))
		cdns, _, err := cl.GetCDNs(opts)
		assert.NoError(t, err, "Error deleting CDN '%s': %v - alerts: %+v", cdn.Name, err, cdns.Alerts)
		assert.Equal(t, 0, len(cdns.Response), "Expected CDN '%s' to be deleted", cdn.Name)
	}
}

func CreateTestCoordinates(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, coordinate := range td.Coordinates {
		resp, _, err := cl.CreateCoordinate(coordinate, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Could not create coordinate: %v - alerts: %+v", err, resp.Alerts)
	}
}

func DeleteTestCoordinates(t *testing.T, cl *toclient.Session) {
	coordinates, _, err := cl.GetCoordinates(toclient.RequestOptions{})
	assert.NoError(t, err, "Cannot get Coordinates: %v - alerts: %+v", err, coordinates.Alerts)
	for _, coordinate := range coordinates.Response {
		alerts, _, err := cl.DeleteCoordinate(coordinate.ID, toclient.RequestOptions{})
		assert.NoError(t, err, "Unexpected error deleting Coordinate '%s' (#%d): %v - alerts: %+v", coordinate.Name, coordinate.ID, err, alerts.Alerts)
		// Retrieve the Coordinate to see if it got deleted
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.Itoa(coordinate.ID))
		getCoordinate, _, err := cl.GetCoordinates(opts)
		assert.NoError(t, err, "Error getting Coordinate '%s' after deletion: %v - alerts: %+v", coordinate.Name, err, getCoordinate.Alerts)
		assert.Equal(t, 0, len(getCoordinate.Response), "Expected Coordinate '%s' to be deleted, but it was found in Traffic Ops", coordinate.Name)
	}
}

func CreateTestDeliveryServiceRequestComments(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, comment := range td.DeliveryServiceRequestComments {
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("xmlId", comment.XMLID)
		resp, _, err := cl.GetDeliveryServiceRequests(opts)
		assert.NoError(t, err, "Cannot get Delivery Service Request by XMLID '%s': %v - alerts: %+v", comment.XMLID, err, resp.Alerts)
		assert.Equal(t, len(resp.Response), 1, "Found %d Delivery Service request by XMLID '%s, expected exactly one", len(resp.Response), comment.XMLID)
		assert.NotNil(t, resp.Response[0].ID, "Got Delivery Service Request with xml_id '%s' that had a null ID", comment.XMLID)

		comment.DeliveryServiceRequestID = *resp.Response[0].ID
		alerts, _, err := cl.CreateDeliveryServiceRequestComment(comment, toclient.RequestOptions{})
		assert.NoError(t, err, "Could not create Delivery Service Request Comment: %v - alerts: %+v", err, alerts.Alerts)
	}
}

func DeleteTestDeliveryServiceRequestComments(t *testing.T, cl *toclient.Session) {
	comments, _, err := cl.GetDeliveryServiceRequestComments(toclient.RequestOptions{})
	assert.NoError(t, err, "Unexpected error getting Delivery Service Request Comments: %v - alerts: %+v", err, comments.Alerts)

	for _, comment := range comments.Response {
		resp, _, err := cl.DeleteDeliveryServiceRequestComment(comment.ID, toclient.RequestOptions{})
		assert.NoError(t, err, "Cannot delete Delivery Service Request Comment #%d: %v - alerts: %+v", comment.ID, err, resp.Alerts)

		// Retrieve the delivery service request comment to see if it got deleted
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.Itoa(comment.ID))
		comments, _, err := cl.GetDeliveryServiceRequestComments(opts)
		assert.NoError(t, err, "Unexpected error fetching Delivery Service Request Comment %d after deletion: %v - alerts: %+v", comment.ID, err, comments.Alerts)
		assert.Equal(t, len(comments.Response), 0, "Expected Delivery Service Request Comment #%d to be deleted, but it was found in Traffic Ops", comment.ID)
	}
}

func CreateTestDeliveryServices(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, ds := range td.DeliveryServices {
		ds = ds.RemoveLD1AndLD2()
		if ds.XMLID == nil {
			t.Error("Found a Delivery Service in testing data with null or undefined XMLID")
			continue
		}
		resp, _, err := cl.CreateDeliveryService(ds, toclient.RequestOptions{})
		assert.NoError(t, err, "Could not create Delivery Service '%s': %v - alerts: %+v", *ds.XMLID, err, resp.Alerts)
	}
}

func DeleteTestDeliveryServices(t *testing.T, cl *toclient.Session) {
	dses, _, err := cl.GetDeliveryServices(toclient.RequestOptions{})
	assert.NoError(t, err, "Cannot get Delivery Services: %v - alerts: %+v", err, dses.Alerts)

	for _, ds := range dses.Response {
		delResp, _, err := cl.DeleteDeliveryService(*ds.ID, toclient.RequestOptions{})
		assert.NoError(t, err, "Could not delete Delivery Service: %v - alerts: %+v", err, delResp.Alerts)
		// Retrieve Delivery Service to see if it got deleted
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.Itoa(*ds.ID))
		getDS, _, err := cl.GetDeliveryServices(opts)
		assert.NoError(t, err, "Error deleting Delivery Service for '%s' : %v - alerts: %+v", *ds.XMLID, err, getDS.Alerts)
		assert.Equal(t, 0, len(getDS.Response), "Expected Delivery Service '%s' to be deleted", *ds.XMLID)
	}
}

func GetDeliveryServiceId(t *testing.T, cl *toclient.Session, xmlId string) func() int {
	return func() int {
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("xmlId", xmlId)

		resp, _, err := cl.GetDeliveryServices(opts)
		assert.RequireNoError(t, err, "Get Delivery Service Request failed with error: %v", err)
		assert.RequireEqual(t, 1, len(resp.Response), "Expected delivery service response object length 1, but got %d", len(resp.Response))
		assert.RequireNotNil(t, resp.Response[0].ID, "Expected id to not be nil")

		return *resp.Response[0].ID
	}
}

func CreateTestDeliveryServicesRequiredCapabilities(t *testing.T, cl *toclient.Session, td TrafficControl) {
	// Assign all required capability to delivery services listed in `tc-fixtures.json`.
	for _, dsrc := range td.DeliveryServicesRequiredCapabilities {
		dsId := GetDeliveryServiceId(t, cl, *dsrc.XMLID)()
		dsrc = tc.DeliveryServicesRequiredCapability{
			DeliveryServiceID:  &dsId,
			RequiredCapability: dsrc.RequiredCapability,
		}
		resp, _, err := cl.CreateDeliveryServicesRequiredCapability(dsrc, toclient.RequestOptions{})
		assert.NoError(t, err, "Unexpected error creating a Delivery Service/Required Capability relationship: %v - alerts: %+v", err, resp.Alerts)
	}
}

func DeleteTestDeliveryServicesRequiredCapabilities(t *testing.T, cl *toclient.Session) {
	// Get Required Capabilities to delete them
	dsrcs, _, err := cl.GetDeliveryServicesRequiredCapabilities(toclient.RequestOptions{})
	assert.NoError(t, err, "Error getting Delivery Service/Required Capability relationships: %v - alerts: %+v", err, dsrcs.Alerts)

	for _, dsrc := range dsrcs.Response {
		alerts, _, err := cl.DeleteDeliveryServicesRequiredCapability(*dsrc.DeliveryServiceID, *dsrc.RequiredCapability, toclient.RequestOptions{})
		assert.NoError(t, err, "Error deleting a relationship between a Delivery Service and a Capability: %v - alerts: %+v", err, alerts.Alerts)
	}
}

func GetTypeId(t *testing.T, cl *toclient.Session, typeName string) int {
	opts := toclient.NewRequestOptions()
	opts.QueryParameters.Set("name", typeName)
	resp, _, err := cl.GetTypes(opts)

	assert.RequireNoError(t, err, "Get Types Request failed with error: %v", err)
	assert.RequireEqual(t, 1, len(resp.Response), "Expected response object length 1, but got %d", len(resp.Response))
	assert.RequireNotNil(t, &resp.Response[0].ID, "Expected id to not be nil")

	return resp.Response[0].ID
}

func execSQL(db *sql.DB, sqlStmt string) error {
	var err error

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("transaction begin failed %v %v ", err, tx)
	}

	res, err := tx.Exec(sqlStmt)
	if err != nil {
		return fmt.Errorf("exec failed %v %v", err, res)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit failed %v %v", err, res)
	}
	return nil
}

func CreateTestDeliveryServicesRegexes(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, dsRegex := range td.DeliveryServicesRegexes {
		dsID := GetDeliveryServiceId(t, cl, dsRegex.DSName)()
		typeId := GetTypeId(t, cl, dsRegex.TypeName)
		dsRegexPost := tc.DeliveryServiceRegexPost{
			Type:      typeId,
			SetNumber: dsRegex.SetNumber,
			Pattern:   dsRegex.Pattern,
		}
		alerts, _, err := cl.PostDeliveryServiceRegexesByDSID(dsID, dsRegexPost, toclient.RequestOptions{})
		assert.NoError(t, err, "Could not create Delivery Service Regex: %v - alerts: %+v", err, alerts)
	}
}

func DeleteTestDeliveryServicesRegexes(t *testing.T, cl *toclient.Session, td TrafficControl, db *sql.DB) {
	for _, regex := range td.DeliveryServicesRegexes {
		err := execSQL(db, fmt.Sprintf("DELETE FROM deliveryservice_regex WHERE deliveryservice = '%v' and regex ='%v';", regex.DSID, regex.ID))
		assert.RequireNoError(t, err, "Unable to delete deliveryservice_regex by regex %v and ds %v: %v", regex.ID, regex.DSID, err)

		err = execSQL(db, fmt.Sprintf("DELETE FROM regex WHERE Id = '%v';", regex.ID))
		assert.RequireNoError(t, err, "Unable to delete regex %v: %v", regex.ID, err)
	}
}

func CreateTestDivisions(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, division := range td.Divisions {
		resp, _, err := cl.CreateDivision(division, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Could not create Division '%s': %v - alerts: %+v", division.Name, err, resp.Alerts)
	}
}

func DeleteTestDivisions(t *testing.T, cl *toclient.Session) {
	divisions, _, err := cl.GetDivisions(toclient.RequestOptions{})
	assert.NoError(t, err, "Cannot get Divisions: %v - alerts: %+v", err, divisions.Alerts)
	for _, division := range divisions.Response {
		alerts, _, err := cl.DeleteDivision(division.ID, toclient.RequestOptions{})
		assert.NoError(t, err, "Unexpected error deleting Division '%s' (#%d): %v - alerts: %+v", division.Name, division.ID, err, alerts.Alerts)
		// Retrieve the Division to see if it got deleted
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.Itoa(division.ID))
		getDivision, _, err := cl.GetDivisions(opts)
		assert.NoError(t, err, "Error getting Division '%s' after deletion: %v - alerts: %+v", division.Name, err, getDivision.Alerts)
		assert.Equal(t, 0, len(getDivision.Response), "Expected Division '%s' to be deleted, but it was found in Traffic Ops", division.Name)
	}
}

func CreateTestFederationResolvers(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, fr := range td.FederationResolvers {
		fr.TypeID = util.UIntPtr(uint(GetTypeId(t, cl, *fr.Type)))
		resp, _, err := cl.CreateFederationResolver(fr, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Failed to create Federation Resolver %+v: %v - alerts: %+v", fr, err, resp.Alerts)
	}
}

func DeleteTestFederationResolvers(t *testing.T, cl *toclient.Session) {
	frs, _, err := cl.GetFederationResolvers(toclient.RequestOptions{})
	assert.RequireNoError(t, err, "Unexpected error getting Federation Resolvers: %v - alerts: %+v", err, frs.Alerts)
	for _, fr := range frs.Response {
		alerts, _, err := cl.DeleteFederationResolver(*fr.ID, toclient.RequestOptions{})
		assert.NoError(t, err, "Failed to delete Federation Resolver %+v: %v - alerts: %+v", fr, err, alerts.Alerts)
		// Retrieve the Federation Resolver to see if it got deleted
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.Itoa(int(*fr.ID)))
		getFR, _, err := cl.GetFederationResolvers(opts)
		assert.NoError(t, err, "Error getting Federation Resolver '%d' after deletion: %v - alerts: %+v", *fr.ID, err, getFR.Alerts)
		assert.Equal(t, 0, len(getFR.Response), "Expected Federation Resolver '%d' to be deleted, but it was found in Traffic Ops", *fr.ID)
	}
}

func CreateTestJobs(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, job := range td.InvalidationJobs {
		job.StartTime = time.Now().Add(time.Minute).UTC()
		resp, _, err := cl.CreateInvalidationJob(job, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Could not create job: %v - alerts: %+v", err, resp.Alerts)
	}
}

func DeleteTestJobs(t *testing.T, cl *toclient.Session) {
	jobs, _, err := cl.GetInvalidationJobs(toclient.RequestOptions{})
	assert.NoError(t, err, "Cannot get Jobs: %v - alerts: %+v", err, jobs.Alerts)

	for _, job := range jobs.Response {
		alerts, _, err := cl.DeleteInvalidationJob(job.ID, toclient.RequestOptions{})
		assert.NoError(t, err, "Unexpected error deleting Job with ID: (#%d): %v - alerts: %+v", job.ID, err, alerts.Alerts)
		// Retrieve the Job to see if it got deleted
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.Itoa(int(job.ID)))
		getJobs, _, err := cl.GetInvalidationJobs(opts)
		assert.NoError(t, err, "Error getting Job with ID: '%d' after deletion: %v - alerts: %+v", job.ID, err, getJobs.Alerts)
		assert.Equal(t, 0, len(getJobs.Response), "Expected Job to be deleted, but it was found in Traffic Ops")
	}
}

func CreateTestOrigins(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, origin := range td.Origins {
		resp, _, err := cl.CreateOrigin(origin, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Could not create Origins: %v - alerts: %+v", err, resp.Alerts)
	}
}

func DeleteTestOrigins(t *testing.T, cl *toclient.Session) {
	origins, _, err := cl.GetOrigins(toclient.RequestOptions{})
	assert.NoError(t, err, "Cannot get Origins : %v - alerts: %+v", err, origins.Alerts)

	for _, origin := range origins.Response {
		assert.RequireNotNil(t, origin.ID, "Expected origin ID to not be nil.")
		assert.RequireNotNil(t, origin.Name, "Expected origin ID to not be nil.")
		assert.RequireNotNil(t, origin.IsPrimary, "Expected origin ID to not be nil.")
		if !*origin.IsPrimary {
			alerts, _, err := cl.DeleteOrigin(*origin.ID, toclient.RequestOptions{})
			assert.NoError(t, err, "Unexpected error deleting Origin '%s' (#%d): %v - alerts: %+v", *origin.Name, *origin.ID, err, alerts.Alerts)
			// Retrieve the Origin to see if it got deleted
			opts := toclient.NewRequestOptions()
			opts.QueryParameters.Set("id", strconv.Itoa(*origin.ID))
			getOrigin, _, err := cl.GetOrigins(opts)
			assert.NoError(t, err, "Error getting Origin '%s' after deletion: %v - alerts: %+v", *origin.Name, err, getOrigin.Alerts)
			assert.Equal(t, 0, len(getOrigin.Response), "Expected Origin '%s' to be deleted, but it was found in Traffic Ops", *origin.Name)
		}
	}
}

func CreateTestParameters(t *testing.T, cl *toclient.Session, td TrafficControl) {
	alerts, _, err := cl.CreateMultipleParameters(td.Parameters, toclient.RequestOptions{})
	assert.RequireNoError(t, err, "Could not create Parameters: %v - alerts: %+v", err, alerts)
}

func DeleteTestParameters(t *testing.T, cl *toclient.Session) {
	parameters, _, err := cl.GetParameters(toclient.RequestOptions{})
	assert.RequireNoError(t, err, "Cannot get Parameters: %v - alerts: %+v", err, parameters.Alerts)

	for _, parameter := range parameters.Response {
		alerts, _, err := cl.DeleteParameter(parameter.ID, toclient.RequestOptions{})
		assert.NoError(t, err, "Cannot delete Parameter #%d: %v - alerts: %+v", parameter.ID, err, alerts.Alerts)

		// Retrieve the Parameter to see if it got deleted
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.Itoa(parameter.ID))
		getParameters, _, err := cl.GetParameters(opts)
		assert.NoError(t, err, "Unexpected error fetching Parameter #%d after deletion: %v - alerts: %+v", parameter.ID, err, getParameters.Alerts)
		assert.Equal(t, 0, len(getParameters.Response), "Expected Parameter '%s' to be deleted, but it was found in Traffic Ops", parameter.Name)
	}
}

func CreateTestPhysLocations(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, pl := range td.PhysLocations {
		resp, _, err := cl.CreatePhysLocation(pl, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Could not create Physical Location '%s': %v - alerts: %+v", pl.Name, err, resp.Alerts)
	}
}

func DeleteTestPhysLocations(t *testing.T, cl *toclient.Session) {
	physicalLocations, _, err := cl.GetPhysLocations(toclient.RequestOptions{})
	assert.NoError(t, err, "Cannot get Physical Locations: %v - alerts: %+v", err, physicalLocations.Alerts)

	for _, pl := range physicalLocations.Response {
		alerts, _, err := cl.DeletePhysLocation(pl.ID, toclient.RequestOptions{})
		assert.NoError(t, err, "Unexpected error deleting Physical Location '%s' (#%d): %v - alerts: %+v", pl.Name, pl.ID, err, alerts.Alerts)
		// Retrieve the PhysLocation to see if it got deleted
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.Itoa(pl.ID))
		getPL, _, err := cl.GetPhysLocations(opts)
		assert.NoError(t, err, "Error getting Physical Location '%s' after deletion: %v - alerts: %+v", pl.Name, err, getPL.Alerts)
		assert.Equal(t, 0, len(getPL.Response), "Expected Physical Location '%s' to be deleted, but it was found in Traffic Ops", pl.Name)
	}
}

func GetProfileID(t *testing.T, cl *toclient.Session, profileName string) func() int {
	return func() int {
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("name", profileName)
		resp, _, err := cl.GetProfiles(opts)
		assert.RequireNoError(t, err, "Get Profiles Request failed with error: %v", err)
		assert.RequireEqual(t, 1, len(resp.Response), "Expected response object length 1, but got %d", len(resp.Response))
		return resp.Response[0].ID
	}
}

func CreateTestProfileParameters(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, profile := range td.Profiles {
		profileID := GetProfileID(t, cl, profile.Name)()

		for _, parameter := range profile.Parameters {
			assert.RequireNotNil(t, parameter.Name, "Expected parameter name to not be nil.")
			assert.RequireNotNil(t, parameter.Value, "Expected parameter value to not be nil.")
			assert.RequireNotNil(t, parameter.ConfigFile, "Expected parameter configFile to not be nil.")

			parameterOpts := toclient.NewRequestOptions()
			parameterOpts.QueryParameters.Set("name", *parameter.Name)
			parameterOpts.QueryParameters.Set("configFile", *parameter.ConfigFile)
			parameterOpts.QueryParameters.Set("value", *parameter.Value)
			getParameter, _, err := cl.GetParameters(parameterOpts)
			assert.RequireNoError(t, err, "Could not get Parameter %s: %v - alerts: %+v", *parameter.Name, err, getParameter.Alerts)
			if len(getParameter.Response) == 0 {
				alerts, _, err := cl.CreateParameter(tc.Parameter{Name: *parameter.Name, Value: *parameter.Value, ConfigFile: *parameter.ConfigFile}, toclient.RequestOptions{})
				assert.RequireNoError(t, err, "Could not create Parameter %s: %v - alerts: %+v", parameter.Name, err, alerts.Alerts)
				getParameter, _, err = cl.GetParameters(parameterOpts)
				assert.RequireNoError(t, err, "Could not get Parameter %s: %v - alerts: %+v", *parameter.Name, err, getParameter.Alerts)
				assert.RequireNotEqual(t, 0, len(getParameter.Response), "Could not get parameter %s: not found", *parameter.Name)
			}
			profileParameter := tc.ProfileParameterCreationRequest{ProfileID: profileID, ParameterID: getParameter.Response[0].ID}
			alerts, _, err := cl.CreateProfileParameter(profileParameter, toclient.RequestOptions{})
			assert.NoError(t, err, "Could not associate Parameter %s with Profile %s: %v - alerts: %+v", parameter.Name, profile.Name, err, alerts.Alerts)
		}
	}
}

func DeleteTestProfileParameters(t *testing.T, cl *toclient.Session) {
	profileParameters, _, err := cl.GetProfileParameters(toclient.RequestOptions{})
	assert.NoError(t, err, "Cannot get Profile Parameters: %v - alerts: %+v", err, profileParameters.Alerts)

	for _, profileParameter := range profileParameters.Response {
		alerts, _, err := cl.DeleteProfileParameter(GetProfileID(t, cl, profileParameter.Profile)(), profileParameter.Parameter, toclient.RequestOptions{})
		assert.NoError(t, err, "Unexpected error deleting Profile Parameter: Profile: '%s' Parameter ID: (#%d): %v - alerts: %+v", profileParameter.Profile, profileParameter.Parameter, err, alerts.Alerts)
	}
	// Retrieve the Profile Parameters to see if it got deleted
	getProfileParameter, _, err := cl.GetProfileParameters(toclient.RequestOptions{})
	assert.NoError(t, err, "Error getting Profile Parameters after deletion: %v - alerts: %+v", err, getProfileParameter.Alerts)
	assert.Equal(t, 0, len(getProfileParameter.Response), "Expected Profile Parameters to be deleted, but %d were found in Traffic Ops", len(getProfileParameter.Response))
}

func CreateTestProfiles(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, profile := range td.Profiles {
		resp, _, err := cl.CreateProfile(profile, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Could not create Profile '%s': %v - alerts: %+v", profile.Name, err, resp.Alerts)
	}
}

func DeleteTestProfiles(t *testing.T, cl *toclient.Session) {
	profiles, _, err := cl.GetProfiles(toclient.RequestOptions{})
	assert.NoError(t, err, "Cannot get Profiles: %v - alerts: %+v", err, profiles.Alerts)
	for _, profile := range profiles.Response {
		alerts, _, err := cl.DeleteProfile(profile.ID, toclient.RequestOptions{})
		assert.NoError(t, err, "Cannot delete Profile: %v - alerts: %+v", err, alerts.Alerts)
		// Retrieve the Profile to see if it got deleted
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.Itoa(profile.ID))
		getProfiles, _, err := cl.GetProfiles(opts)
		assert.NoError(t, err, "Error getting Profile '%s' after deletion: %v - alerts: %+v", profile.Name, err, getProfiles.Alerts)
		assert.Equal(t, 0, len(getProfiles.Response), "Expected Profile '%s' to be deleted, but it was found in Traffic Ops", profile.Name)
	}
}

func CreateTestRegions(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, region := range td.Regions {
		resp, _, err := cl.CreateRegion(region, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Could not create Region '%s': %v - alerts: %+v", region.Name, err, resp.Alerts)
	}
}

func DeleteTestRegions(t *testing.T, cl *toclient.Session) {
	regions, _, err := cl.GetRegions(toclient.RequestOptions{})
	assert.NoError(t, err, "Cannot get Regions: %v - alerts: %+v", err, regions.Alerts)

	for _, region := range regions.Response {
		alerts, _, err := cl.DeleteRegion(region.Name, toclient.RequestOptions{})
		assert.NoError(t, err, "Unexpected error deleting Region '%s' (#%d): %v - alerts: %+v", region.Name, region.ID, err, alerts.Alerts)
		// Retrieve the Region to see if it got deleted
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.Itoa(region.ID))
		getRegion, _, err := cl.GetRegions(opts)
		assert.NoError(t, err, "Error getting Region '%s' after deletion: %v - alerts: %+v", region.Name, err, getRegion.Alerts)
		assert.Equal(t, 0, len(getRegion.Response), "Expected Region '%s' to be deleted, but it was found in Traffic Ops", region.Name)
	}
}

func CreateTestRoles(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, role := range td.Roles {
		_, _, err := cl.CreateRole(role, toclient.RequestOptions{})
		assert.NoError(t, err, "No error expected, but got %v", err)
	}
}

func DeleteTestRoles(t *testing.T, cl *toclient.Session) {
	roles, _, err := cl.GetRoles(toclient.RequestOptions{})
	assert.NoError(t, err, "Cannot get Roles: %v - alerts: %+v", err, roles.Alerts)
	for _, role := range roles.Response {
		// Don't delete active roles created by test setup
		if role.Name == "admin" || role.Name == "disallowed" || role.Name == "operations" || role.Name == "portal" || role.Name == "read-only" || role.Name == "steering" || role.Name == "federation" {
			continue
		}
		_, _, err := cl.DeleteRole(role.Name, toclient.NewRequestOptions())
		assert.NoError(t, err, "Expected no error while deleting role %s, but got %v", role.Name, err)
		// Retrieve the Role to see if it got deleted
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("name", role.Name)
		getRole, _, err := cl.GetRoles(opts)
		assert.NoError(t, err, "Error getting Role '%s' after deletion: %v - alerts: %+v", role.Name, err, getRole.Alerts)
		assert.Equal(t, 0, len(getRole.Response), "Expected Role '%s' to be deleted, but it was found in Traffic Ops", role.Name)
	}
}

func CreateTestServerCapabilities(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, sc := range td.ServerCapabilities {
		resp, _, err := cl.CreateServerCapability(sc, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Unexpected error creating Server Capability '%s': %v - alerts: %+v", sc.Name, err, resp.Alerts)
	}
}

func DeleteTestServerCapabilities(t *testing.T, cl *toclient.Session) {
	serverCapabilities, _, err := cl.GetServerCapabilities(toclient.RequestOptions{})
	assert.NoError(t, err, "Cannot get Server Capabilities: %v - alerts: %+v", err, serverCapabilities.Alerts)

	for _, serverCapability := range serverCapabilities.Response {
		alerts, _, err := cl.DeleteServerCapability(serverCapability.Name, toclient.RequestOptions{})
		assert.NoError(t, err, "Unexpected error deleting Server Capability '%s': %v - alerts: %+v", serverCapability.Name, err, alerts.Alerts)
		// Retrieve the Server Capability to see if it got deleted
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("name", serverCapability.Name)
		getServerCapability, _, err := cl.GetServerCapabilities(opts)
		assert.NoError(t, err, "Error getting Server Capability '%s' after deletion: %v - alerts: %+v", serverCapability.Name, err, getServerCapability.Alerts)
		assert.Equal(t, 0, len(getServerCapability.Response), "Expected Server Capability '%s' to be deleted, but it was found in Traffic Ops", serverCapability.Name)
	}
}

func CreateTestServerChecks(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, servercheck := range td.Serverchecks {
		resp, _, err := cl.InsertServerCheckStatus(servercheck, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Could not insert Servercheck: %v - alerts: %+v", err, resp.Alerts)
	}
}

// Need to define no-op function as TCObj interface expects a delete function
// There is no delete path for serverchecks
func DeleteTestServerChecks(*testing.T, *toclient.Session) {
	return
}

func GetServerID(t *testing.T, cl *toclient.Session, hostName string) func() int {
	return func() int {
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("hostName", hostName)
		serversResp, _, err := cl.GetServers(opts)
		assert.RequireNoError(t, err, "Get Servers Request failed with error:", err)
		assert.RequireEqual(t, 1, len(serversResp.Response), "Expected response object length 1, but got %d", len(serversResp.Response))
		assert.RequireNotNil(t, serversResp.Response[0].ID, "Expected id to not be nil")
		return *serversResp.Response[0].ID
	}
}

func CreateTestServerServerCapabilities(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, ssc := range td.ServerServerCapabilities {
		assert.RequireNotNil(t, ssc.Server, "Expected Server to not be nil.")
		assert.RequireNotNil(t, ssc.ServerCapability, "Expected Server Capability to not be nil.")
		serverID := GetServerID(t, cl, *ssc.Server)()
		ssc.ServerID = &serverID
		resp, _, err := cl.CreateServerServerCapability(ssc, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Could not associate Capability '%s' with server '%s': %v - alerts: %+v", *ssc.ServerCapability, *ssc.Server, err, resp.Alerts)
	}
}

func DeleteTestServerServerCapabilities(t *testing.T, cl *toclient.Session) {
	sscs, _, err := cl.GetServerServerCapabilities(toclient.RequestOptions{})
	assert.RequireNoError(t, err, "Cannot get server server capabilities: %v - alerts: %+v", err, sscs.Alerts)
	for _, ssc := range sscs.Response {
		assert.RequireNotNil(t, ssc.Server, "Expected Server to not be nil.")
		assert.RequireNotNil(t, ssc.ServerCapability, "Expected Server Capability to not be nil.")
		alerts, _, err := cl.DeleteServerServerCapability(*ssc.ServerID, *ssc.ServerCapability, toclient.RequestOptions{})
		assert.NoError(t, err, "Could not remove Capability '%s' from server '%s' (#%d): %v - alerts: %+v", *ssc.ServerCapability, *ssc.Server, *ssc.ServerID, err, alerts.Alerts)
	}
}

func CreateTestServers(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, server := range td.Servers {
		resp, _, err := cl.CreateServer(server, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Could not create server '%s': %v - alerts: %+v", *server.HostName, err, resp.Alerts)
	}
}

func DeleteTestServers(t *testing.T, cl *toclient.Session) {
	servers, _, err := cl.GetServers(toclient.RequestOptions{})
	assert.NoError(t, err, "Cannot get Servers: %v - alerts: %+v", err, servers.Alerts)

	for _, server := range servers.Response {
		delResp, _, err := cl.DeleteServer(*server.ID, toclient.RequestOptions{})
		assert.NoError(t, err, "Could not delete Server: %v - alerts: %+v", err, delResp.Alerts)
		// Retrieve Server to see if it got deleted
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.Itoa(*server.ID))
		getServer, _, err := cl.GetServers(opts)
		assert.RequireNotNil(t, server.HostName, "Expected server host name to not be nil.")
		assert.NoError(t, err, "Error deleting Server for '%s' : %v - alerts: %+v", *server.HostName, err, getServer.Alerts)
		assert.Equal(t, 0, len(getServer.Response), "Expected Server '%s' to be deleted", *server.HostName)
	}
}

func CreateTestServiceCategories(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, serviceCategory := range td.ServiceCategories {
		resp, _, err := cl.CreateServiceCategory(serviceCategory, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Could not create Service Category: %v - alerts: %+v", err, resp.Alerts)
	}
}

func DeleteTestServiceCategories(t *testing.T, cl *toclient.Session) {
	serviceCategories, _, err := cl.GetServiceCategories(toclient.RequestOptions{})
	assert.NoError(t, err, "Cannot get Service Categories: %v - alerts: %+v", err, serviceCategories.Alerts)

	for _, serviceCategory := range serviceCategories.Response {
		alerts, _, err := cl.DeleteServiceCategory(serviceCategory.Name, toclient.RequestOptions{})
		assert.NoError(t, err, "Unexpected error deleting Service Category '%s': %v - alerts: %+v", serviceCategory.Name, err, alerts.Alerts)
		// Retrieve the Service Category to see if it got deleted
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("name", serviceCategory.Name)
		getServiceCategory, _, err := cl.GetServiceCategories(opts)
		assert.NoError(t, err, "Error getting Service Category '%s' after deletion: %v - alerts: %+v", serviceCategory.Name, err, getServiceCategory.Alerts)
		assert.Equal(t, 0, len(getServiceCategory.Response), "Expected Service Category '%s' to be deleted, but it was found in Traffic Ops", serviceCategory.Name)
	}
}

func CreateTestStatuses(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, status := range td.Statuses {
		resp, _, err := cl.CreateStatus(status, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Could not create Status: %v - alerts: %+v", err, resp.Alerts)
	}
}

func DeleteTestStatuses(t *testing.T, cl *toclient.Session, td TrafficControl) {
	opts := toclient.NewRequestOptions()
	for _, status := range td.Statuses {
		assert.RequireNotNil(t, status.Name, "Cannot get test statuses: test data statuses must have names")
		// Retrieve the Status by name, so we can get the id for the Update
		opts.QueryParameters.Set("name", *status.Name)
		resp, _, err := cl.GetStatuses(opts)
		assert.RequireNoError(t, err, "Cannot get Statuses filtered by name '%s': %v - alerts: %+v", *status.Name, err, resp.Alerts)
		assert.RequireEqual(t, 1, len(resp.Response), "Expected 1 status returned. Got: %d", len(resp.Response))
		respStatus := resp.Response[0]

		delResp, _, err := cl.DeleteStatus(respStatus.ID, toclient.RequestOptions{})
		assert.NoError(t, err, "Cannot delete Status: %v - alerts: %+v", err, delResp.Alerts)

		// Retrieve the Status to see if it got deleted
		resp, _, err = cl.GetStatuses(opts)
		assert.NoError(t, err, "Unexpected error getting Statuses filtered by name after deletion: %v - alerts: %+v", err, resp.Alerts)
		assert.Equal(t, 0, len(resp.Response), "Expected Status '%s' to be deleted, but it was found in Traffic Ops", *status.Name)
	}
}

func CreateTestStaticDNSEntries(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, staticDNSEntry := range td.StaticDNSEntries {
		resp, _, err := cl.CreateStaticDNSEntry(staticDNSEntry, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Could not create Static DNS Entry: %v - alerts: %+v", err, resp.Alerts)
	}
}

func DeleteTestStaticDNSEntries(t *testing.T, cl *toclient.Session) {
	staticDNSEntries, _, err := cl.GetStaticDNSEntries(toclient.RequestOptions{})
	assert.NoError(t, err, "Cannot get Static DNS Entries: %v - alerts: %+v", err, staticDNSEntries.Alerts)

	for _, staticDNSEntry := range staticDNSEntries.Response {
		alerts, _, err := cl.DeleteStaticDNSEntry(staticDNSEntry.ID, toclient.RequestOptions{})
		assert.NoError(t, err, "Unexpected error deleting Static DNS Entry '%s' (#%d): %v - alerts: %+v", staticDNSEntry.Host, staticDNSEntry.ID, err, alerts.Alerts)
		// Retrieve the Static DNS Entry to see if it got deleted
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("host", staticDNSEntry.Host)
		getStaticDNSEntry, _, err := cl.GetStaticDNSEntries(opts)
		assert.NoError(t, err, "Error getting Static DNS Entry '%s' after deletion: %v - alerts: %+v", staticDNSEntry.Host, err, getStaticDNSEntry.Alerts)
		assert.Equal(t, 0, len(getStaticDNSEntry.Response), "Expected Static DNS Entry '%s' to be deleted, but it was found in Traffic Ops", staticDNSEntry.Host)
	}
}

func CreateTestSteeringTargets(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, st := range td.SteeringTargets {
		st.TypeID = util.IntPtr(GetTypeId(t, cl, *st.Type))
		st.DeliveryServiceID = util.UInt64Ptr(uint64(GetDeliveryServiceId(t, cl, string(*st.DeliveryService))()))
		st.TargetID = util.UInt64Ptr(uint64(GetDeliveryServiceId(t, cl, string(*st.Target))()))
		resp, _, err := cl.CreateSteeringTarget(st, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Creating steering target: %v - alerts: %+v", err, resp.Alerts)
	}
}

func DeleteTestSteeringTargets(t *testing.T, cl *toclient.Session, td TrafficControl) {
	dsIDs := []uint64{}
	for _, st := range td.SteeringTargets {
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("xmlId", string(*st.DeliveryService))
		respDS, _, err := cl.GetDeliveryServices(opts)
		assert.RequireNoError(t, err, "Deleting steering target: getting ds: %v - alerts: %+v", err, respDS.Alerts)
		assert.RequireEqual(t, 1, len(respDS.Response), "Deleting steering target: getting ds: expected 1 delivery service")
		assert.RequireNotNil(t, respDS.Response[0].ID, "Deleting steering target: getting ds: nil ID returned")

		dsID := uint64(*respDS.Response[0].ID)
		st.DeliveryServiceID = &dsID
		dsIDs = append(dsIDs, dsID)

		opts.QueryParameters.Set("xmlId", string(*st.Target))
		respTarget, _, err := cl.GetDeliveryServices(opts)
		assert.RequireNoError(t, err, "Deleting steering target: getting target ds: %v - alerts: %+v", err, respTarget.Alerts)
		assert.RequireEqual(t, 1, len(respTarget.Response), "Deleting steering target: getting target ds: expected 1 delivery service")
		assert.RequireNotNil(t, respTarget.Response[0].ID, "Deleting steering target: getting target ds: not found")

		targetID := uint64(*respTarget.Response[0].ID)
		st.TargetID = &targetID

		resp, _, err := cl.DeleteSteeringTarget(int(*st.DeliveryServiceID), int(*st.TargetID), toclient.RequestOptions{})
		assert.NoError(t, err, "Deleting steering target: deleting: %v - alerts: %+v", err, resp.Alerts)
	}

	for _, dsID := range dsIDs {
		sts, _, err := cl.GetSteeringTargets(int(dsID), toclient.RequestOptions{})
		assert.NoError(t, err, "deleting steering targets: getting steering target: %v - alerts: %+v", err, sts.Alerts)
		assert.Equal(t, 0, len(sts.Response), "Deleting steering targets: after delete, getting steering target: expected 0 actual %d", len(sts.Response))
	}
}

func CreateTestTenants(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, tenant := range td.Tenants {
		resp, _, err := cl.CreateTenant(tenant, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Could not create Tenant '%s': %v - alerts: %+v", tenant.Name, err, resp.Alerts)
	}
}

func DeleteTestTenants(t *testing.T, cl *toclient.Session) {
	opts := toclient.NewRequestOptions()
	opts.QueryParameters.Set("sortOrder", "desc")
	tenants, _, err := cl.GetTenants(opts)
	assert.NoError(t, err, "Cannot get Tenants: %v - alerts: %+v", err, tenants.Alerts)

	for _, tenant := range tenants.Response {
		if tenant.Name == "root" {
			continue
		}
		alerts, _, err := cl.DeleteTenant(tenant.ID, toclient.RequestOptions{})
		assert.NoError(t, err, "Unexpected error deleting Tenant '%s' (#%d): %v - alerts: %+v", tenant.Name, tenant.ID, err, alerts.Alerts)
		// Retrieve the Tenant to see if it got deleted
		opts.QueryParameters.Set("id", strconv.Itoa(tenant.ID))
		getTenants, _, err := cl.GetTenants(opts)
		assert.NoError(t, err, "Error getting Tenant '%s' after deletion: %v - alerts: %+v", tenant.Name, err, getTenants.Alerts)
		assert.Equal(t, 0, len(getTenants.Response), "Expected Tenant '%s' to be deleted, but it was found in Traffic Ops", tenant.Name)
	}
}

func CreateTestServerCheckExtensions(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, ext := range td.ServerCheckExtensions {
		resp, _, err := cl.CreateServerCheckExtension(ext, toclient.RequestOptions{})
		assert.NoError(t, err, "Could not create Servercheck Extension: %v - alerts: %+v", err, resp.Alerts)
	}
}

func DeleteTestServerCheckExtensions(t *testing.T, cl *toclient.Session, td TrafficControl) {
	extensions, _, err := cl.GetServerCheckExtensions(toclient.RequestOptions{})
	assert.RequireNoError(t, err, "Could not get Servercheck Extensions: %v - alerts: %+v", err, extensions.Alerts)

	for _, extension := range extensions.Response {
		alerts, _, err := cl.DeleteServerCheckExtension(*extension.ID, toclient.RequestOptions{})
		assert.NoError(t, err, "Unexpected error deleting Servercheck Extension '%s' (#%d): %v - alerts: %+v", *extension.Name, *extension.ID, err, alerts.Alerts)
		// Retrieve the Server Extension to see if it got deleted
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.Itoa(*extension.ID))
		getExtension, _, err := cl.GetServerCheckExtensions(opts)
		assert.NoError(t, err, "Error getting Servercheck Extension '%s' after deletion: %v - alerts: %+v", *extension.Name, err, getExtension.Alerts)
		assert.Equal(t, 0, len(getExtension.Response), "Expected Servercheck Extension '%s' to be deleted, but it was found in Traffic Ops", *extension.Name)
	}
}

func CreateTestTopologies(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, topology := range td.Topologies {
		resp, _, err := cl.CreateTopology(topology, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Could not create Topology: %v - alerts: %+v", err, resp.Alerts)
	}
}

func DeleteTestTopologies(t *testing.T, cl *toclient.Session) {
	topologies, _, err := cl.GetTopologies(toclient.RequestOptions{})
	assert.NoError(t, err, "Cannot get Topologies: %v - alerts: %+v", err, topologies.Alerts)

	for _, topology := range topologies.Response {
		alerts, _, err := cl.DeleteTopology(topology.Name, toclient.RequestOptions{})
		assert.NoError(t, err, "Cannot delete Topology: %v - alerts: %+v", err, alerts.Alerts)
		// Retrieve the Topology to see if it got deleted
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("name", topology.Name)
		resp, _, err := cl.GetTopologies(opts)
		assert.NoError(t, err, "Unexpected error trying to fetch Topologies after deletion: %v - alerts: %+v", err, resp.Alerts)
		assert.Equal(t, 0, len(resp.Response), "Expected Topology '%s' to be deleted, but it was found in Traffic Ops", topology.Name)
	}
}

func CreateTestTypes(t *testing.T, cl *toclient.Session, td TrafficControl, db *sql.DB) {
	defer func() {
		err := db.Close()
		assert.NoError(t, err, "unable to close connection to db, error: %v", err)
	}()
	dbQueryTemplate := "INSERT INTO type (name, description, use_in_table) VALUES ('%s', '%s', '%s');"

	for _, typ := range td.Types {
		if typ.UseInTable != "server" {
			err := execSQL(db, fmt.Sprintf(dbQueryTemplate, typ.Name, typ.Description, typ.UseInTable))
			assert.RequireNoError(t, err, "could not create Type using database operations: %v", err)
		} else {
			alerts, _, err := cl.CreateType(typ, toclient.RequestOptions{})
			assert.RequireNoError(t, err, "could not create Type: %v - alerts: %+v", err, alerts.Alerts)
		}
	}
}

func DeleteTestTypes(t *testing.T, cl *toclient.Session, td TrafficControl, db *sql.DB) {
	dbDeleteTemplate := "DELETE FROM type WHERE name='%s';"

	types, _, err := cl.GetTypes(toclient.RequestOptions{})
	assert.NoError(t, err, "Cannot get Types: %v - alerts: %+v", err, types.Alerts)

	for _, typ := range types.Response {
		if typ.Name == "CHECK_EXTENSION_BOOL" || typ.Name == "CHECK_EXTENSION_NUM" || typ.Name == "CHECK_EXTENSION_OPEN_SLOT" {
			continue
		}

		if typ.UseInTable != "server" {
			err := execSQL(db, fmt.Sprintf(dbDeleteTemplate, typ.Name))
			assert.RequireNoError(t, err, "cannot delete Type using database operations: %v", err)
		} else {
			delResp, _, err := cl.DeleteType(typ.ID, toclient.RequestOptions{})
			assert.RequireNoError(t, err, "cannot delete Type using the API: %v - alerts: %+v", err, delResp.Alerts)
		}

		// Retrieve the Type by name to see if it was deleted.
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("name", typ.Name)
		types, _, err := cl.GetTypes(opts)
		assert.NoError(t, err, "error fetching Types filtered by presumably deleted name: %v - alerts: %+v", err, types.Alerts)
		assert.Equal(t, 0, len(types.Response), "expected Type '%s' to be deleted", typ.Name)
	}
}

func CreateTestUsers(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, user := range td.Users {
		resp, _, err := cl.CreateUser(user, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Could not create user: %v - alerts: %+v", err, resp.Alerts)
	}
}

// ForceDeleteTestUsers forcibly deletes the users from the db.
// NOTE: Special circumstances!  This should *NOT* be done without a really good reason!
// Connects directly to the DB to remove users rather than going through the client.
// This is required here because the DeleteUser action does not really delete users,  but disables them.
func ForceDeleteTestUsers(t *testing.T, cl *toclient.Session, td TrafficControl, db *sql.DB) {
	var usernames []string
	for _, user := range td.Users {
		usernames = append(usernames, `'`+user.Username+`'`)
	}

	// there is a constraint that prevents users from being deleted when they have a log
	q := `DELETE FROM log WHERE NOT tm_user = (SELECT id FROM tm_user WHERE username = 'admin')`
	err := execSQL(db, q)
	assert.RequireNoError(t, err, "Cannot execute SQL: %v; SQL is %s", err, q)

	q = `DELETE FROM tm_user WHERE username IN (` + strings.Join(usernames, ",") + `)`
	err = execSQL(db, q)
	assert.NoError(t, err, "Cannot execute SQL: %v; SQL is %s", err, q)
}

// this resets the IDs of things attached to a DS, which needs to be done
// because the WithObjs flow destroys and recreates those object IDs
// non-deterministically with each test - BUT, the client method permanently
// alters the DSR structures by adding these referential IDs. Older clients
// got away with it by not making 'DeliveryService' a pointer, but to add
// original/requested fields you need to sometimes allow each to be nil, so
// this is a problem that needs to be solved at some point.
// A better solution _might_ be to reload all the test fixtures every time
// to wipe any and all referential modifications made to any test data, but
// for now that's overkill.
func resetDS(ds *tc.DeliveryServiceV4) {
	if ds == nil {
		return
	}
	ds.CDNID = nil
	ds.ID = nil
	ds.ProfileID = nil
	ds.TenantID = nil
	ds.TypeID = nil
}

func CreateTestDeliveryServiceRequests(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, dsr := range td.DeliveryServiceRequests {
		resetDS(dsr.Original)
		resetDS(dsr.Requested)
		respDSR, _, err := cl.CreateDeliveryServiceRequest(dsr, toclient.RequestOptions{})
		assert.NoError(t, err, "Could not create Delivery Service Requests: %v - alerts: %+v", err, respDSR.Alerts)
	}
}

func DeleteTestDeliveryServiceRequests(t *testing.T, cl *toclient.Session) {
	resp, _, err := cl.GetDeliveryServiceRequests(toclient.RequestOptions{})
	assert.NoError(t, err, "Cannot get Delivery Service Requests: %v - alerts: %+v", err, resp.Alerts)
	for _, request := range resp.Response {
		alert, _, err := cl.DeleteDeliveryServiceRequest(*request.ID, toclient.RequestOptions{})
		assert.NoError(t, err, "Cannot delete Delivery Service Request #%d: %v - alerts: %+v", request.ID, err, alert.Alerts)

		// Retrieve the DeliveryServiceRequest to see if it got deleted
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.Itoa(*request.ID))
		dsr, _, err := cl.GetDeliveryServiceRequests(opts)
		assert.NoError(t, err, "Unexpected error fetching Delivery Service Request #%d after deletion: %v - alerts: %+v", *request.ID, err, dsr.Alerts)
		assert.Equal(t, len(dsr.Response), 0, "Expected Delivery Service Request #%d to be deleted, but it was found in Traffic Ops", *request.ID)
	}
}
