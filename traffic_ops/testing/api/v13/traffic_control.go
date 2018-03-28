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

package v13

import (
	tcapi "github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc/v13"
)

// TrafficControl - maps to the tc-fixtures.json file
type TrafficControl struct {
	ASNs                           []tcapi.ASN                           `json:"asns"`
	CDNs                           []v13.CDN                             `json:"cdns"`
	CacheGroups                    []tcapi.CacheGroup                    `json:"cachegroups"`
	DeliveryServiceRequests        []tcapi.DeliveryServiceRequest        `json:"deliveryServiceRequests"`
	DeliveryServiceRequestComments []tcapi.DeliveryServiceRequestComment `json:"deliveryServiceRequestComments"`
	DeliveryServices               []tcapi.DeliveryService               `json:"deliveryservices"`
	Divisions                      []tcapi.Division                      `json:"divisions"`
	Profiles                       []tcapi.Profile                       `json:"profiles"`
	Parameters                     []tcapi.Parameter                     `json:"parameters"`
	PhysLocations                  []tcapi.PhysLocation                  `json:"physLocations"`
	Regions                        []tcapi.Region                        `json:"regions"`
	Servers                        []v13.Server                          `json:"servers"`
	Statuses                       []tcapi.Status                        `json:"statuses"`
	Tenants                        []tcapi.Tenant                        `json:"tenants"`
	Types                          []tcapi.Type                          `json:"types"`
}
