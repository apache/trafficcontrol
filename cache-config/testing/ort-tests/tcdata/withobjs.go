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
	"database/sql"
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-tc/totest"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	toclient "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

type TCObj int

const (
	ASN TCObj = iota
	CacheGroups
	CacheGroupsDeliveryServices
	// CacheGroupParameters
	CDNs
	CDNFederations
	Coordinates
	DeliveryServices
	DeliveryServicesRegexes
	DeliveryServiceRequests
	DeliveryServiceRequestComments
	DeliveryServicesRequiredCapabilities
	Divisions
	FederationResolvers
	FederationUsers
	Jobs
	Origins
	Parameters
	PhysLocations
	Profiles
	ProfileParameters
	Regions
	Roles
	ServerCapabilities
	ServerChecks
	ServerServerCapabilities
	Servers
	Statuses
	ServiceCategories
	StaticDNSEntries
	SteeringTargets
	Tenants
	ServerCheckExtensions
	Topologies
	Types
	Users
)

type TCObjFuncs struct {
	Create func(t *testing.T)
	Delete func(t *testing.T)
}

// WrapClient is a temp func that creates a func(t *testing.T) from a *toclient.Session and a func(t *testing.T, toClient *toclient.Session).
func WrapClient(
	cl **toclient.Session,
	clF func(t *testing.T, toClient *toclient.Session),
) func(t *testing.T) {
	return func(t *testing.T) {
		clF(t, *cl)
	}
}

// WrapClient is a temp func that creates a func(t *testing.T) from a *toclient.Session, totest.TrafficControl, and a func(t *testing.T, toClient *toclient.Session, dat totest.TrafficControl).
func WrapClientDat(
	cl **toclient.Session,
	dat *totest.TrafficControl,
	clF func(t *testing.T, toClient *toclient.Session, dat totest.TrafficControl),
) func(t *testing.T) {
	return func(t *testing.T) {
		clF(t, *cl, *dat)
	}
}

// WrapNewClientDat is like WrapClientDat but creates a new custom TO client instead of using the global variable.
func WrapNewClientDat(
	r *TCData,
	dat *totest.TrafficControl,
	user *string,
	pass *string,
	clF func(t *testing.T, toClient *toclient.Session, dat totest.TrafficControl),
) func(t *testing.T) {
	return func(t *testing.T) {
		cl := utils.CreateV5Session(t, r.Config.TrafficOps.URL, *user, *pass, r.Config.Default.Session.TimeoutInSecs)
		clF(t, cl, *dat)
	}
}

func WrapClientDatDB(
	r *TCData,
	cl **toclient.Session,
	dat *totest.TrafficControl,
	clF func(t *testing.T, toClient *toclient.Session, dat totest.TrafficControl, db *sql.DB),
) func(t *testing.T) {
	return func(t *testing.T) {
		db, err := r.OpenConnection()
		assert.RequireNoError(t, err, "Cannot open db: %v", err)
		defer func() {
			err := db.Close()
			assert.NoError(t, err, "Unable to close connection to db: %v", err)
		}()
		clF(t, *cl, *dat, db)
	}
}

