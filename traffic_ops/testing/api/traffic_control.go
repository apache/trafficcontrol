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

package api

import tc "github.com/apache/incubator-trafficcontrol/lib/go-tc"

// TrafficControl - maps to the tc-fixtures.json file
type TrafficControl struct {
	ASNs             []tc.ASN             `json:"asns"`
	CDNs             []tc.CDN             `json:"cdns"`
	Cachegroups      []tc.CacheGroup      `json:"cachegroups"`
	DeliveryServices []tc.DeliveryService `json:"deliveryservices"`
	Divisions        []tc.Division        `json:"divisions"`
	Regions          []tc.Region          `json:"regions"`
	Tenants          []tc.Tenant          `json:"tenants"`
}
