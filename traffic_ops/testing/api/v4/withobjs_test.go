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

package v4

import (
	"database/sql"
	"testing"

	totest "github.com/apache/trafficcontrol/v8/lib/go-tc/totestv4"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	toclient "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"
)

// WithObjs creates the objs in order, runs f, and defers deleting the objs in the same order.
//
// Because deletion is deferred, using this ensures objects will be cleaned up if f panics or calls t.Fatal, as much as possible.
//
// Note that f itself may still create things which are not cleaned up properly, and likewise, the object creation and deletion tests themselves may fail.
// All tests in the Traffic Ops API Testing framework use the same Traffic Ops instance, with persistent data. Because of this, when any test fails, all subsequent tests should be considered invalid, irrespective whether they pass or fail. Users are encouraged to use `go test -failfast`.
func WithObjs(t *testing.T, objs []TCObj, f func()) {
	defer func() {
		for index := len(objs) - 1; index >= 0; index-- {
			obj := objs[index]
			withFuncs[obj].Delete(t)
		}
	}()
	for _, obj := range objs {
		withFuncs[obj].Create(t)
	}
	f()
}

type TCObj int

const (
	ASN TCObj = iota
	CacheGroups
	CacheGroupsDeliveryServices
	CDNs
	CDNLocks
	CDNFederations
	CDNNotifications
	Coordinates
	DeliveryServices
	DeliveryServicesRegexes
	DeliveryServiceRequests
	DeliveryServiceRequestComments
	DeliveryServicesRequiredCapabilities
	DeliveryServiceServerAssignments
	Divisions
	FederationDeliveryServices
	FederationResolvers
	FederationFederationResolvers
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
	ServiceCategories
	Statuses
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

func WrapClient(
	cl **toclient.Session,
	clF func(t *testing.T, toClient *toclient.Session),
) func(t *testing.T) {
	return func(t *testing.T) {
		clF(t, *cl)
	}
}

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
	dat *totest.TrafficControl,
	user *string,
	pass *string,
	clF func(t *testing.T, toClient *toclient.Session, dat totest.TrafficControl),
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("WrapNewClientDat user '%v' pass '%v'\n", *user, *pass)
		cl := utils.CreateV4Session(t, Config.TrafficOps.URL, *user, *pass, Config.Default.Session.TimeoutInSecs)
		clF(t, cl, *dat)
	}
}

func WrapClientDatDB(
	cl **toclient.Session,
	dat *totest.TrafficControl,
	clF func(t *testing.T, toClient *toclient.Session, dat totest.TrafficControl, db *sql.DB),
) func(t *testing.T) {
	return func(t *testing.T) {
		db, err := OpenConnection()
		assert.RequireNoError(t, err, "Cannot open db: %v", err)
		defer func() {
			err := db.Close()
			assert.NoError(t, err, "Unable to close connection to db: %v", err)
		}()
		clF(t, *cl, *dat, db)
	}
}

func NewDB(t *testing.T) *sql.DB {
	db, err := OpenConnection()
	assert.RequireNoError(t, err, "Cannot open db: %v", err)
	return db
}