// WithObjs creates the objs in order, runs f, and defers deleting the objs in the same order.
//
// Because deletion is deferred, using this ensures objects will be cleaned up if f panics or calls t.Fatal, as much as possible.
//
// Note that f itself may still create things which are not cleaned up properly, and likewise, the object creation and deletion tests themselves may fail.
// All tests in the Traffic Ops API Testing framework use the same Traffic Ops instance, with persistent data. Because of this, when any test fails, all subsequent tests should be considered invalid, irrespective whether they pass or fail. Users are encouraged to use `go test -failfast`.
func (r *TCData) WithObjs(t *testing.T, objs []TCObj, f func()) {
	var withFuncs = map[TCObj]TCObjFuncs{
		ASN:                         {WrapClientDat(&TOSession, r.TestData, totest.CreateTestASNs), WrapClient(&TOSession, totest.DeleteTestASNs)},
		CacheGroups:                 {WrapClientDat(&TOSession, r.TestData, totest.CreateTestCacheGroups), WrapClientDat(&TOSession, r.TestData, totest.DeleteTestCacheGroups)},
		CacheGroupsDeliveryServices: {WrapClient(&TOSession, totest.CreateTestCachegroupsDeliveryServices), WrapClient(&TOSession, totest.DeleteTestCachegroupsDeliveryServices)},
		// CacheGroupParameters:                 {WrapClientDat(&TOSession, r.TestData, totest.CreateTestCacheGroupParameters), WrapClient(&TOSession, totest.DeleteTestCacheGroupParameters)},
		CDNs:                                 {WrapClientDat(&TOSession, r.TestData, totest.CreateTestCDNs), WrapClient(&TOSession, totest.DeleteTestCDNs)},
		CDNFederations:                       {WrapClientDat(&TOSession, r.TestData, totest.CreateTestCDNFederations), WrapClient(&TOSession, totest.DeleteTestCDNFederations)},
		Coordinates:                          {WrapClientDat(&TOSession, r.TestData, totest.CreateTestCoordinates), WrapClient(&TOSession, totest.DeleteTestCoordinates)},
		DeliveryServices:                     {WrapClientDat(&TOSession, r.TestData, totest.CreateTestDeliveryServices), WrapClient(&TOSession, totest.DeleteTestDeliveryServices)},
		DeliveryServicesRegexes:              {WrapClientDat(&TOSession, r.TestData, totest.CreateTestDeliveryServicesRegexes), WrapClientDatDB(r, &TOSession, r.TestData, totest.DeleteTestDeliveryServicesRegexes)},
		DeliveryServiceRequests:              {WrapClientDat(&TOSession, r.TestData, totest.CreateTestDeliveryServiceRequests), WrapClient(&TOSession, totest.DeleteTestDeliveryServiceRequests)},
		DeliveryServiceRequestComments:       {WrapClientDat(&TOSession, r.TestData, totest.CreateTestDeliveryServiceRequestComments), WrapClient(&TOSession, totest.DeleteTestDeliveryServiceRequestComments)},
		DeliveryServicesRequiredCapabilities: {WrapClientDat(&TOSession, r.TestData, totest.CreateTestDeliveryServicesRequiredCapabilities), WrapClient(&TOSession, totest.DeleteTestDeliveryServicesRequiredCapabilities)},
		Divisions:                            {WrapClientDat(&TOSession, r.TestData, totest.CreateTestDivisions), WrapClient(&TOSession, totest.DeleteTestDivisions)},
		FederationUsers:                      {WrapClient(&TOSession, totest.CreateTestFederationUsers), WrapClient(&TOSession, totest.DeleteTestFederationUsers)},
		FederationResolvers:                  {WrapClientDat(&TOSession, r.TestData, totest.CreateTestFederationResolvers), WrapClient(&TOSession, totest.DeleteTestFederationResolvers)},
		Jobs:                                 {WrapClientDat(&TOSession, r.TestData, totest.CreateTestJobs), WrapClient(&TOSession, totest.DeleteTestJobs)},
		Origins:                              {WrapClientDat(&TOSession, r.TestData, totest.CreateTestOrigins), WrapClient(&TOSession, totest.DeleteTestOrigins)},
		Parameters:                           {WrapClientDat(&TOSession, r.TestData, totest.CreateTestParameters), WrapClient(&TOSession, totest.DeleteTestParameters)},
		PhysLocations:                        {WrapClientDat(&TOSession, r.TestData, totest.CreateTestPhysLocations), WrapClient(&TOSession, totest.DeleteTestPhysLocations)},
		Profiles:                             {WrapClientDat(&TOSession, r.TestData, totest.CreateTestProfiles), WrapClient(&TOSession, totest.DeleteTestProfiles)},
		ProfileParameters:                    {WrapClientDat(&TOSession, r.TestData, totest.CreateTestProfileParameters), WrapClient(&TOSession, totest.DeleteTestProfileParameters)},
		Regions:                              {WrapClientDat(&TOSession, r.TestData, totest.CreateTestRegions), WrapClient(&TOSession, totest.DeleteTestRegions)},
		Roles:                                {WrapClientDat(&TOSession, r.TestData, totest.CreateTestRoles), WrapClient(&TOSession, totest.DeleteTestRoles)},
		ServerCapabilities:                   {WrapClientDat(&TOSession, r.TestData, totest.CreateTestServerCapabilities), WrapClient(&TOSession, totest.DeleteTestServerCapabilities)},
		ServerChecks:                         {WrapNewClientDat(r, r.TestData, &r.Config.TrafficOps.Users.Extension, &r.Config.TrafficOps.UserPassword, totest.CreateTestServerChecks), WrapClient(&TOSession, totest.DeleteTestServerChecks)},
		ServerServerCapabilities:             {WrapClientDat(&TOSession, r.TestData, totest.CreateTestServerServerCapabilities), WrapClient(&TOSession, totest.DeleteTestServerServerCapabilities)},
		Servers:                              {WrapClientDat(&TOSession, r.TestData, totest.CreateTestServers), WrapClient(&TOSession, totest.DeleteTestServers)},
		ServiceCategories:                    {WrapClientDat(&TOSession, r.TestData, totest.CreateTestServiceCategories), WrapClient(&TOSession, totest.DeleteTestServiceCategories)},
		Statuses:                             {WrapClientDat(&TOSession, r.TestData, totest.CreateTestStatuses), WrapClientDat(&TOSession, r.TestData, totest.DeleteTestStatuses)},
		StaticDNSEntries:                     {WrapClientDat(&TOSession, r.TestData, totest.CreateTestStaticDNSEntries), WrapClient(&TOSession, totest.DeleteTestStaticDNSEntries)},
		SteeringTargets:                      {WrapNewClientDat(r, r.TestData, util.StrPtr("steering"), util.StrPtr("pa$$word"), totest.CreateTestSteeringTargets), WrapNewClientDat(r, r.TestData, util.StrPtr("steering"), util.StrPtr("pa$$word"), totest.DeleteTestSteeringTargets)},
		Tenants:                              {WrapClientDat(&TOSession, r.TestData, totest.CreateTestTenants), WrapClient(&TOSession, totest.DeleteTestTenants)},
		ServerCheckExtensions:                {WrapNewClientDat(r, r.TestData, &r.Config.TrafficOps.Users.Extension, &r.Config.TrafficOps.UserPassword, totest.CreateTestServerCheckExtensions), WrapNewClientDat(r, r.TestData, &r.Config.TrafficOps.Users.Extension, &r.Config.TrafficOps.UserPassword, totest.DeleteTestServerCheckExtensions)},
		Topologies:                           {WrapClientDat(&TOSession, r.TestData, totest.CreateTestTopologies), WrapClient(&TOSession, totest.DeleteTestTopologies)},
		Types:                                {WrapClientDatDB(r, &TOSession, r.TestData, totest.CreateTestTypes), WrapClientDatDB(r, &TOSession, r.TestData, totest.DeleteTestTypes)},
		Users:                                {WrapClientDat(&TOSession, r.TestData, totest.CreateTestUsers), WrapClientDatDB(r, &TOSession, r.TestData, totest.ForceDeleteTestUsers)},
	}

	for _, obj := range objs {
		defer withFuncs[obj].Delete(t)
		withFuncs[obj].Create(t)
	}
	f()
}
