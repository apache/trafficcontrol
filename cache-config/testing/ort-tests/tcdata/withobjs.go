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
	"testing"
)

type TCObj int

const (
	CacheGroups TCObj = iota
	CacheGroupsDeliveryServices
	CacheGroupParameters
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

// WithObjs creates the objs in order, runs f, and defers deleting the objs in the same order.
//
// Because deletion is deferred, using this ensures objects will be cleaned up if f panics or calls t.Fatal, as much as possible.
//
// Note that f itself may still create things which are not cleaned up properly, and likewise, the object creation and deletion tests themselves may fail.
// All tests in the Traffic Ops API Testing framework use the same Traffic Ops instance, with persistent data. Because of this, when any test fails, all subsequent tests should be considered invalid, irrespective whether they pass or fail. Users are encouraged to use `go test -failfast`.
func (r *TCData) WithObjs(t *testing.T, objs []TCObj, f func()) {
	var withFuncs = map[TCObj]TCObjFuncs{
		CacheGroups:                          {r.CreateTestCacheGroups, r.DeleteTestCacheGroups},
		CacheGroupsDeliveryServices:          {r.CreateTestCachegroupsDeliveryServices, r.DeleteTestCachegroupsDeliveryServices},
		CacheGroupParameters:                 {r.CreateTestCacheGroupParameters, r.DeleteTestCacheGroupParameters},
		CDNs:                                 {r.CreateTestCDNs, r.DeleteTestCDNs},
		CDNFederations:                       {r.CreateTestCDNFederations, r.DeleteTestCDNFederations},
		Coordinates:                          {r.CreateTestCoordinates, r.DeleteTestCoordinates},
		DeliveryServices:                     {r.CreateTestDeliveryServices, r.DeleteTestDeliveryServices},
		DeliveryServicesRegexes:              {r.CreateTestDeliveryServicesRegexes, r.DeleteTestDeliveryServicesRegexes},
		DeliveryServiceRequests:              {r.CreateTestDeliveryServiceRequests, r.DeleteTestDeliveryServiceRequests},
		DeliveryServiceRequestComments:       {r.CreateTestDeliveryServiceRequestComments, r.DeleteTestDeliveryServiceRequestComments},
		DeliveryServicesRequiredCapabilities: {r.CreateTestDeliveryServicesRequiredCapabilities, r.DeleteTestDeliveryServicesRequiredCapabilities},
		Divisions:                            {r.CreateTestDivisions, r.DeleteTestDivisions},
		FederationUsers:                      {r.CreateTestFederationUsers, r.DeleteTestFederationUsers},
		FederationResolvers:                  {r.CreateTestFederationResolvers, r.DeleteTestFederationResolvers},
		Origins:                              {r.CreateTestOrigins, r.DeleteTestOrigins},
		Parameters:                           {r.CreateTestParameters, r.DeleteTestParameters},
		PhysLocations:                        {r.CreateTestPhysLocations, r.DeleteTestPhysLocations},
		Profiles:                             {r.CreateTestProfiles, r.DeleteTestProfiles},
		ProfileParameters:                    {r.CreateTestProfileParameters, r.DeleteTestProfileParameters},
		Regions:                              {r.CreateTestRegions, r.DeleteTestRegions},
		Roles:                                {r.CreateTestRoles, r.DeleteTestRoles},
		ServerCapabilities:                   {r.CreateTestServerCapabilities, r.DeleteTestServerCapabilities},
		ServerChecks:                         {r.CreateTestServerChecks, r.DeleteTestServerChecks},
		ServerServerCapabilities:             {r.CreateTestServerServerCapabilities, r.DeleteTestServerServerCapabilities},
		Servers:                              {r.CreateTestServers, r.DeleteTestServers},
		ServiceCategories:                    {r.CreateTestServiceCategories, r.DeleteTestServiceCategories},
		Statuses:                             {r.CreateTestStatuses, r.DeleteTestStatuses},
		StaticDNSEntries:                     {r.CreateTestStaticDNSEntries, r.DeleteTestStaticDNSEntries},
		SteeringTargets:                      {r.SetupSteeringTargets, r.DeleteTestSteeringTargets},
		Tenants:                              {r.CreateTestTenants, r.DeleteTestTenants},
		ServerCheckExtensions:                {r.CreateTestServerCheckExtensions, r.DeleteTestServerCheckExtensions},
		Topologies:                           {r.CreateTestTopologies, r.DeleteTestTopologies},
		Types:                                {r.CreateTestTypes, r.DeleteTestTypes},
		Users:                                {r.CreateTestUsers, r.DeleteTestUsers},
	}

	for _, obj := range objs {
		defer withFuncs[obj].Delete(t)
		withFuncs[obj].Create(t)
	}
	f()
}