var withFuncs = map[TCObj]TCObjFuncs{
	ASN:                                  {WrapClientDat(&TOSession, &testData, totest.CreateTestASNs), WrapClient(&TOSession, totest.DeleteTestASNs)},
	CacheGroups:                          {WrapClientDat(&TOSession, &testData, totest.CreateTestCacheGroups), WrapClientDat(&TOSession, &testData, totest.DeleteTestCacheGroups)},
	CacheGroupsDeliveryServices:          {WrapClient(&TOSession, totest.CreateTestCachegroupsDeliveryServices), WrapClient(&TOSession, totest.DeleteTestCachegroupsDeliveryServices)},
	CDNs:                                 {WrapClientDat(&TOSession, &testData, totest.CreateTestCDNs), WrapClient(&TOSession, totest.DeleteTestCDNs)},
	CDNLocks:                             {CreateTestCDNLocks, DeleteTestCDNLocks},
	CDNNotifications:                     {CreateTestCDNNotifications, DeleteTestCDNNotifications},
	CDNFederations:                       {WrapClientDat(&TOSession, &testData, totest.CreateTestCDNFederations), WrapClient(&TOSession, totest.DeleteTestCDNFederations)},
	Coordinates:                          {WrapClientDat(&TOSession, &testData, totest.CreateTestCoordinates), WrapClient(&TOSession, totest.DeleteTestCoordinates)},
	DeliveryServices:                     {WrapClientDat(&TOSession, &testData, totest.CreateTestDeliveryServices), WrapClient(&TOSession, totest.DeleteTestDeliveryServices)},
	DeliveryServicesRegexes:              {WrapClientDat(&TOSession, &testData, totest.CreateTestDeliveryServicesRegexes), WrapClientDatDB(&TOSession, &testData, totest.DeleteTestDeliveryServicesRegexes)},
	DeliveryServiceRequests:              {WrapClientDat(&TOSession, &testData, totest.CreateTestDeliveryServiceRequests), WrapClient(&TOSession, totest.DeleteTestDeliveryServiceRequests)},
	DeliveryServiceRequestComments:       {WrapClientDat(&TOSession, &testData, totest.CreateTestDeliveryServiceRequestComments), WrapClient(&TOSession, totest.DeleteTestDeliveryServiceRequestComments)},
	DeliveryServicesRequiredCapabilities: {WrapClientDat(&TOSession, &testData, totest.CreateTestDeliveryServicesRequiredCapabilities), WrapClient(&TOSession, totest.DeleteTestDeliveryServicesRequiredCapabilities)},
	DeliveryServiceServerAssignments:     {CreateTestDeliveryServiceServerAssignments, DeleteTestDeliveryServiceServers},
	Divisions:                            {WrapClientDat(&TOSession, &testData, totest.CreateTestDivisions), WrapClient(&TOSession, totest.DeleteTestDivisions)},
	FederationDeliveryServices:           {CreateTestFederationDeliveryServices, WrapClient(&TOSession, totest.DeleteTestCDNFederations)},
	FederationUsers:                      {WrapClient(&TOSession, totest.CreateTestFederationUsers), WrapClient(&TOSession, totest.DeleteTestFederationUsers)},
	FederationResolvers:                  {WrapClientDat(&TOSession, &testData, totest.CreateTestFederationResolvers), WrapClient(&TOSession, totest.DeleteTestFederationResolvers)},
	FederationFederationResolvers:        {CreateTestFederationFederationResolvers, DeleteTestFederationFederationResolvers},
	Jobs:                                 {WrapClientDat(&TOSession, &testData, totest.CreateTestJobs), WrapClient(&TOSession, totest.DeleteTestJobs)},
	Origins:                              {WrapClientDat(&TOSession, &testData, totest.CreateTestOrigins), WrapClient(&TOSession, totest.DeleteTestOrigins)},
	Parameters:                           {WrapClientDat(&TOSession, &testData, totest.CreateTestParameters), WrapClient(&TOSession, totest.DeleteTestParameters)},
	PhysLocations:                        {WrapClientDat(&TOSession, &testData, totest.CreateTestPhysLocations), WrapClient(&TOSession, totest.DeleteTestPhysLocations)},
	Profiles:                             {WrapClientDat(&TOSession, &testData, totest.CreateTestProfiles), WrapClient(&TOSession, totest.DeleteTestProfiles)},
	ProfileParameters:                    {WrapClientDat(&TOSession, &testData, totest.CreateTestProfileParameters), WrapClient(&TOSession, totest.DeleteTestProfileParameters)},
	Regions:                              {WrapClientDat(&TOSession, &testData, totest.CreateTestRegions), WrapClient(&TOSession, totest.DeleteTestRegions)},
	Roles:                                {WrapClientDat(&TOSession, &testData, totest.CreateTestRoles), WrapClient(&TOSession, totest.DeleteTestRoles)},
	ServerCapabilities:                   {WrapClientDat(&TOSession, &testData, totest.CreateTestServerCapabilities), WrapClient(&TOSession, totest.DeleteTestServerCapabilities)},
	ServerChecks:                         {WrapNewClientDat(&testData, &Config.TrafficOps.Users.Extension, &Config.TrafficOps.UserPassword, totest.CreateTestServerChecks), WrapClient(&TOSession, totest.DeleteTestServerChecks)},
	ServerServerCapabilities:             {WrapClientDat(&TOSession, &testData, totest.CreateTestServerServerCapabilities), WrapClient(&TOSession, totest.DeleteTestServerServerCapabilities)},
	Servers:                              {WrapClientDat(&TOSession, &testData, totest.CreateTestServers), WrapClient(&TOSession, totest.DeleteTestServers)},
	ServiceCategories:                    {WrapClientDat(&TOSession, &testData, totest.CreateTestServiceCategories), WrapClient(&TOSession, totest.DeleteTestServiceCategories)},
	Statuses:                             {WrapClientDat(&TOSession, &testData, totest.CreateTestStatuses), WrapClientDat(&TOSession, &testData, totest.DeleteTestStatuses)},
	StaticDNSEntries:                     {WrapClientDat(&TOSession, &testData, totest.CreateTestStaticDNSEntries), WrapClient(&TOSession, totest.DeleteTestStaticDNSEntries)},
	SteeringTargets:                      {WrapNewClientDat(&testData, util.StrPtr("steering"), util.StrPtr("pa$$word"), totest.CreateTestSteeringTargets), WrapNewClientDat(&testData, util.StrPtr("steering"), util.StrPtr("pa$$word"), totest.DeleteTestSteeringTargets)},
	Tenants:                              {WrapClientDat(&TOSession, &testData, totest.CreateTestTenants), WrapClient(&TOSession, totest.DeleteTestTenants)},
	ServerCheckExtensions:                {WrapNewClientDat(&testData, &Config.TrafficOps.Users.Extension, &Config.TrafficOps.UserPassword, totest.CreateTestServerCheckExtensions), WrapNewClientDat(&testData, &Config.TrafficOps.Users.Extension, &Config.TrafficOps.UserPassword, totest.DeleteTestServerCheckExtensions)},
	Topologies:                           {WrapClientDat(&TOSession, &testData, totest.CreateTestTopologies), WrapClient(&TOSession, totest.DeleteTestTopologies)},
	Types:                                {WrapClientDatDB(&TOSession, &testData, totest.CreateTestTypes), WrapClientDatDB(&TOSession, &testData, totest.DeleteTestTypes)},
	Users:                                {WrapClientDat(&TOSession, &testData, totest.CreateTestUsers), WrapClientDatDB(&TOSession, &testData, totest.ForceDeleteTestUsers)},
}
