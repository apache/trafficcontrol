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

package client_tests

type TrafficControl struct {
	ASNs []struct {
		CachegroupName string `json:"cachegroupName"`
		Name           string `json:"name"`
	} `json:"asns"`
	Cachegroups []struct {
		Latitude                interface{} `json:"latitude"`
		Longitude               interface{} `json:"longitude"`
		Name                    string      `json:"name"`
		ParentCacheGroupName    interface{} `json:"parentCacheGroupName"`
		SecondaryCacheGroupName string      `json:"secondaryCacheGroupName"`
		ShortName               string      `json:"shortName"`
	} `json:"cachegroups"`
	CDNs []struct {
		DNSSECEnabled string `json:"dnssecEnabled"`
		DomainName    string `json:"domainName"`
		Name          string `json:"name"`
	} `json:"cdns"`
	DeliveryServices []struct {
		Active     bool   `json:"active"`
		DSCP       int64  `json:"dscp"`
		TenantName string `json:"tenantName"`
		XmlId      string `json:"xmlId"`
	} `json:"deliveryServices"`
	Divisions []struct {
		Name string `json:"name"`
	} `json:"divisions"`
	Regions []struct {
		DivisionName string `json:"divisionName"`
		Name         string `json:"name"`
	} `json:"regions"`
	Tenants []struct {
		Active           bool        `json:"active"`
		Name             string      `json:"name"`
		ParentTenantName interface{} `json:"parentTenantName"`
	} `json:"tenants"`
}
