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

package fixtures

import "github.com/apache/incubator-trafficcontrol/traffic_ops/client"

// Cachegroups returns a default CacheGroupResponse to be used for testing.
func Cachegroups() *client.CacheGroupResponse {
	return &client.CacheGroupResponse{
		Response: []client.CacheGroup{
			client.CacheGroup{
				Name:        "edge-philadelphia",
				ShortName:   "phila",
				Type:        "EDGE_LOC",
				Longitude:   -5.66,
				Latitude:    55.67,
				ParentName:  "mid-northeast",
				LastUpdated: "2014-03-28 15:23:50",
			},
		},
	}
}
