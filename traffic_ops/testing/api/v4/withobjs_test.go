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
	"testing"
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
	CDNFederations
	CDNNotifications
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
	ServerServerCapabilitiesForTopologies
	Servers
	ServiceCategories
	Statuses
	StaticDNSEntries
	SteeringTargets
	Tenants
	ServerCheckExtensions
	Topologies
	TopologyBasedDeliveryServiceRequiredCapabilities
	Types
	Users
)

type TCObjFuncs struct {
	Create func(t *testing.T)
	Delete func(t *testing.T)
}

var withFuncs = map[TCObj]TCObjFuncs{
	ASN:                                   {CreateTestASNs, DeleteTestASNs},
	CacheGroups:                           {CreateTestCacheGroups, DeleteTestCacheGroups},
	CacheGroupsDeliveryServices:           {CreateTestCachegroupsDeliveryServices, DeleteTestCachegroupsDeliveryServices},
	CDNs:                                  {CreateTestCDNs, DeleteTestCDNs},
	CDNNotifications:                      {CreateTestCDNNotifications, DeleteTestCDNNotifications},
	CDNFederations:                        {CreateTestCDNFederations, DeleteTestCDNFederations},
	Coordinates:                           {CreateTestCoordinates, DeleteTestCoordinates},
	DeliveryServices:                      {CreateTestDeliveryServices, DeleteTestDeliveryServices},
	DeliveryServicesRegexes:               {CreateTestDeliveryServicesRegexes, DeleteTestDeliveryServicesRegexes},
	DeliveryServiceRequests:               {CreateTestDeliveryServiceRequests, DeleteTestDeliveryServiceRequests},
	DeliveryServiceRequestComments:        {CreateTestDeliveryServiceRequestComments, DeleteTestDeliveryServiceRequestComments},
	DeliveryServicesRequiredCapabilities:  {CreateTestDeliveryServicesRequiredCapabilities, DeleteTestDeliveryServicesRequiredCapabilities},
	Divisions:                             {CreateTestDivisions, DeleteTestDivisions},
	FederationUsers:                       {CreateTestFederationUsers, DeleteTestFederationUsers},
	FederationResolvers:                   {CreateTestFederationResolvers, DeleteTestFederationResolvers},
	Origins:                               {CreateTestOrigins, DeleteTestOrigins},
	Parameters:                            {CreateTestParameters, DeleteTestParameters},
	PhysLocations:                         {CreateTestPhysLocations, DeleteTestPhysLocations},
	Profiles:                              {CreateTestProfiles, DeleteTestProfiles},
	ProfileParameters:                     {CreateTestProfileParameters, DeleteTestProfileParameters},
	Regions:                               {CreateTestRegions, DeleteTestRegions},
	Roles:                                 {CreateTestRoles, DeleteTestRoles},
	ServerCapabilities:                    {CreateTestServerCapabilities, DeleteTestServerCapabilities},
	ServerChecks:                          {CreateTestServerChecks, DeleteTestServerChecks},
	ServerServerCapabilities:              {CreateTestServerServerCapabilities, DeleteTestServerServerCapabilities},
	ServerServerCapabilitiesForTopologies: {CreateTestServerServerCapabilities, DeleteTestServerServerCapabilitiesForTopologies},
	Servers:                               {CreateTestServers, DeleteTestServers},
	ServiceCategories:                     {CreateTestServiceCategories, DeleteTestServiceCategories},
	Statuses:                              {CreateTestStatuses, DeleteTestStatuses},
	StaticDNSEntries:                      {CreateTestStaticDNSEntries, DeleteTestStaticDNSEntries},
	SteeringTargets:                       {SetupSteeringTargets, DeleteTestSteeringTargets},
	Tenants:                               {CreateTestTenants, DeleteTestTenants},
	ServerCheckExtensions:                 {CreateTestServerCheckExtensions, DeleteTestServerCheckExtensions},
	Topologies:                            {CreateTestTopologies, DeleteTestTopologies},
	TopologyBasedDeliveryServiceRequiredCapabilities: {CreateTestTopologyBasedDeliveryServicesRequiredCapabilities, DeleteTestDeliveryServicesRequiredCapabilities},
	Types: {CreateTestTypes, DeleteTestTypes},
	Users: {CreateTestUsers, ForceDeleteTestUsers},
}
