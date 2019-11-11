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

package v14

import (
	"testing"
)

// WithObjs creates the objs in order, runs f, and defers deleting the objs in the same order.
//
// Because deletion is deferred, using this ensures objects will be cleaned up if f panics or calls t.Fatal, as much as possible.
//
// Note that f itself may still create things which are not cleaned up properly, and likewise, the object creation and deletion tests themselves may fail.
// All tests in the Traffic Ops API Testing framework use the same Traffic Ops instance, with persistent data. Because of this, when any test fails, all subsequent tests should be considered invalid, irrespective whether they pass or fail. Users are encouraged to use `go test -failfast`.
func WithObjs(t *testing.T, objs []TCObj, f func()) {
	for _, obj := range objs {
		defer withFuncs[obj].Delete(t)
		withFuncs[obj].Create(t)
	}
	f()
}

type TCObj int

const (
	CacheGroups TCObj = iota
	CacheGroupsDeliveryServices
	CacheGroupParameters
	CDNs
	CDNFederations
	Coordinates
	DeliveryServices
	DeliveryServiceRequests
	DeliveryServiceRequestComments
	DeliveryServicesRequiredCapabilities
	Divisions
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
	Statuses
	StaticDNSEntries
	SteeringTargets
	Tenants
	TOExtensions
	Types
	Users
	UsersDeliveryServices
)

type TCObjFuncs struct {
	Create func(t *testing.T)
	Delete func(t *testing.T)
}

var withFuncs = map[TCObj]TCObjFuncs{
	CacheGroups:                          {CreateTestCacheGroups, DeleteTestCacheGroups},
	CacheGroupsDeliveryServices:          {CreateTestCachegroupsDeliveryServices, DeleteTestCachegroupsDeliveryServices},
	CacheGroupParameters:                 {CreateTestCacheGroupParameters, DeleteTestCacheGroupParameters},
	CDNs:                                 {CreateTestCDNs, DeleteTestCDNs},
	CDNFederations:                       {CreateTestCDNFederations, DeleteTestCDNFederations},
	Coordinates:                          {CreateTestCoordinates, DeleteTestCoordinates},
	DeliveryServices:                     {CreateTestDeliveryServices, DeleteTestDeliveryServices},
	DeliveryServiceRequests:              {CreateTestDeliveryServiceRequests, DeleteTestDeliveryServiceRequests},
	DeliveryServiceRequestComments:       {CreateTestDeliveryServiceRequestComments, DeleteTestDeliveryServiceRequestComments},
	DeliveryServicesRequiredCapabilities: {CreateTestDeliveryServicesRequiredCapabilities, DeleteTestDeliveryServicesRequiredCapabilities},
	Divisions:                            {CreateTestDivisions, DeleteTestDivisions},
	FederationUsers:                      {CreateTestFederationUsers, DeleteTestFederationUsers},
	Origins:                              {CreateTestOrigins, DeleteTestOrigins},
	Parameters:                           {CreateTestParameters, DeleteTestParameters},
	PhysLocations:                        {CreateTestPhysLocations, DeleteTestPhysLocations},
	Profiles:                             {CreateTestProfiles, DeleteTestProfiles},
	ProfileParameters:                    {CreateTestProfileParameters, DeleteTestProfileParameters},
	Regions:                              {CreateTestRegions, DeleteTestRegions},
	Roles:                                {CreateTestRoles, DeleteTestRoles},
	ServerCapabilities:                   {CreateTestServerCapabilities, DeleteTestServerCapabilities},
	ServerChecks:                         {CreateTestServerChecks, DeleteTestServerChecks},
	ServerServerCapabilities:             {CreateTestServerServerCapabilities, DeleteTestServerServerCapabilities},
	Servers:                              {CreateTestServers, DeleteTestServers},
	Statuses:                             {CreateTestStatuses, DeleteTestStatuses},
	StaticDNSEntries:                     {CreateTestStaticDNSEntries, DeleteTestStaticDNSEntries},
	SteeringTargets:                      {SetupSteeringTargets, DeleteTestSteeringTargets},
	Tenants:                              {CreateTestTenants, DeleteTestTenants},
	TOExtensions:                         {CreateTestTOExtensions, DeleteTestTOExtensions},
	Types:                                {CreateTestTypes, DeleteTestTypes},
	Users:                                {CreateTestUsers, ForceDeleteTestUsers},
	UsersDeliveryServices:                {CreateTestUsersDeliveryServices, DeleteTestUsersDeliveryServices},
}
