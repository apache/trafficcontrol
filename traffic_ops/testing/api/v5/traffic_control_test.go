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

package v5

import (
	"github.com/apache/trafficcontrol/lib/go-tc"
)

// TrafficControl - maps to the tc-fixtures.json file
type TrafficControl struct {
	ASNs                                              []tc.ASNV5                              `json:"asns"`
	CDNs                                              []tc.CDN                                `json:"cdns"`
	CDNLocks                                          []tc.CDNLock                            `json:"cdnlocks"`
	CacheGroups                                       []tc.CacheGroupNullable                 `json:"cachegroups"`
	Capabilities                                      []tc.Capability                         `json:"capability"`
	Coordinates                                       []tc.CoordinateV5                       `json:"coordinates"`
	DeliveryServicesRegexes                           []tc.DeliveryServiceRegexesTest         `json:"deliveryServicesRegexes"`
	DeliveryServiceRequests                           []tc.DeliveryServiceRequestV5           `json:"deliveryServiceRequests"`
	DeliveryServiceRequestComments                    []tc.DeliveryServiceRequestCommentV5    `json:"deliveryServiceRequestComments"`
	DeliveryServices                                  []tc.DeliveryServiceV5                  `json:"deliveryservices"`
	DeliveryServicesRequiredCapabilities              []tc.DeliveryServicesRequiredCapability `json:"deliveryservicesRequiredCapabilities"`
	DeliveryServiceServerAssignments                  []tc.DeliveryServiceServers             `json:"deliveryServiceServerAssignments"`
	TopologyBasedDeliveryServicesRequiredCapabilities []tc.DeliveryServicesRequiredCapability `json:"topologyBasedDeliveryServicesRequiredCapabilities"`
	Divisions                                         []tc.Division                           `json:"divisions"`
	Federations                                       []tc.CDNFederation                      `json:"federations"`
	FederationResolvers                               []tc.FederationResolverV5               `json:"federation_resolvers"`
	Jobs                                              []tc.InvalidationJobCreateV4            `json:"jobs"`
	Origins                                           []tc.Origin                             `json:"origins"`
	Profiles                                          []tc.Profile                            `json:"profiles"`
	Parameters                                        []tc.Parameter                          `json:"parameters"`
	ProfileParameters                                 []tc.ProfileParameter                   `json:"profileParameters"`
	PhysLocations                                     []tc.PhysLocation                       `json:"physLocations"`
	Regions                                           []tc.Region                             `json:"regions"`
	Roles                                             []tc.RoleV4                             `json:"roles"`
	Servers                                           []tc.ServerV4                           `json:"servers"`
	ServerServerCapabilities                          []tc.ServerServerCapability             `json:"serverServerCapabilities"`
	ServerCapabilities                                []tc.ServerCapabilityV5                 `json:"serverCapabilities"`
	ServiceCategories                                 []tc.ServiceCategoryV5                  `json:"serviceCategories"`
	Statuses                                          []tc.StatusNullable                     `json:"statuses"`
	StaticDNSEntries                                  []tc.StaticDNSEntry                     `json:"staticdnsentries"`
	StatsSummaries                                    []tc.StatsSummaryV5                     `json:"statsSummaries"`
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
